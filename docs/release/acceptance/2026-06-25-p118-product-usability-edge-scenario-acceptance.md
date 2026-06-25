# P118 Product Usability Edge Scenario Acceptance

> Date: 2026-06-25  
> Change: `p118-product-usability-edge-scenario-acceptance`  
> Status: passed; pending archive only  
> Boundary: local product-use edge scenario acceptance only. This explicitly excludes Docker/install/upgrade/uninstall/release package/GitHub Release/physical second-machine validation. It does not claim broker execution, external push, fresh external provider cleanliness, fresh real LLM quality, automatic trading, automatic confirmation, automatic rule application, prediction accuracy or return guarantees.

## Scope

The user asked to skip release/install/upgrade scenarios and continue with the remaining product-use scenarios. P118 therefore covers:

- Long-cycle durability after accumulated local history.
- Abnormal import and transaction recovery.
- Data-quality degradation with explicit scoped resolution.
- Decision-quality interpretation across rising, falling and volatile contexts.
- Multi-account/household ledger facts as local notes, not broker accounts.
- Cross-page and mobile readback after accumulated state.
- Safety negative evidence.

## Runner Result

Command:

```bash
bash scripts/p118-product-usability-edge-scenario-acceptance.sh
```

Result:

- Status: `passed`.
- Runner scenarios: `18/18`.
- Fresh pass: `16`.
- Scoped pass: `2`.
- Browser console errors: `0`.
- Browser page errors: `0`.
- Browser API 5xx responses: `0`.
- Restart probe: `passed`.

Evidence:

- Final summary: `docs/release/ui-audit-assets/2026-06-25-p118-product-usability-edge-scenario-acceptance/p118-edge-usability-summary.json`.
- API/SQLite summary: `docs/release/ui-audit-assets/2026-06-25-p118-product-usability-edge-scenario-acceptance/api_sqlite/p118-api-sqlite-summary.json`.
- Restart summary: `docs/release/ui-audit-assets/2026-06-25-p118-product-usability-edge-scenario-acceptance/restart/p118-restart-summary.json`.
- Browser summary: `docs/release/ui-audit-assets/2026-06-25-p118-product-usability-edge-scenario-acceptance/browser/p118-browser-results.json`.
- Screenshots: `docs/release/ui-audit-assets/2026-06-25-p118-product-usability-edge-scenario-acceptance/browser/*.png`.

## Product Usability Interpretation

P118 adds evidence beyond P117's seven-day story:

- Long-cycle use is usable: accumulated daily reports, transactions, notifications, risk alerts and audit rows remain readable.
- Recovery is usable: invalid import and invalid transaction attempts are rejected without partial transaction writes.
- Data-quality degradation is usable within a scoped boundary: stale/missing source facts require explicit local resolution and do not become clean external-data claims.
- Decision interpretation is usable as local evidence: rising, falling and volatile contexts produce different seeded verdicts with explicit prohibited actions and disclaimers.
- Household ledger is usable as local facts: multiple household account notes and fund categories remain traceable without broker account integration.
- Persistence is usable: backend restart reads accumulated local history from the same SQLite database.
- Mobile is usable within checked scope: 390px accumulated-state pages render without console/page/API 5xx failures.
- Safety boundary holds: no broker/order/push tables, automatic confirmations or auto rule-apply audit events were found.

## SQLite Readback Highlights

- `daily_discipline_reports`: `31` rows, including `3` degraded/insufficient rows.
- `position_transactions`: `13` rows across `7` distinct symbols.
- `positions`: `7` current rows, with local household/family notes.
- `operation_confirmations`: `14` rows.
- `error_cases`: `1` marked-error case.
- `risk_alerts`: `8` rows.
- `notifications`: `26` rows.
- `audit_events`: `67` rows.
- `data_quality_gate_resolutions`: retired scoped resolution present.

## Safety Counters

- `forbidden_broker_order_push_tables = 0`.
- `auto_confirmation_rows = 0`.
- `auto_rule_apply_audit_events = 0`.
- `automatic_trading_affordances = 0`.
- `return_guarantee_claims = 0`.
- `secret_or_raw_prompt_leaks_on_primary_ui = 0`.

## Notes

During development, two Playwright assertions failed before the final pass:

- `银华日利` matched two visible elements. The assertion was narrowed to the exact table cell.
- Decision detail pages use `决策详情` as the H1 for existing decisions. The assertion was corrected from `主动咨询` to `决策详情`.

Both were test assertion issues, not product runtime failures. The final P118 runner passed after those corrections.

## Regression Gates

- `bash scripts/p118-product-usability-edge-scenario-acceptance.sh`: passed.
- `openspec validate p118-product-usability-edge-scenario-acceptance --strict`: passed.
- `openspec validate --all --strict`: passed, 39 items.
- `go test ./...`: passed, with sqlite-vec macOS deprecation warnings.
- `go vet ./...`: passed, with sqlite-vec macOS deprecation warnings.
- `npm --prefix web test -- --run`: passed, 53 files / 191 tests.
- `npm --prefix web run build`: passed.
- `python3 scripts/p92_final_requirement_audit.py --check`: passed.
- `python3 scripts/p93_code_reality_audit.py --check`: stale, `docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`.
- `git diff --check`: passed.

P93 stale is expected after P114-P118 source and evidence changes. P118 does not treat that stale result as a fresh P93 pass.
