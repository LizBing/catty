# ADR-0016: Multi-engine execution with AOT as the primary product path

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** [ADR-0007](./0007-reflection-dynamic-tiered.md)

## Context

catty has three execution-related assets with different long-term purposes:

- a bytecode interpreter that can execute Java behavior without first lowering
  every path to Go;
- an IR used for analysis, validation, lowering, and optimization; and
- an AOT path that translates Java business code into Go source and lets the Go
  toolchain produce native code.

ADR-0007 described these as a reflection-driven tiered system in which dynamic
features live in the interpreter and statically resolvable hot code lives in
AOT. That retained the interpreter, but it coupled several independent choices:
reflection policy, engine selection, runtime promotion, and the long-term AOT
architecture. It also suggested a linear interpreter-to-IR-to-AOT progression.

That progression is not intrinsic to catty. AOT normally happens at build time,
IR can serve several consumers without being a runtime tier, and reflection is
only one reason that a path may require runtime resolution or fallback.

catty's product thesis is stronger and more specific: while preserving the
declared Java semantics, as much stable business logic as practical should
dissolve into Go-native control flow and data flow. The interpreter remains
necessary so that incomplete AOT coverage or dynamic behavior does not require
unsafe compilation or silent semantic approximation.

## Decision

catty adopts a **multi-engine execution architecture**:

1. **AOT is the primary product path.** Its purpose is to translate as much
   eligible Java business code as practical into Go-native code, removing
   bytecode dispatch, simulated operand-stack traffic, and avoidable runtime
   representation boundaries. The Go toolchain supplies downstream
   optimization and native-code generation.
2. **The interpreter is a permanent semantic fallback and reference execution
   path.** It is not a temporary development tier to be deleted when AOT
   coverage grows. Unsupported AOT lowering, runtime-loaded code, and other
   explicitly dynamic paths may execute through it.
3. **IR is shared compiler infrastructure, not a mandatory runtime tier.** It
   may support verification, control-flow and data-flow analysis, semantic
   lowering, AOT generation, optimization, an IR executor, and a possible
   future JIT. The current IR schema and executor are replaceable
   implementations, not a stable cross-engine ABI.
4. **A future JIT is permitted but not committed.** It may consume the shared IR
   if evidence justifies it. This ADR does not authorize a specific JIT design,
   runtime compiler, or promotion policy.
5. **Engine selection is policy, not Java-visible semantics.** Build-time
   closure analysis, available lowering, dynamic loading, reflection,
   `invokedynamic`, profiles, code-size constraints, or other evidence may
   influence selection. Reflection does not own or uniquely drive the engine
   architecture.

The engines must compose through shared runtime semantics:

- Interpreter, IR execution, and AOT obey the same declared Java semantic
  contract.
- A Java object retains its identity, class state, and observable contents
  across engine boundaries; fallback must not create a second object world.
- Calls, returns, exceptions, class initialization, and runtime state crossing
  AOT/interpreter/helper boundaries use explicit, testable protocols.
- A fallback is an advertised capability state, not an accidental recovery
  from incorrect AOT output. Unsupported behavior must not be silently
  approximated.
- Capability claims report each relevant engine as `Supported`, `Fallback`, or
  `Not implemented`. Success in one engine does not establish support in the
  others.

Specific IR forms, AOT eligibility rules, transition ABIs, deoptimization,
reflection implementation, and JIT promotion require workstream evidence and,
where they create durable cross-package constraints, separate ADRs.

## Consequences

### Positive

- The product architecture matches catty's main performance thesis: business
  code can become native Go without deleting the runtime path needed for Java's
  dynamic surface.
- AOT coverage can grow incrementally. Missing lowering can use an explicit
  interpreter fallback instead of forcing incorrect compilation.
- The interpreter remains useful for semantic comparison, diagnosis, and code
  that was not available at AOT build time.
- IR investment can be shared by AOT, validation, optimization, and a possible
  future JIT without committing the project to a fixed tier sequence.
- A unified runtime and object world make mixed-engine programs an intentional
  product mode rather than an ad hoc bridge.

### Negative

- A permanent interpreter and retained runtime metadata increase binary size
  and maintenance cost relative to a pure closed-world AOT system.
- Cross-engine transitions can dominate performance if AOT coverage is shallow
  or generated code calls runtime helpers too frequently.
- Correct mixed-engine exceptions, initialization, object identity, and dynamic
  dispatch require explicit protocols and broad differential testing.
- Go-native output is not by itself proof of a performance gain. Workstreams
  must measure coverage and boundary costs rather than infer them from the
  selected engine.

## Evidence policy

Performance work should report, where relevant:

- AOT coverage of methods or exercised paths;
- the amount of JVM frame, slot, and dispatch simulation eliminated;
- cross-engine transition count and cost;
- generated code's runtime-helper density;
- semantic parity across the applicable engines and the pinned Java reference;
- build time, generated source size, binary size, and execution measurements.

These are evaluation dimensions, not fixed release gates. Each workstream
selects exact gates appropriate to its scope.

## Supersession

This ADR supersedes ADR-0007 in full. It retains ADR-0007's durable conclusion
that the interpreter must remain available, but replaces the reflection-driven
tier model and its concrete dynamic-feature mechanisms. Reflection,
`invokedynamic`, dynamic proxies, runtime class loading, dispatch registries,
and deoptimization remain undecided capabilities until governed by accepted
ADRs and workstream contracts.
