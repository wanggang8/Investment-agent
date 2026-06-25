# P116 Multi-Fund Transaction Ledger Acceptance

> Date: 2026-06-25  
> Change: `p116-multi-fund-transaction-ledger-acceptance`  
> Result: `passed_scoped_local_multi_fund_ledger_acceptance`

## Summary

P116 executed an isolated local backend, Vite frontend, temporary SQLite database and Playwright browser journey for richer real-use scenario acceptance. It specifically expands P115 beyond one fixed fund by validating a multi-fund local transaction ledger with several fund symbols, multiple dated operations, invalid input rejection, decision confirmation, review/audit readback and safety boundaries.

Final merged evidence:

- `docs/release/ui-audit-assets/2026-06-25-p116-multi-fund-transaction-ledger-acceptance/p116-scenario-summary.json`
- `docs/release/ui-audit-assets/2026-06-25-p116-multi-fund-transaction-ledger-acceptance/api_sqlite/p116-api-sqlite-summary.json`
- `docs/release/ui-audit-assets/2026-06-25-p116-multi-fund-transaction-ledger-acceptance/browser/p116-browser-results.json`

## Coverage

- Scenario count: 16.
- Fresh pass: 14.
- Scoped pass: 2 (`L11` data-quality seeded gate lifecycle, `L12` aggregate readback through local governance linkage).
- Browser status: passed.
- Browser console errors: 0.
- Browser page errors: 0.
- Browser API 5xx responses: 0.
- Symbols covered: `510300`, `159915`, `588000`, `512000`, `110022`, `161725`.

## Important Real-Use Paths

- Initial portfolio recorded five funds through the API and read back from SQLite.
- Offline ledger recorded buy/sell/reduce operations across `159915`, `510300`, `512000`, `588000`, with fees and dates.
- Mixed import validation rejected invalid rows and did not write transactions.
- Valid import committed a new holding and a transaction row.
- Invalid operations were rejected: insufficient cash, oversell, future execution time, negative fees, missing symbol and invalid `position_state`.
- Holding edit, holding removal, correction audit and rebalance review all wrote local facts.
- Decision manual execution and `marked_error` produced confirmation/audit evidence.
- Risk alert, notification and data-quality gate lifecycle paths were exercised.
- Browser journey covered `/positions`, `/decisions/:id`, `/decision-loop`, `/risk-alerts`, `/notifications`, `/data-quality`, `/`, `/workbench`, `/review`, `/audit` and 390px `/positions`.

## Safety Boundary

P116 remains local ledger acceptance only. It does not claim broker integration, real exchange execution, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, fresh external provider quality, fresh real LLM output, release package refresh or physical second-machine validation.

Safety counters in the final summary:

- Forbidden broker/order/push tables: 0.
- Auto confirmation rows: 0.
- Auto rule-apply audit events: 0.
- Automatic trading affordances: 0.
- Return guarantee claims: 0.
- Secret/raw prompt leaks: 0.

## Notes

During the first browser run, the form test input used an inconsistent total asset value. The backend rejected it and the UI showed a read-failure state, proving the screen was connected to real validation. The Playwright input was corrected to match cash plus market value, and the full P116 runner then passed.

P93 remains a potentially stale historical audit after later P114/P115/P116 work; P116 must not be cited as fresh P93 code-reality pass.

## Verification

- `bash scripts/p116-multi-fund-transaction-ledger-acceptance.sh`: passed.
- `openspec validate p116-multi-fund-transaction-ledger-acceptance --strict`: passed.
- `openspec validate --all --strict`: passed, 37 items.
- `go test ./...`: passed, with sqlite-vec macOS deprecation warnings.
- `go vet ./...`: passed, with sqlite-vec macOS deprecation warnings.
- `npm --prefix web test -- --run`: passed, 53 files / 191 tests.
- `npm --prefix web run build`: passed.
- `python3 scripts/p92_final_requirement_audit.py --check`: passed.
- `python3 scripts/p93_code_reality_audit.py --check`: stale, `docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`.
- `git diff --check`: passed.
