# Catty — Gates 1–5

Catty builds one Java 21 application into one standalone Go executable. A Go
AOT emitter turns verified JVM bytecode into real Go functions, an interpreter
retains the methods outside the AOT lowering profile, and both engines share the
same object model, dispatch, class initialization, reflection and exception
semantics.

This repository implements the Gate 2 delivery in [GATE2.md](GATE2.md) against
[ARCHITECTURE.md](ARCHITECTURE.md); the Gate 1 contract remains in
[ACCEPTANCE.md](ACCEPTANCE.md) and [IMPLEMENTATION.md](IMPLEMENTATION.md).
`tests/gate2.py` is the Gate 2 differential suite and runs the Gate 1 suite as
its regression check.

Gate 3 adds runtime plugin admission. Build with `--enable-plugins` to select
the `catty-jvm21-gate3` profile and image ABI `3`; without the flag the
published Gate 2 profile `catty-jvm21-gate2` and ABI `2` stay unchanged. See the
[Gate 3 implementation specification](GATE3.md) and
[Gate 3 test plan](GATE3-TESTING.md).

Gate 4 adds Java platform threads, monitors, wait/notify/interrupt, shared
memory ordering and concurrent plugin admission. Build with `--enable-threads`
to select the `catty-jvm21-gate4` profile and image ABI `4`; it includes the
Gate 3 plugin admission surface. See the
[implementation specification](GATE4.md) and
[test plan](GATE4-TESTING.md), and run the differential suite with
`python3 -B tests/gate4.py --java-home /path/to/jdk-21 --race`.

Gate 5 adds application shutdown, process diagnostics and workload measurements.
Build with `--enable-lifecycle` to select the `catty-jvm21-gate5` profile and
image ABI `5`; it inherits Gate 4 threading and Gate 3 plugin admission. See the
[implementation specification and plan](GATE5.md), [acceptance and measurement
plan](GATE5-TESTING.md), and the checked-in
[Gate 5 evidence](out/gate5/EVIDENCE.md).

The [Catty 0.1 preview plan](CATTY-PREVIEW.md) covers packaging, installation,
the first pilot workload, and release acceptance.

## Requirements

- Go 1.21 or newer (validated with **go1.27.1** on darwin/arm64).
- **OpenJDK 21** providing `java`, `javac`, `javap` and `jar` (validated with
  Temurin **21.0.13.1**). Select it with `--java-home` or `JAVA_HOME`.
- Python 3.9+ for the acceptance runner.

On this machine a non-login tool shell does not resolve `go` automatically;
prefix commands with `PATH=/usr/local/go/bin:$PATH` when needed.

## Run the Gate 3 suite

```sh
python3 -B tests/gate3.py --java-home /path/to/jdk-21 --prepare-only
python3 -B tests/gate3.py --java-home /path/to/jdk-21 --case late-override
python3 -B tests/gate3.py --java-home /path/to/jdk-21
```

The runner compiles the compile-only `catty.runtime` API stub, the shared
application classes and every plugin version, builds each Catty executable
before its JARs exist, records the OpenJDK 21 oracle, then runs all-interpreted,
mixed and forced-AOT images against the same bytes. It checks archive and JAR
policy directly through the Java registration intrinsic, exercises rollback and
standalone execution, and prints the temporary artifact directory. A full run
ends with `Gate 3 passed.`

Build a plugin-enabled image by hand:

```sh
out/bin/catty build --enable-plugins --classpath app-classes --main app.Main \
  --mode auto --report report.json --output app
```

## Run the Gate 4 suite

```sh
python3 -B tests/gate4.py --java-home /path/to/jdk-21 --prepare-only
python3 -B tests/gate4.py --java-home /path/to/jdk-21 --case wait-set
python3 -B tests/gate4.py --java-home /path/to/jdk-21 --race
```

The runner compiles the thread fixtures, records the OpenJDK 21 oracle, and
runs all-interpreted, mixed and forced-AOT images against the same class files
with tracing enabled. With `--race` it also builds and runs race-instrumented
executables in all three modes and checks their Go metadata for `-race=true`.
A full run ends with `Gate 4 differential tests passed.`

Build a threaded image by hand:

```sh
out/bin/catty build --enable-threads --classpath app-classes --main app.Main \
  --mode auto --report report.json --output app
```

## Run the Gate 5 suite

```sh
python3 -B tests/gate5.py --java-home /path/to/jdk-21 --prepare-only
python3 -B tests/gate5.py --java-home /path/to/jdk-21
python3 -B tests/gate5.py --java-home /path/to/jdk-21 --race
python3 -B tests/gate5_bench.py --java-home /path/to/jdk-21 --smoke
python3 -B tests/gate5_bench.py --java-home /path/to/jdk-21
```

`tests/gate5.py` compiles each fixture once with `javac --release 21 -g`, checks
the OpenJDK 21 oracle, and runs the same class files in all-interpreted, mixed
and forced-AOT images under the Gate 5 profile. It covers startup arguments and
main-class initialization, natural non-daemon drain, hook registration and
concurrency, explicit exit statuses and encoding, exit contention, halt,
late-admitted plugin hooks, diagnostics setup failures, a host-fault fixture,
and the tracing and metrics contracts. It also checks CLI precedence, the exact
profile/ABI/feature pair, earlier-profile build refusal, earlier-profile late
admission refusal, and race-instrumented executables. A successful run ends
with `Gate 5 differential tests passed.`

`--prepare-only` runs only the JDK path. `--case NAME` runs one lifecycle case.
The full measurement runner is described under
[Application measurements](#application-measurements).

Build a Gate 5 executable by hand:

```sh
out/bin/catty build --enable-lifecycle --classpath app-classes --main app.Main \
  --mode auto --report report.json --output app
```

### Gate 5 lifecycle profile

`--enable-lifecycle` selects `catty-jvm21-gate5`, image ABI `5`, and inherits
`--enable-threads` and `--enable-plugins`. In descending precedence,
`--enable-lifecycle`, `--enable-threads` and `--enable-plugins` select Gates 5, 4
and 3; the default remains Gate 2. The image and report serialize
`lifecycle: true`, and installation rejects a mismatched profile, ABI or
feature flag before Java runs.

The profile publishes `System.exit(int)`, `System.nanoTime()`,
`Runtime.getRuntime()`, `Runtime.exit`, `Runtime.halt`,
`Runtime.addShutdownHook`, `Runtime.removeShutdownHook`, `Runtime.gc()` and
`java.lang.IllegalStateException`. The process has one running state, one
shutdown state with a frozen identity-keyed hook snapshot, and one terminal
state with the selected status. Natural shutdown commits when main and every
live non-daemon peer first reach zero; explicit exit may start while other
threads are live; halt skips or cuts short graceful shutdown. Successful exit
and halt calls never return into Java: the requesting thread and its monitors
stay observable until process termination, so `finally` bodies do not run.
Exit statuses keep their full signed value and are encoded by the host; on POSIX
`System.exit(-3)` is observed as `253` and `halt(19)` as `19`. Setup failures
(decoding, image installation, trace or metrics output) exit `2` without
entering Java.

`CATTY_TRACE` JSON lines gain `app_start`, `shutdown_begin` and a terminal
`process_exit` that seals the trace, while `CATTY_METRICS=<path>` writes one JSON
object with publication, execution, reflection, exception, monitor and thread
counters plus process-wide Go heap and GC observations. Tracing and metrics are
off by default and never write to application stdout.

Signals and OS logoff hooks, security managers, custom uncaught-exception
handlers, hooks on virtual threads, thread subclasses that override `start`, and
other `Runtime` APIs remain outside Gate 5. Gate 4 platform-thread restrictions
continue to apply.

### Application measurements

`tests/gate5_bench.py` compiles one deterministic Java workload and compares
OpenJDK 21 with Catty's interpreted, mixed and forced-AOT configurations using
the same class files, operation counts, worker counts, seed and inputs. It
validates engine placement and result checksums before timing, then collects
build time/size, cold time to first verified result, steady throughput with
p50/p99 normalized batch latency, peak process RSS, and a retained-graph heap
observation. The full run uses seven fresh processes per configuration and 30
post-warmup batches per process; `--smoke` uses a reduced matrix. Raw evidence
and `SUMMARY.json`/`SUMMARY.txt` are written under
[out/gate5/bench](out/gate5/bench); the checked-in summary is
[out/gate5/EVIDENCE.md](out/gate5/EVIDENCE.md).

## Run the Gate 2 suite

```sh
python3 -B tests/gate2.py --java-home /path/to/jdk-21
```

The suite compiles every corpus fixture, records the OpenJDK 21 oracle output,
builds Catty, and runs all ten cases in forced-AOT, mixed and all-interpreted
modes against the same class files. It also checks the malformed `long`-slot
fixture, the profile/image inventory, wrong-profile and wrong-ABI images, a
missing bootstrap member, generated Go for arithmetic/switch/exception offsets,
the mixed executable with no JDK available, and finally the whole Gate 1 suite.
A successful run ends with:

```text
Gate 2 passed.
```

The positive Gate 2 corpus can also run under the Gate 3 profile to exercise
the shared runtime helpers:

```sh
python3 -B tests/gate2.py --java-home /path/to/jdk-21 --enable-plugins
```

It ends `Gate 2 positive corpus passed with Gate 3 enabled.` and skips the
Gate 2-specific wrong-profile/negative checks and the separately-run Gate 1
suite.

`--prepare-only` compiles the fixtures and checks the OpenJDK results without
building Catty. `--case NAME` runs one corpus case without the full gate checks:

```sh
python3 -B tests/gate2.py --java-home /path/to/jdk-21 --prepare-only
python3 -B tests/gate2.py --java-home /path/to/jdk-21 --case long-crossing
```

The runner prints the temporary directory holding its class files, reports,
generated Go, traces and logs. The checked-in handoff evidence is under
[out/gate2/EVIDENCE.md](out/gate2/EVIDENCE.md); Gate 1 evidence remains under
[out/EVIDENCE.md](out/EVIDENCE.md).

## Build a fixture by hand

```sh
export PATH=/usr/local/go/bin:$PATH
mkdir -p out/bin out/classes out/reports out/generated
go build -o out/bin/catty ./cmd/catty
javac --release 21 -g -d out/classes testdata/gate2/long-crossing/LongCrossing.java

out/bin/catty build \
  --classpath "$PWD/out/classes" --main LongCrossing --mode auto \
  --report out/reports/long.json \
  --output out/bin/long-crossing \
  --generated-source out/generated

CATTY_TRACE="$PWD/out/traces/long.jsonl" ./out/bin/long-crossing
```

The build report records the profile, image ABI, the supported opcode and
bootstrap inventories, per-method engine selections, retained interpreter
bytecode, and any application references retained for lazy linkage. The image
carries the same profile and ABI; `rt.Install` rejects a wrong image before any
initializer or `main` runs.

## Repository layout

| Path | Contents |
| --- | --- |
| `cmd/catty` | CLI entry point |
| `internal/classfile` | binary reader, constant pool, attributes, all admitted instructions, stack maps, `ConstantValue`/nest metadata |
| `internal/verify` | typed worklist verification, category-two local/operand slots, stack maps, logical stack shapes |
| `internal/aot` | JVM bytecode → Go basic-block lowering (every admitted instruction except `irem`) |
| `internal/build` | ordered inputs, JAR/directory loading, deferred-linkage retention, engine selection, report/image generation, staging and `go build` |
| `rt` | image installation, class/object/array model, dispatch and defaults, initialization, reflection, boxing/numeric helpers, interpreter, launcher, trace, lifecycle coordinator, metrics |
| `testdata` | checked-in Java fixtures, including the ten-case Gate 2 corpus and the Gate 5 lifecycle/benchmark fixtures |
| `tests` | Gate 2–5 differential runners, Gate 1 acceptance runner, class-file patch helper, application benchmark runner |
| `out/gate2` | durable Gate 2 evidence: reports, traces, generated Go and test logs |
| `plugin` | narrow Gate 3 launcher adapter connecting admission to the runtime |
| `internal/admission` | bounded ZIP/JAR validation, policy checks and staging |
| `testdata/gate3` | Gate 3 API stub, OpenJDK adapter and plugin fixtures |
| `out/gate3` | durable Gate 3 acceptance evidence |
| `out/gate4` | durable Gate 4 acceptance evidence |
| `out/gate5` | durable Gate 5 acceptance evidence, benchmark summary and raw samples |

## Profile limits (Gate 2)

The profile remains deliberately closed and fails closed; the report records
the selected engine and reason for every admitted method.

- **Execution** — one logical Java thread. Values include references, `void`
  returns, the `int` family, `long`, `float` and `double`; ordinary arithmetic,
  comparisons, conversions, bitwise operations, shifts, conditional branches,
  `goto`/`goto_w`, `tableswitch`, `lookupswitch`, legal stack manipulation and
  `wide` local indexes. No `invokedynamic`/condy, method handles, monitors,
  `synchronized` methods, native methods, legacy subroutines, reserved opcodes,
  threads, JNI, agents or shutdown hooks.
- **Objects and arrays** — construction, fields, casts, `instanceof`, all eight
  primitive array types, reference and multidimensional arrays, typed loads and
  stores, `arraylength`, and partially allocated arrays. Java's check order
  (null, bounds, then store type) is preserved.
- **Calls and errors** — `invokestatic`, `invokespecial`, `invokevirtual` and
  `invokeinterface`; class, member, package, protected-receiver and nestmate
  access checks; constructor and super calls; interface defaults with
  maximally-specific selection; linkage errors; `athrow` and exception handlers.
- **Initialization** — `ConstantValue` applied by the JVMS 5.5 procedure,
  superclass and required superinterface order, reentry, failed initialization,
  `ExceptionInInitializerError` cause identity, and later
  `NoClassDefFoundError` for erroneous classes.
- **Reflection** — class mirrors, primitive mirrors and the eight boxed
  primitives (including `TYPE`, `valueOf` and unboxing), `Object.getClass`,
  public inherited/overriding/covariant method lookup through `Class.getMethod`,
  and `Method.invoke` with caller-sensitive access checks, JLS 5.3 primitive
  widening, target dispatch, boxing and `InvocationTargetException` cause
  identity.
- **Inputs** — class file version 65.0 only; ordered class directories and
  single-release JARs. Multi-release JARs, named modules and `module-info.class`
  are rejected. The bootstrap namespace is reserved.
- **Linkage** — unresolved application symbols whose absence does not prevent
  verification are retained in the report and resolved lazily at the referring
  instruction; unavailable bootstrap facilities fail the build. Resolution
  outcomes, including failures, are cached.
- **AOT lowering** — every admitted instruction except `irem`, which retains the
  Gate 1 interpreter fallback and forced-AOT rejection. Dispatch decisions are
  resolved at run time rather than through precomputed vtables.
- **Lifecycle** — the Gate 2 profile covers application arguments, normal
  `main` return, uncaught Java throwable exit status and `CATTY_TRACE` evidence.
  The opt-in Gate 5 profile adds `System.exit`/`Runtime.exit`/`Runtime.halt`,
  shutdown hooks, the frozen-registry policy, process trace events and metrics;
  it retains the Gate 2 behavior when a build omits `--enable-lifecycle`.
- **Plugins (Gate 3 only)** — one unsigned single-release classpath JAR per
  registration, bounded at 16 MiB compressed, 64 MiB expanded and 64 published
  loaders. Parent-first loaders, loader-aware dispatch and class identity,
  reflective lookup, and rejected-transaction rollback are implemented. The
  runner documents exact assertions and remaining limitations.

## Verification notes

`go vet ./...` and `go test -race -count=1 ./...` are clean. The Gate 5
acceptance, race-enabled acceptance and measurement suites pass on OpenJDK 21,
and the Gate 1–4 regressions pass under their original profiles. Exact
commands, results, traces, metrics and measurements are recorded in
[out/gate5/EVIDENCE.md](out/gate5/EVIDENCE.md); earlier evidence remains under
`out/gate2/`, `out/gate3/` and `out/gate4/`.
