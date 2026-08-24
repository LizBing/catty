# 技术债登记簿

每项格式：ID / Description / Impact / Urgency(H/M/L) / Estimated Cost(S/M/L) / Suggested Fix。
Orchestrator 在空闲期择"高收益/低成本"清偿；清偿后移入表格底部归档区并注明解决方式。

## 活跃

| ID | Description | Impact | Urgency | Cost | Suggested Fix |
|---|---|---|---|---|---|
| ~~DEBT-0001~~ ✅ 缓解落地 | classfile 解析器无 fuzz harness | 解析崩溃/安全 | — | — | `make fuzz`（FuzzParse，种子=全部 fixture；30s 401 万次执行无 panic）。长会话深度模糊化与 CI 接入仍为可选项 |
| DEBT-0002 ⏸ deferred | 无自动化许可证/provenance 扫描，目前靠 review 手查 | IP 污染风险 | M | M | 引入 CI 步骤校验 PROVENANCE 完整性与许可证清单白名单 |
| DEBT-0003 ⏸ deferred | R7（内核 import 边界）尚无 arch test 强制 | 分层腐化 | L | S | go.mod 划分包后写 import 边界测试进 make check |
| DEBT-0004 ⏸ deferred | libcore/Harmony 逐文件许可甄别未开始（ADR-0006 前置） | 阻塞 L2 移植 | M | M | M1 启动前产出甄别清单报告（research artifact） |
| DEBT-0008 | cmd/catty 已接线目录 -cp 与 main args；**JAR 加载仍缺** | 可用性 | L | M | 随 ClassLoader 抽象（M2+）落地 zip 读取与并行加载 |

## 归档

| ID | Description | Resolution |
|---|---|---|
| ~~DEBT-0005~~ ✅ | 字节码验证器缺失 | 结构层(P-0003)+数据流类型检查(P-0005)；未知引用类对放行残留于 DEV-0001 |
| ~~DEBT-0006~~ ✅ | 监视器不可重入、无 wait/notify/interrupt | P-0003 Monitor 重写 |
| ~~DEBT-0007~~ ✅ | SOE 深度计量 | P-0004 Options.MaxFrames；发射器侧经 ThreadRegistry.FrameEnter 同样计量 |
| ~~DEBT-0009~~ ✅ | 验证器缺数据流类型检查 | P-0005 专项轮；残余保守点登记 DEV-0001 |
| ~~DEBT-0010~~ ✅ | Class 元对象缺失 | P-0005：ldc-class/getClass/静态 synchronized 解锁 |
| ~~DEBT-0011~~ ✅ | socket 读阻塞不可被中断唤醒 | SetDeadline-on-interrupt；回归 TestSocketReadInterruptible |
| ~~DEBT-0012~~ ✅ | 数据流验证器三次顺延 | P-0005 专项整段预算闭环；教训：此类任务不得并行挤压 |
| ~~DEBT-0014~~ ✅ | AOT 异常处理器路径槽位错位 | computeCanonicalDepths handler-entry 规范深度恒 1（JVMS §2.6） |
| ~~DEBT-0015~~ ✅ | AOT 深度模拟与发射器槽位一致性（解释器侧效果表缺行） | 补全解释器 netStackEffect 并以 javap 逐 pc 推演验证；全 fixture AOT=JVM 逐字节 |
| ~~DEBT-0016~~ ✅ | 验证器/发射器栈深保真族（switch 合并、dup 族、init 顺序、char[] 宽度泄露等 10 项） | switch 克隆序、dup_x* 语义、类初始化顺序(DEV-0002 关闭)、uninitThis 归属、Writer/String/Math/List 迭代面逐项修复；**残余 char[] 元素宽度泄露根修**：natStringCharAt 返回裸 uint16（解释器 push 归一化掩盖、AOT .(int32) 断言必炸），改回 int32 并以 TestNativeReturnWidths 钉扎边界契约 |
| ~~DEBT-0017~~ ✅ | invoke 返回 cat2 时发射器深度与 canonDepths 一致性未统一 | 根治方式升级为**单一事实源**：classfile.StackEffect 全 256 opcode 总分类（EffectFixed/Descriptor/Illegal）+ 漂移钉扎测试；发射器 netStackEffect 委托之。nanoTime(J) 路径解锁（TestAOTCat2NativeReturn）；Bench 暂留 tickMillis 以保 R-0005 可比性，M4 重评 |
| ~~DEBT-0018~~ ✅ | 同进程第二次 Json.parse(String) 返回 null | Route C 类初始化顺序根修（DEV-0002，super 腿先标记 in-progress）顺带治愈；双路径回归钉扎 TestTwoParseDoubleParseBothPaths（TwoParse fixture 双 parse len=94） |
| ~~DEBT-0019~~ ✅ | minimal-json JsonDriver AOT 端到端 NPE "array store into null" | 两层根修：①诊断基建——Java 堆栈回填（kernel.InvokeAs 统一帧追踪 + Throwable 构造点快照 + fillInStackTrace 式 <init> 裁剪 + FormatUncaught 渲染）；②堆栈直指发射器效果表缺整个数组存储族（0x4f–0x56 落 default 0）→ aastore 后深度漂移 → dup 复制错误槽位。修复并入 classfile.StackEffect。JsonDriver AOT 输出与 JVM oracle 逐字节一致（TestAOTJsonDriverEndToEnd）。附带修复：genrt 方法缓存跨内核泄露（InstallKernel 时失效化） |
