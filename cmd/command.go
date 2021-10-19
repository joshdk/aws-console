// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package cmd contains functionality for supporting the aws-console cli.
package cmd

import (
	"fmt"
	"time"

	"github.com/joshdk/aws-console/console"
	"github.com/joshdk/aws-console/credentials"
	"github.com/spf13/cobra"
)

type flags struct {
	// duration is how long the AWS Console session should last before expiring.
	duration time.Duration

	// federateName is the identifier used for temporary security credentials
	// when federating an IAM user.
	federateName string

	// federatePolicy is the policy ARN to attach when when federating an IAM
	// user.
	federatePolicy string

	// profile is the name of profile used for retrieving credentials from the
	// AWS cli config files.
	profile string
}

// Command returns a complete handler for the aws-console cli.
func Command() *cobra.Command { // nolint:funlen
	var flags flags

	cmd := &cobra.Command{
		Use:  "aws-console [flags…] [profile]",
		Long: "aws-console - Generate temporary login URLs for the AWS Console",

		SilenceUsage:  true,
		SilenceErrors: true,

		Args: cobra.MaximumNArgs(1),
		PreRun: func(_ *cobra.Command, args []string) {
			if len(args) >= 1 {
				// As a convenience, determine profile name from cli args here.
				flags.profile = args[0]
			}
		},

		RunE: func(*cobra.Command, []string) error {
			// Retrieve credentials from the AWS cli config files.
			creds, err := credentials.FromConfig(flags.profile)
			if err != nil {
				return err
			}

			// If the named profile was configured with user credentials
			// (opposed to a role), then the user must be federated before an
			// AWS Console login url can be generated.
			creds, err = credentials.FederateUser(creds, flags.federateName, flags.federatePolicy, flags.duration)
			if err != nil {
				return err
			}

			// Generate a login URL for the AWS console.
			url, err := console.GenerateLoginURL(creds, flags.duration)
			if err != nil {
				return err
			}

			// Print the login url!
			fmt.Println(url.String()) // nolint:forbidigo

			return nil
		},
	}

	// Define -d/--duration flag.
	cmd.Flags().DurationVarP(&flags.duration, "duration", "d",
		12*time.Hour, // nolint:gomnd
		"session duration")

	// Define -n/--name flag.
	cmd.Flags().StringVarP(&flags.federateName, "name", "n",
		"aws-console",
		"name used for federated user session")

	// Define -p/--policy flag.
	cmd.Flags().StringVarP(&flags.federatePolicy, "policy", "p",
		"arn:aws:iam::aws:policy/AdministratorAccess",
		"policy ARN attached to federated user session")

	return cmd
}
