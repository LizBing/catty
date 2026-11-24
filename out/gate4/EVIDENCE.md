# Gate 4 acceptance evidence

Gate 4 (`catty-jvm21-gate4`, image ABI `4`) is implemented and accepted on
`codex/gate-4`. Both engines share the same logical thread identities,
monitors, wait sets, class-initialization protocol and memory operations; the
Gate 2 and Gate 3 profiles and their suites remain unchanged and pass.

Environment (full versions in [TOOLCHAINS.txt](TOOLCHAINS.txt)):

- go1.27.1 darwin/arm64, CGO enabled, Go race detector available.
- Temurin OpenJDK 21.0.13.1+2 (`java`, `javac`, `jar`).
- Python 3.9.6; all acceptance runs use a scrubbed environment with no JDK on
  `PATH` for Catty executables.

## Command statuses

Every command below was run from the repository root with the toolchain above.
Full logs are under [logs/](logs/); the status table is
[STATUS.txt](STATUS.txt).

| # | Command | Status | Log |
| - | - | - | - |
| 1 | `go test -race -count=1 ./...` | 0 | [go-test-race.log](logs/go-test-race.log) |
| 2 | `go vet ./...` | 0 | [go-vet.log](logs/go-vet.log) |
| 3 | `python3 -B tests/gate4.py --java-home <jdk-21> --prepare-only` | 0 | [gate4-reference.log](logs/gate4-reference.log) |
| 4 | `python3 -B tests/gate4.py --java-home <jdk-21> --race` | 0 | [gate4-race.log](logs/gate4-race.log) |
| 5 | `python3 -B tests/gate3.py --java-home <jdk-21>` | 0 | [gate3.log](logs/gate3.log) |
| 6 | `python3 -B tests/gate2.py --java-home <jdk-21>` | 0 | [gate2.log](logs/gate2.log) |
| 7 | `go test -race -count=1 -v -run '<follow-up timer regex>' ./rt/` | 0 | [followup-timer-tests.log](logs/followup-timer-tests.log) |
| 8 | `go test -race -count=1 -v -run '<blocker regressions>' ./rt/` | 0 | [blocker-regressions.log](logs/blocker-regressions.log) |

The Gate 4 runner compiled the fixtures, recorded the OpenJDK 21 oracle and
exercised the same class files through all-interpreted, mixed and forced-AOT
images, then repeated all three modes with race-instrumented executables and
tracing disabled. Case results (from [gate4-race.log](logs/gate4-race.log)):

```text
lifecycle passed
monitors passed
wait-set passed
initialization passed
memory passed
launcher passed
plugin passed
Gate 4 differential tests passed.
```

The runner also passed the profile negatives it owns: the bad-monitor-operand
class file is rejected at the mutated method and PC without publishing build
artifacts, and both earlier profiles still reject the synchronized-method and
explicit-monitor fixtures. `--prepare-only` independently re-passed every
fixture against OpenJDK 21.

Race instrumentation is present in the generated executables themselves, not
only in `go test`; every race image reports `-race=true` in Go build metadata:

```text
images/monitors-interp-race/app: go1.27.1
	build	-race=true
	build	GOARCH=arm64
	build	GOOS=darwin
```

## Image and profile evidence

- Build reports for all 42 successful images (7 cases x 3 modes, plus the same
  modes with race instrumentation) are in [reports/](reports/); the rejected
  `bad-monitor-operand-interp` build is covered by the runner and publishes no
  report. Each records profile
  `catty-jvm21-gate4`, image ABI `4`, `threads: true`, `plugins: true`, the
  Gate 3 plugin limits, capabilities
  `plugin-admission, loader-identity, guarded-dispatch, threads, monitors,
  synchronized-memory`, and an opcode inventory containing exactly the Gate 2
  and Gate 3 opcodes plus `monitorenter` and `monitorexit`.
- Representative method selections show main AOT-compiled with the mixed
  policy's forced interpreter method, for example `Monitors.main` compiled
  (`aot: lowered to Go`) and `Monitors.nested` retained as
  `policy: forced interp`.
- Canonical generated Go for the Gate 4 images is under
  [generated/](generated/). `monitors-aot` contains the monitor lowering
  (`t.MonitorEnter(m, pc, stack[...])` and `t.MonitorExit(m, pc, stack[...])`)
  with the original instruction PCs feeding `t.HandleThrow`.
- [generated/INSPECTION.txt](generated/INSPECTION.txt) is the output of the
  existing AST inspector extended with a `monitors` check:

  ```sh
  go run tests/inspect_generated.go cattyMethod9 monitors \
    out/gate4/generated/monitors-aot/generated.go
  ```

  It verifies that every monitor helper call carries the original PC and that
  the same PC feeds the shared exception-transfer path.
- Class-file and JAR SHA-256 digests, generated executable digests (including
  the race images) and the built launcher digest are in
  [DIGESTS.txt](DIGESTS.txt); the compiled fixture classes and JARs are copied
  under [classes/tree](classes/tree) and [jars/](jars/).

## Standalone plugin isolation

After building the plugin image, the binary and its JAR were copied into an
isolated directory and executed with an empty environment and no fixture class
directories or adapter on any path:

```text
status=0
before=6da2523db0e489354f1f8f64f33c95c880b9e5bfe7e08ebe3a1f9ffba15daca2
after =6da2523db0e489354f1f8f64f33c95c880b9e5bfe7e08ebe3a1f9ffba15daca2
stdout=plugin
```

Artifacts: [standalone/](standalone/) (`app`, `plugin-worker.jar`, trace and
outputs). The executable digest is unchanged by the run.

## Trace excerpts

Full traces are under [traces/](traces/). They are JSON lines in emission
order; each file is a complete trace from one image.

### Thread lifecycle (`traces/lifecycle-interp-normal.jsonl`)

Main is a root; workers carry `parentThread`; every non-daemon thread has one
start and one end:

```json
{"daemon": false, "event": "thread_start", "name": "main", "pid": 51707, "thread": 2}
{"daemon": false, "event": "thread_start", "name": "counted", "parentThread": 13, "pid": 51707, "thread": 15}
{"event": "thread_end", "pid": 51707, "result": "terminated", "thread": 15}
{"event": "thread_end", "pid": 51707, "result": "terminated", "thread": 14}
{"event": "thread_end", "pid": 51707, "result": "terminated", "thread": 2}
```

The launcher case trace (`traces/launcher-interp-worker-uncaught.jsonl`) ends
with the worker's `result: "uncaught"` and a throwable identity while main
ends `terminated`; the main-uncaught trace ends main with
`result: "uncaught"`.

### Nested monitors (`traces/monitors-interp-excerpt.jsonl`)

Explicit and implicit acquisitions share one monitor identity, depth grows
through synchronized methods and the explicit block, and the throwable
identity is preserved through unwind and catch:

```json
{"depth": 4, "event": "monitor", "method": {"descriptor": "()V", "loaderId": 1, "name": "nested", "owner": "Monitors"}, "object": 77, "pc": 3, "phase": "enter", "pid": 51766, "thread": 2}
{"depth": 2, "descriptor": "()V", "event": "unwind", "loaderId": 1, "name": "nested", "owner": "Monitors", "pid": 51766, "thread": 2, "throwable": 1, "throwableLoaderId": 0}
{"descriptor": "([Ljava/lang/String;)V", "event": "catch", "handlerPC": 22, "loaderId": 1, "name": "main", "owner": "Monitors", "pid": 51766, "thread": 2, "throwPC": 9, "throwable": 1, "throwableLoaderId": 0}
```

### Wait outcomes (`traces/wait-set-interp-normal.jsonl`)

One resume per waiter, each with a single reason, and full recursion depth
restored (depth 2 waiters resume at depth 2):

```json
{"depth": 2, "event": "monitor", "object": 91, "phase": "wait", "pid": 51825, "thread": 7}
{"depth": 2, "event": "monitor", "object": 91, "phase": "wait", "pid": 51825, "thread": 8}
{"depth": 1, "event": "monitor", "object": 91, "phase": "wait", "pid": 51825, "thread": 2}
{"depth": 2, "event": "monitor", "object": 91, "phase": "resume", "pid": 51825, "reason": "interrupted", "thread": 7}
{"depth": 2, "event": "monitor", "object": 91, "phase": "resume", "pid": 51825, "reason": "notified", "thread": 8}
{"depth": 1, "event": "monitor", "object": 91, "phase": "resume", "pid": 51825, "reason": "notified", "thread": 2}
```

### Class initialization (`traces/initialization-interp-normal.jsonl`)

The owner starts, a competing thread records a wait with the owner identity,
and completion is published before the waiter proceeds:

```json
{"class": "Initialization$Ready", "event": "class_init", "loaderId": 1, "phase": "start", "pid": 51887, "thread": 5}
{"class": "Initialization$Ready", "event": "class_init_wait", "loaderId": 1, "ownerThread": 5, "pid": 51887, "thread": 6}
{"class": "Initialization$Ready", "event": "class_init", "loaderId": 1, "phase": "complete", "pid": 51887, "thread": 5}
{"class": "Initialization$Broken", "event": "class_init", "failed": true, "loaderId": 1, "phase": "complete", "pid": 51887, "thread": 9, "throwable": 2}
```

### Concurrent plugin admission (`traces/plugin-interp-normal.jsonl`)

Two Java workers register identical bytes concurrently; the runtime assigns
distinct loader IDs and publishes each class once. Their plugin threads use
distinct class mirrors as the static synchronized monitor:

```json
{"deferredReferences": 0, "definedClasses": 1, "event": "plugin_admission", "interpretedMethods": 5, "jarSHA256": "d969024fe13f71d978e21b020d089ff62ad069cd7dc0edc8619d8a4c7cf3a671", "loaderId": 2, "parentLoaderId": 1, "pid": 52082, "result": "accepted", "shadowedClasses": 0, "thread": 5}
{"class": "ext/Worker", "event": "class_define", "jarSHA256": "d969024fe13f71d978e21b020d089ff62ad069cd7dc0edc8619d8a4c7cf3a671", "loaderId": 2, "pid": 52082, "thread": 5}
{"deferredReferences": 0, "definedClasses": 1, "event": "plugin_admission", "interpretedMethods": 5, "jarSHA256": "d969024fe13f71d978e21b020d089ff62ad069cd7dc0edc8619d8a4c7cf3a671", "loaderId": 3, "parentLoaderId": 1, "pid": 52082, "result": "accepted", "shadowedClasses": 0, "thread": 6}
{"class": "ext/Worker", "event": "class_define", "jarSHA256": "d969024fe13f71d978e21b020d089ff62ad069cd7dc0edc8619d8a4c7cf3a671", "loaderId": 3, "pid": 52082, "thread": 6}
{"depth": 1, "event": "monitor", "method": {"descriptor": "()V", "loaderId": 3, "name": "record", "owner": "ext/Worker"}, "object": 80, "phase": "enter", "pid": 52082, "thread": 9}
{"depth": 1, "event": "monitor", "method": {"descriptor": "()V", "loaderId": 2, "name": "record", "owner": "ext/Worker"}, "object": 79, "phase": "enter", "pid": 52082, "thread": 10}
```

## Deterministic Go coverage

`go test -race -count=1 ./...` covers the implementation-level transitions in
the test plan without depending on wall-clock scheduling:

- `rt/gate4_runtime_test.go`: reentrant monitor ownership and contention,
  non-owner exit, null monitor operands, monitored-object collectability,
  implicit method entry and unwind, same-frame catch while holding the monitor,
  per-receiver and per-loader static lock identity, wait enrollment with
  notify/notifyAll/timeout/interrupt winners, late interrupts, restored
  recursion depth and unrelated monitors, timer overloads and saturation,
  single-deadline timed joins with unrelated wakeups, sleep interruption with
  monitors retained, first-start and setDaemon/start races, multiple joiners,
  uncaught worker diagnostics and trace, descendant accounting after main,
  class initialization contention/failure/wrapper/parent-failure/order and
  per-loader identity, resolution and dispatch cache races, metadata interning,
  concurrent admission quota and rollback, detached admission snapshots, all
  shared-memory slot/array operations, start/join publication, monitor release
  on thread termination, and trace close during daemon activity.
- `internal/verify/gate4_test.go`: monitor operand typing, earlier-profile
  rejection, synchronized constructor/abstract/interface/native rejection,
  valid synchronized static methods, `<clinit>` synchronization-flag handling,
  and final+volatile field rejection under Gate 4 only.
- `internal/aot/emit_test.go`: generated monitor lowering uses the shared
  helpers with original PCs and the shared handler-transfer path.
- Follow-up timer tests: early-notified and interrupted long waits, long joins
  ended by target termination, timeout cleanup, a stale deadline racing a later
  wait, bounded wait/join segments, Java timeouts beyond Go's duration range,
  interrupted sleep cleanup, termination publication under the peer monitor,
  and 32 real-clock notified waits showing no goroutine accumulation.
- Second follow-up regressions: wait-set backing-array tail clearing for middle
  and end removals, direct deadline arithmetic just beyond Go's duration
  ceiling, and timed wait/join expiry at that deadline.

Timing-dependent tests use the injectable monotonic clock
(`Runtime.clock`); interleavings use barriers and channels.

## Follow-up: timer ownership and requirement audit

The merge-blocking review finding was confirmed in `rt/monitor.go`.
`startWaitTimer` spawned a goroutine that waited only for its deadline. When
notification or interruption woke the Java wait first, the waiter record was
removed, but that goroutine remained blocked on the timer and retained the
waiter and monitor until the Java deadline. Timed joins used the same helper.
A related leak existed for `Thread.sleep`: its helper goroutine exited on
interruption, but the uncanceled `time.After` timer stayed armed until its
deadline.

The fix makes every Java timer own a cancelable host timer:

- `runtimeClock` now exposes `NewTimer` returning a `runtimeTimer` with `C`
  and `Stop`, implemented with `time.Timer` for production and a manual
  stop/fire timer for deterministic tests.
- `waiter` carries its current timer segment and a `timerDone` channel. The
  single winning transition (`wakeWaiterLocked`) stops the segment, closes
  `timerDone`, and drops both references. The timer goroutine exits through
  `timerDone`; a stale callback can never observe a pending record, so it
  cannot affect a later wait.
- Timed waits and joins use one absolute millisecond `waitTimer`. The
  remaining duration is recomputed from that deadline on every segment and
  wakeup, so a timed join keeps one monotonic deadline across unrelated
  notifications and cannot extend under repeated wakeups. Each host timer is
  capped at `maxTimerSegmentMillis` (production: the largest Go duration;
  tests lower it), which walks Java timeouts beyond Go's representable range
  in bounded segments.
- `sleepFor` now selects directly on its own host timer and the interrupt
  channel; an interrupted sleep stops the timer, and a timer that fires first
  makes a later interrupt only leave its status pending.
- `publishTermination` releases the terminating thread's monitors, silently
  acquires the target's peer monitor, publishes the terminated state and
  notifies all joiners while owning it, then releases it. Publication is
  silent because the trace lifetime already ended, and it is shared by worker
  and main termination.

Focused tests added in `rt/gate4_runtime_test.go` (all run under `-race`):

- `TestGate4EarlyWaitReleasesTimers` (32 long notified waits)
- `TestGate4RealClockEarlyWaitsDoNotParkTimers` (32 production-clock waits)
- `TestGate4EarlyJoinReleasesTimers` (16 long joins ended by termination)
- `TestGate4InterruptedWaitReleasesTimers`,
  `TestGate4InterruptedJoinReleasesTimers`,
  `TestGate4TimeoutReleasesTimers`,
  `TestGate4InterruptedSleepReleasesTimer`
- `TestGate4StaleTimerDoesNotAffectLaterWait`
- `TestGate4BoundedWaitAndJoinSegments`,
  `TestGate4HugeWaitSegmentsAndCancellation`
- `TestGate4TerminationPublishesUnderPeerMonitor`

The follow-up log records the original leak measurement now being clean:

```text
gate4_runtime_test.go:2593: 32 notified long waits: armed timers=0 timer workers=0
gate4_runtime_test.go:3036: 32 real-clock notified waits: goroutines before=2 after=2 timer workers=0
```

### Related requirement audit

- **Bounded timer segments.** The previous wait/join path saturated positive
  timeouts at `maxTimerDuration` (about 292 years), so a Java timeout larger
  than Go's duration range could fire early. The new `waitTimer` keeps the full
  int64 millisecond timeout and walks it in bounded segments; a positive
  nanosecond remainder rounds up to one millisecond with saturation, while
  zero stays indefinite. `TestGate4BoundedWaitAndJoinSegments` and
  `TestGate4HugeWaitSegmentsAndCancellation` exercise the segment path with a
  lowered cap, and `waitTimeoutMillis` covers the arithmetic boundaries.
- **Termination publication under the peer monitor.** The previous
  implementation woke joiners under `thrMu` without owning the target's peer
  monitor. `publishTermination` now acquires that monitor for the publication
  and notification window, and `TestGate4TerminationPublishesUnderPeerMonitor`
  proves a thread that holds the peer monitor keeps the target observable as
  started and prevent a joiner from returning until it releases the monitor.

### Follow-up commands

All commands below were run on the follow-up tree and exited 0; full output is
in the linked logs.

```sh
go test -race -count=1 ./...
go vet ./...
go test -race -count=1 -v \
  -run 'TestGate4EarlyWaitReleasesTimers|TestGate4RealClockEarlyWaitsDoNotParkTimers|TestGate4EarlyJoinReleasesTimers|TestGate4InterruptedWaitReleasesTimers|TestGate4InterruptedJoinReleasesTimers|TestGate4TimeoutReleasesTimers|TestGate4StaleTimerDoesNotAffectLaterWait|TestGate4BoundedWaitAndJoinSegments|TestGate4HugeWaitSegmentsAndCancellation|TestGate4InterruptedSleepReleasesTimer|TestGate4TerminationPublishesUnderPeerMonitor' \
  ./rt/
python3 -B tests/gate4.py --java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home --prepare-only
python3 -B tests/gate4.py --java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home --race
python3 -B tests/gate3.py --java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home
python3 -B tests/gate2.py --java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home
```

## Second follow-up: wait-set retention and overflow deadlines

### Removed waiters release the backing array

`removeWaiterLocked` previously shortened the wait-set slice with
`append(waiters[:i], waiters[i+1:]...)`, which shifted the later records but
left the last pointer in the unused tail. A long-lived monitor therefore kept
the completed waiter and its Java thread reachable through the backing array
until the slice was replaced.

The fix copies the surviving tail, writes `nil` into the vacated last slot, and
then shortens the slice. `TestGate4RemoveWaiterReleasesBackingArray` scans the
full capacity after removing from the middle and from the end, and rejects any
retained pointer to the removed record.

### Deadlines beyond Go's duration range

`waitTimer.remaining` computed elapsed time with `now.Sub(start)`, which
saturates at `MaxInt64` nanoseconds (about 292 years). For a Java timeout 10 ms
beyond the largest whole-millisecond Go duration, advancing the manual clock to
the exact deadline still returned `remaining=10ms, expired=false`, so the wait
could never expire at its real deadline.

The fix keeps the full int64 millisecond timeout. Elapsed time uses the
monotonic `Sub` reading while that reading is representable, and falls back to
whole Unix seconds plus nanoseconds (`wallElapsedMillis`) once the monotonic
difference saturates. The deadline remains one absolute Java time; join wakeups
recompute from the same start and millisecond timeout, so they cannot extend it.
This preserves monotonic timing for every representable wait, and exact
millisecond accounting for the full Java range.

Regressions, all passing under `-race`:

- `TestGate4WaitTimerCrossesDurationCeiling` — at the exact age of
  `largestWholeMillisecondTimeout()+10ms`, `remaining` is `0` and
  `expired` is true (the pre-fix behavior was `10ms,false`).
- `TestGate4TimedWaitCrossesDurationCeiling` — a timed `Object.wait` with that
  timeout is still waiting at the duration ceiling, then resumes at its exact
  original deadline.
- `TestGate4TimedJoinCrossesDurationCeiling` — the same deadline for a timed
  join against a live target, with no early exit and no timer work left after
  expiry.

### Commands for this follow-up

All exited 0; full output is in the linked logs and `STATUS.txt`.

```sh
go test -race -count=1 ./...
go vet ./...
go test -race -count=1 -v \
  -run 'TestGate4RemoveWaiterReleasesBackingArray|TestGate4WaitTimerCrossesDurationCeiling|TestGate4TimedWaitCrossesDurationCeiling|TestGate4TimedJoinCrossesDurationCeiling' \
  ./rt/
python3 -B tests/gate4.py --java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home --prepare-only
python3 -B tests/gate4.py --java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home --race
python3 -B tests/gate3.py --java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home
python3 -B tests/gate2.py --java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home
```

## Known limits / remaining gaps

- Gate 4 stays opt-in (`--enable-threads`); Gate 2 remains the default profile
  and Gate 3 stays selectable with `--enable-plugins`. This is the published
  contract, not a gap.
- Outside the profile by specification: virtual threads, thread groups, thread
  locals, executors, park/unpark, Java-level thread inspection, native methods,
  agents, JDK synchronization utilities, and VM exit/shutdown-hook/application
  measurement (Gate 5).
- The `monitors` AST check is exercised here from the evidence step and by the
  emitter unit test; the checked-in Python runner does not invoke the inspector
  itself.
- A failed admission may advance the internal unobservable identity counter
  (Gate 4 permits it); loader IDs and quota accounting never advance on
  rejection.
- Class-initialization start/wait/completion trace records are emitted while
  the per-class guard is held so the file demonstrates the required order; the
  trace connection itself takes only the trace mutex and never a memory,
  metadata or thread mutex.
