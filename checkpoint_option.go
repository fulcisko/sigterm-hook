package sigterm

// WithCheckpoint enables checkpoint tracking for the manager.
// When enabled, each handler's outcome is recorded in the CheckpointTracker
// accessible via Manager.Checkpoints().
func WithCheckpoint() ManagerOption {
	return func(m *Manager) {
		m.checkpoints = newCheckpointTracker()
	}
}

// ManagerOption is a functional option for configuring a Manager.
// It is already defined in hook.go; this file extends its usage.
// (No redeclaration needed — options are applied via NewManager.)
