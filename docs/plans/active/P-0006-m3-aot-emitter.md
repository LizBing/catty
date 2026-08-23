# P-0006: M3 — AOT 发射器

- 状态：active（本文件即 M3 规划基线；实现按 Phase 推进）
- 风险评分：High（项目核心卖点；架构级不可逆决策集中地）
- 决策纪律：调用约定与异常机制走 Research→Design→Critic→ADR 全流程（协议 §7）

## Goal

将 javac 产物编译为 Go 源码并整体构建为静态二进制；同一 fixture 在
AOT 与解释两条路径下行为一致；产出首批启动/运行对照数据。
**M3 完成定义 = 正确性闭环 + 对照表，不含吞吐优化（M4 范畴）。**

## Current State

解释器内核完成（M2 收口）：统一对象表示 kernel.Instance、完整线程/监视器/
验证器、四个真实 fixture 验收。R-0001 已证明：混合表示相对统一表示在分派上
有 5–11% 优势——但见下方"表示法判断"，该结论不直接适用于发射器 v1。

## 核心技术判断（规划期预决策，实现前经 ADR 复核）

### 判断一：v1 发射器采用统一表示（kernel.Instance），非混合表示

R-0001 比较的是"两种表示各自写成手写 Go"的场景；而 AOT 相对解释器的
收益来源是**执行模型差异**（消除 opcode dispatch、操作数栈封箱、指令
边界检查），估算为一个数量级以上的解释开销消除，远大于表示法的 5–11%。
v1 选统一表示的理由：
1. **AOT/解释混合互操作零成本**（同一 fixture 双路径一致性验收的前提）；
2. 内核 natives/线程/监视器全部面向 Instance，避免双表示桥接层；
3. R-0001 的混合表示优势记录在案，作为发射器 v2 的 accessor 内联优化方向
   重评（届时 ADR-0004 正式收口）。

### 判断二：异常机制由实验裁决，裁决前隔离

ADR-0005 遗留决策在本阶段必须收口。启动 R-0002 基准
（panic/recover vs 显式旗标返回，热点+异常路径双测），数据回写后以
ADR-0009 定稿。发射器所有方法签名经由统一的调用约定模块生成，
切换成本被限制在约定层内部。

### 判断三：v1 不设独立 IR

classfile 语义分析（栈映射已由验证器完成）后直接发射 Go 源码。
理由：MVP 复杂度最低；优化 pass（去虚化已由类型系统涌现、常量折叠等）
在 v2 引入时再抽 IR，避免过早抽象。

## Tasks（DAG）

```
Phase 0 研究（并行）
R-0002 异常机制基准 → ADR-0009 输入        [Researcher 后台]
R-0003 Go 编译器对生成代码模式的优化行为
       （内联预算/逃逸分析/接口去虚化实证）  [Researcher 后台]

Phase 1 规格与 ADR
T1 调用约定规格 specifications/emitter-abi.md
   （命名混淆规则、签名形态、异常通道、clinit/synchronized/SOE 映射）
   ↑ R-0002
T2 ADR-0009 异常机制定稿                    ↑ R-0002
T3 ADR-0010 表示法与互操作边界定稿          ↑ 本计划判断一评审

Phase 2 实现
T4 发射器骨架：类→包布局、常量池固化、字符串/数组字面量、静态字段
T5 方法体发射：v1 指令子集（对齐 CollectionsDemo 所需集合）
T6 运行时桥：AOT main 入口、双路径一致性 harness
   ↑ T4..T6 互相依赖，顺序实施

Phase 3 验收
T7 CollectionsDemo/ThreadsDemo 纯 AOT 运行，输出与解释路径及参考 JVM 一致
T8 启动/首响对照表（AOT 二进制 vs 解释器 vs java），方法学按 performance.md
```

## Validation

- `make check` 绿；全仓 -race 绿
- 双路径（AOT/解释）对全部既有 fixture 输出逐字节一致
- 对照表落盘 docs/research/R-0004（含方法学与原始数据）

## Risks

- Go 编译时间随生成代码量增长 → MVP 只发射验收所需类，量化后设阈值
- panic 跨 goroutine 语义（若选 panic 方案）→ ADR-0009 明确线程边界行为
- 生成代码可读性差导致调试困难 → 发射时保留源行号注释（LineNumberTable）

## Progress

- [x] 计划建立（2026-08-23）
- [x] R-0002 完成：推荐 F 方案（旗标返回），临界频率 0.6%
- [x] R-0003 完成：10 条发射器形态建议，Confidence 0.82
- [x] T1 emitter-abi.md v0.1 落盘
- [x] T2 ADR-0009 定稿（F 方案 v1；P 降级 v2）
- [ ] T3 ADR-0010 表示法边界 → 下轮评审
- [ ] T4..T8 发射器实现 → **下轮整段专注预算**
      当前进度：genrt 桥 ✓ / 命名混淆 ✓ / 方法发射器骨架 ✓
      handler-entry 深度缺口已关闭（2026-08-23）：异常接收者推送改到
      handler 标签之后发射（JVMS §2.6 栈清空 + 异常入栈，规范深度恒 1，
      不再依赖 StackMapFrame 命中）；CollectionsDemo AOT 输出与解释路径/
      参考 JVM 一致（回归：internal/emitter/body_test.go、
      internal/vm/aot_acceptance_test.go）
      余项：T8 启动/首响对照表
