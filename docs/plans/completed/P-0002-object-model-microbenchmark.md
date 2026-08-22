# P-0002: 对象模型表示法微基准（裁决 ADR-0004）

- 状态：completed（2026-08-23）
- 优先级：高（截止：M3 发射器动工前）
- 风险评分：Low（纯实验，不改产品代码）
- 成本预估：~0.5 agent 日

## Goal

用实测数据裁决 ADR-0004 的"具体类型 + 接口混合表示" vs "统一 `*Object`"，把 ADR-0004 从 Provisional 变为 Accepted 或 Superseded。

## Hypothesis

- H1（混合表示）：虚分派开销接近 itab 直查，字段访问零装箱，且 Go 编译器能穿透窄类型做逃逸分析；
- H2（统一表示）：每次访问需类型断言/字段间接寻址，深继承链下放大。

## Method

1. 手写两个独立 Go 原型（不进主仓，放 `docs/research/assets/objmodel/`）：
   - A：embedding 继承链（深度 3 与 6 两档）+ interface 提升；
   - B：统一 `*Object` 头 + 断言访问；
2. 四组基准：单态调用点、双态/ megamorphic 调用点、字段读写、接口断言下行转换；
3. `-benchmem -count=10`，benchstat 出 median/p95；
4. 附带观察：两版原型经 `go build -gcflags=-m` 的逃逸分析差异（定性记录）。

## Success Criteria

- benchstat 结果表进入 research 报告（模板见 docs/research/TEMPLATE.md）；
- ADR-0004 更新 Evidence 与 confidence（≥0.8 方可转 Accepted）；
- 结论直接影响发射器调用约定的样板设计。

## Constraints

- 不修改 production 代码；不阻塞 P-0001（解释器阶段两种表示法不影响验收）。

## Progress

- [x] 完成（2026-08-23）：R-0001 报告 + objbench 原型落盘；ADR-0004 证据已回写，
      状态维持 Provisional（0.72 < 0.8），收口条件见 ADR

## Open Questions

- 是否需要第三方案 C（代码生成展开的伪 vtable）——由 Researcher 判断成本后决定。
