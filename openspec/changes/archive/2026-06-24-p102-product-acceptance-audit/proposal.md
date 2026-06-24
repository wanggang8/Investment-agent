# P102 Product Acceptance Audit

## Why

After real OpenAI-compatible LLM access was restored, the project needs a product-level acceptance audit that verifies whether a user can complete the core investment-discipline workflow, not only whether machine tests pass.

## What Changes

- Run local-source product acceptance using the configured real LLM, local backend, Vite frontend, SQLite, and VecLite.
- Capture current-run screenshots for key product routes and workflows.
- Assess usability, product design reasonableness, feature completeness, data/readback trust, accessibility risks, and safety boundaries.
- Add a product acceptance report and screenshot evidence under release UI audit assets.

## Scope Boundaries

- Does not modify runtime product behavior unless a release-blocking defect is discovered and explicitly handled under a follow-up change.
- Does not validate Docker, install/upgrade/uninstall, GitHub Release, distribution package refresh, or physical second-machine runs.
- Does not commit or print API keys.
- Does not claim broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.
