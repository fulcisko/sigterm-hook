package sigterm

import (
	"time"
)

// reportBuilder accumulates per-handler results during a shutdown run
// and produces a final ShutdownReport.
type reportBuilder struct {
	start    time.Time
	entries  []HandlerReport
	errors   []error
}

func newReportBuilder() *reportBuilder {
	return &reportBuilder{start: time.Now()}
}

// record appends the result of a single handler execution.
func (rb *reportBuilder) record(name string, phase Phase, dur time.Duration, retries int, err error) {
	rb.entries = append(rb.entries, HandlerReport{
		Name:     name,
		Phase:    phase,
		Duration: dur,
		Retries:  retries,
		Err:      err,
	})
	if err != nil {
		rb.errors = append(rb.errors, err)
	}
}

// build finalises and returns the completed ShutdownReport.
func (rb *reportBuilder) build() *ShutdownReport {
	now := time.Now()
	handlers := make([]HandlerReport, len(rb.entries))
	copy(handlers, rb.entries)
	errs := make([]error, len(rb.errors))
	copy(errs, rb.errors)
	return &ShutdownReport{
		StartedAt:  rb.start,
		FinishedAt: now,
		Duration:   now.Sub(rb.start),
		Handlers:   handlers,
		Errors:     errs,
	}
}
