package sigterm

import (
	"context"
	"fmt"
	"time"
)

// timeoutChain applies both a per-handler timeout and the global deadline,
// returning whichever context expires first along with a combined cancel.
func timeoutChain(parent context.Context, handlerTimeout, globalDeadline time.Duration) (context.Context, context.CancelFunc) {
	if handlerTimeout <= 0 && globalDeadline <= 0 {
		return context.WithCancel(parent)
	}

	ctx := parent
	cancels := make([]context.CancelFunc, 0, 2)

	if globalDeadline > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, globalDeadline)
		cancels = append(cancels, cancel)
	}

	if handlerTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, handlerTimeout)
		cancels = append(cancels, cancel)
	}

	combined := func() {
		for _, c := range cancels {
			c()
		}
	}

	return ctx, combined
}

// deadlineLabel returns a human-readable label for which deadline was hit.
func deadlineLabel(handlerTimeout, globalDeadline time.Duration) string {
	if handlerTimeout > 0 && globalDeadline > 0 {
		return fmt.Sprintf("handler(%s) or global(%s)", handlerTimeout, globalDeadline)
	}
	if handlerTimeout > 0 {
		return fmt.Sprintf("handler(%s)", handlerTimeout)
	}
	if globalDeadline > 0 {
		return fmt.Sprintf("global(%s)", globalDeadline)
	}
	return "none"
}
