// Package sigterm provides graceful shutdown orchestration with
// dependency-aware teardown ordering.
//
// # Shutdown Context
//
// When a shutdown is initiated, the library enriches the context passed to
// every registered handler with a [ShutdownContext] value.  Handlers can
// retrieve it via [ShutdownContextFrom] to inspect why the shutdown was
// triggered and which handlers have already completed.
//
// Example:
//
//	func myHandler(ctx context.Context) error {
//		if sc := sigterm.ShutdownContextFrom(ctx); sc != nil {
//			log.Printf("shutdown reason: %s", sc.Reason())
//			log.Printf("already done: %v", sc.CompletedHandlers())
//		}
//		// … perform cleanup …
//		return nil
//	}
package sigterm
