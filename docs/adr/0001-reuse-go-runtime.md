# ADR-0001: Reuse the Go runtime (no custom GC, scheduler, or JIT)

- **Status:** Accepted
- **Date:** 2026-07-11

## Context

A "complete" JVM is dominated by three heavy, interdependent subsystems: a
garbage collector, a thread scheduler, and a JIT compiler. Together they are
the bulk of HotSpot's engineering — hundreds of thousands of lines, years of
tuning. Building them from scratch is what makes "write a JVM" a multi-year
project rather than a weekend one.

The project's stated goal is an *experimental* JVM that deliberately **sits on
top of the Go runtime**. Go already provides a concurrent, low-latency GC; an
M:N goroutine scheduler (GMP); and a TCMalloc-style allocator — all production
hardened.

## Decision

Do not implement any of the three heavy subsystems. Instead:

- **GC**: Java objects are Go heap allocations (`*rtda.Object`). Go's GC traces
  them natively; catty writes no mark/sweep, no write barriers, no card tables.
- **Scheduler**: MVP is single-threaded. The concurrency arc maps
  `java.lang.Thread` to a goroutine, inheriting the GMP scheduler.
- **JIT**: none in Phase 1. The Phase 2 plan (ADR-0001's successor in
  `ROADMAP.md` Theme A) replaces the interpreter with a bytecode→Go-source AOT
  transpiler, letting `go build` act as the optimizing backend.

## Consequences

**Positive**
- The implementation is small (~3300 LOC for a working interpreter MVP) because
  the three hardest subsystems are inherited. The MVP floor collapses to "a
  bytecode interpreter plus a class loader" — bounded, deliverable in days.
- Java objects get a real, concurrent GC for free; native methods interoperate
  trivially (a Java `String` can carry a Go `string` in its `extra` payload).
- The AOT arc is uniquely enabled by this choice: lowering bytecode to Go source
  turns Go's toolchain into the JIT backend, with no SSA/register-allocation
  work.

**Negative**
- A hard performance ceiling: without a JIT, catty is ~87× slower than HotSpot
  JIT on `fib(35)` (and ~7× slower than HotSpot's own interpreter — see
  ADR-0002). "Optimal performance" is a Phase 2 outcome, not Phase 1.
- Go's GC is non-generational; allocation-heavy Java workloads (many short-lived
  objects) won't be collected as efficiently as HotSpot's young generation.
- The Java Memory Model is not Go's memory model. `volatile`/happens-before
  semantics can only be approximated once concurrency lands — a correctness
  hazard to document, not hide.
- Some JVMS features that assume control over the stack/GC (stack walking for
  exceptions, `sun.misc.Unsafe`, precise GC pinning) are awkward or impossible
  to express through the Go runtime.
