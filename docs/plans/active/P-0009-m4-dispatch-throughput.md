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
T1 单态内联缓存：CallVirtual 接收者动态类 [class→Method] 单槽缓存
   ↑ 基线：先重跑 R-0005 四基准固化 pre 数据（tickMillis 口径不变以保可比）
T2 natives 白名单直调展开：HashMap.get 等高频方法发射符号引用绕过字符串解析
T3 StringBuilder 五连模式折叠（new+dup+init+append*+toString）
T4 复测与报告：R-0006 落盘 docs/research/，含方法学、原始数据与噪声分析
T5 （择机）Bench 切回 nanoTime(J) 口径重校 DEBT-0017 解锁后的长时钟数据
```

## Validation

- `make check` 绿；全仓 -race 绿
- 全部既有钉扎测试（fixture 双路径、JsonDriver e2e、TraceProbe 堆栈形态）不回退
- T4 报告附 benchstat 显著性检验

## Risks

- 内联缓存的失效语义（类卸载不存在，但 Install 多内核场景需失效化——
  同 genrt resolved 缓存教训）
- 直调展开绕过 ResolveMethod 的动态性：仅限 final/static/已证单态调用点
