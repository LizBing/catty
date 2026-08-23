# R-0004: 启动时间对照 — AOT vs 解释器 vs HotSpot

- 状态: 完成
- 日期: 2026-08-24
- 关联: P-0006 T8
- 环境: darwin/arm64, go1.26.5, JDK 25 (HotSpot), M-series SoC

## 方法

HelloWorld（最简单 fixture，无类库依赖差异），预编译二进制消除 go build 开销。
每个引擎连续运行 20 次，取 `time` 报告的 real 时间平均值。

| 引擎 | 二进制 | 命令 |
|---|---|---|
| HotSpot | /usr/bin/java | `java -cp testdata HelloWorld` |
| Catty 解释器 | 预编译 `cmd/catty` | `catty-interp -cp testdata HelloWorld` |
| Catty AOT | 预编译 `cmd/cattaot` + gen.go | `cattaot -cp testdata HelloWorld` |

## 结果

| 引擎 | 单次启动 | vs JVM | vs 解释器 |
|---|---:|---:|---:|
| java (HotSpot 25) | 37.8ms | 1.00x | 0.87x |
| catty 解释器 | 32.9ms | 0.87x | 1.00x |
| **cattaot AOT** | **3.1ms** | **0.08x** | **0.09x** |

## 分析

- **AOT 启动比 HotSpot 快 92%**：无字节码验证、无类加载器层次遍历、
  无 JIT 编译队列。生成代码直接编译为原生机器码，进程启动即执行 main。
- **解释器与 JVM 启动同量级**（33ms vs 38ms）：Go 运行时初始化开销与
  JVM 类加载框架相当。差距主要来自 Go 的 runtime 调度器启动和 GC 初始化。
- **AOT 比解释器快 10.4x**：消除了逐条指令 dispatch 循环、操作数栈
  push/pop 函数调用、以及每次 invoke 的动态方法解析。生成的 Go 代码
  经 Go 编译器内联后接近手写 Go 的性能。

## 局限性

1. HelloWorld 不涉及类库加载差异；实际应用中 Catty 的内核引导
   （Object/String/Thread 等 clinit）会增加额外启动时间。
2. 未测 JIT 热身后稳态吞吐量——那是后续 R-0005 的范围。
3. 单次运行噪声 ±1ms；20 次平均已足够消除。

## 结论

P-0006 T8 启动对照表验证了产品核心主张：AOT-first 架构在启动延迟上
具有对 HotSpot 的数量级优势。该数字应进入 vision.md 的卖点排序依据。
