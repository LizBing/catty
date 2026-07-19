# R3 Java SE concat and lambda compatibility slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Optional Java SE Compatibility Profile
**Roadmap item:** Phase R3 — Reflection & dynamic features, compatibility track
**Governing ADRs:** ADR-0016, ADR-0019 through ADR-0023, ADR-0025,
ADR-0027 through ADR-0034
**Prerequisites:** `r3-java-se-invoke-bootstrap-slice` and
`r3-generated-class-kernel-slice` Done; separate acceptance anchor fixed

## Outcome

Interpreter and IR match Temurin 25 on StringConcatIndy, StatelessLambda,
CapturingLambda, and MethodReference as optional Java SE compatibility
capabilities. StringConcatFactory preserves Java conversion order and ADR-0027
UTF-16 values. Bounded LambdaMetafactory links a factory CallSite, captures
typed values, defines a loader-aware SAM implementation Class, and invokes the
required implementation handles. AOT precisely rejects these capabilities.

## Scope

- `javac 25` StringConcatFactory recipe/constants/types exercised by the fixed
  concat row.
- Bounded standard LambdaMetafactory for one SAM interface, stateless and
  capturing lambdas, and required static/virtual/interface/special references.
- Consumer-owned generated-lambda Class cache over ADR-0033 kernel.
- Primitive/reference/category-2 capture, adaptation, return, and abrupt paths.
- UTF-16 concat probes and concurrent first-link/cache stress.

## Non-scope

Catty Runtime bootstrap API; `altMetafactory` flags; serialization;
marker/bridge lists; arbitrary concat recipes/bootstrap methods; lambda object
singleton identity; hidden classes; Proxy; broad `java.lang.invoke`; AOT
execution/fallback; or JIT.

## Semantic constraints

These are optional Java SE profile capabilities. Linkage, capture, and
invocation remain distinct. Lambda object identity is unspecified; captured
Java value/object identity is preserved. Generated Classes use the declared
defining loader and ordinary Class/monitor/heap/init/exception semantics.
Concat is lossless in UTF-16 code units and handles null/primitives according
to the declared Java 25 surface.

## Acceptance

| Gate | Required result |
|---|---|
| Fixed profile rows | four concat/lambda rows plus prerequisite custom-bootstrap row Match Temurin in Interpreter/IR |
| Concat | primitive/reference/null/constant/recipe/UTF-16 supplemental matrix Pass |
| Lambda | stateless/capture/category-2/method-kind/abrupt/concurrent tests Pass |
| Generated Class | one race-free Class per declared lambda cache key; no lambda-object singleton promise |
| Scope/AOT | Proxy and broad factories remain Not implemented; profile rows precisely reject in AOT |
| Regression/governance | core/R2/kernel/race/isolation/`git diff --check` gates Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18 as an optional Java SE Compatibility Profile
contract. Outcome, Scope, Non-scope, Semantic constraints, Acceptance gates,
profile classification, and owner review are frozen. Implementation remains
unauthorized until prerequisites are Done and this contract is fixed in an
acceptance-anchor commit.
