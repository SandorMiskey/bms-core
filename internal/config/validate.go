// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config validation.
// This file defines ValidateConfig, which checks resolved configs for required
// relationships and enum constraints and returns aggregated field errors.
// Validation is deterministic, runs after the merge pipeline, and does not
// touch external systems.

package config

import "time"

// Config validation. {{{
// This block defines ValidateConfig and section-specific validators that
// append field errors without short-circuiting on the first failure.

// ValidateConfig validates a resolved config and returns aggregated field errors.
func ValidateConfig(config Config) error {
	var errs ValidationErrors

	validateDatabaseConfig(config.Database, &errs)
	validateAuthConfig(config.Auth, config.Server, &errs)
	validateSyncConfig(config.Sync, &errs)
	validateAuthDurations(config.Auth, &errs)

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateDatabaseConfig(database DatabaseConfig, errs *ValidationErrors) {
	if database.Driver != DriverSQLite && database.Driver != DriverPostgres {
		appendFieldError(errs, "database.driver", "must be sqlite or postgres")
		return
	}
	if database.DSN == "" {
		appendFieldError(errs, "database.dsn", "is required when database.driver is set")
	}
}

func validateAuthConfig(auth AuthConfig, server ServerConfig, errs *ValidationErrors) {
	if auth.Enabled && !auth.KeyAuth.Enabled && !auth.PasswordAuth.Enabled {
		appendFieldError(errs, "auth.enabled", "requires auth.key_auth.enabled or auth.password_auth.enabled")
	}
	if auth.Mode != "" && auth.Mode != AuthModeLocal && auth.Mode != AuthModeRemote && auth.Mode != AuthModeHybrid {
		appendFieldError(errs, "auth.mode", "must be local, remote, or hybrid")
	}
	if auth.Mode == AuthModeRemote && auth.Remote.Endpoint == "" {
		appendFieldError(errs, "auth.remote.endpoint", "is required when auth.mode is remote")
	}
	if auth.LocalTrust.Enabled {
		if server.Environment != EnvLocal || auth.Mode != AuthModeLocal {
			appendFieldError(errs, "auth.local_trust.enabled", "requires server.environment=local and auth.mode=local")
		}
	}
}

func validateSyncConfig(sync SyncConfig, errs *ValidationErrors) {
	if sync.Enabled && sync.Mode == "" {
		appendFieldError(errs, "sync.mode", "is required when sync.enabled is true")
	}
	if sync.Mode != "" && sync.Mode != SyncModeLocal && sync.Mode != SyncModeRemote {
		appendFieldError(errs, "sync.mode", "must be local or remote")
	}
}

func validateAuthDurations(auth AuthConfig, errs *ValidationErrors) {
	if auth.TokenTTL == "" {
		appendFieldError(errs, "auth.token_ttl", "must be a duration string")
	} else if _, err := time.ParseDuration(auth.TokenTTL); err != nil {
		appendFieldError(errs, "auth.token_ttl", "must be a valid duration")
	}
	if auth.RefreshBeforeExpiry < 0 || auth.RefreshBeforeExpiry > 1 {
		appendFieldError(errs, "auth.refresh_before_expiry", "must be between 0 and 1")
	}
}

func appendFieldError(errs *ValidationErrors, path string, message string) {
	*errs = append(*errs, FieldError{Path: path, Message: message})
}

// }}}

// vim: set ts=4 sw=4 noet:
