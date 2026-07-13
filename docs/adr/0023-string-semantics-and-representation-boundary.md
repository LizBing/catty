# ADR-0023: String semantics precede String representation

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** [ADR-0014](./0014-synthetic-string-extra-payload.md)

## Context

R1's synthetic String stores a Go string in `Extra()` and is valuable evidence
that a Go-native representation can simplify constant and I/O paths. Its
rune-based `length` and `charAt` behavior, however, is not Java String's
UTF-16-code-unit contract and cannot be a Java 25 compatibility claim.

## Decision

Any supported String capability obeys Java 25 String semantics, including its
UTF-16 code-unit model and the ability to preserve arbitrary code-unit
sequences. Go-native String representation remains an intended optimization
space, but no particular backing format, synthetic-class boundary, intrinsic
set, or classfile split is fixed by this ADR.

R1's existing String implementation is retained only as an explicitly limited
baseline. It does not establish support for supplementary-character,
unpaired-surrogate, reflection-layout, or broad java.base String behavior.
Such behavior must not be claimed supported until a research workstream
provides Java 25 and Temurin 25 evidence.

## Consequences

- Go `string`, `[]uint16`, compact dual encodings, and hybrid representations
  remain open research options.
- String method additions require semantic tests, not only successful Go
  standard-library calls.
