## ADDED Requirements

### Requirement: P89 real provider and dynamic probability closure

After P88, P89 SHALL resolve or explicitly preserve the 10 remaining full-release-required rows by adding real provider evidence for structured data categories and real UI/API/SQLite evidence for dynamic expected-return monitoring.

#### Scenario: P89 row inventory starts from the P88 remainder

- **GIVEN** P89 starts from the P88 matrix
- **WHEN** the inventory gate runs
- **THEN** it SHALL find exactly 10 full-release-required non-`real_pass` rows
- **AND** the row IDs SHALL be `REQ-04-016`, `REQ-05-003`, `REQ-05-004`, `REQ-05-005`, `REQ-08-004`, `REQ-08-023`, `REQ-09-004`, `REQ-09-023`, `REQ-09-024`, and `REQ-09-025`.

#### Scenario: Structured field rows require verified real providers

- **GIVEN** P89 evaluates capital-flow, margin-financing, or constituent-financial rows
- **WHEN** it claims a row as `real_pass`
- **THEN** it SHALL prove a no-login/no-paid/no-authorization/no-Level2/no-high-frequency runtime provider was verified and used
- **AND** it SHALL prove the required fields were persisted in SQLite and read back through product APIs or UI
- **AND** parser-only, fixture, stub, accepted-local, or manually seeded evidence SHALL NOT upgrade those rows.

#### Scenario: Extreme fear locks active trading and shows history

- **GIVEN** an extreme-fear sentiment state exists for a held symbol
- **WHEN** the product evaluates that symbol through real UI/API acceptance
- **THEN** it SHALL lock active trading advice
- **AND** it SHALL display historical similar-scenario context
- **AND** the decision and context SHALL be readable from SQLite.

#### Scenario: Dynamic expected-return monitoring is proven end to end

- **GIVEN** baseline expected-return probabilities and supporting assumptions exist
- **WHEN** valuation, fundamentals, market state, assumptions, or actual path data change
- **THEN** P89 SHALL prove a rerun through real UI/API/SQLite lowers affected probabilities when applicable
- **AND** it SHALL prove periodic assumption checks, a two-month below-expectation downshift warning, and a one-month pessimistic-path manual probability-adjustment suggestion.

#### Scenario: P89 final claims remain evidence gated

- **GIVEN** P89 writes final release materials
- **WHEN** any full-release-required row remains partial, blocked, scoped-only, unsupported, or unverified
- **THEN** P89 SHALL NOT claim full original-requirement pass
- **AND** it SHALL list exact remaining rows and blockers.
