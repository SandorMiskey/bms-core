// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

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
	DriverPostgres DatabaseDriver = "postgres"
	DriverSQLite   DatabaseDriver = "sqlite"
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
	DSN        string         `toml:"dsn"`        // Connection string for the selected driver.
	Driver     DatabaseDriver `toml:"driver"`     // Database engine (`sqlite` or `postgres`).
	Migrations string         `toml:"migrations"` // Migrations directory.
}

// }}}

// vim: set ts=4 sw=4 noet:
