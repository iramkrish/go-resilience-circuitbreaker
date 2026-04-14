package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/iramkrish/go-resilience-circuitbreaker/pkg/circuitbreaker"
)

func TestCircuitBreaker_OpenState(t *testing.T) {
	cb := circuitbreaker.New(
		circuitbreaker.WithConsecutiveFailures(2),
		circuitbreaker.WithTimeout(1*time.Second),
	)

	fail := func() (any, error) {
		return nil, errors.New("fail")
	}

	_, err := cb.Execute(fail)
	if err == nil {
		t.Fatalf("expected error")
	}

	_, err = cb.Execute(fail)
	if err == nil {
		t.Fatalf("expected error")
	}

	if cb.State() != circuitbreaker.Open {
		t.Fatalf("expected OPEN state, got %v", cb.State())
	}
}
