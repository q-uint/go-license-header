// SPDX-License-Identifier: MPL-2.0

// Copyright (c) 2026 Quint Daenen.
// This file is part of go-license-header.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at https://mozilla.org/MPL/2.0/.

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
