# 技术债登记簿

每项格式：ID / Description / Impact / Urgency(H/M/L) / Estimated Cost(S/M/L) / Suggested Fix。
Orchestrator 在空闲期择"高收益/低成本"清偿；清偿后移入表格底部归档区并注明解决方式。

| ID | Description | Impact | Urgency | Cost | Suggested Fix |
|---|---|---|---|---|---|
| DEBT-0001 | classfile 解析器无 fuzz harness；解析器是最大外部输入攻击面 | 解析崩溃/安全 | M | S | M0 落地后加 go-fuzz/native fuzzing 目标进 make check 可选项 |
| DEBT-0002 | 无自动化许可证/provenance 扫描，目前靠 review 手查 | IP 污染风险 | M | M | 引入 CI 步骤校验 PROVENANCE 完整性与许可证清单白名单 |
| DEBT-0003 | R7（内核 import 边界）尚无 arch test 强制 | 分层腐化 | L | S | go.mod 划分包后写 import 边界测试进 make check |
| DEBT-0004 | libcore/Harmony 逐文件许可甄别未开始（ADR-0006 前置） | 阻塞 L2 移植 | M | M | M1 启动前产出甄别清单报告（research artifact） |
| DEBT-0005 | ~~字节码验证器缺失~~ 已部分完成：结构层验证落地（P-0003），数据流类型检查未做 | 安全/健壮性 | H | M | 拆出 DEBT-0009 跟踪类型层 |
| DEBT-0006 | 监视器不可重入、无 wait/notify/interrupt | 正确性（DEV-0003） | H | M | M1 重写 kernel Monitor：计数+等待队列+中断通道 |
| DEBT-0007 | ~~SOE 深度计量~~ 已完成（P-0004，Options.MaxFrames） | — | — | — | 关闭；发射器侧序言计量待 AOT 后评估 |
| DEBT-0010 | ~~Class 元对象缺失~~ 已完成（P-0005）：ldc-class/getClass/静态 synchronized 解锁 | — | — | — | 关闭；Class 反射 API（字段/方法）仍属远期反射面 |
| DEBT-0011 | socket 读阻塞不可被中断唤醒（DEV-0008） | 线程中断完整性 | M | S | SetDeadline-on-interrupt：Interrupt 时对该线程登记的 conn 设置最近 deadline |
| DEBT-0012 | 数据流类型验证器仍未实施（连续第三次排期顺延） | 安全边界完整性 | H | L | 下一轮整段专注预算，禁止与其它主线并行；输入仍限定可信 javac 产物 |
| DEBT-0008 | ~~cmd/catty 未接线 main args 与 classpath~~ 已完成目录 -cp；args/JAR 仍缺 | 可用性 | L | S | JAR 支持随 ClassLoader 抽象（M2） |
| DEBT-0009 | 验证器缺数据流类型检查（类别不匹配/栈深线性模拟） | 安全边界完整性 | H | L | M2 高优：复用 vm 操作码表做线性模拟，SM 帧为合并规范形；未知引用类对保守放行需同步收敛 |

## 归档
