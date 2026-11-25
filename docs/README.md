# Catty developer documentation

The [project README](../README.md) covers installation, usage, and CLI flags.
Run commands in these documents from the repository root unless specified
otherwise.

- [Architecture](ARCHITECTURE.md)
- [Catty 0.1 preview and installation plan](CATTY-PREVIEW.md)
- [Development and validation guide](DEVELOPMENT.md)

| Gate | Specification | Test plan | Evidence |
| --- | --- | --- | --- |
| 1 | [Vertical slice](IMPLEMENTATION.md) | [Acceptance](ACCEPTANCE.md) | [Results](../out/EVIDENCE.md) |
| 2 | [Java 21 compatibility profile](GATE2.md) | [Gate checks](GATE2.md#6-gate-checks) | [Results](../out/gate2/EVIDENCE.md) |
| 3 | [Runtime plugins](GATE3.md) | [Test plan](GATE3-TESTING.md) | [Results](../out/gate3/EVIDENCE.md) |
| 4 | [Threads and monitors](GATE4.md) | [Test plan](GATE4-TESTING.md) | [Results](../out/gate4/EVIDENCE.md) |
| 5 | [Application lifecycle and diagnostics](GATE5.md) | [Test plan](GATE5-TESTING.md) | [Results](../out/gate5/EVIDENCE.md) |
