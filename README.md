# Catty — Gate 1

Catty builds one Java 21 application into one standalone Go executable. A Go
AOT emitter turns verified JVM bytecode into real Go functions, an interpreter
retains the methods outside the AOT lowering profile, and both engines share the
same object model, dispatch, class initialization, reflection and exception
semantics.

This repository implements the Gate 1 delivery in [IMPLEMENTATION.md](IMPLEMENTATION.md)
against [ARCHITECTURE.md](ARCHITECTURE.md); [ACCEPTANCE.md](ACCEPTANCE.md) is the
evidence contract and `tests/acceptance.py` is the suite.

## Requirements

- Go 1.21 or newer (validated with **go1.27.1** on darwin/arm64).
- A JDK accepting `javac --release 21` (validated with **Temurin 25.0.4.1**),
  including `java`, `javac`, `javap` and `jar`.
- Python 3.9+ for the acceptance runner.

On this machine a non-login tool shell does not resolve `go` automatically;
prefix build commands with `PATH=/usr/local/go/bin:$PATH` when needed. Login
shells and the acceptance command run from the repository root work as written.

## Run the acceptance suite

```sh
python3 -B tests/acceptance.py
```

Run it from a shell whose `PATH` includes Go. On this machine the Go
toolchain is installed at `/usr/local/go/bin`; a non-login tool shell must
prefix `PATH=/usr/local/go/bin:$PATH`, while a login shell already resolves
`go`. The checked-in evidence log was produced by the bare command above.

The suite builds `./cmd/catty`, compiles the fixture classes in a fresh work
directory, runs the OpenJDK oracles, builds and launches the application in all
three engine modes, checks the traces, and runs `go test -count=1 -json ./...`
including the two runtime acceptance tests. A successful run ends with:

```text
Ran 8 tests in ...s
OK
```

`--prepare-only` compiles fixtures and checks the OpenJDK oracle outcomes
without building Catty.

```sh
python3 -B tests/acceptance.py --prepare-only
```

## Build and run the principal fixture by hand

```sh
mkdir -p out/bin out/classes out/reports out/generated out/traces
go build -o out/bin/catty ./cmd/catty
javac --release 21 -g -d out/classes testdata/crossing/Crossing.java testdata/crossing/Main.java

out/bin/catty build \
  --classpath "$PWD/out/classes" --main Main --mode auto \
  --policy testdata/crossing/mixed.json \
  --report out/reports/mixed.json \
  --output out/bin/crossing-mixed \
  --generated-source out/generated

CATTY_TRACE="$PWD/out/traces/mixed.jsonl" ./out/bin/crossing-mixed
```

The mixed evidence used by the handoff is checked in under `out/`; see
[out/EVIDENCE.md](out/EVIDENCE.md) for the report, generated-source and trace
line references, including the engine-selection, nested-callback and throwable
identity events.

## Repository layout

| Path | Contents |
| --- | --- |
| `cmd/catty` | CLI entry point |
| `internal/classfile` | binary reader, constant pool, attribute and instruction decoder, stack maps |
| `internal/verify` | structural/type verification, stack-height analysis, policy capability context |
| `internal/aot` | JVM bytecode → Go basic-block lowering |
| `internal/build` | ordered inputs, JAR/directory loading, engine selection, report/image generation, staging and `go build` |
| `rt` | image installation, class/object model, dispatch, initialization, reflection intrinsics, interpreter, launcher, trace |
| `testdata` | checked-in Java fixtures |
| `tests` | acceptance runner, class-file patch helper, generated-source inspector |
| `out` | durable three-build evidence, reports, generated source, traces and test logs |

## Profile limits (Gate 1)

The profile is intentionally narrow and fails closed; the report records the
selected engine and reason for every admitted method.

- **Execution** — one logical Java thread, category-one int/integer-family and
  reference values only. No `long`, `float`, `double`, primitive arrays,
  `invokedynamic`/condy, method handles, monitors, `synchronized` methods,
  native methods, threads, JNI, agents, or shutdown hooks.
- **Inputs** — class file version 65.0 only; class directories and
  single-release JARs. Multi-release JARs, named modules and `module-info.class`
  are rejected. The bootstrap namespace is reserved.
- **Dispatch** — one application loader and a fixed bootstrap inventory.
  Dynamic class loading/plugin admission (Gate 3) is not implemented.
- **Verification** — Java 21 verification rules for the admitted instruction
  subset, including stack maps, local/operand types, receiver assignability,
  `uninitializedThis` and each `new` alias; the full JVM verifier and every
  category-two rule are out of scope.
- **AOT lowering** — every admitted instruction except `irem`; auto mode falls
  back to the interpreter for a method containing `irem` and reports the reason.
  Dispatch decisions are resolved at run time rather than through precomputed
  vtables, which preserves correctness at some cost.
- **Bootstrap library** — `Object`, throwable/exception hierarchies, `System.out`,
  `PrintStream.println(I)`, `Class.getMethod`, `Method.invoke`,
  `InvocationTargetException.getCause`, and `Integer` boxing/unboxing for `int`.
  Reflection covers public method lookup/invocation, inherited lookup, argument
  conversion and target-throwable wrapping; constructors, fields, access checks,
  annotations, generic signatures and non-int boxing are not part of Gate 1.
- **Class state** — single-threaded preparation/initialization, constant-field
  and failure semantics used by the fixtures; `ConstantValue` attributes are not
  used for static initialization, and line-number metadata is parsed but not
  printed in uncaught stack traces.
- **Lifecycle** — application arguments, normal main return, uncaught Java
  throwable exit status and `CATTY_TRACE` evidence are implemented. `System.exit`,
  non-daemon thread completion and shutdown hooks are later gates.

## Verification notes

`go vet ./...` and `go test -count=1 ./...` are clean. The acceptance evidence
and exact test logs are stored under `out/logs/`.
