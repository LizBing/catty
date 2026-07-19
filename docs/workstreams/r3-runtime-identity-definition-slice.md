# R3 runtime identity and typed class-definition slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Catty JVMS Core shared kernel; no public ClassLoader API
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0022, ADR-0025,
ADR-0028 through ADR-0031, ADR-0033, and ADR-0034
**Prerequisites:** `r3-reflection-dynamic-research` Done;
`r3-metadata-slice` Done; acceptance anchor fixed

## Outcome

Every runtime Class has canonical defining-loader-aware identity. Lookup,
initiation, and definition return typed results rather than Java-reachable Go
panics, and concurrent definition publishes one fully linked Class or one
terminal linkage failure. Existing Class mirrors in Interpreter and IR refer to
the same canonical runtime identities.

This slice creates no `Class.forName`, declared-member, annotation, arbitrary
`ClassLoader.defineClass`, or generated-class capability.

## Scope

- Explicit defining-loader identity on runtime Class and loader-owned canonical
  `(defining loader, binary name)` definition state.
- Separate typed lookup/initiation and definition services, including Java
  linkage failure results and invariant-only panic boundaries.
- Canonical reference, array, primitive, and void internal type identities;
  audit existing Class-mirror producing paths against those identities.
- Atomic concurrent first definition, delegation-result identity, partial-link
  exclusion, and race-free cache publication.
- Adapters that preserve existing R2 initialization, Thread, monitor, heap, and
  exception semantics without adding profile APIs.

## Non-scope

Java `Class.forName`; declared members; annotations; arbitrary ClassLoader
subclasses or byte-array `defineClass`; modules/packages/sealing; generated
classes; unloading; weak caches; reflection invocation; InvokeDynamic; or AOT
dynamic loading/fallback.

## Semantic constraints

Class identity is defining loader plus binary name; initiating loader is not a
second Class identity. Arrays derive loader identity from their component
types; primitives and void use canonical VM identities. Supported classfile
lookup/definition failures are typed Java results. No caller observes a
partially linked Class, duplicate canonical Class, or alternate engine-specific
type world.

## Acceptance

| Gate | Required result |
|---|---|
| Identity | loader/name, delegation, array/component, primitive/void unit matrix Pass |
| Typed failure | missing, duplicate, malformed dependency, and linkage paths return typed results without Java-reachable panic |
| Concurrency | concurrent lookup/definition stress publishes one identity or one terminal failure; race/deadlock checks Pass |
| Mirror continuity | all existing Class-producing paths in Interpreter/IR refer to canonical runtime identities |
| Capability honesty | no Class/annotation/reflection fixture is newly claimed Supported; AOT unsupported paths remain exact rejection |
| Regression/governance | core/R2, unit/race, evidence-isolation, and `git diff --check` gates Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18. Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, and owner review are
frozen. Implementation remains unauthorized until prerequisites are Done and
this contract is fixed in an acceptance-anchor commit.
