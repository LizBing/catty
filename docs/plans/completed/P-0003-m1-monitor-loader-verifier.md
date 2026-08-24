# P-0003: M1 — Monitor 重写 + 类加载器 + 字节码验证器

- 状态：completed（里程碑验收通过，2026-08-24 归档）
- 风险评分：Medium-High（并发原语与安全边界）
- 前置：P-0001 完成（M0）

## Goal

补齐 M0 明确欠下的三个正确性缺口（DEV-0002/0003、DEBT-0005..0006），使内核具备
多线程语义的地基和不可信输入的第一道防线。

## Current State

- M0 完成：解释器跑通真实 javac 产物；监视器为不可重入 sync.Mutex 占位；
  无验证器；无 classpath 概念（单文件显式加载）。

## Target State

1. **Monitor 重写**：可重入计数、FIFO 唤醒的锁竞争队列、`Object.wait(timeout)/
   notify/notifyAll` 完整语义（释放全部深度再重取）、线程身份经 `OwnerKey`
   贯穿 vm→kernel；中断钩子留接口（M2 接 Thread.interrupt）。
2. **类加载器**：目录 classpath + 懒式超类解析（循环守卫）；`catty run -cp`。
3. **验证器**：StackMapTable 解析 + 数据流检查（类型类别/栈深/合并点），
   默认对 CF 类开启；引用类型归并的保守规则登记 ledger。

## Tasks（DAG）

```
T1 Monitor 重写（含 Object.wait/notify natives、vm monitorenter/exit 改造）
T2 类加载器（classpath 目录 + 懒超类 + CLI -cp）
T3 StackMapTable 解析（classfile 层）
   ↑ 无依赖，可与 T1 并行
T4 验证器数据流检查 + 集成 + 负面测试
   ↑ T3
T5 多类 fixture 验收（Main↔Helper 经 -cp 加载，含 synchronized/wait-notify
   单线程可测部分）+ 文档回写
   ↑ T1 + T2 (+T4)
并行轨道：P-0002 微基准（Researcher 后台执行中，见 docs/research/R-0001）
```

## Validation

- `make check` 绿（vet/build/test 全仓）
- monitor 并发单测在 `-race` 下通过（Orchestrator 本地额外跑一次）
- 被篡改字节码的负面测试被验证器拒绝
- 多类程序经 classpath 正常运行

## Risks

- Go 无 goroutine ID → 线程身份用 vm.Thread 显式 ID 贯穿（OwnerKey 接口），
  该重构波及 CallContext —— 一次性做对，避免 M2 再返工。
- 验证器归并点的引用类型子类型查询需要已加载类图 → 保守规则：
  名称相同或可证子类才兼容，未知名一律放行并计数（trusted-input 期妥协，
  DEV-0006 登记）。
- Jackson 验收从 M1 移至 M1-stretch/M2（依赖面远超预估：reflect/io/nio.charset），
  本计划以多类自编 Java 程序作为 M1 验收，Jackson 决策记入 brief。

## Progress

- [x] 计划建立（2026-08-22）
- [x] P-0002 已委派 Researcher（后台，进行中）
- [x] T1 Monitor 重写（24fc480）：可重入/FIFO 竞争队列/wait 全深度释放/
      中断感知等待（atomic outcome）；OwnerKey 身份贯穿；-race ×2 绿
- [x] T3 SM 解析 + T4 结构层验证器（37a356a）：SM 全帧型解码；
      分支/处理点帧存在性、池合法性、边界检查；负面测试×5
- [x] T2 类加载器（37a356a）：目录 classpath、懒式依赖、循环守卫、CLI -cp
- [x] T5 多类验收（app.Main/Greeter/Helper）：invokeinterface +
      instanceof/checkcast(接口) + synchronized 重入 + 跨类加载，
      输出与参考 JVM 逐字节一致；字段默认值语义修复随本验收发现
- [ ] Jackson 反序列化验收 → 移交 M1-stretch/M2 决策（brief 披露）

## Open Questions

- wait 的中断唤醒通道形态（channel per waiter vs 共享 broadcast）——T1 实现
  时由 Worker 定并在代码注释记录理由。
