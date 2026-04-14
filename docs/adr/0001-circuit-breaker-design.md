# ADR 0001: Circuit Breaker Design Approach

## Status

Accepted

---

## Context

We need a reusable, production-grade circuit breaker for Go services to prevent cascading failures when interacting with unreliable downstream dependencies.

Key requirements:

* Fail fast under sustained failure
* Support realistic failure patterns (not just consecutive failures)
* Be extensible for future resilience strategies
* Maintain a clean and minimal public API

---

## Options Considered

### 1. Use existing library (e.g., sony/gobreaker)

**Pros**

* Battle-tested
* Minimal implementation effort

**Cons**

* Limited extensibility
* No clear abstraction for failure detection strategies
* Hard to adapt for advanced use cases (sliding window, adaptive logic)

---

### 2. Wrap existing library

**Pros**

* Faster development
* Some customization possible

**Cons**

* Leaks underlying constraints
* Limited control over core behavior
* Adds unnecessary abstraction layer

---

### 3. Build custom implementation

**Pros**

* Full control over:

  * State transitions
  * Failure detection logic
  * Extensibility
* Enables clean architecture
* Better alignment with system-specific requirements

**Cons**

* Increased maintenance cost
* Requires careful testing

---

## Decision

Build a custom circuit breaker implementation with a **pluggable strategy pattern for failure detection**.

---

## Key Design Choice: Strategy Pattern

Introduce a strategy interface:

```go
type Strategy interface {
    Record(success bool)
    ShouldOpen() bool
    Reset()
}
```

---

## Rationale

### Separation of Concerns

* **Circuit breaker** → manages state transitions
* **Strategy** → decides when to open

This avoids coupling failure logic with state management.

---

### Extensibility

Supports multiple failure detection models:

* Consecutive failures (simple baseline)
* Sliding window failure rate (production-ready)
* Future:

  * Time-based windows
  * Latency-based triggers
  * Adaptive circuit breaker

---

### Real-World Alignment

Consecutive failures are insufficient:

```
F S F F S F S F F
```

Sliding window captures **failure patterns**, not just streaks.

---

## HALF-OPEN Design Decision

HALF-OPEN is treated as a **strict probe phase**.

### Rules:

* Limited requests allowed
* Successes move toward recovery
* **Any failure immediately transitions back to OPEN**

### Rationale:

* Prevents unstable dependencies from receiving full traffic
* Ensures fast re-failure detection
* Avoids oscillation between states

---

## Consequences

### Positive

* High flexibility and extensibility
* Clean separation of logic
* Realistic failure detection
* Production-ready behavior

---

### Negative

* Increased code complexity
* Higher maintenance burden vs using a library
* Requires careful testing and validation

---

## Alternatives Rejected

### Use sony/gobreaker directly

Rejected because:

* Limited extensibility
* No strategy abstraction
* Hard to evolve beyond basic use cases

---

## Future Considerations

* Time-based sliding window (rolling buckets)
* Latency-based circuit breaker (p95/p99 thresholds)
* Adaptive circuit breaker (dynamic thresholds)
* Distributed/shared circuit breaker state

---

## Summary

This decision prioritizes:

* Extensibility over simplicity
* Realistic failure handling over naive models
* Clear architecture over quick implementation

The resulting design provides a strong foundation for building advanced resilience patterns in distributed systems.
