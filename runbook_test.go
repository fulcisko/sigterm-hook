package sigterm

import (
	"strings"
	"testing"
	"time"
)

func TestRunbookEmpty(t *testing.T) {
	rb := newRunbook()
	if got := rb.Entries(); len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestRunbookRecord(t *testing.T) {
	rb := newRunbook()
	rb.record(RunbookEntry{Name: "db", Description: "close db", Phase: PhasePrimary})
	rb.record(RunbookEntry{Name: "cache", Phase: PhaseSecondary})

	entries := rb.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Name != "db" {
		t.Errorf("expected first entry name 'db', got %q", entries[0].Name)
	}
	if entries[1].Phase != PhaseSecondary {
		t.Errorf("expected second entry phase %v, got %v", PhaseSecondary, entries[1].Phase)
	}
}

func TestRunbookEntriesIsCopy(t *testing.T) {
	rb := newRunbook()
	rb.record(RunbookEntry{Name: "svc"})

	a := rb.Entries()
	a[0].Name = "mutated"

	b := rb.Entries()
	if b[0].Name == "mutated" {
		t.Error("Entries() should return a copy; original was mutated")
	}
}

func TestRunbookWriteTo(t *testing.T) {
	rb := newRunbook()
	rb.record(RunbookEntry{
		Name:              "db",
		Description:       "Closes DB connections",
		EstimatedDuration: 3 * time.Second,
		Dependencies:      []string{"cache"},
		Phase:             PhasePrimary,
	})

	var sb strings.Builder
	n, err := rb.WriteTo(&sb)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes written")
	}
	out := sb.String()
	for _, want := range []string{"db", "Closes DB connections", "3s", "cache", "primary"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q\ngot:\n%s", want, out)
		}
	}
}

func TestWithEstimatedDuration(t *testing.T) {
	opts := defaultRegisterOptions("svc")
	applyRegisterOptions(&opts, WithEstimatedDuration(5*time.Second))

	if opts.runbookEntry == nil {
		t.Fatal("expected runbookEntry to be set")
	}
	if opts.runbookEntry.EstimatedDuration != 5*time.Second {
		t.Errorf("expected 5s, got %v", opts.runbookEntry.EstimatedDuration)
	}
}

func TestWithRunbookEntry(t *testing.T) {
	opts := defaultRegisterOptions("svc")
	applyRegisterOptions(&opts, WithRunbookEntry(RunbookEntry{
		Description: "test handler",
		Phase:       PhaseFinal,
	}))

	if opts.runbookEntry == nil {
		t.Fatal("expected runbookEntry to be set")
	}
	if opts.runbookEntry.Description != "test handler" {
		t.Errorf("unexpected description: %q", opts.runbookEntry.Description)
	}
	if opts.runbookEntry.Phase != PhaseFinal {
		t.Errorf("unexpected phase: %v", opts.runbookEntry.Phase)
	}
}
