// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

package config

// GRPCConfig configures the gRPC server listener. {{{

type GRPCConfig struct {
	Address string `toml:"address"` // gRPC bind address.
}

// }}}
// RESTConfig configures the REST gateway listener. {{{

type RESTConfig struct {
	Address string `toml:"address"` // REST bind address.
}

// }}}
// WebsocketConfig configures the WebSocket listener. {{{

type WebsocketConfig struct {
	Address string `toml:"address"` // WebSocket bind address.
}

// }}}
// vim: set ts=4 sw=4 noet:
