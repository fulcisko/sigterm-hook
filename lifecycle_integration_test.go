package sigterm_test

import (
	"context"
	"testing"

	sigterm "github.com/example/sigterm-hook"
)

func TestLifecycleHookFiredOnShutdown(t *testing.T) {
	var events []sigterm.LifecycleEvent

	mgr := sigterm.NewManager(
		sigterm.WithLifecycleHook(func(e sigterm.LifecycleEvent, _ sigterm.Phase) {
			events = append(events, e)
		}),
	)

	_ = mgr.Register("svc", func(_ context.Context) error { return nil })

	if err := mgr.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}

	if len(events) < 2 {
		t.Fatalf("expected at least 2 lifecycle events, got %d", len(events))
	}
	if events[0] != sigterm.LifecycleShutdownStarted {
		t.Errorf("first event should be LifecycleShutdownStarted, got %v", events[0])
	}
	if events[len(events)-1] != sigterm.LifecycleShutdownDone {
		t.Errorf("last event should be LifecycleShutdownDone, got %v", events[len(events)-1])
	}
}

func TestLifecycleHookPhaseEvents(t *testing.T) {
	var phaseStarts, phaseDones int

	mgr := sigterm.NewManager(
		sigterm.WithLifecycleHook(func(e sigterm.LifecycleEvent, _ sigterm.Phase) {
			switch e {
			case sigterm.LifecyclePhaseStarted:
				phaseStarts++
			case sigterm.LifecyclePhaseDone:
				phaseDones++
			}
		}),
	)

	_ = mgr.Register("a", func(_ context.Context) error { return nil })

	if err := mgr.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}

	if phaseStarts == 0 || phaseStarts != phaseDones {
		t.Errorf("phase start/done counts mismatch: starts=%d dones=%d", phaseStarts, phaseDones)
	}
}
