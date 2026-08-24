# P-0004: M2 第一阶段 — 线程映射 + 中断语义 + SOE

- 状态：completed（里程碑验收通过，2026-08-24 归档）
- 风险评分：High（并发原语 × 语义保真；按 §28 需评审与人类可见检查点）

## Goal

把 goroutine 底座升级为可用的 java.lang.Thread：生命周期（start/join/isAlive）、
线程身份（currentThread）、中断语义（sleep/Object.wait 可被中断并抛
InterruptedException，标志位 peek/clear 两态），以及 StackOverflowError 的
深度计量。验收 = 多线程 fixture 输出与参考 JVM 一致。

## Current State

M1 完成：Monitor 具备完整所有权/等待/中断唤醒数据通路（waiter outcome 原子化，
InterruptWaiters 就绪但无人调用）；vm.Thread 有身份键但无 Java 对象绑定；
无 Thread 类；无 SOE。

## Target State

1. `kernel/thread.go`：JThread 记录（状态机/中断标志/join 门闩）+
   ThreadRegistry（key→记录、Object.wait 挂起登记、睡眠登记）。
2. Monitor.Wait 发布等待记录到 registry → Interrupt(key) 可跨监视器唤醒。
3. bootstrap `java/lang/Thread`（start/join/sleep/currentThread/interrupt/
   isAlive/getName/setName/isInterrupted/static interrupted）+ Runnable 标记 +
   VirtualMachineError/StackOverflowError 链。
4. vm：SpawnHook 安装（goroutine 起新 vm.Thread 绑定 Thread 对象、uncaught
   异常打印）、主线程注册、exec 深度计量 → StackOverflowError。
5. fixture ThreadsDemo（双线程计数 join / 中断 sleep / 守护输出）与真机对照。

## Tasks（DAG）

```
T1 kernel 线程注册表 + Monitor.Wait 发布改造
T2 Thread natives/bootstrap + SpawnHook(vm)
   ↑ T1
T3 SOE 深度计量（独立，可与 T1/T2 并行）
T4 fixture + 对照验收 + -race 全绿
   ↑ T2 (+T3)
后续任务（本计划外，M2 续）：DEBT-0009 数据流验证器、nio/net 映射与 http echo。
排序理由：验证器是独立的安全边界工程，需要整段专注预算（上一轮已证明赶工
不可行）；http echo 依赖 socket natives + charset 面。二者不阻塞线程主线。
```

## Validation

- make check 绿；`go test -race ./internal/...` 绿
- ThreadsDemo 与 `java -cp` 参考输出逐字节一致
- interrupt 三路径测试：睡着的 sleep / Object.wait 中的等待者 / 运行中标志位

## Risks

- goroutine 泄漏（join 未达/守护线程）→ 测试用带超时 join；daemon 语义 M2 后期
- SpawnHook 循环依赖 → kernel 定义钩子字段，vm 安装实现（同 Invoker 模式）
- 中断竞态（标志置位 vs 进入阻塞的窗口）→ 先查标志后挂载登记的顺序 +
  登记后复查标志的双检模式

## Progress

- [x] 计划建立（2026-08-23）
- [x] T1 注册表 + Wait 发布（registry.publishWait/retractWait）
- [x] T2 Thread natives/bootstrap + SpawnHook + 主线程注册
- [x] T3 SOE 深度计量（MaxFrames，默认 4096，解释帧计数）
- [x] T4 ThreadsDemo 验收：500/interrupted ok/main done main/w1 alive=false
      与参考 JVM 逐字节一致；`go test -race ./internal/...` 全绿
- 验证期修复的关键缺陷：①ctx.Invoke 丢失线程身份（sleep 变不可中断裸睡）
  ②k.Invoker 全局单槽被并发覆盖导致 <clinit> 重入检查失效自死锁
  （改为 owner 路由执行）③Class.State 无锁读（改 atomic）
- Jackson/http echo/DEBT-0009 数据流验证器 → P-0005（M2 续）

## Open Questions

- Thread.stop/deprecated 语义：不做（现代 JDK 已移除）。
- 守护线程对进程退出的影响：v1 不模拟 setDaemon 退出规则（登记 ledger 待办）。
