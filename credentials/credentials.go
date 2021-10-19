// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package credentials contains functionality for retrieving and manipulating
// AWS credentials.
package credentials

import (
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
