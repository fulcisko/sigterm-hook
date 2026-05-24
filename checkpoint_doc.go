// Package sigterm provides graceful shutdown orchestration.
//
// # Checkpoint Tracking
//
// The checkpoint subsystem records the pass/fail outcome of every shutdown
// handler so callers can inspect the shutdown result after the Manager has
// finished.
//
// Enable it when constructing the manager:
//
//	m := sigterm.NewManager(sigterm.WithCheckpoint())
//
// After Shutdown returns, iterate the results:
//
//	for _, entry := range m.Checkpoints().Entries() {
//		fmt.Printf("%s -> %v (took %s)\n", entry.HandlerName, entry.Status, entry.Duration)
//	}
//
// AllPassed is a convenience helper that returns true only when every
// registered handler completed without error.
package sigterm
