## ADDED Requirements

### Requirement: P81 dynamic source field coverage closure

After P80, dynamic source field coverage rows SHALL NOT be marked `real_pass` unless fresh acceptance proves the current product obtains, evaluates, displays, or safely blocks the relevant data fields for a user-selected symbol through real local product paths.

#### Scenario: P81 row inventory is complete before execution

- **GIVEN** P81 starts from the P80 evidence matrix
- **WHEN** execution begins
- **THEN** the P81 plan SHALL enumerate exactly 59 dynamic source field coverage rows
- **AND** the plan SHALL preserve the previous status and target evidence type for each row.

#### Scenario: User-selected symbol drives source evidence

- **GIVEN** P81 evaluates data source coverage
- **WHEN** a browser or API scenario requests readiness or analysis for a symbol
- **THEN** the evidence SHALL show that the requested symbol drives the source/readiness result
- **AND** hard-coded `510300`-only evidence SHALL NOT be sufficient for P81 `real_pass`.

#### Scenario: Formal evidence unavailable

- **GIVEN** an external or built-in data category is missing, degraded, stale, or background-only
- **WHEN** P81 evaluates impacted features
- **THEN** the product SHALL safely degrade, qualify, or block affected claims
- **AND** the acceptance result SHALL NOT mark formal evidence requirements as passed by background knowledge alone.

#### Scenario: P81 claims remain bounded

- **GIVEN** P81 passes some or all rows
- **WHEN** release materials are updated
- **THEN** they SHALL state the exact upgraded rows and remaining non-`real_pass` count
- **AND** they SHALL NOT claim full original-requirement pass, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real DB overwrite, return promises, paid/login/authorized source, Level2 source, or high-frequency source.

