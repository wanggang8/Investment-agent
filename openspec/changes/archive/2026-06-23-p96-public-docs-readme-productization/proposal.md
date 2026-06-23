# P96 Public Docs README Productization

## Why

The repository is now public, but it does not yet have a root `README.md`. The existing `docs/README.md` is useful as internal release history, but too long and phase-heavy for a new open-source reader. The project needs a product-quality documentation front door that explains what Investment Agent is, what it does not do, how to run it locally, and where to find deeper technical material.

## What Changes

- Add a root `README.md` with product overview, diagrams, feature summary, safety boundaries, quickstart, CI/release status, and documentation links.
- Add a concise `docs/product-overview.md` for user-facing product concepts and real-use workflows.
- Add a concise `docs/quickstart.md` for local installation, configuration, start, upgrade, uninstall, and troubleshooting.
- Slim `docs/README.md` into a documentation map and move the long phase history into a release/history document.
- Reuse existing diagram assets instead of inventing unrelated visuals.

## Out Of Scope

- Changing L1 requirement semantics in `docs/requirements.md`.
- Changing runtime behavior, API behavior, SQLite schema, CI gates, or release package behavior.
- Rewriting all historical release acceptance records.
- Broker integration, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.

## Acceptance

P96 is acceptable only if:

- GitHub root README exists and gives a new reader a clear first five-minute understanding.
- README uses existing diagrams/screenshots where appropriate and renders with relative Markdown links.
- `docs/README.md` becomes a concise navigation map, not a phase log.
- The moved history remains discoverable and does not erase release caveats.
- `docs/requirements.md` remains the L1 requirements truth source and is not rewritten as marketing copy.
- Documentation links are checked for obvious broken references.
