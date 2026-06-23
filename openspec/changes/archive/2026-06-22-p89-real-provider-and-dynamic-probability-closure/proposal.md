# P89 Real Provider And Dynamic Probability Closure

## Summary

P89 closes or explicitly preserves the 10 full-release-required rows that remained `partial` after P88:

- `REQ-04-016`
- `REQ-05-003`
- `REQ-05-004`
- `REQ-05-005`
- `REQ-08-004`
- `REQ-08-023`
- `REQ-09-004`
- `REQ-09-023`
- `REQ-09-024`
- `REQ-09-025`

P89 focuses on two evidence gaps:

1. Real no-login/no-paid/no-Level2/no-high-frequency public provider collection and SQLite readback for capital flow, margin financing, and constituent financial fields.
2. Real UI/API/SQLite dynamic expected-return monitoring for valuation/fundamental/market-state changes, extreme-fear historical similar-scenario display, two-month assumption downshift, and one-month pessimistic-path manual probability adjustment.

## Why

P88 upgraded 17 of the final 27 P86 blockers but intentionally kept 10 rows `partial`. The remaining rows cannot be honestly upgraded by parser contracts, accepted-local evidence, fixture data, route smoke, screenshots, or Go-only tests. They need real provider proof or real user-facing dynamic-monitoring evidence.

## In Scope

- Build an inventory gate that proves P89 owns exactly the 10 P88 remaining full-release-required non-`real_pass` rows.
- Verify candidate public providers before production collector use, including authority, access shape, fields, update frequency, legal/access limits, rate assumptions, failure behavior, and target SQLite path.
- Implement only read-only, low-frequency public collection for the three structured categories if a candidate source is verified:
  - capital flow: `date`, `net_inflow`, `net_outflow`
  - margin financing: `date`, `margin_balance`, `balance_change_rate`
  - constituent financial: `revenue`, `net_profit`, `growth`, `disclosure_date`
- If a source is unavailable, blocked, unstable, login/paid/authorization-only, Level2, or high-frequency, preserve the row as `partial` with exact blocker evidence.
- Add real UI/API/SQLite acceptance for dynamic expected-return monitoring and extreme-fear lock behavior.
- Generate a P89 closure matrix and acceptance record.

## Out Of Scope

- Broker integration, order placement, one-click trading, delegated trading, automatic trading, automatic confirmation, external push, automatic rule application, automatic repair, automatic migration, automatic recovery, or overwriting a real user database.
- Login, paid, authorization-only, Level2, or high-frequency sources.
- Future provider availability promises.
- Future return accuracy promises or market direction promises.
- P76/package refresh, remote release, Git tag, physical second-machine validation, or full original-requirement pass claim unless all full-release-required rows are `real_pass`.

## Acceptance

P89 is acceptable only if:

- Inventory gate proves exactly 10 owned rows from the P88 matrix.
- Every upgraded row has direct fresh P89 evidence.
- Structured data rows are upgraded only with verified non-mock runtime provider collection plus SQLite readback.
- Dynamic probability rows are upgraded only with real UI/API/SQLite before/after proof.
- Final matrix lists all remaining blockers exactly if any row remains non-`real_pass`.
- Subagent plan and final reviews report no Critical or Important findings.
- `openspec validate --all --strict`, P89 runner/checkers, Go tests, frontend tests/build, and `git diff --check` pass.
