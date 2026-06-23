## ADDED Requirements

### Requirement: P96 SHALL provide a public documentation front door

P96 SHALL make the public repository understandable to a new reader without requiring them to start from archived phase logs.

#### Scenario: Root README introduces the product honestly

- **GIVEN** a user opens the GitHub repository
- **WHEN** they read the root `README.md`
- **THEN** they SHALL see the product purpose, supported local workflows, architecture/data-flow visuals, installation entrypoint, documentation map, CI/release status, and safety boundaries
- **AND** the README SHALL NOT claim broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.

#### Scenario: Documentation map is concise

- **GIVEN** a maintainer opens `docs/README.md`
- **WHEN** they use it as a navigation page
- **THEN** it SHALL point to product, architecture, API, data model, workflow, frontend, deployment, governance, release evidence, and history documents
- **AND** it SHALL NOT require reading the full P0-P96 phase log to find normal documentation.

#### Scenario: Historical release evidence remains available

- **GIVEN** P96 moves or summarizes phase history
- **WHEN** a reader needs release caveats or acceptance history
- **THEN** the history SHALL remain discoverable under `docs/release/`
- **AND** P96 SHALL NOT erase historical limitations, scoped claims, or Not Claimed boundaries.

#### Scenario: Requirements truth source remains stable

- **GIVEN** `docs/requirements.md` is the L1 product requirement truth source
- **WHEN** P96 adds public-facing docs
- **THEN** public docs SHALL link to requirements for full details
- **AND** P96 SHALL NOT rewrite L1 requirement semantics as marketing copy.
