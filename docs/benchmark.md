# Benchmarking & Performance

## Goal

Measure:

* Throughput
* Latency overhead
* Contention under concurrency

---

## Benchmark Setup

### Command

```
go test -bench=. -benchmem ./...
```

---

## Example Benchmark

```go
func BenchmarkCircuitBreaker(b *testing.B) {
	cb := circuitbreaker.New(
		circuitbreaker.WithSlidingWindow(50, 0.5, 10),
	)

	fn := func() (any, error) {
		return "ok", nil
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cb.Execute(fn)
		}
	})
}
```

---

## What to Measure

### 1. ns/op

* Execution overhead
* Should be minimal (few hundred ns)

### 2. allocs/op

* Ideally near zero
* Avoid unnecessary allocations

### 3. contention

* Mutex lock impact under load

---

## Expected Characteristics

| Metric           | Expected                         |
| ---------------- | -------------------------------- |
| Latency overhead | Low                              |
| Allocations      | Minimal                          |
| Throughput       | High                             |
| Contention       | Moderate under heavy concurrency |

---

## Optimization Notes

### Current Design

* Mutex-based synchronization
* O(n) window scan (sliding window)

---

## Potential Improvements

### 1. O(1) Window

* Maintain rolling failure count

### 2. Lock Optimization

* Use atomic counters (advanced)

### 3. Sharded Breakers

* Reduce contention in high-QPS systems

---

## Real-World Considerations

* Circuit breaker overhead must be negligible vs network call
* Avoid premature optimization
* Focus on correctness first

---

## Conclusion

Current implementation is:

* efficient enough for production
* safe under concurrency
* extensible for optimization
