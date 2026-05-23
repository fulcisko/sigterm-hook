package sigterm

import "context"

// HandlerFunc is the signature every shutdown handler must satisfy.
type HandlerFunc func(ctx context.Context) error

// registeredHandler is the internal representation of a handler that has been
// registered with the Manager. It bundles the callable function together with
// its resolved options and dependency metadata.
type registeredHandler struct {
	// name is the unique identifier supplied at registration time.
	name string

	// fn is the shutdown function to invoke.
	fn HandlerFunc

	// deps holds the names of other handlers that must complete before this
	// one is started.
	deps []string

	// opts captures all resolved per-handler options (phase, priority, …).
	opts registerOptions
}

// newRegisteredHandler constructs a registeredHandler from its constituent
// parts after options have been applied.
func newRegisteredHandler(name string, fn HandlerFunc, deps []string, opts registerOptions) *registeredHandler {
	return &registeredHandler{
		name: name,
		fn:   fn,
		deps: deps,
		opts: opts,
	}
}
