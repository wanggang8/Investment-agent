## ADDED Requirements

### Requirement: P94 GitHub CI and release hardening

Before the repository is published publicly, the project SHALL provide GitHub Actions quality gates for PR/main changes, tag-based release packaging, and independent security scanning.

#### Scenario: PR and main pushes run full quality gates

- **GIVEN** a pull request or push to `main`
- **WHEN** GitHub Actions runs CI
- **THEN** it SHALL run OpenSpec validation, `go vet`, bounded `golangci-lint`, Go tests, frontend lint, frontend tests, frontend build, P91 deployment checks, P92/P93 audit checks, release package smoke verification, and whitespace checks.

#### Scenario: Version tags build release artifacts

- **GIVEN** a tag matching `v*` is pushed
- **WHEN** the release workflow runs
- **THEN** it SHALL run release preflight checks before packaging
- **AND** it SHALL upload the local deployment package and manifest as release artifacts without embedding runtime LLM secrets.

#### Scenario: Security scan remains separate and repeatable

- **GIVEN** a PR/main push or scheduled weekly run
- **WHEN** the security workflow runs
- **THEN** it SHALL run Go vulnerability scanning, frontend production dependency audit, and P93 code reality / secret checks
- **AND** it SHALL NOT require broker connectivity, trading, paid/login sources, Level2 data, high-frequency data, or runtime LLM keys.
