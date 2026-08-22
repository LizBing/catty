# 技术债登记簿

每项格式：ID / Description / Impact / Urgency(H/M/L) / Estimated Cost(S/M/L) / Suggested Fix。
Orchestrator 在空闲期择"高收益/低成本"清偿；清偿后移入表格底部归档区并注明解决方式。

| ID | Description | Impact | Urgency | Cost | Suggested Fix |
|---|---|---|---|---|---|
| DEBT-0001 | classfile 解析器无 fuzz harness；解析器是最大外部输入攻击面 | 解析崩溃/安全 | M | S | M0 落地后加 go-fuzz/native fuzzing 目标进 make check 可选项 |
| DEBT-0002 | 无自动化许可证/provenance 扫描，目前靠 review 手查 | IP 污染风险 | M | M | 引入 CI 步骤校验 PROVENANCE 完整性与许可证清单白名单 |
| DEBT-0003 | R7（内核 import 边界）尚无 arch test 强制 | 分层腐化 | L | S | go.mod 划分包后写 import 边界测试进 make check |
| DEBT-0004 | libcore/Harmony 逐文件许可甄别未开始（ADR-0006 前置） | 阻塞 L2 移植 | M | M | M1 启动前产出甄别清单报告（research artifact） |
| DEBT-0005 | 字节码验证器缺失（StackMapTable 推理验证，JVMS §4.10） | 安全/健壮性 | H | M | M1 实现 StackMapTable 验证；在此之前仅信任本地 javac 输入 |
| DEBT-0006 | 监视器不可重入、无 wait/notify/interrupt | 正确性（DEV-0003） | H | M | M1 重写 kernel Monitor：计数+等待队列+中断通道 |
| DEBT-0007 | 无 StackOverflowError 深度计量 | 语义完整性（DEV-0004） | M | S | M2 发射序言/解释循环加低频深度检查，映射 -Xss |
| DEBT-0008 | cmd/catty 未接线 main args 与 classpath 目录加载 | 可用性 | L | S | M1 引入 ClassLoader 抽象时一并处理 |

## 归档
