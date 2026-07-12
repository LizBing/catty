# ADR-0011: Adopt Go memory model (not JMM/JSR-133)

- **Status:** Superseded by ADR-0016
- **Date:** 2026-07-12

## Context

> **Supersession note (2026-07-12):** This proposal was not adopted. ADR-0016
> replaces it with protected Java memory semantics plus measured,
> explicitly-approved deviation gates.

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
