package sigterm

import (
	"context"
	"testing"
	"time"
)

func TestTimeoutChainNoBounds(t *testing.T) {
	ctx, cancel := timeoutChain(context.Background(), 0, 0)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Fatal("context should not be done without timeout")
	default:
	}
}

func TestTimeoutChainHandlerTimeout(t *testing.T) {
	ctx, cancel := timeoutChain(context.Background(), 20*time.Millisecond, 0)
	defer cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected handler timeout to fire")
	}
}

func TestTimeoutChainGlobalDeadline(t *testing.T) {
	ctx, cancel := timeoutChain(context.Background(), 0, 20*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("expected global deadline to fire")
	}
}

func TestTimeoutChainShorterWins(t *testing.T) {
	start := time.Now()
	ctx, cancel := timeoutChain(context.Background(), 30*time.Millisecond, 500*time.Millisecond)
	defer cancel()

	<-ctx.Done()
	elapsed := time.Since(start)

	if elapsed > 200*time.Millisecond {
		t.Fatalf("expected shorter handler timeout to win, elapsed=%s", elapsed)
	}
}

func TestTimeoutChainParentCancellation(t *testing.T) {
	parent, parentCancel := context.WithCancel(context.Background())
	ctx, cancel := timeoutChain(parent, 500*time.Millisecond, 500*time.Millisecond)
	defer cancel()

	parentCancel()

	select {
	case <-ctx.Done():
		// expected: parent cancellation propagates
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected parent cancellation to propagate")
	}
}

func TestDeadlineLabel(t *testing.T) {
	cases := []struct {
		handler  time.Duration
		global   time.Duration
		expected string
	}{
		{0, 0, "none"},
		{5 * time.Second, 0, "handler(5s)"},
		{0, 10 * time.Second, "global(10s)"},
		{5 * time.Second, 10 * time.Second, "handler(5s) or global(10s)"},
	}

	for _, tc := range cases {
		got := deadlineLabel(tc.handler, tc.global)
		if got != tc.expected {
			t.Errorf("deadlineLabel(%s, %s) = %q, want %q", tc.handler, tc.global, got, tc.expected)
		}
	}
}
