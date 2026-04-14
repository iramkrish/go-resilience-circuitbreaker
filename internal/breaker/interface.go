package breaker

import "github.com/iramkrish/go-resilience-circuitbreaker/internal/core"

type Breaker interface {
	Execute(fn func() (any, error)) (any, error)
	State() core.State
}
