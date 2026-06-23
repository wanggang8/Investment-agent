## 1. OpenSpec 与范围

- [x] 1.1 确认 P33 只覆盖本地账户/持仓录入、校准、线下交易记录、批量导入、错误修正和首次使用引导。
- [x] 1.2 确认 P33 不接券商 API、不自动交易、不外部推送、不登录/付费/授权/Level2/高频源、不承诺收益、不预测确定涨跌。
- [x] 1.3 对齐 `docs/api.md`、`docs/data-model.md`、`docs/frontend-contract.md`、`docs/workflow.md` 中现有 portfolio / confirmation / daily report 契约。

## 2. 后端模型与持久化

- [x] 2.1 盘点现有 `portfolio_snapshots`、`positions`、`position_snapshots`、`position_transactions`、`operation_confirmations`、`audit_events` 是否足够表达 P33；若不足，新增 migration delta。
- [x] 2.2 定义账户初始化、持仓编辑、线下交易流水、批量导入校验、确认写入和错误修正 DTO。
- [x] 2.3 实现统一校验：金额非负、数量非负、必填 symbol/name/cost/buy reason、交易日期合法、卖出不超过当前本地数量、总资产一致性容差。
- [x] 2.4 实现事务服务：初始化、持仓编辑/移除、线下买入/卖出/减少、批量导入确认、错误修正均在同一事务写账户、持仓、快照、流水和审计。
- [x] 2.5 增加 repository/service tests，覆盖成功写入、校验失败不写入、事务回滚、历史快照保留和审计事件。

## 3. HTTP API

- [x] 3.1 扩展或新增账户初始化 API，返回 snapshot、position count、audit ids 和缺前提状态变化。
- [x] 3.2 新增或扩展持仓新增/编辑/移除 API，要求用户确认和 reason，并返回新的本地事实引用。
- [x] 3.3 新增线下交易流水录入 API，支持 buy/sell/reduce，写 `position_transactions` 并更新本地快照。
- [x] 3.4 新增批量导入校验 API，返回逐行 validation results，不写事实。
- [x] 3.5 新增批量导入确认 API，仅在校验通过后单事务写入事实。
- [x] 3.6 新增错误修正 API，记录 before/after、correction reason 和 audit。
- [x] 3.7 增加 handler tests，覆盖成功、validation error、not found、conflict 和 non-trading safety text。

## 4. 前端产品体验

- [x] 4.1 新增账户初始化向导页面或 Portfolio 页内 onboarding 模式，覆盖现金、总资产、持仓、成本、买入原因、资产标签和风险偏好基础信息。
- [x] 4.2 增强 Portfolio 页，支持新增、编辑、移除当前持仓，并解释历史事实不会被物理删除。
- [x] 4.3 新增线下交易流水录入表单，文案明确“仅记录你已在线下完成的动作”。
- [x] 4.4 新增批量表格录入或 CSV 粘贴体验，展示逐行校验错误和确认写入按钮。
- [x] 4.5 新增错误修正流程，展示 before/after、修正原因和审计引用。
- [x] 4.6 更新 Dashboard 和今日纪律报告缺前提状态，账户/持仓缺失时引导到 P33 onboarding。
- [x] 4.7 增加前端 tests，覆盖空库引导、初始化成功、校验错误、线下交易记录、批量导入校验和禁止自动交易文案。

## 5. 文档与验收

- [x] 5.1 在 P33 delta 中记录待归档合并到 `docs/api.md` 的 API、请求/响应、错误码和事务边界。
- [x] 5.2 在 P33 delta 中记录待归档合并到 `docs/data-model.md` 的新增或复用数据模型、追加式历史和修正语义。
- [x] 5.3 在 P33 delta 中记录待归档合并到 `docs/frontend-contract.md` 的 onboarding、portfolio edit、offline transaction、batch import 和 correction view model。
- [x] 5.4 更新 `docs/development-plan.md` 和 `openspec/PROGRESS.md` 的 P33 状态与验收命令。
- [x] 5.5 运行 `go test ./...`。
- [x] 5.6 运行 `npm --prefix web test -- --run`。
- [x] 5.7 运行 `npm --prefix web run build`。
- [x] 5.8 运行 P33 定向 smoke 或扩展现有 E2E smoke，验证空库到完成账户初始化的浏览器路径。
- [x] 5.9 运行 `openspec validate p33-account-position-onboarding --strict`。
- [x] 5.10 运行 `openspec validate --all --strict`。
- [x] 5.11 运行 `git status --short`，确认只包含预期修改且无临时产物。
