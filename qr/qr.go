// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package qr contains functionality for generating and displaying QR codes.
package qr

import (
	"encoding/binary"
	"os"

	"github.com/skip2/go-qrcode"
)

// Render generates a QR code from the given image and streams it to STDOUT.
// The stream is intended to then be redirected to a file.
func Render(url string, size int) error {
	png, err := qrcode.Encode(url, qrcode.Low, size)
	if err != nil {
		return err
	}

	// Write the raw PNG bytes to STDOUT.
	return binary.Write(os.Stdout, binary.LittleEndian, png)
}
