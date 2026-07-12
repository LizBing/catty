# catty 项目进度快照

**日期**: 2026-07-12
**最后提交**: `ece4412` — feat: R1.0 — exception handling in the interpreter
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

## 当前状态

- **11 Go 包, ~8000 LOC**
- **9 e2e 夹具全绿** (三引擎: java / tree-walker / IR executor)
- **vet clean**
- **13 个 ADR** (0001-0013)
- **解释器操作码覆盖**: ~140/201 (缺 invokeinterface, invokedynamic, wide, multianewarray)
- **AOT 发射器覆盖**: 全部解释器支持的操作码

## 下一步: R1.1 — 剩余操作码

已规划（在 plan file 中），尚未执行。核心内容:
1. **invokeinterface** — 镜像 invokevirtual + 跳过 2 操作数字节。最高杠杆率。
2. **wide** — u2 索引扩展。简单。
3. **multianewarray** — 多维数组递归创建。~15 行。
4. lowering 的 decodeInst 支持这三个操作码的解码（长度 5/4-6/4）。

里程碑: InterfaceTest (Comparable + sort) 三引擎一致。

## 战略路线图 (R1-R6)

- ✅ R1.0 异常处理
- ⬜ R1.1 剩余操作码 (invokeinterface/wide/multianewarray)
- ⬜ R1.2 bootstrap classpath (加载真实 java.base)
- ⬜ R1.3 原生层扩展 (~20 个 bootstrapping 类)
- ⬜ R2 多线程 (Thread=goroutine)
- ⬜ R3 反射 + invokedynamic
- ⬜ R4 I/O 集成
- ⬜ R5 AOT 覆盖扩展
- ⬜ R6 性能调优

## 关键架构决策 (ADR)

| ADR | 决策 |
|---|---|
| 0001 | 套用 Go runtime (不自建 GC/调度/JIT) |
| 0002 | switch dispatch (Go 无 computed goto) |
| 0003 | tagged 16B Slot |
| 0004 | 原生核心类 (不加载 JRE) |
| 0005 | 懒 <clinit> |
| 0006 | predecode 不划算 — AOT 是性能路径 |
| 0007 | 反射/动态: 分层, 保留解释器 |
| 0008 | AOT-first (解释器是开发层) |
| 0009 | 混合类库 (~50 原生 + ~7000 JDK) |
| 0010 | Thread = goroutine |
| 0011 | Go 内存模型 |
| 0012 | escape analysis 替代分代 GC |
| 0013 | 直接 Go runtime 集成 |
