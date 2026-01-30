// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Config loading.
// This file provides file-based loaders for both runtime configs and overlay
// configs, validating that the path exists and points to a file before decode.
// The helper functions return the resolved path so callers can log which file
// was used without duplicating discovery logic.

package config

import "os"

// LoadConfig reads and decodes the config file at the given path. {{{

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

// }}}
// LoadConfigFromDefault resolves and loads a config path using overrides. {{{

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
// LoadConfigOverlay reads and decodes the config overlay at the given path. {{{

func LoadConfigOverlay(path string) (ConfigOverlay, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ConfigOverlay{}, ErrConfigNotFound
		}
		return ConfigOverlay{}, err
	}
	if info.IsDir() {
		return ConfigOverlay{}, ErrConfigPathIsDir
	}

	file, err := os.Open(path)
	if err != nil {
		return ConfigOverlay{}, err
	}
	defer file.Close()

	return DecodeConfigOverlay(file)
}

// }}}
// LoadConfigOverlayFromDefault resolves and loads an overlay config file. {{{

func LoadConfigOverlayFromDefault(override string) (ConfigOverlay, string, error) {
	path, err := ResolveConfigPath(override)
	if err != nil {
		return ConfigOverlay{}, "", err
	}

	overlay, err := LoadConfigOverlay(path)
	if err != nil {
		return ConfigOverlay{}, path, err
	}

	return overlay, path, nil
}

// }}}

// vim: set ts=4 sw=4 noet:
