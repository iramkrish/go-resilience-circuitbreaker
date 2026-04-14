package unit

import (
	"testing"

	"github.com/iramkrish/go-resilience-circuitbreaker/pkg/circuitbreaker"
)

func BenchmarkCircuitBreaker(b *testing.B) {
	cb := circuitbreaker.New(
		circuitbreaker.WithSlidingWindow(50, 0.5, 10),
	)

	fn := func() (any, error) {
		return "ok", nil
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cb.Execute(fn)
		}
	})
}
