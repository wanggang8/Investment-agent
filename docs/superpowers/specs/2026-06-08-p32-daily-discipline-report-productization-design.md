# P32 每日纪律报告产品化设计

## 目标

P32 将 Daily workflow 从“可触发任务”产品化为用户每天打开即可阅读、追踪和回看的“每日纪律报告”。P32 建立报告聚合 API、轻量报告索引、今日报告页面、历史报告列表和详情入口，让 P31 每日自动运行产物进入稳定的用户阅读面。

## 范围选择

采用轻量产品化方案：新增 Daily Report 聚合模型和轻量索引，但复用现有 `decision_records` 作为正式建议事实源。

P32 做：

- 今日报告状态与摘要。
- 最近报告历史列表。
- 报告详情视图。
- 自动运行、手动运行和失败状态统一到同一报告视图。
- 最小趋势摘要：最近成功、降级、失败、证据不足次数。
- 空库初始化引导，不生成正式建议。

P32 不做：

- 账户初始化、持仓编辑或交易流水录入，这属于 P33。
- 真实数据覆盖扩展，这属于 P34。
- 风险预警 SOP 中心，这属于 P35。
- 规则进化效果验证，这属于 P36。
- 外部推送、券商接口或自动执行交易。

## 数据模型

新增轻量索引表 `daily_discipline_reports`，只保存每日阅读入口和历史索引，不复制完整决策内容。

核心字段：

- `report_id`：稳定 ID，例如 `daily_report:2026-06-08:holdings:<symbol_hash>:v1`。
- `local_date`：本地日期。
- `scope`：初期固定 `holdings`。
- `symbol_set_hash`：与 P31 一致。
- `source_type`：`auto_run` 或 `manual`。
- `source_id`：daily auto-run idempotency key 或手动 request id。
- `decision_id`：关联 daily discipline 生成的 `decision_records`。
- `status`：`not_started`、`running`、`success`、`degraded`、`failed`、`insufficient_data`。
- `summary`：报告可读摘要。
- `failure_code`、`failure_reason`：失败或缺前提诊断。
- `created_at`、`updated_at`。

幂等规则：同一 `local_date + scope + symbol_set_hash` 复用同一报告索引。需要重跑版本化时以后再单独立项，不在 P32 扩大。

## API 设计

新增每日纪律报告 API：

- `GET /api/v1/daily-discipline/reports/today`
  - 返回今日是否完成、状态、摘要、关联决策、关联自动运行、缺失项、风险/证据/规则摘要和安全说明。
- `GET /api/v1/daily-discipline/reports`
  - 返回最近报告列表，默认 30 条，支持状态筛选和 limit 上限。
- `GET /api/v1/daily-discipline/reports/{report_id}`
  - 返回报告详情，包含决策、证据、审计、通知和每日自动运行状态链接。

P32 不新增手动触发 POST API。手动运行继续沿用现有本地任务入口，避免扩大行为边界。

## 聚合逻辑

新增 `DailyDisciplineReportService` 或在现有 `QueryService` 中实现聚合方法。职责是把现有事实拼成每日报告视图：

1. 优先读取今日 `daily_discipline_reports` 索引。
2. 若没有报告索引，查询最新 `daily_auto_run_states`。
3. 若今日 auto-run 失败，返回 `failed` 或 `insufficient_data` 报告视图，并展示失败原因。
4. 若没有运行记录，返回 `not_started` 和初始化/运行提示。
5. 若报告关联 `decision_id`：
   - 从 `decision_records` 读取最终裁决、触发规则、禁止动作、可选动作、expected return 和 source verification 状态。
   - 从 `evidence_refs` 聚合证据数量、独立信源数、高等级信源数和证据摘要。
   - 从 `audit_events` 聚合运行步骤、失败步骤和诊断摘要。

状态映射：

- daily auto-run success + decision completed → `success`。
- daily auto-run degraded 或 workflow degraded → `degraded`。
- missing prerequisites/no data → `insufficient_data`。
- timeout/source failure → `failed`。
- running state → `running`。

`decision_records` 仍是正式建议事实源；报告只是每日阅读、追踪和回看的聚合视图。

## 前端设计

升级当前 `/` 今日纪律页为报告入口，而不是新增割裂的首页。

今日纪律页展示：

- 今日报告状态卡片：已完成、未开始、运行中、失败、数据不足。
- 组合摘要：资产、现金、持仓数量、仓位结构。
- 今日裁决：最终裁决、禁止动作、可选动作、是否需人工确认。
- 风险与规则：触发规则、红线说明、降级原因。
- 证据覆盖：证据数、独立信源、高等级信源、缺失项。
- 追踪链接：决策详情、审计详情、通知、每日自动运行状态。

新增页面：

- `/daily-discipline/reports`
  - 最近 30 次报告列表。
  - 支持 success/degraded/failed/insufficient_data 简单筛选。
  - 空列表展示初始化引导。
- `/daily-discipline/reports/:reportId`
  - 报告详情页。
  - 复用现有决策详情组件和证据摘要组件，避免重复 UI。

趋势最小版：today/detail API 返回最近 N 次状态摘要，包括最近成功、降级、失败、证据不足次数。前端用已有图表或简单列表展示。复杂趋势留给后续阶段。

## 空库体验

空库或无持仓时，今日报告页显示：

- 尚未完成账户/持仓初始化。
- 每日纪律报告不会生成正式建议。
- 引导入口：持仓页、每日自动运行状态、配置说明。
- 明确说明系统不会自动执行交易。

空库状态不创建 `decision_records`。

## 测试与验收

后端测试：

- `daily_discipline_reports` migration/repository：upsert、幂等读取、latest/list/detail。
- 聚合服务：空库、P31 failed auto-run、success decision、degraded/failed audit、同日重复聚合。
- handler：today、list、detail、limit 校验、未知 report 404。

前端测试：

- 今日纪律页 success、insufficient_data、failed 三类状态。
- 历史报告列表、空列表、状态标签、详情链接。
- 报告详情页关联决策/审计/通知链接。
- service 层 API 路径。
- 不出现自动执行交易或确定收益文案。

E2E smoke：

- smoke seed 写入一个 P32 daily report index。
- 验证 `/` 展示“今日纪律报告”。
- 验证 `/daily-discipline/reports` 展示历史报告。
- 验证报告详情能跳转到现有决策详情或审计详情。
- 验证页面保留非自动执行交易和人工复核边界。

## OpenSpec change

创建 `p32-daily-discipline-report-productization`：

- proposal：说明 P32 是 P31 后的报告产品化层。
- tasks：模型/API、聚合逻辑、前端页面、E2E smoke、文档同步、验收。
- specs：新增 `daily-discipline-report` spec，或扩展相邻 spec。

## 安全边界

P32 不新增交易执行、券商接口、外部推送、登录源、付费源或高频抓取。报告中的行动项只允许人工复核、查看证据、补齐数据、记录线下计划。报告不得承诺收益、不得预测确定涨跌、不得覆盖规则裁决。
