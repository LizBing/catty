# Architecture Decision Records

ADRs record *why* a decision was made — the part of a codebase hardest to
reconstruct later. Each is short and immutable: superseding a decision adds a
new ADR that references the old one rather than editing it.

Format (Michael Nygard's template): Status, Context, Decision, Consequences.

Only **Accepted** ADRs govern implementation. Proposed ADRs are discussion
inputs: age, Roadmap placement, or repeated references do not make them
accepted. `Withdrawn` closes a proposal that never became authoritative;
`Superseded` replaces a formerly Accepted decision. Reversing an Accepted
decision requires a superseding ADR.

| # | Title | Status |
|---|---|---|
| [0001](./0001-reuse-go-runtime.md) | Reuse the Go runtime (no custom GC, scheduler, or JIT) | Superseded by ADR-0018 |
| [0002](./0002-switch-dispatch.md) | Switch dispatch for the interpreter (Go has no computed goto) | Superseded by ADR-0024 |
| [0003](./0003-tagged-slot.md) | Tagged 16-byte Slot (HotSpot stack-word model) | Superseded by ADR-0020 |
| [0004](./0004-native-core-classes.md) | Synthetic native core classes instead of a JRE | Superseded by ADR-0019 |
| [0005](./0005-lazy-clinit.md) | Lazy `<clinit>` via frame push at JVMS §5.5 points | Superseded by ADR-0021 |
| [0006](./0006-predecode-no-speedup.md) | Predecode in the interpreter doesn't pay off — AOT is the perf path | Superseded by ADR-0024 |
| [0007](./0007-reflection-dynamic-tiered.md) | Reflection & dynamic features: tiered — keep the interpreter | Superseded by ADR-0016 |
| [0008](./0008-aot-first.md) | AOT-first architecture (interpreter is the dev tier, not production) | Withdrawn |
| [0009](./0009-hybrid-class-library.md) | Hybrid class library (~50 native Go + ~7000 interpreted from real JDK) | Withdrawn |
| [0010](./0010-thread-equals-goroutine.md) | Thread = goroutine (virtual threads from day one) | Withdrawn |
| [0011](./0011-go-memory-model.md) | Adopt Go memory model (not JMM/JSR-133) | Withdrawn |
| [0012](./0012-escape-analysis-replaces-generational-gc.md) | Escape analysis replaces generational GC | Withdrawn |
| [0013](./0013-direct-go-runtime-integration.md) | Direct Go runtime integration (no JVM abstraction layer) | Withdrawn |
| [0014](./0014-synthetic-string-extra-payload.md) | Synthetic String with a Go-string Extra() payload | Superseded by ADR-0023 |
| [0015](./0015-bootstrap-class-boundary.md) | The bootstrap-class boundary (6 irreducible synthetic classes) | Superseded by ADR-0022 |
| [0016](./0016-multi-engine-execution.md) | Multi-engine execution with AOT as the primary product path | Accepted |
| [0017](./0017-java-25-semantic-baseline.md) | Java 25 semantic baseline for supported capabilities | Accepted |
| [0018](./0018-go-runtime-infrastructure-boundary.md) | Reuse Go runtime infrastructure without substituting Java semantics | Accepted |
| [0019](./0019-go-native-dissolution-policy.md) | Go-native dissolution and representation policy | Accepted |
| [0020](./0020-representation-domains-and-engine-boundaries.md) | Representation domains and typed engine boundaries | Accepted |
| [0021](./0021-class-initialization-boundary.md) | Class initialization is a shared Java runtime boundary | Accepted |
| [0022](./0022-bootstrap-kernel-and-java-facades.md) | Bootstrap kernel capabilities and Java-visible facades | Accepted |
| [0023](./0023-string-semantics-and-representation-boundary.md) | String semantics precede String representation | Accepted |
| [0024](./0024-interpreter-policy-and-evidence-driven-optimization.md) | Interpreter simplicity and evidence-driven optimization | Accepted |
| [0025](./0025-class-initialization-state-machine.md) | Java 25 class and interface initialization state machine | Accepted |
| [0027](./0027-string-utf16-representation.md) | UTF-16 String kernel backing | Accepted |
