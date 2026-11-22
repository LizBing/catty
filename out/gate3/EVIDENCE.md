# Gate 3 acceptance evidence

Baseline: `4860987`. Round-two fixes: `2047adf`
(`Fix Gate 3 loader dispatch, access, handler verification, and work limits`),
with test-only follow-up `d85fbbf`
(`Add work-limit refusal usability regression test`), on `codex/gate-3`. Default builds remain `catty-jvm21-gate2` / ABI `2`;
the plugin runs below use the opt-in `catty-jvm21-gate3` profile / ABI `3`.

## Toolchains

| Tool | Version |
| --- | --- |
| Go | `go1.27.1 darwin/arm64` |
| OpenJDK | Microsoft build `21.0.6+7-LTS` |
| javac / javap / jar | `21.0.6` |
| Python | `3.9.6` |

Exact text: [logs/toolchains.txt](logs/toolchains.txt).

## Commands and statuses

| Command | Status | Evidence |
| --- | --- | --- |
| `go vet ./...` | 0 | [logs/go-vet.txt](logs/go-vet.txt), [logs/go-vet-status.txt](logs/go-vet-status.txt) |
| `go test -count=1 ./...` | 0 | [logs/go-test.txt](logs/go-test.txt), [logs/go-test-status.txt](logs/go-test-status.txt) |
| `python3 -B tests/gate3.py --java-home JDK21` | 0 | [logs/full-gate.log](logs/full-gate.log) |
| `python3 -B tests/gate3.py --java-home JDK21 --prepare-only` | 0 | [logs/prepare-only.log](logs/prepare-only.log) |
| `python3 -B tests/gate2.py --java-home JDK21` | 0 | [logs/gate2-full.log](logs/gate2-full.log), [gate2/logs](gate2/logs) |
| `python3 -B tests/gate2.py --java-home JDK21 --enable-plugins` | 0 | [logs/gate2-gate3-positive.log](logs/gate2-gate3-positive.log), [gate2-gate3/logs](gate2-gate3/logs) |
| `go test ./internal/admission -run='^$' -fuzz=FuzzStageJar -fuzztime=30s -parallel=1` | 0 | [logs/fuzz.txt](logs/fuzz.txt) |

The full Gate 3 run ends with 18 `PASS` lines:

```text
PASS late-override
PASS name-lookup
PASS parent-first
PASS loader-identity
PASS replacement
PASS crossing-throw
PASS initialization
PASS deferred-linkage
PASS uncaught-plugin
PASS reflection-identity
PASS access
PASS rollback
PASS package-access
PASS package-dispatch
PASS handler-verify
PASS admission-policy
PASS standalone
PASS profile-and-packaging
Gate 3 passed.
```

Gate 2 ends `Gate 2 passed.` and includes `PASS Gate 1 regression`.
The Gate 3-enabled Gate 2 positive run ends
`Gate 2 positive corpus passed with Gate 3 enabled.` and includes
`PASS generated-go checks`; it intentionally skips the Gate 2-only
wrong-profile/negative checks and the separately-run Gate 1 suite.

## Corrected review reproducers

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

1. **Runtime loading-constraint conflict is unit-level only** for the reason
   above; there is no end-to-end Java reproducer that reaches the conflict
   through normal parent-first registration. The preparation path is exercised
   by `rt` unit tests, not by a public `Plugins.register` fixture.
2. **Work units are deterministic Catty units**, not a byte-for-byte model of
   every HotSpot verifier allocation. The boundary tests establish the profile
   contract (`Used` accepted, one less `LIMIT`), not HotSpot equivalence.
3. No end-to-end Java fixture injects a small verification-work limit (the
   production profile ceiling is 16,777,216 and would require a huge class).
   Work-limit refusal is covered at admission/verifier unit level; runtime
   rollback is covered by the `rollback` case and `rt` unit tests.
4. Handler-edge verification uses StackMapTable frames. Legacy classes without
   frame tables retain the pre-Gate-3 inference path; class-file version 65
   requires frames for the fixtures used here.
5. The exact 128/129 hierarchy-depth boundary and exact 16 MiB compressed-input
   boundary are not constructed as fixtures; the checks enforce inclusive
   ceilings and one-unit-over tests at smaller injected limits.
6. `readJar` uses Unix `O_NONBLOCK`/regular-file checks; non-Unix behavior was
   not exercised.
7. The fuzz target uses reduced test ceilings (1 MiB compressed, 4 MiB
   expanded, 128 entries, 32 classes) so it spends its budget in validation;
   production ceilings are tested separately.
8. The OpenJDK reference is Microsoft build 21.0.6, not Temurin.

## Artifact map

- [DIGESTS.txt](DIGESTS.txt): SHA-256 for saved binaries, JARs and reports.
- [reports/](reports): all case/mode reports.
- [traces/](traces): dispatch, admission, class-init, reflection, handler and
  rollback traces.
- [generated/](generated): generated Go, including `catty/plugin` adapter.
- [bin/](bin): preserved mixed binary, standalone copy and launcher.
- [jars/](jars): all plugin, malformed and boundary JARs.
- [gate2/](gate2): default Gate 2/Gate 1 logs and traces.
- [gate2-gate3/](gate2-gate3): positive Gate 2 corpus logs/traces with Gate 3
  enabled.
