# K2 Candidate Evidence

**Candidate branch**: `codex/r3-runtime-identity-definition-slice`
**Candidate commits**: `8907948` (base) + fix commits (see below)
**Date**: 2026-07-20
**Status**: All gates pass

## Issues Fixed

### Issue 1: Array Class Identity (`rtda/class.go`, `rtda/vm_types.go`)

**Problem**: `GetArrayClass()` generated wrong array names and primitive classes lacked correct kind.

**Fix**:
- Added `primitiveInfo` map mapping primitive names to kind and descriptor byte
- `PrimitiveDescriptor(name)` returns JVM descriptor byte
- `makeVMPrimitive` sets `componentKind` from `primitiveInfo`
- `GetArrayClass()` generates correct JVM descriptors:
  - Primitive: `"[" + descriptor` (e.g. `[I`)
  - Reference: `"[L" + name + ";"` (e.g. `[Ljava/lang/String;`)
  - Array: `"[" + name` (e.g. `[[I`)
- `componentKind` propagation: primitive arrays inherit kind; ref/array arrays get `kindNone`

**Tests** (`rtda/vm_types_test.go`):
- `TestGetArrayClassNamePrimitive` — all 8 primitive array names
- `TestGetArrayClassNameReference` — reference array name
- `TestGetArrayClassNameDeepNesting` — 3-level nesting for primitive and reference
- `TestGetArrayClassIdentityStable` — pointer stability
- `TestGetArrayClassConcurrentSameIdentity` — 32 concurrent callers, same pointer
- `TestVMPrimitiveComponentKind` — all primitive kinds correct
- `TestPrimitiveDescriptor` — descriptor byte mapping
- `TestVMPrimitiveForNameLazyInit` — safe before `InitVMTypes()`

### Issue 2: Cross-Session Definition Deadlock (`classloader/classloader.go`)

**Problem**: Two sessions defining cross-referencing classes could deadlock.

**Fix**:
- `defMu` serialises ALL top-level definitions (acquired in `LoadClassResult`)
- `sessionLoader` bypasses `defMu` for recursive calls within the same session
- `Session.seen` map detects same-session circularity (no goroutine identity dependency)

**Tests** (`classloader/classloader_test.go`):
- `TestLoadClassConcurrentDefMuSerialisation` — two concurrent goroutines, defMu serialises correctly
- `TestLoadClassCircularDependency` — A→B→A circular dependency returns `FailureCircularity`
- `TestConcurrentLoadSingleIdentity` — 32 goroutines see the same pointer

### Issue 3: Typed Dependency Failures (`rtda/build.go`)

**Problem**: `NewClass` lost typed failure info (returned nil).

**Fix**:
- Added `BuildResult{Class, Failure}` and `BuildClass()` for typed propagation
- Superclass/interface failures wrapped with dependency context
- `NewClass` retained as legacy wrapper for bootstrap callers

**Tests** (`rtda/build_test.go`):
- `TestBuildClassSuperclassFailure` — superclass failure propagates `FailureNotFound`
- `TestBuildClassInterfaceFailure` — interface failure propagates `FailureLinkage`
- `TestBuildClassSuccess` — successful build returns class with correct superclass
- `TestNewClassReturnsNilOnFailure` — legacy `NewClass` returns nil on failure

### Issue 4: Separate Lookup/Initiation from Definition (`classloader/classloader.go`)

**Problem**: `defRecords` contained delegated classes; lookup conflated with definition.

**Fix**:
- `initiatingCache` may contain delegated classes (FastPath)
- `defRecords` tracks only THIS loader's definitions
- Delegated classes (already have `DefiningLoader`) are removed from `defRecords` after caching
- `FailureDuplicateDefinition` removed
- Classfile name validation added in `ClasspathProvider`

**Tests** (`classloader/classloader_test.go`):
- `TestLoadClassDelegatedIdentity` — delegated class enters `initiatingCache` but NOT `defRecords`
- `TestLoadClassNameMismatch` — name mismatch returns `FailureFormat`
- `TestLoadClassCachesSingleEntry` — same-name loads return same pointer
- `TestLoadClassProviderAllMiss` — all providers miss returns `FailureNotFound`
- `TestLoadClassResultTypedFailure` — typed failure cached on `defRecord`

### Issue 5: Catch-Type Typed Resolution (`interpreter/interpreter.go`, `interpreter/invoke.go`)

**Problem**: Exception handler and clinit catch-type used `LoadClass` (must-load, panics on failure).

**Fix**:
- `handleException`: uses `resolveClass()` for typed catch-type resolution
  - On failure: new throwable (LinkageError/NoClassDefFoundError) replaces original
  - Returns immediately; interpreter loop re-enters `handleException` for cascading exceptions
- `runClinit`: uses `resolveClass()` with labeled `frameLoop` restart
  - On failure: replaces `thrown`, restarts handler search from current frame
- `Loop`: changed `if thread.HasException()` to `for thread.HasException()` for cascading exception handling

**Verification**: Build passes, all tests pass. Catch-type resolution failure behavior is verified via code review (requires runtime classload failure setup).

### Issue 6: VM Type and Mirror Lifecycle (`rtda/vm_types.go`, `rtda/class.go`)

**Problem**: VM primitive accessors could read nil before `InitVMTypes()`; mirror creation lacked mutual exclusion.

**Fix**:
- `VMPrimitiveForName` and `VMPrimitiveForKind` call `InitVMTypes()` for lazy self-initialisation
- `Class.classObjectMu` mutex added for double-checked locking in `ClassObject()`
  - At most one factory invocation per Class
  - Factory called under lock (loader-protected creation)

**Tests** (`rtda/vm_types_test.go`):
- `TestVMPrimitiveForNameLazyInit` — calls before `InitVMTypes` work correctly
- `TestConcurrentClassMirrorIdentity` — existing test still passes with new mutex

### Issue 7: Comprehensive Tests

**New test files**:
- `rtda/vm_types_test.go` — 9 tests: array class naming, identity, kind, descriptor, lazy init
- `rtda/build_test.go` — 4 tests: typed superclass/interface failure, success, legacy nil
- `classloader/classloader_test.go` — expanded: +8 tests for typed failure, delegation, circularity, concurrency

**Total new tests**: 21

## Gate Results

### gofmt
```
# New/modified files: clean
rtda/build_test.go     — ok
rtda/vm_types_test.go  — ok
rtda/class.go          — ok
interpreter/interpreter.go — ok
interpreter/invoke.go  — ok
classloader/classloader.go — ok
classloader/classloader_test.go — ok
rtda/vm_types.go       — ok
rtda/build.go          — ok
```
Pre-existing formatting issues in other files (not K2-related).

### go vet
```
No issues.
```

### go test
```
ok  catty/classfile    (cached)
ok  catty/classloader  0.617s
ok  catty/lowering     2.685s
ok  catty/native       1.120s
ok  catty/rtda         3.304s
ok  catty/transpile    14.695s
```

### go test -race
```
ok  catty/classfile   4.233s
ok  catty/classloader 1.265s
ok  catty/lowering    3.143s
ok  catty/native      2.259s
ok  catty/rtda        4.156s
ok  catty/transpile   15.815s
```

### e2e (interpreter)
```
HelloWorld: PASS
ExceptionTest: PASS
```

Reproduce: `cd tests && ./run.sh`

## Fix Commit Summary

Commits after `8907948`:

1. Fix Issue 1: array class identity (primitiveInfo, GetArrayClass descriptor generation)
2. Fix Issues 2-4: classloader protocol rewrite (defMu serialisation, BuildClass, name validation, delegated class defRecord cleanup)
3. Fix Issue 5: catch-type typed resolution (handleException, runClinit, Loop cascading)
4. Fix Issue 6: VM type and mirror lifecycle (lazy init, classObjectMu double-checked locking)
5. Fix Issue 7,9: comprehensive tests and evidence (vm_types_test, build_test, classloader tests)
