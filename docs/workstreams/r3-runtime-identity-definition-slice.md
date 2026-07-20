# R3 runtime identity and typed class-definition slice

**Status:** In Progress
**Type:** implementation
**Review:** owner
**Profile:** Catty JVMS Core shared kernel; no public ClassLoader API
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0022, ADR-0025,
ADR-0028 through ADR-0031, ADR-0033, and ADR-0034
**Prerequisites:** `r3-reflection-dynamic-research` Done;
`r3-metadata-slice` Done; acceptance anchor fixed at `0fcf316`
**Acceptance anchor / actual base:** `0fcf316` / `397b407`
**Branch:** `codex/r3-runtime-identity-definition-slice`

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

| Step | State |
|---|---|
| Owner accepts frozen contract | Complete |
| K1 prerequisite and acceptance anchor | Complete — `0fcf316` |
| Implementation preflight | Complete — selected as the sole active K2 workstream |
| Map current loader, identity, mirror, and failure paths | Complete |
| Fix K2 service boundaries, implementation order, and test matrix | Complete |
| Implement typed lookup/definition and atomic publication | Complete |
| Implement canonical runtime types and mirror continuity | Complete |
| Migrate Java-reachable consumers and run contract gates | Complete |
| Runtime identity and typed lookup/definition implementation | Ready |
| Owner review and integration | Pending |

## Acceptance record

Accepted by Owner on 2026-07-18. Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, and owner review are
frozen. The prerequisite is now satisfied by K1 integration at `0fcf316`; K2
is the selected active workstream, with implementation to begin from that
anchor.

## Implementation investigation — 2026-07-20

This investigation refines the implementation order without changing the
frozen Outcome, Scope, Non-scope, Semantic constraints, Acceptance gates,
profile classification, or owner review requirement.

### Existing path

The current production path is:

```text
Loader.LoadClass(name) *Class
  -> initiating loader cache lookup by name
  -> ordered ClassProvider chain
  -> rtda.NewClass / NewArrayClass / native synthetic builder
  -> recursive LoadClass for superclass, interfaces, and array component
  -> first completed candidate wins the cache entry
```

`rtda.Loader` exposes only the panicking `LoadClass` method. `ClassProvider`
returns either a Class or nil, so nil currently conflates provider miss,
classpath I/O failure, malformed classfile, and dependency failure. Runtime
`Class` retains no defining-loader identity.

The current cache is concurrency-safe but not a definition state machine.
Concurrent misses may all execute the provider and construct separate linked
Classes before the final cache double-check returns one winner. The existing
concurrency test supplies the same prebuilt pointer to every caller and
therefore does not detect duplicate construction, duplicated side effects, or
partial-link work.

### Observed gaps

| Area | Current behavior | K2 requirement |
|---|---|---|
| Loader identity | Class identity is inferred from a loader-local name cache | Class stores immutable defining-loader identity plus binary name |
| Initiation vs definition | One cache represents both concepts | Initiating lookup may cache a Class defined by another loader without changing identity |
| Provider failures | nil discards parse/I/O/dependency cause | miss, format, circularity, duplicate definition, and linkage failure remain typed |
| Publication | providers build outside the lock; first completed candidate wins | one definition attempt publishes one fully linked Class or one terminal failure |
| Recursive dependencies | nested `LoadClass` has no explicit load context | same-chain circularity is detected without stack overflow; waiters cannot observe partial state |
| Primitive/void identity | primitive lookup is mapped to array names; void maps to unsupported `V` | nine canonical VM type identities independent of initiating loader |
| Array identity | each initiating loader builds/cache-indexes an array by descriptor | reference arrays derive defining identity from the component Class; primitive arrays derive from VM primitive identity |
| Class mirrors | `ClassObject(factory)` is canonical only if every producer uses it | Class owns mirror materialization; all Interpreter/IR/native producers converge on it |
| Native Class paths | superclass and primitive helpers allocate `Class` objects directly | `getSuperclass` and primitive lookup return the canonical mirror |
| Java-reachable failure | missing class normally panics in Go | Interpreter/IR adapters translate typed lookup/linkage failures to Java throwable transport |
| Classpath errors | `CompositeEntry` converts every entry error into not-found | a real miss is distinguishable from unreadable/corrupt classpath input |

### Chosen service boundary

The typed vocabulary belongs in `rtda`, because `rtda` owns the cycle-breaking
Loader interface and cannot import `classloader`. K2 will add:

- an opaque, comparable `LoaderIdentity` allocated once for each defining
  loader, plus the canonical VM identity used by primitives and void;
- a `ClassLoadResult` carrying exactly one fully linked Class or one immutable
  `ClassLoadFailure`;
- failure kinds for not-found, format, circularity, duplicate definition, and
  dependency/linkage failure, retaining the symbolic class name and cause;
- a typed lookup method on `Loader`, while preserving a must-load convenience
  method only for bootstrap invariants and legacy callers proven unreachable
  from supported classfiles.

`classloader` owns lookup/initiation, definition state, provider selection, and
publication. Providers will return an explicit miss, definition candidate, or
typed failure instead of nil. A delegated result that is already bound to a
different defining loader is recorded only in the initiating cache; it is not
rebound or duplicated.

### Atomic definition protocol

Each loader will maintain separate initiating-cache and definition-state
records. A definition record has unresolved, defining, defined, or failed
state and publishes its terminal Class/failure to all waiters.

The first implementation will serialize first-time definition work per loader
and pass an explicit internal load session through the provider-facing Loader
view. Nested superclass/interface/component lookup reuses that session. This
gives deterministic same-chain circularity detection and avoids cross-name
lock cycles without relying on goroutine identity. Already defined and
initiated cache hits remain concurrent reads.

No Class enters either public cache until superclass/interfaces, field layout,
methods, native binding, defining identity, and required structural linkage
are complete. A failed definition publishes no partial Class and wakes every
waiter with the same terminal failure category.

### Runtime identity and mirrors

`rtda.Class` will retain its defining Loader and `LoaderIdentity`. One-time
binding is permitted during construction; rebinding after definition is an
invariant failure. Equality remains pointer identity, with tests proving the
pointer is canonical for `(defining loader, binary name)`.

Primitive and void Classes will be VM-canonical singletons. Reference array
Class creation will be owned by the canonical component Class, using its
existing unused `arrayClass` cache as a race-free publication point. Primitive
array Classes will be owned by the corresponding VM primitive identity.
Initiating loaders may cache these results but cannot redefine them.

Mirror allocation will move behind one Class-owned operation that selects the
canonical `java/lang/Class` Class through the defining context. Interpreter
class literals, `Object.getClass`, synchronized static methods, native
`getSuperclass`, primitive lookup, and IR paths will use that operation. The
existing caller-supplied factory API will be reduced to internal/test support
or removed after migration so a caller cannot accidentally manufacture a
second mirror path.

### Failure transport boundary

Class loading and linking first return semantic failures, not Java facade
objects. Interpreter and IR share a small adapter that maps Java-reachable
symbolic-resolution failures to the existing pending-exception transport.
Internal bootstrap must-load helpers may still panic only when a Catty
invariant is broken. K2 does not add `Class.forName` or its
`ClassNotFoundException` policy.

The AOT bridge keeps its existing contract and unsupported classifications in
this slice. K2 will not turn an unsupported dynamic-loading path into a
built-then-panic fallback.

### File-level implementation order

| Order | Area | Intended change |
|---|---|---|
| 1 | `rtda/thread.go`, focused new loader-result file, `rtda/class.go` | Add LoaderIdentity, typed result/failure, defining-loader ownership, and immutable accessors |
| 2 | `classpath/entry*.go` | Preserve a typed classpath miss while retaining real I/O/corruption errors |
| 3 | `classloader/classloader.go` | Add explicit provider results, initiating/definition state, load sessions, atomic terminal publication, and must-load wrapper |
| 4 | `rtda/build.go`, native synthetic builders | Propagate typed superclass/interface/component failures and bind definitions exactly once |
| 5 | primitive/array runtime identity path | Add VM primitive/void identities and component-owned canonical array Classes |
| 6 | `rtda.Class` mirror service plus Interpreter/IR/native callers | Remove direct Class facade allocation and prove mirror continuity |
| 7 | Interpreter/IR resolution adapters | Convert Java-reachable typed failures to shared throwable transport; preserve invariant panics only |
| 8 | classloader/rtda/engine tests and repository gates | Run focused identity/failure/concurrency tests, fixed R3 baseline, R2 regression, then all contract gates |

### Test matrix

| Case | Required observation |
|---|---|
| Same loader/name | repeated and concurrent lookup returns one Class pointer and provider/build count is exactly one |
| Different defining loaders | same binary name produces distinct Classes and mirrors |
| Delegation/initiation | child lookup returns the parent's Class identity; child initiation cache does not redefine or rebind it |
| Duplicate definition | second definition returns one typed duplicate failure and leaves the original Class unchanged |
| Missing/malformed dependency | typed terminal failure, no panic, no partial cache entry, and all waiters observe the same category |
| Circular hierarchy | deterministic typed circularity/linkage failure without stack overflow or deadlock |
| Primitive/void | boolean, byte, char, short, int, long, float, double, and void are canonical across initiating loaders |
| Arrays | repeated, concurrent, multidimensional, reference, and primitive array lookup follows component/VM identity rules |
| Mirrors | class literal, object Class, superclass, primitive, array, synchronized-static, Interpreter, and IR paths return one mirror per Class |
| Initialization continuity | canonical identity still owns one ADR-0025/0029 initialization state and Class monitor behavior |
| Capability honesty | fixed R3 Class/annotation/reflection rows remain unsupported; AOT classification does not improve |

Focused synthetic providers remain the primary negative/concurrency corpus
because they can count definition attempts, block exact phases, inject typed
failures, and model delegation deterministically. Existing Java fixtures remain
the capability-honesty and mirror-continuity evidence; they do not expand the
K2 public profile.

### Principal risks and controls

- **Loader-interface blast radius:** many tests provide small Loader mocks.
  Introduce the typed vocabulary first and migrate mocks mechanically before
  changing production resolution paths.
- **Bootstrap recursion:** `java/lang/Class` mirror materialization can recurse
  through synthetic Object/Class loading. Keep Class construction separate
  from mirror construction and test the Class mirror of `java/lang/Class`.
- **Definition deadlock:** nested and concurrent definitions cannot rely on
  goroutine identity. Use explicit load sessions and terminal wakeup tests with
  timeouts.
- **Over-broad Java SE policy:** K2 returns semantic loading/linkage failures;
  it does not implement `Class.forName`, arbitrary ClassLoader, module/package,
  or Java reflection behavior.
- **AOT regression:** retain the current exact build rejection. Typed
  Thread-aware cross-engine loading/fallback remains outside K2.
- **Historical evidence mutation:** new K2 evidence uses an isolated directory;
  K1 and R2 evidence stays unchanged.

### Investigation baseline

Before production changes, the focused suites pass:
`go test ./classloader ./rtda ./interpreter ./native ./runtime`. Interpreter
and runtime currently have no package-local tests, so K2 must add focused
engine continuity coverage rather than treating their empty suites as evidence.

## K2 implementation evidence (2026-07-20)

### Design summary

**Defining-loader-aware identity.** `rtda.LoaderIdentity` is an opaque
pointer-comparable type produced by `NewLoaderIdentity()`. `VMIdentity`
(singleton, id=0) owns the nine canonical VM types. Every `rtda.Class` carries
an immutable `definingLoader` field set exactly once by `BindLoader`.

**Separate initiation vs definition caches.** `ClassLoader.initiatingCache`
maps name→Class for fast lookup (may return delegated Classes). `defRecords`
maps name→*defRecord for definition state owned by THIS loader. Delegated
classes enter only the initiating cache.

**Atomic definition protocol.** `defMu` serialises first-time definition per
loader. `defRecord` publishes one terminal state (Defined or Failed) via
`sync.Cond.Broadcast` to all concurrent waiters. No goroutine identity
dependency.

**Load sessions for circularity detection.** A `loadSession` with a `seen`
map tracks in-flight names within one top-level `LoadClassResult` call.
Session-aware loader views (`sessionLoader`) are threaded through recursive
provider calls. Same-name circularity is caught deterministically without
stack overflow or deadlock.

**Component-owned array identity.** `Class.GetArrayClass()` creates and
publishes the array Class via atomic CAS. Reference arrays inherit the
component's defining loader; primitive arrays use `VMIdentity`.

**VM canonical primitives/void.** `rtda.InitVMTypes()` creates and publishes
nine singleton Classes (`boolean`, `byte`, `char`, `short`, `int`, `long`,
`float`, `double`, `void`) owned by `VMIdentity`. `VMPrimitiveForName` and
`IsVMPrimitive` intercept lookups before the provider chain.

**Typed failure vocabulary.** `FailureKind` enum: `NotFound`, `Format`,
`Circularity`, `DuplicateDefinition`, `Linkage`. `ClassLoadResult` is a
sum type (Class or Failure). `classpath.ErrNotFound` distinguishes true
misses from I/O errors.

**Mirror continuity.** All mirror production paths converge on
`Class.ClassObject()` with CAS publication. Paths verified: interpreter
`getClassObject`, native `getClassObject`, `classGetSuperclass`,
`classGetPrimitiveClass`, `objectGetClass`, static synchronized monitor,
IR `ldc` class.

**Interpreter/IR typed failure adapter.** `resolveClass(thread, pc, name)`
calls `LoadClassResult` and maps `FailureKind` to Java throwables:
`NoClassDefFoundError`, `ClassFormatError`, `ClassCircularityError`, or
`LinkageError`. All bytecode/IR handlers use this adapter. Bootstrap
invariant paths (String, Class, exception types) keep `LoadClass` panics.

### Files

| File | Change |
|---|---|
| `rtda/loader_identity.go` | New — LoaderIdentity type, VMIdentity singleton |
| `rtda/class_load_result.go` | New — FailureKind, ClassLoadFailure, ClassLoadResult |
| `rtda/vm_types.go` | New — canonical VM primitives, component-owned GetArrayClass |
| `rtda/func_loader.go` | New — test helper implementing full Loader interface |
| `interpreter/resolve.go` | New — resolveClass adapter, FailureKind→exception mapping |
| `rtda/class.go` | Modified — definingLoader field, BindLoader, arrayClass CAS |
| `rtda/thread.go` | Modified — Loader interface adds LoadClassResult + LoaderIdentity |
| `rtda/build.go` | Modified — typed failure propagation, component-owned arrays |
| `classloader/classloader.go` | Rewritten — typed providers, defRecord, loadSession, atomic publish |
| `classpath/entry.go` | Rewritten — typed ErrNotFound, errors.As check in CompositeEntry |
| `classpath/entry_dir.go` | Modified — os.IsNotExist → typed ErrNotFound |
| `classpath/entry_zip.go` | Modified — zip miss → typed ErrNotFound |
| `native/system.go` | Modified — classGetSuperclass/classGetPrimitiveClass use ClassObject |
| `native/exceptions.go` | Modified — added ClassFormatError, ClassCircularityError synthetics |
| `interpreter/interpreter.go` | Modified — all Java-reachable LoadClass→resolveClass |
| `interpreter/ir.go` | Modified — all Java-reachable LoadClass→resolveClass |
| `interpreter/invoke.go` | Modified — pushConstant uses resolveClass, returns bool |
| `interpreter/helpers.go` | Modified — newRefArray uses resolveClass, accepts pc |
| `launch/launch.go` | Modified — InitVMTypes before loader creation |
| `runtime/runtime.go` | Modified — InitVMTypes in Bootstrap |
| `transpile/build.go` | Modified — InitVMTypes in BuildProgram |
| `classloader/classloader_test.go` | Rewritten — countProvider verifying atomic definition |
| `rtda/frame_test.go` | Modified — testLoader implements full Loader interface |
| `rtda/bootstrap_test.go` | Modified — recordingMockLoader implements full Loader interface |
| `native/thread_test.go` | Modified — simpleLoader implements full Loader interface |

### Gate results

| Gate | Command | Result |
|---|---|---|
| gofmt | `gofmt -l .` | Clean (no unformatted files in K2 scope) |
| go vet | `go vet ./...` | Clean (no warnings) |
| go test | `go test ./...` | All pass |
| go test -race | `go test -race ./...` | All pass (no data races) |
| e2e fixtures | `bash tests/run.sh` | 10/10 pass |
| git diff --check | `git diff --check` | Clean (no whitespace errors) |

### Test coverage

| Package | Test |
|---|---|
| `classloader` | `TestConcurrentLoadSingleIdentity` — 32 goroutines, countProvider verifies single Provide invocation, same pointer returned to all waiters |
| `classloader` | Additional identity and failure tests retained from K1 |
| `rtda` | All existing frame, bootstrap, and monitor tests pass with updated Loader interface |
| `native` | All existing native method tests pass with updated Loader interface |
| `lowering` | All existing lowering tests pass |
| `transpile` | All existing emit tests pass |
| e2e | HelloWorld, Fibonacci, Factorial, ArraySum, OOPDemo, StaticFields, SwitchDemo, EmptyMain, ExceptionTest, InterfaceTest — all pass |

### Capability classification (unchanged)

K2 does NOT implement: `Class.forName`, reflection facade, annotation API,
arbitrary `ClassLoader.defineClass`, generated classes, `InvokeDynamic`,
or AOT dynamic fallback. No new Java-visible R3 Supported rows are declared.
The typed failure vocabulary is an internal Go contract; it is not exposed
to Java code through new public API.

### Remaining risks

1. **Catch-type resolution in exception handlers** (`entry.CatchType()`) still
   uses `LoadClass` (must-load). If a class file's exception table references
   an unresolvable catch type, the interpreter panics rather than skipping
   the handler. This follows JVMS §5.4.3 semantics and is an existing behaviour;
   converting it to typed resolution is deferred.
2. **No ClassLoader API for Java.** The identity model, typed failures, and
   definition protocol are Go-internal. Java classes cannot yet implement or
   extend `ClassLoader`. This is explicit Non-scope.
3. **AOT build-time rejection unchanged.** Transpile still rejects any class
   with a `<clinit>` trigger. Cross-engine exception propagation is a separate
   future workstream.
4. **Interpreter/runtime packages have no package-local tests.** Regression
   coverage relies on the e2e fixture suite. Focused engine continuity tests
   for the typed failure paths (e.g., `resolveClass` with a missing class
   injected via a test provider) remain a follow-on.

### Status

K2 implementation is **Ready** for owner review. All contract gates pass.
A candidate commit has been created locally.

**未 merge、未 push、未发布。** No remote, main-branch, or release operation
was performed. The work remains on branch `codex/r3-runtime-identity-definition-slice`
as a local candidate only.
