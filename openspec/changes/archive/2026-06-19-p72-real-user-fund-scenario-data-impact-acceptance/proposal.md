# P72 Real User Fund Scenario Data Impact Acceptance

## Why

P71 proved the strict release gate: current data true pass, healthy/fresh VecLite, real local UI operation, and real LLM-backed consultation all passed. The user now requires the next level of confidence: realistic user scenarios must prove that actual product operations affect the right local data, audit events, derived views, and follow-up analysis.

P72 adds a scenario-and-data-impact acceptance layer. It does not replace P71. It verifies that a real fund/ETF user journey behaves correctly across UI, API, SQLite, audit, notifications, daily discipline, risk alerts, decision detail, decision loop, and review surfaces.

## What Changes

- Establish P72 as the real user fund scenario and data-impact acceptance stage after P71.
- Create a reviewed scenario matrix before execution.
- Add a P72 strict acceptance script that runs against a temporary local SQLite database, real backend, Vite frontend, real LLM config, `use_stub=false`, and real public collector configuration.
- Add browser automation for a realistic `510300` fund/ETF journey:
  - portfolio calibration and holding maintenance,
  - local knowledge import and VecLite rebuild,
  - market/current-data refresh and quality gate check,
  - daily discipline and risk alert generation,
  - real LLM consultation and decision detail,
  - manual confirmation / offline transaction record,
  - decision loop / review / notifications / audit readback,
  - safety-boundary checks.
- Add read-only SQLite data-impact verification after UI operation, writing a sanitized JSON summary under P72 evidence artifacts.
- Update release and governance materials only after P72 evidence passes.

## In Scope

- OpenSpec change, scenario matrix, acceptance record, UI evidence artifacts, scripts, Playwright specs, and release/governance/progress updates.
- Temporary local databases and configs under `tmp/`.
- Read-only SQLite verification of expected table changes and numeric consistency.
- Accuracy checks for deterministic calculations: position market value, unrealized profit ratio, cash / total assets, report and alert links, confirmation status, audit counts, and retrieval/index state.
- Analysis-quality checks that are appropriate for LLM and investment material: source traceability, parsed/quality-passed analyst reports, rule-based final verdict, safety wording, and no investment-return guarantee.

## Out of Scope

- No guarantee of future market direction, future returns, or investment judgment correctness.
- No broker interface, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic upgrade, automatic restore, automatic overwrite of real databases, login-gated source, paid source, authorization-gated source, Level2 data, or high-frequency data.
- No physical second-machine repeat unless a separate change is opened.
- No new public-provider capability beyond already configured read-only collectors.
- No new product API solely to expose internal SQLite tables.

## Impact

P72 may add acceptance scripts, browser tests, evidence records, and focused hardening if the real scenario exposes product bugs. If an external public source or real LLM provider is unavailable, P72 records a blocked status rather than converting the run into mock-only pass evidence.
