# 升级策略（协议 §19 的项目本地化配置）

- 地位：协议 E1–E7 在本项目的具体数字与流程。协议文本本身不改。

## 本地数字

| 条目 | 项目配置 |
|---|---|
| E4 预算阈值 | 单一目标累计 >$50 compute，**或** >8h agent runtime，**或**需新增付费服务/基础设施 |
| E6 重复失败 | 同一问题 **3 次**不同方案均失败（每次须留 Failure Protocol 记录） |
| E7 置信度坍塌 | 重大决策 confidence < **0.6** 且无可行实验降低不确定性 |

## 渠道

1. 主渠道：当前人类会话（DeepSeek Harness GUI）；
2. 备渠道：仓库 issue，标签 `decision-required`。

## 升级消息格式

必须使用协议 §20 的 `DECISION REQUIRED` 结构（Problem / Why human input is required / Option A/B… / Recommendation / Cost of delaying / Default if no decision）。禁止只说"遇到问题请问怎么办"。

## 超时与默认行为

- 决策可逆：人类 48h 未响应 → 按 Recommendation 继续，并在下一份 brief 披露；
- 决策不可逆（E2 类）：无响应则任务移入 `docs/plans/blocked/`，不得以默认值推进。

## 静默原则适用范围

正常进展不通知；仅上报：重要结果、异常、需要决策、里程碑、每日简报。
