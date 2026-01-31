// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Logging context helpers.
// This file defines context helpers for request and trace identifiers used by
// logging and diagnostics. Identifiers are stored as strings for now; typed
// wrappers may be introduced later when tracing is integrated.

package logging

import "context"

// Logging context helpers. {{{

type contextKey struct {
	name string
}

var (
	requestIDKey = contextKey{name: FieldRequestID}
	traceIDKey   = contextKey{name: FieldTraceID}
)

// WithRequestID returns a context that carries the request identifier.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext extracts the request identifier if present.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(requestIDKey)
	if value == nil {
		return "", false
	}
	requestID, ok := value.(string)
	return requestID, ok
}

// WithTraceID returns a context that carries the trace identifier.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// TraceIDFromContext extracts the trace identifier if present.
func TraceIDFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(traceIDKey)
	if value == nil {
		return "", false
	}
	traceID, ok := value.(string)
	return traceID, ok
}

// }}}

// vim: set ts=4 sw=4 noet:
