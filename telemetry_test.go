package sigterm

import (
	"errors"
	"testing"
	"time"
)

func TestTelemetryRecorderEmpty(t *testing.T) {
	rec := newTelemetryRecorder()
	if len(rec.Events()) != 0 {
		t.Fatal("expected empty events")
	}
	if rec.TotalDuration() != 0 {
		t.Fatal("expected zero total duration")
	}
}

func TestTelemetryRecorderRecord(t *testing.T) {
	rec := newTelemetryRecorder()
	now := time.Now()
	rec.Record(TelemetryEvent{
		Handler:    "db",
		StartedAt:  now,
		FinishedAt: now.Add(50 * time.Millisecond),
	})
	events := rec.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Handler != "db" {
		t.Errorf("unexpected handler name: %s", events[0].Handler)
	}
	if events[0].Duration != 50*time.Millisecond {
		t.Errorf("unexpected duration: %v", events[0].Duration)
	}
}

func TestTelemetryRecorderFailed(t *testing.T) {
	rec := newTelemetryRecorder()
	rec.Record(TelemetryEvent{Handler: "ok"})
	rec.Record(TelemetryEvent{Handler: "bad", Err: errors.New("boom")})

	failed := rec.Failed()
	if len(failed) != 1 {
		t.Fatalf("expected 1 failed event, got %d", len(failed))
	}
	if failed[0].Handler != "bad" {
		t.Errorf("wrong failed handler: %s", failed[0].Handler)
	}
}

func TestTelemetryRecorderTotalDuration(t *testing.T) {
	rec := newTelemetryRecorder()
	now := time.Now()
	rec.Record(TelemetryEvent{StartedAt: now, FinishedAt: now.Add(100 * time.Millisecond)})
	rec.Record(TelemetryEvent{StartedAt: now, FinishedAt: now.Add(200 * time.Millisecond)})

	if rec.TotalDuration() != 300*time.Millisecond {
		t.Errorf("unexpected total: %v", rec.TotalDuration())
	}
}

func TestTelemetryEventsIsCopy(t *testing.T) {
	rec := newTelemetryRecorder()
	rec.Record(TelemetryEvent{Handler: "a"})
	snap := rec.Events()
	snap[0].Handler = "mutated"

	original := rec.Events()
	if original[0].Handler != "a" {
		t.Error("Events() should return a copy, not a reference")
	}
}

func TestNewTelemetryRecorderPublic(t *testing.T) {
	rec := NewTelemetryRecorder()
	if rec == nil {
		t.Fatal("expected non-nil recorder")
	}
}
