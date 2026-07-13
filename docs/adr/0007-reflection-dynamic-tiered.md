# ADR-0007: Reflection & dynamic features — tiered, keep the interpreter

- **Status:** Superseded by [ADR-0016](./0016-multi-engine-execution.md)
- **Date:** 2026-07-11

## Context

Java's dynamic surface — reflection (`java.lang.reflect`), dynamic class loading
(`Class.forName`, custom `ClassLoader`, `ServiceLoader`), `invokedynamic` /
lambdas / string concat, and dynamic proxies — is inherently **open-world**:
programs resolve classes, methods, and fields by *name at run time*, including
classes that may not exist at compile time.

This collides head-on with catty's performance path: **AOT transpilation**
(`transpile.Emit` → `go build`) is inherently **closed-world** — it must know
every class/method/field at compile time to emit Go source for it. This is the
same conflict every AOT JVM faces (GraalVM native-image most prominently).

Two structural facts about catty shape the answer:

1. **catty retains full runtime metadata.** `rtda.Class` carries the class name,
   its `Method`s (name + descriptor + access flags + code), its `Field`s (name +
   descriptor + access flags + slot offset), and the constant pool. `Object.Fields()`
   exposes slot storage. That is ~90% of what reflection needs — kept live because
   the interpreter uses the same data.
2. **the interpreter is always present.** The classloader loads from the classpath
   at run time; any loaded class is interpretable. (A0 confirmed we keep the
   interpreter even as AOT matures — it is the IR-execution validation tier and
   the cold path.)

## Decision

Dynamic features live in the **interpreter tier**; the **AOT tier handles
statically-resolvable hot code**. They coexist (tiered), like HotSpot
(C1/C2 + interpreter) and ART (AOT + JIT) — **not** like GraalVM native-image
(pure AOT + reachability metadata).

Concretely:

- **Reflection API** (`Method.invoke`, `Field.get/set`, `getDeclaredFields`,
  `Class.forName`) → native methods in `native/` operating on the retained
  `rtda` metadata. `Field.get/set` reads/writes `Object.Fields()[slotID]`.
  `Method.invoke` runs the method through the interpreter (or, for an AOT'd
  method, through a name → Go-func dispatch registry). **No reachability
  configuration is required** — the metadata is already live.
- **Dynamic class loading** → the existing `classloader.LoadClass` (runtime,
  classpath-sourced). Loaded classes are interpretable. Classes absent at AOT
  time can only be interpreted (the AOT path cannot compile what did not exist).
- **`invokedynamic` / lambdas** → model `CallSite` + bootstrap methods. The
  bootstrap runs as native/interpreted code and caches a target; the interpreter
  dispatches through it; AOT'd callers go through a runtime call-site table.
  `LambdaMetafactory` synthesizes an adapter `rtda.Class` (as core classes are
  synthesized today) implementing the functional interface.
- **Dynamic proxies** → synthesize an `rtda.Class` implementing the interfaces,
  delegating to an `InvocationHandler`.
- **The interpreter is never deleted.** When AOT matures it remains the home of
  dynamic features, the cold path, and the validation tier.

## Consequences

**Positive**
- Reflection maps to native methods on already-retained metadata — no
  reachability-config burden, no "unregistered reflection fails" failure mode.
- Cold start stays fast (the interpreter is light; ~3 ms today).
- Graceful by construction: anything the AOT path can't handle (dynamic call
  sites, runtime-loaded classes) degrades to the interpreter instead of failing.

**Negative**
- **World-transition cost.** When AOT-compiled Go code hits a dynamic feature it
  must call back into the interpreter/runtime (analogous to HotSpot deopt /
  GraalVM reverse-embedding). Dynamic paths are slower than static ones.
- **Open-world classes are interpreted only.** A program that `Class.forName`s a
  class not present at AOT time cannot AOT it — it runs in the interpreter.
- **Speculative AOT invalidation.** If closed-world assumptions baked into
  emitted Go (e.g. a devirtualized call target) are invalidated by a later-loaded
  class, the affected method must deopt to the interpreter — a future mechanism
  to build.

## Alternatives considered
- **Pure AOT + reachability metadata** (GraalVM native-image style): require a
  config listing reflectively-accessed members, emit reflection stubs. Rejected
  for catty: needs per-program config, fails on unregistered reflection, and
  discards the "metadata already live" advantage the interpreter gives us.
- **Delete the interpreter when AOT matures**: rejected — it is the home for
  dynamic features, the cold path, and the validation tier (A0).
- **Runtime code generation via `plugin`** for dynamically-loaded classes: a
  stretch goal only — Go's `plugin` is platform-limited and version-fragile; the
  interpreter covers the same need more robustly today.

## Roadmap placement
Theme C (spec coverage), interpreter tier. Prerequisites: **exceptions** (the
reflection API throws `NoSuchMethodException` / `IllegalAccessException` etc., so
reflection semantics can't be correct without try/catch) and
**`invokedynamic`/lambdas**. Reflection proper comes after those; the AOT↔
interpreter dispatch registry is the integration step (Theme A, A4).
