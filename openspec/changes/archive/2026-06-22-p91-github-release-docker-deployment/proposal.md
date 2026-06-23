# P91 GitHub Release And Docker Deployment

## Why

P90 closed the known product acceptance blockers, but the project still lacks a GitHub-ready release and deployment layer. Users need a repeatable way to download a release package, provide local secrets such as `DEEPSEEK_API_KEY`, initialize local SQLite/VecLite data paths safely, start the application, upgrade without overwriting data, and uninstall while preserving data by default.

## What Changes

- Add Docker/Compose runtime packaging for the Go backend and built React frontend.
- Add `.env.example`, Docker config, entrypoint, healthcheck, and local data volume conventions.
- Add deployment scripts: `install.sh`, `upgrade.sh`, `uninstall.sh`, `backup.sh`, `status.sh`, and `doctor.sh`.
- Add GitHub CI and release workflows that run validation and package release artifacts without committing secrets.
- Add deployment documentation and P91 acceptance checks.

## Out Of Scope

- Physical second-machine validation.
- Broker integration, trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.
- Managed cloud hosting, domain/TLS automation, app-store packaging, notarization, or desktop app packaging.
- Automatic deletion or overwrite of real user data.

## Acceptance

P91 is acceptable only if:

- A fresh checkout can pass the P91 deployment check.
- Docker Compose deployment files exist and keep database, VecLite, backups, logs, and config in persistent local paths.
- `install.sh` detects first install vs upgrade and never deletes data.
- `upgrade.sh` backs up before changing containers.
- `uninstall.sh` preserves data by default and requires explicit confirmation for `--purge`.
- Runtime config accepts LLM key/base/model/timeout through environment variables.
- GitHub workflows run CI and create release artifacts without secrets in the package.
- Final validation includes OpenSpec, Go tests, frontend tests/build, P91 deployment checks, and `git diff --check`.
