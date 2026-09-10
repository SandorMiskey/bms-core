// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Database connections.
// This file opens sqlite/libsql/postgres connections, applies pool settings,
// and injects SQLite connection pragmas into the DSN based on the resolved
// configuration.

package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"github.com/SandorMiskey/bms-core/internal/config"
	"github.com/SandorMiskey/bms-core/internal/errtext"
)

// Database connections. {{{

func Open(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	db, err := openWithDriver(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errtext.ErrDatabaseOpenFailed, err)
	}

	if err := configurePool(db, cfg); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("%s: %w", errtext.ErrDatabaseOpenFailed, err)
	}

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("%s: %w", errtext.ErrDatabasePingFailed, err)
	}

	return db, nil
}

func openWithDriver(ctx context.Context, cfg config.DatabaseConfig) (*sql.DB, error) {
	switch cfg.Driver {
	case config.DriverSQLite:
		return openSQLite(cfg)
	case config.DriverLibSQL:
		return openLibSQL(ctx, cfg)
	case config.DriverPostgres:
		return openPostgres(cfg)
	default:
		return nil, fmt.Errorf("%s: %s", errtext.ErrDatabaseDriverUnsupported, cfg.Driver)
	}
}

func openSQLite(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn, err := sqliteDSN(cfg.DSN, cfg.BusyTimeout)
	if err != nil {
		return nil, err
	}
	return sql.Open("sqlite", dsn)
}

func openPostgres(cfg config.DatabaseConfig) (*sql.DB, error) {
	return sql.Open("pgx", cfg.DSN)
}

// }}}

// SQLite connection pragmas. {{{
// Connection pragmas must travel in the DSN. `database/sql` hands out pooled
// connections, so issuing `PRAGMA busy_timeout` through ExecContext configures
// only whichever single connection happens to serve that statement; every other
// connection in the pool keeps the SQLite default of 0 and fails instantly on
// lock contention. Encoding the pragma in the DSN makes the driver apply it to
// each connection it opens.

const (
	sqliteMemoryDSN    = ":memory:"
	sqliteMemoryURIDSN = "file::memory:"
)

func sqliteDSN(dsn string, busyTimeout string) (string, error) {
	if busyTimeout == "" {
		return dsn, nil
	}

	timeout, err := time.ParseDuration(busyTimeout)
	if err != nil {
		return "", err
	}
	if timeout < 0 {
		return "", fmt.Errorf("%s: %s", errtext.ErrDatabaseBusyTimeoutInvalid, busyTimeout)
	}

	// The bare `:memory:` form carries no query component; its URI equivalent does.
	if dsn == sqliteMemoryDSN {
		dsn = sqliteMemoryURIDSN
	}

	separator := "?"
	if strings.Contains(dsn, "?") {
		separator = "&"
	}

	return fmt.Sprintf("%s%s_pragma=busy_timeout(%d)", dsn, separator, timeout.Milliseconds()), nil
}

// }}}

// Pool settings. {{{

func configurePool(db *sql.DB, cfg config.DatabaseConfig) error {
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != "" {
		value, err := time.ParseDuration(cfg.ConnMaxLifetime)
		if err != nil {
			return err
		}
		db.SetConnMaxLifetime(value)
	}
	if cfg.ConnMaxIdleTime != "" {
		value, err := time.ParseDuration(cfg.ConnMaxIdleTime)
		if err != nil {
			return err
		}
		db.SetConnMaxIdleTime(value)
	}
	return nil
}

// }}}

// vim: set ts=4 sw=4 noet:
