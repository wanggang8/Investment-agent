# P90 Capital-Flow Provider Closure

## Result

- P90 evaluated the 2 full-release-required rows that remained after P89: `REQ-04-016` and `REQ-05-003`.
- P90 upgraded both rows to `real_pass` using fresh public-provider, real Settings UI, market snapshot API, and SQLite readback evidence.
- P90 conclusion: `release_ready_full_original_requirement_real_pass_candidate_with_p90_capital_flow_closure`.
- No P89-chain full-release-required row is known to remain non-`real_pass` after P90.
- P90 does not refresh the P76 package and does not claim physical second-machine validation, remote release, Git tag, broker integration, trading, external push, automatic confirmation, automatic rule application, paid/login/auth source, Level2 source, high-frequency source, or return guarantee.

## Evidence

- Acceptance summary: `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/p90-acceptance-summary.json`
- Matrix: `docs/release/acceptance/2026-06-22-p90-capital-flow-provider-matrix.md`
- Command: `bash scripts/p90-capital-flow-provider-acceptance.sh`
- Source preverification: `python3 scripts/p90_source_preverification.py --check`
- SQLite readback: `python3 scripts/p90_sqlite_readback_check.py <sqlite> <browser-results.json> <p90-source-preverification.json>`

## Verified Capital-Flow Fields

- Runtime snapshot: `market_c9077a3593f8f039f0adf7e174229ce510c4376572d90ef5d8b4efd7c3612e5c`
- `date`: `2026-06-22`
- `net_inflow`: `11895999.0`
- `net_outflow`: `0.0`
- `raw_net_flow`: `11895999.0`
- Directional mapping: positive raw value maps to `net_inflow`; negative raw value maps to `net_outflow`; raw value is preserved.

## Safety Boundary

- P90 only performs low-frequency read-only public data collection and local fact persistence.
- If the H5 provider becomes unavailable, the product must degrade/block dependent claims and must not synthesize capital-flow values.
