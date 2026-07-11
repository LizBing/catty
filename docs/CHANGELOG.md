# Changelog

A running work log for catty. Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versions are project-local (no published releases yet).

The plan that governs this work lives in `plans/go-jvm-go-mvp-humming-bonbon.md`.

## [Unreleased]

Nothing yet. In-progress work is tracked in [`ROADMAP.md`](./ROADMAP.md).

## [0.1.0] — 2026-07-11

The interpreter MVP (Phase 1). A single-threaded, switch-dispatched bytecode
interpreter running a minimal Java subset, end-to-end byte-identical to `java`
on the test corpus. Sits entirely on the Go runtime — no custom GC, scheduler,
or JIT.

### Added — packages
- **`classfile/`** — full `.class` parser (JVMS §4): `ClassReader` big-endian
  primitives, `ConstantPool` with all standard tags, `MemberInfo` for
  fields/methods, `CodeAttribute` + exception table, and a modified-UTF-8
  decoder (`mutf8.go`). `Parse([]byte) (*ClassFile, error)`.
- **`classpath/`** — directory, zip/jar, and composite entries; `Classpath.Parse`
  for the `-cp` option.
- **`classloader/`** — load/link/cache with three-way routing (array types,
  native core classes, user classes); implements `rtda.Loader`.
- **`rtda/`** — runtime data areas: `Slot`, `Frame`, `Thread`, `Class`,
  `Method`, `Field`, `Object`, `Array`; method-descriptor parsing; class
  construction from both `classfile.ClassFile` (`NewClass`) and scratch
  (`NewSyntheticClass`, `NewArrayClass`).
- **`interpreter/`** — `Loop` dispatch + 136 opcode `case`s covering constants,
  typed load/store, full int/long/float/double arithmetic + conversions, shifts,
  comparisons, all conditional branches, `tableswitch`/`lookupswitch`, stack
  manipulation, object/array creation and access, field access,
  `invoke{virtual,special,static}`, `checkcast`/`instanceof`, and returns.
- **`native/`** — synthetic core classes implemented in Go: `java.lang.Object`,
  `String`, `StringBuilder`, `java.lang.System`, `java.io.PrintStream`.
- **`cmd/jvm/`** — launcher: `-cp` flag, finds `main`, builds a `Thread`,
  triggers `<clinit>`, enters the interpreter.

### Added — verification
- **`tests/run.sh`** — e2e harness: compiles `tests/fixtures/*.java` with
  `javac -source 8`, runs each main through `java` and catty, diffs stdout.
- **Test corpus** (all passing): `HelloWorld`, `Fibonacci` (recursion),
  `Factorial` (long arithmetic + StringBuilder), `ArraySum` (arrays),
  `OOPDemo` (fields/constructors/virtual dispatch), `StaticFields` (`<clinit>`),
  `SwitchDemo` (`tableswitch` + `lookupswitch`), `BenchFib` (perf).
- **Unit tests**: `classfile/classfile_test.go` (parse a real compiled class),
  `rtda/frame_test.go` (two-slot long/double encoding, `PopRef` GC hygiene).

### Added — documentation
- `README.md` hub; `docs/ARCHITECTURE.md`, `docs/DEVELOPMENT.md`,
  `docs/ROADMAP.md`, `docs/CHANGELOG.md`, `docs/adr/` (5 ADRs).

### Performance — baseline on `fib(35)` (~29M recursive calls)
| Engine | Time | Relative |
|---|---|---|
| catty (interpreter) | 4.34 s | 1× |
| `java -Xint` (HotSpot interpreter) | 0.61 s | ~7× catty |
| `java` (HotSpot JIT) | 0.05 s | ~87× catty |

The ~7× gap to `java -Xint` is interpreter headroom (switch vs. computed-goto,
16 B slots vs. machine words, per-call frame allocation); the ~87× gap to JIT
is the JIT headroom only the planned AOT transpiler can close. Both are
discussed in `docs/ROADMAP.md` Themes B and A.

### Known limitations
- Single-threaded; `monitorenter`/`monitorexit` are nops (concurrency deferred).
- No exceptions/try-catch; exception tables are parsed but unused.
- No `invokedynamic` (compile fixtures with `-source 8`); `invokeinterface`
  panics.
- No reflection, JNI, `sun.misc.Unsafe`, or the real class library.
- `main`'s `args` is passed as `null` (programs that read `args` are unsupported).
- Java memory model not modeled (moot while single-threaded).

### Design decisions (recorded as ADRs)
- ADR-0001 — Reuse the Go runtime; write no GC, scheduler, or JIT.
- ADR-0002 — Switch dispatch (Go has no computed goto).
- ADR-0003 — Tagged 16-byte `Slot` (HotSpot stack-word model).
- ADR-0004 — Native core classes instead of shipping a JRE.
- ADR-0005 — Lazy `<clinit>` triggered at JVMS §5.5 init points via frame push.
