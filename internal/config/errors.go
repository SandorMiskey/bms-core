// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config error types.
// This file defines errors used by the config pipeline and validation helpers.
// It provides file/loader errors plus aggregated field validation errors for
// reporting multiple config issues at once.

package config

import (
	"errors"
	"strings"
)

// Config errors and validation types. {{{

// ErrConfigNotFound indicates the config file does not exist.
var ErrConfigNotFound = errors.New("config file not found")

// ErrConfigPathIsDir indicates the config path points to a directory.
var ErrConfigPathIsDir = errors.New("config path is a directory")

type FieldError struct {
	Message string `toml:"-"` // Validation error message.
	Path    string `toml:"-"` // Config field path.
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
