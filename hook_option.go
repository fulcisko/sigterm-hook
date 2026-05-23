package sigterm

// RegisterOption configures how a handler is registered.
type RegisterOption func(*registerOptions)

// registerOptions holds all per-handler configuration.
type registerOptions struct {
	priority int
	phase    Phase
	timeout  *HandlerTimeoutConfig
	retry    *RetryConfig
}

// defaultRegisterOptions returns sensible defaults.
func defaultRegisterOptions() registerOptions {
	return registerOptions{
		priority: PriorityNormal,
		phase:    PhaseDefault,
	}
}

// WithPriority sets the execution priority for the handler.
// Higher values run first within the same phase.
func WithPriority(p int) RegisterOption {
	return func(o *registerOptions) {
		o.priority = p
	}
}

// WithHandlerPhase sets the lifecycle phase the handler belongs to.
func WithHandlerPhase(p Phase) RegisterOption {
	return func(o *registerOptions) {
		o.phase = p
	}
}

// WithHandlerTimeout overrides the global handler timeout for this specific handler.
func WithHandlerTimeout(cfg HandlerTimeoutConfig) RegisterOption {
	return func(o *registerOptions) {
		o.timeout = &cfg
	}
}

// WithHandlerRetry attaches a retry policy to this specific handler.
func WithHandlerRetry(cfg RetryConfig) RegisterOption {
	return func(o *registerOptions) {
		o.retry = &cfg
	}
}

// applyRegisterOptions merges a slice of RegisterOption into a registerOptions struct.
func applyRegisterOptions(opts []RegisterOption) registerOptions {
	ro := defaultRegisterOptions()
	for _, o := range opts {
		o(&ro)
	}
	return ro
}
