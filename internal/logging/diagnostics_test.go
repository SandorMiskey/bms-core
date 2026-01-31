// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Logging diagnostics tests.
// This file validates the warning formatting behavior for JSON and text logs
// to ensure warnings are emitted as structured lists or readable strings.

package logging

import (
	"testing"

	"github.com/SandorMiskey/bms-core/internal/config"
)

// Logging diagnostics tests. {{{

func TestFormatWarningsJSON(t *testing.T) {
	warnings := config.WarningList{
		{Path: "first.path", Message: "first message"},
		{Path: "second.path", Message: "second message"},
	}

	value := formatWarnings(config.LogFormatJSON, warnings)
	entries, ok := value.([]warningEntry)
	if !ok {
		t.Fatalf("expected warningEntry slice, got %T", value)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Path != "first.path" || entries[0].Message != "first message" {
		t.Fatalf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Path != "second.path" || entries[1].Message != "second message" {
		t.Fatalf("unexpected second entry: %+v", entries[1])
	}
}

func TestFormatWarningsText(t *testing.T) {
	warnings := config.WarningList{
		{Path: "first.path", Message: "first message"},
		{Path: "second.path", Message: "second message"},
	}

	value := formatWarnings(config.LogFormatText, warnings)
	text, ok := value.(string)
	if !ok {
		t.Fatalf("expected string, got %T", value)
	}
	if text != warnings.String() {
		t.Fatalf("unexpected warnings string: %s", text)
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
