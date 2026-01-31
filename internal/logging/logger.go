// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Logging initialization.
// This file defines NewLogger, which constructs a slog.Logger from logging
// config values, applies defaults when the config is unset, and attaches
// base fields (component, server_id, environment) for consistent output.

package logging

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/SandorMiskey/bms-core/internal/config"
	"github.com/SandorMiskey/bms-core/internal/errtext"
)

// Logger defaults. {{{

type DefaultFields struct {
	Component   Component
	Environment string
	ServerID    string
}

type LoggerDefaults struct {
	Fields DefaultFields
	Format config.LogFormat
	Level  config.LogLevel
}

// NewLogger builds a slog.Logger from config with fallbacks and base fields.
func NewLogger(cfg config.LoggingConfig, defaults LoggerDefaults) (*slog.Logger, error) {
	format, err := resolveLogFormat(cfg.Format, defaults.Format)
	if err != nil {
		return nil, err
	}
	level, err := resolveLogLevel(cfg.Level, defaults.Level)
	if err != nil {
		return nil, err
	}
	if defaults.Fields.Component != "" && !ValidComponent(defaults.Fields.Component) {
		return nil, fmt.Errorf("%s: %q", errtext.ErrInvalidLogComponent, defaults.Fields.Component)
	}

	logger := slog.New(newHandler(format, level))
	logger = applyDefaultFields(logger, defaults.Fields)

	return logger, nil
}

// }}}
// Logger helpers. {{{

func resolveLogFormat(format config.LogFormat, fallback config.LogFormat) (config.LogFormat, error) {
	if format == "" {
		format = fallback
	}
	if format == "" {
		return "", fmt.Errorf("%s", errtext.ErrLogFormatRequired)
	}

	switch format {
	case config.LogFormatJSON, config.LogFormatText:
		return format, nil
	default:
		return "", fmt.Errorf("%s: %q", errtext.ErrInvalidLogFormat, format)
	}
}

func resolveLogLevel(level config.LogLevel, fallback config.LogLevel) (slog.Level, error) {
	if level == "" {
		level = fallback
	}
	if level == "" {
		return 0, fmt.Errorf("%s", errtext.ErrLogLevelRequired)
	}

	switch level {
	case config.LogLevelDebug:
		return slog.LevelDebug, nil
	case config.LogLevelInfo:
		return slog.LevelInfo, nil
	case config.LogLevelWarn:
		return slog.LevelWarn, nil
	case config.LogLevelError:
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("%s: %q", errtext.ErrInvalidLogLevel, level)
	}
}

func newHandler(format config.LogFormat, level slog.Level) slog.Handler {
	options := &slog.HandlerOptions{Level: level}
	if format == config.LogFormatText {
		return slog.NewTextHandler(os.Stdout, options)
	}
	return slog.NewJSONHandler(os.Stdout, options)
}

func applyDefaultFields(logger *slog.Logger, fields DefaultFields) *slog.Logger {
	attrs := make([]any, 0, 6)
	if fields.Component != "" {
		attrs = append(attrs, FieldComponent, string(fields.Component))
	}
	if fields.ServerID != "" {
		attrs = append(attrs, FieldServerID, fields.ServerID)
	}
	if fields.Environment != "" {
		attrs = append(attrs, FieldEnvironment, fields.Environment)
	}
	if len(attrs) == 0 {
		return logger
	}
	return logger.With(attrs...)
}

// }}}

// vim: set ts=4 sw=4 noet:
