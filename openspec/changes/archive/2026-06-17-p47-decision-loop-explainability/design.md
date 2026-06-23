# P47 Design

## Overview

P47 提供“决策闭环解释”只读聚合层。它不替代决策详情、复盘摘要或审计时间线，而是在一个页面中串联现有事实，让用户能看见：系统当时给了什么裁决、用户后来记录了什么、是否有线下记录、哪些风险/复盘/审计事实与这条决策有关，以及当前闭环缺口在哪里。

## API Shape

### `GET /api/v1/decision-loops`

查询参数：

- `symbol`：可选，按标的过滤。
- `limit`：可选，默认 10，最大 50。

响应 `data`：

- `items`: `DecisionLoopItem[]`
- `total`: 返回条数
- `safety_note`: 固定只读安全边界

### `GET /api/v1/decision-loops/{decision_id}`

响应 `data` 为单条 `DecisionLoopItem`。

## DTO

`DecisionLoopItem`：

- `decision_id`
- `symbol`
- `generated_at`
- `final_verdict_status`
- `final_verdict_text`
- `confirmation_status`
- `loop_status`: `open / planned / recorded / reviewed / incomplete`
- `stages`: `DecisionLoopStage[]`
- `manual_actions`: `DecisionLoopManualAction[]`
- `risk_links`: `DecisionLoopLink[]`
- `review_links`: `DecisionLoopLink[]`
- `audit_links`: `DecisionLoopLink[]`
- `missing_links`: `string[]`
- `safety_note`

`DecisionLoopStage`：

- `stage`: `recommendation / confirmation / manual_record / risk_review / review`
- `status`: `complete / pending / not_required / missing / degraded`
- `label`
- `summary`
- `ref_type`
- `ref_id`
- `at`

`DecisionLoopManualAction`：

- `confirmation_id`
- `confirmation_type`
- `operation_type`
- `symbol`
- `quantity`
- `price`
- `fees`
- `executed_at`
- `transaction_ids`
- `note_preview`

`DecisionLoopLink`：

- `type`
- `id`
- `label`
- `href`
- `status`

## Backend Aggregation

Add read-only DTO/service/handler:

- `internal/application/dto/decision_loop.go`
- `internal/application/service/decision_loop.go`
- `internal/application/handler/decision_loop_handler.go`

Repository additions are read-only only:

- `DecisionRepository.ListOperationConfirmations(ctx, decisionID string)`
- `DecisionRepository.ListPositionTransactionsByConfirmation(ctx, confirmationID string)`

Aggregation rules:

1. Load decisions from `ListDecisionRecords`, filter by `symbol`, cap by `limit`.
2. For each item, load full `DecisionRecord` via `GetDecisionRecord` to access final verdict text, record type and JSON fields.
3. Load confirmations for that decision and transactions for each confirmation.
4. Load `ListErrorCases`, `ListRiskAlerts`, and `ListAuditEvents` once; filter in memory by decision, confirmation, error case or symbol.
5. Build stages:
   - recommendation is complete when decision exists.
   - confirmation is not_required when decision status is `not_required`, pending when `pending`, complete when a matching confirmation exists.
   - manual_record is complete when a confirmation has transaction ids, not_required for planned/watch/marked_error, missing if status is `executed_manually` but no transaction exists.
   - risk_review is complete when risk links exist, otherwise pending for degraded/high-risk statuses.
   - review is complete when error case or review/audit links exist; otherwise pending.
6. Build `missing_links` from incomplete expected stages.

## Frontend

Add:

- `web/src/types/decisionLoop.ts`
- `web/src/services/decisionLoop.ts`
- `web/src/pages/DecisionLoopPage.tsx`
- `web/src/pages/DecisionLoopPage.test.tsx`

Page sections:

- 状态概览：闭环条数、未闭合条数、最近决策。
- 闭环列表：每条决策的阶段、缺口、人工记录、风险/复盘/审计链接。
- 空态/错误态：沿用 `StatusNotice` 和只读文案。

Navigation:

- Add `/decision-loop` route.
- Add nav item “决策闭环”.
- Add Workbench and review summary links to `/decision-loop`.

## Guardrails

- P47 API is read-only; no transaction, confirmation, rule, notification or risk state writes.
- Frontend page has no action buttons that mutate facts.
- Do not display raw payload JSON, SQL, private paths, complete keys, prompts, broker state or order ids.
- Safety text should explain that the page is an explanation surface only, using wording that avoids smoke forbidden body-scan terms.
