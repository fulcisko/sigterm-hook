package sigtermhook_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	sigtermhook "github.com/your-org/sigterm-hook"
)

func TestNotifierReceivesStartedAndCompleted(t *testing.T) {
	var mu sync.Mutex
	var events []sigtermhook.ShutdownEvent

	m := sigtermhook.NewManager(
		sigtermhook.WithNotifierFunc(func(n sigtermhook.ShutdownNotification) {
			mu.Lock()
			events = append(events, n.Event)
			mu.Unlock()
		}),
	)

	m.Register("svc", func(ctx context.Context) error { return nil })

	if err := m.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	want := []sigtermhook.ShutdownEvent{
		sigtermhook.EventShutdownStarted,
		sigtermhook.EventHandlerStarted,
		sigtermhook.EventHandlerCompleted,
		sigtermhook.EventShutdownCompleted,
	}
	if len(events) != len(want) {
		t.Fatalf("got %d events, want %d: %v", len(events), len(want), events)
	}
	for i, e := range events {
		if e != want[i] {
			t.Errorf("event[%d] = %v, want %v", i, e, want[i])
		}
	}
}

func TestNotifierReceivesHandlerFailed(t *testing.T) {
	var got []sigtermhook.ShutdownEvent
	sentinel := errors.New("boom")

	m := sigtermhook.NewManager(
		sigtermhook.WithNotifierFunc(func(n sigtermhook.ShutdownNotification) {
			got = append(got, n.Event)
		}),
	)

	m.Register("bad", func(ctx context.Context) error { return sentinel })
	_ = m.Shutdown(context.Background())

	failed := false
	for _, e := range got {
		if e == sigtermhook.EventHandlerFailed {
			failed = true
		}
	}
	if !failed {
		t.Error("expected EventHandlerFailed to be emitted")
	}
}

func TestMultipleNotifiers(t *testing.T) {
	var countA, countB int

	m := sigtermhook.NewManager(
		sigtermhook.WithNotifierFunc(func(_ sigtermhook.ShutdownNotification) { countA++ }),
		sigtermhook.WithNotifierFunc(func(_ sigtermhook.ShutdownNotification) { countB++ }),
	)

	m.Register("svc", func(ctx context.Context) error { return nil })
	_ = m.Shutdown(context.Background())

	if countA == 0 || countB == 0 {
		t.Errorf("expected both notifiers to be called, got countA=%d countB=%d", countA, countB)
	}
	if countA != countB {
		t.Errorf("notifiers received different counts: %d vs %d", countA, countB)
	}
}

func TestNotifierHandlerName(t *testing.T) {
	var names []string

	m := sigtermhook.NewManager(
		sigtermhook.WithNotifierFunc(func(n sigtermhook.ShutdownNotification) {
			if n.Handler != "" {
				names = append(names, n.Handler)
			}
		}),
	)

	m.Register("my-service", func(ctx context.Context) error { return nil })
	_ = m.Shutdown(context.Background())

	for _, name := range names {
		if name == "my-service" {
			return
		}
	}
	t.Error("expected handler name 'my-service' to appear in notifications")
}
