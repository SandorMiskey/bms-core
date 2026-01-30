// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config root structure.
// This file defines the top-level Config struct used by the config pipeline.
// It groups server and client sections and maps them to TOML sections for
// decoding, merging, and validation.

package config

// Config holds server and client configuration sections. {{{

type Config struct {
	Auth         AuthConfig         `toml:"auth"`         // Authentication settings.
	Database     DatabaseConfig     `toml:"database"`     // Database connectivity settings.
	GRPC         GRPCConfig         `toml:"grpc"`         // gRPC listener configuration.
	Integrations IntegrationsConfig `toml:"integrations"` // External integration settings.
	Logging      LoggingConfig      `toml:"logging"`      // Logging output configuration.
	Plugins      PluginsConfig      `toml:"plugins"`      // Plugin runtime settings.
	REST         RESTConfig         `toml:"rest"`         // REST listener configuration.
	Server       ServerConfig       `toml:"server"`       // Server instance settings.
	Sync         SyncConfig         `toml:"sync"`         // Sync settings.
	Telemetry    TelemetryConfig    `toml:"telemetry"`    // Telemetry reporting settings.
	Websocket    WebsocketConfig    `toml:"websocket"`    // WebSocket listener configuration.
	Client       ClientConfig       `toml:"client"`       // Client-side settings.
}

// }}}

// vim: set ts=4 sw=4 noet:
