// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config tests.
// This file exercises strict TOML decoding, overlay merge semantics, and
// validation error aggregation for the config package.
// Tests use in-memory TOML strings and explicit configs to cover edge cases.

package config

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/SandorMiskey/bms-core/internal/errtext"
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
		t.Fatal("expected invalid config keys error")
	}
	if !strings.Contains(err.Error(), errtext.ErrInvalidConfigKeys) {
		t.Fatalf("expected invalid config keys error, got: %v", err)
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

func TestResolveConfigDiagnosticsReturnsWarnings(t *testing.T) {
	file, err := os.CreateTemp("", "bms-config-*.toml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		_ = os.Remove(file.Name())
	}()

	input := `
[database]
driver = "sqlite"
dsn = "file:bms.db"

[auth.key_storage]
allow_unencrypted = true
`
	if _, err := file.WriteString(input); err != nil {
		_ = file.Close()
		t.Fatalf("failed to write temp config: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("failed to close temp config: %v", err)
	}

	config, path, warnings, err := ResolveConfigDiagnostics(file.Name(), ConfigOverlay{}, ConfigOverlay{})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if path == "" {
		t.Fatal("expected resolved path")
	}
	if config.Database.DSN == "" {
		t.Fatal("expected database.dsn to be set")
	}
	if len(warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(warnings))
	}
	if warnings[0].Path != "auth.key_storage.allow_unencrypted" {
		t.Fatalf("unexpected warning path: %s", warnings[0].Path)
	}
}

func boolPointer(value bool) *bool {
	return &value
}

// }}}

// vim: set ts=4 sw=4 noet:
