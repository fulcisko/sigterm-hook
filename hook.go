package sigterm

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// ShutdownFunc is the signature for a handler that performs cleanup work.
type ShutdownFunc func(ctx context.Context) error

// handler holds metadata about a registered shutdown handler.
type handler struct {
	name  string
	fn    ShutdownFunc
	deps  []string
	retry *RetryConfig
}

// Manager orchestrates graceful shutdown of registered handlers.
type Manager struct {
	mu       sync.Mutex
	handlers map[string]*handler
	timeout  TimeoutConfig
	logger   Logger
	metrics  *shutdownMetrics
}

// NewManager creates a Manager with the provided options applied.
func NewManager(opts ...Option) *Manager {
	m := &Manager{
		handlers: make(map[string]*handler),
		timeout:  DefaultTimeoutConfig(),
		logger:   newNoopLogger(),
		metrics:  newShutdownMetrics(),
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Register adds a named shutdown handler with optional dependencies and retry config.
func (m *Manager) Register(name string, fn ShutdownFunc, deps []string, retry *RetryConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.handlers[name]; exists {
		return fmt.Errorf("handler %q already registered", name)
	}
	if retry != nil {
		fn = wrapWithRetry(fn, *retry)
	}
	m.handlers[name] = &handler{name: name, fn: fn, deps: deps, retry: retry}
	return nil
}

// Unregister removes a previously registered handler by name.
// Returns an error if no handler with that name exists.
func (m *Manager) Unregister(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.handlers[name]; !exists {
		return fmt.Errorf("handler %q not found", name)
	}
	delete(m.handlers, name)
	return nil
}

// Shutdown runs all registered handlers in dependency order.
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	order, err := topoSort(m.handlers)
	m.mu.Unlock()
	if err != nil {
		return err
	}
	m.metrics.recordStart()
	defer m.metrics.recordFinish()
	return runInOrder(ctx, order, m.handlers, m.timeout, m.logger, m.metrics)
}

// ListenAndServe blocks until SIGTERM or SIGINT is received, then shuts down.
func (m *Manager) ListenAndServe(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(quit)
	select {
	case <-quit:
		m.logger.Info("signal received, starting shutdown")
	case <-ctx.Done():
		m.logger.Info("context done, starting shutdown")
	}
	return m.Shutdown(ctx)
}

// Option is a functional option for Manager.
type Option func(*Manager)

// WithTimeoutConfig sets the timeout configuration on the Manager.
func WithTimeoutConfig(tc TimeoutConfig) Option {
	return func(m *Manager) { m.timeout = tc }
}

// WithLogger sets a custom logger on the Manager.
func WithLogger(l Logger) Option {
	return func(m *Manager) { m.logger = l }
}
