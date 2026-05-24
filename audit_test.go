package hook

import (
	"errors"
	"testing"
	"time"
)

func TestAuditLogEmpty(t *testing.T) {
	log := newAuditLog()
	if log.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", log.Len())
	}
	if len(log.Entries()) != 0 {
		t.Fatal("expected empty entries slice")
	}
	if len(log.Failed()) != 0 {
		t.Fatal("expected empty failed slice")
	}
}

func TestAuditLogRecord(t *testing.T) {
	log := newAuditLog()
	now := time.Now()
	log.Record(ShutdownAuditEntry{
		HandlerName: "db",
		Phase:       PhaseDefault,
		Priority:    PriorityNormal,
		StartedAt:   now,
		FinishedAt:  now.Add(10 * time.Millisecond),
		Duration:    10 * time.Millisecond,
	})
	if log.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", log.Len())
	}
	entries := log.Entries()
	if entries[0].HandlerName != "db" {
		t.Errorf("unexpected handler name: %s", entries[0].HandlerName)
	}
}

func TestAuditLogFailed(t *testing.T) {
	log := newAuditLog()
	log.Record(ShutdownAuditEntry{HandlerName: "ok"})
	log.Record(ShutdownAuditEntry{HandlerName: "bad", Err: errors.New("boom")})
	log.Record(ShutdownAuditEntry{HandlerName: "panic", Panicked: true})

	if log.Len() != 3 {
		t.Fatalf("expected 3 total entries, got %d", log.Len())
	}
	failed := log.Failed()
	if len(failed) != 2 {
		t.Fatalf("expected 2 failed entries, got %d", len(failed))
	}
	names := map[string]bool{failed[0].HandlerName: true, failed[1].HandlerName: true}
	if !names["bad"] || !names["panic"] {
		t.Errorf("unexpected failed handler names: %v", names)
	}
}

func TestAuditLogEntriesIsCopy(t *testing.T) {
	log := newAuditLog()
	log.Record(ShutdownAuditEntry{HandlerName: "svc"})
	snap := log.Entries()
	snap[0].HandlerName = "mutated"

	orig := log.Entries()
	if orig[0].HandlerName == "mutated" {
		t.Error("Entries() returned a reference to internal slice")
	}
}

func TestAuditLogRetries(t *testing.T) {
	log := newAuditLog()
	log.Record(ShutdownAuditEntry{HandlerName: "flaky", Retries: 3})
	entries := log.Entries()
	if entries[0].Retries != 3 {
		t.Errorf("expected 3 retries, got %d", entries[0].Retries)
	}
}
