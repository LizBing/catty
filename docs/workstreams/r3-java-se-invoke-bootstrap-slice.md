# R3 Java SE InvokeDynamic bootstrap compatibility slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Optional Java SE Compatibility Profile
**Roadmap item:** Phase R3 — Reflection & dynamic features, compatibility track
**Governing ADRs:** ADR-0016, ADR-0019 through ADR-0022, ADR-0025,
ADR-0028 through ADR-0032, and ADR-0034
**Prerequisites:** `r3-runtime-identity-definition-slice`,
`r3-typed-invocation-kernel-slice`, and `r3-invokedynamic-kernel-slice` Done;
separate acceptance anchor fixed

## Outcome

Interpreter and IR provide the bounded Java-visible Lookup, MethodType,
MethodHandle, and CallSite facade/protocol needed by the fixed
`BootstrapFailureOnce` fixture. A supported Java bootstrap receives the
required arguments, is resolved once per Method+PC, and publishes the JVMS
success or `BootstrapMethodError` failure result without re-execution. This is
a bounded Java SE compatibility capability, not general arbitrary-bootstrap or
complete `java.lang.invoke` support.

## Scope

- Minimal Java-visible facade identities and native/kernel adapters required by
  the fixed custom-bootstrap row.
- Bounded caller Lookup/access context and ordered name/type/static bootstrap
  arguments.
- Exact target/null/type validation and JVMS bootstrap failure translation.
- Interpreter and IR execution of the fixed custom bootstrap through the shared
  linkage and typed-invocation kernels.
- Concurrent/repeated failure and facade-identity supplemental probes.

## Non-scope

Broad Lookup API; arbitrary user bootstrap signatures; MethodHandle
combinators; `invokeWithArguments`; ConstantDynamic; Mutable/VolatileCallSite;
StringConcatFactory; LambdaMetafactory; generated classes; Proxy; Catty Runtime
bootstrap API; or AOT execution/fallback.

## Semantic constraints

The profile facade cannot create a second linkage state or object/Class world.
Per-site resolution, descriptor/target validation, publication, Thread/caller
context, and terminal failures reuse ADR-0032. The bootstrap executes once per
Method+PC. Throwable reference identity need not be reused where JVMS permits,
but failure category and cause behavior match the declared Java contract.

## Acceptance

| Gate | Required result |
|---|---|
| Fixed profile row | `BootstrapFailureOnce` Match Temurin in Interpreter and IR |
| Bootstrap protocol | name/type/Lookup/static args, null/type mismatch, thrown failure, and target invocation supplemental matrix Pass |
| Publication | repeated/concurrent resolution executes bootstrap once and publishes one terminal state without deadlock/race |
| Scope honesty | capability is bounded to declared protocol; concat/lambda and broad `java.lang.invoke` remain Not implemented |
| AOT/regression | precise AOT rejection plus core/R2/kernel gates Pass |
| Governance | evidence isolated and `git diff --check` Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18 as an optional Java SE Compatibility Profile
contract. Outcome, Scope, Non-scope, Semantic constraints, Acceptance gates,
profile classification, and owner review are frozen. Implementation remains
unauthorized until prerequisites are Done and this contract is fixed in an
acceptance-anchor commit.
