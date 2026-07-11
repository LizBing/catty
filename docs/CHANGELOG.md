# Changelog

A running work log for catty. Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versions are project-local (no published releases yet).

The plan that governs this work lives in `plans/go-jvm-go-mvp-humming-bonbon.md`.

## [Unreleased]

### A1 — AOT emitter: bytecode → Go source (run a method natively)

The first executable proof of the performance thesis: lower a method's IR to Go
source and let the Go toolchain compile it, so `go build` is the optimizing
backend and the Go runtime stays the GC/scheduler.

### Added (A1)
- **`transpile/`** — `Emit(method, ir) (string, error)`: turns one method's
  `lowering.IR` into a Go function. The operand stack is eliminated into Go
  locals (`sK`); JVM locals are `lK` (the first `ArgSlotCount` are the params);
  control flow is `goto`/labels. Go's source rules are handled by construction:
  all slot/extra-local declarations precede any label (no goto-over-decl), labels
  appear only at branch/switch targets (no unused labels), and a trailing sink +
  `return 0` marks every slot used and satisfies the missing-return check.
- **`transpile/emit_test.go`** — lowers `Fibonacci.fib`, emits a `main` package,
  `go build`s + runs it, and asserts correctness and native-class speed.

### Performance — emitted `fib(35)` (the headline A1 result)
| Engine | Time |
|---|---|
| emitted Go (`go build`) | **~44 ms** |
| `java` (HotSpot JIT) | ~50 ms |
| `java -Xint` | ~600 ms |
| catty tree-walker / IR | ~4.5 s |

Emitted `fib` runs at native Go speed — ~100× faster than the interpreter and on
par with HotSpot's JIT. The thesis holds: bytecode → Go source → `go build`
delivers native-class performance, which no interpreter variant could.

### Scope (A1)
Int-only, static methods, the `fib` opcode subset (loads, int arithmetic,
shifts, `iinc`, conditional/unconditional branches, `ireturn`, `invokestatic`).
Non-int types, the object model (`new`/fields/`invokevirtual`), arrays, and
switches are A1.5/A2; integration into `cmd/jvm` (hot-method selection) is A4.

### A0 — bytecode IR + stack elimination (the AOT keystone)

The AOT keystone: a `lowering` pass that converts a method's stack-based
bytecode into a register-form IR (the operand stack eliminated into slot-indexed
virtual registers), plus an executor that runs the lowered form. This de-risks
the Phase 2 AOT transpiler by proving the stack-elimination transform is feasible
and correct *in isolation*, and builds the lowering infrastructure the emitter
will extend.

### Added
- **`opcode/`** — leaf package of JVMS opcode constants, shared by the
  interpreter and the lowering pass (extracted from `interpreter/` to break the
  would-be import cycle).
- **`lowering/`** — the IR types (`IR`, `IRInst`, `SwitchTable`), per-opcode
  instruction-length and slot-effect tables (field/invoke effects resolved from
  the constant-pool descriptor — no Loader needed), and the three-pass `Lower`:
  linear decode → forward depth dataflow over the CFG → vreg (Uses/Defs)
  assignment. Depth-only; no type inference, no SSA, no phis (position-stable
  slot indices are sound for execution because JVMS guarantees equal depth at
  every merge).
- **`interpreter.LoopIR`** — a second dispatch (`execIR`) that runs the lowered
  IR. Pure ops (arithmetic, constants, load/store, branches, comparisons,
  conversions, array access) address inputs/outputs by the precomputed
  Uses/Defs slot indices; complex ops (invoke, field, ldc, new, switches,
  returns, shuffles) reuse the tree-walker's helpers against a seeded
  operand-stack pointer.
- **`-ir` flag** on the launcher to select the IR executor.
- **`rtda.Frame`** indexed operand-stack accessors (`StackSlotNum`/`Ref`,
  `SetStackSlotNum`/`Ref`, `CopyStackSlot`, `SetStackTop`) for the IR executor.

### Changed
- `tests/run.sh` now runs each fixture through **three** engines — `java`, the
  tree-walker (`Loop`), and the IR executor (`-ir`) — and requires all three to
  be byte-identical. All 7 fixtures pass.
- `rtda.MethodDescriptor.argSlotCount` → exported `ArgSlots` (used by lowering).

### Performance — BenchFib (fib(35), ~29M recursive calls)
| Engine | Time |
|---|---|
| tree-walker (`Loop`) | 4.5 s |
| IR executor (`-ir`) | 4.8 s (~6% slower) |
| `java -Xint` | 0.6 s |

**Honest finding:** the IR executor is *not* faster than the tree-walker — the
per-instruction overhead (stack-pointer seeding, IR-instruction dereference,
slot-accessor calls) outweighs the predecode savings. Predecode inside a Go
interpreter does not pay off; this *confirms* that the performance arc must go
through the AOT transpiler (emit Go source, let `go build` optimize) rather than
a fancier interpreter. The IR executor's value is **validation**, not speed.
Recorded as ADR-0006.

### Validation gates (all green)
- `go test ./lowering/` — 22 methods across 8 classes lower with consistent
  dataflow; hand-checked `HelloWorld.<init>` Uses/Defs; all slot indices in range.
- `tests/run.sh` — 7/7 fixtures byte-identical across `java`, `Loop`, `LoopIR`.
- `go vet ./...` clean.


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
