# P95 Architecture API Engineering Hardening Acceptance

> Date: 2026-06-23

## Scope

P95 hardened public engineering validation, API route governance, SQLite local runtime settings, Docker secret handling, and architecture documentation. It did not add investment runtime capability, broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.

## Implemented Changes

- Added `scripts/go-packages.sh` and wired CI/release Go vet/test gates to project-owned backend packages instead of raw `./...`.
- Updated `golangci-lint` package scope to `./cmd/... ./internal/... ./pkg/...`.
- Made P93 code reality audit build its production scan inventory from Git tracked plus nonignored untracked release-relevant files, so new source files are scanned while ignored runtime artifacts such as `cmd/agent/tmp/veclite` do not affect the report.
- Added `scripts/api_route_contract_check.py` and wired it into CI/release preflight.
- Documented missing API routes:
  - `GET /api/v1/daily-auto-run/status`
  - `POST /api/v1/portfolio/rebalance-review`
  - `POST /api/v1/rule-proposals/sop-addendum`
- Added SQLite connection-level configuration hook for foreign keys, busy timeout, and file-backed WAL so every newly opened pooled SQLite connection is configured.
- Added `DEEPSEEK_API_KEY_FILE` support, Compose environment pass-through, entrypoint warning alignment, `.env.example` documentation, and deployment documentation.
- Updated `docs/architecture.md` to reflect the current directory layout and engineering gates.
- Updated release package manifest verification commands to use `scripts/go-packages.sh`.

## Validation

Passed in this environment:

- `openspec validate --all --strict`: passed.
- `bash scripts/go-packages.sh`: selected only `cmd`, `internal`, and `pkg` packages; no `web/node_modules` or `investment-agent/web` package selected.
- `go vet $(bash scripts/go-packages.sh)`: passed.
- `go test ./internal/infrastructure/config ./internal/infrastructure/persistence/sqlite ./internal/application/handler ./internal/domain/rule ./pkg/httputil`: passed.
- `go test ./internal/infrastructure/persistence/sqlite`: passed with pooled multi-connection PRAGMA coverage.
- `go test $(bash scripts/go-packages.sh) -run '^$'`: passed, full backend package test compile gate.
- `go test $(bash scripts/go-packages.sh)`: passed after unrestricted local permissions were restored.
- `npm --prefix web run lint`: passed.
- `npm --prefix web test -- --run`: 48 files passed, 176 tests passed.
- `npm --prefix web run build`: passed.
- `python3 scripts/p91_deployment_check.py`: passed.
- `python3 scripts/p92_final_requirement_audit.py --check`: passed.
- `python3 scripts/p93_code_reality_audit.py --check`: passed.
- `python3 scripts/api_route_contract_check.py`: passed, 57 routes.
- `bash scripts/local-release-package.sh --release-label p95-smoke-fix --output-dir tmp/p95-release-fix`: package created.
- `bash scripts/local-release-package.sh --verify tmp/p95-release-fix/20260623T030924Z/investment-agent-p95-smoke-fix.tar.gz --output-dir tmp/p95-release-fix`: passed.
- `release-manifest.json` verification commands include the copyable command `go test $(bash scripts/go-packages.sh)` without an escaped dollar sign.
- `git diff --check`: passed.

Passed in GitHub Actions on commit `054d2708440da925e2e4cc4ae65065ac002b8905`:

- CI run: `https://github.com/wanggang8/Investment-agent/actions/runs/28001447118`
- CI conclusion: `success`, completed `2026-06-23T04:09:08Z`.
- CI `Go tests` step: `success`, completed `2026-06-23T04:08:13Z`.
- CI also passed OpenSpec validation, Go vet, Go lint, frontend lint/tests/build, P91/P92/P93 checks, API route contract check, release package smoke, and whitespace check.
- Security Scan run: `https://github.com/wanggang8/Investment-agent/actions/runs/28001447137`
- Security Scan conclusion: `success`, completed `2026-06-23T04:07:30Z`; included Go vulnerability scan, frontend production dependency audit, and P93 code reality / secret scan.

## Review Closure

Subagent review identified four blockers before archive: SQLite PRAGMAs were pool-scoped instead of per-connection, the release manifest escaped the backend test command, `docs/architecture.md` described P93 as tracked-only while the script scans tracked plus nonignored untracked source files, and full backend execution evidence was unavailable in this sandbox. The first three blockers are closed by the P95 patch and the validation above.

## Environment-Limited Checks

Before GitHub CI evidence was available, `go test $(bash scripts/go-packages.sh)` was attempted in the earlier local sandbox and was blocked by local listening-socket restrictions used by `httptest.NewServer`, causing tests in `cmd/agent` and `internal/application/workflow` to panic with `listen tcp6 [::1]:0: bind: operation not permitted`. This was an environment permission limit, not a test assertion failure. The archive blocker is now closed by the successful GitHub CI full backend test evidence above.

## Boundary

P95 is an engineering hardening stage. It does not change the release claim into a new product capability claim and does not publish a tag or GitHub Release.
