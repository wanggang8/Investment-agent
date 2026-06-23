# P67 Current Data Gate Resolution Acceptance

> Date: 2026-06-18
> Change: `p67-current-data-gate-resolution-workflow`
> Status: `resolved_with_scope_exclusion`

## Scope

P67 adds a local manual resolution workflow for the P66 current data policy gate. It records a local `scope_exclusion` or `waiver`, exposes the release claim state through API/UI/CLI, and keeps P66 policy semantics unchanged.

## Evidence

| Check | Result | Evidence |
| --- | --- | --- |
| P66 strict gate before/after P67 | blocked | `policy=blocked:gate=block:cases=1:degraded=1:failed=0:no_auto_trading` |
| P67 initial resolution check | requires resolution | `claim_state=requires_resolution:policy=blocked:gate=block:resolution=none:clean_data_claim=false:no_auto_trading` |
| Real UI operation | pass | Started local server and Vite app, opened `/data-quality`, entered symbol `000300`, clicked `检查门禁处置`, clicked `记录处置`, and observed `已排除 current data clean claim`, `clean data claim：不允许`, `范围排除 · active`. Screenshot: `tmp/p67-data-quality-resolution-ui-000300.png`. |
| P67 CLI after UI record | pass | `claim_state=resolved_with_scope_exclusion:policy=blocked:gate=block:fingerprint=e63f5ffed4a4307a4f791a4d68a9fc4a6f37bd34188dcacea1fd3d1f9ac4da48:resolution=scope_exclusion:clean_data_claim=false:no_auto_trading` |
| P66 strict gate after resolution | still blocked | P67 did not convert policy to pass. |

The first parallel rerun of P66/P67 commands hit a SQLite busy lock. The commands were rerun sequentially; the sequential evidence above is authoritative.

## Allowed Claim

- It is acceptable to state that current local data health has been excluded from the clean release claim through a local `scope_exclusion`.

## Not Claimed

- Current local data is not clean.
- Current local data is not healthy.
- P66 `policy=blocked` is not passed.
- P67 did not refresh data, repair providers, call external sources, call LLM providers, create broker connectivity, trade, push externally, auto-confirm, auto-apply rules, or promise future provider availability or returns.

## Verification Commands

```bash
go test ./internal/infrastructure/persistence/sqlite ./internal/application/service ./internal/application/handler ./cmd/agent
go test ./...
npm --prefix web test
npm --prefix web run build
bash scripts/e2e-smoke.sh
go run ./cmd/agent --task data-source-quality-resolution-check --symbol 000300
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate
openspec validate p67-current-data-gate-resolution-workflow --strict
openspec validate --all --strict
git diff --check
```
