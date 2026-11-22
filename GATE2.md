# Gate 2 implementation specification

Gate 2 extends the Gate 1 executable to a published, single-threaded Java 21 profile. A single standalone application must behave the same when its methods execute as generated Go, interpreted bytecode, or a mix of the two. The Java 21 reference suite defines the observable result for each case. [Architecture and delivery gates](ARCHITECTURE.md#11-delivery-gates) · [Gate 1 acceptance contract](ACCEPTANCE.md)

## 1. Compatibility profile

The target profile is `catty-jvm21-gate2` with image ABI `2`. Accept ordinary non-preview class files at version `65.0` from ordered class directories and single-release classpath JARs. Preserve the first-definition and bootstrap namespace rules from Gate 1. The executable has one application loader, one bootstrap loader, and one logical Java thread. It includes its image, Go runtime and interpreter. Both engines use the same Java references, class mirrors, statics and throwables.

| Area | Gate 2 requirements |
| --- | --- |
| Values and locals | References, `void` method return, the `int` family (`boolean`, `byte`, `char`, `short`, `int`), `long`, `float`, and `double`. Support constants, local loads/stores, `iinc`, legal stack manipulation forms, and local indexes encoded by `wide`. |
| Arithmetic and branches | Arithmetic, comparisons, conversions, bitwise operations, shifts, typed returns, conditional branches, `goto`/`goto_w`, `tableswitch`, and `lookupswitch`. Record supported opcodes individually in the profile inventory. |
| Objects and arrays | Object construction, field access, casts, `instanceof`, reference and all eight primitive array types, partially allocated multidimensional arrays, typed array loads and stores, and `arraylength`. |
| Calls and errors | `invokestatic`, `invokespecial`, `invokevirtual`, and `invokeinterface`. Include ordinary access checks, constructor and super calls, interface defaults, linkage errors, `athrow`, and exception handlers. |
| Initialization | `ConstantValue`, active use of the declaring class, superclass and required superinterface order, reentry on the owning thread, failed initialization, and later use of an erroneous class. |
| Reflection | Class mirrors, `Object.getClass`, public method lookup with `Class.getMethod`, and `Method.invoke`. Support runtime method names, inherited methods, covariant return selection, caller-sensitive access, target dispatch, argument conversion, primitive result boxing, and cause identity. |

Retain Gate 1 bootstrap bindings. Add primitive mirrors and the eight boxed primitive classes with their `TYPE` fields, `valueOf` methods and respective unboxing methods. Add the `Cloneable` and `Serializable` marker interfaces for array assignability. Include `Throwable.getCause`, supported linkage and initialization error classes, and the `Float`/`Double` raw-bit conversion methods used by the numeric fixtures. Publish exact owner, name and descriptor entries for each available bootstrap member. Every additional member admitted by the builder needs a matching definition or binding in the runtime.

All admitted opcodes must run in the interpreter. Generated Go must support the admitted instructions except `irem`, which retains Gate 1's automatic fallback and forced-AOT rejection as a regression case. Auto mode selects interpretation for the entire method before any part of it executes. Forced AOT reports an error with the method and bytecode offset and preserves any existing output executable. The build report records the engine, selection reason and original retained bytecode size for each method.

Reject `invokedynamic`, dynamically computed constants, general method handles, monitors, synchronized methods, unbound native methods, legacy subroutines, reserved opcodes, preview classes, named modules and multi-release JARs. Reject language features whose class-file semantics or bootstrap library support fall outside the published inventory. Declared-member, constructor and field reflection, accessibility overrides, generic metadata, and annotations remain outside this gate. Gate 3 adds run-time JAR admission. Gate 4 adds threads and monitors. Gate 5 adds full shutdown behavior and application measurements.

## 2. Reader, verifier, and image

1. Extend `internal/classfile` to carry typed numeric constants and their raw bits. Validate constant-pool types, the reserved second pool entry for `long` and `double`, `ldc` versus `ldc2_w`, `ConstantValue` descriptors and member attributes used for access checks. Keep the original PC and length for every decoded instruction, including the `wide` prefix and variable-length switches. Bound switch padding and table reads, check switch keys and ordering, and verify every target lands on an actual instruction.
2. Extend `internal/verify` to track JVM types and widths in the same frame analysis consumed by both engines. Count `long` and `double` as two slots for `max_stack`, `max_locals`, invoke argument counts and descriptor limits. Treat the second local slot as unusable and invalidate overlapping locals on writes. Check stack-map expansion, merges, each legal `dup*` and `pop2` form, uninitialized `this` and each allocation site's aliases. Transfer to exception handlers with the correct local state and a stack that holds only the throwable. Report the owner, method, descriptor and original bytecode offset for instruction-level verification failures.
3. Represent each primitive or reference as one logical `rt.Value` across the internal call interface. Use `Bits` for primitive payloads and preserve `Ref` as a Go-visible pointer. A long or double occupies one argument in descriptor order and two Java local slots in the callee. Use the verified stack shape to index logical values in the interpreter and generated Go. Keep `Outcome` and exact throwable identity on every return path.
4. Serialize new primitive constant, field and instruction data under ABI 2. Derive the builder's profile validation, bootstrap inventory and runtime installation expectations from the same capability definitions. Expose the inventory as structured build-report data alongside profile ID and image ABI. Tests must check report and image agreement, admit a current image and reject wrong profile and ABI values before startup or `main` runs.

[JVM values and frames](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-2.html) · [Class-file format and verification](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-4.html)

## 3. Interpreter and Go emitter

Implement each instruction family against the verified form in `rt/interpreter.go` and `internal/aot`. Use ordinary Go operations for arithmetic with matching semantics. Share helpers for operations that can throw and for numerical behavior that requires explicit rounding, saturation or overflow handling. Implement masked shifts, two's-complement overflow, signed division and remainder edge cases, Java narrowing and widening, and the floating-point rules for NaN, signed zero, infinities and conversions to integral values. Preserve rounding at Java instruction boundaries. The reference corpus should compare defined bits for zeros and infinities and test NaN by classification where its payload is unspecified.

Keep exception checks and evaluation order at the original Java offsets. Preserve class-file handler order and half-open ranges. Record the call-site PC before entering another method so an unwinding throwable reaches the correct handler in either engine. The Go emitter must generate its own operation and branch code, with shared runtime calls for observable operations such as dispatch and allocation.

Use typed primitive storage and retain GC-visible reference slots for arrays. Derive array class identity and assignability from element kind, component class and rank. Recursively check reference array covariance and keep primitive leaf types invariant. A multidimensional allocation checks every supplied size and allocates only the requested dimensions, leaving unspecified trailing dimensions null. Test Java null, bounds, cast, negative-size and array-store failures against the reference, including cases where more than one check could fail. [JVM instruction semantics](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html)

## 4. Resolution, initialization, and reflection

**Resolution and dispatch.** Keep symbolic resolution separate from selection on the actual receiver. Use the referring class, declaring class, receiver and opcode to apply class, member, package, protected-receiver and nestmate access checks within the supported profile. Resolve static and special calls using their opcode rules. For virtual and interface calls, prefer applicable class overrides and select maximally specific interface defaults where required. Exercise abstract targets and incompatible invocation shapes in the corpus. Route both engines through the shared runtime implementation.

**Deferred linkage.** Reject missing required verification dependencies and unavailable bootstrap facilities during admission. Retain an unresolved application symbol when its absence does not prevent verification and Java resolves it on use. Report these references with owner, name, descriptor and bytecode offset. Cache resolution results at the referring constant-pool entry, including persistent failure outcomes where required. The reference case should compile against a valid application dependency and run both engines and OpenJDK against the same classpath after the dependency changes. Check that an unused bad reference permits startup and that executing it produces the specified Java error in order.

**Initialization.** Prepare static fields with type defaults. Track uninitialized, initializing with an owner, initialized and erroneous states. On triggering active use, apply static constant values at the specified point, initialize required superclass and superinterfaces in the specified order, and run `<clinit>` once. Class literals and array creation leave component classes uninitialized. The owning thread can reenter and observe the current state of statics. Preserve an `Error` thrown by an initializer and wrap another throwable in `ExceptionInInitializerError` with its original cause. Later active use yields `NoClassDefFoundError`. If prerequisite initialization fails, mark the dependent class erroneous and propagate the specified failure.

**Reflection.** Resolve `getMethod` by public name and the actual parameter class identities. Retain full descriptors in method IDs and account for inherited, overridden and covariant-return members. For `invoke`, use the logical Java caller for access checks, validate the receiver and arguments, unbox and allow primitive widening, and initialize declaring classes on valid static calls. An instance call selects from the actual receiver. Box primitive returns and return null for void. Wrap a throwable from the invoked target while retaining cause identity. Report access, argument and initialization failures with their own Java outcomes. [JVM resolution and initialization](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-5.html) · [Java 21 reflection contract](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/reflect/Method.html#invoke(java.lang.Object,java.lang.Object...))

## 5. Delivery sequence

| Step | Deliverable | Evidence for review |
| --- | --- | --- |
| 0 | Land this profile and the reference runner on the Gate 2 branch. Seed the corpus with `long-crossing` and `initialization`. | OpenJDK 21 accepts the valid fixtures and rejects a load from the second slot of a `long`. The full gate reports the outstanding Catty failures and missing cases. |
| 1 | Establish ABI 2 and the shared primitive and verifier representation. Lower `long` across both engines, then finish floats, doubles, other ordinary numeric instructions, switches and wide local access. | AOT → interpreter → AOT carries the same nontrivial input and return bits. The branch stack map verifies. Add malformed width, stack shape, and control-flow cases. |
| 2 | Implement arrays, object and method resolution, constructor access, and initialization semantics. | Java 21 differential cases cover covariance, check order, hidden fields, default methods, constant-field initialization, reentry and both uses of a failed initializer. |
| 3 | Complete reflection, deferred linkage and profile enforcement. | The corpus checks inherited, overridden and bridge methods, invalid conversions and access, target cause identity, sticky application resolution errors, correct diagnostics and ABI mismatches. |
| 4 | Run the full gate and review generated Go and representative traces. | Record toolchains, command logs, build reports and traces. The standalone mixed binary and Gate 1 regression suite pass. |

## 6. Gate checks

Run these from the repository root with Go, Python 3.9 or later, and OpenJDK 21. `JAVA_HOME` can select the JDK in place of `--java-home`.

```sh
python3 -B tests/gate2.py --java-home /path/to/jdk-21 --prepare-only
python3 -B tests/gate2.py --java-home /path/to/jdk-21 --case long-crossing
python3 -B tests/gate2.py --java-home /path/to/jdk-21
```

Preparation compiles once per fixture with `javac --release 21`, checks the bytecode with `javap`, and runs the reference outcomes with `java -Xverify:all`. Known markers pin each starter oracle, including result bits, order of side effects, failure class and cause identity. The Catty runs use the exact same class files. For cases with caught exceptions, the marker stream encodes required types, causes and event order. Uncaught cases should compare exit status and the specified Java diagnostic fields independently of JDK-specific presentation.

The full command must meet all of the following conditions.

- Every required corpus case in `tests/gate2.py` runs. Add each fixture and its assertions as its corresponding implementation lands. A focused `--case` run checks only that fixture.
- Interpreter-only, mixed and forced-AOT builds match Java 21 for AOT-capable fixtures. For deliberate fallback, prove auto mode selects interpretation before entry and forced AOT rejects the whole build with the expected diagnostic.
- Build reports advertise the correct profile and ABI. Traces agree with the reports, balance frames, and demonstrate actual crossings on one logical Java thread. Review representative generated Go for arithmetic, switch control flow, and original exception offsets.
- Invalid class files and wrong profile or ABI images are rejected. Packaging and missing-bootstrap failures leave no new executable, and a failed rebuild preserves the prior executable.
- Copy the representative mixed executable into an empty directory, move its input classes aside and run it with an empty `PATH` and cleared Java home settings. It must produce the same observable result and engine trace.
- The entire existing Gate 1 suite passes on the same JDK 21 toolchain. This covers the original Crossing fixture, JAR/classpath behavior, handler boundaries, mandatory runtime invariants and the retained `irem` fallback case.

The runner prints the temporary location of its class files, logs, reports and traces. Keep handoff evidence with the project only once the full gate passes.
