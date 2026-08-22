# 测试质量要求

- 地位：DoD 的"Tests pass"以本文为准（协议 §17）。

## 测试层级

| 层 | 内容 | 位置 | 引入时机 |
|---|---|---|---|
| L0 单元/golden | 解析器、IR、发射器的单元与 golden 测试 | 主仓 `internal/**/*_test.go` | M0 起 |
| L1 内核符合性 | jtreg 子集、Harmony 测试套件驱动 java.* 行为 | 独立 GPL 外部仓（catty-conformance），结果以 CI 状态呈现 | M1 起 |
| L2 应用级 fixture | Jackson 2.x 反序列化、Netty echo 等 | 主仓 tests/fixtures | M1–M2 |
| L3 工作负载 | DaCapo / SPECjvm2008 | 外部，基准报告引用 | M2+ |

## 规则

1. 失败的测试不可绕过（硬规则，AGENTS.md）;
2. flaky 政策：标记 quarantine + 开 debt 项，48h 修复窗口，超期升级 Orchestrator;
3. 覆盖率目标：内核包 ≥80%（代码出现后由 Orchestrator 决定阈值写入 make check 的时机）;
4. 发射器必须用"字节码输入 → Go 源码输出"golden 文件测试，禁止快照式黑盒;
5. 符合性测试发现的行为分歧 → 先登记 deviation-ledger 或修引擎，不得改测试期望迁就实现。
