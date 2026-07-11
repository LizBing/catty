# ADR-0005: Lazy `<clinit>` via frame push at JVMS §5.5 points

- **Status:** Accepted
- **Date:** 2026-07-11

## Context

JVMS §5.5 specifies that a class is *initialized* (its `<clinit>` run) at
well-defined points: before `new`, `getstatic`/`putstatic` on a class's own
static fields, and `invokestatic` — and before a class's `main` is invoked.
Initialization must run superclass initialization first, and must not re-run.

Two implementation choices:

1. **Eager**: run `<clinit>` inside `classloader.LoadClass`, right after linking.
   Simple, but couples the classloader to the interpreter (running bytecode needs
   a `Thread` + the dispatch loop) — reintroducing the import cycle that
   `rtda.Loader` exists to break.
2. **Lazy**: the classloader only loads + links; the interpreter triggers
   initialization the first time a §5.5 instruction runs.

## Decision

Lazy initialization, owned by the interpreter. `interpreter.ensureInitialized`
runs `<clinit>` at the §5.5 points (`new`, `getstatic`, `putstatic`,
`invokestatic`), and the launcher calls it for the main class before `main`. It
marks the class started, recurses into the superclass first, then — if a
`<clinit>` method exists — **pushes its frame** and returns. The dispatch loop
runs that frame to completion before the caller's next instruction.

## Consequences

**Positive**
- No classloader→interpreter dependency; the cycle stays broken.
- Initialization happens exactly when the spec says, for free, on the
  instruction that needs it. The `StaticFields` fixture (which reads a static
  field initialized by `compute()` in `<clinit>`) verifies the path end-to-end.
- Superclass-first ordering and "run once" are handled by the same recursion +
  `initStarted` flag.

**Negative**
- The push-then-return shape has a subtlety the launcher must respect:
  `InitClass` needs a caller frame already on the stack, so `main`'s frame is
  pushed *before* `InitClass` (which then pushes `<clinit>` on top, so it runs
  first). Getting this order wrong panics on an empty stack.
- `<clinit>` that throws or runs very long blocks the triggering instruction;
  acceptable for MVP, but exception support (Theme C) will need to route those
  failures properly.
- Initialization triggered mid-instruction means the caller's partially-decoded
  state (operands already read) must be safely resumable — it is, because each
  §5.5 opcode reads all its operands before calling `ensureInitialized`.

## Future direction
Once exceptions land, `<clinit>` failures must propagate as
`ExceptionInInitializerError`; the current panic-on-error behavior is a
placeholder.
