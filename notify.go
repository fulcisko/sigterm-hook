package sigtermhook

// ShutdownEvent represents a lifecycle event emitted during shutdown.
type ShutdownEvent int

const (
	// EventShutdownStarted is emitted when shutdown begins.
	EventShutdownStarted ShutdownEvent = iota
	// EventHandlerStarted is emitted before a handler runs.
	EventHandlerStarted
	// EventHandlerCompleted is emitted after a handler finishes successfully.
	EventHandlerCompleted
	// EventHandlerFailed is emitted when a handler returns an error.
	EventHandlerFailed
	// EventShutdownCompleted is emitted when all handlers have finished.
	EventShutdownCompleted
)

// ShutdownNotification carries information about a lifecycle event.
type ShutdownNotification struct {
	Event   ShutdownEvent
	Handler string
	Err     error
}

// Notifier receives shutdown lifecycle notifications.
type Notifier interface {
	Notify(n ShutdownNotification)
}

// NotifierFunc is a function that implements Notifier.
type NotifierFunc func(n ShutdownNotification)

// Notify implements Notifier.
func (f NotifierFunc) Notify(n ShutdownNotification) { f(n) }

// notifierBus fans out notifications to multiple Notifiers.
type notifierBus struct {
	notifiers []Notifier
}

func newNotifierBus(notifiers []Notifier) *notifierBus {
	return &notifierBus{notifiers: notifiers}
}

func (b *notifierBus) emit(n ShutdownNotification) {
	for _, notifier := range b.notifiers {
		notifier.Notify(n)
	}
}

func (b *notifierBus) hasNotifiers() bool {
	return len(b.notifiers) > 0
}
