package sigterm

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// ShutdownReport summarises the outcome of a shutdown sequence.
type ShutdownReport struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   time.Duration
	Handlers   []HandlerReport
	Errors     []error
}

// HandlerReport holds the result for a single shutdown handler.
type HandlerReport struct {
	Name     string
	Phase    Phase
	Duration time.Duration
	Retries  int
	Err      error
}

// Succeeded returns true when no handler produced an error.
func (r *ShutdownReport) Succeeded() bool {
	return len(r.Errors) == 0
}

// FailedHandlers returns the subset of handlers that produced an error.
func (r *ShutdownReport) FailedHandlers() []HandlerReport {
	var failed []HandlerReport
	for _, h := range r.Handlers {
		if h.Err != nil {
			failed = append(failed, h)
		}
	}
	return failed
}

// WriteTo writes a human-readable summary to w.
func (r *ShutdownReport) WriteTo(w io.Writer) (int64, error) {
	var sb strings.Builder
	status := "OK"
	if !r.Succeeded() {
		status = fmt.Sprintf("FAILED (%d error(s))", len(r.Errors))
	}
	fmt.Fprintf(&sb, "Shutdown report — %s\n", status)
	fmt.Fprintf(&sb, "  started : %s\n", r.StartedAt.Format(time.RFC3339))
	fmt.Fprintf(&sb, "  finished: %s\n", r.FinishedAt.Format(time.RFC3339))
	fmt.Fprintf(&sb, "  duration: %s\n", r.Duration.Round(time.Millisecond))
	fmt.Fprintf(&sb, "  handlers: %d\n", len(r.Handlers))
	for _, h := range r.Handlers {
		result := "ok"
		if h.Err != nil {
			result = "ERR: " + h.Err.Error()
		}
		fmt.Fprintf(&sb, "    [phase=%d] %-30s %6s retries=%d  %s\n",
			h.Phase, h.Name, h.Duration.Round(time.Millisecond), h.Retries, result)
	}
	if len(r.Errors) > 0 {
		fmt.Fprintf(&sb, "  errors:\n")
		for _, e := range r.Errors {
			fmt.Fprintf(&sb, "    - %s\n", e)
		}
	}
	n, err := io.WriteString(w, sb.String())
	return int64(n), err
}
