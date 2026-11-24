# Gate 5 acceptance and measurement plan

Compile each case once with OpenJDK 21 and `javac --release 21 -g`. Verify those class files with `java -Xverify:all`, then use the same files in Catty `interp`, mixed `auto` with an explicit policy, and forced `aot` images. Assert the case's expected outcome on the oracle as well as on Catty. Capture stdout, stderr, exit status, and timeout for every run. [Gate 5 implementation specification](GATE5.md)

## Lifecycle acceptance

Add fixtures under `testdata/gate5/` and a `tests/gate5.py` runner that reuses the existing acceptance tooling. Give each case a subprocess watchdog. Use the Gate 4 monitor and thread primitives to synchronize critical orderings. For concurrent hook output, assert required lines, multiplicity, and happens-before constraints. Keep order-independent output out of exact stdout comparisons. Use pre-exit markers to establish progress in cases that deliberately block.

| Case | Required observation |
| --- | --- |
| `startup` | Run with no arguments, then with an empty string, whitespace, a leading `-`, and Unicode. Assert argument count, order, and Java string contents. Assert main-class initialization before the first main body output. In the mixed run, AOT main calls an interpreted method that calls an AOT callback with the same live object. Observe the result and one PID and logical thread ID through the full nested sequence. The all-AOT and all-interpreted reports and traces prove their assigned engines. |
| `natural` | Main returns or throws while a non-daemon worker is live. The worker spawns a non-daemon descendant that completes before the first hook body. A blocked daemon stays live at exit. Exercise both main outcomes with named, checked stderr and statuses `0` and `1`. Exercise a main joiner that sees main's termination before doing its last work. |
| `registration` | Verify the runtime singleton, successful add and remove, repeated remove, null rejection, duplicate identity rejection, live peer rejection, and registration of a retired peer according to the reference behavior. Use distinct hooks whose `equals` returns true to verify identity. Check hook probes for add and remove after freeze, including null arguments. Ensure removed and already retired peers never execute at shutdown. |
| `hooks` | Use two hooks that must reach a barrier before either can finish. Include a daemon hook, a hook with a virtual `run` override, and a named hook that throws. Assert hook identity through `Thread.currentThread`, one uncaught diagnostic with Java frames, and the original exit status. A hook starts an ordinary non-daemon worker that stays blocked. The process still finishes when its hook snapshot retires. |
| `exit` | Use `System.exit` from a worker while main is blocked and `Runtime.exit` via reflection from mixed main. Include exit during main-class initialization. Test status `0`, a positive status, and a negative or wide value to verify host encoding. Hooks run exactly once. Catch, `finally`, after-exit, and post-startup body markers stay absent where applicable. One exit caller holds a monitor and a hook observes that the caller remains live, before a separate hook uses halt to resolve the fixture. |
| `exit-contention` | Hold the shutdown snapshot open with a hook and have a second worker attempt a different explicit exit. Complete the hook and check that the first request's status wins. Separately have a hook call exit, assert that its pre-exit marker appears and that no code after exit executes, then observe the expected watchdog timeout or release it with a coordinated halt in another variant. |
| `halt` | Halt before shutdown with a registered hook and during shutdown with an active hook. Assert selected status, absence of skipped hook or cleanup markers, and immediate launcher progress despite blocked Java threads. |
| `late-hook` | Build the app before a plugin JAR exists, then admit the JAR and register a hook whose `run` or `Runnable` target is supplied by the plugin. Use the Gate 3 OpenJDK adapter in the reference run. Prove the target executes interpreted with the published plugin loader ID, the peer has a separate logical thread, and the same process terminates with the chosen status. Forced AOT applies to image methods. |
| `setup-failure` | Rebuild generated `main.go` with an invalid image pair or feature flag and exercise an unusable diagnostics output path. The executable reports the failing setup stage, exits `2`, and produces no Java entry or hook output. Use a Go subprocess fixture to inject a host panic while Java frames are active and require both Java frame information and the original Go panic stack. |

Use `tests/inspect_generated.go` for representative generated bodies and check report-selected engines, including the mixed crossing and hook path. An isolated run starts the copied executable with class directories, JDK, and adapter unavailable. Supply only explicit application inputs such as the late JAR. Check the executable digest before and after.

### Profile and concurrency checks

- Cover CLI flag precedence, Gate 5's exact profile/ABI pair, inherited opcodes and plugin limits, and the bootstrap descriptors in `GATE5.md`. Build separate images that refer to the new APIs and assert earlier profile refusal without artifacts. Test the same admission boundary with a reference in an otherwise valid late JAR. Confirm install fails for ABI/profile/flag mismatches and that older images decode with the lifecycle flag omitted.
- Add deterministic Go tests for main and last-worker retirement racing a new non-daemon start, add/remove versus snapshot, ordinary start versus snapshot-owned start, simultaneous exit requests, halt versus terminal publication, and hook completion versus registration cleanup. Control these with channels or existing runtime test hooks. Assert one terminal publication and verify that a worker published before the natural transition participates in the drain.
- Validate traces against the subprocess PID, per-thread stacks, loader identity, event ordering and exit encoding. Require one `app_start` and one final `process_exit` on completed runs, one `shutdown_begin` on graceful runs, correct requester and hook thread links, and no events after the terminal record. Accept open frame prefixes for parked exit callers and live peers when the process terminates. Check Java diagnostics and trace counts under concurrent output.
- Verify telemetry with a small deterministic fixture that exercises an AOT/interpreter crossing, reflection success and lookup miss, matched exceptions, monitor contention, one late admission and one rejection, thread publication, and the explicit GC observation. Assert counter deltas where the base includes bootstrap definitions. Compare externally visible behavior with tracing and metrics disabled and enabled.
- Run all lifecycle cases uninstrumented, then repeat them with race-enabled generated executables where supported. Inspect `go version -m` to establish that those executables contain race instrumentation, and use untraced runs as well. Treat any detected Go data race as a failure. Run the Gate 4 race-enabled suite, Gate 3, and Gate 2 with its Gate 1 regression under their original profiles.

## Application measurement

Check in workload sources under `testdata/gate5/benchmark/` and a reproducible `tests/gate5_bench.py` runner. Build the workload classes once, digest them, and use identical class files, operation counts, worker counts, seed, inputs, and expected checksums in all four configurations. The mixed policy forces a designated hot method to interpretation while its caller and callback remain AOT. The whole workload runs in separate JVM and Catty processes with the same externally visible work and results. Include a correctness run that proves engine placement before taking uninstrumented timing and RSS samples.

### Workload and samples

| Workload | Required comparison |
| --- | --- |
| Application | Call, field, and allocation work over a deterministic object graph with an engine-crossing path and a checked result. A separate command mode retains a known graph through the explicit GC observation. |
| Calls | Static/direct, virtual dispatch with more than one receiver type, and the mixed-policy crossing, timed as distinct kernels. Use a result checksum and inspect the selected targets. |
| Reflection | Public member lookup by name and repeated cached `Method.invoke` measured separately. Keep lookup outside the cached-invoke window. |
| Memory and threads | Allocation with escaping objects, a field- and array-heavy shared-memory kernel, and a contended monitor kernel with a fixed number of platform threads and verified completed work. Run field-heavy work at one and multiple threads to show the cost of the shared memory mutex. |

- **Build and size.** Measure javac duration and Catty builder wall time per mode with the Go toolchain cache policy recorded. Report compilation and packaging boundaries, median build time, executable bytes and SHA-256, and embedded image JSON bytes before Go string escaping. The common class directory or JAR is the input artifact for OpenJDK. Record its size and report executable and image sizes only for the Catty configurations.
- **Cold result.** Measure from process launch to receipt of the first complete, verified result line with a monotonic clock. Validate final status and checksum as well. Use fresh processes for each sample and document filesystem-cache conditions. Keep build time out of launch timing.
- **Steady result.** Warm up each process for a fixed documented amount of work. Use `System.nanoTime()` around batches with a known operation count. Keep the clock read, checksum reporting, and output outside the measured operation loop. Report aggregate operations/s and distributions of batch latency normalized to ns/op, including p50 and p99 with a documented quantile rule. Record the actual batch sizes and avoid interpreting normalized batch times as individual operation latency.
- **Process memory.** Measure peak RSS for each dedicated sample process using a per-child OS high-water reading, with the OS source and units recorded. Run the retained-graph workload separately. Take Catty's post-`Runtime.gc()` Go heap observation with telemetry enabled. Take OpenJDK's post-collection Java heap occupancy at the same workload checkpoint using GC logging and a configuration that performs the explicit collection. Keep the root live through the measurement. Record retained item count, collection confirmation, and raw bytes or source units. Label the two heap metrics with their different runtime scopes. If the available platform cannot provide a required metric, record the failed method and blocker in the evidence.

The evidence run uses at least seven fresh processes per configuration for cold, steady, and memory runs and at least 30 post-warmup batches per steady process. Interleave the configurations and keep raw per-process and per-batch samples. Use the same CPU/OS, worker counts, and benchmark input for each run. Record execution order, repository revision, ordered input digests, Go and OpenJDK versions, target OS/architecture, CPU, Go build settings, `GOMAXPROCS`, `GOGC`, `GOMEMLIMIT`, JVM heap and GC flags, warmup, sample count, and quantile method. Use a smaller smoke run only to verify the runner. Keep trace, optional counters, GC probes, and race instrumentation out of measured latency and throughput runs. Report variation, workload-specific observations, and scope of measurements. [Architecture evaluation](ARCHITECTURE.md#9-diagnostics-and-evaluation)

## Commands and evidence

These are the required entry points after implementation. `--prepare-only` runs the JDK path without a Gate 5 build, so fixture work can begin on the planning baseline. The benchmark runner should expose a smoke option with the full matrix and reduced repeats.

```sh
go test -race -count=1 ./...
go vet ./...
python3 -B tests/gate5.py --java-home /path/to/jdk-21 --prepare-only
python3 -B tests/gate5.py --java-home /path/to/jdk-21 --race
python3 -B tests/gate5_bench.py --java-home /path/to/jdk-21 --smoke
python3 -B tests/gate5_bench.py --java-home /path/to/jdk-21
python3 -B tests/gate4.py --java-home /path/to/jdk-21 --race
python3 -B tests/gate3.py --java-home /path/to/jdk-21
python3 -B tests/gate2.py --java-home /path/to/jdk-21
```

Record actual commands, exit statuses, runtime and fixture digests, report inventories, representative trace and generated-code evidence, and benchmark summaries in `out/gate5/EVIDENCE.md`. Link raw benchmark data and the logs needed to rerun or audit each claim. Retain full results for any case that is blocked, disagrees with the oracle, or times out unexpectedly. Gate 5 completes when the acceptance and regression suites pass and all four configurations have usable measurements.
