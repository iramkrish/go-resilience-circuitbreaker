package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/iramkrish/go-resilience-circuitbreaker/pkg/circuitbreaker"
)

func TestHalfOpen_ClosesOnSuccess(t *testing.T) {
	cb := circuitbreaker.New(
		circuitbreaker.WithConsecutiveFailures(1),
		circuitbreaker.WithTimeout(100*time.Millisecond),
	)

	fail := func() (any, error) { return nil, errors.New("fail") }
	success := func() (any, error) { return "ok", nil }

	_, _ = cb.Execute(fail) // OPEN

	time.Sleep(150 * time.Millisecond)

	_, _ = cb.Execute(success)
	_, _ = cb.Execute(success)

	if cb.State() != circuitbreaker.Closed {
		t.Fatalf("expected CLOSED, got %v", cb.State())
	}
}
