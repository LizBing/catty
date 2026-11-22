# Catty Gate 1 evidence

This directory contains artifacts from a real local run on darwin/arm64. The
three Crossing executables were built from the same class directory and policy,
and the OpenJDK oracle prints `2` for the same fixture.

## Toolchains

| Tool | Version |
| --- | --- |
| Go | go1.27.1 darwin/arm64 |
| JDK | Temurin OpenJDK 25.0.4.1 (`javac --release 21` emits v65.0) |
| Python | 3.9.6 |
| OS | Darwin arm64 |

`out/logs/toolchains.txt` records the version commands.

## Artifacts

| Artifact | Path |
| --- | --- |
| Catty builder | [`bin/catty`](bin/catty) |
| Compiled fixture classes | [`classes/`](classes) |
| Mixed build (auto + policy) | [`bin/crossing-mixed`](bin/crossing-mixed) |
| All-AOT build | [`bin/crossing-aot`](bin/crossing-aot) |
| All-interpreted build | [`bin/crossing-interp`](bin/crossing-interp) |
| Mixed report | [`reports/mixed.json`](reports/mixed.json) |
| All-AOT report | [`reports/aot.json`](reports/aot.json) |
| All-interpreted report | [`reports/interp.json`](reports/interp.json) |
| Mixed generated source | [`generated/generated.go`](generated/generated.go) |
| Mixed generated module | [`generated/main.go`](generated/main.go), [`generated/go.mod`](generated/go.mod) |
| Mixed trace | [`traces/mixed.jsonl`](traces/mixed.jsonl) |
| All-AOT trace | [`traces/aot.jsonl`](traces/aot.jsonl) |
| All-interpreted trace | [`traces/interp.jsonl`](traces/interp.jsonl) |
| Standalone trace | [`traces/standalone.jsonl`](traces/standalone.jsonl) |
| Test logs | [`logs/`](logs) |

Every run printed exactly `2\n`, exited zero, and every trace came from one OS
process and one logical Java thread (`thread: 1`) with balanced frame depth.

## Engine selection

Mixed build report ([`reports/mixed.json`](reports/mixed.json)):

| Method | Engine | Retained bytes | Generated symbol |
| --- | --- | --- | --- |
| `Crossing$Box.bump()V` | aot | 0 | `cattyMethod1` |
| `Crossing$Step.apply(LCrossing$Box;)V` | abstract | 0 | — |
| `Crossing$Worker.apply(LCrossing$Box;)V` | interp | 9 | — |
| `Crossing.probe()I` | aot | 0 | `cattyMethod5` |

All-AOT ([`reports/aot.json`](reports/aot.json)) selects `aot` for every
concrete application method, including `Worker.apply` and `bump`; all-interpreted
([`reports/interp.json`](reports/interp.json)) selects `interp` and retains each
method's original Code length (for example `probe` retains 99 bytes, `apply` 9,
`bump` 11). The three reports share the same profile (`catty-jvm21-gate1`),
image ABI (`1`), input origin and digest.

## Generated Go review

### `Box.bump` — real field mutation

`generated/generated.go` lines [`50`–`80`](generated/generated.go#L50-L80):
`cattyMethod1` emits its own lowered body, not an interpreter handoff. It reads
`count` from the receiver, increments an int32 expression, writes the same field
on the same receiver, and returns:

- line [`63`](generated/generated.go#L63): `v, thrown := t.GetField(... "Crossing$Box", "count", "I")`
- line [`70`](generated/generated.go#L70): `rt.IntValue(rt.AsInt(...) + rt.AsInt(...))`
- line [`72`](generated/generated.go#L72): `t.PutField(stack[0], "Crossing$Box", "count", "I", stack[1])`

### `Crossing.probe` — lowered exception control flow

`generated/generated.go` lines [`151`–`545`](generated/generated.go#L151-L545):
`cattyMethod5` switches over Go basic blocks generated from the verified CFG.
The direct interface call is a real Go call to shared dispatch
([`247`](generated/generated.go#L247)) and every potentially throwing operation
has a generated transfer to the correct original handler block:

- handler `25` is block case `1` at [`268`](generated/generated.go#L268)
- handler `79` is block case `5` at [`466`](generated/generated.go#L466)
- the `invokeinterface` at bytecode offset 18 transfers to handler 25
  ([`250`](generated/generated.go#L250)–[`265`](generated/generated.go#L265))
- the reflective `Method.invoke` at bytecode offset 72 transfers to handler 79
  ([`447`](generated/generated.go#L447)–[`462`](generated/generated.go#L462))

The generated bodies call `rt` helpers for allocation, fields, dispatch,
exception matching and trace; there is no opcode dispatch loop inside either
function. `Worker.apply` is deliberately absent from the generated source
because the mixed policy retains it as bytecode for the interpreter.

## Mixed-engine trace

[`traces/mixed.jsonl`](traces/mixed.jsonl) has 37 JSON lines, one OS PID and one
`thread` value. Key events (line numbers in that file):

| Line | Event | Evidence |
| --- | --- | --- |
| 1 | enter `Main.main` | depth 0, engine `aot` |
| 2 | enter `Crossing.probe` | depth 1, engine `aot` |
| 15 | dispatch | `Crossing$Step.apply` → `Crossing$Worker.apply` (`invokeinterface`) |
| 16 | enter `Crossing$Worker.apply` | depth 2, engine `interp` |
| 17 | enter `Crossing$Box.bump` | depth 3, engine `aot` — AOT callback under interpreted frame |
| 19 | unwind `Crossing$Worker.apply` | original throwable ID `1` |
| 20 | catch `Crossing.probe` | throwPC 18, handlerPC 25, throwable `1` |
| 24 | reflect | `Crossing$Step.apply` → `Crossing$Worker.apply` |
| 25 | enter `Crossing$Worker.apply` | depth 3, engine `interp` (second entry) |
| 26 | enter `Crossing$Box.bump` | depth 4, engine `aot` (second callback) |
| 28 | unwind `Crossing$Worker.apply` | throwable `1` (same object field thrown again) |
| 29 | wrap | wrapper `2`, cause `1` |
| 31 | catch `Crossing.probe` | throwPC 72, handlerPC 79, throwable `2` (the wrapper) |
| 34 | return `Crossing.probe` | depth 1 |
| 37 | return `Main.main` | depth 0 |

This establishes:

1. the AOT caller enters an interpreted interface implementation
   (`dispatch` then `enter` engine `interp`);
2. the interpreted body calls back into AOT code (`bump`, depth +1, engine `aot`);
3. the first throwable reaches the AOT catch with the original ID (`1`), and
   the second throw propagates after the reflective call;
4. reflection selects the same interpreted target, and the reflective wrapper
   preserves the original throwable as its cause (`wrapper 2`, `cause 1`), which
   is exactly what the second catch observes;
5. frames balance back to depth 0 after both exceptional paths.

## Standalone launch

`bin/crossing-mixed` was copied alone into `standalone/`; it was launched with
`env -i PATH="" JAVA_HOME="" CATTY_TRACE=... ./crossing-mixed`. The class
directory was not present in that working directory and the binary prints `2`
([`logs/run-standalone.stdout`](logs/run-standalone.stdout)); its trace is
[`traces/standalone.jsonl`](traces/standalone.jsonl).

## Test and build logs

| Log | Result |
| --- | --- |
| [`logs/acceptance.txt`](logs/acceptance.txt) | `Ran 8 tests ... OK` (acceptance=0) |
| [`logs/go-test.txt`](logs/go-test.txt) | `go test -count=1 ./...` all packages pass |
| [`logs/go-test-runtime.txt`](logs/go-test-runtime.txt) | `TestCattyAcceptanceImageABI` and `TestCattyAcceptanceInterpreterOperandGC` pass |
| [`logs/go-vet.txt`](logs/go-vet.txt) | `go vet ./...` clean |
| [`logs/build-mixed.log`](logs/build-mixed.log), [`build-aot.log`](logs/build-aot.log), [`build-interp.log`](logs/build-interp.log) | all three builds exit 0 |
| [`logs/run-status.txt`](logs/run-status.txt) | mixed, aot, interp, standalone all exit 0 and print `2` |

## Reproduction

```sh
export PATH=/usr/local/go/bin:$PATH
go build -o out/bin/catty ./cmd/catty
javac --release 21 -g -d out/classes testdata/crossing/Crossing.java testdata/crossing/Main.java
out/bin/catty build --classpath "$PWD/out/classes" --main Main --mode auto \
  --policy testdata/crossing/mixed.json --report out/reports/mixed.json \
  --output out/bin/crossing-mixed --generated-source out/generated
CATTY_TRACE="$PWD/out/traces/mixed.jsonl" ./out/bin/crossing-mixed
python3 -B tests/acceptance.py
```
