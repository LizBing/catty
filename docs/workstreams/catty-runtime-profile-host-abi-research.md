# Catty Runtime Profile and Host ABI research

**Status:** Accepted
**Type:** research
**Review:** owner
**Proposed base:** the governance commit that fixes Accepted ADR-0034
**Acceptance anchor:** required after Owner accepts this contract
**Roadmap item:** ADR-0034 cross-cutting profile boundary; prepares Phase R3
profile reconciliation and Phase R4 host services
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0025, ADR-0027 through
ADR-0030, and ADR-0034
**Prerequisites:** ADR-0034 fixed in a governance commit; the active
`r3-reflection-dynamic-research` workstream reviewed and marked Done; no other
active workstream

## Outcome

Establish an evidence-backed v0 contract for the Catty Runtime Profile and its
typed Host ABI. The research SHALL define the minimum profile needed to build a
controlled Java rule program into a standalone or Go-embedded Catty artifact,
without making Java SE, `java.base`, JNI, or a security sandbox part of the
core product contract.

Deliver a current-runtime inventory, a profile/capability model, a Host ABI
design, a build-time validation and failure model, and an ordered set of
Proposed implementation contracts. This workstream does not add production
profile, native, launcher, embedding, or AOT capability.

## Reference use case

Use a bounded “Catty Rules” program as the design driver:

```text
Java classfiles + declared profile/capabilities
                  |
          build-time validation
                  |
       Interpreter/IR or Catty AOT
                  |
     standalone binary or Go embedding
                  |
       typed Host ABI providers
```

The reference program accepts typed scalar, String, array, and structured
inputs; performs deterministic computation; may throw a Java exception; and
may use explicitly granted logging, clock, configuration, and random services.
File, network, process, graphics, arbitrary native-library loading, and Java SE
facades are negative capability probes rather than v0 requirements.

The reference use case informs the minimum profile. It does not permanently
limit Catty to rules engines or authorize a public product/API name.

## Scope

### 1. Current runtime and native-boundary inventory

- Map every current synthetic/bootstrap class, native registration, native
  fallback/stub, launcher `java.base` decision, and Interpreter/IR/AOT native
  call path to its owning package and observable behavior.
- Classify each item as JVMS Core mechanism, Catty Runtime Profile candidate,
  Java SE Compatibility Profile behavior, host provider, temporary R1
  compatibility mechanism, or defect/implicit approximation.
- Identify all paths that expose `rtda.Frame`, `Slot`, Go pointer/layout, Go
  panic, default zero/null return, or process-global mutable state across the
  intended Host ABI boundary.
- Record the current behavior when a declared native method has no provider,
  independently for Interpreter, IR, and AOT build/run paths.

### 2. Runtime Profile v0 capability model

- Define the minimum Java-visible bootstrap and runtime surface required by the
  reference use case. Every class/method-level capability SHALL have an owner,
  profile, version, engine state, dependency closure, and failure behavior.
- Separate language/runtime necessities such as Object, Class identity,
  String, arrays, Throwable transport, initialization, and Thread context from
  optional host services and Java SE-named facades.
- Define how a program selects a profile and declares host capabilities, how
  transitive dependencies are computed, and which metadata is embedded in a
  standalone or Go-embedded artifact.
- Define profile/provider version negotiation and deterministic provider
  precedence. A provider registration SHALL NOT silently replace another
  provider or expand the program's declared capabilities.
- Determine whether the first profile declaration is a manifest, CLI/build
  option, classfile-adjacent descriptor, or a combination. The research may
  recommend a format but does not implement it.

### 3. Typed Host ABI

- Compare at least three boundary designs: generated typed bindings, a logical
  typed dynamic-value/result boundary, and the current Frame/Slot registry as
  a measured migration baseline.
- Define logical representations for void, all Java primitive values, object
  references, String/array access, structured values, Java throwable results,
  and explicit execution/Thread context without exposing engine layouts.
- Define provider identity, registration, lookup, lifecycle, concurrency,
  reentrancy, cancellation/interruption interaction, panic containment, and
  normal/exception result transport.
- Define the first synchronous ABI. Asynchronous/streaming providers may be
  analyzed for forward compatibility but SHALL NOT expand the v0 contract.
- Map adapters required by Interpreter, IR, emitted AOT code, and Go embedding.
  AOT linkage must preserve explicit capability and failure reporting; a
  built-then-panic path is not a supported fallback.
- Specify the Java-visible failure categories for absent provider, denied
  capability, signature/version mismatch, provider panic/internal failure, and
  unsupported engine path. The research must recommend exact throwable
  mappings for later Owner acceptance.

### 4. Build-time validation and artifact contract

- Define a closed-world validation model for class, method, opcode, profile,
  native/provider, dynamic-linkage, and engine dependencies.
- Establish which facts can be proven statically and which require an explicit
  runtime check or Interpreter fallback. Reflection, `invokedynamic`, generated
  classes, and late loading SHALL be classified rather than assumed absent.
- Define deterministic diagnostics and exit classifications for unsupported
  profile use, undeclared host capability, unavailable provider, invalid
  signature, and AOT-incompatible dynamic behavior.
- Define the metadata needed to reproduce an artifact's Catty version, profile
  version, engine policy, provider requirements, and capability grant set.
- Explicitly distinguish capability declaration/validation from adversarial
  sandboxing. The v0 model SHALL NOT claim security isolation unless a later
  Accepted ADR and threat-model workstream establish it.

### 5. Fixed research probes

Before baseline execution, freeze a bounded probe manifest with exactly 16
cases across these families:

| Family | Count | Required questions |
|---|---:|---|
| Core values and execution | 5 | primitives, String/UTF-16, arrays, object identity, Java exception result |
| Candidate v0 host services | 4 | logging, clock, configuration, deterministic/random provider semantics |
| Binding and failure | 4 | absent provider, denied capability, signature mismatch, provider panic containment |
| Engine/artifact boundary | 3 | Interpreter/IR parity, AOT classification, Go-embedding value/exception round trip |

Probes may use research-only mock providers and harness code under this
workstream's evidence directories. They do not become production APIs. The
manifest records source/hash, intended profile, declared grants, expected
classification, engine attempts, stdout/stderr, exit status, and timeout.
Changing the denominator after contract acceptance requires an Accepted
amendment; optional observations do not alter it.

### 6. Required reports and proposals

- `runtime-native-inventory.md` — current provider/stub/engine/launcher map and
  every Frame/Slot/panic/zero-value boundary leak.
- `runtime-profile-v0.md` — proposed capability taxonomy, minimum API/facade
  surface, dependency model, profile selection, versioning, and provider
  precedence.
- `host-abi-options.md` — candidate comparison, selected typed boundary,
  lifecycle/concurrency rules, adapters, exception transport, and rejected
  alternatives.
- `artifact-validation-contract.md` — closed-world analysis, diagnostics,
  embedded metadata, fallback classification, and non-security boundary.
- `catty-rules-reference.md` — the bounded reference program, success criteria,
  gaps, and why each proposed v0 capability is necessary.
- `profile-host-abi-slice-plan.md` — ordered production slices, prerequisites,
  profile/engine state, gates, migration strategy, and explicit non-scope.
- Any additional Proposed ADR needed for a durable decision not already fixed
  by ADR-0034, plus independently reviewable Proposed implementation
  workstreams. The first implementation proposal SHALL be the smallest slice
  that establishes a typed provider call and explicit missing-provider failure
  without changing launcher defaults or claiming the full v0 profile.

## Non-scope

- Production changes to classfile, classloader, `rtda`, native, Interpreter,
  IR, runtime, transpiler, launcher, embedding APIs, tests, CI, or packaging.
- Accepting or implementing a public Catty Runtime Profile API, Host ABI,
  manifest format, provider SDK, Catty Rules product, or Go embedding API.
- Removing current `java.base` auto-detection, deleting synthetic classes,
  migrating existing native registrations, or changing current engine
  behavior.
- Java SE Compatibility Profile implementation, complete `java.lang.*`, JCK,
  JNI, FFM API compatibility, dynamic native libraries, agents,
  instrumentation, arbitrary ClassLoader support, or module compatibility.
- File, socket, selector, process, graphics, UI, audio, database, or operating-
  system provider implementation.
- Adversarial-code sandboxing, process isolation, resource quotas, denial-of-
  service resistance, cryptographic trust, package signing, or a security
  certification claim.
- Broad performance claims. Bounded overhead measurements may compare ABI
  options but cannot select a semantically incomplete boundary.
- Rewriting the completed R3 research evidence. Reconciliation SHALL preserve
  its historical Java 25 observations and propose profile labels separately.

## Semantic and architectural constraints

- ADR-0034's JVMS Core, Catty Runtime Profile, and Java SE Compatibility
  Profiles remain parallel scopes. The Catty and Java SE profiles may share
  kernels/providers but cannot require each other's public APIs.
- Profile capability and engine capability are independent axes. Every claim
  records profile plus Interpreter/IR/AOT state as `Supported`, `Fallback`, or
  `Not implemented`.
- Java-visible values, object/Class identity, initialization, Thread context,
  monitors, heap publication, and exceptions reuse accepted shared runtime
  semantics; a provider cannot create an alternate Java object world.
- The Host ABI is typed and engine-neutral. Frame, Slot, IR register, heap-cell
  bits, Go pointer layout, and emitted AOT layout are adapter details.
- Unsupported or unauthorized behavior fails explicitly. A Go panic, process
  crash, silent provider selection, default zero/null result, or built-then-
  panic AOT artifact is never a supported outcome.
- Research mock providers and prototypes remain isolated evidence. Passing a
  probe cannot promote them into production.
- Repository evidence outranks the reference-product narrative. Every proposed
  v0 capability must be justified by a fixed probe, an irreducible runtime
  dependency, or an explicit forward-compatibility constraint.

## Acceptance

| Gate | Required evidence | Result |
|---|---|---|
| Fixed probe manifest | Exactly 16 reviewed probes with hashes, declared profile/grants, bounded timeouts, and expected classifications | Not run |
| Current inventory | Complete synthetic/native/stub/launcher and three-engine call-path map; all Frame/Slot/panic/zero-value leaks classified | Not run |
| Profile v0 contract | Minimum capability graph, API/facade ownership, selection/versioning/provider-precedence rules, and parallel-profile boundary | Not run |
| Host ABI decision evidence | At least three candidates compared; typed value/result, lifecycle, concurrency, panic, failure, and per-engine adapters specified | Not run |
| Artifact validation | Static/runtime/fallback boundary, deterministic diagnostics, reproducibility metadata, and explicit non-security statement | Not run |
| Probe baseline | All 16 probes attempted on applicable current/research paths; no omitted engine/provider row or unbounded process | Not run |
| Reference-use-case closure | Every proposed v0 capability traced to the Catty Rules reference, an irreducible runtime dependency, or a documented forward constraint | Not run |
| Decision and slice coverage | Required Proposed ADRs plus ordered Proposed implementation contracts; first slice is independently acceptable and explicit-failure bounded | Not run |
| Governance | Research-only changes isolated; historical evidence unchanged; `git diff --check` passes | Not run |

Results use only `Pass`, `Fail`, `Not run`, or `Not implemented`. Current Catty
failures are valid baseline data, but missing rows, fabricated compatibility,
or silently skipped paths fail the research gate.

## Plan

| Step | State |
|---|---|
| Owner reviews and accepts research contract | Complete |
| Fix acceptance anchor and record research preflight | Pending |
| Freeze 16 probes and capture current boundary behavior | Pending |
| Produce runtime/native inventory and profile capability graph | Pending |
| Compare Host ABI candidates and run bounded research prototypes | Pending |
| Define artifact validation and reference-use-case closure | Pending |
| Draft required ADRs and ordered implementation contracts | Pending |
| Owner reviews research conclusions | Pending |

Status uses `Pending`, `In progress`, or `Complete`.

## Amendments

None.

## Acceptance record

Accepted by Owner on 2026-07-18. This acceptance freezes Outcome, Scope,
Non-scope, Semantic and architectural constraints, the 16-probe denominator,
Acceptance gates, and owner review. It authorizes research deliverables only
after the prerequisites are satisfied and this record is fixed in an
acceptance-anchor commit; it does not authorize production implementation or
other integration actions.
