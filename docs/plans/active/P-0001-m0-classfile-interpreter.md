# P-0001: M0 — 仓库工程化 + classfile 解析器 + 最小解释器（HelloWorld）

- 状态：active
- 优先级：最高（当前唯一主线）
- 风险评分：Medium（新代码库，无不可逆外部承诺）

## Goal

达成 M0 验收：解释器在 Catty 内核上跑通 `javac --release 8` 编译的 HelloWorld，以及一个包含集合操作与异常抛接的小程序。

## Current State

- 治理与文档骨架完成（2026-08-22 本轮）：协议保存、AGENTS.md、ADR-0001..0008、make check 门。
- 无任何 Go 代码。

## Target State

Go module + 可运行的 v52 解释器内核 + 测试 + 验收脚本；`make check` 覆盖 go vet/build。

## Tasks（DAG）

```
T1 go.mod + 目录约定(internal/…) + Makefile 扩展(go vet/build/test)
T2 classfile v52 解析器
   ↑ T1        （规格笔记落 docs/specifications/classfile-v52.md）
T3 运行时内核最小对象模型(类元数据/对象头部/监视器占位)
   ↑ T2 接口冻结后可并行实现
T4 解释器核心(帧/操作数栈/dispatch loop/异常表)
   ↑ T2 + T3
T5 最小 java.lang 子集(Object/String/StringBuilder/System/
   基本类型包装类/Throwable 层级骨架/<clinit> 状态机 v0)
   ↑ T3
T6 集成验收(HelloWorld + 集合/异常 demo) + 文档回写
   ↑ T4 + T5
```

可并行对：T2 ∥ T3（接口先行）；T5 在 T3 后即启动。

## Dependencies

无外部依赖。javac 用于重编译 fixtures；CI 若无 JDK 则使用仓库内预编译 .class fixture。

## Validation

- `make check` 绿；
- `go test ./...` 绿；
- 验收脚本运行两个 fixture 程序输出正确。

## Risks

- 常量池/属性长尾解析量超预期 → 只解析 M0 所需子集，未识别属性按 JVMS 跳过；
- monitor 过早做全 → v0 用 sync.Mutex 占位，wait/notify 推迟到 M1。

## Progress

- [x] 治理 bootstrap（2026-08-22）
- [ ] T1 … T6

## Open Questions

- 解释器内部表示是否直接采用统一 `*Object`（倾向是，见 ADR-0004 rollback 路径）——Worker 自行决定并记录；
- fixtures 是否纳入 CI 重编译校验。
