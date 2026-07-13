# R2 String UTF-16 slice

**Status:** Accepted
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
| String differential matrix | `bash docs/workstreams/r2-string-fixtures/run-string-diff.sh` over the fixed eight fixtures; nonzero on any required mismatch or missing tool/fixture | Not run |
| Interpreter / IR | All eight fixtures match Temurin 25 stdout, stderr, and exit code | Not run |
| AOT | `HashDivergence`, `StringSubstringUnits`, `SupplementaryChar`, `LoneSurrogateLiteral`, and `StringNativeSurface` match; the other three are explicit `Not implemented` | Not run |
| Kernel/unit invariants | `go test ./...` includes lossless MUTF-8 units, immutable-copy boundaries, native exception return suppression, and host-output adapter tests | Not run |
| Core regression | `go vet ./... && go test -race ./... && bash tests/run.sh` | Not run |
| Evidence isolation | Historical `docs/workstreams/r2-evidence/{matrix.md,run-r2-results.txt}` unchanged; candidate results archived under `docs/workstreams/r2-string-evidence/<candidate>/` | Not run |
| Governance | `git diff --check 298b723..<candidate>` | Not run |

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
| A — kernel value, MUTF-8 unit decoding, copy/host adapters | Pending | Focused classfile and kernel unit tests |
| B — migrate native/runtime/interpreter/rtda producers and consumers | Pending | Package tests plus repository-wide producer/consumer audit |
| C — String API correctness and bounded native exception propagation | Pending | Eight-fixture Interpreter/IR differential |
| D — AOT lossless literal/bridge path and explicit unsupported matrix | Pending | Five AOT matches; three classified `Not implemented` |
| E — full gates, immutable candidate evidence, docs and self-review | Pending | Candidate evidence directory and gate summary |

---

## Handoff

- **Branch / candidate:** not started; implementation base `298b723` on `main`
- **Dirty files:** governance documents updated for acceptance; no implementation files
- **Last location:** accepted contract, before implementation
- **Checks run / not run:** `git diff --check` Pass; implementation gates not run
- **Blocker:** None; Active Agent may begin within the frozen contract
- **Next action:** Active Agent creates the work branch, marks the workstream In Progress, and starts Plan slice A
- **Non-derivable context:** local Temurin 25.0.3 UTF-8 PrintStream probe encoded an isolated high surrogate between ASCII sentinels as byte `0x3f`

## Review

**owner** — only the Owner may accept this workstream or mark it Done.

## Acceptance record

Accepted by Owner on 2026-07-13. The frozen contract authorizes implementation only within
this document's Outcome, Scope, Non-scope, Semantic constraints, and Acceptance gates.
