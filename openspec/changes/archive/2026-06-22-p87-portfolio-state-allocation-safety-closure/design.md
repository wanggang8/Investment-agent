# P87 Design

## Evidence Strategy

P87 should run an integrated real-user portfolio scenario against a real local backend, frontend, and temporary SQLite database:

- Setup and inspect a portfolio with core, satellite, and cash assets.
- Exercise holding-state transitions for normal, sell-only, and frozen-watch cases.
- Trigger data-insufficient and multi-source-insufficient paths and verify visible safe degradation.
- Run quarterly allocation/rebalance calculations and verify deterministic readback.
- Confirm or reject user-visible proposals manually and verify audit/readback.
- Check dashboard, decision, review, audit, and data-quality surfaces after each material action.
- Scan for forbidden broker, trade, push, auto-confirmation, auto-rule-application, install/upgrade, migration, repair, and restore behavior.

## Real-Pass Rule

A row may become `real_pass` only when the fresh evidence proves the complete row text. Adjacent behavior, route smoke, screenshots without data readback, seeded-only records, or single-field proof are insufficient for broad allocation, attribution, audit-history, or release-safety rows.
