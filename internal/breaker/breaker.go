package breaker

import (
	"sync"
	"time"

	"github.com/iramkrish/go-resilience-circuitbreaker/internal/core"
)

type breaker struct {
	mu sync.Mutex

	state core.State

	successes int

	lastFailureTime  time.Time
	halfOpenRequests int

	cfg *core.Config
}

func New(cfg *core.Config) Breaker {
	return &breaker{
		state: core.Closed,
		cfg:   cfg,
	}
}

func (b *breaker) State() core.State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

func (b *breaker) Execute(fn func() (any, error)) (any, error) {
	b.mu.Lock()

	switch b.state {
	case core.Open:
		if time.Since(b.lastFailureTime) > b.cfg.Timeout {
			b.transition(core.HalfOpen)
			b.halfOpenRequests = 0
			b.cfg.Strategy.Reset()
		} else {
			b.mu.Unlock()
			return nil, core.ErrCircuitOpen
		}

	case core.HalfOpen:
		if b.halfOpenRequests >= b.cfg.MaxHalfOpenRequests {
			b.mu.Unlock()
			return nil, core.ErrCircuitOpen
		}
		b.halfOpenRequests++
	}

	b.mu.Unlock()

	res, err := fn()

	success := err == nil || !b.cfg.ErrorFilter(err)

	b.mu.Lock()
	defer b.mu.Unlock()

	b.cfg.Strategy.Record(success)

	if !success {
		b.lastFailureTime = time.Now()

		if b.cfg.Metrics != nil {
			b.cfg.Metrics.OnFailure()
		}

		if b.state == core.HalfOpen {
			b.transition(core.Open)
			return nil, err
		}

		if b.cfg.Strategy.ShouldOpen() {
			b.transition(core.Open)
		}

		return nil, err
	}

	if b.cfg.Metrics != nil {
		b.cfg.Metrics.OnSuccess()
	}

	if b.state == core.HalfOpen {
		b.successes++
		if b.successes >= b.cfg.SuccessThreshold {
			b.successes = 0
			b.cfg.Strategy.Reset()
			b.transition(core.Closed)
		}
		return res, nil
	}

	return res, nil
}

func (b *breaker) transition(to core.State) {
	from := b.state
	b.state = to

	if b.cfg.Metrics != nil {
		b.cfg.Metrics.OnStateChange(from, to)
	}
}
