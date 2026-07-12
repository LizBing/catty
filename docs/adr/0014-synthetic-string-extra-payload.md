# ADR-0014: Synthetic java.lang.String with a Go-string Extra() payload

- **Status:** Accepted
- **Date:** 2026-07-12

## Context

`java.lang.String` is one of the most-used classes in any Java program. catty
has two competing concerns:

1. **The Go↔Java bridge depends on a Go string.** `ldc "hello"` produces a
   String object; `PrintStream.println(String)` writes it to stdout; string
   constants in the constant pool are Go strings at parse time. The cheapest
   representation is to store the Go string directly in the object's `Extra()`
   payload — no `char[]`/`byte[]` to allocate and keep in sync.

2. **The real JDK's String is backed by a `byte[]` + `coder`** (compact strings,
   Java 9+). Its `hashCode`, `equals`, `substring`, `charAt`, … are bytecode
   that reads `this.value[]` and `this.coder`. Loading the real `String.class`
   would inherit all those methods for free.

These are in tension. If catty loads the real String, every `ldc` must build a
real `byte[]` and set `coder`, *and* keep `Extra()` in sync (because catty's
native `println`/`StringBuilder` read `Extra()`). If catty keeps String
synthetic with `Extra()` = Go string, the real-JDK bytecode methods can't run
(they'd read a non-existent `value[]`) — so every content method (`equals`,
`hashCode`, `substring`, …) must be re-implemented as a native Go function
operating on the `Extra()` string.

## Decision

**String stays synthetic, with the Go string in `Extra()`.** Content methods
(`length`, `charAt`, `equals`, `hashCode`, `isEmpty`, `substring`,
`concat`, `indexOf`, `startsWith`, `endsWith`, `compareTo`, `toCharArray`) are
native Go functions operating on `Extra()`. The `String(byte[], coder)`
constructor decodes bytes → Go string so JDK methods that build strings
(`Integer.toHexString`, …) produce correct output.

Rejected alternative: load real `String.class` and bridge `Extra()` ↔ internal
`byte[]` at every entry/exit. The synchronisation surface is large (every
`ldc`, every native call, every interop point) and the bridge itself becomes a
correctness hazard. The cost of the chosen path is bounded and local to
`native/lang.go`.

## Consequences

**Positive**
- `ldc` and `println` stay O(1) — no array allocation, no sync.
- String content operations are pure Go (`strings.HasPrefix`, `utf8.RuneCountInString`,
  …) — fast and obviously correct.
- The bridge payload model is uniform across all bootstrap classes (Object's
  `Extra()` for Class objects, String's for the Go string, …).

**Negative**
- **Long-tail maintenance**: every String method real JDK code calls that we
  haven't implemented is a `NoSuchMethodError`. The set grows as java.base
  coverage widens (e.g. `String.format`, `String.matches`, `replaceAll`).
- `length`/`charAt` are rune-based (BMP-correct); surrogate pairs (emoji) are
  not handled — deferred to R6.
- Two String worlds exist mentally: the Go-string payload (fast path) and the
  byte-array contract (what real bytecode expects). The `([BB)V` constructor is
  the seam where the second is decoded into the first.

This is the conscious cost of keeping the bridge simple. Revisit if the
content-method maintenance burden grows past ~30 methods or if a JDK hot path
needs direct `byte[]` access.
