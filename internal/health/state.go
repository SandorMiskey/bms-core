// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

// Health state tracking.
// This file defines a readiness state used by health handlers to report
// whether the service is ready to accept traffic.

package health

import "sync"

// Health state. {{{

type State struct {
	mu    sync.RWMutex
	ready bool
}

func NewState() *State {
	return &State{}
}

func (state *State) SetReady(ready bool) {
	state.mu.Lock()
	defer state.mu.Unlock()
	state.ready = ready
}

func (state *State) IsReady() bool {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.ready
}

// }}}

// vim: set ts=4 sw=4 noet:
