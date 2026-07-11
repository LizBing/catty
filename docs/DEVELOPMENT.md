# Development guide

How to build, run, test, navigate, and extend catty. Read this before changing
code; for the *why* behind the structure, see [`ARCHITECTURE.md`](./ARCHITECTURE.md).

## Prerequisites

| Tool | Version | Why |
|---|---|---|
| Go | 1.22+ (developed on 1.26) | the `for range <int>` loop form is used |
| JDK (`javac`, `java`) | any modern JDK (developed on Temurin 25) | compiles test fixtures and provides the reference output to diff against |

Both must be on `PATH`. Confirm with `go version` and `javac -version`.

## Build & run

```sh
go build -o catty ./cmd/jvm          # build the launcher
./catty -cp <classpath> <MainClass>  # run a program

# examples
./catty -cp tests/fixtures HelloWorld
go run ./cmd/jvm -cp tests/fixtures Fibonacci
```

`-cp` is colon-separated directories/jars (defaults to `.`). `<MainClass>`
accepts dots or slashes: `pkg.Main` or `pkg/Main`.

Unhandled VM errors (unsupported opcode, `NullPointerException`, etc.) print
`catty: <message>` and exit 1. Set `CATTY_DEBUG=1` to also dump the Go stack
trace — useful when adding new opcodes:

```sh
CATTY_DEBUG=1 ./catty -cp tests/fixtures MyBuggyProgram
```

## Test

Two independent layers:

```sh
go test ./...        # unit tests: classfile parsing, Slot/Frame encoding
./tests/run.sh       # end-to-end: compile fixtures, diff catty output vs java
```

**`tests/run.sh`** is the source of truth. It compiles every `tests/fixtures/*.java`
with `javac -source 8`, runs each main class through both `java` and the catty
binary, and reports `PASS`/`FAIL` per class — a fixture passes only when
catty's stdout is byte-identical to java's. Build the binary and clean up
`.class` files itself.

To add an e2e case, drop a `.java` with a `main` into `tests/fixtures/` and add
its class name to `MAIN_CLASSES` in `tests/run.sh`.

## Project layout

```
catty/
├── cmd/jvm/            launcher (flag parsing, main-method entry)
├── classfile/          .class binary → structs (JVMS §4)
├── classpath/          locate .class in dirs/jars/zips
├── classloader/        load + link + cache; implements rtda.Loader
├── rtda/               runtime data areas + class construction
├── interpreter/        switch dispatch loop + opcode handlers
├── native/             synthetic core classes + native Go methods
├── tests/
│   ├── fixtures/       *.java programs (the e2e corpus)
│   └── run.sh          the e2e verification harness
└── docs/               this documentation
```

## Coding conventions

- **Go style**: standard `gofmt`/`go vet`; comments on exported identifiers.
  Run `go vet ./...` before committing — it should be clean.
- **Error handling**: malformed class files and unsupported bytecode are fatal
  during loading/execution and use `panic` with a `catty:`-prefixed message;
  `cmd/jvm`'s `recover` turns them into a clean exit. Do **not** add silent
  fallbacks for unsupported opcodes — a loud panic is more useful than wrong
  results.
- **Slots are unexported**: `rtda.Slot`'s `num`/`ref` fields are touched only
  inside `rtda`. Cross-package code uses the `Num()`/`Ref()`/`SetNum()`/`SetRef()`
  accessors (for fields/arrays) or the `Frame` typed methods (for the operand
  stack / locals).
- **No import cycles**: see ARCHITECTURE §3. New run-time resolution must go
  through `rtda.Loader`, not by importing `classloader`.
- **Tests**: prefer a `tests/fixtures` e2e case over a Go unit test for
  interpreter behavior — the java-diff is the strongest signal. Use Go unit
  tests for pure data-structure invariants (e.g. `rtda/frame_test.go` guards the
  two-slot long/double encoding).

## How to extend

### Add a bytecode opcode

Most opcodes are ~3 lines. Example pattern, using a made-up `opFoo`:

1. **`interpreter/opcodes.go`** — add the mnemonic constant with its JVMS §6.5
   byte value:
   ```go
   opFoo = 0xNN
   ```
2. **`interpreter/interpreter.go`** — add a `case` to the dispatch `switch`,
   grouped under the right section comment. Read inline operands via
   `frame.ReadUint8/ReadUint16/ReadInt16/ReadInt32` (they advance `pc`), and
   manipulate the stack/locals via the `Frame` typed methods:
   ```go
   case opFoo:
       // iadd, as a template:
       b := frame.PopInt()
       a := frame.PopInt()
       frame.PushInt(a + b)
   ```
3. For branches, compute the target against `opcodePc` (the opcode's address):
   ```go
   case opIfFoo:
       if /* condition */ {
           branch(frame, opcodePc, int(frame.ReadInt16()))
       } else {
           frame.ReadInt16() // consume the offset even when not taken
       }
   ```
4. Add a `tests/fixtures/*.java` that exercises it and run `./tests/run.sh`.

Unsupported opcodes fall through to the `default` case, which panics with the
opcode byte — so a missing opcode fails loudly the first time a program hits it.

### Add a native method

Native methods live on synthetic core classes built by the `native` package.
The registry is just `native.NativeClass`'s `switch` on class name.

1. **`native/<file>.go`** — write the Go function with signature
   `func(*rtda.Frame)`. Read arguments from the frame's **locals** (not the
   stack): `locals[0]` is `this` for instance methods, args follow. Push any
   return value onto the frame's **operand stack** via `PushInt`/`PushRef`/…
   The interpreter transfers it to the caller.
   ```go
   func mathAbs(f *rtda.Frame) {
       x := f.GetInt(0)
       if x < 0 { x = -x }
       f.PushInt(x)
   }
   ```
2. Register it on the class in the relevant `buildXxxClass` builder:
   ```go
   c.AddMethod(rtda.NativeMethod(c, "abs", "(I)I", mathAbs))
   ```
   `NativeMethod` derives `argSlotCount`, `maxLocals`, `maxStack` from the
   descriptor — `this` is added by the interpreter, not counted here.

### Add a core class

If a program references a JDK class catty doesn't ship (e.g. `java.util.List`),
add a builder:

1. **`native/registry.go`** — add a `case "java/util/List": return buildListClass(loader)` arm.
2. **`native/<file>.go`** — write `buildListClass(loader rtda.Loader) *rtda.Class`:
   ```go
   func buildListClass(loader rtda.Loader) *rtda.Class {
       c := rtda.NewSyntheticClass("java/util/List", loader.LoadClass("java/lang/Object"))
       c.AddMethod(rtda.NativeMethod(c, "size", "()I", listSize))
       // ...static fields via c.AddStaticField(...).SetStaticRef(...)
       return c
   }
   ```
   `NewSyntheticClass(name, super)` gives an empty class; `AddMethod` /
   `AddStaticField` populate it. If the class has a `<clinit>`-equivalent
   initializer, do it inline in the builder (like `buildSystemClass` sets `out`
   /`err`).

### Add a test fixture

1. Write `tests/fixtures/MyCase.java` with a `public static void main`.
2. Add `MyCase` to `MAIN_CLASSES` in `tests/run.sh`.
3. `./tests/run.sh`.

Keep fixtures small and focused on one feature; they double as living examples.

## Debugging tips

- **`CATTY_DEBUG=1`** prints the Go stack trace on a VM panic.
- **`go test -run TestParseHelloWorld -v ./classfile/`** confirms parsing before
  you trust the loader/interpreter.
- For a hang, suspect an infinite loop in your new branch logic or a method
  that never returns. Temporarily log `frame.PC()` / `op` at the top of `Loop`.
- For wrong numeric results, check the two-slot long/double encoding first
  (the `rtda/frame_test.go` unit tests guard it).
