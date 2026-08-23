# Question

M3 AOT 发射器（P-0006 判断三：classfile 语义分析后直接发射 Go 源码）将产出大量模式化 Go 代码。这些代码的运行效率取决于 gc（Go 编译器）对**这些特定模式**的优化行为。本报告用实验回答：哪些手写习惯在生成代码中**会失效或生效**，从而让发射器直接产出最优形态。

逐项验证六个问题（go 1.26.5，GOARCH=arm64 / Apple M1）：

1. **接口去虚化**：当 interface 变量的具体类型可静态证明时，gc 是否内联？内联预算（80 cost units）对多大方法体生效？
2. **方法值与闭包**：用闭包捕获局部变量模拟字节码局部性，逃逸成本如何？
3. **大函数**：单 Go 函数对应一个大 Java 方法（数百字节码折叠），gc 内联/寄存器分配是否退化？有无规模上限？
4. **字符串拼接链**：连续 `+` 拼接 vs `[]byte` append / `strings.Builder`，应选哪种？
5. **类型 switch 与断言链**：多类型分派（模拟 instanceof 链）的最优形态。
6. **边界检查消除（BCE）**：下标循环写成什么形状可稳定触发 BCE。

# Context

关联 `docs/plans/active/P-0006-m3-aot-emitter.md`。P-0006 判断一选统一表示（`kernel.Instance`），判断三 v1 不设独立 IR、直接发射 Go 源码——因此**发射器产出的 Go 代码形态就是最终被 gc 编译的形态**，编译器对这些形态的优化行为是硬性输入，而非可事后调优项。协议 §26 要求性能主张必须有 benchmark 证据；本报告所有结论均由 `-gcflags=-m/-S`、`-d=ssa/check_bce` 与 `-benchmem` 基准佐证。

# Candidates

每问给出可选的发射形态（Candidates），实验裁决：

1. 接口直调 vs 断言后直调 vs 手写伪 vtable vs 接受 itab 间接调用。
2. 闭包捕获局部变量 vs 显式结构体状态 vs 方法值。
3. 单大函数 vs 拆分为多个小函数。
4. `+` 链 vs `strings.Builder` vs 预分配/裸 `[]byte` append。
5. `switch v.(type)`（绑定/不绑定）vs if-else 断言链 vs int 标签 switch。
6. `for i := range` vs `for i := 0; i < len(s); i++` vs 手写 len 缓存 vs 逆向循环。

# Evidence

所有原始数据在 `docs/research/assets/gcopt/results/`（编译报告 + `bench_all_count10.txt` + `benchstat_summary.txt`），复现命令见该目录 `README.md`。工具链：`go1.26.5 darwin/arm64`（Apple M1）。噪声说明沿用 `assets/objmodel` 方法论：3–5ns 级样本 ±15–25% 波动，结论以**确定性编译报告 + 精确 alloc 计数**为主，timing 差异 <5% 视为噪声。

---

## Q1 接口去虚化与内联预算

**设置**：`Adder` 接口（`Add(a,b int64) int64`）+ 具体类型 `C1/C2`；探针函数覆盖四种调用形态；`M00..M40`（0–40 条 `x = x*2 + c` 语句）定位预算边界。

**证据摘录**（`results/q1_inline_m.txt`、`results/q1_asm.txt`）：

```
q1_devirt/types.go:42:14: devirtualizing v.Add to *C1   // 直接赋值后接口调用
q1_devirt/types.go:42:14: inlining call to (*C1).Add
q1_devirt/types.go:49:15: inlining call to (*C1).Add    // 断言后调用
（ProbeOpaqueParam / ProbeField 无任何 devirtualizing 消息）

q1_devirt/budget_gen.go:116:6: can inline (*C1).M11 with cost 79
q1_devirt/budget_gen.go:132:6: cannot inline (*C1).M12: cost 86 exceeds budget 80
```

反汇编佐证：`ProbeInterfaceLocal` 16 字节纯算术（`ADD; ADD $1; RET`，无 itab 加载、无 CALL、字段常量折叠）；`ProbeAssertThenCall` 一条 `CMP itab; BNE` + 内联体；`ProbeOpaqueParam` 为 `MOVD 24(R0),R4; CALL (R4)`（真间接调用，`rel … t=R_CALLIND`）。

**基准**（M1，count=10，0 allocs 全项）：Direct 2.25n / InterfaceLocal 2.45n / AssertThenCall 2.17n / OpaqueMono 2.11n / OpaqueBi 2.13n / Field 2.15n / BudgetUnder(M05) 2.77n / BudgetOver(M30) 7.14n。

**结论**：
1. gc 在**具体类型可静态证明**时去虚化并内联，且发生在接口调用点上（不只是断言之后）。两个可证明场景：①同函数内 `var v I = &C{n:1}; v.M()`（"devirtualizing"）；②`if c, ok := v.(*C1); ok { c.M() }`。接口跨函数边界（参数）或存于 struct 字段时**不去虚化**，保持 itab 间接调用。
2. 内联预算恰为 80 cost units；本形状每条简单语句 ≈7 units，故 **~11 条简单语句是内联上限**（M11=79 内联，M12=86 拒绝）。`Add`（1 条语句 + 字段读）cost 7。
3. 微基准全部落在 ~2ns 的全局 `sink` 存储延迟地板上，去虚化/间接调用的绝对差异被掩盖——但确定性 `-m`/反汇编证据是精确的，且 alloc 计数全为 0（无装箱逃逸）。

---

## Q2 闭包与方法值逃逸

**设置**：闭包按值/按引用捕获局部变量、同步调用 vs 返回逃逸；具体/接口方法值；noinline 同包消费者 `callFunc`。

**证据摘录**（`results/q2_escape_m.txt`）：

```
q2_closure/types.go:24:10: inlining call to ProbeClosureByValue.func1  // 同步闭包内联
q2_closure/types.go:30:10: inlining call to ProbeClosureByRef.func1
q2_closure/types.go:38:9:  func literal escapes to heap    // 返回闭包 → 逃逸
q2_closure/types.go:42:2:  moved to heap: x                 // 按引用捕获且逃逸 → x 上堆
q2_closure/types.go:84:9:  func literal escapes to heap
q2_closure/types.go:91:18: moved to heap: x                 // makeCounter 按引用捕获
q2_closure/types.go:75:15: f does not escape                 // callFunc 参数（noinline 同包）不逃逸
```

关键细节：`//go:noinline` 只阻止内联（代码复制），**不阻止逃逸分析**——同包内 noinline 函数同步调用闭包时，闭包仍被证明不逃逸、留在栈上。

**基准**（0 allocs 栏为同步形态；逃逸形态见 alloc）：

| 形态 | sec/op | B/op | allocs/op |
|---|---|---|---|
| DirectCall | 2.07n | 0 | 0 |
| ClosureByValue（同步） | 2.07n | 0 | 0 |
| ClosureByRef（同步） | 2.40n | 0 | 0 |
| ClosureHeap（noinline 同包调用） | 2.15n | 0 | 0 |
| MethodValueConcrete | 2.08n | 0 | 0 |
| MethodValueInterface | 2.11n | 0 | 0 |
| ClosureEscapeByValue | 11.5n | 16 | **1** |
| ClosureEscapeByRef | 22.2n | 24 | **2** |

**结论**：
1. 同步闭包（模拟字节码局部性）**零成本**：内联、0 分配、被捕获局部变量留在栈上（按值与按引用皆然）。
2. 逃逸成本只在闭包**逃出作用域**（返回 / 存入堆对象 / 跨包传递）时发生：按值捕获 1 次分配（+~9n），按引用捕获 2 次分配（闭包 + 被 box 的变量，+~20n）。
3. 方法值：具体类型方法值同步使用 0 分配且可内联；接口方法值 0 分配但调用保持虚分派（捕获 itab）。

---

## Q3 大函数

**设置**：生成 `BigChainN`（N 条非线性 `x = x*1103515245 + c; x ^= x >> 16` 语句的依赖链，N=50..4000）与 `SplitChain200`（8 个 25 步小函数顺序组合）。非线性步骤用于阻止 gc 把仿射链折叠（见下）。

**证据摘录**（`results/q3_inline_m.txt`、`results/q3_asm.txt`）：

```
q3_bigfunc/big_gen.go:4118:6: cannot inline BigChain4000: cost 48002 exceeds budget 80
q3_bigfunc/big_gen.go:4118:6: function BigChain4000 considered 'big'; reducing max cost of inlinees
```

反汇编：`BigChain4000 STEXT size=63744 … LEAF|NOFRAME`（locals=0，**无栈帧、无溢出**）；`BigChain50 size=816`。代码体积线性：816B（50 步）→ 63.7KB（4000 步）。

**基准**（线性、0 allocs）：

| 步数 | sec/op | ns/步 |
|---|---|---|
| 50 | 47.4n | 0.95 |
| 200 | 324.7n | 1.62 |
| 1000 | 1853n | 1.85 |
| 4000 | 7528n | 1.88 |
| SplitChain200（8 次非内联调用） | 321.2n | — |

**结论**：
1. gc 编译**数千条语句的直链函数无性能悬崖**：运行时间随步数线性增长，无栈溢出（leaf、无 frame），寄存器分配不退化。每步成本从 0.95n（50 步）升到 1.88n（4000 步）约 2 倍，源于取指/前端吞吐而非寄存器压力。
2. 80 预算只约束**跨函数内联**，不约束函数自身编译。大方法照常编译，只是不会被内联进调用者；且命中 "considered 'big'" 启发式会**下调其内部被内联对象的预算**。
3. 拆分 200 步方法为 8 个非内联 25 步调用 ≈ 单大函数（321n vs 325n，噪声内）：8 次调用开销 <10% 且被流水线掩盖。
4. **附带发现**：纯仿射链（常量乘数 LCG）被 gc 折叠成单条乘法——发射器生成的循环算术可能被大幅化简（通常是收益，非风险）。

---

## Q4 字符串拼接链

**设置**：`loopN=1000` 动态拼接四形态；固定 4/8 段拼接 `+` vs `Builder`。

**证据摘录**（`results/q4_asm.txt`）：`BenchmarkFixedPlus4` → `runtime.concatstring4`；`BenchmarkFixedPlus8` → `runtime.concatstrings`；`BenchmarkLoopPlus` → 每次迭代 `runtime.concatstring2`。

**基准**：

| 形态 | sec/op | B/op | allocs/op |
|---|---|---|---|
| LoopPlus（`s = s + piece`） | 215.1µ | 2.31MB | 999 |
| LoopBuilder（Grow） | 3.40µ | 9.25KB | 2 |
| LoopAppendPre（预分配 cap） | 2.90µ | 5.25KB | 1 |
| LoopAppendNaive | 3.91µ | 17.5KB | 13 |
| FixedPlus4 | 20.4n | 0 | **0** |
| FixedBuilder4 | 36.0n | 36 | 1 |
| FixedPlus8 | 45.3n | 48 | 1 |
| FixedBuilder8 | 58.4n | 96 | 2 |

**结论**：
1. 循环内 `+` 是**平方级灾难**：215µ vs 3.4µ（约 63×）、999 次分配。
2. **固定少量段**（≤5 操作数）时 `+` 链是零分配且最快：走 `concatstringN`（栈缓冲）。4 段 `+`（20.4n/0 alloc）反快于 Builder（36n/1 alloc）。阈值在 5 操作数：>5 段改走堆分配的 `concatstrings`（FixedPlus8 = 1 alloc）。
3. 动态长度下 `Builder+Grow` 与**预分配** `[]byte` append 同量级（2.9–3.4µ）；预分配 append 分配最少（1 alloc）。裸 append 无预分配多出 12 次分配（13 allocs）。

---

## Q5 类型 switch 与断言链

**设置**：8 个具体类型 + 接口 `I`；三形态：`switch c := v.(type)`（绑定）、if-else 断言链（绑定）、`switch v.(type)` 不绑定后调 `v.Val()`、int 标签 switch。

**证据摘录**（`results/q5_asm.txt`）：

- 类型 switch 编译为**哈希跳转表**（O(1)）：`MOVWU 16(R4) → AND $15 → MOVD (R7)(R6<<3) → JMP`，逐项 itab 校验。
- if-else 断言链编译为**线性 itab 比较**（O(n)），8 路 480 字节。
- 绑定形态 `c.Val()` 内联 → **0 次间接调用**；不绑定形态 `v.Val()` → **8 次间接调用**（每个 case 重新虚分派）。

**基准**（0 allocs 全项）：TypeSwitch2 2.11n / AssertChain2 2.08n / TypeSwitch4 2.09n / AssertChain4 2.09n / **TypeSwitch8 2.09n / TypeSwitchNoBind8 2.58n / AssertChain8 2.14n / TagSwitch8 2.48n**。

**结论**：
1. 类型 switch = O(1) 跳转表，断言链 = O(n) 线性；n≤8 时两者 timing 在噪声内（分支可预测 + 全局存储地板），但类型 switch 渐近更优。
2. **绑定是决定性细节**：不绑定 `switch v.(type)` 后再调 `v.Method()` 每个 case 重新虚分派，8 路慢 ~23%（2.09n → 2.58n）。
3. int 标签 switch（统一表示）2.48n，与 Go 类型 switch 同量级、不更慢——表示法选择不强制分派惩罚。

---

## Q6 边界检查消除（BCE）

**设置**：七个求和函数覆盖 range / 值 range / len 条件 / len 缓存 / 逆向 / 参数 n / 不透明索引；`-d=ssa/check_bce/debug=1` 输出确定性报告。

**证据摘录**（`results/q6_bce.txt`，仅 3 处保留检查）：

```
q6_bce/probe.go:83:14: Found IsInBounds   // SumGlobalCallBefore：访问前有可变调用
q6_bce/probe.go:94:9:  Found IsInBounds   // SumWithN：n 为外部参数，无法关联 len(s)
q6_bce/probe.go:104:9: Found IsInBounds   // SumOpaqueIndex：不透明索引
```

即 **range / 值 range / `i < len(s)` / `n := len(s)` 缓存 / 逆向循环全部证明消除**。

**基准**（0 allocs 全项）：SumRangeIndex 92.9n / SumRangeValue 94.0n / SumForLen 100.8n / SumForCache 93.8n / SumReverse 94.7n / **SumWithN 93.4n**（保留检查）/ SumOpaqueIndex 574n（6×，主要来自 noinline 助手调用，非检查本身）。

**结论**：
1. `for i := range s { s[i] }`、`for i := 0; i < len(s); i++`、`n := len(s)` 缓存、逆向循环都**稳定触发 BCE**（精确证据）。
2. 保留一次标量边界检查在 arm64 上**几乎无成本**（SumWithN 93.4n ≈ SumForLen 100.8n）——BCE 的标量收益在噪声内，其真正价值是清理代码/便于后续向量化，而非标量 ns。
3. 发射器应避免的两种确定破坏 BCE 的形状：①用不透明/外部来源的索引；②在循环体内、数组访问**之前**插入可能改变切片/全局的可变调用。

---

# Findings

1. gc 的去虚化/内联/逃逸分析在**同函数内可见具体类型**时非常激进且精确：接口直接赋值调用去虚化+内联+常量折叠到 16 字节；同步闭包零逃逸；len 缓存与逆向循环全 BCE。发射器应**顺应而非对抗**这些能力。
2. 内联预算 80 units 是硬门槛（~11 条简单语句），跨过即真实函数调用；但这是**跨函数**内联预算，不影响大函数自身编译。
3. `//go:noinline` 阻止内联但不阻止逃逸分析；闭包只有**逃出作用域**（返回/存字段/跨包）才上堆。
4. 字符串 `+` 存在精确分界：≤5 固定段零分配（栈），>5 段与循环场景必须换 Builder/append。
5. 类型 switch 是 O(1) 跳转表；**绑定变量**是正确性+性能的双重要求。
6. 标量边界检查在 arm64 上成本可忽略，但破坏 BCE 的形状（不透明索引、访问前可变调用）应主动避免。

# Recommendation

## 给发射器的直接建议清单

> 证据强度：**强**（确定性编译报告 + 精确 alloc 计数，可复现）＞**中**（确定性报告 + timing，或 timing 受噪声地板影响）＞**弱**（timing 在噪声内或单点观察）。

1. **强** — 热调用点的接收者若是局部可证明的具体类型，直接发射接口调用或"断言后具体调用"；gc 会自动去虚化+内联。不要为这类调用点手写伪 vtable / itab 副本。
2. **强** — 把 hot 方法（accessor、小算术）控制在 ~10 条简单语句内（≤80 cost），使其可内联；不要为"可读性"把一次性逻辑堆进热方法。
3. **强** — 跨函数边界/结构体字段传接口时不要指望去虚化：要么在发射端先断言成具体类型再调用，要么接受 itab 间接调用（其单态预测成本在 M1 上低）。
4. **强** — 用闭包模拟字节码局部性**只要闭包不逃逸**就是零成本的；一旦该"局部"会被返回/存入字段/进 goroutine，改用显式结构体状态，避免按引用捕获变量被搬到堆（+1 分配，慢 ~10×）。
5. **强** — 字符串：固定 ≤5 段（javac 的 `a+b+c`）发射纯 `+` 链（零分配）；动态/循环长度发射 `strings.Builder`（Grow）或预分配 `[]byte` append，**严禁**循环内 `s += x`（平方级、999 分配）。
6. **强** — 多类型分派（instanceof 链/虚分派）发射**绑定变量的类型 switch**（`switch c := v.(type)`）并在 case 内调 `c.M()`；绝不 `switch v.(type)` 后调 `v.M()`（每 case 重新虚分派，慢 ~23%）。
7. **中** — 大量分支（>~8 类型）优先类型 switch（O(1) 跳转表）而非 if-else 断言链（O(n)）。
8. **强** — 数组/切片循环统一发射 `for i := range s` 或 `for i := 0; i < len(s); i++`（含 `n := len(s)` 缓存与逆向循环），稳定 BCE；避免不透明索引与"访问前插入可变调用"两种破坏形状。
9. **中** — 不必为"性能"把一个大 Java 方法拆成多个 Go 函数：数千条语句的直链函数编译/运行均无悬崖、无栈溢出；拆分反而增加调用边界。只有可复用且 ≤10 语句的辅助逻辑才值得拆（拆出后可内联，零成本）。
10. **弱** — 统一表示若用 int 标签分派，与 Go 类型 switch 同量级，不构成换表示法的理由（分派成本不是决策变量）。

# Confidence

0.82

确定性编译报告（Q1 去虚化/预算、Q2 逃逸、Q3 代码体积/无溢出、Q5 跳转表、Q6 BCE 报告）与精确 alloc 计数（Q2 逃逸 1/2 分配、Q4 拼接 0/1/999 分配）是**精确可复现**的，支撑主要建议（强度"强"项）。Confidence 未达 0.9 的原因：①微基准 timing 受 M1 全局 `sink` 存储地板与频率缩放影响，n≤8 分派、标量 BCE 等"量级相同"的结论依赖该地板而非绝对数；②未直接测寄存器压力（>~28 个同时存活局部变量）与向量化交互。

# Unknowns

1. **寄存器压力/溢出行为**：大函数依赖链只用 1 个存活值；未测"数十个同时存活局部变量"（模拟操作数栈多槽位存活）是否会触发溢出与退化。建议后续用 N 个独立累加器（N>28）的 wide 函数探针验证。
2. **BCE 与向量化交互**：标量边界检查成本可忽略；但检查是否阻碍 gc 对可向量化循环（连续 int32/float 求和）的 NEON 自动向量化，未测。
3. **跨包边界对逃逸/去虚化的影响**：本实验全部同包；发射器生成的 AOT 代码与 kernel 运行时跨包时，闭包/接口的逃逸与去虚化行为可能更保守，未验证。
4. **编译时间随生成代码量增长**：观察到大函数编译正常、代码体积线性，但未量化 4000+ 语句方法的编译时长（P-0006 风险项）。
5. **x86-64 与 amd64 目标**：全部数据在 arm64；x86 的 itab 间接调用成本、`concatstringN` 阈值行为、跳转表生成可能不同，未测。
