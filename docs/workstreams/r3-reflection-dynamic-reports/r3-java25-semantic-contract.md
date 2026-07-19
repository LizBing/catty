# R3 Java 25 semantic contract

**Research baseline:** Temurin 25.0.3
**Acceptance anchor:** `6cf3636`
**Fixture surface:** the fixed 24-fixture R3 matrix

This report records the observable rules that an R3 implementation contract
must preserve. It is a bounded contract for the fixed matrix, not a claim of
complete Java reflection, MethodHandle, or proxy compatibility.

Primary references are the Java SE 25 API and JVMS 25:

- [Class API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/Class.html)
- [Method API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/reflect/Method.html)
- [Field API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/reflect/Field.html)
- [AnnotatedElement API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/reflect/AnnotatedElement.html)
- [Proxy API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/reflect/Proxy.html)
- [InvocationHandler API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/reflect/InvocationHandler.html)
- [java.lang.invoke package](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/invoke/package-summary.html)
- [LambdaMetafactory API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/invoke/LambdaMetafactory.html)
- [StringConcatFactory API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/invoke/StringConcatFactory.html)
- [JVMS §4 classfile format](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-4.html)
- [JVMS §5.4 resolution](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-5.html#jvms-5.4)
- [JVMS invokedynamic instruction](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-6.html#jvms-6.5.invokedynamic)

## 1. Class identity, loading, and initialization

- One defining loader plus one binary name identifies a class or interface.
  Repeated successful lookup through the same defining context returns the
  same runtime Class and canonical `java.lang.Class` object. Primitive and
  array Class mirrors are canonical too.
- `Object.getClass`, class literals, `Class.forName`, member type queries, and
  descriptor/MethodType resolution must converge on that same mirror. A native
  helper may not allocate an observably distinct mirror.
- `Class.forName(String)` loads, links, and initializes the named class and
  throws Java `ClassNotFoundException` when it cannot be located. The
  three-argument form uses the supplied loader and initializes only when its
  boolean argument is true. Array lookup loads but does not initialize the
  element type merely because `initialize` is true.
- Metadata discovery (`X.class`, `getDeclaredField`, `getDeclaredMethod`) does
  not by itself initialize the declaring class. Reflective access to a static
  field and reflective invocation of a static method initialize the actual
  declaring class through the accepted ADR-0025/ADR-0029 service.
- Loading, linkage, initialization, and target exceptions remain distinct Java
  failure categories. A Go panic is never the Java-visible result.

## 2. Class and member discovery

- Declared-member APIs return members declared directly by the represented
  class/interface, including non-public and compiler-generated members where
  specified, and exclude inherited members. Public lookup APIs apply hierarchy
  and interface rules separately.
- The Java API does not promise a stable ordering for returned member arrays.
  Fixtures therefore query members by exact name and parameter type and only
  assert counts where the source fixes them.
- A reflected Method/Field/Constructor has stable logical equality based on
  declaring Class, name, descriptor/signature, and member kind. The runtime may
  cache facade identity, but Java code must not depend on reference identity of
  independently returned member facade objects unless a later contract proves
  that requirement.
- Primitive, void, array, interface, superclass, component-type, modifier, and
  assignability queries return canonical Class mirrors and Java modifier bits.

## 3. Reflective invocation and field access

- Instance Method invocation validates a non-null compatible receiver and then
  uses ordinary virtual dispatch, so an override selected from the receiver's
  runtime class executes. Static Method invocation ignores the receiver and
  initializes the declaring class before execution.
- Constructor invocation allocates one instance of the declaring class, applies
  access and argument checks, invokes `<init>`, and returns the initialized
  object. Abrupt constructor completion produces the specified reflection
  exception rather than a partially successful result.
- Primitive arguments are unboxed and may undergo Java method-invocation
  widening; narrowing is not silently performed. Reference arguments must be
  assignment-compatible. Varargs are represented by the final array argument;
  Method.invoke does not invent an additional language-level varargs packing
  rule beyond its Object-array call surface.
- Primitive results and field values are boxed. A void Method returns null.
  Primitive arrays remain primitive arrays rather than arrays of wrappers.
- Instance Field access validates the receiver. Static Field access ignores the
  receiver and initializes the field's declaring class. Writes use identity or
  permitted unboxing+widening conversions and then the existing race-free heap
  cell API.
- Access failure throws `IllegalAccessException`; incompatible receivers,
  arity, unboxing, or conversion throw `IllegalArgumentException`; null instance
  receivers throw `NullPointerException`. An exception thrown by the invoked
  method or constructor is transported as the cause of
  `InvocationTargetException` without changing its Java object identity.

## 4. Runtime-visible annotations

- Only attributes with runtime visibility participate in the fixed reflection
  surface. Class, field, method, and parameter declaration annotations retain
  their annotation type and element-value tree losslessly.
- Element values cover primitives, String, Class, enum constants, nested
  annotations, and arrays. Missing explicit values are obtained from the
  annotation method's `AnnotationDefault` attribute.
- An annotation facade implements its annotation interface plus the Annotation
  contract. Its element methods return defensive array copies where required;
  `annotationType`, logical equality, hash code, and string form are
  Java-visible behavior even if the first implementation matrix exercises only
  a bounded subset.
- `@Inherited` affects class queries through superclasses only. It does not
  inherit annotations from interfaces or apply to non-class elements.
- Repeatable lookup unwraps the declared container annotation according to
  `@Repeatable`; directly, indirectly, present, associated, and declared
  queries must not be collapsed into one undifferentiated lookup.

## 5. InvokeDynamic and call-site linkage

- Each invokedynamic instruction is a distinct resolution site even when two
  instructions refer to the same constant-pool entry. Its symbolic reference
  supplies the call-site name and MethodType; the BootstrapMethods entry
  supplies a bootstrap MethodHandle and static arguments.
- Resolution obtains canonical Class/MethodType/MethodHandle values, invokes
  the bootstrap method, validates a non-null CallSite and a target MethodHandle
  whose type equals the call-site descriptor, and binds the instruction to
  that CallSite.
- Successful resolution is published once for that instruction. Concurrent
  executions observe the same completed linkage result. Bootstrap execution
  may race internally only as permitted by JVMS resolution; catty must publish
  one Java-visible result and cannot expose partially initialized linkage
  state.
- A failed resolution is memoized per instruction. Every later execution fails
  with the same resolution error and does not execute the bootstrap again, as
  required by JVMS §5.4.3.6. Throwable reference identity is not promised; the
  fixed Temurin reference produces distinct `NoClassDefFoundError` objects.
- Invocation pops arguments according to the symbolic MethodType and invokes
  the current CallSite target. A ConstantCallSite target is permanent. Mutable
  and volatile CallSite semantics are outside the first implementation unless
  a later accepted slice explicitly includes them.
- Signature-polymorphic MethodHandle invocation and adaptation require typed
  logical values. Interpreter frame Slot layout is not the stable API or
  cross-engine ABI.

## 6. String concatenation and lambdas

- StringConcatFactory linkage receives the call-site MethodType, recipe, and
  constants. Invocation converts arguments in Java order and produces the
  Java String result without losing UTF-16 code units. The bounded fixture
  covers references, null, int, long, literals, and supplementary characters.
- LambdaMetafactory has three distinct phases: linkage produces a factory
  CallSite; capture invokes that factory and may allocate or reuse a function
  object; interface invocation calls the implementation MethodHandle with
  captured and invocation arguments after required adaptation.
- Lambda object reference identity is intentionally unspecified. Tests may
  assert implemented interface, runtime class reuse observed in the pinned
  reference, and behavior, but implementation contracts must not promise
  singleton object identity for stateless lambdas.
- Captures preserve Java values and object identity. Method references obey
  their resolved invocation kind and ordinary null/access/error behavior.

## 7. Dynamic proxies

- A proxy class is defined for an ordered interface list and defining loader,
  cached according to an explicit runtime rule, implements those interfaces,
  and associates each instance with exactly one InvocationHandler.
- Interface calls are boxed into the handler boundary; primitive results are
  unboxed and checked before return. `equals`, `hashCode`, and `toString` are
  dispatched to the handler as Object-declared Method facades; other Object
  methods retain ordinary Object behavior.
- A RuntimeException or Error from the handler propagates. A checked exception
  assignable to a declared interface exception propagates. Any other checked
  exception is wrapped in `UndeclaredThrowableException`.
- Default interface methods reach their Java body only through the explicit
  InvocationHandler `invokeDefault` operation, with invokespecial-equivalent
  selection on an interface implemented by the proxy class.
- Proxy class, Method facade, handler, argument, and thrown Java object
  identities remain in the shared Interpreter/IR/AOT object world.

## 8. Engine and concurrency contract

- Interpreter and IR must produce byte-identical output and equal exit status
  to Temurin for every fixture claimed Supported. IR may share the same dynamic
  services; it may not dereference a missing lowered target or approximate
  unsupported operations.
- AOT may use explicit Interpreter fallback only after the typed call,
  exception, initialization, Thread-context, and object-identity transition is
  accepted and tested. Otherwise the relevant fixture is build-rejected as
  Not implemented. A built binary that panics is a failure.
- Class/member/annotation facade caches, call-site resolution state, and
  generated-class caches must be race-free. Publication uses existing R2 SC
  heap and initialization guarantees plus explicit internal synchronization;
  Java monitor state is not reused as a VM metadata lock.

## 9. Explicit first-R3 boundary

The fixed contract does not establish generic signatures, type-use
annotations, module access, records, sealed classes, hidden classes, arbitrary
MethodHandle combinators, ConstantDynamic execution, mutable/volatile
CallSites, arbitrary bootstrap methods, serialization, custom ClassLoader API
compatibility, or complete java.lang.reflect behavior. Those remain
Not implemented unless an accepted later contract adds exact fixtures and
engine gates.
