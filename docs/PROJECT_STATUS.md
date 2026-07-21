# Project status

**As of:** 2026-07-21
**Stable baseline:** R3 K2 runtime identity and typed class-definition slice complete; R2 concurrency milestone remains complete
**Baseline commit:** `100e29a` (K2 runtime identity/typed definition local closure; not merged or pushed by governance)
**Governance/research anchor:** `f685526` (R3 research Done); K2 acceptance anchor: `0fcf316`
**Active workstream:** None. The next R3 implementation successor must fix its own acceptance anchor before production work begins.
**Current phase:** R3 research is Done. K1 dynamic metadata and K2 runtime
identity/typed class-definition shared-kernel slices are complete. No
Java-visible reflection, InvokeDynamic, generated-class, or arbitrary
ClassLoader capability is claimed by K1/K2. Timed `wait`/`join`, `Unsafe`,
virtual threads, `ThreadGroup`/`ThreadLocal`, and `java.util.concurrent`
remain out of scope.

This is the single model-neutral current-state entry. Strategy lives in
[`ROADMAP.md`](./ROADMAP.md); decisions live in [`adr/`](./adr/); scoped work
lives in [`workstreams/`](./workstreams/).

## Verified capability

- Interpreter: approximately 145 opcodes, exceptions, interface dispatch,
  multidimensional arrays, and class initialization.
- Class loading: provider chain plus real `java.base` auto-detection through a
  JDK-extracted image.
- Native/bootstrap layer: the current R1 implementation has six synthetic
  bootstrap classes, additional synthetic fallbacks, and approximately 40
  native registrations; ADR-0022 does not treat that class list as permanent.
- AOT: standalone Go binary path; `fib(35)` recorded at approximately 40–60 ms
  on the development machine.
- Regression baseline: unit tests, three-engine fixture comparison, and the
  real `java.base` smoke path in CI.
- R2 initialization: bounded Java 25 single-execution-context class/interface
  initialization at `new`, resolved `getstatic`/`putstatic`, and resolved
  `invokestatic`; 13/13 differential fixtures match in Interpreter and IR.
  AOT supports the constant-field path and explicitly rejects the remaining
  tested initialization paths pending cross-engine exception propagation.
- R2 String: immutable UTF-16 code-unit backing for the bounded synthetic/native
  String surface. All eight differential fixtures match Temurin 25 in Interpreter
  and IR; AOT supports five fixtures and explicitly reports three as Not implemented.
- R2 concurrency (Slices A/B/C/D/E): race-free SC heap cells (ADR-0030),
  concurrency-safe class loading with canonical Class mirrors, stable Java
  Thread identity/lifecycle with one goroutine carrier per started platform
  Thread, VM daemon liveness, interruptible wait/join/sleep, Java object
  monitors and wait sets with notify/interrupt ordering (ADR-0029),
  `holdsLock`/`wait` argument validation, and per-Class `initMu`/`initCond`
  with JVMS §5.5 cross-context initialization protocol (other-owner wait,
  terminal publication + notify-all, unchanged interrupt status of init
  waiters). The full 19-fixture concurrency matrix matches Temurin 25 in
  Interpreter and IR (1× and race-built 100× stress); AOT reports all
  concurrency fixtures as `Not implemented`/build rejection. Timed
  `wait`/`join`, `Unsafe`, virtual threads, `ThreadGroup`/`ThreadLocal`,
  and `java.util.concurrent` remain out of scope and are governed by later
  phases.
- R3 K1 metadata kernel: classfile dynamic metadata required by later R3
  kernels is retained without adding Java-visible reflection or InvokeDynamic
  support.
- R3 K2 runtime identity and typed class-definition kernel: every runtime
  Class has defining-loader-aware identity; lookup/initiation/definition use
  typed results; concurrent definition publishes one fully linked Class or one
  terminal failure; primitive/void and array runtime identities are canonical;
  existing Class mirror paths in Interpreter and IR converge on canonical
  runtime identities. K2 adds no `Class.forName`, reflection facade,
  annotation API, arbitrary `ClassLoader.defineClass`, generated class,
  InvokeDynamic, or AOT dynamic fallback capability.

## Governance validation

K2 closure validated locally on 2026-07-21 at candidate `100e29a`:

- `go vet ./...` — Pass
- `go test ./...` — Pass
- `go test -race ./...` — Pass
- `bash tests/run.sh` — Pass, 10/10 fixtures
- R3 24-row baseline — Pass, 24/24 rows complete; Interpreter 0/24 Match,
  IR 0/24 Match, AOT 24/24 NO-BUILD
- 19-fixture R2 concurrency matrix (1×) — Pass, 19/19 Interpreter + IR Match,
  19/19 AOT NO-BUILD
- 19-fixture R2 concurrency matrix (`R2_CONCURRENCY_STRESS=100`, race-built)
  — Pass, 19/19 Match, no races
- Evidence isolation — Pass (historical evidence directories unchanged)
- Governance `git diff --check` — Pass

## Explicit boundary

catty does not claim timed `wait`/`join`, `Unsafe`/`VarHandle`, broad
reflection, Java-visible reflection facades, `invokedynamic`, generated
classes, arbitrary `ClassLoader.defineClass`, broad I/O/networking, arbitrary
`java.base` application compatibility, cross-engine AOT exception propagation,
AOT concurrency, virtual threads, `ThreadGroup`/`ThreadLocal`,
`java.util.concurrent`, or a complete Java String API. The bounded Java 25
concurrency surface (Slices A–E) is implemented in Interpreter and IR only;
AOT reports all concurrency fixtures as `Not implemented`.
`Integer/Long.toString`, `Double.parseDouble`, and representative `HashMap`
behavior remain blocked by unresolved runtime/library dependencies.

## Decision state

ADRs 0016 and 0018–0034 (excluding unused 0026) are Accepted.
ADRs 0001–0007, 0014–0015, and 0017 are superseded;
ADRs 0008–0013 are withdrawn. ADR-0017 established the earlier Java 25
supported-capability semantic baseline; ADR-0034 supersedes it with a JVMS Core
plus parallel Catty Runtime and optional Java SE Compatibility Profiles over a
typed Host ABI.
ADR-0016 fixes AOT as the primary product path with a
permanent interpreter fallback. ADR-0025 is implemented by the completed,
bounded class/interface-initialization workstream; ADR-0027 is implemented by the
completed bounded UTF-16 String workstream. ADR-0028 through ADR-0030 govern and are implemented
by the completed, bounded Thread/monitor/init Slices A–E. ADR-0031 through
ADR-0033 govern the R3 shared-kernel sequence; K1 and K2 implement the
metadata-retention and runtime-identity/typed-definition prerequisites, while
typed invocation, InvokeDynamic linkage, and generated classes remain future
workstreams. Bootstrap capability mapping, Unsafe, and allocation remain
deferred.

## Next action

Select the next R3 shared-kernel successor. The accepted
[`r3-typed-invocation-kernel-slice`](./workstreams/r3-typed-invocation-kernel-slice.md)
now has its K2 prerequisite satisfied, but implementation remains unauthorized
until its acceptance anchor is fixed in a commit. No Java-visible R3 row may be
newly claimed Supported without a later accepted workstream and evidence.
