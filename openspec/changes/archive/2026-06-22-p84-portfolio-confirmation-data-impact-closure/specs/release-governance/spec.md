## ADDED Requirements

### Requirement: P84 portfolio confirmation data-impact closure

After P83, portfolio and manual-confirmation rows SHALL NOT be marked `real_pass` unless real local UI workflows prove user action, local data mutation, downstream readback, deterministic value accuracy where applicable, and safety boundaries.

#### Scenario: P84 row inventory is complete before execution

- **GIVEN** P84 starts from the P83 evidence matrix
- **WHEN** execution begins
- **THEN** the P84 plan SHALL enumerate exactly 35 portfolio/confirmation rows
- **AND** each row SHALL map to before/after data impact and downstream readback evidence.

#### Scenario: Portfolio mutation is local and manual

- **GIVEN** a P84 scenario changes portfolio-related state
- **WHEN** the change is accepted
- **THEN** the evidence SHALL show local user-driven action, API response, SQLite before/after delta, audit event, and UI readback
- **AND** it SHALL NOT depend on broker synchronization, automatic trading, order placement, or automatic confirmation.

#### Scenario: Derived values are checked independently

- **GIVEN** P84 evidence includes market value, cost, quantity, ratio, cash, or profit/loss values
- **WHEN** the row is evaluated
- **THEN** P84 SHALL compare product output with independently computed expectations
- **AND** future return or market-direction accuracy SHALL NOT be claimed.
