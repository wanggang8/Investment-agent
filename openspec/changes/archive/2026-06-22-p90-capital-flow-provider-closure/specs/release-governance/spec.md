## ADDED Requirements

### Requirement: P90 capital-flow provider closure

After P89, P90 SHALL resolve or explicitly preserve the two remaining capital-flow related full-release-required rows using real public provider evidence and product UI/API/SQLite readback.

#### Scenario: P90 row inventory starts from the P89 remainder

- **GIVEN** P90 starts from the P89 matrix
- **WHEN** the inventory gate runs
- **THEN** it SHALL find exactly two full-release-required non-`real_pass` rows
- **AND** the row IDs SHALL be `REQ-04-016` and `REQ-05-003`.

#### Scenario: Capital-flow rows require a verified public runtime provider

- **GIVEN** P90 evaluates capital-flow fields
- **WHEN** it claims a row as `real_pass`
- **THEN** it SHALL prove a no-login/no-paid/no-authorization/no-Level2/no-high-frequency runtime provider was verified and used
- **AND** it SHALL prove `date`, `net_inflow`, and `net_outflow` were persisted in SQLite and read back through product APIs or UI
- **AND** parser-only, fixture, stub, accepted-local, or manually seeded evidence SHALL NOT upgrade those rows.

#### Scenario: Directional net-flow semantics are explicit

- **GIVEN** the public H5 capital-flow history endpoint exposes a daily net-flow value
- **WHEN** P90 stores the value
- **THEN** positive net flow SHALL map to `net_inflow`
- **AND** negative net flow SHALL map to `net_outflow`
- **AND** the raw daily net-flow value SHALL be preserved as `raw_net_flow`.

#### Scenario: P90 final claims remain evidence gated

- **GIVEN** P90 writes final release materials
- **WHEN** any full-release-required row remains partial, blocked, scoped-only, unsupported, or unverified
- **THEN** P90 SHALL NOT claim full original-requirement pass
- **AND** it SHALL list exact remaining rows and blockers.
