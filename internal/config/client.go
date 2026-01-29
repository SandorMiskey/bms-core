// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

// ClientConfig configures client-side settings. {{{

type ClientConfig struct {
	Auth    ClientAuthConfig    `toml:"auth"`    // Authentication settings.
	Keymap  ClientKeymapConfig  `toml:"keymap"`  // Keymap selection.
	Offline ClientOfflineConfig `toml:"offline"` // Offline mode settings.
	Plugins ClientPluginsConfig `toml:"plugins"` // Client plugin settings.
	Server  ClientServerConfig  `toml:"server"`  // Server endpoints.
	Theme   ClientThemeConfig   `toml:"theme"`   // Theme selection.
}

// }}}
// Client supporting configuration structs. {{{

// ClientServerConfig configures client endpoints.
type ClientServerConfig struct {
	Address string `toml:"address"` // gRPC endpoint.
	REST    string `toml:"rest"`    // REST endpoint.
}

// ClientAuthConfig configures client-side auth persistence.
type ClientAuthConfig struct {
	RefreshBeforeExpiry float64 `toml:"refresh_before_expiry"` // Token refresh threshold.
	StoreToken          bool    `toml:"store_token"`           // Persist auth token locally.
	Token               string  `toml:"token"`                 // Auth token value.
}

// ClientThemeConfig selects the theme.
type ClientThemeConfig struct {
	Name string `toml:"name"` // Theme name.
}

// ClientKeymapConfig selects the keymap.
type ClientKeymapConfig struct {
	Name string `toml:"name"` // Keymap name.
}

// ClientPluginsConfig configures client plugins.
type ClientPluginsConfig struct {
	Enabled bool   `toml:"enabled"` // Toggle client plugins.
	Path    string `toml:"path"`    // Plugin directory.
}

// ClientOfflineConfig configures offline behavior.
type ClientOfflineConfig struct {
	Enabled bool `toml:"enabled"` // Toggle offline mode.
}

// }}}
// vim: set ts=4 sw=4 noet:
