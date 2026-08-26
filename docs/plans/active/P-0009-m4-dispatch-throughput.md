# P-0009: M4 — 分发层吞吐优化与三方基准报告

- 状态：active（T0–T6 首轮完成；SB 折叠/装箱优化与 p99 采样待下轮）
- 风险评分：High（性能主张必须有 benchmark 证据，协议 §26；噪声范围内
  变化不得宣称优化成功）
- 决策纪律：内联缓存涉及分派语义 → 设计前查 ADR-0010 §2 与 emitter-abi §5

## Goal

在保持双路径逐字节一致的前提下，按 R-0005 已排序的 ROI 清单削减运行时
分发成本，产出第二批吞吐对照数据，兑现卖点②③⑤的首批证据。

## Tasks（DAG）— 首轮完成情况

```
T0 基线固化        重跑 R-0005 四基准（口径不变，保可比）                 ✅
T0.5 分发链剖析    cattaot CATTY_PPROF 门控；vcall 热循环 55% 在分发链    ✅
T1 分发链削减      单态内联缓存(发射期烘焙槽位)+FrameMeter 免锁计量+
                   EmitBody 免分配直调                                    ✅
                   连带根修 6 个潜伏 bug：getfield 模拟效应漏 -1、
                   dreturn 误归 cat1、线性深度模拟跨 goto 泄漏(重写为
                   CFG 工作表传播)、字段/数组访问 cat2 纪律、验证器
                   cmp2 popCat2、类加载器并发重复加载误判循环；
                   懒安装钩子修复嵌套/依赖类静默走解释器的历史问题
T2 natives 直调    被 #1 吸收：IC 后解析成本已消，瓶颈移至分配侧          ✅(降级)
                   （pprof 证据与决策记录见 R-0006 §决策记录）
T3 复测 + R-0006   vcall 1.10×→2.95×；mapops 1.08×→1.50×（压线达标）      ✅
T4 并发首证        ConcurrentBench：wall 不随线程数劣化（goroutine 底座    ✅
                   承载成立）；连带抓出并根修加载器竞态等并发缺陷
T5 可观测性        CATTY_PPROF 常驻能力；-race 构建 = Java 数据竞争检测器  ✅
                   （解释/AOT 双路径触发；DEV-0007 覆盖面限制已登记）
T6 捎带            make fuzz（DEBT-0001 缓解：30s 401 万次无 panic）；     ✅
                   叶帧行号未做 → 移交 U3
```

## Success Criteria 达成情况

- [x] make check 绿 / 全仓 -race 绿 / 钉扎零回退
- [x] vcall ≥1.5×（实测 **2.95×**）；mapops ≥1.5×（实测 **1.50×**，压线）
- [x] R-0006 落盘（docs/research/R-0006-dispatch-throughput.md）
- [x] U1 达标：mapops 再降 ≥20%（实测 -39.7%，R-0007）

## 下轮余项

```
U1 分配侧组合拳      JString hashCode/UTF-8 惰性缓存、jkey 结构键、          ✅
                     CallContext 池化、SB 折叠(增长链守卫)、abuf 参数缓冲
                     → mapops -39.7%（2.49× 于解释器，R-0007）
U2 strcat 波动归因   O(n²) 复制流量 × GC 相位；折叠 23× 回归反向印证        ✅(记录)
U3 堆栈叶帧行号       SetLine 行段发射+调用点捕获+源文件渲染；                ✅
                     TraceProbe 双路径与 JVM oracle 逐字节一致
U7 反射最小面         Class.forName/getDeclared*/Field/Method/Constructor     ✅
                     （P-0011 单列轮次完成，迷你序列化器=JVM 逐字节）
                     v2: getFields/getMethods/getConstructors 继承遍历        ✅
U8 JAR 加载           .jar 条目索引+流式读取（DEBT-0008 关闭）                 ✅
U9 装箱快路径         26ns vs 57-67ns（-55%，微观基准 BoxedInt*）             ✅
U10 gson 战役         P1-P3 全通（forName/toJson/fromJson），P5 注解边界      ✅
                     P6 泛型边界；R-0009 边界报告落盘                        
U4 p99 持续负载采样   P99Bench 单线程分位 + embeddemo 嵌入服务形态           ✅
                     （4/8/16 worker 吞吐 43k→65k req/s，p99<1ms，
                      R-0007 附录②）
U5 Integer 装箱削减   需逃逸分析级手段或 workload 层接受                      ⬜
U6 float/double 算术组对调修复 + Math.floor/ceil/sqrt/rint 面                 ✅
   （P99 fixture 首踩浮点 AOT 路径暴露）
```

## Risks（存档）

- 内联缓存失效语义已按 genrt resolved 缓存教训处理（InstallKernel 全量失效化，
  TestInlineCacheHitAndReinstall 钉扎）
- spike 无负结果：两项均正向兑现
