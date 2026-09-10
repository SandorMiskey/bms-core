//go:build cgo && !windows
// +build cgo,!windows

// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// libSQL driver support.
// This file wires the go-libsql driver into database/sql for supported builds.

package db

import (
	"context"
	"database/sql"
	"net/url"

	_ "github.com/tursodatabase/go-libsql"

	"github.com/SandorMiskey/bms-core/internal/config"
)

// libSQL driver support. {{{

func openLibSQL(_ context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn, err := withLibSQLAuthToken(cfg.DSN, cfg.AuthToken)
	if err != nil {
		return nil, err
	}
	return sql.Open("libsql", dsn)
}

func withLibSQLAuthToken(base string, token string) (string, error) {
	if token == "" {
		return base, nil
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	query := parsed.Query()
	query.Set("authToken", token)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

// }}}

// vim: set ts=4 sw=4 noet:
