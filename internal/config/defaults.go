// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Default configuration values.
// This file defines DefaultConfig, which builds the baseline Config used by the
// merge pipeline before any file, env, or CLI overrides are applied.
// Defaults are explicit, deterministic, and limited to documented values; all
// other fields remain zero-valued so overlays can control final behavior.

package config

// Default values. {{{

const (
	defaultAuthTokenTTL           = "168h"
	defaultRefreshBeforeExpiry    = 0.8
	defaultServerAuthTokenStorage = AuthTokenStorageKeychain
)

// DefaultConfig returns the baseline configuration defaults.
func DefaultConfig() Config {
	return Config{
		Auth: AuthConfig{
			RefreshBeforeExpiry: defaultRefreshBeforeExpiry,
			TokenStorage:        defaultServerAuthTokenStorage,
			TokenTTL:            defaultAuthTokenTTL,
		},
		Client: ClientConfig{
			Auth: ClientAuthConfig{
				RefreshBeforeExpiry: defaultRefreshBeforeExpiry,
			},
		},
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
