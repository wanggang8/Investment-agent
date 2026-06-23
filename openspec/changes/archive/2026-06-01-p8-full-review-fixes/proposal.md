## Why

P0–P8 全仓复审发现若干跨层契约偏差：决策确认、规则提案、账户快照、DTO 展示、治理进度与测试门禁仍有不一致。需要在进入 P9 前统一修复，避免后续自动化交付建立在不完整状态上。

## What Changes

- 修复决策记录类型、确认状态、详情接口与用户确认链路的一致性。
- 修复主动咨询能力圈判断、规则提案用户确认到守门人审计与最终确认链路。
- 修复手动执行后的现金、总资产、持仓与快照一致性。
- 补齐决策详情、市场快照、证据列表、前端状态映射、错误态与空态展示。
- 同步 P6/P7/P8 治理进度、测试门禁与 OpenSpec 摘要能力。
- 补充覆盖上述契约的后端、前端与文档验证。

## Capabilities

### New Capabilities
- `p0-p8-review-fixes`: 覆盖 P0–P8 全仓复审后的修复行为与验收标准。

### Modified Capabilities
- `real-data-integration`: 修正真实数据、证据、市场与规则提案链路中的一致性要求。
- `e2e-hardening`: 扩展端到端验证、事务回滚与质量门禁要求。
- `frontend-experience-tests`: 扩展前端中文映射、错误态、空态与关键交互测试要求。

## Impact

- 后端：workflow、handler、service、repository、SQLite migration/model/DTO、测试。
- 前端：dashboard、decision detail、portfolio、evidence、audit、rules、settings、services、types、测试。
- 文档治理：development plan、testing plan、openspec specs、progress/agent 指引。
