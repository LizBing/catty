# Java↔Go 互操作（interop）v1

- 状态：Accepted — ADR-0011 的产品面；实现于 `internal/interop`
- 调用税：**~82ns/次**（BenchmarkInteropNativeCall，R-0008 §2），与虚调用
  74ns 同量级——边界近乎免费。

## 快速上手（Go 宿主侧）

```go
k := kernel.New(kernel.Options{Stdout: os.Stdout})
th := vm.New(k)                       // 解释/AOT 执行桥

err := interop.Bind(k, interop.Spec{
    Class: "com/app/GoBridge",        // Java 侧类名（内部名）
    Funcs: map[string]any{
        "md5Hex":   func(s string) string { ... },
        "fetchLen": func(url string) (int64, error) { ... }, // (T,error) ⇒ 可抛
    },
})

// 之后照常：loader.Load → gen.Install → th.Call(...)
```

Java 规则代码直接调用静态方法：

```java
String digest = com.app.GoBridge.md5Hex("input");
```

编译期可提供一个同名 stub 类供 javac 引用；运行时合成注册表先于 classpath
解析器命中，stub 永不加载（TestInteropFromAOT 钉扎）。

## 类型映射表（v1，ADR-0011 固定）

| Go | Java 描述符 | 备注 |
|---|---|---|
| string | Ljava/lang/String; | 内容转换；非 interned |
| int32 / int64 | I / J | 值拷贝 |
| int | J | 平台字长=64；跨平台代码请用显式 int64 |
| float32 / float64 | F / D | NaN/±Inf 直通 |
| bool | Z | 非 0 即 true |
| kernel.Value | Ljava/lang/Object; | 不透明句柄透传，身份保持 |
| error（末位返回） | — | 非 nil ⇒ RuntimeException(err.Error()) |

不支持（注册期 fail-fast）：struct、slice/map/chan、变参、多非 error 返回值。
扩展须新 ADR。

## 错误传播

```
error 返回      → java.lang.RuntimeException(err.Error())
ctx.Throw(...)  → 任意 bootstrap 异常（DefineClass 逃生舱）
Go panic        → 引擎错误（fatal）——bug 语义与内核一致
null 标量参数   → RuntimeException("interop: null argument N …")（Value 参数允许 null）
```

## 线程契约（DEV-0009）

- 绑定函数同步运行于调用方 goroutine；
- 函数内阻塞**不可被 Thread.interrupt 唤醒**——interrupt 只覆盖内核阻塞
  原语（sleep/wait/socket）；
- 长任务模式：embedder 自行启动 goroutine 并立即返回句柄/零值，完成后再
  回投 Java（经 ctx.Invoke）——v1 仅文档化，不提供调度器。

## 嵌入清单（checklist）

1. `kernel.New` + `vm.New`（执行桥）；
2. `interop.Bind`（必须在首次分派前，与 `gen.Install` 同窗口）；
3. loader + `SetResolver`（如需 classpath 规则类）；
4. 解析规则入口方法并 `th.Call`；
5. stdout/stderr 经 Options 注入——Java 输出即宿主输出。

## 已知限制汇总

- 类型集边界见映射表（struct/slice 等以 Value 句柄模式手工处理）；
- 阻塞不可中断（DEV-0009）；
- null 标量参数抛错而非转零值（显式优于静默）；
- 绑定函数内不得长期持有 CallContext（池化复用）。
