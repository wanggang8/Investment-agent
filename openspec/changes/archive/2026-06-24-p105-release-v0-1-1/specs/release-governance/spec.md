## MODIFIED Requirements

### Requirement: Initial release version marker

The project SHALL preserve `v0.1.0` as the historical initial local release version marker recorded by P99, while allowing the current repository release marker to advance in later release changes.

#### Scenario: Initial version history is preserved

- **GIVEN** release history or P99 acceptance materials are inspected
- **WHEN** the initial local release marker is described
- **THEN** the materials SHALL identify `v0.1.0` as the initial local release version recorded by P99
- **AND** the materials SHALL NOT treat that historical marker as proof that the current root `VERSION` file must remain `v0.1.0`.

## ADDED Requirements

### Requirement: P105 current release version v0.1.1

The repository SHALL declare `v0.1.1` as the current local source release version after P100-P104 validation has passed and P105 release gates have completed.

#### Scenario: Current version metadata is synchronized

- **GIVEN** P105 release validation has passed
- **WHEN** a user or release operator inspects version metadata
- **THEN** the root `VERSION` file SHALL contain `v0.1.1`
- **AND** `web/package.json` SHALL declare version `0.1.1`
- **AND** the root package entry in `web/package-lock.json` SHALL declare version `0.1.1`.

#### Scenario: P105 release claims stay bounded

- **GIVEN** `v0.1.1` is described in release materials
- **WHEN** release readiness is communicated
- **THEN** the project MAY claim local source product acceptance through P104 and current source version metadata synchronization
- **AND** it SHALL NOT claim Docker installation validation, package refresh, GitHub Release workflow success, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns unless separately validated.
