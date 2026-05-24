package sigterm

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestWrapWithRecoverNoPanic(t *testing.T) {
	expected := errors.New("normal error")
	h := wrapWithRecover("test", func(_ ShutdownContext) error {
		return expected
	})

	ctx := newShutdownContext(context.Background(), "test")
	got := h(ctx)
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestWrapWithRecoverPanicString(t *testing.T) {
	h := wrapWithRecover("test", func(_ ShutdownContext) error {
		panic("something went wrong")
	})

	ctx := newShutdownContext(context.Background(), "test")
	err := h(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var pe *PanicError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PanicError, got %T", err)
	}
	if pe.Value != "something went wrong" {
		t.Fatalf("unexpected panic value: %v", pe.Value)
	}
	if len(pe.Stack) == 0 {
		t.Fatal("expected non-empty stack trace")
	}
}

func TestWrapWithRecoverPanicNonString(t *testing.T) {
	h := wrapWithRecover("test", func(_ ShutdownContext) error {
		panic(42)
	})

	ctx := newShutdownContext(context.Background(), "test")
	err := h(ctx)

	var pe *PanicError
	if !errors.As(err, &pe) {
		t.Fatalf("expected *PanicError, got %T", err)
	}
	if pe.Value != 42 {
		t.Fatalf("expected 42, got %v", pe.Value)
	}
}

func TestPanicErrorMessage(t *testing.T) {
	pe := &PanicError{Value: "oops", Stack: []byte("stack trace")}
	msg := pe.Error()
	if !strings.Contains(msg, "oops") {
		t.Fatalf("error message missing panic value: %q", msg)
	}
}

func TestWrapWithRecoverSuccess(t *testing.T) {
	h := wrapWithRecover("test", func(_ ShutdownContext) error {
		return nil
	})

	ctx := newShutdownContext(context.Background(), "test")
	if err := h(ctx); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
