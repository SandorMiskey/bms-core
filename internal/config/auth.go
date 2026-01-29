// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

// Auth enums and constants. {{{

type AuthMode string

const (
	AuthModeHybrid AuthMode = "hybrid"
	AuthModeLocal  AuthMode = "local"
	AuthModeRemote AuthMode = "remote"
)

type AuthTokenStorage string

const (
	AuthTokenStorageConfig   AuthTokenStorage = "config"
	AuthTokenStorageFile     AuthTokenStorage = "file"
	AuthTokenStorageKeychain AuthTokenStorage = "keychain"
)

// }}}
// AuthConfig configures authentication and session handling. {{{

type AuthConfig struct {
	DevicePairing       AuthDevicePairingConfig `toml:"device_pairing"`		// Device registration settings.
	Enabled             bool                    `toml:"enabled"`			// Toggle auth on or off.
	KeyAuth             AuthKeyAuthConfig       `toml:"key_auth"`		// Key-based login settings.
	KeyStorage          AuthKeyStorageConfig    `toml:"key_storage"`		// Local key storage options.
	LocalTrust          AuthLocalTrustConfig    `toml:"local_trust"`		// Local-only passwordless access.
	Mode                AuthMode                `toml:"mode"`			// Runtime auth mode.
	PasswordAuth        AuthPasswordAuthConfig  `toml:"password_auth"`		// Password login settings.
	Recovery            AuthRecoveryConfig      `toml:"recovery"`		// Recovery code settings.
	RefreshBeforeExpiry float64                 `toml:"refresh_before_expiry"`	// Token rotation threshold.
	Remote              AuthRemoteConfig        `toml:"remote"`			// Delegated auth endpoint settings.
	TokenStorage        AuthTokenStorage        `toml:"token_storage"`		// Token persistence target.
	TokenTTL            string                  `toml:"token_ttl"`			// Token lifetime (duration string).
}

// }}}
// Auth supporting configuration structs. {{{

// Device pairing settings.
type AuthDevicePairingConfig struct {
	Enabled      bool `toml:"enabled"`		// Toggle pairing on or off.
	QR           bool `toml:"qr"`		// Enable QR-based pairing hints.
	RequireLocal bool `toml:"require_local"`	// Require local-only pairing.
}

// Key-based login settings.
type AuthKeyAuthConfig struct {
	Enabled bool `toml:"enabled"`	// Toggle key auth on or off.
}

// Local private key storage settings.
type AuthKeyStorageConfig struct {
	AllowUnencrypted bool `toml:"allow_unencrypted"`	// Allow unencrypted storage (explicit opt-in).
	Encrypted        bool `toml:"encrypted"`		// Enable encryption of private keys.
}

// Local-only passwordless access settings.
type AuthLocalTrustConfig struct {
	Enabled bool `toml:"enabled"`	// Toggle local trust on or off.
}

// Password login settings.
type AuthPasswordAuthConfig struct {
	Enabled bool `toml:"enabled"`	// Toggle password auth on or off.
}

// Recovery code settings.
type AuthRecoveryConfig struct {
	Codes   int  `toml:"codes"`		// Number of recovery codes to issue.
	Enabled bool `toml:"enabled"`	// Toggle recovery codes on or off.
}

// Delegated auth endpoint settings.
type AuthRemoteConfig struct {
	Endpoint string `toml:"endpoint"`	// Delegated auth endpoint.
}

// }}}
// vim: set ts=4 sw=4 noet:
