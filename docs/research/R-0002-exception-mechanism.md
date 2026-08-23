# R-0002: 异常传播机制基准 — panic/recover vs 旗标返回，裁决 ADR-0009 的证据

- 日期：2026-08-23
- 关联计划：[P-0006](docs/plans/active/P-0006-m3-aot-emitter.md)（判断二：异常机制由实验裁决）
- 关联决策：[ADR-0005](docs/decisions/ADR-0005.md)（当前"延后"即决策）、ADR-0009（将定稿）
- 原型位置：`docs/research/assets/excbench/`（独立 module `excbench`）
- 原始数据：`docs/research/assets/excbench/results/`

---

# Question

AOT 发射器为每个生成方法选择异常传播机制：方案 P（panic/recover，哨兵值）与方案 F（旗标返回二元组）。需要回答：在"4 层调用深度、d 在循环中按概率抛出"的真实发射器形态下，(1) 正常路径吞吐、(2) 异常路径延迟、(3) 深传播三场景谁更快、快多少；并把性能与发射器实现复杂度、语义风险合在一起给出**单一方推荐**与**临界异常频率**。

# Context

ADR-0005 把"异常机制"延后到有数据再定，唯一硬约束是"方法调用必须经统一调用约定抽象层"使机制可切换。P-0006 判断二要求本基准（R-0002）回写数据后以 ADR-0009 定稿。发射器的 v1 目标是**正确性闭环**（M3 完成定义），不含吞吐优化（M4 范畴）；因此语义风险与"发射正确"的难度权重应不低于性能。

# Candidates

- **方案 P（panic/recover）**：Java `athrow` = `panic(*JException 哨兵)`；每个拥有 handler 表项（try/catch/finally/synchronized）的方法边界安装 `defer`，`recover()` 后按类标匹配——匹配则运行 handler，不匹配则 `panic(e)` 重抛（即 Java 异常表"此处不处理，继续传播"）；非哨兵 panic（Go 运行时 panic）同样重抛以严格区分两类 panic。
- **方案 F（旗标返回）**：所有可能抛异常的方法返回 `(result, *JException)`；`athrow` = `return 0, sentinel`；调用点显式 `if err != nil { return err }`。

两者都沿用 ADR-0005 的哨兵类型 `*JException`（只建模类标，消息/cause/Java 栈帧正交于机制，本基准省略）。正常路径上 P 的成本谱从"无 handler 方法（零 defer，真零成本）"到"有 catch 方法（defer+recover+匹配）"连续分布，F 则统一。为公平，两候选做逐位一致的 `common.Work` 算术、预分配哨兵、链方法全部 `//go:noinline`（见 README 方法学）。

# Evidence

## 方法

原型为独立 module `docs/research/assets/excbench/`，`panicpkg`（P）与 `flagpkg`（F）两子包，`common` 保证算术一致。四层链 `a→b→c→d`，d 做 `Work` 并按场景抛出（无 RNG，避免 RNG 成本掩盖机制）。三场景 + P 成本谱各一层：

1. 正常路径（零异常稳态循环）：P 三个变体——全 catch-handler / 全 finally（纯 defer，无 recover）/ 全无 handler；F 一个变体。
2. 异常路径延迟（高频 throw-catch）：d 每轮抛，c 一层外捕获（1 边界）。
3. 深传播（异常穿越 4 层）：P 两个变体——自由展开（中间层无 handler）与逐层重抛（中间层各有一个不匹配 handler）；F 一个变体。

正确性护栏（`go test`）：重抛/吞没逻辑、跨层误捕获（各层 handler 结果值互异）、非哨兵 panic 原样穿透，全部断言通过。

命令（可复现）：

```sh
cd docs/research/assets/excbench
go build ./... && go vet ./... && go test ./...          # 构建 + 正确性护栏

go test ./panicpkg ./flagpkg -run '^$' -bench . -benchmem -count=10   # 主数据 + 分配观测
# 交错计时（Apple Silicon 热漂移控制）：15 轮交替两包，奇偶轮交换顺序
for round in $(seq 1 15); do …; done                     # 见 results/interleaved.txt
```

benchstat 已安装（`/Users/lizbing/go/bin/benchstat`），用于交叉验证；主数据按计划采用 `-count=10` + 手工中位数。硬件 Apple M1，darwin/arm64，go1.26.5。

## 测量结果（median ns/op，数值越小越快；0 B/op / 0 allocs/op 全表成立）

`-count=10 -benchmem` 中位数，15 轮交错运行独立复验（两轮方向一致、偏差 ≤3%）：

| 场景 | P（panic） | F（旗标） | 方向与幅度 |
|---|---:|---:|---|
| 正常路径·无 handler（P 最优点 / F 常态） | **3.6** | 4.6 | P 快 ~21%（−0.97 ns） |
| 正常路径·finally 纯 defer（P） | **~7**（6.6–7.4） | 4.6 | F 快 ~1.4–1.6× |
| 正常路径·catch handler（P） | **13.2** | 4.6 | F 快 ~2.9× |
| throw-catch（1 边界，高频） | **170** | 2.2 | F 快 ~77× |
| 深传播·自由展开 | **208** | 2.9 | F 快 ~72× |
| 深传播·逐层重抛（P 最坏点） | **1595** | 2.9 | F 快 ~550× |

分配：两方案三场景全部 0 B/op、0 allocs/op——机制本身（含 panic 记账）在本微基准中不产生堆分配；预分配哨兵把"每次 throw 的 Java 异常对象分配"排除在外（该分配两方案一致、正交）。

benchstat 交叉验证（`/Users/lizbing/go/bin/benchstat`，对 `-count=10` 数据做 A/B）：NormalPath F 较 P −65.4%（p=0.000）、ThrowShallow −98.8%（p=0.000）、全基准 geomean −93.4%，方向与上表手工中位数完全一致。

逃逸/内联（`-gcflags=-m` 定性）：链方法 `//go:noinline` 显式隔离了"每边界机制成本"，故上述差异即发射器无法靠内联消掉的每方法边界成本，而非函数调用开销本身。P 的 defer 在现代 Go 中 open-coded，`recover` 不可内联是 P catch 路径贵的结构性原因。

## 附注 A：panic 跨 goroutine 行为（读 runtime 源码，非实验）

源码位置 `runtime/panic.go`、`runtime/stack.go`（go1.26.5）：

1. **panic/recover 严格 goroutine 本地**：`gopanic` 用 `getg()` 把 `_panic` 挂在当前 g 上，`gorecover` 也从 `getg()._panic` 取——跨 goroutine 不可能 recover 另一个 goroutine 的 panic。
2. **未恢复的 panic = 整进程退出**：`gopanic` 耗尽 defers 后走 `fatalpanic` → `exit(2)`/`crash()`。这与 Java 语义关键不同：Java 未捕获异常交给线程的 `UncaughtExceptionHandler`（默认打印、该线程死、JVM 继续），而 Go 直接杀进程。**映射要求：发射器必须在每个 Java 线程（goroutine）顶层包一层 recover，把哨兵转译为 Java uncaught-exception 分发；否则任何漏网异常即宕机。**
3. **recover 只能在 defer 内直接调用**：`gorecover` 要求 gopanic 与 gorecover 之间"恰好一个非 wrapper 帧"（源码注释 + unwinder 逻辑）。发射器的 recover 必须写在 defer 闭包里，不能经 helper 转调。
4. **栈溢出是致命 `throw`，非可恢复 panic**：`runtime/stack.go` 栈增长超限走 `throw("stack overflow")`（"goroutine stack exceeds N-byte limit"），**recover 接不住**。Java 的 `StackOverflowError` 是可捕获异常——这是两方案共有的语义缺口，必须用发射器注入的栈深 guard 补，不能靠 recover。
5. **非哨兵 panic 必须重抛**：nil 解引用/越界/类型断言产生的 `runtime.Error` panic，若被 recover 顺手当 Java 异常吞掉，会掩盖真实 bug。P 的 recover 必须 `else { panic(e) }`（本基准即如此，护栏已验）。
6. 次要：`panic(nil)` 自 Go 1.21 起变 `*PanicNilError`（Java `throw null` 是 NPE，应映射为 NPE 哨兵而非 panic(nil)）；`runtime.Goexit` 期间 recover 返回 nil（Goexit 非 panic）。

## 附注 B：Java 语义映射完备性检查（两方案发射复杂度，定性）

| Java 构造 | P（panic/recover） | F（旗标返回） |
|---|---|---|
| 无 handler 方法（多数） | 零成本，零样板（panic 自由展开） | 签名污染 `(T,err)` + 每调用点检查，无法豁免 |
| try/catch 且 catch 后继续同方法 | **需要函数拆分**：recover 使函数"返回给调用者"，catch 块后置代码必须挪到外层或把 try 块拆成返回 `(T,err)` 的 helper——拆分边界又退化成 F | 原生：`if err != nil { handle; 继续 }` |
| finally | 原生：`defer`（自动覆盖正常+展开所有出口） | 需枚举所有出口点重复 cleanup（漏一个 return 即泄漏），易错 |
| try-with-resources | defer + 嵌套 catch；suppressed 记账手工 | 显式 close + 检查；suppressed 记账显式，但样板更多 |
| synchronized（monitor exit） | 原生：defer 释放监视器 | 需在每个出口显式释放 |
| 异常对象（cause/suppressed/栈帧） | 正交，两方案同载 `*JException`，成本一致 | 同左 |
| 与 Go 运行时 panic 的区分 | **必须**哨兵类型断言 + 重抛（纪律风险） | 无需——错误是普通值，Go panic 保持为真 panic（可调试性更好） |

净结论：**P 在 finally/synchronized 上更省样板，F 在 catch-then-continue 上更省样板**；但 P 独有"哨兵纪律 + 重抛正确性"这一运行时正确性风险，F 独有的"签名污染"是机械性、编译期可见的成本。ADR-0005 的隔离层能锁住两方案切换成本，本维度不构成压倒性判据，但"发射正确"难度上 F 略低。

# Findings

**测量事实**
1. P 的"正常路径零成本"**只对无 defer 方法成立**（3.6 ns，等于纯调用链）。一旦方法拥有 handler/finally/synchronized，P 反而变慢：finally 纯 defer ~7 ns（比 F 慢 ~1.5×），catch 的 defer+recover 13.2 ns（比 F 慢 ~2.9×）。**"零成本"不是 P 的默认形态，是它的一个端点。**
2. F 的正常路径检查极便宜：+0.24 ns/方法边界（4.6 − 3.6 除以 4 层），远低于 P 的 finally（+0.7 ns/层）与 catch（+2.4 ns/层）开销。F 的"多一次分支"几乎免费。
3. 异常路径 F 全面碾压：throw-catch 快 ~77×、深传播（自由展开）快 ~72×。P 的异常路径贵在 recover 的 unwinder 走查 + 类型断言 + 标匹配 + 重抛；逐层重抛（P 最坏形态）更贵到 1595 ns（~550×），证明"每层 recover 后重 panic"这一 Java 异常表字面语义在 Go 里代价极高，应避免为中间层无 handler 的方法生成 defer。
4. 两方案稳态零堆分配。

**推断（由确定性证据支撑）**
5. P 的 catch 路径贵在 recover 不可内联 + unwinder 帧走查（源码 `gorecover` 的 unwinder 逻辑），是结构性成本，非实现细节。
6. 发射器若"无差别给每个方法加 defer+recover"会把 P 变成 2.9× 慢于 F 的灾难；P 的可行性严格依赖"只在真有 handler 的方法上装 defer"这一发射纪律。

**临界频率与 99.9% 正常路径的答案**（用链级数字，见下表）

总成本 = (1−f)·正常路径成本 + f·异常路径成本。临界频率 f* = (F_norm − P_norm) / (F_norm − P_norm + P_throw − F_throw)。

| 模型 | 数字（ns） | 临界异常频率 f* | 99.9% 正常（f=0.1%）时更优方 |
|---|---|---|---|
| 无 handler 方法 + 浅捕获 | P_norm 3.6 / F_norm 4.6 / P_throw 170 / F_throw 2.2 | **≈0.6%** | P（3.77 vs 4.60，省 ~18%） |
| 无 handler 方法 + 深传播(自由展开) | P_throw 208 / F_throw 2.9 | ≈0.5% | P（省 ~18%） |
| 无 handler 方法 + 逐层重抛 | P_throw 1595 | ≈0.06% | P（边界） |
| finally/synchronized 方法 | P_norm ~7 | **0%（任何频率 F 都赢）** | **F** |
| catch handler 方法 | P_norm 13.2 | **0%（任何频率 F 都赢）** | **F** |

**明确回答**：若正常路径占比 99.9%（典型服务器负载），**只有"方法完全不带 try/catch/finally/synchronized"的那部分**，P 才便宜约 18%；任何带 finally/synchronized/catch 的方法，F 在任何异常频率下总成本都更低。临界异常频率 ≈ **0.6%（浅捕获）到 0.06%（逐层重抛），典型值远低于 1%**——异常一旦超过这个频率，连无 handler 方法也转为 F 占优。

# Recommendation

**v1 发射器采用方案 F（旗标返回）**，P（panic/recover）作为 v2 对"已证实热点、且方法无 defer"的内层链路的可选优化，受 ADR-0005 隔离层约束、按 M4 再评。三维修据：

1. **性能**：P 的唯一优势（无 handler 正常路径省 ~18%）是**窄幅、且被"方法必须无 defer"严格限制**的；真实 Java 代码里 synchronized/finally/catch 并不罕见，每出现一处就把优势翻转为 F 占优。F 的代价（+0.24 ns/边界）几乎可忽略，且换来异常路径 ~70–80× 更快。异常频率只要不是"极低 + 全链路无 defer"的理想化情形，F 的期望总成本更低、且**可预测**（P 的成本是二元的，取决于方法内容）。
2. **发射器复杂度**：P 在 finally/synchronized 上省样板，但**catch-then-continue 必须函数拆分**（recover 语义是"返回给调用者"，不是"跳回 catch 后继续"），且背负"哨兵类型断言 + 非哨兵重抛"的运行时正确性纪律——一个漏重抛就静默吞异常。F 的签名污染是机械、编译期可见、集中在 ABI 层的成本，且 catch-then-continue 原生表达。对"正确性闭环"的 v1，F 的发射更不易错。
3. **语义风险**：P 独有的跨 goroutine 进程宕机陷阱（漏网异常 = 杀进程，必须每个线程顶层 recover 兜底）与"误吞 Go 运行时 panic"风险，F 均无（错误是值，Go panic 保持真 panic）。栈溢出两方案都无法 recover（Go 是致命 throw），需栈深 guard 补齐，此项正交。

# Confidence

**0.68**

性能结论（三场景中位数、方向、临界频率数量级）在两轮独立运行 + 正确性护栏下高度一致，单独置信 ~0.85；但综合置信度被两处未决因素拉低：(a) **方法构成**——"带 defer 方法占比"是决定 P 还是 F 更优的关键变量，本报告未对 Catty 目标 fixture 实测该占比；(b) 发射器复杂度与语义风险是定性判断。0.68 是"据此直接定稿 ADR-0009"的综合置信度。

# Unknowns

- **方法构成未实测**：Catty 目标代码（CollectionsDemo 等）中带 try/catch/finally/synchronized 的方法占比未知，而它直接决定 P 的 18% 优势能否兑现。建议 Orchestrator 在 M3 发射器输出真实 AOT 产物后，对该占比做一次静态统计。
- 本基准是合成微基准，非真实发射器产物；真实 AOT 的寄存器分配、内联预算、逃逸分析可能与手写链不同（尤其 P 的"handler 写命名返回值"形态可能触发逃逸）。
- 未测多层嵌套 try/catch、多 handler 类型（instanceof 多类标匹配）、finally 覆盖在飞异常（Java 的"finally 抛新异常压制旧异常"）等深层语义的发射复杂度——仅定性。
- 临界频率的链级 vs 每边界两种口径（0.6% vs ~0.14%）相差 ~4×，取决于异常通常被捕获的深度；本报告以"浅捕获"为最常用情形取 0.6% 作头数，深层传播会显著拉低它。
- 未测 `GOMAXPROCS>1` 下 panic 的全局锁竞争（`paniclk`、`runningPanicDefers` 计数器）；高并发多线程高频异常时 P 可能更差。
