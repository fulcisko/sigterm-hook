package sigtermhook_test

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	sigtermhook "github.com/example/sigterm-hook"
)

// captureLogger records log messages for assertions in tests.
type captureLogger struct {
	mu   sync.Mutex
	logs []string
}

func (c *captureLogger) Infof(format string, args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logs = append(c.logs, fmt.Sprintf(format, args...))
}

func (c *captureLogger) Errorf(format string, args ...interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logs = append(c.logs, fmt.Sprintf(format, args...))
}

func (c *captureLogger) contains(substr string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, l := range c.logs {
		if strings.Contains(l, substr) {
			return true
		}
	}
	return false
}

// all returns a copy of all recorded log messages.
func (c *captureLogger) all() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := make([]string, len(c.logs))
	copy(result, c.logs)
	return result
}

func TestCustomLogger(t *testing.T) {
	logger := &captureLogger{}
	m := sigtermhook.NewManager().WithLogger(logger)

	m.Register("svc", func(_ context.Context) error { return nil })
	m.Shutdown(context.Background())

	if !logger.contains("svc") {
		t.Errorf("expected logger to record shutdown of 'svc'; got logs: %v", logger.all())
	}
}

func TestNoopLogger(t *testing.T) {
	// Should not panic or produce output.
	m := sigtermhook.NewManager().WithNoopLogger()
	m.Register("quiet", func(_ context.Context) error { return nil })
	if errs := m.Shutdown(context.Background()); len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}

func TestNilLoggerFallsBackToNoop(t *testing.T) {
	m := sigtermhook.NewManager().WithLogger(nil)
	m.Register("x", func(_ context.Context) error { return nil })
	if errs := m.Shutdown(context.Background()); len(errs) != 0 {
		t.Errorf("unexpected errors: %v", errs)
	}
}
