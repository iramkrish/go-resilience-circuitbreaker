# Failure Timeline — Sliding Window Behavior

## Scenario

Configuration:

* Window size = 10
* Failure rate threshold = 60%
* Minimum requests = 5

---

## Timeline

```
Request #:   1   2   3   4   5   6   7   8   9   10
Result:      F   F   S   F   S   F   F   S   F   F
```

Failures = 6 / 10 = 60%

---

## State Transition

```
CLOSED
  |
  | (failure rate exceeds threshold)
  v
OPEN
```

---

## Recovery Timeline

```
Time passes (timeout)

OPEN → HALF_OPEN

Probe:
  Request 1 → Success
  Request 2 → Success

→ CLOSED
```

---

## Failure During Recovery

```
HALF_OPEN:
  Request 1 → Success
  Request 2 → Failure

→ immediately OPEN
```

---

## Key Insight

Sliding window captures **patterns**, not spikes.

Bad pattern:

```
F S F F S F F S F F → OPEN
```

Good pattern:

```
S S F S S S S F S S → stays CLOSED
```

---

## Why This Matters

Real systems fail intermittently — not in clean streaks.

Sliding window:

* avoids false positives
* reacts to sustained degradation
