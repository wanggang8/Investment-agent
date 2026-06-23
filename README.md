# Investment Agent

Investment Agent is a local-first investment discipline cockpit for personal research, portfolio maintenance, evidence review, and manual decision governance. It combines public evidence collection, SQLite-backed local state, VecLite retrieval, deterministic rules, and LLM-assisted analysis into a workflow that keeps the final action with the user.

![Investment Agent cockpit](docs/diagrams/ui-cockpit-ia.png)

## What It Helps With

| Area | What the app does |
| --- | --- |
| Daily discipline | Shows today's status, pending manual actions, risk signals, notifications, and review work. |
| Portfolio maintenance | Records local account calibration, holdings, cash buckets, buy dates, watch states, and manual corrections. |
| Evidence review | Collects and indexes supported public evidence, local knowledge, market snapshots, and source health records. |
| Decision consultation | Runs a governed analysis flow with deterministic rules, RAG context, expected-return checks, and LLM analyst reports when configured. |
| Manual confirmation | Records the user's explicit final decision and downstream audit trail without placing orders. |
| Review and governance | Supports monthly/quarterly review, error marking, rule proposals, gatekeeper checks, and release evidence traceability. |

## Safety Boundary

Investment Agent is a research and decision-support system. It does not connect to brokers, trade automatically, provide one-click trading, place delegated orders, push external notifications, auto-confirm decisions, auto-apply rule changes, guarantee returns, require paid/login/authorization-only market data, use Level2 feeds, or perform high-frequency data collection.

## Product Flow

![Decision flow](docs/diagrams/ui-agent-decision-flow.png)

1. Maintain local portfolio and account facts.
2. Refresh supported public data and local evidence.
3. Review source health, risk alerts, and readiness gates.
4. Ask for a consultation when evidence is sufficient.
5. Read the analysis, rule checks, assumptions, and audit trail.
6. Confirm or reject manually; the app records the decision locally.
7. Use review and governance pages to learn from outcomes.

## Architecture At A Glance

![System architecture](docs/diagrams/system-architecture.png)

The application is a local Go + SQLite + VecLite backend with a React/Vite frontend. Docker Compose is the easiest way to run it on a personal machine, NAS, or VPS with user-controlled storage. More implementation detail lives in [docs/architecture.md](docs/architecture.md); P95 owns detailed architecture corrections, so this README intentionally stays high level.

## Quick Start

```bash
cp .env.example .env
# Edit .env and set DEEPSEEK_API_KEY if you want LLM-backed analysis.
bash scripts/install.sh
```

Open the URL printed by the installer, usually:

```text
http://127.0.0.1:4173
```

The default Docker Compose ports bind to `127.0.0.1`. See [docs/quickstart.md](docs/quickstart.md) for install, upgrade, uninstall, health checks, and troubleshooting. See [docs/deployment.md](docs/deployment.md) for the release deployment contract.

## Documentation

| Start here | Purpose |
| --- | --- |
| [Product overview](docs/product-overview.md) | User workflows, concepts, and product safety boundaries. |
| [Quickstart](docs/quickstart.md) | Docker Compose setup, local config, operations, upgrade, uninstall, and troubleshooting. |
| [Documentation map](docs/README.md) | Concise index for product, architecture, API, operations, governance, and release evidence. |
| [Requirements](docs/requirements.md) | L1 product requirements truth source. |
| [API contract](docs/api.md) | HTTP API contract. |
| [Frontend contract](docs/frontend-contract.md) | Frontend route and interaction contract. |
| [Release materials](docs/release/README.md) | Acceptance records, caveats, and release evidence. |
| [Release history](docs/release/history.md) | Moved phase history and historical caveats from the former docs README. |

## CI And Release Status

The repository includes GitHub Actions gates for OpenSpec validation, Go vet/tests, bounded Go lint, frontend lint/test/build, deployment checks, P92/P93 audit checks, release package smoke, whitespace checks, and security scans. Tag workflows package local deployment artifacts after preflight. These gates do not publish trading capability or change the product safety boundary.

## Governance

The authoritative docs live under `docs/`. Contract-level changes go through OpenSpec change packages and are merged back to `docs/` when archived. Read [docs/GOVERNANCE.md](docs/GOVERNANCE.md) and [openspec/project.md](openspec/project.md) before changing behavior, API contracts, schemas, workflow semantics, or frontend contracts.
