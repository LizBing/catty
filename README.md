# Catty

一个以 Go 运行时为底座的 Java 语言运行时（研究项目，目标可产品化）：

- **寄生战略**——GC、调度器、netpoller 直接复用 Go 运行时，不自建；
- **AOT-first**——字节码在构建期翻译为 Go 源码整体编译；内置字节码解释器兜底动态类加载与反射；
- **JDK 8 基线**起步，IR 版本无关，保留升级到现代 JDK 的路径。

## 快速开始

```sh
make check   # 机械校验门（当前阶段校验文档结构完整性）
```

## 文档入口

新会话的 AI Agent 请从 [AGENTS.md](AGENTS.md) 开始；人类请从 [docs/vision.md](docs/vision.md) 开始。

## License

Apache-2.0（主仓库）。第三方来源许可政策见 [docs/quality/architecture-rules.md](docs/quality/architecture-rules.md)。
