// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Overlay decoding helpers.
// This file defines DecodeConfigOverlay, which parses TOML into pointer-based
// overlay structs so that "unset" values remain nil and can be merged safely
// onto a base Config without losing explicit zero-value overrides.
// Parsing is strict and rejects unknown keys using the same metadata checks
// as the runtime config decoder.

package config

import (
	"io"

	"github.com/BurntSushi/toml"
)

// Config overlay decoding. {{{

func DecodeConfigOverlay(reader io.Reader) (ConfigOverlay, error) {
	var overlay ConfigOverlay
	decoder := toml.NewDecoder(reader)
	meta, err := decoder.Decode(&overlay)
	if err != nil {
		return ConfigOverlay{}, err
	}
	if err := checkUndecodedKeys(meta); err != nil {
		return ConfigOverlay{}, err
	}

	return overlay, nil
}

// }}}

// vim: set ts=4 sw=4 noet:
