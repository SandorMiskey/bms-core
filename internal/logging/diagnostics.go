// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config diagnostics logging.
// This file logs startup diagnostics for configuration resolution, including
// the resolved config (redacted), the config path, and any non-fatal warnings.
// Warnings are emitted as a list for JSON logs and a single string for text logs
// to keep both formats readable without losing detail.

package logging

import (
	"log/slog"

	"github.com/SandorMiskey/bms-core/internal/config"
)

// Config diagnostics logging. {{{
// This block defines LogConfigDiagnostics, which emits config_loaded and
// config_warnings events for startup diagnostics.

const (
	eventConfigLoaded   = "config_loaded"
	eventConfigWarnings = "config_warnings"
)

// LogConfigDiagnostics emits startup diagnostics for a resolved config.
func LogConfigDiagnostics(logger *slog.Logger, format config.LogFormat, cfg config.Config, path string, warnings config.WarningList) {
	redacted := config.RedactConfig(cfg)

	logger.Info(
		"config loaded",
		FieldComponent, string(ComponentConfig),
		FieldEvent, eventConfigLoaded,
		FieldConfigPath, path,
		FieldRedacted, true,
		FieldWarningsCount, len(warnings),
		"config", redacted,
	)

	if len(warnings) == 0 {
		return
	}

	logger.Warn(
		"config warnings",
		FieldComponent, string(ComponentConfig),
		FieldEvent, eventConfigWarnings,
		FieldWarningsCount, len(warnings),
		"warnings", formatWarnings(format, warnings),
	)
}

// }}}
// Warning formatting helpers. {{{

type warningEntry struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func formatWarnings(format config.LogFormat, warnings config.WarningList) any {
	if format == config.LogFormatJSON {
		return warningEntries(warnings)
	}
	return warnings.String()
}

func warningEntries(warnings config.WarningList) []warningEntry {
	entries := make([]warningEntry, 0, len(warnings))
	for _, warning := range warnings {
		entries = append(entries, warningEntry{
			Path:    warning.Path,
			Message: warning.Message,
		})
	}
	return entries
}

// }}}

// vim: set ts=4 sw=4 noet:
