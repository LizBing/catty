# Codex 独立复盘：clear 前后工作是否顺利衔接

## 1. 评价依据

本复盘基于两份 Claude Code 原始 JSONL、Git 提交、当前代码、CI、测试脚本、公开
文档和 `.claude/progress.md`，不是只根据最终截图或提交信息推断。

评价问题：

1. clear 事故后是否恢复了真实上下文？
2. 是否沿着原架构主线继续，而不是另起炉灶？
3. 是否引入隐含技术债或错误成功信号？
4. R1 Block D 后是否真的“回到正轨”？

## 2. 一句话结论

**不是单纯歪打正着。更准确的说法是：交接开局有损、R1.2 中段有过度乐观，
Block A–D 通过主动架构清理和证据加固把项目真正拉回了正轨。**

clear 本身造成了损失，但也迫使项目第一次系统审计文档、类加载边界、java.base
回归路径和测试口径。最终收益不是偶然落下来的，而是第二份 session 后半段通过
调查、复现、根因修复和 CI 加固获得的。

综合评分：

| 阶段 | 评价 |
|---|---:|
| 初始交接与回退 | 5/10 |
| 上下文重建 | 8/10 |
| R1.2 native/bootstrap 收尾 | 7/10 |
| Block A–C 架构与功能整理 | 8/10 |
| Block D 加固 | 9/10 |
| clear 后整体衔接 | **8/10** |

## 3. 衔接得好的地方

### 3.1 从干净 checkpoint 重建，而不是假装记得

第二份 session 开始时，模型只看见 4 个 dirty 文件。它没有直接猜代码意图，而是
承认缺失上下文。用户随后选择回退到 `a1b991b`。

这个决定牺牲了少量未提交工作，但恢复了一个可验证基线。对于 JVM 这类跨包状态
复杂的项目，干净 checkpoint 比在未知 dirty diff 上继续堆改动更可靠。

### 3.2 全量阅读文档和核心代码后才继续

新模型通读 README、ARCHITECTURE、ROADMAP、CHANGELOG、13 个 ADR、解释器、
classloader、native、rtda、AOT 和测试，然后提出具体问题：

- 先同步文档还是直接继续 R1.2？
- R1.2 是 classpath、delegation 还是 native filling？
- Proposed ADR 是否应当 Accepted？

这些问题说明它成功恢复了项目结构和决策边界，而不是只恢复了“下一个 TODO”。

### 3.3 文档同步作为第一 checkpoint

`c6cf3f1` 先把 R1.0、R1.1、R1.2 和 progress 对齐。模型切换后先更新仓库内记忆，
是正确的长期协作方式。这也直接缓解了第一次交接依赖聊天上下文的问题。

### 3.4 Block A 的 Provider chain 是真正的架构修复

clear 前 classloader 中“合成类优先还是 java.base 优先”主要靠硬编码。Provider
chain 把来源选择变成显式策略，launch 包则把 CLI 与运行时装配分离。

这不是为了让测试通过的局部补丁，而是为后续：

- bootstrap/application loader 分层；
- 模块来源；
- 动态类加载；
- 测试 provider；

建立了可演进边界。

### 3.5 Block D 修的是根因，不是症状

InterfaceTest 最初被描述为 invokeinterface/IR 问题。通过最小复现，Claude 发现：

- 单次 invokeinterface 正常；
- 循环中的 invokevirtual 也错；
- 真正根因是通用 `aload/astore` 对 local ≥4 错用 num 访问器，导致引用丢失。

这是非常好的调试过程：逐步缩小变量，推翻原假设，找到更一般的系统性 Bug。

修复后删除了测试容忍 hack，而不是继续把失败标记为“已知问题”。这是项目回到正轨
最有说服力的证据之一。

### 3.6 CI 开始验证真正重要的 java.base 路径

此前 RealBaseSmoke 只在本地有 `CATTY_BOOT` 时运行，CI 实际没有保护 R1 的核心
成果。Block D 将 JDK 固定到 25，使用 `jimage extract`，再将 java.base 路径纳入
端到端回归。

这比“18/18”本身更重要，因为它把一次成功变成了持续约束。

### 3.7 性能数字经过复测

R1 增加异常、provider 和 native 分发后，重新测量得到：

- tree-walker：约 4.54 s；
- AOT：约 40–60 ms；
- HotSpot `-Xint`：约 0.58 s；
- HotSpot JIT：约 0.05 s。

结果与原文档基本一致，说明 R1 没有悄悄破坏核心性能论点。

## 4. 衔接不理想的地方

### 4.1 回退前没有保存 dirty patch

用户要求回退到最近提交，Claude 直接丢弃 4 个文件的修改。用户明确授权了回退，
所以操作本身不违规；但更稳妥的交接流程应当先：

```text
git diff > handoff.patch
或创建 backup branch / stash
```

然后再清理工作区。那 4 个文件中的思路后来被重新实现，但这是可避免的信息损失。

### 4.2 R1.2 曾过早宣布“native 填充完整”

`78d92f8` 后 Claude 一度称 R1.2 可以收工，但用户追问“是否跑通所有 java.base”后，
压力测试立即暴露新的初始化和 native 链问题。

这说明当时的完成标准是“列出的 native 已补”，而不是“代表性 java.base 程序通过”。
完成定义偏实现清单，缺少行为验收。

### 4.3 默认零值 native stub 会制造静默错误

缺失 native 返回 0/null 能帮助依赖探索，但不满足 Java 语义。它会把：

```text
明确失败
```

变成：

```text
程序继续运行，但结果可能错误
```

这是当前最危险的隐含风险之一。研究模式可保留，但默认 JRE 模式应抛
`UnsatisfiedLinkError`，只允许显式白名单 stub。

### 4.4 “18/18”是筛选后的成功集合

Claude 自己在第二份 session 中承认，HashMap、Double.parseDouble、
Integer.toString 被移出 smoke test，因为它们会进入 DecimalDigits → Unsafe 级联。

因此 18/18 是真实结果，但不是随机或代表性覆盖率。它证明“已选择能力运行正确”，
不能证明“普通 java.base 应用大体可用”。

### 4.5 String 策略带来双表示债务

合成 String 使用 Go string 存在 `Object.Extra()`，真实 JDK String 则有 compact
string 的 byte[]/coder 表示。R1 选择继续使用合成 String，并手写内容方法。

这让启动和常用路径简单，但风险包括：

- 每出现一个真实 JDK String bytecode 方法调用，可能还需 native 补丁；
- UTF-16 code unit 与 Go rune/UTF-8 语义不同；
- 反射查看字段时与 OpenJDK String 布局不一致；
- 序列化、intern、substring 边界可能出现双表示不一致。

ADR-0014 已经诚实记录这一债务，但它仍是 R3 反射前必须解决的架构问题。

### 4.6 Unsafe 不是简单的“R2 前置 stub”

JDK 25 的 Integer、Double、HashMap 进入 Unsafe 级联。Unsafe 不只是几十个 native
方法，它连接：

- field offset；
- CAS；
- volatile/fence；
- 数组基址与步长；
- park/unpark；
- 直接内存；
- JMM observable semantics。

如果只补返回零值的 stub，可能让并发容器静默错误。R2 必须先定义最小语义合同，
而不是以“方法数量”衡量完成度。

## 5. 最深的战略矛盾

用户当前目标是“保持所有 JRE 语义”。但 ADR-0011 明确决定采用 Go memory model，
而非 JMM，并接受约 0.1% 偏差。

两者不能同时成立。

正确的架构表述应当是：

> catty 使用 Go 的 mutex、atomic、goroutine 和 channel 实现 JMM 的可观察保证，
> 而不是以 Go memory model 替代 JMM。

类似地，Thread = goroutine 是实现策略，不代表 Thread 语义自动成立。仍需实现：

- start/join happens-before；
- interrupt 状态和 InterruptedException；
- wait/notify monitor ownership；
- daemon 生命周期；
- ThreadLocal；
- uncaught exception handler；
- park/unpark；
- 类初始化同步。

因此 R2 开始前必须修订 ADR-0011，或明确降低产品目标为“高兼容 Java-like
runtime”。

## 6. 是否“歪打正着”

### 有歪打正着的成分

1. clear 事故迫使项目第一次做完整上下文审计和文档同步。
2. 原以为是 invokeinterface 的 Bug，最终找到通用 `aload/astore` 引用丢失。
3. 原本为了证明 R1 的 smoke test，反而暴露 CI 根本没跑 java.base。
4. 重复 PASS 计数在清理 IR 容忍逻辑时被顺带发现。

这些收益确实带有偶然性。

### 但最终回正不是偶然

Block A–D 的成果来自明确行动：

- provider abstraction；
- 最小复现；
- 根因分析；
- 删除容忍 hack；
- CI 回归门禁；
- 性能复测；
- ADR 与公开文档同步。

如果只是歪打正着，项目会停在“测试暂时绿”；现在它具备了较清晰的边界和持续回归
机制。因此我认为：**clear 后先偏航、随后通过工程纪律真正回到正轨。**

## 7. 当前成果应如何定性

合理表述：

> R1 已完成并加固：catty 能在受控边界内加载真实 java.base，执行一组经过明确
> 筛选的基础程序，并由 CI 同时保护合成类与 java.base 路径。

不宜表述：

> catty 已经兼容 java.base，或已经是完整 JRE。

目前更接近“有真实 java.base 路径的实验性 Java execution platform”。

## 8. R2 前建议的硬门槛

1. **推送本地提交**：当前本地 `main` 比 `origin/main` 超前 9 个提交。
2. **修订 ADR-0011**：目标改为用 Go 原语实现 JMM observable semantics。
3. **native strict mode**：默认缺失 native 抛 `UnsatisfiedLinkError`；探索模式才允许
   zero stub。
4. **Unsafe contract**：先列出真正需要的操作、语义、调用者和测试，不按 native
   方法数量盲补。
5. **并发差分测试**：start/join、volatile publication、CAS、monitor、wait/notify、
   interrupt、park/unpark。
6. **代表性失败集常驻 CI**：HashMap、Integer.toString、Double.parseDouble 即使尚未
   通过，也应作为 expected-failure 或 capability test 留在仓库，避免被遗忘。
7. **会话 checkpoint 协议**：模型切换前保存 HEAD、diff、测试状态、阻塞栈和下一步
   单一任务。

## 9. 最终评价

clear 事故后的工作衔接不是完美的：丢了 dirty diff，早期完成口径偏乐观，也通过
筛选测试弱化了 Unsafe/String 风险。但第二份 session 没有继续堆补丁到失控；它最终
完成了架构清理、根因修复、CI 加固、性能复测和文档收口。

所以我的最终判断是：

> **R1 的代码和工程流程已经回到正轨；JRE 全语义目标尚未回到正轨。**

前者允许项目进入 R2 规划，后者要求先解决 JMM、native strictness 和 Unsafe
语义合同三个问题。
