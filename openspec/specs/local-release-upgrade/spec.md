# local-release-upgrade Specification

## Purpose
定义本地发布与升级体验的只读检查能力：在升级前后提供版本意图、备份提醒、迁移文件预检、升级后 smoke 命令和脱敏诊断汇总，同时保持不自动升级、不运行迁移、不覆盖真实库、不交易、不外推和不自动修复的安全边界。
## Requirements
### Requirement: P49 SHALL provide a local release and upgrade readiness report

The system SHALL provide a local release and upgrade readiness report that helps an operator review version intent, backup readiness, migration precheck status, and post-upgrade smoke commands before changing a local installation.

#### Scenario: Upgrade readiness check is requested
- **WHEN** a user runs the P49 local release/upgrade check
- **THEN** the system SHALL report current version, target version, backup reminder, migration file precheck status, post-upgrade smoke commands, and safety note
- **AND** it SHALL make clear that the check did not perform an upgrade, did not run migrations, did not create a backup, and did not restore or overwrite a database.

#### Scenario: Target version is missing
- **WHEN** a user runs the release/upgrade check without a target version
- **THEN** the report SHALL return a warning status that asks the user to provide a target version or release label
- **AND** it SHALL still provide backup and smoke checklist guidance without modifying local data.

### Requirement: P49 SHALL keep release and upgrade diagnostics local, read-only, and sanitized

P49 SHALL expose release and upgrade diagnostics through local CLI and script summaries without leaking sensitive local details or creating side effects.

#### Scenario: Diagnostics JSON is written
- **WHEN** the release/upgrade check writes a diagnostics JSON file
- **THEN** the file SHALL contain only sanitized status, commands, counts, basenames or placeholders, and safety notes
- **AND** it SHALL NOT contain complete API keys, private absolute paths, raw SQL, full prompts, raw HTTP exchanges, private keys, supplier raw responses, account secrets, or broker credentials.

#### Scenario: Local install diagnostics includes release upgrade check
- **WHEN** the local install diagnostics script is explicitly run with release/upgrade inclusion enabled
- **THEN** it SHALL include the release/upgrade check as an additional local step
- **AND** the default local install diagnostics behavior SHALL remain unchanged unless the user explicitly enables that step.

### Requirement: P49 SHALL preserve upgrade safety boundaries

P49 SHALL remain a local upgrade planning and verification aid and SHALL NOT expand automation, trading, data acquisition, or repair boundaries.

#### Scenario: Upgrade readiness detects missing backup or migration uncertainty
- **WHEN** the readiness report detects a missing SQLite file, missing target version, missing migration files, or another warning condition
- **THEN** it SHALL classify the report as warning or blocked according to the check severity
- **AND** it SHALL NOT automatically create backups, run migrations, repair files, restore data, mark the installation healthy, execute trades, apply rules, confirm operations, call LLMs, access public endpoints, or promise returns.

#### Scenario: Post-upgrade smoke plan is shown
- **WHEN** the report lists post-upgrade smoke commands
- **THEN** those commands SHALL remain manual commands for the user to run explicitly
- **AND** they SHALL NOT include broker APIs, external push channels, one-click trading, delegated order placement, login-gated data sources, paid data sources, authorization-gated sources, Level2 feeds, or high-frequency polling.

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

### Requirement: P65 SHALL repeat release acceptance from an isolated package workspace

P65 SHALL provide a local cross-machine-equivalent repeat acceptance flow that starts from a P65 candidate archive generated by the P64 release package workflow and verifies the package from an isolated extracted workspace.

#### Scenario: Package repeat acceptance is started

- **WHEN** the operator runs the P65 repeat acceptance command with a package archive
- **THEN** the command SHALL verify the package archive and adjacent sidecar manifest before extraction
- **AND** it SHALL extract the package into a fresh directory under project `tmp/`
- **AND** it SHALL execute repeat commands from the extracted package root rather than from the active repository checkout
- **AND** it SHALL write a sanitized repeat summary with package identity, release label, source commit, source status, command results, output paths, skip flags, known caveats, and safety note.

#### Scenario: Isolated repeat commands run

- **WHEN** the repeat flow executes validation commands
- **THEN** it SHALL run OpenSpec validation, Go tests, frontend dependency installation, frontend tests, frontend build, and local E2E smoke from the extracted package root
- **AND** it SHALL use temporary local SQLite, VecLite, server, and web ports for smoke execution
- **AND** it SHALL NOT write to real user databases, private configuration directories, home directories, or active repository source files.

#### Scenario: Repeat acceptance is handed off

- **WHEN** P65 updates release materials
- **THEN** the handoff SHALL reference the package archive, sidecar manifest, repeat summary, command matrix, known caveats, and physical second-machine follow-up commands
- **AND** the handoff SHALL NOT claim remote publishing, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic restore, automatic repair, real database overwrite, broker connectivity, automatic trading, one-click trading, order delegation, delegated order placement, external push, automatic confirmation, automatic rule application, login-gated sources, paid sources, authorization-gated sources, Level2 data, high-frequency data, future provider availability, or investment returns.

### Requirement: P69 SHALL refresh final package evidence from a clean tree

P69 SHALL regenerate package evidence from a clean committed source tree after P65-P68 so final package claims are not based on stale dirty candidate archives.

#### Scenario: Clean package is generated from the P68 source commit

- **GIVEN** P68 has been committed and the main repository has a clean HEAD
- **WHEN** P69 generates final package evidence
- **THEN** the package SHALL be generated from a clean detached worktree or equivalent clean checkout at the P68 source commit
- **AND** the package manifest SHALL record `source_status=clean`
- **AND** the package manifest SHALL record the source commit used for the package
- **AND** P69 release materials SHALL NOT claim the generated archive includes P69 documentation unless a later package refresh is performed after P69 commit.

#### Scenario: Clean package is verified and repeated

- **WHEN** the P69 package archive is generated
- **THEN** package verification SHALL confirm archive checksum consistency, required entries, forbidden path exclusions, and manifest safety boundaries
- **AND** repeat acceptance SHALL run from an extracted package workspace rather than from the active repository checkout
- **AND** repeat acceptance SHALL cover OpenSpec validation, Go tests, frontend dependency installation, frontend tests, frontend build, and local E2E smoke.

#### Scenario: Package evidence is handed off

- **WHEN** P69 updates release materials
- **THEN** the materials SHALL include package identity, source commit, `source_status`, checksum, archive entry count, verify result, repeat command matrix, known caveats, and Not Claimed boundaries
- **AND** the materials SHALL preserve P68 `release_ready_limited_current_data_scope`
- **AND** the materials SHALL NOT claim P66 current-data policy passed, physical second-machine execution, remote publishing, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic restore, automatic repair, real database overwrite, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, provider availability, or investment returns.
