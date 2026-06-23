## Why

P34 已补齐真实公开数据覆盖、freshness 与失败分类，但这些信号仍主要停留在数据与报告展示层。P35 需要把规则裁决、数据健康、持仓状态和每日纪律结果组织成可追踪的风险预警与 SOP 状态流转，让用户能看到风险为何触发、当前处于什么处理阶段，以及哪些动作被禁止。

## What Changes

- 新增本地风险预警与 SOP 编排能力，覆盖估值高位、买入逻辑破坏、流动性不足、情绪极端、仓位超限、证据不足和数据降级。
- 新增风险预警事实模型、仓储、服务和 API，支持触发、持续观察、升级、解除和归档。
- 将风险预警写入应用内通知、审计事件和每日纪律报告关联上下文。
- 前端新增风险预警中心与报告/驾驶舱风险摘要，展示触发证据、当前状态、建议人工动作和禁止动作。
- 保持本地只读与人工复核边界；不新增交易执行、券商 API、外部推送或自动规则应用。

## Capabilities

### New Capabilities
- `risk-alert-sop`: 定义风险预警、SOP 状态流转、触发依据、解除/归档和非交易边界。

### Modified Capabilities
- `daily-discipline-report`: 每日纪律报告需要展示关联风险预警摘要和 SOP 状态。
- `real-data-integration`: 工作流与前端状态需要消费 P34 数据 freshness/source health 作为风险预警输入，而不把缺失数据误当完整证据。

## Impact

- 后端：新增风险预警领域模型、SQLite migration、repository/service、HTTP handler、DTO、审计与通知写入路径。
- 工作流：DailyDisciplineGraph / RuleArbitration 结果后增加风险预警编排服务，复用 market、portfolio、decision、source health、evidence 状态。
- 前端：新增风险预警中心页面，扩展 Dashboard、每日纪律报告详情、通知/审计跳转展示。
- 文档：通过本 change 的 delta 记录待归档合并内容；不直接修改 L1 契约正文作为实现真源。
- 验收：Go tests、前端 tests/build、OpenSpec strict validation、风险场景 smoke、非交易边界检查。
