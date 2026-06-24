## ADDED Requirements

### Requirement: P107 SHALL provide bilingual README discovery

P107 SHALL make the public README entrypoint understandable to both English and Simplified Chinese readers.

#### Scenario: Root README exposes a language switch

- **GIVEN** a reader opens the root `README.md`
- **WHEN** they look for language options
- **THEN** the README SHALL provide a visible link to `README.zh-CN.md`
- **AND** the root README SHALL continue to provide the existing English product introduction and documentation map.

#### Scenario: Chinese README introduces the product honestly

- **GIVEN** a Simplified Chinese reader opens `README.zh-CN.md`
- **WHEN** they read the project overview
- **THEN** they SHALL see the product purpose, supported feature areas, product flow, architecture summary, quickstart path, documentation links, CI/release status, governance notes, and safety boundaries
- **AND** the Chinese README SHALL NOT claim broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.

#### Scenario: README local links remain valid

- **GIVEN** the README language switch and documentation links are updated
- **WHEN** local Markdown links are checked
- **THEN** linked local files and diagram assets SHALL resolve within the repository.
