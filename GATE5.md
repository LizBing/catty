# Gate 5 implementation specification

Gate 5 delivers application shutdown and process diagnostics for one standalone Catty executable, plus measurements of the same Java workload across the interpreter, AOT, mixed execution, and OpenJDK 21. The delivery follows the startup, evaluation, and acceptance contracts in the architecture. [Architecture §§8–11](ARCHITECTURE.md#8-application-startup-fallback-and-exit) · [Acceptance and measurement plan](GATE5-TESTING.md) · [Implementation-agent prompt](GATE5-IMPL-PROMPT.md)

## Profile and bootstrap surface

Add `catty build --enable-lifecycle` to select `catty-jvm21-gate5` and image ABI `5`. In descending precedence, `--enable-lifecycle`, `--enable-threads`, and `--enable-plugins` select Gates 5, 4, and 3. The default remains Gate 2. Gate 5 inherits Gate 4's opcodes, threading, and Gate 3 plugin limits. Serialize `lifecycle: true` in the Gate 5 image and build report. Treat omission as false and accept omitted flags on older images. Publish canonical capabilities including `application-lifecycle`, `shutdown-hooks`, and `application-diagnostics`. Reject an inconsistent profile, ABI, or feature flag before Java runs. Use the selected profile for the bootstrap inventories, ordinary build admission, plugin verification and admission, and image installation. Keep the descriptors and calling convention of `Entry`, `Outcome`, and `Value`.

Extend the Gate 4 bootstrap inventory without modifying its definitions. Publish the new classes and members through the existing inventory, linkage, access, and reflection paths. The additional surface is as follows. [Java 21 System API](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/System.html) · [Java 21 Runtime API](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/Runtime.html)

| Owner | Member and descriptor | Metadata and behavior |
| --- | --- | --- |
| `java/lang/System` | `exit(I)V`, `nanoTime()J` | Public static. Exit uses the shared shutdown request. `nanoTime` uses an arbitrary process-local monotonic origin and is for elapsed-time differences. |
| `java/lang/Runtime` | `getRuntime()Ljava/lang/Runtime;` | Public static. Return the same non-null bootstrap instance throughout each runtime. The instance has private construction. |
| `java/lang/Runtime` | `exit(I)V`, `halt(I)V` | Public instance. Use the process termination paths below. |
| `java/lang/Runtime` | `addShutdownHook(Ljava/lang/Thread;)V`, `removeShutdownHook(Ljava/lang/Thread;)Z` | Public instance. Use the receiver-independent registry for that runtime. |
| `java/lang/Runtime` | `gc()V` | Public instance. Request a Go GC cycle. When metrics are enabled, record the process-wide Go heap observation immediately after the cycle for the retained-heap workload. |
| `java/lang/IllegalStateException` | Existing bootstrap throwable-constructor subset | Public subclass of `RuntimeException` for mutations to the hook registry during shutdown. |

The profile covers application-initiated and natural shutdown. OS signal and logoff integration, security managers, custom uncaught-exception handlers, hooks on virtual threads, thread subclasses that override `start` as hooks, and other `Runtime` APIs remain outside Gate 5. Gate 4 platform thread restrictions continue to apply. The evaluation records measurements without setting performance targets. [Gate 4 bootstrap inventory](GATE4.md#bootstrap-inventory)

## Lifecycle and shutdown

The generated launcher installs the embedded image and passes `os.Args[1:]` into one Java `String[]`. It creates the canonical main peer and uses `Thread.Call` for the resolved entry and nested synchronous calls. Maintain one process and logical thread identity across synchronous engine crossings. Main class initialization, arguments, class lookup and initialization ordering, and standalone packaging follow the earlier gates.

Model a single running state, a shutdown state with a frozen hook snapshot, and a terminal state with a selected Java or host status. Serialize lifecycle state, hook registry mutations, and live thread publication with the thread lifecycle guard or an equivalent lock protocol with the same linearization points. Keep Java calls, blocking waits, diagnostic and trace writes, and process finalization outside the shared state lock. Allow the launcher to receive a terminal result even while Java main is running or blocked. Keep OS process exit at the generated launcher boundary. [JVMS termination](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-5.html#jvms-5.7) · [Java 21 Runtime shutdown contract](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/Runtime.html#shutdown)

| Trigger | Transition and status |
| --- | --- |
| Natural termination | Count main and all live non-daemon peers. Publish main's normal `0` or uncaught `1` status and any uncaught diagnostic before it retires. When the count first reaches zero in the running state, freeze the registry and commit natural shutdown with that status. A non-daemon descendant published before that transition extends the drain. A daemon alone never delays the transition. |
| `System.exit(status)` or `Runtime.exit(status)` | The first request to acquire shutdown ownership snapshots the registry and selects the exact signed Java `int`. Explicit shutdown can begin while main and other non-daemons remain live. Later exit calls park until process termination, including an exit from a hook. Keep the selected status if a hook throws or a later exit call requests a different status. A hook calling exit can indefinitely delay graceful shutdown. |
| `Runtime.halt(status)` | A halt that wins terminal publication skips the hook snapshot or cuts short an active graceful shutdown. Select its exact status and signal the launcher without waiting for Java execution. Further requests cannot change a committed terminal result. |
| Setup failure | Decode, image installation, and launcher diagnostics setup failures use status `2`, before Java entry or hook execution. A Go panic remains a host fault. |

Keep the full signed value through terminal diagnostics and pass it to `os.Exit(int(status))`. The host defines process-status representation. On POSIX the parent observes the low eight bits of a normal exit. Successful exit and halt calls never return or throw into Java, so surrounding catch and `finally` bodies do not execute on that calling thread. Park a requesting logical thread without unwinding its frames or releasing its Java monitors. Keep that thread observable as live until process termination. The shutdown coordinator progresses independently of an exiting or halted main, including when main exits from class initialization. Use host parking that allows a deliberately stuck Java shutdown to remain live even when no Go work is runnable. [Java 21 `Runtime.exit` and `halt`](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/Runtime.html#exit(int))

### Hook registration and execution

- Key registration by Java object identity. While running, adding a null hook throws `NullPointerException`. Adding a duplicate or live peer throws `IllegalArgumentException`. Follow the OpenJDK 21 behavior for a previously terminated peer. Registration may succeed, and the peer is skipped at shutdown because it has already run. Removing a registered peer returns `true`, removing an absent peer returns `false`, and removing null throws `NullPointerException`. Once shutdown claims the registry, both operations throw `IllegalStateException` before validating the hook argument. [OpenJDK 21 hook registry](https://github.com/openjdk/jdk21u/blob/master/src/java.base/share/classes/java/lang/ApplicationShutdownHooks.java)
- Freeze and detach the registry once. Start all eligible hooks before waiting for any. Reuse the Gate 4 new-to-started transition, Java peer, logical thread, virtual `run` selection, and ordinary termination path. Hooks are concurrent with each other and existing threads. They can be daemon peers or non-daemon peers, and all started hooks count toward completion of the snapshot. Use the same peer synchronization for join visibility. A hook started by Java after registration must never run twice. If its manual start wins, count that execution only while it is still live and registered in the snapshot. Hook start order is unspecified. [OpenJDK 21 hook startup and joins](https://github.com/openjdk/jdk21u/blob/master/src/java.base/share/classes/java/lang/ApplicationShutdownHooks.java#L88-L112)
- Keep ordinary `Thread.start`, interrupt, join, monitors, class initialization, and plugin linkage operational during graceful shutdown. Finish when the hook snapshot has retired, even if a regular thread started by a hook remains live. Diagnose an uncaught hook exception through the standard named-thread path before publishing its completion. Hook exceptions leave the selected status intact. Release registry and completion references as peers cease to participate. [Java 21 shutdown sequence](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/Runtime.html#shutdown)

## Process diagnostics and telemetry

Extend `CATTY_TRACE` JSON lines with additive process events. Keep existing event shapes, logical thread lifetimes, loader IDs, engines, throwable IDs, and bytecode PCs. A process event always has `pid`. Include `requesterThread` when a Java thread requested a transition. A VM-started hook has its own `thread_start` with a `shutdownHook` marker and no ordinary `parentThread`. Use the existing `thread_end` for hooks that actually retire. The following events are required for an installed application.

| Event | Required fields and order |
| --- | --- |
| `app_start` | `profile`, `imageABI`, `argc`, and `mainThread`. Emit before the first Java entry. Record argument count without values. |
| `shutdown_begin` | `reason` (`natural` or `exit`), selected `status`, `hookCount` for the frozen registry, and `requesterThread` for explicit exit. Emit once before any hook start attributed to the snapshot and before `process_exit`. A halt in the running state skips this event. |
| `process_exit` | Terminal `reason` (`natural`, `exit`, or `halt`), signed `status`, and `requesterThread` for a Java-requested terminal action. For graceful shutdown, emit after the participating hooks' last events. Seal the trace at this event so any remaining live thread has a well-formed trace prefix. |

Use the existing named-thread Java diagnostic on stderr. For an internal panic during Java execution, obtain the logical thread's Java frame snapshot while its frames are still available, then preserve the original panic and Go stack. Keep Go faults out of Java exception handling and ordinary shutdown completion. The trace sink must remain race-safe when a daemon or unjoined ordinary Java thread outlives the trace. [Architecture diagnostics](ARCHITECTURE.md#9-diagnostics-and-evaluation)

Add opt-in `CATTY_METRICS=<path>` for Gate 5. Publish one JSON object with schema version, PID, profile, terminal reason, signed status, and the following counters. Use snapshot-safe updates independent of tracing. With metrics disabled there is no counter work beyond the feature check. Keep telemetry off application stdout. An output-open failure is a setup error. Report an output-write failure on stderr and preserve the selected terminal status.

| Counter | Counting point |
| --- | --- |
| `classesLoaded` | Successfully published application, bootstrap, plugin, primitive, and array class definitions, once per runtime identity. Rejected admission leaves counts intact. |
| `methodsByEngine` | Published method definitions by their assigned engine, including bootstrap intrinsics and late interpreted members. |
| `crossEngineCalls` | Adjacent synchronous Java execution frames crossing between AOT and interpretation. Treat intervening intrinsic frames as adapters for this counter and break the chain at a new logical thread. |
| `reflectiveInvocations`, `reflectionLookupMisses` | Reflective calls that reach their selected target and lookup failures for absent members. |
| `javaExceptions` | Java throwable escapes from `Thread.Call`, counted per escaping frame. Catch, wrap, and uncaught events retain their existing identity semantics. |
| `monitorContentions` | Monitor acquisitions that must wait for another owner, once per contending acquisition. |
| `liveJavaThreads`, `peakJavaThreads` | Active main, workers, and hooks at the snapshot and the high-water mark, including daemons and excluding unstarted peers. |
| `goRuntime` | Process-wide Go heap allocation bytes, cumulative allocation bytes, completed GC count, and an optional post-`Runtime.gc()` heap observation. Label the latter as live Go-managed heap for the process, including runtime structures. |

## Implementation plan

| Step | Implementation | Exit check |
| --- | --- | --- |
| 1 | Wire the CLI and image/report schema in `internal/build`, add the profile and bootstrap inventory in `rt`, and audit profile-scoped checks in `internal/classfile`, `internal/verify`, and `internal/admission`. | Gate 5 accepts and advertises its published surface. Gates 2–4 keep their existing profile IDs, ABIs, and admission behavior. |
| 2 | Add the shutdown coordinator and registration operations. Integrate start and termination publication in `rt/thread_java.go` and process handoff in `rt/launcher.go` and generated `main.go`. | Runtime transition tests control the first zero count, competing exit calls, hook snapshot races, and a halt during hook execution. |
| 3 | Add process trace events, fault diagnostics, and optional metrics at the shared runtime paths. | Subprocess evidence shows one terminal result, synchronized diagnostics, and matching process and thread IDs. Counter tests cover instrumentation on and off. |
| 4 | Add the OpenJDK differential cases and profile checks. Exercise interpreted, mixed, forced-AOT, late plugin, standalone, and generated race-enabled paths. | Gate 5 acceptance and the Gate 1–4 regressions in the test plan pass. |
| 5 | Add the measurement workload and runner, then collect four-configuration results with toolchain and environment metadata. | `out/gate5/EVIDENCE.md` links commands, raw samples, summaries, and trace evidence. [Acceptance and measurement plan](GATE5-TESTING.md) |
