# Gate 3 acceptance evidence

Gate 3 was implemented on `codex/gate-3` and verified on the remote
`darwin/arm64` host. Default builds remain the Gate 2 profile
(`catty-jvm21-gate2`, ABI `2`); every result below is from the opt-in
`catty-jvm21-gate3` profile with ABI `3`.

## Toolchains

| Tool | Version |
| --- | --- |
| Go | `go1.27.1 darwin/arm64` |
| OpenJDK | Microsoft build `21.0.6+7-LTS` |
| javac / javap / jar | `21.0.6` |
| Python | `3.9.6` |

Exact text: [logs/toolchains.txt](logs/toolchains.txt).

## Commands and results

| Command | Status | Evidence |
| --- | --- | --- |
| `go vet ./...` | 0 | [logs/go-vet.txt](logs/go-vet.txt), [logs/go-vet-status.txt](logs/go-vet-status.txt) |
| `go test -count=1 ./...` | 0 | [logs/go-test.txt](logs/go-test.txt), [logs/go-test-status.txt](logs/go-test-status.txt) |
| `python3 -B tests/gate3.py --java-home ... --prepare-only` | 0 | prepares all 12 cases and the OpenJDK oracle without building Catty |
| `python3 -B tests/gate3.py --java-home ...` | 0 | [logs/full-gate.log](logs/full-gate.log) |
| `go test ./internal/admission -run='^$' -fuzz=FuzzStageJar -fuzztime=30s -parallel=1` | 0 | [logs/fuzz.txt](logs/fuzz.txt) |
| `python3 -B tests/gate2.py --java-home ...` (includes Gate 1 regression) | 0 | [logs/gate2-full.log](logs/gate2-full.log), [gate2/logs](gate2/logs) |

The full Gate 3 log ends with:

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
PASS admission-policy
PASS standalone
PASS profile-and-packaging
Gate 3 passed.
```

Gate 2 ends with `Gate 2 passed.`; its embedded Gate 1 regression reports
`PASS Gate 1 regression`.

## What the run proves

- The application executable for each case is built before its plugin JAR
  paths exist. The runner asserts the JAR path is absent at each of the three
  image builds, and that the report contains no plugin class/entry names.
- Every positive case is executed from the same class files and JAR bytes by
  three images: all-interpreted, mixed (policy forces `app/Main.main` to the
  interpreter while the wrappers and `Box.add` stay AOT), and forced-AOT.
  OpenJDK 21 with `-Xverify:all` is the oracle; stdout bytes and empty stderr
  are compared exactly.
- The late-override fixture uses one `invokevirtual` site and one
  `invokeinterface` site in image wrappers. The trace first shows the image
  target, then an interpreted plugin target, then the second plugin loader,
  then the original plugin target again. Cache transitions include both misses
  and hits at each site. See
  [traces/late-override-mixed.jsonl](traces/late-override-mixed.jsonl).
- The generated mixed/aot Go passes the original bytecode offsets and
  constant-pool indexes into the shared runtime, for example:

  ```go
  out := t.InvokeVirtual(m, 2, 7, stack[0], "app/Box", "add", "(I)I", stack[1:2])
  out := t.InvokeInterface(m, 2, 13, stack[0], "app/Step", "value", "(Lapp/Box;)I", stack[1:2])
  ```

  Source: [generated/late-override-mixed/generated.go](generated/late-override-mixed/generated.go).
- The Gate 3 report/API pair advertises:

  ```text
  profile=catty-jvm21-gate3 imageABI=3 plugins=true
  capabilities=plugin-admission,loader-identity,guarded-dispatch
  bootstrap: 50 classes, 116 methods
  limits: 16 MiB compressed, 64 MiB expanded, 4 MiB/entry, 4096 entries,
          1024 classes, depth 128, 16,777,216 work, 64 loaders, 128 MiB classes
  ```

  Report: [reports/late-override-mixed.json](reports/late-override-mixed.json).
- `irem` still falls back to the interpreter in auto mode and is rejected by
  forced AOT under the Gate 3 profile (runner `profile-and-packaging` check).
- The standalone check copies the mixed executable and late JARs into an empty
  directory and runs them with empty `PATH`, unset `JAVA_HOME`/`CLASSPATH`, and
  an absolute path. Stdout, status and trace relationships match, and the
  executable digest is unchanged. See [logs/standalone.json](logs/standalone.json)
  and [traces/standalone.jsonl](traces/standalone.jsonl).

### Admission and rollback evidence

The Java registration intrinsic is exercised directly for each policy
variant. Representative variants and stable codes:

| Variant | Code |
| --- | --- |
| non-ZIP / truncated central directory / truncated entry / bad CRC | `ARCHIVE` |
| duplicate entry / file-directory conflict / empty class set | `ARCHIVE` |
| dot-dot / backslash / absolute / drive-prefixed names | `ARCHIVE` |
| encrypted / unsupported compression | `ARCHIVE` |
| entry/name mismatch / malformed class | `CLASS_FORMAT` |
| unsupported major version 66 | `PROFILE` |
| class under `META-INF` / module-info / multi-release / versioned path | `PACKAGE` |
| nonempty `Class-Path` / sealed / signature metadata / reserved `catty/runtime` | `PACKAGE` |
| missing superclass from the image | `DEPENDENCY` |
| entry size above 4 MiB, class count, entry count, expanded total | `LIMIT` |

The rollback case prints `true` for the `CLASS_FORMAT` code, keeps the warmed
image and previously published plugin working, cannot find the rejected
`ext/Alpha` definition through the application or old plugin loader, then
registers corrected bytes and executes `ext/Alpha` in the same process and
logical thread. The rejection event has no `loaderId`; the later acceptance has
a fresh loader.

Trace excerpt:

```json
{"code":"CLASS_FORMAT","entry":"ext/Broken.class","event":"plugin_admission","jarSHA256":"...","message":"bad class file magic 0x00000000","parentLoaderId":1,"result":"rejected","thread":1}
```

Source: [traces/rollback-mixed.jsonl](traces/rollback-mixed.jsonl).

### Loader identity and cache scope

Go tests prove sibling loaders define the same name independently, arrays are
interned by immediate component class, primitive arrays are shared, runtime
package equality includes the defining loader, failed resolution outcomes are
sticky per referring runtime class and CP entry, failed admission leaves loader
IDs/quotas/registry unchanged, and plugin methods run in the interpreter. See
[rt/gate3_test.go](../../rt/gate3_test.go) and
[internal/admission/admission_test.go](../../internal/admission/admission_test.go).

## Artifact manifest

- [DIGESTS.txt](DIGESTS.txt): SHA-256 for every saved binary, JAR and report.
- [reports/](reports): all build reports for the 12 cases in three modes plus
  the policy app and `irem` fixture.
- [traces/](traces): 37 JSON-lines traces covering accepted/rejected admission,
  class definition/initialization, dispatch cache transitions, reflection,
  crossing and rollback.
- [generated/](generated): generated Go for every mixed/aot image (including
  `main.go` with the `catty/plugin` adapter import and `SetAdmission`).
- [bin/late-override-mixed](bin/late-override-mixed): the preserved mixed
  executable used by the standalone run (also copied as
  [bin/standalone-app](bin/standalone-app)).
- [jars/](jars): the plugin and malformed/policy JARs used by the run.
- [gate2/](gate2): Gate 2/Gate 1 regression logs and traces.

## Remaining unverified or approximate checks

These are reported rather than claimed:

1. **Loading-constraint table.** The runtime does not maintain the full JVMS
   5.3.4 initiating-loader/equality-constraint table for every field and method
   descriptor. It detects conflicting reference descriptor types for an
   override during plugin preparation and rejects them as `VERIFY`; a runtime
   conflict arising later from an unrelated initiating loader is not modeled.
2. **Verification-work accounting.** The admission ceiling is charged from
   instruction/stack/frame counts before verification; it is an aggregate
   bound, not an exact reproduction of every internal verifier work unit.
3. **Exact hierarchy-depth boundary.** The 128-edge ceiling is enforced and
   tested structurally at low injected limits, but no fixture exercises exactly
   128 versus 129 class/interface edges.
4. **Compressed-input exact boundary.** Entry/expanded/count boundaries are
   tested at and one unit beyond the relevant ceiling. The exact 16 MiB
   compressed-input boundary is not constructed as a fixture; the runtime
   rejects over-cap input before copying.
5. **Encrypted/unsupported-compression fixtures** are produced by patching ZIP
   headers/flags; they prove the policy path, not every possible central
   directory malformation (the fuzz target covers additional shapes).
6. **File-handle blocking.** `readJar` opens with `O_NONBLOCK` and verifies a
   regular file before reading on Unix; non-Unix behavior was not exercised.
7. **Fuzz ceilings.** The fuzz target uses reduced test limits (1 MiB
   compressed, 4 MiB expanded, 128 entries, 32 classes) so the target spends
   its budget in validation rather than in huge ZIP allocations. Production
   limits are separately tested at their entry/expanded/count boundaries.
8. **Oracle JDK.** The oracle is Microsoft OpenJDK 21.0.6, not the Temurin
   build used for Gate 2; both are JDK 21 and the tests pin the tool major.
