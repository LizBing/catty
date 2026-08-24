# Catty 能力状态

> 持续更新的能力矩阵。每次里程碑收口时刷新；与 `deviation-ledger.md`（行为偏差）、
> `debt/register.md`(已知债务) 互为对照。最后更新：2026-08-24（P-0009 分发层优化轮）。

## 一句话现状

**M3 收口 + M4 首轮完成**：javac 字节码翻译为 Go 源码整体编译，全部 fixture
与真实三方库 minimal-json 在纯 AOT 下输出与参考 JVM **逐字节一致**；
冷启动快 HotSpot 92%（R-0004）；分发链优化后调用密集吞吐从持平跃至
**2.4–3.0× 于解释器**、mapops 1.50×（R-0006）；goroutine 并发承载不随线程数
劣化；pprof/-race 直接观测"Java"程序（卖点⑤双兑现）。解释器内核持续兜底。

体量：生产 Go ≈33k 行（含生成代码），测试 50+ 例，`-race` 干净。

## 能力矩阵

### ✅ 可用

| 域 | 内容 |
|---|---|
| 执行引擎 | v52 指令集：解释执行 + **AOT 发射双路径**（invokedynamic 构建期脱糖；jsr/ret 非法） |
| AOT 发射器 | genemit→gen.go 整体编译（48 类）；installTable+懒加载钩子混合执行零成本回退；异常通道=旗标返回（ADR-0009）；统一表示（ADR-0010） |
| 分发性能 | 单态内联缓存（烘焙槽位）+免锁帧计量+EmitBody 直调：vcall 2.95×、mapops 1.50× 于解释器（R-0006） |
| 栈深正确性 | `classfile.StackEffect` 全 256 opcode 单一事实源 + 总分类哨兵测试；CFG 工作表深度传播 |
| 堆栈回填 | Java 层 Throwable 堆栈：InvokeAs 统一帧追踪、构造点捕获、<init> 链裁剪、叶帧在前渲染（行号 Unknown Source 待 U3） |
| 异常 | Throwable 家族、异常表分发、隐式抛出（NPE/越界/除零/负长/强转）、SOE（双路径计量）、uncaught 报告 |
| 对象模型 | 身份语义、继承链 embedding、接口分派、字段默认值、数组、UTF-16 String |
| 并发 | Thread 全生命周期 + 中断三路径（含 socket 读 SetDeadline 唤醒）；可重入监视器；wait/notify |
| 类元对象 | ldc <class>、getClass()、静态 synchronized |
| 网络 | ServerSocket/Socket 最小映射，纯 Java HTTP echo 可运行 |
| 类加载 | 目录 classpath（多段）、懒式依赖、结构层+数据流验证器默认开启 |
| CLI | `catty [-cp dir] run <Main>`（args 透传）/ `cattaot`（AOT 优先）/ `genemit` |
| 类库面 | 集合(ArrayList/HashMap/HashSet)、StringBuilder、String 宽方法、包装类、Math/fill/parse-radix、Reader/Writer、regex split |

### ❌ 缺失（按解锁价值排序）

| 缺口 | 债务/计划 | 解锁什么 |
|---|---|---|
| SB 五连折叠 + 装箱分配削减 | P-0009 U1（R-0006 归因：mapops 新首位瓶颈） | mapops 吞吐第二刀 |
| p99 持续负载采样 | P-0009 U4 | 卖点②正式证据 |
| 反射 API / 注解 / MethodHandle | M2+ 远期 | Jackson/Spring 级生态 |
| JAR 加载 | DEBT-0008 | 部署形态 |
| 堆栈行号（逐帧调用点 pc） | P-0009 U3 | 诊断体验对齐 JVM |
| JNI | ADR-0007 范围内未开工 | 原生库互操作 |
| Reference 四件套 / -Xmx 映射 | M3+ | GC 语义完整 |
| provenance 扫描 | DEBT-0002 deferred | IP 纪律 |

## 行为偏差

见 [deviation-ledger.md](specifications/deviation-ledger.md)。活跃项：
volatile 标志忽略（DEV-0007）、守护线程规则未实现（DEV-0006）、
SOE 数值模型不同（DEV-0004）、invokespecial 简化（DEV-0005）、
未知引用类对验证放行（DEV-0001 残留）。
已关闭：类初始化顺序(0002)、socket 读中断唤醒(0008)、数据流验证(0001 主项)、监视器重入(0003)。

## 里程碑位置

```
M0 ████████ 完成（解释器 + HelloWorld）
M1 ██████████ 完成（Monitor/加载器/验证器结构层）
M2 ██████████ 完成（线程/SOE/Class 元对象/net+echo/数据流验证器）
M3 ██████████ 完成（AOT 正确性闭环 + R-0004/R-0005 对照表 + minimal-json assault）
M4 ████░░░░░░ 进行中（分发层优化首轮达标 R-0006；余 SB 折叠/装箱/p99 采样）
```

## 质量纪律快照

- 机械门：`make check`（文档结构 + go vet/build/test）
- 并发门：全仓 `-race`
- 符合性：fixture 输出与参考 JVM 对照（oracle = 本机 JDK 25 `--release 8`）；
  JsonDriver/TwoParse/TraceProbe 双路径钉扎于 internal/vm
- 信誉资产：偏差账本公开；债务登记活跃 5 项（4 项 deferred）；ADR 10 份
