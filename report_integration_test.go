package sigterm_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	sigterm "github.com/example/sigterm-hook"
)

func TestReportBuilderIntegration(t *testing.T) {
	rb := sigterm.NewReportBuilder()

	rb.Record("alpha", sigterm.PhaseNormal, 15*time.Millisecond, 0, nil)
	rb.Record("beta", sigterm.PhaseNormal, 30*time.Millisecond, 1, errors.New("oops"))

	rep := rb.Build()

	if rep.Succeeded() {
		t.Fatal("expected report to indicate failure")
	}
	if len(rep.Handlers) != 2 {
		t.Fatalf("expected 2 handler entries, got %d", len(rep.Handlers))
	}
}

func TestReportWriteToIntegration(t *testing.T) {
	_ = context.Background() // keep import tidy in broader suite

	rb := sigterm.NewReportBuilder()
	rb.Record("cleanup", sigterm.PhaseLast, 5*time.Millisecond, 0, nil)
	rep := rb.Build()

	var buf bytes.Buffer
	_, err := rep.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo failed: %v", err)
	}
	if !strings.Contains(buf.String(), "OK") {
		t.Error("expected OK status in report output")
	}
	if !strings.Contains(buf.String(), "cleanup") {
		t.Error("expected handler name in report output")
	}
}
