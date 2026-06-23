# P76 Post-P75 Package Refresh

> Date: 2026-06-21  
> Change: `p76-post-p75-final-package-refresh`  
> Result: `package_refresh_passed_for_release_ready_scoped_with_traceability_gaps`  
> Release status preserved: `release_ready_scoped_with_traceability_gaps`

## Scope

P76 refreshed the final local source handoff package after P75. It does not broaden P75 into `release_ready_full_requirements_traceable`; P75 remains a scoped release claim with explicit traceability gaps.

P76 did not add runtime product capability, SQLite schema, HTTP API, Eino workflow, frontend product behavior, provider calls, LLM calls, broker connectivity, trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic upgrade, automatic migration, automatic restore, automatic repair, real database overwrite, future provider availability, or investment-return claims.

During package repeat acceptance, the existing browser smoke exposed stale UI locators around the evidence summary and audit references after P73/P75 evidence additions. P76 fixed those Playwright assertions to target the visible P30 evidence row and P30 audit item directly. This is an acceptance-harness correction only; it does not relax the UI requirement and does not change runtime behavior.

## Package Identity

| Field | Value |
| --- | --- |
| Release label | `p76-post-p75-final` |
| Source commit | `8a317f25917b8ff18ec9b5049e6a6188206a22d3` |
| Source status | `clean` |
| Generated at | `20260621T030713Z` |
| Archive | `tmp/p76-final-release/20260621T030713Z/investment-agent-p76-post-p75-final.tar.gz` |
| Manifest | `tmp/p76-final-release/20260621T030713Z/release-manifest.json` |
| SHA-256 | `7540429d0b6c3cdd09dad2ebb10e2356580faf0b05e6acd92bc3bd9763a3dcb7` |
| Archive entries | `1417` |
| Archive size | `3.0M` |
| Verify summary | `tmp/p76-final-release/20260621T030723Z-verify/verify-summary.json` |
| Verify status | `passed` |
| Repeat summary | `tmp/p76-final-repeat/20260621T030727Z/repeat-summary.json` |
| Repeat status | `passed` |

The package was generated from a detached clean worktree. P76 release evidence itself is package-after-the-fact documentation and is not claimed to be inside this archive.

## Included Source Evidence

Direct package file-list checks confirmed the archive includes the committed P72-P75 acceptance Markdown and OpenSpec archives:

| Evidence | Package path |
| --- | --- |
| P72 acceptance | `docs/release/acceptance/2026-06-18-p72-real-user-fund-scenario.md` |
| P73 acceptance | `docs/release/acceptance/2026-06-19-p73-product-effectiveness-ux-validation.md` |
| P74 acceptance | `docs/release/acceptance/2026-06-19-p74-built-in-knowledge-and-data-readiness.md` |
| P75 real-use closure | `docs/release/acceptance/2026-06-20-p75-real-use-closure.md` |
| P75 traceability matrix | `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md` |
| P72 OpenSpec archive | `openspec/changes/archive/2026-06-19-p72-real-user-fund-scenario-data-impact-acceptance/` |
| P73 OpenSpec archive | `openspec/changes/archive/2026-06-19-p73-product-effectiveness-ux-validation/` |
| P74 OpenSpec archive | `openspec/changes/archive/2026-06-19-p74-built-in-knowledge-and-data-readiness/` |
| P75 OpenSpec archive | `openspec/changes/archive/2026-06-21-p75-requirements-traceability-and-real-use-closure/` |

`docs/release/ui-audit-assets/` remains excluded by the package safety contract. The archive includes `release-excluded-list.txt`, where the excluded screenshot/asset paths are recorded.

## Verification

| Gate | Command | Status | Evidence |
| --- | --- | --- | --- |
| Package generation | `bash scripts/local-release-package.sh --release-label p76-post-p75-final --output-dir tmp/p76-final-release` | passed | archive and manifest generated |
| Package verify | `bash scripts/local-release-package.sh --verify tmp/p76-final-release/20260621T030713Z/investment-agent-p76-post-p75-final.tar.gz --output-dir tmp/p76-final-release` | passed | checksum matched; required entries present; forbidden paths absent |
| Repeat OpenSpec | `openspec validate --all --strict` | passed | repeat log under `tmp/p76-final-repeat/20260621T030727Z/logs/` |
| Repeat Go tests | `go test ./...` | passed | duration 11s |
| Repeat npm install | `npm --prefix web ci` | passed | duration 2s; npm audit reported existing high-severity advisories |
| Repeat frontend tests | `npm --prefix web test` | passed | duration 5s |
| Repeat frontend build | `npm --prefix web run build` | passed | duration 3s |
| Repeat browser smoke | `env E2E_SERVER_PORT=18165 E2E_WEB_PORT=14265 bash scripts/e2e-smoke.sh` | passed | duration 27s |

Repeat acceptance was executed from the extracted package workspace, not from the active repository checkout.

## Known Caveats

- This is a local isolated, cross-machine-equivalent repeat, not a physical second-machine execution.
- Public-source and model-provider availability are not guaranteed by this package repeat.
- The package manifest's legacy `acceptance_references` field still lists older P63-era references, but the package file list and tar inspection confirm P72-P75 committed evidence is included under `docs/` and `openspec/`.
- P76 does not change P75 traceability counts: P75 remains 291 `partial`, 33 `scoped_pass`, 17 `deterministic_local_evidence`, and 0 `real_pass`.

## Not Claimed

P76 does not claim remote publishing, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic restore, automatic repair, real database overwrite, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, login-gated sources, paid sources, authorization-gated sources, Level2 data, high-frequency data, future provider availability, physical second-machine verification, future investment returns, or full original-requirement pass.
