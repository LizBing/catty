# Architecture Decision Records

ADRs record *why* a decision was made — the part of a codebase hardest to
reconstruct later. Each is short and immutable: superseding a decision adds a
new ADR that references the old one rather than editing it.

Format (Michael Nygard's template): Status, Context, Decision, Consequences.

Only **Accepted** ADRs govern implementation. Proposed ADRs are discussion
inputs: age, Roadmap placement, or repeated references do not make them
accepted. Reversing an Accepted decision requires a superseding ADR.

| # | Title | Status |
|---|---|---|
| [0001](./0001-reuse-go-runtime.md) | Reuse the Go runtime (no custom GC, scheduler, or JIT) | Accepted |
| [0002](./0002-switch-dispatch.md) | Switch dispatch for the interpreter (Go has no computed goto) | Accepted |
| [0003](./0003-tagged-slot.md) | Tagged 16-byte Slot (HotSpot stack-word model) | Accepted |
| [0004](./0004-native-core-classes.md) | Synthetic native core classes instead of a JRE | Accepted |
| [0005](./0005-lazy-clinit.md) | Lazy `<clinit>` via frame push at JVMS §5.5 points | Accepted |
| [0006](./0006-predecode-no-speedup.md) | Predecode in the interpreter doesn't pay off — AOT is the perf path | Accepted |
| [0007](./0007-reflection-dynamic-tiered.md) | Reflection & dynamic features: tiered — keep the interpreter | Accepted |
| [0008](./0008-aot-first.md) | AOT-first architecture (interpreter is the dev tier, not production) | Proposed |
| [0009](./0009-hybrid-class-library.md) | Hybrid class library (~50 native Go + ~7000 interpreted from real JDK) | Proposed |
| [0010](./0010-thread-equals-goroutine.md) | Thread = goroutine (virtual threads from day one) | Proposed |
| [0011](./0011-go-memory-model.md) | Adopt Go memory model (not JMM/JSR-133) | Proposed |
| [0012](./0012-escape-analysis-replaces-generational-gc.md) | Escape analysis replaces generational GC | Proposed |
| [0013](./0013-direct-go-runtime-integration.md) | Direct Go runtime integration (no JVM abstraction layer) | Proposed |
| [0014](./0014-synthetic-string-extra-payload.md) | Synthetic String with a Go-string Extra() payload | Accepted |
| [0015](./0015-bootstrap-class-boundary.md) | The bootstrap-class boundary (6 irreducible synthetic classes) | Accepted |
