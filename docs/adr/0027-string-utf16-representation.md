# ADR-0027: UTF-16 String kernel backing

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** None (resolves the representation choice left open by [ADR-0023](./0023-string-semantics-and-representation-boundary.md))
- **Governing research workstream:** `r2-runtime-semantics-research` (Slice C)
- **Differential evidence:** `docs/workstreams/r2-evidence/reports/r2-string-matrix.md`,
  `docs/workstreams/r2-evidence/reports/r2-string-representation.md`

## Context

ADR-0023 requires Java 25 String semantics, including UTF-16 code units and arbitrary
code-unit sequences. R1 stores a Go `string` in `Object.Extra()` and is observably wrong
for supplementary characters, lone surrogates, indexing, and hashing. Its classfile MUTF-8
path also currently converts literal data through a Go-string representation, which is not
a sufficient lossless boundary for Java code units.

The decision needed here is the semantic backing of a Java String value. It is not a
decision that `java/lang/String` must remain synthetic, that a particular Java field layout
is permanent, or that every Go call boundary must expose a `[]uint16` ABI.

## Decision

The canonical kernel backing of every Java `String` value SHALL be a sequence of UTF-16
code units. The first conforming implementation SHALL use **`[]uint16`** for that backing.

### Semantic invariants

- `length`, `charAt`, `hashCode`, equality, ordering, substring indices, search, and
  `toCharArray` SHALL use UTF-16 code-unit semantics for every supported method.
- The backing SHALL preserve every `uint16` value, including unpaired high and low
  surrogates.
- Java String values are immutable. No constructor, `substring`, `StringBuilder`,
  `toCharArray`, or adapter may expose a mutable alias that can change an existing String's
  logical code units. Backing sharing is permitted only when it is not Java-observable.
- String literal materialization SHALL obtain code units losslessly from classfile MUTF-8.
  A conversion through Go runes or ordinary UTF-8 decoding is not an acceptable general
  mechanism for that boundary.

### Boundaries that remain open

- The current synthetic `java/lang/String` facade may place the backing in `Object.Extra()`.
  This is a current implementation technique, not a decision that future ordinary
  `java.base` String bytecode or a Java-visible field layout cannot participate.
- Interpreter, IR, and AOT must observe the same code units for the same Java String object,
  but this ADR does not prescribe an AOT Go function signature. A bridge may materialize
  literals from a Go `string` only through an explicit, lossless conversion boundary.
- Conversion between Java UTF-16 units and host text (`stdout`, files, OS APIs) must be an
  explicit adapter with documented behavior for unpaired surrogates. A host Go `string` is
  not a second canonical Java String representation.

### Compact encoding deferred

A Latin-1/UTF-16 dual encoding remains a future optimization decision. It must preserve the
same semantic invariants and be supported by measured benefit under ADR-0024; it is not
required for the first conforming implementation.

## Consequences

- The R2 String implementation must revise all existing String-producing and String-
  consuming native paths, not only `length`, `charAt`, and `hashCode`. This includes String
  constructors, `concat`, search/comparison methods, `toCharArray`, StringBuilder,
  PrintStream, System/native helpers, and AOT materialization.
- `String(char[])` must copy UTF-16 units; `toCharArray` must return a separate array.
- The classfile package needs a lossless MUTF-8-to-code-unit API for String constants;
  current Go-string accessors remain suitable only where their text semantics are sufficient
  (for example, ordinary ASCII names/descriptors).
- The future implementation workstream must test bounds exceptions, unit-indexed substring,
  char-array round trips, literals, valid supplementary text, and lone-surrogate retention
  across every claimed engine.
- Reflection layout, ordinary java.base String participation, compact encoding, and broad
  String API coverage remain separate decisions.
