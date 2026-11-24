# Gate 4 test plan

Gate 4 requires the checked-in Java cases to pass in OpenJDK 21 and Catty, plus deterministic Go tests for runtime transitions and admission state. Run the same class files through all-interpreted, mixed, and forced-AOT Catty images. Test memory-model invariants and permitted outcomes without depending on scheduling order. [Gate 4 implementation specification](GATE4.md)

## Executable suite

[tests/gate4.py](tests/gate4.py) compiles the fixtures in [testdata/gate4](testdata/gate4) with `javac --release 21 -g`, runs them with `java -Xverify:all`, and compares Catty exit status and stdout against the same fixture. The checked-in expectations validate the oracle too. Successful cases require empty stderr. Uncaught cases check the Java exception class, thread name, and application frame. The process watchdog treats a hang as a failure and keeps command logs with timeout details.

Each mixed policy forces the method named below to interpretation while main is AOT-compiled. The runner checks the profile, opcode inventory, capabilities, method selections, and engine use. Normal Catty binaries run with trace enabled and without a JDK on `PATH`. With `--race`, the suite also executes race-instrumented binaries in all three modes with tracing disabled. It checks Go binary metadata for `-race=true`, so a Go package test alone cannot satisfy the generated-code requirement. [Go race detector](https://go.dev/doc/articles/race_detector)

| Case | Mixed interpreter method | Required observations |
| --- | --- | --- |
| `lifecycle` | `Lifecycle$Worker.run` | Canonical current thread, stable positive thread ID, synchronous direct `run`, virtual child entry, base `Runnable` dispatch, one winning first start, restart rejection, join on new/live/terminated peers, timed join, interrupted join retaining unrelated locks, and daemon configuration. |
| `monitors` | `Monitors.nested` | Reentrant explicit and implicit acquisitions see the same receiver. Instance and static monitors release during unwind. Reflection wraps the exact original throwable. Contended entry preserves interrupt status. Both workers contribute to the exact protected counter. Null monitor entry throws at the bytecode site. |
| `wait-set` | `WaitSet$Trial$Waiter.run` | Wait and notification require ownership. Argument errors and pending interrupts follow Java rules. Timed wait, normal notification, and interrupted wait restore the full lock depth. Either interrupt/notify ordering permits one normal wake and leaves each participant with a legal interrupt outcome. Long sleep remains interruptible and retains unrelated monitors. |
| `initialization` | `Initialization$Ready.<clinit>` | Class literals leave initialization dormant. Recursive active use sees the default. Competing users observe the initialized value or the erroneous state. The owner gets a wrapper with the original cause and a later use gets `NoClassDefFoundError`. The class initializes once and a competing thread retains its interrupt. |
| `memory` | `Memory$Publisher.run` | A volatile flag publishes an ordinary field, a long array element, a reference array element, and constructor writes. Joining publishes worker results. Volatile store-buffering excludes `(0, 0)`. Racy reference reads permit absence and can see the default ordinary field on OpenJDK, while each observed object's final field must be correct. After join, the ordinary field is fully published. |
| `launcher` | `Launcher$Worker.run` | After main terminates, the worker prints and starts a non-daemon descendant. The process waits for the descendant's output after the worker terminates. On normal main return, it reports the worker exception and exits `0`. On uncaught main, it reports the main exception and exits `1`. A daemon remains waiting at process exit, and a thread constructed by that daemon inherits daemon status. |
| `plugin` | `app/Main$Registrar.run` | Two Java workers concurrently register identical late JAR bytes. Independent loaders define the same worker name, each with its own class mirror, static storage, and monitor. Each plugin subclass starts on a new thread and runs interpreted, then joins and exposes its state to its registrar. |

The plugin image builds before the plugin class directory and JAR are created. The compile-only stub supplies the API at javac time. The Gate 3 OpenJDK adapter supplies an independent loader during reference runs. Catty sees the app classes at build time and the JAR only at execution. The runner asserts absence of the late worker from the image method inventory, and correlates registration loader IDs, registrar thread IDs, and depth-zero worker runs.

The suite includes a class-file mutation that puts an integer on the operand stack at the original `monitorenter` PC. OpenJDK rejects it with `VerifyError`. Gate 4 must reject at that method and PC without publishing artifacts. The original explicit-monitor class and the preexisting synchronized-method fixture must still be rejected by both earlier profiles.

The runner checks that all trace records carry the launched PID, frames follow separate logical thread stacks, all non-daemon threads retire, and the waiting daemon may leave an open frame prefix. It also checks monitor reentry and bytecode locations, wait outcomes and restored depth, identity through unwind and catch, initializer ownership and outcome, and plugin loader and worker relationships.

## Deterministic Go tests to add with the runtime

Use in-memory test threads, Go barriers, and an injectable monotonic clock where timing affects behavior. Wait for an observable transition before releasing the next participant. Restrict wall-clock deadlines to watchdogs. Java fixture handshakes can show an attempt but cannot ensure a competing Go goroutine has reached its actual blocking point. These tests must force those transitions.

| Area | Required schedules and assertions |
| --- | --- |
| Monitor ownership | Reenter to depth three, queue a distinct thread, and prove the contender stays blocked until the final exit. Exit from a non-owner leaves the owner and depth intact. Check null and unowned exit outcomes in both engines using verified fixtures or direct helper tests. Keep monitored objects collectible after the last Java reference is dropped. |
| Method entry | Check static lock identity by defining class and loader, inherited static selection, per-receiver instance locks, abstract/missing method errors before lock acquisition, successful return and uncaught unwind, and a same-frame catch while the implicit monitor is held. Test implicit monitor acquisition for reentrant Java callbacks using the existing `Thread.Call` ABI. |
| Wait-set transitions | Prove enrollment and full release are atomic with notification. Force notify, notify-all, timeout, and interrupt individually. Force notify then interrupt before reacquisition, interrupt then notify with a second waiter, and timeout versus interrupt in both orders. Assert that exactly one reason wins, no waiter is stranded, unrelated monitors stay owned, late interrupts remain pending on normal return, and the original recursion depth is restored in every branch. |
| Timers and overloads | Check no-wait versus indefinite cases, negative and overflowing inputs, nanos boundaries, a stale timer firing after a later wait begins, delayed reacquisition after timeout, and a timed join with repeated unrelated wakeups. Confirm one monotonic join deadline. |
| Thread accounting | Race two starters and race `setDaemon` with start. Check peer identity, start and termination publication, multiple joiners, a pending interrupt on active join and sleep, and clean termination of an uncaught worker. Force an older non-daemon worker to start a new non-daemon child after main finishes and before the live count reaches zero. Verify that the launcher waits for both. |
| Class initialization | Hold the owner in its initializer until a second thread reaches the initialization wait. Check recursion, interrupt preservation, superclass/default-interface order, and publication of ordinary static writes. Force success, an original `Error`, non-`Error` wrapping, and parent failure. Confirm exact owner outcomes, terminal state, trace ordering, and a fresh `NoClassDefFoundError` for each subsequent use. Repeat with the same binary name under two plugin loaders. |
| Runtime metadata | Race first resolution, cached failure, compiled dispatch, reflection method and argument snapshots, array classes, mirrors, interned strings, and boxed-cache entries. Verify canonical identities and loader qualification. While resolution and successful admission run, fail a second registration after staging constraints and arrays. Confirm rollback preserves the other transaction and never rewinds object IDs. Assert a shared JAR byte-array snapshot remains detached and a concurrent quota race allows only the registration inside the limit. |
| Java memory access | Exercise all four ordinary and volatile static and instance field operations and every array backing kind through both engines. Race Java reads, stores, and reflection argument snapshots. Assert valid reference and wide primitive values, no torn writes, constructor publication, and monitor/start/join happens-before. Verify that tracing can close safely during daemon activity. |

## Build and regression completion

Add implementation-level coverage for all profile/ABI pairings, feature flag mismatches, missing or conflicting inventories, synchronized method flags in retained images, and rejected legacy bootstrap references. Inspect generated monitor code using the existing Go AST inspector pattern in `tests/inspect_generated.go`. Verify monitor helpers and original instruction PCs feed the shared exception-transfer path. Add equivalent verifier checks for an invalid `monitorexit` operand, synchronized class-file flag combinations, and a valid synchronized static method. Keep a strict forced-AOT rejection case for the existing unsupported image lowering.

Run the full Gate 3 suite against its original profile and the Gate 2 suite with its Gate 1 regression. Reuse their positive corpora with Gate 4 enabled once the build flag exists. Include the Gate 3 rollback, loader-identity, and late-override cases. In the standalone Gate 4 run, hide the fixture class directories and adapter after building and allow the executable only its explicit JAR input. Confirm the executable digest stays the same.

## Commands

From the repository root, select OpenJDK 21 and a Go host with race support for full acceptance.

```sh
python3 -B tests/gate4.py --java-home /path/to/jdk-21 --prepare-only
python3 -B tests/gate4.py --java-home /path/to/jdk-21 --case wait-set
python3 -B tests/gate4.py --java-home /path/to/jdk-21 --race --keep
go test -race -count=1 ./...
go vet ./...
python3 -B tests/gate3.py --java-home /path/to/jdk-21
python3 -B tests/gate2.py --java-home /path/to/jdk-21
```

The reference-only command works on the Gate 3 baseline. The Catty paths require the planned Gate 4 profile and runtime. Completion evidence belongs under `out/gate4/` and includes toolchain versions, command statuses, class and JAR digests, case results, build reports, representative generated code, and trace excerpts for thread lifecycle, nested monitors, wait outcomes, class initialization, and concurrent plugin admission. Include the untraced generated-binary race runs, Go race tests, and prior-gate regressions.
