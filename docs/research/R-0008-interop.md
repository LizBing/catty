# R-0008: Java↔Go 互操作产品化研究 — 能力审计、可行性 spike 与调用税

- 状态: 完成
- 关联: P-0009 后续；vision.md 卖点④（差异化最强）；ADR-0011 的定稿输入
- 环境: darwin/arm64 · go1.26.5 · JDK 25 HotSpot

## 1. 现状审计（能力已存在，产品面为零）

| 构件 | 现状 |
|---|---|
| `Kernel.DefineClass(ClassDef)` | ✅ 合成类注册，出生即初始化；全部 ~60 个 bootstrap 类由此运作，机制经 48 类 AOT 混跑验证 |
| `MethodDef.Native` | ✅ Go 函数挂为方法实现，解释/AOT 双路径同权分派 |
| `NativeFunc(ctx,recv,args)` | ⚠ 原始 Value 层签名——embedder 需手写描述符、类型断言、槽语义 |
| `CallContext` | ✅ 提供 K / Owner / Stringify / Invoke(嵌套 Java 调用) / NewStringGo / Throw |
| 注册表遮蔽 | ✅ 合成类先于 classpath 解析器命中（stub 类可安全共存） |

**结论**：Java→Go 方向"今天就能跑"（spike A 证明）；缺的是 embedder 人体工学
与契约文档，不是新机制。

## 2. 可行性 Spike

### Spike A — 正确性（TestInteropFromAOT，已钉扎）

Go 侧 raw DefineClass 注册三函数，AOT 发射的 GoCall.main 调用：

| 用例 | 结果 |
|---|---|
| add(20,22) → int | 42 ✓ 双向 int32 转换 |
| greet("catty") → String | "hello, catty!" ✓ JString↔Go string |
| fail() 抛错 | catch 到 RuntimeException("go-side failure") ✓ error→Java 异常 |

### Spike B — 跨边界调用税（BenchmarkInteropNativeCall）

完整链路 CallVirtualIC → invokeChecked(FrameMeter) → InvokeAs(帧追踪+
ctxPool) → NativeFunc 包装 → Go 函数：

```
BenchmarkInteropNativeCall-8   81.17 / 82.73 / 83.80 ns/op
```

对照 R-0006 虚调用 74ns：**边界本身近乎零附加成本**（差异在 ctx 池化取还与
包装层）。设计含义：互操作不需要任何"快路径"特化；要控的成本只在注册期
（反射构建 wrapper 一次性），不在调用期。

## 3. 设计空间分析

### 3.1 API 形态（三选一）

| 方案 | 形态 | 优点 | 缺点 |
|---|---|---|---|
| A 现状 | 手写 ClassDef | 全功率 | 人体工学差，无文档，易错 |
| **B 反射绑定** | `interop.Bind(k, Class, map[string]any{...})`，reflect 构建 wrapper | 一行注册；签名即文档；wrapper 闭包持有类型化函数→调用期零反射 | 反射魔法面；受限类型集需明确 |
| C 显式描述符+助手 | 仍传 desc，但提供 Value 转换助手 | 可预测 | 比 B 略好于 A，仍繁琐 |

**推荐 B 为默认 + A 作为逃生舱**（Bind 底层就是 DefineClass）。

### 3.2 类型映射表（v1）

| Go | Java | 备注 |
|---|---|---|
| string ↔ String | 内容转换（MakeJStringFromGo / Go()） | 非 interned，身份语义保持 |
| int32/int64 → I/J | 值拷贝 | int 平台字长映射 J 并在文档声明 |
| float32/float64 → F/D | 值拷贝 | NaN/Inf 直通 |
| bool → Z | 0/1 | |
| error（末位返回） | 非 nil ⇒ Throw RuntimeException(msg) v1；专用 go/GoError v2 候选 | msg 承载 err.Error() |
| *kernel.Value 级透传 | 不透明句柄对象 | 允许 embedder 存取 Java 对象引用而不解构 |
| struct/指针/slice/map | **v1 不映射** | 明确报错；v2 讨论句柄注册表 |

多返回 `(T, error)` 即 Java 方法形态；`(T)` 视为不抛错；仅 `error` 即 void。

### 3.3 线程契约（v1 文档化，不做机制）

- 所有调用同步发生在调用方 goroutine（=Java 线程）上；
- Go 侧阻塞会阻塞该 Java 线程，且 **Thread.interrupt 无法唤醒**（登记为
  interop 已知限制，DEV 序号随实现轮分配）；
- 异步模式：embedder 自行起 goroutine 并立即返回（fire-and-forget），回投
  Java 需经 ctx.Invoke——v1 仅文档化模式，不做调度器。

### 3.4 错误传播层级

```
Go panic（非 error）      → 引擎错误（fatal，进程级）——bug 语义不变
error 返回                → Java 异常（RuntimeException v1）
ctx.Throw(类名,msg)       → 任意 bootstrap 异常（A 方案逃生舱已有）
Java 异常传入 Go 函数参数  → 以 Value 透传，不解构
```

### 3.5 与 AOT/懒加载的交互

合成类在注册表中先于解析器命中 → 发射体 CallStatic 直接解析到 Native ✓；
Install 懒加载钩子不受影响。唯一注意点：Bind 必须在首次分派前完成（与
gen.Install 同一时段），文档化为嵌入清单步骤。

## 4. 实施计划指向（P-0010）

M1: `internal/interop` 包 —— Bind/BindSpec + 类型映射表 + 错误传播 +
    与 DefineClass 的桥接；钉扎测试覆盖映射表全格。
M2: 线程/阻塞文档 + DEV 登记；示例 cmd/demo（Go HTTP server 内嵌 Java 规则
    或反向）。
M3: 调用税进 R-0006 表族；emitter-abi/新 interop.md 文档。

估算：M1 S-M、M2 S、M3 S——总量一轮内可完成（研究轮已消解全部设计不确定性）。

## 5. 风险

- 反射绑定对未支持类型的失败必须发生在**注册期**（fail-fast），不得延到调用期；
- ctx 池化的 CallContext 在用户函数里被长期持有会串数据——Bind 生成的
  wrapper 在函数返回前完成全部读取（文档 + wrapper 设计保证）；
- stub 类遮蔽依赖注册表优先序——若未来解析顺序变化需保留回归测试。
