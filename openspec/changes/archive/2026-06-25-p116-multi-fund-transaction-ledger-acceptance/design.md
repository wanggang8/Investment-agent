# P116 Design

## Goal

P116 verifies whether the product can handle a realistic local investment ledger: several funds/ETFs, several dates, buy/sell/reduce operations, import errors, corrections, decision confirmations, and downstream readback. It is an acceptance change, not a runtime feature change.

## Evidence Model

The runner writes three evidence layers under `docs/release/ui-audit-assets/2026-06-25-p116-multi-fund-transaction-ledger-acceptance/`:

- `api_sqlite`: API calls and SQLite field-level readback.
- `browser`: Playwright UI screenshots, DOM assertions, console/page/API health.
- `final`: merged `p116-scenario-summary.json`.

Every scenario records:

- `scenario_id`, `title`, `status`, `expected_eligibility`, `classification_reason`.
- Runtime boundary: `config_mode`, `runtime_mode`, `use_stub`, `provider_mode`, `llm_mode`.
- `symbols`, `api_evidence`, `sqlite_evidence`, `browser_evidence`, `downstream_evidence`, `side_effects`, `safety_counters`.

## Scenario Set

- `L01` Fresh local runtime and empty ledger.
- `L02` Multi-fund initial portfolio.
- `L03` Multi-date offline buy/sell/reduce ledger.
- `L04` Mixed batch import with valid and invalid rows.
- `L05` Invalid transaction rejection with no partial writes.
- `L06` Holding edit/remove and local correction.
- `L07` Decision-to-manual-execution confirmation.
- `L08` Marked-error review loop.
- `L09` Quarterly rebalance review across core/satellite/cash buckets.
- `L10` Risk alert lifecycle and notification readback.
- `L11` Data-quality gate resolution and retirement.
- `L12` Dashboard/workbench/review/audit aggregate readback.
- `L13` Browser multi-fund positions path.
- `L14` Browser decision confirmation path.
- `L15` Mobile portfolio rendering.
- `L16` Safety negative evidence.

## Data Shape

The seeded local ledger uses ETF/fund-like symbols:

- `510300` 沪深300ETF, core.
- `159915` 创业板ETF, satellite.
- `588000` 科创50ETF, satellite.
- `512000` 券商ETF, satellite.
- `110022` 易方达消费行业混合, fund-like local holding.

The runner exercises buy, sell, and reduce operations through local offline transaction APIs and manual decision confirmation. All execution is recorded as user-completed offline facts. P116 must not create broker/order/push tables, automatic confirmations, or automatic rule-apply audit events.

## Browser Strategy

Browser coverage is intentionally narrower than API/SQLite but proves the core product surface:

- `/positions`: initialize a multi-fund portfolio, import another fund row, record an offline transaction, run rebalance, verify visible multi-fund table.
- `/decisions/decision_p116_browser_execute`: record a manual execution confirmation.
- `/decision-loop`, `/`, `/workbench`, `/review`, `/audit`, `/risk-alerts`, `/notifications`, `/data-quality`: verify downstream surfaces render after the ledger writes.
- 390px `/positions`: verify mobile portfolio rendering.

## Boundary

P116 is local-source functional reality evidence. It does not claim external provider availability, fresh real LLM output, broker integration, auto trading, auto confirmation, auto rule application, release package freshness, or P93 freshness after P114/P115/P116 changes.

