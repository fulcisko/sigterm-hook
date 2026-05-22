package sigtermhook

import (
	"context"
	"fmt"
	"time"
)

// TimeoutConfig holds configuration for shutdown timeout behavior.
type TimeoutConfig struct {
	// GlobalTimeout is the maximum time allowed for the entire shutdown sequence.
	GlobalTimeout time.Duration
	// HandlerTimeout is the maximum time allowed for each individual handler.
	HandlerTimeout time.Duration
}

// DefaultTimeoutConfig returns a TimeoutConfig with sensible defaults.
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		GlobalTimeout:  30 * time.Second,
		HandlerTimeout: 10 * time.Second,
	}
}

// WithGlobalTimeout sets the global shutdown timeout on the Manager.
func (m *Manager) WithGlobalTimeout(d time.Duration) *Manager {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.timeoutCfg.GlobalTimeout = d
	return m
}

// WithHandlerTimeout sets the per-handler timeout on the Manager.
func (m *Manager) WithHandlerTimeout(d time.Duration) *Manager {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.timeoutCfg.HandlerTimeout = d
	return m
}

// wrapWithTimeout wraps a ShutdownFunc with a per-handler deadline.
func wrapWithTimeout(name string, fn ShutdownFunc, timeout time.Duration) ShutdownFunc {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		type result struct {
			err error
		}
		ch := make(chan result, 1)
		go func() {
			ch <- result{err: fn(ctx)}
		}()

		select {
		case res := <-ch:
			return res.err
		case <-ctx.Done():
			return fmt.Errorf("handler %q timed out after %s", name, timeout)
		}
	}
}
