# P90 Capital Flow Provider Closure

## Summary

P90 closes the two full-release-required rows left `partial` after P89:

- `REQ-04-016`
- `REQ-05-003`

P90 replaces the blocked Eastmoney `push2` capital-flow path with the publicly reachable Eastmoney H5 capital-flow endpoint, then proves `date`, `net_inflow`, and `net_outflow` through real Settings UI market refresh, market snapshot API readback, and SQLite readback.

## Why

P89 upgraded margin-financing and constituent-financial rows but preserved the capital-flow rows because the `push2`/`push2his` endpoint returned curl exit 52 in this environment. The Eastmoney H5 page exposes a separate public endpoint (`/dc/ZJLX/getDBHistoryData`) used by the mobile capital-flow page and reachable without login, payment, authorization, Level2, or high-frequency access.

## In Scope

- Verify Eastmoney H5 capital-flow endpoint availability and fields.
- Implement runtime read-only capital-flow collection for `date`, `net_inflow`, and `net_outflow`.
- Keep directional net-flow semantics explicit: a positive daily net flow maps to `net_inflow`, a negative daily net flow maps to `net_outflow`; values are not synthesized from unrelated fields.
- Prove product path through Settings UI market refresh, market snapshot API readback, and SQLite readback.
- Generate P90 matrix and closure record that upgrades only the two directly evidenced rows.

## Out Of Scope

- Broker integration, order placement, one-click trading, delegated trading, automatic trading, automatic confirmation, external push, automatic rule application, automatic repair, automatic migration, automatic recovery, or overwriting a real user database.
- Login, paid, authorization-only, Level2, or high-frequency sources.
- Future provider availability promises.
- Return accuracy promises or market direction promises.
- P76/package refresh, remote release, Git tag, or physical second-machine validation.

## Acceptance

P90 is acceptable only if:

- P90 inventory proves it owns exactly `REQ-04-016` and `REQ-05-003`.
- Source preverification records the public H5 endpoint, request URL, page/JS evidence, fields, update frequency, access limits, rate assumptions, and failure behavior.
- Runtime collector fetches capital-flow data without stubs, fixtures, accepted-local data, manual seed, login, paid, authorization, Level2, or high-frequency source.
- Real browser acceptance triggers product market refresh and reads capital-flow fields through UI/API/SQLite.
- Final matrix has both P90 rows `real_pass`.
- Subagent final review reports no Critical or Important findings.
- `openspec validate --all --strict`, P90 runner/checkers, Go/frontend tests/build, and `git diff --check` pass.
