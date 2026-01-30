// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Environment override helpers.
// This file maps documented environment variables into a ConfigOverlay and
// applies them to a base Config. Empty environment values are ignored so that
// only explicitly set variables override the base configuration.

package config

import "os"

// Env overrides. {{{

const (
	envAuthMode       = "BMS_AUTH_MODE"
	envDatabaseDSN    = "BMS_DATABASE_DSN"
	envDatabaseDriver = "BMS_DATABASE_DRIVER"
	envServerID       = "BMS_SERVER_ID"
	envSyncMode       = "BMS_SYNC_MODE"
)

// ApplyEnvOverrides merges environment overrides into the base config.
func ApplyEnvOverrides(base Config) (Config, error) {
	overlay := ConfigOverlay{}

	if value := os.Getenv(envDatabaseDSN); value != "" {
		overlay.Database = ensureDatabaseOverlay(overlay.Database)
		overlay.Database.DSN = stringPointer(value)
	}
	if value := os.Getenv(envDatabaseDriver); value != "" {
		overlay.Database = ensureDatabaseOverlay(overlay.Database)
		overlay.Database.Driver = databaseDriverPointer(value)
	}
	if value := os.Getenv(envServerID); value != "" {
		overlay.Server = ensureServerOverlay(overlay.Server)
		overlay.Server.ID = stringPointer(value)
	}
	if value := os.Getenv(envAuthMode); value != "" {
		overlay.Auth = ensureAuthOverlay(overlay.Auth)
		overlay.Auth.Mode = authModePointer(value)
	}
	if value := os.Getenv(envSyncMode); value != "" {
		overlay.Sync = ensureSyncOverlay(overlay.Sync)
		overlay.Sync.Mode = syncModePointer(value)
	}

	return ApplyOverlay(base, overlay), nil
}

// }}}
// Env overlay helpers. {{{

func ensureAuthOverlay(overlay *AuthConfigOverlay) *AuthConfigOverlay {
	if overlay != nil {
		return overlay
	}
	return &AuthConfigOverlay{}
}

func ensureDatabaseOverlay(overlay *DatabaseConfigOverlay) *DatabaseConfigOverlay {
	if overlay != nil {
		return overlay
	}
	return &DatabaseConfigOverlay{}
}

func ensureServerOverlay(overlay *ServerConfigOverlay) *ServerConfigOverlay {
	if overlay != nil {
		return overlay
	}
	return &ServerConfigOverlay{}
}

func ensureSyncOverlay(overlay *SyncConfigOverlay) *SyncConfigOverlay {
	if overlay != nil {
		return overlay
	}
	return &SyncConfigOverlay{}
}

func stringPointer(value string) *string {
	return &value
}

func authModePointer(value string) *AuthMode {
	mode := AuthMode(value)
	return &mode
}

func databaseDriverPointer(value string) *DatabaseDriver {
	driver := DatabaseDriver(value)
	return &driver
}

func syncModePointer(value string) *SyncMode {
	mode := SyncMode(value)
	return &mode
}

// }}}

// vim: set ts=4 sw=4 noet:
