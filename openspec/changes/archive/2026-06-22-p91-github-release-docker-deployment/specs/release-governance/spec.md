## ADDED Requirements

### Requirement: P91 GitHub release and Docker deployment

After P90, the project SHALL provide a GitHub-ready release and Docker Compose deployment path that can initialize and run the product without embedding secrets or overwriting user data.

#### Scenario: Install script detects first install versus upgrade

- **GIVEN** a user runs `bash scripts/install.sh`
- **WHEN** no local deployment state or data directory exists
- **THEN** the script SHALL initialize local deployment directories, create `.env` from `.env.example` when needed, start Docker Compose, and run health checks
- **AND** when existing deployment state or data is present it SHALL route through the upgrade path instead of deleting or reinitializing user data.

#### Scenario: Runtime secrets are supplied outside the package

- **GIVEN** the release package and Docker image are built
- **WHEN** the user configures LLM credentials
- **THEN** `DEEPSEEK_API_KEY`, base URL, model, and timeout SHALL come from `.env` or environment variables
- **AND** no complete API key SHALL be committed, baked into the image, or written to release manifests.

#### Scenario: Uninstall preserves data by default

- **GIVEN** a deployed instance has local SQLite, VecLite, backup, log, and `.env` data
- **WHEN** the user runs `bash scripts/uninstall.sh`
- **THEN** containers and networks MAY be removed
- **AND** local data SHALL be preserved by default
- **AND** deleting local data SHALL require `--purge` and an exact confirmation phrase.

#### Scenario: GitHub release automation remains evidence gated

- **GIVEN** GitHub Actions creates a release artifact
- **WHEN** CI or release packaging runs
- **THEN** it SHALL run OpenSpec validation, Go tests, frontend tests/build, deployment checks, and package verification
- **AND** it SHALL upload release artifacts without claiming physical second-machine validation, broker connectivity, trading, automatic confirmation, external push, Level2 data, paid/login sources, or return guarantees.
