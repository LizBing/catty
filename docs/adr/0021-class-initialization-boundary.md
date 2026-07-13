# ADR-0021: Class initialization is a shared Java runtime boundary

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** [ADR-0005](./0005-lazy-clinit.md)

## Context

ADR-0005 correctly separated loading/linking from initialization, but made an
interpreter frame-push technique part of the decision. The current R1 code has
already changed to synchronous execution, and Java 25 class/interface
initialization requires ordering, recursion, failure, interface, and
cross-thread rules beyond a single `initStarted` flag.

## Decision

Loading, linking, and initialization are distinct. Classloader operations do
not execute Java `<clinit>` merely because a class was loaded or linked.
Initialization is a shared runtime responsibility invoked at Java 25-required
use points; no execution engine may define incompatible initialization
semantics.

The concrete Java 25 class/interface initialization state machine remains a
research decision. Until it is accepted, R1 class initialization is a limited
existing capability and does not authorize concurrency, interface-initialization,
or full failure-semantics claims.

## Consequences

- Frame push, synchronous nested execution, and any future AOT barrier are
  replaceable implementation strategies.
- R2 research must propose the detailed state machine before concurrency or
  broad class-library compatibility is implemented.
