package sigterm_test

import (
	"context"
	"errors"
	"testing"

	sigterm "github.com/example/sigterm-hook"
)

func TestCheckpointAllPassedIntegration(t *testing.T) {
	m := sigterm.NewManager(sigterm.WithCheckpoint())

	m.Register("alpha", func(_ context.Context) error { return nil })
	m.Register("beta", func(_ context.Context) error { return nil })

	_ = m.Shutdown(context.Background())

	cp := m.Checkpoints()
	if cp == nil {
		t.Fatal("expected checkpoint tracker to be non-nil")
	}
	if !cp.AllPassed() {
		t.Errorf("expected all checkpoints to pass, failed: %v", cp.FailedNames())
	}
	if len(cp.Entries()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(cp.Entries()))
	}
}

func TestCheckpointRecordsFailureIntegration(t *testing.T) {
	m := sigterm.NewManager(sigterm.WithCheckpoint())

	wantErr := errors.New("handler boom")
	m.Register("failing", func(_ context.Context) error { return wantErr })
	m.Register("ok", func(_ context.Context) error { return nil })

	_ = m.Shutdown(context.Background())

	cp := m.Checkpoints()
	if cp.AllPassed() {
		t.Error("expected AllPassed=false")
	}
	failed := cp.FailedNames()
	if len(failed) != 1 || failed[0] != "failing" {
		t.Errorf("unexpected failed names: %v", failed)
	}
}

func TestCheckpointNilWhenNotEnabled(t *testing.T) {
	m := sigterm.NewManager() // no WithCheckpoint
	m.Register("svc", func(_ context.Context) error { return nil })
	_ = m.Shutdown(context.Background())

	if m.Checkpoints() != nil {
		t.Error("expected nil checkpoint tracker when feature not enabled")
	}
}
