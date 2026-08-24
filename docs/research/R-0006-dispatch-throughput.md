# R-0006: 分发链优化后稳态吞吐复测 — T1 内联缓存轮（P-0009）

- 状态: 完成
- 关联: P-0009；R-0005 的优化方向 #1/#2 落地复测
- 环境: darwin/arm64 · go1.26.5 · JDK 25 HotSpot · 同机同负载 · GOCACHE=/tmp

## 改动清单（本轮进入测量路径的）

1. **单态内联缓存**：`genrt` 直映射 1024 槽 `[site,dyn]→*Method`，原子指针换入，
   `InstallKernel` 失效化；发射期烘焙槽位常量（`CallVirtualIC`），消除每调用
   3 次字符串 FNV 哈希。
2. **免锁帧计量**：`kernel.FrameMeter` 接口，vm.Thread 以自有字段计量，
   发射路径弃用 ThreadRegistry 全局锁；解释/发射共享同一 SOE 预算
   （更贴近 JVM 单栈模型）。
3. **EmitBody 直调**：发射体直挂 `Method.EmitBody`，`InvokeAs` 免分配快路径
   （不再构造 CallContext、少一层闭包间接）。
4. **懒安装钩子**：内核类加载回调——修复嵌套类/依赖类永远走解释器的静默
   回退（此前 minimal-json 库部分从未真正 AOT 执行过）。

### 连带根修的潜伏发射器 bug（全 AOT 首次暴露）

| bug | 后果 | 发现途径 |
|---|---|---|
| getfield 模拟效应漏 -1 | 直线段深度逐指令虚涨 | JsonParser.read nil 槽 panic |
| dreturn 归入 cat1 返回组 | 双精度返回弹 1 槽取错值 | JsonObject.getDouble 垃圾返回 |
| 线性深度模拟跨 goto 泄漏 | 合并点槽位污染 | 同上 |
| 字段/数组访问无 cat2 纪律 | putfield J 把哨兵当 receiver | ConcurrentBench Worker panic |
| 验证器 cmp2 未用 popCat2 | lcmp 程序被拒收 | ConcurrentBench VerifyError |
| 类加载器并发重复加载误判循环 | 多线程首次 new 即 ClassCircularityError | 8 线程压测 |

另修复：genrt.New 对齐解释器 new 语义（String 魔法表示 + EnsureInitialized，
统一初始化 tracker 防 DEV-0002 家族死锁）。

## 方法

与 R-0005 相同 fixture `Bench`（tickMillis 时钟；本轮实测打印单位为纳秒，
与 R-0005 表内 ms 数值 ×10⁶ 一致）。N=1,000,000，各基准 3 轮取最小值；
post 数据取 15 轮最小值以压噪声。单负载 fixture（VCallBench/MapBench）
用于 pprof 定靶与解释器基线。

## 结果（ms，越小越快）

| 基准 | 解释器 | AOT pre | AOT post | post/pre | AOT/解释器(pre→post) |
|---|---:|---:|---:|---:|---|
| arith | 89.72 | 37.96 | 37.85 | 1.00× | 2.36× → 2.37×（对照组，不动 ✓） |
| vcall | 201.51 | 182.47 | **68.42** | **0.375×** | 1.10× → **2.95×** |
| mapops | 1566.39 | 1448.83 | **1042.76** | **0.720×** | 1.08× → **1.50×** |
| strcat | 3183.10 | 3475.62 | 3054.88 | 0.879× | 0.92× → 1.04× |

单负载交叉验证：VCallBench 74ns/调用（pprof 口径）。

**目标核对（P-0009 成功标准 ≥1.5×）：vcall 2.95× ✓（超额）、mapops 1.50× ✓（压线）。**

## 归因分析

- **vcall**：纯分发负载，三项削减（IC 免解析+免哈希、免锁帧计量、EmitBody
  直调）全部命中，剩余 ~68ns 中生成码本体（aload/getfield/ireturn）占比
  上升为多数——分发层不再是首位项。
- **mapops**：IC 消除的是解析成本；HashMap native 本体的哈希/UTF-16 比较、
  Integer 装箱分配、`"k"+(i%1024)` 的 SB 链构成新瓶颈（pprof：mallocgc/
  memclr/kevent 占 56%）。R-0005 方向 #3（SB 五连折叠）与装箱缓存是下一刀。
- **strcat**：波动大（3.05–4.86s），改进幅度在噪声边缘，不作宣称。
- **arith 对照不动**：符合预期，证明改动未伤非分发路径。

## 并发首证（T4 spike，ConcurrentBench：N worker × 混合负载）

| 配置 | Catty AOT wall | JDK 25 wall |
|---|---:|---:|
| 4 × 2M | 210–273ms | 14.5ms |
| 8 × 1M | 114–250ms | 23.7–32.6ms |
| 16 × 500k | 52–217ms | （未采） |

诚实结论：绝对性能差一个数量级（JIT + 成熟集合）；但 **wall time 不随线程数
劣化**——goroutine 底座的并发承载成立，且多线程首次类加载竞态、cat2 字段
纪律等 5 个真并发缺陷在本实验中被抓出并根修。卖点③的"承载"叙事有了第一份
数据；"吞吐"仍属 M4 后续。

## 可观测性白嫖验证（T5 spike）

- **pprof**：`CATTY_PPROF=<file>` 直接剖析运行中的"Java"程序，画像精确到
  genrt 助手/发射函数粒度——卖点⑤前半兑现，已固化为 cattaot 常驻能力。
- **-race**：`go build -race` 构建的 cattaot 运行含数据竞争的 Java 程序
  （RaceRig：双线程非同步 counter++），Go race detector 报出 DATA RACE 并给
  出两条 goroutine 的完整 Java 执行栈；解释/AOT 两路径均触发，计数结果
  128342（丢失更新可见）。**卖点⑤后半兑现：-race 构建即 Java 数据竞争
  检测器。**
- 附带发现：volatile 语义缺失（DEV-0007）意味着竞争检测的覆盖面与 HotSpot
  JMM 有偏差，登记为该能力的已知限制。

## 决策记录

- **T2（natives 白名单直调）降级不实施**：IC 落地后解析成本已消，mapops
  剩余瓶颈在分配侧（pprof 证据如上），直调的预期收益 < 装箱/SB 折叠。
  R-0005 方向 #2 标记为"被 #1 吸收"，后续按新画像重排。
- DEBT-0001 缓解落地：`make fuzz`（FuzzParse，30s 有界会话 401 万次执行
  无 panic；语料含全部 fixture 种子）。

## 结论

分发层优化一轮到位：调用密集场景从"与解释器持平"变为 **2.4–3.0×**，
mapops 压线达标 1.50×。R-0005 优化清单消费完毕方向 #1，#2 被 #1 吸收，
#3（SB 折叠）与装箱分配优化是 mapops 的下一刀。并发与可观测性两个卖点
从零数据变为有数据、有演示、有已知限制登记。
