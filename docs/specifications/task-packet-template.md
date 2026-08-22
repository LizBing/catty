# Task Packet Template（协议 §10）

Orchestrator 派发 Worker 时按此填写，随任务存档于对应 plan 的 Tasks 条目或 PR 描述。

```markdown
TASK
一句话动词开头的任务描述。

CONTEXT
为什么做：关联 plan 任务 ID、相关 ADR、上游依赖产物。

SUCCESS CRITERIA
可验证的完成条件（测试命令 + 期望输出）。

CONSTRAINTS
不可违反项（引用 architecture-rules 编号即可）+ Scope Expansion Rule 提醒。

RELEVANT FILES
起点文件清单（最小充分）。

DEPENDENCIES
前置任务/产物及其位置。

MANDATORY TESTS
必须新增/保持通过的测试。
```

Worker 回报格式（协议 §10）：Result / Files changed / Tests / Important decisions / Remaining risks / Unexpected findings。禁止只回 "Done."
