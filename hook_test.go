package sigterm

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRegisterAndShutdown(t *testing.T) {
	m := NewManager()
	var called int32

	_ = m.Register("svc", func(ctx context.Context) error {
		atomic.AddInt32(&called, 1)
		return nil
	})

	if err := m.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("handler was not called")
	}
}

func TestDependencyOrder(t *testing.T) {
	m := NewManager()
	var seq []string
	var mu sync.Mutex // use stdlib sync directly in test

	append_ := func(s string) func(context.Context) error {
		return func(_ context.Context) error {
			mu.Lock()
			seq = append(seq, s)
			mu.Unlock()
			return nil
		}
	}

	_ = m.Register("db", append_("db"))
	_ = m.Register("cache", append_("cache"))
	_ = m.Register("api", append_("api"), "db", "cache")

	if err := m.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if seq[len(seq)-1] != "api" {
		t.Errorf("expected api last, got %v", seq)
	}
}

func TestCyclicDependency(t *testing.T) {
	m := NewManager()
	_ = m.Register("a", func(_ context.Context) error { return nil }, "b")
	_ = m.Register("b", func(_ context.Context) error { return nil }, "a")

	if err := m.Shutdown(context.Background()); err == nil {
		t.Fatal("expected cyclic dependency error")
	}
}

func TestHandlerError(t *testing.T) {
	m := NewManager()
	expected := errors.New("shutdown failed")
	_ = m.Register("failing", func(_ context.Context) error { return expected })

	err := m.Shutdown(context.Background())
	if err == nil {
		t.Fatal("expected error from handler")
	}
	if !errors.Is(err, expected) {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestContextCancellation(t *testing.T) {
	m := NewManager()
	_ = m.Register("slow", func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
			return nil
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := m.Shutdown(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
