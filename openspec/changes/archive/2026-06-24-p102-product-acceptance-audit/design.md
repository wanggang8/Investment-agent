# Design

## Acceptance Method

P102 combines three evidence types:

1. Machine gates: OpenSpec, Go, frontend tests/build, P92/P93, and real LLM smoke where appropriate.
2. Real browser product journey: capture screenshots of core routes and workflows against a local Go backend and Vite frontend.
3. Product audit: evaluate the captured evidence against the user's goal of maintaining a disciplined, evidence-backed investment workflow.

## Product Lenses

The audit uses these lenses:

- Task entry and discoverability.
- Information architecture and navigation clarity.
- Core workflow completion from portfolio setup to consultation, decision review, confirmation, risk, governance, and audit readback.
- Trust signals: evidence, source status, SQLite/readback, audit trail, and safe degradation.
- Accessibility and responsive risks visible from screenshots and browser checks.
- Boundary honesty: no trading, no broker/order affordances, no return promises.

## Evidence Location

Screenshots and notes are stored under:

```text
docs/release/ui-audit-assets/2026-06-24-product-acceptance-audit/
```

The acceptance summary is stored under:

```text
docs/release/acceptance/2026-06-24-p102-product-acceptance-audit.md
```
