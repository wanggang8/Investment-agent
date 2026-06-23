## Why

P32 已把每日纪律结果产品化为今日报告、历史报告和详情回看，但真实用户从空库开始仍缺少完整账户与持仓录入路径。没有可靠的账户、持仓、成本、现金和线下交易事实，后续每日纪律、风险预警、全路径 E2E 都只能依赖 seed 或人工写库。

P33 将补齐账户初始化、持仓维护、线下交易流水、一致性校验、错误修正和首次使用引导，使用户能在本地 UI 中完成从空库到可生成每日纪律报告的基础数据准备。

## What Changes

- 新增账户初始化向导，录入现金、总资产、持仓、成本、买入原因、资产标签和风险偏好基础信息。
- 增强持仓新增、编辑、删除前置确认和校验；历史事实继续通过追加式快照和审计记录表达。
- 支持线下交易流水录入后更新本地持仓、现金、账户快照和审计，且只记录用户已在线下完成的动作。
- 支持批量导入或表格化录入持仓与历史交易，并提供逐行校验结果。
- 增加录入错误修正流程，避免静默覆盖历史数据。
- 前端展示首次使用引导，明确距离可运行每日纪律还缺哪些账户、持仓、行情或证据前提。
- 不新增券商连接、自动下单、外部推送、收益承诺或确定性涨跌预测。

## Capabilities

### New Capabilities
- `account-position-onboarding`: 覆盖账户初始化、持仓录入/维护、线下交易流水、本地一致性校验、错误修正和首次使用引导。

### Modified Capabilities
- `daily-discipline-report`: 今日报告缺前提状态应能引用 P33 的账户/持仓初始化引导入口。

## Impact

- 影响后端：portfolio handlers、repositories、service/query、transaction handling、audit events、可能新增 import validation DTO。
- 影响数据模型：如现有表不足，需通过 delta 明确新增或扩展本地录入/修正记录字段；不得破坏追加式历史。
- 影响前端：Portfolio page、Dashboard empty/onboarding state、API service/types、表单与批量录入 UI。
- 影响测试：repository/service/handler tests、frontend tests、P39 前置 smoke seed 或用户旅程准备。
- 影响文档：`docs/api.md`、`docs/data-model.md`、`docs/frontend-contract.md`、`docs/development-plan.md`、`openspec/PROGRESS.md`。
