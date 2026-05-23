package sigtermhook

import (
	"testing"
)

func TestNotifierBusEmitCallsAll(t *testing.T) {
	var calls int
	fn := NotifierFunc(func(_ ShutdownNotification) { calls++ })
	bus := newNotifierBus([]Notifier{fn, fn, fn})

	bus.emit(ShutdownNotification{Event: EventShutdownStarted})

	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestNotifierBusHasNotifiers(t *testing.T) {
	empty := newNotifierBus(nil)
	if empty.hasNotifiers() {
		t.Error("expected hasNotifiers to return false for empty bus")
	}

	bus := newNotifierBus([]Notifier{NotifierFunc(func(_ ShutdownNotification) {})})
	if !bus.hasNotifiers() {
		t.Error("expected hasNotifiers to return true")
	}
}

func TestNotifierBusEmptyEmitIsNoop(t *testing.T) {
	bus := newNotifierBus(nil)
	// should not panic
	bus.emit(ShutdownNotification{Event: EventShutdownCompleted})
}

func TestNotifierFuncImplementsNotifier(t *testing.T) {
	called := false
	var n Notifier = NotifierFunc(func(_ ShutdownNotification) { called = true })
	n.Notify(ShutdownNotification{})
	if !called {
		t.Error("NotifierFunc.Notify did not call the underlying function")
	}
}

func TestShutdownEventConstants(t *testing.T) {
	events := []ShutdownEvent{
		EventShutdownStarted,
		EventHandlerStarted,
		EventHandlerCompleted,
		EventHandlerFailed,
		EventShutdownCompleted,
	}
	seen := make(map[ShutdownEvent]bool)
	for _, e := range events {
		if seen[e] {
			t.Errorf("duplicate ShutdownEvent value: %d", e)
		}
		seen[e] = true
	}
}
