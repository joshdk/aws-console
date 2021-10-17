// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package credentials contains functionality for retrieving and manipulating
// AWS credentials.
package credentials

import (
	"github.com/aws/aws-sdk-go/aws"
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
