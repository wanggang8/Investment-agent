# Quickstart

> Updated: 2026-06-23

This guide covers the local Docker Compose path and the most common operations. The detailed deployment contract is [deployment.md](deployment.md).

## Requirements

- Docker with Compose support.
- A local checkout or extracted release package.
- Optional: `DEEPSEEK_API_KEY` for LLM-backed analysis. If it is empty, LLM-backed analysis degrades safely where required.

## Configure

```bash
cp .env.example .env
```

Edit `.env`:

```dotenv
INVESTMENT_AGENT_WEB_PORT=4173
INVESTMENT_AGENT_SERVER_PORT=8080
INVESTMENT_AGENT_DATA_DIR=.investment-agent
DEEPSEEK_API_KEY=
DEEPSEEK_BASE_URL=https://api.deepseek.com
DEEPSEEK_MODEL=deepseek-chat
DEEPSEEK_TIMEOUT_SECONDS=30
INVESTMENT_AGENT_USE_STUB_DATA=false
```

Keep `.env` local. Do not commit real keys.

## Start With Docker Compose

```bash
bash scripts/install.sh
```

`install.sh` creates runtime directories, runs diagnostics, builds and starts the Compose service, waits for health, and prints next steps. If it detects an existing install, it routes to `scripts/upgrade.sh`.

Open the UI at the printed URL, usually:

```text
http://127.0.0.1:4173
```

Compose binds the web and backend ports to `127.0.0.1` by default.

## Health And Operations

```bash
bash scripts/status.sh
bash scripts/doctor.sh
bash scripts/backup.sh
```

Backups are written under:

```text
.investment-agent/backups/
```

The main health endpoint is:

```text
http://127.0.0.1:8080/api/v1/health
```

## Upgrade

```bash
bash scripts/upgrade.sh
```

Upgrade preserves `.env`, SQLite, VecLite, logs, and backups. It runs:

```bash
bash scripts/backup.sh --reason upgrade
```

before recreating containers.

## Uninstall

Stop containers while preserving local data:

```bash
bash scripts/uninstall.sh
```

Delete local data only with explicit purge:

```bash
bash scripts/uninstall.sh --purge
```

The purge path requires typing this exact confirmation phrase:

```text
DELETE INVESTMENT AGENT DATA
```

## Source Development

For backend work:

```bash
go test $(bash scripts/go-packages.sh)
go run ./cmd/server
```

For frontend work:

```bash
npm --prefix web install
npm --prefix web run dev
npm --prefix web test
npm --prefix web run build
```

Use the contracts in [api.md](api.md), [frontend-contract.md](frontend-contract.md), and [data-model.md](data-model.md) when changing behavior.

## Troubleshooting

| Symptom | Check |
| --- | --- |
| Installer says Docker is unavailable | Run `bash scripts/doctor.sh` and confirm Docker/Compose are installed and running. |
| UI does not open | Run `bash scripts/status.sh`; confirm `INVESTMENT_AGENT_WEB_PORT` is not already used. |
| Backend health fails | Open `http://127.0.0.1:8080/api/v1/health`; check Compose logs if unhealthy. |
| LLM analysis is degraded | Confirm `DEEPSEEK_API_KEY`, `DEEPSEEK_BASE_URL`, `DEEPSEEK_MODEL`, and timeout settings in `.env`. |
| Data seems missing after reinstall | Default uninstall preserves `.investment-agent`; confirm `INVESTMENT_AGENT_DATA_DIR` points to the intended directory. |

## Safety Boundary

The local runtime does not add broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, paid/login/authorization-only sources, Level2 data, high-frequency data, or return guarantees.
