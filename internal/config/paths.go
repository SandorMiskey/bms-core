// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

import (
	"os"
	"path/filepath"
	"strings"
)

// Config path discovery. {{{

const (
	configDirName  = "bms"
	configEnvVar   = "BMS_CONFIG"
	configFileName = "config.toml"
)

// DefaultConfigPath returns the default config.toml location.
func DefaultConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, configDirName, configFileName), nil
}

// ResolveConfigPath picks the config path using overrides or defaults.
func ResolveConfigPath(override string) (string, error) {
	if override != "" {
		return expandUserPath(override)
	}

	envOverride := os.Getenv(configEnvVar)
	if envOverride != "" {
		return expandUserPath(envOverride)
	}

	return DefaultConfigPath()
}

func expandUserPath(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	if path == "~" {
		return os.UserHomeDir()
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}

// }}}
// vim: set ts=4 sw=4 noet:
