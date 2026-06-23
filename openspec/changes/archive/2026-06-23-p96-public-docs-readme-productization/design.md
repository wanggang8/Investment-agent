# P96 Design

P96 creates a public documentation front door while preserving existing governance truth sources.

The root `README.md` should serve GitHub readers first. It should answer: what this is, who it is for, what it can and cannot do, how to run it, what the architecture looks like, what data stays local, how CI/release works, and where to read more. It should use existing diagrams from `docs/diagrams/` so the page is visual without adding generated assets or stale screenshots.

`docs/product-overview.md` should explain real user workflows: daily discipline, portfolio maintenance, evidence review, decision consultation, manual confirmation, review/audit, and rule governance. It should keep the safety boundary prominent: no broker connection, no auto-trading, no return promise.

`docs/quickstart.md` should be operational and short. It should point to Docker Compose as the easiest path, then mention source development commands. It should cover `.env.example`, `DEEPSEEK_API_KEY`, local data directory, install/upgrade/uninstall scripts, and health checks.

`docs/README.md` should become a navigable map with sections for product, architecture, API, operations, governance, release evidence, and history. The current long phase log should move to `docs/release/history.md` with wording that preserves caveats and acceptance history.

P96 may update documentation only. Any architecture detail that conflicts with P95 should link to `docs/architecture.md` rather than duplicating implementation claims.
