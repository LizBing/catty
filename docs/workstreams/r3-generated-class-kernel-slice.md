# R3 generated-class kernel slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Shared runtime kernel; no public Catty or Java SE generation API
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0022, ADR-0025,
ADR-0028 through ADR-0031, ADR-0033, and ADR-0034
**Prerequisites:** `r3-runtime-identity-definition-slice` and
`r3-typed-invocation-kernel-slice` Done; acceptance anchor fixed

## Outcome

The runtime can atomically define a loader-owned generated Class with ordinary
Class/Method/Field identities and typed executable bodies understood by
Interpreter and IR. Generated Classes participate in canonical identity,
assignability, initialization, heap storage, monitors, Thread context, and
exception transport without fabricating mutable classfile bytes or exposing a
host Go type as Java identity.

This slice creates no LambdaMetafactory, Proxy, arbitrary ClassLoader, hidden
class, or public generated-class capability.

## Scope

- Internal generated-class request/result model with explicit consumer
  provenance, defining loader, name policy, fields/methods, and typed bodies.
- Atomic definition/publication through the ADR-0033 loader service.
- Interpreter and IR execution of bounded typed generated bodies, including
  normal/throwable results and shared initialization.
- Consumer-owned cache hooks that do not impose one universal cache key or
  object-singleton policy.
- Test-only generated consumers for identity, assignability, dispatch,
  initialization, failure, and concurrent publication.

## Non-scope

Lambda, Proxy, InvocationHandler, byte-array `defineClass`, public Catty
generation API, hidden/anonymous classes, modules/packages/sealing, unloading,
weak caches, serialization, agents/instrumentation, arbitrary generated
bytecode, or AOT generated-class execution/fallback.

## Semantic constraints

Only fully linked Classes are published. Each consumer explicitly owns its
cache key, lifetime, generated-name policy, and terminal failures. Generated
objects share the ordinary Java object/Class world and cannot be substituted by
profile-visible Go objects. Generated method bodies use ADR-0031 typed values
and preserve Thread, initialization, heap, monitor, and throwable semantics.

## Acceptance

| Gate | Required result |
|---|---|
| Definition | generated identity, fields/methods, assignability, dispatch, and initialization kernel matrix Pass |
| Typed bodies | Interpreter/IR normal, primitive/reference, abrupt, and caller/Thread round trips Pass |
| Publication | concurrent same-request and conflicting-definition stress is race-free and publishes one Class or terminal failure |
| Consumer separation | two test consumers use distinct cache policies without kernel collision or singleton assumptions |
| Capability honesty | no Lambda/Proxy/ClassLoader fixture is claimed Supported; AOT generated paths precisely reject |
| Regression/governance | core/R2, unit/race, isolation, and `git diff --check` gates Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18. Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, and owner review are
frozen. Implementation remains unauthorized until prerequisites are Done and
this contract is fixed in an acceptance-anchor commit.
