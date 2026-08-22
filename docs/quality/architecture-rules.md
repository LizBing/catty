# 架构红线（Architecture Rules）

- 地位：可机械强制的规则写成程序，其余进 review checklist（协议 §16）。
- Agent 无权宣布违反以下规则"没关系"。

| # | 规则 | 理由 | 强制方式 |
|---|---|---|---|
| R1 | cgo 只允许存在于 JNI 桥的 `jnion` 构建标签内；纯 Go 构建（默认 tag 集）必须永远可用 | 静态链接与部署故事 | review + 未来 arch test |
| R2 | 主仓禁止引入 GPL 代码；非自写来源必须逐文件 PROVENANCE 登记 | IP 姿态 ADR-0008 | review + DEBT-0002 CI 计划 |
| R3 | IR 与 classfile 解析模型必须版本无关；版本特有逻辑只允许在前端 reader 与独立 pass | JDK 8→21 演进路径 ADR-0002 | review + golden tests |
| R4 | 一切阻塞操作（IO/锁等待/sleep）必须经过可中断封装层 | Thread.interrupt 语义完整性 architecture.md §7 | review checklist |
| R5 | Java 异常必须用哨兵类型承载，禁止与普通 Go panic 混用；机制本身待 ADR-0005 定稿 | 边界语义清晰 | review checklist |
| R6 | generated 文件（发射器产物等）禁手工修改，必须改生成器 | 可重现性 | 文件头 marker + 未来 script 检查 |
| R7 | 内核包不得 import 执行引擎/发射器/类库上层包 | 分层稳定 | DEBT-0003 arch test |
| R8 | 行为与 JVMS/JDK8 规范不符时：修引擎或登记 deviation-ledger，二选一，不允许静默偏差 | 可审计性 | review checklist |

## Review checklist 补充项（不可机械化的部分）

- 并发风险：新增共享状态是否有明确 happens-before 说明；
- 失败模式与回滚路径；
- 观测性：出错时能否从日志/pprof 定位。
