# Architecture

This document describes catty's design: the premise it is built on, the
execution pipeline, the responsibilities of each package, how a program flows
through the system, and the key data structures. For *why* the decisions were
made, see [`adr/`](./adr/). For *how to change the code*, see
[`DEVELOPMENT.md`](./DEVELOPMENT.md).

## 1. The premise

catty is a JVM that **sits on top of the Go runtime**. The three heaviest JVM
subsystems are deliberately *not* implemented:

| JVM subsystem | catty's approach |
|---|---|
| Garbage collector | Java objects are Go heap allocations; Go's GC traces them natively. No write barriers, no mark/sweep code. |
| Thread scheduler | MVP is single-threaded. A future concurrency arc may use goroutines as carriers, but Java Thread identity, lifecycle, monitors, interruption, and memory semantics require Accepted decisions. |
| JIT compiler | None directly. Instead, a bytecode→Go-source AOT transpiler (`catty build`) hands optimization to the Go compiler — `go build` IS the optimizing backend. |

This trade is the whole point of the project: by not writing a GC, scheduler, or
JIT, the implementation is small (~7500 LOC) and the MVP floor is bounded — the
work collapses to "a bytecode interpreter plus a class loader" (Phase 1) and "a
lowering pass + Go-source emitter" (A0–A4). The AOT path reaches native speed
(`fib(35)` in ~44 ms, on par with HotSpot JIT); the interpreter is the fallback
for anything the emitter can't handle (see §7).

## 2. The pipeline

catty has two execution paths — **interpret** and **AOT build** — sharing the
same class-loading front end:

### Interpret path (`catty -cp . MainClass`)

```
   .class bytes
        │
        ▼
 ┌──────────────┐   ┌─────────────┐   ┌───────────────┐   ┌───────┐
 │  classpath   │──▶│  classfile  │──▶│  classloader  │──▶│ rtda  │
 │ locate bytes │   │ parse bytes │   │ link + cache  │   │ Class │
 └──────────────┘   └─────────────┘   └───────────────┘   └───┬───┘
                                                             │
                                       ┌─────────────────────┘
                                       ▼
                               ┌────────────────┐     ┌────────┐
                               │  interpreter   │◀───▶│ native │
                               │ switch dispatch│     │ core   │
                               └────────────────┘     └────────┘
                                       ▲
                                       │
                                 ┌─────────┐
                                 │ cmd/jvm │
                                 └─────────┘
```

### AOT build path (`catty build -cp . MainClass`)

```
   .class bytes
        │
        ▼
 ┌──────────────┐   ┌─────────────┐   ┌───────────────┐   ┌───────┐
 │  classpath   │──▶│  classfile  │──▶│  classloader  │──▶│ rtda  │
 └──────────────┘   └─────────────┘   └───────────────┘   └───┬───┘
                       │ (reachability)                        │
                       ▼                                       ▼
                 ┌───────────┐   ┌──────────┐   ┌────────────────┐
                 │ lowering  │──▶│ transpile │──▶│  go build       │──▶ native binary
                 │ IR + types│   │ Emit Go   │   │ (embedded runtime│
                 └───────────┘   └──────────┘   │ + interpreter)  │
                                                └────────────────┘
```

The emitted binary embeds the entire catty runtime (`classloader` + `interpreter`
+ `native` + `runtime` bridge). At run time: AOT'd methods call direct Go /
the bridge; un-AOT'd methods are served by the interpreter via the bridge — the
current implementation of the **multi-engine model** (ADR-0016).

### Stage responsibilities

1. **`classpath/`** — find `<name>.class` on disk (directories, jars, zips).
2. **`classfile/`** — decode raw bytes into `*classfile.ClassFile` per JVMS §4
   (constant pool, members, Code attribute, StackMapTable, modified-UTF-8).
3. **`classloader/`** — load + link + cache. Converts a `classfile.ClassFile`
   into a runtime `rtda.Class`: resolves supers, computes field offsets, builds
   `rtda.Method`s. Routes array/core/user classes.
4. **`rtda/`** — runtime data areas: `Slot`, `Frame`, `Thread`, `Class`,
   `Method`, `Field`, `Object`, `Array`. Pure data + class construction.
5. **`opcode/`** — leaf package of JVMS opcode constants, shared by the
   interpreter and the lowering pass (no import cycle).
6. **`lowering/`** — converts a method's bytecode into a register-form IR:
   decode → depth dataflow → vreg (Uses/Defs) assignment + type tracking
   (via StackMapTable). The stack is eliminated into slot-indexed vregs.
7. **`interpreter/`** — switch-dispatch loop over bytecode (the interpreter)
   and over the lowered IR (`LoopIR`, the validation tier). Also the
   interpreted-target bridge (`RunMethod`) for the AOT path.
8. **`transpile/`** — the AOT emitter: `Emit` (one method → Go source) and
   `BuildProgram` (whole-program: reachability + pass1/2 + assemble).
9. **`runtime/`** — the AOT bridge: `Bootstrap`, `GetStatic`,
   `InvokeVirtual`/`Special`/`Static`, `NewObject`/`NewString`, `FloatMod`/
   `DoubleMod`. Resolves targets by (class, name, desc) at run time.
10. **`native/`** — core JDK classes (`Object`, `String`, `StringBuilder`,
    `System`, `PrintStream`) implemented in Go rather than loaded from a JRE.
11. **`cmd/jvm/`** — launcher: interpret (`catty -cp . Main`) or build
    (`catty build -cp . Main`).

## 3. Package dependency graph

```
classfile          (stdlib only)
classpath          (stdlib only)
opcode             (stdlib only)
rtda        ──▶ classfile
native      ──▶ rtda
classloader ──▶ classfile, classpath, rtda, native
lowering   ──▶ classfile, opcode, rtda
interpreter ──▶ classfile, opcode, lowering, rtda
runtime    ──▶ classloader, classpath, interpreter, rtda
transpile ──▶ classfile, lowering, opcode, rtda
launch     ──▶ classloader, classpath, interpreter, rtda, transpile
cmd/jvm    ──▶ launch                (CLI only: arg parsing + java.base detection)
```

There are **no import cycles**. The one place a cycle would naturally appear —
`rtda.Class` needing to load other classes (superclass, array components) while
`classloader` builds `rtda.Class` — is broken by the `rtda.Loader` interface:

- `rtda` declares `type Loader interface { LoadClass(name string) *Class }`.
- `classloader.ClassLoader` implements it.
- `rtda.NewClass` / `rtda.NewArrayClass` take a `Loader` and call back into it
  for supers and components.
- The interpreter resolves classes through `thread.Loader()` (typed
  `rtda.Loader`), so it never imports `classloader`.

This keeps `rtda` — the most depended-upon package — free of upward imports.

## 4. Key data structures

### `rtda.Slot` — the atomic cell  *(slot.go)*

```go
type Slot struct {
    num int32    // byte/char/short/int/boolean, float-as-bits, returnAddress
    ref *Object  // object or array reference; nil == Java null
}
```

Both the operand stack and local variable arrays are `[]Slot`. A category-2
value (long/double) occupies **two** consecutive slots: high word, then low
word. This is the HotSpot "stack word" model. `PopRef` nils the freed slot so
the GC can reclaim the object — a subtle but real correctness property.

### `rtda.Frame` — one stack frame  *(frame.go)*

```go
type Frame struct {
    lower    *Frame   // caller frame (thread keeps frames as a slice)
    thread   *Thread
    method   *Method
    code     []byte   // cached method.code, the dispatch loop reads here
    locals   []Slot
    stack    []Slot
    stackTop int      // operand-stack size = next free index
    pc       int      // index into code of the next instruction
}
```

`Frame` owns both the locals and the operand stack as `[]Slot`, and exposes
typed push/pop (`PushInt`/`PopInt`/`PushLong`/…) and typed local access
(`GetInt`/`SetRef`/…) so instruction handlers stay decoupled from the slot
layout. Operand decoders (`ReadUint8`/`ReadUint16`/`ReadInt16`/`ReadInt32`)
advance `pc` past the bytes they consume.

### `rtda.Class` / `Method` / `Field` / `Object`

- `Class` holds metadata + a `methodTable map[string]*Method` (key = name +
  descriptor), instance/static field slices, static var storage, and array
  metadata (`isArray`, `componentClass`, `componentKind`). It also carries a
  reference to its `*classfile.ConstantPool` for run-time resolution of
  `ldc`/fieldref/methodref.
- `Method` carries bytecode (`code []byte`), `maxStack`/`maxLocals`, the parsed
  exception table, and — if native — a `func(*Frame)`.
- `Object` is `{class *Class, fields []Slot, extra any}`. `fields` backs both
  instance fields (indexed by `Field.SlotID()`) **and** array elements (1 slot
  per category-1 element, 2 for long[]/double[]). `extra` is a native payload —
  e.g. a synthetic `java.lang.String`'s immutable `rtda.StringValue` UTF-16
  code-unit backing, or a `PrintStream`'s `io.Writer`.

### `rtda.Thread` — execution context  *(thread.go)*

```go
type Thread struct {
    stack  []*Frame
    loader Loader
}
```

The MVP runs a single `Thread`. Pushing a frame per call and popping on return
gives the JVM stack. ADR-0018 permits goroutines as future execution carriers but
does not equate them with Java Thread objects or define the concurrency contract.

## 5. Class loading, in detail

`classloader.ClassLoader.LoadClass(name)` is the single entry point used both at
startup and at run time (via `thread.Loader()`). It memoizes into a `cache`
map. For a cache miss it walks a **chain of `ClassProvider`s** — the first that
returns a non-nil class wins. `New(cp)` wires the standard chain; the order is
the loading strategy:

| # | Provider | Serves |
|---|---|---|
| 1 | `ArrayProvider` | `[...` array types → `rtda.NewArrayClass` |
| 2 | `BootstrapProvider` | the 6 irreducible bootstrap classes (Object, String, Class, System, Thread, Throwable) — always synthetic, never from a class file (they carry catty's Go↔Java bridge payloads in `Extra()`) |
| 3 | `SyntheticProvider` | the rest of the synthetic registry (StringBuilder, PrintStream, exception subclasses, Comparable, …) — synthetic fallback when no real JDK is present |
| 4 | `ClasspathProvider` | real `.class` files via `classpath.ReadClass` → `classfile.Parse` → `rtda.NewClass` |

So with java.base on the classpath, a user class like `ArrayList` misses
providers 1–3 and is served real bytecode by `ClasspathProvider`; the bootstrap
classes are still synthetic (provider 2 wins before provider 4). The strategy is
**configured by the order of providers**, not hardcoded in `LoadClass` — custom
strategies use `NewCustom(providers...)`.

The synthetic registry itself is a `map[string]builderFunc` populated by
`init()` blocks across `native/*.go` via `registerSynthetic(name, fn)`. The
bootstrap-class boundary (`native.BootstrapClasses`) is the current R1 set that
`BootstrapProvider` claims; ADR-0022 makes its capability boundary revisable.

After a class is provided, `resolveNativeMethods` patches any `ACC_NATIVE`
methods whose `(class, name, desc)` is in the global `RegisterNative` table
(`native/system.go` `init()`) with a Go implementation; unregistered natives
keep a zero-return stub (graceful `NoSuchMethodError`-free degradation).

`rtda.NewClass` is where linking happens: it loads the superclass and interfaces
recursively, computes instance field offsets starting from the superclass's
count (so subclasses inherit offsets), allocates static var slots, and builds a
`Method` per declared method. Class **initialization** (`<clinit>`) is not done
here — see ADR-0021.

### 5a. java.base auto-detection (CLI layer)

`cmd/jvm/main.go` (not the runtime) detects a JDK and prepends java.base to the
user's `-cp`. Detection order: `$CATTY_BOOT` → `$JAVA_HOME` →
`/usr/libexec/java_home` (macOS) → well-known install paths. catty expects a
**pre-extracted** java.base directory (produced by the JDK's own `jimage
extract`); it deliberately ships no runtime jimage parser to stay lean. Pass
`--no-boot` to force pure-synthetic mode.

The runtime itself (`launch.Interpret`) knows nothing about java.base — it
receives only a fully-formed classpath string. The "dissolve into Go runtime"
principle applies: bootstrap/environment concerns live in the CLI, not in the
runtime packages.

## 6. An execution trace — `System.out.println("Hello, World!")`

1. `cmd/jvm` loads `HelloWorld` (caches it; `rtda.NewClass` recursively loads the
   native `java.lang.Object`, `System`, `PrintStream`, `String`).
2. A `Thread` is created; a frame for `main` is pushed with `locals[0] = null`
   (the `args` array); `InitClass` runs `HelloWorld.<clinit>` if present (pushed
   on top, so it runs before `main`).
3. `interpreter.Loop` enters its dispatch cycle. `main`'s bytecode is roughly:
   ```
   getstatic   System.out        ; push the PrintStream for stdout
   ldc         "Hello, World!"   ; push a String object (extra = UTF-16 units)
   invokevirtual PrintStream.println(String)
   ...
   ```
4. `getstatic` resolves the fieldref via the owner class's constant pool, loads
   `java.lang.System`, looks up `out`, and pushes the `PrintStream` object that
   `native.buildSystemClass` stored (its `extra` is `os.Stdout`).
5. `ldc` decodes the classfile MUTF-8 literal losslessly to UTF-16 code units
   and builds a `String` object whose `extra` holds an immutable `StringValue`.
6. `invokevirtual` peeks the receiver (the `PrintStream`), does dynamic dispatch
   on `receiver.Class()`, finds `println(String)` (native), and calls
   `invokeNative`: it copies args into a throwaway frame, runs
   `native.printlnString`, which does `fmt.Fprintln(os.Stdout, value)`.
7. `return` pops `main`'s frame; the stack drains; `Loop` exits.

## 7. Where the performance comes from (and where it doesn't)

The current interpreter is a single dense `switch` on the opcode byte; this is
an implementation choice governed by ADR-0024.
Measured on `fib(35)` (~29M recursive calls):

| Engine | Time | Relative |
|---|---|---|
| catty (interpreter) | 4.34 s | 1× |
| `java -Xint` (HotSpot interpreter) | 0.61 s | ~7× catty |
| `java` (HotSpot JIT) | 0.05 s | ~87× catty |

The ~7× gap to HotSpot's interpreter is **interpreter headroom**: switch
dispatch vs. computed-goto threaded dispatch, 16-byte `Slot` vs. machine words,
and per-call `Frame` allocation. All tunable. The ~87× gap to the JIT is
**JIT headroom** — only the Phase 2 AOT transpiler can close it, by lowering
bytecode to Go source and letting `go build` optimize.

## 8. Lowering & IR (A0)

The `lowering` package converts a method's stack-based bytecode into a
register-form IR — the operand stack eliminated into slot-indexed virtual
registers. It exists to **de-risk the AOT transpiler** (ROADMAP Theme A): it
proves the hardest transform — static stack elimination — is feasible and
correct, in isolation, before the emitter is built.

The pipeline (`lowering.Lower`):

1. **Decode** — walk instructions from pc 0, predecoding operands and resolving
   branch/switch targets to absolute pcs (no operand parsing at run time).
2. **Depth dataflow** — a forward worklist over the control-flow edges computes
   the operand-stack depth (in *slots*) on entry to each instruction. Each
   opcode's slot effect is statically known; field/invoke effects are read from
   the constant-pool descriptor (so lowering needs no `Loader` — it's a pure
   function of the method bytecode + its cp).
3. **vreg assignment** — turn each instruction's (entry depth, slot effect) into
   concrete `Uses`/`Defs` slot indices.

A0 is **depth-only** for vreg assignment and needs **no SSA or phis**: vregs are
position-stable slot indices, and JVMS guarantees equal stack depth at every
merge point, so a single path's definitions are always the live ones at
execution time. **A1.5 adds type tracking** on top: a linear pass propagates
operand-stack slot types (`IRInst.InTypes`), pinning merges at the
`StackMapTable`'s frames — no type-lattice/merge logic needed. Loads are
opcode-derived, so only the stack is tracked. This is what A2's emitter consumes
to declare the right Go type per slot.

`interpreter.LoopIR` runs the lowered form (opt-in via `-ir`). Every instruction
seeds the operand-stack pointer from the IR's known depth; pure ops read/write
the precomputed `Uses`/`Defs` slots, while complex ops reuse the tree-walker's
helpers. `tests/run.sh` requires `java`, `Loop`, and `LoopIR` to be
byte-identical — the equivalence gate that proves the lowering is
semantics-preserving.

**The IR executor is not faster than the tree-walker** (~6% slower on
`BenchFib`): predecode savings are smaller than the IR dispatch overhead in Go.
This is expected and fine — the IR's job is validation and as the emitter's
input; ADR-0024 keeps its performance role evidence-driven while ADR-0016
defines AOT as the primary product path.

## 9. AOT transpiler (A1)

`transpile.Emit` turns a method's `lowering.IR` into Go source — the executable
proof that bytecode → Go source → `go build` reaches native-class speed. Each
operand-stack slot becomes a Go local `sK`; each JVM local is `lK` (the first
`ArgSlotCount` are the function's parameters); bytecode control flow becomes
`goto`/labels. The Go toolchain compiles it, with the Go runtime as GC/scheduler.

The Go-source rules shape the emitter by construction: all slot/extra-local
declarations precede any label (so `goto` never crosses a `var`); a `pcNN:`
label appears only at branch/switch targets (no unused labels); a trailing
`_ = sK` sink plus `return 0` marks every slot used and satisfies the
missing-return check.

**Result on `fib(35)`:** emitted Go runs in ~44 ms — native speed, ~100× the
interpreter and on par with HotSpot's JIT (see the A1 changelog entry). A1 is
scoped to int-only static methods and the `fib` opcode subset; non-int types,
the object model, and runtime integration are A1.5/A2/A4.

## 10. What catty does *not* model (yet)

Kept out of scope deliberately; each is a documented future work item in
[`ROADMAP.md`](./ROADMAP.md):

- **Concurrency** — single-threaded; no `synchronized`/`wait`/`notify`
  (`monitorenter`/`monitorexit` are nops).
- **Cross-engine AOT exception propagation** — Interpreter and IR implement
  `try`/`catch`/`athrow` and bounded native exception propagation, but emitted
  AOT code still explicitly rejects paths requiring exception propagation.
- **`invokedynamic` / lambdas** — panics; compile fixtures with `-source 8` to
  get StringBuilder-based string concat instead.
- **Reflection, JNI, `sun.misc.Unsafe`** — not modeled.
- **Java memory model** — moot while single-threaded; a correctness concern for
  the concurrency arc (Java and Go memory models are related but not identical).
