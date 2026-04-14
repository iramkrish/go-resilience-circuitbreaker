# ADR 0002: State Management Strategy

## Status

Accepted

---

## Context

The circuit breaker requires safe and predictable state transitions under concurrent access.

Key requirements:

* Ensure correctness of state transitions (CLOSED, OPEN, HALF_OPEN)
* Avoid race conditions under concurrent execution
* Keep implementation simple and maintainable
* Minimize performance overhead

---

## Options Considered

### 1. Lock-Free (Atomic Operations)

**Pros**

* High performance under concurrency
* Avoids mutex contention

**Cons**

* Complex to implement correctly
* Hard to reason about multi-field state transitions
* Error-prone (especially with HALF-OPEN semantics)

---

### 2. Channel-Based State Machine

**Pros**

* Clear event-driven model
* Strong guarantees via message passing

**Cons**

* Adds architectural complexity
* Introduces additional latency
* Overkill for this use case

---

### 3. Mutex-Based Synchronization

**Pros**

* Simple and easy to reason about
* Ensures correctness across multiple fields
* Idiomatic in Go for shared mutable state

**Cons**

* Potential contention under high concurrency
* Slight performance overhead

---

## Decision

Use **mutex-based synchronization** to manage circuit breaker state.

---

## Rationale

### Correctness Over Micro-Optimization

State transitions involve multiple fields:

* current state
* success counters
* half-open request limits
* timestamps

Atomic operations would not safely cover all transitions.

Mutex ensures consistency across all fields.

---

### Simplicity and Maintainability

* Easier to understand and debug
* Lower cognitive load
* Reduced risk of subtle concurrency bugs

---

### Controlled Critical Sections

The design minimizes lock contention by:

* Locking only for state checks and updates
* Executing downstream calls **outside the lock**

This ensures:

* minimal blocking
* good throughput

---

## Execution Model

```id="qk2d8l"
lock → check state → unlock
execute downstream call (no lock)
lock → update state → unlock
```

---

## HALF-OPEN Considerations

HALF-OPEN requires careful synchronization:

* Limit number of probe requests
* Track success count
* Transition immediately on failure

Mutex ensures:

* no race conditions in probe counting
* consistent transition behavior

---

## Consequences

### Positive

* Strong correctness guarantees
* Predictable behavior under concurrency
* Easier debugging and maintenance

---

### Negative

* Potential contention under very high load
* Slight latency overhead per request

---

## When This Becomes a Bottleneck

* Extremely high QPS (>50k requests/sec per instance)
* Ultra low-latency systems (sub-millisecond requirements)

---

## Future Considerations

If contention becomes significant:

### 1. Hybrid Approach

* Combine mutex with atomic counters

### 2. Sharded Circuit Breakers

* Partition load across multiple instances

### 3. Lock-Free Optimizations

* Only if justified by benchmarks

---

## Summary

The chosen approach prioritizes:

* correctness
* simplicity
* maintainability

over premature optimization.

This aligns with the principle:

> “Make it correct first, then make it fast if necessary.”
