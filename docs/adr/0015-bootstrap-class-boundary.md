# ADR-0015: The bootstrap-class boundary

- **Status:** Accepted
- **Date:** 2026-07-12

## Context

ADR-0009 proposes a hybrid class library: ~50 classes native in Go, ~7000
loaded as real bytecode from the JDK. But it doesn't say *which* classes must be
native, or *when* a native class overrides a real one. Without a rule, the
`classloader` had a hardcoded "synthetic wins" order — meaning `java.util.List`
would never load from java.base even when it was available, and the boundary
between "catty's Go implementation" and "the JDK's bytecode" was implicit.

The Provider-chain refactor (Block A) made the loading strategy configurable,
but still needs a principled answer: which classes are **irreducibly synthetic**
— they can never be replaced by bytecode, regardless of what's on the classpath?

## Decision

A small set — **`BootstrapClasses`** — is always served synthetically, ahead of
any classpath:

| Class | Why irreducible |
|---|---|
| `java/lang/Object` | Root of every class; must exist before any class links. |
| `java/lang/String` | Carries the Go-string `Extra()` bridge payload (ADR-0014). |
| `java/lang/Class` | `Extra()` holds the `*rtda.Class` it represents; reflection bridge. |
| `java/lang/System` | `out`/`err` initialised eagerly to wrap `os.Stdout`/`os.Stderr`. |
| `java/lang/Thread` | Will map to a goroutine (ADR-0010); the bridge starts here. |
| `java/lang/Throwable` | Root of the exception hierarchy; `detailMessage` accessed by the interpreter's exception path. |

`BootstrapProvider` claims exactly these six; `SyntheticProvider` serves the
rest of the registry (`StringBuilder`, `PrintStream`, exception subclasses,
`Comparable`, …) only when `ClasspathProvider` can't — i.e. when no real JDK is
present. With java.base on the classpath, `ArrayList`/`Math`/`Integer` load as
real bytecode; the six bootstrap classes stay synthetic.

**Criteria for adding to `BootstrapClasses`**: the class carries a catty-native
payload in `Extra()` that real bytecode cannot satisfy, OR it must exist before
the classloader can parse any class file. If neither holds, prefer loading the
real bytecode (let `SyntheticProvider` be the fallback only).

## Consequences

**Positive**
- The hybrid boundary is explicit and auditable (`native.BootstrapClasses`).
- Real java.base classes load as bytecode where possible — semantic
  compatibility with OpenJDK for free.
- The six bridge classes can't be accidentally shadowed by a stray `.class`
  file on the classpath.

**Negative**
- Adding a 7th bootstrap class is a deliberate act (edit the set + write a
  builder) — friction by design.
- Real-JDK improvements to the six bootstrap classes (e.g. a new `String`
  method) must be re-implemented natively — the maintenance tail noted in
  ADR-0014.

This sharpens ADR-0009 from "~50 native" into "6 irreducible + as many
fallbacks as the no-JDK path needs".
