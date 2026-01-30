// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config warnings.
// This file defines warning types and a collection helper that reports
// non-fatal configuration issues which should be surfaced to operators.
// Warnings are returned as data so callers can log or display them.

package config

import "strings"

// Config warnings. {{{

type FieldWarning struct {
	Message string `toml:"-"`	// Warning message.
	Path    string `toml:"-"`	// Config field path.
}

func (warn FieldWarning) String() string {
	if warn.Path == "" {
		return warn.Message
	}
	if warn.Message == "" {
		return warn.Path
	}
	return warn.Path + ": " + warn.Message
}

type WarningList []FieldWarning

func (warnings WarningList) String() string {
	if len(warnings) == 0 {
		return ""
	}

	var builder strings.Builder
	for index, warning := range warnings {
		if index > 0 {
			builder.WriteString("; ")
		}
		builder.WriteString(warning.String())
	}
	return builder.String()
}

// CollectConfigWarnings returns non-fatal config warnings for operator review.
func CollectConfigWarnings(config Config) WarningList {
	var warnings WarningList

	if config.Auth.KeyStorage.AllowUnencrypted {
		warnings = append(warnings, FieldWarning{
			Path:    "auth.key_storage.allow_unencrypted",
			Message: "allows unencrypted key storage; review before enabling",
		})
	}

	return warnings
}

// }}}

// vim: set ts=4 sw=4 noet:
