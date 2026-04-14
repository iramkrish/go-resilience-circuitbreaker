package unit

import (
	"errors"
	"sync"
	"testing"

	"github.com/iramkrish/go-resilience-circuitbreaker/pkg/circuitbreaker"
)

func TestConcurrencySafety(t *testing.T) {
	cb := circuitbreaker.New(
		circuitbreaker.WithSlidingWindow(20, 0.5, 5),
	)

	var wg sync.WaitGroup

	fn := func() (any, error) {
		return nil, errors.New("fail")
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = cb.Execute(fn)
		}()
	}

	wg.Wait()

	// No panic = success
}
