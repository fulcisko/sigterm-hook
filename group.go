package sigterm

import "sync"

// HandlerGroup represents a collection of handlers that execute concurrently
// within the same phase and priority level.
type HandlerGroup struct {
	mu       sync.Mutex
	handlers []*registeredHandler
}

// newHandlerGroup creates an empty HandlerGroup.
func newHandlerGroup() *HandlerGroup {
	return &HandlerGroup{
		handlers: make([]*registeredHandler, 0),
	}
}

// add appends a handler to the group.
func (g *HandlerGroup) add(h *registeredHandler) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.handlers = append(g.handlers, h)
}

// all returns a snapshot of the current handlers.
func (g *HandlerGroup) all() []*registeredHandler {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make([]*registeredHandler, len(g.handlers))
	copy(out, g.handlers)
	return out
}

// size returns the number of handlers in the group.
func (g *HandlerGroup) size() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return len(g.handlers)
}

// groupKey uniquely identifies a bucket of handlers sharing phase and priority.
type groupKey struct {
	phase    Phase
	priority Priority
}

// groupHandlers partitions a slice of registered handlers into groups keyed by
// (phase, priority), preserving insertion order within each group.
func groupHandlers(handlers []*registeredHandler) map[groupKey]*HandlerGroup {
	result := make(map[groupKey]*HandlerGroup)
	for _, h := range handlers {
		key := groupKey{
			phase:    h.opts.phase,
			priority: h.opts.priority,
		}
		if _, ok := result[key]; !ok {
			result[key] = newHandlerGroup()
		}
		result[key].add(h)
	}
	return result
}
