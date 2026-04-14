# Circuit Breaker Design

## Overview

This library implements a fault-tolerant circuit breaker for Go services, designed to prevent cascading failures in distributed systems.

It acts as a control layer between the caller and downstream dependencies, ensuring system stability under partial failures.

---

## Goals

* Fail fast under sustained failure
* Protect upstream resources (goroutines, threads)
* Provide configurable recovery strategies
* Enable extensibility for future resilience patterns

---

## Non-Goals

* Not a retry mechanism
* Not a rate limiter
* Not a full resilience framework

---

## High-Level Architecture

```
Client → Circuit Breaker → Dependency
```

---

## Layered Design

* **pkg/** → Public API (stable surface)
* **internal/core/** → Shared types and contracts
* **internal/strategy/** → Failure detection logic
* **internal/breaker/** → State machine and orchestration

---

## State Machine

The circuit breaker operates as a finite state machine with three states:

* **CLOSED** → Normal operation
* **OPEN** → Fail-fast mode (no downstream calls)
* **HALF_OPEN** → Recovery probe mode

---

## State Transitions

```
CLOSED → OPEN       (failure condition met)
OPEN → HALF_OPEN    (timeout elapsed)
HALF_OPEN → CLOSED  (success threshold met)
HALF_OPEN → OPEN    (any failure)
```

---

## Execution Flow

1. Check current state
2. If **OPEN** → reject request immediately
3. If **HALF_OPEN** → allow limited probe requests
4. Execute downstream function
5. Record result via strategy
6. Apply transition rules:

   * CLOSED → OPEN (strategy triggers)
   * HALF_OPEN → OPEN (any failure)
   * HALF_OPEN → CLOSED (success threshold met)

---

## Failure Detection Strategy

Failure detection is abstracted via a pluggable strategy interface:

```go
type Strategy interface {
    Record(success bool)
    ShouldOpen() bool
    Reset()
}
```

This separates:

* **Decision logic** (when to open)
* **Execution logic** (state transitions)

---

## Implemented Strategies

### 1. Consecutive Failures

* Trips on back-to-back failures
* Simple and fast
* Less effective for intermittent failures

---

### 2. Sliding Window (Primary Strategy)

* Tracks last N requests
* Calculates failure rate
* Opens circuit when threshold is exceeded

Example configuration:

* Window size: 10
* Failure rate threshold: 60%
* Minimum requests: 5

---

## Why Sliding Window?

Real-world failures are rarely consecutive:

```
F S F F S F S F F
```

Sliding window:

* Detects sustained degradation
* Avoids false positives
* Provides stable behavior

---

## HALF-OPEN Semantics

HALF_OPEN is a **controlled recovery phase**, not a retry phase.

### Rules

* Only a limited number of requests are allowed
* Successes contribute toward recovery
* **Any failure immediately transitions back to OPEN**

---

## Recovery Flow

```
OPEN → (timeout) → HALF_OPEN
HALF_OPEN → (success threshold met) → CLOSED
HALF_OPEN → (any failure) → OPEN
```

---

## Concurrency Model

* Mutex-based synchronization
* Critical sections are minimized
* Downstream execution occurs outside locks

---

## Performance Characteristics

| Component                 | Complexity |
| ------------------------- | ---------- |
| Execute()                 | O(1)       |
| Strategy (Sliding Window) | O(n)       |
| State transitions         | O(1)       |

---

## Trade-offs

### Simplicity vs Accuracy

* Sliding window improves accuracy over consecutive failures
* Introduces minor computational overhead

---

### Locking Strategy

* Mutex ensures correctness
* Potential contention under high concurrency

---

## Failure Modes

* Misconfigured thresholds (too aggressive or too lenient)
* Incorrect error classification
* High concurrency leading to contention

---

## Extensibility

New strategies can be added without modifying the breaker:

* Time-based sliding window
* Latency-based circuit breaker (p95/p99)
* Adaptive circuit breaker (dynamic thresholds)

---

## Observability

Hooks supported:

* OnSuccess
* OnFailure
* OnStateChange

Metrics systems (Prometheus, etc.) can be integrated without coupling.

---

## Design Principles

* Separation of concerns (strategy vs orchestration)
* Minimal public API surface
* Dependency inversion via interfaces
* Extensibility without core modification

---

## Summary

This circuit breaker implementation provides:

* Robust failure detection
* Controlled recovery behavior
* Extensible architecture
* Production-ready concurrency model

It balances correctness, simplicity, and flexibility, making it suitable for real-world distributed systems.
