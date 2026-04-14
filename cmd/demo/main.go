package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/iramkrish/go-resilience-circuitbreaker/pkg/circuitbreaker"
)

func main() {
	cb := circuitbreaker.New(
		circuitbreaker.WithSlidingWindow(10, 0.6, 5),
		circuitbreaker.WithTimeout(3*time.Second),
	)

	for i := 0; i < 25; i++ {
		res, err := cb.Execute(simulatedCall)

		fmt.Printf("State=%s | Result=%v | Err=%v\n",
			cb.State(), res, err)

		time.Sleep(300 * time.Millisecond)
	}
}

func simulatedCall() (any, error) {
	if rand.Intn(10) < 7 {
		return nil, errors.New("random failure")
	}
	return "success", nil
}
