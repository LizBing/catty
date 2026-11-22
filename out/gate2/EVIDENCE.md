# Catty Gate 2 evidence

This directory contains artifacts from a real local run of
`tests/gate2.py` on 2026-11-22 (Asia/Shanghai). The run used OpenJDK 21 for
every reference execution, the same compiled class files for the reference,
interpreted, mixed and forced-AOT runs, and the repository's current
`codex/gate-2` tree.

- Runner: `python3 -B tests/gate2.py --java-home <Temurin 21>`
- Result: **exit 0 — `Gate 2 passed.`**
- Full log: [`logs/full-gate.log`](logs/full-gate.log)
- Toolchains: [`logs/toolchains.txt`](logs/toolchains.txt)
- Temporary run directory retained by the runner:
  `/private/var/folders/3_/pzk1z58n29jclr9vvlhhhtyr0000gp/T/catty-gate2-9fedbnvr`

## Toolchains

| Tool | Version |
| --- | --- |
| Go | go1.27.1 darwin/arm64 |
| JDK | Temurin OpenJDK 21.0.13.1+2 (class files `65.0`) |
| Python | 3.9.6 |
| OS | Darwin arm64 |

## Published profile

| Field | Value |
| --- | --- |
| Profile | `catty-jvm21-gate2` |
| Image ABI | `2` |
| Bootstrap classes / methods / fields | 44 / 80 / 9 |
| Admitted opcode names | 196 |

Every build report in [`reports/`](reports) carries the same profile, image ABI,
opcode inventory and bootstrap inventory. The image embedded by a build carries
the identical values; `rt.Install` checks them before any startup initializer or
`main` method runs.

## Differential corpus

The runner compiled each fixture once, recorded the OpenJDK 21 stdout, and then
required byte-identical stdout, stderr and exit status from Catty's forced-AOT,
mixed and all-interpreted builds. 13 OpenJDK runs and one verifier rejection
were pinned during preparation.

| Case | Fixed runs | Mixed engine crossings | Report |
| --- | --- | --- | --- |
| `long-crossing` | 2 | `main` AOT → `bridge` interpreter → `leaf` AOT (retains 23 bytes) | [report](reports/long-crossing-auto.json) |
| `initialization` | 1 | `WithDefault.<clinit>` interpreter, rest AOT (retains 9 bytes) | [report](reports/initialization-auto.json) |
| `numeric-edge-cases` | 2 | `leaf` interpreter, 32 AOT methods | [report](reports/numeric-edge-cases-auto.json) |
| `switch-and-wide` | 1 | `dense` interpreter (88 bytes), 15 AOT methods including `wide` locals | [report](reports/switch-and-wide-auto.json) |
| `arrays` | 2 | `lengthOf` interpreter, typed/multidimensional arrays AOT | [report](reports/arrays-auto.json) |
| `dispatch-and-constructors` | 1 | `note` interpreter; virtual/interface/super dispatch AOT | [report](reports/dispatch-and-constructors-auto.json) |
| `constant-values-and-reentry` | 1 | `mark` interpreter; initialization order and failures AOT | [reports](reports) |
| `reflection` | 1 | `leaf` interpreter; getMethod/invoke AOT | [report](reports/reflection-auto.json) |
| `deferred-linkage` | 1 | `callGood` interpreter; unresolved references retained | [report](reports/deferred-linkage-auto.json) |
| `profile-and-image` | 1 | `mix` interpreter; inventory checked | [report](reports/profile-and-image-auto.json) |

The mixed reports show genuine crossings rather than a whole-program engine
choice. Traces agree with the reports, stay on one logical Java thread
(`thread: 1`) and one OS process, and balance every entry with a return or
unwind:

| Trace | Entries | Engines | Threads / PIDs |
| --- | --- | --- | --- |
| [`long-crossing-auto-0`](traces/long-crossing-auto-0.jsonl) | 5 | AOT, interpreter, intrinsic | 1 / 1 |
| [`numeric-edge-cases-auto-0`](traces/numeric-edge-cases-auto-0.jsonl) | 322 | AOT, interpreter, intrinsic | 1 / 1 |
| [`switch-and-wide-auto-0`](traces/switch-and-wide-auto-0.jsonl) | 61 | AOT, interpreter, intrinsic | 1 / 1 |
| [`initialization-auto-0`](traces/initialization-auto-0.jsonl) | 22 | AOT, interpreter, intrinsic | 1 / 1 |
| [`reflection-auto-0`](traces/reflection-auto-0.jsonl) | 241 | AOT, interpreter, intrinsic | 1 / 1 |

The `long-crossing` mixed trace nests `LongCrossing.main` (AOT, depth 0),
`LongCrossing.bridge` (interpreter, depth 1) and `LongCrossing.leaf` (AOT,
depth 2), which is the crossing the acceptance runner asserts. Engine identity
in each `enter` event is compared against the build report.

## Representative generated Go

These are the files actually compiled into the standalone executables. The
acceptance runner also parses them on every full run and checks arithmetic,
switch control flow and original exception offsets.

### Arithmetic and remainder — [`generated/numeric-edge-cases/generated.go`](generated/numeric-edge-cases/generated.go)

`cattyMethod27` lowers `NumericEdgeCases.remainders(JIFD)I` (lines
[`4646`–`5286`](generated/numeric-edge-cases/generated.go#L4646-L5286)):

```go
// line 4680: lrem with a zero check that transfers to the original handler
stack[0] = rt.LongValue(rt.AsLong(stack[0]) % d)
// line 4709: frem
stack[0] = rt.FloatValue(rt.FRem(rt.AsFloat(stack[0]), rt.AsFloat(stack[1])))
// line 4717: drem
stack[0] = rt.DoubleValue(rt.DRem(rt.AsDouble(stack[0]), rt.AsDouble(stack[1])))
```

Long shifts in `cattyMethod4` use ordinary Go shifts with the masked Java
distance ([`line 233`](generated/numeric-edge-cases/generated.go#L233)):

```go
stack[0] = rt.LongValue(rt.AsLong(stack[0]) << rt.LongShiftDistance(rt.AsInt(stack[1])))
```

### Switch control flow — [`generated/switch-and-wide/generated.go`](generated/switch-and-wide/generated.go)

`cattyMethod1` lowers `SwitchWide.offsetDense(I)I` (lines
[`28`–`109`](generated/switch-and-wide/generated.go#L28-L109)). The generated body owns a
Go `switch` over the case keys, not an interpreter handoff
([`lines 39`–`47`](generated/switch-and-wide/generated.go#L39-L47)):

```go
switch rt.AsInt(stack[0]) {
case 10:
	block = 1
	continue
case 11:
	block = 2
	continue
...
default:
	block = 11
```

`SwitchWide.wideInts` and `wideLongs` exercise `wide` local indexes; the
interpreter and AOT use the same retained local slot numbers.

### Original exception offsets — [`generated/initialization/generated.go`](generated/initialization/generated.go)

`cattyMethod14` lowers `InitSequence.main([Ljava/lang/String;)V` (lines
[`372`–`837`](generated/initialization/generated.go#L372-L837)). The class file exception
table for that method is `[27,36)→39` and `[63,72)→75`; generated transfers use
the throwing instruction's original bytecode offset and the original handler
offset ([`lines 442`–`451`](generated/initialization/generated.go#L442-L451)):

```go
if hpc, ok := t.HandleThrow(m, 18, thrown); ok {
	stack = [3]rt.Value{}
	stack[0] = rt.RefValue(thrown)
	switch hpc {
	case 39:
		block = 4
	case 75:
		block = 9
	default:
		panic("catty: generated handler transfer has no target block")
	}
	continue
}
```

The full-run check parses `InitSequence.class`, extracts handler offsets 39 and
75, requires both generated `case` arms, and requires at least one
`HandleThrow(m, pc, ...)` whose `pc` lies inside a protected exception-table
range.

## Deferred linkage

`Link` compiles against `Dep` with `good()I` and `value:I`. The gate runs with
a classpath whose `Dep.class` no longer has either member. Startup prints `1`
before `Dep` is touched, proving verification and loading do not force method
resolution. Later uses throw `NoSuchMethodError` twice and `NoSuchFieldError`
once, in Java order, then the class remains usable (`Link.finish()`):

[`report`](reports/deferred-linkage-auto.json) ·
[`trace`](traces/deferred-linkage-auto-0.jsonl) · runtime output `1\n2\n3\n4\n5\n`.

The report retains both application references with the referring method and
bytecode offset:

```json
{"kind":"method","owner":"Dep","name":"good","descriptor":"()I",
 "referrer":"Link","referrerMethod":"unusedCall","referrerDescriptor":"()V","pc":0}
{"kind":"field","owner":"Dep","name":"value","descriptor":"I",
 "referrer":"Link","referrerMethod":"unusedField","referrerDescriptor":"()V","pc":0}
```

## Negative and profile evidence

- **Malformed `long` slot** — the suite patches `lload_3` in
  `LongCrossing.bridge(JI)J` to read the second slot of the incoming `long`.
  Build output:
  `LongCrossing.bridge(JI)J: pc 10: local index 1 is the unusable second slot of the long at local 0`.
  See [`logs/bad-long-slot-build.json`](logs/bad-long-slot-build.json).
- **Malformed stack shape** — the runner changes the first `iadd` in
  `NumericEdgeCases.intArith(II)I` to `ladd`. The verifier reports
  `NumericEdgeCases.intArith(II)I: pc 2: expected long on the operand stack, found int`
  ([build log](logs/bad-stack-shape-build.json)).
- **Malformed control flow** — the runner redirects the default target of
  `SwitchWide.dense(I)I`'s `tableswitch` into the middle of an instruction. The
  verifier rejects the switch with
  `SwitchWide.dense(I)I: pc 1: stack map has no frame at switch default 58`
  ([build log](logs/bad-control-flow-build.json)).
- **Wrong image ABI** — a generated image mutated to ABI `1` exits `2` with
  `catty: unsupported image ABI 1 (runtime ABI 2)` and no Java output:
  [`logs/image-wrong-abi.json`](logs/image-wrong-abi.json).
- **Wrong image profile** — profile `catty-jvm21-gate1` exits `2` with
  `catty: unsupported image profile "catty-jvm21-gate1" (runtime profile "catty-jvm21-gate2")`:
  [`logs/image-wrong-profile.json`](logs/image-wrong-profile.json).
- **Missing bootstrap member** — a fixture using `System.nanoTime()` fails the
  build with
  `needs unavailable bootstrap method java/lang/System.nanoTime()J`, releases no
  report, and a prior executable at the same output path remains byte-identical:
  [`logs/missing-bootstrap-build.json`](logs/missing-bootstrap-build.json).
- **Gate 1 fallback preserved** — the Gate 1 regression run passes all eight
  tests, including automatic interpretation for `irem` and forced-AOT rejection
  with method and bytecode offset: [`logs/gate1-regression.txt`](logs/gate1-regression.txt).

## Standalone mixed executable

The runner copies the mixed `long-crossing` executable into an empty directory,
renames its input class directory, and launches it with `PATH=""`,
`JAVA_HOME=""` and `JDK_HOME=""`. It prints `270544960\n1348497537\n`, exits `0`
and produces the same engine trace:

- binary: [`bin/long-crossing-mixed`](bin/long-crossing-mixed)
- launch record: [`logs/standalone.json`](logs/standalone.json)
- trace: [`traces/standalone.jsonl`](traces/standalone.jsonl)

## Runtime unit checks

`go vet ./...` and `go test -count=1 ./...` are clean
([vet](logs/go-vet.txt), [test](logs/go-test.txt)). Besides the Gate 1 runtime
tests, the Gate 2 tree adds focused checks for:

- JVMS 5.5 `ConstantValue` preparation: a superclass initializer observes its
  subclass's constant field as 42 before the subclass initializer runs.
- Interpreter stack manipulation (`swap`) and `goto_w`, which javac does not
  normally emit.

The corpus supplies the differential coverage for the other stack forms,
`lrem`/`frem`/`drem`, all primitive and multidimensional arrays, switches,
`wide` locals, initialization failures and reflection.

## Gate 1 regression

The full Gate 2 command runs the untouched Gate 1 acceptance suite on the same
JDK 21 toolchain. Its log records `Ran 8 tests ... OK` for:

`test_class_initialization_and_hidden_fields`, `test_crossing_and_standalone`,
`test_dispatch_and_reflection`, `test_exceptions_and_exclusive_range`,
`test_fallback_and_failed_replacement`, `test_invalid_inputs`,
`test_packaging_and_classpath`, `test_runtime_and_go_packages`.

See [`logs/gate1-regression.txt`](logs/gate1-regression.txt).

## Reproduce

```sh
export PATH=/usr/local/go/bin:$PATH
export JAVA_HOME=/path/to/jdk-21
python3 -B tests/gate2.py --java-home "$JAVA_HOME"
```

The runner prints a fresh temporary artifact directory. This checked-in
evidence is a copy of what the command produced; per-method reports for every
mode, all corpus traces and the full generated Go files remain in that
directory for detailed review.

## Scope retained for later gates

Runtime JAR loading and plugin admission remain Gate 3, concurrency and thread
lifecycle remain Gate 4, and full application shutdown (`System.exit`,
shutdown hooks, complete non-daemon thread drain) remains Gate 5. Gate 2 keeps
one logical Java thread and one application loader.
