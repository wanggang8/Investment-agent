# P88 Remaining Full Release Blockers Closure

## Why

P86 reduced the original-requirement matrix to 303 `real_pass`, 11 `reference_only`, and 27 remaining full-release-required `partial` rows. Those 27 rows are no longer broad acceptance bookkeeping. They require concrete product evidence or implementation across source-verified state transitions, structured public data fields, expected-return historical/probability logic, quarterly rebalance, and SOP addendum proposal generation.

P88 is the next row-specific closure change. Its goal is to resolve the 27 P86 blockers with real UI/API/SQLite/workflow evidence, or to honestly preserve any row that still cannot be proven without forbidden capabilities or unavailable public data.

## What Changes

- Create a P88 inventory gate from the P86 matrix and fail if the remaining blocker set is not exactly the 27 known full-release-required rows.
- Add row-specific acceptance for source-verified buy-logic-break and multi-source-insufficient transitions into `sell_only` and `frozen_watch`.
- Add source preverification and structured data readback evidence for capital-flow, margin-financing, and constituent-financial fields.
- Add expected-return historical sample/probability, dynamic scenario update, assumption tracking, and low-sample supplement-data evidence.
- Add quarterly +/-15% rebalance UI/API/SQLite readback evidence.
- Add high-frequency uncovered-scenario SOP addendum proposal generation evidence.
- Generate a P88 matrix and acceptance record that either reaches full original-requirement pass or lists exact remaining rows and blockers.

## In Scope

The P88 owned rows are exactly:

- Source-verified state transitions: `REQ-02-022`, `REQ-02-025`, `REQ-06-023`, `REQ-06-024`, `REQ-17-015`.
- Structured data/source readiness: `REQ-04-016`, `REQ-04-025`, `REQ-05-003`, `REQ-05-004`, `REQ-05-005`.
- Extreme fear and scenario update: `REQ-08-004`, `REQ-08-023`.
- Expected return: `REQ-09-001`, `REQ-09-003`, `REQ-09-004`, `REQ-09-006`, `REQ-09-007`, `REQ-09-008`, `REQ-09-009`, `REQ-09-010`, `REQ-09-013`, `REQ-09-023`, `REQ-09-024`, `REQ-09-025`, `REQ-09-027`.
- Rebalance and SOP proposal: `REQ-10-004`, `REQ-13-010`.

## Out of Scope

- No broker connection, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real DB overwrite, paid/login/authorized source, Level2 source, high-frequency source, or return promise.
- No paid, login, authorized, scraped-behind-access-control, Level2, or high-frequency data source.
- No full original-requirement pass unless every full-release-required row is `real_pass`. Any row reclassification requires explicit L1/OpenSpec rationale and must be reported separately from a 27/27 real-pass outcome.
- No row upgrade based only on seeded decisions, route smoke, screenshots, fixture/mock/stub data, or broad narrative.
- No P76 package refresh unless a later explicit packaging change is created.
