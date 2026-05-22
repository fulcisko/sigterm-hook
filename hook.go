package sigterm

import (
	"context"
	"fmt"
	"sync"
)

// Hook represents a shutdown handler with optional dependencies.
type Hook struct {
	name     string
	deps     []string
	handler  func(ctx context.Context) error
}

// Manager orchestrates graceful shutdown with dependency-aware teardown ordering.
type Manager struct {
	mu    sync.Mutex
	hooks map[string]*Hook
}

// NewManager creates a new shutdown Manager.
func NewManager() *Manager {
	return &Manager{
		hooks: make(map[string]*Hook),
	}
}

// Register adds a named shutdown hook with optional dependencies.
// Dependencies are names of other hooks that must complete before this one runs.
func (m *Manager) Register(name string, handler func(ctx context.Context) error, deps ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if name == "" {
		return fmt.Errorf("hook name must not be empty")
	}
	if handler == nil {
		return fmt.Errorf("handler must not be nil")
	}
	if _, exists := m.hooks[name]; exists {
		return fmt.Errorf("hook %q already registered", name)
	}

	m.hooks[name] = &Hook{
		name:    name,
		deps:    deps,
		handler: handler,
	}
	return nil
}

// Shutdown executes all registered hooks in dependency-aware order.
// Hooks whose dependencies have completed are run concurrently.
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	hooks := make(map[string]*Hook, len(m.hooks))
	for k, v := range m.hooks {
		hooks[k] = v
	}
	m.mu.Unlock()

	order, err := topoSort(hooks)
	if err != nil {
		return fmt.Errorf("shutdown ordering failed: %w", err)
	}

	return runInOrder(ctx, order, hooks)
}
