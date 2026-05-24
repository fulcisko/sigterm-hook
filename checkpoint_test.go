package sigterm

import (
	"errors"
	"testing"
	"time"
)

func TestCheckpointTrackerEmpty(t *testing.T) {
	ct := newCheckpointTracker()
	if !ct.AllPassed() {
		t.Error("empty tracker should report AllPassed=true")
	}
	if len(ct.Entries()) != 0 {
		t.Error("expected no entries")
	}
}

func TestCheckpointTrackerRecordSuccess(t *testing.T) {
	ct := newCheckpointTracker()
	ct.Record("db", 10*time.Millisecond, nil)

	entries := ct.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Status != CheckpointPassed {
		t.Errorf("expected Passed, got %v", entries[0].Status)
	}
	if entries[0].HandlerName != "db" {
		t.Errorf("unexpected name %q", entries[0].HandlerName)
	}
}

func TestCheckpointTrackerRecordFailure(t *testing.T) {
	ct := newCheckpointTracker()
	err := errors.New("timeout")
	ct.Record("cache", 5*time.Millisecond, err)

	if ct.AllPassed() {
		t.Error("expected AllPassed=false after failure")
	}
	failed := ct.FailedNames()
	if len(failed) != 1 || failed[0] != "cache" {
		t.Errorf("unexpected failed names: %v", failed)
	}
}

func TestCheckpointTrackerEntriesIsCopy(t *testing.T) {
	ct := newCheckpointTracker()
	ct.Record("svc", time.Millisecond, nil)

	a := ct.Entries()
	a[0].HandlerName = "mutated"

	b := ct.Entries()
	if b[0].HandlerName == "mutated" {
		t.Error("Entries should return a copy, not a reference")
	}
}

func TestCheckpointTrackerMixed(t *testing.T) {
	ct := newCheckpointTracker()
	ct.Record("a", time.Millisecond, nil)
	ct.Record("b", time.Millisecond, errors.New("fail"))
	ct.Record("c", time.Millisecond, nil)

	if ct.AllPassed() {
		t.Error("expected AllPassed=false")
	}
	failed := ct.FailedNames()
	if len(failed) != 1 || failed[0] != "b" {
		t.Errorf("expected only 'b' to fail, got %v", failed)
	}
}

func TestCheckpointStatusConstants(t *testing.T) {
	if CheckpointPending == CheckpointPassed {
		t.Error("Pending and Passed must differ")
	}
	if CheckpointPassed == CheckpointFailed {
		t.Error("Passed and Failed must differ")
	}
}
