# R2 String UTF-16 slice

**Status:** Ready
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

---

## Amendments

*A1 — Owner review round 1 (candidate `83accd9` rejected). 2026-07-13.*

1. **Fixture expansion.** The contract is upgraded from six fixtures to eight.
   `docs/workstreams/r2-string-fixtures/` holds two new fixtures:
   - `LoneSurrogateLiteral` — lone surrogate from classfile literal (MUTF-8 path).
   - `StringNativeSurface` — null contract (NPE / "null" output) + PrintStream +
     StringBuilder + lone surrogate `println(char)` → `?`.
   All eight fixtures are promoted to regression evidence.

2. **Null contract.** Correct Java null semantics are required:
   - NPE on: `new String((String)null)`, `new String((char[])null)`,
     `concat(null)`, `startsWith(null)`, `endsWith(null)`, `compareTo(null)`.
   - `"null"` output: `PrintStream.print/println((String)null)` and
     `StringBuilder.append((String)null)`.

3. **Host-text adaptation.** Valid surrogate pairs → UTF-8 scalar; lone surrogates in
   `GoString()` and `println(char)` → `?` (0x3f), matching Temurin 25.0.3.

4. **Mutable alias removal.** No `RawUnits()` accessor; no no-copy constructor. All
   access paths return defensive copies. Tests must not treat aliasing as correct.

5. **Harness.** `run-string-diff.sh` replaces `run-r2-diff.sh`, running all eight
   fixtures across Interpreter, IR, and AOT with output + exit-code comparison.
   Harness is fail-closed (non-zero exit on mismatch or missing toolchain).

6. **Repository siting.** `docs/workstreams/r2-string-fixtures/` is part of this
   workstream, not of the evidence directory. Evidence stays under
   `docs/workstreams/r2-string-evidence/<candidate>/`.

7. **Engine completion state (frozen AOT matrix).**

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

   **AOT: 5 Supported, 3 Not implemented.**

---

*A2 — Owner review round 2 (candidate `7c1dc04` / `da3e7d0` revisions). 2026-07-13.*

1. **System.getProperty(null) → NPE.** `System.getProperty(null)` must throw
   NullPointerException. Differential coverage is asserted in `StringNativeSurface`
   fixture alongside the existing six NPE cases.

2. **Evidence preservation.** All prior candidate evidence directories under
   `docs/workstreams/r2-string-evidence/` must be preserved in git history.
   No silent deletion of any candidate evidence is permitted.

3. **Workstream contract immutability.** The owner-accepted frozen workstream
   contract (this document above the `---` line) is restored from base commit
   `298b723` as-is, changing only Status from Proposed to In Progress. All
   new requirements are appended as Amendments rather than rewriting the document.

4. **Harness AOT expectation enforcement.** `run-string-diff.sh` must explicitly
   encode the AOT Supported / Not implemented matrix:
   - **Supported (must match Temurin 25):** HashDivergence, StringSubstringUnits,
     SupplementaryChar, LoneSurrogateLiteral, StringNativeSurface.
   - **Not implemented (must NO-BUILD):** LoneSurrogate, StringBounds,
     StringCharArrayRoundTrip.
   - Any other fixture producing NO-BUILD must fail-closed.
   Any fixture expected to be Supported that produces NO-BUILD or MISMATCH must
   fail-closed.

5. **Candidate hygiene.** The final candidate commit and its evidence commit must be
   cleanly recorded. The evidence directory must be bound to the final candidate.
   The worktree must be clean (no staged/unstaged changes) before the workstream is
   marked Ready.

6. **Gate rerun.** Before marking Ready, rerun and record exit codes for:
   - `bash docs/workstreams/r2-string-fixtures/run-string-diff.sh`
   - `go vet ./...`
   - `go test ./...`
   - `go test -race ./...`
   - `bash tests/run.sh`
