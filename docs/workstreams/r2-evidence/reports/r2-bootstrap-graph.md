# Bootstrap kernel capability/dependency graph (Slice B — ADR-0022 mapping)

**Baseline**: R1 at governance baseline commit `ecb086e`. Pure-synthetic mode (no `java.base`
extracted; `$CATTY_BOOT` unset).

ADR-0022 instructs catty to define its bootstrap boundary **by irreducible capabilities,
not a permanent list of whole classes**. The six-class R1 bootstrap set
(`Object`, `String`, `Class`, `System`, `Thread`, `Throwable`) and the broader synthetic
set (~21 classes) are an **implementation baseline, not a permanent architecture
constraint**.

This report maps each ADR-0022 **candidate** capability to its current R1 implementation
and records questions for later work. It does not decide that all seven candidates are
irreducible, that any whole Java class must remain synthetic, or that the current provider
order is the future architecture. Those choices require a dedicated Accepted workstream.

## Capability mapping

| Candidate capability | Current R1 provider (class / file:line) | R1 form | Minimum observed responsibility | Open facade/provider question |
|---|---|---|---|---|---|
| **type/class identity** | `java/lang/Class` builder plus `rtda.Class` | bootstrap synthetic facade + Go runtime record | Engines must observe consistent Java type identity | How much of `java/lang/Class` can be ordinary classfile behavior remains unproven |
| **object allocation** | `java/lang/Object` builder; `rtda.Class.NewObject`; AOT runtime bridge | bootstrap synthetic facade + Go allocation | Engines need a common object/type/layout contract | Native allocation does not prove the whole `Object` facade must remain synthetic |
| **class mirrors** | singleton `ClassObj` associated with each `rtda.Class` | Go identity exposed through synthetic `Class` | One mirror identity per runtime class identity | Reflection surface and ordinary `Class` bytecode participation remain open |
| **string constant materialization** | `newString` / `runtime.NewString`; synthetic `String` with `Extra()` | shared bridge with Go-string R1 payload | `ldc` and engines must create Java Strings with identical semantics | ADR-0027 decides representation; facade and java.base participation are separate decisions |
| **throwable transport** | `rtda.Thread` pending exception plus synthetic `Throwable` hierarchy | Go control signal + Java object | Engines need compatible throw/catch/unwind behavior | AOT exception transport and division between native kernel and Java methods need evidence |
| **execution-context attachment** | `rtda.Thread`; synthetic `java/lang/Thread` | single-context R1 representation | Runtime operations need an owning Java execution context | Thread facade, scheduler mapping, identity and concurrency wait for ADR-0018 research |
| **class-initialization state** | `rtda.Class.initStarted`; shared initialization entry points | Go runtime state | Engines must share Java-visible initialization transitions | ADR-0025 now fixes single-context semantics; concurrency synchronization remains open |

## Additional synthetic classes not in the bootstrap set

These are served by `SyntheticProvider` (`classloader/classloader.go:52-63`), below the
bootstrap threshold but still synthetic in R1:

| Class | File:line | Role | Permanent facade? |
|---|---|---|---|
| `java/lang/StringBuilder` | `native/lang.go:323` (`buildStringBuilderClass`) | `+` concatenation substitute for pre-Java-9; stores `*stringsBuilder{[]byte}` in Extra | Could be replaced by java.base's real StringBuilder classfile if the bootstrap kernel provides char[]/String backing |
| `java/io/PrintStream` | `native/io.go:10` (`buildPrintStreamClass`) | `System.out`/`err`; wraps `io.Writer` in Extra | Could be replaced by java.base's real PrintStream if I/O and charset infrastructure are present (R4 I/O arc) |
| Exception subclasses (11 total) | `native/exceptions.go:7-29` | RuntimeException, NPE, ArithmeticException, ClassCastException, IllegalArgumentException, IndexOutOfBoundsException, ArrayIndexOutOfBoundsException, Error, LinkageError, IncompatibleClassChangeError, NoSuchMethodError | These are thin — a `detailMessage` field + constructors. They could come from java.base classfiles once the bootstrap provides the Throwable transport |
| `java/lang/Comparable` | `native/exceptions.go:32` (`buildComparable`) | Interface, no methods — marker only | Trivially replaceable by java.base |

## Current R1 provider precedence

R1's `ClassLoader` chain is hardcoded: `ArrayProvider → BootstrapProvider → SyntheticProvider →
ClasspathProvider` (`classloader/classloader.go:103-113`). With a real `java.base` extracted:
- User classes skip BootstrapProvider and SyntheticProvider (those names are claimed
  by the bootstrap kernel) → served by ClasspathProvider from java.base.
- The 6 bootstrap classes stay synthetic (BootstrapProvider wins first for those names).
- SyntheticProvider serves the non-bootstrap synthetics (StringBuilder, exceptions, etc.)
  — **these shadow real java.base classes**. For example, `java/lang/StringBuilder` from
  java.base would never be reached because SyntheticProvider wins first.

This order is evidence about R1, not an accepted migration design. Removing a synthetic
registration, retaining a bootstrap class-name claim, or composing native methods with an
ordinary classfile are distinct future options. A bootstrap workstream must test them and
define precedence per capability rather than infer permanence from the current six names.

## Open: the point where java.base classfiles can participate

Observed or known decision gates include:
1. **`jdk.internal.misc.Unsafe` and dependent runtime services** reached by researched
   java.base paths. The required supported subset must be established by reachability and
   differential evidence; this report does not authorize broad Unsafe implementation.
2. **Thread/monitor/JMM** — needed for any java.base code path that uses
   `synchronized`, `Thread` APIs, or `volatile` fields. Deferred past the first R2
   implementation slice (concurrency needs separate research per ADR-0018).
3. **Assertion-status service** — Java initialization determines desired assertion status;
   ordinary assertion bytecode is not a separate initialization trigger.
4. **Full String representation** — java.base's `String` classfile uses a `byte[]`
   value field with `coder` (COMPACT_STRINGS). Until catty's String representation
   is Java-visible (ADR-0023), java.base's String classfile cannot replace the
   synthetic one.

The **class-init state machine** (ADR-0025, Accepted; proposed implementation slice)
removes no java.base blockers directly but is the prerequisite for any code path
that exercises `java.base`'s own `<clinit>` methods — many java.base classes have
non-trivial static initializers. Without correct init semantics, even a
java.base-loaded String classfile would be incorrectly initialized.

## Closure decision

This report closes the R2 mapping task but intentionally creates no ADR-0026. ADR-0022
already governs the direction; the evidence here does not yet select a concrete kernel/
facade/provider design. Number 0026 remains unused until a future bootstrap workstream has
a real decision to propose. ADR numbering need not be contiguous.
