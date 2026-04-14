package main

import (
	"errors"
	"net/http"

	"github.com/iramkrish/go-resilience-circuitbreaker/pkg/circuitbreaker"
)

func main() {
	cb := circuitbreaker.New(
		circuitbreaker.WithSlidingWindow(5, 0.5, 2),
	)

	cb.Execute(func() (any, error) {
		resp, err := http.Get("https://example.com")
		if err != nil {
			return nil, err
		}

		if resp.StatusCode >= 500 {
			return nil, errors.New("server error")
		}

		return resp, nil
	})
}
