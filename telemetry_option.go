package sigterm

// WithTelemetry attaches a TelemetryRecorder to the Manager so that
// structured per-handler timing and error data is collected during
// every shutdown run.
//
// Example:
//
//	rec := sigterm.NewTelemetryRecorder()
//	m := sigterm.NewManager(sigterm.WithTelemetry(rec))
//	// after shutdown:
//	for _, ev := range rec.Events() { ... }
func WithTelemetry(rec *TelemetryRecorder) Option {
	return func(m *Manager) {
		m.telemetry = rec
	}
}

// NewTelemetryRecorder is a convenience constructor exposed at the
// package level so callers do not need to import internal types.
func NewTelemetryRecorder() *TelemetryRecorder {
	return newTelemetryRecorder()
}
