# Catty

A Java language runtime parasitic on the Go runtime.

**GC · scheduler · netpoller: borrowed. Bytecode → Go source: ours.**

No JVM. No JIT. No warmup. Single static binary.

## What It Does

Catty executes compiled Java bytecode on top of the Go runtime. Objects, threads and GC are **parasitic** — they live directly on Go's allocator, scheduler and garbage collector. The execution engine is **AOT-first**: bytecode is translated to Go source at build time and compiled into a single static binary. An interpreter handles dynamic class loading and reflection as fallback.

**JDK 8 baseline.** IR version-agnostic with an upgrade path to modern JDKs.

## Performance

Measured on darwin/arm64, go1.26.5, JDK 25 HotSpot. See [docs/research/](docs/research/) for methodology.

| Benchmark | Catty AOT | Catty interpreter | vs interpreter |
|---|---:|---:|---|
| Arithmetic loop | 38ms | 90ms | **2.4×** |
| Virtual call (inline cache) | 68ms | 202ms | **2.9×** |
| HashMap get/put | 629ms | 1566ms | **2.5×** |
| Cold start | <10ms | ~15ms | **>92% faster than JDK** |

Sustained load (embedded, 16 workers): **65k req/s**, p99 = 0.96ms.

## Status

```
M0 ████████ Interpreter + HelloWorld
M1 ██████████ Monitor rewrite, classloader, structural verifier
M2 ██████████ Threads/SOE/Class metadata/net+echo/dataflow verifier
M3 ██████████ AOT emitter closed loop — all fixtures byte-identical to JVM
M4 ██████░░░░ Dispatch optimization + allocation reduction
```

**Real-world validation:** gson 2.8.9 (214 classes) loads from JAR, serializes and deserializes POJOs with output byte-identical to HotSpot. Reflection minimal surface (Class/Field/Method/Constructor with inherited traversal) enables framework-shaped code.

## Quick Start

### Build

```sh
go build ./cmd/cattaot    # AOT-first runner
go build ./cmd/catty      # pure interpreter runner
```

### Run Java

```sh
javac --release 8 -d classes src/HelloWorld.java
./cattaot -cp classes run HelloWorld
```

### Embed in Go

```go
k := kernel.New(kernel.Options{Stdout: os.Stdout})
th := vm.New(k)

// Bind Go functions as Java static methods
interop.Bind(k, interop.Spec{
    Class: "com/app/Bridge",
    Funcs: map[string]any{
        "md5Hex": func(s string) string { /* ... */ return "" },
    },
})

loader := kernel.NewClassPathLoader(k, []string{"app/classes"})
k.SetResolver(loader.Load)
cls, _ := loader.Load("com/app/Main")
main, _ := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
argsArr, _ := k.NewArray("Ljava/lang/String;", 0)
th.Call(main, nil, []kernel.Value{argsArr})
```

See [cmd/embeddemo](cmd/embeddemo/main.go) for a working example.

## Architecture

```
┌─────────────────────────────────────────────┐
│           Your Go Application               │
│  ┌───────────────────────────────────────┐  │
│  │         Catty Runtime                 │  │
│  │  ┌──────────┐  ┌───────────────┐      │  │
│  │  │ AOT Code │  │  Interpreter   │      │  │
│  │  │ (gen.go) │  │  (fallback)    │      │  │
│  │  └────┬─────┘  └──────┬────────┘      │  │
│  │       │    InvokeAs    │              │  │
│  │       ▼               ▼              │  │
│  │  ┌───────────────────────────┐       │  │
│  │  │     Kernel (objects,      │       │  │
│  │  │     threads, monitors)    │       │  │
│  │  └────────────┬──────────────┘       │  │
│  └───────────────┼──────────────────────┘  │
│                  ▼                          │
│         Go Runtime (GC, sched, net)         │
└─────────────────────────────────────────────┘
```

Both engines share the same object model (`kernel.Instance`), the same dispatch entry point (`Kernel.InvokeAs`), and the same stack-trace infrastructure. Unemitted classes transparently fall back to the interpreter.

## Key Design Decisions

| Decision | ADR | Summary |
|---|---|---|
| Exception mechanism | [ADR-0009](docs/decisions/ADR-0009.md) | Flag-return (not panic/recover) — measured 0.6% crossover frequency |
| Object representation | [ADR-0010](docs/decisions/ADR-0010.md) | Unified kernel.Instance; mixed representation deferred to v2 |
| Opcode stack effects | N/A | Single source of truth in `classfile.StackEffect` — total classification over 256 opcodes |

## Known Limitations

- No annotations metadata retention (getAnnotation returns null)
- No ParameterizedType / generic type resolution (TypeToken won't work)
- getMethods() excludes interface default methods
- Access control not enforced (private fields settable without setAccessible)
- Blocking inside interop-bound Go functions is not interrupt-wakeable
- No JAR signing, no security manager, no JPMS modules

Full list: [docs/specifications/deviation-ledger.md](docs/specifications/deviation-ledger.md)

## Development

```sh
make check    # mechanical gate: docs structure + go vet/build/test
make fuzz     # classfile parser fuzzing (30s bounded session)
go test -race ./...  # full race-detector suite
```

## Project Layout

```
cmd/
  catty/       pure interpreter CLI
  cattaot/     AOT-first CLI (emitted bodies take precedence)
  genemit/     bytecode→Go source translator
  embeddemo/   embedding example (Go host + Java rules)
internal/
  classfile/   .class file parser + opcode stack-effect table
  verify/      structural + dataflow verifier
  kernel/      object model, threads, monitors, natives, reflection
  vm/          bytecode interpreter
  emitter/     AOT source generator
  genrt/       runtime bridge between generated code and kernel
  gen/         generated Go source (DO NOT EDIT)
  interop/     Bind API for embedders
docs/
  vision.md    product intent and selling-point ranking
  status.md    capability matrix (refreshed at milestone close)
  research/    benchmark and probe reports (R-0001…R-0009)
  decisions/   architecture decision records (ADR-0001…ADR-0011)
  plans/       execution plans (active + completed)
  debt/        technical debt register
```

## License

Apache-2.0. Third-party provenance policy: [docs/quality/architecture-rules.md](docs/quality/architecture-rules.md).
