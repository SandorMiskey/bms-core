// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Canonical error text.
// This file defines shared error message fragments used across packages so
// errors remain consistent, searchable, and easy to maintain.
// Callers format these constants with context details and wrapped errors.

package errtext

// Error text constants. {{{

const (
	ErrConfigResolutionFailed = "config resolution failed"
	ErrConfigValidationFailed = "config validation failed"
	ErrInvalidConfigKeys      = "invalid config keys"
	ErrInvalidLogComponent    = "invalid log component"
	ErrInvalidLogFormat       = "invalid log format"
	ErrInvalidLogLevel        = "invalid log level"
	ErrLoggerInitFailed       = "logger init failed"
	ErrLogFormatRequired      = "log format is required"
	ErrLogLevelRequired       = "log level is required"
	ErrOpenConfig             = "open config"
	ErrOpenConfigOverlay      = "open config overlay"
	ErrStatConfig             = "stat config"
	ErrStatConfigOverlay      = "stat config overlay"
)

// }}}

// vim: set ts=4 sw=4 noet:
