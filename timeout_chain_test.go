package sigterm

import (
	"context"
	"testing"
	"time"
)

func TestTimeoutChainNoBounds(t *testing.T) {
	ctx, cancel, label := timeoutChain(context.Background(), 0, time.Time{})
	defer cancel()

	if label != "none" {
		t.Fatalf("expected label 'none', got %q", label)
	}
	if _, ok := ctx.Deadline(); ok {
		t.Fatal("expected no deadline")
	}
}

func TestTimeoutChainHandlerTimeout(t *testing.T) {
	handlerTimeout := 500 * time.Millisecond
	ctx, cancel, label := timeoutChain(context.Background(), handlerTimeout, time.Time{})
	defer cancel()

	if label != "handler" {
		t.Fatalf("expected label 'handler', got %q", label)
	}
	dl, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline")
	}
	if time.Until(dl) > handlerTimeout {
		t.Fatal("deadline exceeds handler timeout")
	}
}

func TestTimeoutChainGlobalDeadline(t *testing.T) {
	globalDeadline := time.Now().Add(2 * time.Second)
	ctx, cancel, label := timeoutChain(context.Background(), 0, globalDeadline)
	defer cancel()

	if label != "global" {
		t.Fatalf("expected label 'global', got %q", label)
	}
	dl, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline")
	}
	if !dl.Equal(globalDeadline) {
		t.Fatalf("expected deadline %v, got %v", globalDeadline, dl)
	}
}

func TestTimeoutChainShorterWins(t *testing.T) {
	globalDeadline := time.Now().Add(5 * time.Second)
	handlerTimeout := 1 * time.Second

	ctx, cancel, label := timeoutChain(context.Background(), handlerTimeout, globalDeadline)
	defer cancel()

	if label != "handler" {
		t.Fatalf("expected label 'handler', got %q", label)
	}
	dl, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline")
	}
	if time.Until(dl) > handlerTimeout+50*time.Millisecond {
		t.Fatal("expected handler deadline to win")
	}
}

func TestTimeoutChainParentCancellation(t *testing.T) {
	parent, parentCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer parentCancel()

	// Handler timeout is longer — parent should win.
	ctx, cancel, label := timeoutChain(parent, 10*time.Second, time.Time{})
	defer cancel()

	if label != "parent" {
		t.Fatalf("expected label 'parent', got %q", label)
	}

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context should have been cancelled by parent timeout")
	}
}

func TestDeadlineLabelNoDeadline(t *testing.T) {
	label := deadlineLabel(context.Background())
	if label != "no deadline" {
		t.Fatalf("expected 'no deadline', got %q", label)
	}
}

func TestDeadlineLabelWithDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	label := deadlineLabel(ctx)
	if label == "no deadline" {
		t.Fatal("expected a deadline label")
	}
}
