# R2 String UTF-16 slice

**Status:** Proposed
**Type:** implementation
**Review:** owner
**Base commit:** `ecb086e`
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
- Define and test explicit host-text adaptation for supported output paths.
- Update Interpreter, IR, and AOT materialization/bridge paths without creating a second
  canonical Java String backing.
- Promote the six String fixtures in `docs/workstreams/r2-evidence/fixtures/` to regression
  evidence: `SupplementaryChar`, `HashDivergence`, `LoneSurrogate`, `StringBounds`,
  `StringSubstringUnits`, and `StringCharArrayRoundTrip`.

## Non-scope

- Changing the Java String facade to a classfile-backed implementation, reflection layout,
  compact strings, broad java.base compatibility, character-set/I/O completeness, or a
  complete Java String API.
- Changing Slot/object layout, concurrency/JMM, class initialization, or unrelated AOT
  refactoring.

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| String differential matrix | `bash docs/workstreams/r2-evidence/run-r2-diff.sh` | Not run |
| Interpreter / IR | All six String fixtures match Temurin 25 | Not run |
| AOT | All supported String fixtures match; unsupported paths explicit `Fallback` / `Not implemented` | Not run |
| Core regression | `go vet ./... && go test ./... && go test -race ./... && bash tests/run.sh` | Not run |
| Governance | `git diff --check` | Not run |

## Review

**owner** — only the Owner may accept this workstream or mark it Done.
