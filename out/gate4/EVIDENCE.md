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
before=51b27daebabcdd89ee971abf7fb7e4d8c1e3e4e0e4b1cf9f68ecf81e0bb48e99
after =51b27daebabcdd89ee971abf7fb7e4d8c1e3e4e0e4b1cf9f68ecf81e0bb48e99
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
{"daemon": false, "event": "thread_start", "name": "main", "pid": 30769, "thread": 2}
{"daemon": false, "event": "thread_start", "name": "counted", "parentThread": 14, "pid": 30769, "thread": 15}
{"event": "thread_end", "pid": 30769, "result": "terminated", "thread": 15}
{"event": "thread_end", "pid": 30769, "result": "terminated", "thread": 14}
{"event": "thread_end", "pid": 30769, "result": "terminated", "thread": 2}
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
{"depth": 4, "event": "monitor", "method": {"descriptor": "()V", "loaderId": 1, "name": "nested", "owner": "Monitors"}, "object": 77, "pc": 3, "phase": "enter", "pid": 30828, "thread": 2}
{"depth": 2, "descriptor": "()V", "event": "unwind", "loaderId": 1, "name": "nested", "owner": "Monitors", "pid": 30828, "thread": 2, "throwable": 1, "throwableLoaderId": 0}
{"descriptor": "([Ljava/lang/String;)V", "event": "catch", "handlerPC": 22, "loaderId": 1, "name": "main", "owner": "Monitors", "pid": 30828, "thread": 2, "throwPC": 9, "throwable": 1, "throwableLoaderId": 0}
```

### Wait outcomes (`traces/wait-set-interp-normal.jsonl`)

One resume per waiter, each with a single reason, and full recursion depth
restored (depth 2 waiters resume at depth 2):

```json
{"depth": 2, "event": "monitor", "object": 91, "phase": "wait", "pid": 30887, "thread": 7}
{"depth": 2, "event": "monitor", "object": 91, "phase": "wait", "pid": 30887, "thread": 8}
{"depth": 1, "event": "monitor", "object": 91, "phase": "wait", "pid": 30887, "thread": 2}
{"depth": 2, "event": "monitor", "object": 91, "phase": "resume", "pid": 30887, "reason": "interrupted", "thread": 7}
{"depth": 2, "event": "monitor", "object": 91, "phase": "resume", "pid": 30887, "reason": "notified", "thread": 8}
{"depth": 1, "event": "monitor", "object": 91, "phase": "resume", "pid": 30887, "reason": "notified", "thread": 2}
```

### Class initialization (`traces/initialization-interp-normal.jsonl`)

The owner starts, a competing thread records a wait with the owner identity,
and completion is published before the waiter proceeds:

```json
{"class": "Initialization$Ready", "event": "class_init", "loaderId": 1, "phase": "start", "pid": 30946, "thread": 5}
{"class": "Initialization$Ready", "event": "class_init_wait", "loaderId": 1, "ownerThread": 5, "pid": 30946, "thread": 6}
{"class": "Initialization$Ready", "event": "class_init", "loaderId": 1, "phase": "complete", "pid": 30946, "thread": 5}
{"class": "Initialization$Broken", "event": "class_init", "failed": true, "loaderId": 1, "phase": "complete", "pid": 30946, "thread": 9, "throwable": 2}
```

### Concurrent plugin admission (`traces/plugin-interp-normal.jsonl`)

Two Java workers register identical bytes concurrently; the runtime assigns
distinct loader IDs and publishes each class once. Their plugin threads use
distinct class mirrors as the static synchronized monitor:

```json
{"deferredReferences": 0, "definedClasses": 1, "event": "plugin_admission", "interpretedMethods": 5, "jarSHA256": "417ea4a59f8c500afa0c49232d2ea75e61388eccd2d54a2abbd9bb1706a18481", "loaderId": 2, "parentLoaderId": 1, "pid": 31134, "result": "accepted", "shadowedClasses": 0, "thread": 6}
{"class": "ext/Worker", "event": "class_define", "jarSHA256": "417ea4a59f8c500afa0c49232d2ea75e61388eccd2d54a2abbd9bb1706a18481", "loaderId": 2, "pid": 31134, "thread": 6}
{"deferredReferences": 0, "definedClasses": 1, "event": "plugin_admission", "interpretedMethods": 5, "jarSHA256": "417ea4a59f8c500afa0c49232d2ea75e61388eccd2d54a2abbd9bb1706a18481", "loaderId": 3, "parentLoaderId": 1, "pid": 31134, "result": "accepted", "shadowedClasses": 0, "thread": 5}
{"class": "ext/Worker", "event": "class_define", "jarSHA256": "417ea4a59f8c500afa0c49232d2ea75e61388eccd2d54a2abbd9bb1706a18481", "loaderId": 3, "pid": 31134, "thread": 5}
{"depth": 1, "event": "monitor", "method": {"descriptor": "()V", "loaderId": 3, "name": "record", "owner": "ext/Worker"}, "object": 80, "phase": "enter", "pid": 31134, "thread": 10}
{"depth": 1, "event": "monitor", "method": {"descriptor": "()V", "loaderId": 2, "name": "record", "owner": "ext/Worker"}, "object": 79, "phase": "enter", "pid": 31134, "thread": 8}
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

Timing-dependent tests use the injectable monotonic clock
(`Runtime.clock`); interleavings use barriers and channels.

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
