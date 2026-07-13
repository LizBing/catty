# String representation trade-off report (Slice C — ADR-0023)

**Evidence:** `r2-string-matrix.md` and `run-r2-results.txt`, Temurin 25.0.3 against
catty R1 in pure-synthetic mode.

## Requirement

Supported Java String behavior is defined in UTF-16 code units and must preserve arbitrary
unit sequences, including unpaired surrogates. The requirement applies to the Java value,
not necessarily to every host-text boundary.

## Evidence and current failure modes

- R1 stores a Go `string` in `Object.Extra()` and native methods use Go byte/rune operations.
  This is not a canonical UTF-16 representation.
- `SupplementaryChar` and `HashDivergence` mismatch on every current engine.
- `LoneSurrogate` crashes because `String(char[])` is absent.
- The classfile MUTF-8 decoder converts constants to Go strings. Its surrogate-pair
  recombination masks the first surrogate payload byte with `0x3f`, where the leading
  nibble must be excluded; the resulting literal path is not lossless. This is an
  implementation bug to correct as part of materialization, not a reason to choose a
  different Java representation.
- Existing native methods also clamp/return values instead of Java bounds exceptions, and
  several String-producing consumers remain Go-string based.

## Candidate comparison

| Candidate | UTF-16/lone-surrogate correctness | Complexity | Decision |
|---|---|---|---|
| Go `string` | No canonical unit representation; rune/UTF-8 adapters lose the required model | Low | Reject |
| `[]uint16` | Direct one-unit-per-element representation | Low/moderate | Select |
| Latin-1/UTF-16 compact dual encoding | Semantically possible | High: tagged layouts and branchy engine boundary | Defer pending measurements |
| Go-string/UTF-16 hybrid | Semantically possible only with careful tags/adapters | High: two canonical paths | Reject for first conforming slice |

## Recommendation

Use `[]uint16` as the canonical kernel backing. It directly supports Java unit operations
and unpaired surrogates. The backing must remain immutable from Java's perspective; copies
are required at mutable boundaries such as `String(char[])` and `toCharArray`.

This does not decide whether Java String stays synthetic, uses `Object.Extra()` forever, or
uses a particular AOT ABI. Those are facade and bridge choices. Go strings may be used only
as explicit adapters for host text or as an input representation that is converted losslessly
at a defined boundary; they cannot be a second canonical Java String representation.
