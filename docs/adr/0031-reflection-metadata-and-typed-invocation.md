# ADR-0031: Runtime metadata, canonical type identity, and typed dynamic invocation

- **Status:** Accepted
- **Date:** 2026-07-17
- **Revised:** 2026-07-18 for ADR-0034 profile boundaries
- **Governing workstream:** `r3-reflection-dynamic-research`
- **Evidence:** `docs/workstreams/r3-reflection-dynamic-reports/{r3-java25-semantic-contract,r3-classfile-metadata-gap,r3-runtime-boundary-map,r3-engine-capability-matrix}.md`

## Context

R2 provides canonical Class mirrors, Java 25 initialization, exception
transport, and race-free heap cells. R3 baseline evidence nevertheless shows
0/24 matches in Interpreter and IR. Classfile discards annotation and member
attributes, Class-producing native paths do not all preserve mirror identity,
class loading reports absence as a Go panic, and there is no engine-neutral
dynamic value, invocation, or exception-result boundary.

ADR-0020 rejects interpreter Slot as a universal representation. ADR-0022
requires bootstrap capabilities to define kernel and facade responsibilities.
ADR-0034 now separates the Catty JVMS Core, Catty Runtime Profile, and optional
Java SE Compatibility Profiles. The R3 evidence remains valid, but exact
`java.lang.reflect` and annotation-facade behavior cannot define the core
runtime contract.

## Decision

### Profile ownership

This ADR governs three shared kernel capabilities:

1. immutable classfile/runtime metadata ownership;
2. canonical runtime type/member identity; and
3. an engine-neutral typed dynamic invocation boundary.

The Catty Runtime Profile may expose a small metadata or invocation API over
these capabilities. A Java SE Compatibility Profile may expose
`java.lang.Class`, `java.lang.reflect.*`, and annotation facades with the exact
Java SE behavior it declares. Those public APIs are parallel profile surfaces;
neither is required to implement the other.

### Metadata ownership

Classfile SHALL validate and immutably retain every attribute required by an
accepted capability. Core structural and dynamic-linkage metadata is retained
for its runtime consumers. Annotation, parameter, exception, signature, and
other discovery metadata is retained only when a declared Catty or Java SE
profile capability requires it. Parsing SHALL not load Java classes or
allocate profile-visible facade objects.

Runtime Class, Method, and Field identities SHALL retain or reference the
immutable declared metadata required by accepted consumers. Lazy resolution,
runtime identity, caches, and terminal Java failure belong to runtime services,
not classfile structs. Structurally valid unrecognized attributes are ignored
for execution as required by JVMS; an unavailable metadata view is reported as
an unsupported profile capability rather than silently becoming a reflection
claim.

### Canonical runtime identity and facades

A non-primitive runtime Class identity SHALL be its defining loader plus binary
name, consistent with ADR-0033 if accepted. Arrays, primitives, and void SHALL
have canonical internal type identities. Every engine and profile facade SHALL
refer to those shared identities rather than allocate an independent type
world.

Profile-visible Class/member/annotation objects are adapters over runtime
metadata services. Their reference identity, ordering, equality, access, and
exception behavior are governed by the profile capability that exposes them.
In particular, Java SE member objects, annotation equality/hash behavior,
`Class.forName`, and declared-member discovery are Java SE Compatibility
Profile obligations, not automatic consequences of retaining metadata.

### Typed dynamic invocation kernel

The runtime SHALL define a logical typed dynamic-value/result boundary
supporting void, every Java primitive as one full value, Java references, and
an explicit normal-result-or-Java-throwable result. It SHALL carry the current
execution/Thread context and preserve object/Class identity, initialization,
heap publication, and exception transport.

Interpreter Slot, IR register layout, Go heap-cell bits, provider ABI layout,
and current AOT bridge Slot arrays are adapters, not the boundary's stable
representation. The vocabulary may be shared with the typed Host ABI under
ADR-0034, but profile invocation and host-provider binding remain separate
services with explicit adapters and authority checks.

Core invocation SHALL provide type-checked receiver, arity, reference
assignment, primitive-value transport, dispatch, construction, field access,
initialization, and Java throwable propagation needed by accepted runtime
consumers. Profile-specific conversion, accessibility, varargs, boxing,
exception wrapping, and caller-discovery rules SHALL be layered policies.
Exact Java reflection conversions and `InvocationTargetException` are Java SE
Compatibility Profile obligations.

### Loading and initialization

The runtime SHALL provide typed, non-panicking lookup and definition results
suitable for profile adapters and linkage services. Internal must-load helpers
may panic only for Catty invariants that supported classfiles cannot trigger.

Metadata discovery does not by itself initialize a Class. Each profile API and
dynamic consumer SHALL declare its initialization triggers and invoke the
shared ADR-0025/ADR-0029 initialization service. Java SE `Class.forName` and
reflective static access use Java SE trigger and failure rules only when that
compatibility capability is declared.

### Engine boundary

Interpreter and IR SHALL share metadata, identity, lookup, initialization,
typed-value, heap, and exception services while retaining separate acceptance
evidence. AOT capabilities SHALL be reported per consumer. Until a later
accepted workstream supplies typed Thread-aware normal/exception transitions,
unsupported AOT dynamic invocation is precisely build-rejected; a generic
built-then-panic bridge is forbidden.

## Consequences

- Classfile and runtime packages gain explicit metadata ownership without
  forcing the complete Java reflection object model into Catty JVMS Core.
- Runtime Class gains defining-loader-aware identity; facades in different
  profiles can share identity without sharing public APIs.
- Reflection, dynamic linkage, generated classes, Host ABI adapters, and future
  embedding can reuse a type-correct invocation vocabulary without making Slot
  a cross-engine ABI.
- Java SE reflection fixtures remain valuable compatibility evidence, but do
  not gate completion of metadata or typed-invocation core capabilities.
- Each public facade needs a profile-specific contract for identity, ordering,
  access, conversion, initialization, and failure behavior.
- Metadata retention and resolution caches add bounded memory and
  synchronization cost only for declared consumers.

## Non-scope

Accepting a Catty Runtime Profile API; implementing `java.lang.reflect`, Java
annotation facades, generic signatures, type-use annotations, modules, records,
sealed or hidden classes, broad `setAccessible`, Java SE caller-sensitive
behavior, arbitrary custom ClassLoader APIs, serialization, JNI, MethodHandle
combinators, generated classes, InvokeDynamic linkage, or AOT dynamic
execution.

## Acceptance record

Accepted by Owner on 2026-07-18 after reconciliation with ADR-0034. Acceptance
fixes the shared metadata, canonical identity, and typed dynamic-invocation
kernel boundaries. It does not accept a Catty Runtime Profile API, a Java SE
reflection/annotation compatibility surface, or any production implementation
contract.
