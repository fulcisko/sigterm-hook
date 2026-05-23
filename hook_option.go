package sigterm

import (
	"time"
)

// registerOptions holds the configuration for a single registered handler.
type registerOptions struct {
	priority Priority
	phase    Phase
	timeout  time.Duration
	retry    *RetryConfig
	drainer  *Drainer
}

// defaultRegisterOptions returns sensible defaults for handler registration.
func defaultRegisterOptions() registerOptions {
	return registerOptions{
		priority: PriorityNormal,
		phase:    PhaseDefault,
	}
}

// RegisterOption is a functional option for Register.
type RegisterOption func(*registerOptions)

// WithPriority sets the execution priority of the handler.
func WithPriority(p Priority) RegisterOption {
	return func(o *registerOptions) {
		o.priority = p
	}
}

// WithHandlerPhase assigns the handler to a specific shutdown phase.
func WithHandlerPhase(p Phase) RegisterOption {
	return func(o *registerOptions) {
		o.phase = p
	}
}

// WithHandlerTimeout sets a per-handler timeout that overrides the global one.
func WithHandlerTimeout(d time.Duration) RegisterOption {
	return func(o *registerOptions) {
		o.timeout = d
	}
}

// WithHandlerRetry attaches a RetryConfig to the handler.
func WithHandlerRetry(cfg RetryConfig) RegisterOption {
	return func(o *registerOptions) {
		o.retry = &cfg
	}
}

// WithDrainer attaches a Drainer that must reach zero in-flight operations
// before the handler is invoked. The drainer is closed (no new acquisitions
// allowed) and the handler waits for it to drain within the handler's timeout.
func WithDrainer(d *Drainer) RegisterOption {
	return func(o *registerOptions) {
		o.drainer = d
	}
}
