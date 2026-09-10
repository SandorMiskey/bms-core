//go:build !cgo || windows
// +build !cgo windows

// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// libSQL driver stub.
// This file provides a build-time stub when libSQL is unavailable.

package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/SandorMiskey/bms-core/internal/config"
	"github.com/SandorMiskey/bms-core/internal/errtext"
)

// libSQL driver stub. {{{

func openLibSQL(_ context.Context, _ config.DatabaseConfig) (*sql.DB, error) {
	return nil, fmt.Errorf("%s: libsql driver requires cgo and is unavailable on this platform", errtext.ErrDatabaseDriverUnavailable)
}

// }}}

// vim: set ts=4 sw=4 noet:
