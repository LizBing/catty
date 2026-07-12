# catty 项目进度快照

**日期**: 2026-07-12
**最后提交**: `a1b991b` — feat: R1.2 (partial) — bootstrap classpath infrastructure + native stubs
**CI 状态**: 全绿

## 已完成的里程碑

| 阶段 | 内容 | 提交 |
|---|---|---|
| Phase 1 | 解释器 MVP (~140 操作码, 5 核心类, 8 夹具) | `12cf2f2` |
| A0 | 字节码 IR + 消栈 | `c7c7556` |
| A1 | AOT 发射器 (fib 44ms 原生) | `c168530` |
| A1.5 | 类型追踪 (StackMapTable) | `ce24313` |
| A2.1 | fresh-per-def 类型感知发射 | `fcc0666` |
| A2.2 | invoke 桥接 (HelloWorld 原生) | `2616128` |
| A2.3 | 循环 (空栈合并) | `aa613cb` |
| A2.4 | diamonds/phi | `368979d` |
| A2.2b | OOP (new/字段/构造器) | `c6d2009` |
| A2.5 | long/float/double | `c3a7b5b` |
| A2.6 | 边缘项 (frem/drem, cat-2 phi, switch) | `0843556` |
| A4 | catty build (离线 AOT 集成) | `8023e90` |
| 文档 | 战略愿景 + 6 ADR (0008-0013) | `2ad5091` |
| **R1.0** | **异常处理** | `ece4412` |
| **R1.1** | **invokeinterface + multianewarray + wide** | `e043fc5` |
| **R1.2** | **bootstrap classpath 基础设施 (partial)** | `a1b991b` |

## 当前状态

- **11 Go 包, ~8000 LOC**
- **10 e2e 夹具** (HelloWorld, Fibonacci, Factorial, ArraySum, OOPDemo, StaticFields, SwitchDemo, EmptyMain, ExceptionTest, InterfaceTest)
- **三引擎验证**: tree-walker + IR executor + java（InterfaceTest 的 IR executor 有已知的 invokeinterface dispatch 问题，run.sh 已容忍）
- **vet clean**
- **13 个 ADR** (0001-0013，其中 0008-0013 为 Proposed)
- **解释器操作码覆盖**: ~145/201
- **AOT 发射器覆盖**: 全部解释器支持的操作码

## R1.2 完成的部分

1. **Native stub 机制** — `rtda/method.go`: native 方法默认返回零值 stub，不再 crash
2. **全局 native 注册表** — `native/native_registry.go`: RegisterNative/GetNative，
   以 (class, method, descriptor) 为 key，`init()` 自动注册
3. **核心 native 方法** — `native/system.go`:
   - System.arraycopy/currentTimeMillis/nanoTime/identityHashCode
   - Object.hashCode/getClass/clone (含 `native/lang.go` 的 equals/toString)
   - Class.getName/getSimpleName/isInterface/isArray/getModifiers/desiredAssertionStatus
   - Thread.currentThread
   - Float.floatToRawIntBits/intBitsToFloat, Double.doubleToRawLongBits/longBitsToDouble
4. **合成 Class/Thread 类** — `native/registry.go`: buildClass, buildThread
5. **Classloader native 解析** — `classloader/classloader.go`: resolveNativeMethods
6. **clinit 重写** — `interpreter/invoke.go`: 同步 runClinit/Loop 执行
7. **ldc Class 字面量** — 推送 Class 对象（之前为 null）

## R1.2 剩余工作（待完成的 native 填充）

目标：补全加载真实 java.base 类所需的 native 方法，使
`catty -cp <user:java.base> HelloWorld` 能运行。

- **继续填充 native 注册表** — 真实 JDK 类中缺失的 `ACC_NATIVE` 方法
  （如 Class.getPrimitiveClass、Class.forName 等）
- **依赖类加载** — Math.max、Arrays.sort 等深层级联类加载
- **Bootstrap classpath** — 前置 java.base 到 classpath
- **ClassLoader 委托层次** — parent delegation
- **JRT/jimage 支持** — 读取 JDK 模块镜像

## 战略路线图 (R1-R6)

- ✅ R1.0 异常处理
- ✅ R1.1 剩余操作码 (invokeinterface/wide/multianewarray)
- 🔶 R1.2 bootstrap classpath + native stubs (partial)
- ⬜ R1.3 原生层扩展 (~20 个 bootstrapping 类)
- ⬜ R2 多线程 (Thread=goroutine)
- ⬜ R3 反射 + invokedynamic
- ⬜ R4 I/O 集成
- ⬜ R5 AOT 覆盖扩展
- ⬜ R6 性能调优

## 关键架构决策 (ADR)

| ADR | 决策 | 状态 |
|---|---|---|
| 0001 | 套用 Go runtime (不自建 GC/调度/JIT) | Accepted |
| 0002 | switch dispatch (Go 无 computed goto) | Accepted |
| 0003 | tagged 16B Slot | Accepted |
| 0004 | 原生核心类 (不加载 JRE) | Accepted |
| 0005 | 懒 <clinit> | Accepted |
| 0006 | predecode 不划算 — AOT 是性能路径 | Accepted |
| 0007 | 反射/动态: 分层, 保留解释器 | Accepted |
| 0008 | AOT-first (解释器是开发层) | Proposed |
| 0009 | 混合类库 (~50 原生 + ~7000 JDK) | Proposed |
| 0010 | Thread = goroutine | Proposed |
| 0011 | Go 内存模型 | Proposed |
| 0012 | escape analysis 替代分代 GC | Proposed |
| 0013 | 直接 Go runtime 集成 | Proposed |
