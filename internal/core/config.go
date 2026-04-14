package core

import "time"

type Metrics interface {
	OnSuccess()
	OnFailure()
	OnStateChange(from, to State)
}

type Strategy interface {
	Record(success bool)
	ShouldOpen() bool
	Reset()
}

type Config struct {
	SuccessThreshold    int
	Timeout             time.Duration
	MaxHalfOpenRequests int
	ErrorFilter         func(error) bool
	Metrics             Metrics
	Strategy            Strategy
}
