// SPDX-License-Identifier: MPL-2.0

// Copyright (c) 2026 Quint Daenen.
// This file is part of go-license-header.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at https://mozilla.org/MPL/2.0/.

package wrap

import (
	"fmt"
	"strings"
)

func Comment(s string) string {
	var sb strings.Builder
	for l := range strings.SplitSeq(strings.TrimSpace(s), "\n") {
		if l == "" {
			fmt.Fprintf(&sb, "//\n")
		} else {
			fmt.Fprintf(&sb, "// %s\n", strings.TrimSpace(l))
		}
	}
	return sb.String()
}
