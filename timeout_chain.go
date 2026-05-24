package sigterm

import (
	"context"
	"fmt"
	"time"
)

// timeoutChain resolves the effective deadline for a handler execution by
// composing the handler-level timeout with the global shutdown deadline.
// The shorter of the two wins; if neither is set the parent context is
// returned unchanged.
func timeoutChain(
	parent context.Context,
	handlerTimeout time.Duration,
	globalDeadline time.Time,
) (context.Context, context.CancelFunc, string) {
	now := time.Now()

	var (
		effective time.Time
		label    string
	)

	// Start with the global deadline when it is in the future.
	if !globalDeadline.IsZero() && globalDeadline.After(now) {
		effective = globalDeadline
		label = "global"
	}

	// Prefer the handler timeout when it produces an earlier deadline.
	if handlerTimeout > 0 {
		handlerDeadline := now.Add(handlerTimeout)
		if effective.IsZero() || handlerDeadline.Before(effective) {
			effective = handlerDeadline
			label = "handler"
		}
	}

	// If a parent deadline already exists and is even earlier, respect it.
	if pd, ok := parent.Deadline(); ok {
		if effective.IsZero() || pd.Before(effective) {
			effective = pd
			label = "parent"
		}
	}

	if effective.IsZero() {
		// No bounds — return a no-op cancel so callers can always defer cancel().
		ctx, cancel := context.WithCancel(parent)
		return ctx, cancel, "none"
	}

	ctx, cancel := context.WithDeadline(parent, effective)
	return ctx, cancel, label
}

// deadlineLabel returns a human-readable description of the remaining time
// until the context deadline, or "no deadline" when none is set.
func deadlineLabel(ctx context.Context) string {
	if dl, ok := ctx.Deadline(); ok {
		remaining := time.Until(dl).Truncate(time.Millisecond)
		return fmt.Sprintf("deadline in %s", remaining)
	}
	return "no deadline"
}
