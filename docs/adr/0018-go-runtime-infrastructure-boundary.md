# ADR-0018: Reuse Go runtime infrastructure without substituting Java semantics

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** [ADR-0001](./0001-reuse-go-runtime.md)

## Context

Go supplies mature allocation, garbage collection, scheduling, synchronization,
tooling, and native-code generation. Reusing those facilities is fundamental to
catty. ADR-0001 also implied that a Java Thread could equal a goroutine and
that Go memory behavior could stand in for Java behavior; those are separate
semantic questions.

## Decision

catty reuses Go runtime and toolchain infrastructure rather than implementing a
custom GC, scheduler, or native-code generator. It may use Go allocations, GC,
goroutines, synchronization primitives, atomics, profiling, race detection, and
the Go compiler as implementation mechanisms.

This decision authorizes no Java-visible substitution. A Java Thread is a Java
semantic object; a goroutine is one possible internal carrier. Go's memory
model, mutexes, channels, GC behavior, and panic behavior do not replace Java
25 Thread, monitor, interruption, class-initialization, exception, reference,
or memory-model requirements.

If Go infrastructure cannot directly supply a required Java guarantee, catty
must implement the required adaptation, explicitly reject the unsupported
capability, or obtain an Accepted deviation ADR. An accidental Go data race in
catty is an implementation defect even when the Java program itself has a data
race.

## Consequences

- Thread-to-goroutine mapping is not an architecture decision and requires
  separate R2 research before any supported concurrency claim.
- Go-native optimization remains encouraged when it preserves ADR-0017.
- ADR-0001 remains historical evidence for the original MVP choice but no
  longer governs implementation.
