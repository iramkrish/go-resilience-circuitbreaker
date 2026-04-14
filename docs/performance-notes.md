# Performance Notes

## Complexity

| Component                 | Complexity |
| ------------------------- | ---------- |
| Execute()                 | O(1)       |
| Sliding window evaluation | O(n)       |
| State transition          | O(1)       |

---

## Bottlenecks

1. Mutex contention
2. Window iteration
3. High QPS systems

---

## When This Becomes a Problem

* > 50k QPS per instance
* extremely low latency systems (sub-ms)

---

## Mitigation Strategies

* shard circuit breakers
* use lock-free structures
* move to time-bucketed window

---

## Trade-off

Simplicity vs performance:

Current design favors:

* correctness
* readability
* extensibility

Over:

* micro-optimizations
