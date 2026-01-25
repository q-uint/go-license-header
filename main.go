// SPDX-License-Identifier: MPL-2.0

// Copyright (c) 2026 Quint Daenen.
// This file is part of go-license-header.
//
// This Source Code Form is subject to the terms of the Mozilla Public License,
// v. 2.0. If a copy of the MPL was not distributed with this file, You can
// obtain one at https://mozilla.org/MPL/2.0/.

package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/q-uint/go-license-header/internal/spdx"
	"github.com/q-uint/go-license-header/internal/wrap"
)

const (
	LICENSE_HEADER_PATH = ".license.header"
	LINE_WIDTH          = 80 - 3 // len("// ")
)

func main() {
	name := filepath.Base(os.Args[0])
	log.SetPrefix(fmt.Sprintf("%s: ", name))
	log.SetFlags(0)

	if len(os.Args) < 2 {
		log.Fatalf("missing command, run '%s help' for more information", name)
	}

	switch cmd := os.Args[1]; cmd {
	case "help":
		if _, err := fmt.Printf(usage, name); err != nil {
			log.Fatalf("failed printing usage information: %v", err)
		}
	case "header":
		flags := flag.NewFlagSet("", flag.ContinueOnError)

		var licenseIdentifier string
		flags.StringVar(&licenseIdentifier, "spdx", "MPL-2.0", "SPDX License Identifier")
		var year uint64
		flags.Uint64Var(&year, "y", uint64(time.Now().Year()), "Year")
		var copyright string
		flags.StringVar(&copyright, "c", "", "Copyright Holder")
		var name string
		flags.StringVar(&name, "n", "", "Project Name")

		if err := flags.Parse(os.Args[2:]); err != nil {
			log.Fatalf("invalid options, run '%s help' for more information", name)
		}

		if copyright == "" {
			log.Fatal("no Copyright Holder (-c) specified")
		}
		if name == "" {
			log.Fatal("no Name (-n) specified")
		}

		header := spdx.LicenseIdentifier(licenseIdentifier).Header(int(year), copyright, name)
		fmt.Println(wrap.Wrap(strings.TrimSpace(header), LINE_WIDTH))
	case "license":
		flags := flag.NewFlagSet("", flag.ContinueOnError)

		var licenseIdentifier string
		flags.StringVar(&licenseIdentifier, "spdx", "MPL-2.0", "SPDX License Identifier")
		var out string
		flags.StringVar(&out, "o", "", "")
		var dryRun bool
		flags.BoolVar(&dryRun, "d", false, "")

		if err := flags.Parse(os.Args[2:]); err != nil {
			log.Fatalf("invalid options, run '%s help' for more information", name)
		}

		if out == "" {
			raw, err := spdx.LicenseIdentifier(licenseIdentifier).License()
			if err != nil {
				log.Fatalf("could not get license file: %v", err)
			}
			fmt.Println(string(raw))
		} else {
			stat, err := os.Stat(out)
			if err != nil {
				log.Fatalf("%q not found: %v", out, err)
			}
			if !stat.IsDir() {
				log.Fatalf("%q is not a dir", out)
			}
			switch id := spdx.LicenseIdentifier(licenseIdentifier); id {
			case spdx.MPL2:
				writeLicense(id, path.Join(out, "LICENSE"), dryRun)
			case spdx.LGPL21, spdx.LGPL21Later:
				writeLicense(id, path.Join(out, "COPYING.LESSER"), dryRun)
				writeLicense(spdx.GPL2, path.Join(out, "COPYING"), dryRun)
			case spdx.LGPL3, spdx.LGPL3Later:
				writeLicense(id, path.Join(out, "COPYING.LESSER"), dryRun)
				writeLicense(spdx.GPL3, path.Join(out, "COPYING"), dryRun)
			case spdx.GPL2, spdx.GPL2Later, spdx.GPL3, spdx.GPL3Later, spdx.AGPL3, spdx.AGPL3Later:
				writeLicense(id, path.Join(out, "COPYING"), dryRun)
			default:
				log.Fatalf("unsupported SPDX License Identifier: %q", licenseIdentifier)
			}
		}
	case "check":
		flags := flag.NewFlagSet("", flag.ContinueOnError)

		var licenseIdentifier string
		flags.StringVar(&licenseIdentifier, "spdx", "MPL-2.0", "")
		var path string
		flags.StringVar(&path, "p", ".", "")
		var recursive bool
		flags.BoolVar(&recursive, "r", false, "")
		var dryRun bool
		flags.BoolVar(&dryRun, "d", false, "")
		var year uint64
		flags.Uint64Var(&year, "y", uint64(time.Now().Year()), "Year")
		var copyright string
		flags.StringVar(&copyright, "c", "", "Copyright Holder")
		var name string
		flags.StringVar(&name, "n", "", "Project Name")

		if err := flags.Parse(os.Args[2:]); err != nil {
			log.Fatalf("invalid options, run '%s help' for more information", name)
		}

		var licenseRef []byte
		stat, err := os.Stat(LICENSE_HEADER_PATH)
		if err != nil {
			if copyright == "" || name == "" {
				log.Fatalf("%q not found: %v", LICENSE_HEADER_PATH, err)
			}
			header := spdx.LicenseIdentifier(licenseIdentifier).Header(int(year), copyright, name)
			licenseRef = []byte(wrap.Wrap(strings.TrimSpace(header), LINE_WIDTH))
		} else {
			if stat.IsDir() {
				log.Fatalf("%q not found", LICENSE_HEADER_PATH)
			}
			if licenseRef, err = os.ReadFile(LICENSE_HEADER_PATH); err != nil {
				log.Fatalf("could not open file %q: %v", LICENSE_HEADER_PATH, err)
			}
		}

		dir, file, err := splitDir(path)
		if err != nil {
			log.Fatalf("failed splitting path: %v", err)
		}

		handleFile := func(path string) {
			if err := checkFile(path, spdx.LicenseIdentifier(licenseIdentifier), licenseRef); err != nil {
				if dryRun {
					log.Printf("%s: %v", path, err)
				} else {
					log.Fatalf("file %q is missing the license header: %v", path, err)
				}
			}
		}

		if file == "" && recursive {
			var fn filepath.WalkFunc = func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					log.Fatalf("failed to walk dir: %v", err)
				}
				ext := filepath.Ext(path)
				if ext != ".go" {
					return nil
				}
				handleFile(path)
				return nil
			}
			if err := filepath.Walk(dir, fn); err != nil {
				log.Fatalf("failed walking dir %q: %v", dir, err)
			}
			return
		} else if file != "" {
			handleFile(path)
		} else {
			log.Fatalf("nothing to do on path %q", path)
		}
	case "run":
		flags := flag.NewFlagSet("", flag.ContinueOnError)

		var licenseIdentifier string
		flags.StringVar(&licenseIdentifier, "spdx", "MPL-2.0", "")
		var path string
		flags.StringVar(&path, "p", ".", "")
		var recursive bool
		flags.BoolVar(&recursive, "r", false, "")
		var dryRun bool
		flags.BoolVar(&dryRun, "d", false, "")
		var year uint64
		flags.Uint64Var(&year, "y", uint64(time.Now().Year()), "Year")
		var copyright string
		flags.StringVar(&copyright, "c", "", "Copyright Holder")
		var name string
		flags.StringVar(&name, "n", "", "Project Name")

		if err := flags.Parse(os.Args[2:]); err != nil {
			log.Fatalf("invalid options, run '%s help' for more information", name)
		}

		var licenseRef []byte
		stat, err := os.Stat(LICENSE_HEADER_PATH)
		if err != nil {
			if copyright == "" || name == "" {
				log.Fatalf("%q not found: %v", LICENSE_HEADER_PATH, err)
			}
			header := spdx.LicenseIdentifier(licenseIdentifier).Header(int(year), copyright, name)
			licenseRef = []byte(wrap.Wrap(strings.TrimSpace(header), LINE_WIDTH))
		} else {
			if stat.IsDir() {
				log.Fatalf("%q not found", LICENSE_HEADER_PATH)
			}
			if licenseRef, err = os.ReadFile(LICENSE_HEADER_PATH); err != nil {
				log.Fatalf("could not open file %q: %v", LICENSE_HEADER_PATH, err)
			}
		}

		dir, file, err := splitDir(path)
		if err != nil {
			log.Fatalf("failed splitting path: %v", err)
		}

		handleFile := func(path string) {
			id := spdx.LicenseIdentifier(licenseIdentifier)
			if err := checkFile(path, id, licenseRef); err != nil {
				if dryRun {
					log.Printf("%s: %v", path, err)
				} else {
					orig, err := os.ReadFile(path)
					if err != nil {
						log.Fatalf("could not read file %q: %v", path, err)
					}

					newContent := append(fmt.Appendf(nil,
						"%s\n%s\n",
						wrap.Comment(id.Identifier()),
						wrap.Comment(string(licenseRef)),
					), orig...)
					if err := os.WriteFile(path, newContent, 0644); err != nil {
						log.Fatalf("could not write file %q: %v", path, err)
					}
				}
			}
		}

		if file == "" && recursive {
			var fn filepath.WalkFunc = func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					log.Fatalf("failed to walk dir: %v", err)
				}
				ext := filepath.Ext(path)
				if ext != ".go" {
					return nil
				}
				handleFile(path)
				return nil
			}
			if err := filepath.Walk(dir, fn); err != nil {
				log.Fatalf("failed walking dir %q: %v", dir, err)
			}
			return
		} else if file != "" {
			handleFile(path)
		} else {
			log.Fatalf("nothing to do on path %q", path)
		}
	default:
		log.Fatalf("unknown command: %q", cmd)
	}
}

func checkFile(path string, id spdx.LicenseIdentifier, licenseRef []byte) error {
	f, err := os.OpenFile(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("could not read file %q: %w", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanLine := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	first, ok := scanLine()
	expectedFirst := "// " + id.Identifier()
	if !ok {
		return fmt.Errorf("%s: missing first line", path)
	}
	if first != expectedFirst {
		return fmt.Errorf("%s: expected %q, got %q", path, expectedFirst, first)
	}

	second, ok := scanLine()
	if !ok {
		return fmt.Errorf("%s: missing second line", path)
	}
	if strings.TrimSpace(second) != "" {
		return fmt.Errorf("%s: expected empty second line, got %q", path, second)
	}

	refScanner := bufio.NewScanner(bytes.NewReader(licenseRef))
	for refScanner.Scan() {
		expected := refScanner.Text()

		actual, ok := scanLine()
		if !ok {
			return fmt.Errorf("%s: file ended early", path)
		}

		if expected == "" {
			if actual != "//" {
				return fmt.Errorf("%s: expected %q, got %q", path, "//", actual)
			}
			continue
		}

		actual = strings.TrimPrefix(actual, "// ")
		if actual != expected {
			return fmt.Errorf("%s: expected %q, got %q", path, expected, actual)
		}
	}

	if err := refScanner.Err(); err != nil {
		return fmt.Errorf("scanning license reference: %w", err)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanning file %q: %w", path, err)
	}

	return nil
}

func writeLicense(id spdx.LicenseIdentifier, out string, dryRun bool) {
	raw, err := id.License()
	if err != nil {
		log.Fatalf("could not get license file: %v", err)
	}
	if dryRun {
		fmt.Println(out)
	} else {
		if err := os.WriteFile(out, raw, fs.ModePerm); err != nil {
			log.Fatalf("failed writing license to %q: %v", out, err)
		}
	}
}

func splitDir(path string) (dir string, file string, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", "", err
	}
	if stat.IsDir() {
		return path, "", nil
	}
	dir = filepath.Dir(path)
	file = filepath.Base(path)
	return
}
