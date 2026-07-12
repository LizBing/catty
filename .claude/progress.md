# catty 项目进度快照

**日期**: 2026-07-12
**最后提交**: Block D — R1 加固完成
**CI 状态**: 全绿（含 java.base 烟雾测试）

## R1：✅ 完成并加固

### 提交链
| 提交 | 内容 |
|---|---|
| `ece4412` R1.0 | 异常处理 |
| `e043fc5` R1.1 | invokeinterface + multianewarray + wide |
| `a1b991b`/`78d92f8`/`0a96373` R1.2 | native stubs + bootstrap classpath + 两个根因 bug 修复 |
| `755561b` Block A | classloader Provider 链 + launch 包 |
| `ff6db1d` Block B | java.base 自动检测 + 烟雾测试 |
| `b704b75` Block C | String 内容方法 + System.getProperty |
| Block D | R1 加固：CI 跑 java.base、IR executor bug 修复、文档同步、ADR-0014/0015 |

### Block D 具体成果
- **D1**：CI（`.github/workflows/ci.yml`）pin JDK 25，`jimage extract` 解出 java.base，设 `CATTY_BOOT`，RealBaseSmoke 进 CI 回归
- **D2**：修了 IR executor 一个潜伏两个里程碑的 bug——通用带索引的 `aload`/`astore`（local ≥4 时）错用 int 访问器，引用丢失。修后 InterfaceTest 三引擎全过，`run.sh` 清掉容忍 hack
- **D3**：fib(35) 复测——解释器 4.54s、AOT ~50ms、java -Xint 0.58s、java JIT 0.05s，全部与文档一致，**无性能回归**
- **D4**：ROADMAP（R1=Complete）、ARCHITECTURE（Provider 链、launch、bootstrap 边界、java.base 检测）、README（夹具表 + R1 边界）、DEVELOPMENT（新 layout + registerSynthetic）全部对齐
- **D5**：ADR-0014（合成 String + Extra 策略）、ADR-0015（bootstrap 类边界）
- **D6**：RealBaseSmoke 顶部注释明确标注覆盖范围与刻意排除项（HashMap/Double/Integer.toString → R2 Unsafe 依赖）

## 当前能力（R1 最终）

- 解释器 ~145 操作码 + 异常处理 + 接口分发 + 多维数组
- 6 个不可替代合成类 + ~12 非引导合成类 + ~40 native 注册
- java.base 自动检测（CATTY_BOOT/JAVA_HOME/java_home）
- String 内容方法（equals/hashCode/substring 等 14 个）+ byte[] 构造器解码
- 烟雾测试 18/18（三引擎一致，与 java 输出字节相同）
- AOT：`catty build`，fib(35) ~50ms
- CI 回归保护：纯合成 + java.base 两条路径

## R1 边界（明确属后续阶段）

| 项 | 归属 | 原因 |
|---|---|---|
| Integer/Long.toString、Double.parseDouble、HashMap | R2 | DecimalDigits→Unsafe 级联（~50 native） |
| invokedynamic（0xba） | R3 | LambdaMetafactory |
| Unsafe 全家桶 | R2 | 与 Thread=goroutine、Go 内存模型纠缠 |
| ClassLoader.defineClass | R2/R3 | 动态类加载 |
| JRT 内存直读 | — | 用 jimage extract 工具代替 |
| 完整 UTF-16 代理对 | R6 | rune 感知已覆盖 BMP |

## 战略路线图

- ✅ R1 Run real Java programs（含加固）
- ⬜ R2 多线程（Thread=goroutine）——前置：最小 Unsafe stub 层
- ⬜ R3 反射 + invokedynamic
- ⬜ R4 I/O 集成
- ⬜ R5 AOT 覆盖扩展
- ⬜ R6 性能调优

## ADR 索引

| ADR | 决策 | 状态 |
|---|---|---|
| 0001-0007 | 基础架构（Go runtime/switch dispatch/Slot/原生类/懒 clinit/predecode/分层） | Accepted |
| 0008-0013 | 战略愿景（AOT-first/混合类库/Thread=goroutine/Go 内存模型/escape analysis/Go runtime 集成） | Proposed |
| 0014 | 合成 String + Extra 负载 | Accepted |
| 0015 | bootstrap 类边界（6 个不可替代） | Accepted |
