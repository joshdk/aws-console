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
	"account":    "https://{region}.{console}/billing/home?region={region}#/account",
	"billing":    "https://{region}.{console}/costmanagement/home?region={region}#/home",
	"cloudfront": "https://{region}.{console}/cloudfront/v4/home?region={region}#/distributions",
	"cloudtrail": "https://{region}.{console}/cloudtrailv2/home?region={region}#/dashboard",
	"cloudwatch": "https://{region}.{console}/cloudwatch/home?region={region}#home:",
	"console":    "https://{region}.{console}/console/home?region={region}",
	"ec2":        "https://{region}.{console}/ec2/home?region={region}#Instances:",
	"ecr":        "https://{region}.{console}/ecr/private-registry/repositories?region={region}",
	"ecs":        "https://{region}.{console}/ecs/v2/clusters?region={region}",
	"eip":        "https://{region}.{console}/vpcconsole/home?region={region}#Addresses:",
	"eks":        "https://{region}.{console}/eks/clusters?region={region}",
	"groups":     "https://{region}.{console}/iam/home?region={region}#/groups",
	"home":       "https://{region}.{console}/console/home?region={region}",
	"iam":        "https://{region}.{console}/iam/home?region={region}#/home",
	"kms":        "https://{region}.{console}/kms/home?region={region}#/kms/home",
	"org":        "https://{region}.{console}/organizations/v2/home?region={region}",
	"policies":   "https://{region}.{console}/iam/home?region={region}#/policies",
	"r53":        "https://{region}.{console}/route53/v2/hostedzones?region={region}",
	"rds":        "https://{region}.{console}/rds/home?region={region}#databases:",
	"roles":      "https://{region}.{console}/iam/home?region={region}#/roles",
	"s3":         "https://{region}.{console}/s3/buckets?region={region}",
	"support":    "https://support.{console}/support/home?region={region}#/case/history",
	"users":      "https://{region}.{console}/iam/home?region={region}#/users",
	"vpc":        "https://{region}.{console}/vpcconsole/home?region={region}#vpcs:",
	"vpn":        "https://{region}.{console}/vpcconsole/home?region={region}#ClientVPNEndpoints:",
}

func resolveLocationAlias(alias, consoleDomain, region string) (string, bool) {
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

	// Replace all the placeholders.
	return strings.NewReplacer(
		"{console}", consoleDomain,
		"{region}", region,
	).Replace(template), true
}

// policies is a list of aliases that can be resolved to IAM policy ARNs. Used
// for attaching a policy to a federated user session.
var policies = map[string]string{ //nolint:gochecknoglobals
	"admin":    "arn:{partition}:iam::aws:policy/AdministratorAccess",
	"all":      "arn:{partition}:iam::aws:policy/AdministratorAccess",
	"billing":  "arn:{partition}:iam::aws:policy/job-function/Billing",
	"readonly": "arn:{partition}:iam::aws:policy/ReadOnlyAccess",
	"ro":       "arn:{partition}:iam::aws:policy/ReadOnlyAccess",
}

func resolvePolicyAlias(alias, partition string) string {
	template := alias

	if result, found := policies[alias]; found {
		// Resolve the alias into a URL.
		template = result
	}

	return strings.ReplaceAll(template, "{partition}", partition)
}
