# R3 typed dynamic-invocation kernel slice

**Status:** Done
**Type:** implementation
**Review:** owner
**Profile:** Catty JVMS Core shared kernel; reusable by profiles and Host ABI adapters
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0021, ADR-0025,
ADR-0028 through ADR-0031, ADR-0033, and ADR-0034
**Prerequisites:** `r3-runtime-identity-definition-slice` Done; acceptance
anchor fixed at `6563045`
**Acceptance anchor:** `65630453cdb7579c85f602b21c2d9f459a042f88`; prerequisites
include the completed execution-context and Java Thread facade boundary slice
at the same commit

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

In Progress; prerequisites and acceptance anchor are fixed at `6563045`.
Implementation preflight is recorded at
[`r3-typed-invocation-preflight-20260725/preflight-review.md`](./r3-typed-invocation-preflight-20260725/preflight-review.md).
Owner accepted the preflight and authorized implementation start on 2026-07-25.
The first implementation slice is the typed value/result vocabulary and its
frame/heap adapters; no Java-visible capability is claimed until acceptance
evidence is recorded.

### First implementation record

On 2026-07-25, the candidate introduced the shared `rtda.JavaValue` and
`rtda.DynamicResult` vocabulary plus descriptor-aware Frame-local and typed
HeapCell adapters. It also added resolved instance/static typed field access
with exact descriptor-kind validation; class initialization remains owned by a
later ExecutionContext-aware direct-invocation service. The adapter unit matrix
covers void, every primitive, category-2 values, reference/null identity,
normal/throwable/internal-failure states, field receiver/value validation, and
containment of legacy adapter panics. At that point it did not yet add direct
method/constructor dispatch, Interpreter/IR invocation acceptance, or AOT
dynamic invocation. Local `go test ./rtda` and `go test ./...` passed for this
candidate; full workstream acceptance evidence remains pending.

### Direct invocation implementation record

The candidate now defines `rtda.InvocationRequest` (execution context, caller,
resolved method, receiver, and logical arguments) and an Interpreter adapter
for direct native and bytecode method execution. It validates exact descriptor
kinds and receiver/arity before writing Frame locals; it carries normal return
values, including category-2 values, through `DynamicResult`, and returns an
escaping Java throwable with its original object identity. The initial adapter
required an empty execution-context stack; it now captures returns and abrupt
completion at an explicit entry depth, preserving any ordinary Java caller
frame below the typed boundary. IR direct invocation, Java-level direct-access
failures, class-initialization triggers, and AOT dynamic invocation remain
pending.

### Dispatch and IR implementation record

The candidate now provides matching tree-Interpreter and IR direct adapters.
Resolved direct dispatch covers static, virtual, interface, special, and
constructor selection; static calls initialize the actual declarer's Class
before invocation. Direct field wrappers apply the corresponding static-field
initialization trigger before using the typed field boundary. A null direct
receiver is mapped to a Java `NullPointerException` result when the configured
loader can construct it; malformed adapter input remains an internal failure.
The unit matrix exercises normal and abrupt Interpreter/IR execution,
category-2 transport, dispatch selection, constructor receiver identity,
static initialization, and null-receiver throwable transport. It does not yet
prove broad access checks, every linkage failure mapping, static-field storage
through a full classfile fixture, or nested typed calls. Direct reference-field
assignment now resolves the descriptor target through the typed loader service:
an incompatible non-null object yields `ClassCastException`, and load failure
yields the corresponding linkage throwable. The AOT bridge exposes an explicit
typed-dynamic unsupported result and does not delegate this new boundary to its
legacy Slot fallback.

### Nested invocation boundary record

Typed direct invocation now permits an existing Java caller stack. Each direct
call installs a LIFO result capture at its entry depth; returns crossing that
depth become `JavaValue` results, while ordinary bytecode calls above it still
transfer through their existing Frame/Slot path. Abrupt completion unwinds only
to that boundary and preserves the outer caller frame. Tree Interpreter and IR
tests cover nested normal category-2 return and nested uncaught throwable
identity. This does not convert the legacy bytecode invocation path into a
stable Slot API.

### Caller access record

Dispatch and direct-field adapters now enforce caller-aware public, private,
package-private, and protected access checks. Public members remain callable
without a caller identity; non-public field consumers use the explicit
`ReadDirectFieldFrom`/`WriteDirectFieldFrom` variants. Denied access yields a
Java `IllegalAccessError`. The protected rule includes the cross-package
receiver constraint for instance calls. The current tests cover method access
for each visibility and private-field denial; full classfile fixture coverage
for static typed fields is recorded as composite evidence at
[`r3-typed-invocation-evidence-20260725/acceptance-candidate.md`](./r3-typed-invocation-evidence-20260725/acceptance-candidate.md).
Local candidate validation on 2026-07-25 passed `go vet ./...`, `go test ./...`,
`go test -race ./...`, `bash tests/run.sh` (10/10 fixtures), and
`git diff --check`; this is implementation evidence, not final workstream
acceptance.

### Final-gate record

Candidate `a7c94eb8e81fcd613341d2fab5e29c1464576eeb` fixes the current
implementation scope. Its complete gate record is
[`final-gates-a7c94eb.md`](./r3-typed-invocation-evidence-20260725/final-gates-a7c94eb.md).
All gates pass, including the required race-built R2 100× matrix (19/19
Interpreter match, 19/19 IR match, 19/19 AOT NO-BUILD). This workstream is
Ready for Owner review. Per COLLABORATION.md §6.1, only the Owner can accept
the review and mark this workstream Done.

### Owner review correction record

On 2026-07-26, Owner accepted the existing gate evidence but requested changes
before Done. The workstream returned to In Progress for three corrections
within the frozen contract: validate non-null reference arguments against their
declared parameter descriptors, restore the typed invocation entry depth after
an IR internal failure, and restrict constructor resolution to constructors
declared by the target Class. No ADR or new workstream is required.

### Corrected candidate record

Candidate `047ffca5b0126262865e13884dc48528b74a2985` fixes the three Owner-requested
corrections. Its complete final-gate evidence is recorded at
[`final-gates-047ffca.md`](./r3-typed-invocation-evidence-20260725/final-gates-047ffca.md).
All gates pass: `go vet ./...`, `go test ./...`, `go test -race ./...`,
`bash tests/run.sh` (10/10 fixtures), R2 concurrency 1× (19/19 Interpreter + IR
Match, 19/19 AOT NO-BUILD), R2 concurrency 100× race-built (19/19 Interpreter +
IR Match, 19/19 AOT NO-BUILD), evidence isolation, and `git diff --check`.

Per COLLABORATION.md §6.1, Owner accepted the review and marked this workstream
Done on 2026-07-26. The integration commit is recorded in the acceptance record.

## Preflight review

Recorded on 2026-07-25 at actual review base `cc1da5d`, a descendant of the
fixed acceptance anchor `6563045`. The review found the implementation may
start after explicit Owner authorization, with the first slice focused on the
typed value/result vocabulary and adapters. It also recorded two boundary
risks: the existing AOT bridge still exposes `[]rtda.Slot`, and missing
classfile natives still use zero/null default stubs. The former is in scope as
adapter debt for typed invocation; the latter is adjacent Host ABI/native
binding debt and is not authorized for this workstream.

## Acceptance record

Accepted by Owner on 2026-07-18. Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, and owner review are
frozen. Acceptance anchor fixed by Owner on 2026-07-25 at `6563045` after the
execution-context and Java Thread facade boundary prerequisite was closed.
Implementation start explicitly accepted by Owner on 2026-07-25 under this
contract.

Integration accepted by Owner on 2026-07-26 after candidate
`047ffca5b0126262865e13884dc48528b74a2985` passed the complete final-gate set.
