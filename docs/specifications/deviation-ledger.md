# 语义偏差账本（Deviation Ledger）

- 目的：任何与 JVMS / JDK 8 参考行为不一致之处必须在此登记（architecture-rules R8）。这是项目的公开信誉资产。
- 流程：发现偏差 → 登记本表 → Orchestrator 决定"修引擎"或"接受并记录"；接受项需注明影响面与理由。
- 约定：`—` 表示暂无内容。

| ID | 规范条款 | Catty 行为 | 差异原因 | 影响面 | 追踪 |
|---|---|---|---|---|---|
| DEV-0001 | JVMS §4.10 字节码验证 | 结构层已实现（分支/处理点必须有 SM 帧、池引用合法、边界检查）；操作数类型数据流模拟未做，引用归并对未知类保守放行 | 类型检查器按高风险任务单列（DEBT-0009），避免带缺陷的类型检查器造成虚假安全感 | 恶意 .class 仍可致引擎 panic | DEBT-0005→拆分：结构层已完成，数据流=DEBT-0009 |
| DEV-0002 | JVMS §5.5 跨线程循环 <clinit> 死锁破除 | 未实现；跨线程循环初始化会阻塞 | M0 单线程语义 | 多线程类初始化 | M1 loader 工作 |
| DEV-0003 | JVMS §2.2.8 监视器可重入 + wait/notify | sync.Mutex 不可重入，wait/notify 未实现 | M0 占位实现 | synchronized 重入会死锁 | DEBT-0006 |
| DEV-0004 | ~~无 SOE~~ 已实现：解释帧深度计量（默认 4096 帧），超限抛 StackOverflowError；注意与 HotSpot 的字节级深度不同 | 帧预算模型差异 | 低 | DEBT-0007 已关闭 |
| DEV-0006 | Thread.setDaemon / 守护线程退出规则 | 未实现；所有线程等价非守护 | 范围未至 | 进程退出语义 | P-0005 |
| DEV-0005 | JVMS §5.4.3.3 invokespecial 精确接收者规则 | 按引用类沿超类链解析（简化） | M0 简化 | 极端 super调用场景 | review checklist |
