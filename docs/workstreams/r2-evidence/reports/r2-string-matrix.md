# String UTF-16 edge-case matrix (Slice C — ADR-0023)

**Evidence source:** `docs/workstreams/r2-evidence/run-r2-results.txt` (Temurin 25.0.3
versus catty R1 at `ecb086e`, pure-synthetic mode).

## Current R1 findings

| Area | R1 implementation/evidence | Java 25 requirement | Result |
|---|---|---|---|
| Literal materialization | classfile MUTF-8 is decoded through Go `string`; surrogate recombination is not lossless | preserve literal code units | gap |
| `length` / `charAt` | Go rune-oriented native path | UTF-16 unit count and unit lookup | `SupplementaryChar` mismatch |
| `hashCode` | native iterates Go rune data | `31*h + UTF-16-unit` | `HashDivergence` mismatch; reference `1772899` |
| `String(char[])` | constructor is unregistered | defensive copy of every char unit | `LoneSurrogate` crash |
| `substring` | rune-indexed and clamps invalid indices | unit indices and Java bounds exceptions | `StringSubstringUnits` / `StringBounds` mismatch |
| `toCharArray` | rune conversion | fresh char array containing exact units | `StringCharArrayRoundTrip` mismatch |
| comparison/search/concat | Go-string operations | UTF-16 unit semantics | not yet independently probed; in implementation scope |
| host output | PrintStream passes Go string to host writer | explicit adapter, no internal unit loss | policy not yet implemented |

## Fixture set

The String regression set is `SupplementaryChar`, `HashDivergence`, `LoneSurrogate`,
`StringBounds`, `StringSubstringUnits`, and `StringCharArrayRoundTrip`. The first three
establish the original representation gap; the latter three pin failure behavior, unit
substrings, and defensive char-array copying before implementation begins.

The matrix is evidence of R1 divergence, not a claim that all Java String API methods are
supported. An Accepted implementation workstream must state support or explicit fallback
per engine under ADR-0016.
