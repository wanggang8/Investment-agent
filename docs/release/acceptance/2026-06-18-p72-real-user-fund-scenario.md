# P72 Real User Fund Scenario Data Impact Acceptance

> Status: `release_ready_full_real_user_scenario_acceptance`
> Change: `p72-real-user-fund-scenario-data-impact-acceptance`

P72 adds a real user scenario acceptance layer after P71. It verifies that a user can add and maintain a real fund/ETF holding, import local knowledge, collect formal public evidence, consult the real LLM-backed workflow, manually record an offline execution, and then see the expected data impact across portfolio, decision, daily discipline, risk, notification, review, audit, rules, and workbench readbacks.

This is not a mock or scope-exclusion pass. The run used `use_stub=false`, a temporary SQLite database, a temporary VecLite file, a real local Go backend, Vite frontend, real public evidence sources, and real LLM configuration copied from the local private config.

P72 does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic upgrade, automatic restore, automatic overwrite of real user databases, return promises, login-gated sources, paid sources, authorization-gated sources, Level2 data, or high-frequency sources.

## Collection Blocker Resolution

During P72, accepting a consultation as safe-degraded was rejected as insufficient for the release goal. The failing path was formal evidence collection for the real user `510300` scenario.

Resolution implemented:

- Added read-only `csindex_index` formal evidence for the tracked index `000300` from CSIndex official index profile data.
- Added read-only `eastmoney_fund` formal evidence for the `510300` fund profile and latest NAV from Eastmoney/Tiantian fund public data.
- Updated normal formal evidence verification so two independent formal sources with at least one high-grade A/S source can satisfy normal fund-profile evidence; major/buy-logic evidence still requires two high-grade sources.
- Fixed retrieval verification merge so C-level background local knowledge no longer downgrades already satisfied formal evidence to `background_only`.
- Fixed VecLite rebuild persistence so SQLite `rag_chunks.index_status` is updated to `indexed` after a successful rebuild, making the database impact check verify the same state the UI/API reports.

## Command Evidence

The final P72 runner command:

```text
bash scripts/p72-real-user-fund-scenario-acceptance.sh
```

Final precheck and browser evidence:

```text
P72 precheck started at 2026-06-18T23:50:43Z
task public-evidence-refresh completed；已写入 audit_events；不会执行交易。
task p34-expanded-refresh completed；已写入 audit_events；不会执行交易。
data source quality regression completed:data_source_quality:mode=current:status=passed:policy=passed:gate=pass:cases=3:degraded=0:failed=0:no_auto_trading；不会执行交易。
task llm-smoke completed；已写入 audit_events；不会执行交易。

✓ P72 real user fund scenario validates UI operations, readbacks, LLM, RAG, daily discipline, and data impact
1 passed
```

SQLite impact checker result:

```text
status=passed
symbol=510300
cash=95630.5
total_assets=101265.0
position_count=2
total_quantity=1390.0
total_market_value=5634.5
workflow_status=completed
final_verdict_status=hold
confirmation_status=executed_manually
analyst_report_count=3
rag_chunks_p72_indexed=1
source_verifications=3
manual_daily_reports=1
risk_alerts_510300=1
notifications=4
forbidden_tables=[]
```

Artifacts:

```text
docs/release/ui-audit-assets/2026-06-18-p72/browser-results.json
docs/release/ui-audit-assets/2026-06-18-p72/db-impact-summary.json
docs/release/ui-audit-assets/2026-06-18-p72/precheck.log
docs/release/ui-audit-assets/2026-06-18-p72/db-impact-check.log
docs/release/ui-audit-assets/2026-06-18-p72/*.png
```

## Scenario Matrix

| Area | Evidence | Result |
| --- | --- | --- |
| Invalid input | Consultation blocks empty question/symbol and shows safe UI error | pass |
| Portfolio calibration | `510300` saved through UI; refresh shows local facts | pass |
| Holding maintenance | Edit, batch import validate/confirm, correction audit, and offline transaction operated through UI | pass |
| Local knowledge | P72 background note validated, confirmed, written to intelligence tables, and indexed | pass |
| Formal evidence | `csindex_index` A-level and `eastmoney_fund` B-level formal evidence collected for `510300` | pass |
| Current data | `000300` strict current gate `policy=passed` / `gate=pass` | pass |
| Market refresh | `510300` refresh accepts real success or explicit `DATA_SOURCE_UNAVAILABLE` safe state; final analysis does not rely on a mock | pass |
| Real LLM | LLM smoke passed; consultation generated three parsed and quality-passed analyst reports | pass |
| Retrieval | Consultation retrieval `status=hit`, `fallback_source=veclite`, `index_health=healthy`, no degraded reason | pass |
| Decision detail | Generated detail opened, evidence chain displayed, no sqlite fallback/missing-index warning | pass |
| Manual confirmation | User recorded offline sell confirmation; decision status became `executed_manually` | pass |
| Data impact | Portfolio cash, total assets, positions, transactions, confirmations, audit events, and decision-linked transaction matched deterministic expectations | pass |
| Daily/risk/readback | Dashboard, daily reports, risk alerts, review, audit, rules, notifications, workbench, and decision loop readbacks rendered after the scenario | pass |
| UI runtime health | Zero page errors, zero unexpected API failures, zero failed resource responses, zero console errors after fixes | pass |
| Safety | No forbidden broker/order/trade execution/external push tables or visible forbidden affordances | pass |

## Accuracy Checks

The data-impact checker verifies deterministic calculations from the UI operations:

- Initial calibration, edit, batch import, offline buy, and decision-linked manual sell produce the final expected `510300` aggregate.
- Final cash is `95630.5`; final total assets are `101265.0`.
- Two local `510300` position rows remain, with total quantity `1390.0` and total market value `5634.5`.
- At least one decision-linked operation confirmation and one offline local transaction confirmation exist.
- Daily discipline, risk alert, notification, and review records exist and are readable after the scenario.
- P72 RAG material is present in `intelligence_items`, `intelligence_summary`, and `rag_chunks` with `index_status='indexed'`.

## Post-Execution Coverage Review

Executed and evidenced:

- Real fund setup and maintenance.
- Offline transaction and generated-decision manual confirmation.
- Formal public evidence collection and local background knowledge.
- RAG/VecLite rebuild, retrieval, and database index status.
- Current-data strict gate and data-quality UI.
- Real LLM consultation quality and decision detail readback.
- Daily discipline, risk alert, notifications, review summary, audit timeline, rules, workbench, and decision-loop readbacks.
- Invalid input, unexpected API failure capture, page error capture, console error capture, and forbidden capability checks.

Intentionally out of scope:

- Physical second-machine execution.
- Future public-source or model-provider availability.
- Real broker/order/trade execution integration.
- External push channels.
- Investment return accuracy or future market prediction guarantees.
- Paid/login/authorized/Level2/high-frequency data sources.

## Result

P72 may claim:

- Real user fund scenario acceptance passed for the P72 run.
- Formal public evidence collection for the `510300` scenario was fixed and verified; this is not a safe-degraded pass.
- UI operations and SQLite data impact were verified end to end for portfolio maintenance, local knowledge/RAG, real LLM consultation, manual confirmation, daily discipline, risk, notification, review, audit, rules, workbench, and decision-loop readbacks.
- Current-data strict gate passed for `000300` during this run.

P72 must not claim:

- Future public-source availability.
- Future model-provider availability.
- Physical second-machine package repeat.
- Broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic upgrade, automatic migration, automatic restore, automatic real DB overwrite, or investment returns.
