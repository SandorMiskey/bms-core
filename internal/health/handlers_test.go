// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Health handler tests.
// This file verifies health and readiness handlers return expected status codes
// and response bodies for ready and not-ready states.

package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Health handler tests. {{{

func TestHealthzHandler(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, healthzPath, nil)

	HealthzHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if recorder.Body.String() != "ok\n" {
		t.Fatalf("unexpected body: %q", recorder.Body.String())
	}
}

func TestReadyzHandler(t *testing.T) {
	state := NewState()
	handler := ReadyzHandler(state)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, readyzPath, nil)
	handler(recorder, request)

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", recorder.Code)
	}
	if recorder.Body.String() != "not ready\n" {
		t.Fatalf("unexpected body: %q", recorder.Body.String())
	}

	state.SetReady(true)
	recorder = httptest.NewRecorder()
	request = httptest.NewRequest(http.MethodGet, readyzPath, nil)
	handler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if recorder.Body.String() != "ready\n" {
		t.Fatalf("unexpected body: %q", recorder.Body.String())
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
