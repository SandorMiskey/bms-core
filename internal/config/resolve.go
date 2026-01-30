// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config resolution pipeline.
// This file defines ResolveConfig and ResolveConfigAndValidate, which build
// the effective runtime Config by applying defaults, file overlays, environment
// overrides, CLI overlays, and server-required overrides in a fixed order.
// The validation wrapper runs ValidateConfig after resolution, and both
// functions propagate loader or override errors without fallback.

package config

// Config resolution pipeline. {{{

// This block defines ResolveConfig and its validation wrapper, which construct
// the effective runtime configuration from defaults and override sources.

// ResolveConfig builds the effective runtime config from all override sources.
// overridePath selects the local config file path (empty uses defaults),
// cliOverlay supplies optional CLI overrides, and serverOverride supplies
// allowlisted server-required overrides applied at the end.
// It returns the resolved config, the resolved path used for loading, and
// any error from loading or applying overrides (no fallback is attempted).

func ResolveConfig(overridePath string, cliOverlay ConfigOverlay, serverOverride ConfigOverlay) (Config, string, error) {
	base := DefaultConfig()
	overlay, path, err := LoadConfigOverlayFromDefault(overridePath)
	if err != nil {
		return Config{}, path, err
	}

	base = ApplyOverlay(base, overlay)

	base, err = ApplyEnvOverrides(base)
	if err != nil {
		return Config{}, path, err
	}

	base = ApplyOverlay(base, cliOverlay)
	base = ApplyServerOverrides(base, serverOverride)

	return base, path, nil
}

// ResolveConfigAndValidate resolves the config and validates the result.
func ResolveConfigAndValidate(overridePath string, cliOverlay ConfigOverlay, serverOverride ConfigOverlay) (Config, string, error) {
	config, path, err := ResolveConfig(overridePath, cliOverlay, serverOverride)
	if err != nil {
		return config, path, err
	}
	if err := ValidateConfig(config); err != nil {
		return config, path, err
	}

	return config, path, nil
}

// }}}

// vim: set ts=4 sw=4 noet:
