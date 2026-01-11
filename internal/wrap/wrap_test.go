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

import "testing"

func TestWrap(t *testing.T) {
	if w := Wrap("x x x\n", 1); w != "x\nx\nx\n" {
		t.Errorf("%q", w)
	}
	if w := Wrap("x x x\n", 2); w != "x\nx\nx\n" {
		t.Errorf("%q", w)
	}
	if w := Wrap("x x x\n", 3); w != "x x\nx\n" {
		t.Errorf("%q", w)
	}
	if w := Wrap("x x x\n", 4); w != "x x\nx\n" {
		t.Errorf("%q", w)
	}

	if w := Wrap(" x x x\n x", 4); w != "x x\nx\nx" {
		t.Errorf("%q", w)
	}

	if w := Wrap("xx x\n", 3); w != "xx\nx\n" {
		t.Errorf("%q", w)
	}

	if w := Wrap("x xx\n", 3); w != "x\nxx\n" {
		t.Errorf("%q", w)
	}
	if w := Wrap("x xx\n", 4); w != "x xx\n" {
		t.Errorf("%q", w)
	}

	if w := Wrap("xxx", 3); w != "xxx" {
		t.Errorf("%q", w)
	}
	if w := Wrap("xxx\n", 3); w != "xxx\n" {
		t.Errorf("%q", w)
	}
}
