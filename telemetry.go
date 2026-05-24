package sigterm

import (
	"sync"
	"time"
)

// TelemetryEvent represents a single recorded event during shutdown.
type TelemetryEvent struct {
	Handler   string
	Phase     Phase
	StartedAt time.Time
	FinishedAt time.Time
	Duration  time.Duration
	Err       error
	Retries   int
}

// TelemetryRecorder collects structured telemetry for each handler
// executed during a shutdown sequence.
type TelemetryRecorder struct {
	mu     sync.Mutex
	events []TelemetryEvent
}

func newTelemetryRecorder() *TelemetryRecorder {
	return &TelemetryRecorder{}
}

// Record appends a completed handler event to the recorder.
func (t *TelemetryRecorder) Record(e TelemetryEvent) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if e.Duration == 0 && !e.StartedAt.IsZero() && !e.FinishedAt.IsZero() {
		e.Duration = e.FinishedAt.Sub(e.StartedAt)
	}
	t.events = append(t.events, e)
}

// Events returns a snapshot copy of all recorded telemetry events.
func (t *TelemetryRecorder) Events() []TelemetryEvent {
	t.mu.Lock()
	defer t.mu.Unlock()
	copy := make([]TelemetryEvent, len(t.events))
	for i, e := range t.events {
		copy[i] = e
	}
	return copy
}

// Failed returns only the events where an error was recorded.
func (t *TelemetryRecorder) Failed() []TelemetryEvent {
	t.mu.Lock()
	defer t.mu.Unlock()
	var out []TelemetryEvent
	for _, e := range t.events {
		if e.Err != nil {
			out = append(out, e)
		}
	}
	return out
}

// TotalDuration sums the Duration field of all recorded events.
func (t *TelemetryRecorder) TotalDuration() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	var total time.Duration
	for _, e := range t.events {
		total += e.Duration
	}
	return total
}
