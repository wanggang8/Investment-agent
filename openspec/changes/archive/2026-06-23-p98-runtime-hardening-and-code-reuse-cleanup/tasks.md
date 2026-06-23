# P98 Tasks

## 1. Scope And Change Setup

- [x] Confirm no active OpenSpec change conflicts with P98.
- [x] Create P98 proposal, design, tasks, and release-governance delta.
- [x] Validate the active P98 change.

## 2. Runtime Mode Guardrail

- [x] Add failing config tests for `runtime.mode=release` with `data_sources.use_stub=true` and `runtime.mode=release` with real collectors.
- [x] Add runtime mode config parsing/defaulting/validation.
- [x] Update Docker config and `.env.example` to set release mode.
- [x] Update configuration/deployment docs for runtime mode.

## 3. Shared Frontend Redaction

- [x] Add failing tests for a shared frontend redaction utility.
- [x] Implement `web/src/shared/utils/redaction.ts` and export it.
- [x] Replace duplicated redaction logic in `ErrorState`, `LocalInstallPage`, and `DataQualityPage`.
- [x] Keep existing page/component redaction tests passing.

## 4. Acceptance

- [x] Run `openspec validate --all --strict`.
- [x] Run `go test $(bash scripts/go-packages.sh)`.
- [x] Run `npm --prefix web test`.
- [x] Run `npm --prefix web run build`.
- [x] Run `python3 scripts/api_route_contract_check.py`.
- [x] Run `python3 scripts/p91_deployment_check.py --check`.
- [x] Run `python3 scripts/p92_final_requirement_audit.py --check`.
- [x] Run `python3 scripts/p93_code_reality_audit.py --check`.
- [x] Generate P98 acceptance record.
- [x] Archive P98 only after validation passes.
