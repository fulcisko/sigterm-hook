package sigterm

// WithPhase returns a RegisterOption that assigns the given shutdown phase
// to the handler being registered. Handlers in higher-priority phases are
// shut down before those in lower-priority phases.
//
// Example:
//
//	m.Register("http-server", httpShutdown, WithPhase(PhaseApplication))
//	m.Register("postgres", dbClose, WithPhase(PhaseInfrastructure))
func WithPhase(phase Phase) RegisterOption {
	return func(cfg *registerConfig) {
		cfg.phase = &phase
	}
}

// RegisterOption is a functional option applied when registering a handler.
type RegisterOption func(*registerConfig)

// registerConfig holds per-handler configuration collected from options.
type registerConfig struct {
	phase *Phase
}

// applyRegisterOptions applies a slice of RegisterOption values to a new
// registerConfig and returns the resulting configuration.
func applyRegisterOptions(opts []RegisterOption) registerConfig {
	cfg := registerConfig{}
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
