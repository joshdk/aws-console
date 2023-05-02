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
	"account":  "https://console.aws.amazon.com/billing/home#/account",
	"billing":  "https://console.aws.amazon.com/billing/home",
	"console":  "https://{{.region}}.console.aws.amazon.com/console/home?region={{.region}}",
	"ec2":      "https://{{.region}}.console.aws.amazon.com/ec2/home?region={{.region}}#Home:",
	"ecr":      "https://{{.region}}.console.aws.amazon.com/ecr/repositories?region={{.region}}",
	"eks":      "https://{{.region}}.console.aws.amazon.com/eks/home?region={{.region}}#/clusters",
	"groups":   "https://console.aws.amazon.com/iamv2/home#/groups",
	"home":     "https://{{.region}}.console.aws.amazon.com/console/home?region={{.region}}",
	"iam":      "https://console.aws.amazon.com/iamv2/home#/home",
	"kms":      "https://{{.region}}.console.aws.amazon.com/kms/home?region={{.region}}#/kms/keys",
	"org":      "https://console.aws.amazon.com/organizations/v2/home/dashboard",
	"policies": "https://console.aws.amazon.com/iamv2/home#/policies",
	"r53":      "https://console.aws.amazon.com/route53/v2/hostedzones#",
	"rds":      "https://{{.region}}.console.aws.amazon.com/rds/home?region={{.region}}#databases:",
	"roles":    "https://console.aws.amazon.com/iamv2/home#/roles",
	"s3":       "https://s3.console.aws.amazon.com/s3/buckets?region={{.region}}",
	"support":  "https://support.console.aws.amazon.com/support/home",
	"users":    "https://console.aws.amazon.com/iamv2/home#/users",
	"vpn":      "https://{{.region}}.console.aws.amazon.com/vpc/home?region={{.region}}#ClientVPNEndpoints:",
}

func resolveLocationAlias(alias, region string) (string, bool) {
	var template string

	if strings.HasPrefix(alias, "https://") {
		// Use the alias directly as it was already a URL.
		template = alias
	} else if result, found := locations[alias]; found {
		// Resolve the alias into a URL.
		template = result
	} else {
		// Alias could not be resolved
		return "", false
	}

	// Replace all region placeholders.
	return strings.ReplaceAll(template, "{{.region}}", region), true
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
