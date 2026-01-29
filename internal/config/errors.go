// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

import "strings"

// Validation error types. {{{

type FieldError struct {
	Message string `toml:"-"`	// Validation error message.
	Path    string `toml:"-"`	// Config field path.
}

func (err FieldError) Error() string {
	if err.Path == "" {
		return err.Message
	}
	if err.Message == "" {
		return err.Path
	}
	return err.Path + ": " + err.Message
}

type ValidationErrors []FieldError

func (errs ValidationErrors) Error() string {
	if len(errs) == 0 {
		return ""
	}

	var builder strings.Builder
	for index, err := range errs {
		if index > 0 {
			builder.WriteString("; ")
		}
		builder.WriteString(err.Error())
	}
	return builder.String()
}

// }}}
// vim: set ts=4 sw=4 noet:
