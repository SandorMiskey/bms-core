// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

// LogFormat defines the log output format. {{{

type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

// }}}
// LogLevel defines the minimum log severity. {{{

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// }}}
// LoggingConfig configures structured logging output. {{{

type LoggingConfig struct {
	Format LogFormat `toml:"format"` // Log format (`json` or `text`).
	Level  LogLevel  `toml:"level"`  // Minimum log level (`debug`, `info`, `warn`, `error`).
}

// }}}

// vim: set ts=4 sw=4 noet:
