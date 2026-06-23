# Tasks: P76 Post-P75 Final Package Refresh

## 1. Setup And Scope

- [x] 1.1 Confirm `docs/GOVERNANCE.md`, `openspec/project.md`, and `openspec/PROGRESS.md` show no active change after P75.
- [x] 1.2 Confirm P72-P75 acceptance and archive materials are committed in the post-P75 source commit.
- [x] 1.3 Create OpenSpec change `p76-post-p75-final-package-refresh`.
- [x] 1.4 State P76 is package/release-evidence only and does not broaden P75 beyond `release_ready_scoped_with_traceability_gaps`.
- [x] 1.5 Mark P76 active in governance/progress materials.
- [x] 1.6 Run `openspec validate p76-post-p75-final-package-refresh --strict`, `openspec validate --all --strict`, and `git diff --check`.

## 2. Package Generation

- [x] 2.1 Record post-P75 source commit.
- [x] 2.2 Confirm source tree is clean before P76 package generation.
- [x] 2.3 Run `bash scripts/local-release-package.sh --release-label p76-post-p75-final --output-dir tmp/p76-final-release`.
- [x] 2.4 Parse archive path, manifest path, SHA-256, source commit, source status, archive size, and entry count.
- [x] 2.5 Run package verify with `bash scripts/local-release-package.sh --verify <archive> --output-dir tmp/p76-final-release`.
- [x] 2.6 Confirm package includes committed P72-P75 acceptance Markdown and OpenSpec archives.
- [x] 2.7 Confirm package excludes `docs/release/ui-audit-assets/` under the established package safety contract.

## 3. Package Repeat Acceptance

- [x] 3.1 Run `bash scripts/local-release-repeat-acceptance.sh --archive <archive> --output-dir tmp/p76-final-repeat`.
- [x] 3.2 Record repeat summary: commands, durations, status, source commit, source status, skip flags, and caveats.
- [x] 3.3 Confirm repeat command matrix covers OpenSpec validation, Go tests, npm ci, frontend tests, frontend build, and E2E smoke.
- [x] 3.4 Confirm repeat artifacts remain under `tmp/` and no package archive, manifest sidecar, logs, extracted workspace, node_modules, dist, SQLite DB, or traces are committed.

## 4. Release Materials

- [x] 4.1 Add `docs/release/acceptance/2026-06-21-p76-post-p75-package-refresh.md`.
- [x] 4.2 Update `docs/release/release-packaging-2026-06-18.md` with P76 package identity and freshness.
- [x] 4.3 Update `docs/release/release-handoff-2026-06-18.md` with P76 package handoff status and package boundaries.
- [x] 4.4 Update `docs/release/README.md` and `docs/release/acceptance-repeatability.md`.
- [x] 4.5 Update `docs/development-plan.md`, `docs/README.md`, `docs/GOVERNANCE.md`, `AGENTS.md`, `openspec/project.md`, and `openspec/PROGRESS.md`.

## 5. Verification And Review

- [x] 5.1 Run `openspec validate p76-post-p75-final-package-refresh --strict`.
- [x] 5.2 Run `openspec validate --all --strict`.
- [x] 5.3 Run `git diff --check`.
- [x] 5.4 Run release wording scan for stale package claims and overbroad full-pass claims.
- [x] 5.5 Run forbidden capability scan for broker/trading/push/auto-confirm/auto-rule/auto-repair/return/provider claims.
- [x] 5.6 Run subagent review; Critical/Important findings must be fixed before archive.

## 6. Archive And Completion

- [x] 6.1 Archive P76 and merge release-governance delta.
- [x] 6.2 Confirm no active change remains after archive.
- [x] 6.3 Run final `openspec validate --all --strict` and `git diff --check`.
- [x] 6.4 Commit P76 release package refresh materials.
