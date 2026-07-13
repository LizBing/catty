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

## Acceptance record

Accepted by Owner on 2026-07-14. This acceptance authorizes production work only
within this frozen contract after a new Active Agent starts from the acceptance
anchor and records the required implementation preflight. It does not authorize
scope expansion, integration beyond the Owner's stated authority, or changes to
the frozen sections without an accepted amendment.

---

## Implementation preflight

Before production implementation record:

- **Acceptance anchor / actual base:** `<commits; worktree must descend from anchor>`
- **Historical evidence check:** `<exact command and exit status>`
- **Candidate evidence destination:** `docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/`
- **Harness output policy:** explicit candidate required; never writes research baseline or shared/latest path

Any missing item keeps the workstream `Accepted`; it may not become `In Progress`.

---

## Plan

| Slice | Status | Evidence |
|---|---|---|
| A — SC heap cells, concurrency-safe loader, and canonical Class mirrors | Pending | — |
| B — stable Thread facade/context, lifecycle, carriers, join, and VM liveness | Pending | — |
| C — monitors, synchronized methods, wait sets, and interruption | Pending | — |
| D — concurrent ADR-0025 initialization and full Interpreter/IR fixture matrix | Pending | — |
| E — AOT fail-closed rejection, race stress, regression, evidence, and docs | Pending | — |

Status uses `Pending`, `In progress`, or `Complete`.

---

## Handoff

- **Branch / candidate:** `main`; no implementation candidate yet
- **Acceptance anchor / base:** this 2026-07-14 governance commit / research baseline `63d5658`
- **Dirty files:** research and planning artifacts only
- **Historical evidence check:** research baseline preserved additively; final authoritative run is `run-concurrency-results-v5.txt`
- **Candidate evidence path:** not created
- **Last location:** contract drafted from completed research findings
- **Checks run / not run:** implementation gates Not run
- **Blocker:** implementation preflight and a new Active Agent are required before production work
- **Next action:** start a new Active Agent from this acceptance anchor; record its resolved SHA/base and historical-evidence check before production implementation
- **Non-derivable context:** the 19-fixture denominator includes explicit daemon and non-daemon liveness, all three interruptible blocking points, and the producer-consumer milestone
