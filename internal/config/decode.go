// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config decoding helpers.
// This file parses TOML into the runtime Config struct and enforces strict
// decoding by rejecting any undecoded keys returned by the TOML metadata.
// A small helper formats invalid keys so callers get a clear error message.

package config

import (
	"fmt"
	"io"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/SandorMiskey/bms-core/internal/errtext"
)

// Config decoding. {{{

func DecodeConfig(reader io.Reader) (Config, error) {
	var config Config
	decoder := toml.NewDecoder(reader)
	meta, err := decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	if err := checkUndecodedKeys(meta); err != nil {
		return Config{}, err
	}

	return config, nil
}

func checkUndecodedKeys(meta toml.MetaData) error {
	keys := meta.Undecoded()
	if len(keys) == 0 {
		return nil
	}

	return fmt.Errorf("%s: %s", errtext.ErrInvalidConfigKeys, formatUndecodedKeys(keys))
}

func formatUndecodedKeys(keys []toml.Key) string {
	var builder strings.Builder
	for index, key := range keys {
		if index > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(key.String())
	}
	return builder.String()
}

// }}}

// vim: set ts=4 sw=4 noet:
