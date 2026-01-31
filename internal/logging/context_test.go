// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Logging context tests.
// This file verifies request and trace identifiers are stored in and retrieved
// from context correctly, and that missing IDs return a false indicator.

package logging

import (
	"context"
	"testing"
)

// Logging context tests. {{{

func TestRequestIDContextRoundTrip(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-123")
	value, ok := RequestIDFromContext(ctx)
	if !ok {
		t.Fatal("expected request ID to be present")
	}
	if value != "req-123" {
		t.Fatalf("unexpected request ID: %s", value)
	}
}

func TestTraceIDContextRoundTrip(t *testing.T) {
	ctx := WithTraceID(context.Background(), "trace-456")
	value, ok := TraceIDFromContext(ctx)
	if !ok {
		t.Fatal("expected trace ID to be present")
	}
	if value != "trace-456" {
		t.Fatalf("unexpected trace ID: %s", value)
	}
}

func TestMissingIDs(t *testing.T) {
	requestID, requestOk := RequestIDFromContext(context.Background())
	if requestOk {
		t.Fatal("expected request ID to be missing")
	}
	if requestID != "" {
		t.Fatalf("expected empty request ID, got: %s", requestID)
	}

	traceID, traceOk := TraceIDFromContext(context.Background())
	if traceOk {
		t.Fatal("expected trace ID to be missing")
	}
	if traceID != "" {
		t.Fatalf("expected empty trace ID, got: %s", traceID)
	}
}

// }}}

// vim: set ts=4 sw=4 noet:
