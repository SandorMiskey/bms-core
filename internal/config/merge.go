// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Overlay merge helpers.
// This file applies pointer-based overlay structs to the runtime Config.
// Each helper copies only non-nil overlay fields so explicit zero values are
// preserved, while missing fields leave the base configuration untouched.
// The merge functions are explicit and section-scoped to keep override rules
// readable, auditable, and easy to extend.

package config

// Overlay merge helpers. {{{

// This block defines ApplyOverlay, which merges top-level overlay sections into a base
// Config and delegates to section-specific helpers when the overlay section is non-nil.
// All merges treat non-nil overlay fields as overrides and keep base values otherwise.

// ApplyOverlay merges an overlay into a base Config.
func ApplyOverlay(base Config, overlay ConfigOverlay) Config {
	if overlay.Auth != nil {
		base.Auth = mergeAuthConfig(base.Auth, *overlay.Auth)
	}
	if overlay.Database != nil {
		base.Database = mergeDatabaseConfig(base.Database, *overlay.Database)
	}
	if overlay.GRPC != nil {
		base.GRPC = mergeGRPCConfig(base.GRPC, *overlay.GRPC)
	}
	if overlay.Integrations != nil {
		base.Integrations = mergeIntegrationsConfig(base.Integrations, *overlay.Integrations)
	}
	if overlay.Logging != nil {
		base.Logging = mergeLoggingConfig(base.Logging, *overlay.Logging)
	}
	if overlay.Plugins != nil {
		base.Plugins = mergePluginsConfig(base.Plugins, *overlay.Plugins)
	}
	if overlay.REST != nil {
		base.REST = mergeRESTConfig(base.REST, *overlay.REST)
	}
	if overlay.Server != nil {
		base.Server = mergeServerConfig(base.Server, *overlay.Server)
	}
	if overlay.Sync != nil {
		base.Sync = mergeSyncConfig(base.Sync, *overlay.Sync)
	}
	if overlay.Telemetry != nil {
		base.Telemetry = mergeTelemetryConfig(base.Telemetry, *overlay.Telemetry)
	}
	if overlay.Websocket != nil {
		base.Websocket = mergeWebsocketConfig(base.Websocket, *overlay.Websocket)
	}
	if overlay.Client != nil {
		base.Client = mergeClientConfig(base.Client, *overlay.Client)
	}

	return base
}

// }}}
// Server merge helpers. {{{

// This block merges server-side sections (server, database, logging, transport,
// integrations, plugins, sync, telemetry). Each helper accepts a base section
// plus its overlay and returns the updated section, applying only non-nil fields.

func mergeServerConfig(base ServerConfig, overlay ServerConfigOverlay) ServerConfig {
	if overlay.Environment != nil {
		base.Environment = *overlay.Environment
	}
	if overlay.ID != nil {
		base.ID = *overlay.ID
	}

	return base
}

func mergeDatabaseConfig(base DatabaseConfig, overlay DatabaseConfigOverlay) DatabaseConfig {
	if overlay.DSN != nil {
		base.DSN = *overlay.DSN
	}
	if overlay.Driver != nil {
		base.Driver = *overlay.Driver
	}
	if overlay.Migrations != nil {
		base.Migrations = *overlay.Migrations
	}

	return base
}

func mergeLoggingConfig(base LoggingConfig, overlay LoggingConfigOverlay) LoggingConfig {
	if overlay.Format != nil {
		base.Format = *overlay.Format
	}
	if overlay.Level != nil {
		base.Level = *overlay.Level
	}

	return base
}

func mergeGRPCConfig(base GRPCConfig, overlay GRPCConfigOverlay) GRPCConfig {
	if overlay.Address != nil {
		base.Address = *overlay.Address
	}

	return base
}

func mergeRESTConfig(base RESTConfig, overlay RESTConfigOverlay) RESTConfig {
	if overlay.Address != nil {
		base.Address = *overlay.Address
	}

	return base
}

func mergeWebsocketConfig(base WebsocketConfig, overlay WebsocketConfigOverlay) WebsocketConfig {
	if overlay.Address != nil {
		base.Address = *overlay.Address
	}

	return base
}

func mergeIntegrationsConfig(base IntegrationsConfig, overlay IntegrationsConfigOverlay) IntegrationsConfig {
	if overlay.Clublog != nil {
		base.Clublog = mergeClublogConfig(base.Clublog, *overlay.Clublog)
	}
	if overlay.LoTW != nil {
		base.LoTW = mergeLoTWConfig(base.LoTW, *overlay.LoTW)
	}
	if overlay.QRZ != nil {
		base.QRZ = mergeQRZConfig(base.QRZ, *overlay.QRZ)
	}

	return base
}

func mergeClublogConfig(base ClublogConfig, overlay ClublogConfigOverlay) ClublogConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}

	return base
}

func mergeLoTWConfig(base LoTWConfig, overlay LoTWConfigOverlay) LoTWConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}

	return base
}

func mergeQRZConfig(base QRZConfig, overlay QRZConfigOverlay) QRZConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}

	return base
}

func mergePluginsConfig(base PluginsConfig, overlay PluginsConfigOverlay) PluginsConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}
	if overlay.Path != nil {
		base.Path = *overlay.Path
	}

	return base
}

func mergeSyncConfig(base SyncConfig, overlay SyncConfigOverlay) SyncConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}
	if overlay.Mode != nil {
		base.Mode = *overlay.Mode
	}

	return base
}

func mergeTelemetryConfig(base TelemetryConfig, overlay TelemetryConfigOverlay) TelemetryConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}
	if overlay.Endpoint != nil {
		base.Endpoint = *overlay.Endpoint
	}

	return base
}

// }}}
// Auth merge helpers. {{{

// This block merges auth sections. Each helper accepts a base auth section and
// its overlay, returning an updated struct that applies only non-nil overrides.

func mergeAuthConfig(base AuthConfig, overlay AuthConfigOverlay) AuthConfig {
	if overlay.DevicePairing != nil {
		base.DevicePairing = mergeAuthDevicePairingConfig(base.DevicePairing, *overlay.DevicePairing)
	}
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}
	if overlay.KeyAuth != nil {
		base.KeyAuth = mergeAuthKeyAuthConfig(base.KeyAuth, *overlay.KeyAuth)
	}
	if overlay.KeyStorage != nil {
		base.KeyStorage = mergeAuthKeyStorageConfig(base.KeyStorage, *overlay.KeyStorage)
	}
	if overlay.LocalTrust != nil {
		base.LocalTrust = mergeAuthLocalTrustConfig(base.LocalTrust, *overlay.LocalTrust)
	}
	if overlay.Mode != nil {
		base.Mode = *overlay.Mode
	}
	if overlay.PasswordAuth != nil {
		base.PasswordAuth = mergeAuthPasswordAuthConfig(base.PasswordAuth, *overlay.PasswordAuth)
	}
	if overlay.Recovery != nil {
		base.Recovery = mergeAuthRecoveryConfig(base.Recovery, *overlay.Recovery)
	}
	if overlay.RefreshBeforeExpiry != nil {
		base.RefreshBeforeExpiry = *overlay.RefreshBeforeExpiry
	}
	if overlay.Remote != nil {
		base.Remote = mergeAuthRemoteConfig(base.Remote, *overlay.Remote)
	}
	if overlay.TokenStorage != nil {
		base.TokenStorage = *overlay.TokenStorage
	}
	if overlay.TokenTTL != nil {
		base.TokenTTL = *overlay.TokenTTL
	}

	return base
}

func mergeAuthDevicePairingConfig(base AuthDevicePairingConfig, overlay AuthDevicePairingConfigOverlay) AuthDevicePairingConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}
	if overlay.QR != nil {
		base.QR = *overlay.QR
	}
	if overlay.RequireLocal != nil {
		base.RequireLocal = *overlay.RequireLocal
	}

	return base
}

func mergeAuthKeyAuthConfig(base AuthKeyAuthConfig, overlay AuthKeyAuthConfigOverlay) AuthKeyAuthConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}

	return base
}

func mergeAuthKeyStorageConfig(base AuthKeyStorageConfig, overlay AuthKeyStorageConfigOverlay) AuthKeyStorageConfig {
	if overlay.AllowUnencrypted != nil {
		base.AllowUnencrypted = *overlay.AllowUnencrypted
	}
	if overlay.Encrypted != nil {
		base.Encrypted = *overlay.Encrypted
	}

	return base
}

func mergeAuthLocalTrustConfig(base AuthLocalTrustConfig, overlay AuthLocalTrustConfigOverlay) AuthLocalTrustConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}

	return base
}

func mergeAuthPasswordAuthConfig(base AuthPasswordAuthConfig, overlay AuthPasswordAuthConfigOverlay) AuthPasswordAuthConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}

	return base
}

func mergeAuthRecoveryConfig(base AuthRecoveryConfig, overlay AuthRecoveryConfigOverlay) AuthRecoveryConfig {
	if overlay.Codes != nil {
		base.Codes = *overlay.Codes
	}
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}

	return base
}

func mergeAuthRemoteConfig(base AuthRemoteConfig, overlay AuthRemoteConfigOverlay) AuthRemoteConfig {
	if overlay.Endpoint != nil {
		base.Endpoint = *overlay.Endpoint
	}

	return base
}

// }}}
// Client merge helpers. {{{

// This block merges client sections. Each helper accepts a base client section
// and its overlay, returning the updated struct with non-nil overrides applied.

func mergeClientConfig(base ClientConfig, overlay ClientConfigOverlay) ClientConfig {
	if overlay.Auth != nil {
		base.Auth = mergeClientAuthConfig(base.Auth, *overlay.Auth)
	}
	if overlay.Keymap != nil {
		base.Keymap = mergeClientKeymapConfig(base.Keymap, *overlay.Keymap)
	}
	if overlay.Offline != nil {
		base.Offline = mergeClientOfflineConfig(base.Offline, *overlay.Offline)
	}
	if overlay.Plugins != nil {
		base.Plugins = mergeClientPluginsConfig(base.Plugins, *overlay.Plugins)
	}
	if overlay.Server != nil {
		base.Server = mergeClientServerConfig(base.Server, *overlay.Server)
	}
	if overlay.Theme != nil {
		base.Theme = mergeClientThemeConfig(base.Theme, *overlay.Theme)
	}

	return base
}

func mergeClientServerConfig(base ClientServerConfig, overlay ClientServerConfigOverlay) ClientServerConfig {
	if overlay.Address != nil {
		base.Address = *overlay.Address
	}
	if overlay.REST != nil {
		base.REST = *overlay.REST
	}

	return base
}

func mergeClientAuthConfig(base ClientAuthConfig, overlay ClientAuthConfigOverlay) ClientAuthConfig {
	if overlay.RefreshBeforeExpiry != nil {
		base.RefreshBeforeExpiry = *overlay.RefreshBeforeExpiry
	}
	if overlay.StoreToken != nil {
		base.StoreToken = *overlay.StoreToken
	}
	if overlay.Token != nil {
		base.Token = *overlay.Token
	}

	return base
}

func mergeClientThemeConfig(base ClientThemeConfig, overlay ClientThemeConfigOverlay) ClientThemeConfig {
	if overlay.Name != nil {
		base.Name = *overlay.Name
	}

	return base
}

func mergeClientKeymapConfig(base ClientKeymapConfig, overlay ClientKeymapConfigOverlay) ClientKeymapConfig {
	if overlay.Name != nil {
		base.Name = *overlay.Name
	}

	return base
}

func mergeClientPluginsConfig(base ClientPluginsConfig, overlay ClientPluginsConfigOverlay) ClientPluginsConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}
	if overlay.Path != nil {
		base.Path = *overlay.Path
	}

	return base
}

func mergeClientOfflineConfig(base ClientOfflineConfig, overlay ClientOfflineConfigOverlay) ClientOfflineConfig {
	if overlay.Enabled != nil {
		base.Enabled = *overlay.Enabled
	}

	return base
}

// }}}

// vim: set ts=4 sw=4 noet:
