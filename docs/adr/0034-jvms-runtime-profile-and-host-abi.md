# ADR-0034: JVMS core, runtime profiles, and typed host interoperability

- **Status:** Accepted
- **Date:** 2026-07-18
- **Supersedes:** [ADR-0017](./0017-java-25-semantic-baseline.md)

## Context

catty began as an experimental JRE-oriented project and currently uses a
real `java.base` image for a bounded compatibility smoke path. Its accepted
runtime decisions already separate Java semantics from Go implementation
mechanisms: AOT and the Interpreter share one Java object world (ADR-0016), Go
infrastructure is not a semantic waiver (ADR-0018), Go-native representations
may dissolve behind Java-visible contracts (ADR-0019), and `Slot` is not a
general ABI (ADR-0020).

The intended product direction is narrower and cleaner: a JVMS-compatible
runtime platform with its own small runtime profile and native host services.
It should be able to offer selected Java SE APIs as compatibility layers, but
must not accidentally promise a complete JRE, an OpenJDK-derived `java.base`,
or JNI. The existing native registry also needs a durable policy: an
unimplemented native method cannot silently return a zero or null value.

ADR-0017 correctly requires honest, capability-scoped semantics, but makes
Java SE API contracts part of the universal precedence chain. That is too
broad for a platform whose core contract is JVMS and whose standard-library
surface can be intentionally profile-specific.

## Decision

### Product contract and profiles

catty SHALL describe every externally visible capability by both engine state
and runtime profile. It defines three distinct scopes:

1. **Catty JVMS Core** — execution of the explicitly supported JVMS 25
   classfile and runtime semantics. JVMS is authoritative; JLS is consulted
   only where a claimed source-language behavior requires it.
2. **Catty Runtime Profile** — a small, versioned Catty-owned set of
   Java-visible bootstrap facades and APIs needed to expose the core and host
   services. Its API contract is specified by Catty ADRs, workstreams, and
   versioned documentation, not presumed Java SE compatible.
3. **Java SE Compatibility Profiles** — opt-in, individually versioned
   compatibility layers that provide named Java SE packages/APIs and obey the
   applicable Java SE API contract for their declared surface. An extracted
   OpenJDK `java.base` image may participate only in such a profile or in
   differential testing.

No profile is implied by another. In particular, loading a real `java.base`,
passing a Java-source fixture, or offering a class with a JDK name does not
claim Java SE or JRE compatibility beyond the declared profile capability.
The Catty Runtime Profile and Java SE Compatibility Profiles are parallel API
profiles over the shared Catty JVMS Core and typed Host ABI; neither profile's
public API is a required implementation dependency of the other.

### Semantic precedence and deviations

For Catty JVMS Core claims, precedence is JVMS 25, then any relevant JLS 25
rule, then Catty's accepted documentation for deliberately exposed profile
behavior. For Java SE Compatibility Profile claims, Java SE 25 API contracts
are additionally authoritative; pinned Temurin 25 remains the differential
reference for implementation-permitted choices.

Unsupported classfile behavior, APIs, host services, and native bindings SHALL
be rejected or reported explicitly. Go behavior, test convenience, or a
default return value is never an implicit semantic deviation. Intentional
deviations from a claimed profile require a separate Accepted ADR with the
observable behavior, affected scope, evidence, benefit, and review condition.

### Typed Host ABI and native linkage

catty SHALL expose host functionality through an internal, typed Host ABI.
Provider calls receive logical Java values/references and execution context,
and return either a typed normal result or an explicit Java throwable result.
`rtda.Frame`, interpreter `Slot`, Go pointer layout, and AOT bridge layout are
adapters and SHALL NOT be the stable provider ABI.

Host providers may implement Catty Runtime Profile services or Java SE
compatibility facades. They do not require JNI compatibility, dynamic native
library loading, or an operating-system ABI shared with existing JVMs. JNI is
outside the default product contract.

A native method with no applicable provider SHALL fail explicitly at binding or
invocation through a Java-visible unsatisfied-link/unsupported-capability
failure. It SHALL NOT execute a successful-looking zero/null stub. The exact
throwable mapping and binding lifecycle require an implementation workstream.

### Class-library and bootstrap policy

The bootstrap kernel remains capability-based under ADR-0022. A Java-visible
facade may be Catty-supplied, host-backed, synthetic, or provided by an opted-in
compatibility profile, provided it preserves the semantics claimed for that
profile. A real `java.base` image is neither a mandatory bootstrap dependency
nor the permanent default definition of the Catty Runtime Profile.

R3 decisions SHALL separate JVMS dynamic linkage and generated-class identity
from optional Java SE reflection, annotation, lambda, and proxy facades. An
implementation workstream must state the applicable profile for each claimed
R3 capability.

## Consequences

- Catty can pursue a small runtime and host-service model without inheriting
  Java SE/JRE or JNI as a blanket compatibility obligation.
- Current Java 25 differential fixtures and real-`java.base` smoke tests stay
  valuable evidence, but prove only their named compatibility scope.
- Existing native registrations can be migrated incrementally behind adapters;
  missing bindings become explicit failures rather than fabricated results.
- Roadmap and workstreams must label profile scope in addition to Interpreter,
  IR, and AOT capability state.
- ADR-0031 through ADR-0033 require reconciliation before acceptance so that
  their Java SE facade obligations are expressly optional compatibility work.
- This ADR supersedes ADR-0017's universal Java SE API precedence while
  preserving its requirements for explicit support claims, differential
  evidence, and documented deviations.

## Non-scope

Implementing the Host ABI, changing the launcher default, removing current
`java.base` auto-detection, providing a Catty API design, executing arbitrary
Java SE libraries, Java SE certification/JCK conformance, JNI, FFM API
compatibility, dynamic native-library loading, graphics/I/O providers, or
changing any existing accepted workstream contract.

## Acceptance record

Accepted by Owner on 2026-07-18. Acceptance establishes the profile and Host
ABI direction and supersedes ADR-0017. It does not authorize production
implementation; non-trivial profile, native-boundary, launcher, or R3 changes
require an Accepted workstream.
