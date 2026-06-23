# P84 Portfolio Confirmation Data Impact Closure

## Why

Thirty-five rows still need stronger evidence that portfolio/account setup, allocation bands, position maintenance, offline transaction records, manual confirmations, review linkage, and data impact readback work as real product workflows rather than isolated demos.

## What Changes

- Define the P84 row set as the portfolio/confirmation data-impact batch from the latest matrix, including P82-deferred allocation and quarterly rebalance rows.
- Execute real browser workflows that create or modify account/holding/transaction/confirmation state and verify affected downstream data.
- Verify UI, API, SQLite, audit, risk, review, and decision-loop readbacks.
- Produce P84 evidence and update remaining-row governance.

## In Scope

- 35 rows: REQ-01-001, REQ-01-006, REQ-02-006, REQ-02-022, REQ-02-024, REQ-02-025, REQ-02-031, REQ-02-033, REQ-03-004, REQ-03-005, REQ-03-006, REQ-04-003, REQ-04-008, REQ-04-016, REQ-04-025, REQ-05-010, REQ-06-023, REQ-06-024, REQ-07-006, REQ-07-015, REQ-08-018, REQ-08-020, REQ-10-001, REQ-10-002, REQ-10-003, REQ-10-004, REQ-11-002, REQ-11-005, REQ-11-019, REQ-14-005, REQ-14-007, REQ-16-028, REQ-16-033, REQ-17-015, REQ-17-024.

## Out of Scope

- No broker interface, order placement, one-click trading, automatic confirmation, or real brokerage account synchronization.
- Manual local records remain the only accepted path for portfolio mutation.
- No full original-requirement pass claim from this batch alone.
