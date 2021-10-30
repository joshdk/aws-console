// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package cmd

// policies is a list of aliases that can be resolved to IAM policy ARNs. Used
// for attaching a policy to a federated user session.
var policies = map[string]string{ // nolint:gochecknoglobals
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

// redirects is a list of aliases that can be resolved to URLs in the AWS
// Console. Used for quickly redirecting the user to the desired service after
// logging in.
var redirects = map[string]string{ // nolint:gochecknoglobals
	"billing": "https://console.aws.amazon.com/billing/home",
	"console": "https://console.aws.amazon.com/console/home",
	"ec2":     "https://console.aws.amazon.com/ec2/v2/home",
	"eks":     "https://console.aws.amazon.com/eks/home",
	"home":    "https://console.aws.amazon.com/console/home",
	"iam":     "https://console.aws.amazon.com/iam/home",
	"s3":      "https://s3.console.aws.amazon.com/s3/home",
}

func resolveRedirectAlias(alias string) string {
	if result, found := redirects[alias]; found {
		return result
	}

	return alias
}
