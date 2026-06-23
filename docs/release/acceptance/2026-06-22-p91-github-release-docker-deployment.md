# P91 GitHub Release And Docker Deployment Acceptance

## Status

`passed`

## Scope

P91 adds a GitHub Release and Docker Compose deployment path with first-install, upgrade, uninstall, backup, status, and doctor scripts. It does not add runtime investment behavior.

## Required Evidence

- `python3 scripts/p91_deployment_check.py` -> passed.
- `bash -n scripts/deploy-lib.sh scripts/install.sh scripts/upgrade.sh scripts/uninstall.sh scripts/backup.sh scripts/status.sh scripts/doctor.sh docker/entrypoint.sh docker/healthcheck.sh` -> passed.
- `go test ./internal/infrastructure/config -run TestLoadDeepSeekDeploymentEnvOverrides -count=1` -> passed.
- `go test ./internal/infrastructure/config -count=1` -> passed.
- `go test ./...` -> passed.
- `openspec validate p91-github-release-docker-deployment --strict` -> passed.
- `openspec validate --all --strict` -> passed.
- `npm --prefix web test` -> 48 files / 176 tests passed.
- `npm --prefix web run build` -> passed.
- `bash scripts/local-release-package.sh --release-label p91-github-docker-deployment --output-dir tmp/p91-release-final` -> passed.
- `bash scripts/local-release-package.sh --verify tmp/p91-release-final/*/investment-agent-p91-github-docker-deployment.tar.gz --output-dir tmp/p91-release-final` -> passed.
- `git diff --check` -> passed.

## Real Docker Deployment Evidence

Temporary acceptance configuration:

- `COMPOSE_PROJECT_NAME=p91-deploy-acceptance`
- `INVESTMENT_AGENT_WEB_PORT=19173`
- `INVESTMENT_AGENT_SERVER_PORT=19080`
- `INVESTMENT_AGENT_DATA_DIR=tmp/p91-deploy-acceptance`
- `DEEPSEEK_API_KEY=` blank, verifying safe LLM degradation path for deployment without embedding secrets.

Executed real local Docker workflow:

- `bash scripts/install.sh` on empty deployment state -> first install started, image built, container started, health passed at `http://127.0.0.1:19173`.
- second `bash scripts/install.sh` -> detected existing install and routed through `scripts/upgrade.sh`.
- upgrade generated backup under `tmp/p91-deploy-acceptance/backups/*-upgrade`.
- `bash scripts/status.sh` -> container healthy; Compose port output showed `127.0.0.1:19173->4173/tcp` and `127.0.0.1:19080->8080/tcp`.
- `bash scripts/uninstall.sh` -> container/network removed and `tmp/p91-deploy-acceptance/release-state.json` remained, confirming data preservation by default.

## Final Package Evidence

Initial P91 package smoke:

- archive: `tmp/p91-release-final/20260622T101815Z/investment-agent-p91-github-docker-deployment.tar.gz`
- sha256: `ddd62f6ec1a5674cdd4bb7bb1691aba04b0c9942c6e82ed17722abc9f90d8589`
- verify status: `passed`

After P91 archive, a final package refresh was generated and verified. The exact archive checksum is recorded in the generated `release-manifest.json`; operators should use the manifest inside the selected package as the checksum truth because changing this acceptance record changes the archive hash.

## Safety Boundary

P91 does not claim physical second-machine validation, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, Level2 data, paid/login sources, future provider availability, or investment returns.
