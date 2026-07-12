# catty 项目进度快照

**日期**: 2026-07-12
**最后提交**: `0a96373` — fix: clinit stack corruption + native method access flags
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
| **R1.2** | **bootstrap classpath 基础设施 + native stubs** | `a1b991b` / `78d92f8` / `0a96373` |

## 当前状态

- **解释器操作码覆盖**: ~145/201（缺 invokedynamic、jsr/ret、以及部分 java.base 深层类所需的 opcode）
- **10 e2e 夹具**，三引擎验证（IR executor 的 invokeinterface 已知问题，run.sh 已容忍）
- **go vet clean**，所有 Go 测试通过
- **13 个 ADR**（0001-0007 Accepted，0008-0013 Proposed）
- **AOT 发射器覆盖**: 全部解释器支持的操作码

## R1.2 最终成果

### Native 基础设施
- **Native stub 机制** — 零值 stub，不再 panic
- **全局 native 注册表** — RegisterNative/GetNative，(class, method, descriptor) 索引
- **Classloader native 解析** — resolveNativeMethods

### 合成类及其 native 方法
- **Object**: hashCode, getClass, clone, equals, toString, notify, notifyAll, wait, registerNatives
- **String**: <init> (x2), length, charAt, intern
- **StringBuilder**: <init>, append(String/I/J/Z/C), toString
- **Class**: getName, getSimpleName, isInterface, isArray, getModifiers, desiredAssertionStatus, isInstance, isAssignableFrom, getSuperclass, isHidden, getPrimitiveClass, registerNatives
- **System**: out/err 静态字段; arraycopy, currentTimeMillis, nanoTime, identityHashCode, mapLibraryName
- **Thread**: <init>, currentThread, holdsLock, registerNatives
- **PrintStream**: println (String/I/J/Z/C/void), print (String/I)
- **Throwable 层级**: Throwable, Exception, RuntimeException, Error, LinkageError, IncompatibleClassChangeError, NoSuchMethodError, NPE, ArithmeticException, AIOOBE, CCE, IAE
- **Float/Double**: floatToRawIntBits, intBitsToFloat, doubleToRawLongBits, longBitsToDouble
- **Runtime**: availableProcessors, freeMemory, totalMemory, maxMemory, gc
- **AccessController**: doPrivileged (x2, stub)

### Bug 修复
- **clinit 栈腐败** — runClinit 不再劫持调用方帧
- **NativeMethod accessFlags** — 默认 instance + SetStatic()，避免静态方法弹出多余 this
- **InterpretedMethod maxLocals** — 对 native 方法补足最小局部变量槽
- **invokevirtual/invokeinterface nil 检查** — 方法查找失败时 throw NoSuchMethodError，不再 crash

## R1 阶段状态总结

ROADMAP 定义的 R1 子任务：

| 子任务 | 状态 | 备注 |
|---|---|---|
| Exceptions (try/catch/athrow) | ✅ | R1.0 |
| 剩余 opcodes (invokeinterface, wide, multianewarray) | ✅ | R1.1; invokedynamic 延后至 R3 |
| Bootstrap classpath (java.base, delegation, JRT/jimage) | 🔶 | 可加载 java.base 类，ArrayList/Math/Integer 等已验证；JRT 直接读取 jimage 尚未实现（当前通过 jimage extract 导出文件） |
| Native layer expansion (~20 类) | 🔶 | 合成类约 18 个，native 注册表覆盖 ~25 方法；ROADMAP 要求的"与真实 java.base 的 System.out 打通"已实现 |

## 战略路线图 (R1-R6)

- ✅ R1.0 异常处理
- ✅ R1.1 剩余操作码
- ✅ R1.2 bootstrap classpath + native stubs（核心目标已达成）
- ⬜ R1.3 原生层扩展（剩余的 bootstrap 类补齐 + JRT 直读）
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
