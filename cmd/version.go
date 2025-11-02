// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

package cmd

import (
	_ "embed"

	"github.com/joshdk/buildversion"
)

var (
	version   string
	revision  string
	timestamp string

	_ = buildversion.Override(version, revision, timestamp)

	//go:embed files/usage-template.tmpl
	usageTemplate string

	//go:embed files/version-template.tmpl
	versionTemplate string

	//go:embed files/examples.txt
	exampleText string
)
