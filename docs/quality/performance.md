# 性能质量要求

- 地位：协议 §26（Benchmark Protocol）的项目本地化。

## 流程（强制）

```
Baseline → Hypothesis → Implementation → Measurement → Comparison → Decision
```

性能主张必须按 Evidence-before-opinion 格式陈述（协议 §1.3）。噪声范围内的变化**不得**宣称优化成功。

## 统计规则

- 微基准：`-benchmem -count=10` + benchstat；报告 median 与 p95;
- 启动基准：hyperfine，n≥30，报告 mean±sd;
- 噪声带：变化幅度 <±3% 视为噪声;跨机器比较禁止，同机同电源策略;
- 基准结果持久化到 PR 描述或 `docs/briefs/`，不做本地即焚。

## 基准集路线

| 里程碑 | 基准 |
|---|---|
| M1 | 异常机制双版本微基准（ADR-0005）、内核热路径微基准 |
| M2 | JSON HTTP echo：吞吐/延迟分布/内存 |
| M3+ | 冷启动、常驻内存、稳态吞吐三方对比（OpenJDK / native-image / Catty），方法学公开可复现 |

## 目标值（Target，非承诺）

冷启动 <10ms；v1 稳态吞吐达 native-image 水平；p99 无预热抖动。目标修订须走 vision.md + ADR。

## 回归门

M3 起评估将关键阈值（如启动时间上界）写入 make check 的可行性；在此之前靠 brief 中的趋势披露。
