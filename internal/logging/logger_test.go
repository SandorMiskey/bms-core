// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Logger validation tests.
// This file verifies NewLogger rejects invalid formats, levels, and components
// and returns clear error messages for operators.

package logging

import (
	"strings"
	"testing"

	"github.com/SandorMiskey/bms-core/internal/config"
	"github.com/SandorMiskey/bms-core/internal/errtext"
)

// Logger validation tests. {{{

func TestNewLoggerInvalidFormat(t *testing.T) {
	cfg := config.LoggingConfig{Format: config.LogFormat("invalid")}
	defaults := LoggerDefaults{Format: config.LogFormatJSON, Level: config.LogLevelInfo}

	_, err := NewLogger(cfg, defaults)
	if err == nil {
		t.Fatal("expected error for invalid log format")
	}
	if !strings.Contains(err.Error(), errtext.ErrInvalidLogFormat) {
		t.Fatalf("expected invalid log format error, got: %v", err)
	}
}

func TestNewLoggerInvalidLevel(t *testing.T) {
	cfg := config.LoggingConfig{Level: config.LogLevel("invalid")}
	defaults := LoggerDefaults{Format: config.LogFormatJSON, Level: config.LogLevelInfo}

	_, err := NewLogger(cfg, defaults)
	if err == nil {
		t.Fatal("expected error for invalid log level")
	}
	if !strings.Contains(err.Error(), errtext.ErrInvalidLogLevel) {
		t.Fatalf("expected invalid log level error, got: %v", err)
	}
}

func TestNewLoggerInvalidComponent(t *testing.T) {
	defaults := LoggerDefaults{
		Fields: DefaultFields{Component: Component("invalid")},
		Format: config.LogFormatJSON,
		Level:  config.LogLevelInfo,
	}

	_, err := NewLogger(config.LoggingConfig{}, defaults)
	if err == nil {
		t.Fatal("expected error for invalid log component")
	}
	if !strings.Contains(err.Error(), errtext.ErrInvalidLogComponent) {
		t.Fatalf("expected invalid log component error, got: %v", err)
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
