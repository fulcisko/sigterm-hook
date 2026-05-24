package sigterm

import (
	"context"
	"sync"
)

// shutdownKey is the key type for shutdown context values.
type shutdownKey struct{}

// ShutdownContext carries metadata about an in-progress shutdown.
type ShutdownContext struct {
	mu       sync.RWMutex
	reason   string
	handlers []string
}

// Reason returns the human-readable reason the shutdown was initiated.
func (s *ShutdownContext) Reason() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.reason
}

// CompletedHandlers returns a snapshot of handler names that have finished.
func (s *ShutdownContext) CompletedHandlers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.handlers))
	copy(out, s.handlers)
	return out
}

// markHandlerDone records that a handler finished successfully.
func (s *ShutdownContext) markHandlerDone(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers = append(s.handlers, name)
}

// newShutdownContext creates a derived context carrying ShutdownContext metadata.
func newShutdownContext(parent context.Context, reason string) (context.Context, *ShutdownContext) {
	sc := &ShutdownContext{reason: reason}
	return context.WithValue(parent, shutdownKey{}, sc), sc
}

// ShutdownContextFrom retrieves the ShutdownContext from ctx, if present.
// Returns nil when ctx was not created by newShutdownContext.
func ShutdownContextFrom(ctx context.Context) *ShutdownContext {
	v, _ := ctx.Value(shutdownKey{}).(*ShutdownContext)
	return v
}
