# ADR-0033: Defining-loader identity and generated-class kernel

- **Status:** Accepted
- **Date:** 2026-07-17
- **Revised:** 2026-07-18 for ADR-0034 profile boundaries
- **Governing workstream:** `r3-reflection-dynamic-research`
- **Evidence:** `docs/workstreams/r3-reflection-dynamic-reports/{r3-java25-semantic-contract,r3-runtime-boundary-map,r3-engine-capability-matrix,r3-slice-plan}.md`

## Context

The standard ClassLoader currently owns a concurrency-safe name cache, but
runtime Class does not retain defining-loader identity and the Loader interface
only supports must-load-by-name. Synthetic classes can be assembled through Go
helpers, but there is no typed, loader-aware, race-free class-definition
protocol.

Defining-loader identity is a core class-identity requirement. Several future
consumers may also need runtime-created classes, including Catty Runtime
services, lambda factories, dynamic proxies, and later embedding facilities.
Those consumers can share safe class-definition machinery without sharing
cache keys, generated method bodies, public APIs, or compatibility contracts.

ADR-0034 makes Java SE LambdaMetafactory and Proxy optional compatibility
features. They cannot define the generated-class kernel or make their exact
dispatch/cache behavior a Catty JVMS Core obligation.

## Decision

### Class and loader identity

Every non-primitive runtime Class SHALL retain one explicit defining-loader
identity. A Class identity is defining loader plus binary name. Arrays derive
their defining context from component-type rules, and primitives/void belong
to canonical VM type identities.

The loader service SHALL distinguish lookup/initiation from definition. A
lookup may return a Class defined by another loader after delegation. A
definition atomically installs one Class in its defining loader or returns a
typed Java linkage failure; no caller may observe a partially linked Class.
Class lookup or definition failures reachable from a supported classfile SHALL
not escape as Go panics.

This identity and definition contract belongs to Catty JVMS Core. Public
`java.lang.ClassLoader`, arbitrary byte-array definition, module/package rules,
and caller-sensitive Java APIs belong to declared compatibility profiles.

### Generated-class kernel

The runtime SHALL provide an internal, loader-owned generated-class definition
kernel for accepted consumers. It SHALL construct ordinary runtime
Class/Method/Field identities that participate in canonical type identity,
assignability, initialization, heap storage, monitors, Thread context,
exceptions, and Interpreter/IR dispatch.

Generated provenance and executable bodies may use runtime-native kernels, but
the resulting runtime Class and objects SHALL obey the same declared semantics
as other Classes in their profile. Generated bodies require an explicit typed
executable form understood by participating engines; they SHALL NOT depend on
mutable fabricated classfile bytes or expose host Go types as Java identity.

The generated-class kernel is not a public arbitrary ClassLoader API. Each
consumer SHALL declare its profile, defining-loader selection, generated name
policy, cache key, cache lifetime, concurrency/publication behavior,
initialization triggers, method-body form, and terminal failures in its own
accepted capability contract.

### Consumer separation

Catty Runtime Profile consumers and Java SE Compatibility Profile consumers
may share the kernel but SHALL maintain independent public contracts and
caches.

In particular:

- lambda implementation classes and capture objects are governed by a
  separately declared lambda/metafactory compatibility capability;
- Proxy class caching, ordered interface keys, InvocationHandler dispatch,
  Object methods, default methods, and exception wrapping are governed by a
  separately declared Java SE Proxy capability; and
- an internal Catty-generated class is not evidence that either Java SE
  capability is supported.

Consumer caches SHALL be race-free and publish only fully constructed Classes.
The kernel SHALL not impose one universal cache key or object-singleton policy
on all consumers.

### Engine boundary

The first generated-class capability SHALL name its participating engines.
Interpreter and IR may share the same generated Class and typed executable
bodies while retaining separate evidence. AOT construction or execution of
runtime-generated classes remains `Not implemented` and precisely
build-rejected until an Accepted design defines direct generation or a typed
engine transition preserving values, exceptions, initialization, and Thread
context.

## Consequences

- Loader and Class gain explicit identity, lookup, and atomic-definition
  responsibilities that broad custom loading can later extend.
- Generated-class infrastructure can support Catty and compatibility profiles
  without importing Lambda or Proxy policy into the core runtime.
- Lambda and Proxy evidence remains useful for their future Java SE profile
  contracts but does not gate defining-loader or generated-class kernel work.
- Each generated-class consumer must specify cache, lifecycle, dispatch, and
  failure semantics; there is intentionally no universal generated-class
  policy.
- Cache lifetime may retain generated Classes for at least the defining
  loader's supported lifetime; unloading and weak caches remain deferred.

## Non-scope

Accepting a public Catty generated-class API; arbitrary Java ClassLoader
subclasses or byte-array `defineClass`; loader unloading; weak caches; hidden
or anonymous classes; modules, packages, or sealing; agents; instrumentation;
serialization; LambdaMetafactory; Java lambda compatibility; dynamic Proxy;
InvocationHandler; broad `altMetafactory`; or AOT generated-class execution.

## Acceptance record

Accepted by Owner on 2026-07-18 after reconciliation with ADR-0034. Acceptance
fixes defining-loader identity, typed atomic class definition, and the shared
generated-class kernel. It does not accept LambdaMetafactory, dynamic Proxy,
arbitrary ClassLoader compatibility, or production implementation.
