# ADR-0012: Escape analysis replaces generational GC

- **Status:** Withdrawn
- **Date:** 2026-07-12
- **Withdrawn:** 2026-07-13 — the current generic object representation does
  not establish the claimed escape-analysis results, while field reordering and
  value conversion affect identity, reflection, synchronization, and layout.
  These ideas belong in a measured future performance research workstream.

## Context

Go's GC is non-generational concurrent mark-sweep. Traditional JVMs use
generational GC (young/old) to efficiently handle Java's typical allocation
pattern (many short-lived objects). catty doesn't implement its own GC, but it
can leverage Go's compiler escape analysis to reduce allocations at the source —
the AOT path.

## Decision

catty's GC strategy is **"reduce allocations" rather than "collect fast"**. Three
mechanisms:

1. **AOT path**: transpiled Go code goes through Go's compiler escape analysis —
   Java objects that don't escape (locally created, locally used, not stored in
   globals) are stack-allocated by Go, never reaching the GC.
2. **Object layout optimization**: AOT analyzes field access patterns and
   reorders fields for cache line utilization (Java fields are
   declaration-ordered; Go can freely reorder).
3. **Value-type identification**: if a Java class's objects are never synchronized
   on, never have `System.identityHashCode` called, never compared by `==`
   identity — they can be emitted as Go value types (stack-allocated structs)
   rather than pointer types (heap-allocated).

## Consequences

**Positive**
- Short-lived objects (Java's dominant allocation pattern) eliminated by escape
  analysis — no generational GC needed.
- Go GC pause times <1ms (measured), an advantage for latency-sensitive Java
  applications.
- `new Point(3,4)` in an inner loop may not allocate at all — Go's compiler
  discovers non-escaping and stack-allocates or register-allocates.

**Negative**
- For workloads with many long-lived objects (large caches), Go GC throughput is
  lower than G1/ZGC (no young-gen fast path).
- Value-type identification is conservative (false negatives — safe objects
  treated as references).
