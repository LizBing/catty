# Changelog

A running work log for catty. Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
versions are project-local (no published releases yet).

The plan that governs this work lives in `plans/go-jvm-go-mvp-humming-bonbon.md`.

## [Unreleased]

### Strategic vision: catty as an experimental JRE platform

6 new ADRs (0008–0013) defining the architectural breakthrough: catty compiles
Java programs into Go programs that run natively on Go's runtime. Key decisions:
AOT-first (no JIT), Thread=goroutine (virtual threads from day one), Go memory
model (not JMM), escape analysis replaces generational GC, direct Go runtime
integration (no JNI layer), hybrid class library (~50 native + ~7000 from JDK).

Updated ROADMAP with phased plan (R1–R6).

### A2.5 — long / float / double type support (all primitives covered)

The emitter now handles all Java primitive types: long (int64), float (float32),
double (float64) — the last **type** gap. After this, the emitter covers int/
long/float/double/ref + arrays, all control flow (loops/diamonds), OOP, and the
native+interpreted invoke bridge.

### Added (A2.5)
- **`defTempCat2`** — allocates one Go temp for a category-2 def (long/double
  spanning 2 JVM slots → one int64/float64 Go value); binds both slots to it.
- **~50 new opcode cases** (float/long/double families): loads/stores, constants,
  arithmetic (add/sub/mul/div/rem for int64; add/sub/mul/div for float32/float64),
  shifts (long), comparisons (lcmp/fcmp/dcmp → int via 3-line if chain),
  conversions (i2l/l2i/i2f/f2i/... all 12), returns, array access, `ldc2_w`.
- **`buildLocalMap` + `totalParamSlots`** — maps JVM local slots to Go param
  names, accounting for category-2 (a double at slots 2-3 is one Go param "l1").
- **`invokestatic`** uses logical params (one Go arg per param, not per JVM slot).
- **`rtda.Object`** typed array accessors: `GetLongElement/SetLongElement`,
  `GetFloatElement/SetFloatElement`, `GetDoubleElement/SetDoubleElement`.
- `descToGo`/`goTypeOf`/`localTypes`/`storeLocalType` extended for J/F/D.

### Fixes
- `localName` switch dropped `Aload/Astore` _n variants during the A2.5 rewrite
  — restored (OOP's `astore_1` was resolving to `l0` instead of `l1`).

### Validation
- `TestEmitFact`: emitted `Factorial.fact(long)` compiles + runs natively;
  `fact(10) == 3628800` (long recursive — lload/lconst/lcmp/ifgt/lsub/lmul/
  invokestatic/lreturn).
- `TestEmitFloatDouble`: `fadd(1.5, 2.5) == 4`, `dmul(1.5, 2.5) == 3.75`
  (float cat1 + double cat2 arithmetic).
- fib/first/HelloWorld/sum/OOP/max still pass; `go vet` clean; e2e 8/8.

### Scope (A2.5)
All primitive types + refs. `frem`/`drem` (Go has no float `%`) and cat-2 merge
temps (long/double phi) deferred. `cmd/jvm` integration is A4.

### A2.4 — Diamonds / ternary (phi via copy-insertion)

The last control-flow gap: merges that leave a value on the operand stack
(diamonds — `cond ? x : y`, if/else yielding a value). A2.3 handled empty-stack
merges (loops) via mutable locals; A2.4 handles non-empty-stack merges by
inserting the equivalent of SSA **phi nodes** as **copy-insertion**: a per-slot
merge temp assigned at each predecessor edge (branch-to-merge and
fall-through-into-merge), read after the join.

### Changed (A2.4)
- **`transpile/emit.go`** — `cfgAnalysis` (merges + fall-through edges);
  `allocMergeTemps` (one temp per stack slot at each non-empty-stack merge, typed
  by `InTypes`; refuses long/float/double merge slots — the only remaining gate);
  copy-insertion at branch edges and fall-through; `slotTemp`/`slotType` reset to
  the merge temps after the join. The A2.3 blanket refusal of non-empty-stack
  merges is removed.
- `TestEmitDiamondGate` (which expected `max` to error) → `TestEmitMax` (emits +
  runs `a > b ? a : b`; both orderings return `7`).

### Validation
- `TestEmitMax`: emitted `max` compiles + runs natively; the merge temp is
  assigned on both branches and returned at the join — correct phi. fib/first/
  HelloWorld/sum/OOP still pass; `go vet` clean; e2e 8/8.

### Scope (A2.4)
int/ref diamonds. long/float/double merge slots, switch joins (switches aren't
emitted yet), and `cmd/jvm` integration (A4) are later. Control-flow is now
complete (straight-line / loops / diamonds).

### A2.2b — OOP: new / fields / invokespecial (+ interpreted-target bridge)

The AOT path now handles objects. The bridge extended from native-only to **native
or interpreted** targets (constructors are bytecode via `invokespecial`; user
methods via `invokevirtual` are interpreted), and the emitter learned `new`/`dup`/
`invokespecial`/`getfield`/`putfield`. Milestone: an OOP program — `new Box();
b.v = 21; println(b.v + b.doubled())` — transpiled and run natively, printing `63`.

### Added (A2.2b)
- **Interpreted-target bridge**: `rtda.Thread` gains a bridge-return slot; the 5
  return helpers write it when the stack is empty (the bridge's outermost callee);
  `interpreter.RunMethod` runs an interpreted method and captures its return.
- **`catty/runtime`**: `InvokeSpecial`; `InvokeVirtual`/`InvokeSpecial` dispatch
  native→`runNative`, interpreted→`interpreter.RunMethod`; `NewObject` (allocate,
  no init). 
- **`transpile/emit.go`**: `Emit(method, ir, loader)` (field-offset resolution);
  `new`/`dup`/`invokespecial`/`getfield`/`putfield`; extra-locals typed from their
  store opcodes (`astore`→`*rtda.Object`, `istore`→`int32`); `slotType` tracking
  for `dup`.
- `tests/fixtures/OOPAot.java` — `Box{int v; int doubled(){return v+v;}}` + main.

### Fixes
- `getfield`/`iadd`-class bug: read use-temps before allocating the def-temp (a
  def often reuses an operand's slot).

### Validation
- `TestEmitOOP`: emitted `OOPAot.main` compiles + runs natively, prints `63`
  (== interpreter/`java`). Exercises new + interpreted `<init>`, putfield,
  getfield, interpreted `invokevirtual doubled`, native `println`. fib/first/
  HelloWorld/sum still pass; `go vet` clean; e2e 8/8.

### Scope (A2.2b)
Straight-line OOP. Diamonds/ternary (phi), `long`/`float`/`double` in the emitter,
and interpreted `long`/`double` returns through the bridge are later; `cmd/jvm`
integration is A4.

### A2.3 — Loops (empty-stack merges)

The AOT path now handles **loops** — the dominant pattern in compute-intensive
Java. The key realization: loops need **no phi insertion**, because loop state
(counters, accumulators) lives in **locals**, which the emitter already maps to
mutable Go vars; the operand stack is empty at loop heads, so no operand-stack
value crosses the merge either. Empty-stack merges therefore work with the
existing fresh-per-def + mutable-locals machinery once the gate lets them through.

### Changed (A2.3)
- **`transpile/emit.go`** — the merge gate now refuses only merges with a
  **non-empty operand stack** (a value crosses the merge → phi needed, e.g.
  diamonds `cond ? x : y`), checked via `IRInst.InTypes` length at merge pcs.
  Empty-stack merges (loops) are allowed.
- **Fix**: `iinc`'s local index is in `IncIndex`, not `Index` — `localName` now
  handles `iinc` (previously emitted the wrong local; latent, unexercised until a
  loop used `iinc`).
- **`runtime.NewIntArray`** — builds a Java `int[]` for transpiled test programs.
- `ArrayOps.max` (ternary) added for the diamond-gate test.

### Validation
- `TestEmitSum`: emitted `ArrayOps.sum` (for-loop over an array) compiles + runs
  natively, returns `15` for `[1,2,3,4,5]`.
- `TestEmitDiamondGate`: `max(a,b)` (ternary) → `Emit` errors (non-empty-stack
  merge). fib/first/HelloWorld still pass; `go vet` clean; e2e 8/8.

### Scope (A2.3)
Empty-stack merges (loops). Diamonds / non-empty-stack merges (phi via
copy-insertion), `new`/fields/OOP, interpreted-target bridge, and `cmd/jvm`
integration are later.

### A2.2 — The invoke bridge: HelloWorld via AOT (first full-program native run)

The AOT path can now serve a real program: emitted Go calls back into catty's
runtime (the "world transition" of ADR-0007) to resolve classes/fields/methods
and run native code. Milestone: **HelloWorld transpiled and run natively**,
output byte-identical to `java`.

### Added (A2.2)
- **`catty/runtime`** — the bridge package: `Bootstrap` (load main class + deps,
  run `<clinit>`), `GetStatic`, `InvokeVirtual` (dynamic dispatch on the
  receiver), `NewString`, `runNative`/`popReturn`. Targets are resolved by
  (class, name, descriptor) at run time via the loader — no method registry.
- **`rtda.RefSlot`/`IntSlot`** — slot constructors (fields are unexported; emitted
  code boxes call args with these).
- **`transpile/emit.go`** — `getstatic` → `runtime.GetStatic`, `ldc`-String →
  `runtime.NewString`, `invokevirtual` → `runtime.InvokeVirtual` (typed arg slots
  + return extract); void `return`.

### Validation
- `TestEmitHelloWorld`: emitted `HelloWorld.main` (getstatic / ldc-String /
  invokevirtual-println / int math) compiles + runs natively, printing
  `Hello, World!\n42\n` (== `java`). fib still runs natively; `go vet` clean;
  e2e 8/8; all tests green.

### Scope (A2.2)
The bridge runs **native** targets (covers `println`). Interpreted targets via
the bridge (catcher frame), `new`/`invokespecial`/fields (OOP), and merges/loops
(phi) are later. `invokestatic` stays a direct Go call (emitted targets).

### A2.1 — Fresh-per-def, type-aware emitter (refs + arrays)

The emitter rewrites from A1's position-stable slots (all `int32`) to
**fresh-per-def temps**: each def is a new typed Go local; uses reference the
defining temp. This resolves the slot-type-reuse that broke ref methods under A1
(`aload; arraylength` writes a `*Object` then an `int32` into the *same* slot).
Refs, arrays, and typed params/returns now emit correctly — the foundation for
transpiling real (object-using) programs.

### Changed (A2.1)
- **`transpile/emit.go`** rewritten: typed params + return (from the descriptor),
  fresh-per-def temps declared at the top by type (goto never crosses a `var`),
  trailing sink + `return <zero>` terminator, and a **merge-free gate** (methods
  with control-flow merges — loops/diamonds, which need phi insertion — return an
  error instead of wrong code; that's A2.3).
- New emitted opcodes: `aload`/`astore`/`areturn`, `aconst_null`, `arraylength`,
  `iaload`/`iastore`/`aaload`/`aastore` (typed array access); `invokestatic`
  emits a direct Go call to the mangled (emitted) target. Emitted code imports
  `catty/rtda` for `*rtda.Object` + array accessors.
- `tests/fixtures/ArrayOps.java` — `first(int[])I` (merge-free, ref+array) and
  `sum(int[])I` (loop, for the merge gate).

### Scope (A2.1)
Static, merge-free methods; int + ref(+array) types. Long/float/double, fields
(`getfield`/`putfield`), the invoke bridge (native/interpreted targets), and
merges/loops are later (A2.1b / A2.2 / A2.3).

### Validation
- `fib` still executes natively (`fib(35)==9227465`, ~43 ms) — the rewrite
  didn't break int code.
- `first` emits correct ref+array code **and it compiles** (compile-checked; it
  isn't executed because running it needs the runtime bridge — A2.2).
- `sum` (loop) → `Emit` returns a "has merges" error.
- `go vet` clean; e2e 8/8; all unit tests green.

### A1.5 — Type tracking in the lowering (the A2 enabler)

The lowering now carries, per IR instruction, the operand-stack slot types
(`IRInst.InTypes`) — the prerequisite for A2's emitter to declare correct Go
types (int/long/float/double/ref) instead of int-only. Lowering-only; no
runtime/emitter change, no visible speedup yet.

### Added (A1.5)
- **`classfile/stackmap.go`** — parses the `StackMapTable` attribute (JVMS
  §4.7.4): delta-encoded frames reconstructed via `Reconstruct(initialLocals)`,
  exposed through `CodeAttribute.StackMapTable()`.
- **`lowering/types.go`** — `SlotType` (Int/Long/Float/Double/Ref/Top), a
  `typeDataflow` linear pass that propagates operand-stack types per opcode,
  pinning merges at `StackMapTable` frames. Loads are opcode-derived (the verifier
  guarantees them), so only the operand stack is tracked, not locals.

### Design notes
- **No type lattice / merge logic**: the `StackMapTable` gives the exact merged
  frame at every branch target (Java 6+), so catty only propagates types within
  basic blocks and resets at frames.
- The pass **always runs** (linear propagation); branch-free methods without a
  `StackMapTable` (e.g. `HelloWorld.main`) still get correct `InTypes`. Methods
  with branches but no table (pre-Java-6) may get imprecise types — acceptable;
  A2 will just skip AOT for them.

### Validation
- `classfile` test: `fib`'s table parses to one SAME frame at pc 7 (vs `javap`).
- `lowering` test: `fib` is all-`TypeInt`; `HelloWorld.main` shows `TypeRef`
  (System.out/String); `Factorial.fact` shows `TypeLong`. Existing depth/Uses/Defs
  tests unchanged; `go vet` clean; e2e 8/8.

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
