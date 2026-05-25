package sigterm

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestShutdownReportSucceeded(t *testing.T) {
	r := &ShutdownReport{}
	if !r.Succeeded() {
		t.Fatal("expected Succeeded() == true for empty error list")
	}
	r.Errors = []error{errors.New("boom")}
	if r.Succeeded() {
		t.Fatal("expected Succeeded() == false when errors present")
	}
}

func TestShutdownReportWriteTo(t *testing.T) {
	r := &ShutdownReport{
		StartedAt:  time.Now().Add(-2 * time.Second),
		FinishedAt: time.Now(),
		Duration:   2 * time.Second,
		Handlers: []HandlerReport{
			{Name: "db", Phase: PhaseNormal, Duration: 100 * time.Millisecond, Retries: 0, Err: nil},
			{Name: "cache", Phase: PhaseNormal, Duration: 50 * time.Millisecond, Retries: 1, Err: errors.New("timeout")},
		},
		Errors: []error{errors.New("timeout")},
	}
	var buf bytes.Buffer
	n, err := r.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo error: %v", err)
	}
	if n == 0 {
		t.Fatal("expected non-zero bytes written")
	}
	out := buf.String()
	if !strings.Contains(out, "FAILED") {
		t.Error("expected FAILED in output")
	}
	if !strings.Contains(out, "db") {
		t.Error("expected handler name 'db' in output")
	}
	if !strings.Contains(out, "timeout") {
		t.Error("expected error message 'timeout' in output")
	}
}

func TestReportBuilderBuild(t *testing.T) {
	rb := newReportBuilder()
	rb.record("svc-a", PhaseNormal, 10*time.Millisecond, 0, nil)
	rb.record("svc-b", PhaseLast, 20*time.Millisecond, 2, errors.New("fail"))

	rep := rb.build()
	if len(rep.Handlers) != 2 {
		t.Fatalf("expected 2 handlers, got %d", len(rep.Handlers))
	}
	if rep.Succeeded() {
		t.Fatal("expected report to be failed")
	}
	if rep.Duration <= 0 {
		t.Fatal("expected positive duration")
	}
}

func TestReportBuilderIsolation(t *testing.T) {
	rb := newReportBuilder()
	rb.record("x", PhaseFirst, 5*time.Millisecond, 0, nil)
	rep := rb.build()
	rep.Handlers[0].Name = "mutated"
	// rebuild should not be affected
	rep2 := rb.build()
	if rep2.Handlers[0].Name == "mutated" {
		t.Fatal("report builder entries should be isolated from returned report")
	}
}
