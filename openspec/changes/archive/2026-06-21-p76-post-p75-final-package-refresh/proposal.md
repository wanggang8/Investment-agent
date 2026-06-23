# P76 Post-P75 Final Package Refresh

## Why

P75 completed and archived the original-requirement traceability and real-use closure work, but the latest local distribution package evidence still predates P72-P75 acceptance records and source changes. A final handoff package should be regenerated from the clean post-P75 commit so the archive contains the P72, P73, P74, and P75 acceptance documents and OpenSpec archives that are now part of the committed source.

P76 is a packaging and release-evidence refresh only. It does not broaden the P75 conclusion: the current release remains `release_ready_scoped_with_traceability_gaps`, not a full original-requirement pass.

## What Changes

- Record the post-P75 clean source commit.
- Correct stale Playwright smoke locators if package repeat acceptance exposes acceptance-harness drift after P73/P75 evidence additions.
- Generate a local release package from that clean source commit using the existing package workflow.
- Verify the package archive and adjacent manifest.
- Run cross-machine-equivalent local repeat acceptance from the extracted package workspace.
- Update release packaging, handoff, repeatability, release index, governance, and progress materials with the package identity, checksum, source commit, verify result, repeat result, and Not Claimed boundaries.

## In Scope

- OpenSpec change files, release docs, package/repeat acceptance record, governance/progress docs, and acceptance-harness-only locator corrections required for repeat acceptance.
- Local generated package/repeat artifacts under `tmp/` only.
- Existing package and repeat scripts.

## Out of Scope

- No runtime feature work, SQLite schema changes, HTTP API changes, Eino workflow changes, frontend UI changes, provider calls, LLM calls, source refresh, data repair, migrations, restores, remote publication, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic repair, broker interface, trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, return promise, login-gated source, paid source, authorization-gated source, Level2 data, or high-frequency source.
- No claim that the package includes screenshots or UI audit asset directories, because the established P64 package safety contract excludes `docs/release/ui-audit-assets/`.
- No claim that a separate physical second machine ran the package.
- No claim that P75 became `release_ready_full_requirements_traceable`.

## Impact

Docs/OpenSpec and local `tmp/` package artifacts only. The committed package evidence should identify the generated archive and repeat result; the archive itself, manifest sidecar, logs, extracted workspace, node modules, build output, SQLite databases, and traces remain local artifacts under `tmp/` and are not committed.
