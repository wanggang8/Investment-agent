# P115 Real User Scenario Acceptance Matrix

> Date: 2026-06-25  
> Change: `p115-real-user-scenario-acceptance`  
> Status: draft matrix for execution  
> Boundary: This matrix defines acceptance coverage only. It does not add broker integration, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, automatic repair, return guarantees, paid/login/authorized sources, Level2 data, high-frequency data, Docker/package refresh, or physical second-machine validation.

## Evidence Rules

Each scenario must produce or reference:

- UI/browser evidence or HTTP API evidence.
- SQLite field-level readback.
- Downstream readback from another page or API when the operation should affect product state.
- Audit, notification, confirmation, risk, or rule side effect when applicable.
- Safety negative evidence for forbidden automation.
- Expected eligibility: `fresh_pass`, `scoped_pass`, or `degraded_expected`.
- Actual status starts as `pending` and can only become `fresh_pass`, `scoped_pass`, `degraded_expected`, or `blocked` after execution evidence exists.
- Evidence derived from P104 or deterministic local seed must be labeled `local_seeded_linkage`; it cannot support external provider or real LLM fresh claims.

## Scenario Coverage

| ID | Scenario | Main Surfaces | Required Operations | Required Evidence | Expected Eligibility | Actual Status |
| --- | --- | --- | --- | --- | --- | --- |
| S01 | First launch and local capability boundary | `/local-install`, `/settings`, `/api-diagnostics`, `/` | Health check, system settings read, capability settings read, diagnostics render | API health, settings/capability JSON, browser DOM, redacted UI, no trading claims | `fresh_pass` | `pending` |
| S02 | Empty account onboarding | `/positions`, `/workbench`, `/` | Open empty DB, inspect empty-state guidance | Browser screenshot/DOM, no fake position rows, dashboard/workbench empty-state | `fresh_pass` | `pending` |
| S03 | Account initialization and allocation calibration | `/positions` | Initialize total assets, cash, holdings, asset tags | Browser form evidence, portfolio API response, `portfolio_snapshots`, `positions`, dashboard/review readback | `fresh_pass` | `pending` |
| S04 | Holding add/edit/remove | `/positions` | Add holding, edit holding, remove holding | Browser holding controls, holding APIs, `positions`, `position_transactions`, audit readback | `fresh_pass` | `pending` |
| S05 | Batch portfolio import | `/positions` | Validate import, confirm import, reject invalid rows | Browser import controls, import APIs, committed batch, positions delta, invalid-row error | `fresh_pass` | `pending` |
| S06 | Offline transaction recording | `/positions`, decision confirmation flow | Record offline buy/sell/rebalance result | Browser or API operation, cash/position delta, transaction row, no broker/order rows | `fresh_pass` | `pending` |
| S07 | Local fact correction and audit | `/positions`, `/audit` | Correct a portfolio fact with reason | Correction API, correction row, audit event, review/dashboard readback | `fresh_pass` | `pending` |
| S08 | Quarterly rebalance review | `/positions`, `/review` | Run rebalance review and inspect suggested manual action | Rebalance API, review result, audit event, no auto adjustment | `fresh_pass` | `pending` |
| S09 | Active consultation full path | `/consultation`, `/decisions` | Submit symbol, assumptions, target return, previous baseline | Browser consultation, decision API, persisted decision, expected-return data, evidence refs, decision list | `scoped_pass` | `pending` |
| S10 | Decision detail evidence and rule explanation | `/decisions/:id` | Open decision detail, inspect verdict, evidence, rule chain, optional actions | Decision detail API, DOM readback, rule/evidence linkage | `fresh_pass` | `pending` |
| S11 | Manual confirmation and execution record | `/decisions/:id`, `/decision-loop`, `/audit` | Confirm accept/reject/watch and record offline outcome | Browser confirmation, confirmation API, `operation_confirmations`, decision loop, audit, no auto confirmation | `fresh_pass` | `pending` |
| S11B | Decision error marking and review loop | `/decisions/:id`, `/review`, `/audit` | Mark decision as `marked_error`, capture root cause and lesson learned | Confirmation API, `operation_confirmations.confirmation_type=marked_error`, review/audit readback, no auto rule mutation | `fresh_pass` | `pending` |
| S12 | Decision loop traceability | `/decision-loop` | List loops, filter/open one decision, inspect input-to-confirmation chain | Browser loop view, decision-loop APIs, linked decision/confirmation/audit rows | `fresh_pass` | `pending` |
| S13 | Evidence refresh/list/verification | `/evidence` | Refresh evidence, list evidence, inspect verification | Browser evidence controls, evidence APIs, evidence rows, verification status, no fake evidence on failure | `scoped_pass` | `pending` |
| S14 | RAG/VecLite rebuild and knowledge readiness | `/evidence`, `/local-knowledge` | Rebuild index, read readiness | Browser rebuild/readiness state, rebuild API, readiness API, chunk/index status, degradation message when unavailable | `scoped_pass` | `pending` |
| S15 | Local knowledge import governance | `/local-knowledge` | Validate import, redacted preview, confirm import | Browser import/redaction evidence, local knowledge APIs, fact/chunk rows, redaction evidence | `fresh_pass` | `pending` |
| S16 | Market refresh and source health | `/data-quality`, `/settings` | Refresh market data, read source health and latest snapshot | Browser quality/source health view, market APIs, snapshot/source health rows, provider/degraded classification | `scoped_pass` | `pending` |
| S17 | Data-quality regression and gate resolution | `/data-quality` | Read regression, create resolution, retire resolution | Browser resolution controls, DQ APIs, resolution rows, retired state, dashboard/workbench readback | `fresh_pass` | `pending` |
| S18 | Risk alert SOP lifecycle | `/risk-alerts`, `/risk-alerts/:id` | Open alert, update lifecycle/SOP state | Browser risk detail/action, risk APIs, risk row update, audit event, dashboard/review readback | `fresh_pass` | `pending` |
| S19 | Rule current version and proposal confirmation | `/rules` | Read current rule, create SOP addendum, confirm, final confirm | Browser proposal controls, rule APIs, proposal/rule rows, audit events, no auto rule apply | `fresh_pass` | `pending` |
| S20 | Rule effect validation and tracking | `/rules`, `/review` | Refresh effect validation, inspect tracking | Effect APIs, validation/tracking rows, sample-insufficient degradation | `scoped_pass` | `pending` |
| S21 | Notification center | `/notifications`, topbar | List notifications, mark one read, mark all read | Browser notification controls, notification APIs, read-state rows, unread count/readback | `fresh_pass` | `pending` |
| S22 | Daily discipline report | `/daily-discipline/reports`, detail route | Read today, list reports, open detail | Browser reports view, report APIs, report rows, insufficient-data frozen/needs-action state | `scoped_pass` | `pending` |
| S23 | Daily auto-run readonly status | `/daily-auto-run` | Read schedule/status/last run | Browser status view, auto-run status API, no automatic trade/confirmation side effects | `fresh_pass` | `pending` |
| S24 | Dashboard and workbench aggregate status | `/`, `/workbench` | Reopen after operations, inspect next actions and summaries | Dashboard API, browser DOM, aggregate values match prior facts | `fresh_pass` | `pending` |
| S25 | Review summary | `/review` | Inspect monthly/quarterly review summary | Browser review view, review API, cross-read decisions/confirmations/risk/portfolio facts | `fresh_pass` | `pending` |
| S26 | Audit event trace | `/audit` | Inspect audit trail after all writes | Browser audit table, audit API, event count/type/subject rows for key operations | `fresh_pass` | `pending` |
| S27 | Settings update | `/settings` | Read/update system and capability settings | Browser settings controls, settings APIs, persisted values, redacted sensitive fields | `fresh_pass` | `pending` |
| S28 | API diagnostics | `/api-diagnostics` | View static diagnostics page; runner separately calls backend health API | Static diagnostic DOM, independent health API evidence, no secret/path/raw payload leakage | `fresh_pass` | `pending` |
| S29 | Mobile core operation path | 390px viewport | Initialize/edit portfolio, consult, confirm, resolve risk/DQ, mark notification read | Mobile screenshots, DOM no-overflow/touch target, API/SQLite parity | `fresh_pass` | `pending` |
| S30 | Failure, degradation, and safety反证 | All core routes | Invalid input, no evidence, missing key, source failure, unavailable index, missing id | Error APIs/DOM, no partial fake writes, forbidden table/action counts are zero | `fresh_pass` | `pending` |
| S31 | Settings forbidden rule/SOP mutation | `/settings`, `/rules` | Submit `rule_thresholds` or `sop_config` through settings API | 400 rejection, no rule/SOP SQLite mutation, no auto rule apply audit | `fresh_pass` | `pending` |
| S32 | Local install diagnostic summary and redaction | `/local-install` | Inspect diagnostic summary and redaction behavior | Browser DOM and optional diagnostic artifact, no key/path/raw stack/prompt leakage | `fresh_pass` | `pending` |
| S33 | Browser-level interaction parity | Interactive product routes | Execute or assert one real browser interaction per interactive route | Desktop/mobile screenshots, console count, API/SQLite parity for `/positions`, `/evidence`, `/local-knowledge`, `/data-quality`, `/risk-alerts`, `/rules`, `/notifications`, `/settings` | `fresh_pass` | `pending` |

## Required Safety Negative Evidence

P115 must record all of the following:

- `forbidden_broker_order_push_tables = 0`.
- `auto_confirmation_rows = 0`.
- `auto_rule_apply_audit_events = 0`.
- `automatic_trading_affordances = 0`.
- `return_guarantee_claims = 0`.
- `secret_or_raw_prompt_leaks_on_primary_ui = 0`.

## Execution Notes

- P104 runner coverage may be reused for S03-S08, S11-S12, S17-S18, S21, S24-S26, and S30 only as `local_seeded_linkage` evidence when the API/SQLite effect is equivalent.
- P115 must add explicit browser coverage for S01-S05, S09-S19, S21-S29, and S32-S33.
- P115 must add explicit coverage for S11B and S31 before execution is considered complete.
- External provider and LLM-dependent scenarios may be `scoped_pass` or `degraded_expected` when the environment lacks a valid key or provider response; the acceptance record must not describe those as future provider guarantees.
