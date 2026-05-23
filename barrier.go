package sigterm

import (
	"context"
	"sync"
)

// Barrier blocks shutdown from proceeding until all registered waiters
// have signalled readiness. This is useful when components need to finish
// in-flight work before teardown begins.
type Barrier struct {
	mu      sync.Mutex
	waiters int
	ready   chan struct{}
	closed  bool
}

// newBarrier creates a new Barrier instance.
func newBarrier() *Barrier {
	return &Barrier{
		ready: make(chan struct{}),
	}
}

// Add registers n additional waiters that must call Done before the barrier
// is released. Add must not be called after the barrier has been closed.
func (b *Barrier) Add(n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		panic("sigterm: Barrier.Add called after barrier is closed")
	}
	b.waiters += n
}

// Done signals that one waiter has completed. When all waiters are done,
// the barrier is released and Wait unblocks.
func (b *Barrier) Done() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.waiters <= 0 {
		panic("sigterm: Barrier.Done called more times than Add")
	}
	b.waiters--
	if b.waiters == 0 && !b.closed {
		b.closed = true
		close(b.ready)
	}
}

// Release immediately unblocks the barrier regardless of remaining waiters.
// Useful for forcing shutdown when context is cancelled.
func (b *Barrier) Release() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.closed {
		b.closed = true
		close(b.ready)
	}
}

// Wait blocks until all waiters call Done, the barrier is released, or the
// context is cancelled. Returns ctx.Err() if context expires first.
func (b *Barrier) Wait(ctx context.Context) error {
	b.mu.Lock()
	if b.waiters == 0 && !b.closed {
		b.closed = true
		close(b.ready)
	}
	b.mu.Unlock()

	select {
	case <-b.ready:
		return nil
	case <-ctx.Done():
		b.Release()
		return ctx.Err()
	}
}
