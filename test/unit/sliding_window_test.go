package unit

import (
	"errors"
	"testing"

	"github.com/iramkrish/go-resilience-circuitbreaker/pkg/circuitbreaker"
)

func TestSlidingWindow_OpensOnFailureRate(t *testing.T) {
	cb := circuitbreaker.New(
		circuitbreaker.WithSlidingWindow(5, 0.6, 3),
	)

	fail := func() (any, error) { return nil, errors.New("fail") }
	success := func() (any, error) { return "ok", nil }

	// Pattern: F F S F F → 4/5 failures = 80%
	_, _ = cb.Execute(fail)
	_, _ = cb.Execute(fail)
	_, _ = cb.Execute(success)
	_, _ = cb.Execute(fail)
	_, _ = cb.Execute(fail)

	if cb.State() != circuitbreaker.Open {
		t.Fatalf("expected OPEN, got %v", cb.State())
	}
}
