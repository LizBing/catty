# Slice A — Candidate Evidence

**Candidate:** `9576828`
**Baseline:** `a0288be` (acceptance anchor, 2026-07-14 governance commit)
**Worktree base:** `a0288be`
**Slice scope:** SC heap cells, concurrency-safe loader, canonical Class mirrors
**Date:** 2026-07-14

## Scope summary (relative to baseline)

22 files changed, 1306 insertions, 259 deletions.

| File | Change |
|---|---|
| `rtda/cell.go` | HeapCell typed getters/setters; ToSlot panics on long/double; CopyHeapCells; NewHeapCells |
| `rtda/object.go` | Removed `Cells()`; added typed per-cell accessors; CopyObjectCells (overlap-safe, kind-dispatch); CloneObject; CloneCells |
| `rtda/class.go` | Removed `StaticCells()`; added typed static accessors; StaticCellToSlot; ClassObject (CAS-once, canonical mirrors) |
| `rtda/init.go` | Migrated to typed accessors |
| `rtda/build.go` | SetStaticRef (existing, used by native); componentKind support |
| `interpreter/fields.go` | loadInstanceField/storeInstanceField/loadStaticField/storeStaticField — all dispatch on desc[0] to typed accessors |
| `interpreter/interpreter.go` | Migrated all .Cells()/.StaticCells() calls to typed accessors |
| `interpreter/ir.go` | Same migration |
| `interpreter/helpers.go` | readTwoSlots/writeTwoSlots use typed accessors |
| `native/system.go` | systemArrayCopy uses CopyObjectCells; objectClone uses CloneObject |
| `native/lang.go` | Migrated to typed static accessors |
| `native/exceptions.go` | Migrated to typed accessors |
| `native/native_registry.go` | getClassObject uses ClassObject CAS-once pattern |
| `classloader/classloader.go` | Thread-safe LoadClass with RWMutex + double-checked first-wins |
| `runtime/runtime.go` | GetStatic uses StaticCellToSlot; GetStaticLong/GetStaticDouble for 64-bit bridges |
| `transpile/emit.go` | cellGetterName/cellSetterName dispatch; fixed AOT types (long→GetLongCell, float→Float32bits); getstatic emits GetStaticLong/GetStaticDouble |
| `transpile/emit_test.go` | Updated to match new typed accessor pattern |
| `rtda/cell_test.go` | **NEW** — HeapCell typed access, 64-bit preservation, ToSlot dispatch, concurrent access |
| `rtda/object_test.go` | **NEW** — Typed accessors, CopyObjectCells overlap (forward/reverse/long/double), CloneObject, concurrent access |
| `classloader/classloader_test.go` | **NEW** — Concurrent loader single identity, canonical Class mirror identity |

## Verification gates

| Gate | Command | Result |
|---|---|---|
| Build | `go build ./...` | **Pass** |
| Vet | `go vet ./...` | **Pass** |
| Unit tests | `go test ./...` | **Pass** (all packages) |
| Race tests | `go test -race ./...` | **Pass** (all packages) |
| Whitespace diff | `git diff --check a0288be..9576828` | **Pass** (no whitespace errors) |
| Slice A directed tests | `go test -race -run "TestHeapCell\|TestNewHeapCells\|TestCopyHeapCells\|TestToSlot\|TestObjectTypedAccessors\|TestCopyObjectCells\|TestCloneObject\|TestConcurrentObjectAccess" ./rtda/` | **Pass** (28/28 subtests) |
| Loader identity test | `go test -race -run "TestConcurrentLoadSingleIdentity\|TestConcurrentClassMirrorIdentity\|TestLoadClassCachesSingleEntry" ./classloader/` | **Pass** (3/3 tests) |

## Not-yet-run contract acceptance gates

These gates belong to the full workstream contract, not Slice A alone:

| Gate | Status |
|---|---|
| Fixed candidate differential matrix | **Not run** (requires Slices B–E) |
| Interpreter / IR 19-fixture match | **Not run** (requires Slices B–E) |
| AOT rejection matrix | **Not run** (requires Slices B–E) |
| Race-enabled concurrency stress | **Not run** (requires Slices B–E) |
| Kernel/unit invariants (monitor, lifecycle, init) | **Not run** (Slice A kernel tests pass; B–D kernel tests pending) |
| Core regression `bash tests/run.sh` | **Not run** (test script not yet present) |
| Evidence isolation check | **Not run** (historical baseline and `baseline-63d5658/` unchanged — pending verification) |

## Evidence isolation

- Historical evidence under `docs/workstreams/r2-concurrency-evidence/baseline-63d5658/` is **untouched**.
- Research baseline files are **unchanged**.
- All new evidence written to `docs/workstreams/r2-concurrency-candidate-evidence/9576828/`.
- `docs/workstreams/r2-concurrency-fixtures/` is **unchanged** (test harness not yet written — belongs to Slice E).

## Slice B starting point

- **Branch:** `worktree-r2-thread-monitor-foundation`
- **Candidate:** `9576828`
- **What Slice A delivers:** Race-free SC HeapCells with typed accessors, concurrency-safe classloader with single identity, canonical Class mirror identity via CAS-once. No mutable heap slice escapes.
- **What Slice B starts with:** Stable `*Object` heap storage for all instance/static fields and arrays. Thread-safe class loading. `Class.ClassObject()` for canonical mirrors.
- **Slice B needs:** Attach synthetic Thread facade to `rtda.Thread` execution context, race-free execution-context IDs, frame stack confinement, Thread lifecycle (start/join/isAlive/interrupt/sleep), goroutine carrier per platform Thread, VM daemon liveness supervision.
