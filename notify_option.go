package sigtermhook

// WithNotifier registers a Notifier that receives lifecycle events during
// shutdown. Multiple notifiers may be added; they are called in registration
// order for every event.
func WithNotifier(n Notifier) ManagerOption {
	return managerOptionFunc(func(cfg *managerConfig) {
		cfg.notifiers = append(cfg.notifiers, n)
	})
}

// WithNotifierFunc is a convenience wrapper around WithNotifier that accepts a
// plain function instead of a Notifier implementation.
func WithNotifierFunc(fn func(ShutdownNotification)) ManagerOption {
	return WithNotifier(NotifierFunc(fn))
}
