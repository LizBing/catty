# Project status

**As of:** 2026-07-15
**Stable baseline:** R2 initialization and bounded UTF-16 String slices complete
**Baseline commit:** `8171361` (integration; String candidate `00327d6`, evidence `9008b00`)
**Active workstream:** Accepted [`r2-thread-monitor-foundation-slice`](./workstreams/r2-thread-monitor-foundation-slice.md)
**Current phase:** R2 Slice A, B, and C accepted (Slice C monitors/wait sets/interruption complete on `eea253d`); Slice D (concurrent ADR-0025 initialization and full 19-fixture Interpreter/IR matrix) is the next bounded implementation slice

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
- R2 concurrency (Slices A/B/C): race-free SC heap cells (ADR-0030),
  concurrency-safe class loading with canonical Class mirrors, stable Java
  Thread identity/lifecycle with one goroutine carrier per started platform
  Thread, VM daemon liveness, interruptible wait/join/sleep, Java object
  monitors and wait sets with notify/interrupt ordering (ADR-0029), and
  `holdsLock`/`wait` argument validation. The bounded 11-fixture Slice C matrix
  matches Temurin 25 in Interpreter and IR (1× and race-built 20× stress);
  AOT reports all concurrency fixtures as `Not implemented`/build rejection.
  Cross-thread class initialization, timed `wait`/`join`, `Unsafe`, virtual
  threads, `ThreadGroup`/`ThreadLocal`, and `java.util.concurrent` remain out
  of scope and are governed by later slices.

## Governance-reset validation

Revalidated locally on 2026-07-13:

- `go vet ./...` — Pass
- `go test ./...` — Pass
- `go test -race ./...` — Pass
- `bash tests/run.sh` — Pass, 10/10 fixtures

## Explicit boundary

catty does not claim timed `wait`/`join`, `Unsafe`/`VarHandle`, broad
reflection, `invokedynamic`, broad I/O/networking, arbitrary `java.base`
application compatibility, cross-thread class initialization behavior,
cross-engine AOT exception propagation, AOT concurrency, virtual threads,
`ThreadGroup`/`ThreadLocal`, `java.util.concurrent`, or a complete Java String
API. The bounded Java 25 concurrency surface (Slices A/B/C) is implemented in
Interpreter and IR only; AOT reports all concurrency fixtures as `Not
implemented`. `Integer/Long.toString`, `Double.parseDouble`, and representative
`HashMap` behavior remain blocked by unresolved runtime/library dependencies.

## Decision state

ADRs 0016–0030 (excluding unused 0026) are Accepted. ADRs 0001–0007 and 0014–0015 are superseded;
ADRs 0008–0013 are withdrawn. ADR-0017 fixes Java 25 as the supported-capability
semantic baseline; ADR-0016 fixes AOT as the primary product path with a
permanent interpreter fallback. ADR-0025 is implemented by the completed,
bounded class/interface-initialization workstream; ADR-0027 is implemented by the
completed bounded UTF-16 String workstream. ADR-0028 through ADR-0030 now
govern the bounded Thread/monitor/class-init/heap direction. ADR-0028 through
ADR-0030 are implemented by the completed, bounded Thread/monitor SC-heap
Slices A/B/C. Bootstrap capability mapping, Unsafe, and allocation remain
deferred. The Proposed Thread/monitor implementation contract is accepted. Its
acceptance-anchor commit is the required base before production implementation
begins.

## Next action

Start the next Active Agent for Slice D from the accepted Slice C candidate
lineage (`eea253d`). Slice C added Java object monitors, synchronized methods,
wait sets, and interruption with `holdsLock`/`wait` argument validation; the
runtime now supports the bounded 11-fixture Slice C matrix in Interpreter and
IR (race-built stress verified) and AOT remains `Not implemented` for those
fixtures. Slice D must preserve the accepted Slice A/B/C evidence and draft a
bounded contract for concurrent ADR-0025 initialization and the full 19-fixture
Interpreter/IR matrix before production implementation begins.
