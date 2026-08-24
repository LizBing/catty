# P-0009: M4 — 分发层吞吐优化与三方基准报告

- 状态：active
- 风险评分：High（性能主张必须有 benchmark 证据，协议 §26；噪声范围内
  变化不得宣称优化成功）
- 决策纪律：内联缓存涉及分派语义 → 设计前查 ADR-0010 §2 与 emitter-abi §5

## Goal

在保持双路径逐字节一致的前提下，按 R-0005 已排序的 ROI 清单削减运行时
分发成本，产出第二批吞吐对照数据（AOT vs 解释器 vs HotSpot vs GraalVM
native-image 若可行），兑现卖点②"零预热抖动"的首个证据环。

## Tasks（DAG）

```
T0 基线固化      重跑 R-0005 四基准（tickMillis 口径不变，保可比）
T0.5 分发链剖析  cattaot 加 env 门控 pprof；vcall/mapops 热点画像定靶
T1 分发链削减    单态内联缓存（[site,dyn]→Method 原子单槽）+
                 帧深度计量去全局锁（按证据决定是否入本轮）
T2 natives 直调  白名单方法发射符号引用绕过字符串解析
                 （static/final/已证单态调用点准入）
T3 复测 + R-0006 benchstat 显著性检验；噪声内变化不得宣称优化成功（§26）
T4 并发首证(spike) N 线程工作负载 vs JDK 平台线程；持续负载 p99 采样
T5 可观测性(spike) pprof 演示资产；-race 构建检测 Java 数据竞争实验
T6 捎带(S 成本)   DEBT-0001 fuzz 目标进可选 make 目标；
                  堆栈叶帧行号（发射期 LNT 烘焙，视主线余量）
```

## Success Criteria

- `make check` 绿；全仓 `-race` 绿；既有钉扎零回退
  （fixture 双路径 / JsonDriver e2e / TraceProbe 堆栈形态 / 效果表哨兵）
- mapops/vcall 的 AOT÷解释器比值从 1.0–1.14× 提升至 **≥1.5×**；
  未达标则 R-0006 如实记录并给出瓶颈归因
- 产出：docs/research/R-0006（方法学 + 原始数据 + 显著性）、
  卖点⑤演示资产或精确登记的偏差/债务

## Risks

- 内联缓存失效语义：Install 多内核场景需失效化（genrt resolved 缓存教训已内化）
- 直调展开绕过 ResolveMethod 的动态性：仅限 final/static/已证单态调用点
- spike 结论可能为负（如 -race 检不出 Java 层竞争）：负结果同样落盘登记

## Risks

- 内联缓存的失效语义（类卸载不存在，但 Install 多内核场景需失效化——
  同 genrt resolved 缓存教训）
- 直调展开绕过 ResolveMethod 的动态性：仅限 final/static/已证单态调用点
