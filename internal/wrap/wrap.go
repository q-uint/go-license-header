// SPDX-License-Identifier: MPL-2.0

// Copyright (c) 2026 Quint Daenen.
// This file is part of go-license-header.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at https://mozilla.org/MPL/2.0/.

package wrap

import (
	"bytes"
	"unicode"
)

func Wrap(s string, width uint) string {
	b := bytes.NewBuffer(make([]byte, 0, len(s)))
	var lsize uint

	var wb bytes.Buffer
	var wsize uint

	var sb bytes.Buffer
	var ssize uint

	nl := func() {
		b.WriteRune('\n')
		lsize = 0 // reset line

		sb.Reset()
		ssize = 0 // trailing wsp
	}

	ww := func() {
		if lsize == 0 {
			sb.Reset()
			ssize = 0 // trim leading wsp
		}

		if width < lsize+ssize+wsize {
			nl() // buffer does not fit
		}

		sb.WriteTo(b)
		lsize += ssize
		ssize = 0

		wb.WriteTo(b)
		lsize += wsize
		wsize = 0
	}

	for _, c := range s {
		switch {
		case c == '\n':
			if wsize != 0 {
				ww() // write buffer
			}

			nl()
		case unicode.IsSpace(c):
			if wsize != 0 {
				ww()
			}

			sb.WriteRune(c)
			ssize++
		default:
			wb.WriteRune(c)
			wsize++
		}
	}

	ww()

	return b.String()
}
