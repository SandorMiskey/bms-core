// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

import "os"

// Config loading. {{{

// LoadConfig reads and decodes the config file at the given path.
func LoadConfig(path string) (Config, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, ErrConfigNotFound
		}
		return Config{}, err
	}
	if info.IsDir() {
		return Config{}, ErrConfigPathIsDir
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	return DecodeConfig(file)
}

// LoadConfigFromDefault resolves and loads a config path using overrides.
func LoadConfigFromDefault(override string) (Config, string, error) {
	path, err := ResolveConfigPath(override)
	if err != nil {
		return Config{}, "", err
	}

	config, err := LoadConfig(path)
	if err != nil {
		return Config{}, path, err
	}

	return config, path, nil
}

// }}}
// vim: set ts=4 sw=4 noet:
