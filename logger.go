package sigtermhook

import "log"

// Logger defines the interface for shutdown event logging.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// defaultLogger is a Logger backed by the standard library log package.
type defaultLogger struct{}

func (l *defaultLogger) Infof(format string, args ...interface{}) {
	log.Printf("[sigtermhook] INFO  "+format, args...)
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[sigtermhook] ERROR "+format, args...)
}

// noopLogger silently discards all log messages.
type noopLogger struct{}

func (l *noopLogger) Infof(_ string, _ ...interface{})  {}
func (l *noopLogger) Errorf(_ string, _ ...interface{}) {}

// WithLogger sets a custom logger on the Manager.
func (m *Manager) WithLogger(l Logger) *Manager {
	m.mu.Lock()
	defer m.mu.Unlock()
	if l == nil {
		m.logger = &noopLogger{}
	} else {
		m.logger = l
	}
	return m
}

// WithNoopLogger disables all logging for the Manager.
func (m *Manager) WithNoopLogger() *Manager {
	return m.WithLogger(&noopLogger{})
}
