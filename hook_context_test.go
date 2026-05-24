package sigterm

import (
	"context"
	"testing"
)

func TestShutdownContextFrom_Missing(t *testing.T) {
	sc := ShutdownContextFrom(context.Background())
	if sc != nil {
		t.Fatal("expected nil for plain context")
	}
}

func TestShutdownContextFrom_Present(t *testing.T) {
	ctx, sc := newShutdownContext(context.Background(), "SIGTERM")
	got := ShutdownContextFrom(ctx)
	if got == nil {
		t.Fatal("expected non-nil ShutdownContext")
	}
	if got != sc {
		t.Fatal("expected same pointer")
	}
}

func TestShutdownContextReason(t *testing.T) {
	_, sc := newShutdownContext(context.Background(), "manual")
	if sc.Reason() != "manual" {
		t.Fatalf("expected reason 'manual', got %q", sc.Reason())
	}
}

func TestShutdownContextCompletedHandlers_Empty(t *testing.T) {
	_, sc := newShutdownContext(context.Background(), "test")
	if len(sc.CompletedHandlers()) != 0 {
		t.Fatal("expected empty completed handlers")
	}
}

func TestShutdownContextMarkHandlerDone(t *testing.T) {
	_, sc := newShutdownContext(context.Background(), "test")
	sc.markHandlerDone("db")
	sc.markHandlerDone("cache")

	handlers := sc.CompletedHandlers()
	if len(handlers) != 2 {
		t.Fatalf("expected 2 handlers, got %d", len(handlers))
	}
	if handlers[0] != "db" || handlers[1] != "cache" {
		t.Fatalf("unexpected handler order: %v", handlers)
	}
}

func TestShutdownContextSnapshotIsolation(t *testing.T) {
	_, sc := newShutdownContext(context.Background(), "test")
	sc.markHandlerDone("a")

	snap := sc.CompletedHandlers()
	sc.markHandlerDone("b")

	if len(snap) != 1 {
		t.Fatalf("snapshot should not reflect later mutations, got len=%d", len(snap))
	}
}
