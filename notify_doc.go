// Package sigtermhook provides lifecycle notifications via the Notifier
// interface, allowing callers to observe shutdown events without coupling to
// the shutdown logic itself.
//
// # Lifecycle events
//
// Events are emitted in the following order during a normal shutdown:
//
//  1. EventShutdownStarted  — shutdown has been initiated
//  2. EventHandlerStarted   — a handler is about to run (one per handler)
//  3. EventHandlerCompleted — handler finished without error
//     OR EventHandlerFailed — handler returned an error (includes the error)
//  4. EventShutdownCompleted — all handlers have finished
//
// # Registering notifiers
//
// Pass one or more notifiers when constructing the Manager:
//
//	m := sigtermhook.NewManager(
//	    sigtermhook.WithNotifierFunc(func(n sigtermhook.ShutdownNotification) {
//	        log.Printf("event=%d handler=%s err=%v", n.Event, n.Handler, n.Err)
//	    }),
//	)
//
// Multiple notifiers are called in registration order for every event.
package sigtermhook
