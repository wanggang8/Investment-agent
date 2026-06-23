# Investment Agent 功能需求拆分

> 文档版本：v1.0  
> 最后更新：2026-05-27  
> 适用范围：功能拆分、验收标准、开发计划前置输入。  
> 上游文档：`docs/requirements.md`、`docs/architecture.md`。  
> 契约文档：`docs/data-model.md`、`docs/api.md`、`docs/workflow.md`、`docs/frontend-contract.md`、`docs/ui-design.md`、`docs/ui-flow.md`。

## 1. 文档目标

本文档把产品需求拆成可开发、可测试、可验收的功能单元。每个功能单元明确：

- 用户场景。
- 前置条件。
- 用户动作。
- 系统行为。
- 数据写入。
- API / DTO。
- UI 展示。
- 异常状态。
- 验收标准。
- 禁止行为。

本文档不替代 `docs/data-model.md`、`docs/api.md` 和 `docs/workflow.md`。字段、枚举、状态流转以这些契约文档为准。

## 2. 全局边界

### 2.1 系统定位

Investment Agent 是个人投资纪律辅助系统。它帮助用户理解当前状态、证据、规则和可选动作，不替用户交易，不预测未来涨跌。

### 2.2 全局禁止行为

| 禁止行为 | 说明 |
| --- | --- |
| 自动交易 | 不接券商交易 API，不提供一键买卖，不代替用户执行。 |
| 主动推荐标的 | 只分析用户主动提出且在能力圈内的标的。 |
| 承诺收益 | 收益相关内容只能作为概率和情景分析。 |
| 使用 C 级信源作为正式裁决依据 | C 级信源只能作为 `background` 背景材料展示。 |
| 大模型生成最终裁决 | DeepSeek 只生成分析材料，最终裁决由规则引擎给出。 |
| 自动应用规则提案 | 守门人审计通过后仍需用户最终确认。 |

### 2.3 全局数据写入规则

`decision_records` 必须区分记录类型：

| record_type | 使用场景 | 是否允许用户确认 | 必须写入 | 禁止行为 |
| --- | --- | --- | --- | --- |
| `formal_trade_advice` | 能力圈内、数据充足、生成正式交易类建议 | 是，`confirmation_status=pending` | `decision_records`、`evidence_refs`、`audit_events` | 不得缺少规则版本、账户快照、正式证据引用 |
| `rejection_record` | 能力圈外、规则明确拒绝分析 | 否，`confirmation_status=not_required` | `decision_records`、`audit_events` | 不得展示交易确认动作 |
| `non_trade_record` | 证据不足、数据过期、VecLite 与 SQLite 均不足等信息不足场景 | 否，`confirmation_status=not_required` | 可写 `decision_records`、必须写 `audit_events` | 不得展示交易类建议 |

首次使用、账户缺失、持仓缺失时不得创建 `decision_records`，只返回 `DATA_REQUIRED` 或 `dashboard_state=first_use`，并可写 `audit_events`。

| 事件 | 必须写入 | 说明 |
| --- | --- | --- |
| 生成正式建议 | `decision_records(record_type=formal_trade_advice)`、`evidence_refs`、`audit_events` | 保存建议、证据引用和执行审计。 |
| 用户记录计划 | `operation_confirmations`、`audit_events` | 不更新账户。 |
| 用户记录已手动执行 | `operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 只记录用户已在线下完成的动作。 |
| 用户标记待观察 | `operation_confirmations`、`audit_events` | 不更新账户。 |
| 用户标记错误 | `operation_confirmations`、`error_cases`、`audit_events` | 生成错误案例。 |
| 情报入库 | `intelligence_items`、`intelligence_summary`、`rag_chunks` | VecLite 可由 `rag_chunks` 重建。 |
| 多源验证 | `source_verifications`、`audit_events` | 保存验证结果。 |
| 生成规则提案 | `rule_proposals`、`audit_events` | 不修改正式规则。 |
| 守门人审计 | `gatekeeper_audits`、`audit_events` | 不修改正式规则。 |
| 用户最终确认规则 | `rule_versions`、`rule_proposals`、`audit_events` | 提案状态变为 `applied`。 |

### 2.4 全局状态枚举

| 类型 | 可选值 | 来源 |
| --- | --- | --- |
| 页面状态 | `first_use` / `normal` / `insufficient_data` / `frozen_watch` / `high_risk` | `docs/frontend-contract.md` |
| 工作流状态 | `completed` / `degraded` / `failed` | `docs/api.md` |
| 记录类型 | `formal_trade_advice` / `non_trade_record` / `rejection_record` | `docs/data-model.md` |
| 持仓状态 | `normal` / `sell_only` / `frozen_watch` | `docs/api.md` |
| 验证状态 | `satisfied` / `failed` / `background_only` | `docs/api.md` |
| 用户确认状态 | `not_required` / `pending` / `planned` / `executed_manually` / `watch` / `marked_error` | `docs/api.md` |
| 用户确认动作 | `planned` / `executed_manually` / `watch` / `marked_error` | `docs/api.md` |
| 最终裁决状态 | `buy_allowed` / `hold` / `reduce` / `sell_only` / `frozen_watch` / `rejected` / `insufficient_data` | `docs/api.md` |
| 规则提案状态 | `draft` / `pending_user_confirm` / `under_gatekeeper_audit` / `pending_final_confirm` / `rejected` / `applied` | `docs/api.md` |
| 守门人审计结果 | `approved` / `rejected` / `needs_user_review` | `docs/api.md` |

## 3. 功能模块总览

| 编号 | 模块 | 核心目标 |
| --- | --- | --- |
| F01 | 首次使用与基础设置 | 让系统获得账户、持仓、能力圈和默认规则。 |
| F02 | 今日纪律驾驶舱 | 用户每日查看纪律状态、风险、证据和建议。 |
| F03 | 持仓与账户管理 | 维护当前账户状态和持仓状态。 |
| F04 | 决策咨询 | 用户主动询问标的或场景，系统生成结构化裁决。 |
| F05 | 证据与情报管理 | 采集、清洗、验证、检索和展示证据。 |
| F06 | 用户确认与线下动作记录 | 记录用户计划、已手动执行、待观察和错误。 |
| F07 | 错误案例与规则进化 | 从错误案例生成规则提案，并受控应用。 |
| F08 | 规则与纪律管理 | 展示规则版本、阈值、SOP 和能力圈。 |
| F09 | 审计事件与复盘 | 展示建议、动作、节点、规则变更的审计链路。 |
| F10 | 系统配置、市场数据与数据维护 | 管理本地配置、数据源、索引和数据刷新状态。 |

---

## F01. 首次使用与基础设置

### F01.1 首次进入系统

**用户场景**  
用户首次打开系统，还没有录入账户、持仓或能力圈配置。

**前置条件**

- `portfolio_snapshots` 为空，或 `positions` 为空。
- `capability_configs` 为空，或未启用。
- `rule_versions` 至少有默认 active 版本 `v3.0`。如果缺失，系统进入高危状态。

**用户动作**

- 打开今日纪律页。

**系统行为**

1. 前端请求 `GET /api/v1/dashboard/today`。
2. 后端检查账户、持仓、能力圈和规则版本。
3. 缺少账户或持仓时返回错误码 `DATA_REQUIRED`，或返回 `dashboard_state=first_use`。
4. 前端展示首次使用引导。

**数据写入**

- 仅查看首页时不写业务表。
- 如果系统生成检查事件，可写入 `audit_events`，`action=generate_decision`，`status=degraded`。

**API / DTO**

- `GET /api/v1/dashboard/today`
- 错误码：`DATA_REQUIRED`
- 页面状态：`first_use`

**UI 展示**

- 主文案：需要先录入账户和持仓，系统才能生成纪律报告。
- 主要动作：录入持仓、配置能力圈。
- 不展示交易建议。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 规则版本缺失 | 返回 `RULE_VERSION_MISSING`，页面状态为 `high_risk`。 |
| 配置文件读取失败 | 返回 `INTERNAL_ERROR`，展示系统配置异常。 |

**验收标准**

- 无账户数据时，页面不展示买入、卖出、持有建议。
- 前端能显示缺失项：账户、持仓、能力圈。
- 后端不会创建 `decision_records` 正式建议。

**禁止行为**

- 不允许用默认虚拟账户生成正式建议。
- 不允许跳过能力圈配置进入交易类分析。

### F01.2 录入初始账户与持仓

**用户场景**  
用户录入现金、总资产、持仓代码、名称、数量、成本价、买入理由和资产标签。

**前置条件**

- 用户处于 `first_use`。
- 已存在 active `rule_versions`。

**用户动作**

- 在持仓或首次使用引导页提交账户与持仓表单。

**系统行为**

1. 校验现金、总资产、持仓数量、成本价为非负数。
2. 校验每条持仓必须有 `symbol`、`name`、`quantity`、`cost_price`、`buy_reason`。
3. 创建 `portfolio_snapshots`。
4. 创建或更新 `positions` 当前态。
5. 创建完整 `position_snapshots` 集合。
6. 写入 `audit_events`。

**数据写入**

- `portfolio_snapshots`
- `positions`
- `position_snapshots`
- `audit_events`

**API / DTO**

- `POST /api/v1/portfolio/init` 已在 `docs/api.md` 定义。
- 响应必须包含 `request_id`、`snapshot_id`、`position_count`。

**UI 展示**

- 表单必须说明：这里只记录本地账户状态，不连接券商账户。
- 提交后返回今日纪律页。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 总资产小于现金 | 返回 `BAD_REQUEST`。 |
| 持仓缺少买入理由 | 返回字段级错误。 |
| 写入持仓成功但快照失败 | 整个事务回滚。 |

**验收标准**

- 提交后能查到 1 条 `portfolio_snapshots`。
- 每条当前持仓都有对应 `position_snapshots`。
- 审计事件记录用户录入动作。

**禁止行为**

- 不自动从券商同步真实账户。
- 不把首次录入解释为买入动作。

---

## F02. 今日纪律驾驶舱

### F02.1 查看今日纪律状态

**用户场景**  
用户每天打开系统，查看当前账户、市场、风险红线、证据和今日建议。

**前置条件**

- 已有账户快照和持仓快照。
- 已有 active `rule_versions`。
- 市场数据可用或能返回明确过期状态。

**用户动作**

- 打开今日纪律页。
- 点击刷新今日状态。

**系统行为**

1. 执行 DailyDisciplineGraph。
2. StateSnapshotNode 读取账户、持仓、市场和规则版本。
3. EvidenceRetrievalNode 读取 `intelligence_summary`、`rag_chunks`、VecLite 检索结果和 `source_verifications`。
4. ValueAnalystNode 和 TrendRiskOfficerNode 调用 DeepSeek 生成分析材料。
5. RuleArbitrationNode 根据规则裁决最终建议。
6. DecisionRecordNode 写入建议、证据引用和审计事件。
7. 返回驾驶舱 DTO。

**数据写入**

- 成功生成正式建议：`decision_records(record_type=formal_trade_advice)`、`evidence_refs`、`audit_events`。
- 信息不足且需要详情页复现：可写 `decision_records(record_type=non_trade_record)` 和 `audit_events`，但 `confirmation_status=not_required`，不得展示交易确认动作。
- 信息不足且不需要详情页复现：只写 `audit_events`。

**API / DTO**

- `GET /api/v1/dashboard/today`
- 返回字段：`dashboard_state`、`discipline_status`、`portfolio_summary`、`market_summary`、`triggered_rules`、`decision_summary`、`evidence_summary`。

**UI 展示**

首屏顺序：

1. 当前纪律状态。
2. 风险红线和禁止事项。
3. 今日建议。
4. 账户摘要。
5. 证据摘要。

**异常状态**

| 情况 | 错误码 / 状态 | UI 行为 |
| --- | --- | --- |
| 账户缺失 | `DATA_REQUIRED` / `first_use` | 引导录入。 |
| 行情过期 | `DATA_STALE` / `insufficient_data` | 展示更新时间和暂停原因。 |
| 证据不足 | `EVIDENCE_NOT_FOUND` / `insufficient_data` | 暂停交易类建议。 |
| VecLite 不可用且 SQLite 摘要不足 | `VECTOR_INDEX_UNAVAILABLE` / `insufficient_data` | 展示索引不可用。 |
| VecLite 不可用但 SQLite 摘要充足 | `workflow_status=degraded` | 可降级展示证据摘要。 |
| 多源验证失败 | `SOURCE_VERIFICATION_FAILED` / `frozen_watch` | 展示等待更多证据。 |
| 规则版本缺失 | `RULE_VERSION_MISSING` / `high_risk` | 暂停裁决。 |

**验收标准**

- 正常状态下可看到规则、证据、裁决和建议编号。
- 信息不足时不会生成交易类建议。
- DeepSeek 不可用时，系统可返回规则裁决降级结果，`workflow_status=degraded`。
- 所有节点状态写入 `audit_events.node_action`。

**禁止行为**

- 不以收益榜、涨跌刺激作为首页主信息。
- 不把聊天输入框作为主界面中心。
- 不展示一键交易按钮。

### F02.2 高危与冻结观察展示

**用户场景**  
系统发现买入逻辑破坏、多源验证不足、情绪极端或流动性风险。

**前置条件**

- 已有持仓和市场快照。
- 已有相关证据或缺失证据结论。

**用户动作**

- 查看今日纪律页或决策详情页。

**系统行为**

- 买入逻辑破坏：裁决为 `sell_only`，持仓状态可进入 `sell_only`。
- 多源验证不足：裁决为 `frozen_watch`。
- 情绪极端：暂停主动交易建议。
- 流动性风险：禁止市价式大额操作，只允许提示分批、限价或暂停。
- PE/PB 分位规则：任一核心估值指标高于 80% 时禁止新增买入；50%-80% 为观察区，只允许持有或观察；30%-50% 为舒适区，可在其他规则满足时按计划定投；低于 30% 为低估区，可在买入逻辑完好且情绪非极端时提示分批配置。
- 移动止盈规则：浮盈达到 20% 时可提示卖出 30% 并启动移动止盈；浮盈达到 30% 且 20% 阶段已处理时可提示再卖出 30%；剩余仓位从阶段高点回撤 10% 时触发减仓或卖出评估。
- 现金冗余规则：`cash_ratio<0.05` 时限制新增买入；`0.05<=cash_ratio<=0.10` 为正常冗余区间，不因现金规则单独禁止交易。
- 核心-卫星规则：核心资产目标 60%-70%，卫星资产目标 20%-30%；卫星资产超过 30% 或偏离目标超过 15% 时优先提示再平衡，止盈资金优先回归核心资产。

**数据写入**

- `decision_records.risk_reason_code`
- `evidence_refs`
- `audit_events`

**API / DTO**

- `final_verdict.status`
- `prohibited_actions`
- `optional_actions`
- `triggered_rules`

**UI 展示**

- 高危信息优先于分析观点。
- 明确展示禁止事项。
- 展示触发规则和证据来源。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 风险原因无法结构化 | 返回 `INTERNAL_ERROR`，不展示正式建议。 |
| 证据为 C 级 | 只能展示为背景，不进入正式裁决。 |

**验收标准**

- `sell_only` 状态下不出现新增买入建议。
- `frozen_watch` 状态下展示等待条件。
- 风险状态可追溯到规则和证据。
- PE/PB 分位、移动止盈、现金冗余、核心-卫星仓位命中的规则必须出现在 `triggered_rules`。
- 任一核心估值指标高危时，`prohibited_actions` 必须包含新增买入。
- 卫星资产超上限时，`optional_actions` 必须优先展示再平衡，而不是新增卫星资产。

**禁止行为**

- 不把分析师观点放在最终裁决之前。
- 不用模糊文案暗示用户立即交易。

---

## F03. 持仓与账户管理

### F03.1 查看当前持仓

**用户场景**  
用户查看当前本地账户状态、持仓、成本、仓位比例、买入理由和状态。

**前置条件**

- 已录入持仓。

**用户动作**

- 打开持仓页。

**系统行为**

1. 读取 `positions` 当前态。
2. 读取最新 `portfolio_snapshots`。
3. 返回持仓列表和账户摘要。

**数据写入**

- 只读，不写业务表。
- 可写访问审计，但不是必须。

**API / DTO**

- `GET /api/v1/portfolio/current` 已在 `docs/api.md` 定义。

**UI 展示**

- 展示 symbol、name、quantity、cost_price、current_price、market_value、position_state、buy_reason。
- `buy_reason` 不得默认折叠到完全不可见。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 没有持仓 | 展示首次使用引导。 |
| 最新价格缺失 | 展示价格缺失，不自动估算。 |

**验收标准**

- 持仓页展示当前态，不用于复现历史建议。
- 决策详情使用的是 `position_snapshots`，不是当前 `positions`。

**禁止行为**

- 不允许在持仓页提供真实下单入口。
- 不允许把编辑持仓误写成交易流水。

### F03.2 手动校准账户状态

**用户场景**  
用户发现本地账户状态与真实账户不一致，需要手动校准。

**前置条件**

- 已有当前账户和持仓。

**用户动作**

- 在持仓页编辑现金、持仓数量、成本价或买入理由。

**系统行为**

1. 校验输入。
2. 更新 `positions` 当前态。
3. 创建新的 `portfolio_snapshots`。
4. 创建完整 `position_snapshots`。
5. 写 `audit_events`。

**数据写入**

- `positions`
- `portfolio_snapshots`
- `position_snapshots`
- `audit_events`

**API / DTO**

- `POST /api/v1/portfolio/adjustments` 已在 `docs/api.md` 定义。
- 请求必须包含调整原因。

**UI 展示**

- 文案必须说明：这是本地账户校准，不代表实际交易。
- 调整原因必填。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 调整后总资产为负 | 返回 `BAD_REQUEST`。 |
| 调整原因为空 | 返回字段级错误。 |
| 快照写入失败 | 事务回滚。 |

**验收标准**

- 校准后最新快照包含完整持仓集合。
- 校准不会创建 `position_transactions`。
- 审计事件能区分账户校准与已手动执行。

**禁止行为**

- 不把账户校准当成系统建议执行结果。
- 不生成买卖建议。

---

## F04. 决策咨询

### F04.1 用户主动咨询标的

**用户场景**  
用户输入“这个标的还能不能持有？”或“是否值得继续研究？”。

**前置条件**

- 用户已配置能力圈。
- 已有账户状态。
- 标的在能力圈内，或系统能明确判断能力圈外。

**用户动作**

- 在决策咨询页提交问题、标的代码和场景。

**系统行为**

1. 调用 `POST /api/v1/decisions/consult`。
2. 执行 ConsultationGraph。
3. CapabilityCheckNode 判断能力圈。
4. 能力圈外时返回 `final_verdict.status=rejected`，只输出研究前置问题。
5. 能力圈内时继续证据检索、分析师观点、预期收益情景评估、规则裁决。
6. 预期收益评估基于历史类似样本、估值分位和市场快照生成上行/基准/下行情景。
7. 样本数量决定 `precision_status`：`>=20` 为 `available`，`5~19` 为 `insufficient`，`<5` 或市场快照缺失为 `unavailable`。
8. 写入正式建议和审计事件。

**数据写入**

- 能力圈内且生成正式建议：`decision_records(record_type=formal_trade_advice)`、`evidence_refs`、`audit_events`。
- 能力圈外：必须写入 `decision_records(record_type=rejection_record)` 和 `audit_events`，用于复盘能力圈边界；不得写入交易类建议或可执行动作。
- 证据不足 / 数据不足：不生成正式交易建议；若需要详情页复现，可写入 `decision_records(record_type=non_trade_record)`，`final_verdict.status=insufficient_data`，且 `confirmation_status=not_required`。

**API / DTO**

- `POST /api/v1/decisions/consult`
- `GET /api/v1/decisions/{decision_id}`
- DTO 包含：`capability_check`、`evidence_chain`、`analyst_reports`、`expected_return_scenarios`、`arbitration_chain`、`final_verdict`、`user_confirmation`。
- `final_verdict.status=rejected` 或 `insufficient_data` 时允许返回 `decision_id`，但 `user_confirmation.confirmation_status` 必须为 `not_required`，前端不得展示交易确认动作。

**UI 展示**

展示顺序：

1. 用户问题。
2. 能力圈检查。
3. 信息核查。
4. Agent 观点。
5. 预期收益情景。
6. 规则裁决。
7. 最终建议。
8. 用户确认区。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 问题为空 | `BAD_REQUEST`。 |
| 标的为空且场景需要标的 | `BAD_REQUEST`。 |
| 能力圈外 | `final_verdict.status=rejected`。 |
| 证据不足 | `final_verdict.status=insufficient_data`。 |
| LLM 不可用 | `workflow_status=degraded`，规则裁决可继续。 |
| 历史类似样本少于 20 个 | 不输出精确概率，只展示样本不足。 |
| 历史类似样本少于 5 个 | 不生成收益区间，只输出定性风险提示。 |

**验收标准**

- 能力圈外不会进入交易类建议。
- Agent 观点不等于最终裁决。
- 决策详情能复现账户快照、市场快照、证据、规则和审计链。
- `decision_records` 必须保存或等价引用当次 `portfolio_snapshot_id`、`market_snapshot_id`、`rule_version`；信息不足导致无有效市场快照时，需在 `errors_json` 中说明原因。
- 预期收益只作为情景分析展示，不覆盖最终规则裁决。

**禁止行为**

- 不主动推荐用户未提出的标的。
- 不输出确定性涨跌判断。
- 不承诺收益，不把预期收益作为买卖保证。

### F04.2 查看历史决策详情

**用户场景**  
用户重新打开某条历史建议，查看当时依据。

**前置条件**

- 存在 `decision_records`。

**用户动作**

- 在复盘页或驾驶舱点击建议详情。

**系统行为**

1. 根据 `decision_id` 读取决策记录。
2. 读取关联账户快照、持仓快照、市场快照、证据引用和审计事件。
3. 返回详情 DTO。

**数据写入**

- 默认只读。

**API / DTO**

- `GET /api/v1/decisions/{decision_id}`

**UI 展示**

- 展示当时快照，不使用当前持仓覆盖历史数据。
- 展示证据 hash、来源等级和发布时间。
- 展示规则版本。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 决策不存在 | `NOT_FOUND`。 |
| 证据原文已清理 | 展示摘要、来源和 hash。 |

**验收标准**

- 历史详情不受当前持仓变化影响。
- 能看到 `rule_version`。

**禁止行为**

- 不重新生成历史裁决覆盖旧记录。

---

## F05. 证据与情报管理

### F05.1 情报采集与入库

**用户场景**  
系统定时或手动刷新公告、新闻、研报摘要，用于后续决策。

**前置条件**

- 数据源已配置。
- SQLite 可写。

**用户动作**

- 用户手动刷新证据，或本地任务触发刷新。

**系统行为**

1. 从 S/A/B 级信源获取内容。
2. C 级内容可以保留为背景材料，但不得进入正式裁决。
3. 标准化来源、发布时间、原文 URL、hash。
4. 生成结构化摘要。
5. 写入 RAG 文本块。
6. 更新 VecLite 索引。

**数据写入**

- `intelligence_items`
- `intelligence_summary`
- `rag_chunks`
- `source_verifications`
- VecLite 索引文件
- `audit_events`

**API / DTO**

- `POST /api/v1/evidence/refresh` 已在 `docs/api.md` 定义。
- 查询使用：`GET /api/v1/evidence`。

**UI 展示**

- 展示信源等级、来源名称、发布时间、摘要、hash。
- C 级背景材料弱化展示，并标明不可作为裁决依据。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 来源缺少 URL 或发布时间 | 标记为不可用，不进入正式证据。 |
| VecLite 写入失败 | 保留 SQLite 数据，标记索引不可用。 |
| 内容重复 | 根据 hash 去重。 |

**验收标准**

- 同一内容重复刷新不会重复创建正式证据。
- VecLite 不可用时可从 SQLite 重建。
- 同步刷新接口成功时必须写入 `source_verifications`，或在同一请求中明确返回验证失败原因。

**禁止行为**

- 不保存无法追溯来源的正式证据。
- 不让 C 级信源参与正式裁决。

### F05.2 多源验证

**用户场景**  
涉及重大利好、重大利空或买入逻辑破坏时，系统需要验证是否有足够独立信源。

**前置条件**

- 已有 `intelligence_summary`。
- 至少有一个相关标的或事件。

**用户动作**

- 系统在决策前自动触发，或用户从证据页查看验证结果。

**系统行为**

1. 聚合同一事件下的独立来源。
2. 统计 S/A/B 来源数量，并区分是否为重大事件。
3. 普通正式证据允许 S/A/B 级来源进入 `formal` 证据链，但 C 级只能作为 `background`。
4. 重大利好、重大利空、买入逻辑破坏等重大事件必须满足至少 2 个 A 或 S 级独立信源；不满足时 `verification_status=failed`。
5. 写入 `source_verifications`。
6. 决策时读取验证结果。

**数据写入**

- `source_verifications`
- `audit_events`

**API / DTO**

- `GET /api/v1/evidence`
- DTO 包含 `verification_status`、`evidence_role`。

**UI 展示**

- `satisfied`：可作为正式证据。
- `failed`：进入冻结观察。
- `background_only`：只作为背景材料。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 独立来源不足 | `SOURCE_VERIFICATION_FAILED`。 |
| 只有 C 级来源 | `background_only`。 |
| 来源冲突 | 展示冲突说明，并进入冻结观察。 |

**验收标准**

- 买入逻辑破坏、重大利好或重大利空必须能追溯到至少 2 个 A 或 S 级独立信源，或明确进入冻结观察。
- 验证结果可被决策详情引用。

**禁止行为**

- 不用单一来源确认重大事件。

---

## F06. 用户确认与线下动作记录

### F06.0 用户确认状态转移矩阵

`confirmation_type` 是用户本次提交动作，`confirmation_status` 是 `decision_records` 当前确认状态。每次提交成功都创建新的 `operation_confirmations`；`decision_records.confirmation_status` 保存最近一次有效状态。

| 当前 confirmation_status | 允许动作 | 目标 confirmation_status | 是否创建 confirmation | 是否更新账户 | 说明 |
| --- | --- | --- | --- | --- | --- |
| `not_required` | 无 | 不变 | 否 | 否 | 非交易型记录、拒绝型记录不得确认。 |
| `pending` | `planned` | `planned` | 是 | 否 | 记录计划，不代表成交。 |
| `pending` | `executed_manually` | `executed_manually` | 是 | 是 | 用户已在线下执行，必须填写成交字段。 |
| `pending` | `watch` | `watch` | 是 | 否 | 标记待观察。 |
| `pending` | `marked_error` | `marked_error` | 是 | 否 | 写错误案例。 |
| `planned` | `executed_manually` | `executed_manually` | 是 | 是 | 允许从计划升级为已手动执行。 |
| `planned` | `watch` | `watch` | 是 | 否 | 允许放入观察。 |
| `planned` | `marked_error` | `marked_error` | 是 | 否 | 计划被认为错误时写错误案例。 |
| `watch` | `planned` | `planned` | 是 | 否 | 观察后重新记录计划。 |
| `watch` | `executed_manually` | `executed_manually` | 是 | 是 | 观察后用户在线下执行。 |
| `watch` | `marked_error` | `marked_error` | 是 | 否 | 观察后标记错误。 |
| `executed_manually` | 任意确认动作 | 不变，返回 `BAD_REQUEST` | 否 | 否 | 已手动执行是确认终态，不允许二次确认避免重复账户变更。 |
| `marked_error` | 任意确认动作 | 不变，返回 `BAD_REQUEST` | 否 | 否 | 已标记错误是确认终态。 |

约束：只有 `record_type=formal_trade_advice` 且 `confirmation_status` 不是 `not_required` 时，前端才展示确认区。所有确认动作必须写 `audit_events(action=confirm_operation)`；`marked_error` 额外写 `audit_events(action=mark_error)`。

### F06.1 记录计划

**用户场景**  
用户看到建议后，只想记录计划，不改变账户。

**前置条件**

- 存在 `decision_records`。
- 当前记录必须满足 `record_type=formal_trade_advice`，且状态转移必须符合 F06.0 矩阵。
- 建议允许用户确认。

**用户动作**

- 在确认区选择“记录计划”。

**系统行为**

1. 写入 `operation_confirmations`，`confirmation_type=planned`。
2. 更新建议确认状态为 `planned`。
3. 写入 `audit_events`。
4. 不更新账户和持仓。

**数据写入**

- `operation_confirmations`
- `audit_events`

**API / DTO**

- `POST /api/v1/decisions/{decision_id}/confirmations`

**UI 展示**

- 文案：只记录计划，不改变本地账户。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 决策不存在 | `NOT_FOUND`。 |
| 已标记错误后再次计划 | 返回 `BAD_REQUEST` 或要求用户先查看历史确认。 |

**验收标准**

- 不新增 `position_transactions`。
- 不新增 `portfolio_snapshots`。

**禁止行为**

- 不把计划记录当成实际成交。

### F06.2 记录已手动执行

**用户场景**  
用户已经在线下完成交易，需要记录结果以校准本地账户。

**前置条件**

- 存在 `decision_records`。
- 当前记录必须满足 `record_type=formal_trade_advice`，且状态转移必须符合 F06.0 矩阵。
- 用户在线下实际完成交易。

**用户动作**

- 选择“已手动执行”，填写 operation_type、symbol、quantity、price、executed_at。

**系统行为**

1. 校验 `operation_type` 只能是 `buy / sell / reduce`。
2. 校验 `quantity`、`price`、`executed_at` 必填。
3. 写入 `operation_confirmations`，`confirmation_type=executed_manually`。
4. 写入 `position_transactions`。
5. 更新 `positions` 当前态。
6. 创建新的 `portfolio_snapshots`。
7. 创建完整 `position_snapshots`。
8. 写入 `audit_events`。

**数据写入**

- `operation_confirmations`
- `position_transactions`
- `positions`
- `portfolio_snapshots`
- `position_snapshots`
- `audit_events`

**API / DTO**

- `POST /api/v1/decisions/{decision_id}/confirmations`

**UI 展示**

- 明确文案：仅记录你已在线下完成的交易，系统不会执行交易。
- 必填成交数量、价格和时间。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 缺少成交数量 | `BAD_REQUEST`。 |
| 卖出数量超过当前持仓 | `BAD_REQUEST`。 |
| 任一写入失败 | 整个事务回滚。 |

**验收标准**

- 事务成功后，6 类数据全部可查。
- 事务失败时，不留下部分交易流水。
- 审计事件包含前后状态。

**禁止行为**

- 不调用任何外部交易接口。
- 不允许 `watch` 作为成交类型。

### F06.3 标记待观察

**用户场景**  
用户暂不决定执行结果，希望后续跟踪。

**前置条件**

- 存在 `decision_records`。
- 当前记录必须满足 `record_type=formal_trade_advice`，且状态转移必须符合 F06.0 矩阵。

**用户动作**

- 选择“待观察”。

**系统行为**

- 写入 `operation_confirmations`，`confirmation_type=watch`。
- 写入 `audit_events`。
- 不更新账户。

**数据写入**

- `operation_confirmations`
- `audit_events`

**API / DTO**

- `POST /api/v1/decisions/{decision_id}/confirmations`

**UI 展示**

- 显示观察原因和下次复查时间输入。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 观察原因为空 | 允许提交，但 UI 提示建议填写。 |

**验收标准**

- 不写 `position_transactions`。
- 决策状态显示为待观察。

**禁止行为**

- 不把待观察解释为持有建议。

### F06.4 标记错误

**用户场景**  
用户认为某条建议事后证明有问题，标记错误并记录原因。

**前置条件**

- 存在历史 `decision_records`。
- 当前记录必须满足 `record_type=formal_trade_advice`，且状态转移必须符合 F06.0 矩阵。

**用户动作**

- 选择“标记错误”，填写实际结果、错误原因、经验记录。

**系统行为**

1. 写入 `operation_confirmations`，`confirmation_type=marked_error`。
2. 写入 `error_cases`。
3. 将 `error_case_id` 回写到确认记录。
4. 写入 `audit_events`。
5. 不直接生成或应用正式规则。

**数据写入**

- `operation_confirmations`
- `error_cases`
- `audit_events`

**API / DTO**

- `POST /api/v1/decisions/{decision_id}/confirmations`
- 响应必须返回 `error_case_id`。

**UI 展示**

- 要求填写错误原因。
- 展示该错误可能进入规则提案，但不会自动修改规则。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 缺少错误原因 | `BAD_REQUEST`。 |
| 错误案例写入失败 | 事务回滚。 |

**验收标准**

- `marked_error` 必须同一事务写入三类数据。
- 返回体包含 `error_case_id`。

**禁止行为**

- 不因单个错误案例自动改规则。

---

## F07. 错误案例与规则进化

### F07.1 生成规则提案

**用户场景**  
系统从多个错误案例中发现模式，生成规则优化提案。

**前置条件**

- 存在 `error_cases`。
- 样本数量达到提案所需条件；少于 3 个时必须标记样本不足。

**用户动作**

- 用户在复盘页查看系统生成的规则提案。
- 或本地任务周期性生成提案。

**系统行为**

1. 聚合错误案例。
2. 提取模式。
3. 生成 `rule_proposals`，初始状态为 `draft` 或 `pending_user_confirm`。
4. 写入 `audit_events`。
5. 不修改 `rule_versions`。

**数据写入**

- `rule_proposals`
- `audit_events`

**API / DTO**

- `GET /api/v1/rule-proposals`

**UI 展示**

- 展示提案来源、关联错误案例、变更前后文本、样本数量和风险说明。
- 样本不足必须明确展示。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 样本不足 | 允许生成提案，状态只能是 `draft` 或 `pending_user_confirm`，并在 `risk_notes_json` 写入样本不足说明；除非后续 EvolutionGraph 或受控内部任务生成满足样本条件的新提案版本，否则不得进入 `pending_final_confirm` 或 `applied`。 |
| 提案与根本规则冲突 | 标记为高风险，等待守门人审计。 |

**验收标准**

- 生成提案不会创建新的 active `rule_versions`。
- 提案能追溯到错误案例。

**禁止行为**

- 不让大模型直接改正式规则。

**规则提案状态机**

| 当前状态 | 动作 / 审计结果 | 目标状态 | 必须写入 | 禁止写入 |
| --- | --- | --- | --- | --- |
| 无 | EvolutionGraph 生成提案且样本不足 | `draft` | `rule_proposals`、`audit_events` | `rule_versions`、`pending_final_confirm`、`applied` |
| 无 | EvolutionGraph 生成提案且满足送审前条件 | `pending_user_confirm` | `rule_proposals`、`audit_events` | `rule_versions` |
| `draft` | 系统或用户补齐提案说明后提交给用户确认 | `pending_user_confirm` | `rule_proposals`、`audit_events` | `rule_versions` |
| `pending_user_confirm` | 用户 `confirm=false` 放弃送审 | `rejected` | `rule_proposals`、`audit_events` | `gatekeeper_audits`、`rule_versions` |
| `pending_user_confirm` | 用户 `confirm=true` 且 `sample_count>=3` | `under_gatekeeper_audit` | `rule_proposals`、`audit_events` | `rule_versions` |
| `pending_user_confirm` | 用户 `confirm=true` 且 `sample_count<3` | 不变，返回 `BAD_REQUEST` | `audit_events` | `gatekeeper_audits`、`rule_versions` |
| `under_gatekeeper_audit` | `audit_result=approved` | `pending_final_confirm` | `gatekeeper_audits`、`rule_proposals`、`audit_events` | `rule_versions` |
| `under_gatekeeper_audit` | `audit_result=rejected` | `rejected` | `gatekeeper_audits`、`rule_proposals`、`audit_events` | `rule_versions` |
| `under_gatekeeper_audit` | `audit_result=needs_user_review` | `pending_user_confirm` | `gatekeeper_audits`、`rule_proposals`、`audit_events` | `rule_versions` |
| `pending_final_confirm` | 用户 `confirm=true` 且 `sample_count>=3` | `applied` | `rule_versions`、`rule_proposals`、`audit_events` | - |
| `pending_final_confirm` | 用户 `confirm=true` 且 `sample_count<3` | 不变，返回 `BAD_REQUEST` | `audit_events` | `rule_versions` |
| `pending_final_confirm` | 用户 `confirm=false` | `rejected` | `rule_proposals`、`audit_events` | `rule_versions` |
| `rejected` / `applied` | 任意确认动作 | 不变，返回 `BAD_REQUEST` | `audit_events` | `rule_versions` |

说明：`approved / rejected / needs_user_review` 只属于 `gatekeeper_audits.audit_result`，不是 `rule_proposal.status`。当前版本不提供用户编辑提案接口；`needs_user_review` 后用户可放弃或重新送审；如需修改，只能由后续 EvolutionGraph 或受控内部任务生成新提案版本。
### F07.2 用户确认提案并进入守门人审计

**用户场景**  
用户认为某条提案值得审计，确认进入守门人审计。

**前置条件**

- 提案状态为 `pending_user_confirm`。
- `sample_count>=3`；样本不足提案不得进入守门人审计。

**用户动作**

- 点击确认提案。

**系统行为**

1. 状态从 `pending_user_confirm` 变为 `under_gatekeeper_audit`。
2. 执行 GatekeeperAuditGraph。
3. 写入 `gatekeeper_audits`。
4. 审计结果为 `approved` 时，提案状态变为 `pending_final_confirm`。
5. 审计结果为 `rejected` 时，提案状态变为 `rejected`。
6. 审计结果为 `needs_user_review` 时，提案回到 `pending_user_confirm`；当前版本不提供前端手动编辑提案接口，如需修改，只能由后续 EvolutionGraph 或受控内部任务生成新提案版本，用户也可放弃或重新送审。
7. 写入 `audit_events`。

**数据写入**

- `rule_proposals`
- `gatekeeper_audits`
- `audit_events`

**API / DTO**

- `POST /api/v1/rule-proposals/{proposal_id}/confirm`

**UI 展示**

- 展示审计结果：通过、否决、需要用户复核。
- 审计通过后显示“等待最终确认”。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 提案状态不允许确认 | `BAD_REQUEST`。 |
| 样本不足 | `BAD_REQUEST`，提案状态不变，不写 `gatekeeper_audits` 或 `rule_versions`。 |
| 守门人审计失败 | `workflow_status=failed`，保留原状态并写审计事件。 |

**验收标准**

- 审计通过不会写 `rule_versions`。
- `approved` 只作为 `gatekeeper_audits.audit_result`，不作为提案状态。

**禁止行为**

- 不在守门人审计通过后自动应用规则。

### F07.3 用户最终确认应用规则

**用户场景**  
守门人审计通过后，用户最终确认是否应用规则。

**前置条件**

- 提案状态为 `pending_final_confirm`。
- 最近一次 `gatekeeper_audits.audit_result=approved`。

**用户动作**

- 点击最终确认应用，或拒绝应用。

**系统行为**

- `confirm=true`：当 `sample_count>=3` 时写入新的 `rule_versions`，提案状态变为 `applied`；当 `sample_count<3` 时返回 `BAD_REQUEST`，不得写 `rule_versions`。
- `confirm=false`：提案状态变为 `rejected`，不写 `rule_versions`。
- 两种情况都写入 `audit_events`。

**数据写入**

- `rule_versions`，仅 `confirm=true`。
- `rule_proposals`
- `audit_events`

**API / DTO**

- `POST /api/v1/rule-proposals/{proposal_id}/final-confirm`

**UI 展示**

- 必须展示变更前后内容。
- 必须展示守门人审计摘要。
- 最终确认按钮文案应为“确认应用到正式规则”，不能暗示自动优化。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 提案不是 `pending_final_confirm` | `BAD_REQUEST`。 |
| 样本不足 | `BAD_REQUEST`，不得写 `rule_versions`，提案状态不变。 |
| 规则版本写入失败 | 事务回滚，提案保持原状态。 |

**验收标准**

- 应用后同一时间只有一个 active `rule_versions`。
- 提案记录 `applied_rule_version`。
- 审计事件记录前后状态。

**禁止行为**

- 不允许绕过最终确认。

---

## F08. 规则与纪律管理

### F08.1 查看当前规则版本

**用户场景**  
用户查看当前生效的根本规则、阈值、SOP 和大师智慧映射。

**前置条件**

- 存在 active `rule_versions`。

**用户动作**

- 打开规则与纪律页。

**系统行为**

- 读取 active `rule_versions`。
- 返回规则版本、规则内容、阈值和生效时间。

**数据写入**

- 只读。

**API / DTO**

- `GET /api/v1/rules/current` 已在 `docs/api.md` 定义。

**UI 展示**

- 展示当前规则版本，如 `v3.0`。
- 展示根本规则、阈值、SOP。
- 展示规则来源和生效时间。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| active 规则缺失 | `RULE_VERSION_MISSING`。 |
| 多个 active 版本 | `INTERNAL_ERROR`，进入高危状态。 |

**验收标准**

- 每条建议详情显示当时使用的规则版本。
- 规则页显示当前 active 版本。

**禁止行为**

- 不允许直接编辑 active 规则并立即生效。

### F08.2 管理能力圈

**用户场景**  
用户配置系统允许分析的资产类型、标的、排除项和策略范围。

**前置条件**

- 用户已进入设置或规则页。

**用户动作**

- 新增、修改能力圈配置。

**系统行为**

1. 保存 `capability_configs`。
2. 写入 `audit_events`。
3. 后续 ConsultationGraph 必须读取能力圈配置。

**数据写入**

- `capability_configs`
- `audit_events`

**API / DTO**

- `GET /api/v1/settings/capability`、`PUT /api/v1/settings/capability` 已在 `docs/api.md` 定义。

**UI 展示**

- 展示资产类型范围、纳入标的、排除标的、策略范围。
- 修改时提示：能力圈会影响系统是否拒绝交易类分析。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 纳入和排除同一标的 | 返回 `BAD_REQUEST`。 |
| 配置为空 | 允许保存，但咨询时会进入能力圈未知或拒绝状态。 |

**验收标准**

- 能力圈外咨询返回 `rejected`。
- 能力圈配置变更写审计事件。

**禁止行为**

- 不允许能力圈外仍输出交易类建议。

---

## F09. 审计事件与复盘

### F09.0 审计字段语义

审计字段语义固定如下：

| 字段 | 含义 | 必填条件 | 示例 / 约束 |
| --- | --- | --- | --- |
| `action` | 业务动作，描述用户或系统在做什么 | 必填 | generate_decision、confirm_operation、mark_error、create_proposal、audit_rule_change、update_rule、refresh_market_data、update_settings、update_capability、rebuild_index |
| `node_name` | Eino 节点名称，描述由哪个节点执行 | 工作流节点事件必填；纯用户动作可为空 | StateSnapshotNode、RuleArbitrationNode、MarketSnapshotRecordNode |
| `node_action` | 节点动作枚举，描述节点内标准动作 | 工作流节点事件必填；纯用户动作可为空 | load_state_snapshot、arbitrate_rule、refresh_market_data |
| `status` | 本条审计事件状态 | 必填 | success / degraded / failed；不等同于 `workflow_status` |
| `actor` | 操作者 | 必填 | system / user / gatekeeper |
| `decision_id` | 关联决策 | 与建议、确认、错误案例有关时必填 | 非决策类事件可为空 |
| `proposal_id` | 关联规则提案 | 提案、守门人、规则应用事件必填 | 非提案事件为空 |
| `confirmation_id` | 关联确认记录 | 用户确认事件必填 | 非确认事件为空 |
| `error_code` | 错误码 | `status=failed` 时必填，降级时按需填写 | DATA_STALE、MARKET_SNAPSHOT_WRITE_FAILED |
| `input_ref_type` / `input_ref` | 输入引用类型和摘要 | 工作流节点事件必填 | workflow_context、decision_record、rule_proposal、market_refresh_request |
| `output_ref_type` / `output_ref` | 输出引用类型和记录 ID / 摘要 | 产生输出时必填 | decision_record、market_snapshot、failed_symbols、rule_version |
| `before_state` / `after_state` | 状态变更前后 | 状态机、设置变更、规则变更、能力圈变更必填 | 不得包含密钥明文 |

敏感字段处理：DeepSeek API Key、数据源密钥、完整本地文件路径不得写入 `input_ref`、`output_ref`、`before_state`、`after_state`；如需记录，只保存脱敏摘要。

### F09.1 查看审计事件

**用户场景**  
用户查看某条建议、确认动作、规则提案或工作流节点的审计链路。

**前置条件**

- 存在 `audit_events`。

**用户动作**

- 打开复盘与审计页。
- 按决策、规则提案、时间筛选。

**系统行为**

- 查询 `audit_events`。
- 支持按 `decision_id`、`proposal_id`、`request_id`、时间范围筛选。
- 返回业务动作和节点动作。

**数据写入**

- 只读。

**API / DTO**

- `GET /api/v1/audit-events`

**UI 展示**

- `action` 展示业务动作，如 generate_decision、confirm_operation、mark_error、refresh_market_data、update_settings、update_capability、rebuild_index。
- `node_name` 展示工作流节点名称，如 StateSnapshotNode、RuleArbitrationNode。
- `node_action` 展示节点动作枚举，如 load_state_snapshot、arbitrate_rule。
- 输入输出引用默认折叠，支持展开。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 没有审计事件 | 展示空状态。 |
| 查询参数非法 | `BAD_REQUEST`。 |

**验收标准**

- 审计页能区分 `action`、`node_name` 与 `node_action`。
- 每条正式建议至少有一条审计事件。

**禁止行为**

- 不物理删除审计事件。

### F09.2 月度和季度复盘

**用户场景**  
用户查看纪律执行情况、错误案例、规则变化和操作确认历史。

**前置条件**

- 存在决策记录或用户确认记录。

**用户动作**

- 打开复盘页，选择时间范围。

**系统行为**

- 汇总 `decision_records`、`operation_confirmations`、`error_cases`、`rule_proposals`、`audit_events`。
- 展示纪律执行、错误分布、规则提案状态。

**数据写入**

- 默认只读。

**API / DTO**

- `GET /api/v1/review/summary` 已在 `docs/api.md` 定义。

**UI 展示**

- 展示建议数量、确认动作数量、错误案例数量、规则提案数量。
- 展示已手动执行与计划记录的区别。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 时间范围无数据 | 展示空状态，不生成虚构总结。 |

**验收标准**

- 复盘数据可点击回到原始建议或审计事件。
- 不把未执行计划计入实际交易。

**禁止行为**

- 不基于复盘自动推荐新标的。

---

## F10. 系统配置、市场数据与数据维护

### F10.1 本地配置管理

**用户场景**  
用户或开发者配置 SQLite、VecLite、DeepSeek、数据源和日志级别。

**前置条件**

- 本地配置文件或环境变量可用。

**用户动作**

- 修改配置文件或环境变量。
- 打开设置页查看状态。

**系统行为**

- 读取配置。
- 校验 SQLite 路径、VecLite 路径、DeepSeek Key、数据源开关。
- 配置异常时返回系统状态错误。
- 设置更新按三类入口处理：通知、页面偏好、普通数据源通过 `PUT /api/v1/settings` 写 `user_settings`；能力圈通过 `PUT /api/v1/settings/capability` 写 `capability_configs`；根本规则、裁决优先级、核心阈值和 SOP 必须生成 `rule_proposals`。

**数据写入**

- 可写 `user_settings`。
- 可写 `audit_events` 记录配置变更。

**API / DTO**

- `GET /api/v1/settings/system`、`PUT /api/v1/settings` 已在 `docs/api.md` 定义。

**UI 展示**

- 展示本地数据库状态、索引状态、数据源状态、模型连接状态。
- 不展示完整 API Key。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| SQLite 路径不可写 | 系统不可用，提示修复配置。 |
| VecLite 路径不可写 | 证据检索降级或信息不足。 |
| DeepSeek Key 缺失 | 分析节点不可用，规则裁决可降级。 |

**验收标准**

- 配置错误有明确错误码和 UI 文案。
- `PUT /api/v1/settings` 只能保存通知、页面偏好和普通数据源配置。
- 能力圈配置只能通过 `PUT /api/v1/settings/capability` 保存。
- 规则类设置变更不得直接写入 active 规则，必须生成规则提案。
- 不泄露密钥。

**禁止行为**

- 不把密钥写入审计事件明文。
- 不通过 `PUT /api/v1/settings` 直接修改根本规则、裁决阈值或能力圈配置。

### F10.2 市场数据刷新与快照入库

**用户场景**  
用户手动刷新行情、估值、流动性和情绪指标，供每日纪律和决策咨询使用。

**前置条件**

- 数据源配置可用。
- SQLite 可写。

**用户动作**

- 在设置页或驾驶舱点击刷新市场数据。

**系统行为**

1. 读取已配置的数据源。
2. 拉取行情、估值、成交额、换手率、波动率、融资余额、媒体热度等指标。
3. 计算 `liquidity_state`、`sentiment_state`、PE/PB 分位等字段。
4. 写入 `market_snapshots`。
5. 写入 `audit_events`，`node_action=refresh_market_data`。
6. 刷新失败时，DailyDisciplineGraph 返回 `DATA_STALE` 或沿用已有最新快照并标记数据时间。
7. 单次刷新支持全部成功、部分成功、全部失败和写入失败四类结果。
8. 全部失败中，数据源连接失败、接口不可用或所有标的拉取失败返回 `DATA_SOURCE_UNAVAILABLE`；数据可拉取但交易日期、估值或关键指标过期返回 `DATA_STALE`。
9. 部分成功时接口返回 200，成功标的写入 `market_snapshots`，失败标的进入 `failed_symbols`，审计事件状态为 `degraded`。
10. 快照写入失败时回滚本次市场快照写入，但仍需使用独立事务写入失败审计事件。

**数据写入**

- `market_snapshots`
- `audit_events`

**API / DTO**

- `POST /api/v1/market/refresh`
- `GET /api/v1/market/snapshots/latest?symbol=510300`

**UI 展示**

- 展示行情更新时间、估值分位、流动性状态、情绪状态。
- 数据过期时展示“信息不足”，不生成交易类建议。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| 数据源不可用 | 返回 `DATA_SOURCE_UNAVAILABLE`。 |
| 数据过期 | 返回 `DATA_STALE`，页面展示信息不足。 |
| 指标缺失 | 写入可用字段，缺失字段进入 `market_metrics_json` 的 missing 列表。 |
| 快照写入失败 | 返回 `MARKET_SNAPSHOT_WRITE_FAILED`，不展示正式建议。 |

**验收标准**

- 市场刷新成功后，`market_snapshots` 至少新增 1 条。
- 全部成功时 `audit_events.status=success`，`failed_symbols` 为空。
- 部分成功时返回 200，`audit_events.status=degraded`，`failed_symbols` 写明失败标的与原因。
- 全部失败时返回 `DATA_SOURCE_UNAVAILABLE` 或 `DATA_STALE`：数据源连接失败、接口不可用或所有标的拉取失败为 `DATA_SOURCE_UNAVAILABLE`；数据可拉取但交易日期、估值或关键指标过期为 `DATA_STALE`。
- 快照写入失败时返回 `MARKET_SNAPSHOT_WRITE_FAILED`，本次市场快照事务回滚，不留下部分快照，但仍需使用独立事务写入失败审计事件。
- DailyDisciplineGraph 读取最新 `market_snapshots`。
- 行情过期时，页面状态为 `insufficient_data`。

**禁止行为**

- 不用 LLM 主观判断市场情绪。
- 不在市场数据过期时输出交易类建议。

### F10.3 VecLite 索引维护

**用户场景**  
VecLite 索引损坏、版本不兼容或缺失，需要重建。

**前置条件**

- SQLite 中存在 `rag_chunks` 和 `intelligence_summary`。

**用户动作**

- 在设置页点击重建索引，或系统启动时检测后提示。

**系统行为**

1. 读取 `rag_chunks`。
2. 重建 VecLite 向量索引和 BM25 索引。
3. 更新索引状态。
4. 写入 `audit_events`。

**数据写入**

- VecLite 索引文件。
- `rag_chunks.index_status`。
- `audit_events`。

**API / DTO**

- `POST /api/v1/evidence/rebuild-index` 已在 `docs/api.md` 定义。

**UI 展示**

- 展示索引状态：正常、不可用、重建中、失败。
- 重建期间允许查看 SQLite 摘要，不允许误报证据完整。

**异常状态**

| 情况 | 处理 |
| --- | --- |
| `rag_chunks` 为空 | 返回信息不足。 |
| 重建失败 | 保留原索引文件或标记不可用。 |

**验收标准**

- 删除 VecLite 文件后，可从 SQLite 重建。
- 重建失败不影响 SQLite 事实数据。

**禁止行为**

- 不把 VecLite 当作唯一事实来源。

## 4. API 覆盖清单

以下接口已纳入 `docs/api.md`。开发计划应按这些接口生成任务和契约测试。

| 功能 | 建议接口 | 用途 |
| --- | --- | --- |
| 初始账户录入 | `POST /api/v1/portfolio/init` | 首次录入账户与持仓。 |
| 今日纪律驾驶舱 | `GET /api/v1/dashboard/today` | 查看今日纪律状态。 |
| 当前账户查询 | `GET /api/v1/portfolio/current` | 持仓页读取当前账户状态。 |
| 账户校准 | `POST /api/v1/portfolio/adjustments` | 手动校准本地账户。 |
| 发起决策咨询 | `POST /api/v1/decisions/consult` | 运行 ConsultationGraph 并返回决策详情。 |
| 决策详情 | `GET /api/v1/decisions/{decision_id}` | 重新打开或刷新历史建议。 |
| 历史决策列表 | `GET /api/v1/decisions` | 复盘页查询历史建议。 |
| 用户确认记录 | `POST /api/v1/decisions/{decision_id}/confirmations` | 记录计划、已手动执行、待观察或标记错误。 |
| 情报刷新 | `POST /api/v1/evidence/refresh` | 手动触发证据采集、摘要、验证和索引。 |
| 证据列表 | `GET /api/v1/evidence` | 查询证据卡片、信源等级和证据角色。 |
| 多源验证结果 | `GET /api/v1/evidence/verification` | 查询事件或标的的多源验证状态。 |
| 市场数据刷新 | `POST /api/v1/market/refresh` | 手动刷新行情、估值、流动性和情绪指标。 |
| 最新市场快照 | `GET /api/v1/market/snapshots/latest` | 读取最新市场快照。 |
| 当前规则 | `GET /api/v1/rules/current` | 查看 active 规则版本。 |
| 能力圈读取 | `GET /api/v1/settings/capability` | 查看能力圈配置。 |
| 能力圈更新 | `PUT /api/v1/settings/capability` | 更新能力圈配置。 |
| 审计事件查询 | `GET /api/v1/audit-events` | 查询工作流、用户动作和规则变更审计。 |
| 复盘汇总 | `GET /api/v1/review/summary` | 月度、季度复盘。 |
| 系统设置状态 | `GET /api/v1/settings/system` | 查看本地配置和依赖状态。 |
| 用户设置更新 | `PUT /api/v1/settings` | 保存非规则类设置；规则类变更必须生成提案。 |
| 重建索引 | `POST /api/v1/evidence/rebuild-index` | 从 SQLite 重建 VecLite 索引。 |

## 5. 开发前验收清单

开发计划重写前，需确认以下事项：

- [x] `docs/api.md` 补齐第 4 节列出的接口。
- [x] `docs/data-model.md` 已包含 `capability_configs`、`user_settings` 和 `rag_chunks.index_status` 字段。
- [x] `docs/frontend-contract.md` 已补齐持仓页、设置页、复盘页的 DTO 映射。
- [x] `docs/development-plan.md` 已引用本文档中的 F 编号，并将 P0-P6 拆成细任务。
- [ ] `docs/testing-plan.md` 覆盖 F01-F10 的核心验收场景。

## 6. 可测试验收断言

| 编号 | 场景 | 断言 |
| --- | --- | --- |
| A01 | 首次使用 | 无账户数据时，`GET /api/v1/dashboard/today` 返回 `DATA_REQUIRED` 或 `dashboard_state=first_use`，且不创建 `decision_records`。 |
| A02 | 正常每日纪律 | 生成正式建议后，`decision_records` 增加 1 条，`evidence_refs` 至少 1 条，`audit_events` 至少包含 StateSnapshotNode、EvidenceRetrievalNode、RuleArbitrationNode、DecisionRecordNode。 |
| A03 | 证据不足 | `EVIDENCE_NOT_FOUND` 时，`final_verdict.status=insufficient_data`，前端不展示交易类建议。 |
| A04 | VecLite 不可用 | SQLite 摘要充足时 `workflow_status=degraded`；SQLite 摘要不足时页面状态为 `insufficient_data`。 |
| A05 | 能力圈外 | 咨询能力圈外标的时，`final_verdict.status=rejected`，且不调用 ValueAnalystNode 与 TrendRiskOfficerNode。 |
| A06 | 记录计划 | `planned` 只写 `operation_confirmations` 和 `audit_events`，不写 `position_transactions`，不新增账户快照。 |
| A07 | 已手动执行 | `executed_manually` 成功后，同时写入 `operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events`。 |
| A08 | 已手动执行失败 | 交易流水或快照任一写入失败时，事务回滚，不留下部分确认记录。 |
| A09 | 标记错误 | `marked_error` 成功后，同一事务写入 `operation_confirmations`、`error_cases`、`audit_events`，响应返回 `error_case_id`。 |
| A10 | C 级信源 | C 级证据只能以 `evidence_role=background` 返回，不得出现在正式裁决的 `formal` 证据引用中。 |
| A11 | LLM 不可用 | DeepSeek 调用失败时，返回 `ANALYST_UNAVAILABLE`，最终裁决来自规则引擎，`workflow_status=degraded`。 |
| A12 | 守门人审计通过 | 提案状态从 `under_gatekeeper_audit` 变为 `pending_final_confirm`，不写 `rule_versions`；`sample_count<3` 的提案不得进入守门人审计，返回 `BAD_REQUEST`。 |
| A13 | 规则最终确认 | `confirm=true` 且 `sample_count>=3` 后创建新 active `rule_versions`，旧 active 归档，提案状态为 `applied`；`sample_count<3` 的提案最终确认返回 `BAD_REQUEST`，不得写 `rule_versions`。 |
| A14 | 审计事件 | 审计页返回字段同时包含 `action`、`node_name` 与 `node_action`，前端分别展示业务动作、节点名称和节点动作。 |
| A15 | 禁止自动交易 | API 列表中不存在买入、卖出、撤单、改单接口；前端确认区不出现一键交易文案。 |
| A16 | 市场数据刷新 | `POST /api/v1/market/refresh` 全部成功时新增 `market_snapshots` 且 `audit_events.status=success`；部分成功时返回 200、写入成功标的、返回 `failed_symbols` 且 `audit_events.status=degraded`；全部失败时返回 `DATA_SOURCE_UNAVAILABLE` 或 `DATA_STALE`；快照写入失败时返回 `MARKET_SNAPSHOT_WRITE_FAILED` 且不留下部分写入。 |
| A17 | 预期收益评估 | 输出只能作为情景概率展示，不得覆盖最终规则裁决，不得承诺收益；`available` 必须包含 upside/base/downside 且可返回概率，`insufficient` 不得返回精确概率且必须写样本不足说明，`unavailable` 必须返回空 `scenarios` 并写定性原因。 |

## 7. F 编号到开发计划映射

| 功能编号 | 功能 | 开发计划阶段 |
| --- | --- | --- |
| F01 | 首次使用与基础设置 | P1 数据底座、P4 HTTP API、P5 前端驾驶舱 |
| F02 | 今日纪律驾驶舱 | P2 领域规则、P3 工作流、P4 HTTP API、P5 前端驾驶舱 |
| F03 | 持仓与账户管理 | P1 数据底座、P4 HTTP API、P5 前端持仓页 |
| F04 | 决策咨询 | P2 领域规则、P3 工作流、P4 HTTP API、P5 决策详情页 |
| F05 | 证据与情报管理 | P1 数据底座、P3 证据工作流、P4 证据 API、P5 证据页 |
| F06 | 用户确认与线下动作记录 | P1 事务写入、P4 确认 API、P5 确认区 |
| F07 | 错误案例与规则进化 | P1 数据底座、P3 Evolution / Gatekeeper、P4 规则提案 API、P5 规则提案页 |
| F08 | 规则与纪律管理 | P1 规则版本、P4 规则与能力圈 API、P5 规则页 |
| F09 | 审计事件与复盘 | P1 审计表、P4 审计与复盘 API、P5 审计页 |
| F10 | 系统配置、市场数据与数据维护 | P0 配置、P3 市场刷新工作流、P4 设置/市场/索引 API、P5 设置页 |

## 8. 优先级建议

| 优先级 | 功能 | 理由 |
| --- | --- | --- |
| 必须先完成 | F01、F02、F03、F06、F09 | 没有账户、驾驶舱、确认和审计，系统无法形成可追溯闭环。 |
| 第二批 | F04、F05、F08 | 决策咨询、证据和规则是 Agent 能力核心。 |
| 第三批 | F07、F10 | 规则进化和系统维护依赖前面数据与审计能力。 |

说明：虽然功能有批次顺序，但数据模型、API 契约和状态机应一次性按完整系统设计，避免后续迁移成本。
