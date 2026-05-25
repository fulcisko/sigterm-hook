package sigterm

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ShutdownBudget tracks remaining time across the shutdown pipeline and
// provides a child context that expires when the budget is exhausted.
type ShutdownBudget struct {
	mu        sync.Mutex
	total     time.Duration
	start     time.Time
	spent     time.Duration
	cancelled bool
}

// newShutdownBudget creates a ShutdownBudget with the given total duration.
// A zero or negative duration means no budget is enforced.
func newShutdownBudget(total time.Duration) *ShutdownBudget {
	return &ShutdownBudget{
		total: total,
		start: time.Now(),
	}
}

// Remaining returns the time left in the budget. If the budget is unlimited
// (total <= 0), it returns a large sentinel value.
func (b *ShutdownBudget) Remaining() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.total <= 0 {
		return 24 * time.Hour
	}
	elapsed := time.Since(b.start)
	rem := b.total - elapsed
	if rem < 0 {
		return 0
	}
	return rem
}

// Spend records that d of budget was consumed by a handler.
func (b *ShutdownBudget) Spend(d time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.spent += d
}

// Exhausted reports whether the budget has been fully consumed.
func (b *ShutdownBudget) Exhausted() bool {
	return b.Remaining() == 0
}

// WithContext returns a derived context that is cancelled when the budget
// expires. If the budget is unlimited the parent context is returned unchanged.
func (b *ShutdownBudget) WithContext(parent context.Context) (context.Context, context.CancelFunc) {
	rem := b.Remaining()
	if b.total <= 0 {
		return parent, func() {}
	}
	if rem == 0 {
		ctx, cancel := context.WithCancel(parent)
		cancel()
		return ctx, cancel
	}
	return context.WithDeadline(parent, time.Now().Add(rem))
}

// String returns a human-readable summary of the budget.
func (b *ShutdownBudget) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.total <= 0 {
		return "budget: unlimited"
	}
	return fmt.Sprintf("budget: total=%s remaining=%s spent=%s",
		b.total, b.Remaining(), b.spent)
}
