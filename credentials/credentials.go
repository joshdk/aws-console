// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package credentials contains functionality for retrieving and manipulating
// AWS credentials.
package credentials

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// FromConfig retrieves credentials from the AWS cli config files, typically
// ~/.aws/credentials and ~/.aws/config. Credentials for the named profile are
// returned, or the default profile if no name is given. Additionally, the
// value of $AWS_PROFILE will be used if it is set.
func FromConfig(profile string) (*sts.Credentials, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	value, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	return &sts.Credentials{
		AccessKeyId:     aws.String(value.AccessKeyID),
		SecretAccessKey: aws.String(value.SecretAccessKey),
		SessionToken:    aws.String(value.SessionToken),
	}, nil
}

// FromReader retrieves credentials from given io.Reader, typically os.Stdin.
// Expects JSON data in one of two possible formats. The first is returned by
// several STS operations (assume-role/get-session-token/etc) and looks like:
//
//     {
//         "AssumedRoleUser": {...},
//         "Credentials": {
//             "AccessKeyId":     "...",
//             "SecretAccessKey": "...",
//             "SessionToken":    "..."
//             "Expiration":      "...",
//         }
//     }
//
// The second is returned by various AWS cli credential exec plugins, and looks
// like:
//
//     {
//         "AccessKeyId":     "...",
//         "SecretAccessKey": "...",
//         "SessionToken":    "...",
//         "Expiration":      "...",
//         "Version":         1
//     }
//
// See https://docs.aws.amazon.com/cli/latest/reference/sts/assume-role.html#output.
// See https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-sourcing-external.html.
func FromReader(reader io.Reader) (*sts.Credentials, error) {
	// Read the entire body, as it will be potentially parsed multiple times.
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	type creds struct {
		Credentials struct {
			AccessKeyID     string `json:"AccessKeyId"`
			SecretAccessKey string `json:"SecretAccessKey"`
			SessionToken    string `json:"SessionToken"`
		} `json:"Credentials"`
	}

	var result creds

	if err := json.Unmarshal(body, &result); err == nil && result.Credentials.AccessKeyID != "" && result.Credentials.SecretAccessKey != "" {
		// Credentials were unmarshaled into the entire struct.
		return &sts.Credentials{
			AccessKeyId:     aws.String(result.Credentials.AccessKeyID),
			SecretAccessKey: aws.String(result.Credentials.SecretAccessKey),
			SessionToken:    aws.String(result.Credentials.SessionToken),
		}, nil
	} else if err := json.Unmarshal(body, &result.Credentials); err == nil && result.Credentials.AccessKeyID != "" && result.Credentials.SecretAccessKey != "" {
		// Credentials were unmarshaled into part of the struct.
		return &sts.Credentials{
			AccessKeyId:     aws.String(result.Credentials.AccessKeyID),
			SecretAccessKey: aws.String(result.Credentials.SecretAccessKey),
			SessionToken:    aws.String(result.Credentials.SessionToken),
		}, nil
	}

	// Credentials could not be fully unmarshaled.
	return nil, fmt.Errorf("failed to parse credentials") // nolint:goerr113
}

// FederateUser will federate the given user credentials by calling STS
// GetFederationToken. If the given credentials are not for a user (like
// credentials for a role) then they are returned unmodified.
func FederateUser(creds *sts.Credentials, name, policy string, duration time.Duration, userAgent string) (*sts.Credentials, error) {
	// Only federate if user credentials were given.
	if aws.StringValue(creds.SessionToken) != "" {
		return creds, nil
	}

	// Create a new session given the static user credentials.
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			aws.StringValue(creds.AccessKeyId),
			aws.StringValue(creds.SecretAccessKey),
			aws.StringValue(creds.SessionToken),
		),
	})
	if err != nil {
		return nil, err
	}

	// The minimum value for the DurationSeconds parameter is 15 minutes.
	// See https://docs.aws.amazon.com/STS/latest/APIReference/API_GetFederationToken.html#API_GetFederationToken_RequestParameters.
	const minDuration = 15 * time.Minute
	if duration < minDuration {
		duration = minDuration
	}

	input := sts.GetFederationTokenInput{
		DurationSeconds: aws.Int64(int64(duration.Seconds())),
		Name:            aws.String(name),
		PolicyArns: []*sts.PolicyDescriptorType{{
			Arn: aws.String(policy),
		}},
	}

	// Configure client.
	client := sts.New(sess)
	client.Handlers.Build.PushBack(request.WithSetRequestHeaders(map[string]string{"User-Agent": userAgent}))

	// Federate the user.
	result, err := client.GetFederationToken(&input)
	if err != nil {
		return nil, err
	}

	return result.Credentials, nil
}
