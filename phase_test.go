package sigterm

import (
	"testing"
)

func TestPhaseConstants(t *testing.T) {
	if PhasePreShutdown.Priority <= PhaseApplication.Priority {
		t.Errorf("expected PreShutdown priority > Application priority")
	}
	if PhaseApplication.Priority <= PhaseInfrastructure.Priority {
		t.Errorf("expected Application priority > Infrastructure priority")
	}
	if PhaseInfrastructure.Priority <= PhaseCleanup.Priority {
		t.Errorf("expected Infrastructure priority > Cleanup priority")
	}
}

func TestPhaseRegistryDefaultPhase(t *testing.T) {
	r := newPhaseRegistry()
	p := r.Get("unknown-handler")
	if p != PhaseApplication {
		t.Errorf("expected default phase to be PhaseApplication, got %v", p)
	}
}

func TestPhaseRegistryAssignAndGet(t *testing.T) {
	r := newPhaseRegistry()
	r.Assign("db", PhaseInfrastructure)
	p := r.Get("db")
	if p != PhaseInfrastructure {
		t.Errorf("expected PhaseInfrastructure, got %v", p)
	}
}

func TestGroupByPhaseOrdering(t *testing.T) {
	r := newPhaseRegistry()
	r.Assign("cache", PhaseInfrastructure)
	r.Assign("http", PhaseApplication)
	r.Assign("pre", PhasePreShutdown)
	r.Assign("cleanup", PhaseCleanup)

	handlers := []string{"cache", "http", "pre", "cleanup"}
	buckets := r.GroupByPhase(handlers)

	if len(buckets) != 4 {
		t.Fatalf("expected 4 phase buckets, got %d", len(buckets))
	}

	// First bucket must contain the highest-priority phase handler.
	if buckets[0][0] != "pre" {
		t.Errorf("expected first bucket to contain 'pre', got %v", buckets[0])
	}
	// Last bucket must contain the lowest-priority phase handler.
	if buckets[len(buckets)-1][0] != "cleanup" {
		t.Errorf("expected last bucket to contain 'cleanup', got %v", buckets[len(buckets)-1])
	}
}

func TestGroupByPhaseSamePhaseSameBucket(t *testing.T) {
	r := newPhaseRegistry()
	r.Assign("svc-a", PhaseApplication)
	r.Assign("svc-b", PhaseApplication)

	buckets := r.GroupByPhase([]string{"svc-a", "svc-b"})
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket for same phase, got %d", len(buckets))
	}
	if len(buckets[0]) != 2 {
		t.Errorf("expected 2 handlers in bucket, got %d", len(buckets[0]))
	}
}

func TestGroupByPhaseEmpty(t *testing.T) {
	r := newPhaseRegistry()
	buckets := r.GroupByPhase([]string{})
	if len(buckets) != 0 {
		t.Errorf("expected 0 buckets for empty input, got %d", len(buckets))
	}
}
