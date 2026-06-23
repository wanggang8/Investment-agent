# Design: P2 领域规则

## Domain model

- 模型放在 `internal/domain/model/`。
- 枚举集中在 `enums.go`，每类枚举提供 `Valid() bool` 或 `Validate() error`。
- `WorkflowContext` 放在 `decision.go` 或独立结构中，字段对齐 `docs/workflow.md`。
- 模型只表达领域状态，不依赖 SQLite、HTTP、Eino 或前端 DTO。

## Rule engine

- 规则实现放在 `internal/domain/rule/`。
- 入口建议：`Evaluate(ctx model.WorkflowContext, input EvaluationInput) model.RuleVerdict`。
- 各策略拆分：
  - `capability_policy.go`：能力圈规则
  - `source_policy.go`：信源与多源验证规则
  - `risk_policy.go`：估值、现金、核心-卫星、移动止盈、情绪规则
  - `expectation_engine.go`：预期收益情景
  - `gatekeeper_logic.go`：规则提案状态机
  - `rules_engine.go`：组合裁决与优先级

## Priority

规则优先级：

1. 能力圈外 → `rejected`
2. 数据/证据不足 → `insufficient_data`
3. 多源验证失败或重大事件证据不足 → `frozen_watch`
4. 买入逻辑破坏 → `sell_only`
5. 高危估值 / 现金不足 / 仓位超限 → 禁止新增买入
6. 估值舒适区或低估区 → 只给分批、按计划类可选动作
7. 预期收益仅作为分析材料，不覆盖最终裁决

## Testing

- 模型枚举测试覆盖所有枚举合法/非法值。
- 规则测试按 `docs/development-plan.md` P2 场景表逐项覆盖。
- Gatekeeper 状态机测试覆盖送审、放弃、通过、拒绝、复核、最终确认、样本不足、终态重复操作。

## Constraints

- 不写数据库。
- 不调用 LLM。
- 不生成 HTTP response。
- 不做自动交易动作。
