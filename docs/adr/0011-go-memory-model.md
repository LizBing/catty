# ADR-0011: Adopt Go memory model (not JMM/JSR-133)

- **Status:** Withdrawn
- **Date:** 2026-07-12
- **Withdrawn:** 2026-07-13 — the compatibility percentage is unsupported and
  several claimed Go happens-before/final-field guarantees are incorrect.
  Go is an implementation mechanism, not a replacement for Java observable
  memory semantics. R2 requires a new Java-semantics-first decision.

## Context

Java Memory Model (JSR-133) and Go Memory Model are both happens-before based,
but differ in details (`volatile` semantics, `final` field initialization
guarantees). Maintaining 100% JMM compliance requires extra barrier logic and
spec-level precision that adds implementation complexity.

## Decision

catty declares it follows **Go's memory model**, not JSR-133.

- `volatile` → Go atomic (consistent semantics).
- `synchronized` → `sync.Mutex` (implies happens-before).
- `final` field → Go's initial assignment (Go guarantees initialization
  visibility).
- `Thread.start`/`join` → goroutine creation/wait (Go guarantees
  happens-before).

Deviations are documented: a small fraction of programs relying on fine-grained
JMM semantics (lock-free data structures with specific memory ordering) may
behave differently. For 99.9% of programs, fully compatible.

## Consequences

**Positive**
- Simpler implementation (no extra barrier insertion logic).
- Performance: Go's atomic and mutex compile directly to hardware instructions,
  no JVM abstraction.
- Honest declaration: catty is not JCK-compliant.

**Negative**
- Not 100% Java-spec compliant for the most aggressive lock-free code.
- Must document deviations clearly.
