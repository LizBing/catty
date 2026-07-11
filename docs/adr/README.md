# Architecture Decision Records

ADRs record *why* a decision was made — the part of a codebase hardest to
reconstruct later. Each is short and immutable: superseding a decision adds a
new ADR that references the old one rather than editing it.

Format (Michael Nygard's template): Status, Context, Decision, Consequences.

| # | Title | Status |
|---|---|---|
| [0001](./0001-reuse-go-runtime.md) | Reuse the Go runtime (no custom GC, scheduler, or JIT) | Accepted |
| [0002](./0002-switch-dispatch.md) | Switch dispatch for the interpreter (Go has no computed goto) | Accepted |
| [0003](./0003-tagged-slot.md) | Tagged 16-byte Slot (HotSpot stack-word model) | Accepted |
| [0004](./0004-native-core-classes.md) | Synthetic native core classes instead of a JRE | Accepted |
| [0005](./0005-lazy-clinit.md) | Lazy `<clinit>` via frame push at JVMS §5.5 points | Accepted |
