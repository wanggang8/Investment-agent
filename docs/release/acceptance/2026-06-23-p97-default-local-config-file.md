# P97 Default Local Config File Acceptance

> Date: 2026-06-23

## Scope

P97 fixes local startup configuration semantics. `configs/config.yaml` is now the preferred default local runtime config, while `configs/config.example.yaml` remains the committed template and fresh-checkout fallback.

P97 does not add investment runtime capability, broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.

## Implemented Changes

- Updated `config.Load("")` default resolution order to `INVESTMENT_AGENT_CONFIG` -> `configs/config.yaml` -> `configs/config.example.yaml`.
- Updated `cmd/agent` config path display/help to match the shared loader.
- Added tests for default local config, example fallback, and env override precedence.
- Added `configs/config.yaml` to `.gitignore`.
- Updated local startup/config docs to instruct users to copy `configs/config.example.yaml` to `configs/config.yaml`.
- Updated local diagnostic script defaults to prefer `configs/config.yaml` when present.

## Validation

- `go test ./internal/infrastructure/config`: passed.
- `go test $(bash scripts/go-packages.sh)`: passed.
- `openspec validate --all --strict`: passed.
- `python3 scripts/p93_code_reality_audit.py --check`: passed after refreshing the P93 source inventory report for P97 files.
- `git check-ignore -q configs/config.yaml`: passed.
- `git diff --check`: passed.
