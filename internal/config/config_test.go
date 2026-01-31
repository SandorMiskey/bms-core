// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config tests.
// This file exercises strict TOML decoding, overlay merge semantics, redaction,
// loader error handling, and validation aggregation for the config package.
// Tests use in-memory TOML strings and explicit configs to cover edge cases.

package config

import (
	"errors"
	"os"
	"path/filepath"
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

func TestRedactConfig(t *testing.T) {
	config := Config{
		Auth: AuthConfig{
			Remote: AuthRemoteConfig{Endpoint: "https://user:pass@example.com"},
		},
		Client: ClientConfig{
			Auth: ClientAuthConfig{Token: "token"},
		},
		Database: DatabaseConfig{DSN: "file:secret.db"},
		Logging:  LoggingConfig{Level: LogLevelInfo},
	}

	redacted := RedactConfig(config)
	if redacted.Database.DSN != "[redacted]" {
		t.Fatalf("expected database.dsn to be redacted, got: %s", redacted.Database.DSN)
	}
	if redacted.Auth.Remote.Endpoint != "[redacted]" {
		t.Fatalf("expected auth.remote.endpoint to be redacted, got: %s", redacted.Auth.Remote.Endpoint)
	}
	if redacted.Client.Auth.Token != "[redacted]" {
		t.Fatalf("expected client.auth.token to be redacted, got: %s", redacted.Client.Auth.Token)
	}
	if redacted.Logging.Level != LogLevelInfo {
		t.Fatalf("expected logging.level to remain unchanged, got: %s", redacted.Logging.Level)
	}
}

func TestLoadConfigErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.toml")
	if _, err := LoadConfig(missing); !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got: %v", err)
	}
	if _, err := LoadConfigOverlay(missing); !errors.Is(err, ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound for overlay, got: %v", err)
	}

	dir := t.TempDir()
	if _, err := LoadConfig(dir); !errors.Is(err, ErrConfigPathIsDir) {
		t.Fatalf("expected ErrConfigPathIsDir, got: %v", err)
	}
	if _, err := LoadConfigOverlay(dir); !errors.Is(err, ErrConfigPathIsDir) {
		t.Fatalf("expected ErrConfigPathIsDir for overlay, got: %v", err)
	}
}

func TestResolveConfigPrecedence(t *testing.T) {
	file, err := os.CreateTemp("", "bms-config-*.toml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		_ = os.Remove(file.Name())
	}()

	input := `
[server]
id = "file-id"
environment = "local"

[auth]
mode = "local"
enabled = true
key_auth = { enabled = true }
`
	if _, err := file.WriteString(input); err != nil {
		_ = file.Close()
		t.Fatalf("failed to write temp config: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("failed to close temp config: %v", err)
	}

	t.Setenv("BMS_SERVER_ID", "env-id")
	t.Setenv("BMS_DATABASE_DSN", "env.db")
	t.Setenv("BMS_AUTH_MODE", "local")

	cliOverlay := ConfigOverlay{
		Auth:     &AuthConfigOverlay{Mode: authModePointer(string(AuthModeHybrid))},
		Database: &DatabaseConfigOverlay{DSN: stringPointer("cli.db")},
		Server:   &ServerConfigOverlay{ID: stringPointer("cli-id")},
	}
	serverOverride := ConfigOverlay{
		Auth:   &AuthConfigOverlay{Mode: authModePointer(string(AuthModeRemote))},
		Server: &ServerConfigOverlay{ID: stringPointer("server-id")},
	}

	result, path, err := ResolveConfig(file.Name(), cliOverlay, serverOverride)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if path != file.Name() {
		t.Fatalf("expected resolved path to match file, got: %s", path)
	}
	if result.Server.Environment != EnvLocal {
		t.Fatalf("expected server.environment to be local, got: %s", result.Server.Environment)
	}
	if result.Server.ID != "cli-id" {
		t.Fatalf("expected server.id to be cli-id, got: %s", result.Server.ID)
	}
	if result.Database.DSN != "cli.db" {
		t.Fatalf("expected database.dsn to be cli.db, got: %s", result.Database.DSN)
	}
	if result.Auth.Mode != AuthModeRemote {
		t.Fatalf("expected auth.mode to be remote, got: %s", result.Auth.Mode)
	}
}

func boolPointer(value bool) *bool {
	return &value
}

// }}}

// vim: set ts=4 sw=4 noet:
