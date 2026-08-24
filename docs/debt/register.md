# 技术债登记簿

每项格式：ID / Description / Impact / Urgency(H/M/L) / Estimated Cost(S/M/L) / Suggested Fix。
Orchestrator 在空闲期择"高收益/低成本"清偿；清偿后移入表格底部归档区并注明解决方式。

| ID | Description | Impact | Urgency | Cost | Suggested Fix |
|---|---|---|---|---|---|
| ~~DEBT-0001~~ ⏸ deferred | classfile 解析器无 fuzz harness；解析器是最大外部输入攻击面 | 解析崩溃/安全 | M | S | M0 落地后加 go-fuzz/native fuzzing 目标进 make check 可选项 |
| ~~DEBT-0002~~ ⏸ deferred | 无自动化许可证/provenance 扫描，目前靠 review 手查 | IP 污染风险 | M | M | 引入 CI 步骤校验 PROVENANCE 完整性与许可证清单白名单 |
| ~~DEBT-0003~~ ⏸ deferred | R7（内核 import 边界）尚无 arch test 强制 | 分层腐化 | L | S | go.mod 划分包后写 import 边界测试进 make check |
| ~~DEBT-0004~~ ⏸ deferred | libcore/Harmony 逐文件许可甄别未开始（ADR-0006 前置） | 阻塞 L2 移植 | M | M | M1 启动前产出甄别清单报告（research artifact） |
| DEBT-0005 | ~~字节码验证器缺失~~ 已完成：结构层(P-0003)+数据流类型检查(P-0005) | — | — | — | 关闭；未知引用类对放行残留于 DEV-0001 |
| DEBT-0006 | ~~监视器不可重入、无 wait/notify/interrupt~~ 已完成（P-0003 Monitor 重写） | — | — | — | 关闭 |
| DEBT-0007 | ~~SOE 深度计量~~ 已完成（P-0004，Options.MaxFrames） | — | — | — | 关闭；发射器侧序言计量待 AOT 后评估 |
| DEBT-0010 | ~~Class 元对象缺失~~ 已完成（P-0005）：ldc-class/getClass/静态 synchronized 解锁 | — | — | — | 关闭；Class 反射 API（字段/方法）仍属远期反射面 |
| ~~DEBT-0011~~ ✅ 已关闭 | socket 读阻塞不可被中断唤醒 | JThread 增加 netConn 停靠槽；natStreamReadB 读前清 stale deadline 并登记 conn；Interrupt 路径 SetDeadline(now) 强制唤醒；超时+中断标志 → InterruptedIOException（新增 java/io/IOException 家族）；中断标志保持置位符合 JDK 语义 | — | — | — | 回归：TestSocketReadInterruptible（真实 loopback，断言 <1s 唤醒 + 异常类型 + 标志保持）|
| DEBT-0012 | ~~数据流验证器三次顺延~~ 已闭环（P-0005 专项） | — | — | — | 关闭；教训记录：该类任务必须整段预算，前三次并行挤压均产出废稿 |
| DEBT-0008 | ~~cmd/catty 未接线 main args 与 classpath~~ 已完成目录 -cp；args/JAR 仍缺 | 可用性 | L | S | JAR 支持随 ClassLoader 抽象（M2） |
| DEBT-0009 | ~~验证器缺数据流类型检查~~ 已完成（P-0005 专项轮） | — | — | — | 关闭；残余保守点（未知引用类对放行）登记于 DEV-0001 |

## 归档

| ~~DEBT-0014~~ ✅ 已解决 | AOT 异常处理器路径槽位错位 | 线性发射器地址序遍历在 try→catch 切换处深度追踪偏移；CollectionsDemo/ThreadsDemo 首个 println 正确但 handler 内 StringBuilder 链的 recv 槽错位 | computeCanonicalDepths 地址序模拟未正确处理 goto→handler 转换后的 SM 帧重同步 | 中 | 下轮：per-basic-block 发射或从 verifier 数据流结果直接导出每 pc 规范深度 |
| ~~DEBT-0015~~ ✅ 已关闭 | AOT 深度模拟与发射器槽位一致性 | 根因链全修：①netStackEffect 缺 aaload 族/pop(0x57)（cat2 重构时误删）②lstore/dstore cat2 应 -2 ③invoke 返回 J/D 应 +2④lstore 需 popRefCat2；WordCount+CollectionsDemo+ThreadsDemo+app.Main+HelloWorld 全部经纯 AOT 输出与 JVM 一致（CD 仅 DEV-0009 消息格式） | 地址序线性模拟本身可行，错在效果表不完整——已补全并用 javap 逐 pc 手工推演验证 | — | — | 回归测试：internal/vm TestAOT* 双路径逐字节钉扎 |
| DEBT-0016 | 验证器 readValue 型 switch 汇合点栈深保真 | minimal-json JsonParser.readValue()：多 case goto 同一 join，某前驱边携带多余 [int]，模拟与 SM 帧冲突（"merge from pc=4: stack depth 1, frame says 0"）；已修死锁(类初始化顺序)+uninitThis 归属后暴露的下一个数据流保真缺口 | 需对 switch/异常边逐前驱核对效果表；或引入 verify 单元级最小复现用例逐指令对照 javap -v | 中 | 进展：switch 克隆序、dup_x1/dup2_x* 语义与 append 别名污染、类初始化顺序(DEV-0002)、uninitThis 归属、String/Writer/Double/Math/List 迭代协议等 9 项已修；当前边界=解释器 char[] 元素类型泄露(uint16 vs int32) panic @JsonParser 缓冲链 |
| DEBT-0017 | invoke 返回 cat2 时发射器深度与 canonDepths 一致性未完全统一 | Bench 用 tickMillis(int) 绕行后全部基准可 AOT；nanoTime(J) 路径仍受限 | 与 0015 同源：需单一事实源 | 中 | 并入 0015 收尾方案 |

| DEBT-0018 | 同一进程内第二次 Json.parse(String) 返回 null | 六个渐进文档探针（含嵌套/浮点/科学计数法）全部经解释器解析正确；随后对同一文档再次 parse 得 null（"Not an object: null"）。JVM oracle 全部正常 | 初判 StringReader EOF 边界或 JsonParser fillBuffer 与 read() 返回约定的残余偏差；需最小二例复现（同 doc 连续两次 parse）逐 read 对照 | 中 | 下轮：新增双 parse 回归用例，instrument StringReader.read 序列 vs JDK |

| DEBT-0019 | AOT 路径 anewarray→aastore 链 NPE | JsonDriver.main docs 数组：生成序列 NewRefArray+AStoreChecked 文本正确，运行时报 array store into null；解释器同文档正常 | NewRefArray 内 K.ArrayClassOf 可能返回 nil → Throw NPE；或 genrt.K 时序 | 高 | 下轮首项：edit 工具落地 [nra] 探针定位 |
