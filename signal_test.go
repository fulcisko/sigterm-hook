package sigtermhook

import (
	"context"
	"syscall"
	"testing"
	"time"
)

func TestNewSignalWatcherDefaultSignals(t *testing.T) {
	m := NewManager()
	sw := NewSignalWatcher(m)
	if len(sw.signals) != len(DefaultShutdownSignals) {
		t.Fatalf("expected %d default signals, got %d", len(DefaultShutdownSignals), len(sw.signals))
	}
}

func TestNewSignalWatcherCustomSignals(t *testing.T) {
	m := NewManager()
	sw := NewSignalWatcher(m, syscall.SIGUSR1)
	if len(sw.signals) != 1 || sw.signals[0] != syscall.SIGUSR1 {
		t.Fatalf("expected SIGUSR1, got %v", sw.signals)
	}
}

func TestSignalWatcherContextCancelled(t *testing.T) {
	m := NewManager()
	sw := NewSignalWatcher(m)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := sw.Watch(ctx)
	if err != context.Canceled {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestSignalWatcherTriggersShutdown(t *testing.T) {
	m := NewManager()
	shutdownCalled := false
	_ = m.Register("svc", func(ctx context.Context) error {
		shutdownCalled = true
		return nil
	})

	sw := NewSignalWatcher(m, syscall.SIGUSR2)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	errCh := sw.WatchBackground(ctx)

	// Give the goroutine time to start listening.
	time.Sleep(20 * time.Millisecond)

	// Send the signal to ourselves.
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGUSR2); err != nil {
		t.Fatalf("failed to send signal: %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("unexpected shutdown error: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("timed out waiting for shutdown")
	}

	if !shutdownCalled {
		t.Fatal("expected shutdown handler to be called")
	}
}

func TestWatchBackgroundReturnsChannel(t *testing.T) {
	m := NewManager()
	sw := NewSignalWatcher(m)

	ctx, cancel := context.WithCancel(context.Background())
	ch := sw.WatchBackground(ctx)
	cancel()

	select {
	case err := <-ch:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for WatchBackground result")
	}
}
