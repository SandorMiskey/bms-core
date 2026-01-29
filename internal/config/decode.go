// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

import (
	"io"

	"github.com/BurntSushi/toml"
)

// Config decoding. {{{

func DecodeConfig(reader io.Reader) (Config, error) {
	var config Config
	decoder := toml.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if _, err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// }}}
// vim: set ts=4 sw=4 noet:
