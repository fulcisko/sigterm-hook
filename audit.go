package hook

import (
	"sync"
	"time"
)

// ShutdownAuditEntry records the outcome of a single handler execution.
type ShutdownAuditEntry struct {
	HandlerName string
	Phase       Phase
	Priority    int
	StartedAt   time.Time
	FinishedAt  time.Time
	Duration    time.Duration
	Err         error
	Panicked    bool
	Retries     int
}

// ShutdownAuditLog collects audit entries produced during a shutdown run.
type ShutdownAuditLog struct {
	mu      sync.Mutex
	entries []ShutdownAuditEntry
}

func newAuditLog() *ShutdownAuditLog {
	return &ShutdownAuditLog{}
}

// Record appends an entry to the audit log.
func (a *ShutdownAuditLog) Record(entry ShutdownAuditEntry) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = append(a.entries, entry)
}

// Entries returns a snapshot of all recorded entries in the order they were
// appended. The returned slice is a copy and safe to use after shutdown.
func (a *ShutdownAuditLog) Entries() []ShutdownAuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]ShutdownAuditEntry, len(a.entries))
	copy(out, a.entries)
	return out
}

// Failed returns only the entries whose handler returned an error or panicked.
func (a *ShutdownAuditLog) Failed() []ShutdownAuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	var out []ShutdownAuditEntry
	for _, e := range a.entries {
		if e.Err != nil || e.Panicked {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the total number of recorded entries.
func (a *ShutdownAuditLog) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.entries)
}
