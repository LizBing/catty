# R2 initialization slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Base commit:** `ecb086e`
**Roadmap item:** Phase R2 — Runtime semantics and concurrency planning
**Governing ADRs:** ADR-0016 through [ADR-0025](../adr/0025-class-initialization-state-machine.md) (Accepted)
**Prerequisites:** ADR-0025 Accepted; research evidence recorded

## Outcome

Implement the bounded Java 25 single-execution-context initialization contract in
ADR-0025 across Interpreter, IR, and AOT, with differential evidence for all initialization
fixtures. This workstream does not implement Java concurrency.

## Scope

- Replace `initStarted bool` with the four semantic states and an initializing-owner
  execution-context identity behind one shared initialization service.
- Implement requests at `new`, `getstatic` (except constant variables), `putstatic`, and
  `invokestatic`, using the resolved member's actual declaring class/interface.
- During class initialization, initialize the superclass and recursively enumerated
  default-bearing superinterfaces in JVMS §5.5 order. Do not initialize superinterfaces
  merely when an interface itself initializes.
- Make recursive same-owner requests return normally without re-running `<clinit>`.
- Implement predecessor-failure propagation, `ExceptionInInitializerError` wrapping,
  erroneous state, and later `NoClassDefFoundError` behavior.
- Add the minimal synthetic exception support needed for those observable semantics.
- Preserve equivalent behavior in Interpreter, IR, and AOT. Add an AOT `invokestatic`
  guard; do not add initialization guards to virtual, special, or interface invocation.
- Promote all 9 initialization fixtures in the research matrix to a permanent regression
  gate: the original 6 plus `ConstantFieldNoInit`, `RecursiveInitialization`, and
  `SuperclassInitializationFailure`.

## Non-scope

- Per-class Go locks, cross-Java-thread waiting, monitor behavior, deadlock detection, JMM
  visibility, volatile/final-field publication, or thread-to-goroutine mapping.
- `invokevirtual`, ordinary `invokespecial`, `invokeinterface`, or `assert` as initialization
  triggers; Java 25 does not define them as such.
- Reflection, method handles, VM-startup initialization, `invokedynamic`, broad `java.base`
  compatibility, String/Slot/object-layout changes, JIT, or unrelated AOT refactoring.
- CI wiring of the differential harness.

## Semantic constraints

- Java 25 is the baseline and Temurin 25.0.3 is the differential reference.
- One state record belongs to one runtime class/interface identity, including defining
  loader identity.
- Unsupported engine paths are explicit `Fallback` or `Not implemented`; no silent
  approximation qualifies as compatibility.
- No engine owns a separate initialization state machine.

## Required completion state by engine

| Capability | Interpreter | IR | AOT |
|---|---|---|---|
| Shared four-state service + owner identity | Required | Required | Required |
| `new` / `getstatic` / `putstatic` / `invokestatic` requests | Required | Required | Required |
| Declarer-owner and constant-field rules | Required | Required | Required |
| Superclass + default-bearing-superinterface order | Required | Required | Required or explicit Fallback |
| Recursive request and failure semantics | Required | Required | Required or explicit Fallback |
| AOT `invokestatic` guard | N/A | N/A | Required |

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| Interpreter initialization matrix | `bash docs/workstreams/r2-evidence/run-r2-diff.sh` → all 11 initialization fixtures match | Pass |
| IR initialization matrix | Same harness → all 11 initialization fixtures match | Pass |
| AOT initialization matrix | Same harness → match where supported; every unsupported path explicitly classified | Pass (7 match, 4 NI) |
| R1 regression | `go vet ./... && go test ./... && go test -race ./... && bash tests/run.sh` | Pass |
| Governance consistency | `git diff --check` | Pass |

Results key: `Pass` / `Fail` / `Not run` / `Not implemented`.

## Handoff (2026-07-13)

### Interpreter & IR

All 11 initialization fixtures match Temurin 25: ClinitOrder, ClinitThrows,
ConstantFieldNoInit, DirectInvokeStaticInit, GetstaticOwner,
InheritedStaticInit, InterfaceDefaultInit, InterfaceNoInitOnImpl,
InvokeStaticInit, RecursiveInitialization, SuperclassInitializationFailure.

### AOT

| Fixture | AOT Result | Classification |
|---|---|---|
| ConstantFieldNoInit | match | Supported |
| DirectInvokeStaticInit | match | Supported |
| GetstaticOwner | match | Supported |
| InheritedStaticInit | match | Supported |
| InterfaceDefaultInit | match | Supported |
| InvokeStaticInit | match | Supported |
| RecursiveInitialization | match | Supported |
| ClinitThrows | transpile refusal (exception handlers in main) | Not implemented |
| ClinitOrder | transpile refusal (unsupported opcodes) | Not implemented |
| InterfaceNoInitOnImpl | transpile refusal (unsupported opcodes) | Not implemented |
| SuperclassInitializationFailure | transpile refusal (unsupported opcodes) | Not implemented |

### Key decisions

- **Declarer-owner rule**: `getstatic`/`putstatic` initialize and access the
  field's actual declaring class (`field.Owner()`), not the referenced class
  from the constant pool. Applied in Interpreter, IR, and AOT bridge.
- **State machine ordering**: `MarkInitInProgress` is called BEFORE predecessor
  recursion (JVMS §5.5 step 6 precedes step 7). Same-owner recursive requests
  during predecessor init detect `initInProgress` + matching `ecID` and return
  normally.
- **AOT clinit failure — interpreter fallback**: When class initialization
  fails during AOT execution, the bridge panics with a typed `InitFailure`
  sentinel. The emitted program's `main()` recovers this sentinel and calls
  `FallbackToInterpreter` to re-run the entire program through the interpreter
  — providing real semantic fallback where Java exception handlers can observe
  ExceptionInInitializerError / NoClassDefFoundError. Methods with their own
  exception handlers that could intercept clinit failure are explicitly
  refused at build time (NO-BUILD) until AOT exception propagation is wired.
- **Direct AOT invokestatic init guard**: Before every direct Go call to an
  AOT'd static method, the transpiler emits `runtime.EnsureInit` targeting
  the resolved method's actual declaring class (not the constant-pool
  referenced class, in case of inherited static methods). A new
  `DirectInvokeStaticInit` fixture proves `<clinit>` runs exactly once even
  when the called method does not read/write static state.
- **Inherited static method init**: The `InvokeStatic` bridge resolves the
  method before requesting initialization, using `method.Owner()` for the
  init target. A new `InheritedStaticInit` fixture proves the declaring
  ancestor is initialized, not the constant-pool referenced subclass.
- **NCDFE detailMessage fix**: `setDetailMessage` now allocates a real
  `java/lang/String` via the loader instead of allocating from the wrong
  class in the Throwable hierarchy.
- **AOT InvokeStatic bridge**: Args are now wrapped with `slotConstructor`
  (e.g. `rtda.IntSlot(t3)`, `rtda.RefSlot(t4)`) — fixes a compilation error
  in emitted Go code. Nil LookupField/LookupMethod results now panic
  descriptively instead of nil-dereference.

## Review

**owner** — only the Owner may accept this workstream or mark it Done.

## Acceptance record

Accepted by Owner on 2026-07-13. The frozen contract authorizes implementation only within
this document's Scope, Non-scope, Semantic constraints, and Acceptance gates.
