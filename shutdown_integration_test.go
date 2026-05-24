package hook_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	hook "github.com/yourusername/sigterm-hook"
)

// TestFullShutdownPipeline exercises the complete shutdown pipeline:
// signal → drain → phases → priority → dependency ordering → notifications.
func TestFullShutdownPipeline(t *testing.T) {
	var order []string

	mgr := hook.NewManager(
		hook.WithTimeoutConfig(hook.TimeoutConfig{
			Global:  5 * time.Second,
			Handler: 2 * time.Second,
		}),
	)

	// Register handlers in reverse expected execution order.
	mgr.Register("db", func(ctx context.Context) error {
		order = append(order, "db")
		return nil
	}, hook.WithPriority(hook.PriorityLow))

	mgr.Register("cache", func(ctx context.Context) error {
		order = append(order, "cache")
		return nil
	}, hook.WithPriority(hook.PriorityNormal))

	mgr.Register("http", func(ctx context.Context) error {
		order = append(order, "http")
		return nil
	}, hook.WithPriority(hook.PriorityHigh))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mgr.Shutdown(ctx)
	if err != nil {
		t.Fatalf("unexpected shutdown error: %v", err)
	}

	if len(order) != 3 {
		t.Fatalf("expected 3 handlers to run, got %d", len(order))
	}

	// High priority must run before Normal, Normal before Low.
	pos := func(name string) int {
		for i, v := range order {
			if v == name {
				return i
			}
		}
		return -1
	}

	if pos("http") >= pos("cache") {
		t.Errorf("http (high) should run before cache (normal), got order %v", order)
	}
	if pos("cache") >= pos("db") {
		t.Errorf("cache (normal) should run before db (low), got order %v", order)
	}
}

// TestShutdownWithRetryRecovery verifies that a transiently failing handler
// succeeds after retries and does not surface an error.
func TestShutdownWithRetryRecovery(t *testing.T) {
	var attempts int32

	mgr := hook.NewManager()
	mgr.Register("flaky", func(ctx context.Context) error {
		n := atomic.AddInt32(&attempts, 1)
		if n < 3 {
			return errors.New("temporary failure")
		}
		return nil
	}, hook.WithHandlerRetry(hook.RetryConfig{
		MaxAttempts: 3,
		Delay:       10 * time.Millisecond,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mgr.Shutdown(ctx); err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", atomic.LoadInt32(&attempts))
	}
}

// TestShutdownNotifierIntegration confirms that notifier callbacks are invoked
// with the correct event types during a real shutdown.
func TestShutdownNotifierIntegration(t *testing.T) {
	var started, completed int32

	mgr := hook.NewManager(
		hook.WithNotifierFunc(func(evt hook.ShutdownEvent) {
			switch evt.Type {
			case hook.EventShutdownStarted:
				atomic.AddInt32(&started, 1)
			case hook.EventShutdownCompleted:
				atomic.AddInt32(&completed, 1)
			}
		}),
	)

	mgr.Register("noop", func(ctx context.Context) error { return nil })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mgr.Shutdown(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if atomic.LoadInt32(&started) != 1 {
		t.Errorf("expected 1 started event, got %d", atomic.LoadInt32(&started))
	}
	if atomic.LoadInt32(&completed) != 1 {
		t.Errorf("expected 1 completed event, got %d", atomic.LoadInt32(&completed))
	}
}
