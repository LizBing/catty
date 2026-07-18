# R3 classfile and metadata gap

**Repository:** acceptance anchor `6cf3636`
**Reference format:** [JVMS 25 §4](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-4.html)

## Summary

Catty parses enough structural metadata for execution, but not enough lossless
metadata for reflection or dynamic linkage. The parser currently materializes
only `Code` and nested `StackMapTable`; every other attribute becomes an empty
`UnparsedAttribute`. MethodHandle, MethodType, and InvokeDynamic constant-pool
entries are recognized structurally, but their operands cannot be resolved by
runtime consumers. ConstantDynamic is rejected as an unknown tag.

The first implementation slice must preserve metadata without making
classfile structs the permanent Java facade or runtime ABI. Parsed immutable
metadata belongs to the classfile domain; `rtda` owns resolved runtime identity
and lazy facade/linkage caches.

## Attribute inventory

| Attribute | Anchor state | Fixed-matrix consumer | Required first-R3 state |
|---|---|---|---|
| `Code` | Parsed | ordinary execution, generated lambda bodies | Retain unchanged |
| `StackMapTable` | Parsed | IR lowering/type tracking | Retain unchanged |
| `BootstrapMethods` | Discarded | all five InvokeDynamic fixtures | Parse losslessly at ClassFile scope; expose bootstrap handle index and ordered static-argument indexes |
| `RuntimeVisibleAnnotations` | Discarded | class, field, method annotations | Parse declaration annotation trees on ClassFile/field/method |
| `RuntimeVisibleParameterAnnotations` | Discarded | `MemberAnnotation` | Parse per-parameter ordered annotation sets |
| `AnnotationDefault` | Discarded | `AnnotationDefaults` | Parse element-value tree on annotation-interface methods |
| `MethodParameters` | Discarded | future parameter names/modifiers | Preserve as optional metadata; names are not required by the fixed output |
| `Exceptions` | Discarded | proxy declared-exception policy, reflected exception types | Parse ordered checked-exception class indexes on methods |
| `InnerClasses` | Discarded | broader declaring/enclosing queries | Preserve for a later bounded query slice; not required by fixed output |
| `EnclosingMethod` | Discarded | broader local/anonymous class queries | Preserve later; not required by fixed output |
| `Signature` | Discarded | generic reflection | Explicitly Not implemented in first R3 |
| `RuntimeVisibleTypeAnnotations` | Discarded | type-use reflection | Explicitly Not implemented in first R3 |
| `ConstantValue` | Discarded | static constant materialization | Existing execution boundary must be audited; reflective access must observe the same initialized storage |
| `NestHost` / `NestMembers` | Discarded | modern access checks | Required if access rules include nestmate private access; otherwise first reflective access slice must exclude and build fixtures accordingly |
| `Record`, `PermittedSubclasses`, module attributes | Discarded | records/sealed/modules | Explicitly Not implemented |

Unknown attributes must continue to be safely skipped by declared length. A
known attribute that the accepted runtime consumes must never be silently
converted to an empty placeholder.

## Annotation element-value representation

The parser needs an immutable tagged value tree matching JVMS §4.7.16.1:

| Tag family | Stored payload | Runtime resolution |
|---|---|---|
| primitive/String constants | constant-pool index plus declared tag | Java wrapper/String value |
| enum | type-name and constant-name indexes | canonical enum Class and constant |
| Class | descriptor index | canonical primitive/reference/array Class mirror |
| nested annotation | annotation tree | annotation facade implementing nested interface |
| array | ordered child values | defensive Java array result |

The raw tree must not eagerly load annotation types during classfile parsing.
Resolution occurs in a Java execution context so loading failures and class
initialization behavior can use Java exception transport.

## Constant-pool inventory

| Entry | Anchor state | Required first-R3 state |
|---|---|---|
| `CONSTANT_MethodType` (16) | Tag and descriptor index parsed; no descriptor accessor or pool association | Expose validated descriptor and lazily resolve canonical MethodType |
| `CONSTANT_MethodHandle` (15) | Reference kind/index parsed; no accessors/resolution | Expose kind/index; validate JVMS kind/tag rules; lazily resolve direct typed handle with access context |
| `CONSTANT_InvokeDynamic` (18) | Bootstrap index/name-and-type index parsed; no pool association/accessors | Expose bootstrap index, name, descriptor; associate each bytecode instruction with per-site linkage state |
| `CONSTANT_Dynamic` (17) | Unknown tag panic | Parse structurally so Java 25 classfiles do not crash the parser; execution remains Not implemented unless separately accepted |
| `CONSTANT_Module` / `CONSTANT_Package` (19/20) | Unknown tag panic | Parse structurally or reject with typed unsupported diagnostics when broader module classfiles enter scope; no first-R3 module semantics |

InvokeDynamic resolution cannot be stored solely on the constant-pool entry:
JVMS resolution is per invokedynamic instruction, and different instructions
may reference the same entry while retaining separate resolution state.

## Runtime metadata inventory

| Runtime type | Retained at anchor | Missing for fixed matrix |
|---|---|---|
| `rtda.Class` | name, super, interfaces internally, access flags, CP, declared fields/methods internally, init state, canonical Class mirror | defining loader, canonical primitive Class model, exported declared views, annotation trees/caches, member facade cache, generated-class provenance |
| `rtda.Method` | owner, name, descriptor, flags, bytecode, handlers, stack map, native target | declared exceptions, parameter metadata/annotations, annotation/default metadata, constructor/member kind facade, typed dynamic invocation entry |
| `rtda.Field` | owner, name, descriptor, flags internally, heap cell id | public access flags, annotation metadata, typed reflection get/set service |
| `rtda.Object` | runtime Class, SC heap cells, native payload, monitor | typed payload contracts for Class/member/annotation/MethodHandle/CallSite/proxy facades |
| loader | provider chain, concurrency-safe name cache per loader object | non-panicking lookup result, defining-loader identity on Class, define-generated-class operation, loader-aware proxy/lambda cache |

## Canonicality audit

`rtda.Class.ClassObject` is the correct canonical mirror service. The following
anchor paths bypass it and allocate fresh Class objects:

- `native.classGetSuperclass`;
- `native.classGetPrimitiveClass`.

All R3 Class-producing paths must route through one canonical service. The
primitive model also needs one canonical runtime Class per primitive/void type;
mapping `int.class` to the int-array Class is not a valid long-term identity.

Member and annotation facade canonicality should be intentionally weaker:
cache resolved metadata and underlying identity, but do not promise stable
facade reference identity unless the Java 25 contract or measured compatibility
requires it.

## Ownership and mutation rules

- Parsed attribute/constant-pool metadata is immutable after ClassFile parse.
- `rtda.Class` owns resolved Java identities, facade caches, annotation
  resolution caches, and per-site dynamic linkage state.
- Loader owns name/definition uniqueness and generated-class caches.
- Java heap fields remain in ADR-0030 HeapCells; reflection never receives a
  mutable backing slice.
- Lazy caches use race-free publication and cache Java-visible terminal failure
  when the Java specification requires repeatable failure behavior.
- No parsed byte slice, constant-pool internal index, interpreter Slot, or Go
  function pointer becomes a Java-visible or cross-engine stable ABI.

## Minimum parser gates

The metadata implementation slice should not claim completion until unit tests
cover:

1. every element-value tag including nested arrays and annotations;
2. annotations at class, field, method, and parameter locations;
3. annotation defaults and declared exceptions;
4. BootstrapMethods with zero/multiple static arguments;
5. MethodHandle kinds 1–9, MethodType, InvokeDynamic name/descriptor, and
   structural ConstantDynamic parsing;
6. malformed indexes/tags/lengths failing with typed parser errors rather than
   slice bounds panic;
7. immutability and no mutation of historical classfile execution behavior.
