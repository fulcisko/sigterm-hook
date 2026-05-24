package sigterm

// WithLifecycleHook appends a LifecycleHook to the Manager.
// Multiple calls accumulate hooks; they are invoked in registration order.
//
// Example:
//
//	mgr := sigterm.NewManager(
//		sigterm.WithLifecycleHook(func(e sigterm.LifecycleEvent, p sigterm.Phase) {
//			log.Printf("lifecycle event=%d phase=%d", e, p)
//		}),
//	)
func WithLifecycleHook(h LifecycleHook) Option {
	return optionFunc(func(m *Manager) {
		m.lifecycle.hooks = append(m.lifecycle.hooks, h)
	})
}
