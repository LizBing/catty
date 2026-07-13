# ADR-0017: Java 25 semantic baseline for supported capabilities

- **Status:** Accepted
- **Date:** 2026-07-13

## Context

catty uses Go-native representations and execution mechanisms, but the project
needs one external semantic contract before it can make compatibility or
performance claims. A reference implementation is also needed for differential
testing without confusing one implementation detail with the Java specification.

## Decision

For MVP, catty defines the meaning of every declared supported capability
against Java SE 25, JLS 25, and JVMS 25. The repository-pinned Eclipse Temurin
25 toolchain is the differential-test reference implementation.

This is a capability-scoped contract, not a claim of complete Java 25 language,
class-file, JRE API, or JVM-feature coverage. Source precedence is:

1. JVMS;
2. JLS;
3. Java SE API contracts;
4. pinned Temurin 25 behavior where the specification permits choices; then
5. catty documentation only for supported scope and Owner-approved deviations.

An implementation may select a subset of Java-permitted nondeterministic
outcomes, but it must not produce Java-forbidden behavior or remove required
guarantees. Unsupported behavior must be explicitly reported or rejected;
Go implementation convenience is not an implicit semantic deviation.

An intentional deviation requires a separate Accepted ADR that records the
affected Java requirement, observable catty behavior, affected programs,
deterministic and differential evidence, measured benefit, Owner acceptance,
and a re-evaluation condition.

## Consequences

- Go runtime, Go compiler, and Go standard-library behavior are implementation
  mechanisms, not semantic waivers.
- Engine claims must identify Interpreter, IR, and AOT status independently as
  `Supported`, `Fallback`, or `Not implemented`.
- Temurin behavior resolves permitted choices; it does not override Java 25
  specification requirements.
- Upgrading the pinned Temurin build is a reproducibility maintenance change;
  changing the Java 25 semantic baseline requires a new ADR.
