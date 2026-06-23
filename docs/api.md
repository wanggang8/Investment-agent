# Investment Agent HTTP API 契约

> 文档版本：v1.0
> 最后更新：2026-06-17
> 适用范围：`cmd/server` HTTP 服务、前端 Web 控制台、应用层 handler DTO。

## 1. API 定位

HTTP API 只服务本地 Web 控制台和本地自动任务。系统默认不暴露公网服务，不接入券商交易 API，不提供自动下单能力。本地自动任务只用于刷新数据、生成报告和写入审计日志，不执行任何交易动作。

边界规则：

- API 统一使用 JSON。
- API 路径统一以 `/api/v1` 开头。
- 前端不得直接访问 SQLite、VecLite 或本地文件。
- 所有会改变账户状态的接口必须来自用户确认动作。
- 系统建议不会直接改变账户状态。
- 所有响应必须包含 `request_id`。

## 2. 通用响应结构

成功响应：

```json
{
  "request_id": "req_20260522_000001",
  "data": {},
  "meta": {
    "generated_at": "2026-05-22T09:30:00+08:00",
    "rule_version": "v3.0"
  }
}
```

失败响应：

```json
{
  "request_id": "req_20260522_000001",
  "error": {
    "code": "DATA_STALE",
    "message": "行情数据已过期，暂停生成交易类建议。",
    "detail": "latest_market_date=2026-05-20"
  }
}
```

## 3. 通用错误码

HTTP API 层必须从统一应用错误生成错误响应。未知错误返回 `INTERNAL_ERROR` 与 HTTP 500，不得暴露 SQL、文件路径或外部服务原始错误。应用错误响应必须包含 `request_id`、`error.code` 和 `error.message`。

| 错误码 | HTTP 状态 | 含义 | 前端处理 |
| --- | --- | --- | --- |
| BAD_REQUEST | 400 | 请求参数错误 | 展示字段级错误 |
| NOT_FOUND | 404 | 资源不存在 | 展示空状态 |
| DATA_REQUIRED | 409 | 缺少账户或持仓数据 | 引导录入 |
| DATA_STALE | 409 | 行情或估值过期 | 展示信息不足状态 |
| DATA_SOURCE_UNAVAILABLE | 503 | 市场或行情数据源不可用 | 展示数据源异常状态，不生成交易类建议 |
| MARKET_SNAPSHOT_WRITE_FAILED | 500 | 市场快照写入失败 | 展示刷新失败，不生成交易类建议 |
| EVIDENCE_NOT_FOUND | 409 | 未找到有效证据 | 暂停交易类建议 |
| SOURCE_VERIFICATION_FAILED | 409 | 多源验证不满足 | 进入冻结观察 |
| VECTOR_INDEX_UNAVAILABLE | 409 | VecLite 索引不可用 | SQLite 摘要充足时降级展示；摘要不足时展示信息不足状态 |
| ANALYST_UNAVAILABLE | 503 | LLM 分析节点不可用 | 使用规则裁决降级输出 |
| RULE_VERSION_MISSING | 409 | 规则版本缺失 | 暂停裁决 |
| DECISION_RECORD_FAILED | 409 | 决策记录保存失败 | 不展示正式建议 |
| CONFLICT | 409 | 数据冲突或约束冲突 | 展示冲突状态 |
| INVALID_STATE | 409 | 非法状态流转 | 阻止本次操作 |
| INTERNAL_ERROR | 500 | 内部错误 | 展示通用失败状态 |

## 3.1 P4 HTTP API 响应约束

P4 业务 API handler 必须使用第 2 节通用响应信封。

约束：

- 成功响应必须包含 `request_id` 和 `data`。
- 成功响应可在适用场景包含 `meta.generated_at` 与 `meta.rule_version`。
- 失败响应必须包含 `request_id`、`error.code` 和 `error.message`。
- 未知内部错误必须返回 `INTERNAL_ERROR` 与 HTTP 500。
- 失败响应不得暴露 SQL、文件路径或外部服务原始错误文本。

## 3.2 P4 核心 API 表面

P4 只实现本文档和 `docs/development-plan.md` P4.2 列出的 HTTP API。完成后必须存在 dashboard、decision、portfolio、evidence、market、rule、settings、audit、review、notification 十组 handler，且各 handler 返回的 DTO 字段必须对齐 `docs/frontend-contract.md`。

## 3.3 P4 确认状态边界

`POST /api/v1/decisions/{decision_id}/confirmations` 必须执行本文档约束的状态流转。若决策已处于 `executed_manually` 或 `marked_error`，再次确认必须返回 `BAD_REQUEST`，且不得重复写入账户快照、交易流水或错误案例。

## 4. 核心枚举

| 枚举 | 可选值 | 说明 |
| --- | --- | --- |
| dashboard_state | first_use / normal / insufficient_data / frozen_watch / high_risk | 驾驶舱页面状态 |
| workflow_status | completed / degraded / failed | 工作流状态 |
| position_state | normal / sell_only / frozen_watch | 持仓状态 |
| verification_status | satisfied / failed / background_only | 多源验证状态 |
| confirmation_status | not_required / pending / planned / executed_manually / watch / marked_error | 用户确认状态 |
| confirmation_type | planned / executed_manually / watch / marked_error | 用户确认动作类型 |
| final_verdict.status | buy_allowed / hold / reduce / sell_only / frozen_watch / rejected / insufficient_data | 最终裁决状态 |
| operation_type | buy / sell / reduce | 线下成交操作类型，仅用于 `confirmation_type=executed_manually` |
| audit_result | approved / rejected / needs_user_review | 守门人审计结果 |
| rule_proposal.status | draft / pending_user_confirm / under_gatekeeper_audit / pending_final_confirm / rejected / applied | 规则提案状态 |
| audit.actor | system / user / gatekeeper | 审计操作者 |
| audit.action | generate_decision / confirm_operation / mark_error / create_proposal / audit_rule_change / update_rule / refresh_market_data / update_settings / update_capability / rebuild_index / run_local_task / risk_alert | 审计动作 |
| risk_type | valuation_high / buy_thesis_broken / liquidity_danger / sentiment_extreme / position_limit_breach / insufficient_evidence / data_degraded | 风险预警类型 |
| risk_severity | info / warning / critical | 风险预警严重程度 |
| risk_sop_status | triggered / active / observing / escalated / resolved / archived | 风险预警 SOP 状态 |
| rule_effect_validation.status | not_evaluated / insufficient / passed / failed / needs_more_samples / needs_user_review | 规则效果验证状态 |
| rule_effect_overfit_risk | low / medium / high | 规则效果验证过拟合风险 |
| rule_effect_replay_result | passed / failed / mixed / unknown | 历史回放验证结果 |
| rule_effect_guardrail_decision | passed / rejected / needs_user_review | 效果验证门禁结论 |
| rule_effect_trend_direction | improved / flat / worsened / unknown | 应用后规则效果趋势 |
| audit.status | success / degraded / failed | 审计节点状态 |
| liquidity_state | normal / warning / danger | 市场流动性状态 |
| sentiment_state | cold / neutral / hot / extreme | 市场情绪状态 |
| precision_status | available / insufficient / unavailable | 预期收益精度状态 |
| ops_status.data_source_status | success / degraded / failed / empty / unknown | 复盘运维状态中的数据源状态 |
| ops_status.index_status | success / degraded / failed / missing / unknown | 复盘运维状态中的索引状态 |
| ops_status.review_status | success / degraded / failed / empty / unknown | 复盘运维状态中的复盘任务状态 |

## 4.1 健康检查

`GET /api/v1/health`

用途：探测 HTTP 服务是否存活，供本地开发与部署健康检查使用。无需认证。

响应 HTTP 200，JSON 体：

```json
{"status":"ok"}
```

说明：P0 骨架阶段返回简单 JSON；后续业务 API 仍使用第 2 节通用信封（含 `request_id`）。

## 5. 驾驶舱 API

### 5.1 获取今日纪律驾驶舱

`GET /api/v1/dashboard/today`

用途：获取今日纪律状态、账户摘要、触发规则、证据摘要和裁决结果。

字段约束：`evidence_summary` 在证据数据完整时返回；当行情过期、证据不足或 VecLite 索引不可用且 SQLite 摘要不完整时可缺失。前端收到缺失字段时展示“信息不足”或隐藏证据摘要模块。

响应：

```json
{
  "request_id": "req_20260522_000001",
  "data": {
    "dashboard_state": "normal",
    "discipline_status": "观察",
    "data_updated_at": "2026-05-22T09:00:00+08:00",
    "portfolio_summary": {
      "total_assets": 120000.00,
      "cash_ratio": 0.08,
      "high_risk_ratio": 0.18,
      "position_count": 5
    },
    "market_summary": {
      "pe_percentile": 63.0,
      "pb_percentile": 58.0,
      "sentiment_state": "neutral",
      "liquidity_state": "normal"
    },
    "triggered_rules": [
      {
        "rule_id": "R-3",
        "rule_name": "不超过仓位上限",
        "severity": "warning",
        "description": "高风险资产仓位接近上限。"
      }
    ],
    "decision_summary": {
      "decision_id": "dec_20260522_0001",
      "verdict": "暂停新增买入，继续观察。",
      "final_verdict_status": "hold",
      "prohibited_actions": ["新增买入"],
      "optional_actions": ["继续观察", "查看证据"],
      "action_required": false,
      "confirmation_status": "not_required"
    },
    "evidence_summary": {
      "source_count": 3,
      "highest_source_level": "A",
      "verification_status": "satisfied"
    }
  },
  "meta": {
    "generated_at": "2026-05-22T09:30:00+08:00",
    "rule_version": "v3.0"
  }
}
```

## 6. 每日纪律报告 API

### 6.1 获取今日每日纪律报告

`GET /api/v1/daily-discipline/reports/today`

用途：获取当前本地日期的每日纪律报告。报告来源于 DailyDisciplineGraph 手动任务或每日自动运行结果，只作为本地阅读、复核和追踪入口，不产生交易执行副作用。

响应：

```json
{
  "request_id": "req_20260608_000001",
  "data": {
    "report_id": "daily_report:2026-06-08:holdings:hash:v1",
    "local_date": "2026-06-08",
    "scope": "holdings",
    "status": "success",
    "summary": "今日纪律报告已生成",
    "source_type": "auto_run",
    "source_id": "daily_auto_run:2026-06-08:holdings:hash:v1",
    "decision_id": "dec_20260608_0001",
    "decision_link": "/decisions/dec_20260608_0001",
    "auto_run_link": "/daily-auto-run",
    "audit_link": "/audit?source_id=daily_auto_run:2026-06-08:holdings:hash:v1",
    "notification_link": "/notifications?source_id=daily_auto_run:2026-06-08:holdings:hash:v1",
    "failure_code": "",
    "failure_reason": "",
    "missing_action": "",
    "missing_categories": [],
    "final_verdict": "暂停新增买入，继续观察。",
    "verdict_status": "hold",
    "evidence": {
      "evidence_count": 3,
      "independent_source_count": 2,
      "high_grade_independent_source_count": 2
    },
    "trend": {
      "success_count": 4,
      "degraded_count": 1,
      "failed_count": 0,
      "insufficient_data_count": 0
    },
    "safety_note": "每日纪律报告只用于人工复核，不会自动执行交易。",
    "updated_at": "2026-06-08T09:30:00+08:00"
  }
}
```

缺少前提时，`status` 返回 `insufficient_data` 或 `not_started`，并通过 `missing_categories`、`failure_code`、`failure_reason` 说明缺少账户、持仓、行情、证据、规则或配置等前提。此时不得伪造报告摘要、证据、预期收益或交易指令。

### 6.2 获取每日纪律报告列表

`GET /api/v1/daily-discipline/reports?status=success&limit=30`

用途：按更新时间倒序获取本地历史每日纪律报告列表。`status` 可选；`limit` 有服务端上限。

响应：

```json
{
  "request_id": "req_20260608_000002",
  "data": {
    "reports": [
      {
        "report_id": "daily_report:2026-06-08:holdings:hash:v1",
        "local_date": "2026-06-08",
        "scope": "holdings",
        "status": "success",
        "summary": "今日纪律报告已生成",
        "source_type": "auto_run",
        "source_id": "daily_auto_run:2026-06-08:holdings:hash:v1",
        "decision_id": "dec_20260608_0001",
        "decision_link": "/decisions/dec_20260608_0001",
        "auto_run_link": "/daily-auto-run",
        "audit_link": "/audit?source_id=daily_auto_run:2026-06-08:holdings:hash:v1",
        "notification_link": "/notifications?source_id=daily_auto_run:2026-06-08:holdings:hash:v1",
        "failure_code": "",
        "failure_reason": "",
        "missing_action": "",
        "missing_categories": [],
        "final_verdict": "暂停新增买入，继续观察。",
        "verdict_status": "hold",
        "evidence": {
          "evidence_count": 3,
          "independent_source_count": 2,
          "high_grade_independent_source_count": 2
        },
        "trend": {
          "success_count": 4,
          "degraded_count": 1,
          "failed_count": 0,
          "insufficient_data_count": 0
        },
        "safety_note": "每日纪律报告只用于人工复核，不会自动执行交易。",
        "updated_at": "2026-06-08T09:30:00+08:00"
      }
    ]
  }
}
```

### 6.3 获取每日纪律报告详情

`GET /api/v1/daily-discipline/reports/{report_id}`

用途：获取单份每日纪律报告详情。未知 `report_id` 返回 404 与统一错误信封。

详情响应字段与 6.1 一致。前端可通过 `decision_link`、`audit_link`、`notification_link`、`auto_run_link` 跳转到关联材料；所有链接保持本地只读语义。

报告详情可包含 `risk_alerts` 数组。每个元素包含 `alert_id`、`risk_type`、`severity`、`sop_status`、`symbol`、`trigger_summary`、`prohibited_actions`、`suggested_actions`、`link` 和 `safety_note`。该摘要只用于本地追踪和人工复核；仅展示 triggered / active / observing / escalated 风险，不展示 resolved / archived 历史预警。

## 6.4 风险预警 API

### 6.4.1 获取风险预警列表

`GET /api/v1/risk-alerts?status=active,escalated&symbol=510300`

用途：读取本地风险预警事实，供风险预警中心展示。响应使用统一信封，`data` 为 `PageResult<RiskAlertDTO>`。

### 6.4.2 获取风险预警详情

`GET /api/v1/risk-alerts/{alert_id}`

用途：查看单条本地风险预警的触发依据、SOP 状态、禁止动作、建议人工动作和关联材料。未知 `alert_id` 返回 `NOT_FOUND`。

### 6.4.3 更新风险预警 SOP 生命周期

`POST /api/v1/risk-alerts/{alert_id}/lifecycle`

请求：

```json
{
  "status": "observing",
  "reason": "前端人工 SOP 操作：继续观察"
}
```

用途：执行本地 SOP 状态动作。允许的目标状态包括 observing、escalated、resolved、archived；终态 `resolved` / `archived` 不得再转回非终态，非法流转返回 `INVALID_STATE`。

`RiskAlertDTO` 字段：

| 字段 | 说明 |
| --- | --- |
| alert_id | 风险预警 ID |
| risk_type | 风险类型，见第 4 节枚举 |
| severity | info / warning / critical |
| sop_status | triggered / active / observing / escalated / resolved / archived |
| symbol | 影响标的 |
| trigger_summary | 触发摘要 |
| trigger_context | 触发上下文，可包含估值、流动性、情绪、证据或 source health |
| prohibited_actions | 禁止动作列表 |
| suggested_actions | 建议人工动作列表 |
| related_decision_id / decision_link | 关联决策 |
| related_report_id / report_link | 关联每日纪律报告 |
| related_notification_id / notification_link | 关联通知 |
| related_audit_event_id / audit_link | 关联审计 |
| last_triggered_at / resolved_at / resolution_reason | 生命周期字段 |
| safety_note | 固定安全提示，声明只用于本地人工复核 |
| created_at / updated_at | 创建与更新时间 |

事务边界：触发或更新风险预警时，系统可在同一本地事务中写入 `risk_alerts`、`notifications` 和 `audit_events`。不得写入或修改 positions、portfolio_snapshots、operation_confirmations、position_transactions、rule_versions、broker state、orders 或 external notifications。

错误分类：请求格式或非法枚举返回 `BAD_REQUEST`；未知预警返回 `NOT_FOUND`；非法状态流转返回 `INVALID_STATE`；唯一约束冲突返回 `CONFLICT`；仓储异常返回 `INTERNAL_ERROR` 且不得暴露底层错误文本。

## 7. 决策 API

### 7.1 发起决策咨询

`POST /api/v1/decisions/consult`

请求：

```json
{
  "question": "这个标的还能不能持有？",
  "symbol": "510300",
  "scenario": "hold_review"
}
```

响应策略：接口同步执行 ConsultationGraph，并直接返回可渲染的完整决策详情。前端可立即展示裁决、证据链、Agent 观点、预期收益情景和裁决链；`GET /api/v1/decisions/{decision_id}` 用于刷新或重新打开详情页。

响应：

```json
{
  "request_id": "req_20260522_000002",
  "data": {
    "decision_id": "dec_20260522_0002",
    "question": "这个标的还能不能持有？",
    "symbol": "510300",
    "capability_check": {
      "status": "in_scope",
      "reason": "该标的属于用户配置的指数 ETF 能力圈。"
    },
    "workflow_status": "completed",
    "account_snapshot": {
      "snapshot_id": "snap_20260522_0900",
      "cash_ratio": 0.08,
      "high_risk_ratio": 0.18
    },
    "triggered_rules": [
      {
        "rule_id": "R-3",
        "rule_name": "不超过仓位上限",
        "severity": "warning",
        "description": "高风险资产仓位接近上限。"
      }
    ],
    "evidence_chain": [
      {
        "evidence_id": "ev_001",
        "source_name": "巨潮资讯",
        "source_level": "A",
        "published_at": "2026-05-22T08:00:00+08:00",
        "summary": "相关公告未显示买入逻辑破坏。"
      }
    ],
    "analyst_reports": [
      {
        "agent_name": "价值分析师",
        "conclusion": "估值处于观察区，暂不新增买入。",
        "key_reasons": ["PE 分位处于 50%-80%", "买入逻辑未破坏"],
        "risk_warnings": ["继续关注估值上行风险"],
        "confidence": "medium",
        "evidence_ids": ["ev_001"],
        "prompt_version": "p37-analyst-v1",
        "model": "gpt-5.4-mini",
        "input_summary": "value 510300 ...",
        "output_summary": "估值处于观察区...",
        "parse_status": "parsed",
        "quality_status": "passed"
      }
    ],
    "expected_return_scenarios": {
      "sample_count": 36,
      "sample_window": "2014-01-01~2026-05-22",
      "screening_condition": "PE 分位 50%-80%，流动性 normal",
      "precision_status": "available",
      "scenarios": [
        {"scenario": "upside", "return_range": "8%~15%", "probability": 0.28, "trigger": "当前价格进入上行情景下沿时复核移动止盈"},
        {"scenario": "base", "return_range": "0%~8%", "probability": 0.52, "trigger": "当前价格突破基准情景上沿时复核分批止盈"},
        {"scenario": "downside", "return_range": "-12%~0%", "probability": 0.20, "trigger": "当前价格跌破下行情景下沿时复核买入逻辑"}
      ],
      "sell_evaluation": {
        "status": "review_required",
        "triggers": ["当前价格突破基准情景上沿"],
        "prompts": ["人工复核是否需要分批止盈"],
        "actions": ["记录人工计划或继续观察"],
        "non_trading_disclaimer": "卖出评估仅提示人工复核，不会自动交易。"
      },
      "reassessment_trigger": {
        "reason": "基准情景中枢下移超过 15%",
        "boundary": "base_midpoint_downshift_gt_15pct",
        "current_value": 0.16
      },
      "disclaimer": "仅为历史样本情景分析，不承诺收益；最终裁决以规则裁决为准。"
    },
    "arbitration_chain": [
      {
        "priority": 3,
        "rule_id": "R-3",
        "result": "限制新增买入"
      }
    ],
    "final_verdict": {
      "status": "hold",
      "display_text": "买入逻辑未破坏，继续持有并观察估值变化。",
      "prohibited_actions": ["新增买入"],
      "optional_actions": ["继续持有", "查看证据链"]
    },
    "user_confirmation": {
      "confirmation_status": "pending",
      "available_actions": ["planned", "executed_manually", "watch", "marked_error"]
    }
  },
  "meta": {
    "generated_at": "2026-05-22T09:35:00+08:00",
    "rule_version": "v3.0"
  }
}
```

`expected_return_scenarios` 字段约束：

| precision_status | scenarios | probability | reason |
| --- | --- | --- | --- |
| available | 必须包含 upside / base / downside，可包含 `trigger` | 可返回精确概率 | 可为空 |
| insufficient | 可包含收益区间和 `trigger`，但不得返回精确概率，`probability` 必须为空 | 空 | 必须说明样本不足 |
| unavailable | 必须为空数组 | 空 | 必须说明无法生成收益区间的原因 |

新增 P28 字段约束：

- `sample_count`、`sample_window`、`screening_condition` 用于解释样本来源；consultation response 当前从本地持仓、最新市场快照和已保存公开市场元数据派生可解释样本数，缺少历史样本时不得伪造成完整历史回测；样本不足或不可用时也必须返回可解释原因。
- `scenarios[].trigger` 是情景边界说明，只用于人工复核提示；后端动态卖出评估按情景区间边界计算，例如进入 upside 下沿、突破 base 上沿、跌破 downside 下沿，而不是按单点收益率直接触发。
- `sell_evaluation.status/triggers/prompts/actions/non_trading_disclaimer` 只表达动态卖出评估材料；`not_applicable` 表示缺当前价格、缺持仓成本/基准价或缺情景边界，不创建交易、不更新账户、不触发确认或通知。
- `reassessment_trigger.reason/boundary/current_value` 用于说明何时需要重新评估买入逻辑或情景基准；旧决策 JSON 缺失这些字段时前端按空/不适用展示。

### 7.2 获取决策详情

`GET /api/v1/decisions/{decision_id}`

响应示例：

```json
{
  "request_id": "req_20260522_000004",
  "data": {
    "decision_id": "dec_20260522_0002",
    "generated_at": "2026-05-22T09:35:00+08:00",
    "account_snapshot": {
      "snapshot_id": "snap_20260522_0900",
      "cash": 9600.00,
      "total_assets": 120000.00,
      "cash_ratio": 0.08,
      "high_risk_ratio": 0.18
    },
    "triggered_rules": [
      {
        "rule_id": "R-3",
        "rule_name": "不超过仓位上限",
        "severity": "warning",
        "description": "高风险资产仓位接近上限。"
      }
    ],
    "evidence_chain": [
      {
        "evidence_id": "ev_001",
        "source_name": "巨潮资讯",
        "source_level": "A",
        "published_at": "2026-05-22T08:00:00+08:00",
        "summary": "相关公告未显示买入逻辑破坏。"
      }
    ],
    "analyst_reports": [
      {
        "agent_name": "价值分析师",
        "conclusion": "估值处于观察区，暂不新增买入。",
        "key_reasons": ["PE 分位处于 50%-80%", "买入逻辑未破坏"],
        "risk_warnings": ["继续关注估值上行风险"],
        "confidence": "medium",
        "evidence_ids": ["ev_001"]
      }
    ],
    "expected_return_scenarios": {
      "sample_count": 36,
      "sample_window": "2014-01-01~2026-05-22",
      "screening_condition": "PE 分位 50%-80%，流动性 normal",
      "precision_status": "available",
      "scenarios": [
        {"scenario": "upside", "return_range": "8%~15%", "probability": 0.28, "trigger": "当前价格进入上行情景下沿时复核移动止盈"},
        {"scenario": "base", "return_range": "0%~8%", "probability": 0.52, "trigger": "当前价格突破基准情景上沿时复核分批止盈"},
        {"scenario": "downside", "return_range": "-12%~0%", "probability": 0.20, "trigger": "当前价格跌破下行情景下沿时复核买入逻辑"}
      ],
      "sell_evaluation": {
        "status": "review_required",
        "triggers": ["当前价格突破基准情景上沿"],
        "prompts": ["人工复核是否需要分批止盈"],
        "actions": ["记录人工计划或继续观察"],
        "non_trading_disclaimer": "卖出评估仅提示人工复核，不会自动交易。"
      },
      "reassessment_trigger": {
        "reason": "基准情景中枢下移超过 15%",
        "boundary": "base_midpoint_downshift_gt_15pct",
        "current_value": 0.16
      },
      "disclaimer": "仅为历史样本情景分析，不承诺收益；最终裁决以规则裁决为准。"
    },
    "arbitration_chain": [
      {
        "priority": 3,
        "rule_id": "R-3",
        "result": "限制新增买入"
      }
    ],
    "final_verdict": {
      "status": "hold",
      "display_text": "继续持有，暂停新增买入。",
      "prohibited_actions": ["新增买入"],
      "optional_actions": ["继续观察", "查看证据"]
    },
    "user_confirmation": {
      "confirmation_status": "pending",
      "available_actions": ["planned", "executed_manually", "watch", "marked_error"]
    }
  },
  "meta": {
    "generated_at": "2026-05-22T09:36:00+08:00",
    "rule_version": "v3.0"
  }
}
```

### 7.3 查询历史决策

`GET /api/v1/decisions?from=2026-05-01&to=2026-05-22&status=pending`

查询参数：`status` 按 `confirmation_status` 过滤；`from` 与 `to` 使用 `YYYY-MM-DD`，非法日期或 `from > to` 返回 `BAD_REQUEST`。

响应：

```json
{
  "request_id": "req_20260522_000003",
  "data": {
    "items": [
      {
        "decision_id": "dec_20260522_0002",
        "display_title": "2026-05-22 第 2 条建议",
        "symbol": "510300",
        "final_verdict": "继续持有",
        "triggered_rule_ids": ["R-3"],
        "confirmation_status": "pending",
        "generated_at": "2026-05-22T09:35:00+08:00"
      }
    ],
    "total": 1
  },
  "meta": {
    "generated_at": "2026-05-22T09:36:00+08:00",
    "rule_version": "v3.0"
  }
}
```

## 8. 持仓与账户 API

### 8.1 初始账户录入

`POST /api/v1/portfolio/init`

用途：首次录入本地账户、现金、总资产和完整持仓集合。该接口只初始化本地事实库，不连接券商账户。

请求：

```json
{
  "cash": 9600.00,
  "total_assets": 120000.00,
  "positions": [
    {
      "symbol": "510300",
      "name": "沪深300ETF",
      "quantity": 1000,
      "cost_price": 3.850,
      "current_price": 4.120,
      "buy_date": "2026-05-01",
      "position_state": "normal",
      "buy_reason": "低估区分批配置核心资产",
      "asset_tag": "core"
    }
  ]
}
```

响应示例：

```json
{
  "request_id": "req_20260522_000005",
  "data": {
    "snapshot_id": "snap_20260522_0900",
    "position_count": 1,
    "position_snapshot_count": 1,
    "audit_event_ids": ["audit_005"]
  },
  "meta": {
    "generated_at": "2026-05-22T09:30:00+08:00",
    "rule_version": "v3.0"
  }
}
```

事务要求：必须在同一事务写入 `portfolio_snapshots + positions + position_snapshots + audit_events`。任一写入失败时，全部回滚。

### 8.2 获取当前账户与持仓

`GET /api/v1/portfolio/current`

用途：持仓页读取当前本地账户状态。响应聚合最新账户快照和当前持仓列表，前端不再分别请求账户快照与持仓列表。

响应示例：

```json
{
  "request_id": "req_20260522_000006",
  "data": {
    "snapshot": {
      "snapshot_id": "snap_20260522_0900",
      "snapshot_time": "2026-05-22T09:00:00+08:00",
      "cash": 9600.00,
      "total_assets": 120000.00,
      "cash_ratio": 0.08,
      "high_risk_ratio": 0.18,
      "position_count": 1
    },
    "positions": [
      {
        "position_id": "pos_510300",
        "symbol": "510300",
        "name": "沪深300ETF",
        "quantity": 1000,
        "cost_price": 3.850,
        "current_price": 4.120,
        "market_value": 4120.00,
        "unrealized_profit_ratio": 0.0701,
        "position_state": "normal",
        "buy_date": "2026-05-01",
        "buy_reason": "低估区分批配置核心资产",
        "asset_tag": "core"
      }
    ]
  },
  "meta": {
    "generated_at": "2026-05-22T09:30:00+08:00",
    "rule_version": "v3.0"
  }
}
```

### 8.3 手动校准账户

`POST /api/v1/portfolio/adjustments`

用途：用户校准本地账户状态。该接口不代表实际交易，不写 `position_transactions`。

请求：

```json
{
  "cash": 9800.00,
  "total_assets": 120500.00,
  "adjust_reason": "与券商账户手动核对后校准",
  "positions": [
    {
      "symbol": "510300",
      "name": "沪深300ETF",
      "quantity": 1000,
      "cost_price": 3.850,
      "current_price": 4.150,
      "buy_date": "2026-05-01",
      "position_state": "normal",
      "buy_reason": "低估区分批配置核心资产",
      "asset_tag": "core"
    }
  ]
}
```

响应示例：

```json
{
  "request_id": "req_20260522_000012",
  "data": {
    "snapshot_id": "snap_20260522_0930",
    "position_count": 1,
    "audit_event_ids": ["audit_012"]
  },
  "meta": {
    "generated_at": "2026-05-22T09:35:00+08:00",
    "rule_version": "v3.0"
  }
}
```

事务要求：必须在同一事务写入 `positions + portfolio_snapshots + position_snapshots + audit_events`。不得写入 `position_transactions`。

### 8.4 本地账户 onboarding 与维护 API（P33）

以下接口只记录用户本地输入或线下已完成动作，不连接券商账户，不创建订单，不代表系统执行交易。

#### 8.4.1 新增或编辑当前持仓

`POST /api/v1/portfolio/holdings`

请求：

```json
{
  "position_id": "pos_510300",
  "reason": "用户本地持仓编辑",
  "confirmation": "我确认这是本地账户事实记录，不代表系统交易。",
  "position": {
    "symbol": "510300",
    "name": "沪深300ETF",
    "quantity": 1000,
    "cost_price": 3.850,
    "current_price": 4.120,
    "buy_date": "2026-05-01",
    "position_state": "normal",
    "buy_reason": "低估区分批配置核心资产",
    "asset_tag": "core"
  }
}
```

响应 `data` 至少包含 `snapshot_id`、`position_id`、`audit_event_ids` 和 `safety_statement`。后端必须在同一事务写入当前持仓、账户快照、持仓快照和审计事件；历史快照不得被物理删除。

#### 8.4.2 移除当前持仓

`POST /api/v1/portfolio/holdings/remove`

请求：

```json
{
  "position_id": "pos_510300",
  "reason": "用户确认不再作为当前持仓展示",
  "confirmation": "我确认只改变当前本地视图，不删除历史事实。"
}
```

移除只改变当前聚合态并生成新的本地事实；历史 `position_snapshots`、交易流水和审计仍必须可查。

#### 8.4.3 记录线下交易

`POST /api/v1/portfolio/offline-transactions`

请求：

```json
{
  "operation_type": "sell",
  "symbol": "510300",
  "name": "沪深300ETF",
  "quantity": 100,
  "price": 4.250,
  "fees": 1.20,
  "executed_at": "2026-05-22T10:00:00+08:00",
  "note": "用户补记线下交易",
  "buy_reason": "低估区分批配置核心资产",
  "asset_tag": "core"
}
```

`operation_type` 只能是 `buy / sell / reduce`。卖出或减仓不得超过当前本地数量。成功时必须写入 `position_transactions`，并在同一事务更新当前持仓、账户快照、持仓快照和审计。

#### 8.4.4 批量导入校验

`POST /api/v1/portfolio/imports/validate`

请求：

```json
{
  "rows": [
    {
      "row_number": 1,
      "row_type": "holding",
      "symbol": "510300",
      "name": "沪深300ETF",
      "quantity": 1000,
      "cost_price": 3.850,
      "current_price": 4.120,
      "buy_reason": "低估区分批配置核心资产"
    }
  ]
}
```

响应 `data` 必须包含 `import_batch_id`、`summary.row_count`、`summary.valid_count`、`summary.invalid_count` 和逐行校验结果。校验阶段不得写账户、持仓、交易或快照事实；但必须保存导入批次 metadata、校验摘要和已校验 rows hash。

#### 8.4.5 批量导入确认

`POST /api/v1/portfolio/imports/confirm`

请求：

```json
{
  "import_batch_id": "import_20260612_0001",
  "confirm_reason": "确认导入本地账户事实",
  "rows": []
}
```

确认阶段必须验证 `import_batch_id` 已存在、状态为 `validated`、无无效行、rows hash 非空且与当前提交 rows 匹配。成功后在同一事务写入导入事实、账户快照、持仓快照、交易流水和审计事件；失败不得留下部分成功状态。

#### 8.4.6 记录错误修正审计

`POST /api/v1/portfolio/corrections`

请求：

```json
{
  "target_type": "position",
  "target_id": "pos_510300",
  "before_json": "{\"quantity\":1000}",
  "after_json": "{\"quantity\":900}",
  "correction_reason": "用户发现录入数量有误"
}
```

修正接口只写修正事实和审计。若修正需要改变当前持仓、现金或快照，用户必须通过持仓编辑或线下交易记录生成新的快照状态。

### 8.5 记录用户操作确认

`POST /api/v1/decisions/{decision_id}/confirmations`

用途：记录用户在线下或第三方交易系统中已经执行的操作，或记录计划、待观察、错误标注。该接口不执行交易。路径中的 `decision_id` 是唯一建议来源，请求体不得再传另一个 `decision_id`。

门禁约束：只有 `decision_records.record_type=formal_trade_advice` 且当前 `confirmation_status` 不为 `not_required` 时允许提交确认；`record_type!=formal_trade_advice`、`confirmation_status=not_required`、`executed_manually` 或 `marked_error` 终态再次确认时，均返回 `BAD_REQUEST`，不得写入账户快照、交易流水或错误案例。

请求：

```json
{
  "confirmation_type": "executed_manually",
  "operation_type": "sell",
  "symbol": "510300",
  "quantity": 300,
  "price": 4.250,
  "executed_at": "2026-05-22T10:00:00+08:00",
  "note": "用户在线下手动卖出。"
}
```

`confirmation_type` 可选值：

| 值 | 说明 | 是否更新账户 |
| --- | --- | --- |
| planned | 仅记录计划 | 否 |
| executed_manually | 已在线下手动执行 | 是 |
| watch | 标记待观察 | 否 |
| marked_error | 标记错误 | 否 |

`confirmation_type` 是用户提交动作，`confirmation_status` 是建议当前状态。提交动作成功后，后端根据动作和业务结果更新建议状态。

当 `confirmation_type=executed_manually` 时，`operation_type` 只能是 `buy / sell / reduce`，后端必须在同一事务写入 `operation_confirmations + position_transactions + positions + portfolio_snapshots + position_snapshots + audit_events`。任一写入失败时，整笔确认回滚。

当 `confirmation_type=marked_error` 时，请求必须包含错误标注字段：

```json
{
  "confirmation_type": "marked_error",
  "actual_outcome": "实际走势与建议方向相反，三日后跌幅超过 8%。",
  "root_cause_tag": "evidence_missed",
  "lesson_learned": "后续同类事件必须检查公告后的成交量变化。",
  "note": "该案例进入错误案例库，后续只能生成规则提案，不能自动改规则。"
}
```

`root_cause_tag` 可选值：`evidence_missed / rule_threshold_issue / analyst_error / user_context_missing / market_exception`。

| confirmation_type | 更新后的 confirmation_status | 说明 |
| --- | --- | --- |
| planned | planned | 仅记录计划，不更新账户 |
| executed_manually | executed_manually | 用户已在线下执行，更新账户 |
| watch | watch | 标记待观察，不更新账户 |
| marked_error | marked_error | 标记错误，写入错误案例库 |

响应示例：

```json
{
  "request_id": "req_20260522_000010",
  "data": {
    "confirmation_id": "conf_20260522_0001",
    "decision_id": "dec_20260522_0002",
    "confirmation_status": "executed_manually",
    "error_case_id": "",
    "transaction_ids": ["tx_20260522_0001"],
    "snapshot_id": "snap_20260522_1000",
    "audit_event_ids": ["audit_010"]
  },
  "meta": {
    "generated_at": "2026-05-22T10:01:00+08:00",
    "rule_version": "v3.0"
  }
}
```

当 `confirmation_type=marked_error` 时，后端必须在同一事务中写入 `operation_confirmations + error_cases + audit_events`，响应中的 `error_case_id` 必须返回新建错误案例 ID。任一写入失败时，整笔确认回滚。

## 9. 证据 API

### 9.1 刷新证据

`POST /api/v1/evidence/refresh`

用途：手动触发证据采集、摘要生成、RAG 文本块写入、多源验证和 VecLite 索引更新。该接口用于本地刷新数据，不生成交易建议。

请求：

```json
{
  "symbol": "510300",
  "refresh_scope": "symbol",
  "include_background": false
}
```

`refresh_scope` 可选值：`symbol / market / all`。

响应示例：

```json
{
  "request_id": "req_20260522_000013",
  "data": {
    "intelligence_item_count": 12,
    "summary_count": 8,
    "rag_chunk_count": 24,
    "verification_count": 3,
    "index_status": "indexed",
    "audit_event_ids": ["audit_013"]
  },
  "meta": {
    "generated_at": "2026-05-22T10:10:00+08:00",
    "rule_version": "v3.0"
  }
}
```

写入要求：写入 `intelligence_items + intelligence_summary + rag_chunks + source_verifications + audit_events`。新建 `rag_chunks` 初始为 `pending`；VecLite 写入成功后标记为 `indexed`；VecLite 索引失败时不得回滚 SQLite 事实数据，应将相关 `rag_chunks.index_status` 标记为 `failed`。

### 9.2 查询证据列表

`GET /api/v1/evidence?symbol=510300&level=A&verification_status=satisfied`

响应字段：

| 字段 | 说明 |
| --- | --- |
| evidence_id | 证据编号 |
| source_name | 来源名称 |
| source_level | S / A / B / C | 默认返回 S/A/B；C 级只可在 `evidence_role=background` 时返回 |
| evidence_role | formal / background |
| published_at | 发布时间 |
| captured_at | 抓取时间 |
| summary | 摘要 |
| original_url | 原文链接 |
| content_hash | 内容 hash |
| time_weight | 时效权重（0–1，越近越高） |
| relevance_score | 相关度得分（0–1） |
| independent_source_count | 独立信源数量 |
| high_grade_independent_source_count | S/A 级独立信源数量（重大事件须 ≥ 2） |

### 9.3 获取多源验证状态

`GET /api/v1/evidence/verification?symbol=510300&event_id=evt_001`

查询参数：`symbol` 与 `event_id` 均可单独使用；同时提供时返回同时匹配两者的最新验证记录；无匹配记录返回 `NOT_FOUND`。

响应字段：

| 字段 | 说明 |
| --- | --- |
| verification_status | satisfied / failed / background_only |
| independent_source_count | 独立信源数量 |
| high_grade_independent_source_count | S/A 级独立信源数量 |
| highest_source_level | 最高信源等级 |
| latest_published_at | 最新发布时间 |
| evidence_ids | 证据编号列表 |

### 9.4 本地知识导入治理 API（P46）

P46 提供本地研究材料的两阶段导入：先校验和脱敏预览，再由用户显式写入本地背景事实。该能力只复用 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 和 `audit_events`，不新增数据库表，不生成正式交易建议。

#### 9.4.1 校验本地知识导入

`POST /api/v1/local-knowledge/imports/validate`

请求：

```json
{
  "source_label": "local_research_notes",
  "default_symbol": "510300",
  "rows": [
    {
      "title": "510300 估值观察",
      "text": "本地研究记录正文",
      "symbol": "510300",
      "as_of_date": "2026-06-17",
      "tags": ["估值", "本地研究"]
    }
  ]
}
```

响应：

```json
{
  "request_id": "req_20260617_000001",
  "data": {
    "import_batch_id": "lk_batch_abc123",
    "summary": {
      "total_count": 1,
      "valid_count": 1,
      "warning_count": 0,
      "blocking_count": 0
    },
    "rows": [
      {
        "row_number": 1,
        "status": "valid",
        "symbol": "510300",
        "title_preview": "510300 估值观察",
        "text_preview": "本地研究记录正文",
        "content_hash": "hash_abc123",
        "risks": []
      }
    ],
    "index_plan": {
      "rag_chunk_count": 1,
      "index_status": "pending"
    },
    "safety_note": "本地知识导入仅写入本地背景材料，不接券商、不交易、不外部推送、不自动确认、不自动应用规则。"
  }
}
```

约束：

- `import_batch_id` 必须由规范化后的 `source_label`、`default_symbol`、行顺序和内容 hash 派生，保证确认阶段可重算。
- 返回的 `title_preview` 与 `text_preview` 必须经过脱敏；不得返回完整 key、私有路径、原始 SQL、私钥或完整 prompt。
- `blocking_count > 0` 时，前端不得启用确认写入。
- `index_plan.index_status` 初始为 `pending`，只表达后续需要重建或刷新索引。

#### 9.4.2 确认本地知识导入

`POST /api/v1/local-knowledge/imports/confirm`

请求：

```json
{
  "import_batch_id": "lk_batch_abc123",
  "confirm_reason": "人工确认导入为本地背景材料",
  "source_label": "local_research_notes",
  "default_symbol": "510300",
  "rows": [
    {
      "title": "510300 估值观察",
      "text": "本地研究记录正文",
      "symbol": "510300",
      "as_of_date": "2026-06-17",
      "tags": ["估值", "本地研究"]
    }
  ]
}
```

响应：

```json
{
  "request_id": "req_20260617_000002",
  "data": {
    "import_batch_id": "lk_batch_abc123",
    "intelligence_item_count": 1,
    "summary_count": 1,
    "rag_chunk_count": 1,
    "verification_count": 1,
    "audit_event_ids": ["audit_046_001"],
    "index_status": "pending",
    "safety_note": "本地知识导入仅写入本地背景材料，不接券商、不交易、不外部推送、不自动确认、不自动应用规则。"
  }
}
```

写入要求：

- confirm 必须重新校验 rows，并用请求中的 `source_label`、`default_symbol` 和 rows 重算 `import_batch_id`；不匹配时返回 `CONFLICT` 或 `INVALID_STATE`。
- 任一行存在 blocking 风险时不得写入。
- 成功写入必须处于同一事务：`intelligence_items + intelligence_summary + rag_chunks + source_verifications + audit_events` 整体成功或整体回滚。
- 默认 `source_level=C`、`evidence_role=background`、`source_verifications.verification_status=background_only`、`rag_chunks.index_status=pending`。
- 本地知识只作为背景材料，不得提升为 formal 证据，不得绕过多源验证、规则裁决或用户人工复核。

## 10. 规则与审计 API

### 10.1 获取当前规则

`GET /api/v1/rules/current`

响应字段：

| 字段 | 说明 |
| --- | --- |
| rule_version | 规则版本 |
| fundamental_rules | 根本规则 |
| arbitration_priority | 裁决优先级 |
| market_thresholds | 市场状态阈值 |
| sop_list | SOP 列表 |

### 10.2 查询规则提案

`GET /api/v1/rule-proposals?status=pending`

生成入口说明：本系统不提供前端手动创建规则提案接口。规则提案只能由错误案例分析、EvolutionGraph 或受控内部任务生成，写入 `rule_proposals + notifications + audit_events`；`PUT /api/v1/settings` 不得创建或绕过规则提案。

响应示例：

```json
{
  "request_id": "req_20260522_000008",
  "data": {
    "items": [
      {
        "proposal_id": "rp_20260522_0001",
        "status": "pending_user_confirm",
        "source_error_case_id": "err_20260520_0001",
        "title": "提高单信源重大利空的审慎等级",
        "reason": "历史错误案例显示，单信源利空在高热度阶段容易被低估。",
        "before_rule": {
          "rule_id": "R-7",
          "version": "v3.0",
          "content": "单信源重大信息只能作为背景材料。"
        },
        "after_rule": {
          "rule_id": "R-7",
          "version": "v3.1-proposal",
          "content": "单信源重大利空且媒体热度异常时，进入冻结观察。"
        },
        "impact_scope": ["consultation", "daily_discipline"],
        "sample_count": 2,
        "effect_validation": {
          "validation_status": "insufficient",
          "sample_count": 2,
          "sample_window": "2026-Q2",
          "representativeness_status": "needs_more_samples",
          "overfit_risk": "high",
          "replay_result": "unknown",
          "guardrail_decision": "rejected",
          "validation_link": "/rule-effect-validations/val_20260616_0001",
          "safety_note": "规则效果验证只用于本地规则治理，不会自动应用规则或执行交易。"
        },
        "risk_notes": ["样本数量不足，当前不得进入守门人审计或最终应用"],
        "created_at": "2026-05-22T11:00:00+08:00"
      }
    ],
    "total": 1
  },
  "meta": {
    "generated_at": "2026-05-22T11:05:00+08:00",
    "rule_version": "v3.0"
  }
}
```

### 10.3 用户确认规则提案进入审计

`POST /api/v1/rule-proposals/{proposal_id}/confirm`

用途：用户第一次确认提案可以进入守门人审计。该接口不应用正式规则。

送审门禁：`sample_count<3` 的样本不足提案不得进入守门人审计，接口返回 `BAD_REQUEST`，提案保持 `draft` 或 `pending_user_confirm`；除非后续 EvolutionGraph 或受控内部任务生成满足样本条件的新提案版本。若存在 P36 规则效果验证结果，`validation_status=insufficient/failed/needs_more_samples`、`overfit_risk=high`、`replay_result=failed` 或 `guardrail_decision=rejected` 也必须阻止送审或返回 `needs_user_review`，不得视为可应用规则。

请求：

```json
{
  "confirm": true,
  "note": "同意进入守门人审计。"
}
```

响应示例（守门人返回 `needs_user_review` 后，最终响应状态回到 `pending_user_confirm`）：

```json
{
  "request_id": "req_20260522_000009",
  "data": {
    "proposal_id": "rp_20260522_0001",
    "status": "pending_user_confirm",
    "audit_result": "needs_user_review",
    "audit_summary": "提案具备历史案例依据，但影响范围涉及每日纪律与咨询模式，需要用户复核后再生效。",
    "rule_change_diff": {
      "rule_id": "R-7",
      "before_version": "v3.0",
      "after_version": "v3.1-proposal",
      "before_text": "单信源重大信息只能作为背景材料。",
      "after_text": "单信源重大利空且媒体热度异常时，进入冻结观察。"
    },
    "audit_events": ["audit_002", "audit_003"]
  },
  "meta": {
    "generated_at": "2026-05-22T11:10:00+08:00",
    "rule_version": "v3.0"
  }
}
```

状态流转：

- `confirm=true` 且 `sample_count>=3`：`pending_user_confirm -> under_gatekeeper_audit`。
- `confirm=true` 且 `sample_count<3`：返回 `BAD_REQUEST`，提案保持 `draft` 或 `pending_user_confirm`，不写 `gatekeeper_audits`。
- `confirm=true` 且 P36 效果验证缺失或未通过：返回 `BAD_REQUEST` 或 `needs_user_review`，不得进入最终应用路径；缺失验证只能展示为 `not_evaluated`，不得假定安全。
- 守门人审计通过：`under_gatekeeper_audit -> pending_final_confirm`。
- 守门人审计拒绝：`under_gatekeeper_audit -> rejected`。
- 守门人审计需要用户复核：`under_gatekeeper_audit -> pending_user_confirm`。
- `audit_result=approved / rejected / needs_user_review` 是审计结果，不是提案状态。

### 10.4 用户最终确认应用规则

`POST /api/v1/rule-proposals/{proposal_id}/final-confirm`

用途：用户在守门人审计通过后，最终确认是否把提案写入正式规则版本。

请求：

```json
{
  "confirm": true,
  "note": "确认按审计后的提案更新规则。"
}
```

响应示例：

```json
{
  "request_id": "req_20260522_000011",
  "data": {
    "proposal_id": "rp_20260522_0001",
    "status": "applied",
    "applied_rule_version": "v3.1",
    "final_confirmed_at": "2026-05-22T11:30:00+08:00",
    "audit_events": ["audit_004"],
    "created_rule_version": "v3.1"
  },
  "meta": {
    "generated_at": "2026-05-22T11:30:00+08:00",
    "rule_version": "v3.1"
  }
}
```

状态流转与事务要求：

- `confirm=true`：`pending_final_confirm -> applied`，同时写入 `rule_versions`。
- `confirm=false`：`pending_final_confirm -> rejected`，不得写入正式规则。
- `sample_count<3` 的样本不足提案不得执行最终确认；即使状态异常进入 `pending_final_confirm`，接口也必须返回 `BAD_REQUEST`，不得写入 `rule_versions`。
- P36 效果验证不通过、过拟合高、历史回放不利或验证版本与当前提案版本不一致时，不得执行最终确认；接口必须返回稳定错误，提案状态保持不变，不写入 `rule_versions`。
- `confirm=true` 必须在同一事务更新 `rule_proposals`、归档旧 active `rule_versions`、创建新 active `rule_versions`、写入 `audit_events`。
- 任一写入失败时，提案保持 `pending_final_confirm`。

规则提案完整状态机：

| 事件 | 前态 | 审计结果 / 参数 | 后态 | 写入表 | 说明 |
| --- | --- | --- | --- | --- | --- |
| 生成提案 | 无 | - | `draft` 或 `pending_user_confirm` | `rule_proposals`、`notifications`、`audit_events` | 不写正式规则；待用户处理的提案必须生成应用内通知。 |
| 用户确认送审 | `pending_user_confirm` | `confirm=true` 且 `sample_count>=3` | `under_gatekeeper_audit` | `rule_proposals`、`audit_events` | 随后执行守门人审计。 |
| 样本不足提案 | `draft` / `pending_user_confirm` | `sample_count<3` | 原状态不变 | `audit_events` | 返回 `BAD_REQUEST`，不得送审；如需修改，由后续 EvolutionGraph 或受控内部任务生成满足样本条件的新提案版本。 |
| 用户拒绝送审 | `pending_user_confirm` | `confirm=false` | `rejected` | `rule_proposals`、`audit_events` | 不写正式规则。 |
| 守门人通过 | `under_gatekeeper_audit` | `approved` | `pending_final_confirm` | `gatekeeper_audits`、`rule_proposals`、`audit_events` | 等待用户最终确认。 |
| 守门人否决 | `under_gatekeeper_audit` | `rejected` | `rejected` | `gatekeeper_audits`、`rule_proposals`、`audit_events` | 不写正式规则。 |
| 需要用户复核 | `under_gatekeeper_audit` | `needs_user_review` | `pending_user_confirm` | `gatekeeper_audits`、`rule_proposals`、`audit_events` | 用户可放弃或重新送审；如需修改，由后续 EvolutionGraph 或受控内部任务生成新提案版本。 |
| 最终确认应用 | `pending_final_confirm` | `confirm=true` | `applied` | `rule_versions`、`rule_proposals`、`audit_events` | 创建新 active 规则版本。 |
| 最终拒绝应用 | `pending_final_confirm` | `confirm=false` | `rejected` | `rule_proposals`、`audit_events` | 不写正式规则。 |
| 终态重复操作 | `rejected` / `applied` | 任意确认动作 | 原状态不变 | `audit_events` | 返回 `BAD_REQUEST`，不写 `rule_versions`。 |

### 10.5 获取规则提案效果验证

`GET /api/v1/rule-proposals/{proposal_id}/effect-validation`

用途：读取指定规则提案最近一次本地效果验证事实。无验证结果时返回 `NOT_FOUND` 或 typed empty state；前端必须展示 `not_evaluated`，不得伪造验证指标。

响应字段：

| 字段 | 说明 |
| --- | --- |
| validation_id | 验证事实 ID |
| proposal_id | 规则提案 ID |
| candidate_rule_version | 候选规则版本 |
| validation_status | not_evaluated / insufficient / passed / failed / needs_more_samples / needs_user_review |
| sample_count | 样本数量 |
| sample_window | 样本窗口 |
| representativeness_status | 样本代表性状态 |
| overfit_risk | low / medium / high |
| replay_result | passed / failed / mixed / unknown |
| guardrail_decision | passed / rejected / needs_user_review |
| source_explanation | 来源解释 JSON，包含错误案例、决策、风险预警等本地事实线索 |
| metrics | 命中、误判、缺证据、降级、风险预警等指标快照 |
| risk_notes | 风险说明 |
| related_error_cases / related_decision_ids / related_risk_alert_ids / related_audit_event_ids | 本地追踪 ID 列表 |
| validation_link | 本地详情链接 |
| safety_note | 固定安全文案，说明不会自动应用规则或执行交易 |
| created_at / updated_at | 创建和更新时间 |

### 10.6 刷新规则提案效果验证

`POST /api/v1/rule-proposals/{proposal_id}/effect-validation`

用途：基于本地规则提案、错误案例、决策记录、复盘、风险预警和审计事实刷新效果验证。请求体可包含 `sample_window`；服务端不得信任前端传入的样本指标来绕过门禁。

事务边界：刷新验证时，系统只能写入 `rule_effect_validations` 和对应 `audit_events`。不得写入或修改 `positions`、`portfolio_snapshots`、`operation_confirmations`、`position_transactions`、`rule_versions`、broker state、orders 或 external notifications。

### 10.7 查询应用后规则效果追踪

`GET /api/v1/rule-effect-tracking?rule_version=v3.2&proposal_id=rp_001&period=2026-Q3`

用途：查询已应用规则版本的后续效果追踪事实。查询参数均可选；返回按本地事实过滤后的追踪列表。

响应字段：

| 字段 | 说明 |
| --- | --- |
| tracking_id | 追踪事实 ID |
| applied_rule_version | 已应用规则版本 |
| proposal_id | 来源提案 ID，可为空 |
| period | 追踪周期 |
| hit_count / misjudgment_count / missing_evidence_count / degraded_count / risk_alert_count | 效果指标 |
| trend_direction | improved / flat / worsened / unknown |
| metrics | 指标快照 JSON |
| related_proposal_ids / related_audit_event_ids / related_risk_alert_ids | 本地追踪 ID 列表 |
| safety_note | 固定安全文案，说明不会自动回滚规则或执行交易 |
| created_at / updated_at | 创建和更新时间 |

事务边界：应用后追踪只作为复盘和规则治理事实；不得自动创建、回滚或替换 active `rule_versions`，不得执行交易、连接券商或外部推送。

### 10.8 获取审计事件

`GET /api/v1/audit-events?from=2026-05-01&to=2026-05-22&type=decision`

响应示例：

```json
{
  "request_id": "req_20260522_000007",
  "data": {
    "items": [
      {
        "audit_event_id": "audit_001",
        "event_id": "audit_001",
        "request_id": "req_20260522_000002",
        "workflow_type": "consultation",
        "node_name": "RuleArbitrationNode",
        "actor": "system",
        "action": "generate_decision",
        "node_action": "arbitrate_rule",
        "status": "success",
        "before_state": "pending",
        "after_state": "pending",
        "rule_version": "v3.0",
        "snapshot_id": "snap_20260522_0900",
        "input_ref": "ctx_req_20260522_000002",
        "output_ref": "dec_20260522_0002",
        "error_code": "",
        "created_at": "2026-05-22T09:35:00+08:00"
      }
    ],
    "total": 1
  },
  "meta": {
    "generated_at": "2026-05-22T09:40:00+08:00",
    "rule_version": "v3.0"
  }
}
```

说明：响应中的 `audit_event_id` 与 `event_id` 值相同，均对应数据模型 `audit_events.audit_event_id`；保留 `event_id` 作为前端展示别名。`action` 表示业务动作，`node_action` 表示 Eino 节点动作。

## 11. 设置 API

### 11.1 获取能力圈配置

`GET /api/v1/settings/capability`

响应示例：

```json
{
  "request_id": "req_20260522_000014",
  "data": {
    "capability_id": "cap_default",
    "asset_types": ["ETF", "index_fund", "bond_fund"],
    "symbols": ["510300"],
    "excluded_symbols": [],
    "strategy_scope": ["long_term_allocation", "discipline_review"],
    "updated_at": "2026-05-22T09:00:00+08:00"
  },
  "meta": {
    "generated_at": "2026-05-22T10:20:00+08:00",
    "rule_version": "v3.0"
  }
}
```

### 11.2 更新能力圈配置

`PUT /api/v1/settings/capability`

请求：

```json
{
  "asset_types": ["ETF", "index_fund", "bond_fund"],
  "symbols": ["510300"],
  "excluded_symbols": [],
  "strategy_scope": ["long_term_allocation", "discipline_review"],
  "note": "仅分析长期配置范围内的 ETF。"
}
```

写入要求：写入 `capability_configs + audit_events`。如果同一标的同时出现在 `symbols` 与 `excluded_symbols`，返回 `BAD_REQUEST`。审计事件必须使用 `action=update_capability`，并在 `before_state`、`after_state` 中记录变更前后摘要。

### 11.3 获取系统设置状态

`GET /api/v1/settings/system`

用途：查看本地配置和依赖状态，不返回完整密钥。

响应字段：

| 字段 | 说明 |
| --- | --- |
| sqlite_status | ok / failed |
| sqlite_path | SQLite 文件路径，可脱敏 |
| veclite_status | ok / unavailable / rebuilding / failed |
| veclite_path | VecLite 文件路径，可脱敏 |
| deepseek_status | configured / missing / unavailable |
| deepseek_model | DeepSeek 或兼容模型名，可为空；不得包含 key |
| data_sources | 数据源开关与状态 |
| log_level | 当前日志级别 |

### 11.4 更新用户设置

`PUT /api/v1/settings`

约束：

- 非规则配置：通知、页面偏好、普通数据源配置，只能通过 `PUT /api/v1/settings` 保存，写入 `user_settings + audit_events`。
- 能力圈配置：只通过 `PUT /api/v1/settings/capability` 保存，写入 `capability_configs + audit_events`，不得混入 `PUT /api/v1/settings`。
- 规则类配置：根本规则、裁决优先级、核心阈值、SOP、规则版本内容，必须生成规则提案并进入审计，不得直接生效。
- 不得在响应或审计事件中明文返回 DeepSeek API Key。
- P37 起，分析报告可携带 `prompt_version`、`model`、`input_summary`、`output_summary`、`parse_status`、`quality_status` 等可选字段；LLM 错误分类使用 missing_key、timeout、http_error、empty_response、parse_error、quality_failed、unavailable 等稳定值进入审计或降级诊断，不返回明文 key。

### 11.5 市场数据刷新

`POST /api/v1/market/refresh`

用途：同步刷新行情、估值、流动性和情绪指标，并写入 `market_snapshots`。

请求：

```json
{
  "symbols": ["510300"]
}
```

说明：`as_of_date` 字段已预留但当前不支持指定历史交易日刷新；请求中提供合法或非法 `as_of_date` 都返回 `BAD_REQUEST`，避免静默忽略。

响应字段：

| 字段 | 说明 |
| --- | --- |
| refreshed_count | 成功写入的快照数量 |
| failed_symbols | 失败标的列表 |
| latest_snapshot_ids | 最新市场快照 ID 列表 |
| audit_event_ids | 审计事件 ID 列表 |

部分成功策略：只要至少 1 个标的成功写入 `market_snapshots`，接口返回 200，并通过 `failed_symbols` 返回失败标的及原因；所有标的都失败时，按失败类型返回 `DATA_SOURCE_UNAVAILABLE`、`DATA_STALE` 或 `MARKET_SNAPSHOT_WRITE_FAILED`。

四类结果约束：

| 结果 | HTTP / 错误码 | 写入要求 | 审计要求 | 旧快照展示 |
| --- | --- | --- | --- | --- |
| 全部成功 | 200 | 所有请求标的写入 `market_snapshots` | `audit_events.status=success` | 展示最新快照 |
| 部分成功 | 200 | 成功标的写入，失败标的不写快照 | `audit_events.status=degraded`，`output_ref` 记录 `failed_symbols` | 成功标的展示最新快照，失败标的可展示旧快照并标记时间 |
| 全部失败：数据源不可用 | 503 / `DATA_SOURCE_UNAVAILABLE` | 不写 `market_snapshots` | `audit_events.status=failed` | 可展示旧快照，但必须标记 stale 或 missing |
| 全部失败：数据过期 | 409 / `DATA_STALE` | 不写 `market_snapshots` | `audit_events.status=failed` | 可展示旧快照，但不得生成交易类建议 |
| 写入失败 | 500 / `MARKET_SNAPSHOT_WRITE_FAILED` | 市场快照事务回滚，不留下部分快照；失败审计事件仍需写入，且应使用独立事务避免被快照事务回滚 | `audit_events.status=failed` | 不展示本次刷新结果 |

写入要求：写入 `market_snapshots + audit_events`。数据源不可用返回 `DATA_SOURCE_UNAVAILABLE`；数据过期返回 `DATA_STALE`；快照写入失败返回 `MARKET_SNAPSHOT_WRITE_FAILED`。上述失败状态均不得生成交易类建议。

### 11.6 获取最新市场快照

`GET /api/v1/market/snapshots/latest?symbol=510300`

用途：读取指定标的最新市场快照，供驾驶舱、决策详情和设置页展示数据状态。

响应字段：

| 字段 | 说明 |
| --- | --- |
| market_snapshot_id | 市场快照 ID |
| symbol | 标的代码 |
| trade_date | 交易日期 |
| market_metrics | 原始指标 |
| close_price | 收盘价或净值 |
| turnover_rate | 换手率 |
| pe_percentile | PE 分位 |
| pb_percentile | PB 分位 |
| liquidity_state | normal / warning / danger |
| sentiment_state | cold / neutral / hot / extreme |
| data_status | fresh / stale / missing |

结构化字段：当真实公开 provider 已通过市场刷新写入时，`market_metrics.metadata.p88_structured_fields.capital_flow` 可包含 `date`、`net_inflow`、`net_outflow`、`raw_net_flow`。其中 `raw_net_flow` 是公开 H5 历史资金流向的日净流向原值；正值映射到 `net_inflow`，负值绝对值映射到 `net_outflow`。

## 12. 通知 API

### 12.1 获取应用内通知列表

`GET /api/v1/notifications`

用途：读取本地应用内通知中心状态，供前端轮询展示。该接口只读取本地 `notifications` 表，不发送邮件、短信、系统 Push、Webhook 或 WebSocket 消息。成功响应使用统一信封，`data` 包含以下字段。

响应字段：

| 字段 | 说明 |
| --- | --- |
| items | 通知列表，按 `created_at` 倒序 |
| unread_count | 未读通知数量 |

`items` 元素字段：

| 字段 | 说明 |
| --- | --- |
| notification_id | 通知 ID |
| type | 通知类型，如 `data_source_failure`、`vector_index_failure`、`rule_proposal_pending` |
| severity | info / warning / critical |
| title | 通知标题 |
| message | 通知正文 |
| source_type | 来源类型，用于去重和追踪 |
| source_id | 来源 ID；同一 `type/source_type/source_id` 的未读通知应去重 |
| read_at | 已读时间；为空表示未读 |
| created_at | 创建或最近刷新时间 |

### 12.2 标记单条通知已读

`POST /api/v1/notifications/{notification_id}/read`

用途：把指定本地通知标记为已读。成功响应使用统一信封，`data` 为 `{"ok": true}`。

### 12.3 标记全部通知已读

`POST /api/v1/notifications/read-all`

用途：把所有本地未读通知标记为已读。成功响应使用统一信封，`data` 为 `{"ok": true}`。

安全边界：通知 API 只改变本地通知已读状态，不执行交易，不自动应用规则，不触发任何外部推送。

## 13. 复盘 API

### 13.1 获取复盘汇总

`GET /api/v1/review/summary?period=monthly`

用途：按 `monthly` 或 `quarterly` 周期汇总建议、用户确认、错误案例、规则提案和审计事件。未传或传空时默认 `monthly`。当复盘存在降级或缺失证据时，接口会写入本地应用内通知 `review_degraded`；该通知只用于 UI 提醒，不触发外部推送、交易或规则应用。

响应字段：

| 字段 | 说明 |
| --- | --- |
| period | 周期类型：monthly / quarterly |
| decision_count | 建议数量 |
| confirmation_count | 用户确认数量 |
| executed_manually_count | 已手动执行数量 |
| planned_count | 记录计划数量 |
| error_case_count | 错误案例数量 |
| rule_proposal_count | 规则提案数量 |
| audit_event_count | 审计事件数量 |
| rule_hit_count | 规则命中次数 |
| misjudgment_count | 误判案例数量 |
| missing_evidence_count | 缺失证据数量 |
| degraded_count | 降级记录数量 |
| ops_status | 运维状态对象（见下） |
| rule_suggestions | 规则建议列表（只读展示，不可自动应用） |
| attribution_summaries | 可追溯归因摘要，来自本地决策、证据状态和确认状态 |
| recurring_error_tags | 错误案例 root_cause_tag 计数 |
| missing_evidence_themes | 缺证据状态主题计数 |
| rule_proposal_outcomes | 规则提案状态与最新守门人审计结果 |
| degraded_workflows | 降级工作流追踪列表 |
| rule_effect_tracking | 应用后规则效果追踪列表；只读展示，不自动回滚或应用规则 |
| tracking_links | 追踪入口列表，指向审计、决策、规则提案、错误案例或规则效果追踪 |

`ops_status` 字段结构：

| 子字段 | 说明 |
| --- | --- |
| data_source_status | success / degraded / failed / empty / unknown |
| index_status | success / degraded / failed / missing / unknown |
| review_status | success / degraded / failed / empty / unknown |
| explanation | 可选说明文本 |

`rule_suggestions` 列表元素字段：

| 子字段 | 说明 |
| --- | --- |
| proposal_id | 规则提案 ID |
| title | 提案标题 |
| status | 提案当前状态 |
| reason | 建议原因 |
| can_auto_apply | 始终为 false；规则变更须守门人审计与用户最终确认 |

`rule_effect_tracking` 列表元素字段同 `GET /api/v1/rule-effect-tracking` 的追踪 DTO。若复盘周期内样本不足，接口必须展示缺少事实或 `unknown` 趋势，不得宣称规则改善或恶化。

`tracking_links` 列表元素字段：

| 子字段 | 说明 |
| --- | --- |
| type | 记录类型：rule_proposal / decision / error_case / audit_event / rule_effect_tracking |
| id | 记录 ID |
| label | 展示标签 |

### 13.2 决策闭环解释 API（P47）

P47 新增只读决策闭环解释接口，用于把建议、用户确认、本地线下记录、风险线索、复盘线索和审计线索串成同一条解释链。接口只读取现有本地事实，不写入 `decision_records`、`operation_confirmations`、`position_transactions`、`risk_alerts`、`audit_events`、`notifications`、规则版本或任何外部通道。

#### 13.2.1 列出最近决策闭环

`GET /api/v1/decision-loops?symbol=510300&limit=10`

查询参数：

| 参数 | 说明 |
| --- | --- |
| symbol | 可选，按标的过滤。 |
| limit | 可选，默认 10，最大 50；非正整数返回 `BAD_REQUEST`。 |

响应字段：

| 字段 | 说明 |
| --- | --- |
| items | 决策闭环列表，按决策时间倒序。 |
| total | 本次返回条数。 |
| safety_note | 固定只读边界说明。 |

`items[]` 字段：

| 字段 | 说明 |
| --- | --- |
| decision_id | 决策 ID。 |
| symbol | 标的，可为空。 |
| generated_at | 决策生成时间。 |
| final_verdict_status | 最终裁决状态。 |
| final_verdict_text | 最终裁决文本；只展示安全摘要。 |
| confirmation_status | 当前确认状态。 |
| loop_status | open / planned / recorded / reviewed / incomplete。 |
| stages | recommendation / confirmation / manual_record / risk_review / review 五类阶段。 |
| manual_actions | 用户确认和本地流水摘要，不返回 raw payload。 |
| risk_links | 风险预警只读导航链接。 |
| review_links | 错误案例或复盘只读导航链接。 |
| audit_links | 审计事件只读导航链接。 |
| missing_links | 缺口说明，如缺少用户确认、本地流水、风险线索或复盘线索。 |
| safety_note | 单条闭环的只读边界说明。 |

`stages[]` 字段：

| 字段 | 说明 |
| --- | --- |
| stage | recommendation / confirmation / manual_record / risk_review / review。 |
| status | complete / pending / not_required / missing / degraded。 |
| label | 页面展示标签。 |
| summary | 安全摘要。 |
| ref_type | 可选引用类型。 |
| ref_id | 可选引用 ID。 |
| at | 可选时间。 |

`manual_actions[]` 字段只允许包含 `confirmation_id`、`confirmation_type`、`operation_type`、`symbol`、`quantity`、`price`、`fees`、`executed_at`、`transaction_ids` 和 `note_preview`。`payload_json`、完整 key、私有路径、原始 SQL、完整 prompt、供应商原始响应和外部订单类信息不得返回。

#### 13.2.2 获取单条决策闭环

`GET /api/v1/decision-loops/{decision_id}`

用途：返回单条 `DecisionLoopItem`，字段同列表 `items[]`。不存在的 `decision_id` 返回 `NOT_FOUND`。该接口与列表接口一样只读，不创建确认、不写交易流水、不改变风险 SOP、不创建通知、不应用规则。

### 13.3 数据源质量回归 API（P48）

P48 新增本地数据源质量回归接口，用于验证 source health/freshness 分类、失败类别和脱敏摘要。接口不触发 collector，不刷新市场快照，不重建索引，不调用 LLM，不创建通知，不修改账户、确认、风险 SOP 或规则。

`GET /api/v1/data-source-quality/regression?mode=fixture&symbol=000300`

查询参数：

| 参数 | 说明 |
| --- | --- |
| mode | 可选，`fixture` 或 `current`；默认 `fixture`。`fixture` 使用确定性本地样本，`current` 只读评估已有市场快照中的 P34 source health。 |
| symbol | 可选，仅 `current` 模式用于读取指定标的的最新市场快照。 |

响应字段：

| 字段 | 说明 |
| --- | --- |
| mode | 实际运行模式：fixture / current |
| status | passed / degraded / failed |
| generated_at | 生成时间 |
| summary | 紧凑摘要 |
| cases | 回归 case 列表 |
| missing_categories | 降级或缺失的数据类别 |
| safety_note | 固定安全边界说明 |

`cases[]` 字段：

| 字段 | 说明 |
| --- | --- |
| case_id | case 标识 |
| source_name | 来源名称 |
| source_level | A / B / C 等来源等级 |
| source_type | 来源类型 |
| data_category | 数据类别 |
| expected_freshness | 期望 freshness |
| actual_freshness | 实际 freshness |
| status | passed / degraded / failed |
| data_date | 可选数据日 |
| failure_category | 可选失败类别 |
| affected_symbols | 影响标的 |
| diagnostic_preview | 脱敏诊断预览 |

安全约束：`diagnostic_preview`、`summary` 和 `safety_note` 不得包含完整 key、私有路径、原始 SQL、完整 prompt、raw HTTP、HTTP status line、private key 或供应商原始响应。`fixture` 模式不得访问公网；`current` 模式只读取 `market_snapshots.market_metrics_json.metadata.p34_source_health`。

### 13.4 当前数据门禁处置 API（P67）

P67 在 P66 current data policy 之上新增本地人工处置记录。该 API 不改变 P66 policy verdict，不刷新数据、不修复源、不调用外部 provider、不创建通知、不修改账户、确认、风险 SOP 或规则。

`GET /api/v1/data-source-quality/gate-resolution?symbol=000300`

返回当前 P66 policy、canonical `policy_fingerprint`、`release_claim_state`、`clean_data_claim_allowed`、active resolution 和固定 allowed/prohibited claims。GET 只读，不写 `audit_events`。

`GET /api/v1/data-source-quality/resolutions?symbol=000300&status=active`

返回本地处置记录列表，按 `created_at DESC` 排序。支持按 `symbol` 与 `status` 过滤。

`POST /api/v1/data-source-quality/resolutions`

请求：

```json
{
  "symbol": "000300",
  "resolution_type": "scope_exclusion",
  "scope": "本次 release clean claim 排除 current local data health",
  "reason": "当前本地数据源存在降级，发布材料只声明有限范围",
  "release_impact": "不得声明 current data healthy",
  "evidence_ref": "docs/release/acceptance/p67"
}
```

规则：

- `policy=passed` 不允许创建 resolution。
- `policy=blocked` 只允许 `resolution_type=scope_exclusion`。
- `policy=waiver_required` 允许 `waiver` 或 `scope_exclusion`。
- 同一 `symbol + policy_fingerprint` 只允许一个 active resolution；同类型重复请求复用 active record，不同类型 active 冲突返回 `CONFLICT`。
- `scope`、`reason`、`release_impact` 必填，所有文本必须脱敏。

`POST /api/v1/data-source-quality/resolutions/{resolution_id}/retire`

只将本地 resolution 标记为 `retired`，不改变 source health、P66 policy、账户、规则或发布事实。POST create/retire 成功时写入脱敏 `audit_events`，`action=run_local_task`，`node_action=data_quality_gate_resolution_create|retire`。

`release_claim_state` 取值：

| 值 | 含义 |
| --- | --- |
| `pass` | P66 policy pass，可声明当前本地数据门禁通过。 |
| `requires_resolution` | 当前 policy 需要人工处置，不能声明 current data clean。 |
| `resolved_with_waiver` | 已记录 waiver；仍不得称为 clean pass。 |
| `resolved_with_scope_exclusion` | 已记录范围排除；只能声明 current local data health 已排除在 clean claim 外。 |

## 14. 索引维护 API

### 14.1 重建 VecLite 索引

`POST /api/v1/evidence/rebuild-index`

用途：从 SQLite 的 `rag_chunks` 与 `intelligence_summary` 重建 VecLite 索引。

请求：

```json
{
  "scope": "stale_only"
}
```

`scope` 可选值：`stale_only / all`。

响应字段：

| 字段 | 说明 |
| --- | --- |
| indexed_count | 成功索引文本块数量 |
| skipped_count | 跳过文本块数量（内容未变更） |
| last_rebuild_at | 最近重建时间 |
| index_health | 索引健康对象（见下） |
| audit_event_ids | 审计事件 ID 列表 |

`index_health` 字段结构：

| 子字段 | 说明 |
| --- | --- |
| status | healthy / missing / corrupted / incompatible / degraded |
| path | 索引文件或目录路径 |
| version | 索引版本号 |
| chunk_count | 已索引文本块数量 |
| rebuildable | 是否可从 SQLite 重建 |
| degraded_reason | 降级原因说明（可选） |

约束：VecLite 重建失败不得影响 SQLite 事实数据；失败块必须更新 `rag_chunks.index_status`。

## 15. P74 知识与数据准备度 API

`GET /api/v1/knowledge-readiness?symbol=510300`

用途：只读返回指定标的的内置知识引用、数据依赖矩阵、功能影响、LLM 上下文摘要和安全边界。该接口不刷新外部数据、不调用 LLM、不写审计、不改变规则、账户、持仓或确认记录。

响应 `data` 字段：

| 字段 | 说明 |
| --- | --- |
| `symbol` | 请求标的；为空时默认本地主路径 `510300` |
| `status` | `ready / degraded / blocked` |
| `symbol_profile` | 标的画像；未知标的必须 `known=false` |
| `knowledge_references` | 内置知识引用列表，包含稳定 ID、类别、规则映射、LLM eligibility 和 formal evidence boundary |
| `data_dependencies` | 数据依赖矩阵 |
| `feature_impacts` | 降级或阻断对功能的影响 |
| `llm_context_summary` | 可进入 LLM 的脱敏摘要；不得包含完整 prompt、密钥、路径或 raw provider payload |
| `safety_notes` | 安全边界说明 |

`data_dependencies[].category` 当前包含：`symbol_profile`、`fund_profile`、`tracked_index`、`market_price`、`valuation_percentiles`、`liquidity`、`sentiment_proxy`、`active_rule`、`formal_evidence`、`rag_index`、`llm_context`。

规则：

- `symbol_profile` 未知时整体 `status=blocked`，不得伪造画像。
- `active_rule` 缺失时整体至少 `degraded`，不得生成交易确认。
- `formal_evidence` 必须满足多源高等级验证才可 `ready`；C/background 或单一信源必须 `degraded`。
- `valuation_percentiles` 缺失时不得声明安全边际、低估或高估结论。
- 内置知识条目中的 `formal_evidence_allowed` 必须为 false，除非未来独立 change 另行定义正式外部证据来源。

## 16. 本地访问与安全边界

- 默认监听 `127.0.0.1:8080`。
- 不提供公网部署配置。
- 不保存券商账号、交易密码、交易 token。
- 不提供买入、卖出、撤单、改单接口。
- 所有确认类接口必须明确显示“用户已在线下手动执行”或“仅记录计划”。

## 17. 与其他文档关系

- Eino 工作流见 `docs/workflow.md`。
- 前端字段映射见 `docs/frontend-contract.md`。
- UI 展示规范见 `docs/ui-design.md`。
- 数据和规则需求见 `docs/requirements.md`。
- 分层和目录规范见 `docs/architecture.md`。
