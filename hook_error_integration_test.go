package sigterm_test

import (
	"context"
	"errors"
	"testing"

	sigterm "github.com/example/sigterm-hook"
)

// TestShutdownSurfacesHandlerErrors verifies that errors returned by registered
// handlers are wrapped in a ShutdownError and accessible via errors.As.
func TestShutdownSurfacesHandlerErrors(t *testing.T) {
	m := sigterm.NewManager()

	wantErr := errors.New("handler boom")
	m.Register("failing", func(_ context.Context) error {
		return wantErr
	})

	err := m.Shutdown(context.Background())
	if err == nil {
		t.Fatal("expected shutdown error, got nil")
	}

	var se *sigterm.ShutdownError
	if !errors.As(err, &se) {
		t.Fatalf("expected *ShutdownError, got %T: %v", err, err)
	}
	if len(se.Errors) == 0 {
		t.Fatal("expected at least one wrapped error")
	}
	if !errors.Is(err, wantErr) {
		t.Error("errors.Is should find original error through ShutdownError")
	}
}

// TestShutdownNoErrorWhenAllSucceed verifies that a nil error is returned when
// every handler completes successfully.
func TestShutdownNoErrorWhenAllSucceed(t *testing.T) {
	m := sigterm.NewManager()

	m.Register("ok1", func(_ context.Context) error { return nil })
	m.Register("ok2", func(_ context.Context) error { return nil })

	if err := m.Shutdown(context.Background()); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}
