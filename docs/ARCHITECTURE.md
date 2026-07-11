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
| Thread scheduler | MVP is single-threaded; the concurrency arc maps `java.lang.Thread` to a goroutine and uses the Go GMP scheduler. |
| JIT compiler | None in Phase 1. Phase 2 replaces the interpreter with a bytecode→Go AOT transpiler that hands optimization to the Go compiler. |

This trade is the whole point of the project: by not writing a GC, scheduler, or
JIT, the implementation is small (~3300 LOC) and the MVP floor is bounded — the
work collapses to "a bytecode interpreter plus a class loader." The cost is a
performance ceiling the interpreter cannot break (see §7).

## 2. The pipeline

A catty run is a linear pipeline, one Go package per stage:

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
                                 │ cmd/jvm │  launcher: wires it all together
                                 └─────────┘
```

Stage responsibilities:

1. **`classpath/`** — find `<name>.class` on disk (directories, jars, zips).
2. **`classfile/`** — decode raw bytes into a `*classfile.ClassFile` per JVMS §4
   (constant pool, members, Code attribute, modified-UTF-8).
3. **`classloader/`** — load + link + cache. Converts a `classfile.ClassFile`
   into a runtime `rtda.Class`: resolves the superclass and interfaces,
   computes field slot offsets, builds `rtda.Method`s. Routes three class kinds:
   array types, natively-implemented core classes, and ordinary user classes.
4. **`rtda/`** — the runtime data areas (JVMS §2.5): `Slot`, `Frame`, `Thread`,
   `Class`, `Method`, `Field`, `Object`, `Array`. Pure data + the class
   construction logic.
5. **`interpreter/`** — a switch-dispatch loop over bytecode. Resolves constant
   pool refs at run time, manages the operand stack and locals, invokes methods.
6. **`native/`** — the core JDK classes (`java.lang.Object`, `String`,
   `StringBuilder`, `System`, `java.io.PrintStream`) implemented in Go rather
   than loaded from a JRE.
7. **`cmd/jvm/`** — parses `-cp` and the main class, finds `main(...)`, builds
   a `Thread`, and enters the interpreter loop.

## 3. Package dependency graph

```
classfile          (stdlib only)
classpath          (stdlib only)
rtda        ──▶ classfile
native      ──▶ rtda
classloader ──▶ classfile, classpath, rtda, native
interpreter ──▶ classfile, rtda
cmd/jvm     ──▶ classpath, classloader, interpreter, rtda
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
  e.g. a `java.lang.String`'s Go `string` value, or a `PrintStream`'s
  `io.Writer`.

### `rtda.Thread` — execution context  *(thread.go)*

```go
type Thread struct {
    stack  []*Frame
    loader Loader
}
```

The MVP runs a single `Thread`. Pushing a frame per call and popping on return
gives the JVM stack. The concurrency arc promotes this to a goroutine.

## 5. Class loading, in detail

`classloader.ClassLoader.LoadClass(name)` is the single entry point used both at
startup and at run time (via `thread.Loader()`). It memoizes into a `cache`
map. For a cache miss it routes by name:

| Name shape | Source |
|---|---|
| `[...` (array) | `rtda.NewArrayClass(name, loader)` — parses the descriptor, sets `componentClass` (object/component arrays resolve through the loader) or `componentKind` (primitive). |
| a core class | `native.NativeClass(loader, name)` — returns a synthetic `*rtda.Class` built from scratch (no `.class` file). |
| anything else | `loadFileClass` → `classpath.ReadClass` → `classfile.Parse` → `rtda.NewClass(cf, loader)`. |

`rtda.NewClass` is where linking happens: it loads the superclass and interfaces
recursively, computes instance field offsets starting from the superclass's
count (so subclasses inherit offsets), allocates static var slots, and builds a
`Method` per declared method. Class **initialization** (`<clinit>`) is not done
here — see ADR-0005.

## 6. An execution trace — `System.out.println("Hello, World!")`

1. `cmd/jvm` loads `HelloWorld` (caches it; `rtda.NewClass` recursively loads the
   native `java.lang.Object`, `System`, `PrintStream`, `String`).
2. A `Thread` is created; a frame for `main` is pushed with `locals[0] = null`
   (the `args` array); `InitClass` runs `HelloWorld.<clinit>` if present (pushed
   on top, so it runs before `main`).
3. `interpreter.Loop` enters its dispatch cycle. `main`'s bytecode is roughly:
   ```
   getstatic   System.out        ; push the PrintStream for stdout
   ldc         "Hello, World!"   ; push a String object (extra = the Go string)
   invokevirtual PrintStream.println(String)
   ...
   ```
4. `getstatic` resolves the fieldref via the owner class's constant pool, loads
   `java.lang.System`, looks up `out`, and pushes the `PrintStream` object that
   `native.buildSystemClass` stored (its `extra` is `os.Stdout`).
5. `ldc` builds a `String` object whose `extra` holds the Go string.
6. `invokevirtual` peeks the receiver (the `PrintStream`), does dynamic dispatch
   on `receiver.Class()`, finds `println(String)` (native), and calls
   `invokeNative`: it copies args into a throwaway frame, runs
   `native.printlnString`, which does `fmt.Fprintln(os.Stdout, value)`.
7. `return` pops `main`'s frame; the stack drains; `Loop` exits.

## 7. Where the performance comes from (and where it doesn't)

The interpreter is a single dense `switch` on the opcode byte (see ADR-0002).
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

## 8. What catty does *not* model (yet)

Kept out of scope deliberately; each is a documented future work item in
[`ROADMAP.md`](./ROADMAP.md):

- **Concurrency** — single-threaded; no `synchronized`/`wait`/`notify`
  (`monitorenter`/`monitorexit` are nops).
- **Exceptions** — no `try`/`catch`, no `athrow` handling (exception tables are
  parsed but unused).
- **`invokedynamic` / lambdas** — panics; compile fixtures with `-source 8` to
  get StringBuilder-based string concat instead.
- **Reflection, JNI, `sun.misc.Unsafe`** — not modeled.
- **Java memory model** — moot while single-threaded; a correctness concern for
  the concurrency arc (Java and Go memory models are related but not identical).
