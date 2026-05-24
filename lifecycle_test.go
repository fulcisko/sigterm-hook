package sigterm

import (
	"sync"
	"testing"
)

func TestLifecycleBusEmitCallsAll(t *testing.T) {
	var mu sync.Mutex
	var got []LifecycleEvent

	h := func(e LifecycleEvent, _ Phase) {
		mu.Lock()
		got = append(got, e)
		mu.Unlock()
	}

	bus := newLifecycleBus(h, h)
	bus.emit(LifecycleShutdownStarted, PhaseNormal)

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(got))
	}
}

func TestLifecycleBusHas(t *testing.T) {
	if newLifecycleBus().has() {
		t.Fatal("empty bus should return false for has()")
	}
	if !newLifecycleBus(func(LifecycleEvent, Phase) {}).has() {
		t.Fatal("non-empty bus should return true for has()")
	}
}

func TestLifecycleBusEmitPassesPhase(t *testing.T) {
	var gotPhase Phase
	bus := newLifecycleBus(func(_ LifecycleEvent, p Phase) {
		gotPhase = p
	})
	bus.emit(LifecyclePhaseStarted, PhaseCritical)
	if gotPhase != PhaseCritical {
		t.Fatalf("expected PhaseCritical, got %v", gotPhase)
	}
}

func TestLifecycleEventConstants(t *testing.T) {
	events := []LifecycleEvent{
		LifecycleShutdownStarted,
		LifecyclePhaseStarted,
		LifecyclePhaseDone,
		LifecycleShutdownDone,
	}
	seen := make(map[LifecycleEvent]bool)
	for _, e := range events {
		if seen[e] {
			t.Fatalf("duplicate lifecycle event value: %d", e)
		}
		seen[e] = true
	}
}

func TestLifecycleEmptyBusIsNoop(t *testing.T) {
	bus := newLifecycleBus()
	// Must not panic.
	bus.emit(LifecycleShutdownDone, PhaseNormal)
}
