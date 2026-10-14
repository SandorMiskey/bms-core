// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Server entry point placeholder.
// The previous bmsd implementation was removed (R7-23) because none of it is
// reused: the Slice 1 server is rebuilt from the approved decisions. Until
// then this placeholder keeps the module buildable, so the just recipes for
// build, vet, test and lint run cleanly. It prints a notice and exits with
// status 1, because it serves nothing yet.

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "bmsd: not implemented yet; the Slice 1 server is being rebuilt")
	os.Exit(1)
}

// vim: set ts=4 sw=4 noet:
