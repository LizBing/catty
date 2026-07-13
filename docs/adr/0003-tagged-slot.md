# ADR-0003: Tagged 16-byte Slot (HotSpot stack-word model)

- **Status:** Superseded by [ADR-0020](./0020-representation-domains-and-engine-boundaries.md)
- **Date:** 2026-07-11

## Context

The JVSM operand stack and local-variable array are arrays of *slots*: each
slot holds one category-1 value (byte/char/short/int/boolean/float/reference/
returnAddress) and a category-2 value (long/double) occupies two. catty must
choose a Go representation for a slot.

Options:

1. **`interface{}` slot** — `[]any`. Simple, but every numeric push boxes an
   allocation; catastrophically slow for arithmetic.
2. **Tagged struct** — `type Slot struct { num int32; ref *Object }`. One cell
   holds either an int-ish value (covers all category-1 numeric types and
   float-as-bits) or a reference. Two slots per long/double. This is the HotSpot
   "stack word" model.
3. **Parallel typed arrays** — `[]int32` for numerics + `[]*Object` for refs,
   indexed identically. Halves memory traffic on numeric code but complicates
   every handler (each opcode must know which array to touch).
4. **NaN-boxing** — pack refs and doubles into a single `uint64` (V8/SpiderMonkey
   style). Most compact, but fiddly and platform-sensitive.

## Decision

Use the tagged struct (option 2): `type Slot struct { num int32; ref *Object }`,
two slots per category-2 value, high word first.

## Consequences

**Positive**
- No boxing on numeric operations — arithmetic stays allocation-free.
- One storage type uniformly backs operand stack, locals, instance fields, and
  array elements, so a single set of typed `Frame` accessors and field/array
  helpers serves everywhere.
- GC hygiene is explicit and cheap: `PopRef` nils the freed slot so the GC can
  reclaim the object; `rtda/frame_test.go` guards this invariant.

**Negative**
- 16 bytes per slot (8-byte-aligned `int32` + pointer), vs. a machine word for
  real JVMs. That doubles memory bandwidth on numeric-heavy loops and is a real
  contributor to the ~7× gap vs. `java -Xint`.
- The `int32 num` "is" a float via `math.Float32bits` — a convention handlers
  must honor, not something the type system enforces.

## Future direction
Option 3 (parallel `[]int32`/`[]*Object`) is the documented interpreter-tuning
step (`ROADMAP.md` Theme B) once AOT is not yet available and interpreter
throughput matters. It touches every handler, so it should be done in one
measured pass with `BenchFib` as the yardstick.
