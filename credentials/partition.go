// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package credentials

import (
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var partitionURLs = map[string]struct {
	consoleDomain string
	federationURL string
}{
	"aws": {
		consoleDomain: "console.aws.amazon.com",
		federationURL: "https://signin.aws.amazon.com/federation",
	},
	"aws-cn": {
		// This partition has not been tested.
		consoleDomain: "console.amazonaws.cn",
		federationURL: "https://signin.amazonaws.cn/federation",
	},
	"aws-us-gov": {
		consoleDomain: "console.amazonaws-us-gov.com",
		federationURL: "https://signin.amazonaws-us-gov.com/federation",
	},
}

// ResolveRegionPartition uses the given AWS region to determine the corresponding AWS partition, Console URL, and federation URL.
func ResolveRegionPartition(region string) (string, string, string, bool) {
	partition := "aws"
	if endpoint, err := sts.NewDefaultEndpointResolver().ResolveEndpoint(region, sts.EndpointResolverOptions{}); err == nil {
		partition = endpoint.PartitionID
	}

	if urls, ok := partitionURLs[partition]; ok {
		return partition, urls.consoleDomain, urls.federationURL, true
	}

	return "", "", "", false
}
