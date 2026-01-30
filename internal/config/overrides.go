// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Server-required overrides.
// This file applies a constrained subset of overlay fields that the server
// is allowed to enforce. Only allowlisted auth and sync fields are honored,
// so clients cannot override server-required runtime behavior.

package config

// Server-required overrides. {{{

// ApplyServerOverrides applies allowlisted server overrides to a base Config.
func ApplyServerOverrides(base Config, override ConfigOverlay) Config {
	sanitized := ConfigOverlay{}

	if override.Auth != nil {
		authOverlay := AuthConfigOverlay{}
		if override.Auth.Enabled != nil {
			authOverlay.Enabled = override.Auth.Enabled
		}
		if override.Auth.Mode != nil {
			authOverlay.Mode = override.Auth.Mode
		}
		if authOverlay.Enabled != nil || authOverlay.Mode != nil {
			sanitized.Auth = &authOverlay
		}
	}

	if override.Sync != nil {
		syncOverlay := SyncConfigOverlay{}
		if override.Sync.Enabled != nil {
			syncOverlay.Enabled = override.Sync.Enabled
		}
		if override.Sync.Mode != nil {
			syncOverlay.Mode = override.Sync.Mode
		}
		if syncOverlay.Enabled != nil || syncOverlay.Mode != nil {
			sanitized.Sync = &syncOverlay
		}
	}

	return ApplyOverlay(base, sanitized)
}

// }}}

// vim: set ts=4 sw=4 noet:
