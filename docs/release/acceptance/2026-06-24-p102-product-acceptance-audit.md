# P102 Product Acceptance Audit

Date: 2026-06-24

Scope: local-source product acceptance using `configs/config.yaml`, local Go backend, Vite frontend, SQLite, VecLite, and the configured OpenAI-compatible real model. Docker/install/release packaging/physical second-machine verification were intentionally out of scope.

## Verdict

`product_acceptance_pass_with_non_blocking_ux_findings`

From a product perspective, the current local application is usable for the intended investment-discipline workflow:

1. A user can open the main workbench and understand the next manual actions.
2. A user can create or recover local portfolio facts.
3. The system rejects inconsistent account data instead of silently accepting bad facts.
4. A real-model consultation can run through the OpenAI-compatible path and produce a decision detail.
5. The user can record a manual plan, and the decision loop, audit, review, workbench, and notifications surfaces reflect that state.
6. Safety boundaries are visible and consistent: no broker connection, no automatic trading, no delegated orders, no automatic confirmation, and no automatic rule application are claimed.

The audit found no product release-blocking defect in the local-source configured-real-model scope. It did find several UX risks that should be handled in a follow-up polish change before treating the experience as fully smooth for a new non-technical user.

## Evidence Bundle

Asset folder:

`docs/release/ui-audit-assets/2026-06-24-product-acceptance-audit/`

Key evidence:

- Desktop route screenshots: `01-dashboard.png` through `17-decision-loop.png`
- Mobile screenshots: `18-mobile-dashboard.png` through `21-mobile-data-quality.png`
- Capture summaries: `capture-summary.json`, `mobile-capture-summary.json`
- Real model consult summary: `consult-summary.json`
- API readbacks:
  - `api-portfolio-adjust-response.json`
  - `api-portfolio-current-after-adjust.json`
  - `api-confirm-planned-response.json`
  - `api-decision-loop-after-confirmation.json`
- DOM readbacks after real workflow:
  - `22-positions-after-calibration.dom.txt`
  - `23-decision-detail-real-model.dom.txt`
  - `24-decision-loop-after-confirmation.dom.txt`
  - `26-audit-after-confirmation.dom.txt`
  - `27-review-after-confirmation.dom.txt`
  - `28-workbench-after-confirmation.dom.txt`
  - `29-notifications-after-confirmation.dom.txt`
- SQLite readback: `sqlite-readback-summary.json`

## Product Journey

### 1. Entry And Navigation

The dashboard and workbench present the product as a local discipline cockpit rather than a trading app. The first-screen hierarchy is clear: daily status, prohibited actions, optional manual actions, account/risk/data-quality entry points, and governance links.

All captured desktop routes loaded without browser console errors during the screenshot pass. The tested surfaces included dashboard, workbench, consultation, positions, risk alerts, data quality, evidence, local knowledge, rules, audit, review, notifications, daily reports, daily auto-run, local install, settings, and decision loop.

### 2. Portfolio Facts

Initial empty local portfolio state produced `GET /api/v1/portfolio/current` 404 with `record not found`. On mobile `/positions`, this surfaced as a generic "读取失败 / 系统暂时无法处理请求，请稍后重试" message above the first-use guidance.

The product recovered once local account facts were entered through the product API:

- Cash: `80000`
- Total assets: `126500`
- Positions: `510300` normal core holding, `159915` frozen-watch satellite holding

The first attempted write with inconsistent total assets was rejected with a clear error: `total_assets 与现金和持仓市值不一致`. The corrected write succeeded and was read back from API and SQLite.

Product assessment: the validation logic is good and protects local facts. The empty-state copy should be improved because a first-time user may interpret the generic 404 state as system failure.

### 3. Real Model Consultation

Real consultation request:

- Endpoint: `POST /api/v1/decisions/consult`
- Request ID: `req_p102_product_consult_real_model`
- Symbol: `510300`
- Scenario: `hold_review`
- Model observed in analyst reports: `gpt-5.4-mini`

Result:

- Decision ID: `decision_6f9fa7db5afe919a`
- Workflow status: `degraded`
- Capability: `in_scope`
- Final verdict: `hold`
- Display text: `按纪律观察`
- Confirmation status before user action: `pending`
- Available manual actions: `planned`, `watch`, `executed_manually`, `marked_error`

The three analyst reports were parsed and quality-passed. The product also correctly exposed the degraded data context instead of hiding it.

Product assessment: the real-model path is usable and the final decision is appropriately bounded. The decision detail is information-rich, but the raw analyst text is dense; it would benefit from progressive disclosure or stronger section collapsing for normal users.

### 4. Manual Confirmation And Decision Loop

Manual planned confirmation succeeded:

- Confirmation ID: `confirm_63322851fcc52a6a`
- Confirmation type: `planned`
- Note: user chose continued observation and no automatic trade
- Audit event: `audit_373ccc83e4e53534`

Decision loop readback showed:

- Recommendation: complete
- User record: complete
- Manual transaction record: not required
- Risk link: pending, with explicit missing risk clue
- Review link: complete through audit clue

Product assessment: the closed-loop model is coherent. It records manual intent without inventing a transaction. The missing risk clue is explicit, which is preferable to pretending the loop is fully complete.

### 5. Governance, Audit, Review, Notifications

After confirmation:

- Audit showed the confirmation transition `pending -> planned`.
- Review summary counted the confirmation, current decision, degraded workflows, missing evidence themes, and rule proposal status.
- Workbench reflected the updated portfolio assets, active risk, rules/review status, and next actions.
- Notifications surfaced local warning items without external push claims.

Product assessment: governance pages are useful and mostly well connected. They support traceability and manual follow-up rather than presenting the system as an autonomous trading engine.

## Findings

### P1 UX: Empty Portfolio Is Presented As Generic Failure

Evidence: `20-mobile-positions.png`; `GET /api/v1/portfolio/current` returned 404 before account calibration.

Impact: a first-time user can see "读取失败" before the first-use guidance, which makes normal empty setup look like a system error.

Recommendation: map missing portfolio snapshot to a first-use empty state. Keep the API strict if desired, but the UI should say "尚未录入本地账户" and place calibration guidance first.

### P2 UX: Decision Detail Is Too Dense For Routine Review

Evidence: `23-decision-detail-real-model.dom.txt`.

Impact: real model output is long and useful, but the page reads like a large report dump. This makes the final verdict and action boundary harder to scan.

Recommendation: keep the current detail for auditability, but add collapsed analyst sections, a compact "why this verdict" summary, and a persistent manual-action panel near the top.

### P2 UX: Decision Loop Query Deep Link Does Not Act Like A Focused Detail View

Evidence: `24-decision-loop-after-confirmation.dom.txt`; navigating with `decision_id` still rendered the loop list with the target latest decision at top.

Impact: acceptable when the target is recent, but less reliable for older decisions.

Recommendation: make `/decision-loop?decision_id=<id>` focus or filter to that decision, or change links to a supported route.

### P2 Product State: Data Quality Degradation Is Correctly Visible But Blocks "Clean Data" Claims

Evidence: `06-data-quality.png`; `28-workbench-after-confirmation.dom.txt`; consultation workflow status `degraded`.

Impact: this does not block local product use, but it prevents any honest claim that the current data state is fully clean.

Recommendation: keep release wording scoped: real model works and product flow works, while current data quality still includes degraded items.

## Strengths

- The product consistently frames itself as a local discipline and review tool, not a broker or auto-trading system.
- The navigation model is understandable across workbench, portfolio, decision, evidence, governance, and operations.
- Data validation is real and user-facing; bad portfolio totals are rejected.
- Real model output is connected to local facts, market snapshot, rule version, evidence, expected-return scenarios, audit, and confirmation.
- Safety boundaries are repeated in the right places without feeling hidden in legal text.
- Mobile routes reflow into a usable single-column layout, though empty-state wording needs work.

## Scope Boundaries

This audit does not claim:

- Docker install/upgrade/uninstall validation.
- GitHub Release validation.
- Distribution package refresh.
- Physical second-machine verification.
- Broker connectivity.
- Automatic trading, one-click trading, delegated orders, external push, automatic confirmation, or automatic rule application.
- Paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.

## Conclusion

The product is usable in the local-source configured-real-model scope. The core workflow from local facts to real-model decision, manual confirmation, audit, review, and notification is coherent and readback-backed. The remaining issues are product polish and clarity issues rather than functional blockers.
