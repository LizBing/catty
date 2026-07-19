# R3 Java SE reflective-execution compatibility slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Optional Java SE Compatibility Profile
**Roadmap item:** Phase R3 — Reflection & dynamic features, compatibility track
**Governing ADRs:** ADR-0016, ADR-0019 through ADR-0023, ADR-0025,
ADR-0027 through ADR-0031, and ADR-0034
**Prerequisites:** `r3-class-annotation-facade-slice` and
`r3-typed-invocation-kernel-slice` Done; separate acceptance anchor fixed

## Outcome

Interpreter and IR match Temurin 25 on the fixed six reflective member rows
through Java SE facade adapters over ADR-0031 typed invocation. Method,
Constructor, and Field operations preserve Java receiver, conversion, access,
dispatch, initialization, heap, and exception semantics. AOT precisely rejects
the compatibility capability.

## Scope

- Java wrapper/Object-array adapters over logical typed values; no stable Slot
  ABI.
- Bounded `Method.invoke`, `Constructor.newInstance`, and Field
  get/set/getX/setX.
- Java SE receiver/arity/access checks, boxing/unboxing/widening, reference
  assignment, varargs-array behavior, virtual dispatch, static initialization,
  HeapCells, and `InvocationTargetException`.
- Unit conversion matrix and race tests for reflective static access/init.

## Non-scope

Catty Runtime Profile invocation API; InvokeDynamic; public MethodHandle
combinators; module access; broad `setAccessible`; final-field mutation
guarantees; serialization; JNI; generated classes; Proxy; or AOT
reflection/fallback.

## Semantic constraints

This is an optional Java SE API contract. Target Java throwable identity is
preserved as `InvocationTargetException` cause. No Go panic crosses a
Java-visible path. Static operations initialize the declaring class through
ADR-0025. Field access uses typed ADR-0030 APIs. Virtual invocation dispatches
on receiver Class. Only identity and Java method-invocation widening
conversions are permitted.

## Acceptance

| Gate | Required result |
|---|---|
| Fixed profile rows | six reflective rows and prerequisite ten facade rows Match Temurin in Interpreter and IR |
| Conversions | exhaustive primitive/reference/receiver/arity normal+failure table Pass |
| Init/exceptions | static init, abrupt target, constructor failure, access, and wrong receiver match Java 25 |
| Profile isolation | typed kernel tests remain profile-neutral; no Catty Runtime reflection API is implied |
| Remaining/AOT | eight dynamic rows retain prior classification; all reflective profile rows precisely reject in AOT |
| Race/regression/governance | reflective stress, core/R2, isolation, and `git diff --check` gates Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18 as an optional Java SE Compatibility Profile
contract. Outcome, Scope, Non-scope, Semantic constraints, Acceptance gates,
profile classification, and owner review are frozen. Implementation remains
unauthorized until prerequisites are Done and this contract is fixed in an
acceptance-anchor commit.
