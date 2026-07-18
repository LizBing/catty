# R3 profile-separated implementation slice plan

**Status:** contracts Accepted by Owner; execution remains prerequisite- and anchor-gated
**Research anchor:** `6cf3636`
**Reconciled:** 2026-07-18 under Accepted ADR-0031 through ADR-0034

## Recommendation

The original six-step compatibility pipeline is replaced by ten independently
Accepted contracts: five shared-kernel slices and five optional Java SE
Compatibility Profile slices. No Catty Runtime Profile public API is selected
by this plan.

The shared-kernel track may complete without implementing reflection,
annotations, `java.lang.invoke`, StringConcatFactory, LambdaMetafactory, or
Proxy. Compatibility slices reuse the kernels but do not become prerequisites
for kernel completion. Interpreter, IR, and AOT state is reported separately
for every accepted capability.

### Shared-kernel track

| ID | Contract | Outcome | Prerequisites | Estimate |
|---|---|---|---|---:|
| K1 | `r3-metadata-slice` | BootstrapMethods and dynamic constant-pool metadata | Accepted ADRs/research Done | 0.75–1.25 weeks |
| K2 | `r3-runtime-identity-definition-slice` | defining-loader identity, typed lookup/definition, canonical runtime types | K1 | 1–1.5 weeks |
| K3 | `r3-typed-invocation-kernel-slice` | engine-neutral typed values/results and direct invocation | K2 | 1.5–2.5 weeks |
| K4 | `r3-invokedynamic-kernel-slice` | per-site linkage and logical MethodHandle/CallSite kernel | K1 + K3 | 1.5–2.5 weeks |
| K5 | `r3-generated-class-kernel-slice` | loader-owned generated Classes and typed bodies | K2 + K3 | 1–2 weeks |

K4 and K5 are independent after K3. Kernel completion is approximately
6–10 weeks for one Active Agent including evidence and Owner review; it does
not imply any Java SE dynamic-feature row is Supported.

### Optional Java SE Compatibility Profile track

| ID | Contract | Outcome | Prerequisites | Frozen rows | Estimate |
|---|---|---|---|---:|---:|
| C1 | `r3-class-annotation-facade-slice` | Class/type queries, declared members, runtime annotations | K1 + K2 | 10 | 1.5–2.5 weeks |
| C2 | `r3-reflective-execution-slice` | Method/Constructor/Field Java SE execution | C1 + K3 | 6 | 1.5–2.5 weeks |
| C3 | `r3-java-se-invoke-bootstrap-slice` | bounded Java-visible bootstrap protocol | K2 + K3 + K4 | 1 | 1–1.5 weeks |
| C4 | `r3-concat-lambda-slice` | StringConcatFactory and bounded LambdaMetafactory | C3 + K5 | 4 | 1.5–2.5 weeks |
| C5 | `r3-dynamic-proxy-compatibility-slice` | bounded Java SE dynamic Proxy | C1 + C2 + K5 | 3 | 1–2 weeks |

C1/C2, C3/C4, and C5 are optional capability branches. C5 does not depend on
C3/C4; concat/lambda does not depend on Java reflection. Completing all five
compatibility slices is approximately 6.5–11 weeks after their kernel
prerequisites, but each capability may be accepted or deferred independently.

## Dependency graph

```text
K1 metadata
  |
  +--> K2 identity/definition --> K3 typed invocation
  |        |                       |       |
  |        |                       |       +--> K4 InvokeDynamic kernel --> C3 Java SE bootstrap --> C4 concat/lambda
  |        |                       |
  |        +-----------------------+----------> K5 generated-class kernel
  |                                                        |
  +--> C1 Class/annotation --> C2 reflection --------------+--> C5 Proxy
```

The graph records semantic prerequisites, not permission to parallelize active
workstreams. Project governance still permits only one Active Agent/workstream
unless the Owner explicitly changes coordination.

## Kernel slice boundaries

### K1 — dynamic metadata

Retain BootstrapMethods and dynamic constant-pool structures without parsing
annotation/profile metadata or claiming Java-visible behavior. Candidate gates
are parser correctness, immutability, typed malformed-input failure, unchanged
24-row capability classification, and full regression/isolation checks.

### K2 — runtime identity and definition

Introduce defining-loader-aware Class identity, typed lookup/initiation/
definition, atomic publication, and canonical type identities. Java
`Class.forName`, arbitrary ClassLoader APIs, and generated Classes remain out
of scope. Candidate gates emphasize identity, failure, concurrency, existing
mirror continuity, and no new facade claim.

### K3 — typed invocation kernel

Introduce engine-neutral logical values and explicit normal/Java-throwable
results with Interpreter and IR adapters. Java reflection conversions,
MethodHandle adaptation, Host ABI provider policy, and public embedding remain
out of scope. Candidate gates exhaust value categories, dispatch, identity,
Thread/caller context, initialization, and layout isolation.

### K4 — InvokeDynamic linkage kernel

Implement per-Method/per-PC linkage state and logical MethodType,
direct-MethodHandle, and immutable-target CallSite kernels using an internal
test bootstrap. Kernel readiness is not Java SE bootstrap or opcode support.
Candidate gates verify linkage-once publication, recursion/concurrency,
descriptor/target failures, engine parity, and unchanged Java SE rows.

### K5 — generated-class kernel

Implement loader-owned atomic generated-Class definition and typed executable
bodies for test consumers. Lambda and Proxy caches remain consumer policy.
Candidate gates verify ordinary Class/heap/init/monitor/exception participation,
consumer separation, concurrent publication, and explicit AOT rejection.

## Compatibility slice boundaries

### C1 — Java SE Class and annotation facades

Parse only the annotation/member metadata required by the fixed ten rows and
expose Java SE facades over K2 identities. This is an optional profile
capability; it is not Core metadata completion.

### C2 — Java SE reflective execution

Layer Java wrappers, access/conversion/varargs policies, Field operations, and
`InvocationTargetException` over K3. No Catty Runtime Profile reflection API is
implied.

### C3 — Java SE bootstrap protocol

Provide the bounded Java-visible Lookup/MethodType/MethodHandle/CallSite
protocol needed by `BootstrapFailureOnce`, including JVMS failure translation.
This separates linkage-kernel completion from public `java.lang.invoke`.

### C4 — StringConcatFactory and LambdaMetafactory

Implement the four fixed concat/lambda rows over C3 and K5. Factory-specific
recipes, captures, generated-Class cache policy, and UTF-16 behavior remain a
named Java SE compatibility capability.

### C5 — dynamic Proxy

Implement the three fixed Proxy rows over C1/C2 and K5. Proxy owns its ordered
interface cache, InvocationHandler dispatch, Object/default methods, and
exception policy. It neither closes all R3 work nor depends on concat/lambda.

## Fixture and closure policy

- The original 24 fixtures remain immutable Java 25 compatibility evidence.
- Kernel slices use focused unit/integration evidence and must not relabel a
  Java SE row as Supported merely because an internal service exists.
- Each compatibility contract owns only its frozen rows and prerequisite
  continuity. No single optional contract is an automatic “R3 closure” gate.
- A future project-level milestone may aggregate all 24 rows only if the Owner
  explicitly accepts that Java SE Compatibility Profile target.
- AOT remains capability-specific `Not implemented`/precise build rejection
  until later Accepted work; a built-then-panic path is never fallback.

## Re-estimation points

Re-estimate after K2 because loader identity changes all later ownership, after
K3 because typed exception/value transport is the largest shared boundary,
after K4 because bootstrap linkage determines compatibility-adapter cost, and
after K5 because generated-body execution determines Lambda/Proxy feasibility.
Broad modules, arbitrary loaders/bootstraps, AOT dynamic fallback, or public
Catty Runtime APIs are scope changes requiring new contracts or amendments.
