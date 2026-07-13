# R1 class-initialization behavior vs JVMS §5.5

**Evidence:** `docs/workstreams/r2-evidence/run-r2-results.txt` and `matrix.md`, Temurin
25.0.3 versus catty R1 in pure-synthetic mode.

## Current R1 mechanism

R1 records only `initStarted bool` on `rtda.Class`. `ensureInitialized` returns immediately
for interfaces, sets the boolean before executing `<clinit>`, initializes the superclass,
and has no initialized/erroneous distinction or initializing-owner record. Interpreter and
IR call it at `new`, `getstatic`, `putstatic`, and `invokestatic`. The AOT bridge calls it
from `Bootstrap` and `GetStatic`, but its `InvokeStatic` path lacks the equivalent request.

## Java 25 contract and measured gaps

| Area | Java 25 requirement | R1 finding | Evidence |
|---|---|---|---|
| State | not-initialized, initializing, initialized, erroneous | one boolean; no erroneous state | code inspection; `ClinitThrows` |
| Recursive request | request by the initializing execution context completes normally without re-running `<clinit>` | no owner identity; behavior not represented explicitly | `RecursiveInitialization` |
| Direct bytecode triggers | `new`, `getstatic`, `putstatic`, `invokestatic` | Interpreter/IR wired; AOT `InvokeStatic` gap | code inspection; `InvokeStaticInit` |
| Member owner | initialize the class/interface that actually declares the resolved static member | referenced class and wrong static storage can be used | `GetstaticOwner` |
| Constants | `getstatic` of a constant variable does not request initialization | no explicit distinction in the access path; ordinary javac use is inlined | `ConstantFieldNoInit` pins observable baseline |
| Class predecessors | superclass, then recursively selected default-bearing superinterfaces, before the class | superclass only; interfaces skipped unconditionally | `InterfaceDefaultInit` |
| Interface predecessors | initializing an interface does not initialize its superinterfaces | interfaces never initialize at all | code inspection |
| Failure | wrap non-`Error` abrupt completion in EIIE; mark erroneous; later request throws NCDFE | raw exception; no erroneous transition | `ClinitThrows` |
| Predecessor failure | mark subclass erroneous and propagate the predecessor's abrupt reason; do not run subclass `<clinit>` | not represented | `SuperclassInitializationFailure` |

## Correct trigger boundary

The earlier draft incorrectly described `invokevirtual`, `invokespecial`, default-method
`invokeinterface`, and `assert` as independent initialization triggers. They are not.

Default-bearing interfaces enter the observable sequence because initializing a **class**
recursively initializes the superinterfaces that declare non-abstract, non-static methods.
`InterfaceDefaultInit` therefore prints the interface marker before `after-new`; the later
default-method invocation adds no initialization request.

Method handles, reflection, and VM startup can also request initialization under JVMS §5.5,
but those mechanisms are outside the bounded R2 implementation slice and must later reuse
the shared service.

## Failure and concurrency boundaries

The later erroneous-class request must throw `NoClassDefFoundError`; the current evidence
does not justify prescribing a particular cause object or cause chain. The contract pins
the observable Java 25 result instead.

JVMS initialization also specifies locking, waiting, and happens-before behavior for
multiple Java threads. R2 has no accepted concurrency model, so ADR-0025 preserves a future
synchronization boundary but claims only one execution context. A Go lock by itself would
not establish the full Java-visible concurrency contract.

## R2 classification

| Severity | Item | Required correction |
|---|---|---|
| Bug | static declarer-owner/storage mismatch | resolve and initialize the actual declarer |
| Major gap | no interface initialization | implement class predecessor enumeration |
| Major gap | no erroneous state or EIIE/NCDFE behavior | implement four-state transitions and failure transport |
| Major gap | no explicit recursive-request semantics | record initializing owner and return normally on recursion |
| AOT gap | `invokestatic` path lacks initialization request | add the same shared request at the AOT boundary |
| Deferred | cross-thread locking/visibility | decide under ADR-0018; do not claim in this slice |
