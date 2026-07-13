# ADR-0004: Synthetic native core classes instead of shipping a JRE

- **Status:** Superseded by [ADR-0019](./0019-go-native-dissolution-policy.md)
- **Date:** 2026-07-11

## Context

Even "Hello, World!" touches `java.lang.Object` (every `<init>` chain), `String`
(`ldc`), and `System.out` (`getstatic` + `invokevirtual` on `java.io.PrintStream`).
A real JVM loads these from `rt.jar` / `lib/modules`. catty has no JRE on the
classpath and no intent to ship one — that would be enormous and require
interpreting megabytes of JDK bytecode just to print a line.

## Decision

Implement the handful of core classes catty needs **natively in Go**, as
synthetic `rtda.Class` values built without any `.class` file. `native.NativeClass`
is the dispatch table: given an internal class name, it returns a built class or
`nil` (meaning "not a core class; read it from the classpath").

Each builder (`buildObjectClass`, `buildStringClass`, `buildSystemClass`, …)
calls `rtda.NewSyntheticClass(name, super)` and adds native methods via
`rtda.NativeMethod` (which takes a `func(*rtda.Frame)`) and static fields via
`AddStaticField`/`SetStaticRef`. Native methods read args from the frame's
locals and push return values onto its operand stack; the interpreter transfers
the return.

## Consequences

**Positive**
- Zero dependency on a JRE; the project is self-contained and the MVP stays
  tiny.
- Core classes can carry a Go payload in `Object.extra` — a `String`'s value is
  a Go `string`, a `PrintStream`'s sink is an `io.Writer` — so native operations
  are direct Go, not interpreted bytecode.
- Adding a needed core method is a few lines of Go (see `DEVELOPMENT.md`).

**Negative**
- The coverage is *only* what catty implements. Programs referencing unshipped
  classes/methods (`java.util.*`, reflection, `Math.abs`, …) fail until the
  relevant builder is extended. This is the most common source of "method not
  found" during real-program bring-up.
- Semantics can drift from the JDK: `String` here is a Go-string wrapper, not a
  `char[]`/`byte[]`, so `charAt`/`substring` aren't supported yet.
- It is a *choice to diverge from spec*, kept honest by the e2e diff against
  real `java` in `tests/run.sh`.

## Future direction
Migrate frequently-used core methods to interpreted JDK bytecode if a JRE ever
becomes available, keeping native implementations only where the Go interop
(e.g. I/O) genuinely needs them.
