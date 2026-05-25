package sigterm

import (
	"context"
	"testing"
	"time"
)

func TestBudgetUnlimited(t *testing.T) {
	b := newShutdownBudget(0)
	if b.Exhausted() {
		t.Fatal("unlimited budget should never be exhausted")
	}
	rem := b.Remaining()
	if rem != 24*time.Hour {
		t.Fatalf("expected sentinel 24h, got %s", rem)
	}
}

func TestBudgetRemaining(t *testing.T) {
	b := newShutdownBudget(5 * time.Second)
	rem := b.Remaining()
	if rem <= 0 || rem > 5*time.Second {
		t.Fatalf("unexpected remaining: %s", rem)
	}
}

func TestBudgetExhausted(t *testing.T) {
	b := newShutdownBudget(1 * time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if !b.Exhausted() {
		t.Fatal("budget should be exhausted after sleep")
	}
}

func TestBudgetSpend(t *testing.T) {
	b := newShutdownBudget(0)
	b.Spend(200 * time.Millisecond)
	b.Spend(300 * time.Millisecond)
	if b.spent != 500*time.Millisecond {
		t.Fatalf("expected 500ms spent, got %s", b.spent)
	}
}

func TestBudgetWithContextUnlimited(t *testing.T) {
	b := newShutdownBudget(0)
	ctx, cancel := b.WithContext(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		t.Fatal("unlimited budget context should not be cancelled")
	default:
	}
}

func TestBudgetWithContextExhausted(t *testing.T) {
	b := newShutdownBudget(1 * time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	ctx, cancel := b.WithContext(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		// expected
	case <-time.After(50 * time.Millisecond):
		t.Fatal("exhausted budget context should be cancelled immediately")
	}
}

func TestBudgetWithContextDeadline(t *testing.T) {
	b := newShutdownBudget(500 * time.Millisecond)
	ctx, cancel := b.WithContext(context.Background())
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("expected a deadline on the context")
	}
	if time.Until(deadline) <= 0 {
		t.Fatal("deadline should be in the future")
	}
}

func TestBudgetString(t *testing.T) {
	b := newShutdownBudget(10 * time.Second)
	s := b.String()
	if s == "" {
		t.Fatal("String() should not be empty")
	}
}

func TestBudgetStringUnlimited(t *testing.T) {
	b := newShutdownBudget(0)
	if b.String() != "budget: unlimited" {
		t.Fatalf("unexpected string: %s", b.String())
	}
}
