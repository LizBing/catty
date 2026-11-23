# Gate 3 acceptance evidence

Review baseline `6a18827` on `codex/gate-3`. Final implementation:
`be33262` (`Fix Gate 3 pending constraint merges and preparation staging`).
This round is limited to the two loading-constraint gaps: rejecting a conflict
when two pending equality groups are joined, and staging pending equalities
during plugin preparation. The round-three inherited-interface dispatch, exact
hierarchy-depth and verification-budget fixes and their reproducers are
preserved; all Gate 3, Gate 2, Gate 1, Go and fuzz checks were re-run on the
final commit. Earlier round evidence remains below as historical context.
Default builds stay `catty-jvm21-gate2` / ABI `2`; plugin runs use the opt-in
`catty-jvm21-gate3` profile / ABI `3`.

## Toolchains

| Tool | Version |
| --- | --- |
| Go | `go1.27.1 darwin/arm64` |
| OpenJDK | Temurin `21.0.13.1+2` (LTS) |
| javac / javap / jar | `21.0.13.1` |
| Python | `3.9.6` |

Exact text: [round4/logs/toolchains.txt](round4/logs/toolchains.txt).

## Commands and statuses

| Command | Status | Evidence |
| --- | --- | --- |
| `go test ./rt -run 'TestLoadingConstraintMergedGroups|TestLoadingConstraintStageRollsBack|TestPluginPreparation|TestOverridePreparation' -count=1 -v` | 0 | [round4/logs/focused-tests.txt](round4/logs/focused-tests.txt) |
| `go vet ./...` | 0 | [round4/logs/go-vet.txt](round4/logs/go-vet.txt) |
| `go test -count=1 ./...` | 0 | [round4/logs/go-test.txt](round4/logs/go-test.txt) |
| `python3 -B tests/gate3.py --java-home JDK21` | 0 | [round4/logs/gate3-full.log](round4/logs/gate3-full.log) |
| `python3 -B tests/gate2.py --java-home JDK21` | 0 | [round4/logs/gate2-full.log](round4/logs/gate2-full.log) |
| `python3 -B tests/gate2.py --java-home JDK21 --enable-plugins` | 0 | [round4/logs/gate2-gate3-positive.log](round4/logs/gate2-gate3-positive.log) |
| `go test ./internal/admission -run='^$' -fuzz=FuzzStageJar -fuzztime=30s -parallel=1` | 0 | [round4/logs/fuzz.txt](round4/logs/fuzz.txt) |

The full Gate 3 run contains 17 Java cases and 3 policy/standalone/profile
checks (20 `PASS` lines):

```text
PASS late-override
PASS name-lookup
PASS parent-first
PASS loader-identity
PASS handler-verify
PASS package-access
PASS package-dispatch
PASS inherited-interface
PASS hierarchy-depth
PASS replacement
PASS crossing-throw
PASS initialization
PASS deferred-linkage
PASS uncaught-plugin
PASS reflection-identity
PASS access
PASS rollback
PASS admission-policy
PASS standalone
PASS profile-and-packaging
Gate 3 passed.
```

Gate 2 ends `Gate 2 passed.` and includes `PASS Gate 1 regression`. The
Gate 3-enabled Gate 2 positive run ends
`Gate 2 positive corpus passed with Gate 3 enabled.` The fuzz target completed
1,299,992 executions in 30 seconds without a panic or published loader.

## Round-four corrected behaviors

### 1. Pending equality-group merges are validated as a whole

`rt/constraints.go` now stages equality edges in a `constraintTx` and validates
the complete equality component after every edge addition instead of rechecking
only the new endpoints. Joining two previously disjoint pending groups is
refused when already-resolved members of the merged component disagree, even
when both new endpoints are still unresolved. A refused staging operation leaves
the transaction unchanged, so an operation that adds several descriptor leaves
publishes none of them on conflict. `recordInitiation` uses the same component
validation, so later lookups keep lazy JVMS 5.3.4 behavior and still produce
`LinkageError` on a conflict.

The reviewed scenario is exactly A initiates `X1`, D initiates a distinct `X2`
with the same binary name, `A = B` and `C = D` are pending, and `B = C` joins the
groups. On the baseline it grows the graph from two edges to three; on `be33262`
it is refused:

- [round4/logs/baseline-repro.txt](round4/logs/baseline-repro.txt) records the
  identical runtime test failing on `6a18827` and passing on `be33262`.

Focused tests:

- `TestLoadingConstraintMergedGroupsConflict` — the reviewed scenario; the
  refused merge publishes no edge, leaves both initiation records as they were,
  and B can still resolve to A's class.
- `TestLoadingConstraintMergedGroupsSameClass` — a compatible merge where both
  groups resolve to the same class publishes three edges; a disagreeing late
  initiation is refused and agreeing late initiations complete the component.
- `TestLoadingConstraintStageRollsBackWholeDescriptor` — a conflict on a later
  descriptor type leaves the edge staged for an earlier type unpublished.
- The existing pending, transitive-component, shared-class, array-leaf and
  primitive-descriptor tests continue to pass.

### 2. Plugin preparation stages pending equalities in the candidate transaction

`rt/plugin.go` no longer resolves descriptor leaves eagerly in
`checkOverrideConstraints`. It stages one equality for every reference leaf of
every declared override and inherited interface implementation, using the same
transaction staging validation as member resolution. Pending equalities are kept
without forcing a descriptor leaf lookup or class initialization; a conflict
already visible from initiation records is reported as `VERIFY` and rejects the
candidate. `installPluginTx` commits the transaction only at the publication
point, after the last fallible preparation step: a rejected candidate discards
every pending equality it staged, while a successful admission publishes them
with the loader definitions and quota accounting. A conflict during later
resolution still surfaces as the existing sticky `LinkageError` through
`recordInitiation`.

Focused tests (runtime unit tests use manually linked loaders because fixed
parent-first visibility makes conflicting mixed-loader configurations
unreachable through the Java API):

- `TestPluginPreparationStagesPendingEqualities` — the real `installPlugin`
  preparation path publishes pending reference and reference-array equalities,
  adds no edge for primitive descriptors, and a pending equality is enforced
  when a later initiation resolves the other side to a different class.
- `TestPluginPreparationPendingEqualityDiscardedOnRejection` — a candidate whose
  first class stages a pending equality and whose later class conflicts is
  rejected with `VERIFY`; constraints, initiations, registry and loader ID stay
  unchanged, and a control admission of the first class alone publishes the
  pending equality.
- `TestOverridePreparationUsesLeafTypes` and
  `TestOverridePreparationInheritedInterface` now record the conflicting
  initiations explicitly and assert that the rejected preparation transaction
  publishes no edges.

## Round-three review fixes

### 1. Inherited interface method selection

`rt/thread.go` now separates the receiver-scoped JVMS 5.4.6 selection rule from
preparation's stronger override test. `selectVirtual` searches the receiver and
its superclass chain with `selectionOverrides`, the JVMS 5.4.5 can-override
relation, so a public superclass implementation satisfies an interface
declaration even when that superclass does not implement the interface.
Package-private class declarations still require same-runtime-package or
transitive overriding, and the compiled dispatch cache still requires the
identical resolved declaration and exact receiver `*Class`. Reflection uses the
same general selector. `methodOverrides` (which additionally requires the
declaring types to be related) is retained for override-preparation loading
constraints. `invokeinterface` also rejects a selected method that is neither
public nor private with `IllegalAccessError`, matching JVMS 6.5.

New differential fixture `testdata/gate3/inherited-interface`: image
`api.Base.value()` returns `31`, image `api.Step` declares `value()`, plugin
`ext.Child extends api.Base implements api.Step` inherits `value()` and exposes
`static Step create()`. The runner builds the image before packaging the JAR
and runs all three image modes.

| Implementation | Result |
| --- | --- |
| OpenJDK 21 oracle (`-Xverify:all`) | status 0, stdout `31` four times |
| Catty `4860987` | status 0, stdout `31` four times |
| Catty `6b77273` | status 1, `Exception in thread "main" java.lang.AbstractMethodError` |
| Catty `57332b8` interp / mixed / aot | status 0, stdout `31` four times in every mode |

The dispatch trace from the final interp run names the inherited image
superclass as the target:

```json
{"event":"dispatch","opcode":"invokeinterface","declaration":{"owner":"api/Step","name":"value","loaderId":1},"receiverLoaderId":2,"target":{"owner":"api/Base","name":"value","loaderId":1},"cache":{"miss":true}}
```

The third `31` is the reflective `Step.class.getMethod("value").invoke(child)`
call; the first two are compiled `invokeinterface` sites. Exact commands and
outputs: [round3/inherited-interface-comparison.txt](round3/inherited-interface-comparison.txt);
traces: [round3/traces/](round3/traces).

### 2. Exact hierarchy depth

New shared package `internal/hierarchy` computes the longest path through the
effective superclass/interface graph with a memoized depth-first traversal. It
detects cycles, charges one unit per newly visited node through a callback, and
returns `DepthError` when a root's longest path exceeds the inclusive ceiling.
Memoization keeps the full remaining path length for every ancestor, so visit
order (including ancestor-first archive order) no longer affects enforcement.

Admission applies the rule to effective definitions and reads parent-loader
metadata; runtime revalidation applies the same rule to the candidate loader's
effective graph through `Loader.Lookup`. A path spanning a published parent
loader counts every edge. Depth failures keep code `LIMIT` with
`limit=hierarchyDepth`; cycles remain `VERIFY`.

The production-boundary fixture generates five JARs (386 class files) with
ancestor-first archive order:

| JAR | Definition | Result |
| --- | --- | --- |
| `big128.jar` | `chain.C000..C127`, 128-edge longest path | accepted, loader 2, `depth()` runs and returns 128 |
| `big129.jar` | `chain.C000..C128`, 129-edge longest path | rejected `LIMIT/hierarchyDepth` at `chain/C128`, no loader |
| `parent1.jar` | `span.P000..P119`, 120-edge path | accepted, loader 3 |
| `parent2.jar` | `span.P120..P127` extending parent `P119`, combined 128-edge path | accepted, loader 4 |
| `parent3.jar` | `span.P128` extending parent `P127`, combined 129-edge path | rejected `LIMIT/hierarchyDepth` at `span/P128`, no loader |

Final run output, identical in interp, mixed and aot images:

```text
128
120
128
LIMIT
LIMIT
true
```

Trace events:
[round3/hierarchy-depth-boundary.txt](round3/hierarchy-depth-boundary.txt) and
[round3/traces/hierarchy-depth-interp.jsonl](round3/traces/hierarchy-depth-interp.jsonl).
Go coverage includes the exact 128/129 boundary in both orders, interface
paths, cycles, admission with synthetic class files, and runtime revalidation
with a parent-spanning chain.

### 3. Verification work refusals stop the work

`internal/verify/types.go` now returns a `WorkError` from `copyState` and
`statesEqual` before copying or comparing, and `assignable`,
`classAssignable`, `classIsInterface`, `classImplements`, `interfaceExtends`
and `isThrowable` all propagate a refused charge instead of continuing the
query. `applyStep` records the first successor-copy refusal and returns it
before exposing partial successors. `analyzeMethod` charges handler-range
inspection as well as matching exception edges, charges frame allocations and
copies, and does not publish heights/shapes when any step is refused.
Admission still classifies the propagated `WorkError` as `LIMIT` with
`limit=verificationWork` and never publishes candidate definitions.

New tests in `internal/verify/budget_test.go` exercise a refused copy, a
refused comparison, a refused hierarchy visit at an exact step, and a sweep of
every verifier budget boundary: each failure names the refused operation, does
not consume the refused units, stops at the expected step, and leaves the
analysis maps empty. `TestWorkLimitRefusalLeavesRunningLoaderUsable` (existing)
continues to prove the runtime remains usable after a work-limit admission
refusal.

### 4. Pending and transitive loading constraints

`rt/constraints.go` now stages an equality edge for every reference leaf type
without resolving the descriptor leaf. When either loader has not initiated the
leaf yet the edge is pending; `recordInitiation` walks the full connected
component of equality edges for that binary name and compares every
already-resolved peer, including peers reached through unresolved endpoints.
Updates are built on copies and published atomically only after the whole
staging operation succeeds. Reference-array leaves use the element class name,
and primitive/void descriptors create no edge.

The previous test that manually inserted into the private constraint map was
replaced with staging-API coverage: pending equality then later initiation,
transitive `A = B = C` conflict through an unresolved endpoint, successful
shared parent-visible class, failed staging/initiation preserving prior state,
array leaves, and primitive descriptors outside the graph. Runtime collection is
still unit-level because fixed parent-first visibility makes a conflicting
public-API setup unreachable; the preparation path remains covered by
`TestOverridePreparationUsesLeafTypes` and
`TestOverridePreparationInheritedInterface`.

## Round-two corrected review reproducers (preserved)

### 1. Package-private virtual dispatch across defining loaders

Fixture `testdata/gate3/package-dispatch`:

- public image `demo.Base` has package-private `value()` returning `11`
  and public `call()` invoking `value()`;
- plugin `demo.Child` extends `Base` and declares public `value()` returning
  `22`;
- public image `demo.PublicBase.value()` returning `11` and plugin
  `demo.PublicChild` returning `22` retain the legitimate override.

The same `Base.call()` site is executed with image and plugin receivers in
interpreted, mixed and forced-AOT builds. After the fix the trace shows:

```json
{"caller":{"owner":"demo/Base","name":"call","loaderId":1},"declaration":{"owner":"demo/Base","name":"value"},"receiverLoaderId":1,"target":{"owner":"demo/Base","name":"value"},"cache":{"miss":true}}
{"caller":{"owner":"demo/Base","name":"call","loaderId":1},"declaration":{"owner":"demo/Base","name":"value"},"receiverLoaderId":2,"target":{"owner":"demo/Base","name":"value"},"cache":{"miss":true}}
{"caller":{"owner":"demo/PublicBase","name":"call"},"declaration":{"owner":"demo/PublicBase","name":"value"},"receiverLoaderId":2,"target":{"owner":"demo/PublicChild","name":"value"}}
```

The plugin receiver still enters `demo/Base.call` but dispatch selects
`demo/Base.value` (11); the public override selects `demo.PublicChild.value`
(22). OpenJDK 21 produces the same `11`, `11`, `11`, `22` output. Selection now
uses the JVMS override relation (`methodOverrides`) including runtime package
identity and transitive overriding; the cached compiled site remains guarded by
the resolved declaration and exact receiver class.

### 2. Runtime package access during admission

Fixture `testdata/gate3/package-access`:

- plugin `demo.Child` extends package-private image `demo.HiddenParent`;
- plugin `demo.Impl` implements package-private image `demo.HiddenIface`;
- plugin `demo.LocalChild` extends a package-private parent defined in the
  same plugin loader, and `demo.LocalImpl` implements a same-loader
  package-private interface.

OpenJDK 21 reports `IllegalAccessError` for both cross-loader classes. Catty
rejects both with `VERIFY` and retains no loader; it then registers the
same-loader control and executes it. Trace:

```json
{"event":"plugin_admission","result":"rejected","code":"VERIFY","class":"demo/Child","entry":"demo/Child.class","message":"superclass is not accessible from the plugin loader: demo/HiddenParent"}
{"event":"plugin_admission","result":"rejected","code":"VERIFY","class":"demo/Impl","entry":"demo/Impl.class","message":"interface is not accessible from the plugin loader: demo/HiddenIface"}
{"event":"plugin_admission","result":"accepted","definedClasses":4,"loaderId":2}
```

`internal/admission/checks.go` now resolves the effective parent-first
definition and treats a package-private target from a different defining loader
as inaccessible even when the textual package matches. Go tests cover both
rejections, the same-loader control, and the failure path; the runtime quota
and rollback unit tests cover registry/array-cache restoration.

### 3. Exception-handler input frames

`internal/verify/verify.go` now checks every handler edge whose protected range
covers the instruction: the incoming locals and exactly the caught throwable
(or `Throwable` for catch-all) must be assignable to the handler stack-map
frame. Shared targets and uninitialized values are covered by focused Go tests
(`internal/verify/handler_test.go`).

Hand-assembled fixtures (`testdata/gate3/handler-frame/generate.py`) provide a
valid handler and two invalid ones. OpenJDK 21:
`good.Handler` loads and returns `2`; `bad.Handler` fails with
`VerifyError: Stack map does not match the one at exception handler 16`;
`badstack.Handler` also fails with `VerifyError: Stack map does not match the one at exception handler 16`.
Catty rejects both at registration and the valid class registered after the
rejections still executes:

```json
{"event":"plugin_admission","result":"rejected","code":"VERIFY","class":"badstack/Handler","method":"f","descriptor":"(I)I","pc":16,"message":"badstack/Handler.f(I)I: handler pc 16: stack map declares int, not a Throwable"}
{"event":"plugin_admission","result":"rejected","code":"VERIFY","class":"bad/Handler","method":"f","descriptor":"(I)I","pc":8,"message":"bad/Handler.f(I)I: exception edge at pc 8 to handler pc 16: local 1 is java/lang/Object, not assignable to java/lang/String"}
```

`ClassifyVerifyError` recovers the method, descriptor and bytecode offset from
the verifier's stable diagnostic prefix, so the Java registration trace carries
the method and location.

### 4. Verification work accounting

`internal/verify/work.go` charges deterministic units before frame storage,
frame copies/comparisons, instruction and control-flow edges, switch successors,
exception edges, hierarchy/interface visits and class-file table elements.
`internal/admission` charges the class/field/method/ref/instruction/handler
tables and hierarchy visits before verification. Work is staged in
`AdmissionResult.WorkUsed`; the diagnostic is `LIMIT` with
`limit=verificationWork`.

`TestVerificationWorkAccountingExceedsOldEstimate` computes the previous
estimate for `valid-stored.jar`, proves the new `WorkUsed` is larger, then
proves exact acceptance at `WorkUsed` and `LIMIT` at `WorkUsed-1`. The same test
re-runs a normal admission after the refusal. `TestVerifierWorkExactBoundary`
does the same at the verifier level. `TestWorkLimitRefusalLeavesRunningLoaderUsable`
injects a LIMIT admission refusal after a plugin loader is published and proves
the existing loader, class object and mirror remain usable; the rollback
fixture separately proves a rejected registration leaves the running image and
prior plugin usable.

### 5. Loading constraints

`rt/constraints.go` tracks loader-specific initiations and equality constraints
for reference leaf types (including nested array leaves), with staged updates
committed only on success, sticky resolution errors and `LinkageError` on a
conflict. Member resolution imposes constraints between the caller loader and
the declaring loader for field and method descriptors; class/array lookups
record initiations and check previously registered constraints. Override
preparation walks the superclass chain and complete interface closure, compares
leaf types for actual overrides, and rejects a conflict before commit.

Go coverage:

- `TestLoadingConstraintTransition` — conflict leaves both maps unchanged;
  a later lookup observes a pending constraint and is not published.
- `TestLoadingConstraintPositiveAndArrayLeaf` — shared parent-visible leaves
  and shared array leaves commit; loader-specific initiations are distinct.
- `TestLoadingConstraintArrayLeafConflict` — nested array leaf conflict.
- `TestOverridePreparationUsesLeafTypes` — preparation-time leaf conflict.
- `TestOverridePreparationInheritedInterface` — inherited interface
  implementation conflict.
- `TestDeferredDescriptorLeavesArrayCacheUntouched` — unresolved descriptor
  types stay deferred and do not pollute the array cache.

**Invariant:** because lookup is strictly parent-first, two loaders that share
a parent-visible name resolve the same class, so a proposed runtime conflict
between such loaders is unreachable through the public API. The conflict
transition is therefore covered at the runtime unit level, and the invariant is
documented here rather than exercised by an artificial Java fixture. An absent
class in the parent's view keeps ordinary resolution timing; an unresolved
descriptor type is retained lazily (see the new `ext.Sig` fixture in
`deferred-linkage`) rather than rejected at registration.

## Other acceptance points preserved

- Every image is built before its plugin JAR exists; the runner asserts the JAR
  path is absent at all three image builds and that reports contain no plugin
  class/JAR names.
- The late-override fixture still proves image target, interpreted plugin
  target, second loader target and cache hits/misses at the same bytecode site.
- Generated mixed/AOT Go still passes the original PC and CP entry to
  `InvokeVirtual`/`InvokeInterface`; the standalone run uses empty `PATH` and
  unset `JAVA_HOME`/`CLASSPATH` with the same output and trace relationships.
- `irem` auto fallback and forced-AOT rejection remain checked.
- Admission policy variants, limits, rollback and report inventory remain
  checked.

## Remaining gaps / evidence limits

1. **Runtime loading-constraint conflicts remain unit-level only.** With fixed
   parent-first visibility, two published loaders that share a parent-visible
   name resolve the same class, and an effective plugin class may not shadow a
   parent definition. The reviewed conflicting pending-merge and preparation-time
   conflict transitions therefore use manually linked runtime loaders; no Java
   fixture can reach them. Pending equalities are exercised through the real
   `installPlugin` preparation path, and execution-time conflicts remain covered
   by the existing `LinkageError` unit tests.
2. **Published preparation equalities are latent.** A pending constraint
   committed by a successful admission is only enforced when the affected class
   name is resolved later, so no end-to-end fixture observes an immediate
   Java-visible effect at registration time. The focused tests assert the
   published transaction state directly and then exercise the later initiation.
3. **Work units are deterministic Catty units**, not a byte-for-byte model of
   HotSpot verifier allocations. The tests establish the profile contract
   (`Used` accepted, one less `LIMIT`) and the stop-before-work behavior, not
   HotSpot equivalence.
4. No end-to-end Java fixture injects a small verification-work limit (the
   production ceiling is 16,777,216 and would require a very large class).
   Refusal is covered at the verifier and admission unit boundaries; runtime
   usability after refusal is covered by the existing `rt` test and the
   `rollback` case.
5. Handler-edge verification uses StackMapTable frames. Legacy classes without
   frame tables retain the pre-Gate-3 inference path; class-file version 65
   requires frames for these fixtures.
6. An interface class file names `java/lang/Object` as its `super_class`, so the
   depth traversal counts that artifact edge for interfaces exactly as the
   previous admission/runtime traversals did. The interface-path unit test
   accounts for it; the end-to-end boundary fixture uses class chains.
7. `readJar` uses Unix `O_NONBLOCK`/regular-file checks; non-Unix behavior was
   not exercised.
8. The fuzz target uses reduced test ceilings (1 MiB compressed, 4 MiB
   expanded, 128 entries, 32 classes) so it spends its budget in validation;
   production ceilings are tested separately.
9. The OpenJDK reference is Temurin `21.0.13.1+2`.
10. The package-private inherited-interface edge case is not a valid Java-source
    program; a hand-crafted class file now raises `IllegalAccessError` at the
    invokeinterface site like OpenJDK, but no end-to-end fixture builds that
    invalid class.

## Artifact map

- [DIGESTS.txt](DIGESTS.txt): round-two SHA-256 manifest.
- [reports/](reports), [traces/](traces), [generated/](generated),
  [bin/](bin), [jars/](jars): round-two artifacts.
- [gate2/](gate2): default Gate 2/Gate 1 logs and traces.
- [gate2-gate3/](gate2-gate3): positive Gate 2 corpus logs/traces with Gate 3
  enabled.
- [round3/DIGESTS.txt](round3/DIGESTS.txt): SHA-256 for the round-three JARs,
  reports, traces and generated Go.
- [round3/logs/](round3/logs): final command logs and new-case build/run/oracle
  logs.
- [round3/traces/](round3/traces): inherited-interface and hierarchy-depth
  dispatch/admission traces in all three modes.
- [round3/reports/](round3/reports): the corresponding build reports.
- [round3/generated/](round3/generated): generated Go for the new fixtures.
- [round3/jars/](round3/jars): inherited-interface and hierarchy-boundary JARs.
- [round3/inherited-interface-comparison.txt](round3/inherited-interface-comparison.txt):
  OpenJDK, `4860987`, `6b77273` and final-commit comparison.
- [round3/hierarchy-depth-boundary.txt](round3/hierarchy-depth-boundary.txt):
  exact 128/129 acceptance/rejection and trace fields.
- [round4/logs/](round4/logs): round-four toolchains, focused constraint and
  preparation tests, review-baseline reproduction, full Gate 3/Gate 2/Go/vet and
  fuzz logs. No new Java fixtures were added this round; the existing Gate 3
  case set is unchanged.
