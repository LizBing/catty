# ADR-0030: Race-free sequentially consistent heap cells for R2 concurrency

- **Status:** Accepted
- **Date:** 2026-07-13
- **Governing workstream:** `r2-concurrency-semantics-research`
- **Evidence:** `docs/workstreams/r2-concurrency-evidence/reports/{java25-concurrency-contract,go-mechanism-experiments,runtime-boundary-map}.md`

## Context

ADR-0018 requires Java semantics and treats catty-internal Go races as defects even when
the Java program has a data race. Current instance fields, static fields, and arrays expose
mutable `[]Slot` storage. Interpreter, IR, native code, runtime bridges, and AOT emitted
code read and write it directly. `Slot` is also a frame representation and cannot safely
become a copied Go atomic type.

Java concurrency requires monitor/start/join/volatile ordering and final-field freeze
semantics. Implementing only volatile fields atomically would leave ordinary racy Java
accesses as Go data races and would not supply final-field visibility after racy
publication.

## Proposed decision

For the first supported concurrency boundary, Java instance fields, static fields, and
array elements SHALL use a distinct race-free heap-cell representation under ADR-0020.
Interpreter-frame `rtda.Slot` remains separate.

All heap-cell reads and writes SHALL use sequentially consistent Go atomic operations:

- one atomic 64-bit bits cell for primitive values, including atomic `long`/`double`;
- one atomic Java-object reference cell for reference values; and
- descriptor-aware adapters between heap cells and frame/IR/AOT typed values.

All Java field and array accesses in the supported runtime use this boundary, including
ordinary non-volatile accesses. Volatile metadata remains classified, but its first
implementation may use the same stronger SC operation. Constructor writes and object
reference publication also use the boundary.

This intentionally selects a conservative sequentially consistent subset of Java-allowed
executions. It does not declare the Go memory model to be Java semantics. Stronger
ordering of ordinary fields is accepted because every produced behavior remains a
Java-allowed execution, while volatile and final fields receive at least their required
visibility and catty remains free of Go data races.

Mutable native payloads reachable from supported concurrent programs SHALL be protected
by their own race-free service or excluded with explicit unsupported behavior. Immutable
payloads such as ADR-0027 String values need no additional locking.

### Engine boundary

- Interpreter and IR SHALL stop accessing heap slices directly.
- Native/runtime code SHALL use typed heap APIs, including clone/arraycopy paths.
- AOT emitted code SHALL use the same typed heap boundary before concurrent AOT is
  supported. The first concurrency slice may reject all AOT concurrency fixtures.

## Consequences

- Heap layout is no longer `[]Slot`; this is the domain separation anticipated by
  ADR-0020.
- `long` and `double` become atomically read/written even when non-volatile, a permitted
  stronger choice than Java 25's minimum.
- Ordinary Java data races do not become Go data races in the heap-cell implementation.
- The approach adds atomic cost to all heap accesses. Performance optimization is
  deferred until after semantic acceptance and requires a superseding ADR or an
  evidence-backed mechanism that preserves this contract.
- Bulk operations and native payloads require an explicit concurrency audit; returning a
  mutable backing slice or pointer is no longer permitted.

## Non-scope

Weak-memory optimization, non-SC volatile lowering, VarHandle, Unsafe, final-field
write restrictions outside constructors, tearing optimization, object-layout tuning,
and concurrent AOT code generation.
