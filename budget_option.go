package sigterm

import "time"

// WithShutdownBudget sets a hard wall-clock limit for the entire shutdown
// sequence. Once the budget is exhausted every in-flight handler context is
// cancelled and no new handlers are started.
//
// Example:
//
//	m := sigterm.NewManager(
//		sigterm.WithShutdownBudget(10 * time.Second),
//	)
func WithShutdownBudget(d time.Duration) ManagerOption {
	return managerOptionFunc(func(cfg *managerConfig) {
		cfg.shutdownBudget = d
	})
}
