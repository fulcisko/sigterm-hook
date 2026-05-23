package sigtermhook

import "sync/atomic"

// shutdownState represents the lifecycle state of the shutdown manager.
type shutdownState uint32

const (
	// stateIdle means the manager has not started shutting down.
	stateIdle shutdownState = iota
	// stateShuttingDown means shutdown has been initiated.
	stateShuttingDown
	// stateDone means shutdown has completed.
	stateDone
)

// String returns a human-readable name for the state.
func (s shutdownState) String() string {
	switch s {
	case stateIdle:
		return "idle"
	case stateShuttingDown:
		return "shutting_down"
	case stateDone:
		return "done"
	default:
		return "unknown"
	}
}

// atomicState wraps a shutdownState for concurrent access.
type atomicState struct {
	val atomic.Uint32
}

// load returns the current state.
func (a *atomicState) load() shutdownState {
	return shutdownState(a.val.Load())
}

// transition attempts to move from expected to next.
// Returns true if the transition succeeded.
func (a *atomicState) transition(expected, next shutdownState) bool {
	return a.val.CompareAndSwap(uint32(expected), uint32(next))
}

// isIdle reports whether the manager has not yet started shutting down.
func (a *atomicState) isIdle() bool {
	return a.load() == stateIdle
}

// isShuttingDown reports whether shutdown is in progress.
func (a *atomicState) isShuttingDown() bool {
	return a.load() == stateShuttingDown
}

// isDone reports whether shutdown has completed.
func (a *atomicState) isDone() bool {
	return a.load() == stateDone
}
