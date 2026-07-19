# R3 typed dynamic-invocation kernel slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Catty JVMS Core shared kernel; reusable by profiles and Host ABI adapters
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0021, ADR-0025,
ADR-0028 through ADR-0031, ADR-0033, and ADR-0034
**Prerequisites:** `r3-runtime-identity-definition-slice` Done; acceptance
anchor fixed

## Outcome

Interpreter and IR share one logical typed dynamic-value/result service for
void, all Java primitives, references, normal results, and Java throwable
results. Direct runtime method/field/constructor consumers can invoke through
the service with explicit Thread context without exposing interpreter Slot,
IR register layout, heap-cell bits, or Go panic as the stable boundary.

This slice creates no reflection facade, MethodHandle, InvokeDynamic, Host ABI
provider, generated-class, or Java SE conversion capability.

## Scope

- Logical type/value/result vocabulary with one full value for category-2
  primitives and shared Java reference identity.
- Explicit execution/Thread, caller Class, normal-result, and Java-throwable
  transport.
- Interpreter-frame and IR-value adapters; direct dispatch, construction,
  field access, and shared initialization/HeapCell integration required by
  kernel unit consumers.
- Core receiver, arity, descriptor, primitive exactness, reference assignment,
  dispatch, initialization, and exception checks.
- Panic containment at the adapter boundary for internal defects, with a
  distinct non-success result that cannot be mistaken for a Java return value.

## Non-scope

Java wrappers/Object-array adapters; reflection access, boxing, widening,
varargs, or `InvocationTargetException`; MethodHandle adaptation; InvokeDynamic;
Host provider registration or authorization; generated methods; public
embedding API; or AOT dynamic invocation/fallback.

## Semantic constraints

Slot and engine layouts remain adapters. Java object/Class identity, Thread
context, initialization, heap publication, dispatch, and throwable identity
are preserved. A normal zero/null result is distinguishable from Java throwable
and internal-provider failure results. Profile policy cannot be embedded in the
kernel's value representation.

## Acceptance

| Gate | Required result |
|---|---|
| Value/result model | exhaustive primitive/category-2/reference/void/normal/throwable unit matrix Pass |
| Invocation | static, virtual, interface, special, constructor, field, abrupt, and initialization kernel tests Pass in Interpreter and IR adapters |
| Boundary | no stable API exposes Slot/frame/IR/heap layout; panic and Java throwable results are distinguishable |
| Identity/context | object/Class/Thread/caller and throwable identity survive round trips |
| AOT/capability | unsupported AOT dynamic calls are precisely rejected; no reflection/InvokeDynamic row newly claimed Supported |
| Regression/governance | core/R2, unit/race, isolation, and `git diff --check` gates Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18. Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, and owner review are
frozen. Implementation remains unauthorized until prerequisites are Done and
this contract is fixed in an acceptance-anchor commit.
