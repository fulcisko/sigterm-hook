// Package sigterm — lifecycle hooks
//
// # Lifecycle Hooks
//
// Lifecycle hooks let callers observe the high-level progress of a shutdown
// without registering a full shutdown handler.
//
// Four events are emitted:
//
//   - [LifecycleShutdownStarted] — once, before any handler runs.
//   - [LifecyclePhaseStarted]    — before each phase begins.
//   - [LifecyclePhaseDone]       — after each phase finishes.
//   - [LifecycleShutdownDone]    — once, after all phases complete.
//
// Register a hook via [WithLifecycleHook] when constructing the [Manager].
// Hooks are synchronous and must not block; use a goroutine if needed.
package sigterm
