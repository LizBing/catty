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

### Amendment 2 — R2 Slice C boundary argument validation

Accepted by Owner on 2026-07-15. This amendment adds two bounded Java
argument-validation corrections to the Slice C native surface. Neither changes
the behavior of any of the 11 fixed Slice C fixtures or any parent
acceptance gate; both repair observable Java semantics that were not exercised
by the fixture set.

**Problem.** `Thread.holdsLock(null)` returned `false` instead of throwing
`NullPointerException`, and `Object.wait(long)` / `Object.wait(long, int)`
accepted negative `timeoutMillis` (and, for the two-argument form, ignored the
`nanos` range) instead of throwing `IllegalArgumentException`. The Java 25
API specifies:

- `java.lang.Thread` (class level): "Unless otherwise specified, passing a null
  argument to a constructor or method in this class will cause a
  NullPointerException to be thrown." `holdsLock(Object)` does not specify
  otherwise, so `holdsLock(null)` throws `NullPointerException`.
- `java.lang.Object.wait(long timeoutMillis)`:
  `IllegalArgumentException` - if timeoutMillis is negative.
- `java.lang.Object.wait(long timeoutMillis, int nanos)`:
  `IllegalArgumentException` - if timeoutMillis is negative, or if the value of
  nanos is out of range (`0-999999`).

**Change.**

1. `threadHoldsLock` now throws `NullPointerException` when its argument is null.
2. `objectWait` and `objectWaitJI` now throw `IllegalArgumentException` for a
   negative `timeoutMillis`; `objectWaitJI` also throws it when `nanos` is
   outside `0-999999`. For the two-argument form the argument-range check
   precedes the monitor-ownership check, matching `Object.wait(long,int)`.

**Non-scope.** Timed-wait behavior (actual timeout enforcement), spurious
wakeup handling, nanosecond timer precision, and any change to the monitor
kernel, Thread lifecycle, interrupt protocol, the 11 fixture rows, or the
parent workstream's acceptance gates. Timed wait remains Slice D scope; this
amendment adds only argument-range validation.

**Evidence.** `native/object_test.go` covers `holdsLock(null)` → NPE,
`wait(-1)` → IAE, `wait(-5, 0)` → IAE, `wait(100, 1000000)` → IAE, and
`wait(100, -1)` → IAE; ownership and depth are asserted unchanged after the
rejected call. The 11 Slice C fixtures remain byte-identical to Temurin 25 in
Interpreter and IR (1× and race-built 20×).

### Amendment 3 — Correct stale single-threaded/monitor-as-nop claims in `ARCHITECTURE.md`

Accepted by Owner on 2026-07-15. This amendment corrects `docs/ARCHITECTURE.md`
so its stated capability boundary matches the Accepted Slice A/B/C
implementation and ADRs 0028–0030, instead of the pre-R2 single-threaded MVP.

**Problem.** Slice C implements `monitorenter`/`monitorexit`,
`ACC_SYNCHRONIZED`, `Object.wait`/`notify`/`notifyAll`, interruption, and
Thread lifecycle, and Slices A/B added race-free SC heap, concurrency-safe
loading, canonical Class mirrors, and VM daemon liveness. `ARCHITECTURE.md`
still described the runtime as "MVP single-threaded", listed
`monitorenter`/`monitorexit` as nops and `synchronized`/`wait`/`notify` as
missing, and treated the Java memory model as moot. This contradicted the
repository's Accepted capability evidence and could mislead a future Active
Agent.

**Change.** Three corrections to `ARCHITECTURE.md` only (no code):

1. The Thread-scheduler row of the §1 premise table now states the bounded
   Java 25 concurrency implemented in Interpreter/IR and lists the
   `Not implemented` remainder (timed `wait`/`join`, `Unsafe`, JMM
   optimizations, virtual threads, fairness, deadlock detection, AOT
   concurrency).
2. The `rtda.Thread` section now reflects `ecID`, Java Thread facade,
   lifecycle, one goroutine carrier per started platform Thread, and the
   ADR-0028/0029/0030 implemented surface.
3. §10 "What catty does *not* model (yet)" replaces the
   "Concurrency — single-threaded; no synchronized/wait/notify" and
   "Java memory model — moot while single-threaded" bullets with the
   implemented bounded concurrency surface and its `Not implemented`
   remainder.

**Non-scope.** No change to Slice C Outcome, Scope (production code),
Non-scope, Semantic constraints, the engine matrix, the 11 fixture rows, the
parent workstream's acceptance gates, or further architecture/development
documentation beyond these three corrections. `DEVELOPMENT.md` has no stale
single-threaded/monitor text and is unchanged. Remaining architecture/
development doc updates remain Slice E scope.

**Evidence.** The three `ARCHITECTURE.md` hunks are doc-only; the scoped
`git diff --check f3808b7..HEAD` over `docs/ARCHITECTURE.md` and the
production surface is clean. The 11 Slice C fixtures are unaffected (docs-only
change); no gate is re-run beyond the whitespace and build sanity check.

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
| C — monitors, synchronized methods, wait sets, and interruption | Complete | `docs/workstreams/r2-concurrency-candidate-evidence/eea253d/slice-c/` — `eea253d` (Amendment 1 race-built stress gate + Amendment 2 holdsLock/wait argument validation), 11/11 fixtures Interpreter + IR (1x + race-built 20x stress), all gates Pass; Owner accepted 2026-07-15 |
| D — concurrent ADR-0025 initialization and full Interpreter/IR fixture matrix | In progress | preflight Pass; branch `r2-slice-d-concurrent-init` at `4a91470`; acceptance anchor `c4ddde4` |
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
- **Dirty files:** this contract only (`run-concurrency-slice-c.sh`, `native/system.go`, `native/object_test.go` committed as candidates below)
- **Current Slice C state:** rework candidate `eea253d` Complete; Amendment 1 (race-built stress gate, scoped `git diff --check`), Amendment 2 (`holdsLock(null)` → NPE, `wait`/`wait(long,int)` argument-range validation), and Amendment 3 (`ARCHITECTURE.md` single-threaded/monitor-as-nop corrections) applied; Owner accepted 2026-07-15
- **Slice C implementation candidate (original):** `0a96b59` — monitor/Thread implementation, 11/11 fixtures (1x + non-race 20x)
- **Slice C rework candidate 1:** `f0fc2ca` — runner race-build fix only (Amendment 1)
- **Slice C rework candidate 2 (current):** `eea253d` — adds Amendment 2 boundary validation; 11/11 fixtures (1x + race-built 20x), corrected gates Pass
- **Slice C evidence:** `docs/workstreams/r2-concurrency-candidate-evidence/eea253d/slice-c/` — `results.txt` (1x) and `results-stress-20x.txt` (race-built); `0a96b59/` and `f0fc2ca/` evidence preserved
- **Gates (all run on `eea253d`):**
  - `go vet ./...` — **Pass**
  - `go test ./...` — **Pass**
  - `go test -race ./...` — **Pass**
  - `bash tests/run.sh` — **Pass** (10/10)
  - `bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-slice-c.sh eea253d` — **Pass** 11/11 Interpreter + IR
  - `R2_CONCURRENCY_STRESS=20 ...` (race-built) — **Pass** 11/11 Interpreter + IR
  - `git diff --check f3800b7..eea253d` (scoped per Amendment 1) — **Pass**
  - historical evidence unchanged — **Pass**
  - Amendment 2 unit tests (`native/object_test.go`) under `-race` — **Pass**
- **Next action (Slice C):** accepted by Owner 2026-07-15; Plan marked Complete; proceed to Slice D contract
- **Non-derivable context:** the 19-fixture denominator includes explicit daemon and non-daemon liveness, all three interruptible blocking points, and the producer-consumer milestone

### Slice B acceptance record

Accepted by Owner on 2026-07-14. Slice B is accepted as a completed implementation
slice within this workstream. The full workstream remains open; Slice C requires its
own implementation work and review, and the full 19-fixture, AOT, stress, and final
integration gates remain not run.

### Slice C completion record

Rework candidate `eea253d` consolidates the Slice C implementation under three
accepted amendments. Amendment 1 (rework candidate `f0fc2ca`) makes the Slice C
review gates reproducible: stress uses a race-built catty binary and the
`git diff --check` scope excludes immutable intermediate candidate evidence.
Amendment 2 (`eea253d`) adds bounded argument validation: `Thread.holdsLock(null)`
throws `NullPointerException` and `Object.wait(long)` / `Object.wait(long, int)`
throw `IllegalArgumentException` for negative/out-of-range arguments, with timed
wait behavior itself still deferred to Slice D. Amendment 3 (docs-only) corrects
`ARCHITECTURE.md`'s stale single-threaded/monitor-as-nop claims to match the
Accepted A/B/C capability boundary.

All Slice C review gates Pass on `eea253d`: core regression, the Amendment 2
`native` unit tests under `-race`, the 11-fixture 1× matrix (Interpreter + IR),
race-built 20× stress (Interpreter + IR), the scoped `git diff --check`, and
evidence isolation. The parent workstream's final 19-fixture, AOT rejection,
`STRESS=100`, and integration gates remain not run and are governed by Slices
D and E.

Status: **Complete — accepted by Owner on 2026-07-15.** Slice C is marked
Complete in the Plan table. The next action is to draft the Slice D workstream
contract for concurrent ADR-0025 initialization and the full Interpreter/IR
fixture matrix.

---

## Slice D working contract

**Status:** Accepted
**Type:** implementation within this Accepted parent workstream
**Review:** owner
**Planned base:** `ff691b5` (current `main`; accepted Slice C governance integration)
**Parent acceptance anchor:** `a0288be`
**Acceptance anchor:** `c4ddde4` (`docs(r2-slice-d): accept Slice D working contract, set acceptance anchor`, 2026-07-16). This is the commit that froze the Slice D contract as Accepted. The Active Agent must branch from the current `main` tip (which `c4ddde4` and its governance follow-ups are already part of), record the resolved SHA of `c4ddde4` in the preflight, and treat `c4ddde4` as the frozen contract anchor.
**Governing decisions:** the parent contract and ADR-0017, ADR-0018, ADR-0020,
ADR-0022, ADR-0025, and ADR-0028 through ADR-0030

This checkpoint refines the existing plan without amending the parent's frozen
Outcome, Scope, Non-scope, Semantic constraints, engine matrix, or Acceptance
gates. Its evidence does not replace the fixed 19-fixture denominator or any
final workstream gate. On Owner acceptance this authorizes an Active Agent to
record the implementation preflight, change Slice D from `Pending` to
`In progress`, and implement only the accepted Slice D contract. Production work
remains limited by the parent contract and the Owner's session authority.

### Slice D outcome

Catty's shared ADR-0025 class/interface initialization service becomes
race-free across Java execution contexts: exactly one `<clinit>` runs per
runtime Class identity, concurrent other-owner requests block on a per-Class
condition until a terminal state is published, and publication establishes
release/acquire visibility without changing interrupt status of waiters.
Interpreter and IR pass the full fixed 19-fixture concurrency matrix against
Temurin 25 (combined stdout, stderr, exit code), including
`CrossThreadClassInitialization` under race. AOT build-rejects all 19 fixed
concurrency fixtures as `Not implemented`/`NO-BUILD`. The runtime is now
evidence-backed for the R2 producer-consumer milestone end-to-end in
Interpreter and IR.

### Slice D scope

- Add one dedicated initialization lock/condition per runtime `Class` identity
  behind the existing shared `InitializeClass` service, **distinct from the
  Class mirror's Java monitor** (ADR-0029). Representation choice: a `sync.Mutex`
  plus `sync.Cond` owned by each `Class` (constructor-allocated; no lazy map
  lookup).
- Extend `rtda.InitializeClass` with the JVMS §5.5 cross-context protocol:
  - already `initialized` → return normally (acquire the initialized state);
  - recursive same-owner (`initInProgress` and `initOwner == ecID`) → return
    normally;
  - `erroneous` → return `NoClassDefFoundError`;
  - `initNotStarted` → claim ownership under the lock, release the lock for the
    clinit body, run predecessor initialization and `<clinit>` via the existing
    `ClinitRunner` callback, then under the lock publish the terminal state and
    `notify-all`;
  - `initInProgress` with a different owner → **wait** on the per-Class
    condition without holding the Java monitor; on wake, re-read the state and
    proceed (return if `initialized`, throw `NoClassDefFoundError` if
    `erroneous`, loop if still `initInProgress`). **Interrupt status of the
    waiter is unchanged by this wait** (ADR-0029); interrupt is not reported and
    does not break the wait.
  - Terminal publication (success and every erroneous transition) performs an
    atomic state write under the per-Class lock and `notify-all`, establishing
    release/acquire visibility for the initialized state and any clinit-written
    static fields.
- Add a bounded race-deterministic unit test suite for the kernel:
  multiple execution contexts contending on one Class; `<clinit>` runs exactly
  once; second and later waiters see the published value and no clinit re-run;
  erroneous publication reaches all waiters as `NoClassDefFoundError`; a waiter
  interrupted while waiting still completes the init wait with unchanged
  interrupt status (the wait is not interruptible for init purposes).
  predecessor-ordering with concurrent superinterface requests is exercised at
  least once.
- Create `docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh`
  — the parent's final 19-fixture candidate harness. It asserts exactly the
  fixed 19 fixture rows (the list in `matrix.md`) for Interpreter and IR,
  fail-closed on any missing row, timeout, or stdout/stderr/exit mismatch. It
  additionally runs an AOT column asserting every one of the 19 fixtures is
  build-rejected (`NO-BUILD`) with no built binary, panic, or fallback. It builds
  catty with `-race` when `R2_CONCURRENCY_STRESS>1` (per Amendment 1's runner
  convention) and records the candidate/base/toolchain in the evidence header.
  Candidate evidence is isolated under
  `docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/` (no `slice-c`
  suffix; this is the parent final matrix), and is never written to the research
  baseline or shared/latest path.
- Verify and, where needed, enforce that AOT build-rejects all 19 fixed
  concurrency fixtures. If any fixture currently AOT-builds and later panics or
  falls back, add an explicit AOT reachability rejection for the concurrency
  execution-context boundary in this slice. **No AOT concurrency execution is
  implemented**; AOT remains `Not implemented` for all 19 fixtures (ADR-0028).
- Update `docs/PROJECT_STATUS.md` verified-capability and explicit-boundary text
  if the concurrent initialization cross-thread capability becomes the new stated
  boundary. `ARCHITECTURE.md` is already Slice-correct under Amendment 3; no
  stale single-threaded/monitor/nop claims are re-introduced.

### Slice D non-scope

- Timed `Object.wait`/`Thread.join` enforcement (actual timeout behavior beyond
  the argument validation added by Amendment 2); remains a future bounded slice.
- `Unsafe`/`VarHandle`, JMM weak-memory optimization (the heap stays race-free
  SC per ADR-0030), virtual threads, `ThreadGroup`/`ThreadLocal`, context class
  loaders, `java.util.concurrent`, deadlock detection, lock fscheduling
  fairness, priorities, monitor deflation, or biased locking.
- AOT concurrency, an AOT execution-context ABI, or any AOT `Supported`/`Fallback`
  claim for the 19 concurrency fixtures. All remain explicit `Not implemented`.
- `invokedynamic`, broad reflection, JNI, broad `java.base` replacement, or any
  change to the four ADR-0025 bytecode initialization trigger points (`new`,
  `getstatic`, `putstatic`, `invokestatic`). Array creation (`newarray`,
  `anewarray`, `multianewarray`) and `arraylength`/`aaload` are **not**
  initialization triggers per JVMS §5.5 and Verified-by-investigation on
  Temurin 25; Slice D does not add initialization calls there.
- Re-running the 11-row Slice C runner as the final gate; the 19-row candidate
  harness replaces it for parent acceptance. The sealed Slice C evidence is
  immutable and is not overwritten.
- Changes to the parent workstream's frozen Outcome/Scope/Non-scope/Semantic
  constraints/engine matrix/Acceptance gates except as recorded in the parent's
  Amendments section. Slice D refines the Plan only.

### Slice D semantic constraints

- Java Thread execution-context identity (`ecID` per ADR-0028) is the
  initialization owner key; goroutine identity and Java Thread facade pointer
  identity are not owner keys.
- `<clinit>` executes **at most once** per runtime Class identity across all
  execution contexts, and at least once unless initialization fails. The winner
  is the unique claimant that publishes the terminal state.
- A waiter on an other-owned initialization does not hold any Java monitor while
  blocked, is not spuriously woken without a terminal publication, and observes
  the initialized state under release/acquire on wake. `notify-all` is used on
  every terminal transition so all waiters re-check.
- The init wait is **not** interruptible: an interrupt received while waiting
  for another context's initialization does not throw
  `InterruptedException` from the init path and does not clear the interrupt
  flag; the wait completes when the terminal state is published. (This is the
  ADR-0029 "unchanged interrupt status while waiting" rule, distinct from the
  interruptible `Object.wait`/`Thread.sleep`/`Thread.join` semantics already
  implemented.)
- The four ADR-0025 states and transitions, the superclass/required-
  superinterface ordering, the JSON `Error`/`ExceptionInInitializerError`
  wrapping, the constant-field exemption, and the recursive same-owner rule are
  unchanged; only synchronization and wait/notify/visibility are added.
- Interpreter and IR use the identical shared `InitializeClass` service and the
  identical exception transport. Neither engine may busy-wait, silently drop
  the other-owner wait, or approximate cross-thread publication with a single
  racy read.
- AOT reports all 19 fixtures as `Not implemented`/`NO-BUILD`; a built binary
  that later panics, mismatches, or silently falls back is a Fail.

### Slice D review evidence

These are slice-review checks, not amendments to or substitutes for the parent's
final Acceptance table. Results use only `Pass`, `Fail`, `Not run`, or
`Not implemented`.

| Check | Required evidence before Slice D can be proposed Complete |
|---|---|
| Init kernel race-free | `go test -race ./rtda` covers concurrent single-class contention, exactly-once clinit, other-owner wait → published value with no clinit re-run, erroneous publication → `NoClassDefFoundError` to all waiters, waiter interrupt-status unchanged across init wait, and a predecessor-ordering concurrency case |
| Interpreter/IR fixtures | `bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh <candidate>` asserts exactly 19 rows; all match Temurin 25 combined stdout+stderr+exit in Interpreter and IR; includes `CrossThreadClassInitialization` |
| AOT rejection | The same runner asserts all 19 fixtures report `Not implemented` as `NO-BUILD`; any built binary, panic, mismatch, omitted row, or `Fallback` is Fail |
| Bounded race repetition | `R2_CONCURRENCY_STRESS=100 bash .../run-concurrency-candidate.sh <candidate>` with a race-built catty Interpreter/IR; no Go race, timeout, deadlock, missing iteration, or semantic mismatch. This is the parent's final stress gate |
| Core regression | `go vet ./...`, `go test ./...`, `go test -race ./...`, and `bash tests/run.sh` all Pass |
| Evidence isolation | The runner requires an explicit immutable candidate ID, refuses overwrite, records toolchain/base/fixture list, writes only its candidate `docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/` directory; all historical evidence (research baseline, Slice A/B/C, and the 11-row Slice C runner evidence under `<candidate>/slice-c/`) hashes remain unchanged |
| Scope audit | Diff from `ff691b5` contains no timed-wait enforcement, AOT concurrency claim, fixture denominator change, array-creation initialization trigger, or stale single-threaded claim introduced by Slice D |

The 19-row candidate harness may share code with the Slice C 11-row runner, but
it hard-codes and reports all 19 rows; the 11-row Slice C runner remains as
historical slice evidence and is not the final gate.

### Slice D implementation order

1. per-`Class` `initMu`/`initCond` and atomic state/owner accessors.
2. `InitializeClass` cross-context wait/notify/visibility protocol.
3. Race kernel unit tests (contention, exactly-once, erroneous publish, init-wait
   interrupt-status, predecessor ordering).
4. 19-row candidate harness with AOT build-rejection column.
5. Run full 19-fixture matrix in Interpreter + IR (1× and race-built 100×), fix
   any observed cross-thread gap, confirm AOT `NO-BUILD` for all 19.
6. Owner review candidate; on accept, update Plan to Complete and proceed to
   Slice E (final integration, docs, AOT matrix confirmation, push).

### Slice D acceptance record

Accepted by Owner on 2026-07-16. This acceptance authorizes production work
only within this frozen Slice D contract after a new Active Agent starts from
this acceptance-anchor commit, records the required implementation preflight,
and changes the Plan's Slice D row to `In progress`. It does not authorize
scope expansion, integration beyond the Owner's stated authority, or changes to
the frozen sections without an accepted amendment.

---

## Implementation preflight (Slice D)

Before production implementation record (must be filled by the Active Agent
before Slice D becomes `In progress`):

- **Acceptance anchor / actual base:** `c4ddde4` (resolved `c4ddde425307ab595e1a097808002d99249f8a63`) / worktree at `4a91470`
- **Historical evidence check:** `git diff --name-only a0288be..4a91470 -- docs/workstreams/r2-concurrency-evidence docs/workstreams/r2-evidence docs/workstreams/r2-initialization-evidence docs/workstreams/r2-string-evidence` — **Pass** (empty)
- **Slice A/B/C evidence intact:** `9576828/`, `505d3ee/`, `a0e336c/`, `b0a7b70/`, `0a96b59/`, `36a577c/`, `f0fc2ca/`, `eea253d/` — **Pass**
- **Candidate evidence destination:** `docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/` (parent final 19-fixture matrix; no `slice-c` suffix here)
- **Harness output policy:** explicit candidate required; never writes research baseline, shared/latest, or any `slice-c/` directory

Any missing item keeps the workstream `Accepted`; it may not become `In progress`.

## (Roadmap note — non-binding, for Owner's future consideration)

Owner-directed roadmap candidate, **not part of this workstream**: AOT
concurrency should be planned under Phase R5 (AOT coverage expansion), gated by
a future Proposed→Accepted ADR defining an explicit AOT execution-context ABI,
and sequenced after R5's instance-method AOT and exception handling in emitted
code. Until that ABI is Accepted, all 19 concurrency fixtures remain AOT
build-rejected. This note does not authorize implementation; it records the
investigation conclusion for ROADMAP editing at the Owner's discretion.
