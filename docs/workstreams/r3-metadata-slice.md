# R3 dynamic metadata kernel slice

**Status:** Ready
**Type:** implementation
**Review:** owner
**Profile:** Catty JVMS Core shared kernel; no public profile API
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0024, ADR-0031,
ADR-0032, and ADR-0034
**Prerequisites:** `r3-reflection-dynamic-research` Done; acceptance anchor
fixed before implementation
**Acceptance anchor / actual base:** `f685526` / `f685526`
**Branch:** `codex/r3-metadata-slice`

## Outcome

Catty validates and immutably retains the bounded classfile metadata needed by
the accepted InvokeDynamic kernel: BootstrapMethods, usable MethodType,
MethodHandle, InvokeDynamic operands, and structural ConstantDynamic entries.
Runtime metadata can reference these values without loading classes or
executing bootstrap code during parsing.

This slice creates no Java-visible reflection, annotation, MethodHandle,
CallSite, or InvokeDynamic capability.

## Scope

- Typed BootstrapMethods structures with validated indexes, tags, lengths, and
  immutable accessors.
- Usable constant-pool accessors for MethodType, MethodHandle reference kinds
  1–9, InvokeDynamic name/descriptor/bootstrap index, and structural
  ConstantDynamic.
- Profile-neutral immutable metadata attachment needed by later runtime
  consumers under ADR-0031.
- Explicit handling of structurally valid unrecognized attributes according to
  JVMS rules; malformed known structures return typed parse failures.
- Parser/unit fixtures and exact AOT reachability classification where current
  generic diagnostics can be narrowed without later linkage work.

## Non-scope

Annotation attributes or element trees, declared-member discovery, Exceptions
or MethodParameters reflection metadata, Java facades, class lookup/definition,
MethodHandle execution, opcode `0xba`, CallSite state, generated classes,
concat, lambda, Proxy, or AOT execution/fallback.

## Semantic constraints

Parsing is total and side-effect-free: no Java loading, initialization, facade
allocation, provider lookup, or bootstrap execution. Parsed metadata is
immutable and does not expose reader buffers or constant-pool internals as a
runtime ABI. Structurally valid unknown attributes are ignored for execution;
known accepted structures cannot be silently discarded.

## Acceptance

| Gate | Required result |
|---|---|
| Dynamic metadata parser | MethodHandle kinds 1–9, MethodType, InvokeDynamic, BootstrapMethods, and structural ConstantDynamic pass positive/negative tests |
| Immutability/attachment | Runtime consumers retain exact declared symbolic metadata without eager resolution or mutable reader aliases |
| Parse failure | Malformed indexes/tags/lengths return typed errors; no bounds panic or partial publication |
| Capability honesty | Fixed 24-row R3 baseline remains reproducible; no Java-visible row is newly claimed Supported |
| Regression | `go vet ./...`, `go test ./...`, `go test -race ./...`, and `bash tests/run.sh` Pass |
| Isolation/governance | Historical evidence unchanged; candidate evidence isolated; `git diff --check` Pass |

## Plan

| Step | State |
|---|---|
| Owner accepts frozen contract | Complete |
| Fix acceptance anchor and implementation preflight | Complete |
| Establish a total parser error boundary and attribute locations | Complete |
| Implement typed dynamic constant-pool structures and validation | Complete |
| Parse and validate BootstrapMethods | Complete |
| Attach immutable metadata to runtime classes | Complete |
| Evaluate a narrowly scoped AOT diagnostic improvement | Omitted — no coupling to K4 linkage; generic failure unchanged |
| Run contract gates and fix candidate | Complete |
| Owner reviews K1 candidate | Pending |

## Acceptance record

Accepted by Owner on 2026-07-18. Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, and owner review are
frozen. Implementation authorization takes effect only after prerequisites are
Done and this contract is fixed in an acceptance-anchor commit; the preflight
below records those conditions as satisfied.

## Implementation preflight

- **Acceptance anchor / actual base:** `f685526` / `f685526`
- **Branch:** `codex/r3-metadata-slice`
- **Selected by Owner:** 2026-07-18 as the only active implementation
  workstream
- **R3 predecessor:** `r3-reflection-dynamic-research` Done at candidate
  `f685526`
- **Capability boundary:** parser/runtime metadata kernel only; no
  Java-visible R3 row becomes Supported
- **Historical evidence:** existing R2 and fixed R3 baseline evidence remains
  immutable; K1 candidate evidence uses a new isolated destination

## Implementation research — 2026-07-18

This section records the implementation investigation for the accepted K1
contract. It refines the execution order without changing the frozen Outcome,
Scope, Non-scope, Semantic constraints, or Acceptance gates.

### Existing path and gap

The current path is `classloader` -> `classfile.Parse` -> `rtda.NewClass`.
`rtda.Class` already retains the parsed constant pool, so K1 can extend that
ownership without introducing a resolver or a profile API. The missing pieces
are concentrated in the classfile layer:

- `classfile.Parse` has an error return, but reader underflow, invalid magic,
  unknown constant tags, and some malformed structures currently panic.
- MethodType, MethodHandle, and InvokeDynamic are decoded only as private raw
  indexes. ConstantDynamic is not decoded.
- BootstrapMethods is currently retained only as an unparsed attribute.
- Attribute parsing does not distinguish class, field, method, and Code
  locations.
- Existing constant-pool convenience lookups use zero values for type
  mismatches and are unsuitable as the validation API for K1.

The AOT path discovers class references through the constant pool and later
fails unsupported lowering generically. No change to lowering or linkage is
required for K1.

### Chosen parser error boundary

K1 will keep the public `Parse(data) (*ClassFile, error)` shape and introduce a
typed classfile format error. Bounds-checked reader primitives will raise only
a private parser sentinel; `Parse` will recover that sentinel and return its
typed error while re-raising unrelated programming faults. This provides one
transactional publication boundary without threading an error return through
every nested attribute reader.

All direct reads participating in classfile decoding must use the checked
reader boundary, including constant UTF-8 payloads and attribute bodies.
Invalid magic, unsupported or malformed constant tags, invalid lengths, and
reserved malformed frame encodings become format errors. A failed parse
publishes no partial `ClassFile`.

This mechanism is internal control flow, not a public panic contract. Tests
must assert the returned error type and must separately prove that truncated
inputs and malformed lengths do not escape as bounds panics.

### Attribute location and validation phases

Attribute parsing will carry an internal location value for ClassFile, field,
method, and Code. BootstrapMethods is accepted only at ClassFile location and
at most once. Unknown attributes remain structurally skipped after their body
length is safely consumed; their mere presence is not an error.

Validation is split into two phases:

1. Decode bounded bytes into immutable symbolic structures.
2. After the constant pool and class attributes are available, validate
   cross-references, descriptor categories, MethodHandle kind/tag/name rules,
   bootstrap table references, and bootstrap arguments.

The second phase is necessary because constant-pool references may point
forward. It performs structural validation only: it must not load a class,
resolve a member, invoke a bootstrap method, or reject a structurally valid
ConstantDynamic dependency cycle that is specified to fail only during
resolution.

Version-dependent validation is limited to rules that affect an accepted K1
structure. In particular, MethodHandle reference kinds 6 and 7 may target an
InterfaceMethodref only for classfile version 52 or later.

### Metadata model

The classfile package will add typed, read-only representations for:

- MethodType: descriptor index and validated method descriptor.
- MethodHandle: reference kind 1–9, reference index, and validated symbolic
  target category.
- InvokeDynamic and ConstantDynamic: bootstrap-table index plus the validated
  name and method/field descriptor category respectively.
- BootstrapMethods: ordered entries containing one MethodHandle index and an
  ordered list of loadable constant argument indexes.

New accessors will report mismatch or invalid lookup explicitly instead of
using the legacy empty-string/zero convention. They will not expose the
constant-pool backing slice. Bootstrap argument slices returned to callers are
copies so callers cannot mutate parser-owned metadata.

`ClassFile` will expose the validated BootstrapMethods view. `rtda.NewClass`
will attach that view beside the already retained constant pool, and
`rtda.Class` will expose only read-only symbolic access needed by later kernel
slices. K1 adds no Method, Field, Java facade, resolver, or bootstrap execution
state.

### File-level implementation order

| Order | Area | Intended change |
|---|---|---|
| 1 | `classfile/classreader.go`, `classfile/classfile.go` | Add checked reads, typed format errors, and the single Parse recovery/publication boundary |
| 2 | `classfile/attribute.go`, `classfile/member.go`, `classfile/stackmap.go` | Carry attribute locations and convert malformed known structures to typed failures |
| 3 | `classfile/constant_pool.go` | Decode ConstantDynamic, add typed accessors, and validate MethodType/MethodHandle/dynamic entries |
| 4 | `classfile/attribute.go` plus a focused BootstrapMethods file if useful | Parse uniqueness, entry references, arguments, and safe immutable views |
| 5 | `classfile/classfile.go`, `rtda/class.go`, `rtda/build.go` | Publish validated metadata and attach it to runtime Class without resolution |
| 6 | `transpile` only if independently valuable | Replace a generic unsupported-dynamic failure with a stable narrow diagnostic; do not alter reachability or lowering |
| 7 | package and repository gates | Run focused malformed/positive tests, fixed R3 baseline, then all contract regression commands |

The AOT diagnostic step is deliberately non-critical. It will be omitted if it
would couple K1 to K4 linkage or change capability classification.

### Test matrix

| Case | Required observation |
|---|---|
| MethodHandle kinds 1–9 | Every legal kind/tag/name combination parses; illegal combinations return typed format errors |
| MethodType | A valid method descriptor is retained; wrong tag or field descriptor is rejected |
| InvokeDynamic | NameAndType, method descriptor, and BootstrapMethods index are retained and checked |
| ConstantDynamic | NameAndType, field descriptor, and BootstrapMethods index are retained; structural dependency cycles are not resolved or rejected |
| BootstrapMethods | Zero/multiple arguments preserve order; wrong method-ref tag, invalid argument, missing table, duplicate table, and bad index fail |
| Attribute handling | Unknown attributes are skipped safely; BootstrapMethods at a non-class location fails |
| Truncation/lengths | Truncated constant, attribute, Code, and StackMapTable inputs return typed errors without bounds panic |
| Immutability | Mutating caller-received argument slices cannot change retained metadata |
| Runtime attachment | `rtda.NewClass` retains exact symbolic metadata and performs no loading, resolution, or bootstrap execution |
| Real compiler fixtures | JDK string-concat and capturing-lambda classfiles exercise InvokeDynamic, MethodType, MethodHandle, and bootstrap arguments |
| Capability honesty | The fixed 24-row R3 baseline remains unchanged and no row becomes Supported |

Synthetic byte fixtures remain the primary negative corpus because they can
pin exact malformed indexes and tags. Small compiler-produced fixtures provide
integration coverage but are not acceptance evidence for dynamic execution.

### Principal risks and controls

- **Parser regression:** checked reads touch Code and StackMapTable as well as
  K1 structures. Preserve their positive tests and add truncation cases before
  adding dynamic metadata.
- **Premature linkage:** keep cross-reference checks structural and explicitly
  test that parsing causes no provider or bootstrap activity.
- **Over-validation:** do not treat ConstantDynamic resolution cycles as parse
  errors, and apply classfile-version-dependent MethodHandle rules exactly.
- **Mutable aliases:** copy variable-length metadata at the public boundary and
  keep constant-pool storage unexported.
- **Scope leakage:** classloader propagation improvements belong to K2 and
  opcode/bootstrap execution belongs to K4; K1 tests typed failures directly
  at `classfile.Parse`.
- **Diagnostic coupling:** AOT wording is optional and follows the metadata
  work; the fixed unsupported capability result is the required gate.

### Investigation baseline

Before implementation, the focused existing suites passed:
`go test ./classfile ./rtda ./classloader ./transpile ./lowering`. The current
parser accepts compiler-produced InvokeDynamic fixtures but discards their
BootstrapMethods payload, and it cannot accept ConstantDynamic tag 17. These
observations define the K1 starting point rather than completed capability.

## Evidence / Handoff — 2026-07-19 (Round 5: name semantics fix)

The K1 candidate has completed the name-semantics rework round. All contract
gates have been run and recorded. Implementation candidate `d5ca31f` is fixed
and **ready for Owner review**, not accepted. Merge and push remain Owner
decisions.

### 1. MethodHandle name validation rules (final)

`classfile/constant_pool.go` now distinguishes three name-validation levels:

| Validator | Scope | Rules |
|---|---|---|
| `validateUnqualifiedName` | Field names (MH kinds 1–4, ConstantDynamic) | Non-empty; no `.`, `;`, `[`, `/`. Allows `<` and `>`. |
| `validateOrdinaryMethodName` | Method names (MH kinds 5/6/7/9, InvokeDynamic) | Non-empty; no `.`, `;`, `[`, `/`, `<`, `>`. Rejects `<init>`, `<clinit>`, `<foo>`. |
| Kind 8 explicit check | Constructor name (MH kind 8 only) | Must be exactly `<init>`. |

In `validateMethodHandle`:
- Kinds 1–4 (field): `validateUnqualifiedName(name)` — `<init>` and `<clinit>` are **valid** field names. No longer incorrectly rejected.
- Kinds 5/6/7/9 (method): `validateOrdinaryMethodName(name)` — rejects `<foo>`, `<init>`, `<clinit>`, and any name with `<>`.
- Kind 8 (constructor): explicit `name == "<init>"` check — unchanged.

In `validateDynamicNat`:
- `"method"` (InvokeDynamic): `validateOrdinaryMethodName(name)` — inherently rejects `<init>`, `<clinit>`, `<foo>`.
- `"field"` (ConstantDynamic): `validateUnqualifiedName(name)` + explicit `<init>`/`<clinit>` rejection per JVMS §5.4.3.6. `<foo>` is a valid field name and is **accepted**.

### 2. Nested CP structure validation (unchanged from R4)

Every nested index uses checked helpers: `lookupUtf8`, `lookupClass`,
`lookupNameAndType`. No legacy empty-string accessor used for tag detection.

### 3. Test matrix (actual case counts from `go test -v` output)

All in `classfile/classfile_test.go` using shared `buildMHFixture` /
`buildMHClassfile` / `buildCPBuf` / `buildMHClassfileNested` /
`assertParseFormatError` builders.

| Suite | Subtests | Coverage |
|---|---|---|
| `TestMHTableDrivenPositive` | 14 | All 9 kinds valid, version-52 InterfaceMethodref, major-51 Methodref, kind 8 with params |
| `TestMHFieldAcceptsInitClinit` | 8 | Kinds 1–4 × `<init>` / `<clinit>` field names with field descriptors → Parse passes |
| `TestMHTableDrivenNegativeParse` | 28 | Wrong target tag ×9, version-51 interface ×2, init/clinit ×8, kind-8 wrong name ×2, kind-8 field desc, kind-8 non-V return ×2, method-with-field-desc, field-with-method-desc |
| `TestMHMethodRejectsAngleBracketNames` | 6 | Kinds 5/6/7/9 with `<foo>`, `<bar>`, `<>` → rejected |
| `TestMHNestedIndexErrors` | 11 | Typed acc reference_index: zero, OOB, second-slot. Parse reference_index: OOB, second-slot. Parse nested wrong-tag: class_index→Utf8, nat_index→Utf8, name_index→Integer, descriptor_index→Integer. Typed acc nested wrong-tag: class_index→Utf8, nat_index→Utf8 |
| `TestMHNegativeTypedAccessor` | 4 | kind=0, index=OOB, index=0, index=second-slot |
| `TestMHKind8DescriptorRejection` | 5 | `<init>:I`, `()I`, `(I)I`, `(V)V`, `()Q` |
| `TestMHKind8DescriptorAcceptance` | 3 | `()V`, `(I)V`, `(Ljava/lang/String;)V` |
| `TestInvokeDynamicRejectsSpecialName` | 2 | `<init>`, `<clinit>` |
| `TestInvokeDynamicRejectsAngleBracketNames` | 3 | `<foo>`, `<bar>`, `<>` |
| `TestConstantDynamicRejectsSpecialName` | 2 | `<init>`, `<clinit>` |
| `TestConstantDynamicAcceptsAngleBracketFieldName` | 2 | `<foo>`, `<bar>` → accepted (valid field names) |
| `TestConstantDynamicNatWrongTag` | 1 | nat_index → Utf8 |
| `TestInvokeDynamicNatWrongTag` | 1 | nat_index → Utf8 |

Additional tests: `TestMethodTypeDescriptor` (1), `TestMethodTypeDescriptorInvalid` (1),
`TestInvokeDynamicInfo` (1), `TestInvokeDynamicInfoWrongDescriptor` (1),
`TestConstantDynamicFieldDescriptor` (1), `TestConstantDynamicRejectsMethodDescriptor` (1),
`TestConstantDynamicDecoding` (1), `TestConstantDynamicWrongTag` (1),
BootstrapMethods suite (10).

Each count is the number of `t.Run` subtests within the named top-level test
function, verified via `go test -v -run '...' -count=1 | grep '=== RUN.*/'`.
No aggregate total is stated — per-suite counts are the reproducible unit.

### 4. JDK 25 fixture actual dynamic structures (unchanged from R4)

`TestParseJDK25StringConcat` and `TestParseJDK25Lambda` in
`classfile/classfile_test.go` verify MethodType, MethodHandle, and
bootstrap arguments from real JDK-25-compiled classfiles.

### 5. Runtime loader exact call list (unchanged from R4)

`rtda/bootstrap_test.go`: `TestNoEagerBootstrapLoading`,
`TestBootstrapMethodsAttachment`, `TestBootstrapMethodsNotExecuted`,
`TestBootstrapMethodsAttachmentEmpty`, `TestSyntheticClassBootstrapMethodsNil`.

### 6. Capability baseline (24 rows)

**Harness run, not cited from PROJECT_STATUS.md.** Results:

`R3_RESULTS_DIR=<evidence-dir> bash docs/workstreams/r3-reflection-dynamic-fixtures/run-r3-baseline.sh`

| Engine | Result | Count |
|---|---|---|
| Interpreter | EXIT(1) | 24/24 |
| IR | EXIT(1) | 24/24 |
| AOT | NO-BUILD | 24/24 |

Classification: Interpreter 0/24 Match, IR 0/24 Match, AOT 24/24 NO-BUILD.
Unchanged from research baseline. No Java-visible row changed.

Raw results: `evidence/r3-k1-candidate-d5ca31f-20260719-223348/r3-baseline.txt`.
Derived statistics: `evidence/r3-k1-candidate-d5ca31f-20260719-223348/baseline-stats.txt`
(corrected 2026-07-19; machine-counted from r3-baseline.txt with documented
grep commands).

### 7. Gate results — 2026-07-19

| Gate | Command | Result |
|---|---|---|
| gofmt (this round's files) | `gofmt -l classfile/constant_pool.go classfile/classfile_test.go classfile/attribute.go classfile/classfile.go classfile/classreader.go classfile/member.go classfile/stackmap.go rtda/build.go rtda/class.go rtda/bootstrap_test.go` | **Clean** (no output) |
| gofmt (repo-wide) | `gofmt -l classfile/ rtda/` | Pre-existing issues in untouched files: `classfile/mutf8_test.go`, `rtda/field.go`, `rtda/frame.go`, `rtda/frame_test.go`, `rtda/method.go`, `rtda/monitor.go`, `rtda/monitor_test.go`, `rtda/slot.go`, `rtda/thread.go`. None are in K1's diff. |
| go vet | `go vet ./...` | Pass |
| go test | `go test ./...` | Pass |
| go test -race | `go test -race ./...` | Pass |
| Integration tests | `bash tests/run.sh` | Pass (10/10) |
| Whitespace | `git diff --check` | Pass |
| R3 baseline | `bash docs/workstreams/r3-reflection-dynamic-fixtures/run-r3-baseline.sh` | 24/24 rows complete, baseline unchanged |

### 8. Evidence location

| Artifact | Path |
|---|---|
| Name validation functions | `classfile/constant_pool.go`: `validateUnqualifiedName`, `validateOrdinaryMethodName`, `validateMethodHandle`, `validateDynamicNat` |
| Lookup helpers | `classfile/constant_pool.go`: `lookupUtf8`, `lookupClass`, `lookupNameAndType` |
| Descriptor validation | `classfile/constant_pool.go`: `validateMethodDescriptorReturnV`, `validateMethodDescriptor`, `validFieldDescriptor` |
| BootstrapMethods | `classfile/attribute.go`: `readBootstrapMethods`, `findBootstrapMethods`; `classfile/constant_pool.go`: `validateDynamicPool` |
| Parse boundary | `classfile/classreader.go`, `classfile/classfile.go` |
| Runtime attachment | `rtda/class.go`, `rtda/build.go` |
| All K1 tests | `classfile/classfile_test.go` (~2160 lines), `rtda/bootstrap_test.go` (~210 lines) |
| JDK 25 fixtures | `tests/fixtures/DynLambda.java`, `tests/fixtures/DynStringConcat.java` |
| Isolated candidate evidence | `docs/workstreams/r3-metadata-slice/evidence/r3-k1-candidate-d5ca31f-20260719-223348/` |
| Workstream doc | `docs/workstreams/r3-metadata-slice.md` |

### 9. Git diff scope and HEAD state

```
Candidate:  d5ca31f (feat(classfile): retain dynamic bootstrap metadata)
Branch:     codex/r3-metadata-slice
Base:       f685526 (acceptance anchor)
Review:     Ready; Owner acceptance pending
```

```
M classfile/attribute.go      — attrLocation, readBootstrapMethods, findBootstrapMethods
M classfile/classfile.go      — bootstrapMethods field, phase-2 validation call
M classfile/classfile_test.go  — ~2160 lines: table-driven MH + name-semantics +
                                  nested index errors + JDK 25 fixtures
M classfile/classreader.go     — parsePanic sentinel, panicf, checked reader
M classfile/constant_pool.go   — ~1000 lines: typed accessors, validateMethodHandle,
                                  validateDynamicNat, validateOrdinaryMethodName,
                                  lookup helpers, ConstantDynamic
M classfile/member.go          — attrLocation passed through readMembers
M classfile/stackmap.go        — checked reads for StackMapTable
M docs/workstreams/r3-metadata-slice.md — this file
M rtda/build.go                — NewClass attaches BootstrapMethods
M rtda/class.go                — Class exposes BootstrapMethods
A rtda/bootstrap_test.go       — runtime attachment + no-eager-loading tests
A tests/fixtures/DynLambda.java
A tests/fixtures/DynStringConcat.java
```

`classfile/mutf8_test.go` is clean — reverted to match base.

### 10. Remaining risks and blockers

- **JDK 25 version ceiling:** Existing architectural choice, not a K1 regression.
- **ConstantDynamic resolution cycles:** Intentionally not detected per JVMS.
  Deferred to K4.
- **Bootstrap argument recursion:** ConstantDynamic arguments referencing other
  ConstantDynamic entries are not recursively validated — consistent with the
  structural-only validation boundary.
- **AOT diagnostics:** Generic "unsupported dynamic" failure preserved.
  AOT remains 24/24 NO-BUILD.
- **K2/K4 coupling:** No classloader propagation or bootstrap execution code.
  Tests verify no-eager-loading but don't exercise actual bootstrap invocation.

### Owner decision required

Candidate `d5ca31f` is **Ready for review**, not accepted. The Owner must:
1. Review the name-semantics fixes and test coverage.
2. Decide whether to accept K1 and authorize integration.
3. If accepted, update K1 to Done and integrate it under the project protocol.
