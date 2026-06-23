# P98 Runtime Hardening And Code Reuse Cleanup Acceptance

> Date: 2026-06-23  
> Change: `p98-runtime-hardening-and-code-reuse-cleanup`

## Result

P98 passed. The change adds a release runtime-mode guardrail and consolidates frontend diagnostic redaction without adding investment runtime capabilities.

## Implemented Scope

- Added `runtime.mode` configuration with `development`, `test`, and `release` modes.
- Added validation that rejects `runtime.mode=release` when `data_sources.use_stub=true`.
- Set Docker release config and `.env.example` to release mode while keeping example/local development config stub-friendly.
- Added shared frontend `redactSensitiveText` utility.
- Replaced duplicated redaction logic in `ErrorState`, `LocalInstallPage`, and `DataQualityPage`.
- Refreshed the P93 code-reality report after P98 source/config changes.

## Verification

| check | result |
| --- | --- |
| `openspec validate --all --strict` | passed: 35 items, 0 failed |
| `go test $(bash scripts/go-packages.sh)` | passed |
| `npm --prefix web test` | passed: 49 test files, 178 tests |
| `npm --prefix web run build` | passed |
| `python3 scripts/api_route_contract_check.py` | passed: 57 routes |
| `python3 scripts/p91_deployment_check.py --check` | passed |
| `python3 scripts/p92_final_requirement_audit.py --check` | passed |
| `python3 scripts/p93_code_reality_audit.py --check` | passed |

## Boundary

P98 does not add broker connectivity, trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, return guarantees, paid/login/authorization-only sources, Level2 data, high-frequency data, or new investment strategy behavior.
