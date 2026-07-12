# ADR-0016: Java memory semantics with measured deviation gates

- **Status:** Accepted
- **Date:** 2026-07-12
- **Supersedes:** ADR-0011

## Context

ADR-0011 proposed exposing Go's memory model and accepting an estimated 0.1%
compatibility loss. That percentage had no workload denominator, no litmus
evidence, and no Catty performance comparison. It also conflicted with Catty's
goal of preserving JRE semantics.

Go and Java both provide DRF-SC for correctly synchronized programs, and Go's
sequentially consistent atomics can implement Java volatile operations. But a
goroutine exit does not itself publish to another goroutine, Go has no Java
final-field freeze action, and ordinary racy or multiword Go accesses are not a
drop-in implementation of every JMM-allowed execution.

The performance premise is also unproven. Catty's largest gains may come from
AOT, escape analysis, and thread-local access elimination rather than from
weakening Java's memory semantics.

## Decision

Catty protects the following semantics as compatibility requirements:

- data-race-free Java programs;
- volatile, monitor, start/join, interrupt-observation, and class-initialization
  happens-before edges;
- final-field safe publication;
- non-corrupting reference and long/double access;
- selected Unsafe/VarHandle operations declared by ADR-0019.

Implementation uses a shared semantic heap-access layer. The initial production
backend is conservative and race-free. Interpreter, IR, AOT, native, and Unsafe
paths use that layer; AOT bypasses it only after proving an access is
thread-local or non-escaping.

Catty does not promise every outcome for incorrectly synchronized Java programs
until a three-backend study is complete:

1. **Strict** — conservative ordering and race-free storage;
2. **Go-native** — research-only direct Go storage where safe to investigate;
3. **Hybrid** — protected Java semantics with proven thread-local fast paths.

A deviation is described by semantic class, not “compatibility percentage”. It
is eligible for product consideration only when no protected semantic above is
lost, forbidden outcomes are absent from the supported profile, and it yields
at least one material gain:

- 2× in the affected microbenchmark;
- 15% application throughput; or
- 20% application p99 latency.

Each deviation still requires a separate ADR accepted by LizBing. Research-only
Go data races never enter production merely because a benchmark is faster.

## Evidence plan

- semantic matrix: Exact / Stronger / Adapter / Divergent / Unknown;
- 40–60 JMM/jcstress-derived litmus programs;
- HotSpot vs Strict vs Go-native vs Hybrid outcome comparison;
- representative application static/dynamic synchronization profile;
- micro and macro performance with repeated measurements;
- `go test -race` for every production backend.

Primary references:

- [JLS 25 Chapter 17](https://docs.oracle.com/javase/specs/jls/se25/html/jls-17.html)
- [Go Memory Model](https://go.dev/ref/mem)
- [OpenJDK jcstress](https://openjdk.org/projects/code-tools/jcstress/)

## Consequences

**Positive**

- Catty can pursue Go-native speed without making an unmeasured compatibility
  claim.
- Correctly synchronized Java and common immutable/concurrent code retain a
  clear contract.
- Conservative overhead becomes measurable debt with a defined optimization
  path through escape and thread-local analysis.

**Negative**

- R2 starts slower than direct unsynchronized Go field access.
- The study and semantic access layer add engineering work before performance
  claims can be generalized.
- Catty may ultimately document limited behavior for intentionally racy Java
  programs, so full JCK/JMM conformance remains unclaimed until proven.
