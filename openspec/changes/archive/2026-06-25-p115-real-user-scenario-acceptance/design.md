# Design: P115 真实用户场景全链路验收

## Design Goal

P115 是一个验收型 change，不新增运行时投资能力。它要回答的问题是：用户按真实工作流使用产品时，每个关键动作是否真的联动了后端、数据、审计、跨页面状态和安全边界。

## Evidence Standard

每个验收场景至少记录以下证据：

1. `entry_surface`: 用户从哪个页面、按钮或状态进入。
2. `operation`: 用户执行了什么操作。
3. `api_or_browser_evidence`: HTTP request/response、Playwright trace、截图或 DOM readback。
4. `sqlite_readback`: 关键表和字段的 durable readback。
5. `downstream_readback`: Dashboard、Workbench、Review、Audit、Decision Loop、Data Quality、Risk、Rules 等下游读回。
6. `safety_negative_evidence`: 禁止 broker/order/push/auto-trade/auto-confirm/auto-rule-apply 的负证据。
7. `classification`: `fresh_pass`、`scoped_pass`、`degraded_expected`、`blocked`。

## Classification And Config Boundaries

P115 必须区分“预期可接受分类”和“实际执行状态”。矩阵中的分类只表示 eligibility，实际 status 必须由 runner 或人工验收记录写入，默认从 `pending` 开始。

P115 runner 至少区分两类配置证据：

1. `local_seeded_linkage`: 可使用 development mode、临时 SQLite、deterministic seed 或 P104 等价操作证明本地功能链路、API/SQLite/readback 和安全边界。此类证据不得声明真实外部 provider、真实 LLM 或未来数据可用性。
2. `provider_or_llm_path`: 只有在 `use_stub=false` 且 provider/LLM 调用实际成功时，才可记录外部 provider 或真实模型相关 `fresh_pass`。缺 key、网络失败、provider 失败或只有本地 deterministic evidence 时，必须记录 `scoped_pass` 或 `degraded_expected`。

## Artifact Schema

P115 summary JSON 每个 scenario entry 必须包含：

- `scenario_id`、`title`、`status`、`expected_eligibility`、`classification_reason`。
- `config_mode`、`runtime_mode`、`use_stub`、`provider_mode`、`llm_mode`。
- `api_evidence`: method、path、status、request_id、response excerpt 或 artifact path。
- `browser_evidence`: route、viewport、screenshot path、DOM assertions、console error count。
- `sqlite_evidence`: table、field、before、after、row_count、query label。
- `downstream_evidence`: endpoint/page、assertion、readback value。
- `side_effects`: audit ids、confirmation ids、risk ids、notification ids、rule proposal ids。
- `redaction_evidence`: secret/path/prompt/raw payload scan result。
- `safety_counters`: broker/order/push table count、auto confirmation rows、auto rule apply audit events、automatic trading affordances、return guarantee claims。

## Scenario Matrix

### S01 首次启动与本地能力边界

- 路由/入口：`/local-install`、`/settings`、`/api-diagnostics`、`/`。
- 操作：启动本地 backend/frontend，读取健康状态、系统设置、能力设置、诊断状态。
- 证据：`GET /api/v1/health`、`GET /api/v1/settings/system`、`GET /api/v1/settings/capability`。
- SQLite/readback：配置状态不写入敏感 key；诊断页面不泄露 raw secret。
- 安全边界：release mode 不允许 `use_stub=true`；页面不声明券商交易或自动下单。

### S02 空账户首次进入与引导

- 路由/入口：`/positions`、`/workbench`、`/`。
- 操作：使用空库进入产品，查看空组合状态、下一步动作和错误文案。
- 证据：浏览器截图、Dashboard/Workbench empty-state DOM。
- SQLite/readback：空库不产生虚假持仓。
- 安全边界：不把空账户伪装为已校准组合。

### S03 账户初始化与组合校准

- 路由/入口：`/positions`。
- 操作：初始化总资产、现金、核心/卫星/货币基金仓位。
- 证据：浏览器表单操作截图/DOM、`POST /api/v1/portfolio/init`、`GET /api/v1/portfolio/current`。
- SQLite/readback：`portfolio_snapshots`、`positions`。
- 下游联动：Dashboard 资产快照、Workbench 下一步动作、Review 组合摘要。
- 失败分支：总额不一致必须被拒绝，不得部分写入。

### S04 持仓新增、编辑、删除

- 路由/入口：`/positions`。
- 操作：新增持仓、编辑份额/成本/状态、删除一条持仓。
- 证据：浏览器持仓维护操作截图/DOM、`POST /api/v1/portfolio/holdings`、`POST /api/v1/portfolio/holdings/remove`。
- SQLite/readback：`positions`、`position_transactions`、审计字段。
- 下游联动：当前组合、Dashboard、Review、Audit。

### S05 批量导入持仓

- 路由/入口：`/positions`。
- 操作：上传或粘贴批量持仓，先 validate，再 confirm。
- 证据：浏览器导入操作截图/DOM、`POST /api/v1/portfolio/imports/validate`、`POST /api/v1/portfolio/imports/confirm`。
- SQLite/readback：import batch committed，positions 数量和金额变化。
- 失败分支：非法行、总额不一致、重复标的要有可理解错误。

### S06 线下交易记录

- 路由/入口：`/positions` 或决策确认后的执行记录。
- 操作：记录一次线下买入、卖出或调仓。
- 证据：`POST /api/v1/portfolio/offline-transactions`。
- SQLite/readback：`position_transactions`、最新 snapshot、现金变化。
- 安全边界：只记录线下交易，不出现 broker/order/push 表或自动下单。

### S07 本地事实修正与审计

- 路由/入口：`/positions`。
- 操作：修正持仓事实、记录修正原因。
- 证据：`POST /api/v1/portfolio/corrections`。
- SQLite/readback：correction rows、audit events。
- 下游联动：Audit 可查到修正原因，Review 读到新状态。

### S08 季度再平衡人工复核

- 路由/入口：`/positions`、`/review`。
- 操作：触发再平衡 review，查看建议和人工动作。
- 证据：`POST /api/v1/portfolio/rebalance-review`。
- SQLite/readback：rebalance review result、audit events。
- 安全边界：只给建议，不自动调整持仓。

### S09 主动咨询完整链路

- 路由/入口：`/consultation`。
- 操作：输入标的、假设、目标收益、上一轮基准，发起咨询。
- 证据：`POST /api/v1/decisions/consult`、决策详情截图。
- SQLite/readback：decision、analysis materials、expected return、evidence references。
- 下游联动：`/decisions` 列表、`/decision-loop`、`/evidence`。
- 降级分支：LLM key 缺失或质量失败时必须明确降级，不得冒充真实模型输出。

### S10 决策详情证据与规则解释

- 路由/入口：`/decisions/:decisionId`。
- 操作：查看最终裁决、证据、规则命中、可选动作、风险说明。
- 证据：详情页 DOM、`GET /api/v1/decisions/{id}`。
- SQLite/readback：decision facts、rule/evidence linkage。
- 安全边界：DeepSeek/LLM 不写最终裁决；最终裁决来自本地规则/工作流。

### S11 人工确认与执行记录

- 路由/入口：`/decisions/:decisionId`。
- 操作：确认采纳/拒绝/观察，填写线下执行结果。
- 证据：`POST /api/v1/decisions/{id}/confirmations`。
- SQLite/readback：`operation_confirmations`。
- 下游联动：Decision Loop、Audit、Dashboard 下一步动作。
- 安全边界：没有自动确认记录。

### S11B 决策错误标注与复盘闭环

- 路由/入口：`/decisions/:decisionId`、`/review`、`/audit`。
- 操作：将一条决策确认标记为 `marked_error`，填写 root cause / lesson learned。
- 证据：浏览器确认面板截图/DOM、`POST /api/v1/decisions/{id}/confirmations`。
- SQLite/readback：`operation_confirmations.confirmation_type=marked_error`、错误原因/经验字段、audit events。
- 下游联动：Review 和 Audit 能读到错误案例，不自动修改规则。

### S12 决策闭环追踪

- 路由/入口：`/decision-loop`。
- 操作：按决策筛选，查看输入、证据、裁决、确认、后续状态。
- 证据：`GET /api/v1/decision-loops`、`GET /api/v1/decision-loops/{id}`。
- SQLite/readback：decision + confirmation + audit 串联。
- 失败分支：不存在的 decision id 有明确错误或空态。

### S13 证据刷新、列表与验证

- 路由/入口：`/evidence`。
- 操作：刷新证据、查看证据列表、查看证据验证状态。
- 证据：浏览器证据刷新/列表截图、`POST /api/v1/evidence/refresh`、`GET /api/v1/evidence`、`GET /api/v1/evidence/verification`。
- SQLite/readback：evidence rows、source metadata、verification summary。
- 安全边界：数据源失败时不生成伪证据。

### S14 RAG/VecLite 重建与知识准备度

- 路由/入口：`/evidence`、`/local-knowledge`。
- 操作：重建索引，读取知识准备度。
- 证据：`POST /api/v1/evidence/rebuild-index`、`GET /api/v1/knowledge-readiness`。
- SQLite/readback：rag chunks、index status。
- 降级分支：索引不可用时咨询显示检索降级，不隐藏风险。

### S15 本地知识导入治理

- 路由/入口：`/local-knowledge`。
- 操作：validate 本地知识材料，查看脱敏预览，confirm 导入。
- 证据：浏览器导入/脱敏预览截图、`POST /api/v1/local-knowledge/imports/validate`、`POST /api/v1/local-knowledge/imports/confirm`。
- SQLite/readback：local knowledge facts、rag chunks。
- 安全边界：raw secret、路径、prompt payload 不在首层 UI 泄露。

### S16 市场数据刷新与 source health

- 路由/入口：`/data-quality`、`/settings`。
- 操作：刷新市场快照，查看 source health 和最新 snapshot。
- 证据：`POST /api/v1/market/refresh`、`GET /api/v1/market/source-health`、`GET /api/v1/market/snapshots/latest`。
- SQLite/readback：market snapshots、source health metadata。
- 分类：真实公开源 pass、网络/provider degraded、或 scoped local pass 必须分开记录。

### S17 数据质量回归与人工处置

- 路由/入口：`/data-quality`。
- 操作：查看回归结果、创建 gate resolution、退休 resolution。
- 证据：浏览器 resolution 创建/退休截图、`GET /api/v1/data-source-quality/regression`、`POST /api/v1/data-source-quality/resolutions`、`POST /api/v1/data-source-quality/resolutions/{id}/retire`。
- SQLite/readback：data quality resolutions、retired state。
- 下游联动：Dashboard/Workbench 数据质量状态同步。

### S18 风险预警 SOP 生命周期

- 路由/入口：`/risk-alerts`、`/risk-alerts/:alertId`。
- 操作：查看风险、进入详情、更新 SOP/lifecycle 状态。
- 证据：浏览器风险详情/生命周期操作截图、`GET /api/v1/risk-alerts`、`GET /api/v1/risk-alerts/{id}`、`POST /api/v1/risk-alerts/{id}/lifecycle`。
- SQLite/readback：risk alert state、audit events。
- 下游联动：Dashboard 风险摘要、Review 风险历史。

### S19 规则当前版本与提案确认

- 路由/入口：`/rules`。
- 操作：查看当前规则、查看提案、创建 SOP addendum、初步确认、最终确认。
- 证据：浏览器规则提案确认/最终确认截图、`GET /api/v1/rules/current`、`GET /api/v1/rule-proposals`、`POST /api/v1/rule-proposals/sop-addendum`、`POST /api/v1/rule-proposals/{id}/confirm`、`POST /api/v1/rule-proposals/{id}/final-confirm`。
- SQLite/readback：rule proposal state、rule version、audit events。
- 安全边界：没有自动规则应用。

### S20 规则效果验证与跟踪

- 路由/入口：`/rules`、`/review`。
- 操作：刷新规则效果验证，查看跟踪记录。
- 证据：`GET/POST /api/v1/rule-proposals/{id}/effect-validation`、`GET /api/v1/rule-effect-tracking`。
- SQLite/readback：validation result、tracking rows。
- 降级分支：样本不足必须显示“不足/需观察”，不能展示伪准确率。

### S21 通知中心

- 路由/入口：`/notifications`、顶部状态。
- 操作：读取通知、标记单条已读、全部已读。
- 证据：浏览器通知已读操作截图、`GET /api/v1/notifications`、`POST /api/v1/notifications/{id}/read`、`POST /api/v1/notifications/read-all`。
- SQLite/readback：notification read state。
- 下游联动：顶部未读数或通知摘要变化。

### S22 每日纪律报告

- 路由/入口：`/daily-discipline/reports`、`/daily-discipline/reports/:reportId`。
- 操作：查看今日报告、历史报告、报告详情。
- 证据：`GET /api/v1/daily-discipline/reports/today`、`GET /api/v1/daily-discipline/reports`、`GET /api/v1/daily-discipline/reports/{id}`。
- SQLite/readback：report rows、generated state。
- 降级分支：条件不足时必须冻结观察或提示待处理。

### S23 每日自动运行只读状态

- 路由/入口：`/daily-auto-run`。
- 操作：读取自动运行状态。
- 证据：`GET /api/v1/daily-auto-run/status`。
- SQLite/readback：如果有最近运行记录则读回，否则明确空态。
- 安全边界：自动运行默认关闭，不自动执行交易，不自动确认。

### S24 工作台与首页综合状态

- 路由/入口：`/`、`/workbench`。
- 操作：完成持仓、咨询、确认、风险、数据质量后回到首页/工作台。
- 证据：`GET /api/v1/dashboard/today`、浏览器截图。
- SQLite/readback：多模块状态汇总来自已有事实，不制造新事实。
- 下游联动：下一步人工动作、状态卡、资金快照、证据/规则摘要与前序操作一致。

### S25 复盘总览

- 路由/入口：`/review`。
- 操作：查看月度/季度复盘摘要。
- 证据：`GET /api/v1/review/summary`。
- SQLite/readback：读取决策、确认、风险、组合变化等事实。
- 降级分支：无历史数据时显示空态，不编造归因。

### S26 审计事件

- 路由/入口：`/audit`。
- 操作：查看所有关键操作审计轨迹。
- 证据：`GET /api/v1/audit-events`。
- SQLite/readback：audit event count、event type、subject id。
- 覆盖动作：组合初始化、持仓维护、确认、风险处置、数据质量处置、规则确认、本地知识导入。

### S27 设置变更

- 路由/入口：`/settings`。
- 操作：读取和更新系统设置、能力设置。
- 证据：浏览器设置操作截图、`GET/PUT /api/v1/settings`、`GET/PUT /api/v1/settings/capability`。
- SQLite/readback：settings values。
- 安全边界：敏感 key 脱敏；禁用能力不在 UI 假装可用。

### S28 API 诊断

- 路由/入口：`/api-diagnostics`。
- 操作：查看静态安全诊断导航；runner 另行执行 backend health API。
- 证据：浏览器静态诊断 DOM、独立 `GET /api/v1/health` API evidence。
- 安全边界：诊断信息产品化展示，不泄露本机敏感路径或 secret。

### S29 移动端真实路径

- 路由/入口：390px mobile viewport 覆盖核心路径。
- 操作：账户初始化、持仓编辑、主动咨询、决策确认、风险处置、数据质量处置、通知已读。
- 证据：Playwright mobile screenshots、DOM no-overflow/touch target checks。
- 联动：移动端操作与桌面端 API/SQLite 证据一致。

### S30 失败、降级和安全反证总场景

- 路由/入口：全产品。
- 操作：模拟非法输入、无持仓、无证据、数据源失败、LLM key 缺失、索引不可用、不存在 id。
- 证据：错误响应、UI error/empty/degraded state、SQLite 未写入或写入错误审计。
- 安全边界：broker/order/push 表为 0，auto confirmation 为 0，auto rule apply audit 为 0，自动交易相关 UI affordance 为 0。

### S31 Settings 禁止规则/SOP 直接修改

- 路由/入口：`/settings`、`/rules`。
- 操作：直接向 settings API 提交 `rule_thresholds` 或 `sop_config` 变更。
- 证据：`PUT /api/v1/settings` 返回 400 或等价拒绝；浏览器设置页不暴露直接规则/SOP 修改入口。
- SQLite/readback：规则版本、SOP 配置和 rule proposal 不被修改。
- 安全边界：没有自动规则应用审计事件。

### S32 Local Install 诊断摘要与脱敏

- 路由/入口：`/local-install`。
- 操作：查看本地安装诊断摘要、配置状态和安全提示。
- 证据：浏览器诊断摘要 DOM；如存在诊断 API/脚本输出，则记录脱敏后的 artifact。
- 安全边界：不泄露真实 key、本机敏感路径、raw stack、prompt payload。

### S33 Browser-level 交互覆盖补齐

- 路由/入口：`/positions`、`/evidence`、`/local-knowledge`、`/data-quality`、`/risk-alerts`、`/rules`、`/notifications`、`/settings`。
- 操作：每个有真实交互控件的页面至少执行一个 browser-level 操作或 DOM 状态断言。
- 证据：桌面截图、390px 关键路径截图、console error count、对应 API/SQLite parity。
- 安全边界：browser UI 不出现自动交易、自动确认、自动规则应用、收益承诺或敏感信息泄露。

## Runner Strategy

P115 runner 分三层执行：

1. `api_sqlite_layer`: 通过 HTTP API 执行业务操作，直接查 SQLite 字段级 readback。覆盖 S03-S23、S26-S31。
2. `browser_layer`: 通过 Playwright 操作真实 UI，采集截图/DOM/console。覆盖 S01-S05、S09-S19、S21-S29、S32-S33。
3. `degradation_layer`: 使用隔离配置或输入模拟失败/降级，确认 UI/API 不伪造结果。覆盖 S05、S09、S13-S17、S20、S22-S23、S28、S30-S31。

每层必须输出独立 log/summary：`api_sqlite/`、`browser/`、`degradation/`。P115 可复用 P104 runner 的临时 backend、SQLite、seed 和 forbidden-table checks，但 P104-derived evidence 必须标记为 `local_seeded_linkage`，不得用于 provider/LLM fresh claim。

## Artifact Layout

- `docs/release/acceptance/2026-06-25-p115-real-user-scenario-acceptance-matrix.md`
- `docs/release/acceptance/2026-06-25-p115-real-user-scenario-acceptance.md`
- `docs/release/ui-audit-assets/2026-06-25-p115-real-user-scenario-acceptance/`
- `scripts/p115-real-user-scenario-acceptance.sh`
- `scripts/p115_real_user_scenario_acceptance.py`

## Pass Criteria

- 所有 S01-S33 均有明确 status 和证据路径。
- 核心本地功能场景必须 `fresh_pass`。
- 外部 provider/LLM/network 相关场景允许 `scoped_pass` 或 `degraded_expected`，但必须显示边界，不能扩大声明。
- 任何“自动交易、自动确认、自动规则应用、券商下单、外部推送、收益承诺”正向证据均为 release-blocking failure。
