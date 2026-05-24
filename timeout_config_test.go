package sigterm

import (
	"testing"
	"time"
)

func TestTimeoutConfigValidateZeroValues(t *testing.T) {
	cfg := TimeoutConfig{}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for zero values, got: %v", err)
	}
}

func TestTimeoutConfigValidatePositiveValues(t *testing.T) {
	cfg := TimeoutConfig{
		GlobalTimeout:  10 * time.Second,
		HandlerTimeout: 5 * time.Second,
		DrainTimeout:   3 * time.Second,
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error for positive values, got: %v", err)
	}
}

func TestTimeoutConfigValidateNegativeGlobal(t *testing.T) {
	cfg := TimeoutConfig{GlobalTimeout: -1 * time.Second}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative GlobalTimeout")
	}
	invalid, ok := err.(ErrInvalidTimeout)
	if !ok {
		t.Fatalf("expected ErrInvalidTimeout, got %T", err)
	}
	if invalid.Field != "GlobalTimeout" {
		t.Errorf("expected field GlobalTimeout, got %s", invalid.Field)
	}
}

func TestTimeoutConfigValidateNegativeHandler(t *testing.T) {
	cfg := TimeoutConfig{HandlerTimeout: -500 * time.Millisecond}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative HandlerTimeout")
	}
}

func TestTimeoutConfigValidateNegativeDrain(t *testing.T) {
	cfg := TimeoutConfig{DrainTimeout: -1 * time.Millisecond}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative DrainTimeout")
	}
}

func TestErrInvalidTimeoutMessage(t *testing.T) {
	err := ErrInvalidTimeout{Field: "HandlerTimeout", Value: -2 * time.Second}
	msg := err.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestMergeTimeoutConfigUsesPerHandler(t *testing.T) {
	global := TimeoutConfig{HandlerTimeout: 10 * time.Second}
	result := mergeTimeoutConfig(global, 3*time.Second)
	if result != 3*time.Second {
		t.Errorf("expected 3s, got %v", result)
	}
}

func TestMergeTimeoutConfigFallsBackToGlobal(t *testing.T) {
	global := TimeoutConfig{HandlerTimeout: 10 * time.Second}
	result := mergeTimeoutConfig(global, 0)
	if result != 10*time.Second {
		t.Errorf("expected 10s, got %v", result)
	}
}

func TestMergeTimeoutConfigBothZero(t *testing.T) {
	global := TimeoutConfig{}
	result := mergeTimeoutConfig(global, 0)
	if result != 0 {
		t.Errorf("expected 0, got %v", result)
	}
}
