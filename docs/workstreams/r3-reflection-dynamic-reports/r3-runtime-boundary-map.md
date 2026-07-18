# R3 runtime boundary map

**Repository:** acceptance anchor `6cf3636`
**Purpose:** identify ownership and call boundaries before production design

## Package responsibilities

| Package/domain | Anchor responsibility | R3 research conclusion |
|---|---|---|
| `classfile` | parse executable class structure, Code, StackMapTable, selected CP tags | retain immutable annotation/bootstrap/member metadata; never load classes or allocate Java facades |
| `classloader` | provider chain and one concurrency-safe name cache per loader object | provide typed lookup failure and loader-aware definition uniqueness; own generated-class definition/cache policy |
| `rtda` | runtime Class/Method/Field/Object identity, heap, Thread, init, monitors | own resolved metadata identity, canonical mirrors, typed dynamic values, facade/linkage state contracts |
| `native` | synthetic/bootstrap facades and native registrations | thin Java API facades over shared R3 services; must not become a second reflection/linkage implementation |
| `interpreter` | bytecode execution and Java exception transport | trigger shared reflection/dynamic services with explicit Thread; execute resolved targets and opcode 0xba |
| IR executor | validation executor over lowered IR | use the same resolved services and typed values; unsupported lowering must be explicit, never nil-target dereference |
| `lowering` | bytecode to typed register IR | model InvokeDynamic as a typed dynamic call or explicit service call; IR schema remains replaceable under ADR-0016 |
| `runtime` AOT bridge | mixed AOT/interpreted calls using global R1 loader/Thread and Slot adapters | cannot serve as the R3 ABI until Thread context, typed category-2 returns, and Java exception propagation are explicit |
| `transpile` | reachability and Go emission | classify each R3 family as direct, explicit fallback, or build rejection; no build-then-panic |

## Class lookup and mirror flow

```text
Java Class.forName / descriptor / ldc / reflection query
        |
        v
typed loader lookup(name, initiating/defining context)
        |-- failure --> Java ClassNotFoundException / LinkageError
        v
unique runtime Class (defining loader + binary name)
        |
        +--> ADR-0025 initialization when the operation is a trigger
        |
        +--> canonical Class mirror (one Object identity)
        |
        +--> immutable declared metadata + lazy facade caches
```

The anchor's `Loader.LoadClass(string) *Class` panic contract is insufficient
for Java-visible lookup. The production design should preserve a convenient
must-load helper for internal invariant paths while adding a typed lookup path
that returns a Java-mappable failure. It must not convert every loader caller
to recover from Go panic.

Runtime Class must retain defining-loader identity. An initiating-loader cache
may delegate and return a Class defined by another loader, so Java Class
identity cannot be inferred from the cache map that happened to return it.

## Reflection discovery flow

```text
canonical Class mirror
        v
resolve declared metadata without initialization
        v
apply member-kind + name + descriptor/hierarchy rules
        |-- absent --> NoSuchFieldException / NoSuchMethodException
        v
materialize Java Field / Method / Constructor facade
        v
facade payload references immutable runtime metadata identity
```

Member facade payloads must not point at classfile byte slices or rely on Go
interface type assertions scattered through native methods. One shared R3
service validates payload kind and routes all facade methods.

## Typed dynamic invocation boundary

Reflection, MethodHandle, CallSite, lambda capture, proxy dispatch, and eventual
AOT fallback all need the same logical value vocabulary, but not necessarily
the same adapters or public API.

The minimum vocabulary is:

- void;
- boolean/byte/char/short/int/long/float/double as full language-level values;
- Java object/array reference including null;
- declared JVM type/descriptor for validation and conversion;
- normal result or Java throwable result, never an ambient Go panic.

Adapters are domain-local:

```text
Java Object[] / wrapper objects <-- reflection/proxy adapter
             typed dynamic values
Interpreter locals/stack       <-- frame adapter
IR typed registers             <-- IR adapter
AOT Go values                  <-- future typed bridge adapter
```

This honors ADR-0020: interpreter Slot remains an implementation detail. A
single Slot cannot safely represent category-2 values at a stable boundary,
and the current AOT bridge already documents this limitation.

## Reflective execution flow

```text
Method/Constructor/Field facade + caller Thread
        v
validate payload, receiver, access, arity, conversions
        |-- validation failure --> Java reflection exception
        v
request declaring-class initialization when required
        |-- initialization failure --> existing Java init transport
        v
read/write HeapCell OR invoke resolved target
        v
adapt result to wrapper/reference
        |-- target threw --> InvocationTargetException(target throwable)
        v
return through ordinary Java frame
```

Access checking needs an explicit caller Class/context. Inferring the caller
from a Go stack is invalid. Native facade entry can identify its Java caller
through the current execution context/frame chain; later AOT fallback must pass
that context explicitly.

## InvokeDynamic linkage flow

```text
bytecode method + opcode PC
        v
InvokeDynamic CP entry --> name + descriptor + bootstrap index
        v
Class BootstrapMethods --> bootstrap handle + static CP arguments
        v
per-instruction linkage state on runtime Method
  unresolved | resolving(owner) | linked(CallSite) | failed(Error)
        v
resolve MethodType/MethodHandles with caller lookup context
        v
invoke bootstrap through typed dynamic invocation boundary
        v
validate CallSite + target type
        v
publish linked/failed terminal state to all waiting Threads
        v
invoke current target with instruction arguments
```

The per-site lock/condition is VM metadata synchronization, not the Java
monitor of Class, Method, or CallSite objects. Same-Thread recursive linkage
and cross-Thread wait behavior require explicit rules in the accepted ADR and
unit tests. Successful and failed terminal states are immutable for
ConstantCallSite-first support.

## Lambda and proxy generation boundary

LambdaMetafactory and Proxy both may generate classes, but they have different
keys and semantics:

| Mechanism | Definition/cache key | Instance state | Dispatch |
|---|---|---|---|
| lambda factory | caller lookup Class + invokedynamic site + metafactory arguments | captured values; allocation/reuse unspecified | generated SAM method to implementation MethodHandle |
| proxy | defining loader + ordered interface identities + proxy options | one InvocationHandler | generated interface/Object methods to handler |

Both require a loader-owned `defineGeneratedClass`-equivalent service that:

- creates one runtime Class identity with explicit defining loader;
- links super/interfaces/method tables and allocates canonical Class mirror;
- participates in ordinary assignability, monitor, initialization, heap, and
  exception semantics;
- records generated provenance without pretending a source classfile exists;
- publishes cache entries race-free and never returns a partially built Class.

The service does not authorize Java ClassLoader.defineClass, hidden classes, or
arbitrary byte-array class definition in the first R3 boundary.

## Existing defects exposed by the baseline

These are research findings, not authorized fixes:

- Class-producing native paths for superclass and primitive lookup bypass the
  canonical Class mirror service; `ClassQueries` already observes a false
  superclass identity before failing on the next missing method.
- missing Class facade methods become `NoSuchMethodError` in Interpreter, while
  IR often dereferences a nil resolved method and panics. R3 gates must require
  equal Java-visible failure transport across the two engines.
- Interpreter reports unsupported invokedynamic explicitly; IR lowering also
  reports it explicitly. AOT emits only generic unsupported/no-build output,
  not a deliberate R3 capability classification.
- the AOT bridge owns one package-global Thread and loader. It cannot carry
  reflective caller context or concurrent dynamic linkage safely under the R2
  execution-context contract.

## Required concurrency tests

Later implementation contracts should cover:

1. concurrent canonical Class/member/annotation facade lookup;
2. concurrent successful InvokeDynamic linkage with one published CallSite;
3. concurrent failed linkage with one published failure result and no repeated
   bootstrap execution;
4. recursive bootstrap behavior without Go deadlock;
5. concurrent proxy/lambda class cache lookup with one runtime Class identity;
6. reflective static access racing ordinary initialization;
7. handler/target exceptions retaining Java object identity across engine
   boundaries.
