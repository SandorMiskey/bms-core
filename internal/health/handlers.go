// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Health HTTP handlers.
// This file defines HTTP handlers and a mux for liveness and readiness checks.
// The handlers return plain text responses with appropriate status codes.

package health

import "net/http"

// Health handlers. {{{

const (
	healthzPath = "/healthz"
	readyzPath  = "/readyz"
)

func NewMux(state *State) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(healthzPath, HealthzHandler)
	mux.HandleFunc(readyzPath, ReadyzHandler(state))
	return mux
}

func HealthzHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte("ok\n"))
}

func ReadyzHandler(state *State) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if state != nil && state.IsReady() {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("ready\n"))
			return
		}
		writer.WriteHeader(http.StatusServiceUnavailable)
		_, _ = writer.Write([]byte("not ready\n"))
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
