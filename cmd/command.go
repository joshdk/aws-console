// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package cmd contains functionality for supporting the aws-console cli.
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/joshdk/aws-console/console"
	"github.com/joshdk/aws-console/credentials"
	"github.com/joshdk/aws-console/qr"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"jdk.sh/meta"
)

type flags struct {
	// browser indicates that the login URL should be opened with the system's
	// default browser.
	browser bool

	// clipboard indicates that the login URL should be copied to the system
	// clipboard.
	clipboard bool

	// duration is how long the AWS Console session should last before expiring.
	duration time.Duration

	// federateName is the identifier used for temporary security credentials
	// when federating an IAM user.
	federateName string

	// federatePolicy is the policy ARN to attach when federating an IAM user.
	federatePolicy string

	// profile is the name of profile used for retrieving credentials from the
	// AWS cli config files.
	profile string

	// qr indicates that the login URL should be rendered as a QR code.
	qr bool

	// qrSize is the width in pixels of the rendered QR code.
	qrSize int

	// redirect is the AWS Console page to redirect to after logging in.
	redirect string

	// userAgent is the user agent to use when making API calls.
	userAgent string
}

// Command returns a complete handler for the aws-console cli.
func Command() *cobra.Command {
	var flags flags

	cmd := &cobra.Command{
		Use:     "aws-console [profile|-]",
		Long:    "aws-console - Generate temporary login URLs for the AWS Console",
		Version: "-",

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
			// Obtain credentials from either STDIN or a named AWS cli profile.
			var creds *sts.Credentials
			var err error
			if flags.profile == "-" {
				// Retrieve credentials from JSON via STDIN.
				creds, err = credentials.FromReader(os.Stdin)
			} else {
				// Retrieve credentials from the AWS cli config files.
				creds, err = credentials.FromConfig(flags.profile)
			}
			if err != nil {
				return err
			}

			// If the named profile was configured with user credentials
			// (opposed to a role), then the user must be federated before an
			// AWS Console login url can be generated.
			federatePolicy := resolvePolicyAlias(flags.federatePolicy)
			creds, err = credentials.FederateUser(creds, flags.federateName, federatePolicy, flags.duration, flags.userAgent)
			if err != nil {
				return err
			}

			// Generate a login URL for the AWS console.
			redirect := resolveRedirectAlias(flags.redirect)
			url, err := console.GenerateLoginURL(creds, flags.duration, redirect, flags.userAgent)
			if err != nil {
				return err
			}

			switch {
			case flags.qr:
				// Render the login url as a QR code.
				return qr.Render(os.Stdout, url.String(), flags.qrSize)
			case flags.browser:
				// Open the login url with the default browser.
				return browser.OpenURL(url.String())
			case flags.clipboard:
				// Copy the login url to the system clipboard.
				fmt.Println("Copied AWS Console login URL to clipboard.") // nolint:forbidigo

				return clipboard.WriteAll(url.String())
			default:
				// Print the login url.
				fmt.Println(url.String()) // nolint:forbidigo

				return nil
			}
		},
	}

	// Define -b/--browser flag.
	cmd.Flags().BoolVarP(&flags.browser, "browser", "b",
		false,
		"open login URL with default browser")

	// Define -c/--clipboard flag.
	cmd.Flags().BoolVarP(&flags.clipboard, "clipboard", "c",
		false,
		"copy login URL to clipboard")

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
		"admin",
		"policy ARN attached to federated user session")

	// Define -q/--qr flag.
	cmd.Flags().BoolVarP(&flags.qr, "qr", "q",
		false,
		"render login URL as a QR code")

	// Define -s/--qr-size flag.
	cmd.Flags().IntVarP(&flags.qrSize, "qr-size", "s",
		780, // nolint:gomnd
		"width in pixels of QR code")

	// Define -r/--redirect flag.
	cmd.Flags().StringVarP(&flags.redirect, "redirect", "r",
		"home",
		"console page to redirect to after logging in")

	// Define -A/--user-agent flag.
	cmd.Flags().StringVarP(&flags.userAgent, "user-agent", "A",
		versionFmt("joshdk/aws-console", " %s (%s)", meta.Version(), meta.ShortSHA()),
		"user agent to use for http requests")

	cmd.Example = `  Generate a login url for the default profile:
  $ aws-console

  Generate a login url for the "production" profile:
  $ aws-console production

  Generate a login url from the output of the aws cli:
  $ aws sts assume-role … | aws-console -

  Open url with the default browser:
  $ aws-console --browser

  Redirect to IAM service after logging in:
  $ aws-console --redirect iam

  Display a QR code for the login url:
  $ aws-console --qr

  Save QR code to a file:
  $ aws-console --qr > qr.png`

	// Add a custom usage footer template.
	cmd.SetUsageTemplate(cmd.UsageTemplate() + versionFmt(
		"\nInfo:\n"+
			"  https://github.com/joshdk/aws-console\n",
		"  %s (%s) built on %v\n",
		meta.Version(), meta.ShortSHA(), meta.DateFormat(time.RFC3339),
	))

	// Set a custom version template.
	cmd.SetVersionTemplate(versionFmt(
		"homepage: https://github.com/joshdk/aws-auth\n"+
			"author:   Josh Komoroske\n"+
			"license:  MIT\n",
		"version:  %s\n"+
			"sha:      %s\n"+
			"date:     %s\n",
		meta.Version(), meta.ShortSHA(), meta.DateFormat(time.RFC3339),
	))

	return cmd
}

// versionFmt returns the given literal, as well as a formatted string if
// version metadata is set.
func versionFmt(literal, format string, a ...interface{}) string {
	if meta.Version() == "" {
		return literal
	}

	return literal + fmt.Sprintf(format, a...)
}
