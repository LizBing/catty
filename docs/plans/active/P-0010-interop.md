# P-0010: Java↔Go 互操作产品化（卖点④落地）

- 状态：active
- 风险评分：Medium（新增 API 面，机制已验证；风险集中在映射表完备性与文档）
- 决策依据：ADR-0011（反射绑定默认面 + DefineClass 逃生舱）；R-0008

## Goal

把"Java 调 Go 函数"从 raw DefineClass 提升为一行注册的产品面
`interop.Bind`，钉扎类型映射全表与错误传播，交付一个真实嵌入演示与调用税
数据，兑现卖点④的首个证据环。

## Tasks（DAG）

```
M1 interop 包      Bind/BindSpec + reflect wrapper 构建 + 注册期 fail-fast；  ✅
                   类型映射全格钉扎（nil 抛错/空串/NaN±Inf/Value 透传身份/
                   (T,error) 两形态/fail-fast 不留半注册类）
M2 错误与线程契约   error→RuntimeException 传播钉扎；DEV-0009 登记；           ✅
                   异步 fire-and-forget 模式文档化
M3 演示资产         cmd/embeddemo：Go 宿主起本地 HTTP + md5 + fetch，          ✅
                   Java 规则消费三项能力（md5(catty)=b4edb5…实测通过）
M4 文档与基准       docs/interop.md（映射表+契约+嵌入清单）；调用税 ~82ns      ✅
                   已入 R-0008 §2 与 interop.md 头部；ledger DEV-0009
```

## Validation

- make check 绿；全仓 -race 绿；既有钉扎零回退
- TestInteropFromAOT 扩展为映射表全格矩阵测试
- demo 可复现脚本化运行，输出与预期一致

## Risks

- 反射 wrapper 的 CallContext 使用越界 → wrapper 设计在返回前完成读取，
  测试覆盖嵌套 ctx.Invoke 场景
- stub 遮蔽依赖注册表优先序 → 回归测试已在（TestInteropFromAOT）
