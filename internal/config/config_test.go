// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config tests.
// This file exercises strict TOML decoding, overlay merge semantics, and
// validation error aggregation for the config package.
// Tests use in-memory TOML strings and explicit configs to cover edge cases.

package config

import (
	"errors"
	"strings"
	"testing"
)

// Config tests. {{{

func TestDecodeConfigRejectsUnknownKeys(t *testing.T) {
	input := `
[server]
id = "local"
unknown = "value"
`

	_, err := DecodeConfig(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected unknown key error")
	}
	if !strings.Contains(err.Error(), "unknown config keys") {
		t.Fatalf("expected unknown keys error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "server.unknown") {
		t.Fatalf("expected server.unknown in error, got: %v", err)
	}
}

func TestApplyOverlayAllowsZeroOverride(t *testing.T) {
	base := Config{
		Sync: SyncConfig{
			Enabled: true,
			Mode:    SyncModeRemote,
		},
	}
	overlay := ConfigOverlay{
		Sync: &SyncConfigOverlay{
			Enabled: boolPointer(false),
		},
	}

	result := ApplyOverlay(base, overlay)
	if result.Sync.Enabled {
		t.Fatalf("expected sync.enabled to be overridden to false")
	}
	if result.Sync.Mode != SyncModeRemote {
		t.Fatalf("expected sync.mode to remain remote, got: %s", result.Sync.Mode)
	}
}

func TestValidateConfigAggregatesErrors(t *testing.T) {
	config := DefaultConfig()
	config.Database.Driver = DriverSQLite
	config.Database.DSN = ""
	config.Auth.Enabled = true
	config.Auth.KeyAuth.Enabled = false
	config.Auth.PasswordAuth.Enabled = false
	config.Auth.Mode = AuthModeRemote
	config.Auth.Remote.Endpoint = ""
	config.Sync.Enabled = true
	config.Sync.Mode = ""

	err := ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var errs ValidationErrors
	if !errors.As(err, &errs) {
		t.Fatalf("expected ValidationErrors, got: %T", err)
	}

	if len(errs) != 4 {
		t.Fatalf("expected 4 validation errors, got %d", len(errs))
	}

	expected := map[string]bool{
		"database.dsn":         false,
		"auth.enabled":         false,
		"auth.remote.endpoint": false,
		"sync.mode":            false,
	}
	for _, fieldErr := range errs {
		if _, ok := expected[fieldErr.Path]; ok {
			expected[fieldErr.Path] = true
		}
	}
	for path, seen := range expected {
		if !seen {
			t.Fatalf("expected validation error for %s", path)
		}
	}
}

func boolPointer(value bool) *bool {
	return &value
}

// }}}

// vim: set ts=4 sw=4 noet:
