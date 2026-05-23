package sigterm

import (
	"testing"
	"time"
)

func TestDefaultRegisterOptions(t *testing.T) {
	ro := defaultRegisterOptions()
	if ro.priority != PriorityNormal {
		t.Errorf("expected PriorityNormal, got %d", ro.priority)
	}
	if ro.phase != PhaseDefault {
		t.Errorf("expected PhaseDefault, got %v", ro.phase)
	}
	if ro.timeout != nil {
		t.Error("expected nil timeout")
	}
	if ro.retry != nil {
		t.Error("expected nil retry")
	}
}

func TestWithPriority(t *testing.T) {
	ro := applyRegisterOptions([]RegisterOption{WithPriority(PriorityCritical)})
	if ro.priority != PriorityCritical {
		t.Errorf("expected PriorityCritical, got %d", ro.priority)
	}
}

func TestWithHandlerPhase(t *testing.T) {
	ro := applyRegisterOptions([]RegisterOption{WithHandlerPhase(PhasePreShutdown)})
	if ro.phase != PhasePreShutdown {
		t.Errorf("expected PhasePreShutdown, got %v", ro.phase)
	}
}

func TestWithHandlerTimeout(t *testing.T) {
	cfg := HandlerTimeoutConfig{Timeout: 5 * time.Second}
	ro := applyRegisterOptions([]RegisterOption{WithHandlerTimeout(cfg)})
	if ro.timeout == nil {
		t.Fatal("expected non-nil timeout")
	}
	if ro.timeout.Timeout != 5*time.Second {
		t.Errorf("expected 5s, got %v", ro.timeout.Timeout)
	}
}

func TestWithHandlerRetry(t *testing.T) {
	cfg := RetryConfig{MaxAttempts: 3}
	ro := applyRegisterOptions([]RegisterOption{WithHandlerRetry(cfg)})
	if ro.retry == nil {
		t.Fatal("expected non-nil retry")
	}
	if ro.retry.MaxAttempts != 3 {
		t.Errorf("expected 3 attempts, got %d", ro.retry.MaxAttempts)
	}
}

func TestApplyRegisterOptionsMultiple(t *testing.T) {
	cfg := HandlerTimeoutConfig{Timeout: 2 * time.Second}
	ro := applyRegisterOptions([]RegisterOption{
		WithPriority(PriorityLow),
		WithHandlerPhase(PhasePostShutdown),
		WithHandlerTimeout(cfg),
	})
	if ro.priority != PriorityLow {
		t.Errorf("expected PriorityLow, got %d", ro.priority)
	}
	if ro.phase != PhasePostShutdown {
		t.Errorf("expected PhasePostShutdown, got %v", ro.phase)
	}
	if ro.timeout == nil || ro.timeout.Timeout != 2*time.Second {
		t.Error("unexpected timeout config")
	}
}
