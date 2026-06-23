## ADDED Requirements

### Requirement: P64 SHALL produce a local release package manifest

P64 SHALL provide a local release package manifest that identifies the release label, source commit, package archive, checksum, included roots, excluded patterns, verification commands, acceptance references, known degradations, Not Claimed boundaries, and safety note.

#### Scenario: Local release package is generated

- **WHEN** the operator runs the P64 local release package command with a release label
- **THEN** the command SHALL stage release-safe tracked project files under `tmp/`
- **AND** it SHALL write a sanitized `release-manifest.json`
- **AND** it SHALL create a compressed local archive and SHA-256 checksum
- **AND** it SHALL NOT include local private config, temporary SQLite databases, VecLite local indexes, logs, traces, `.cursor/`, `tmp/`, `cmd/agent/tmp/`, `docs/release/ui-audit-assets/`, `web/node_modules/`, `web/dist/`, complete API keys, private paths, complete prompts, raw SQL dumps, or raw vendor payloads.

#### Scenario: Local release package is verified

- **WHEN** the operator runs package verification against the generated archive
- **THEN** verification SHALL parse the manifest, check archive checksum consistency, confirm required package entrypoints are present, and reject forbidden paths or file patterns
- **AND** verification SHALL NOT run migrations, restore data, overwrite databases, call public providers, call LLM providers, execute trades, push notifications, apply rules, or repair files automatically.

#### Scenario: Release package is handed off

- **WHEN** P64 updates release materials
- **THEN** the handoff SHALL reference the package manifest, package verification command, P63 acceptance evidence, known non-blocking degradations, and repeat verification entrypoints
- **AND** the handoff SHALL NOT claim future provider availability, investment returns, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic upgrade, automatic migration, real database overwrite, login sources, paid sources, authorized sources, Level2 data, or high-frequency data.
