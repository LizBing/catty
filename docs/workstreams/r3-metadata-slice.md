# R3 dynamic metadata kernel slice

**Status:** Accepted
**Type:** implementation
**Review:** owner
**Profile:** Catty JVMS Core shared kernel; no public profile API
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018 through ADR-0024, ADR-0031,
ADR-0032, and ADR-0034
**Prerequisites:** `r3-reflection-dynamic-research` Done; acceptance anchor
fixed before implementation

## Outcome

Catty validates and immutably retains the bounded classfile metadata needed by
the accepted InvokeDynamic kernel: BootstrapMethods, usable MethodType,
MethodHandle, InvokeDynamic operands, and structural ConstantDynamic entries.
Runtime metadata can reference these values without loading classes or
executing bootstrap code during parsing.

This slice creates no Java-visible reflection, annotation, MethodHandle,
CallSite, or InvokeDynamic capability.

## Scope

- Typed BootstrapMethods structures with validated indexes, tags, lengths, and
  immutable accessors.
- Usable constant-pool accessors for MethodType, MethodHandle reference kinds
  1–9, InvokeDynamic name/descriptor/bootstrap index, and structural
  ConstantDynamic.
- Profile-neutral immutable metadata attachment needed by later runtime
  consumers under ADR-0031.
- Explicit handling of structurally valid unrecognized attributes according to
  JVMS rules; malformed known structures return typed parse failures.
- Parser/unit fixtures and exact AOT reachability classification where current
  generic diagnostics can be narrowed without later linkage work.

## Non-scope

Annotation attributes or element trees, declared-member discovery, Exceptions
or MethodParameters reflection metadata, Java facades, class lookup/definition,
MethodHandle execution, opcode `0xba`, CallSite state, generated classes,
concat, lambda, Proxy, or AOT execution/fallback.

## Semantic constraints

Parsing is total and side-effect-free: no Java loading, initialization, facade
allocation, provider lookup, or bootstrap execution. Parsed metadata is
immutable and does not expose reader buffers or constant-pool internals as a
runtime ABI. Structurally valid unknown attributes are ignored for execution;
known accepted structures cannot be silently discarded.

## Acceptance

| Gate | Required result |
|---|---|
| Dynamic metadata parser | MethodHandle kinds 1–9, MethodType, InvokeDynamic, BootstrapMethods, and structural ConstantDynamic pass positive/negative tests |
| Immutability/attachment | Runtime consumers retain exact declared symbolic metadata without eager resolution or mutable reader aliases |
| Parse failure | Malformed indexes/tags/lengths return typed errors; no bounds panic or partial publication |
| Capability honesty | Fixed 24-row R3 baseline remains reproducible; no Java-visible row is newly claimed Supported |
| Regression | `go vet ./...`, `go test ./...`, `go test -race ./...`, and `bash tests/run.sh` Pass |
| Isolation/governance | Historical evidence unchanged; candidate evidence isolated; `git diff --check` Pass |

## Plan

Accepted; waiting for prerequisites and acceptance anchor. No implementation
is in progress.

## Acceptance record

Accepted by Owner on 2026-07-18. Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, and owner review are
frozen. Implementation remains unauthorized until prerequisites are Done and
this contract is fixed in an acceptance-anchor commit.
