// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config redaction helpers.
// This file defines RedactConfig, which returns a sanitized copy of a resolved
// Config suitable for summary logging. Sensitive values (DSNs and tokens) are
// replaced with a constant placeholder so the full config shape can be logged
// without exposing secrets.

package config

// Config redaction helpers. {{{

const redactedValue = "[redacted]"

// RedactConfig returns a sanitized copy of config with sensitive fields removed.
func RedactConfig(config Config) Config {
	redacted := config
	redacted.Database.DSN = redactValue(redacted.Database.DSN)
	redacted.Client.Auth.Token = redactValue(redacted.Client.Auth.Token)

	return redacted
}

func redactValue(value string) string {
	if value == "" {
		return ""
	}
	return redactedValue
}

// }}}

// vim: set ts=4 sw=4 noet:
