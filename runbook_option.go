package sigterm

import "time"

// WithRunbookEntry attaches a RunbookEntry to a handler at registration time.
// It records the handler's human-readable description, estimated shutdown
// duration, and any dependency names so that operators can generate a
// runbook for the service's shutdown procedure.
//
// Example:
//
//	manager.Register("db", dbHandler,
//		sigterm.WithRunbookEntry(sigterm.RunbookEntry{
//			Description:       "Closes all active database connections.",
//			EstimatedDuration: 2 * time.Second,
//		}),
//	)
func WithRunbookEntry(entry RunbookEntry) RegisterOption {
	return registerOptionFunc(func(o *registerOptions) {
		o.runbookEntry = &entry
	})
}

// WithEstimatedDuration is a convenience option that sets only the estimated
// shutdown duration for the runbook without requiring a full RunbookEntry.
func WithEstimatedDuration(d time.Duration) RegisterOption {
	return registerOptionFunc(func(o *registerOptions) {
		if o.runbookEntry == nil {
			o.runbookEntry = &RunbookEntry{}
		}
		o.runbookEntry.EstimatedDuration = d
	})
}
