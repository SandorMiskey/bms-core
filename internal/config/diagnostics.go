// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config diagnostics pipeline.
// This file defines ResolveConfigDiagnostics, which resolves configuration,
// collects non-fatal warnings, and validates the result for startup logging.
// It returns the resolved config, resolved path, and warning list so callers
// can emit diagnostics without re-running merge steps.

package config

// Config diagnostics pipeline. {{{

// ResolveConfigDiagnostics resolves configuration and returns warnings with validation.
func ResolveConfigDiagnostics(overridePath string, cliOverlay ConfigOverlay, serverOverride ConfigOverlay) (Config, string, WarningList, error) {
	config, path, err := ResolveConfig(overridePath, cliOverlay, serverOverride)
	if err != nil {
		return config, path, nil, err
	}

	warnings := CollectConfigWarnings(config)
	if err := ValidateConfig(config); err != nil {
		return config, path, warnings, err
	}

	return config, path, warnings, nil
}

// }}}

// vim: set ts=4 sw=4 noet:
