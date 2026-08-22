# Catty — Java 语言运行时（Go 宿主 · AOT-first）

用 Go 实现的 Java 语言运行时：对象、线程、GC 寄生在 Go 运行时之上；执行引擎以"字节码 AOT 翻译为 Go 源码并整体编译"为主路径，内置字节码解释器兜底动态特性。

## Read order（最小充分上下文，按需加载）

1. `docs/vision.md` — 产品意图与卖点排序。产品决策前必读。
2. `docs/architecture.md` — 分层架构与关键映射。
3. `docs/decisions/` — 架构决策（ADR）。做设计前先查是否已有决策。
4. `docs/plans/active/` — 当前执行计划。
5. 相关源码。

## Governance

本仓库受 [docs/protocol/autonomous-development-protocol.md](docs/protocol/autonomous-development-protocol.md)（v0.1）约束：

- 默认自治执行；升级条件见协议 §19 与项目本地化数字见 `docs/quality/escalation-policy.md`。
- 中大型设计必须走 Research → 多方案 Design → Critic → ADR（协议 §7）。
- Agent 之间通过 artifact 交流（research report / ADR / plan / task packet / commit / review），不依赖聊天记忆（协议 §31）。
- 性能主张必须有 benchmark 证据（协议 §26）；噪声范围内的变化不得宣称优化成功。

## Mandatory validation

宣告任何任务完成前必须通过：

```sh
make check
```

## Rules（硬规则）

- 永远不绕过失败的测试。
- 不手工修改 generated 文件。
- 架构变更需要 ADR；大变更需要 `docs/plans/active/` 下的执行计划。
- 主仓库禁止引入 GPL 代码；一切第三方来源必须登记 provenance（见 `docs/quality/architecture-rules.md`）。
- cgo 只允许存在于 JNI 桥的构建标签内。
- 所有阻塞操作必须经过可中断封装层（线程中断语义）。

## Where things are

| 内容 | 位置 |
|---|---|
| 治理协议 | `docs/protocol/autonomous-development-protocol.md` |
| 产品意图 | `docs/vision.md` |
| 能力状态 | `docs/status.md`（里程碑收口时刷新） |
| 架构基线 | `docs/architecture.md` |
| 决策记录 | `docs/decisions/ADR-xxxx.md` |
| 执行计划 | `docs/plans/{active,blocked,completed}/` |
| 质量要求 | `docs/quality/`（含 escalation-policy.md） |
| 技术债 | `docs/debt/register.md` |
| 研究报告 | `docs/research/`（模板：TEMPLATE.md） |
| 项目简报 | `docs/briefs/` |
| 规格文档 | `docs/specifications/` |
| 机械校验 | `Makefile` + `tools/check_docs.sh` |
