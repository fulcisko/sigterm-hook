package sigterm

import (
	"sync"
	"time"
)

// CheckpointStatus represents the result of a handler checkpoint.
type CheckpointStatus int

const (
	CheckpointPending CheckpointStatus = iota
	CheckpointPassed
	CheckpointFailed
)

// CheckpointEntry records the status of a single handler at shutdown time.
type CheckpointEntry struct {
	HandlerName string
	Status      CheckpointStatus
	Duration    time.Duration
	Err         error
	Timestamp   time.Time
}

// checkpointTracker collects pass/fail outcomes for each handler during shutdown.
type checkpointTracker struct {
	mu      sync.Mutex
	entries []CheckpointEntry
}

func newCheckpointTracker() *checkpointTracker {
	return &checkpointTracker{}
}

// Record stores the outcome of a handler execution.
func (c *checkpointTracker) Record(name string, d time.Duration, err error) {
	status := CheckpointPassed
	if err != nil {
		status = CheckpointFailed
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = append(c.entries, CheckpointEntry{
		HandlerName: name,
		Status:      status,
		Duration:    d,
		Err:         err,
		Timestamp:   time.Now(),
	})
}

// Entries returns a snapshot of all recorded checkpoints.
func (c *checkpointTracker) Entries() []CheckpointEntry {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]CheckpointEntry, len(c.entries))
	copy(out, c.entries)
	return out
}

// AllPassed reports whether every recorded checkpoint passed.
func (c *checkpointTracker) AllPassed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, e := range c.entries {
		if e.Status != CheckpointPassed {
			return false
		}
	}
	return true
}

// FailedNames returns the names of handlers that failed their checkpoint.
func (c *checkpointTracker) FailedNames() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var names []string
	for _, e := range c.entries {
		if e.Status == CheckpointFailed {
			names = append(names, e.HandlerName)
		}
	}
	return names
}
