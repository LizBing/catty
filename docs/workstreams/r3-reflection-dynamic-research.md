# R3 reflection and dynamic-features research

**Status:** In Progress
**Type:** research
**Review:** owner
**Proposed base:** `ccc6046`
**Acceptance anchor:** this 2026-07-17 governance commit; the research branch
records its resolved SHA and actual base before creating fixtures or evidence
**Roadmap item:** Phase R3 — Reflection & dynamic features
**Governing ADRs:** ADR-0016 through ADR-0025, ADR-0027, and ADR-0029 through ADR-0030
**Prerequisites:** R2 complete; no active implementation workstream

## Outcome

Establish an evidence-backed Java 25 semantic and engine contract for bounded
reflection, runtime-visible annotations, `invokedynamic`, and dynamic proxies.
The research fixes a differential baseline, maps the classfile/runtime/loading
gaps, proposes the durable ADR decisions required by ADR-0016 and ADR-0022,
and produces an ordered set of independently acceptable implementation slices.

This workstream does not add R3 production capability. Its deliverables are
reports, frozen fixtures and baseline evidence, Proposed ADRs, and Proposed
implementation contracts for later Owner review.

## Pre-contract repository survey

The 2026-07-17 planning survey used repository `ccc6046`, Temurin/JDK 25.0.3,
and Go 1.26.5. It found the following implementation facts.

### Existing foundations

- `rtda.Class` retains class name, hierarchy, access flags, fields, methods,
  descriptors, bytecode, and a constant pool; Interpreter and IR share these
  runtime identities.
- Each runtime Class has a concurrency-safe canonical `java.lang.Class` mirror,
  and R2 supplies Java 25 class initialization, object identity, monitors, and
  race-free heap cells.
- The classloader is runtime-available and concurrency-safe, and ADR-0016 keeps
  the Interpreter as a permanent semantic fallback for explicitly dynamic code.
- The runtime already has descriptor parsing and interpreted/native invocation
  machinery that can inform, but does not yet define, a reflection invocation
  boundary.

### Metadata and facade gaps

- `classfile` decodes only `Code` and `StackMapTable`; every other attribute is
  discarded as `UnparsedAttribute`. This includes `BootstrapMethods`,
  `RuntimeVisibleAnnotations`, parameter/type annotations, `AnnotationDefault`,
  `Signature`, `Exceptions`, `MethodParameters`, and class-relationship
  metadata used by broader reflection.
- Constant-pool tags for MethodHandle, MethodType, and InvokeDynamic are parsed,
  but their payloads have no usable public resolution API. ConstantDynamic is
  not parsed.
- `rtda.Class`, `Method`, and `Field` retain only part of the metadata needed by
  Java reflection. Defining-loader identity and several declared-member,
  interface, modifier, parameter, exception, generic-signature, and annotation
  views are absent or not exposed.
- The synthetic `java.lang.Class` facade implements a small type-query surface,
  but not `forName`, declared-member discovery, annotations, constructors, or
  reflective invocation. Some existing Class-producing natives allocate a new
  mirror instead of using the canonical mirror, which must be covered by R3
  identity fixtures rather than silently inherited as a contract.

### Loading, invocation, and failure gaps

- `rtda.Loader.LoadClass` returns only `*Class`; a missing class becomes a Go
  panic. `Class.forName` needs Java-visible lookup failure, initialization
  control, and defining-loader semantics.
- The current standard loader has one name-keyed cache, while runtime Class
  values do not retain a Java-visible defining loader. Dynamic proxies and
  custom loading therefore require an explicit identity/definition decision.
- Reflection has no generic Java-value adaptation boundary for boxing,
  unboxing, primitive widening, varargs, receiver checks, access checks, or
  `InvocationTargetException`. The AOT bridge still uses interpreter `Slot`
  values for fallback and has known category-2 and exception limitations, so it
  cannot be adopted as the reflection ABI without research under ADR-0020.
- Required Java exception/facade coverage such as `ClassNotFoundException`,
  `NoSuchMethodException`, `IllegalAccessException`, and
  `InvocationTargetException` is not established.

### Dynamic-linkage gaps

- Interpreter and IR do not execute opcode `0xba`; lowering and AOT have no
  CallSite or MethodHandle linkage model. The AOT concurrency scanner merely
  knows the instruction length and is not an R3 capability boundary.
- There is no BootstrapMethods representation, per-call-site resolution state,
  bootstrap failure memoization, MethodHandle adaptation/invocation service,
  or Java facade for MethodType/MethodHandle/CallSite.
- LambdaMetafactory and StringConcatFactory are absent. Runtime class synthesis
  exists only as a Go construction helper for core facades; it is not a
  defining-loader-aware Java class-definition protocol.
- Dynamic proxies are not an `invokedynamic` subtype. They require their own
  generated-class identity, loader, InvocationHandler dispatch, Object-method,
  default-method, exception, and cache semantics.

## Scope

### 1. Fixed Java 25 differential matrix

Freeze the following 24 source fixtures before baseline execution. Each fixture
records combined stdout, stderr, and exit status from Temurin 25.0.3,
Interpreter, IR, and AOT build/run classification.

| Family | Frozen fixtures | Contract questions |
|---|---|---|
| Class/type reflection (6) | `ClassIdentity`, `ClassForNameInit`, `ClassQueries`, `DeclaredMembers`, `PrimitiveAndArrayClass`, `MissingClass` | canonical identity, initialization trigger, hierarchy/type views, member order policy, primitive/array mirrors, Java failure |
| Member access/invocation (6) | `MethodInvoke`, `ConstructorInvoke`, `FieldGetSet`, `StaticReflectiveInit`, `ReflectiveConversions`, `ReflectiveFailures` | dispatch, construction, heap access, initialization, boxing/widening/varargs, access/target exception transport |
| Runtime annotations (4) | `ClassAnnotation`, `MemberAnnotation`, `AnnotationDefaults`, `InheritedRepeatableAnnotation` | retention, element values, defaults, inheritance, repeatable containers, annotation identity/equality |
| InvokeDynamic (5) | `StringConcatIndy`, `StatelessLambda`, `CapturingLambda`, `MethodReference`, `BootstrapFailureOnce` | BootstrapMethods parsing, call-site linkage, capture/adaptation, interface dispatch, one-time failure publication |
| Dynamic proxy (3) | `ProxyDispatch`, `ProxyObjectMethods`, `ProxyFailureAndDefault` | loader/class cache, InvocationHandler, Object methods, default method and exception rules |

The fixture denominator, source bytes, expected output comparison, and engine
classification become frozen when the Owner accepts this contract. Adding or
replacing a fixture after acceptance requires an accepted amendment; optional
probes may not change the denominator.

### 2. Required reports

- `r3-java25-semantic-contract.md`: observable Java 25 rules for the fixed
  surface, including initialization, identity, access, conversion, ordering,
  exception, linkage, and concurrency behavior.
- `r3-classfile-metadata-gap.md`: attribute/constant-pool coverage and the
  minimum lossless metadata representation needed by the matrix.
- `r3-runtime-boundary-map.md`: call and ownership paths across classfile,
  loader, `rtda`, native facades, Interpreter, IR, AOT, initialization, heap,
  and exception transport.
- `r3-engine-capability-matrix.md`: current baseline plus proposed
  `Supported`/`Fallback`/`Not implemented` completion state for each engine and
  family. A built-then-panic path is a failure, not a fallback.
- `r3-slice-plan.md`: ordered implementation slices, prerequisites, proposed
  gates, and explicit non-scope.

### 3. Required architectural decisions

Produce Proposed ADRs, split when the evidence shows independently durable
boundaries, covering at least:

- retained reflection/annotation metadata, canonical Java facades, access and
  invocation semantics, and class-initialization triggers;
- MethodHandle/MethodType/CallSite identity, InvokeDynamic linkage and failure
  memoization, and the multi-engine fallback/AOT policy;
- defining-loader identity and runtime-generated classes for proxies/lambdas,
  including cache and concurrency semantics.

The research may recommend more than one ADR. It must not present ADR-0007's
superseded concrete mechanisms as current decisions.

### 4. Planned implementation sequence

The final report must validate or revise this initial dependency order:

1. lossless classfile metadata and runtime metadata ownership;
2. canonical Class/member/annotation facades and bounded type discovery;
3. reflective construction, field access, and method invocation through a
   typed dynamic-value/exception boundary;
4. MethodType/MethodHandle/CallSite kernel and InvokeDynamic linkage;
5. StringConcatFactory, then LambdaMetafactory;
6. defining-loader-aware generated classes and dynamic proxies;
7. later AOT expansion only after the relevant mixed-engine exception and
   typed-call boundaries are accepted.

## Non-scope

- Production implementation or modification of classfile, classloader, `rtda`,
  native, Interpreter, IR, runtime, transpiler, launcher, tests, or CI.
- A claim of complete `java.lang.reflect`, MethodHandle combinators, arbitrary
  bootstrap methods, custom ClassLoader compatibility, modules, records,
  sealed classes, JNI, serialization, ServiceLoader, Unsafe/VarHandle, hidden
  classes, agents/instrumentation, or JCK compliance.
- AOT implementation of reflection, InvokeDynamic, proxies, dynamic class
  definition, or new cross-engine exception propagation.
- Treating string concatenation, lambdas, and dynamic proxies as one mechanism
  merely because they are grouped in Roadmap Phase R3.
- Changing ADR-0016's AOT-primary product direction, permanent Interpreter
  fallback, or explicit per-engine capability reporting.

## Semantic constraints

- Java 25 and Temurin 25.0.3 are the semantic and differential baselines.
- Class/member/annotation/MethodHandle/CallSite/proxy identity and failure
  behavior are Java-visible semantics and cannot be approximated by Go pointer
  convenience or panics.
- Reflection and dynamic linkage use the existing R2 class-initialization,
  object-identity, monitor, heap-publication, and exception contracts; the
  research must identify every additional trigger or synchronization edge.
- Engine selection is not Java-visible. Interpreter fallback must be explicit,
  share the same object/class world, and preserve return values, exceptions,
  initialization, and Thread context.
- A generic boundary must follow ADR-0020 and cannot expose interpreter Slot
  layout as a stable reflection or dynamic-linkage ABI.
- Synthetic facades and generated classes remain revisable capability
  boundaries under ADR-0019 and ADR-0022, not permanent whole-class decisions.

## Acceptance

| Gate | Required evidence | Result |
|---|---|---|
| Fixture freeze | Pass — exactly 24 reviewed Java sources; `manifest.sha256` verifies all hashes |
| Temurin baseline | Pass — 24/24 compile/run on Temurin 25.0.3, exit 0, bounded time, captured stdout/stderr/exit |
| Catty baseline | Pass — all 24 attempted: Interpreter 0 Match/24 Exit(1), IR 0 Match/24 Exit(1), AOT 24 NO-BUILD; 24/24 rows, no timeout/omission |
| Metadata gap | Pass — required attributes/CP entries mapped to retained, discarded, structural-only, or Not implemented with consumers |
| Runtime boundary | Pass — loading, initialization, invocation, heap, exception, identity, concurrency, and engine transitions mapped to concrete packages |
| Decision coverage | Pass — Accepted ADR-0031 through ADR-0033, revised under ADR-0034, separate metadata/typed invocation, per-site dynamic linkage, and defining-loader/generated-class kernels from optional Java SE facades |
| Slice plan | Pass — ten Accepted contracts separate five shared-kernel slices from five optional Java SE Compatibility Profile slices, with independent prerequisites, engine states, non-scope, and gates |
| Governance | Pass — historical evidence unchanged; R3 files only under `docs/`; `git diff --check` passes; PROJECT_STATUS remains honest |

Results use only `Pass`, `Fail`, `Not run`, or `Not implemented`. Baseline
failures are research data, but a missing row, silently skipped engine, or
unbounded process is a failed research gate.

## Plan

| Step | State |
|---|---|
| Pre-contract repository survey and Proposed contract | Complete |
| Owner reviews/fixes and accepts research contract | Complete |
| Fix acceptance anchor and record research preflight | Complete |
| Freeze fixtures and harness; capture Temurin/catty baselines | Complete |
| Produce metadata and runtime-boundary reports | Complete |
| Draft Proposed ADRs and ordered implementation slices | Complete |
| Reconcile outputs with ADR-0034 profile boundaries | Complete — ADR-0031 through ADR-0033 and ten profile-separated implementation contracts Accepted |
| Owner reviews research conclusions | Pending |
| Mark research Done and select first R3 implementation contract | Pending |

## Acceptance record

Accepted by Owner on 2026-07-17. This acceptance freezes Outcome, Scope,
Non-scope, Semantic constraints, the 24-fixture denominator, Acceptance gates,
and owner review. It authorizes research deliverables after this record is
fixed in a Git acceptance anchor; it does not authorize production
implementation or integration beyond the accepted research scope.

## Research preflight

- **Acceptance anchor / actual base:** `6cf3636` / `6cf3636`
- **Branch:** `codex/r3-reflection-dynamic-research`
- **Toolchain:** Temurin/JDK 25.0.3; Go 1.26.5 darwin/arm64
- **Historical evidence diff:** `git diff --name-only ca42a61..6cf3636 --
  docs/workstreams/r2-evidence docs/workstreams/r2-initialization-evidence
  docs/workstreams/r2-string-evidence docs/workstreams/r2-concurrency-evidence
  docs/workstreams/r2-concurrency-candidate-evidence` — Pass (empty)
- **Historical evidence tree IDs:** R2 research `b4acb756eb8687ffcee2d5ad4231bce448673fe3`;
  initialization `54fdeaa0c7691cf2cd11a45e907c36752079fe8c`;
  String `5f38e9b1083680b881fd58b0cb1ef8bbc37a4106`;
  concurrency baseline `1661c111326adc69291ded3393e7bc161569fcfe`;
  concurrency candidates `c6b691f95841976f24b74ac1fcbbc1ab511ecbab`.
- **Evidence destination:** `docs/workstreams/r3-reflection-dynamic-evidence/baseline-6cf3636/`
- **Production boundary:** research changes are restricted to `docs/`; production
  packages, existing tests, harnesses, evidence, CI, and Roadmap capability
  claims remain unchanged.

## Research reconciliation state

The original research content and eight gates were complete before ADR-0034
changed the governing profile boundary. Reconciliation is now complete:
Decision coverage and Slice plan are Pass, with five shared-kernel and five
optional Java SE Compatibility Profile contracts Accepted. The workstream remains **In
Progress**, not Ready or Done, until the Owner authorizes a fixed candidate
commit. A later candidate Owner review will cover:

- the fixed 24-fixture baseline and its 0/24 current catty capability result;
- the five reports in `r3-reflection-dynamic-reports/`;
- Accepted ADR-0031 through ADR-0033;
- the ten Accepted implementation contracts across the shared-kernel and
  optional Java SE compatibility tracks.

ADR-0031 through ADR-0033 were revised and Accepted by the Owner on 2026-07-18
under ADR-0034. They now govern shared metadata/invocation, InvokeDynamic
linkage, and defining-loader/generated-class kernels while classifying Java SE
reflection, annotation, lambda, and Proxy facades as optional compatibility
work. The implementation contracts now preserve that split and no optional
compatibility slice is an automatic R3 closure gate.

Accepted ADRs and contracts do not make this research Done and do not start a
production workstream. Production implementation remains unauthorized until
the current research is fixed/reviewed Done and the selected first contract has
an acceptance anchor.
