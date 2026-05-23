package sigterm

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBarrierNoWaiters(t *testing.T) {
	b := newBarrier()
	ctx := context.Background()
	if err := b.Wait(ctx); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestBarrierWaitsForDone(t *testing.T) {
	b := newBarrier()
	b.Add(2)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		b.Wait(context.Background()) //nolint
	}()

	time.Sleep(20 * time.Millisecond)
	b.Done()
	b.Done()
	wg.Wait()
}

func TestBarrierContextCancelled(t *testing.T) {
	b := newBarrier()
	b.Add(1)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	err := b.Wait(ctx)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

func TestBarrierRelease(t *testing.T) {
	b := newBarrier()
	b.Add(3)

	done := make(chan error, 1)
	go func() {
		done <- b.Wait(context.Background())
	}()

	time.Sleep(10 * time.Millisecond)
	b.Release()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("expected nil after Release, got %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait did not unblock after Release")
	}
}

func TestBarrierDonePanicOnExcess(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on excess Done call")
		}
	}()
	b := newBarrier()
	b.Add(1)
	b.Done()
	b.Done() // should panic
}

func TestBarrierAddAfterClosePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on Add after close")
		}
	}()
	b := newBarrier()
	b.Wait(context.Background()) //nolint
	b.Add(1)                     // should panic
}
