// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package cmd

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "aws-console [flags…] [profile]",
		Long: "aws-console - Generate temporary login URLs for the AWS Console",

		SilenceUsage:  true,
		SilenceErrors: true,

		RunE: func(*cobra.Command, []string) error {
			return nil
		},
	}

	return cmd
}
