package sigterm

import (
	"context"
	"fmt"
	"sync"
)

// GroupResult holds the outcome of running a single handler inside a group.
type GroupResult struct {
	Name string
	Err  error
}

// runGroup executes all handlers in a HandlerGroup concurrently and waits for
// all of them to finish. It returns a slice of GroupResult — one per handler.
// A cancelled context causes in-flight handlers to be interrupted via their
// own context argument.
func runGroup(ctx context.Context, group *HandlerGroup, metrics *shutdownMetrics, log logger) []GroupResult {
	handlers := group.all()
	results := make([]GroupResult, len(handlers))
	var wg sync.WaitGroup

	for i, h := range handlers {
		wg.Add(1)
		go func(idx int, rh *registeredHandler) {
			defer wg.Done()

			name := rh.name
			if metrics != nil {
				metrics.recordHandlerStart(name)
			}

			handlerFn := rh.fn
			if rh.opts.timeout > 0 {
				handlerFn = wrapWithTimeout(handlerFn, rh.opts.timeout)
			}
			if rh.opts.retry != nil {
				handlerFn = wrapWithRetry(handlerFn, *rh.opts.retry)
			}

			err := handlerFn(ctx)

			if metrics != nil {
				metrics.recordHandlerEnd(name, err)
			}
			if err != nil {
				log.Errorf("handler %q failed: %v", name, err)
			} else {
				log.Infof("handler %q completed successfully", name)
			}

			results[idx] = GroupResult{Name: name, Err: err}
		}(i, h)
	}

	wg.Wait()
	return results
}

// collectErrors aggregates non-nil errors from GroupResult entries into a
// single descriptive error, or returns nil when all handlers succeeded.
func collectErrors(results []GroupResult) error {
	var errs []string
	for _, r := range results {
		if r.Err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", r.Name, r.Err))
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("shutdown errors: %v", errs)
}
