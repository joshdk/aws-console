// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package credentials contains functionality for retrieving and manipulating
// AWS credentials.
package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/processcreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// FromConfig retrieves credentials from the AWS cli config files, typically
// ~/.aws/credentials and ~/.aws/config. Credentials for the named profile are
// returned, or the default profile if no name is given. Additionally, the
// value of $AWS_PROFILE will be used if it is set.
func FromConfig(profile string) (*aws.Credentials, string, error) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return nil, "", err
	}

	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, "", err
	}

	return &creds, cfg.Region, nil
}

// FromReader retrieves credentials from given io.Reader, typically os.Stdin.
// Expects JSON data in one of two possible formats. The first is returned by
// several STS operations (assume-role/get-session-token/etc) and looks like:
//
//	{
//	    "AssumedRoleUser": {...},
//	    "Credentials": {
//	        "AccessKeyId":     "...",
//	        "SecretAccessKey": "...",
//	        "SessionToken":    "..."
//	        "Expiration":      "...",
//	    }
//	}
//
// The second is returned by various AWS cli credential exec plugins, and looks
// like:
//
//	{
//	    "AccessKeyId":     "...",
//	    "SecretAccessKey": "...",
//	    "SessionToken":    "...",
//	    "Expiration":      "...",
//	    "Version":         1
//	}
//
// See https://docs.aws.amazon.com/cli/latest/reference/sts/assume-role.html#output.
// See https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html.
func FromReader(reader io.Reader) (*aws.Credentials, error) {
	// Read the entire body, as it will be potentially parsed multiple times.
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	type creds struct {
		Credentials processcreds.CredentialProcessResponse `json:"Credentials"`
	}

	var result creds

	if err := json.Unmarshal(body, &result); err == nil && result.Credentials.AccessKeyID != "" && result.Credentials.SecretAccessKey != "" {
		// Credentials were unmarshalled into the entire struct.
		return &aws.Credentials{
			AccessKeyID:     result.Credentials.AccessKeyID,
			SecretAccessKey: result.Credentials.SecretAccessKey,
			SessionToken:    result.Credentials.SessionToken,
		}, nil
	} else if err := json.Unmarshal(body, &result.Credentials); err == nil && result.Credentials.AccessKeyID != "" && result.Credentials.SecretAccessKey != "" {
		// Credentials were unmarshalled into part of the struct.
		return &aws.Credentials{
			AccessKeyID:     result.Credentials.AccessKeyID,
			SecretAccessKey: result.Credentials.SecretAccessKey,
			SessionToken:    result.Credentials.SessionToken,
		}, nil
	}

	// Credentials could not be fully unmarshalled.
	return nil, errors.New("failed to parse credentials")
}

// FederateUser will federate the given user credentials by calling STS
// GetFederationToken. If the given credentials are not for a user (like
// credentials for a role) then they are returned unmodified.
func FederateUser(creds *aws.Credentials, region, name, policy string, duration time.Duration, userAgent string) (*aws.Credentials, error) {
	// Only federate if user credentials were given.
	if creds.SessionToken != "" {
		return creds, nil
	}

	client := sts.NewFromConfig(
		aws.Config{
			Credentials: credentials.NewStaticCredentialsProvider(
				creds.AccessKeyID,
				creds.SecretAccessKey,
				creds.SessionToken,
			),
			Region: region,
		},
		func(options *sts.Options) {
			options.APIOptions = append(options.APIOptions, setUserAgent(userAgent))
		},
	)

	input := sts.GetFederationTokenInput{
		Name: aws.String(name),
		PolicyArns: []types.PolicyDescriptorType{{
			Arn: aws.String(policy),
		}},
	}

	// The minimum value for the DurationSeconds parameter is 15 minutes.
	// See https://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html#API_GetFederationToken_RequestParameters.
	const minDuration = 15 * time.Minute
	if duration != 0 && duration < minDuration {
		duration = minDuration
	}

	if duration != 0 {
		input.DurationSeconds = aws.Int32(int32(duration.Seconds()))
	}

	// Federate the user.
	result, err := client.GetFederationToken(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return &aws.Credentials{
		AccessKeyID:     aws.ToString(result.Credentials.AccessKeyId),
		SecretAccessKey: aws.ToString(result.Credentials.SecretAccessKey),
		SessionToken:    aws.ToString(result.Credentials.SessionToken),
	}, nil
}

func setUserAgent(useragent string) func(stack *middleware.Stack) error {
	return func(stack *middleware.Stack) error {
		bm := userAgentMiddleware(useragent)
		stack.Build.Remove(bm.ID()) //nolint

		return stack.Build.Add(&bm, middleware.After)
	}
}

type userAgentMiddleware string

func (userAgentMiddleware) ID() string {
	return "UserAgent"
}

func (u userAgentMiddleware) HandleBuild(ctx context.Context, in middleware.BuildInput, next middleware.BuildHandler) (middleware.BuildOutput, middleware.Metadata, error) {
	if req, ok := in.Request.(*smithyhttp.Request); ok {
		req.Header.Set("User-Agent", string(u))
	}

	return next.HandleBuild(ctx, in)
}
