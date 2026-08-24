# P-0005: M2 第二阶段 — Class 元对象 + 网络映射 + http echo

- 状态：completed（里程碑验收通过，2026-08-24 归档）
- 风险评分：Medium（Class 对象触及分派/监视器；网络为新增边界面）

## Goal

解锁 DEBT-0010 的三特性（ldc-class / 静态 synchronized / getClass），打通
最小 socket 面，并以"纯 Java 写的多线程 HTTP echo 服务器 + 真实 socket
对照请求"作为里程碑验收。

## Current State

M2 第一阶段完成（线程/SOE）。ldc <class> 抛不支持；方法级 synchronized
标志被忽略；无任何 IO。

## Target State

1. `java/lang/Class` 元对象：每个 kernel.Class 惰性持有唯一实例；
   ldc CLASS 推入；`Object.getClass()`；静态 synchronized 以 Class 对象
   监视器实现；**同时修复实例 synchronized 方法标志被忽略的缺口**。
2. `kernel/net.go`：ServerSocket(bind/accept/close/localPort)、
   Socket(读写流/close)、InputStream/OutputStream（payload 挂 Go 句柄）。
3. String 字节转换：getBytes()、new String([BII)、indexOf(String)。
4. fixture HttpEcho：主线程 accept 循环 + 每连接 Thread 处理（读头→回写）。
5. 验收：Go 测试内起服务、真 socket 发请求、断言响应。

## Tasks（DAG）

```
T1 java/lang/Class + 类对象缓存 + getClass + ldc 支持
T2 synchronized 方法标志（实例+静态）——独立可并行
T3 String 字节转换 natives ——独立
T4 net natives（依赖 T1：accept 需要 Class 对象包装）
T5 HttpEcho fixture + 真 socket 验收 + -race
   ↑ T1..T4
后续（本计划外）：DEBT-0009 数据流验证器单列专注任务——连续两次被挤后，
本次明确不并行，待 echo 收口后以整段预算实施。
```

## Validation

- make check 绿、全仓 -race 绿
- curl 或 Go dial 对 HttpEcho 的响应断言通过
- 静态 synchronized 多线程计数精确

## Risks

- 阻塞 IO 与中断语义冲突（R4）：socket 读暂未接 interrupt 唤醒 →
  登记 DEV-0008/DEBT，M2 后期统一 SetDeadline 方案。
- Class 对象身份与用户 == 比较 → 缓存必须稳定（atomic 惰性）。

## Progress

- [x] 计划建立（2026-08-23）
- [x] T1 Class 元对象（惰性缓存/ldc/getClass/NewInstance 拒绝手动实例化）
- [x] T2 synchronized 方法标志（实例接收者 + 静态 Class 监视器）
- [x] T3 String getBytes/new String([BII)/indexOf/Integer.parseInt
- [x] T4 net natives（ServerSocket/Socket/InputStream/OutputStream，
      payload 挂 net.Listener/net.Conn/streamHandle）
- [x] T5 HttpEcho 真 socket 验收：响应体 `77` 与参考 JVM 一致；
      进程内测试走 port 0+getLocalPort+真实 dial 断言
- 验证期修复：验证器分支扫描 off-by-one（`i+3>=len` 误判满装 goto）、
  OutputStream write([B)V 重载缺失、StringBuilder.indexOf 超出合成面
- [x] **DEBT-0009 数据流验证器专项完成（同轮追加，整段专注预算）**：
      表驱动效果表(~30族)+不规则指令专写；SM帧型语义上移验证器
      （full替换/append增量/chop裁剪/same继承+入口前缀仅相对帧）；
      checkcast 方向修正+成功后类型收窄；异常处理帧兼容；
      状态快照防原地突变。全部真实 fixture 过数据流；篡改样本被拒。
      验证期修复：newarray 长度(2B)、s2 符号扩展(二次)、iinc 注册、
      popExpect unknown 语义、setLocalPair 幻影槽、双重栈压入
