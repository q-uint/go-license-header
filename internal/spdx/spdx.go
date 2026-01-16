// SPDX-License-Identifier: MPL-2.0

// Copyright (c) 2026 Quint Daenen.
// This file is part of go-license-header.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at https://mozilla.org/MPL/2.0/.

package spdx

import (
	"embed"
	_ "embed"
	"fmt"
)

// You can find all licenses here:
// - https://www.mozilla.org/media/MPL/2.0/index.txt
// - https://ftp.gnu.org/gnu/Licenses
//
//go:embed licenses
var licenses embed.FS

type LicenseIdentifier string

var (
	None        LicenseIdentifier = ""
	MPL2        LicenseIdentifier = "MPL-2.0"
	GPL2        LicenseIdentifier = "GPL-2.0-only"
	GPL2Later   LicenseIdentifier = "GPL-2.0-or-later"
	GPL3        LicenseIdentifier = "GPL-3.0-only"
	GPL3Later   LicenseIdentifier = "GPL-3.0-or-later"
	LGPL21      LicenseIdentifier = "LGPL-2.1-only"
	LGPL21Later LicenseIdentifier = "LGPL-2.1-or-later"
	LGPL3       LicenseIdentifier = "LGPL-3.0-only"
	LGPL3Later  LicenseIdentifier = "LGPL-3.0-or-later"
	AGPL3       LicenseIdentifier = "AGPL-3.0-only"
	AGPL3Later  LicenseIdentifier = "AGPL-3.0-or-later"
)

func (l LicenseIdentifier) name() (string, error) {
	switch l {
	case MPL2:
		return "Mozilla Public License", nil
	case GPL2, GPL2Later, GPL3, GPL3Later:
		return "GNU General Public License", nil
	case LGPL21, LGPL21Later, LGPL3, LGPL3Later:
		return "GNU Lesser Public License", nil
	case AGPL3, AGPL3Later:
		return "GNU Affero Public License", nil
	default:
		return "", fmt.Errorf("unsupported SPDX License Identifier: %q", l)
	}
}

func (l LicenseIdentifier) version() (string, error) {
	switch l {
	case GPL2:
		return "version 2", nil
	case LGPL21:
		return "version 2.1", nil
	case GPL3, LGPL3, AGPL3:
		return "version 3", nil
	case GPL2Later:
		return "either version 2 of the License, or (at your option) any later version", nil
	case LGPL21Later:
		return "either version 2.1 of the License, or (at your option) any later version", nil
	case GPL3Later, LGPL3Later, AGPL3Later:
		return "either version 3 of the License, or (at your option) any later version", nil
	default:
		return "", fmt.Errorf("unsupported SPDX License Identifier: %q", l)
	}
}

func (l LicenseIdentifier) path() (string, error) {
	switch l {
	case MPL2:
		return "mpl-2.0.txt", nil
	case GPL2, GPL2Later:
		return "gpl-2.0.txt", nil
	case GPL3, GPL3Later:
		return "gpl-3.0.txt", nil
	case LGPL21, LGPL21Later:
		return "lgpl-2.1.txt", nil
	case LGPL3, LGPL3Later:
		return "lgpl-3.0.txt", nil
	case AGPL3, AGPL3Later:
		return "agpl-3.0.txt", nil
	default:
		return "", fmt.Errorf("unsupported SPDX License Identifier: %q", l)
	}
}

func (l LicenseIdentifier) Identifier() string {
	return fmt.Sprintf("SPDX-License-Identifier: %s", l)
}

func (l LicenseIdentifier) Header(year int, copyright, name string) string {
	if l == MPL2 {
		return fmt.Sprintf(`
		Copyright (c) %d %s.
		This file is part of %s.
		
		This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0. If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
	`,
			year, copyright,
			name,
		)
	}

	licenseName, _ := l.name()
	version, _ := l.version()
	return fmt.Sprintf(`
		Copyright (c) %d %s.
		This file is part of %s.
		
		%s is free software: you can redistribute it and/or modify it under the terms of the %s as published by the Free Software Foundation, %s.
		
		%s is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the %s for more details.
		
		You should have received a copy of the %s along with %s. If not, see <https://www.gnu.org/licenses/>.
	`,
		year, copyright,
		name,
		name, licenseName, version,
		name, licenseName,
		licenseName, name,
	)
}

func (l LicenseIdentifier) License() ([]byte, error) {
	p, err := l.path()
	if err != nil {
		return nil, err
	}
	return licenses.ReadFile(fmt.Sprintf("licenses/%s", p))
}
