// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.
// SPDX-License-Identifier: MIT

// Package qr contains functionality for generating and displaying QR codes.
package qr

import (
	"bytes"
	"encoding/binary"
	"image"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/mattn/go-sixel"
	"github.com/skip2/go-qrcode"
)

// Render generates a QR code from the given image and outputs it to the given
// file. If that file is a TTY then the QR code is displayed using sixels, else
// a raw PNG image is written.
func Render(file *os.File, url string, size int) error {
	// Generate a QR code PNG.
	png, err := qrcode.Encode(url, qrcode.Low, size)
	if err != nil {
		return err
	}

	// If the output file is not a TTY (e.g. is being redirected to another
	// file/process) then write the raw PNG bytes.
	if !isatty.IsTerminal(file.Fd()) {
		return binary.Write(file, binary.LittleEndian, png)
	}

	// Decode an image back from the raw PNG data.
	img, _, err := image.Decode(bytes.NewBuffer(png))
	if err != nil {
		return err
	}

	// Render the image to the output file using sixels.
	return sixel.NewEncoder(file).Encode(img)
}
