## ADDED Requirements

### Requirement: P88 remaining full release blockers closure

After P86, P88 SHALL resolve or explicitly preserve the 27 remaining full-release-required rows by adding row-specific implementation and real UI/API/SQLite/workflow evidence for source-verified transitions, structured public-data fields, expected-return historical/probability behavior, quarterly rebalance, and SOP addendum proposals.

#### Scenario: P88 row inventory starts from the P86 remainder

- **GIVEN** P88 starts from the P86 matrix
- **WHEN** the inventory gate runs
- **THEN** it SHALL find exactly 27 full-release-required non-`real_pass` rows
- **AND** the row IDs SHALL be `REQ-02-022`, `REQ-02-025`, `REQ-04-016`, `REQ-04-025`, `REQ-05-003`, `REQ-05-004`, `REQ-05-005`, `REQ-06-023`, `REQ-06-024`, `REQ-08-004`, `REQ-08-023`, `REQ-09-001`, `REQ-09-003`, `REQ-09-004`, `REQ-09-006`, `REQ-09-007`, `REQ-09-008`, `REQ-09-009`, `REQ-09-010`, `REQ-09-013`, `REQ-09-023`, `REQ-09-024`, `REQ-09-025`, `REQ-09-027`, `REQ-10-004`, `REQ-13-010`, and `REQ-17-015`.

#### Scenario: Source-verified state transitions are proven by evidence counts

- **GIVEN** formal source-verification evidence exists for a held symbol
- **WHEN** at least two independent A/S formal sources confirm buy-logic break
- **THEN** P88 SHALL prove the workflow enters `sell_only`, prohibits buy/add actions, and records source-count provenance, API/readback, SQLite facts, and audit evidence
- **AND** when fewer than two independent A/S formal sources exist for buy-logic questioned, major positive, or major negative information, P88 SHALL prove the workflow enters `frozen_watch` with source-count provenance and pause guidance.

#### Scenario: Structured data fields require preverified public-source evidence

- **GIVEN** P88 expands structured data evidence for capital flow, margin financing, and constituent financials
- **WHEN** collectors are used for `real_pass` structured-data evidence
- **THEN** P88 SHALL record source preverification before claiming production readiness
- **AND** it SHALL prove field-level readback for capital-flow date/net-inflow/net-outflow, margin-financing date/balance/change-rate, and constituent-financial revenue/profit/growth/disclosure-date.
- **AND** accepted-local, fixture, stub, or manually seeded evidence SHALL NOT upgrade structured-data collector rows to `real_pass`.

#### Scenario: Expected-return report uses historical/probability evidence and safe degradation

- **GIVEN** expected-return analysis runs through real UI/API/SQLite acceptance
- **WHEN** sufficient historical similar samples exist
- **THEN** P88 SHALL prove probabilities are derived from sample proportions, the base scenario is the highest-frequency path, pessimistic scenario is displayed, and the report shows target name/code, future-12-month ranges, sample metadata, triggers, and disclaimer
- **AND** it SHALL prove a representative holding-class coverage matrix covering broad ETF/index fund, sector/growth ETF or fund, and equity/security-like constituent-financial path before upgrading `REQ-09-001`
- **AND** when samples are fewer than five, P88 SHALL prove no return range is generated and a supplement-data list is displayed.

#### Scenario: Expected-return dynamic monitoring is proven

- **GIVEN** valuation, fundamentals, market state, assumptions, or actual path data change
- **WHEN** the expected-return monitoring path runs
- **THEN** P88 SHALL prove affected scenario probabilities are lowered when applicable
- **AND** it SHALL prove periodic assumption checks, two-month below-expectation downshift warning, and one-month pessimistic-path manual probability-adjustment suggestion.

#### Scenario: Quarterly rebalance remains manual and auditable

- **GIVEN** a portfolio drifts beyond quarterly +/-15% target bands
- **WHEN** the rebalance flow runs
- **THEN** P88 SHALL prove manual buy/sell recommendation amounts through UI/API/SQLite/audit readback
- **AND** it SHALL NOT create broker orders, trades, automatic confirmations, or external push events.

#### Scenario: SOP addendum proposal is generated without automatic rule application

- **GIVEN** repeated review/error-case evidence identifies a high-frequency uncovered scenario
- **WHEN** P88 runs the SOP addendum path
- **THEN** it SHALL create a pending `sop` rule proposal, notification, and audit event
- **AND** it SHALL NOT modify active rules unless the existing gatekeeper and final user-confirmation flow is explicitly completed.

#### Scenario: P88 final claims remain evidence gated

- **GIVEN** P88 writes final release materials
- **WHEN** any full-release-required row remains partial, blocked, scoped-only, unsupported, or unverified
- **THEN** P88 SHALL NOT claim full original-requirement pass
- **AND** it SHALL list exact remaining rows and blockers.
- **AND** any row reclassification SHALL require explicit L1/OpenSpec rationale and SHALL NOT be reported as equivalent to 27/27 `real_pass`.
