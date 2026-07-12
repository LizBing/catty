# Claude 会话恢复档案：catty 从实验性 JVM 到 JRE 平台

> 本文整理自 2026-07-11 至 2026-07-12 的 Claude Code 长会话。原聊天曾被意外
> clear，后从会话记录中恢复。本文不是逐字转录，而是经过核对的工程纪要：保留
> 需求演变、关键推理、实验、错误修正、提交和交接状态。

原始来源：

- Session 1（clear 前）：`26000167-e31b-4725-8d75-b7db3514fbd0.jsonl`
- Session 2（clear 后，R1 Block A–D）：`a8db9be8-552f-4504-a911-ae010bc59c80.jsonl`
- 净化后的可见对话：[CLAUDE_DIALOGUE_CLEAN.md](./CLAUDE_DIALOGUE_CLEAN.md)

## 1. 会话起点

最初目标是：

> 使用 Go 实现一个实验性 JVM，直接套用 Go runtime，并在此前提下追求最优性能；
> 从第一性原则估算最快可交付 MVP 的时间。

初始环境：

- 项目目录为空，greenfield。
- Go 1.26.5，macOS arm64。
- Temurin JDK 25，可用 `java`、`javac`、`jimage`。
- 用户选择的范围：解释器起步、随后 AOT；MVP 只覆盖最小 Java 子集；单线程优先。

第一性原则结论：借用 Go GC、分配器和调度器后，MVP 不需要自行实现 JVM 最重的
GC、线程调度和 JIT；最小解释器 MVP 的理论下界约为一周。实际在首轮会话中完成了
可运行垂直切片。

## 2. 项目愿景的演变

愿景在会话中经历了三次扩展：

1. **实验性 JVM**：能解析 `.class`、加载类、解释字节码。
2. **AOT Java 执行引擎**：bytecode → 类型化 IR → Go 源码 → `go build`。
3. **实验性 JRE 平台**：Java 程序不只是“运行在用 Go 写的 JVM 上”，而是尽可能
   编译为 Go 程序，让 Java 的线程、同步、原子操作、GC 和 I/O 映射到 Go runtime。

最终形成的核心表达是：

> catty 不是简单地“用 Go 重写 HotSpot”，而是让 Java 语义尽可能溶解到 Go
> runtime 中，同时保留可观察的 JRE 语义。

## 3. Phase 1：解释器 MVP

首个根提交：`12cf2f2 Initial commit: catty JVM interpreter MVP (Phase 1)`。

### 3.1 建立的执行管线

```text
.class
  → classpath
  → classfile parser
  → classloader
  → rtda
  → interpreter
  → native core classes
```

首批包：

- `classfile`：ClassReader、常量池、Code 属性、modified UTF-8。
- `classpath`：目录和 jar/zip entry。
- `classloader`：加载、链接、缓存、字段槽位布局。
- `rtda`：Slot、Frame、Thread、Class、Method、Field、Object、Array。
- `interpreter`：密集 `switch` 分派。
- `native`：Object、String、StringBuilder、System、PrintStream。
- `cmd/jvm`：启动入口。

### 3.2 早期关键修复

- 常量池 tag 是 `u1`，最初误读成 `u2`，导致解析错位。
- `CONSTANT_Utf8.length` 是 `u2`，并且内容是 modified UTF-8。
- `putfield` 的弹栈顺序应为先 value、后 objectref。
- `<clinit>` 必须在 main frame 上方执行，不能在没有调用者 frame 时直接调用。
- category-2 值（long/double）占两个 JVM slot。

### 3.3 MVP 验证

首批端到端夹具：

- HelloWorld
- Fibonacci
- Factorial
- ArraySum
- OOPDemo
- StaticFields
- SwitchDemo

所有输出与真实 `java` 逐字节一致。后来加入 EmptyMain 作为启动和空 main 冒烟测试。

### 3.4 首次性能基线

早期 `fib(35)` 数据：

| 引擎 | 约耗时 |
|---|---:|
| catty tree-walker | 4.34 s |
| HotSpot `-Xint` | 0.61 s |
| HotSpot JIT | 0.05 s |

结论：解释器还有约 7 倍工程优化空间，但与 JIT 的数量级差距必须由 AOT 解决。

## 4. 工程化：文档、Git、CI、License

早期即补齐了长期维护基础：

- README、ARCHITECTURE、DEVELOPMENT、ROADMAP、CHANGELOG。
- ADR 体系。
- Git 仓库，`main` 分支。
- 作者信息：`LizBing <lizbing07734@icloud.com>`。
- GitHub Actions：build、vet、test、端到端与 Java 对拍。
- Apache License 2.0。
- 远端：`git@github.com:LizBing/catty.git`。

相关提交：

- `47f2e16`：Apache 2.0。
- `05de869`：EmptyMain。
- `1add842`：修复 `catty build -run` 用裸文件名导致 `$PATH` 查找失败。
- `59340e4`：移除误提交的 `catty_aot` 二进制并更新 `.gitignore`。

## 5. A0：字节码 IR 与消栈

提交：`c7c7556`。

### 5.1 核心洞察

A0 不需要立即构建完整 SSA。对可验证字节码，在控制流合并点，操作数栈深度一致；
可以先用“按栈槽位置编号的虚拟寄存器”表达 Uses/Defs。

管线：

```text
bytecode
  → operand predecode
  → CFG depth dataflow
  → Uses/Defs slot assignment
  → executable IR
```

### 5.2 实验结果

`LoopIR` 正确，但比 tree-walker 慢约 6%。预解码节省的操作数解析成本，被 IR 结构
解引用、栈顶复位和访问器调用抵消。

这形成 ADR-0006：**预解码不是 Go 解释器的性能答案；IR 的价值是验证和 AOT 输入。**

## 6. A1–A1.5：Go 源码发射与类型追踪

### 6.1 A1：首个 AOT 证明

提交：`c168530`。

`Fibonacci.fib` 被降阶后发射为 Go 函数，再由 `go build` 编译。`fib(35)` 约 44 ms，
与 HotSpot JIT 接近，相比解释器约快 100 倍。

这证明了项目的核心性能论点：

```text
Java bytecode → Go source → Go compiler → native code
```

### 6.2 A1.5：类型数据流

提交：`ce24313`。

新增：

- StackMapTable 完整解析与 frame 重构。
- SlotType：Int、Long、Float、Double、Ref、Top。
- IRInst.InTypes。

StackMapTable 在分支合并点提供精确类型，避免一开始就实现完整 verifier 类型格。

## 7. A2：铺宽发射器

### 7.1 A2.1：fresh-per-def

提交：`fcc0666`。

位置稳定的 `s0/s1` 无法处理同一栈槽先装引用、后装整数的情况。发射器改为每个 def
生成新临时变量，use 指向其 reaching def。由此支持 ref 和数组。

### 7.2 A2.2：invoke bridge

提交：`2616128`。

新增 `runtime` 桥接：

- Bootstrap
- GetStatic
- InvokeVirtual
- NewString

HelloWorld 首次整体转译并原生执行，`println` 从 AOT 世界回调 catty native 世界。

### 7.3 A2.3：循环

提交：`aa613cb`。

关键发现：普通 Java 循环的状态通常在 locals 中，循环头操作数栈为空，因此空栈合并
不需要 phi。mutable Go locals 自然承载循环状态。

### 7.4 A2.2b：OOP 与解释态目标桥接

提交：`c6d2009`。

支持：

- new / dup
- invokespecial 构造器
- getfield / putfield
- AOT → interpreted method bridge
- bridge return slot

OOP 夹具成功执行 `new Box(); b.v=21; b.doubled()`。

### 7.5 A2.4：diamond 与 phi copy insertion

提交：`368979d`。

对非空栈合并点分配 merge temp，并在每条前驱边插入 copy；join 后统一读取 merge
temp。三目表达式 `a > b ? a : b` 得到正确结果。

### 7.6 A2.5：所有原始类型

提交：`c3a7b5b`。

支持 long、float、double。category-2 的两个 JVM slot 绑定到一个 Go temp
（int64/float64）。验证了 long 阶乘和浮点运算。

### 7.7 A2.6：边缘项

提交：`0843556`。

- frem/drem → `math.Mod` 桥接。
- category-2 merge phi。
- tableswitch/lookupswitch → Go switch + goto。

## 8. A4：离线 AOT 集成

提交：`8023e90`。

新增 `catty build`：

1. 从 MainClass 做常量池可达性遍历。
2. Pass 1 判断可发射方法。
3. Pass 2 生成直接调用或 runtime bridge。
4. 组装 Go main。
5. 调用 `go build` 生成原生二进制。

HelloWorld 和 Fibonacci 作为完整程序完成 build、run 和 Java 输出对拍。

## 9. 启动性能实验

在缓存预热后，21 次测量得到：

| 场景 | 最小值 | 中位数 |
|---|---:|---:|
| 空进程 `true` | 1.8 ms | 2.0 ms |
| catty EmptyMain | 2.8 ms | 3.0 ms |
| catty HelloWorld | 2.9 ms | 3.6 ms |
| catty HelloWorld `-ir` | 3.2 ms | 3.7 ms |
| Java EmptyMain | 32.0 ms | 34.4 ms |
| Java HelloWorld | 32.9 ms | 34.6 ms |

结论：catty warm startup 约 3 ms，接近原生进程成本，约比 HotSpot 启动快 10 倍。
首次执行曾出现约 580 ms，归因于 OS 冷页缓存和 macOS 首次执行检查，不代表 VM
稳定启动成本。

## 10. 动态语义与 JRE 愿景

关于反射、动态类加载、invokedynamic 和代理，形成的架构结论是：

- 动态功能依赖保留的 Class/Method/Field 元数据。
- AOT 负责 closed-world 热路径。
- 解释器负责运行时新加载、反射和未发射代码。
- 两个世界通过 bridge 转换。
- 不应删除解释器，即使生产路径以 AOT 为主。

相关 ADR：

- ADR-0007：反射与动态功能采用分层模型。
- ADR-0008：AOT-first，解释器是开发与动态 fallback tier。
- ADR-0009：混合类库。
- ADR-0010：Thread = goroutine。
- ADR-0011：Go memory model 映射策略。
- ADR-0012：依靠 Go escape analysis 和 GC。
- ADR-0013：直接映射 Go runtime，不复制传统 JVM 子系统。

用户进一步提出目标：把 catty 扩展为实验性 JRE，在保持 JRE 可观察语义的同时，
大胆丢弃传统 JVM 内部实现包袱。

## 11. R1：运行更真实的 Java 程序

### 11.1 R1.0：异常

提交：`ece4412`。

完成：

- pendingException + throwPC 信号。
- athrow。
- 异常表搜索与跨 frame 展开。
- try/catch/finally。
- NPE、ArithmeticException、ClassCastException、AIOOBE。
- Throwable/Exception/RuntimeException 等合成类。
- lowering 异常 CFG 边，handler entry stack depth = 1。

ExceptionTest 在 java、tree-walker、IR executor 三个引擎上一致。

### 11.2 R1.1：invokeinterface、wide、multianewarray

提交：`e043fc5`。

完成三个重要操作码，并加入 Comparable + 二维数组测试。此时 IR executor 曾存在
通用 `aload/astore` 索引访问错误，聊天中先临时容忍；该错误后来在 Block D 中被
真正修复。

### 11.3 R1.2：java.base bootstrap 基础

聊天末期提交：`a1b991b`。

用户明确选择：**不在 catty 中实现 jimage 解析，使用 JDK 自带 `jimage extract`，
catty 只消费解出的普通 `.class` 目录。**

R1.2 初期完成：

- 未实现 native 方法获得按返回类型的安全 stub，而不是 nil panic。
- RegisterNative/GetNative。
- System.arraycopy、时间、Object/Class/Thread 和位转换相关 natives。
- 加载真实类后替换 native stub。
- ldc Class literal 生成 Class 对象。
- NoSuchMethodError 代替 Go panic。

聊天结束时正在继续补 java.base natives，并调试 Math.max / Arrays.sort 的初始化链。

## 12. 重要踩坑与工程教训

### 12.1 分派方式

Go 的密集 switch 通常会编译为跳转表。改成 `[256]func` 会增加函数调用开销，预期
不会更快。ADR-0006 的 IR 实验也显示“更结构化的分派”并没有带来收益。

### 12.2 `<clinit>` 是高风险区域

会话中多次暴露初始化顺序和 frame 栈问题：

- main frame 与 clinit frame 的相对顺序。
- clinit 同步执行时不能继续消费调用者 frame。
- 在 invokestatic 前触发 class initialization 时，参数必须仍从原调用者 frame 复制。

后续仓库提交 `0a96373` 最终修复了 clinit stack corruption 和 native access flags。

### 12.3 use-before-def

fresh-per-def 发射中，如果 def 与 use 复用同一 JVM slot，必须先读取 use temp，再
更新 slotTemp。该问题曾出现在整数算术和 getfield。

### 12.4 JVM slot 与 Go 参数不是同一个计数体系

long/double 占两个 JVM slot，但对应一个 Go 参数。invokestatic 和 local mapping
必须分别维护逻辑参数索引与 JVM local slot 索引。

### 12.5 真实 JDK 是级联依赖，不是“补一个 native”

Math.max 本身很简单，但 Math `<clinit>` 会触发 Class、Float、Double 等 bootstrap
链。java.base 兼容必须建立清晰的 bootstrap 边界、provider chain 和 smoke test，
不能只依靠遇到一个错误补一个方法。

## 13. Clear 后第二份 Claude session：R1 收尾与 Block A–D

这些工作最初只能从 Git 和进度快照推断；找到第二份 JSONL 后，已确认它们全部发生
在 clear 后的真实 Claude 会话中。

第二份 session 的恢复过程本身值得记录：

1. 新模型启动时只看见 4 个 dirty 文件，没有前序上下文。
2. 用户决定放弃这 4 个未提交修改，回退到 `a1b991b`，从干净 checkpoint 重建认知。
3. Claude 通读公开文档、13 个 ADR、核心代码和测试，并提出三个具体问题。
4. 用户明确：先同步文档；R1.2 剩余工作是批量补 native；ADR 暂缓。
5. 后续逐步完成 R1.2、架构清理、java.base smoke、String 补齐和 Block D 加固。

| 提交 | 后续成果 |
|---|---|
| `c6cf3f1` | 同步 R1.0/R1.1/R1.2 changelog 与进度 |
| `78d92f8` | 补全 java.base bootstrap native 方法 |
| `0a96373` | 修复 clinit 栈污染和 native access flags |
| `9983688` | 更新 R1.2 修复后的进度快照 |
| `755561b` | classloader Provider 链 + launch 包 |
| `ff6db1d` | java.base 自动检测 + smoke test |
| `b704b75` | String 内容方法 + System.getProperty |
| `78670c4` | R1 最终进度快照 |
| `5720147` | R1 Block D 加固 |

Block D 的最终成果是：CI 覆盖 java.base、修复 IR 通用 `aload/astore` 引用丢失、
清理测试重复计数与容忍 hack、性能复测、公开文档同步，以及 ADR-0014/0015。

因此，clear 前 session 的“R1.2 尚未完成”只是第一次会话的断点；第二份 session
已经把 R1 完成并加固。最新权威状态见
[`.claude/progress.md`](../.claude/progress.md)。

## 14. 当前仓库事实状态

截至当前 `main`：

- HEAD：`5720147 feat: R1 Block D — hardening before R2`。
- 解释器：异常、接口分派、多维数组、java.base bootstrap。
- classloader：Provider chain。
- java.base：自动检测，CI 使用 JDK 25 + `jimage extract`。
- String：合成 String + Extra payload，具备内容方法和 byte[] 解码。
- 测试：纯合成与真实 java.base 两条路径；18/18 smoke/e2e。
- 性能复测：解释器约 4.54 s、AOT 约 50 ms、HotSpot `-Xint` 约 0.58 s、
  HotSpot JIT 约 0.05 s。
- 下一阶段：R2，多线程与最小 Unsafe 层。

## 15. 两份会话提交索引

Session 1（项目创建到 clear 前）：

```text
12cf2f2  Phase 1 interpreter MVP
47f2e16  Apache-2.0
c7c7556  A0 IR + stack elimination
c168530  A1 Go-source emitter
05de869  EmptyMain
6d12687  ADR-0007 dynamic tiering
ce24313  A1.5 type tracking
fcc0666  A2.1 fresh-per-def
2616128  A2.2 invoke bridge
aa613cb  A2.3 loops
c6d2009  A2.2b OOP
368979d  A2.4 diamonds/phi
c3a7b5b  A2.5 all primitives
0843556  A2.6 edge cases
8023e90  A4 offline AOT integration
6a3516a  documentation review
1add842  build -run absolute path fix
59340e4  generated binary cleanup
2ad5091  experimental JRE vision + ADR-0008..0013
ece4412  R1.0 exceptions
e043fc5  R1.1 invokeinterface/wide/multianewarray
a1b991b  R1.2 bootstrap infrastructure (partial at chat endpoint)
```

Session 2（clear 后恢复到 R1 Block D）：

```text
c6cf3f1  sync R1.0/R1.1/R1.2 docs and progress
78d92f8  complete java.base bootstrap native filling
0a96373  fix clinit stack corruption + native access flags
9983688  refresh post-R1.2 progress snapshot
755561b  classloader Provider chain + launch package
ff6db1d  java.base auto-detection + smoke test
b704b75  String content methods + System.getProperty
78670c4  final R1 progress snapshot
5720147  R1 Block D hardening
```

## 16. 如何使用这份档案

- 查代码变化：看 [CHANGELOG](./CHANGELOG.md) 和 Git log。
- 查净化后的完整对话：看 [CLAUDE_DIALOGUE_CLEAN](./CLAUDE_DIALOGUE_CLEAN.md)。
- 查 Codex 独立复盘：看 [CODEX_SESSION_REVIEW](./CODEX_SESSION_REVIEW.md)。
- 查当前状态：看 [progress snapshot](../.claude/progress.md)。
- 查架构事实：看 [ARCHITECTURE](./ARCHITECTURE.md)。
- 查“为什么”：看 [ADR index](./adr/README.md)。
- 查原始会话推理和阶段演化：看本文。

本文的作用是恢复上下文，不取代上述权威文档。

## 17. 原会话最后工作片段（恢复版）

下面是意外 clear 前最后一段可恢复的工作记录。用户原话尽量保留；Claude 的部分
按回复和工具调用压缩整理，不伪装成逐字日志。

### 17.1 R1.2 native 补全请求

**用户：**

> 继续补native方法吧，争取在这一步可以补全。

**Claude 随后执行的工作：**

1. 盘点 Float、Double、Integer、Class、System、Thread 的 native 方法。
2. 尝试补充：
   - `Class.getPrimitiveClass`
   - `Class.forName0`
   - `Class.isInstance` / `isAssignableFrom`
   - `Class.getSuperclass` / `initClassName`
   - annotations 空数组
   - `System.mapLibraryName`
   - `Thread.sleepNanos0`
3. 调试 `Math.max` 和 `Arrays.sort` 使用真实 java.base 时的 bootstrap 链。
4. 定位出多个关键问题：
   - `ldc Class` 原先推送 null，导致 Math `<clinit>` 调
     `Class.desiredAssertionStatus()` 时 NPE。
   - `<clinit>` 执行与调用者 frame 之间存在栈污染风险。
   - native 方法注册与 access flags 必须在真实 JDK 类加载后正确绑定。

这一阶段代码尚未形成稳定 checkpoint，后续仓库中的 `78d92f8` 和 `0a96373`
完成了 native 补全和 clinit/access-flags 修复。

### 17.2 模型切换与交接请求

**用户：**

> 我目前是使用glm接入到claude code中使用的，发现我的plan限额了。现在打算接入
> deepseek v4 pro，你认为他的智能足以承担起目前的工作吗？

**Claude 的建议摘要：**

- DeepSeek 可承担边界清晰的小任务和机械编码。
- 当前项目要求长上下文、跨包编辑和 JVMS 细节，切换前应建立稳定 checkpoint。
- 建议用 plan/progress 文档承载架构上下文，避免模型依赖整段聊天。

**用户：**

> 好的，请你完成跟deepseek的交接工作，你说ok之后我重启claude code并切模型。

这里发生了一次明确的流程失误：Claude 应当停止扩展代码、整理当前 diff、恢复构建
和测试、写交接文档，然后回答 OK；但它继续尝试修 native/bootstrap 代码。随后 GLM
额度耗尽，并发生工具并发错误，导致交接没有干净完成。

### 17.3 切换后的失败尝试

用户多次发送：

> 好的，看来glm没有理解我安排交接工作的请求，现在glm的token plan耗光了，
> 该你上场了deepseek！我相信你。

系统连续返回工具并发错误，提示需要 rewind。原会话至此中断。

### 17.4 对这次交接事故的结论

- **正确做法**：一旦收到明确交接请求，应立即冻结实现范围。
- **必须交付**：当前 HEAD、dirty diff、构建状态、失败测试、最后错误栈、下一步最小
  任务、禁止重复踩坑项。
- **项目后来采取的补救**：`.claude/progress.md`、ROADMAP、ADR、Block A–D
  checkpoint 已经把模型切换依赖从“聊天记忆”转成“仓库内事实”。
- **当前风险**：本地 `main` 比 `origin/main` 超前；切换机器或清理本地仓库前必须
  先确认这些提交已推送。
