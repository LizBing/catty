# 发射器调用约定（Emitter ABI）v0.2

- 状态：Accepted — 与实现同步（internal/gen 39 类生产运行）
- 依据：ADR-0009（异常通道 F 方案）、P-0006 判断一（v1 统一 Instance 表示，ADR-0010 收口）
- 本文档定义 AOT 发射器生成的 Go 函数与运行时内核之间的全部接口约定。
  变更须经 ADR。

## 1. 命名与布局

- v0.2 实际布局：单一 `internal/gen/gen.go`（`genemit` 产物，DO-NOT-EDIT 头），
  跨类引用与注册顺序最简；分包布局留待生成量阈值触发（P-0006 风险项）。
- 方法函数名：`Catty_<class_mangled>_<name_mangled>__<desc_mangled>`
  - `/`→`_`、`$`→`00036`（嵌套类转义）、`.`（包段）→`_`；
  - 示例：`Json.<init>()V` → `Catty_com_eclipsesource_json_Json_init____V`。

## 2. 签名形态（异常通道 = F 方案旗标返回，ADR-0009）

```go
func Catty_<mangled>(thr kernel.OwnerKey, recv kernel.Value,
    args []kernel.Value) (kernel.Value, *kernel.Thrown)
```

- 一切方法统一形态：静态方法 recv=nil；cat2 参数占一个 args 槽；
- 注册：`installTable`（cls/name/desc/fn）+ `gen.Install(k)` 将发射体覆写为
  `m.Native`——未覆盖类自动落回解释路径（混合执行零成本）；
- **每个调用点必须检查 exc != nil 并立即转发**；catch 发射为 exc 类型判断链，
  finally 为路径复制，无 recover / defer。

## 3. 对象与字段

- 一切引用值为 `kernel.Instance` / `*JString` / `*ArrayObj`（统一表示）；
- 字段经 `genrt.GetFieldChecked/SetFieldChecked`（含 null-receiver NPE 语义）。

## 4. 运行时映射

| Java 构造 | 发射形态 |
|---|---|
| `<clinit>` | 不发射，走解释器初始化机制 |
| invokes | `genrt.Call{Virtual,Special,Static,Interface}(thr, …)` → `kernel.InvokeAs` |
| 数组/算术等可抛原语 | 带 `thr` 的 Checked 助手；任何助手抛出都回填 Java 栈（§7） |
| throw | `&kernel.Thrown{Obj: …}` 经 exc 通道转发 |
| SOE guard | `genrt.invokeChecked` 经 `ThreadRegistry.FrameEnter/Exit` 计量 |

**栈深效果表 = 单一事实源**：`classfile.StackEffect(op)` 对全部 256 opcode
给出总分类（EffectFixed / EffectDescriptor / EffectIllegal），发射器深度模拟、
哨兵测试与未来引擎一律取自该表。私有效果表是 DEBT-0015/0017/0019 三起事故的
共同根因，禁止再建。

## 5. 与解释器的互操作

- 同一对象空间，双向混跑；分派唯一入口 `kernel.InvokeAs`
  （同时承担 §7 的帧追踪）。
- genrt 层禁止跨内核缓存存活：`InstallKernel` 时失效化全部解析缓存。

## 6. 堆栈回填（stack backfill，DEBT-0019 诊断基建）

语义对齐 HotSpot fillInStackTrace：

1. `kernel.InvokeAs(owner,…)` 在 owner 实现 `kernel.FrameTracker` 时推入/
   弹出 `JavaFrame{Class, Method}`——解释器与发射体共用同一栈；
2. 捕获点：Throwable 构造 natives、`vm.throwNamed`、全部 genrt 可抛助手；
3. 快照裁剪栈顶连续 `<init>` 构造链（用户异常不泄漏构造帧）；
4. 渲染：`kernel.FormatUncaught` 按 JVM 叶帧在前输出；
5. 已知限制：行号恒为 `(Unknown Source)`——逐帧调用点 pc 未追踪
   （需每指令行号登记或返回地址采样，v2 议题）。

## 7. 已知限制（v0.2）

- 无 IR：优化限于发射期局部变换；
- 不支持 jsr/ret（项目规则禁用）、invokedynamic（构建期脱糖产物除外）;
- 堆栈回填无行号（§6.5）；
- ldc 类常量仅支持 String/原始类型，类对象经 ldc-CClass 描述符透传。
