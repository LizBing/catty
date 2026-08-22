# Catty 能力状态

> 持续更新的能力矩阵。每次里程碑收口时刷新；与 `deviation-ledger.md`（行为偏差）、
> `debt/register.md`(已知债务) 互为对照。最后更新：2026-08-23（M2 第一阶段完成）。

## 一句话现状

能解释执行真实 `javac --release 8` 字节码的运行时内核雏形：单/多线程、异常、
集合、字符串、监视器语义正确（真机 JVM 输出逐字节对照）；**AOT 未开始**，
类库面为窄幅合成集，反射与 Class 元对象缺失。

体量：生产 Go ≈7.5k 行，测试 ≈1.7k 行 / 45 例，`-race` 干净。

## 能力矩阵

### ✅ 可用

| 域 | 内容 |
|---|---|
| 执行引擎 | v52 指令集解释执行（invokedynamic 构建期脱糖；jsr/ret 非法） |
| 异常 | Throwable 家族、异常表分发、隐式抛出（NPE/越界/除零/负长/强转）、SOE（帧预算 4096）、uncaught 打印 |
| 对象模型 | 身份语义、继承链 embedding、接口分派 invokeinterface、字段默认值、数组、UTF-16 String（窄方法面） |
| 并发 | Thread 全生命周期 + 中断三路径；可重入监视器；wait/notify/notifyAll；goroutine 底座 |
| 类加载 | 目录 classpath、懒式依赖（循环守卫）、结构层验证（SM 帧/池合法/边界） |
| CLI | `catty [-cp dir] run <File.class \| dotted.Main>` |

### ❌ 缺失（按解锁价值排序）

| 缺口 | 债务/计划 | 解锁什么 |
|---|---|---|
| AOT 发射器 | M3 | 核心卖点全部（启动/吞吐/部署形态） |
| java.lang.Class 元对象 | DEBT-0010 / P-0005 | ldc-class、静态 synchronized、getClass |
| 数据流类型验证 | DEBT-0009 / P-0005 | 不可信输入防线（当前仅信任 javac） |
| nio/net/charset 映射 | P-0005 | socket 程序、http echo 验收 |
| 包装类 Long/Boolean 等 | 类库 backlog | 相应自动装箱 |
| HashMap/Set/迭代器 | 类库 backlog | 常见数据结构程序 |
| String 宽方法面 | 类库 backlog | substring/indexOf/format… |
| 反射/注解/MethodHandle | M2+ 远期 | Jackson/Spring 级生态 |
| JAR 加载、main args | DEBT-0008 残留 | 部署形态 |
| JNI | ADR-0007 范围内未开工 | 原生库互操作 |
| Reference 四件套 / -Xmx 映射 | M3+ | GC 语义完整 |

## 行为偏差

见 [deviation-ledger.md](specifications/deviation-ledger.md)（DEV-0001..0007）。
要点：验证器无类型层、volatile 标志忽略、守护线程规则未实现、SOE 数值模型不同。

## 里程碑位置

```
M0 ████████ 完成（解释器 + HelloWorld）
M1 ██████████ 完成（Monitor/加载器/验证器结构层）
M2 ████░░░░░░ 第一阶段完成（线程/SOE）；剩：Class 元对象、数据流验证器、net 映射
M3 ░░░░░░░░░░ AOT 发射器（未开始）
M4 ░░░░░░░░░░ 三方基准报告（未开始）
```

## 质量纪律快照

- 机械门：`make check`（文档结构 + go vet/build/test）
- 并发门：全仓 `-race`
- 符合性：fixture 输出与参考 JVM 对照（oracle = 本机 JDK 25 `--release 8`）
- 信誉资产：偏差账本 7 条全公开；债务登记 10 项；ADR 8 份
