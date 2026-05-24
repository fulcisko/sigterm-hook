package sigterm

import (
	"errors"
	"strings"
	"testing"
)

var errSentinel = errors.New("sentinel")

func TestHandlerErrorMessage(t *testing.T) {
	he := newHandlerError("db", errSentinel)
	if !strings.Contains(he.Error(), "db") {
		t.Errorf("expected handler name in error, got: %s", he.Error())
	}
	if !strings.Contains(he.Error(), errSentinel.Error()) {
		t.Errorf("expected cause in error, got: %s", he.Error())
	}
}

func TestHandlerErrorUnwrap(t *testing.T) {
	he := newHandlerError("cache", errSentinel)
	if !errors.Is(he, errSentinel) {
		t.Error("errors.Is should find sentinel through HandlerError")
	}
}

func TestShutdownErrorNilWhenEmpty(t *testing.T) {
	if err := newShutdownError(nil); err != nil {
		t.Errorf("expected nil for empty slice, got %v", err)
	}
	if err := newShutdownError([]error{}); err != nil {
		t.Errorf("expected nil for empty slice, got %v", err)
	}
}

func TestShutdownErrorMessage(t *testing.T) {
	errs := []error{
		newHandlerError("db", errors.New("conn refused")),
		newHandlerError("cache", errors.New("timeout")),
	}
	se := newShutdownError(errs).(*ShutdownError)
	if !strings.Contains(se.Error(), "2 error(s)") {
		t.Errorf("unexpected message: %s", se.Error())
	}
}

func TestShutdownErrorIs(t *testing.T) {
	errs := []error{
		newHandlerError("db", errSentinel),
	}
	se := newShutdownError(errs)
	if !errors.Is(se, errSentinel) {
		t.Error("errors.Is should find sentinel through ShutdownError")
	}
}

func TestShutdownErrorUnwrapMultiple(t *testing.T) {
	var (
		err1 = errors.New("e1")
		err2 = errors.New("e2")
	)
	se := &ShutdownError{Errors: []error{err1, err2}}
	unwrapped := se.Unwrap()
	if len(unwrapped) != 2 {
		t.Fatalf("expected 2 unwrapped errors, got %d", len(unwrapped))
	}
}

func TestShutdownErrorNoErrors(t *testing.T) {
	se := &ShutdownError{}
	if !strings.Contains(se.Error(), "no errors") {
		t.Errorf("unexpected message for empty ShutdownError: %s", se.Error())
	}
}
