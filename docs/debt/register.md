# 技术债登记簿

每项格式：ID / Description / Impact / Urgency(H/M/L) / Estimated Cost(S/M/L) / Suggested Fix。
Orchestrator 在空闲期择"高收益/低成本"清偿；清偿后移入表格底部归档区并注明解决方式。

| ID | Description | Impact | Urgency | Cost | Suggested Fix |
|---|---|---|---|---|---|
| DEBT-0001 | classfile 解析器无 fuzz harness；解析器是最大外部输入攻击面 | 解析崩溃/安全 | M | S | M0 落地后加 go-fuzz/native fuzzing 目标进 make check 可选项 |
| DEBT-0002 | 无自动化许可证/provenance 扫描，目前靠 review 手查 | IP 污染风险 | M | M | 引入 CI 步骤校验 PROVENANCE 完整性与许可证清单白名单 |
| DEBT-0003 | R7（内核 import 边界）尚无 arch test 强制 | 分层腐化 | L | S | go.mod 划分包后写 import 边界测试进 make check |
| DEBT-0004 | libcore/Harmony 逐文件许可甄别未开始（ADR-0006 前置） | 阻塞 L2 移植 | M | M | M1 启动前产出甄别清单报告（research artifact） |
| DEBT-0005 | ~~字节码验证器缺失~~ 已完成：结构层(P-0003)+数据流类型检查(P-0005) | — | — | — | 关闭；未知引用类对放行残留于 DEV-0001 |
| DEBT-0006 | ~~监视器不可重入、无 wait/notify/interrupt~~ 已完成（P-0003 Monitor 重写） | — | — | — | 关闭 |
| DEBT-0007 | ~~SOE 深度计量~~ 已完成（P-0004，Options.MaxFrames） | — | — | — | 关闭；发射器侧序言计量待 AOT 后评估 |
| DEBT-0010 | ~~Class 元对象缺失~~ 已完成（P-0005）：ldc-class/getClass/静态 synchronized 解锁 | — | — | — | 关闭；Class 反射 API（字段/方法）仍属远期反射面 |
| DEBT-0011 | socket 读阻塞不可被中断唤醒（DEV-0008） | 线程中断完整性 | M | S | SetDeadline-on-interrupt：Interrupt 时对该线程登记的 conn 设置最近 deadline |
| DEBT-0012 | ~~数据流验证器三次顺延~~ 已闭环（P-0005 专项） | — | — | — | 关闭；教训记录：该类任务必须整段预算，前三次并行挤压均产出废稿 |
| DEBT-0008 | ~~cmd/catty 未接线 main args 与 classpath~~ 已完成目录 -cp；args/JAR 仍缺 | 可用性 | L | S | JAR 支持随 ClassLoader 抽象（M2） |
| DEBT-0009 | ~~验证器缺数据流类型检查~~ 已完成（P-0005 专项轮） | — | — | — | 关闭；残余保守点（未知引用类对放行）登记于 DEV-0001 |

## 归档

| ~~DEBT-0014~~ ✅ 已解决 | AOT 异常处理器路径槽位错位 | 线性发射器地址序遍历在 try→catch 切换处深度追踪偏移；CollectionsDemo/ThreadsDemo 首个 println 正确但 handler 内 StringBuilder 链的 recv 槽错位 | computeCanonicalDepths 地址序模拟未正确处理 goto→handler 转换后的 SM 帧重同步 | 中 | 下轮：per-basic-block 发射或从 verifier 数据流结果直接导出每 pc 规范深度 |
