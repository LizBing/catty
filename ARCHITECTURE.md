# Java on the Go runtime

**Architecture proposal · revision 2**

One Java application builds to one Go executable and runs in one OS process. That executable contains the application's AOT-compiled methods, retained bytecode, a Go interpreter, and the runtime used by both engines. The build tool translates selected JVM bytecode methods to Go source. At execution time each resolved target has an entry that either calls generated Go code or enters the interpreter on the current logical Java thread.

For example, compiled `main` can call an interpreted method that updates an object and calls another compiled method. All three frames use the same object reference, heap, statics, class metadata and monitors. If the innermost method throws, the throwable travels back through those frames to the matching Java handler. This entire call chain runs inside the application's process. Optional support for loading new classes later extends that same application.

## 1. Compatibility contract

Use Java SE 21 as the semantic reference. The initial input contract is ordinary, non-preview class files produced by `javac --release 21` (version 65.0), supplied as directories or single-release classpath JARs. Major/minor version, supported features, available library members, runtime ABI revision, and build target form an explicit **runtime profile**. A supported class-file version alone says little about whether an application will run. The production profile must ship a machine-readable member and feature inventory with the builder and runtime. [JVM class-file format](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-4.html)

| Capability | First mixed-execution milestone | Subsequent profile |
| --- | --- | --- |
| Execution | Constructors, instance and static fields, direct and interface calls, integer operations, object identity, throwing and catching | Complete ordinary instruction families, arrays, floating point, virtual and default interface methods, monitors |
| Loading | Image classes, one bootstrap loader and one application loader | Registered, parent-first plugin loaders that accept additional verified bytecode |
| Reflection | Class literals, public method lookup and invocation, arguments and target exceptions used in §10 | Declared members, constructors, fields, access checks, annotations and generic signatures as their library APIs are implemented |
| Libraries | Explicit bootstrap subset for the fixture, including its error paths and the array and boxing support used by reflection | Expand from dependency reports for named workloads. Publish supported methods and native bindings per release |
| Deployment | One Go executable launches the application | Concurrent Java threads in the same process |

Arbitrary Java-defined loaders, named application modules, JNI, agents, dynamic proxies, reference processing, and general method handles each require their own implementation work. The initial profiles exclude them. In particular, common Java 21 lambdas and string concatenations can use `invokedynamic`. Admission must identify unsupported bootstrap methods even when the class otherwise looks simple. A later extension can recognize specific standard bootstraps and preserve their linkage and error behavior before supporting general bootstrap execution. A threaded production profile also requires the memory-model tests in §7. No Java SE conformance claim is implied by these milestones.

## 2. Build, link, and package

```text
application JARs + bootstrap library + runtime profile + compilation policy
                    |
            parse and verify
                    |
        descriptors and verified method bodies
             /                         \
   lower selected bodies to Go       encode retained bytecode
             \                         /
              image tables and resources
                         |
         go build with interpreter, runtime and launcher
                         |
               Go executable for the target
```

1. **Read inputs.** Preserve classpath order, class names, descriptors, exception tables and bytecode offsets. Record which input supplied each class and resource. Within an ordered application classpath, the first matching class wins. Reject duplicate definitions within a loader registration transaction. Reject unsupported packaging such as multi-release JARs in the initial profile.
2. **Verify once per admitted definition.** Check bounds, constant-pool kinds, branch and handler targets, descriptors, stack-map types, category-two values and uninitialized-object rules. Resolve verification dependencies as required by the verifier. Both engines consume the same verified representation. Keep run-time symbolic references for operations whose linkage failure is observable only on use. Image building may reject inputs outside the declared profile and must distinguish that diagnosis from Java linkage failure. [Class-file verification](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-4.html#jvms-4.10)
3. **Select an implementation per concrete method.** Default to AOT when the backend implements every operation in the method. Otherwise retain the entire body when the interpreter and library profile support it. A policy file may force either engine by owner, method name and descriptor for tests or size control. Reject an explicitly forced AOT method that cannot compile. Abstract declarations carry metadata only. Native declarations bind to registered implementations.
4. **Generate Go.** Translate the verified control-flow graph into typed temporaries and basic blocks. A basic-block dispatch loop is an acceptable first lowering for irregular control flow. Preserve evaluation order and associate every throwing operation with its original bytecode offset. Use shared runtime helpers for observable operations such as invoking an interface, initializing a class and storing a reference. Pure integer computation can stay in Go locals.
5. **Assemble the image.** Emit immutable class and member definitions, field layouts, symbolic constants, line tables when present, handler ranges, retained bodies, resource bytes, and a generated table of Go entry functions. All admitted image classes keep callable bodies or bindings for their concrete methods, including methods only discoverable by name. This retention rule gives reflection and future plugin calls a predictable contract. Stripping members would require a different profile.
6. **Bind artifact to runtime.** The build manifest records ordered input digests, toolchain and library revisions, profile, target, engine choices, and ABI revision. Image tables refer to stable indices within that build. A Go function address is never a serialized image identifier. New AOT code requires building a new executable. New compatible JARs use the interpreter in a running application with plugin loading enabled.

A build report lists every method's selected engine, retained bytecode size, unresolved optional symbolic references, missing platform bindings and unsupported features with owner, descriptor and bytecode offset where applicable. Missing required profile facilities fail the build or plugin admission. An optional application dependency in an unused method can still produce a Java linkage error when the instruction is first executed.

## 3. Runtime ownership, classes, and heap

| Owner | State | Lifetime |
| --- | --- | --- |
| Process | Generated Go functions, immutable definitions and verified image data, Go scheduler and collector | Executable lifetime |
| Application runtime | Bootstrap and application loader state, static fields, class mirrors, interned Java strings, thread registry and I/O bindings | Application startup through shutdown |
| Defining loader | Defined classes, runtime constant pools, resolution caches, plugin dependencies and resources | Application lifetime in the first plugin profile |
| Java thread | Logical ID, Java frame stack, interrupt status, monitor ownership, exception and return state | Java thread lifetime, including synchronous callbacks |
| Java object | Runtime class, identity, instance storage, optional monitor state | Reachable through ordinary Go references |

A runtime class is identified by its **binary name and defining loader** within the application. Both engines see the same `Class` object, initialization state and statics for that class. The bootstrap namespace is reserved. A plugin must resolve a shared API interface through its parent loader to implement the application's version of that interface. Array class identity follows component type and dimensionality. [JVM loading and linking](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-5.html#jvms-5.3)

Use a Go pointer to a runtime-owned object header as the canonical Java reference. Go `nil` represents Java null. A header points to its runtime class and typed payload. The linker assigns inherited and declared instance fields stable slots, keyed by declaring class, name and descriptor. Hidden superclass fields keep separate slots. Use GC-visible reference slots and separate fixed-width primitive storage. Reference arrays keep typed Go pointers and a reified component class for store checks. Primitive arrays keep their element representation. Java strings preserve UTF-16 code units, including isolated surrogates, across internal operations. Any conversion to Go text has an explicit encoding policy.

A Java reference must remain visible to the collector wherever it can live, including interpreter operands, AOT locals, exception outcomes, caches, and registered native bindings. Never encode a live pointer solely in an integer or a byte buffer. Both engines pass the same header pointer across a call. Reference equality compares that pointer, and an identity hash is assigned once independently of addresses. Clearing a retained reference in a native binding drops that GC root. Go decides when to reclaim unreachable objects. Finalization, weak references and class unloading require separate semantics and tests. [Go GC guide](https://go.dev/doc/gc-guide)

## 4. Resolution, initialization, and dispatch

**Resolution.** Cache symbolic resolution in the referring class's runtime constant pool. A method reference resolves using its owner, name, descriptor, loader context and access rules. Selecting an implementation for a receiver is a separate operation. For static and special invocation, implement the opcode-specific selection rules, including constructors and superclass calls. For virtual and interface invocation, consult the receiver class, including inherited implementations and applicable interface defaults. Preserve the specified null, access, abstract-method and incompatible-class-change failures in their required order. Publish resolution outcomes safely to concurrent callers and preserve sticky linkage errors where the specification requires them. [JVM resolution and method selection](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-5.html#jvms-5.4.3) · [JVM instruction semantics](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html)

**AOT optimization boundary.** An exact receiver, a final target, or a statically bound call can justify a direct call after the required checks. A virtual call optimized from image-only knowledge must guard the receiver's runtime class and retain the full dispatch path. An interpreted subclass loaded tomorrow can then override a method called by yesterday's executable. Cache by runtime class identity and call site. Immutable method tables and the exclusion of class redefinition keep that guard valid.

**Initialization.** Preparation allocates static storage with default values. Model each class as uninitialized, initializing with an owning Java thread, initialized, or erroneous. On a triggering use, perform the JVM initialization procedure, including constant-field handling and required superclass and superinterface initialization. The owning thread may reenter. Competing threads wait. Execute Java initializers outside the bookkeeping lock, publish completion under synchronization and wake waiters on either outcome. The first failing use receives the specified error, including wrapping a non-`Error` throwable in `ExceptionInInitializerError` when applicable. Later active uses of the erroneous class fail with `NoClassDefFoundError`. Reading a class literal must not run its initializer. Build-time execution of application initializers is excluded. [JVM initialization procedure](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-5.html#jvms-5.5)

**Plugin admission.** The application supplies JAR bytes and a parent loader. Stage immutable input and validate archive bounds, names, profile features and verification dependencies before publishing a loader handle. Bound decompressed bytes, entry count and verification work. Failed admission publishes no partial loader. Class definition, constant-pool resolution and initialization retain their Java timing after admission. Two plugins with identical binary names may coexist through distinct defining loaders. Version replacement creates a new loader and explicitly rebinds application entry points. Existing objects and executing frames keep their original classes.

## 5. Execution contract

The following Go signatures describe an internal calling convention. The image generator and runtime are versioned together.

```go
type Value struct {
    Kind Kind        // void, int, long, float, double, or reference
    Bits uint64      // primitive payload
    Ref  *Object     // GC-visible reference
}

type Outcome struct {
    Value  Value
    Thrown *Object
}

type Entry func(t *Thread, method *Method, receiver *Object, args []Value) Outcome
```

Arguments are in descriptor order. A long or double occupies one ABI argument and two JVM local-variable slots. An instance receiver is separate, static calls use a nil receiver, and a returned Java null is a successful reference value. A non-nil `Thrown` is the exact Java throwable and takes precedence over the value. Each entry borrows its argument buffer for the duration of the call. Each logical invocation maintains the Java thread, method, current bytecode offset and call site. An AOT fast path may use typed Go arguments if it preserves every observable effect of this convention. Reflection, interpreter calls and host bindings use adapters to the same entries.

**Interpreter.** Decode verified bytecode into instructions with original offsets and indexed operands. A frame contains locals, operand stack, program counter, return site and held synchronized-method monitor if any. Implement instruction families against the same field, dispatch, initialization, allocation, array and monitor helpers used by AOT. A missing AOT lowering chooses interpretation before the method starts. Restarting a partly executed method would repeat side effects and is forbidden.

**Exceptions.** A throwing helper returns a Java throwable. At the throwing offset, scan handler entries in class-file order for a covering half-open range and compatible catch type, including catch-all. A matching handler clears the operand stack and receives only the throwable. Otherwise release a synchronized-method monitor, pop the frame and propagate to the caller at the invocation offset. Explicit synchronized blocks follow their bytecode and compiler-generated cleanup handlers. A native binding uses this outcome protocol too. A Go panic denotes an implementation fault. It cannot enter a Java catch block. Capture Java frames for supported Throwable stack-trace APIs, with source lines when available, so mixed stacks preserve Java call order. [JVM instruction semantics](https://docs.oracle.com/javase/specs/jvms/se21/html/jvms-6.html#jvms-6.5.athrow)

**Arithmetic and failure points.** Preserve integer widths and overflow, masked shift counts, signed and unsigned right shift, integer division by zero, narrowing conversions, floating-point NaN and signed zero. Specify conversion saturation and rounding in helpers where native conversion behavior can differ. Generate explicit Java null, bounds, negative-size, array-store and cast checks in bytecode order. Neither engine may depend on a Go bounds panic to deliver a Java exception. Bound logical Java stack depth before exhausting the Go stack and produce `StackOverflowError` on the supported path.

## 6. Reflection and the platform library

Reflection resolves a member descriptor and invokes its selected implementation. On an image class, `getMethod` can find a public inherited method even if no direct caller mentions it. On a plugin class, the same lookup consults metadata created during definition. Lookup uses the actual parameter classes and Java member lookup rules. The runtime member ID also carries the return descriptor, which matters for bridge methods and covariant returns. Java access checks use the logical caller, independent of its engine. [Class member lookup](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/Class.html#getMethod(java.lang.String,java.lang.Class...))

For `Method.invoke`, check access and receiver, convert arguments under reflection rules, initialize the declaring class on a static invocation, and perform ordinary virtual selection when applicable. Box a primitive result and return null for void. Wrap a throwable from the invoked target in `InvocationTargetException` while retaining the original object as its cause. Lookup failures, adapter argument errors and initialization failures retain their own specified behavior. A descriptor found by name is sufficient to invoke a compiled function through its adapter. [Method.invoke contract](https://docs.oracle.com/en/java/javase/21/docs/api/java.base/java/lang/reflect/Method.html#invoke(java.lang.Object,java.lang.Object...))

Ship a versioned bootstrap library. Start with small Java implementations and explicit Go intrinsics for VM-owned operations such as `Object.getClass`, class mirrors, allocation, reflection adapters and thread operations. Identify every intrinsic by defining bootstrap class, method name and descriptor. Validate signature and profile at link time. Keep Java exceptions and Java string encoding at this boundary. The bootstrap closure for the first acceptance fixture includes `Object`, `Class`, `String`, `Throwable` and its exercised subclasses, `Method`, `InvocationTargetException`, the boxing classes used by the profile and the array types created by javac. Report transitive library requirements from bytecode, generated helpers and intrinsic declarations. General JDK library reuse needs an inventory of native methods and VM-private dependencies before broadening compatibility.

## 7. Threads, monitors, and the Java memory model

Assign a logical Java thread object to each executing Java thread and initially schedule it on a goroutine. Pass that object explicitly through every entry. Synchronous Java → Go → Java callbacks reuse the same logical thread so monitor reentrancy and class initialization ownership survive the boundary. A fresh asynchronous host invocation creates or attaches a distinct Java thread. Thread-local state, interruption and joins use Java thread identity.

A Java monitor records owner, recursion count, entry contention and a wait set. Entry by the owner increments depth. Exit by a non-owner throws `IllegalMonitorStateException`. A waiter must own the monitor, release its full depth, and reacquire that depth before returning or throwing on interrupt. Notification selects eligible waiters while the notifier still holds the lock. Handle timeout, interrupt and notification races using one synchronized waiter state transition. Static synchronized methods lock the class mirror. Instance synchronized methods lock the receiver. [Java threads and locks](https://docs.oracle.com/javase/specs/jls/se21/html/jls-17.html)

For the first threaded implementation, serialize all Java field and array element accesses, including ordinary accesses, under one runtime memory mutex. Use the same helpers in generated code, the interpreter, reflection and host bindings. Keep thread-local operands in Go locals. Release the memory mutex before invoking Java code, blocking on a monitor, or running an initializer. Publish thread start and termination through synchronized handshakes. This gives a deliberately strong, slow reference implementation for volatile order, constructor publication and even racy Java reference accesses without introducing Go data races. Future specialization needs an argument for preserved Java semantics plus differential and race testing. Per-object locks alone do not establish the required global synchronization order for volatile operations. [Java memory model](https://docs.oracle.com/javase/specs/jls/se21/html/jls-17.html#jls-17.4) · [Go memory model](https://go.dev/ref/mem)

## 8. Application startup, fallback and exit

The builder records the application's main class in the image. At process startup the generated Go launcher installs the image tables, creates the shared runtime and initial Java thread, constructs the Java `String[]` arguments, and initializes the main class. It invokes `public static void main(String[] args)` through the method entry defined in §5. The same entry convention applies whether the body of `main` was compiled or retained.

Consider a compiled call to `step.apply(box)`. The caller resolves the interface method and selects the implementation for the receiver. If the selected entry is interpreted, it passes the current thread, receiver and arguments directly to the interpreter. The interpreter pushes a Java frame and stores those references in its locals. If the bytecode then invokes a compiled method, it calls the selected Go entry and consumes its `Outcome` at the current bytecode offset. The original caller resumes after return or follows its exception handler. Crossing between engines requires neither serialization nor a goroutine handoff.

Reflection enters this same path after lookup and argument conversion. A name discovered at run time can resolve to a compiled method. If only a retained bytecode body exists for the selected target, the interpreter executes that body in the same process.

When `main` returns, the launcher waits for remaining non-daemon Java threads. An uncaught Java throwable follows the terminating thread's uncaught-exception path. In the threaded profile, termination of the last non-daemon thread or an explicit `System.exit(status)` starts application shutdown. Once supported shutdown hooks finish, the Go launcher terminates the process with the selected exit status. Supporting an API that registers shutdown hooks requires implementing its lifecycle in the bootstrap library.

## 9. Diagnostics and evaluation

Emit structured events carrying process ID, loader, logical thread, Java method, engine and bytecode offset where applicable. Counters include classes loaded, methods per engine, cross-engine calls, reflective invocations and lookup misses, Java exceptions, monitor contention, and live Java threads. Publish allocator and GC telemetry at its actual scope. Process RSS is a process measurement. A Java stack trace and a Go crash trace answer different operational questions and should both be available for an internal fault.

Measure the same workload in all-interpreted, all-eligible-AOT and mixed configurations, plus the target JDK as an external baseline. Record input digest, toolchain versions, target CPU/OS, Go GC settings, warmup, sample count and distributions. Measure build time, executable and image size, cold time to first successful result, steady throughput and p50/p99 latency. Isolate direct calls, virtual calls, reflection lookup, cached reflective invocation, cross-engine calls, allocation, and monitor contention. Measure peak process RSS and retained heap for the same application under each engine configuration. Report observations before setting performance targets. In particular, the global memory mutex is a correctness baseline whose cost should be measured on a field-heavy workload.

## 10. First acceptance fixture

Save the following source as `Crossing.java` and compile with `javac --release 21`. The application’s test entry point invokes `Crossing.probe()I` (method name `probe`, descriptor `()I`) and expects **2**. Force only `Crossing$Worker.apply(LCrossing$Box;)V` to interpretation in the mixed build. Use the compiled `Box.bump` method as the callback. A test host can invoke the same entry on OpenJDK to obtain the oracle result.

```java
import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

public final class Crossing {
    public interface Step {
        void apply(Box box);
    }

    public static final class Signal extends RuntimeException {}

    public static final class Box {
        int count;
        final Signal signal = new Signal();

        public void bump() {
            count++;
        }
    }

    public static final class Worker implements Step {
        public void apply(Box box) {
            box.bump();
            throw box.signal;
        }
    }

    public static int probe() throws ReflectiveOperationException {
        Box box = new Box();
        Step step = new Worker();
        try {
            step.apply(box);
            return -1;
        } catch (Signal e) {
            if (e != box.signal || box.count != 1) return -2;
        }
        Method method = Step.class.getMethod("apply", Box.class);
        try {
            method.invoke(step, box);
            return -3;
        } catch (InvocationTargetException e) {
            if (e.getCause() != box.signal) return -4;
        }
        return box.count;
    }
}
```

The mixed trace must show an AOT caller reaching interpreted `Worker.apply` twice, an AOT `Box.bump` callback both times, and the original throwable arriving first in the compiled catch and then as the reflection wrapper's cause. Also run with all fixture methods interpreted and with all fixture methods AOT. In the latter run, the reflective call must reach the AOT implementation. Choose an interpreted implementation of `Step` from a separately supplied plugin JAR for a later variant. Build the application executable before that JAR exists and obtain the receiver through the application’s plugin loader so closed-world dispatch cannot accidentally pass the test.

## 11. Delivery gates

| Gate | Deliverable | Acceptance evidence |
| --- | --- | --- |
| 1 · Vertical slice | Class reader and verifier for fixture instructions, bootstrap subset, common object and call ABI, interpreter, Go emitter and application launcher | §10 yields 2 in all three configurations. Mixed traces prove the expected engine transitions. Malformed stack maps and incompatible image revisions are rejected |
| 2 · Semantic parity | Ordinary opcode families, initialization state machine, reflection and error paths for the published profile | Differential corpus against Java 21 covers numeric edge cases, array covariance, hidden fields, constructor rules, handler order, initialization failure and reflective argument conversion. Compare return values, exception type and cause, side effects and observable order |
| 3 · Late classes | Transactional plugin admission and guarded compiled dispatch | An unseen subclass overrides an image method. Equal names from different loaders remain distinct. Failed admission leaves the running application usable. Class lookup by name reaches a method with no static callers |
| 4 · Concurrency | Reentrant monitors, wait sets, Java thread lifecycle, synchronized memory helpers and class initialization locks | Stress recursive and contested initialization, wait/notify/interrupt races, volatile publication, constructor publication and exception unwinding through synchronized calls. Run Go race detection on both engines and verify permitted outcomes against the JLS |
| 5 · Application lifecycle | Startup, main arguments, Java thread termination, exit status and process diagnostics | Launch one executable with a main method that exercises both engines. Record the same process ID and logical thread ID across nested synchronous engine transitions. Verify main arguments, non-daemon thread completion, explicit exit status, and memory and latency measurements from §9 |

Gate 1 demonstrates that a single application can pass a live object through both engines and catch the same throwable on return. Gates 2 through 5 extend the supported semantics and establish the cost of fallback in that application.
