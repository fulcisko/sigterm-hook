package sigterm

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// RunbookEntry describes a single handler's shutdown instructions.
type RunbookEntry struct {
	// Name is the handler identifier.
	Name string
	// Description explains what the handler does during shutdown.
	Description string
	// EstimatedDuration is the expected time for the handler to complete.
	EstimatedDuration time.Duration
	// Dependencies lists handler names this entry depends on.
	Dependencies []string
	// Phase is the shutdown phase this handler belongs to.
	Phase Phase
}

// Runbook holds the collection of registered handler descriptions.
type Runbook struct {
	entries []RunbookEntry
}

// newRunbook creates an empty Runbook.
func newRunbook() *Runbook {
	return &Runbook{}
}

// record appends a RunbookEntry to the runbook.
func (r *Runbook) record(e RunbookEntry) {
	r.entries = append(r.entries, e)
}

// Entries returns a copy of all recorded entries.
func (r *Runbook) Entries() []RunbookEntry {
	out := make([]RunbookEntry, len(r.entries))
	copy(out, r.entries)
	return out
}

// WriteTo writes a human-readable runbook summary to w.
func (r *Runbook) WriteTo(w io.Writer) (int64, error) {
	var sb strings.Builder
	sb.WriteString("=== Shutdown Runbook ===\n")
	for _, e := range r.entries {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", e.Phase, e.Name))
		if e.Description != "" {
			sb.WriteString(fmt.Sprintf("  Description : %s\n", e.Description))
		}
		if e.EstimatedDuration > 0 {
			sb.WriteString(fmt.Sprintf("  Est. Duration: %s\n", e.EstimatedDuration))
		}
		if len(e.Dependencies) > 0 {
			sb.WriteString(fmt.Sprintf("  Depends on  : %s\n", strings.Join(e.Dependencies, ", ")))
		}
	}
	n, err := fmt.Fprint(w, sb.String())
	return int64(n), err
}
