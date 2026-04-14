# go-resilience-circuitbreaker

A production-grade circuit breaker implementation for Go, designed to prevent cascading failures in distributed systems.

---

## Features

- Circuit breaker state machine (Closed / Open / Half-Open)
- Pluggable failure detection strategies:
  - Consecutive failures
  - Sliding window failure rate
- Strict HALF-OPEN semantics (fail-fast on probe failure)
- Configurable thresholds and recovery behavior
- Concurrency-safe
- Extensible architecture (strategy pattern)

---

## Installation

```bash
go get github.com/iramkrish/go-resilience-circuitbreaker
