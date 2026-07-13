# ADR-0022: Bootstrap kernel capabilities and Java-visible facades

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** [ADR-0015](./0015-bootstrap-class-boundary.md)

## Context

R1 names six permanently synthetic classes. That list mixes bootstrap order,
Go-backed payloads, and one withdrawn Thread-to-goroutine assumption. A class
that currently carries a native payload is not thereby proven to require a
permanent whole-class synthetic implementation.

## Decision

catty defines its bootstrap boundary by irreducible capabilities, not a
permanent list of whole classes. Candidate kernel capabilities include type and
class identity, object allocation, class mirrors, string constant
materialization, throwable transport, execution-context attachment, and class
initialization state.

For each required capability, a future workstream must define the Go-native
kernel, Java-visible facade, provider precedence, reflection/type behavior,
cross-engine contract, and the point at which ordinary `java.base` classfiles
can participate. A Go-backed representation or helper does not by itself prove
that its Java facade must be synthetic.

## Consequences

- The current six-class R1 bootstrap set is an implementation baseline, not a
  permanent architecture constraint.
- New bootstrap capabilities require evidence and an Accepted workstream;
  synthetic fallback remains possible under ADR-0019.
