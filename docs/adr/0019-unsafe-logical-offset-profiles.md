# ADR-0019: Unsafe logical offsets and semantic profiles

- **Status:** Accepted
- **Date:** 2026-07-12

## Context

OpenJDK's `jdk.internal.misc.Unsafe` exposes field offsets, array base/scale,
plain and volatile access, atomics, fences, parking, raw memory, class
definition, and VM probes. Treating all native declarations as one “~50 method”
task hides radically different semantics.

Catty objects are Go objects with Slot-based storage. A HotSpot byte offset is
not a Go address, may not match Catty layout, and cannot safely bypass Go's GC.
Exposing Go pointers would undermine the project's central runtime-reuse premise.

## Decision

Unsafe offsets are opaque logical tokens interpreted by Catty. A token denotes
an instance field, static field, or array location; it never denotes a Go
address. Array base/scale APIs may return stable synthetic arithmetic values
that Catty translates back to an element.

Unsafe support is delivered by caller-backed profiles:

- **U0 — bootstrap/layout:** registration, VM constants, field-token lookup,
  array base and scale;
- **U1 — heap access:** plain/volatile primitive and reference get/put over
  logical storage;
- **U2 — atomics:** CAS/exchange and required read-modify-write operations;
- **U3 — fences:** full/acquire/release ordering integrated with ADR-0016;
- **U4 — parking:** one-permit park/unpark tied to ADR-0017 Thread state.

All profiles use the semantic heap layer from ADR-0016. Invalid logical offsets
and unsupported forms fail explicitly.

Initial profiles exclude raw absolute-address access, off-heap allocation,
bulk raw memory, class injection, writeback/cache-line operations, and arbitrary
HotSpot object-header assumptions. These require separate ADRs if demanded by a
representative workload.

The exact pinned-JDK caller graph is checked in at
[`../research/R2_UNSAFE_CALL_GRAPH.md`](../research/R2_UNSAFE_CALL_GRAPH.md).

## Consequences

**Positive**

- Preserves Go GC visibility and avoids pinning Catty to HotSpot object layout.
- Lets real JDK bytecode perform familiar offset arithmetic without exposing a
  physical address.
- Keeps Unsafe implementation proportional to verified callers.

**Negative**

- Code that intentionally relies on HotSpot addresses or object headers is not
  compatible with the initial profile.
- Translation adds overhead to Unsafe access until AOT can specialize known
  logical tokens.
- Reflection-based Field offset APIs remain coupled to later reflection work.
