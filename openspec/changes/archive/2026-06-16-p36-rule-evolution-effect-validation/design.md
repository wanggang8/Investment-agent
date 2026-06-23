## Context

当前系统已有规则提案状态机、守门人审计、复盘摘要、错误案例、每日纪律、风险预警和应用内通知。已有能力能生成提案并通过用户确认与守门人审计进入正式规则，但对“提案为什么可信、样本是否代表真实问题、是否过拟合、应用后是否改善”仍缺少结构化事实。

P36 在本地事实库内完成效果验证与追踪，服务对象是规则治理和复盘页面。它不改变自动交易边界，不绕过守门人审计，不把模型输出作为最终裁决。

## Goals / Non-Goals

**Goals:**

- 为每个可评估规则提案生成来源解释、样本统计、代表性判断、过拟合风险、历史回放指标和风险说明。
- 将验证结果纳入守门人审计门禁：样本不足、代表性不足、过拟合风险高或回放结果不利时不得进入自动应用路径。
- 追踪已应用规则版本在后续复盘中的命中率、误判率、缺证据率、降级情况和相关风险预警变化。
- 提供本地 API 与前端展示，展示验证与追踪结果，但不提供自动应用规则或交易动作。
- 在 P36 delta 中记录归档后应同步到 L1 文档的 API、数据模型、工作流和前端契约。

**Non-Goals:**

- 不实现自动规则应用；所有正式规则变更仍需守门人审计和用户最终确认。
- 不接券商 API、不生成订单、不修改账户持仓、不执行外部推送。
- 不引入付费、登录、授权、Level2 或高频数据源。
- 不做确定性收益承诺或涨跌预测。
- 不用 LLM 覆盖规则裁决或守门人审计结论。

## Decisions

### 1. 使用本地事实驱动的 EffectValidationService

新增应用服务聚合 `rule_proposals`、`error_cases`、`decision_records`、`operation_confirmations`、`daily_discipline_reports`、`risk_alerts`、`audit_events` 等本地事实，输出规则提案验证结果。这样能复用现有 SQLite 事实，不需要外部服务，也能保持测试可控。

备选方案是只在前端动态计算，但会让门禁和审计缺少统一事实来源；因此不采用。

### 2. 验证结果持久化为本地事实

P36 应新增 `rule_effect_validations` 或等价事实存储，保存 proposal_id、rule_version、sample_count、sample_window、representativeness_status、overfit_risk、replay_result、metrics_json、risk_notes_json、created_at、updated_at 等字段。应用后追踪可使用同表的 tracking 类型或单独 `rule_effect_tracking` 表。

这样可以让 API、复盘和审计读取同一事实，也便于归档和回放。若实现时发现现有 `rule_proposals` JSON 字段足够，可作为内部实现简化，但契约仍要求验证结果可持久化、可审计。

### 3. 守门人审计只消费验证结果，不自动放行

守门人审计应读取验证结果并把样本不足、代表性不足、过拟合风险和历史回放不利作为拒绝或需要用户复核的依据。即使验证结果通过，提案仍只能进入 `pending_final_confirm`，等待用户最终确认。

这延续现有安全状态机，避免 P36 变成自动规则进化。

### 4. 应用后追踪并入复盘输出

已应用规则的效果趋势应在 review summary 和规则治理页面中展示，包括命中率、误判率、缺证据率、降级率、相关 risk_alert 数量与趋势。追踪只读展示，不创建新 active rule version；若发现问题，只能生成新的提案或复盘提示。

### 5. 前端采用摘要卡 + 详情展开

规则提案详情显示验证摘要，复盘页显示趋势摘要。复杂 metrics JSON 默认折叠，用户优先看到状态、原因、样本和建议动作。所有动作文案必须说明“需要审计与用户最终确认”。

## Risks / Trade-offs

- 样本量较少导致验证结论不稳定 → 使用 `insufficient` 或 `needs_more_samples` 状态，不让其进入通过结论。
- 历史回放指标可能受数据缺失影响 → 显式记录 missing evidence、degraded data 和 source health，避免伪造完整回测。
- 新表和 API 增加复杂度 → 只保存必要指标和 JSON 快照，避免重建完整回测系统。
- 前端指标过多影响理解 → 摘要优先，详情折叠。
- 验证结果被误解为收益承诺 → 文案固定声明仅作规则治理参考，不构成收益预测，不自动交易。

## Migration Plan

1. 新增 SQLite migration 和 repository，空库直接创建表；旧库迁移不改变现有规则提案状态。
2. 为历史 `applied` 或 `pending_final_confirm` 提案提供按需生成验证/追踪能力；不强制回填全部历史。
3. API 在无验证结果时返回空状态或 `not_evaluated`，前端展示待评估。
4. 归档时将 P36 delta 合并到 L1 文档，并更新进度文件。