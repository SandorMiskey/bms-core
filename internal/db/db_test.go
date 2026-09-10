// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Database connection tests.
// This file verifies SQLite DSN pragma injection, that the busy timeout reaches
// every pooled connection, pool setting validation, and driver dispatch errors.

package db

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/SandorMiskey/bms-core/internal/config"
)

// SQLite DSN pragma tests. {{{

func TestSQLiteDSN(t *testing.T) {
	cases := []struct {
		name        string
		dsn         string
		busyTimeout string
		want        string
		wantErr     bool
	}{
		{
			name:        "no busy timeout leaves the dsn untouched",
			dsn:         "file:bms.db",
			busyTimeout: "",
			want:        "file:bms.db",
		},
		{
			name:        "dsn without a query gains one",
			dsn:         "file:bms.db",
			busyTimeout: "5s",
			want:        "file:bms.db?_pragma=busy_timeout(5000)",
		},
		{
			name:        "dsn with a query is appended to",
			dsn:         "file:bms.db?mode=rwc&cache=shared",
			busyTimeout: "250ms",
			want:        "file:bms.db?mode=rwc&cache=shared&_pragma=busy_timeout(250)",
		},
		{
			name:        "bare memory dsn is rewritten to its uri form",
			dsn:         ":memory:",
			busyTimeout: "1s",
			want:        "file::memory:?_pragma=busy_timeout(1000)",
		},
		{
			name:        "zero busy timeout is allowed",
			dsn:         "file:bms.db",
			busyTimeout: "0s",
			want:        "file:bms.db?_pragma=busy_timeout(0)",
		},
		{
			name:        "unparsable busy timeout fails",
			dsn:         "file:bms.db",
			busyTimeout: "soon",
			wantErr:     true,
		},
		{
			name:        "negative busy timeout fails",
			dsn:         "file:bms.db",
			busyTimeout: "-1s",
			wantErr:     true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := sqliteDSN(testCase.dsn, testCase.busyTimeout)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got dsn %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != testCase.want {
				t.Fatalf("expected dsn %q, got %q", testCase.want, got)
			}
		})
	}
}

// The busy timeout is a per-connection setting. This guards the regression where
// it was applied with ExecContext and therefore reached only one pooled
// connection while the rest kept the SQLite default of 0.
func TestOpenAppliesBusyTimeoutToEveryPooledConnection(t *testing.T) {
	const poolSize = 3

	handle, err := Open(context.Background(), config.DatabaseConfig{
		Driver:       config.DriverSQLite,
		DSN:          "file:busytimeout?mode=memory&cache=shared",
		BusyTimeout:  "7s",
		MaxOpenConns: poolSize,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer handle.Close()

	// Hold every connection open at once so the pool is forced to create each one.
	connections := make([]*sql.Conn, 0, poolSize)
	defer func() {
		for _, connection := range connections {
			connection.Close()
		}
	}()

	for index := 0; index < poolSize; index++ {
		connection, err := handle.Conn(context.Background())
		if err != nil {
			t.Fatalf("connection %d: %v", index, err)
		}
		connections = append(connections, connection)

		var busyTimeout int
		row := connection.QueryRowContext(context.Background(), "PRAGMA busy_timeout")
		if err := row.Scan(&busyTimeout); err != nil {
			t.Fatalf("connection %d pragma read: %v", index, err)
		}
		if busyTimeout != 7000 {
			t.Fatalf("connection %d: expected busy_timeout 7000, got %d", index, busyTimeout)
		}
	}
}

// }}}

// Open and pool tests. {{{

func TestOpenRejectsUnsupportedDriver(t *testing.T) {
	handle, err := Open(context.Background(), config.DatabaseConfig{
		Driver: config.DatabaseDriver("mysql"),
		DSN:    "file:bms.db",
	})
	if err == nil {
		handle.Close()
		t.Fatal("expected an error for an unsupported driver")
	}
	if !strings.Contains(err.Error(), "mysql") {
		t.Fatalf("expected the error to name the driver, got %q", err.Error())
	}
}

func TestOpenRejectsInvalidPoolDurations(t *testing.T) {
	cases := []struct {
		name string
		cfg  config.DatabaseConfig
	}{
		{
			name: "unparsable conn_max_lifetime",
			cfg: config.DatabaseConfig{
				Driver:          config.DriverSQLite,
				DSN:             ":memory:",
				ConnMaxLifetime: "forever",
			},
		},
		{
			name: "unparsable conn_max_idle_time",
			cfg: config.DatabaseConfig{
				Driver:          config.DriverSQLite,
				DSN:             ":memory:",
				ConnMaxIdleTime: "a while",
			},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			handle, err := Open(context.Background(), testCase.cfg)
			if err == nil {
				handle.Close()
				t.Fatal("expected an error")
			}
		})
	}
}

func TestOpenSucceedsWithNilContext(t *testing.T) {
	//nolint:staticcheck // a nil context is defended against on purpose.
	handle, err := Open(nil, config.DatabaseConfig{
		Driver: config.DriverSQLite,
		DSN:    ":memory:",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer handle.Close()

	if err := handle.Ping(); err != nil {
		t.Fatalf("unexpected ping error: %v", err)
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
