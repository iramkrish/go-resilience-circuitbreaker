package circuitbreaker

import "github.com/iramkrish/go-resilience-circuitbreaker/internal/core"

type State = core.State

const (
	Closed   = core.Closed
	Open     = core.Open
	HalfOpen = core.HalfOpen
)
