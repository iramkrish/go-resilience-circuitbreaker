package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/iramkrish/go-resilience-circuitbreaker/pkg/circuitbreaker"
)

func TestHalfOpen_FailureReopensCircuit(t *testing.T) {
	cb := circuitbreaker.New(
		circuitbreaker.WithConsecutiveFailures(1),
		circuitbreaker.WithTimeout(100*time.Millisecond),
	)

	fail := func() (any, error) { return nil, errors.New("fail") }

	_, _ = cb.Execute(fail) // OPEN

	time.Sleep(150 * time.Millisecond)

	_, _ = cb.Execute(fail) // HALF_OPEN → should go back to OPEN

	if cb.State() != circuitbreaker.Open {
		t.Fatalf("expected OPEN, got %v", cb.State())
	}
}
