package sigterm

// LifecycleEvent represents a named moment in the shutdown lifecycle.
type LifecycleEvent int

const (
	// LifecycleShutdownStarted fires once before any handler is invoked.
	LifecycleShutdownStarted LifecycleEvent = iota
	// LifecyclePhaseStarted fires before each phase begins.
	LifecyclePhaseStarted
	// LifecyclePhaseDone fires after each phase completes (with or without errors).
	LifecyclePhaseDone
	// LifecycleShutdownDone fires after all phases have finished.
	LifecycleShutdownDone
)

// LifecycleHook is a callback invoked at a specific lifecycle event.
// phase is only meaningful for LifecyclePhaseStarted and LifecyclePhaseDone.
type LifecycleHook func(event LifecycleEvent, phase Phase)

// lifecycleBus holds zero or more LifecycleHook callbacks and dispatches events.
type lifecycleBus struct {
	hooks []LifecycleHook
}

func newLifecycleBus(hooks ...LifecycleHook) *lifecycleBus {
	return &lifecycleBus{hooks: hooks}
}

// emit calls every registered hook with the given event and phase.
func (b *lifecycleBus) emit(event LifecycleEvent, phase Phase) {
	for _, h := range b.hooks {
		h(event, phase)
	}
}

// has returns true when at least one hook is registered.
func (b *lifecycleBus) has() bool {
	return len(b.hooks) > 0
}
