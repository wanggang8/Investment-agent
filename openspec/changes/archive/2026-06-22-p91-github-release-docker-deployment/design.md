# P91 Design

## Deployment Shape

P91 uses Docker Compose as the first supported deployment target. A release package can be downloaded from GitHub, unpacked, configured through `.env`, and started with `bash scripts/install.sh`.

The runtime uses one application image that contains:

- the Go `cmd/server` binary;
- the built React frontend in `/app/web/dist`;
- `busybox httpd` for static frontend serving;
- `configs/config.docker.yaml` for container-safe defaults.

The backend listens on `0.0.0.0:8080` inside the container. The frontend listens on `0.0.0.0:4173` inside the container and talks to the backend through the same Compose network when run from the browser proxy path used by Vite-built assets.

## Local Data Layout

Deployment scripts default to `.investment-agent/` in the extracted package:

- `.investment-agent/data/sqlite/investment-agent.db`
- `.investment-agent/data/veclite/`
- `.investment-agent/backups/`
- `.investment-agent/logs/`
- `.investment-agent/release-state.json`

The scripts never delete these paths during install or upgrade. `uninstall.sh` removes containers and networks by default, and only deletes data with `--purge` plus exact confirmation text.

## Configuration And Secrets

`.env.example` documents safe defaults. `.env` is user-local and ignored by Git. Required secrets, especially `DEEPSEEK_API_KEY`, are read from environment variables and are not baked into Docker images, release packages, or GitHub workflows.

## Script Responsibilities

- `install.sh`: detect first install vs upgrade, create local directories, ensure `.env`, run doctor, start Compose, wait for health, write release state.
- `upgrade.sh`: backup, rebuild/pull/recreate containers, run health, preserve `.env` and data.
- `uninstall.sh`: stop/remove containers by default; purge only with explicit confirmation.
- `backup.sh`: copy SQLite, VecLite, `.env`, and release state to timestamped backup path.
- `status.sh`: show Compose service state and health endpoint result.
- `doctor.sh`: verify Docker/Compose, `.env`, data directories, config, port variables, and key presence without printing secrets.

## GitHub Release

`.github/workflows/ci.yml` runs OpenSpec, Go tests, frontend tests/build, deployment checks, and package verification. `.github/workflows/release.yml` runs on tags, builds the local release package, verifies it, and uploads package artifacts to GitHub Release. No workflow stores or prints runtime secrets.

## Evidence

- `scripts/p91_deployment_check.py`
- `docs/release/acceptance/2026-06-22-p91-github-release-docker-deployment.md`
- `docs/release/ui-audit-assets/2026-06-22-p91-github-release-docker-deployment/final-validation.log`
