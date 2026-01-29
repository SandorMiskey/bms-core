// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

// ClublogConfig configures Clublog integration. {{{

type ClublogConfig struct {
	Enabled bool `toml:"enabled"`	// Toggle integration on or off.
}

// }}}
// LoTWConfig configures LoTW integration. {{{

type LoTWConfig struct {
	Enabled bool `toml:"enabled"`	// Toggle integration on or off.
}

// }}}
// QRZConfig configures QRZ.com integration. {{{

type QRZConfig struct {
	Enabled bool `toml:"enabled"`	// Toggle integration on or off.
}

// }}}
// IntegrationsConfig groups external integration settings. {{{

type IntegrationsConfig struct {
	Clublog ClublogConfig `toml:"clublog"`	// Clublog integration settings.
	LoTW    LoTWConfig    `toml:"lotw"`		// Logbook of The World settings.
	QRZ     QRZConfig     `toml:"qrz"`		// QRZ.com integration settings.
}

// }}}
// PluginsConfig configures server-side plugins. {{{

type PluginsConfig struct {
	Enabled bool   `toml:"enabled"`	// Toggle plugin loading.
	Path    string `toml:"path"`		// Plugin filesystem path.
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
	Enabled bool     `toml:"enabled"`	// Toggle sync on or off.
	Mode    SyncMode `toml:"mode"`		// Sync runtime mode.
}

// }}}
// TelemetryConfig configures telemetry reporting. {{{

type TelemetryConfig struct {
	Enabled  bool   `toml:"enabled"`	// Toggle telemetry reporting.
	Endpoint string `toml:"endpoint"`	// Optional telemetry endpoint.
}

// }}}
// vim: set ts=4 sw=4 noet:
