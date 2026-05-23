package sigterm

import (
	"testing"
)

func TestNewHandlerGroup(t *testing.T) {
	g := newHandlerGroup()
	if g == nil {
		t.Fatal("expected non-nil HandlerGroup")
	}
	if g.size() != 0 {
		t.Fatalf("expected size 0, got %d", g.size())
	}
}

func TestHandlerGroupAdd(t *testing.T) {
	g := newHandlerGroup()
	g.add(&registeredHandler{name: "a"})
	g.add(&registeredHandler{name: "b"})
	if g.size() != 2 {
		t.Fatalf("expected size 2, got %d", g.size())
	}
}

func TestHandlerGroupAllSnapshot(t *testing.T) {
	g := newHandlerGroup()
	h := &registeredHandler{name: "x"}
	g.add(h)
	snap := g.all()
	if len(snap) != 1 {
		t.Fatalf("expected 1 handler in snapshot, got %d", len(snap))
	}
	// Mutating snapshot must not affect group.
	snap[0] = nil
	if g.all()[0] == nil {
		t.Fatal("snapshot mutation affected internal group state")
	}
}

func TestGroupHandlers(t *testing.T) {
	handlers := []*registeredHandler{
		{name: "a", opts: registerOptions{phase: PhasePre, priority: PriorityHigh}},
		{name: "b", opts: registerOptions{phase: PhasePre, priority: PriorityHigh}},
		{name: "c", opts: registerOptions{phase: PhaseDefault, priority: PriorityNormal}},
	}

	groups := groupHandlers(handlers)

	key1 := groupKey{phase: PhasePre, priority: PriorityHigh}
	key2 := groupKey{phase: PhaseDefault, priority: PriorityNormal}

	if groups[key1].size() != 2 {
		t.Fatalf("expected 2 handlers in PhasePre/High group, got %d", groups[key1].size())
	}
	if groups[key2].size() != 1 {
		t.Fatalf("expected 1 handler in PhaseDefault/Normal group, got %d", groups[key2].size())
	}
}

func TestGroupHandlersEmpty(t *testing.T) {
	groups := groupHandlers(nil)
	if len(groups) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(groups))
	}
}
