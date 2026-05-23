package sigtermhook

import (
	"testing"
)

func TestShutdownStateString(t *testing.T) {
	tests := []struct {
		state    shutdownState
		expected string
	}{
		{stateIdle, "idle"},
		{stateShuttingDown, "shutting_down"},
		{stateDone, "done"},
		{shutdownState(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.state.String(); got != tt.expected {
			t.Errorf("state %d: got %q, want %q", tt.state, got, tt.expected)
		}
	}
}

func TestAtomicStateInitiallyIdle(t *testing.T) {
	var s atomicState
	if !s.isIdle() {
		t.Fatal("expected initial state to be idle")
	}
	if s.isShuttingDown() {
		t.Fatal("expected isShuttingDown to be false initially")
	}
	if s.isDone() {
		t.Fatal("expected isDone to be false initially")
	}
}

func TestAtomicStateTransitionIdleToShuttingDown(t *testing.T) {
	var s atomicState
	if !s.transition(stateIdle, stateShuttingDown) {
		t.Fatal("expected transition to succeed")
	}
	if !s.isShuttingDown() {
		t.Fatal("expected state to be shutting_down")
	}
}

func TestAtomicStateTransitionShuttingDownToDone(t *testing.T) {
	var s atomicState
	s.transition(stateIdle, stateShuttingDown)
	if !s.transition(stateShuttingDown, stateDone) {
		t.Fatal("expected transition to done to succeed")
	}
	if !s.isDone() {
		t.Fatal("expected state to be done")
	}
}

func TestAtomicStateTransitionFailsOnWrongExpected(t *testing.T) {
	var s atomicState
	// Try to go from shutting_down to done when still idle — should fail.
	if s.transition(stateShuttingDown, stateDone) {
		t.Fatal("expected transition to fail when current state does not match expected")
	}
	if !s.isIdle() {
		t.Fatal("state should remain idle after failed transition")
	}
}

func TestAtomicStateTransitionIdempotentFail(t *testing.T) {
	var s atomicState
	s.transition(stateIdle, stateShuttingDown)
	// Second attempt to transition from idle should fail.
	if s.transition(stateIdle, stateShuttingDown) {
		t.Fatal("expected second transition from idle to fail")
	}
	if !s.isShuttingDown() {
		t.Fatal("state should still be shutting_down")
	}
}
