# 双方协作协议

**状态：** Accepted by LizBing，2026-07-13

## 1. 协作模型

catty 始终由两方协作：

| 角色 | 职责 |
|---|---|
| **Project Owner（LizBing）** | 确定方向、接受合同与 ADR、指定验收方式、决定合并与发布 |
| **Active Agent** | 在授权范围内分析、设计、实现、自审、测试、记录证据并交付候选提交 |

Active Agent 是模型无关的临时角色，可以由 Claude、Codex、DeepSeek 或其他工具
承担。更换模型只是更换第二方执行者，不会增加新的长期角色，也不会自动形成多代理
流水线。

同一 workstream 同一时刻只有一个 Active Agent。Owner 可以临时让另一个 agent 做
只读复核，但这是 Owner 选择的辅助动作，不是项目的固定第三方流程。

## 2. 授权与成本边界

Owner 的当前请求决定 Active Agent 的工作模式：

- **讨论/制度模式：** 只讨论或编辑治理文档，不读取或修改实现，除非 Owner 明确要求。
- **只读审查模式：** 可以读取代码、提交和测试证据，但不修改代码、测试、harness、CI，
  不 commit、merge 或 push。
- **实现模式：** 可以在合同范围内修改代码和测试，运行必要门禁并交付候选提交。
- **集成模式：** 只有 Owner 明确授权后，才可以 commit、merge、push、发布或改写远端状态。

授权不会因任务困难或发现顺手修复而自动扩大。Agent 必须优先完成 Owner 指定的结果，
不得为了“做完整”擅自进入更昂贵的实现、全仓审计或多轮验证。

模型订阅额度、运行时间和外部服务费用都是项目约束。Agent 应选择满足合同所需的最小
工作量；可并行的读取可以并行，但不得为了形式完整重复执行已有且可信的证据。

## 3. 事实来源

发生冲突时，按以下顺序判断项目事实：

1. 可执行行为和可复现测试证据
2. Accepted ADR
3. Accepted active workstream contract
4. 项目状态文档
5. 架构、路线图和开发文档
6. Commit message 和 handoff
7. 聊天记录

聊天记录不能覆盖已经接受的仓库决策。Commit 只能定位实现，不能单独证明行为正确。

## 4. 规划体系

项目规划分为四层。每层只回答一种问题，不互相替代：

| 文档 | 回答的问题 | 是否授权实现 |
|---|---|---|
| `ROADMAP.md` | 项目长期往哪里走、阶段和依赖顺序是什么 | 否 |
| ADR | 为什么选择某个长期架构决策 | 只有 Accepted 后才能作为实现依据 |
| Workstream contract | 这一轮具体做什么、受什么约束、如何验收 | Accepted 后授权合同范围内的实现 |
| `PROJECT_STATUS.md` | 当前稳定基线、active workstream、能力边界和下一动作是什么 | 否，只提供当前状态 |

### 4.1 Roadmap

Roadmap 是战略方向和阶段顺序，不是任务列表，也不是实施授权。它可以记录：

- milestone 和目标能力；
- 阶段依赖和大致顺序；
- 明确的 capability boundary；
- 已完成、当前和未来阶段。

Roadmap 不记录 session 任务、Agent 身份、具体文件修改或未经验证的完成声明。Roadmap
中的未来项目必须先经过 ADR/workstream 流程，Agent 不得直接据此开始 non-trivial
implementation。

### 4.2 ADR 在规划中的位置

当 Roadmap 项目涉及长期语义承诺、兼容性偏离、跨包机制、公共边界或反转旧决策时，
必须先解决 ADR：

```text
Proposed → Accepted → Superseded
         ↘ Withdrawn
```

- Proposed ADR 用于讨论，不能作为生产实现依据。
- 只有 Owner 可以接受 ADR。
- Workstream 不得依赖尚未 Accepted 的 ADR 开始正式实现。
- Accepted ADR 只能由新的 superseding ADR 反转。
- Withdrawn 关闭从未生效的提案；Superseded 只用于曾经 Accepted 的决策。

如果决策缺少事实基础，先建立 research workstream 获取数据，再提出或修订 ADR。

### 4.3 Workstream 在规划中的位置

Workstream 将 Roadmap 方向和 Accepted ADR 转换为一次可交付变更。每个合同必须声明：

- 对应的 Roadmap item；
- governing ADRs；
- prerequisites；
- 类型：`research` 或 `implementation`。

Roadmap 解释“为什么现在做”，ADR 规定“受什么长期决策约束”，workstream 定义“这一轮
具体做到哪里”。没有 Roadmap 项时可以由 Owner 直接发起 workstream，但必须说明原因；
这不自动修改长期路线图。

### 4.4 Project Status

`PROJECT_STATUS.md` 是新 session 的唯一当前状态入口，只记录：

- stable baseline 和 commit；
- 当前 milestone；
- active workstream（没有则写 None）；
- 当前 capability boundary；
- 已接受但尚未实施的关键方向；
- 下一协调动作和已知阻塞。

它不复制 workstream 的完整计划，也不把 Ready、Not run 或 Not implemented 的能力写成
Done。

### 4.5 规划流转

```text
Owner 提出产品意图
        ↓
检查 PROJECT_STATUS 和 ROADMAP
        ↓
是否存在未决架构问题？
   ├─ 是 → Proposed ADR → 讨论/研究 → Owner Accepted
   └─ 否
        ↓
Proposed workstream
        ↓
Owner 接受 Outcome / Non-scope / Constraints / Gates / Review
        ↓
Accepted → In Progress → Ready
        ↓
按 self/owner 规则验收
        ↓
Done + 集成
        ↓
更新 PROJECT_STATUS；仅在战略状态变化时更新 ROADMAP
```

### 4.6 Research workstream

研究型 workstream 用于在架构决策前获得调用图、实验、基准或兼容性证据：

- `Type: research`，默认 `Review: owner`；
- Outcome 是报告、数据、原型或 Proposed ADR；
- Non-scope 必须明确不进入生产能力；
- 研究原型不能因测试通过自动成为生产实现；
- Owner 接受结论或 ADR 后，另建 implementation workstream。

### 4.7 完成后的规划更新

Workstream Done 时：

- 总是固定最终证据和 candidate/integration commit；
- active workstream、稳定基线或 capability boundary 变化时更新 `PROJECT_STATUS.md`；
- milestone 完成、阶段顺序或长期范围变化时才更新 `ROADMAP.md`；
- 架构决策变化时新增 ADR 或更新 ADR 状态；
- 普通 commit 不触发 Roadmap 或 ADR 修改。

## 5. 工作分类

### 5.1 Trivial work

满足以下全部条件的工作可以不建 workstream：

- 单 session 可完成；
- 不改变 JVM/JRE 可观察语义；
- 不跨越重要包边界；
- 不需要新的架构决策；
- 回滚简单，验收命令明确。

例如拼写、链接、注释、机械格式化和局部低风险修复。

### 5.2 Non-trivial work

出现以下任一情况必须先建立 `docs/workstreams/<name>.md`：

- 跨 session 或跨包；
- 改变 JVM/JRE 语义、兼容性边界或失败行为；
- 涉及 class loading、AOT、异常、Thread、monitor、JMM、Unsafe 或 native 边界；
- 需要多个实现 slice；
- Owner 需要在实现前确认 scope、non-scope 或验收门。

## 6. Workstream contract

一个 workstream 只维护一个文件，包含稳定合同、动态计划和必要 handoff。

### 6.1 状态

```text
Proposed → Accepted → In Progress → Ready → Done
```

- **Proposed：** Active Agent 起草，尚未授权实现。
- **Accepted：** Owner 接受目标、边界、语义约束和验收方式。
- **In Progress：** Active Agent 正在合同范围内工作。
- **Ready：** Agent 已自审、运行其声明的门禁并固定候选 commit，等待最终验收。
- **Done：** 满足合同规定的验收方式并完成集成。

`Ready` 不等于 `Done`。尚未运行、跳过或仅靠代码目检的 gate 不能标为 Pass。

### 6.2 Review 类型

每个 workstream 在 Accepted 前确定一种方式：

| Review | 使用场景 | Done 权限 |
|---|---|---|
| `self` | 低风险、局部、不改变语义边界 | Active Agent 完成门禁后可标 Done |
| `owner` | 跨包、重要能力、JVM/JRE 语义或 Owner 指定 | 只有 Owner 接受后可标 Done |

Owner 可以在 `owner` review 时临时使用另一个模型复核，但不需要在合同中设置固定的
“Codex reviewer”“Claude implementer”等模型身份。

### 6.3 合同冻结

合同进入 Accepted 后，以下部分冻结：

- Outcome
- Scope
- Non-scope
- Semantic constraints
- Acceptance gates
- Review 类型

需求变化只能在 `Amendments` 中追加日期、原因和影响，并由 Owner 接受。不得在实现完成
后删除、压缩或放宽原验收标准。计划状态、证据结果和 handoff 可以正常更新。

## 7. ADR

以下变更需要 ADR：

- 改变 JVM/JRE 语义承诺或能力边界；
- 引入或替换跨包运行时机制；
- 改变长期架构方向；
- 接受已知语义偏离；
- 反转已有 Accepted ADR。

ADR 状态为 `Proposed → Accepted → Superseded`，Proposed 也可以转为 `Withdrawn`。
只有 Owner 可以接受 ADR。Accepted ADR 不得直接改写以反转决策；必须新增
superseding ADR。

普通实现细节、测试策略、局部优化和合同已覆盖的 opcode/native 增量通常不需要 ADR。

## 8. 分支与集成

- `main` 是唯一稳定集成线。
- Non-trivial work 使用以 workstream 命名的独立分支或 worktree。
- 审查以固定 candidate commit 为准，不以持续变化的脏工作树为准。
- Active Agent 不直接改写 `main` 或远端，除非 Owner 明确授权。
- 合并前必须保护 Owner 的现有修改，不清理或覆盖无关脏文件。
- `PROJECT_STATUS` 只在项目级能力、active workstream 或稳定基线变化后更新。

## 9. 证据与验收

每个 capability claim 必须对应一个明确 gate。有效证据包括：

- 精确测试名称和命令；
- 实际退出状态；
- interpreter / IR / AOT 等引擎矩阵；
- HotSpot/JDK 差分结果；
- CI run；
- race、stress、timeout 或 benchmark 结果；
- 确定性生成产物。

以下不构成通过证据：

- 只有 commit hash；
- “代码看起来已接入”；
- 输出中碰巧包含 `PASS`；
- 缺少依赖后退出 0；
- 没有覆盖相关并发路径的 `go test -race`；
- 把未实现的引擎静默排除后宣称矩阵通过。

结果统一记录为：

- **Pass：** 命令实际运行并满足断言；
- **Fail：** 行为或证据不满足合同；
- **Not run：** 尚未执行或环境缺失；
- **Not implemented：** 能力尚未实现，并由合同明确允许延后。

只有合同允许的 `Not implemented` 可以不阻塞当前 workstream；`Not run` 不能计为 Pass。

## 10. Harness 原则

合同定义 harness 必须证明什么，Active Agent 负责具体实现。最低原则：

- 缺少必需工具、JDK、fixture 或输入时 fail-closed；
- 每个子进程有 wall-clock timeout；
- 保留退出码、stdout 和 stderr；
- 不吞掉 panic、异常、race、deadlock 或非零退出；
- 精确比较预期结果，不用搜索 `PASS` 代替断言；
- interpreter、IR、AOT 分别报告；
- 并发测试使用确定性协调，stress 是补充而非替代；
- 未实现能力显式报告，不能算入通过数；
- CI 运行有界版本，较长 stress 可以独立运行。

## 11. Handoff

Handoff 只为恢复未完成现场服务，记录无法从 Git 和代码直接推导的信息：

- 分支和 candidate/base commit；
- 是否有脏文件；
- 最后工作位置；
- 已运行与未运行的门禁；
- 当前阻塞；
- 下一条具体动作；
- 关键但未落入代码的上下文。

不复制 git log，不写进度时间线，不保留仪式性字段。工作已经 Done 且工作树干净时，
handoff 可以省略或只保留一句最终定位。

## 12. 标准流程

1. Owner 提出目标和当前授权模式。
2. Active Agent 读取仓库事实，判断 trivial 或 non-trivial。
3. Non-trivial work 先提交 Proposed contract；Owner 接受后才能实现。
4. Agent 在独立分支按合同实现，不扩大 scope。
5. Agent 自审并运行准确 gates，逐项记录 Pass/Fail/Not run/Not implemented。
6. Agent 固定 candidate commit，更新为 Ready。
7. `self` review 可由 Agent 完成集成；`owner` review 等待 Owner 决定。
8. 未通过则保持 Ready 或退回 In Progress，不降低验收标准。
9. 通过后集成、标记 Done，并仅在项目级状态变化时更新状态文档。

## 13. 成功标准

这套制度追求四个同时成立的结果：

| 目标 | 制度保障 |
|---|---|
| 安全 | 冻结合同、明确 non-scope、重要工作 owner review、fail-closed 证据 |
| 高效 | 单一 Active Agent、无固定多模型交接、trivial work 免合同 |
| 精确 | capability claim 对应具体 gate，区分 Ready 与 Done |
| 快速 | 短合同、单一执行线、按风险选择 review、最小充分证据 |
