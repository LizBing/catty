# R2 String UTF-16 slice

**Status:** In Progress
**Type:** implementation
**Review:** owner
**Base commit:** `298b723`
**Roadmap item:** Phase R2 — Runtime semantics and concurrency planning
**Governing ADRs:** ADR-0016, ADR-0017, ADR-0019, ADR-0020, ADR-0023, and ADR-0027 (Accepted)
**Prerequisites:** Owner accepts this workstream

## Outcome

For the supported native String surface, catty matches Temurin 25 UTF-16 code-unit behavior
across Interpreter, IR, and AOT or explicitly records an ADR-0016 fallback. Java String
values use the ADR-0027 `[]uint16` kernel backing without exposing mutable aliases.

## Scope

- Add a lossless classfile MUTF-8-to-UTF-16-unit path for String constants and use it at
  String materialization boundaries.
- Migrate all current native String producers and consumers to the canonical units:
  constructors, `length`, `charAt`, `equals`, `hashCode`, `isEmpty`, `substring`, `concat`,
  `indexOf`, prefix/suffix comparison, `compareTo`, `toCharArray`, StringBuilder, PrintStream,
  and applicable System/native helpers.
- Implement `String(char[])` as a defensive copy and ensure `toCharArray` returns a copy.
- Implement correct Java null semantics: NPE on null constructor/method args (`String(String)`,
  `String(char[])`, `concat`, `startsWith`, `endsWith`, `compareTo`); `"null"` output for
  `PrintStream.print/println(String null)` and `StringBuilder.append(String null)`.
- Remove any mutable backing alias: no `RawUnits()` accessor, no no-copy constructor.
- Define and test explicit host-text adaptation for supported output paths: valid surrogate
  pairs → UTF-8 scalar, lone surrogates → `?` (0x3f) in `GoString()` and `println(char)`.
- Update Interpreter, IR, and AOT materialization/bridge paths without creating a second
  canonical Java String backing.
- Promote the eight String fixtures to regression evidence:
  `docs/workstreams/r2-evidence/fixtures/` — `SupplementaryChar`, `HashDivergence`,
  `LoneSurrogate`, `StringBounds`, `StringSubstringUnits`, `StringCharArrayRoundTrip`;
  `docs/workstreams/r2-string-fixtures/` — `LoneSurrogateLiteral`, `StringNativeSurface`.

## Engine completion state (frozen)

| Method | Interpreter | IR | AOT |
|---|---|---|---|
| `ldc` String constant | Supported | Supported | Supported |
| String.length / charAt / isEmpty | Supported | Supported | Supported |
| String.hashCode | Supported | Supported | Supported |
| String.equals | Supported | Supported | Supported |
| String.compareTo | Supported | Supported | Supported |
| String.concat / substring | Supported | Supported | Supported |
| String.startsWith / endsWith / indexOf | Supported | Supported | Supported |
| String.(char[]) / toCharArray / StringBuilder | Supported | Supported | Not implemented |
| Null contract (NPE / "null") | Supported | Supported | Not implemented |
| `println(char)` lone surrogate → `?` | Supported | Supported | Not implemented |

**AOT: 5 Supported, 3 Not implemented**

## Non-scope

- Changing the Java String facade to a classfile-backed implementation, reflection layout,
  compact strings, broad java.base compatibility, character-set/I/O completeness, or a
  complete Java String API.
- Changing Slot/object layout, concurrency/JMM, class initialization, or unrelated AOT
  refactoring.
- AOT lowering for `newarray` (char[]), exception handlers, or null-check branches.
- Interning, compact strings, or String Deduplication.

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| String differential matrix | `bash docs/workstreams/r2-string-fixtures/run-string-diff.sh` | Not run |
| Interpreter / IR | All eight String fixtures match Temurin 25 | Not run |
| AOT | 5 Supported String fixtures match; 3 Not implemented | Not run |
| Core regression | `go vet ./... && go test ./... && go test -race ./... && bash tests/run.sh` | Not run |
| Governance | `git diff --check` | Not run |

## Review

**owner** — only the Owner may accept this workstream or mark it Done.
