// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package cmd contains functionality for supporting the aws-console cli.
package cmd

import (
	"github.com/joshdk/aws-console/credentials"
	"github.com/spf13/cobra"
)

type flags struct {
	// profile is the name of profile used for retrieving credentials from the
	// AWS cli config files.
	profile string
}

// Command returns a complete handler for the aws-console cli.
func Command() *cobra.Command {
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
			_, err := credentials.FromConfig(flags.profile)

			return err
		},
	}

	return cmd
}
