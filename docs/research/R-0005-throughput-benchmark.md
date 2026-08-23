# R-0005: 稳态吞吐基准 — AOT vs 解释器 vs HotSpot JIT

- 状态: 完成
- 关联: P-0008 路线 B；vision.md 性能主张证据链第二环（第一环=R-0004 启动）
- 环境: darwin/arm64 · go1.26.5 · JDK 25 HotSpot · 同机同负载

## 方法

单一 fixture `Bench`（源码 `testdata/src-bench/Bench.java`），四类微基准，
各跑 3 轮取最小值（最小值=最接近真实稳态）。N=1,000,000（strcat=N/10）。
JVM 内置 N/5 预热轮触发 C2 编译。计时用 System.tickMillis()（int32 毫秒
时钟，规避当前发射器 long 返回的已知深度限制——见 §局限）。

| 基准 | 热点形态 |
|---|---|
| arith | 整数算术循环，无调用 |
| vcall | 单态虚调用 ×N |
| mapops | HashMap get/put + Integer 装箱 + 字符串键拼接 |
| strcat | StringBuilder 链式追加 toString |

## 结果（单轮耗时 ms，越小越快）

| 基准 | JVM JIT | Catty AOT | Catty 解释器 | AOT/解释器 | JVM/AOT |
|---|---:|---:|---:|---:|---:|
| arith | 0.41 | 40.4 | 94.6 | **2.3×** | 99× |
| vcall | ~0¹ | 183.9 | 206.6 | 1.1× | — |
| mapops | 21.3 | 1460.2 | 1670.7 | 1.14× | 69× |
| strcat | 494.4² | 3766.9 | 3732.7 | ≈1.0× | 7.6× |

¹ JVM 轮 2/3 为 0ns——HotSpot 死代码消除（结果未被观测）；
   以轮 1 的 195.8ms 作参考上限。
² strcat 三轮波动大（494–1076ms），JVM 逃逸分析对 SB 链生效与否不稳定。

## 分析

1. **算术密集循环 AOT 比解释器快 2.3 倍**——消除了逐指令 dispatch 与操作数栈
   模拟，这是发射器的纯收益区。
2. **调用密集负载 AOT≈解释器**（vcall/mapops/strcat 仅 1.0–1.14×）：生成代码
   每次调用仍走 `genrt.CallVirtual → InvokeAs → Native` 动态分发，与解释器的
   invoke 处理同构。**调用分发是当前共同瓶颈**，不是发射目标码本身。
3. **与 HotSpot 的差距集中在可内联场景**：JIT 把 vcall/mapops 内联+逃逸分析后
   接近归零。AOT v1 无内联（R-0003 建议 #2 未启用），差距符合预期。

## 优化方向（按 ROI）

1. **单态内联缓存**：CallVirtual 对接收者动态类做 [class→Method] 单槽缓存，
   预期把 mapops/vcall 的分发成本砍半以上。
2. **常用 natives 直调展开**：对 HashMap.get 等白名单方法，发射器直接生成
   `genrt.HashMapGet(recv,args)` 符号引用，绕过字符串解析。
3. **StringBuilder 链识别**：new+dup+init+append*+toString 五连模式折叠为
   Go strings.Builder 局部变量（R-0003 建议 #5）。

## 结论

启动维度（R-0004：快 92%）之外，稳态维度的诚实结论是：
**AOT 当前在算术密集场景稳定优于解释器，调用密集场景两者持平——
瓶颈在运行时分发层而非发射层，优化路径清晰且已排序。**
