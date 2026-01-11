package main

var usage = `%s checks and generates license headers in Go files.

Usage:
  %[1]s help
    Prints out this help message.

  %[1]s header [-spdx] -n [-y] -c
    Prints out an example *GPL file header.

    -spdx=<id> SPDX license identifier. (default: LGPL-3.0-or-later)
    -n=<name>  Project name. (required)
    -y=<year>  Year. (default: {current year})
    -c=<copy>  Copyright holder. (required)

  %[1]s license [-spdx] [-o] [-d]
    Prints out the license files linked to the SPDX License Identifier.
    If [-o] is specified, it will also try to write it to that directory.

    -spdx=<id>  SPDX license identifier. (default: LGPL-3.0-or-later)
        -o=<dir>    Output directory.
    -d          Dry-run flag, will print the write locations.

  %[1]s check [-spdx] [-p] [-r] [-d]
    Checks if all files in the given path have a license header.

    -spdx=<id>  SPDX license identifier. (default: LGPL-3.0-or-later)
    -p=<path>   The path to check, can be either a file or directory (with -r).
    -r          Whether to recursively walk the directory.
    -d          Dry-run flag, will print error but not exit(1).

  %[1]s run [-spdx] [-p] [-r] [-d]
    Writes license headers to files that don't have one yet.

    -spdx=<id>  SPDX license identifier. (default: LGPL-3.0-or-later)
    -p=<path>   The path to check, can be either a file or directory (with -r).
    -r          Whether to recursively walk the directory.
    -d          Dry-run flag, will print error but not write the headers.
`
