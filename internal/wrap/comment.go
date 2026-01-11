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
