# R2 initialization slice — candidate evidence

**Candidate commit:** `92e4d1f`
**Base commit:** `ecb086e`
**Branch:** `r2-init-slice`
**Date:** 2026-07-13

## Toolchain

| Tool | Version |
|---|---|
| JDK | openjdk version "25.0.3" 2026-04-21 LTS |
| javac | javac 25.0.3 |
| Go | (recorded in run-r2-results.txt header) |
| catty mode | pure-synthetic (`-no-boot`) |

## Harness

```
R2_RESULTS_DIR=docs/workstreams/r2-initialization-evidence/92e4d1f \
  bash docs/workstreams/r2-evidence/run-r2-diff.sh
```

Exit status: **0**

## Gate summary

| Gate | Command | Result |
|---|---|---|
| go vet | `go vet ./...` | Pass |
| go test | `go test ./...` | Pass |
| go test -race | `go test -race ./...` | Pass |
| R1 regression | `bash tests/run.sh` | Pass (10/10) |
| R2 differential | `run-r2-diff.sh` | 20 fixtures complete |
| git diff --check e21556a..92e4d1f | `git diff --check` | Pass (exit 0) |
| Baseline preserved | `git diff --quiet e21556a -- run-r2-results.txt matrix.md` | Pass (exit 0) |

## 13-fixture acceptance results (Owner-accepted gate)

| # | Fixture | Interp | IR | AOT | AOT classification |
|---|---|---|---|---|---|
| 1 | ClinitOrder | match | match | NO-BUILD | Not implemented |
| 2 | ClinitThrows | match | match | NO-BUILD | Not implemented |
| 3 | ConstantFieldNoInit | match | match | **match** | Supported (constant field) |
| 4 | GetstaticOwner | match | match | NO-BUILD | Not implemented |
| 5 | InterfaceDefaultInit | match | match | **NO-BUILD** | Not implemented (predecessor closure: Iface) |
| 6 | InterfaceNoInitOnImpl | match | match | NO-BUILD | Not implemented |
| 7 | InvokeStaticInit | match | match | NO-BUILD | Not implemented |
| 8 | RecursiveInitialization | match | match | NO-BUILD | Not implemented |
| 9 | SuperclassInitializationFailure | match | match | NO-BUILD | Not implemented |
| 10 | DirectInvokeStaticInit | match | match | NO-BUILD | Not implemented |
| 11 | InheritedStaticInit | match | match | NO-BUILD | Not implemented |
| 12 | SuperInitFailureNoOwnClinit | match | match | NO-BUILD | Not implemented |
| 13 | IfaceInitFailureNoOwnClinit | match | match | NO-BUILD | Not implemented |

**Interpreter: 13/13 match**
**IR: 13/13 match**
**AOT: 1 Supported, 12 Not implemented**

No runtime Go panic stack trace for any fixture.

## AOT classification detail

| Fixture | Rejection path | Classification |
|---|---|---|
| ConstantFieldNoInit | N/A (runs correctly) | Supported |
| InterfaceDefaultInit | `initClosureHasClinit` → Iface (default-bearing interface) | Not implemented |
| ClinitThrows | `initClosureHasClinit` → Bomb | Not implemented |
| DirectInvokeStaticInit | `initClosureHasClinit` → SideEffect | Not implemented |
| GetstaticOwner | `initClosureHasClinit` → Base | Not implemented |
| InheritedStaticInit | `initClosureHasClinit` → Ancestor | Not implemented |
| InvokeStaticInit | `initClosureHasClinit` → Holder | Not implemented |
| RecursiveInitialization | `initClosureHasClinit` → RecursiveInit | Not implemented |
| ClinitOrder | unsupported opcodes in main | Not implemented |
| InterfaceNoInitOnImpl | unsupported opcodes in main | Not implemented |
| SuperclassInitializationFailure | unsupported opcodes in main | Not implemented |
| SuperInitFailureNoOwnClinit | unsupported opcodes in main | Not implemented |
| IfaceInitFailureNoOwnClinit | unsupported opcodes in main | Not implemented |

## Predecessor-closure direct evidence

### Go unit tests (TestInitClosureHasClinit)

Location: `transpile/build_test.go`

| Case | Assertion | Result |
|---|---|---|
| superclass-has-clinit | `initClosureHasClinit(Child)` returns `"Parent"` | PASS |
| interface-has-clinit | `initClosureHasClinit(Impl)` returns `"Iface"` | PASS |
| no-clinit | `initClosureHasClinit(Child)` returns `""` | PASS |
| interface-ignores-superinterface | `initClosureHasClinit(SubIface)` returns `""` | PASS |

### AOT harness evidence (InterfaceDefaultInit)

`InterfaceDefaultInit` is a fixture where:
- `Impl` implements `Iface`, Impl has no `<clinit>`
- `Iface` has both a default method (`m()`) and `<clinit>` (from `int X = mark()`)
- AOT build rejects with: *"triggers class initialization on Iface"*

This proves the predecessor closure check correctly identifies `<clinit>` on
a default-bearing superinterface whose implementing class has no `<clinit>`
of its own. The old (bb02216) check only examined the immediate target and
would have emitted this fixture without an init guard.

### Test helpers (rtda/build.go)

Minimal additions for synthetic interface/class construction:
- `MarkInterface()` — sets ACC_INTERFACE flag
- `AddInterface(iface)` — adds a direct superinterface
