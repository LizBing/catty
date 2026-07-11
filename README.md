# catty

An experimental JVM written in Go that **sits on top of the Go runtime**: it
reuses Go's garbage collector, scheduler, and allocator instead of implementing
its own. Java objects are Go objects (traced natively by Go's GC, with zero
write-barrier code); `java.lang.Thread` is planned to map to a goroutine.

catty has two execution paths:
- **Interpreter** (`catty -cp . MainClass`) — a switch-dispatched bytecode
  interpreter running a minimal-but-real Java subset, byte-identical to
  HotSpot on the test corpus.
- **AOT build** (`catty build -cp . MainClass`) — transpiles bytecode to Go
  source, compiles with `go build`, and produces a **standalone native binary**
  that runs at native speed. On `fib(35)` the emitted Go runs in ~44 ms — on
  par with HotSpot's JIT, ~100× faster than the interpreter.

## Documentation

| Document | What it covers |
|---|---|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | Premise, pipeline (interpreter + AOT), package responsibilities, data structures, traces, performance |
| [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) | Build/run/test, `catty build`, project layout, conventions, extension recipes |
| [docs/ROADMAP.md](docs/ROADMAP.md) | What's done, what's next (exceptions, concurrency, reflection, more) |
| [docs/CHANGELOG.md](docs/CHANGELOG.md) | Versioned work log (Phase 1 → A4) |
| [docs/adr/](docs/adr/) | Architecture Decision Records — the *why* behind each choice |

## What runs

### Interpreter (`catty -cp . MainClass`)

`./tests/run.sh` compiles each fixture with `javac` and diffs catty's stdout
against real `java` (through **three** engines: tree-walker, IR executor, and
`java`). All pass:

| Fixture | Exercises |
|---|---|
| HelloWorld | `ldc` String, `System.out.println`, int multiply |
| Fibonacci | recursion, `if`/`ireturn`, `iadd` |
| Factorial | `long` arithmetic, `long[]`, StringBuilder concat |
| ArraySum | `newarray`/`iaload`/`iastore`, loops |
| OOPDemo | `new`, `<init>`, instance fields, `invokevirtual` dispatch |
| StaticFields | `<clinit>` static initializers |
| SwitchDemo | `tableswitch` (dense) + `lookupswitch` (sparse) |
| EmptyMain | empty main — startup / smoke test |

The interpreter implements ~140 JVMS opcodes: constants, typed loads/stores,
full int/long/float/double arithmetic and conversions, shifts, comparisons,
all conditional branches and switches, stack manipulation, object/array
creation and access, field access, `invoke{virtual,special,static}`,
`checkcast`/`instanceof`, and returns. Core classes (`java.lang.Object`,
`String`, `StringBuilder`, `System`, `java.io.PrintStream`) are implemented
natively in Go rather than loaded from a JRE.

### AOT build (`catty build -cp . MainClass`)

`catty build` transpiles a whole program (all emittable methods via
reachability analysis), compiles with `go build`, and produces a standalone
binary. The emitter covers:

- **All primitive types**: int/long/float/double + ref + arrays
- **All control flow**: straight-line, loops (empty-stack merges), diamonds
  (phi via copy-insertion), switches (tableswitch/lookupswitch)
- **OOP**: `new`/`getfield`/`putfield`/`invokespecial` (constructors run
  interpreted via the bridge)
- **Invoke bridge**: `invokevirtual`/`special`/`static` route to native or
  interpreted targets via `catty/runtime` (the "world transition")
- **frem/drem**: float remainder via `runtime.FloatMod`/`DoubleMod`

Methods the emitter can't handle (unsupported opcodes, instance methods) are
served by the interpreter at runtime via the bridge — the tiered model.

`TestBuildProgram` validates HelloWorld + Fibonacci: both build to standalone
binaries and produce output byte-identical to `java`.

## Quickstart

```sh
go build -o catty ./cmd/jvm                    # build catty
./catty -cp <classpath> <MainClass>            # interpret
./catty -cp <classpath> -ir <MainClass>        # IR executor
./catty build -cp <classpath> [-o output] <MainClass>  # AOT build → native binary
./catty build -cp <classpath> -run <MainClass>        # AOT build + run

./tests/run.sh                                 # e2e: compile fixtures, diff vs java
go test ./...                                  # unit tests
```

Requires Go 1.22+ and a JDK (`javac`/`java`) on `PATH`. `catty build` also
requires running from the catty source tree (the emitted binary imports catty
packages). Set `CATTY_DEBUG=1` for a Go stack trace on a VM error.

## Performance (fib(35), ~29M recursive calls)

| Engine | Time | Relative |
|---|---|---|
| catty AOT (`catty build`) | ~44 ms | **native speed** |
| `java` (HotSpot JIT) | ~50 ms | baseline |
| `java -Xint` (HotSpot interpreter) | ~600 ms | ~14× AOT |
| catty interpreter (`Loop`) | ~4.5 s | ~100× AOT |
| catty IR executor (`-ir`) | ~4.8 s | slightly slower than `Loop` (ADR-0006) |

The AOT path reaches native speed — **~100× the interpreter and on par with
HotSpot's JIT**. The interpreter's ~7× gap to `java -Xint` is dispatch overhead
(switch vs computed goto; 16-byte slots; per-call frame allocation); ADR-0006
shows predecode doesn't close it — only AOT does.

## Vision: an experimental JRE platform

catty is evolving from an AOT transpiler into an **experimental JRE platform**
that compiles Java programs into Go programs — the final product is a native Go
binary running on Go's GC, scheduler, and network stack. Java's
Thread/synchronized/volatile/GC/IO are "dissolved" into Go's
goroutine/mutex/atomic/GC/netpoll at compile time.

Key architectural decisions (ADRs 0008–0013):

- **AOT-first**: interpreter is the dev tier; production runs AOT exclusively
  (no JIT warmup, no safepoints).
- **Thread = goroutine**: virtual threads from day one (no Loom needed).
- **Go memory model**: simpler than JMM; 99.9% compatible.
- **Escape analysis replaces generational GC**: Go's compiler stack-allocates
  non-escaping Java objects.
- **Direct Go runtime integration**: Java I/O compiles to Go netpoll
  (no JNI layer).
- **Hybrid class library**: ~50 critical classes native in Go; ~7000 loaded
  from the real JDK.

See [docs/ROADMAP.md](docs/ROADMAP.md) for the phased plan (R1–R6).

## License

Licensed under the [Apache License, Version 2.0](./LICENSE) (SPDX:
`Apache-2.0`).
