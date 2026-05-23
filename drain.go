package sigterm

import (
	"context"
	"sync"
	"sync/atomic"
)

// DrainConfig configures the connection draining behaviour before shutdown
// handlers are invoked.
type DrainConfig struct {
	// MaxWait is the maximum duration to wait for in-flight requests to finish.
	// Defaults to the global shutdown timeout when zero.
	MaxWait int64 // nanoseconds, use time.Duration cast

	// PollInterval is how often the drain loop checks the in-flight counter.
	PollInterval int64 // nanoseconds
}

// Drainer tracks in-flight work and can block until all work is complete or
// the context is cancelled.
type Drainer struct {
	mu       sync.Mutex
	inflight atomic.Int64
	done     chan struct{}
}

// newDrainer creates a ready-to-use Drainer.
func newDrainer() *Drainer {
	return &Drainer{
		done: make(chan struct{}),
	}
}

// Acquire increments the in-flight counter. It returns false if the drainer
// has already been closed (i.e. shutdown has started).
func (d *Drainer) Acquire() bool {
	select {
	case <-d.done:
		return false
	default:
	}
	d.inflight.Add(1)
	// Double-check: if Close raced us, decrement and refuse.
	select {
	case <-d.done:
		d.inflight.Add(-1)
		return false
	default:
		return true
	}
}

// Release decrements the in-flight counter.
func (d *Drainer) Release() {
	d.inflight.Add(-1)
}

// Close marks the drainer as closed so no new work can be acquired, then
// waits until all in-flight work finishes or ctx is cancelled.
func (d *Drainer) Close(ctx context.Context) error {
	d.mu.Lock()
	select {
	case <-d.done:
	default:
		close(d.done)
	}
	d.mu.Unlock()

	for {
		if d.inflight.Load() == 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

// InFlight returns the current number of in-flight operations.
func (d *Drainer) InFlight() int64 {
	return d.inflight.Load()
}
