# P96 Public Docs README Productization Acceptance

> Date: 2026-06-23

## Scope

P96 productized the public documentation front door only. It did not change runtime code, CI workflows, SQLite schema, API behavior, L1 requirements semantics, P95 engineering files, or release package behavior.

## Files

- `README.md`: new GitHub-facing project overview with existing diagram assets from `docs/diagrams/`.
- `docs/README.md`: concise documentation map.
- `docs/product-overview.md`: user-facing workflow and safety boundary overview.
- `docs/quickstart.md`: Docker Compose and local operations guide aligned with `.env.example` and `docs/deployment.md`.
- `docs/release/history.md`: moved long phase history and caveats from the former docs README.
- `openspec/changes/p96-public-docs-readme-productization/tasks.md`: checkbox status only.

## Boundary Checks

- Existing active change directories at implementation time: `p95-architecture-api-engineering-hardening` and `p96-public-docs-readme-productization`.
- P95 owns detailed architecture/API engineering hardening; P96 only links to architecture material and does not edit `docs/architecture.md`.
- P96 docs preserve the boundary against broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, and high-frequency data.

## Validation

Validation performed after the documentation edits:

- `openspec validate --all --strict`: 36 passed, 0 failed.
- local Markdown link sanity check for changed Markdown files: 54 local links checked, 0 missing.

P96 was not archived in this implementation turn.
