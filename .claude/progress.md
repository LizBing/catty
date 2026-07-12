# catty 项目进度快照

**日期**: 2026-07-12
**最后提交**: `b704b75` — feat: R1 Block C — String content methods + System.getProperty
**CI 状态**: 全绿

## R1 最终状态：✅ 完成

| 子任务 | 提交 | 验证 |
|---|---|---|
| R1.0 异常处理 | `ece4412` | ExceptionTest 三引擎一致 |
| R1.1 剩余操作码 | `e043fc5` | InterfaceTest 通过 tree-walker |
| R1.2 native stubs + bootstrap classpath | `a1b991b`/`78d92f8`/`0a96373` | 合成类 ~18 个，native 注册表 ~25 方法 |
| Block A 架构清理 | `755561b` | Provider 链 + launch 包 |
| Block B java.base 自动检测 | `ff6db1d` | CATTY_BOOT/JAVA_HOME/java_home 检测 |
| Block C String 内容方法 | `b704b75` | 18/18 烟雾测试 |

**里程碑达成**：`catty -cp . HelloWorld` 自动挂载 java.base，一行命令，输出与真实 java 一致。

## 架构概览（R1 结束后）

```
cmd/jvm/main.go          ← CLI：参数解析 + java.base 自动检测 + 委托 launch
launch/launch.go          ← 运行时启动：构建 loader + 线程 + 进入解释器
classloader/              ← Provider 链（Array→Bootstrap→Synthetic→Classpath）
native/registry.go        ← 合成类 map 注册制 + BootstrapClasses 白名单
native/lang.go            ← Object/String/StringBuilder/System 合成工厂 + String 内容方法
native/system.go          ← 全局 native 注册表（RegisterNative）+ 方法实现
native/exceptions.go      ← 异常层级（~13 类）+ Comparable
native/io.go              ← PrintStream
interpreter/              ← Loop exec + invoke 桥 + IR executor
rtda/                     ← 运行时数据（Slot/Frame/Thread/Class/Method/Object）
lowering/ + transpile/     ← AOT 发射器
runtime/                  ← AOT bridge
```

## 当前能力

### 烟雾测试覆盖（18/18，与 java 输出一致）
PrintStream · String.length/charAt · Object.hashCode/toString · Class.getName/isInterface · StringBuilder · Math.max · Integer.MAX_VALUE/parseInt/toHexString · ArrayList · Exception catch · System.identityHashCode · Long.parseLong/MAX/toHexString · String.equals/hashCode/substring/startsWith/endsWith/concat/isEmpty · System.getProperty · String.indexOf

### 已知限制（明确划给后续阶段）
- **Integer/Long.toString** → JDK 25 的 DecimalDigits→Unsafe 级联（R2 需 Unsafe 最小 stub）
- **Double.parseDouble** → 同上
- **HashMap** → 需 String.equals/hashCode（已修复），但 AbstractMap clinit 触发 Unsafe 级联
- **invokedynamic**（0xba）→ R3
- **JRT 内存直读** → 用 jimage extract 工具代替（符合精简运行时原则）
- **完整 UTF-16 代理对** → R6

## 战略路线图

- ✅ R1 Run real Java programs
- ⬜ R2 多线程 (Thread=goroutine)
- ⬜ R3 反射 + invokedynamic
- ⬜ R4 I/O 集成
- ⬜ R5 AOT 覆盖扩展
- ⬜ R6 性能调优

## 关键架构决策

| ADR | 决策 | 状态 |
|---|---|---|
| 0001 | 套用 Go runtime | Accepted |
| 0002 | switch dispatch | Accepted |
| 0003 | tagged 16B Slot | Accepted |
| 0004 | 原生核心类 (String 保持合成) | Accepted |
| 0005 | 懒 <clinit> | Accepted |
| 0006 | predecode 不划算 | Accepted |
| 0007 | 反射/动态: 分层 | Accepted |
| 0008 | AOT-first | Proposed |
| 0009 | 混合类库 | Proposed |
| 0010 | Thread = goroutine | Proposed |
| 0011 | Go 内存模型 | Proposed |
| 0012 | escape analysis 替代 GC | Proposed |
| 0013 | 直接 Go runtime 集成 | Proposed |
