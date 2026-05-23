// Package sigterm provides graceful shutdown orchestration with
// dependency-aware teardown ordering.
//
// # Handler Registration Options
//
// When registering a shutdown handler via [Manager.Register], callers can
// supply zero or more [RegisterOption] values to customise execution
// behaviour:
//
//	// Run this handler early, in the pre-shutdown phase, with a 3-second
//	// per-handler timeout and up to 2 retry attempts on failure.
//	 manager.Register("db", dbHandler,
//	     WithPriority(sigterm.PriorityCritical),
//	     WithHandlerPhase(sigterm.PhasePreShutdown),
//	     WithHandlerTimeout(sigterm.HandlerTimeoutConfig{Timeout: 3 * time.Second}),
//	     WithHandlerRetry(sigterm.RetryConfig{MaxAttempts: 2}),
//	 )
//
// Options are applied in the order they are provided; later options override
// earlier ones for the same field.
//
// # Priority
//
// [WithPriority] accepts any integer.  The package exposes three named
// constants for convenience:
//
//	- [PriorityCritical] — run first
//	- [PriorityNormal]   — default
//	- [PriorityLow]      — run last
//
// # Phase
//
// [WithHandlerPhase] places the handler into a named lifecycle bucket.
// Phases are processed in the order defined by [PhaseOrder].
package sigterm
