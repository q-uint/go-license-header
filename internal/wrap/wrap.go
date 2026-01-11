// SPDX-License-Identifier: LGPL-3.0-or-later

// Copyright (c) 2026 Quint Daenen.
// This file is part of go-license-header.
//
// go-license-header is free software: you can redistribute it and/or modify it
// under the terms of the GNU Lesser Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option) any
// later version.
//
// go-license-header is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser Public License for more
// details.
//
// You should have received a copy of the GNU Lesser Public License along with
// go-license-header. If not, see <https://www.gnu.org/licenses/>.

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
