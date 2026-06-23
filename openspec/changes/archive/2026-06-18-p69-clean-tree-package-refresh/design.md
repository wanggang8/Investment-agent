# Design: P69 Clean Tree Package Refresh

## Context

P64 introduced the package workflow and P65 repeated acceptance from a candidate archive. Those artifacts proved the workflow but were generated from active implementation states. P68 changed the release claim to `release_ready_limited_current_data_scope` and recommended a clean-tree package refresh.

P69 should turn that recommendation into package evidence without changing product behavior. The central design constraint is source cleanliness: once P69 files exist in the main checkout, `scripts/local-release-package.sh` would correctly report `source_status=dirty`. To produce true `source_status=clean` evidence, P69 uses a temporary detached worktree at the committed P68 HEAD.

## Source Snapshot

Use current `HEAD` before P69 edits as the package source commit. At P69 start this should be the P68 commit. Record it explicitly and use the recorded value for later commands:

```bash
P69_SOURCE_COMMIT="$(git rev-parse HEAD)"
git status --short
```

Create a temporary detached worktree under project `tmp/`, for example:

```bash
git worktree add --detach tmp/p69-clean-tree-source "$P69_SOURCE_COMMIT"
```

Run package and repeat commands inside that detached worktree. The detached worktree must remain clean before packaging, and all package/repeat outputs must stay under that worktree's `tmp/` directory or the main project `tmp/` directory. `tmp/` is excluded from the package and from committed outputs.

The package script runs `npm --prefix web run build` before creating the archive, so the detached worktree needs frontend dependencies first:

```bash
(cd tmp/p69-clean-tree-source && npm --prefix web ci)
(cd tmp/p69-clean-tree-source && git status --short)
```

`web/node_modules/` is ignored and excluded from packages; the second command must still be empty.

## Package Commands

From the clean detached worktree:

```bash
npm --prefix web ci
bash scripts/local-release-package.sh --release-label p69-clean-tree --output-dir tmp/p69-release
bash scripts/local-release-package.sh --verify tmp/p69-release/<timestamp>/investment-agent-p69-clean-tree.tar.gz --output-dir tmp/p69-release
bash scripts/local-release-repeat-acceptance.sh --archive tmp/p69-release/<timestamp>/investment-agent-p69-clean-tree.tar.gz --output-dir tmp/p69-repeat
```

The package manifest must show:

- source commit equal to the P68 commit;
- `source_status=clean`;
- release label `p69-clean-tree`;
- archive checksum;
- archive entry count;
- required roots included;
- forbidden paths absent.

## Documentation

P69 should create:

- `docs/release/acceptance/2026-06-18-p69-clean-tree-package-refresh.md`

P69 should update:

- `docs/release/release-packaging-2026-06-18.md`
- `docs/release/acceptance/2026-06-18-p65-cross-machine-repeat.md` only if needed to point to superseding P69 evidence
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

P69 may claim:

- A clean-tree package was generated from the P68 source commit.
- The package archive and manifest verified.
- Package repeat acceptance passed from the generated archive.

P69 must not claim:

- The package includes P69 documentation unless a follow-up post-P69 package refresh is performed.
- P66 current-data policy passed.
- A physical second machine ran the archive.
- Remote publication, Git tagging, installer signing, automatic upgrade/migration/restore/repair, trading, external push, provider availability, or investment returns.

## Cleanup

Temporary worktrees and generated artifacts remain under `tmp/` and must not be committed. If a temporary git worktree is created, remove it with:

```bash
git worktree remove tmp/p69-clean-tree-source
```

or document why it remains if removal is impossible.
