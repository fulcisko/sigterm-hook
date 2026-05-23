package sigterm

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errFlaky = errors.New("transient error")

func TestWrapWithRetrySuccess(t *testing.T) {
	calls := 0
	fn := func(_ context.Context) error {
		calls++
		return nil
	}
	wrapped := wrapWithRetry(fn, DefaultRetryConfig())
	if err := wrapped(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWrapWithRetryEventualSuccess(t *testing.T) {
	calls := 0
	fn := func(_ context.Context) error {
		calls++
		if calls < 3 {
			return errFlaky
		}
		return nil
	}
	cfg := RetryConfig{MaxAttempts: 3, Delay: time.Millisecond}
	wrapped := wrapWithRetry(fn, cfg)
	if err := wrapped(context.Background()); err != nil {
		t.Fatalf("expected success on third attempt, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestWrapWithRetryAllFail(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3, Delay: time.Millisecond}
	wrapped := wrapWithRetry(func(_ context.Context) error { return errFlaky }, cfg)
	err := wrapped(context.Background())
	if err == nil {
		t.Fatal("expected error after all attempts")
	}
	if !errors.Is(err, errFlaky) {
		t.Fatalf("expected wrapped errFlaky, got %v", err)
	}
}

func TestWrapWithRetryContextCancelled(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 5, Delay: 500 * time.Millisecond}
	wrapped := wrapWithRetry(func(_ context.Context) error { return errFlaky }, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	err := wrapped(ctx)
	if err == nil {
		t.Fatal("expected error due to context cancellation")
	}
}

func TestWrapWithRetryOnRetryCallback(t *testing.T) {
	retries := 0
	cfg := RetryConfig{
		MaxAttempts: 3,
		Delay:       time.Millisecond,
		OnRetry:     func(_ int, _ error) { retries++ },
	}
	wrapped := wrapWithRetry(func(_ context.Context) error { return errFlaky }, cfg)
	_ = wrapped(context.Background())
	if retries != 2 {
		t.Fatalf("expected 2 OnRetry calls, got %d", retries)
	}
}

func TestWrapWithRetryMaxAttemptsOne(t *testing.T) {
	calls := 0
	fn := func(_ context.Context) error { calls++; return errFlaky }
	wrapped := wrapWithRetry(fn, RetryConfig{MaxAttempts: 1})
	_ = wrapped(context.Background())
	if calls != 1 {
		t.Fatalf("expected exactly 1 call, got %d", calls)
	}
}
