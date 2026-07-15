# R2 Thread and monitor foundation slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Base commit:** `63d5658`
**Roadmap item:** Phase R2 — multi-threaded producer-consumer milestone
**Governing ADRs:** ADR-0016 through ADR-0025, ADR-0027, and ADR-0028 through ADR-0030
**Prerequisites:** `r2-concurrency-semantics-research` Done; ADR-0028, ADR-0029, and ADR-0030 Accepted
**Acceptance anchor:** this 2026-07-14 governance commit; implementation preflight records its resolved SHA and the actual worktree base

## Outcome

Catty runs the fixed Java 25 R2 concurrency matrix, including the one-slot
producer-consumer milestone, with Temurin-matching behavior in Interpreter and IR. The
runtime has stable Java Thread identity/lifecycle, race-free heap storage, Java object
monitors and wait sets, interruptible wait/join/sleep, VM daemon liveness, canonical Class
mirrors, and synchronized cross-thread class initialization. AOT rejects every fixed
concurrency fixture at build time as `Not implemented` until its explicit execution-
context ABI is implemented.

## Scope

- Replace Java heap use of frame `Slot` storage with ADR-0030 race-free SC heap cells for
  instance/static fields and arrays. Migrate every Interpreter, IR, runtime, native, and
  bulk-operation producer/consumer; no mutable heap slice/cell escapes. Preserve current
  primitive/reference behavior and object identity.
- Make the classloader cache safe for concurrent, recursive loading while preserving one
  runtime Class identity per name/loader. Add one canonical Java Class mirror per runtime
  Class identity.
- Attach one stable synthetic Thread facade object to each `rtda.Thread` execution
  context. Make execution-context IDs race-free and keep each frame stack, bridge return,
  and pending exception confined to its Java Thread.
- Implement the bounded synthetic Thread surface used by the fixed fixtures:
  `currentThread`, `start`, `run` dispatch, `isAlive`, `join()`, `interrupt`,
  `isInterrupted`, `interrupted`, `sleep(long)`, `onSpinWait`, `holdsLock`, `setDaemon`,
  and `isDaemon`, plus start-once and argument/state failures needed by those methods.
- Use one goroutine carrier per successfully started supported platform Thread. Add VM
  supervision so started non-daemon Threads keep the launcher alive and daemon Threads do
  not. Ensure normal and uncaught abrupt completion perform the Thread final action and
  wake joiners.
- Add one lazy monitor sidecar per Java object/array/Class mirror. Implement explicit
  `monitorenter`/`monitorexit`, recursion, ownership, blocking, null/non-owner failures,
  `ACC_SYNCHRONIZED` instance/static method entry, and exactly-once normal/abrupt release.
- Add bounded synthetic Object `wait()`, `notify()`, and `notifyAll()` methods backed by
  the same monitor. Wait fully releases and restores recursion depth. Order notify and
  interrupt atomically so a notification is not lost when interrupt wins.
- Implement interruption status and the fixed matrix's wait, untimed join, and
  `sleep(long)` interruption/clear behavior. Ordinary monitor acquisition remains
  non-interruptible.
- Extend the shared ADR-0025 initialization service with a unique Class initialization
  lock/condition, other-owner wait/retry, same-owner recursion, notify-all on all terminal
  transitions, unchanged interrupt status, and publication of initialized state.
- Add the minimal synthetic exception facades required by the fixed matrix:
  `IllegalThreadStateException`, `IllegalMonitorStateException`, and
  `InterruptedException`, with existing Java exception transport in Interpreter and IR.
- Preserve the fixed 19 Java sources and create a fail-closed candidate harness under
  `docs/workstreams/r2-concurrency-fixtures/`. Candidate evidence is isolated from the
  immutable research baseline.
- Add an AOT reachability/build check that classifies all 19 fixed concurrency fixtures as
  `Not implemented`/`NO-BUILD`. A built binary that later panics or mismatches is a failure.
- Update architecture/development documentation whose single-threaded, heap, Thread,
  monitor, initialization, or AOT boundary becomes stale.

## Non-scope

- Concurrent AOT execution, implicit goroutine-local current Thread, an AOT execution-
  context ABI, or emitted AOT heap/monitor operations. All fixed AOT concurrency fixtures
  remain explicit `Not implemented`.
- Timed `wait`/`join` variants, nanosecond timing, scheduling fairness, priorities,
  ThreadGroup, ThreadLocal/inheritable locals, context class loaders, uncaught-exception
  handler APIs, virtual threads, thread pools, or `java.util.concurrent`.
- Unsafe/VarHandle, broad java.base Thread/Object/Class replacement, reflection, JNI,
  I/O synchronization, or arbitrary mutable native payload concurrency.
- Weak-memory optimization, non-SC heap access, thin/biased locks, monitor deflation,
  deadlock detection, lock elision, or performance/scale claims.
- Broad AOT exception propagation, `invokedynamic`, or unrelated class-library/runtime
  expansion.

## Semantic constraints

- Java 25 is the semantic baseline and Temurin 25.0.3 is the differential reference.
- Java Thread identity is the stable facade object attached to an execution context, not
  a goroutine or numeric ID. Start succeeds once; join and termination detection publish
  the target's prior actions.
- All Java heap access in the supported runtime is Go-race-free and SC. Volatile and final
  visibility is at least Java 25 strength; stronger ordinary-field ordering is intentional
  under ADR-0030.
- Monitor owner and recursion are keyed by Java execution-context identity. Wait restores
  the exact recursion depth before returning or throwing. Every implicit synchronized
  method entry is released exactly once on every exit path.
- Static synchronized methods lock the canonical Class mirror for the declaring runtime
  Class identity. Clones never copy monitor ownership or waiters.
- Notify/interrupt races follow one Java-permitted order and cannot lose a notification.
  Monitor entry is not made interruptible.
- Cross-thread initialization preserves ADR-0025's four states, predecessor/failure rules,
  and shared engine service while adding the JVMS §5.5 synchronization protocol.
- Interpreter and IR report `Supported`; AOT reports `Not implemented`. No engine may
  silently approximate, build-then-panic, or share one mutable Thread across carriers.

## Required completion state by engine

| Capability | Interpreter | IR | AOT |
|---|---|---|---|
| Stable current Thread, lifecycle, start/join/daemon liveness | Required | Required | Not implemented / build rejection |
| SC heap cells and volatile/final fixture behavior | Required | Required | Not implemented for concurrency fixtures |
| Explicit and synchronized-method monitors | Required | Required | Not implemented / build rejection |
| Wait/notify-all and interrupt of wait/join/sleep | Required | Required | Not implemented / build rejection |
| Cross-thread class initialization | Required | Required | Not implemented / build rejection |
| `ProducerConsumer` milestone | Match | Match | Not implemented / build rejection |

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| Fixed candidate differential matrix | `bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh <candidate>` over exactly 19 frozen fixtures; fails closed on missing tool/fixture/row, timeout, output/exit mismatch, or unexpected engine state | Not run |
| Interpreter / IR | All 19 fixtures match Temurin 25 combined stdout, stderr, and exit code | Not run |
| AOT rejection matrix | All 19 fixtures report `Not implemented` as `NO-BUILD`; any built binary, mismatch, panic, omitted row, or fallback is Fail | Not run |
| Race-enabled concurrency stress | `R2_CONCURRENCY_STRESS=100 bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh <candidate>` with a race-built catty Interpreter/IR; no race, timeout, deadlock, missing iteration, or semantic mismatch | Not run |
| Kernel/unit invariants | `go test ./...` includes heap-cell race-free copy/access, loader single identity, canonical mirrors, monitor ownership/depth/wait ordering, lifecycle/start-once/join/daemon, interrupt, and class-init contention/failure tests | Not run |
| Core regression | `go vet ./... && go test -race ./... && bash tests/run.sh` | Not run |
| Evidence isolation | historical `r2-evidence`, initialization evidence, String evidence, and `r2-concurrency-evidence/baseline-63d5658/` unchanged; candidate output only under `docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/` | Not run |
| Governance | `git diff --check <acceptance-anchor>..<candidate>` and documentation matches actual per-engine support | Not run |

Results use only `Pass`, `Fail`, `Not run`, or `Not implemented`.

## Amendments

Accepted changes are appended here after Owner approval; the frozen contract is not
rewritten to reduce gates.

### Amendment 1 — Slice C `git diff --check` scope and stress runner race build

Accepted by Owner on 2026-07-15. This amendment refines how two Slice C review
gates are executed without reducing their intent or the parent workstream's
final acceptance gates.

**Problem.** The intermediate candidate `36a577c` sealed Slice C evidence
(`results.txt`, `results-stress-20x.txt`) was produced by an earlier revision
of `run-concurrency-slice-c.sh` whose table header carried trailing whitespace.
Commit `0a96b59` fixed the runner so later evidence is clean, but the sealed
`36a577c` evidence is immutable candidate evidence (ADR-0029, COLLABORATION.md
§9.1) and must not be rewritten. Consequently the literal command
`git diff --check <acceptance-anchor>..<candidate>` flags the immutable
intermediate artifact even though the current candidate's production code and
its own evidence are clean. Separately, the runner built a plain binary even
under `R2_CONCURRENCY_STRESS>1`, so the "race-built catty binary" requirement
of the Slice C review evidence table was not honored by the harness.

**Refinement.**

1. The `git diff --check` governance check for a Slice C candidate is scoped to
   the candidate's production code and the candidate's own evidence directory:
   `git diff --check <acceptance-anchor>..<candidate> -- interpreter native rtda
   docs/workstreams/r2-thread-monitor-foundation-slice.md
   docs/workstreams/r2-concurrency-fixtures` plus, for the candidate's own
   evidence, `git diff --check <acceptance-anchor>..<candidate> --
   docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/`.
   Previously sealed intermediate candidate evidence directories (e.g.
   `36a577c`, `0a96b59`) are excluded because they are immutable artifacts and
   not part of the current candidate's production surface. This does not weaken
   the parent workstream's final `git diff --check <acceptance-anchor>..<candidate>`
   Governance gate for the final integrated candidate; that remains unchanged.
2. `run-concurrency-slice-c.sh` now builds catty with `-race` when
   `R2_CONCURRENCY_STRESS>1` and records `race-build` in the evidence header,
   so the Slice C stress gate is reproducible as written.

**Non-scope.** No change to Slice C Outcome, Scope, Non-scope, Semantic
constraints, the engine matrix, the 11 fixture rows, the parent workstream's
acceptance gates, or the monitor/Thread implementation.

## Acceptance record

Accepted by Owner on 2026-07-14. This acceptance authorizes production work only
within this frozen contract after a new Active Agent starts from the acceptance
anchor and records the required implementation preflight. It does not authorize
scope expansion, integration beyond the Owner's stated authority, or changes to
the frozen sections without an accepted amendment.

---

## Implementation preflight

Before production implementation record:

- **Parent acceptance anchor:** `a0288be` (2026-07-14 parent governance commit)
- **Slice C acceptance anchor / actual base:** `f3800b7` (2026-07-14 Slice C governance commit) / worktree at `f3800b7`
- **Historical evidence check:** all files under `docs/workstreams/r2-evidence/`, `docs/workstreams/r2-concurrency-evidence/` unchanged from `a0288be` through `d4008c0` — **Pass**
- **Slice A/B evidence:** `9576828/`, `505d3ee/`, `a0e336c/`, `b0a7b70/` intact — **Pass**
- **Required ADRs confirmed Accepted:** 0017, 0018, 0020, 0022, 0028, 0029, 0030 — **Pass**
- **Candidate evidence destination:** `docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/slice-c/` (not yet created)
- **Harness output policy:** explicit candidate required; never writes research baseline or shared/latest path

All preflight items complete. Slice C may proceed to `In progress`.

---

## Plan

| Slice | Status | Evidence |
|---|---|---|
| A — SC heap cells, concurrency-safe loader, and canonical Class mirrors | Complete | `docs/workstreams/r2-concurrency-candidate-evidence/9576828/` — `ec1b398`, 22 files, all gates Pass |
| B — stable Thread facade/context, lifecycle, carriers, join, and VM liveness | Complete | `docs/workstreams/r2-concurrency-candidate-evidence/b0a7b70/` — `b0a7b70` (final), Owner accepted 2026-07-14, all Slice B gates Pass |
| C — monitors, synchronized methods, wait sets, and interruption | Ready | `docs/workstreams/r2-concurrency-candidate-evidence/f0fc2ca/slice-c/` — `f0fc2ca` (rework: stress runner now race-built per Amendment 1), 11/11 fixtures Interpreter + IR (1x + race-built 20x stress), corrected gates Pass; awaiting Owner completion acceptance |
| D — concurrent ADR-0025 initialization and full Interpreter/IR fixture matrix | Pending | — |
| E — AOT fail-closed rejection, race stress, regression, evidence, and docs | Pending | — |

Status uses `Pending`, `In progress`, or `Complete`.

### Slice C technical investigation

Investigated on 2026-07-14 against the accepted Slice B lineage at `d4008c0`,
Java 25 JVMS/JLS/API text, the frozen fixtures, and current Interpreter/IR/native
entry and unwind paths.

The research establishes these implementation facts:

- Java 25 `monitorenter` is reentrant, blocks competing execution contexts, and throws
  `NullPointerException` for null. `monitorexit` throws `NullPointerException` for null
  and `IllegalMonitorStateException` for a non-owner. Ordinary entry is not
  interruptible.
- `ACC_SYNCHRONIZED` is implicit invocation behavior. Instance methods lock the receiver;
  static methods lock the canonical `Class` mirror for the method's actual declaring
  runtime Class. The VM releases only that implicit entry when the invocation completes
  normally or when an exception escapes the method.
- Javac emits explicit `monitorenter`/`monitorexit` plus exception-table cleanup for a
  synchronized block. Slice C must not release arbitrary explicit entries when a frame is
  popped; the shared frame lifecycle owns only the method's implicit entry.
- Current ordinary Interpreter and IR calls share `invokeMethod`, but `Thread.start()` and
  the interpreter bridge can push frames directly. Implicit synchronized entry therefore
  belongs in a shared frame-entry boundary, not only in invoke opcode handlers. Native
  throwaway frames need an equivalent paired cleanup because they are not pushed.
- A pre-interrupted call to `wait()` clears interrupt status and throws without releasing
  the monitor. An interrupt after wait-set enrollment removes the waiter, but the thread
  throws only after reacquiring the monitor and restoring its prior recursion depth.
- Notify and interrupt must be ordered per waiter. If notify wins, wait returns normally
  and the later interrupt remains pending. If interrupt wins, a later notify skips that
  waiter and can select another; a notification cannot be consumed by an already
  interrupted waiter.
- The existing Thread waker channel is sufficient for sleep/join wakeup but cannot alone
  order monitor notification against interruption. Thread needs an atomic active-waiter
  registration boundary; the interrupt path must not hold that registration lock while
  acquiring a monitor state lock.
- The fixed bytecode calls `Object.wait:()V` directly. `SynchronizedMethods` and
  `ProducerConsumer.OneSlot` carry `ACC_SYNCHRONIZED` rather than explicit monitor
  opcodes. The bounded synthetic facade must therefore declare the exact `wait()`,
  `notify()`, and `notifyAll()` descriptors.

Authoritative semantic references are
[JVMS 2.11.10](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-2.html#jvms-2.11.10),
[JVMS monitorenter/monitorexit](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-6.html#jvms-6.5.monitorenter),
[JLS 17.1-17.2](https://docs.oracle.com/javase/specs/jls/se25/html/jls-17.html#jls-17.2),
and the Java 25 [`Object`](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/Object.html)
and [`Thread`](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/Thread.html)
APIs.

### Slice C working contract

**Status:** Accepted by Owner on 2026-07-14

**Type:** implementation within this Accepted parent workstream

**Review:** owner

**Planned base:** `d4008c0` (accepted Slice B governance)

**Parent acceptance anchor:** `a0288be`

**Governing decisions:** the parent contract and ADR-0017, ADR-0018, ADR-0020,
ADR-0022, and ADR-0028 through ADR-0030

This checkpoint refines the existing plan without amending the parent's frozen Outcome,
Scope, Non-scope, Semantic constraints, engine matrix, or Acceptance gates. Its evidence
does not replace the fixed 19-fixture denominator or any final workstream gate. Owner
acceptance authorizes an Active Agent to change Slice C from `Pending` to `In progress`
after it records its implementation preflight. Production work remains limited by the
parent contract and the Owner's session authority.

#### Slice C outcome

Interpreter and IR share one race-free Java object-monitor service covering explicit
monitor bytecodes, implicit synchronized methods, untimed Object wait sets, notification,
and wait interruption. The eight Slice C fixtures match Temurin 25, the three existing
interrupt fixtures remain matching regressions, and the one-slot `ProducerConsumer`
milestone works without moving concurrent class initialization or AOT acceptance into
this slice.

#### Slice C scope

- Add one lazy, CAS-published monitor sidecar to every `rtda.Object`, including arrays and
  canonical Class mirrors. Keep owner execution-context ID, recursion depth, entry
  coordination, and ordered waiters under one monitor state lock. `CloneObject` creates a
  fresh object with no copied monitor state.
- Implement monitor enter, exit, ownership query, untimed wait, notify-one, notify-all,
  and waiter interruption in `rtda`. Entry is reentrant and non-interruptible; notify
  selection need not be fair.
- Add a Thread active-waiter registration protocol that closes the race between the
  pre-wait interrupt check and waiter publication. Preserve the accepted Slice B
  sleep/join waker behavior and clear interrupt status exactly when an
  `InterruptedException` is selected.
- Implement `monitorenter` and `monitorexit` in both Interpreter and IR using the same
  `rtda` service and existing Java exception transport.
- Recognize `ACC_SYNCHRONIZED`. Acquire the instance receiver or the declaring Class's
  canonical mirror at the shared frame-entry boundary. Record exactly one implicit entry
  on the frame and release it exactly once through normal return or exception unwind.
  Apply equivalent paired handling to native throwaway frames.
- Cover every interpreted frame-entry route in scope, including ordinary invocation,
  spawned `Thread.run()` dispatch, `<clinit>` callbacks when applicable, and the existing
  interpreter bridge. This does not claim concurrent AOT support.
- Add bounded synthetic/native Object `wait()`, `notify()`, and `notifyAll()` methods,
  real `Thread.holdsLock`, and the `IllegalMonitorStateException` facade. Preserve
  `InterruptedException` transport from Slice B.
- Add deterministic unit tests for lazy single-monitor identity, exclusion, recursion,
  ownership failures, exact depth restore, pre-interrupted wait, notify-one/all,
  interrupt-before-notify, notify-before-interrupt, no lost notification, frame cleanup,
  static Class-mirror locking, and clone isolation.
- Add a fail-closed Slice C runner whose explicit fixture list and output are isolated
  under `docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/slice-c/`.

#### Slice C non-scope

- Concurrent class initialization or changes to ADR-0025 state/locking; that remains
  Slice D.
- The full 19-fixture candidate matrix, AOT reachability rejection, full stress gate, or
  final architecture/development documentation; those remain later parent slices.
- Timed `wait`/`join`, nanosecond timing, fairness, deliberate spurious wakeups, deadlock
  detection, monitor deflation, thin/biased locking, virtual-thread pinning, or broad
  java.base Object/Thread replacement.
- Interruptible monitor acquisition, Thread stop/suspend/resume, asynchronous exception
  injection, or changes to the accepted Thread lifecycle and VM-liveness model except the
  active-waiter integration required here.
- Concurrent AOT, an AOT execution-context ABI, emitted monitor operations, or any AOT
  `Supported`/`Fallback` claim. Incidental single-context bridge behavior is not evidence
  of AOT concurrency support.
- Reworking heap cells, loader identity, Class-mirror canonicality, or mutable native
  payloads beyond changes strictly required to attach and exercise monitor state.

#### Slice C semantic constraints

- Monitor ownership is keyed only by stable Java execution-context identity. Goroutine
  identity and Java Thread facade pointer identity are not owner keys.
- Monitor state transitions and waiter state transitions are race-free. An unlock is the
  release edge for a later successful lock of the same monitor; the implementation may
  use stronger ordering.
- A pre-interrupted `wait()` throws while retaining the caller's current ownership and
  recursion depth. A waiter interrupted after enrollment reacquires the monitor and
  restores the exact prior depth before its cleared-status `InterruptedException` is
  observed by Java code.
- Waiter states are single-transition: `waiting` becomes `notified` or `interrupted` once.
  Wake signals are private and exactly-once. Notify skips non-waiting entries.
- If notify wins its ordering against interrupt, wait returns normally and interrupt
  status remains set. If interrupt wins, the waiter throws after reacquisition and later
  notify remains available to another eligible waiter.
- A notified or interrupted waiter competes normally for monitor reacquisition; it does
  not resume Java execution while another execution context still owns the monitor.
- Static synchronized methods lock `method.Owner()`'s canonical Class mirror, including
  inherited resolution through another symbolic class. Instance synchronized methods
  lock local 0 (`this`).
- Only the implicit synchronized-method entry is attached to frame cleanup. Explicit
  block entries remain governed by bytecode `monitorexit` and its exception handlers.
  A handler caught inside the synchronized method does not release the implicit entry;
  an exception escaping that frame does.
- Interpreter and IR use identical services and exception classes. Neither engine may
  busy-wait, silently ignore ownership errors, or approximate wait with sleep/polling.

#### Slice C review evidence

These are slice-review checks, not amendments or substitutes for the parent's final
Acceptance table. Results use only `Pass`, `Fail`, `Not run`, or `Not implemented`.

| Check | Required evidence before Slice C can be proposed Complete |
|---|---|
| Monitor kernel | `go test -race ./rtda` covers lazy CAS identity, exclusion, recursion, non-owner exit, wait depth, both notify/interrupt orders, no lost notification, and clone isolation |
| Invocation/unwind | Unit coverage proves instance/static implicit entry, canonical Class mirror identity, all return types, caught versus escaping exceptions, direct frame entry, native cleanup, and exactly-once release |
| Native facade | `go test -race ./native` covers Object ownership failures, wait interruption and status clearing, `holdsLock`, and preservation of Slice B sleep/join interrupt behavior |
| Direct Slice C differential | `bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-slice-c.sh <candidate>` asserts exactly eleven Interpreter/IR rows. The direct rows are `SynchronizedBlocks`, `SynchronizedMethods`, `MonitorNull`, `MonitorOwnership`, `WaitNotify`, `NotifyAll`, `InterruptWait`, and `ProducerConsumer`; each matches Temurin 25 combined stdout, stderr, and exit code |
| Interrupt regression differential | The same runner asserts `InterruptStatus`, `InterruptSleep`, and `InterruptJoin` still match Temurin 25 in Interpreter and IR |
| Bounded race repetition | `R2_CONCURRENCY_STRESS=20 bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-slice-c.sh <candidate>` uses a race-built catty binary for the direct Slice C set; no Go race, timeout, deadlock, missing iteration, or mismatch. This does not replace the parent's later `R2_CONCURRENCY_STRESS=100` full-matrix gate |
| Core regression | `go vet ./...`, `go test ./...`, `go test -race ./...`, and `bash tests/run.sh` all Pass |
| Evidence isolation | The runner requires an explicit immutable candidate ID, refuses overwrite, records toolchain/base/fixture list, and writes only its candidate `slice-c/` directory; all historical evidence hashes remain unchanged |
| Scope audit | Diff from `d4008c0` contains no concurrent class-init implementation, AOT concurrency claim, fixture denominator change, or stale single-threaded claim introduced by Slice C |

The Slice C runner may share code with the future full candidate harness, but it must
hard-code and report the eleven rows above. Additional tests are supplemental and cannot
change the parent's 19-fixture denominator or final pass count.

#### Slice C implementation order

1. Monitor sidecar and kernel invariants.
2. Thread active-waiter ordering and interrupt race tests.
3. Explicit Interpreter/IR monitor opcodes and Java exception mapping.
4. Shared synchronized frame entry/return/unwind, including direct frame and native paths.
5. Object/Thread facades and exception hierarchy.
6. Targeted differential runner, race repetition, regression, isolated evidence, and
   owner review candidate.

#### Slice C acceptance record

Accepted by Owner on 2026-07-14. This acceptance approves the Slice C refinement within
the parent workstream. It does not itself start production implementation, authorize a
commit/integration action, or alter the parent workstream's final acceptance gates.

---

## Handoff

- **Branch / current head:** `worktree-r2-thread-monitor-foundation` / `d4008c0` (accepted Slice B governance)
- **Last implementation candidate:** `b0a7b70` (Slice B final, Owner accepted 2026-07-14)
- **Acceptance anchor / planned Slice C base / research baseline:** `a0288be` / `d4008c0` / `63d5658`
- **Slice A evidence:** `docs/workstreams/r2-concurrency-candidate-evidence/9576828/` — `ec1b398`, accepted by Owner
- **Slice B original:** `docs/workstreams/r2-concurrency-candidate-evidence/505d3ee/` — `505d3ee`
- **Slice B rework 1:** `docs/workstreams/r2-concurrency-candidate-evidence/a0e336c/` — `a0e336c` (3 blockers)
- **Slice B rework 2 (final):** `docs/workstreams/r2-concurrency-candidate-evidence/b0a7b70/` — `b0a7b70` (daemon race fix)
- **Rework 2 scope:** 2 files, +76/−3 — IsDaemon holds configMu
- **Gates (all run on `b0a7b70`):**
  - `go build ./...` — **Pass**
  - `go vet ./...` — **Pass**
  - `go test ./...` — **Pass**
  - `go test -race ./...` — **Pass** (all 8 packages)
  - `bash tests/run.sh` — **Pass** (10/10 fixtures)
  - `git diff --check 3034e05..b0a7b70` — **Pass**
- **New test:** rtda: +1 test function `TestConcurrentSetDaemonAndIsDaemon` (3 subtests) — all under `-race`
- **Contract gates not yet run:** 19-fixture matrix, AOT rejection matrix, race stress, evidence isolation check
- **Slice A scope:** 22 files, +1306/−259 — HeapCell typed accessors, CopyObjectCells overlap-safe, Cells()/StaticCells() removed, classloader CAS/double-check, canonical Class mirrors via ClassObject CAS-once, 34 new `-race` tests
- **Slice B scope (original):** 10 files, +1464/−23 — VM supervisor, Thread lifecycle/interrupt/daemon/sleep, 15 native Thread methods, goroutine carrier, join, DefaultRunLoop callback, 51 new `-race` tests (rtda: 32 thread + 5 vm; native: 14)
- **Dirty files:** this contract and `run-concurrency-slice-c.sh` updated for Slice C rework (Amendment 1); production code unchanged from `0a96b59`
- **Current Slice C state:** rework candidate `f0fc2ca` Ready; stress runner now race-built, `git diff --check` gap resolved by Amendment 1; awaiting Owner completion acceptance
- **Slice C implementation candidate (prior):** `0a96b59` — monitor/Thread implementation, 11/11 fixtures (1x + non-race 20x)
- **Slice C rework candidate:** `f0fc2ca` — runner race-build fix only, 11/11 fixtures (1x + race-built 20x), corrected gates Pass
- **Slice C evidence:** `docs/workstreams/r2-concurrency-candidate-evidence/f0fc2ca/slice-c/` — `results.txt` (1x) and `results-stress-20x.txt` (race-built), both original `0a96b59/slice-c` evidence preserved
- **Gates (all run on `f0fc2ca`):**
  - `go vet ./...` — **Pass**
  - `go test ./...` — **Pass**
  - `go test -race ./...` — **Pass**
  - `bash tests/run.sh` — **Pass** (10/10)
  - `bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-slice-c.sh f0fc2ca` — **Pass** 11/11 Interpreter + IR
  - `R2_CONCURRENCY_STRESS=20 ...` (race-built) — **Pass** 11/11 Interpreter + IR
  - `git diff --check f3800b7..f0fc2ca` (scoped per Amendment 1) — **Pass**
  - historical evidence unchanged — **Pass**
- **Next action (Slice C):** Owner completion acceptance of `f0fc2ca`; on accept, update Plan to Complete and proceed to Slice D contract
- **Non-derivable context:** the 19-fixture denominator includes explicit daemon and non-daemon liveness, all three interruptible blocking points, and the producer-consumer milestone

### Slice B acceptance record

Accepted by Owner on 2026-07-14. Slice B is accepted as a completed implementation
slice within this workstream. The full workstream remains open; Slice C requires its
own implementation work and review, and the full 19-fixture, AOT, stress, and final
integration gates remain not run.

### Slice C completion record

Rework candidate `f0fc2ca` was produced on 2026-07-15 under Amendment 1 to make
the Slice C review gates reproducible. All Slice C review gates Pass on
`f0fc2ca`: core regression, 11-fixture 1× matrix (Interpreter + IR), race-built
20× stress (Interpreter + IR), the scoped `git diff --check`, and evidence
isolation. The parent workstream's final 19-fixture, AOT rejection, `STRESS=100`,
and integration gates remain not run and are governed by Slices D and E.

Status: **Ready — awaiting Owner completion acceptance.** On Owner acceptance,
Slice C is marked Complete in the Plan table and the next action is to draft the
Slice D workstream contract for concurrent ADR-0025 initialization and the full
Interpreter/IR fixture matrix.
