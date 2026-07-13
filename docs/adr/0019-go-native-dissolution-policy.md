# ADR-0019: Go-native dissolution and representation policy

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** [ADR-0004](./0004-native-core-classes.md)

## Context

catty is not intended to reproduce every traditional JVM layer inside Go. A
synthetic Go implementation is not automatically compatible, and a classfile
implementation is not automatically simpler or more correct.

## Decision

catty may represent and execute Java abstractions directly through Go-native
structures and runtime mechanisms when that preserves ADR-0017 and reduces total
runtime complexity or cost. The architectural objective is to minimize semantic
boundaries, duplicate representations, and conversion overhead—not to maximize
either synthetic classes or Java classfile implementations.

Four choices are distinct and must not be conflated:

- Go-backed object representation;
- runtime helper;
- intrinsic replacement for an operation or method; and
- synthetic class implementation.

Each durable choice must state its Java-visible contract, representation
boundary, cross-engine behavior, evidence, and replacement conditions. A
synthetic class is justified by a required capability, not merely by local
implementation convenience. Java-visible exception behavior, validation order,
object identity, type behavior, and state transitions remain subject to
ADR-0017.

## Consequences

- A class may combine classfile metadata or methods with Go-backed state and
  targeted intrinsics.
- A Go-native representation is not a documented semantic deviation by itself.
- The R1 synthetic core remains an implementation baseline, not a permanent
  class-library design.
