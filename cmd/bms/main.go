// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Client entry point placeholder.
// The previous bms implementation was removed (R7-23) because none of it is
// reused: the Slice 1 TUI and CLI are rebuilt from the approved decisions.
// Until then this placeholder keeps the module buildable, so the just recipes
// for build, vet, test and lint run cleanly. It prints a notice and exits with
// status 1, because it does nothing yet.

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "bms: not implemented yet; the Slice 1 client is being rebuilt")
	os.Exit(1)
}

// vim: set ts=4 sw=4 noet:
