# P87 Portfolio State Allocation Safety Closure

- Generated at: `2026-06-22T05:50:40Z`
- Status: `passed`
- Source matrix: `docs/release/acceptance/2026-06-22-p85-expected-return-analysis-accuracy-matrix.md`
- Output matrix: `docs/release/acceptance/2026-06-22-p87-portfolio-state-allocation-safety-matrix.md`
- Summary artifact: `docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety/portfolio-state-allocation-summary.json`
- Browser status: `passed`
- SQLite status: `passed`

## Evidence

- Command: `P87_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety bash scripts/p87-portfolio-state-allocation-acceptance.sh`
- Browser results: `docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety/browser-results.json`
- SQLite readback: `docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety/db-readback-check.log`
- Screenshots: `docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety/p87-*.png`
- Scenarios: core/satellite/cash portfolio UI write/readback, buy-date/state persistence, sell-only decision, frozen-watch decision, information-insufficient decision.

## Row Outcome

- Total rows: `341`
- Counts: `{'partial': 137, 'real_pass': 193, 'reference_only': 11}`
- P87 planned rows: `32`
- P87 upgraded rows: `5`
- Full-release-required rows still non-real-pass: `137`
- Upgraded: `REQ-02-031, REQ-10-001, REQ-10-002, REQ-10-003, REQ-11-005`
- Deferred: `REQ-01-001, REQ-01-006, REQ-02-006, REQ-02-022, REQ-02-024, REQ-02-025, REQ-03-004, REQ-03-005, REQ-03-006, REQ-04-003, REQ-04-008, REQ-04-016, REQ-04-025, REQ-05-010, REQ-06-023, REQ-06-024, REQ-07-006, REQ-07-015, REQ-08-018, REQ-08-020, REQ-10-004, REQ-14-005, REQ-14-007, REQ-16-028, REQ-16-033, REQ-17-015, REQ-17-024`

## Boundary

- P87 does not claim complete quarterly rebalance automation, monthly attribution, full audit-history breadth, proposal accept/reject closure, public collector production readiness, broker connectivity, automatic trading, automatic confirmation, external push, release/upgrade preflight full closure, or full original-requirement pass.
- P87 treats user operation as local fact recording only. Sell-only, frozen-watch, and information-insufficient evidence does not authorize automatic trade execution.
