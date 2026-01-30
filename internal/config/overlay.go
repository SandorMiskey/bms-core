// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Overlay types for config merging. {{{
// This file defines pointer-based overlay structs that mirror the config schema.
// Each field is optional (nil means "no override"), which preserves the
// difference between an unset value and an explicit zero value.
// The overlay structs are decoded from TOML and later applied by merge helpers
// to override a base Config without changing the runtime config types.
// No behavior is implemented here; this file only declares the overlay data model.
// }}}

package config
// Config overlay root. {{{

type ConfigOverlay struct {
	Auth         *AuthConfigOverlay         `toml:"auth"`         // Authentication overrides.
	Database     *DatabaseConfigOverlay     `toml:"database"`     // Database overrides.
	GRPC         *GRPCConfigOverlay         `toml:"grpc"`         // gRPC listener overrides.
	Integrations *IntegrationsConfigOverlay `toml:"integrations"` // Integration overrides.
	Logging      *LoggingConfigOverlay      `toml:"logging"`      // Logging overrides.
	Plugins      *PluginsConfigOverlay      `toml:"plugins"`      // Plugin overrides.
	REST         *RESTConfigOverlay         `toml:"rest"`         // REST listener overrides.
	Server       *ServerConfigOverlay       `toml:"server"`       // Server overrides.
	Sync         *SyncConfigOverlay         `toml:"sync"`         // Sync overrides.
	Telemetry    *TelemetryConfigOverlay    `toml:"telemetry"`    // Telemetry overrides.
	Websocket    *WebsocketConfigOverlay    `toml:"websocket"`    // WebSocket overrides.
	Client       *ClientConfigOverlay       `toml:"client"`       // Client overrides.
}

// }}}
// Server overlay structs. {{{

type ServerConfigOverlay struct {
	Environment *Environment `toml:"environment"` // Runtime mode override.
	ID          *string      `toml:"id"`          // Instance identifier override.
}

type DatabaseConfigOverlay struct {
	DSN        *string         `toml:"dsn"`        // Connection string override.
	Driver     *DatabaseDriver `toml:"driver"`     // Database engine override.
	Migrations *string         `toml:"migrations"` // Migrations directory override.
}

type LoggingConfigOverlay struct {
	Format *LogFormat `toml:"format"` // Log format override.
	Level  *LogLevel  `toml:"level"`  // Minimum log level override.
}

type GRPCConfigOverlay struct {
	Address *string `toml:"address"` // gRPC bind address override.
}

type RESTConfigOverlay struct {
	Address *string `toml:"address"` // REST bind address override.
}

type WebsocketConfigOverlay struct {
	Address *string `toml:"address"` // WebSocket bind address override.
}

type IntegrationsConfigOverlay struct {
	Clublog *ClublogConfigOverlay `toml:"clublog"` // Clublog integration overrides.
	LoTW    *LoTWConfigOverlay    `toml:"lotw"`    // Logbook of The World overrides.
	QRZ     *QRZConfigOverlay     `toml:"qrz"`     // QRZ.com overrides.
}

type ClublogConfigOverlay struct {
	Enabled *bool `toml:"enabled"` // Toggle integration override.
}

type LoTWConfigOverlay struct {
	Enabled *bool `toml:"enabled"` // Toggle integration override.
}

type QRZConfigOverlay struct {
	Enabled *bool `toml:"enabled"` // Toggle integration override.
}

type PluginsConfigOverlay struct {
	Enabled *bool   `toml:"enabled"` // Toggle plugin loading override.
	Path    *string `toml:"path"`    // Plugin filesystem path override.
}

type SyncConfigOverlay struct {
	Enabled *bool     `toml:"enabled"` // Toggle sync override.
	Mode    *SyncMode `toml:"mode"`    // Sync runtime mode override.
}

type TelemetryConfigOverlay struct {
	Enabled  *bool   `toml:"enabled"`  // Toggle telemetry override.
	Endpoint *string `toml:"endpoint"` // Telemetry endpoint override.
}

// }}}
// Auth overlay structs. {{{

type AuthConfigOverlay struct {
	DevicePairing       *AuthDevicePairingConfigOverlay `toml:"device_pairing"`        // Device registration overrides.
	Enabled             *bool                           `toml:"enabled"`               // Toggle auth override.
	KeyAuth             *AuthKeyAuthConfigOverlay       `toml:"key_auth"`              // Key-based login overrides.
	KeyStorage          *AuthKeyStorageConfigOverlay    `toml:"key_storage"`           // Local key storage overrides.
	LocalTrust          *AuthLocalTrustConfigOverlay    `toml:"local_trust"`           // Local-only access overrides.
	Mode                *AuthMode                       `toml:"mode"`                  // Runtime auth mode override.
	PasswordAuth        *AuthPasswordAuthConfigOverlay  `toml:"password_auth"`         // Password login overrides.
	Recovery            *AuthRecoveryConfigOverlay      `toml:"recovery"`              // Recovery code overrides.
	RefreshBeforeExpiry *float64                        `toml:"refresh_before_expiry"` // Token rotation threshold override.
	Remote              *AuthRemoteConfigOverlay        `toml:"remote"`                // Delegated auth endpoint overrides.
	TokenStorage        *AuthTokenStorage               `toml:"token_storage"`         // Token persistence override.
	TokenTTL            *string                         `toml:"token_ttl"`             // Token lifetime override.
}

type AuthDevicePairingConfigOverlay struct {
	Enabled      *bool `toml:"enabled"`       // Toggle pairing override.
	QR           *bool `toml:"qr"`            // QR-based pairing override.
	RequireLocal *bool `toml:"require_local"` // Require local-only pairing override.
}

type AuthKeyAuthConfigOverlay struct {
	Enabled *bool `toml:"enabled"` // Toggle key auth override.
}

type AuthKeyStorageConfigOverlay struct {
	AllowUnencrypted *bool `toml:"allow_unencrypted"` // Allow unencrypted storage override.
	Encrypted        *bool `toml:"encrypted"`         // Enable encryption override.
}

type AuthLocalTrustConfigOverlay struct {
	Enabled *bool `toml:"enabled"` // Toggle local trust override.
}

type AuthPasswordAuthConfigOverlay struct {
	Enabled *bool `toml:"enabled"` // Toggle password auth override.
}

type AuthRecoveryConfigOverlay struct {
	Codes   *int  `toml:"codes"`   // Recovery code count override.
	Enabled *bool `toml:"enabled"` // Toggle recovery override.
}

type AuthRemoteConfigOverlay struct {
	Endpoint *string `toml:"endpoint"` // Delegated auth endpoint override.
}

// }}}
// Client overlay structs. {{{

type ClientConfigOverlay struct {
	Auth    *ClientAuthConfigOverlay    `toml:"auth"`    // Authentication overrides.
	Keymap  *ClientKeymapConfigOverlay  `toml:"keymap"`  // Keymap overrides.
	Offline *ClientOfflineConfigOverlay `toml:"offline"` // Offline mode overrides.
	Plugins *ClientPluginsConfigOverlay `toml:"plugins"` // Plugin overrides.
	Server  *ClientServerConfigOverlay  `toml:"server"`  // Server endpoint overrides.
	Theme   *ClientThemeConfigOverlay   `toml:"theme"`   // Theme overrides.
}

type ClientServerConfigOverlay struct {
	Address *string `toml:"address"` // gRPC endpoint override.
	REST    *string `toml:"rest"`    // REST endpoint override.
}

type ClientAuthConfigOverlay struct {
	RefreshBeforeExpiry *float64 `toml:"refresh_before_expiry"` // Token refresh threshold override.
	StoreToken          *bool    `toml:"store_token"`           // Persist auth token override.
	Token               *string  `toml:"token"`                 // Auth token override.
}

type ClientThemeConfigOverlay struct {
	Name *string `toml:"name"` // Theme name override.
}

type ClientKeymapConfigOverlay struct {
	Name *string `toml:"name"` // Keymap name override.
}

type ClientPluginsConfigOverlay struct {
	Enabled *bool   `toml:"enabled"` // Toggle client plugins override.
	Path    *string `toml:"path"`    // Plugin directory override.
}

type ClientOfflineConfigOverlay struct {
	Enabled *bool `toml:"enabled"` // Toggle offline mode override.
}

// }}}

// vim: set ts=4 sw=4 noet:
