# SDC-001: Design a Ticketing Server

**Difficulty:** Medium

- [SDC-001: Design a Ticketing Server](#sdc-001-design-a-ticketing-server)
  - [Overview](#overview)
  - [Success Criteria](#success-criteria)
  - [Functional API](#functional-api)
  - [Admission Policy](#admission-policy)
  - [Non-Functional Requirements](#non-functional-requirements)
  - [Constraints](#constraints)
  - [Evaluation Signals](#evaluation-signals)
  - [Common Pitfalls](#common-pitfalls)
  - [Optional Extensions](#optional-extensions)
  - [Out of Scope / Assumptions](#out-of-scope--assumptions)
  - [Reflection Prompts](#reflection-prompts)
  - [Time Guidance](#time-guidance)

## Overview

Build a Go service that meters access to a scarce resource under heavy load. The
service behaves like a virtual waiting room: users enter a queue, retain their
relative order, and are admitted at a controlled rate. Treat this as a
standalone practice challenge—work solo or have a partner review your solution
using the signals below.

## Success Criteria

1. Guarantee fairness (strict FIFO) while handling bursty traffic.
2. Keep concurrency control simple, correct, and observable.
3. Separate queue management from admission/token issuance.

## Functional API

| Action       | Endpoint              | Contract                                                     |
| ------------ | --------------------- | ------------------------------------------------------------ |
| Join queue   | `POST /join`          | Assign monotonic position and return a request identifier.   |
| Check status | `GET /status?id=<id>` | Report current position, admission flag, and optional token. |

**Join Queue Response**

```json
{
  "id": "<unique_request_id>",
  "position": <integer>
}
```

**Status Response**

```json
{
  "position": <integer>,
  "admitted": <boolean>,
  "token": "<optional_access_token>"
}
```

Tokens should live 1–5 minutes and be issued at most once per request.

## Admission Policy

- Configure a throughput cap of **N admissions per second**.
- Enforce FIFO ordering—no skipping or reordering.
- Admit each request at most once; revoke/expire tokens after TTL.

## Non-Functional Requirements

- **Concurrency:** Sustain 1k+ `POST /join` calls per second without race
  conditions.
- **Fairness:** Preserve order even during spikes; avoid starvation.
- **Consistency:** Exactly one position and at most one admission per request.
- **Performance:** O(1) join operations, bounded work on status checks,
  responsive under rate limiting.

## Constraints

- Implement in Go, preferring standard library primitives.
- Optimize for a single-node deployment baseline.
- Persistence is optional but trade-offs should be discussed.

## Evaluation Signals

- Correct synchronization (mutexes, channels, atomics) and clean separation of
  queue vs. admission state.
- Readable, tested (when time allows), and defensive Go code.
- Clear articulation of bottlenecks, trade-offs, and failure modes (restarts,
  retries, abuse, backpressure).

## Common Pitfalls

- Breaking FIFO guarantees or over-admitting beyond the rate limit.
- Polling with O(n) scans or global locks that throttle the system.
- Issuing duplicate admissions/tokens or ignoring idempotency for retries.

## Optional Extensions

- Idempotent `POST /join` semantics for retried clients.
- Durable storage across restarts or cross-node coordination.
- Push updates via WebSockets/SSE or backpressure cues for callers.

## Out of Scope / Assumptions

- No upstream user authentication or identity management beyond the admission
  token.
- No payments, ticket inventory management, or UI/front-end components.
- No per-IP fairness policies unless you explicitly add them as a stretch goal.

## Reflection Prompts

- Scaling strategy across multiple nodes while maintaining order.
- Ordering guarantees in distributed systems with partial failures.
- State that must survive restarts and how you would persist it.
- Abuse prevention, retry handling, and acceptable fairness trade-offs.

## Time Guidance

- **Implementation:** 60–90 minutes.
- **Discussion / Review:** 30–45 minutes.
