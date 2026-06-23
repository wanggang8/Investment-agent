# Proposal: P64 Release Packaging And Version Tagging

## Summary

P64 turns the P63 `release_ready` evidence into a local, repeatable release packaging workflow. The phase will define and implement a local release package manifest, version/tag metadata, package verification, and cross-machine repeat verification entrypoints.

P64 is a release-engineering phase. It must not expand Investment Agent runtime product capability, trading capability, external connectivity, or automation boundaries.

## Motivation

P63 refreshed the final release candidate status after full UI regression and real LLM-backed browser validation. The project now has enough evidence to hand off a local release, but it still lacks a concrete package artifact contract:

- what files belong in a local release package;
- which commit/version the package represents;
- how to verify package integrity without unpacking private local state;
- how another machine should repeat the minimum acceptance checks;
- which files must never enter the package.

P64 fills that delivery gap without changing business behavior.

## In Scope

- Add a local release package script that stages a deterministic package directory under `tmp/`.
- Generate a sanitized `release-manifest.json` with commit, release label, package contents, checksums, verification commands, and safety boundaries.
- Produce a compressed local release archive from tracked source, docs, configs, scripts, and web source needed to build/run locally.
- Exclude local secrets, temporary SQLite databases, Playwright traces, logs, `tmp/`, `web/node_modules/`, build outputs, private configs, and raw provider payloads.
- Add package verification that checks manifest presence, checksum consistency, expected entrypoints, and forbidden file patterns.
- Update release docs with P64 packaging, version tagging, repeat verification, and Not Claimed boundaries.
- Run OpenSpec, Go, frontend, E2E smoke, package build, package verify, and sensitive scan gates.

## Out Of Scope

- No broker integration.
- No automatic trading, one-click trading, order delegation, or external push.
- No automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic upgrade, automatic restore, or automatic overwrite of real user databases.
- No new login-gated, paid, authorization-gated, Level2, or high-frequency data source.
- No new HTTP API, SQLite schema, Eino workflow, LLM behavior, or frontend product page unless strictly needed to document packaging status.
- No publishing to package registries, cloud storage, app stores, or remote update channels.
- No claim that public websites, model providers, or future local machines will remain available.

## Deliverables

- `scripts/local-release-package.sh`
- Package verification mode or companion script if the package script is clearer with subcommands.
- `docs/release/release-packaging-2026-06-18.md`
- Updated `docs/release/README.md`, `docs/release/release-handoff-2026-06-18.md`, `docs/development-plan.md`, `docs/README.md`, `docs/GOVERNANCE.md`, `AGENTS.md`, `openspec/project.md`, and `openspec/PROGRESS.md`
- OpenSpec archive for `p64-release-packaging-version-tagging`

## Validation

- Plan review by sub agent before execution.
- `openspec validate p64-release-packaging-version-tagging --strict`
- `openspec validate --all --strict`
- `git diff --check`
- `go test ./...`
- `npm --prefix web test`
- `npm --prefix web run build`
- `bash scripts/e2e-smoke.sh`
- `bash scripts/local-release-package.sh --release-label p64-rc --output-dir tmp/p64-release`
- Package verification command against the generated archive.
- Sensitive and forbidden file scan over the generated manifest and package archive listing.
- Execution review by sub agent before archive.
- Submit review by sub agent before commit.
