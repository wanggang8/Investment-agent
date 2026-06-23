# Design: P76 Post-P75 Final Package Refresh

## Source Snapshot

P76 uses a clean post-P75 package source commit. The repeat package run exposed stale Playwright smoke locators in the evidence and audit UI checks after P73/P75 added additional evidence rows, so P76 first committed an acceptance-harness locator correction:

- Source commit: `8a317f25917b8ff18ec9b5049e6a6188206a22d3`
- Change: target the visible P30 evidence row and P30 audit item directly in `web/e2e/local-smoke.spec.ts`
- Runtime impact: none; no product behavior, API, database, workflow, provider, or LLM changes

The source commit is recorded before package generation:

```bash
P76_SOURCE_COMMIT="$(git rev-parse HEAD)"
git status --short
```

The package source commit is expected to contain the committed P72-P75 acceptance records, P75 archive, and the acceptance-harness locator correction. P76 documentation itself is package-after-the-fact evidence and is not claimed to be included in the generated archive.

## Package Flow

Use the existing local package workflow from the clean source commit:

```bash
bash scripts/local-release-package.sh --release-label p76-post-p75-final --output-dir tmp/p76-final-release
bash scripts/local-release-package.sh --verify tmp/p76-final-release/<timestamp>/investment-agent-p76-post-p75-final.tar.gz --output-dir tmp/p76-final-release
bash scripts/local-release-repeat-acceptance.sh --archive tmp/p76-final-release/<timestamp>/investment-agent-p76-post-p75-final.tar.gz --output-dir tmp/p76-final-repeat
```

The package manifest must show:

- source commit equal to the recorded P76 package source commit;
- `source_status=clean`;
- release label `p76-post-p75-final`;
- archive checksum;
- required source roots included;
- forbidden paths absent.

The package script intentionally excludes `docs/release/ui-audit-assets/`. P76 must state that the package includes acceptance Markdown and OpenSpec records, but not screenshot/asset directories.

## Repeat Acceptance

Repeat acceptance runs from the extracted archive workspace and covers:

- package verify;
- `openspec validate --all --strict`;
- `go test ./...`;
- `npm --prefix web ci`;
- `npm --prefix web test`;
- `npm --prefix web run build`;
- local E2E smoke unless an explicit blocker is recorded.

## Documentation

P76 creates:

- `docs/release/acceptance/2026-06-21-p76-post-p75-package-refresh.md`

P76 updates:

- `docs/release/release-packaging-2026-06-18.md`
- `docs/release/release-handoff-2026-06-18.md`
- `docs/release/README.md`
- `docs/release/acceptance-repeatability.md`
- `docs/development-plan.md`
- `docs/README.md`
- `docs/GOVERNANCE.md`
- `AGENTS.md`
- `openspec/project.md`
- `openspec/PROGRESS.md`

## Claim Boundaries

P76 may claim:

- The final local source package was regenerated from the clean post-P75 commit.
- The final local source package includes the acceptance-harness locator correction required for repeat acceptance.
- The package includes committed P72-P75 acceptance Markdown, OpenSpec archives, runtime source, tests, and package scripts.
- Package verify and cross-machine-equivalent local repeat acceptance passed.

P76 must not claim:

- The archive includes P76 package-after-the-fact evidence.
- The archive includes `docs/release/ui-audit-assets/` screenshots/assets.
- P75 achieved full original-requirement pass.
- Physical second-machine execution, remote publication, Git tag creation, installer signing, automatic upgrade/migration/restore/repair, real database overwrite, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, future provider availability, or investment returns.
