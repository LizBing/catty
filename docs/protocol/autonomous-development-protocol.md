# Autonomous Development Protocol
## AI 自治研发协议 v0.1

### 0. 目标

本协议的目标是：

> 在保证质量、可恢复性和可审计性的前提下，将软件研发中尽可能多的研究、设计、实现、测试、评审、项目管理和维护工作交给 AI 完成。

默认原则：

**AI 自治执行，人类按异常介入。**

人类不负责：
- 日常任务拆分
- 普通技术调研
- 常规代码实现
- 普通代码审查
- 测试补充
- 文档维护
- Issue 更新
- PR 修复
- 常规重构
- 依赖升级
- 技术债整理

人类主要负责：

1. 定义目标
2. 定义不可违反的边界
3. 做高风险或不可逆决策
4. 判断最终产品是否符合意图

---

# 1. 系统基本原则

## 1.1 默认自治

除非满足“升级条件”，Agent 不得因为存在多个可行方案而询问人类。

如果存在：

- A 方案：80 分
- B 方案：75 分
- C 方案：70 分

Agent 应自行选择 A。

不得询问：

> “您更喜欢哪个？”

只有当不同方案涉及明显不同的产品方向、不可逆架构或者重大成本时，才能升级。

---

## 1.2 可逆决策自动做

对于容易撤销的决策：

```text
变量命名
内部 API
模块拆分
小型重构
测试组织
内部数据结构
普通依赖
实现算法
日志格式
```

AI 自行决定。

对于难以撤销的决策：

```text
公共 API
数据格式
数据库 schema
网络协议
长期兼容性承诺
安全模型
大型架构变更
重大第三方依赖
```

需要更严格评审。

---

## 1.3 Evidence before opinion

Agent 不应说：

> “我感觉这种设计更快。”

应该说：

```text
Hypothesis:
方案 A 可能降低 dispatch overhead。

Evidence:
benchmark 当前实现
benchmark A
benchmark B

Result:
A median -18.7%
p95 -13.2%

Decision:
采用 A。
```

对于能够实验的问题：

> 优先实验，而不是讨论。

---

## 1.4 Repository is memory

聊天上下文不是项目真相来源。

项目状态必须进入仓库。

```text
repo/
├── AGENTS.md
├── docs/
├── src/
├── tests/
├── benchmarks/
└── tools/
```

Agent 应当假设：

> 下一次执行自己的 Agent 完全不知道今天发生过什么。

因此重要信息必须写入 repo。

---

# 2. 人类与 AI 的接口

人类提交的最小输入称为：

# Goal Card

格式：

```text
GOAL
希望实现什么。

WHY
为什么需要它。

SUCCESS
怎样算成功。

CONSTRAINTS
不能违反什么。

OPTIONAL
偏好，但允许 AI 自行改变。
```

例如：

```text
GOAL

为 JVM 增加并发类加载能力。

WHY

当前类加载器成为多线程启动时的瓶颈。

SUCCESS

8 线程启动 benchmark 至少提升 30%。

不得出现重复类初始化。

不得引入新的 deadlock。

CONSTRAINTS

保持现有 ClassLoader API。

不能引入重量级全局锁。

OPTIONAL

尽量保持实现简单。
```

除此之外，人类原则上不需要继续指导实现。

---

# 3. 项目知识结构

推荐：

```text
AGENTS.md

docs/
├── vision.md
├── architecture.md
│
├── specifications/
│
├── research/
│
├── decisions/
│   └── ADR-xxxx.md
│
├── plans/
│   ├── active/
│   ├── blocked/
│   └── completed/
│
├── quality/
│   ├── testing.md
│   ├── performance.md
│   ├── security.md
│   └── architecture-rules.md
│
└── debt/
    └── register.md
```

## AGENTS.md

必须保持短小。

主要作用是：

> 项目地图。

而不是把所有规则全部复制进去。

例如：

```text
# Project

Read docs/vision.md before making product decisions.

Architecture:
docs/architecture.md

Architectural decisions:
docs/decisions/

Active plans:
docs/plans/active/

Quality requirements:
docs/quality/

# Mandatory validation

Run:

make check

before declaring a task complete.

# Rules

Never bypass failing tests.

Never modify generated files manually.

Architectural changes require ADR.

Large changes require an execution plan.
```

---

# 4. Agent 类型

不建立：

```text
CEO Agent
CTO Agent
Senior Developer Agent
Junior Developer Agent
```

而建立能力型 Agent。

核心只有：

```text
Orchestrator
Researcher
Planner
Worker
Reviewer
Maintainer
```

---

# 5. Orchestrator

Orchestrator 是整个系统唯一的协调者。

负责：

```text
接收 Goal
理解项目状态
判断工作规模
安排研究
安排设计
拆分任务
建立依赖关系
启动 Worker
追踪结果
处理失败
安排 Review
决定是否合并
决定是否升级给 Human
```

Orchestrator 原则上不直接写大量业务代码。

它主要进行：

> 决策 + 调度。

---

# 6. Researcher

Research Agent 默认只读。

权限：

```text
READ repo
READ internet
READ dependency source
RUN experiments
WRITE docs/research/
```

禁止：

```text
修改 production code
合并 PR
改变架构
```

研究报告模板：

```text
# Question

需要回答什么？

# Context

为什么需要研究？

# Candidates

可能方案有哪些？

# Evidence

论文
文档
源码
benchmark
实验结果

# Findings

发现了什么？

# Recommendation

推荐方案。

# Confidence

0.0 - 1.0

# Unknowns

仍然不知道什么？
```

---

# 7. Architecture Protocol

中大型设计不得由单个 Agent 直接决定。

推荐：

```text
Research
   ↓
Design A
Design B
Design C
   ↓
Critic
   ↓
Decision
```

不同设计应从不同偏好出发：

```text
A：最简单

B：最高性能

C：最佳长期维护性
```

Critic 的职责不是重新设计，而是攻击现有方案。

必须检查：

```text
是否过度设计？
是否存在更简单方案？
是否违反当前架构？
有没有隐含状态？
有没有并发风险？
失败模式是什么？
能否 rollback？
如何测试？
如何观察？
```

重大决策形成 ADR：

```text
docs/decisions/ADR-0042.md
```

格式：

```text
Context

Decision

Alternatives

Evidence

Consequences

Risks

Rollback
```

---

# 8. Execution Plan

中大型任务必须建立：

```text
docs/plans/active/<feature>.md
```

内容：

```text
Goal

Current State

Target State

Tasks

Dependencies

Validation

Risks

Progress

Open Questions
```

Plan 必须不断更新。

不得出现：

> 实际代码已经完全改变，但 Plan 还是三天前的状态。

---

# 9. Task Graph

Planner 将工作拆成 DAG，而不是简单 Todo List。

例如：

```text
T1 API design

T2 loader state machine
   ↑
   T1

T3 concurrent table
   ↑
   T1

T4 cycle detection
   ↑
   T2 + T3

T5 stress test
   ↑
   T4

T6 documentation
   ↑
   all
```

可以并行的任务必须并行。

---

# 10. Worker Protocol

每一个 Worker 必须在隔离工作区工作。

例如：

```text
worktree/task-041
worktree/task-042
worktree/task-043
```

Worker 收到一个：

# Task Packet

```text
TASK

CONTEXT

SUCCESS CRITERIA

CONSTRAINTS

RELEVANT FILES

DEPENDENCIES

MANDATORY TESTS
```

Worker 工作循环：

```text
Understand
↓
Inspect
↓
Implement
↓
Test
↓
Self Review
↓
Commit
↓
Report
```

Worker 不得在任务完成后只说：

> Done.

必须提供：

```text
Result

Files changed

Tests

Important decisions

Remaining risks

Unexpected findings
```

---

# 11. Worker 自治原则

Worker 遇到普通问题时：

```text
发现 bug
→ 修

发现缺测试
→ 加

发现文档错误
→ 修

发现明显的小型设计问题
→ 修

需要轻微重构
→ 做
```

不得频繁询问：

> “我发现这里有一个小 bug，要不要一起修？”

默认：

> 修。

---

# 12. Scope Expansion Rule

为了防止 Agent 无限扩张工作范围：

如果额外工作：

```text
< 原任务预计成本的 20%
```

可以直接处理。

如果：

```text
20% - 100%
```

创建独立 Task。

如果：

```text
> 原任务成本
```

交给 Orchestrator 重新规划。

---

# 13. Review Protocol

代码进入 main 前至少经过：

```text
Worker self-review
↓
Independent review
↓
Mechanical validation
```

对于高风险修改：

```text
Correctness Reviewer

Architecture Reviewer

Security Reviewer

Performance Reviewer
```

但不得让所有 PR 都经过四个模型。

Review 数量应根据风险动态决定。

---

# 14. Reviewer 原则

Reviewer 不问：

> “这代码看起来好吗？”

而检查：

```text
Correctness

是否真的满足 Goal？

Edge cases

哪些输入可能失败？

Architecture

有没有破坏边界？

Complexity

有没有更简单的方法？

Regression

可能破坏什么？

Testing

测试是否覆盖真正风险？

Observability

如果线上失败怎么发现？
```

Reviewer 输出：

```text
APPROVE

或者

REQUEST CHANGES

Blockers:
...

Non-blocking:
...
```

---

# 15. Review Fix Loop

出现问题后自动：

```text
Reviewer
   ↓
Fix Agent
   ↓
Tests
   ↓
Reviewer
```

循环直到：

```text
PASS
```

或者达到：

```text
Escalation Threshold
```

不得第一次 Review 失败就找 Human。

---

# 16. Mechanical Constitution

能写成程序的规则，不写成 prompt。

例如：

```text
cargo test
cargo clippy
cargo fmt

dependency lint

architecture test

API compatibility test

benchmark threshold

coverage threshold

security scan
```

Agent 无权说：

> “这个失败应该没关系。”

除非规则本身明确允许 override。

---

# 17. Definition of Done

任务只有同时满足以下条件才能完成：

```text
Implementation complete

Tests pass

Required benchmark passes

Independent review passes

Docs updated

Plan updated

No unexplained TODO

No known blocking regression

Repository clean
```

---

# 18. Merge Policy

满足以下条件：

```text
CI green

Required reviewers approve

No unresolved blocker

Architecture constraints pass
```

AI 可以自动 merge。

不要求 Human 点击 Merge。

---

# 19. Human Escalation Protocol

这是整套协议最关键的一部分。

Agent 只有以下情况可以打扰 Human。

## E1 — Product ambiguity

例如：

```text
不同方案会产生明显不同的用户体验。
```

---

## E2 — Irreversible decision

例如：

```text
public API

file format

database schema

protocol

compatibility promise
```

---

## E3 — Security / Legal / Privacy

出现明显高风险。

---

## E4 — Budget threshold

预计资源超过预算。

例如：

```text
> $50 compute

> 8h agent runtime

> 重大基础设施成本
```

具体数字由项目配置。

---

## E5 — Goal conflict

例如：

```text
startup <100ms

同时要求：

maximum optimization
```

无法同时满足。

---

## E6 — Repeated failure

例如连续：

```text
3
```

次不同方案均失败。

---

## E7 — Confidence collapse

Agent 对重大决策：

```text
confidence < 0.6
```

并且无法通过实验进一步降低不确定性。

---

# 20. Human Escalation Message

Agent 不得说：

> “遇到了问题，请问怎么办？”

必须提供：

```text
DECISION REQUIRED

Problem:
...

Why human input is required:
...

Option A:
...

Option B:
...

Recommendation:
A

Cost of delaying:
...

Default if no decision:
A
```

也就是说：

> 即使找 Human，也必须尽量减少 Human 的思考成本。

---

# 21. Silence Principle

如果项目正常运行：

> 不通知 Human。

不需要发送：

```text
Task 42 started

Task 42 analyzing

Task 42 coding

Task 42 testing

Task 42 completed
```

默认只汇报：

```text
重要结果
异常
需要决策
里程碑
```

---

# 22. Daily Project Brief

推荐每天最多一次：

```text
PROJECT BRIEF

Completed:
3

In progress:
4

Blocked:
0

Merged:
PR #41
PR #43
PR #44

Performance:
startup -12%

New risks:
none

Human decisions required:
none
```

Human 可以完全不回复。

---

# 23. Background Maintenance

系统应长期运行 Maintainer。

周期扫描：

```text
stale docs

dead code

TODO / FIXME

duplicated logic

dependency drift

flaky tests

benchmark regressions

architecture violations

oversized modules

unused APIs
```

小问题：

> 自动 PR。

大问题：

> 创建技术债任务。

---

# 24. Tech Debt Register

维护：

```text
docs/debt/register.md
```

每项：

```text
ID

Description

Impact

Urgency

Estimated Cost

Suggested Fix
```

Orchestrator 定期选择：

```text
高收益 / 低成本
```

项目进行空闲维护。

---

# 25. Failure Protocol

Agent 工作失败时：

禁止：

```text
retry
retry
retry
retry
```

使用：

```text
Attempt 1
↓
Analyze failure
↓
Form new hypothesis
↓
Attempt 2
```

每次失败必须留下：

```text
What failed?

Why?

What was learned?

What assumption changed?
```

连续失败三次后重新研究或重新设计。

---

# 26. Benchmark Protocol

性能问题禁止仅靠静态推理。

流程：

```text
Baseline

Hypothesis

Implementation

Measurement

Comparison

Decision
```

如果性能变化落入噪声范围：

> 不得宣称优化成功。

---

# 27. Agent 权限等级

定义四级权限。

## L0 Read

```text
research
inspection
analysis
```

## L1 Modify

```text
修改 branch/worktree
运行 tests
commit
```

## L2 Integrate

```text
创建 PR
review
merge
```

## L3 External

```text
deploy
release
修改 production
产生费用
操作真实用户数据
```

默认：

```text
Researcher L0

Worker L1

Reviewer L0/L1

Orchestrator L2
```

L3 应单独配置。

---

# 28. 风险评分

每个任务自动得到：

```text
Risk = Impact × Irreversibility × Uncertainty
```

大致分：

```text
0-3 Low

4-6 Medium

7-9 High
```

Low：

```text
1 Worker
CI
auto merge
```

Medium：

```text
Worker
Reviewer
CI
auto merge
```

High：

```text
Research
multiple designs
multiple reviewers
Human checkpoint
```

---

# 29. Model Routing

Agent 是职责。

Model 是资源。

不要：

```text
Claude = Architect

GPT = Programmer
```

而应该：

```text
简单任务
→ 快模型

复杂推理
→ 强模型

大规模搜索
→ 多个廉价 Agent

关键评审
→ 强模型 / 不同模型

极高风险
→ 多模型交叉验证
```

---

# 30. Context Protocol

Agent 开始任务时：

不要读取整个 repo。

采用：

```text
AGENTS.md
↓
Goal
↓
相关 architecture
↓
相关 plan
↓
相关源码
```

逐步加载。

原则：

> Minimum sufficient context.

---

# 31. Communication Protocol

Agent 之间通过 Artifact 交流。

不要依赖：

```text
Agent A 的聊天记忆
```

应该：

```text
research report

ADR

plan

task packet

commit

PR

review
```

这些都属于可持久化协议对象。

---

# 32. Orchestrator 主循环

整个系统可以抽象成：

```text
while project_active:

    observe_project()

    identify_next_actions()

    prioritize()

    for task in runnable_tasks:
        spawn_agent(task)

    collect_results()

    validate()

    update_project_state()

    if escalation_required:
        ask_human()

    maintain_repository()
```

Human 不处于 loop 内。

Human 是：

```text
exception handler
```

---

# 33. 最终系统状态机

```text
IDEA
 │
 ▼
GOAL
 │
 ▼
RESEARCH ────────┐
 │               │
 ▼               │
DESIGN            │
 │               │
 ▼               │
PLAN              │
 │               │
 ▼               │
IMPLEMENT         │
 │               │
 ▼               │
VALIDATE          │
 │               │
 ▼               │
REVIEW ──fail─────┘
 │
 pass
 │
 ▼
MERGE
 │
 ▼
OBSERVE
 │
 ▼
MAINTAIN
```

任意状态如果：

```text
risk too high

goal conflict

uncertainty too high

repeated failure
```

进入：

```text
ESCALATE
```

然后再回到正常状态机。

---

# 34. 自治系统最重要的指标

不要主要统计：

```text
Token 数量

AI 写了多少行代码
```

重点统计：

```text
Human Minutes / Feature

Human Interruptions / Week

Idea → Merge Time

First-pass Test Rate

Rollback Rate

Escaped Defects

Cost / Accepted Feature

Autonomous Merge Ratio
```

最核心的两个：

```text
Human Minutes / Accepted Feature

Human Interruptions / Week
```

---

# 35. 推荐目标

初期：

```text
50% autonomous PR

< 30 min human / feature
```

成熟后：

```text
90%+ autonomous PR

< 10 min human / feature

< 3 interruptions / week
```

最终目标：

> Human 主要定义 Goal，而不是监督过程。

---

# 36. Anti-patterns

禁止以下模式。

## Prompt Manager

人类每天：

```text
你先看看这个

然后改这个

然后跑一下 test

然后 commit

再看看有没有 bug
```

这只是把键盘换成聊天窗口。

---

## Agent Roleplay Company

```text
CEO
CTO
VP Engineering
Manager
Senior
Junior
```

如果这些 Agent 没有不同权限和输入输出协议，只是在浪费 token。

---

## Infinite Debate

三个 Architect 聊 40 轮。

正确规则：

> 可以实验的问题直接实验。

---

## Giant AGENTS.md

不要构建一个 20,000 行的超级 prompt。

应该：

```text
AGENTS.md → index

docs → knowledge

skills → procedures

tests → laws
```

---

## Human Review Everything

如果人类必须 review 每一行代码：

> 系统实际上没有自治。

---

# 37. 最小可行版本

不要第一天做完整系统。

第一阶段只需要五个角色：

```text
Orchestrator

Researcher

Worker

Reviewer

Maintainer
```

以及：

```text
AGENTS.md

docs/

make check

Git worktree

PR
```

流程：

```text
Human Goal
↓
Orchestrator
↓
Plan
↓
Workers
↓
Reviewer
↓
CI
↓
Merge
```

先让这个循环稳定。

然后逐渐加入：

```text
Research swarm

architecture debate

risk scoring

scheduled maintenance

model routing

automatic release

production observation
```

---

# 38. 核心哲学

最终整个协议可以压缩成六句话：

1. **Human defines intent.**
2. **Repository stores truth.**
3. **Agents make reversible decisions.**
4. **Programs enforce deterministic rules.**
5. **Agents escalate exceptions, not routine work.**
6. **Optimize human attention, not token usage.**

最终目标不是：

> AI 写代码。

而是：

> **AI 维护一个能够持续开发软件的系统。**

Human 从：

```text
programmer
```

逐渐变成：

```text
goal setter
+
product evaluator
+
exception handler
```

这就是本协议希望实现的自治程度。
