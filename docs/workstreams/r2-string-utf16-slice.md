# R2 String UTF-16 slice

**Status:** Done
**Type:** implementation
**Review:** owner
**Base commit:** `298b723`
**Roadmap item:** Phase R2 — Runtime semantics and concurrency planning
**Governing ADRs:** ADR-0016, ADR-0017, ADR-0019, ADR-0020, ADR-0023, and ADR-0027 (Accepted)
**Prerequisites:** R2 initialization slice Done; Owner accepted this workstream on 2026-07-13

## Outcome

For the bounded synthetic/native String surface, catty matches Temurin 25 UTF-16
code-unit behavior in Interpreter and IR. AOT uses the same canonical values on its
currently emittable paths, and every non-emittable fixture is reported as `Not implemented`
rather than silently approximated. Every Java String value uses an immutable ADR-0027
`[]uint16` kernel backing.

## Scope

- Introduce one shared immutable String-kernel value backed by `[]uint16`. Construction and
  read APIs must prevent mutable aliases from escaping; a Go `string` may exist only as an
  explicit host-text adapter result.
- Add a lossless classfile MUTF-8-to-UTF-16-unit accessor for `CONSTANT_String` entries while
  retaining text-oriented accessors for names and descriptors. Interpreter and AOT literal
  materialization must consume the lossless accessor.
- Migrate every current Java String producer and consumer to the canonical value:
  constructors, `length`, `charAt`, `equals`, `hashCode`, `isEmpty`, both `substring`
  overloads, `concat`, `indexOf(int)`, prefix/suffix comparison, `compareTo`, `toCharArray`,
  StringBuilder, PrintStream, exception messages, Object/Class/System helpers, and the AOT
  runtime bridge.
- Add `String(char[])` with a defensive copy; make `toCharArray` return a fresh array; retain
  exact units through substring, concatenation, builder append, and literal materialization.
- Implement Java-visible failure behavior needed by the supported methods, including
  `StringIndexOutOfBoundsException`, null-argument behavior, and propagation from a native
  method back to the invoking Interpreter/IR instruction without attempting to transfer a
  normal return value.
- Define the supported stdout adapter: valid surrogate pairs become their UTF-8 scalar;
  each unpaired surrogate is encoded as `?` (`0x3f`), matching the pinned Temurin 25.0.3
  UTF-8 PrintStream observation. This policy is an output-boundary conversion, never a
  second String backing.
- Preserve the six research fixtures and add two focused fixtures before acceptance:
  `LoneSurrogateLiteral` for lossless classfile literal materialization and
  `StringNativeSurface` for the currently claimed comparison/search/concat/StringBuilder
  surface. Keep the new sources and fail-closed harness under
  `docs/workstreams/r2-string-fixtures/`; candidate evidence is archived separately from
  the historical research matrix.
- Update architecture/development documentation whose String representation or native-call
  exception description becomes stale.

## Non-scope

- A classfile-backed `java/lang/String` facade, Java-visible String field layout, compact
  strings, interning semantics, charset/file APIs, a complete Java String API, or broad
  `java.base` compatibility.
- AOT lowering for `newarray`, exception handlers, `athrow`, or general cross-engine
  exception propagation. Consequently `LoneSurrogate`, `StringCharArrayRoundTrip`, and
  `StringBounds` may remain AOT `Not implemented` in this slice.
- Slot/heap-layout redesign, concurrency/JMM, further class-initialization work, JIT, or
  unrelated AOT refactoring.
- Claiming `Integer/Long.toString`, `Double.parseDouble`, or representative `HashMap`
  compatibility; their remaining runtime/library dependencies are separate work.

## Semantic constraints

- Java 25 is the capability baseline and repository-pinned Temurin 25.0.3 is the
  differential reference.
- `length`, indices, ordering, hashing, equality, substring, and char-array conversion are
  defined over UTF-16 code units. `indexOf(int)` follows the Java API's code-point rules,
  including supplementary-code-point search and UTF-16 result indices.
- Every `uint16` value, including isolated high and low surrogates and MUTF-8 encoded NUL,
  survives all supported Java-to-Java paths exactly.
- `String(char[])`, `toCharArray`, substring/builder results, host adapters, and bridge APIs
  cannot expose a mutable alias that changes an existing String.
- Bounds and null failures use the Java exception type and validation order required by the
  supported API. A pending exception suppresses native return transfer.
- Interpreter, IR, and AOT share the same object/value world. AOT support claims are made
  fixture by fixture as `Supported`, `Fallback`, or `Not implemented`.

## Required completion state by engine

| Capability | Interpreter | IR | AOT |
|---|---|---|---|
| Canonical UTF-16 value and lossless literals | Required | Required | Required on emitted paths |
| Six historical String fixtures | 6/6 match | 6/6 match | Three match; three explicit `Not implemented` allowed by Non-scope |
| `LoneSurrogateLiteral` | Match | Match | Match |
| `StringNativeSurface` | Match | Match | Match |
| Java bounds/null failures in native String methods | Required | Required | `Not implemented` where AOT exception handling is required |

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| String differential matrix | `bash docs/workstreams/r2-string-fixtures/run-string-diff.sh` over the fixed eight fixtures; nonzero on any required mismatch or missing tool/fixture | Pass — 8/8, candidate `00327d6` |
| Interpreter / IR | All eight fixtures match Temurin 25 stdout, stderr, and exit code | Pass — 8/8 each |
| AOT | `HashDivergence`, `StringSubstringUnits`, `SupplementaryChar`, `LoneSurrogateLiteral`, and `StringNativeSurface` match; the other three are explicit `Not implemented` | Pass — 5 Supported / 3 Not implemented |
| Kernel/unit invariants | `go test ./...` includes lossless MUTF-8 units, immutable-copy boundaries, native exception return suppression, and host-output adapter tests | Pass |
| Core regression | `go vet ./... && go test -race ./... && bash tests/run.sh` | Pass — 10/10 end-to-end fixtures |
| Evidence isolation | Historical `docs/workstreams/r2-evidence/{matrix.md,run-r2-results.txt}` unchanged; candidate results archived under `docs/workstreams/r2-string-evidence/<candidate>/` | Pass |
| Governance | `git diff --check 298b723..<candidate>` | Pass |

Results use only `Pass`, `Fail`, `Not run`, or `Not implemented`.

## Amendments

Accepted changes are appended here after Owner approval; the frozen contract is not
rewritten to reduce gates.

---

## Investigation findings

- The current `decodeMUTF8` collapses valid surrogate pairs into a Go UTF-8 string and cannot
  retain an isolated surrogate through String materialization. A separate lossless unit
  accessor is required; replacing all constant-pool text access would unnecessarily affect
  names and descriptors.
- Java Strings are currently created independently in `interpreter`, `runtime`, `native`,
  and `rtda` exception helpers. Direct `Extra().(string)` reads also exist in exception and
  Class/System paths, so changing only `native/lang.go` would leave two representations.
- Native invocation always transfers a declared return after the native function returns.
  A throwing `charAt` would therefore pop a nonexistent return slot before the interpreter
  can dispatch the pending exception. The fix must be bounded to pending-exception return
  suppression and covered in both Interpreter and IR.
- The current AOT emitter handles String literals and the methods used by three historical
  String fixtures, but not primitive-array allocation or exception-handler control flow.
  The explicit five-supported/three-`Not implemented` target above keeps this slice bounded.
- The six research fixtures do not directly prove lossless lone-surrogate literals or the
  comparison/search/concat/StringBuilder paths already named in ADR-0027; the two added
  fixtures close those evidence gaps without expanding to the broad String API.

## Plan

| Slice | Status | Evidence |
|---|---|---|
| A — kernel value, MUTF-8 unit decoding, copy/host adapters | Complete | Classfile and kernel unit tests pass |
| B — migrate native/runtime/interpreter/rtda producers and consumers | Complete | Producer/consumer audit and differential coverage |
| C — String API correctness and bounded native exception propagation | Complete | Eight-fixture Interpreter/IR differential pass |
| D — AOT lossless literal/bridge path and explicit unsupported matrix | Complete | Five AOT matches; three `Not implemented` |
| E — full gates, immutable candidate evidence, docs and self-review | Complete | Evidence `9008b00`; owner review accepted |

---

## Handoff

- **Branch / candidate:** `worktree-r2-string-utf16-slice`; candidate `00327d6`, evidence `9008b00`
- **Dirty files:** none at owner acceptance
- **Last location:** integrated to `main`
- **Checks run / not run:** all acceptance gates Pass
- **Blocker:** none
- **Next action:** proceed to the next bounded R2 workstream
- **Non-derivable context:** local Temurin 25.0.3 UTF-8 PrintStream probe encoded an isolated high surrogate between ASCII sentinels as byte `0x3f`

## Review

**owner** — only the Owner may accept this workstream or mark it Done.

## Acceptance record

Accepted by Owner on 2026-07-13. The frozen contract authorizes implementation only within
this document's Outcome, Scope, Non-scope, Semantic constraints, and Acceptance gates.

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

## Candidate evidence

- **Implementation candidate (C):** `00327d6`
- **Evidence commit (E):** `9008b00`
- **Ready record:** `bfaa6f6`
- **Evidence:** `docs/workstreams/r2-string-evidence/00327d6/`

## Completion record

Owner accepted the Ready candidate and authorized integration on 2026-07-13. Integrated to
`main` after all acceptance gates passed; this workstream is Done.
