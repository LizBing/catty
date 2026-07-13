# Slice A Handoff

**Candidate:** `9576828`
**Branch:** `worktree-r2-thread-monitor-foundation`
**Acceptance anchor:** `a0288be`
**Date:** 2026-07-14

## Scope (relative to `a0288be`)

22 files changed, +1306 / −259.

### Blocker fixes

1. **AOT/runtime bridge HeapCell descriptor dispatch** — `ToSlot`, `cellGetterName`/`cellSetterName`, `slotExtract`/`slotConstructor` all dispatch correctly on descriptor type. Long/double bypass Slot (panic for ToSlot; dedicated bridge for AOT). Float preserves bit patterns.

2. **System.arraycopy semantics** — `CopyObjectCells` dispatch on component kind (int/long/float/double/ref/boolean/byte/char/short). Reverse copy for `srcOff < dstOff` overlap (Java memmove). Full 64-bit long/double. Reference copy via `GetRefCell`/`SetRefCell` (not atomic copy).

3. **HeapCell backing storage escape** — `Object.Cells()` and `Class.StaticCells()` removed. All ~100+ call sites migrated to typed `GetIntCell`/`SetIntCell`, `GetLongCell`, `GetDoubleCell`, `GetRefCell`, `GetStaticInt`, `SetStaticLong`, etc. No mutable backing slice exposed.

4. **Slice A acceptance tests** — 3 new test files (34 subtests) covering typed access, 64-bit preservation, overlap copy, clone, concurrent cell access, concurrent loader single identity, and canonical Class mirror identity. All run with `-race`.

### Test coverage (new, all under `-race`)

| Test file | Tests | Coverage |
|---|---|---|
| `rtda/cell_test.go` | 12 subtests | Typed access, 64-bit preservation, ToSlot dispatch, CopyHeapCells, concurrent cell r/w |
| `rtda/object_test.go` | 19 subtests | Typed accessors, CopyObjectCells overlap (forward/reverse/long/double), CloneObject, concurrent object r/w |
| `classloader/classloader_test.go` | 3 tests | Concurrent loader single identity, canonical Class mirror CAS-once, single-goroutine cache |

## Gate results

| Gate | Result |
|---|---|
| `go build ./...` | **Pass** |
| `go vet ./...` | **Pass** |
| `go test ./...` | **Pass** |
| `go test -race ./...` | **Pass** |
| `git diff --check a0288be..9576828` | **Pass** |
| Slice A directed tests with `-race` | **Pass** (31/31) |
| Contract acceptance gates (matrix, stress, regression) | **Not run** — belong to Slices B–E |

## Evidence directory

```
docs/workstreams/r2-concurrency-candidate-evidence/9576828/
├── evidence.md   # Detailed gate results, scope, per-file summary
└── handoff.md    # This file
```

## Slice B starting point

- **Commit:** `9576828`
- **Stack:** SC HeapCells with typed accessors → concurrency-safe classloader (CAS/double-check) → canonical Class mirrors (`ClassObject` CAS-once)
- **Ready for:** Attach synthetic Thread facade to `rtda.Thread`, execution-context IDs, frame stack confinement, Thread lifecycle, goroutine carriers, VM daemon liveness
- **Non-derivable:** The `componentKind` iota values (kindInt=5, kindLong=6, kindDouble=8, etc.) used by `CopyObjectCells`; typed accessor naming convention; `ClassObject` factory pattern
