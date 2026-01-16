# Go License Header

`go-license-header` is a simple CLI tool for checking and generating license headers in Go files. It supports MPL, GPL, LGPL and AGPL licenses.


```text
go-license-header checks and generates license headers in Go files.

Usage:
  go-license-header help
    Prints out this help message.

  go-license-header header [-spdx] -n [-y] -c
    Prints out an example *GPL file header.

    -spdx=<id> SPDX license identifier. (default: MPL-2.0)
    -n=<name>  Project name. (required)
    -y=<year>  Year. (default: {current year})
    -c=<copy>  Copyright holder. (required)

  go-license-header license [-spdx] [-o] [-d]
    Prints out the license files linked to the SPDX License Identifier.
    If [-o] is specified, it will also try to write it to that directory.

    -spdx=<id>  SPDX license identifier. (default: MPL-2.0)
    -o=<dir>    Output directory.
    -d          Dry-run flag, will print the write locations.

  go-license-header check [-spdx] [-p] [-r] [-d]
    Checks if all files in the given path have a license header.

    -spdx=<id>  SPDX license identifier. (default: MPL-2.0)
    -p=<path>   The path to check, can be either a file or directory (with -r).
    -r          Whether to recursively walk the directory.
    -d          Dry-run flag, will print error but not exit(1).

  go-license-header run [-spdx] [-p] [-r] [-d]
    Writes license headers to files that don't have one yet.

    -spdx=<id>  SPDX license identifier. (default: MPL-2.0)
    -p=<path>   The path to check, can be either a file or directory (with -r).
    -r          Whether to recursively walk the directory.
    -d          Dry-run flag, will print error but not write the headers.
```
