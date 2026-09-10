// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Server and database configuration.
// This file defines runtime environment enums plus server and database config
// structs. The types map to [server] and [database] sections to control
// instance identity and database connectivity.

package config

// Environment selects the server runtime mode. {{{

type Environment string

const (
	EnvLocal  Environment = "local"
	EnvRemote Environment = "remote"
)

// }}}
// DatabaseDriver selects the database engine. {{{

type DatabaseDriver string

const (
	DriverSQLite   DatabaseDriver = "sqlite"
	DriverLibSQL   DatabaseDriver = "libsql"
	DriverPostgres DatabaseDriver = "postgres"
)

// }}}
// ServerConfig holds top-level server settings. {{{

type ServerConfig struct {
	Environment Environment `toml:"environment"` // Runtime mode (`local` or `remote`).
	ID          string      `toml:"id"`          // Instance identifier (ULID string).
}

// }}}
// DatabaseConfig configures database connectivity and migrations. {{{

type DatabaseConfig struct {
	DSN             string         `toml:"dsn"`                // Connection string for the selected driver.
	Driver          DatabaseDriver `toml:"driver"`             // Database engine (`sqlite`, `libsql`, or `postgres`).
	Migrations      string         `toml:"migrations"`         // Migrations directory.
	AuthToken       string         `toml:"auth_token"`         // Authentication token for remote drivers.
	BusyTimeout     string         `toml:"busy_timeout"`       // Busy timeout duration for database locks.
	MaxOpenConns    int            `toml:"max_open_conns"`     // Maximum number of open connections.
	MaxIdleConns    int            `toml:"max_idle_conns"`     // Maximum number of idle connections.
	ConnMaxLifetime string         `toml:"conn_max_lifetime"`  // Maximum lifetime for pooled connections.
	ConnMaxIdleTime string         `toml:"conn_max_idle_time"` // Maximum idle time for pooled connections.
}

// }}}

// vim: set ts=4 sw=4 noet:
