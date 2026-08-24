# P-0008: 稳态吞吐证据链（路线 B）+ Route C 语义保真

- 状态：completed（追溯归档；计划文件补立于 2026-08-24）
- 风险评分：High（性能主张必须有 benchmark 证据，协议 §26）

## Goal

建立 AOT vs 解释器 vs HotSpot JIT 的稳态吞吐对照表（卖点②证据链），
并以真实三方库 minimal-json 为靶进行解释器/AOT 双路径语义保真攻坚
（Route C）。

## Outcome

- R-0005 吞吐报告落盘：算术密集 AOT 2.3× 于解释器；调用密集持平，
  瓶颈定位在运行时分发层（CallVirtual→InvokeAs→Native），优化方向按
  ROI 排序（单态内联缓存 / natives 直调 / SB 链折叠）→ 移交 P-0009；
- netStackEffect 完整表 + 哨兵测试（Route B 奠基）；
- Route C 收口 DEV-0002（类初始化顺序）、dup 族语义、Writer/Math/fill/
  parse-radix 面；暴露 DEBT-0016 残余与 DEBT-0018/0019（均已收口，
  见 debt/register.md 归档区）。
