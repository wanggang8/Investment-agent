# Deployment

> Updated: 2026-06-23

Investment Agent supports a local Docker Compose deployment path. The deployment package is designed for a personal machine, NAS, or VPS where the user controls local storage and secrets.

## Quick Start

```bash
tar -xzf investment-agent-<label>.tar.gz
cd investment-agent-<label>
cp .env.example .env
# Edit .env and set DEEPSEEK_API_KEY, or set DEEPSEEK_API_KEY_FILE.
bash scripts/install.sh
```

Open the UI at the URL printed by `install.sh`, usually:

```text
http://127.0.0.1:4173
```

The default Compose ports bind to `127.0.0.1` only. Put a trusted reverse proxy in front of the service if you intentionally need remote access.

## Configuration

Runtime secrets are local-only. Do not commit `.env`.

Required or common settings:

- `DEEPSEEK_API_KEY`: LLM provider key. If empty, LLM-backed analysis degrades safely where required.
- `DEEPSEEK_API_KEY_FILE`: optional path to a local file containing only the LLM provider key. `DEEPSEEK_API_KEY` takes precedence when both are set. This supports Docker/Compose secret-style mounts without committing keys.
- `DEEPSEEK_BASE_URL`: default `https://api.deepseek.com`.
- `DEEPSEEK_MODEL`: default `deepseek-chat`.
- `DEEPSEEK_TIMEOUT_SECONDS`: default `30`.
- `INVESTMENT_AGENT_WEB_PORT`: default `4173`.
- `INVESTMENT_AGENT_SERVER_PORT`: default `8080`.
- `INVESTMENT_AGENT_DATA_DIR`: persistent local data directory.

Example with a mounted secret file:

```bash
mkdir -p .investment-agent/data
printf '%s' '<your-key>' > .investment-agent/data/deepseek_api_key
chmod 600 .investment-agent/data/deepseek_api_key
DEEPSEEK_API_KEY_FILE=/data/deepseek_api_key docker compose up -d --build
```

When using this pattern, mount the file into the container through a local Compose override or another trusted local secret mechanism. Do not commit secret files.

## Install Versus Upgrade

`scripts/install.sh` is the normal entry point. It detects the local state:

- first install: creates `.env` when needed, creates data directories, starts Docker Compose, waits for health, and writes release state;
- existing install: routes to `scripts/upgrade.sh`.

Upgrade preserves `.env`, SQLite, VecLite, logs, and backups. It runs `scripts/backup.sh --reason upgrade` before recreating containers.

## Uninstall

Default uninstall preserves data:

```bash
bash scripts/uninstall.sh
```

Delete local data only with explicit purge:

```bash
bash scripts/uninstall.sh --purge
```

The purge path requires typing:

```text
DELETE INVESTMENT AGENT DATA
```

## Operations

```bash
bash scripts/status.sh
bash scripts/doctor.sh
bash scripts/backup.sh
```

Backups are written under:

```text
.investment-agent/backups/
```

## GitHub Actions

The public repository is expected to use GitHub Actions as the release gate:

- pull requests and pushes to `main` run OpenSpec validation, backend Go package selection through `scripts/go-packages.sh`, `go vet`, bounded `golangci-lint`, Go tests, frontend lint/test/build, deployment checks, P92/P93 audit checks, API route contract checks, release package smoke, and whitespace checks;
- tags matching `v*` run the release workflow and upload the local deployment package plus manifest as GitHub release artifacts;
- the security scan runs on PR/main pushes and weekly schedule with `govulncheck`, frontend production dependency audit, and P93 code reality / secret checks.

## Safety Boundary

The deployment workflow does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration without preflight, real database overwrite, paid/login/authorization-only sources, Level2 data, high-frequency data, or return guarantees.
