package sigterm

import (
	"sync"
	"time"
)

// ShutdownMetrics holds statistics collected during a shutdown sequence.
type ShutdownMetrics struct {
	mu sync.RWMutex

	HandlerDurations map[string]time.Duration
	HandlerErrors    map[string]error
	TotalDuration    time.Duration
	StartedAt        time.Time
	FinishedAt       time.Time
}

// newShutdownMetrics initialises an empty ShutdownMetrics.
func newShutdownMetrics() *ShutdownMetrics {
	return &ShutdownMetrics{
		HandlerDurations: make(map[string]time.Duration),
		HandlerErrors:    make(map[string]error),
	}
}

// recordStart marks the beginning of the shutdown sequence.
func (m *ShutdownMetrics) recordStart() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.StartedAt = time.Now()
}

// recordFinish marks the end of the shutdown sequence.
func (m *ShutdownMetrics) recordFinish() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FinishedAt = time.Now()
	m.TotalDuration = m.FinishedAt.Sub(m.StartedAt)
}

// recordHandler stores the duration and optional error for a named handler.
func (m *ShutdownMetrics) recordHandler(name string, d time.Duration, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.HandlerDurations[name] = d
	if err != nil {
		m.HandlerErrors[name] = err
	}
}

// Snapshot returns a read-only copy of the current metrics.
func (m *ShutdownMetrics) Snapshot() ShutdownMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	durations := make(map[string]time.Duration, len(m.HandlerDurations))
	for k, v := range m.HandlerDurations {
		durations[k] = v
	}
	errors := make(map[string]error, len(m.HandlerErrors))
	for k, v := range m.HandlerErrors {
		errors[k] = v
	}
	return ShutdownMetrics{
		HandlerDurations: durations,
		HandlerErrors:    errors,
		TotalDuration:    m.TotalDuration,
		StartedAt:        m.StartedAt,
		FinishedAt:       m.FinishedAt,
	}
}
