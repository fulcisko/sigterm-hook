package sigterm

import (
	"fmt"
	"runtime/debug"
)

// PanicError wraps a recovered panic value along with the stack trace
// captured at the point of the panic.
type PanicError struct {
	Value interface{}
	Stack []byte
}

func (p *PanicError) Error() string {
	return fmt.Sprintf("handler panicked: %v", p.Value)
}

// wrapWithRecover wraps a ShutdownHandler so that any panic is recovered
// and returned as a *PanicError instead of crashing the process.
//
// This ensures a misbehaving handler cannot abort the entire shutdown
// sequence and allows other handlers to still run.
func wrapWithRecover(name string, h ShutdownHandler) ShutdownHandler {
	return func(ctx ShutdownContext) error {
		return callWithRecover(name, h, ctx)
	}
}

func callWithRecover(name string, h ShutdownHandler, ctx ShutdownContext) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := debug.Stack()
			err = &PanicError{
				Value: r,
				Stack: stack,
			}
		}
	}()
	return h(ctx)
}
