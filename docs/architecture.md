# Catty 架构基线

- 状态：Accepted（基线 v0.1）
- 日期：2026-08-22
- 地位：本文是架构快照；任何结构性变更必须先产生新 ADR，再回写本文。
- 关联 ADR：ADR-0002（JDK 8 基线）、ADR-0003（执行引擎）、ADR-0004（对象模型·暂定）、ADR-0005（异常机制延后）、ADR-0006（类库策略）、ADR-0007（JNI 范围）

## 1. 分层

```
┌────────────────────────────────────────────────┐
│ 应用 (javac --release 8 产物 .class / jar)      │
├────────────────────────────────────────────────┤
│ AOT 编译器: classfile → IR → 优化 → Go 源码发射   │ ← 构建期
│            (closed-world 可达性分析)             │
├────────────────────────────────────────────────┤
│ 运行时内核 (纯 Go):                              │
│  对象模型 · 监视器/锁 · 引用队列 · 类元数据        │
│  字节码解释器(动态类加载/反射兜底) · JNI 桥(可选)  │
│  java.* mini profile 类库                       │
├────────────────────────────────────────────────┤
│ Go 运行时: GC · 调度器 · netpoller · 连续栈      │
└────────────────────────────────────────────────┘
```

## 2. 设计原则

1. **寄生战略**：凡是 Go 运行时已做对的事绝不重做——不自建 GC、safepoint、调度器、netpoller。工程火力集中于三件事：对象模型语义保真、AOT 翻译质量、java.* 类库。
2. **解释器兜底**：运行期保留完整解释器。未 AOT 的动态加载类、未知反射路径落到解释器——慢但正确，而不是构建期失败（对 native-image 的核心差异化）。
3. **偏差必须登记**：任何与 JVMS/JDK8 规范行为的偏差记入 `docs/specifications/deviation-ledger.md`。

## 3. JRE 语义 ↔ Go 运行时映射

| JRE 语义 | Go 实现 | 备注 |
|---|---|---|
| Java 线程 | goroutine | 平台/虚拟线程统一映射，差异退化为调度提示 |
| GC / safepoint | Go 并发三色标记清除 | 无需自建 safepoint 轮询 |
| `-Xmx` | `GOMEMLIMIT` | **软限制**，行为差异登记 deviation-ledger |
| 调用栈 | goroutine 连续栈 | SOE 需自建深度计量（§7） |
| `synchronized` | 自研监视器（薄锁 + 自定义等待队列） | 不能用 sync.Mutex：需 wait/notify/interrupt |
| `volatile` | `atomic` 操作 | Go 原子操作顺序一致 ≥ Java volatile 要求 |
| 异常 | 待定（panic/recover vs 显式旗标），见 ADR-0005 | 必须用哨兵类型与普通 panic 区分 |
| 弱/软/虚引用 | AddCleanup/SetFinalizer + 自研 ReferenceQueue | 行为偏差登记 |
| 直接内存 | mmap arena + cleanup 登记 | 模拟 MaxDirectMemorySize |
| IO 阻塞 | 可中断封装层（硬规则） | 例：interrupt 触发 `SetReadDeadline(now)` |

## 4. 对象模型（ADR-0004，暂定待证据）

- Java 类 → Go 结构体，embedding 镜像继承链；公共头部 `{cls, monitor(惰性), hashCode(惰性)}`；
- 身份 = 指针相等；数组 = 定长结构体包裹保持稳定身份；
- 引用静态类型不足时提升为 Go interface（itab 分派 ≈ 虚调用）；CHA 收窄后为直接调用；
- 字段遮蔽用限定选择器显式表达；
- String（JDK 8 语义）：单后端 `[]uint16` 起步；升级现代 JDK 时引入 Latin-1 双后端。
- **该决策标记为 provisional**：裁决实验见 docs/plans/active/P-0002，必须在 M3 发射器动工前完成。

## 5. 执行引擎（ADR-0003）

双模式混合，同一对象模型保证互操作性：

```
javac --release 8 → classfile 解析 → IR(版本无关) → 构建期优化 → Go 源码发射 → go build → 静态二进制
                                                          ↘ 未覆盖路径 → 解释器执行
```

- **AOT 主路径**：closed-world 可达性分析；lambda 构建期脱糖（JDK 8 无运行期 indy 需求）；字符串拼接天然是 StringBuilder 字节码。
- **构建期优化清单（按 ROI）**：CHA 去虚化与直接调用 → 内联 → lambda 脱糖 → 常量池固化 → intrinsic（字符串拼接/math）→ 冗余 null check 删除 → 非逃逸 monitor 消除。借力 Go 编译器的逃逸分析与 BCE（把代码塑形成 Go 能证明的形状）。
- **解释器**：v0 不优化，正确性优先；同时服务 M0 验收与长期兜底职责。

## 6. 内存管理与引用语义

- `GOMEMLIMIT` 是软限制：逼近时疯狂 GC 而非抛 OOM；在 GC 钩子里尽力模拟 OutOfMemoryError，差异登记；
- Go GC 不压缩：长命大对象碎片风险由 AOT 逃逸分析缓解（短命对象不上堆）;
- 四种 Reference 的可达性层级与通知顺序手工实现；软引用清除策略以 GOMEMLIMIT 压力信号近似 JDK 时钟-LRU。

## 7. 并发模型三处硬细节

1. **中断纪律（硬规则）**：所有阻塞必须过中断感知封装；goroutine 卡死在不可中断原语上 = 线程模型破洞。
2. **StackOverflowError**：Go 栈溢出是不可恢复 fatal error；发射方法序言插入低频深度计量（每 N 帧），近似 `-Xss`。
3. **`<clinit>` 状态机**：按 JVMS §5.5 自实现（循环初始化处理、ExceptionInInitializerError 传播）；禁用 sync.Once（会静默死锁）。

## 8. 类库三层策略（ADR-0006）

| 层 | 范围 | 来源 | 许可纪律 |
|---|---|---|---|
| L1 内核紧耦合类 | String、包装类、集合骨架、Thread/Object、异常层级、引用类 | clean-room 自写 | 主仓 Apache-2.0 |
| L2 机械大体量类 | BigInteger、regex、java.time 等 | Apache 血统开源移植（Harmony attic / AOSP libcore），落地前逐文件许可甄别 | PROVENANCE 登记 |
| L3 符合性 | 行为正确性证明 | jtreg 子集等 GPL 测试，放独立外部仓，永不进主仓 | 物理隔离 |

v1 目标子集："服务端 mini profile"——够跑一个 JSON HTTP 微服务（java.lang 核心 + java.util + io/nio + net + concurrent + time）。nio/net 大量委托 Go 标准库（共享 netpoller）。

## 9. JNI 边界（ADR-0007）

- v1 范围内、M2.5 并行轨道启动；`jnion` 构建 tag 控制，默认关，静态链接故事永不被 cgo 破坏；
- JNIEnv = 真实 C ABI 函数表（unsafe 布局），函数体指向 Go 实现，回调分发回内核（AOT thunk 或解释器）;
- local/global reference 用句柄表（索引而非裸指针）：符合 JNI 语义并绕开 cgo 指针规则；Go GC 不移动对象，global ref 底层可直接持裸指针；
- AttachCurrentThread：LockOSThread 专用 goroutine 等待任务派发；pending exception 协议映射为 Env 挂起字段；
- v1 裁剪：不支持 GetPrimitiveArrayCritical 等高危接口。

## 10. 可观测性方向

pprof 剖析 Java 线程；`-race` 构建 = 数据竞争检测器；两者作为卖点在 M2 后接入文档与 demo。

## 11. 版本演进路径

当前基线 JDK 8（class file v52）。升级到 17/21 为增量改造（预估 1–2 月，前提：IR 保持版本无关）：indy 运行期链接或库 jar 构建期脱糖、condy 构建期常量折叠、nestmates 访问检查、record 属性解析、String 双后端、虚拟线程 Thread API 皮。

## 12. 未决问题

| 问题 | 追踪位置 |
|---|---|
| 对象模型表示法证据裁决 | plans/active/P-0002 |
| 异常机制（panic vs 旗标） | ADR-0005，M1 末截止 |
| libcore/Harmony 逐文件许可甄别 | ADR-0006 + debt/register |
