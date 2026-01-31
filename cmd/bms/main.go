// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// CLI entry point.
// This file defines the bms main function, which resolves configuration,
// initializes structured logging with CLI defaults, emits startup diagnostics
// (redacted), and exits on configuration errors.

package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"

	"github.com/SandorMiskey/bms-core/internal/config"
	"github.com/SandorMiskey/bms-core/internal/errtext"
	"github.com/SandorMiskey/bms-core/internal/logging"
)

// Main entry point. {{{

func main() {
	configPath := flag.String("config", "", "path to config.toml")
	flag.Parse()

	configResult, path, warnings, err := config.ResolveConfigDiagnostics(*configPath, config.ConfigOverlay{}, config.ConfigOverlay{})
	logger, format := initLogger(configResult, logging.ComponentCLI)

	if err != nil {
		var validationErrors config.ValidationErrors
		if errors.As(err, &validationErrors) {
			logging.LogConfigDiagnostics(logger, format, configResult, path, warnings)
			logger.Error(errtext.ErrConfigValidationFailed, "error", err)
		} else {
			logger.Error(errtext.ErrConfigResolutionFailed, "error", err)
		}
		os.Exit(1)
	}

	logging.LogConfigDiagnostics(logger, format, configResult, path, warnings)
}

func initLogger(cfg config.Config, component logging.Component) (*slog.Logger, config.LogFormat) {
	defaults := logging.LoggerDefaults{
		Fields: logging.DefaultFields{
			Component:   component,
			Environment: string(cfg.Server.Environment),
			ServerID:    cfg.Server.ID,
		},
		Format: config.LogFormatText,
		Level:  config.LogLevelInfo,
	}

	logger, err := logging.NewLogger(cfg.Logging, defaults)
	if err != nil {
		fallback := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
		fallback.Error(errtext.ErrLoggerInitFailed, "error", err)
		logger = fallback
	}

	format := cfg.Logging.Format
	if format == "" || err != nil {
		format = defaults.Format
	}

	return logger, format
}

// }}}

// vim: set ts=4 sw=4 noet:
