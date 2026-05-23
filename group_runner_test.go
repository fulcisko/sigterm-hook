package sigterm

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestRunGroupAllSuccess(t *testing.T) {
	g := newHandlerGroup()
	g.add(&registeredHandler{name: "h1", fn: func(ctx context.Context) error { return nil }, opts: registerOptions{}})
	g.add(&registeredHandler{name: "h2", fn: func(ctx context.Context) error { return nil }, opts: registerOptions{}})

	results := runGroup(context.Background(), g, nil, noopLogger{})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("unexpected error for handler %q: %v", r.Name, r.Err)
		}
	}
}

func TestRunGroupPartialFailure(t *testing.T) {
	g := newHandlerGroup()
	g.add(&registeredHandler{name: "ok", fn: func(ctx context.Context) error { return nil }, opts: registerOptions{}})
	g.add(&registeredHandler{name: "fail", fn: func(ctx context.Context) error { return errors.New("boom") }, opts: registerOptions{}})

	results := runGroup(context.Background(), g, nil, noopLogger{})
	var gotErr bool
	for _, r := range results {
		if r.Name == "fail" && r.Err != nil {
			gotErr = true
		}
	}
	if !gotErr {
		t.Fatal("expected error from 'fail' handler")
	}
}

func TestCollectErrorsNone(t *testing.T) {
	results := []GroupResult{{Name: "a", Err: nil}, {Name: "b", Err: nil}}
	if err := collectErrors(results); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestCollectErrorsMultiple(t *testing.T) {
	results := []GroupResult{
		{Name: "a", Err: errors.New("err-a")},
		{Name: "b", Err: nil},
		{Name: "c", Err: errors.New("err-c")},
	}
	err := collectErrors(results)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "err-a") || !strings.Contains(err.Error(), "err-c") {
		t.Errorf("error message missing expected content: %v", err)
	}
}

func TestRunGroupContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	g := newHandlerGroup()
	g.add(&registeredHandler{
		name: "ctxcheck",
		fn: func(ctx context.Context) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return nil
		},
		opts: registerOptions{},
	})

	results := runGroup(ctx, g, nil, noopLogger{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Err == nil {
		t.Fatal("expected context cancellation error")
	}
}
