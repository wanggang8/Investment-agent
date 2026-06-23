# P71 Real Product Acceptance True Pass

> Status: `release_ready_full_real_product_acceptance`
> Change: `p71-real-product-acceptance-true-pass`

P71 replaces the prior limited current-data release wording with a strict real-product acceptance run. Current data must pass without scope exclusion, VecLite must be healthy before real consultation, and full UI acceptance must run against a real local backend, Vite frontend, and real LLM provider.

P71 does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic upgrade, automatic restore, automatic overwrite of real user databases, return promises, login-gated sources, paid sources, authorization-gated sources, Level2 data, or high-frequency sources.

## Current Data Gate

The real local database was backed up under `tmp/p71-real-local-backup/` before writing new current source-health facts.

Command evidence:

```text
INVESTMENT_AGENT_CONFIG=tmp/p71-real-local/config.real-current.yaml go run ./cmd/agent --task p34-expanded-refresh --source csindex_extended --symbol 000300
task p34-expanded-refresh completed；已写入 audit_events；不会执行交易。

INVESTMENT_AGENT_CONFIG=tmp/p71-real-local/config.real-current.yaml go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate
data source quality regression completed:data_source_quality:mode=current:status=passed:policy=passed:gate=pass:cases=3:degraded=0:failed=0:no_auto_trading；不会执行交易。
```

Result: P66 current-data strict gate is a true pass for this run. This is not based on P67 scope exclusion.

The P71 browser summary confirms the same state through the API:

```text
current_data.status=passed
current_data.policy.verdict=passed
current_data.policy.release_gate=pass
current_data.missing_categories=[]
gate_resolution.release_claim_state=pass
gate_resolution.clean_data_claim_allowed=true
```

## VecLite Gate

P71 fixed `evidence-index` so it rebuilds the configured file VecLite index from SQLite RAG chunks and stamps rebuilt chunks with `indexed_at`, allowing retrieval freshness to be verified as `fresh`.

Command evidence:

```text
INVESTMENT_AGENT_CONFIG=tmp/p71-real-local/config.real-current.yaml go run ./cmd/agent --task evidence-index
task evidence-index completed；已写入 audit_events；不会执行交易。

INVESTMENT_AGENT_CONFIG=tmp/p71-real-local/config.real-current.yaml go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300
task retrieval-quality-smoke completed；已写入 audit_events；不会执行交易。

retrieval_quality:status=hit:topk=9:fallback=veclite:index=healthy:consistency=checked:no_auto_trading
```

The P71 real UI run rebuilt the index through `POST /api/v1/evidence/rebuild-index` and recorded:

```text
indexed_count=1
index_health.status=healthy
chunk_count=1
```

The generated decision detail recorded retrieval quality as `status=hit`, `fallback_source=veclite`, `index_health=healthy`, and `index_freshness=fresh`.

## Real UI Acceptance

P71 added `scripts/p71-real-product-acceptance.sh` and `web/e2e/p71-real-product-acceptance.spec.ts`.

The script uses a temporary SQLite database, real local backend, Vite frontend, `use_stub: false`, enabled `market_collectors/csindex`, real LLM config copied from `configs/config.local.yaml`, and screenshot artifacts under:

```text
docs/release/ui-audit-assets/2026-06-18-p71/
```

The browser run covers the P63 primary routes at 390px, 768px, and 1280px, operates local knowledge import, rebuilds VecLite through the API, calibrates the portfolio, submits real consultation, opens generated decision detail, checks data-quality current gate state, and treats page errors, failed unexpected API responses, retrieval degradation, missing LLM material, and forbidden trading affordances as blockers.

Strict command evidence after LLM recovery:

```text
bash scripts/p71-real-product-acceptance.sh
P71 precheck started at 2026-06-18T10:05:22Z
task p34-expanded-refresh completed；已写入 audit_events；不会执行交易。
data source quality regression completed:data_source_quality:mode=current:status=passed:policy=passed:gate=pass:cases=3:degraded=0:failed=0:no_auto_trading；不会执行交易。
task llm-smoke completed；已写入 audit_events；不会执行交易。

Running 1 test using 1 worker
  ✓  1 [chromium] › e2e/p71-real-product-acceptance.spec.ts:51:1 › P71 real product acceptance requires true current data pass, healthy VecLite, and strict real UI consultation (50.5s)

1 passed (51.1s)
```

Browser summary:

```text
decision_id=decision_a3aed494f6b84ac4
workflow_status=completed
analyst_report_count=3
parse_statuses=parsed,parsed,parsed
quality_statuses=passed,passed,passed
retrieval.status=hit
retrieval.fallback_source=veclite
retrieval.index_health=healthy
retrieval.index_freshness=fresh
llm_displayed=true
```

## Safety And Boundary Scan

P71 acceptance treats the following as blockers:

- Frontend mocks or fixture-only current data used as pass evidence.
- P67 scope exclusion used as P66 policy pass evidence.
- `VECTOR_INDEX_UNAVAILABLE`, stale/unknown VecLite freshness, or fallback to `sqlite_summary` during the real consultation result.
- Missing generated decision detail, missing LLM material, parse failures, quality failures, console errors, page errors, unexpected failed API responses, or mobile/desktop overflow.
- Forbidden trading, broker, order, push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, automatic overwrite, return promise, login-gated source, paid source, authorization-gated source, Level2, or high-frequency affordances.

The P71 run passed those blockers.

## Package Refresh

Post-P70/P71 package refresh passed from the accepted P71 strict-pass commit.

Package evidence:

```text
Archive: tmp/p71-final-release/20260618T101504Z/investment-agent-p71-real-product-acceptance.tar.gz
Manifest: tmp/p71-final-release/20260618T101504Z/release-manifest.json
SHA-256: fa6b857d96719600327e568c0b1a51a520f1d3838bfc516e28c6d77171026bfe
Source commit: 2c195a05cee3b6cdda031e86409d562bcc7ee379
Source status: clean
Verify summary: tmp/p71-final-release/20260618T101512Z-verify/verify-summary.json
Repeat summary: tmp/p71-final-repeat/20260618T101518Z/repeat-summary.json
```

Package verify result:

```text
status=passed
archive_entry_count=1350
errors=[]
warnings=[]
```

Repeat acceptance result:

```text
status=passed
openspec validate --all --strict: passed
go test ./...: passed
npm --prefix web ci: passed
npm --prefix web test: passed
npm --prefix web run build: passed
bash scripts/e2e-smoke.sh: passed
```

The archive includes P69, P70, and P71 acceptance materials that were present in source commit `2c195a05cee3b6cdda031e86409d562bcc7ee379`. It does not include this package-after-the-fact evidence section unless a later package refresh is generated from a later commit.

## Result

P71 may claim:

- Full real product acceptance passed for the P71 run.
- Current-data true pass for `000300`: `policy=passed` / `gate=pass`.
- VecLite acceptance hardening passed with a healthy, fresh file index during real UI consultation.
- All P71 strict real UI operations passed against a real local backend, Vite frontend, and real LLM provider.
- Post-P70/P71 package refresh, package verify, and isolated repeat acceptance passed from clean source commit `2c195a05cee3b6cdda031e86409d562bcc7ee379`.

P71 must not claim:

- Future public-source availability.
- Future model-provider availability.
- Physical second-machine package repeat unless separately performed.
- This package-after-the-fact evidence section is inside the package archive.
- Broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic upgrade, automatic migration, automatic restore, automatic real DB overwrite, or investment returns.
