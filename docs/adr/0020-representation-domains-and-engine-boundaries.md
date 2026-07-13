# ADR-0020: Representation domains and typed engine boundaries

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** [ADR-0003](./0003-tagged-slot.md)

## Context

R1's `rtda.Slot` efficiently stores interpreter locals and operand-stack cells
without numeric boxing. It was also used for fields, arrays, and AOT/runtime
bridge arguments. JVM logical slot rules do not require one physical storage
format across those domains, and that coupling prevents AOT code from fully
dissolving into typed Go values.

## Decision

JVM slots are bytecode-verification and interpreter-frame concepts, not
catty's universal value representation. Interpreter frames, IR values, Java
heap storage, and cross-engine call boundaries may use representations suited
to their domain. Explicit adapters preserve Java types, category rules, object
identity, exceptions, and runtime state.

The current `rtda.Slot` remains an allowed R1 interpreter-frame implementation.
It is not a stable IR schema, Java heap layout, AOT calling ABI, or generic
runtime value contract. IR should model logical typed values; AOT should prefer
typed Go locals and direct typed calls. A generic bridge is reserved for truly
dynamic or fallback boundaries and must represent category-2 values as one
language-level value rather than as two exposed frame cells.

## Consequences

- Existing `Slot` code need not be rewritten outside an Accepted workstream.
- Future heap-layout, bridge, and interpreter changes require separate
  measurement and semantic evidence.
- Shared object identity does not require shared primitive/frame layout.
