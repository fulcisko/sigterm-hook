package sigterm

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewDrainer(t *testing.T) {
	d := newDrainer()
	if d == nil {
		t.Fatal("expected non-nil drainer")
	}
	if d.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight, got %d", d.InFlight())
	}
}

func TestDrainerAcquireRelease(t *testing.T) {
	d := newDrainer()

	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed")
	}
	if d.InFlight() != 1 {
		t.Fatalf("expected 1 in-flight, got %d", d.InFlight())
	}

	d.Release()
	if d.InFlight() != 0 {
		t.Fatalf("expected 0 in-flight after release, got %d", d.InFlight())
	}
}

func TestDrainerCloseWaitsForInFlight(t *testing.T) {
	d := newDrainer()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if !d.Acquire() {
			return
		}
		time.Sleep(20 * time.Millisecond)
		d.Release()
	}()

	// Give goroutine time to acquire.
	time.Sleep(5 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	if err := d.Close(ctx); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	wg.Wait()
}

func TestDrainerAcquireAfterClose(t *testing.T) {
	d := newDrainer()
	ctx := context.Background()
	if err := d.Close(ctx); err != nil {
		t.Fatalf("unexpected error on Close: %v", err)
	}
	if d.Acquire() {
		t.Fatal("expected Acquire to fail after Close")
	}
}

func TestDrainerCloseContextCancelled(t *testing.T) {
	d := newDrainer()

	// Hold a reference so Close blocks.
	if !d.Acquire() {
		t.Fatal("expected Acquire to succeed")
	}
	defer d.Release()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := d.Close(ctx)
	if err == nil {
		t.Fatal("expected context deadline error")
	}
}

func TestDrainerCloseIdempotent(t *testing.T) {
	d := newDrainer()
	ctx := context.Background()
	if err := d.Close(ctx); err != nil {
		t.Fatalf("first Close failed: %v", err)
	}
	if err := d.Close(ctx); err != nil {
		t.Fatalf("second Close failed: %v", err)
	}
}
