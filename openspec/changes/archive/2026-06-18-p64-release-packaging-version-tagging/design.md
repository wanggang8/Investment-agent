# Design: P64 Release Packaging And Version Tagging

## Design Brief

Investment Agent now has a P63 `release_ready` handoff with full UI regression evidence. P64 should make that handoff concrete for local delivery: a reproducible package artifact, a manifest describing exactly what was packaged, and a verification path that another operator can run before trusting the package.

The design is intentionally local and conservative. It does not create remote distribution, auto-update, migration execution, or any trading-related behavior.

## Approach

Use a shell script as the packaging entrypoint, following existing local operations scripts:

- `scripts/local-install-diagnostics.sh`
- `scripts/local-release-upgrade-check.sh`
- `scripts/p63-full-ui-regression.sh`

The package script should run from the repository root, stage files under `tmp/local-release-package/<timestamp>/`, write a sanitized manifest, create an archive, and optionally verify an existing archive. Keeping this in shell matches the current operations surface and avoids adding a new runtime dependency.

## Package Contents

The package should include tracked source and documentation required to build, inspect, and run locally:

- `cmd/`
- `internal/`
- `pkg/`
- `web/` source files, package metadata, Playwright/Vitest config, and public assets
- `configs/config.example.yaml`
- `docs/`
- `openspec/`
- `scripts/`
- `examples/`
- root metadata such as `go.mod`, `go.sum`, `AGENTS.md`, `README` files if present

The package should exclude:

- `.git/`
- `tmp/`
- `configs/config.local.yaml`
- `web/node_modules/`
- `web/dist/`
- Playwright reports, traces, screenshots outside committed docs release assets
- SQLite databases, VecLite local indexes, logs, raw vendor payloads, complete prompts, complete keys, private paths

## Manifest Shape

`release-manifest.json` should contain:

- `release_label`
- `commit`
- `generated_at`
- `package_archive`
- `package_sha256`
- `source_status`
- `included_roots`
- `excluded_patterns`
- `verification_commands`
- `acceptance_references`
- `known_degradations`
- `not_claimed`
- `safety_note`

The manifest must use relative package paths and sanitized placeholders. It must not include local absolute paths except placeholders such as `<repo>` or `<package>`.

## Verification

Verification should inspect the generated archive and manifest without executing packaged runtime behavior. It should check:

- manifest exists;
- manifest JSON parses;
- package checksum matches;
- expected top-level paths exist;
- forbidden paths and file patterns are absent;
- no complete key pattern or known private path pattern appears in manifest text or archive listing.

P64 should continue to rely on P52/P63 acceptance for behavior-level release readiness. Package verification proves package integrity and packaging boundaries, not market correctness or provider availability.

## Safety Boundaries

P64 must preserve all project safety limits:

- no broker interface;
- no trading execution;
- no one-click trading or order delegation;
- no external push;
- no automatic confirmation;
- no automatic rule application;
- no automatic repair, migration, restore, upgrade, or overwrite;
- no return promise;
- no login, paid, authorized, Level2, or high-frequency data source.

## Review Strategy

P64 follows the standard phase cadence:

1. Create change and plan.
2. Sub agent reviews plan.
3. Execute only if no Critical or Important findings remain.
4. Run verification gates.
5. Sub agent reviews execution.
6. Archive.
7. Sub agent reviews submit diff.
8. Commit.
