package sigterm

import (
	"testing"
)

func TestSortByPriorityEmpty(t *testing.T) {
	names := []string{"a", "b", "c"}
	result := sortByPriority(names, nil)
	if len(result) != len(names) {
		t.Fatalf("expected %d names, got %d", len(names), len(result))
	}
}

func TestSortByPriorityDescending(t *testing.T) {
	names := []string{"low", "high", "normal"}
	priorities := map[string]Priority{
		"low":    PriorityLow,
		"high":   PriorityHigh,
		"normal": PriorityNormal,
	}

	result := sortByPriority(names, priorities)

	expected := []string{"high", "normal", "low"}
	for i, name := range expected {
		if result[i] != name {
			t.Errorf("position %d: expected %q, got %q", i, name, result[i])
		}
	}
}

func TestSortByPriorityStableEqualPriority(t *testing.T) {
	names := []string{"first", "second", "third"}
	priorities := map[string]Priority{
		"first":  PriorityNormal,
		"second": PriorityNormal,
		"third":  PriorityNormal,
	}

	result := sortByPriority(names, priorities)

	for i, name := range names {
		if result[i] != name {
			t.Errorf("stable sort broken at %d: expected %q, got %q", i, name, result[i])
		}
	}
}

func TestSortByPriorityMissingUsesNormal(t *testing.T) {
	names := []string{"known", "unknown"}
	priorities := map[string]Priority{
		"known": PriorityCritical,
	}

	result := sortByPriority(names, priorities)

	if result[0] != "known" {
		t.Errorf("expected 'known' (critical) first, got %q", result[0])
	}
	if result[1] != "unknown" {
		t.Errorf("expected 'unknown' (default normal) second, got %q", result[1])
	}
}

func TestPriorityConstants(t *testing.T) {
	if PriorityLow >= PriorityNormal {
		t.Error("PriorityLow should be less than PriorityNormal")
	}
	if PriorityNormal >= PriorityHigh {
		t.Error("PriorityNormal should be less than PriorityHigh")
	}
	if PriorityHigh >= PriorityCritical {
		t.Error("PriorityHigh should be less than PriorityCritical")
	}
}
