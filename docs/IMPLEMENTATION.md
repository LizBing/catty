# Catty prototype implementation specification

**Assignment.** Build Gate 1 of [the architecture](ARCHITECTURE.md), with a working application launcher. One Java application becomes one native executable. Inside that process, a method compiled to Go can call an interpreted Java method, which can call a compiled Java method, mutate the original object, and throw back into the first caller. A reflective invocation must do the same thing through the runtime’s method table.

**Definition of done.** From a fresh checkout, build the tool, compile the supplied Java sources, build three Catty executables, and run the acceptance suite in §7. A review must be able to trace an actual Java bytecode instruction through verification, lowering or interpretation, runtime dispatch, and observable output. The deliverable includes the builder, runtime, fixtures, tests, and a short reproduction command. This document describes proposed interfaces and commands to implement.

## 1. Input, output and reproduction

Use the installed Go toolchain and a JDK that accepts `javac --release 21`. Keep Java tools on the development/build side. The delivered application must run with its own image and the Go runtime after moving it to an empty directory with no JDK on its executable search path. Use the host OS and architecture for Gate 1. Record the Go and JDK versions used. A local verification of the oracle used Go 1.27.1 and Temurin 25.0.4.1 on darwin/arm64. Those exact tool versions are not prerequisites.

Copy the `Crossing.java` fixture from architecture §10 into `testdata/crossing/`. Add this ordinary Java entry point as `testdata/crossing/Main.java`. This makes the reference result and the Catty result observable through the same application interface.

```java
public final class Main {
    public static void main(String[] args) throws ReflectiveOperationException {
        System.out.println(Crossing.probe());
    }
}
```

The following commands are the proposed user interface, run from the repository root after implementation. The test runner should use a temporary work directory so concurrent runs do not overwrite each other.

```sh
mkdir -p out/classes
go build -o out/catty ./cmd/catty
javac --release 21 -g -d out/classes testdata/crossing/Crossing.java testdata/crossing/Main.java
java -Xverify:all -cp out/classes Main
out/catty build --classpath out/classes --main Main --mode auto --policy testdata/crossing/mixed.json --report out/mixed.json --output out/crossing
out/crossing
```

| Input or option | Contract |
| --- | --- |
| `--classpath` | Ordered class directories and single-release JARs separated by the OS path-list separator. Accept class files of version 65.0. Prefer the first classpath entry defining a binary name. Report its origin. Reject attempts to define the reserved bootstrap namespace. |
| `--main` | Java binary name, e.g. `example.Main`. Translate to the JVM internal name for lookup. Require a public static `main([Ljava/lang/String;)V` on the entry class for this profile. |
| `--mode auto\|aot\|interp` | Auto prefers AOT and retains the whole method if only the interpreter supports it. AOT requires lowering every admitted concrete application method. Interp retains all such methods. Abstract declarations and bootstrap intrinsics have their own entry kinds in every mode. |
| `--policy` | Optional JSON for auto mode. Select exact methods by internal owner name, method name and descriptor. An unknown, duplicate, conflicting or impossible selection is a build error. |
| `--report`, `--output` | Write the build report and standalone executable to these paths on success. Complete validation and Go compilation in staging before publishing an executable. A failed build leaves any previous executable usable. |
| Application arguments | Everything after the executable name is a Java `String[]` argument to `main`. Success exits zero. An uncaught Java throwable produces a nonzero exit status with Java class and available frames on stderr. |

Use this complete `mixed.json` for the principal test. All other eligible application bodies then compile to Go.

```json
{
  "methods": [
    {
      "owner": "Crossing$Worker",
      "name": "apply",
      "descriptor": "(LCrossing$Box;)V",
      "engine": "interp"
    }
  ]
}
```

## 2. Implementation boundaries and build order

Suggested Go package layout below. Package names may change. Keep the shared JVM semantics in one place so a fix to interface dispatch, exception matching or field access reaches both engines.

| Package | Owns | Output consumed by the next stage |
| --- | --- | --- |
| `cmd/catty`, `internal/build` | CLI, ordered inputs, method policy, staging, invoking `go build` | Executable and JSON report |
| `internal/classfile` | Binary reader, constant pool, attributes, descriptors and bytecode decoding | Immutable class definitions with original bytecode offsets |
| `internal/verify` | Profile admission, stack maps, instruction type checking, control-flow analysis | Verified methods with typed block entries and handler edges |
| `internal/aot` | Go code generation and registration | Generated functions with the common entry signature, immutable image tables |
| `rt` | Image types, installation, object and class model, JVM semantics, interpreter, intrinsics, launcher | Importable Go package for generated code |
| `testdata/`, Go test packages | Java fixtures, policies, OpenJDK oracle runs and trace assertions | Reproducible end-to-end evidence |

A simple packaging choice is a temporary Go module for the generated `package main` that requires the Catty module and replaces it with the local checkout. Keep the importable runtime outside a Go `internal` directory if using a separate generated module. Embed retained instruction data and metadata as Go literals or embedded data files. Run the Go compiler at build time, then move the completed binary into place. Provide a way for tests to inspect the generated source before staging is cleaned up. Keep method ordering deterministic for identical ordered inputs. [Go module replacement](https://go.dev/ref/mod#go-mod-file-replace)

Implement in these increments. Each has an executable or inspectable checkpoint.

1. **Reader and profile.** Decode the fixture and print/report method descriptors, handler ranges and instruction offsets that agree with `javap -p -c -s`. Reject an invalid branch target and a malformed stack map. Use the same verified method data for both backends.
2. **Runtime and interpreter.** Launch a small class that constructs an object, increments an integer field and prints the value. Then run the full fixture with every application body interpreted. This exercises the semantic helpers before introducing code generation.
3. **AOT.** Lower those same methods to Go. Build and run the all-AOT fixture. Inspect generated code for `Box.bump` and `Crossing.probe`, including branches and exception handlers.
4. **Mixed execution.** Wire engine selection per method and the shared calling convention. Execute the unchanged fixture with the policy above, then exercise unforced fallback using the interpreter-only instruction in §3.
5. **Release check.** Run the comparison suite, the invalid-input and build-failure checks, and the standalone launch in §7. Report results and remaining profile limits.

## 3. Class reader, verifier, linker and image

**Reader.** Read big-endian fields with checked lengths. Decode modified UTF-8 constants without losing Java UTF-16 code units. Parse class hierarchy, access flags, descriptors, fields, methods, `Code`, exception tables, `StackMapTable`, and available source/line metadata. Account for double-slot constant-pool entries even when a method using their values lies outside the executable profile. Preserve constant-pool indices for diagnostics. Skip well-formed optional attributes that do not affect admitted semantics. Diagnose features whose semantics the profile cannot honor. Never execute a Java initializer in the builder.

**Instruction inventory.** Both engines implement the first five groups below for supported operand forms. The interpreter also implements the sixth group. The first four rows cover the observed fixture and launcher. The fifth row gives enough breadth to test application inputs, statics and ordinary branch behavior. The final row makes automatic fallback observable on a real method without a force policy.

| Group | Required operations |
| --- | --- |
| Locals, stack and constants | `aload`, `astore` and their short forms, `dup`, `pop`, `iconst_m1` through `iconst_5`, `bipush`, `ldc` and `ldc_w` for supported int, string and class constants |
| Objects, fields and reference arrays | `new`, `getfield`, `putfield`, `getstatic`, `anewarray`, `aastore` |
| Calls and exits | `invokestatic`, `invokespecial`, `invokevirtual`, `invokeinterface`, `return`, `ireturn`, `athrow` |
| Computation and branches | `iadd`, `if_acmpeq`, `if_acmpne`, `if_icmpeq` |
| Additional ordinary operations | `iload`, `istore` and their short forms, `sipush`, `aconst_null`, `putstatic`, `aaload`, `arraylength`, `areturn`, `goto`, `ifeq`, `ifne`, `if_icmpne` |
| Interpreter-only in this prototype | `irem` with Java signed remainder semantics and `ArithmeticException` on zero divisor. Auto selects interpretation for the enclosing method. Forced AOT reports the unsupported lowering. |

Maintain an explicit capability table for decoding, verification, interpretation and AOT. Reject unsupported executable types, instructions or member features with input origin, class, method descriptor and bytecode offset when available. Apply policy after verifying the whole body. A helper that rejects an unsupported type must never silently read or return a zero value. Support class and interface declarations used by the fixtures, reference arrays, integer/boolean-family category-one values, reference values and void results. Explicitly diagnose unsupported synchronization and dynamic bootstrap execution in this profile.

**Verifier.** Decode instruction boundaries first. Build successors for fallthrough, branch targets and exception handlers. Check constant-pool tags, indices, operands and descriptors, local indices, stack bounds, return types and receiver assignability. Reconstruct stack-map frames from deltas, check incoming types at joins and handler entry, and enforce Java 21 verification rules for the admitted subset. Track `uninitializedThis` and each `new` offset separately, including aliases created by `dup`. A successful legal constructor call initializes its aliases on the normal successor. Preserve the verifier’s exceptional-edge rules. Reject unsupported method shapes before either engine sees them. Keep verifier errors separate from a valid program throwing a linkage error on use. [JVM class-file verification](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-4.html#jvms-4.10.1)

**Linker and image contract.** Assign dense image-local IDs to definitions in a stable order. Symbolic member references keep owner, name and descriptor. A class definition records loader, parent, interfaces, flags, inherited field layout and declared members. A method records flags, descriptor, engine, entry ID or retained body ID, ordered handler table and optional line map. Keep metadata for all admitted members, including reflectively called methods. Record image ABI revision and profile ID. Check both at installation before running `<clinit>` or `main`. Serialize IDs and indices. Install generated Go function values through an entry table in the executable.

A retained instruction should identify its opcode, bytecode offset, typed or indexed operands, and decoded branch destinations. Both the image and AOT code need original offsets for handler lookup. Keep immutable build definitions separate from loaded class state, static storage, class mirrors, resolution caches and initialization status. Distinguish failure to provide a required bootstrap facility (profile/build error) from a missing optional application reference first resolved by an executing instruction (Java linkage error). A dependency required by verification still has to resolve during verification.

**Report.** Emit `profile`, `imageABI`, `main`, toolchain versions, ordered input origins and digests, and one record per method. Include owner, name, descriptor, engine (`aot`, `interp`, `intrinsic`, `abstract`), selection reason, and retained JVM code length in bytes. AOT bodies have zero retained interpreter bytes. Record unsupported-feature errors on stderr with enough context to locate the failing input. The test suite must read the report as data.

## 4. Runtime and engine contracts

Reuse the architecture’s `Value`, `Outcome` and `Entry` ABI. Use one GC-visible `*Object` reference representation in runtime fields, arrays, interpreter locals, generated locals and throwable outcomes. An object header carries its loaded class. Class mirrors and reflective method objects carry typed runtime payloads. Allocate instance slots in superclass-first order and address fields by declaring class plus name and descriptor. Initialize integer slots to zero and reference slots to null. Preserve pointer identity across all four combinations of caller and callee engine. Use Java `int32` arithmetic for Java `int` operations. Keep a single logical Java thread for the prototype launcher and its synchronous calls.

**Common invocation path.** Separate resolution, receiver-based selection and execution. The calling instruction supplies its opcode, symbolic reference, caller and bytecode offset. Resolve with the caller’s loader and access context. Apply opcode-specific initialization, receiver, access and method-selection rules. For `invokeinterface` on the fixture, the referenced declaration is `Step.apply` and the selected executable body belongs to the receiver’s `Worker` class. Cache the declaration separately from any receiver-specific selection. The intrinsic for reflection may start with an already resolved declaration. It must still perform receiver-based selection when required.

After selection, use exactly one owner for frame push/pop and trace enter/exit across all entry kinds. Push the selected method on the same thread, run the generated function, interpreter loop or intrinsic, and return a value or throwable. Instance receiver precedes Java arguments in locals. A cell containing Java null is a successful reference value. Set the caller frame’s current bytecode offset before a potentially throwing instruction. Do not hold a reference into a resizable Go frame slice across a nested call unless the frame representation provides stable identity.

**Interpreter.** Allocate locals and operand slots from verified maxima. Initialize receiver and arguments. Execute verified instructions by instruction index while retaining the original bytecode offset. For a call, pop arguments in reverse stack order and present them in declaration order. For success, push a non-void result. For a throw, look up a handler at the calling instruction. Clear dead reference slots when frames or buffers are reused. Every retained body can call any selected entry, including generated Go code.

**AOT emitter.** Emit one real Go body per selected Java method. Build basic blocks at branch destinations, handler entries and points where control can split. A straightforward first backend uses fixed local/stack variables and Go labels or a switch over **basic blocks**. Emit specialized Go statements for each instruction in the block. For example, `Box.bump` reads its resolved field, adds an `int32(1)`, writes it back on the same receiver, and returns void. An AOT function must execute its own lowered control flow. Routing the body to a bytecode loop or switching over individual bytecode opcodes at runtime does not satisfy this gate. Keep runtime operations such as field access, allocation and dispatch in shared helpers.

**Handler transfer.** On any Java throw, select the first matching handler in class-file order whose range contains the original offset (`start <= pc < end`). Match catch type by loaded-class assignability, with zero catch type accepting all. On a match, preserve locals, replace the operand stack with the exact throwable and jump to the handler. On a miss, return the same throwable to the caller. The interpreter and AOT emitter share matching logic. For AOT, branch to a generated handler block with the caught reference in the expected incoming slot. Treat Go panics as implementation faults. Java error paths, including null dereferences and failed array stores, must return Java throwables through the ABI. [JVM instruction semantics](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html)

**Class state.** One bootstrap loader and one application loader are enough. Install the runtime before creating `String[]` arguments and invoking `main`. Implement single-threaded preparation and initialization for admitted classes. Track uninitialized, initializing with owning thread, initialized and erroneous states. Initialize required parents and execute an admitted `<clinit>` through the same method entry. Respect recursive initialization and failure behavior in architecture §4. Class literals produce canonical mirrors without triggering initialization. Class state and static fields are common to both engines.

## 5. Bootstrap profile and reflection

Implement a small, explicit bootstrap inventory as metadata and Go intrinsics or small verified Java bodies. Use exact JVM member descriptors. Intrinsics may depend on Java platform names. User class names, field names and fixture return values must never be special cases in the builder, resolver, runtime or either engine.

| Platform operation | Required behavior |
| --- | --- |
| `Object.<init>()V`, `RuntimeException.<init>()V` | Permit normal superclass construction. Supply the Throwable, Exception, RuntimeException, ReflectiveOperationException, InvocationTargetException and Error hierarchies, including the specific error classes the supported helpers can produce. |
| `System.out:Ljava/io/PrintStream;`, `PrintStream.println(I)V` | Initialize one non-null stream and print a signed decimal integer followed by a newline. Use this to compare whole applications with the JDK. |
| `Class.getMethod(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;` | Match name and parameter class identities, search public declarations and inherited public methods in the supported hierarchy, and return a method object identifying the declaration. A missing method throws `NoSuchMethodException`. Use metadata to distinguish parameter types even when names are similar. |
| `Method.invoke(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;` | Implement public reference-argument instance calls and supported public static calls. Check access, receiver, arity and assignability. Support zero arguments, including a null args array. Apply normal virtual/interface dispatch. Return null for void and the same reference for reference returns. |
| `InvocationTargetException.getCause()Ljava/lang/Throwable;` | Return the stored original target throwable. Argument, receiver and access errors from the reflection adapter escape directly. Wrap a throwable only when it arose from the invoked target. Preserve static initialization failure behavior. |
| `Integer.valueOf(I)Ljava/lang/Integer;`, `Integer.intValue()I`, `Integer.TYPE:Ljava/lang/Class;` | Supply the int mirror and int boxing/unboxing for a focused primitive reflection test. Box an int method result, support an Integer argument for an int formal, and retain ordinary Java boxing identity guarantees for the supported range. |

Provide canonical mirrors for admitted classes, arrays and `int`, UTF-16 Java strings for literals and launcher arguments, reference arrays with component type and length, and assignability checks for arrays and interfaces. Supply metadata ancestors needed by the supported platform classes and exception types. Include the required cloneable and serializable supertypes when modeling array assignability. List implemented platform descriptors, reflective conversions and unsupported signatures in the prototype profile. New profile features must not be inferred from a class’s presence in bootstrap metadata alone. [Reflection invocation contract](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/reflect/Method.html#invoke(java.lang.Object,java.lang.Object...))

## 6. Evidence and diagnostics

Expose opt-in JSON-lines tracing through `CATTY_TRACE=<path>`. Keep Java stdout for the application. Each Java method entry event records actual OS PID, logical thread ID, depth, owner, name, descriptor and engine. A return or unwind event pairs with its entry. Catch events identify the handler method and the throwing and handler bytecode offsets. A reflection invocation event identifies the declaration and selected target. Record opaque trace-local object IDs for caught, propagated and wrapped throwables so the tests can compare identity across events. A wrapper event carries both wrapper and cause IDs. Assign stable IDs without keeping every Java object alive during normal execution. Keep tracing off by default.

In the mixed trace for the fixture, assert the following ordered observations, allowing unrelated constructor and bootstrap events between them:

1. AOT `Crossing.probe` calls interpreted `Worker.apply` through the interface declaration.
2. Interpreted `Worker.apply` calls AOT `Box.bump`. The first throw reaches the handler in `probe` with the original `Signal` ID.
3. Reflection resolves `Step.apply`, selects `Worker.apply`, calls it interpreted again and reaches AOT `Box.bump` a second time.
4. Reflection wraps the second target throw. Its cause ID equals the first `Signal` ID because this fixture throws the same field twice. The second catch in `probe` sees the wrapper. The process prints `2`.
5. These synchronous events carry the launched executable’s PID and one logical thread ID. The two `bump` entries are nested below the corresponding `apply` entries. Frame depths balance after normal and exceptional returns.

## 7. Acceptance suite and handoff

Use a Go integration test or a single test command covering the matrix. Run Java subprocesses only to prepare fixtures and obtain independent oracle results. Supply the compiler/runtime artifacts, engine reports and traces from real executions for review. Tests should compile their Java inputs in temporary directories and preserve useful diagnostics on failure.

Run `python3 -B tests/acceptance.py` from the repository root. See [ACCEPTANCE.md](ACCEPTANCE.md) for the report and trace field names, generated source export, and Go runtime cases expected by the runner.

| Test | Evidence required |
| --- | --- |
| Three-engine comparison | Build the identical `Crossing` class files in auto with mixed policy, `aot`, and `interp`. Each executable prints exactly `2\n` and exits zero, matching OpenJDK. Inspect report and trace in every mode. In all-AOT the reflective target executes as AOT. In all-interpreted the application launcher invokes interpreted `main`. |
| True AOT and standalone launch | Inspect generated source for `Box.bump` and the catch blocks in `probe`. Copy the mixed binary alone into an empty directory, supply an empty executable search path, launch by absolute path and verify the same output and trace. The built runtime must have no dependency on a JVM process or Java classpath files. |
| Automatic fallback | Compile a second Java program with a static `int rem(int x, int y) { return x % y; }` called by `main`. Auto selects interpretation only for the unsupported AOT body and reports `irem` as its reason. Check positive and negative input and the zero-divisor exception against OpenJDK. Forced AOT fails and preserves a previously built binary. |
| Independent dispatch and reflection | Use different application class and method names. Supply two implementations of one interface, select a receiver from `args.length`, and perform a call. Resolve a public inherited method by a string supplied in `args[0]`. Exercise a missing method, wrong argument type, int boxing and a target throw. Compare output or caught type with OpenJDK. Include a successful reflective call to a method referenced only by name. |
| Class state and fields | Compile a fixture with a static int set by `<clinit>`. Read and update it from opposite engine paths. Check that initialization runs once and class literals leave it uninitialized until active use. Exercise inherited/hidden fields with the same name. Compare results with the JDK. |
| Reader and packaging | Build once from a class directory and once from an ordinary JAR. Compare behavior and engine choices. Exercise classpath precedence, unsupported class version or bootstrap feature, malformed stack map, invalid branch target, unknown policy target and image ABI mismatch before application startup. |
| Exceptions and GC | Compare null, bounds and reference-array store errors, a handler covering a call and a throw at the exclusive end of a handler range. Use a runtime-level test to keep a live object solely in interpreter operands during a nested compiled call that triggers Go GC. Check the object’s field and identity after return. |
| Tool suite | `go test ./...` exercises the core packages. Document one command that runs all oracle and executable tests with a JDK available. Fail the acceptance command if the required JDK is missing. Record exact versions, exit statuses and test results in the handoff. |

When generating secondary fixtures, inspect their actual `javap` output. Choose Java source forms within the prototype profile or add a documented, tested instruction or intrinsic if a gate needs it. Do not weaken an assertion to fit the implementation. For example, verify handler end exclusivity with a small verified-method unit fixture if javac aggregates source lines into inconvenient ranges. A later production verifier and expanded Java library have distinct coverage requirements in the architecture.

**Scope after acceptance.** This assignment ends at the demonstrated single-process prototype. Architecture Gates 2–5 cover full opcode families, late plugin JAR admission, multithreaded monitors and memory ordering, shutdown hooks, and broader performance evaluation. Keep their extension points in the image, loader identity, explicit thread parameter and shared semantics. Report gaps under the published prototype profile.

**Prompt to give the implementing agent**

> Implement the Catty prototype described in docs/IMPLEMENTATION.md in this repository. Use docs/ARCHITECTURE.md revision 2 for shared runtime semantics and its Crossing.java fixture. Deliver the builder and standalone application runtime, execute the §7 acceptance suite, and report the commands, results, mixed-engine trace evidence and remaining profile limitations. The required deployment is one Java application in one Go process.
