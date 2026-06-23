# P96 Tasks

## 1. Scope And Boundary Review

- [x] Confirm no active OpenSpec change conflicts with P96.
- [x] Confirm P96 does not edit runtime code, CI workflows, SQLite schema, or API behavior.
- [x] Confirm P95 owns detailed `docs/architecture.md` engineering corrections.

## 2. Root README

- [x] Create root `README.md` with product positioning, feature map, diagrams, quickstart, safety boundaries, CI/release notes, and documentation links.
- [x] Use existing assets from `docs/diagrams/` with relative paths.
- [x] Avoid claims of broker connectivity, auto-trading, return guarantees, paid/login data sources, Level2, or high-frequency data.

## 3. Product Overview And Quickstart

- [x] Create `docs/product-overview.md` for user-facing workflows and product boundaries.
- [x] Create `docs/quickstart.md` for Docker Compose, local configuration, start, upgrade, uninstall, health check, and troubleshooting.
- [x] Keep operational details consistent with `docs/deployment.md` and `.env.example`.

## 4. Documentation Map And History Split

- [x] Replace `docs/README.md` with a concise documentation map.
- [x] Move the current long phase status narrative into `docs/release/history.md`.
- [x] Ensure release caveats and historical acceptance boundaries remain discoverable.

## 5. Documentation Validation

- [x] Check Markdown links for obvious broken local references.
- [x] Generate P96 acceptance evidence.
- [x] Run OpenSpec validation and documentation-focused checks.
- [x] Archive P96 only after validation passes.
