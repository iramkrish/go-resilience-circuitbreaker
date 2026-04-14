# Sequence Diagram — Request Flow

## Normal Flow (CLOSED)

```
Client
  |
  | Execute()
  v
Circuit Breaker (CLOSED)
  |
  | Forward request
  v
Dependency
  |
  | Response
  v
Circuit Breaker
  |
  | Record success/failure (Strategy)
  v
Client
```

---

## Open State (Fail Fast)

```
Client
  |
  | Execute()
  v
Circuit Breaker (OPEN)
  |
  | -- Reject immediately
  v
Client (ErrCircuitOpen)
```

---

## Half-Open Probe Flow

```
Client
  |
  | Execute()
  v
Circuit Breaker (HALF_OPEN)
  |
  | Limited requests allowed
  v
Dependency
  |
  | Response
  v
Circuit Breaker

IF success:
  -> increment success count
  -> if threshold met → CLOSED

IF failure:
  -> immediate transition → OPEN
```

---

## Key Insight

* CLOSED → full traffic
* OPEN → zero traffic
* HALF_OPEN → controlled probing

---

## Critical Behavior

* No dependency calls in OPEN
* HALF_OPEN allows only a small number of requests
* Any failure in HALF_OPEN → immediate OPEN
