# P119 Full UI Control And Affordance Acceptance

Date: 2026-06-25

Change: `p119-full-ui-control-and-affordance-acceptance`

Status: `passed_pending_archive`

## Summary

P119 adds a repeatable full-route UI control and affordance acceptance layer after P114-P118. It covers all current production routes from `web/src/App.tsx`, inventories visible controls, classifies their interaction meaning, performs desktop/mobile layout checks, exercises upstream/light-toggle interactions with before/after assertions, exercises key UI write actions, and verifies backend/SQLite readback for local facts.

## Runner

Command:

```bash
bash scripts/p119-full-ui-control-and-affordance-acceptance.sh
```

Final evidence:

- `docs/release/ui-audit-assets/2026-06-25-p119-full-ui-control-and-affordance-acceptance/p119-ui-control-summary.json`
- `docs/release/ui-audit-assets/2026-06-25-p119-full-ui-control-and-affordance-acceptance/browser/p119-browser-results.json`

## Result

| Gate | Result |
| --- | --- |
| Production routes visited | 22 |
| Desktop visible controls inventoried | 603 |
| Mobile route layout checks | 8 |
| Unnamed controls | 0 |
| Unclassified controls | 0 |
| Layout issues | 0 |
| Product-copy issues | 0 |
| Upstream/light-toggle interactions exercised | 24 |
| Toggle issues | 0 |
| Browser console errors | 0 |
| Page errors | 0 |
| API 5xx responses | 0 |

Control categories:

| Category | Count |
| --- | ---: |
| navigation | 493 |
| read_action | 30 |
| form_input | 28 |
| write_local_fact | 25 |
| light_interaction | 23 |
| governance_confirm | 3 |
| disabled_expected | 1 |

SQLite readback:

| Evidence | Count |
| --- | ---: |
| positions | 3 |
| portfolio_snapshots | 3 |
| position_transactions | 2 |
| operation_confirmations | 3 |
| error_cases | 1 |
| risk_alert_resolved | 1 |
| unread_p119_notifications | 0 |
| data_quality_resolutions | 1 |
| rule_proposals | 1 |
| intelligence_items | 4 |
| rag_chunks | 3 |
| audit_events | 19 |
| forbidden_broker_order_push_tables | 0 |
| auto_confirmation_rows | 0 |
| auto_rule_apply_events | 0 |

Upstream/light-toggle sweep:

| Area | Interaction evidence |
| --- | --- |
| Global shell | `刷新摘要` reload preserves current route and title; mobile `导航` opens, navigates, and closes. |
| Data quality | Symbol switch updates URL/state; knowledge/audit details summaries open. |
| Decision detail | Evidence chain and analysis-material panels expand and collapse with visible state changes. |
| Evidence and audit | Evidence role filter, evidence row summary, audit status filter, and audit reference expander all update state. |
| Rules/local ops | Rules details, local-knowledge structured-record details, local-install config/command details, settings preflight details, and every visible desktop details summary instance toggle open and closed. |
| Forms | Positions discipline/transaction selects and consultation scenario select update without submitting. |

## Fixes Made During P119

- `web/src/app/AppLayout.tsx`: global `刷新摘要` button now performs a real current-page refresh instead of being an inert fake affordance.
- `web/src/pages/DailyDisciplineReportDetailPage.tsx`: detail page h1 now uses the shared `page-title` class so full-route visual checks cover it consistently.
- `web/src/pages/RulesPage.tsx`: rule proposal confirm/final-confirm actions now preserve the success message after refreshing the proposal list.
- `web/e2e/p119-full-ui-control-and-affordance-acceptance.spec.ts`: P119 now contains a dedicated upstream/light-toggle sweep and an all-visible-details-summary instance sweep, so toggle controls are clicked and asserted, not only inventoried.
- `scripts/p119_full_ui_control_and_affordance_acceptance.py`: final merge now fails if toggle coverage is missing or if any toggle issue is recorded.

## Scope Boundary

P119 validates local UI controls, layout, and backend consistency in an isolated local run. It does not claim install/upgrade/release validation, broker execution, external push, automatic trading, automatic confirmation, automatic rule application, fresh provider coverage, physical second-machine validation, prediction accuracy, or return guarantees.

## Verification

- `bash scripts/p119-full-ui-control-and-affordance-acceptance.sh` passed with 24 upstream/light-toggle interactions and 0 toggle issues.
- `openspec validate p119-full-ui-control-and-affordance-acceptance --strict` passed.
- `openspec validate --all --strict` passed: 40 items.
- `go test ./...` passed with existing sqlite-vec macOS deprecation warnings.
- `go vet ./...` passed with existing sqlite-vec macOS deprecation warnings.
- `npm --prefix web test -- --run` passed: 53 files, 191 tests.
- `npm --prefix web run build` passed.
- `python3 scripts/p92_final_requirement_audit.py --check` passed.
- `python3 scripts/p93_code_reality_audit.py --check` returned expected stale boundary: `docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`.
- `git diff --check` passed.
