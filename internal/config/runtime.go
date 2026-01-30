// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Runtime feature configuration.
// This file defines integration, plugin, sync, and telemetry config structs.
// The types map to [integrations], [plugins], [sync], and [telemetry] sections
// to toggle runtime features and endpoints.

package config

// Clublog, LoTW and QRZ.com integrations. {{{

type ClublogConfig struct {
	Enabled bool `toml:"enabled"` // Toggle integration on or off.
}

type LoTWConfig struct {
	Enabled bool `toml:"enabled"` // Toggle integration on or off.
}

type QRZConfig struct {
	Enabled bool `toml:"enabled"` // Toggle integration on or off.
}

type IntegrationsConfig struct {
	Clublog ClublogConfig `toml:"clublog"` // Clublog integration settings.
	LoTW    LoTWConfig    `toml:"lotw"`    // Logbook of The World settings.
	QRZ     QRZConfig     `toml:"qrz"`     // QRZ.com integration settings.
}

// }}}
// PluginsConfig configures server-side plugins. {{{

type PluginsConfig struct {
	Enabled bool   `toml:"enabled"` // Toggle plugin loading.
	Path    string `toml:"path"`    // Plugin filesystem path.
}

// }}}
// SyncMode selects the sync runtime mode. {{{

type SyncMode string

const (
	SyncModeLocal  SyncMode = "local"
	SyncModeRemote SyncMode = "remote"
)

// }}}
// SyncConfig configures sync behavior. {{{

type SyncConfig struct {
	Enabled bool     `toml:"enabled"` // Toggle sync on or off.
	Mode    SyncMode `toml:"mode"`    // Sync runtime mode.
}

// }}}
// TelemetryConfig configures telemetry reporting. {{{

type TelemetryConfig struct {
	Enabled  bool   `toml:"enabled"`  // Toggle telemetry reporting.
	Endpoint string `toml:"endpoint"` // Optional telemetry endpoint.
}

// }}}

// vim: set ts=4 sw=4 noet:
