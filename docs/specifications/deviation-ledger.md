# 语义偏差账本（Deviation Ledger）

- 目的：任何与 JVMS / JDK 8 参考行为不一致之处必须在此登记（architecture-rules R8）。这是项目的公开信誉资产。
- 流程：发现偏差 → 登记本表 → Orchestrator 决定"修引擎"或"接受并记录"；接受项需注明影响面与理由。
- 约定：`—` 表示暂无内容。

| ID | 规范条款 | Catty 行为 | 差异原因 | 影响面 | 追踪 |
|---|---|---|---|---|---|
| DEV-0001 | ~~JVMS §4.10 数据流验证缺失~~ 已实现：结构层 + 基于SM帧的数据流类型检查（类别/栈深/合并规范化/异常处理帧/checkcast收窄）；残留保守点=未知引用类对放行 | 类型检查器经专项轮实施（DEBT-0012 闭环） | 未知类对的放行仍存在 | 极小 | 后续收紧 resolver 覆盖面 |
| DEV-0002 | JVMS §5.5 跨线程循环 <clinit> 死锁破除 | 未实现；跨线程循环初始化会阻塞 | 单线程时代语义残留 | 多线程类初始化 | P-0005+ |
| DEV-0003 | ~~JVMS §2.2.8 监视器可重入 + wait/notify~~ | 已关闭：M1 Monitor 重写落地（可重入计数/全深度释放/notify/interrupt） | — | — | 24fc480 |
| DEV-0004 | JVM 栈深度 / StackOverflowError | 已实现：解释帧深度计量（默认 4096 帧）超限抛 SOE；数值与 HotSpot 字节级深度不可比 | 帧预算模型差异 | 递归深度语义 | DEBT-0007 已关闭 |
| DEV-0005 | JVMS §5.4.3.3 invokespecial 精确接收者规则 | 按引用类沿超类链解析（简化） | M0 简化 | 极端 super 调用场景 | review checklist |
| DEV-0006 | Thread.setDaemon / 守护线程退出规则 | 未实现；所有线程等价非守护 | 范围未至 | 进程退出语义 | P-0005 |
| DEV-0007 | JVMS §8.3.1.4 volatile 字段 | 字段访问忽略 volatile 标志（普通读写）；解释器下 Go 内存模型恰好覆盖多数场景，但语义上未兑现 | 发射器未上线，解释路径未接 atomic | 可见性敏感代码 | P-0005 评估：解释器接 atomic.Load/Store 或留待 AOT 发射 |
| DEV-0008 | R4 阻塞操作可中断（socket 读） | Socket InputStream.read 阻塞不可被 Thread.interrupt 唤醒（未接 SetDeadline 唤醒通路） | M2 最小网络面范围裁剪 | 中断期间挂死读 | DEBT-0011 |