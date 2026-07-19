# R3 InvokeDynamic linkage-kernel slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Catty JVMS Core shared kernel; no Java SE bootstrap compatibility claim
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0021, ADR-0025,
ADR-0028 through ADR-0032, and ADR-0034
**Prerequisites:** `r3-metadata-slice` and
`r3-typed-invocation-kernel-slice` Done; acceptance anchor fixed

## Outcome

Interpreter and IR share a bounded logical MethodType/direct-MethodHandle/
constant-target CallSite kernel and per-runtime-Method/per-PC InvokeDynamic
resolution service. Successful or failed linkage is race-free, published once,
and invoked through ADR-0031 typed values. AOT precisely rejects unsupported
dynamic sites.

This slice establishes kernel readiness only. No fixed Java SE R3 row or
general opcode support is claimed until a declared bootstrap protocol adapter
is implemented and verified.

## Scope

- Canonical logical MethodType and bounded direct MethodHandles required by
  kernel tests.
- Immutable/constant-target CallSite state and target-type validation.
- Per-instruction unresolved/resolving/linked/failed state with explicit Thread
  owner, wait/retry, bounded recursion policy, and terminal publication.
- Interpreter opcode `0xba` and typed IR service operation behind explicit
  bootstrap-capability lookup.
- Research-derived internal test bootstrap provider exercising successful,
  null/type-mismatch, thrown, recursive, and concurrent linkage without
  becoming a public Catty Runtime Profile API.

## Non-scope

Java-visible Lookup/MethodType/MethodHandle/CallSite facades; arbitrary Java
bootstrap methods; `BootstrapFailureOnce` compatibility fixture;
StringConcatFactory; LambdaMetafactory; ConstantDynamic execution;
mutable/volatile CallSite; generated classes; Proxy; AOT execution/fallback;
or JIT.

## Semantic constraints

Linkage state is keyed by runtime Method identity plus PC, not constant-pool
entry alone. VM synchronization is separate from Java monitors. One terminal
result is published; linkage is not re-executed after success or failure.
Bootstrap and target calls preserve descriptor types, Java identity,
Thread/caller context, initialization, and throwable results. Kernel-only
readiness is never reported as support for an unimplemented profile bootstrap.

## Acceptance

| Gate | Required result |
|---|---|
| Linkage kernel | success, null/type mismatch, thrown, recursion, and descriptor/target validation tests Pass |
| Publication | internal bootstrap executes once per Method+PC; concurrent/race stress publishes one terminal state without deadlock |
| Engine parity | Interpreter and IR adapters produce identical typed normal/throwable kernel results |
| Capability honesty | all 24 Java SE R3 rows retain their prior classification; no named bootstrap is claimed Supported |
| AOT/regression | precise AOT rejection plus core/R2/go vet/test/race gates Pass |
| Governance | evidence isolated and `git diff --check` Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18. Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, and owner review are
frozen. Implementation remains unauthorized until prerequisites are Done and
this contract is fixed in an acceptance-anchor commit.
