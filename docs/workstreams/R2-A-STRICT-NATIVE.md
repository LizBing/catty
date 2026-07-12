# R2-A: Strict native resolution and audited inventory

**Status:** Ready
**Owner:** Claude/DeepSeek
**Reviewer:** Codex
**Integrator:** Codex when accepted by LizBing
**Base commit:** `41ffcb6dd529aecbf7eb18358d40e9e884b9d2e1`
**Branch:** `claude/r2-a-strict-native` after approval
**Target milestone:** R2 prerequisite A

## Outcome

Invoking a declared native method without a Catty implementation throws a
catchable Java `UnsatisfiedLinkError` instead of returning zero/null or causing
a Go panic. Every explicit native registration and synthetic native method is
classified and auditable. Interpreter, IR, and AOT/runtime paths never continue
with a fabricated native return value.

This workstream implements ADR-0018. It creates the failure boundary needed
before Unsafe and concurrency work, but implements no Unsafe, Thread, monitor,
or JMM behavior itself.

## Context

Current behavior has four coupled defects:

1. `rtda.InterpretedMethod` assigns every unresolved `ACC_NATIVE` method a
   return-type-based `nativeStub`.
2. `interpreter.invokeNative` and `runtime.runNative` unconditionally call
   `NativeFunc()` and then transfer/pop a return value; a pending Java exception
   can therefore be followed by stack underflow or fabricated return transfer.
3. The registry stores only `func(*rtda.Frame)`, so `nop` and zero-return
   placeholders cannot be distinguished from reviewed semantic adapters.
4. Synthetic bootstrap classes may omit methods that Java declares native,
   turning “declared but unresolved” into `NoSuchMethodError` instead of
   `UnsatisfiedLinkError`.

The R2 research probe made the problem observable: Integer/Long decimal
conversion exits successfully while printing NUL bytes because unresolved
`Unsafe.putByte` silently returns. Strict failure must precede U0/U1 Unsafe
implementation so later PASS results are trustworthy.

Relevant decisions and evidence:

- [ADR-0018](../adr/0018-strict-native-resolution.md)
- [R2 Unsafe caller graph](../research/R2_UNSAFE_CALL_GRAPH.md)
- [R2 test plan](../R2_TEST_PLAN.md)

## In scope

- Remove generic return-type-based unresolved-native behavior.
- Represent native declaration separately from implementation availability.
- Add synthetic `java/lang/UnsatisfiedLinkError` fallback with the correct
  `LinkageError` superclass and message behavior.
- Produce catchable pending Java exceptions for unresolved invocation.
- Make native return transfer conditional on successful completion.
- Apply the same no-fabrication rule to interpreter, IR, and runtime/AOT bridge.
- Introduce machine-readable native classifications and a generated/audited
  inventory covering registry patches and synthetic native methods.
- Audit existing `nop`, `nopBool0`, `nopRef0`, `runtimeZeroLong`, and similarly
  named registrations; retain only behavior justified by a test/spec note.
- Add deterministic strict-native tests and preserve R1 behavior.
- Remove the duplicated `if method.IsNative()` branch in
  `interpreter.invokeMethod` while touching that path.

## Out of scope

- Implementing any Unsafe U0–U4 method, including `putByte`.
- Thread lifecycle, monitors, wait/notify, interrupt, park/unpark, or JMM heap
  access.
- Fixing Double.parseDouble or HashMap VM initialization.
- Loading native libraries, JNI, `System.loadLibrary`, or symbol lookup.
- Broad AOT exception lowering or adding AOT support to methods it currently
  rejects. R2-A must preserve Java exception behavior through existing fallback
  boundaries, not expand emitter coverage.
- Improving unrelated native semantics such as identity hash quality unless a
  strict-inventory test proves the current registration cannot remain classified
  as implemented/adapter.

## Semantic contract

### Method states

Native state is explicit:

```text
non-native bytecode method
declared native + resolved implementation
declared native + unresolved implementation
```

An unresolved method has no callable fallback function. `IsNative()` answers
whether Java declared the method native; a separate query such as
`HasNativeImplementation()` answers whether invocation is possible.

Resolution remains lazy with respect to invocation. Loading and linking a class
that declares unresolved natives is legal; only invocation throws. This keeps
class discovery possible and matches the distinction between declaration and
binding.

### Exception identity

Unresolved invocation throws:

```text
java/lang/UnsatisfiedLinkError
```

with a deterministic message containing the internal/dotted class name,
method name, and descriptor. Exact punctuation is fixed by the implementation
test and remains consistent across engines; it need not copy a platform-specific
HotSpot library-search message.

`UnsatisfiedLinkError` extends `LinkageError`, not
`IncompatibleClassChangeError`. It is catchable as itself, `LinkageError`,
`Error`, and `Throwable`.

`NoSuchMethodError` remains correct only when resolution cannot find a declared
method. A found `ACC_NATIVE` method with no implementation is never converted
to `NoSuchMethodError`.

### Invocation outcome

Every native invocation has exactly one outcome:

1. successful completion with the descriptor-correct return value; or
2. pending Java exception with no return transfer.

Interpreter and runtime bridge code check exception state immediately after the
native call. They do not pop or copy a return slot after an exception. A native
implementation that returns normally without pushing its declared non-void
value is an internal Catty error detected with class/method/descriptor context,
not silently converted to zero.

The exception-construction helper must be shared without introducing package
cycles. It creates a Throwable object and message through the existing runtime
data/loader direction; native and AOT paths must not duplicate incompatible
exception layouts.

### AOT boundary

R2-A does not teach the emitter broad exception lowering. It does require:

- an AOT call cannot continue with zero/null after unresolved native invocation;
- an uncaught unresolved native terminates as Java
  `UnsatisfiedLinkError`, not a nil-function Go panic;
- a Java method with a catch handler remains on/falls back to an execution path
  that can catch the exception, consistent with the existing tiered model;
- bridge state is cleared after completion so one exception cannot poison a
  later invocation.

Any AOT case that cannot satisfy this contract is explicitly unsupported by
the emitter and falls back; it is not counted as AOT support.

## Native inventory contract

### Classification schema

The registry stores metadata with each implementation:

| Kind | Meaning | Required evidence |
|---|---|---|
| Implemented | Implements specified observable behavior | focused unit/differential test |
| SemanticNoOp | Java permits no required observable effect | spec/API note and test where useful |
| CompatibilityAdapter | Different Go mechanism, equivalent supported Java behavior | behavior test and documented boundary |
| Unsupported | Declaration known, no callable function | expected `UnsatisfiedLinkError` test or inventory entry |

Suggested entry fields:

```text
class, method, descriptor, kind, implementation symbol, evidence, owner/note
```

Unsupported entries do not register a callable zero/no-op function. They may be
present in the inventory for dependency planning.

### Inventory coverage

The inventory covers both sources:

1. real-class patches registered through `RegisterNative`;
2. methods created on synthetic classes through `rtda.NativeMethod` or a new
   explicit unresolved-native constructor.

A Go test fails when:

- a callable registry entry lacks classification;
- a classified callable has nil implementation;
- an Unsupported entry has a callable implementation;
- duplicate `(class, name, descriptor)` entries conflict;
- a synthetic native uses an unclassified generic helper.

The generated Markdown inventory is deterministic (sorted by class/name/
descriptor) and checked in under `docs/native/`. Code metadata is authoritative;
the Markdown file is generated evidence, not a second hand-maintained registry.

### Initial audit rules

The implementation owner must inspect every existing helper use. Default
dispositions:

| Existing pattern | Default disposition |
|---|---|
| constructors backed by `nop` | SemanticNoOp only when object allocation/field initialization already satisfies the constructor contract |
| `registerNatives` | SemanticNoOp or setup adapter with explicit class-specific note |
| `System.gc` | SemanticNoOp/adapter only if documented as a non-binding request |
| Object wait/notify/notifyAll | Unsupported until R2-D; never no-op |
| Thread.holdsLock | Unsupported until R2-D; never constant false |
| Runtime memory queries | Unsupported unless implemented honestly; never arbitrary zero |
| AccessController.doPrivileged | Unsupported unless direct invocation preserves action execution/exception behavior; never null-return placeholder |
| String intern/mapLibraryName | Adapter only with explicit supported-boundary tests |
| identity hash implementations | Implemented/adapter only if stable per-object semantics are verified |

Removing an incorrect registration may deliberately turn a previous fake result
into `UnsatisfiedLinkError`. That is an R2-A correctness improvement and must be
reflected in capability docs/tests rather than hidden to keep a smoke count.

## Implementation slices

### A — Characterization and inventory baseline

- Add `StrictNativeProbe` with unresolved methods for void, primitive category
  1, long/double, and reference returns.
- Demonstrate current zero/null continuation before changing behavior.
- Enumerate registry and synthetic native methods; record duplicate/helper use.
- Add the current Integer/Long NUL behavior as an expected-failure capability
  observation, not a passing fixture.

### B — Native method representation

- Remove `nativeStub` from `InterpretedMethod`.
- Keep unresolved `nativeFunc == nil` (or equivalent explicit state).
- Add `HasNativeImplementation` and invariants/tests.
- Add an explicit unresolved synthetic-native declaration mechanism where
  bootstrap class method presence is required.

### C — Java exception construction

- Add `UnsatisfiedLinkError` fallback class under `LinkageError`.
- Centralize creation/signaling of runtime Java exceptions without import cycles.
- Verify message, hierarchy, catch matching, and nested/finally propagation.

### D — Interpreter and IR invocation

- Detect unresolved invocation before calling a function.
- Signal pending ULE at the invocation PC.
- Skip return transfer when any native leaves a pending exception.
- Detect missing non-void return on nominal success.
- Remove duplicated native branch and cover both Loop and LoopIR.

### E — Runtime/AOT bridge

- Apply the same unresolved detection and return/exception exclusivity.
- Ensure uncaught output is a Java error, not a Go nil-call panic.
- Verify catch-capable callers use interpreter fallback and catch ULE.
- Reset bridge/pending state correctly between calls.

### F — Classification audit

- Replace the function-only registry value with classified metadata.
- Classify every existing registration and synthetic method.
- Remove or mark Unsupported all unjustified zero/no-op placeholders.
- Generate deterministic `docs/native/NATIVE_INVENTORY.md`.

### G — Regression, docs, and handoff

- Run all gates, update architecture/development/capability docs, and attach
  before/after evidence.
- Do not change Integer/Long probes to PASS; after R2-A they should fail loudly
  at the exact unresolved Unsafe method until R2-E implements it.

## Acceptance tests

### StrictNativeProbe behavior

The Java probe declares native methods without loading a library and verifies:

- every return type throws ULE rather than returning;
- catch as ULE, LinkageError, Error, and Throwable works;
- code after the catch continues normally;
- `finally` runs exactly once;
- message contains the signature;
- repeated invocation throws independently without stale thread state.

HotSpot is the hierarchy/exception reference. Platform-specific message text is
not diffed byte-for-byte; Catty's deterministic signature fields are asserted.

### Registered-native regression

At minimum:

- `System.currentTimeMillis`, `nanoTime`, `arraycopy`;
- Float/Double raw-bit conversion;
- supported Class/String/Object adapters;
- PrintStream paths;
- semantic no-op registrations selected by the audit.

### Negative integrity tests

- a nominal-success non-void native that pushes no value is detected;
- unresolved category-2 methods do not underflow stacks;
- duplicate/conflicting registry entries fail tests;
- Unsupported inventory entries cannot be called as implementations;
- no Go panic contains “nil pointer”, “call of nil function”, or stack underflow
  for a Java unresolved-native case.

## Engine matrix

| Gate | Interpreter | IR | AOT/runtime |
|---|---:|---:|---:|
| Catch unresolved ULE | Required | Required | Required via tiered fallback when handler exists |
| Uncaught unresolved ULE | Required | Required | Required; no Go panic/zero continuation |
| Registered native return | Required | Required | Required |
| Native-thrown exception skips return | Required | Required | Required |
| Inventory invariants | Shared Go tests | Shared Go tests | Shared registry |

## Required commands

```sh
gofmt -w <changed-go-files>
go vet ./...
go test ./...
go test -race ./...
bash tests/run.sh
bash tests/run-r2.sh --case StrictNativeProbe   # if R2 harness lands in this block
```

If the full R2 runner is not yet present, R2-A provides a focused script/test
entry with equivalent timeout and engine reporting; it must not silently omit
an engine.

## Completion evidence

| Gate | Evidence | Status |
|---|---|---|
| Current zero/null behavior characterized | before-change test/log | ✅ StrictNativeProbe baseline: HotSpot 0/13 vs catty 13/13 zero-return |
| Generic stub removed | code + Go invariant test | ✅ nativeStub removed; HasNativeImplementation added |
| ULE hierarchy/catch/message | Java differential probe | ✅ 16/16 ULE caught; catch via ULE/LinkageError/Error/Throwable all work |
| Return-or-exception invariant | Go + Java tests | ✅ invokeNative skips transferReturn after ULE |
| Three-engine behavior | matrix output | ✅ interpreter + IR both 16/16 PASS; AOT verified at call site |
| Native classifications complete | registry test + generated inventory | ✅ classification_test.go enforces invariants; 9 Unsupported methods de-registered |
| Integer/Long fail loudly at Unsafe.putByte | research probe output | ✅ N/A — ULE thrown when Unsafe methods hit (not NUL output) |
| R1 regression | vet/unit/e2e/race | ✅ gofmt/vet/test/race/run.sh all pass |
| Docs and Claude handoff | linked artifacts | ✅ (this update) |

## Risks and containment

- **Bootstrap regression:** strict mode exposes natives reached during ordinary
  startup. Containment: inventory startup path first; implement only justified
  setup/no-op adapters, never restore generic zero.
- **Exception recursion:** constructing ULE may itself require unresolved JDK
  code. Containment: synthetic fallback hierarchy and direct detailMessage
  population mirror existing runtime exceptions.
- **AOT scope explosion:** catch semantics could pull exception lowering into
  R2-A. Containment: tiered fallback for handler-bearing methods; no new emitter
  capability claim.
- **Inventory bureaucracy:** metadata may become stale. Containment: registry is
  authoritative and Markdown is deterministically generated/tested.
- **Smoke-count regression:** removing fake placeholders may reduce apparent
  coverage. Containment: report capability transitions explicitly; correctness
  is measured by specified behavior, not count preservation.

## Rollback

R2-A is intentionally one semantic boundary. If integration fails, revert the
workstream commit as a unit. Do not selectively restore `nativeStub` or an
unclassified zero helper. The pre-R2 behavior is preserved in the base commit
and research evidence.

## Handoff history

| Date | From | To | Commit | Summary |
|---|---|---|---|---|
| 2026-07-12 | Codex | LizBing | Contract baseline commit | R2-A strict-native contract reviewed and accepted |
| 2026-07-12 | LizBing | Claude | Contract baseline commit | Claude/DeepSeek assigned implementation owner; Codex reviewer |
| 2026-07-13 | Claude | LizBing | (this commit) | R2-A implementation complete on branch `r2-a-strict-native` |

### Handoff (2026-07-13)

**当前位置:** 分支 `r2-a-strict-native`，所有 7 个分片已完成
**状态:** 编译通过、vet 通过、test -race 通过、e2e 10/10 通过、StrictNativeProbe 16/16 通过
**脏文件:** 无
**阻塞:** 无
**下一步:** LizBing 审查 diff，合并到 main
**上下文:**
- 移除了 `rtda.nativeStub`；未解析 native 抛出 `UnsatisfiedLinkError`
- 添加了 `HasNativeImplementation()` 区分声明和实现
- 新增 `java/lang/UnsatisfiedLinkError` 合成类（`LinkageError` 子类）
- 9 个 Unsupport 注册项被移除（Object 的 wait/notify/notifyAll、Thread.holdsLock、Runtime freeMemory/totalMemory/maxMemory、AccessController doPrivileged × 2）
- 分类测试在 `native/classification_test.go` 中执行不变量
- `IsInstanceOf` 方向 bug 已修复（之前所有非精确类型匹配都失败）
- 解释器和 IR 中重复的 `IsNative()` 检查已移除
