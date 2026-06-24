# P103 Product Acceptance UX Linkage Fixes

## Why

P102 product acceptance found the local-source real-model product flow usable, but identified three non-blocking UX issues that can reduce trust during first use and decision review.

## What Changes

- Map an empty/missing portfolio snapshot to a first-use onboarding state in the portfolio UI instead of a generic system failure.
- Reduce real-model decision detail density by making analyst analysis material progressive-disclosure by default while keeping auditability.
- Make decision loop deep links with `decision_id` focus the requested decision and expose that focus state in the UI.
- Add focused frontend tests and a linked-workflow validation record covering portfolio, consultation, confirmation, loop, audit/review/workbench/notification readback.

## Scope Boundaries

- Does not modify investment rules, final verdict logic, LLM prompts, SQLite schema, or HTTP API contracts.
- Does not add broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, or automatic rule application.
- Does not validate Docker, installation, GitHub Release, package refresh, or physical second-machine execution.
