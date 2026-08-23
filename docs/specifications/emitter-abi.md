# 发射器调用约定（Emitter ABI）v0.1 草案

- 状态：Draft — 实现启动时经 T1 评审转 Accepted
- 依据：R-0002 异常机制基准（F 方案定稿输入）、ADR-0004（统一表示）、判断一（v1 统一 Instance）
- 本文档定义 AOT 发射器生成的 Go 函数与运行时内核之间的全部接口约定。
  变更须经 ADR。

## 1. 命名与布局

- 每个被发射的 Java 类 → 一个独立 Go 文件 `gen/<internal_name>.go`，
  包名 `gen`；文件头注明 generated 标记（R6：禁手改）。
- 方法函数名：`Catty_<class_mangled>_<name>__<desc_mangled>`
  - class/name 内部形式中 `/`→`_`、`$`→`_00036`（嵌套类）；
  - desc_mangled：描述符去括号，`L...;`→`L..._`、`[`→`A`、原始类型字母保留。
  - 示例：`java/util/ArrayList.add(Ljava/lang/Object;)Z`
    → `Catty_java_util_ArrayList_add__Ljava_lang_Object_Z`
- 静态方法与实例方法同为包级函数（无 receiver 魔法），显式首参。

## 2. 签名形态（异常通道 = F 方案，ADR-0009）

```go
// void 静态方法
func Catty_Foo_bar__I_V(n int32) *kernel.Thrown
// 非 void 实例方法（recv 为具体结构或 *kernel.Instance，见 §3）
func Catty_Foo_get__I_J(recv *Catty_Foo, idx int32) (int64, *kernel.Thrown)
```

- 返回二元组 `(result, exc *kernel.Thrown)`；void 方法仅 `(exc)`；
- cat2 参数占一个逻辑参数位（int64/float64 天然表达）；
- **每个调用点必须检查 exc != nil 并立即转发**——发射器保证无遗漏
  （机械检查项：生成代码评审清单 + 未来 lint）；
- catch 块发射为对 exc 的类型判断链（无函数拆分、无 recover）；
- finally 发射为路径复制（返回前内联），不做 defer。

## 3. 对象与字段

- 一切引用值为 `kernel.Instance` / `*JString` / `*ArrayObj`（统一表示，判断一）；
- 字段读写经 kernel 辅助函数（`GetField/SetField` 内联热路径可后续特化）；
- v1 不生成强类型包装结构——混合表示留待发射器 v2 重评（ADR-0004 收口条件不变）。

## 4. 运行时映射

| Java 构造 | 发射形态 |
|---|---|
| `<clinit>` | 普通函数，由 kernel EnsureInitialized 经反射表调用 |
| synchronized 方法 | 序言 Enter / 全部出口 Exit（显式，无 defer；含 exc 转发路径） |
| synchronized 块 | monitorenter/exit 原语直呼 |
| SOE guard | 序言深度计数++、全出口 --；超阈值抛 StackOverflowError（R-0002 发现③） |
| throw | `return nil, &kernel.Thrown{Obj: ...}` |
| catch | exc 类型判断链（IsInstance）|
| String/数组字面量 | 包级惰性初始化变量（首次使用构建，等价 ldc 语义） |

## 5. 与解释器的互操作

- 同一对象空间：AOT 函数可持有/操作解释器创建的 Instance，反之亦然；
- 反向调用（AOT→解释器）：未发射方法经 kernel.Invoker 走解释路径，
  调用点由发射器生成"查表未命中即回退"的桩；
- 正向调用（解释器→AOT）：类元数据的方法表项指向 AOT 函数包装器。

## 6. 已知限制（v0.1）

- 无 IR：优化限于发射期局部变换；
- 不支持 jsr/ret、invokedynamic（后者需构建期脱糖产物）；
- socket 读中断唤醒未接（DEV-0011），发射代码继承该限制。
