# P91 Tasks

## 1. Governance And Red Tests

- [x] Create P91 OpenSpec change and validate it.
- [x] Add failing Go config test for DeepSeek env overrides.
- [x] Add failing P91 deployment checker for Docker/Compose/scripts/GitHub workflow safety.

## 2. Runtime Configuration

- [x] Add env overrides for `DEEPSEEK_BASE_URL`, `DEEPSEEK_MODEL`, and `DEEPSEEK_TIMEOUT_SECONDS`.
- [x] Add `configs/config.docker.yaml` with container-safe paths and no embedded secrets.

## 3. Docker Deployment

- [x] Add `.dockerignore`.
- [x] Add `Dockerfile`.
- [x] Add `docker-compose.yml`.
- [x] Add `docker/entrypoint.sh` and `docker/healthcheck.sh`.

## 4. Install / Upgrade / Uninstall Scripts

- [x] Add shared deployment library `scripts/deploy-lib.sh`.
- [x] Add `scripts/install.sh` that detects first install vs upgrade.
- [x] Add `scripts/upgrade.sh` that backs up before upgrade.
- [x] Add `scripts/uninstall.sh` that preserves data by default and requires explicit confirmation for purge.
- [x] Add `scripts/backup.sh`, `scripts/status.sh`, and `scripts/doctor.sh`.

## 5. GitHub Release Automation

- [x] Add `.github/workflows/ci.yml`.
- [x] Add `.github/workflows/release.yml`.
- [x] Ensure release package workflow includes deployment files and excludes secrets/local data.

## 6. Documentation And Acceptance

- [x] Add deployment documentation.
- [x] Generate P91 acceptance record.
- [x] Run P91 deployment checker.
- [x] Run OpenSpec, Go tests, frontend tests/build, release package verify, and `git diff --check`.
- [x] Archive P91 after final validation passes.
