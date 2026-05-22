package sigtermhook_test

import (
	"context"
	"errors"
	"testing"
	"time"

	sigtermhook "github.com/example/sigterm-hook"
)

func TestHandlerTimeout(t *testing.T) {
	m := sigtermhook.NewManager().
		WithNoopLogger().
		WithHandlerTimeout(50 * time.Millisecond)

	m.Register("slow", func(ctx context.Context) error {
		select {
		case <-time.After(5 * time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errs := m.Shutdown(ctx)
	if len(errs) == 0 {
		t.Fatal("expected timeout error, got none")
	}
}

func TestGlobalTimeout(t *testing.T) {
	m := sigtermhook.NewManager().
		WithNoopLogger().
		WithGlobalTimeout(100 * time.Millisecond).
		WithHandlerTimeout(10 * time.Second)

	m.Register("blocker", func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})

	start := time.Now()
	ctx := context.Background()
	errs := m.Shutdown(ctx)
	elapsed := time.Since(start)

	if elapsed > 500*time.Millisecond {
		t.Errorf("global timeout not respected, elapsed: %s", elapsed)
	}
	if len(errs) == 0 {
		t.Fatal("expected errors due to global timeout")
	}
}

func TestDefaultTimeoutConfig(t *testing.T) {
	cfg := sigtermhook.DefaultTimeoutConfig()
	if cfg.GlobalTimeout != 30*time.Second {
		t.Errorf("expected 30s global timeout, got %s", cfg.GlobalTimeout)
	}
	if cfg.HandlerTimeout != 10*time.Second {
		t.Errorf("expected 10s handler timeout, got %s", cfg.HandlerTimeout)
	}
}

func TestHandlerCompletesBeforeTimeout(t *testing.T) {
	m := sigtermhook.NewManager().
		WithNoopLogger().
		WithHandlerTimeout(500 * time.Millisecond)

	m.Register("fast", func(_ context.Context) error {
		return nil
	})

	errs := m.Shutdown(context.Background())
	for _, err := range errs {
		if errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("unexpected timeout for fast handler: %v", err)
		}
	}
}
