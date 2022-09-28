// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package cmd

import "strings"

// locations is a list of aliases that can be resolved to URLs in the AWS
// Console. Used for quickly redirecting the user to the desired service after
// logging in.
var locations = map[string]string{ //nolint:gochecknoglobals
	"billing":  "https://console.aws.amazon.com/billing/home",
	"console":  "https://console.aws.amazon.com/console/home",
	"ec2":      "https://console.aws.amazon.com/ec2/v2/home",
	"ecr":      "https://console.aws.amazon.com/ecr/repositories",
	"eks":      "https://console.aws.amazon.com/eks/home#/clusters",
	"groups":   "https://console.aws.amazon.com/iamv2/home#/groups",
	"home":     "https://console.aws.amazon.com/console/home",
	"iam":      "https://console.aws.amazon.com/iamv2/home#/home",
	"policies": "https://console.aws.amazon.com/iamv2/home#/policies",
	"r53":      "https://console.aws.amazon.com/route53/v2/hostedzones#",
	"rds":      "https://console.aws.amazon.com/rds/home#databases:",
	"roles":    "https://console.aws.amazon.com/iamv2/home#/roles",
	"s3":       "https://s3.console.aws.amazon.com/s3/home",
	"users":    "https://console.aws.amazon.com/iamv2/home#/users",
}

func resolveLocationAlias(alias string) (string, bool) {
	// The given alias was already a url, so return it unmodified.
	if strings.HasPrefix(alias, "https://") {
		return alias, true
	}

	// Attempt to resolve the given alias into a url.
	if result, found := locations[alias]; found {
		return result, true
	}

	// No url could be resolved.
	return "", false
}

// policies is a list of aliases that can be resolved to IAM policy ARNs. Used
// for attaching a policy to a federated user session.
var policies = map[string]string{ //nolint:gochecknoglobals
	"admin":    "arn:aws:iam::aws:policy/AdministratorAccess",
	"all":      "arn:aws:iam::aws:policy/AdministratorAccess",
	"billing":  "arn:aws:iam::aws:policy/job-function/Billing",
	"readonly": "arn:aws:iam::aws:policy/ReadOnlyAccess",
	"ro":       "arn:aws:iam::aws:policy/ReadOnlyAccess",
}

func resolvePolicyAlias(alias string) string {
	if result, found := policies[alias]; found {
		return result
	}

	return alias
}
