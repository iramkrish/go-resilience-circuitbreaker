package circuitbreaker

import "github.com/iramkrish/go-resilience-circuitbreaker/internal/breaker"

type CircuitBreaker interface {
	Execute(fn func() (any, error)) (any, error)
	State() State
}

type cbAdapter struct {
	inner breaker.Breaker
}

func (c *cbAdapter) Execute(fn func() (any, error)) (any, error) {
	return c.inner.Execute(fn)
}

func (c *cbAdapter) State() State {
	return c.inner.State()
}

func New(opts ...Option) CircuitBreaker {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return &cbAdapter{inner: breaker.New(cfg)}
}
