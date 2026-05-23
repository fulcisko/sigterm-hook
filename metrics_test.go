package sigterm

import (
	"errors"
	"testing"
	"time"
)

func TestNewShutdownMetrics(t *testing.T) {
	m := newShutdownMetrics()
	if m.HandlerDurations == nil {
		t.Fatal("expected HandlerDurations map to be initialised")
	}
	if m.HandlerErrors == nil {
		t.Fatal("expected HandlerErrors map to be initialised")
	}
}

func TestRecordStartAndFinish(t *testing.T) {
	m := newShutdownMetrics()
	before := time.Now()
	m.recordStart()
	time.Sleep(5 * time.Millisecond)
	m.recordFinish()
	after := time.Now()

	if m.StartedAt.Before(before) {
		t.Errorf("StartedAt %v is before test start %v", m.StartedAt, before)
	}
	if m.FinishedAt.After(after) {
		t.Errorf("FinishedAt %v is after test end %v", m.FinishedAt, after)
	}
	if m.TotalDuration < 5*time.Millisecond {
		t.Errorf("TotalDuration %v is less than expected minimum", m.TotalDuration)
	}
}

func TestRecordHandler(t *testing.T) {
	m := newShutdownMetrics()
	sentinelErr := errors.New("handler failed")

	m.recordHandler("db", 10*time.Millisecond, nil)
	m.recordHandler("cache", 20*time.Millisecond, sentinelErr)

	if d := m.HandlerDurations["db"]; d != 10*time.Millisecond {
		t.Errorf("expected db duration 10ms, got %v", d)
	}
	if d := m.HandlerDurations["cache"]; d != 20*time.Millisecond {
		t.Errorf("expected cache duration 20ms, got %v", d)
	}
	if _, ok := m.HandlerErrors["db"]; ok {
		t.Error("expected no error recorded for db")
	}
	if err := m.HandlerErrors["cache"]; !errors.Is(err, sentinelErr) {
		t.Errorf("expected sentinel error for cache, got %v", err)
	}
}

func TestSnapshot(t *testing.T) {
	m := newShutdownMetrics()
	m.recordHandler("svc", 15*time.Millisecond, nil)
	m.recordStart()
	m.recordFinish()

	snap := m.Snapshot()

	// Mutate original; snapshot should be unaffected.
	m.HandlerDurations["svc"] = 99 * time.Second

	if snap.HandlerDurations["svc"] != 15*time.Millisecond {
		t.Errorf("snapshot was mutated; expected 15ms, got %v", snap.HandlerDurations["svc"])
	}
	if snap.TotalDuration <= 0 {
		t.Error("expected positive TotalDuration in snapshot")
	}
}
