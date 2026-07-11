# catty

An experimental JVM written in Go that **sits on top of the Go runtime**: it
reuses Go's garbage collector, scheduler, and allocator instead of implementing
its own. Java objects are Go objects (traced natively by Go's GC, with zero
write-barrier code); the planned concurrency model maps `java.lang.Thread` to a
goroutine.

**Phase 1 is done**: a correct, switch-dispatched bytecode interpreter running a
minimal-but-real Java subset, byte-identical to HotSpot on the test corpus.
Phase 2 (a bytecode → Go-source AOT transpiler) is the open performance arc.

## Documentation

| Document | What it covers |
|---|---|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Premise, execution pipeline, package responsibilities, data structures, an end-to-end trace, and the performance story |
| [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) | Build/run/test, project layout, conventions, and recipes for adding an opcode / native method / core class |
| [docs/ROADMAP.md](docs/ROADMAP.md) | Phase 2 AOT design, interpreter tuning, spec-coverage gaps, performance targets |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | Versioned work log |
| [docs/adr/](docs/adr/) | Architecture Decision Records — the *why* behind each choice |

## Status — what runs

`./tests/run.sh` compiles each fixture with `javac` and diffs catty's stdout
against real `java`. All currently pass:

| Fixture | Exercises |
|---|---|
| HelloWorld | `ldc` String, `System.out.println`, int multiply |
| Fibonacci | recursion, `if`/`ireturn`, `iadd` |
| Factorial | `long` arithmetic, `long[]`, StringBuilder concat |
| ArraySum | `newarray`/`iaload`/`iastore`, loops |
| OOPDemo | `new`, `<init>`, instance fields, `invokevirtual` dispatch |
| StaticFields | `<clinit>` static initializers |
| SwitchDemo | `tableswitch` (dense) + `lookupswitch` (sparse) |

The interpreter implements ~140 JVMS opcodes: constants, typed loads/stores,
full int/long/float/double arithmetic and conversions, shifts, comparisons,
all conditional branches and switches, stack manipulation, object/array
creation and access, field access, `invoke{virtual,special,static}`,
`checkcast`/`instanceof`, and returns. Core classes (`java.lang.Object`,
`String`, `StringBuilder`, `System`, `java.io.PrintStream`) are implemented
natively in Go rather than loaded from a JRE.

There are **two execution engines**, both byte-identical to `java` on the test
corpus: the default tree-walking interpreter (`Loop`) and a stack-eliminated IR
executor (`-ir`, `LoopIR`). The IR executor validates the lowering pass that
underpins the planned AOT transpiler — it is not (yet) faster (ADR-0006).

There is also an experimental **AOT transpiler** (`transpile.Emit`, A1): it
lowers a method's IR to Go source and compiles it natively. On `fib(35)` the
emitted Go runs in ~44 ms — native speed, ~100× the interpreter and on par with
HotSpot's JIT. Currently int-only / single-method (`fib`); see
[docs/ROADMAP.md](docs/ROADMAP.md) Theme A.

## Quickstart

```sh
go build -o catty ./cmd/jvm          # build
./catty -cp <classpath> <MainClass>  # run (tree-walking interpreter)
./catty -cp <classpath> -ir <Main>   # run via the lowered IR executor

./tests/run.sh                       # e2e: compile fixtures, diff catty vs java
go test ./...                        # unit tests
```

Requires Go 1.22+ and a JDK (`javac`/`java`) on `PATH`. `-cp` is
colon-separated dirs/jars (default `.`); `<MainClass>` may use dots or slashes.
Set `CATTY_DEBUG=1` for a Go stack trace on a VM error.

## Performance baseline (fib(35), ~29M recursive calls)

| Engine | Time | Relative |
|---|---|---|
| catty (interpreter) | 4.34 s | 1× |
| `java -Xint` (HotSpot interpreter) | 0.61 s | ~7× catty |
| `java` (HotSpot JIT) | 0.05 s | ~87× catty |

The ~7× gap to `java -Xint` is **interpreter headroom** (switch vs. computed-goto
dispatch; 16-byte slots vs. machine words; per-call frame allocation) — tunable.
The ~87× gap to JIT is **JIT headroom**, closeable only by the Phase 2 AOT
transpiler. See [docs/ROADMAP.md](docs/ROADMAP.md) Themes A and B.

## Out of scope (planned or deferred)

Concurrency (single-threaded today; `Thread`→goroutine + monitors later),
`invokedynamic`/lambdas, exceptions/try-catch, reflection, JNI, and the full
class library. See [docs/ROADMAP.md](docs/ROADMAP.md).

## License

Licensed under the [Apache License, Version 2.0](./LICENSE) (SPDX:
`Apache-2.0`).
