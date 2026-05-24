// Package sigterm provides structured error types for shutdown orchestration.
//
// # Error Hierarchy
//
// During a shutdown sequence, individual handlers may fail. These failures are
// collected and surfaced through two error types:
//
//   - [HandlerError] wraps the error returned (or panic recovered) from a single
//     named handler, preserving the handler name for diagnostics.
//
//   - [ShutdownError] aggregates zero or more [HandlerError] values produced
//     across all phases and groups. It implements the multi-error Unwrap []error
//     interface introduced in Go 1.20, so callers can use errors.Is / errors.As
//     to inspect individual causes.
//
// # Usage
//
//	if err := manager.Shutdown(ctx); err != nil {
//		var se *sigterm.ShutdownError
//		if errors.As(err, &se) {
//			for _, e := range se.Errors {
//				log.Println(e)
//			}
//		}
//	}
package sigterm
