# R3 Java SE Class and annotation compatibility slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Optional Java SE Compatibility Profile
**Roadmap item:** Phase R3 — Reflection & dynamic features, compatibility track
**Governing ADRs:** ADR-0016, ADR-0019 through ADR-0023, ADR-0025,
ADR-0027 through ADR-0031, ADR-0033, and ADR-0034
**Prerequisites:** `r3-metadata-slice` and
`r3-runtime-identity-definition-slice` Done; separate acceptance anchor fixed

## Outcome

Interpreter and IR match Temurin 25 on the fixed six Class/type and four
runtime-annotation fixtures as an explicitly named Java SE Compatibility
Profile capability. Class-producing Java SE facade paths refer to canonical
runtime identities; `Class.forName`, declared-member discovery, and bounded
runtime annotations operate on retained metadata. AOT precisely rejects the
declared compatibility surface.

## Scope

- RuntimeVisible declaration/parameter annotations, AnnotationDefault,
  Exceptions, optional MethodParameters, and lossless annotation element trees
  required by the fixed ten rows.
- Java SE Class facade adapters over defining-loader-aware canonical runtime
  identities, including primitive/void/array/reference mirrors.
- Fixed Class queries, `Class.forName` overloads, declared Field/Method/
  Constructor discovery, and required Java exception types.
- Bounded Class/member/parameter annotation facades with defaults, inherited
  class lookup, repeatable containers, and race-free lazy caches.
- Precise AOT compatibility-profile rejection.

## Non-scope

Catty Runtime Profile metadata APIs; reflective member execution or Field
access; generic/type-use annotations; modules; records; sealed/hidden/local/
enclosing queries; broad custom ClassLoader; `setAccessible`; InvokeDynamic;
lambda; Proxy; or AOT fallback.

## Semantic constraints

This is an optional Java SE API contract, not Catty JVMS Core completion.
Discovery does not initialize a class. `Class.forName` obeys its initialize
flag and Java failure types. Defining loader plus binary name determines Class
identity. Member array order is not promised. Annotation values are lossless;
array values are defensively returned; `Inherited` searches superclasses only.
Facade caches are race-free but member facade reference identity is not a
public promise unless required by the declared Java SE contract.

## Acceptance

| Gate | Required result |
|---|---|
| Fixed profile rows | 10/10 Class+annotation rows Match Temurin in Interpreter and IR |
| Profile isolation | capability is labeled Java SE Compatibility; Catty Runtime Profile and Core make no facade claim |
| Remaining rows | 14/14 retain explicit prior classification; no omission or panic counted as support |
| AOT | all ten profile rows precisely NO-BUILD/Not implemented; no built panic |
| Identity/failure/init | canonical facade identity, missing class/member, and `Class.forName` behavior match Java 25 |
| Race/regression/governance | cache stress, core/R2, isolation, and `git diff --check` gates Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18 as an optional Java SE Compatibility Profile
contract. Outcome, Scope, Non-scope, Semantic constraints, Acceptance gates,
profile classification, and owner review are frozen. Implementation remains
unauthorized until prerequisites are Done and this contract is fixed in an
acceptance-anchor commit.
