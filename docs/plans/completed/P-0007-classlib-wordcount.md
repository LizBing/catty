# P-0007: 类库扩张与真实程序攻坚（WordCount / HashMap 面）

- 状态：completed（追溯归档；计划文件补立于 2026-08-24）
- 风险评分：Medium（类库面扩张触碰分派与装箱语义）

## Goal

补齐真实 Java 程序所需的最小类库面（集合、包装类、String 宽方法），
以 WordCount 端到端 fixture 验收，并为 AOT 分派铺路。

## Outcome

- HashMap/HashSet natives + Integer/Long/Boolean 等包装类 + String 宽方法
  落地（internal/kernel/bootstrap_p7.go / natives_p7.go）；
- WordCount 解释路径与 JVM 逐字节一致；隔离测试 TestHashMapDispatch；
- 暴露 AOT cat2 分派缺口 → 登记 DEBT-0015/0017（均已收口）。
