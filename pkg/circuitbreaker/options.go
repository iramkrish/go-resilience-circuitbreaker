package circuitbreaker

import (
	"time"

	"github.com/iramkrish/go-resilience-circuitbreaker/internal/core"
	"github.com/iramkrish/go-resilience-circuitbreaker/internal/strategy"
)

type Option func(*core.Config)

type Metrics = core.Metrics

func defaultConfig() *core.Config {
	return &core.Config{
		SuccessThreshold:    2,
		Timeout:             5 * time.Second,
		MaxHalfOpenRequests: 2,
		ErrorFilter:         func(err error) bool { return true },
		Strategy:            strategy.NewConsecutive(3), // default
	}
}

func WithConsecutiveFailures(n int) Option {
	return func(c *core.Config) {
		c.Strategy = strategy.NewConsecutive(n)
	}
}

func WithSlidingWindow(size int, rate float64, min int) Option {
	return func(c *core.Config) {
		c.Strategy = strategy.NewSlidingWindow(size, rate, min)
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c *core.Config) { c.Timeout = d }
}
