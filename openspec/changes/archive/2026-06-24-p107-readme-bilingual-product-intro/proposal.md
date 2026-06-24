# P107 README Bilingual Product Intro

## Summary

Add a Chinese README companion and explicit language switch links to the public README so Chinese readers can understand the product purpose, feature set, operating flow, safety boundary, and quickstart path without reading the English-only root README.

## Scope

- Add `README.zh-CN.md` as a Chinese public README.
- Add language switch links between `README.md` and `README.zh-CN.md`.
- Keep the existing English README as the root GitHub landing page.
- Keep product claims aligned with current release boundaries.

## Out Of Scope

- Runtime code changes.
- API, SQLite schema, workflow, frontend contract, or deployment behavior changes.
- Docker or installer validation.
- New release tag or package refresh.
- Any broker connection, automatic trading, one-click trading, delegated order, external push, automatic confirmation, automatic rule application, return guarantee, paid/login/authorized data-source, Level2, or high-frequency-data claim.

## Validation

- OpenSpec validation for the change and the full project.
- Markdown local-link validation for README files.
- Whitespace diff check.
