# 安全质量要求

- 地位：协议 §19-E3 的项目侧展开；安全红线等同架构红线。

## 威胁面排序

1. **classfile 解析器**（最大外部输入攻击面）：恶意/损坏 .class 不得导致 panic 越界或无限循环。fuzz harness 已登记 DEBT-0001；
2. **JNI 边界**（ADR-0007）：native 代码不可信等级最高——句柄表隔离、pending exception 检查、禁止内核热路径 cgo;
3. **AOT 发射器输出**：生成的 Go 源码按 generated 对象对待，禁止手工修补绕过验证。

## 规则

- cgo 默认关闭（R1）；开启 jnion 构建的产品定位为"信任宿主 native 库"的场景，文档必须明示;
- 供应链：依赖最小化原则；引入首个第三方依赖时将 `govulncheck` 纳入 make check;
- 许可证合规属安全范畴：见 architecture-rules R2 与 DEBT-0002;
- 安全问题升级走 E3，不受静默原则约束。
