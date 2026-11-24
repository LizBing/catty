# Gate 5 evidence

Branch: `codex/gate-5`

| Revision | Contents |
| --- | --- |
| `553050a` | Planning baseline (Gate 5 specifications and plans) |
| `f28e8fc` | Gate 5 profile, lifecycle/shutdown, diagnostics, metrics, runtime tests, differential suite |
| `cec09d2` | Gate 5 application benchmark workload and runner |
| `7469e36` | Gate 5 race-enabled acceptance artifacts (this directory's `acceptance/`) |
| `a9dcee0` | Generated-body inspection checks, executable digest stability and their retained evidence |
| `651f706` | Review fixes: live shutdown keeper, shutdown_begin ordering, nil-image setup error, hook daemon capture, contention-at-block |
| `af643e7` | Close the halt-to-hook-start publication race and add its deterministic regression |

The measurement run recorded in §4 was taken from the clean `7469e36` revision
(`environment.repository.dirty = false`). The review fixes in `651f706` change
runtime code after that measurement; the full four-mode benchmark results stay
attributed to `7469e36` and were not rerun. Post-fix validation in §6 used the
benchmark smoke matrix only. The benchmark workload and measured executables in
§4 are unchanged from their recorded run.

## 1. Delivered surface

- `catty build --enable-lifecycle` selects `catty-jvm21-gate5`, image ABI `5`,
  and inherits Gate 4 threading and Gate 3 plugin admission. Precedence is
  `--enable-lifecycle` > `--enable-threads` > `--enable-plugins`; the default
  remains Gate 2.
- The image and build report serialize `lifecycle: true`; installation rejects
  a mismatched profile, ABI or lifecycle/plugin/threads flag before Java runs.
  An omitted flag decodes as `false`, so Gate 2–4 images are unchanged.
- Published lifecycle members: `System.exit(I)V`, `System.nanoTime()J`,
  `Runtime.getRuntime()`, `Runtime.exit/halt/addShutdownHook/
  removeShutdownHook/gc`, and `java.lang.IllegalStateException`.
- Lifecycle: running / frozen-shutdown / terminal states under the thread
  lifecycle guard; natural shutdown at the first zero count of main plus live
  non-daemon peers; explicit exit snapshots the identity-keyed hook registry
  once; halt skips or cuts short graceful shutdown. Successful exit/halt calls
  park without unwinding frames or releasing monitors, and the launcher waits
  for a terminal result independently of Java main.
- Diagnostics: `app_start`, `shutdown_begin`, sealing `process_exit`, VM-started
  hook `thread_start` markers, named-thread uncaught diagnostics, and a host
  fault report that captures Java frames before preserving the original panic.
- `CATTY_METRICS=<path>` writes one terminal JSON document with publication,
  engine-crossing, reflection, exception, monitor, live/peak Java-thread and
  process Go heap/GC counters. Both telemetry channels are opt-in and never use
  application stdout.
- `tests/gate5.py` (OpenJDK 21 differential suite) and
  `tests/gate5_bench.py` (four-configuration measurement runner).

## 2. Validation commands and outcomes

All commands were run from `/Users/services/Projects/catty` with
`PATH=/usr/local/go/bin:$PATH`. Java runs used
`--java-home /Users/services/jdks/jdk-21.0.13.1+2/Contents/Home`
(Temurin OpenJDK 21.0.13.1+2).

| Command | Outcome |
| --- | --- |
| `go vet ./...` | exit 0 |
| `go test -race -count=1 ./...` | `ok catty/rt`, all internal packages `ok` |
| `go test -count=1 ./...` | pass |
| `python3 -B tests/gate5.py --java-home ... --prepare-only` | `Gate 5 OpenJDK 21 reference prepared.` |
| `python3 -B tests/gate5.py --java-home ... --race --keep` | `Gate 5 differential tests passed.`; retained artifacts used here |
| `python3 -B tests/gate5.py --java-home ...` | `Gate 5 differential tests passed.` |
| `python3 -B tests/gate4.py --java-home ... --race` | `Gate 4 differential tests passed.` |
| `python3 -B tests/gate3.py --java-home ...` | `Gate 3 passed.` |
| `python3 -B tests/gate2.py --java-home ...` | `Gate 2 passed.` (includes Gate 1 regression) |
| `python3 -B tests/gate5_bench.py --java-home ... --smoke` | `Gate 5 application benchmark (smoke): OK` |
| `python3 -B tests/gate5_bench.py --java-home ... --out /tmp/gate5-bench-final` | `status: ok`, `failures: none`, 29.8 s |

The complete retained logs for the `--race --keep` acceptance run are in
[`acceptance/logs/`](acceptance/logs/), including every build, OpenJDK oracle,
Catty launch, profile refusal and setup-failure subprocess. The per-case
`--race` run also executes race-instrumented generated executables for every
mode and checks `go version -m` for `-race=true`; no Go data race was reported.

## 3. Acceptance coverage

`tests/gate5.py` compiles each fixture once with `javac --release 21 -g`, checks
the OpenJDK 21 oracle, and runs the same class files in `interp`, mixed `auto`
with an explicit policy, and forced `aot` images. Exact stdout/status,
diagnostics, required markers and trace order are asserted.

| Case | Evidence |
| --- | --- |
| `startup` | args `[]`, `""`, whitespace, `-x`, Unicode; `clinit` before `main`; mixed AOT main → interpreted `Startup.step` → AOT `Startup$Box.bump` on one logical thread |
| `natural` | non-daemon drain with a descendant and a main joiner, blocked daemon left live, both main outcomes (`0`, `1`), hook after last worker; untimed `Object.wait` hook reaches the watchdog timeout with its marker; coordinated halt of that hook ends with status `31` |
| `registration` | singleton, add/remove/absent/duplicate/live/retired peers, identity-equal hooks, frozen-registry add/remove `IllegalStateException` including null |
| `hooks` | two-hook barrier, daemon hook, virtual `run` overrides, identity via `Thread.currentThread`, throwing hook with one named-frame diagnostic, status preserved, hook-started worker left live |
| `exit` | worker and reflective exits, exit during main-class initialization, statuses `0`, `7`, `-3` and `300` (POSIX `0/7/253/44`), no `finally`/after markers, monitor-holding caller observed live, coordinated `halt(23)` |
| `exit-contention` | first exit status wins over a second explicit exit; exit from a hook stalls without post-exit code; a second hook releases it with `halt(21)` |
| `halt` | halt in running state skips hooks/shutdown event (`19`); halt during graceful shutdown cuts it short (`23`) |
| `late-hook` | image built before the plugin JAR exists; plugin hook runs interpreted under the published plugin loader on its own thread and process exits `5`; forced AOT still applies to image methods |
| `setup-failure` | patched lifecycle flag exits `2` before Java; unusable `CATTY_TRACE`/`CATTY_METRICS` path exits `2`; host-fault fixture reports `Fault.main` Java frames plus `panic: injected host fault` Go stack |
| `telemetry` | mixed crossing, reflection success + miss, admission rejection, matched exception, monitor contention, thread publication, explicit GC; metrics counters and terminal fields asserted in all modes; disabled run has identical stdout |

Representative traces (captured from race-instrumented executables):

- [`traces/startup-mixed.jsonl`](traces/startup-mixed.jsonl) — `app_start`
  first, nested AOT/interp/AOT frames on one thread.
- [`traces/natural-mixed-main-uncaught.jsonl`](traces/natural-mixed-main-uncaught.jsonl)
  — natural `shutdown_begin`, hook `thread_start` with `shutdownHook: true`,
  uncaught main `thread_end`, terminal `process_exit`.
- [`traces/registration-mixed.jsonl`](traces/registration-mixed.jsonl) and
  [`traces/hooks-mixed.jsonl`](traces/hooks-mixed.jsonl) — frozen hook registry,
  concurrent hooks, named uncaught hook, live non-hook thread open at exit.
- [`traces/exit-mixed-reflect-neg3.jsonl`](traces/exit-mixed-reflect-neg3.jsonl)
  — reflective exit requested by the main thread, signed trace status `-3`;
  OS status observed as `253`.
- [`traces/exit-mixed-monitor-halt.jsonl`](traces/exit-mixed-monitor-halt.jsonl)
  — exit caller parked with its monitor, `shutdown_begin`, then terminal
  `halt`/`23` published by another hook.
- [`traces/exit-contention-mixed-halt.jsonl`](traces/exit-contention-mixed-halt.jsonl)
  — first exit snapshot, hook calling exit, coordinated halt.
- [`traces/halt-mixed-during.jsonl`](traces/halt-mixed-during.jsonl) and
  [`traces/late-hook-mixed.jsonl`](traces/late-hook-mixed.jsonl) — halt during
  shutdown; plugin hook loader/engine evidence.
- [`traces/gate4-late-admission-rejected.jsonl`](traces/gate4-late-admission-rejected.jsonl)
  — a late class referencing `System.nanoTime` is rejected at the Gate 4
  admission boundary.

Each parsed trace asserts: one `app_start` with profile/ABI/argc/main thread,
per-thread balanced frame stacks, hook starts with no ordinary parent thread,
`shutdown_begin` before attributed hook starts, one final `process_exit` with
signed status encoding, and no event after the terminal record. Open frame
prefixes are accepted for parked exit/halt callers and live ordinary peers.
Trace order is verified against the subprocess PID. Every isolated executable
digest is hashed before and after the run to prove the process did not modify
its own executable; the digest check runs for all modes of every case.

Representative generated bodies and the shared inspector output are retained
under [`generated/`](generated/): the mixed `Startup.main` control flow and
`Startup$Box.bump` field mutation, and the forced-AOT `Hooks$HookA.run`
monitor/exception-transfer lowering and `Hooks$HookB.run` control flow from the
hook path. These use the same `tests/inspect_generated.go` checks as Gates 1–4;
the mixed report also proves `Startup.step` is the interpreted crossing target
while caller and callback are AOT.

Telemetry artifacts are under [`metrics/`](metrics/). Example
`telemetry-mixed-plain.json`: `classesLoaded=74`, `crossEngineCalls=3`,
`reflectiveInvocations=2`, `reflectionLookupMisses=1`, `javaExceptions=2`,
`monitorContentions=1`, `liveJavaThreads=0`, `peakJavaThreads=2`; the rejected
admission leaves `classesLoaded` and `methodsByEngine` equal to the skip run.
Representative build reports are under [`reports/`](reports/) and fixture/class
digests under [`fixtures/`](fixtures/).

Deterministic Go tests in `rt/gate5_runtime_test.go` cover the profile boundary,
image flag consistency, bootstrap intrinsics, first-zero natural shutdown,
first-zero racing a new non-daemon start, registration rules, registration vs.
snapshot, manual vs. VM-owned hook start, simultaneous exits, halt cutting short
graceful shutdown, exit-from-hook stall, parked-caller frame/monitor retention,
metrics counters and trace sealing. They run under `-race`.

## 4. Application measurements (measured for `7469e36`, not smoke)

Runner: `tests/gate5_bench.py`. Evidence:
[`bench/SUMMARY.txt`](bench/SUMMARY.txt), [`bench/SUMMARY.json`](bench/SUMMARY.json),
[`bench/raw/`](bench/raw/) and [`bench/logs/`](bench/logs/).

These are the original full measurements from revision `7469e36`, before the
`651f706` review fixes. They are not a measurement of the current revision.

Protocol: one `javac --release 21 -g` compilation of the checked-in sources,
21 classes, class manifest SHA-256
`01828da24d10ddd1ab23fb7c0647e31491e94f5820b2eefef7f25720d12ab70d`;
identical class directory, operation counts, seed, worker counts and expected
checksums for OpenJDK 21 and Catty `interp`/`mixed`/`aot`. Every configuration
ran 7 fresh processes for cold, steady and memory, with 10 warmup and 30
post-warmup batches per steady process (210 measured batches per kernel per
configuration). Configurations were interleaved in rotating order. Timing runs
had tracing, metrics, `GORACE` and GC logging disabled.

Environment: macOS 27.0.1 arm64, Apple M6, 12 logical CPUs, 32 GiB RAM;
Go 1.27.1 darwin/arm64; Temurin OpenJDK 21.0.13.1+2; Python 3.9.6;
`GOMAXPROCS`/`GOGC`/`GOMEMLIMIT` default; JVM `-Xms256m -Xmx256m
-XX:+UseSerialGC`. Full metadata is in `bench/raw/environment.json` and
`SUMMARY.json`.

Build and size (median of 3 builds; Go build cache warm after the first
builder run, as documented):

| Configuration | Builder median | Executable bytes | Executable SHA-256 (prefix) | Image JSON bytes |
| --- | ---: | ---: | --- | ---: |
| interp | 0.166 s | 5,801,330 | `91735bf7e1220956` | 140,692 |
| mixed | 0.167 s | 5,904,418 | `2f667243dd6d32aa` | 58,195 |
| aot | 0.168 s | 5,904,466 | `b23f436941cc37cf` | 57,789 |

javac: 0.244 s. The mixed build forces `gate5bench.Hot.mix(II)I` to
interpretation while its caller `Workload.process` and callback remain AOT;
the report records the placement per method.

Cold launch to first verified result (median of 7): java 15.14 ms, interp
5.51 ms, mixed 3.85 ms, aot 3.57 ms.

Steady aggregate throughput and normalized batch latency p50/p99 (ns/op,
nearest-rank quantile over pooled batches; batch size documented per kernel):

                         aggregate ops/s                  p50 / p99 ns per op
kernel            java        interp       mixed      aot    | java        interp       mixed      aot
app               16,084,979   1,181,567   1,900,447  2,395,582 | 61.2/69.1  872.2/970.5  491.7/689.2  403.3/562.7
calls-static     552,940,935   4,456,105  21,689,460 21,753,575 |  1.7/2.4   223.6/233.4   46.9/47.8    46.4/47.7
calls-virtual    421,834,897   3,247,809  10,007,962  9,970,544 |  2.3/2.7   307.5/318.6   99.8/106.1  100.2/107.2
calls-mixed      266,931,480   2,604,382   4,889,921 10,293,198 |  1.9/11.4  378.6/418.1  213.4/226.4   95.6/126.8
reflect-lookup     6,635,390   1,730,098   2,523,652  2,646,054 |147.1/231.7 541.5/1022.4 370.7/849.9  365.3/797.8
reflect-invoke     8,890,073   1,468,941   2,989,296  2,976,722 |121.1/130.4 686.6/716.8  325.2/400.2  327.0/404.5
memory-alloc      49,906,304     952,369   1,345,979  1,340,613 |  8.8/78.3 1049.1/1102.4 743.7/782.3  748.4/777.6
shared-1         682,112,315   4,063,647   8,106,547  8,126,160 |  1.4/1.9   246.7/257.1  123.8/127.4  123.6/126.3
shared-4          93,317,127   1,978,897   2,058,016  2,181,809 | 10.5/14.4  505.5/519.2  485.9/502.4  458.8/469.8
monitor           34,112,287   1,198,768   1,424,366  1,430,998 | 25.6/133.7 836.6/889.6  701.5/747.1  699.4/747.6

Peak RSS (median of 7 dedicated memory-workload processes, `wait4`
`ru_maxrss`, macOS bytes): java 43,335,680 B (41.3 MiB), interp 30,703,616 B
(29.3 MiB), mixed 26,460,160 B (25.2 MiB), aot 26,411,008 B (25.2 MiB).

Retained graph (12,000 items, collection confirmed in all 7 samples per
configuration): OpenJDK post-GC Java heap median 1,048,576 B (from
`-Xlog:gc`, MiB granularity); Catty post-`Runtime.getRuntime().gc()`
`goRuntime.postGCHeapBytes` median: interp 6,985,032 B, mixed 6,696,640 B, aot
6,691,848 B. **These scopes differ**: Catty reports process-wide live
Go-managed heap including runtime structures; OpenJDK reports Java heap after a
full collection.

Observations: interpreted execution is the slowest and AOT the fastest for
almost all kernels; mixed recovers much of the AOT cost where the hot path is
AOT while the forced interpreted method remains visible in `calls-mixed`; the
shared-memory mutex cost shows as the gap between `shared-1` and `shared-4`
within each configuration; reflection lookup outside the cached-invoke window
is far more expensive than cached invocation. No performance targets are set
by Gate 5; these are baseline observations from one macOS arm64 host.

## 5. Intentional profile policies and unresolved limits

- **Hook registry semantics follow the OpenJDK 21 observations** used as
  acceptance cases: a terminated peer may be registered and is skipped at
  shutdown; null add/remove throw `NullPointerException` while running but
  `IllegalStateException` once the registry is frozen (validated before the
  argument); duplicate identities and live peers throw
  `IllegalArgumentException`; identity, not `equals`, keys the registry.
- **Exit from a hook intentionally stalls graceful shutdown** (matching the
  tested JDK behavior). A coordinated `halt` is the only way that case ends;
  the suite asserts the watchdog timeout for the unreleased variant.
- **POSIX status encoding is the tested host policy.** Trace events keep the
  full signed status; `os.Exit` encodes the low eight bits on POSIX
  (`-3 → 253`, `300 → 44`). Windows status representation and signal/logoff
  integration are not implemented or tested.
- **Out of scope by Gate 5 specification:** OS signal/logoff hooks, security
  managers, custom uncaught-exception handlers, hooks on virtual threads,
  thread subclasses overriding `start` as hooks, and other `Runtime` APIs.
  Gate 4 platform-thread restrictions still apply.
- **Metrics heap scopes differ** between Catty (`postGCHeapBytes`, Go-managed
  process heap) and OpenJDK (`-Xlog:gc` Java heap); retained-heap values are not
  directly comparable. Tracing/metrics change no externally visible output in
  the tested fixtures.
- **Benchmark scope/limits:** warm OS filesystem cache for cold samples (no root
  cache purge); OpenJDK JIT versus Catty interpreter/AOT is a runtimes-as-
  shipped comparison, not equivalent machine code; thread start/join is inside
  timed shared/monitor batches; normalized ns/op is batch time divided by the
  documented operation count, not per-operation latency. Linux RSS units would
  be KiB via `ru_maxrss`; only macOS arm64 was measured in this evidence.
- **Environment:** only macOS 27.0.1 arm64 with Go 1.27.1 and Temurin
  OpenJDK 21.0.13.1+2 was executed. The repository's system default `java` is
  25.0.4.1; every suite was pinned to the JDK 21 home above. Python is 3.9.6.
- **Smoke versus measured:** the benchmark smoke mode (reduced matrix) was used
  only to validate the runner; its numbers are not reported as evidence. The
  table in §4 is the full 7-process/30-batch measurement from clean revision
  `7469e36`. No benchmark result was fabricated or inferred.


## 6. Review-fix round (`f649a35` → `651f706`)

Revision tested: `651f706` on `codex/gate-5`. The starting revision was
`f649a35`. The observed starting-revision failures are retained verbatim in
[`acceptance/logs/review-reproductions.txt`](acceptance/logs/review-reproductions.txt);
focused fixed-test output is in
[`acceptance/logs/review-regressions.txt`](acceptance/logs/review-regressions.txt).

| Finding | Starting-revision failure | Fix and controlling regression |
| --- | --- | --- |
| 1. Untimed hook wait kills natural shutdown | Catty exited `2` with `fatal error: all goroutines are asleep - deadlock!` after printing `hook-wait`; OpenJDK stayed alive | Start one shutdown keeper with a long timer when the coordinator begins, so a timer remains pending while hooks wait. `natural --case` differential variant `wait` asserts the marker and watchdog timeout; `wait-halt` asserts prompt `halt(31)` termination in all three modes |
| 2. Halt during coordinator startup loses `shutdown_begin` | Trace contained only `process_exit` (`0 begins`) when halt was committed after freeze but before the coordinator emitted the event | Freeze stores a pending begin; terminal publication emits it under a dedicated ordering mutex before publishing the terminal result. `TestGate5HaltRaceKeepsShutdownBegin` pauses the coordinator, commits halt, and requires exactly one `shutdown_begin` (exit status `5`) before final `process_exit` (halt `19`); a running-state halt still emits none |
| 3. `RunApp(nil, nil)` panics | Nil-pointer panic at `rt/launcher.go:118` before returning a status | Reject a nil image before creating the runtime or reading metadata; return `2` and print `catty: nil image`. `TestGate5RunAppNilImage` captures the status and diagnostic |
| 4. Hook daemon flag stale in trace | Trace `thread_start` said `daemon:false` after a successful `setDaemon(true)` between freeze and start, while accounting used the new flag | Capture the daemon flag once at the new-to-started transition and use that value for both trace and live non-daemon accounting. `TestGate5HookDaemonCapturedAtStart` pauses before hook start, performs the successful `setDaemon(true)`, then requires `daemon:true` and `LiveNonDaemon()==0` while the hook runs |
| 5. Contention missed while still blocked | Terminal metrics reported `monitorContentions=0` for a contender blocked at halt | Count the acquisition when it first blocks, once per acquisition. `TestGate5MonitorContentionCountedAtFirstBlock` asserts `1` with the contender still blocked at halt, and `2` for two contenders released through wakeups/retries |

New retained artifacts:

- [`acceptance/logs/review-natural/`](acceptance/logs/review-natural/) — OpenJDK
  and Catty logs for `wait` (watchdog timeout after `main\nhook-wait`, status
  `-9`) and `wait-halt` (status `31`, `main\nhook-wait\nhalter`).
- [`acceptance/traces/launch-natural-*-wait*.jsonl`](acceptance/traces/) — open
  trace prefixes for the waiting hook and terminal halt traces.

Validation commands and outcomes for `651f706`:

```sh
go vet ./...                                                       # exit 0
go test -race -count=1 ./...                                       # ok catty/rt, all packages
go test -race -count=1 -v -run 'TestGate5RunAppNilImage|TestGate5HaltRaceKeepsShutdownBegin|TestGate5HookDaemonCapturedAtStart|TestGate5MonitorContentionCountedAtFirstBlock' ./rt/
                                                                    # PASS (4/4)
python3 -B tests/gate5.py --java-home <jdk21> --case natural --race --keep
                                                                    # natural passed
python3 -B tests/gate5.py --java-home <jdk21> --race               # Gate 5 differential tests passed
python3 -B tests/gate4.py --java-home <jdk21> --race               # Gate 4 differential tests passed
python3 -B tests/gate3.py --java-home <jdk21>                      # Gate 3 passed
python3 -B tests/gate2.py --java-home <jdk21>                      # Gate 2 passed (Gate 1 regression)
python3 -B tests/gate5_bench.py --java-home <jdk21> --smoke        # Gate 5 application benchmark (smoke): OK
```

No required suite failed. The full benchmark table in §4 remains the
`7469e36` measurement; the post-fix benchmark check is smoke-only by design.
Known behavior and environment limits in §5 are unchanged by this round.


## 7. Halt-to-hook-start race (`a65c445` → `af643e7`)

Revision tested: `af643e7` on `codex/gate-5`; starting revision `a65c445`.

`activateShutdown` previously checked the terminal state, released `thrMu`, and
then reacquired it to publish new-to-started hook transitions. A `halt` that
committed in that window let the coordinator mark snapshot hooks started after
the process result was already terminal: the peer left `threadNew`, the body
could run, and live non-daemon accounting was incremented for a terminated
application.

The fix moves the terminal-state check into the same `thrMu` critical section
that publishes the new-to-started transitions. `shutdown_begin` ordering is
unchanged: it is recorded before the publication window and terminal
publication emits the pending begin record before publishing the result.

Reproduction on the starting behaviour (test and coordinator hooks present,
activation fix absent):

```
$ go test -count=1 -run '^TestGate5HaltBeforeHookStartLeavesPeerNew$' ./rt/
--- FAIL: TestGate5HaltBeforeHookStartLeavesPeerNew (0.00s)
    gate5_runtime_test.go:1146: peer state after halt = 1, want new
```

`TestGate5HaltBeforeHookStartLeavesPeerNew` freezes the registry with a
registered, unstarted hook, pauses the coordinator in `BeforeHookStartHook`,
commits `halt(31)`, releases the coordinator, and waits for its completion hook.
It asserts:

- the peer remains `threadNew` and its body never runs;
- `LiveNonDaemon()` is zero before and after the halt;
- the trace has exactly one `shutdown_begin` (the frozen exit status `5`)
  before one final `process_exit` with halt status `31`;
- no `shutdownHook` `thread_start` and no `app/Hook` `enter` event exist.

Validation for `af643e7`:

```sh
go test -race -count=100 -run '^TestGate5HaltBeforeHookStartLeavesPeerNew$' ./rt/
                                                                    # ok
go test -race -count=20 -run 'TestGate5HaltRaceKeepsShutdownBegin|TestGate5HaltCutsShortGracefulShutdown|TestGate5HookDaemonCapturedAtStart|TestGate5SimultaneousExitRequests|TestGate5ExitFromHookStallsUntilHalt|TestGate5TraceSealsAtTerminal|TestGate5ShutdownHookRegistrationRules' ./rt/
                                                                    # ok
go test -race -count=1 ./...                                       # ok catty/rt, all packages
go vet ./...                                                       # exit 0
python3 -B tests/gate5.py --java-home <jdk21> --race               # Gate 5 differential tests passed
```

Retained under [`acceptance/logs/review-halt-publication/`](acceptance/logs/review-halt-publication/):

- `reproduction-a65c445.txt` — starting-revision failure transcript.
- `fixed-stress-new-regression.txt`, `fixed-stress-existing.txt` — 100× and
  20× race stress results.
- `fixed-go-test-race.txt`, `fixed-vet.txt`, `fixed-gate5-race.txt` — required
  single-pass suite, vet and full Gate 5 race outputs.

### Pre-existing Gate 4 timer-arming flake (recorded separately)

A stress run exposed a Gate 4 wait-set test flake unrelated to this fix;
`go test -race -count=5 ./...` failed with
`TestGate4WaitSetWinnerTransitions`/`TestGate4WaitSetTimeoutAndInterrupt`
timing out on a timeout transition. The same command failed in 2 of 3 attempts
on a clean `a65c445` worktree, so it pre-dates the halt-race change. The test
helper closes `ready` before `objectWait` arms the manual-clock timer, so an
advanced clock can miss the timer. Details and logs are in
[`gate4-timer-flake-note.txt`](acceptance/logs/review-halt-publication/gate4-timer-flake-note.txt),
`gate4-timer-flake-a65c445-count5.txt` and
`gate4-timer-flake-fixed-tree-count5.txt`. It is not counted as a remaining
Gate 5 failure: the required single-pass `go test -race -count=1 ./...` passes,
and the serial targeted Gate 4 timeout tests pass.

The full benchmark table in §4 remains the original `7469e36` measurement and
was not rerun for this lifecycle-only race fix.
