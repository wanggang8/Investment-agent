## ADDED Requirements

### Requirement: P97 SHALL default to local config.yaml

The local runtime SHALL treat `configs/config.yaml` as the preferred local configuration file and keep `configs/config.example.yaml` as a committed template.

#### Scenario: Local server uses config.yaml by default

- **GIVEN** `INVESTMENT_AGENT_CONFIG` is unset
- **AND** `configs/config.yaml` exists
- **WHEN** the server or shared config loader starts with no explicit config path
- **THEN** it SHALL load `configs/config.yaml`
- **AND** it SHALL NOT load `configs/config.example.yaml` instead.

#### Scenario: Fresh checkout remains runnable

- **GIVEN** `INVESTMENT_AGENT_CONFIG` is unset
- **AND** `configs/config.yaml` does not exist
- **WHEN** the shared config loader starts with no explicit config path
- **THEN** it SHALL fall back to `configs/config.example.yaml`.

#### Scenario: Local config is not committed

- **GIVEN** a user creates `configs/config.yaml`
- **WHEN** Git status is checked
- **THEN** the file SHALL be ignored by default
- **AND** real local keys SHALL remain outside committed source.
