# Project status

**As of:** 2026-07-20
**Stable baseline:** R3 K1 dynamic metadata kernel complete; R2 concurrency milestone remains complete
**Baseline commit:** `0fcf316` (K1 dynamic metadata kernel integration)
**Governance/research anchor:** `f685526` (R3 research Done); K2 acceptance anchor: `0fcf316`
**Active workstream:** `r3-runtime-identity-definition-slice` — Active on the next implementation branch from anchor `0fcf316`; the only active implementation workstream
**Current phase:** R3 research is Done and K1 dynamic-metadata kernel is integrated. K1 adds no Java-visible reflection or InvokeDynamic capability by contract. K2 is now the active shared-kernel workstream. Timed `wait`/`join`, `Unsafe`, virtual threads, `ThreadGroup`/`ThreadLocal`, and `java.util.concurrent` remain out of scope.

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

## Governance-reset validation

Revalidated locally on 2026-07-17 (Slice E final candidate `ea1f67a`, integrated
as `ca42a61`):

- `go vet ./...` — Pass
- `go test ./...` — Pass
- `go test -race ./...` — Pass
- `bash tests/run.sh` — Pass, 10/10 fixtures
- 19-fixture concurrency matrix (1×) — Pass, 19/19 Interpreter + IR Match, 19/19 AOT NO-BUILD
- 19-fixture concurrency matrix (`R2_CONCURRENCY_STRESS=100`, race-built) — Pass, 19/19 Match, no races
- Evidence isolation — Pass (historical evidence directories unchanged)
- Governance `git diff --check` (scoped per Amendment 1 precedent) — Pass

## Explicit boundary

catty does not claim timed `wait`/`join`, `Unsafe`/`VarHandle`, broad
reflection, `invokedynamic`, broad I/O/networking, arbitrary `java.base`
application compatibility, cross-engine AOT exception propagation, AOT
concurrency, virtual threads, `ThreadGroup`/`ThreadLocal`,
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
by the completed, bounded Thread/monitor/init Slices A–E. Bootstrap capability
mapping, Unsafe, and allocation remain deferred.

## Next action

Implement the Accepted
[`r3-runtime-identity-definition-slice`](./workstreams/r3-runtime-identity-definition-slice.md)
K2 contract from anchor `0fcf316`: defining-loader-aware canonical Class
identity, typed lookup/definition results, atomic publication, and mirror
continuity. K1 evidence remains fixed at 0/24 Match in Interpreter and IR with
24/24 AOT NO-BUILD; no Java-visible R3 row may be newly claimed Supported.
