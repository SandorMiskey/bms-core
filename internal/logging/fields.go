// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Logging field definitions.
// This file defines canonical field keys and component identifiers used in
// structured logging and diagnostics so log output is consistent across
// packages and environments.

package logging

// Logging fields. {{{

const (
	FieldComponent     = "component"
	FieldEvent         = "event"
	FieldRequestID     = "request_id"
	FieldTraceID       = "trace_id"
	FieldServerID      = "server_id"
	FieldEnvironment   = "environment"
	FieldConfigPath    = "config_path"
	FieldWarningsCount = "warnings_count"
	FieldRedacted      = "redacted"
)

type Component string

const (
	ComponentAuth         Component = "auth"
	ComponentCLI          Component = "cli"
	ComponentConfig       Component = "config"
	ComponentDatabase     Component = "database"
	ComponentGRPC         Component = "grpc"
	ComponentIntegrations Component = "integrations"
	ComponentPlugins      Component = "plugins"
	ComponentREST         Component = "rest"
	ComponentServer       Component = "server"
	ComponentSync         Component = "sync"
	ComponentTelemetry    Component = "telemetry"
	ComponentWebsocket    Component = "websocket"
)

func ValidComponent(component Component) bool {
	switch component {
	case ComponentAuth,
		ComponentCLI,
		ComponentConfig,
		ComponentDatabase,
		ComponentGRPC,
		ComponentIntegrations,
		ComponentPlugins,
		ComponentREST,
		ComponentServer,
		ComponentSync,
		ComponentTelemetry,
		ComponentWebsocket:
		return true
	default:
		return false
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
