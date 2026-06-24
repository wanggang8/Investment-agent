# Design

## Acceptance Model

P100 is a local-source final acceptance run. It validates the product by starting from the repository source tree and using local Go, React/Vite, SQLite, VecLite, real browser UI journeys, API responses, SQLite/readback evidence, and audit records.

P100 intentionally excludes deployment packaging. Docker, install/upgrade/uninstall scripts, GitHub Release, package refresh, and physical second-machine repeat acceptance remain out of scope.

## Evidence Layers

The acceptance must combine four evidence layers:

1. Governance and contract evidence: OpenSpec validation, no active unrelated changes, P92 final requirement ledger, and P93 code-reality/design audit.
2. Automated local runtime evidence: Go tests, vet, frontend tests/build, E2E smoke, and P71-P90 local acceptance runners where they are source-runtime based.
3. Real browser product evidence: manual or scripted browser journeys across workbench, positions, settings/data refresh, consultation, decision detail, review, rules, audit, notifications, and data-quality surfaces.
4. Data impact evidence: API response summaries, SQLite readback summaries, and `audit_events` or acceptance JSON proving that UI actions have the intended local data effect.

## Local Runtime Configuration

The run should use an ignored local config file:

```bash
cp configs/config.example.yaml configs/config.yaml
```

The local config used for acceptance must keep release-safety intent:

- `data_sources.use_stub: false` for real-source/provider acceptance claims.
- SQLite and VecLite paths point to an acceptance-specific temporary or disposable local directory unless the task explicitly performs read-only checks.
- LLM-backed claims require a configured test key through environment or local ignored config. Without a valid key, P100 may only claim safe LLM degradation for that path.

## Product Design Rubric

Product design passes only if the real browser journey shows:

- The first screen and each major route make the next user action clear.
- Loading, empty, degraded, error, warning, and success states are visible and recoverable.
- Investment safety boundaries are not hidden or contradicted by UI wording.
- Core conclusions link to evidence, rules, assumptions, source health, or audit records.
- 390px, 768px, and 1280px widths do not show incoherent overlap, clipped controls, or unusable navigation.
- Decision confirmation remains explicit and manual; no UI suggests automatic trading or automatic rule application.

## Blocking Criteria

The P100 run is blocked if any of the following occurs:

- P92 check reports any full-release-required row as non-`real_pass`.
- P93 check reports active release-blocking findings.
- Go tests, Go vet, frontend tests, frontend build, or OpenSpec validation fail without a documented non-release-impact reason.
- A core local source runtime journey cannot be completed from the browser.
- UI action, API readback, SQLite/readback, and audit evidence disagree for a critical path.
- The UI implies forbidden capabilities or hides mandatory manual confirmation.
- Any checked artifact leaks a full API key, private key, raw prompt, raw SQL, raw HTTP payload, private absolute path, broker credential, or supplier raw response.

## Final Acceptance Record

The final record must be created at:

```text
docs/release/acceptance/2026-06-23-p100-local-source-final-acceptance.md
```

The conclusion should be one of:

- `local_source_release_acceptance_passed`
- `local_source_release_acceptance_passed_with_documented_degradation`
- `local_source_release_acceptance_blocked`

The record must explicitly state that Docker, installation scripts, GitHub Release, package refresh, and physical second-machine validation were not part of P100.
