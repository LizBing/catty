# ADR-0009: Hybrid class library (~50 native Go classes + ~7000 interpreted from real JDK)

- **Status:** Withdrawn
- **Date:** 2026-07-12
- **Withdrawn:** 2026-07-13 — the numerical boundary and call-frequency claims
  lack evidence. ADR-0022 governs the bootstrap-capability boundary; future
  native promotion requires scoped behavioral evidence.

## Context

`java.base` has 7363 classes across 195 packages. Reimplementing all natively in
Go is impractical; interpreting all of them is too slow for hot paths. The key is
finding the right native-vs-interpreted boundary.

## Decision

A **hybrid class library**: ~50 critical classes implemented natively in Go;
~7000 classes loaded from the real JDK's `.class` files and interpreted.

Native layer selection criteria:

1. **bootstrapping-essential** — `Object`, `Class`, `ClassLoader`, `String`,
   `Thread`;
2. **performance hotspots** — `StringBuilder`, `Integer`/`Long` wrappers, `Math`,
   `Arrays`;
3. **runtime bridge** — `Throwable`/`Exception` hierarchy (~15 classes),
   `System`;
4. **concurrency primitives** — `Thread`, `ThreadGroup`.

The remaining ~7000 classes load from the real JDK's `java.base`, run via the
interpreter, with hot classes progressively AOT-transpiled as coverage grows.

## Consequences

**Positive**
- Semantic compatibility for the long tail (real JDK bytecode = identical
  behavior to OpenJDK).
- Native layer covers 90%+ of runtime call frequency.
- Controlled bootstrapping (manual `Object`→`Class`→`ClassLoader` load order,
  avoiding circular dependencies).

**Negative**
- JDK version binding (catty tracks the JDK version, currently Temurin 25).
- The native-vs-interpreted boundary must be tuned over time (move hot
  interpreted classes to native).
