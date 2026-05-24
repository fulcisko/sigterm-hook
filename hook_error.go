package sigterm

import (
	"errors"
	"fmt"
	"strings"
)

// ShutdownError aggregates one or more handler errors that occurred during shutdown.
type ShutdownError struct {
	Errors []error
}

// Error implements the error interface, returning a combined message.
func (e *ShutdownError) Error() string {
	if len(e.Errors) == 0 {
		return "shutdown completed with no errors"
	}
	msgs := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		msgs = append(msgs, err.Error())
	}
	return fmt.Sprintf("shutdown encountered %d error(s): %s", len(e.Errors), strings.Join(msgs, "; "))
}

// Unwrap returns the list of wrapped errors for use with errors.Is / errors.As.
func (e *ShutdownError) Unwrap() []error {
	return e.Errors
}

// Is reports whether any error in the aggregated list matches target.
func (e *ShutdownError) Is(target error) bool {
	for _, err := range e.Errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// HandlerError wraps an error returned by a named shutdown handler.
type HandlerError struct {
	Handler string
	Cause   error
}

// Error implements the error interface.
func (e *HandlerError) Error() string {
	return fmt.Sprintf("handler %q failed: %v", e.Handler, e.Cause)
}

// Unwrap returns the underlying cause.
func (e *HandlerError) Unwrap() error {
	return e.Cause
}

// newHandlerError constructs a HandlerError.
func newHandlerError(name string, cause error) *HandlerError {
	return &HandlerError{Handler: name, Cause: cause}
}

// newShutdownError builds a ShutdownError from a slice of errors, returning nil
// when the slice is empty.
func newShutdownError(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return &ShutdownError{Errors: errs}
}
