# ADR 0003: Error Classification Strategy

## Status

Accepted

---

## Context

Not all errors should be treated equally by the circuit breaker.

Examples:

* Transient system failures (timeouts, 5xx) → should count as failures
* Client errors (4xx) → should NOT trip the circuit
* Business logic errors → often should not affect circuit state

Incorrect classification can lead to:

* False positives (circuit opens unnecessarily)
* False negatives (circuit stays closed during real failures)

---

## Problem

If all errors are treated as failures:

* Circuit may open due to user input errors
* System becomes overly sensitive
* Reduced availability

If too many errors are ignored:

* Circuit may fail to protect system
* Cascading failures can occur

---

## Options Considered

### 1. Treat All Errors as Failures

**Pros**

* Simple implementation

**Cons**

* Incorrect for real-world systems
* Causes unnecessary circuit trips

---

### 2. Hardcode Error Types

**Pros**

* More control than naive approach

**Cons**

* Not flexible
* Domain-specific logic leaks into library
* Difficult to maintain

---

### 3. Configurable Error Filter (Chosen)

Allow users to define:

```go
func(error) bool
```

Return:

* `true` → count as failure
* `false` → ignore

---

## Decision

Introduce a configurable **error filter function**:

```go
ErrorFilter func(error) bool
```

This function determines whether an error should contribute to circuit breaker failure.

---

## Rationale

### Flexibility

Different systems have different failure semantics:

* HTTP services → 5xx vs 4xx
* gRPC → status codes
* DB → connection vs query errors

The circuit breaker should not enforce domain-specific rules.

---

### Separation of Concerns

* Circuit breaker → handles resilience logic
* Caller → defines what constitutes a failure

---

### Real-World Alignment

Examples:

#### HTTP

```go
if resp.StatusCode >= 500 {
    return true // failure
}
return false
```

---

#### gRPC

```go
if status.Code(err) == codes.Unavailable {
    return true
}
return false
```

---

## Behavior in System

```id="q3w7dz"
error occurs
   ↓
ErrorFilter(err)
   ↓
true  → record failure → strategy evaluation
false → treated as success (ignored)
```

---

## Consequences

### Positive

* Highly flexible
* Works across domains
* Prevents false circuit trips

---

### Negative

* Slight increase in API complexity
* Incorrect configuration can lead to misuse

---

## Risks

* Misconfigured filters may:

  * Ignore real failures
  * Count non-critical errors as failures

---

## Mitigation

* Provide sensible defaults (treat all errors as failures)
* Document best practices

---

## Future Considerations

* Error categorization helpers (HTTP, gRPC, DB)
* Built-in common filters
* Metrics for classified vs ignored errors

---

## Summary

This decision ensures that:

* The circuit breaker remains generic
* Failure semantics are controlled by the caller
* The system behaves correctly in real-world scenarios

It reinforces the principle:

> “Resilience logic should be generic; failure semantics should be domain-specific.”
