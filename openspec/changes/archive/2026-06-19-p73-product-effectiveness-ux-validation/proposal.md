# P73 Product Effectiveness And UX Validation

## Why

P71/P72 prove that the product can run real local acceptance flows, use real public evidence, use a real LLM provider, rebuild VecLite, operate UI journeys, and persist expected SQLite data impact.

They do not yet prove the higher product goal: whether the product helps a real user execute investment discipline, understand risk/evidence/rules, avoid unsafe or unsupported decisions, and use the UI without confusion.

P73 adds that missing product-effectiveness and UX validation layer.

## What Changes

- Define product-goal effectiveness metrics aligned with `docs/requirements.md`: discipline adherence, evidence sufficiency, traceability, safe blocking, controlled rule evolution, and review usefulness.
- Add a real UX task matrix for first-use, daily discipline, portfolio maintenance, consultation, evidence review, manual confirmation, risk handling, review/audit, and rule proposal understanding.
- Add effect replay checks for discipline scenarios such as insufficient evidence, C-level background-only material, capability/discipline limits, risk escalation, and rule proposal gatekeeping.
- Add a user-facing UX audit record that reviews whether page hierarchy, copy, navigation, and state labels support the actual product goal.
- Produce release evidence that states what passed, what is still gap-qualified, and what must not be claimed.

## Scope

In scope:

- Product effectiveness acceptance materials.
- UX task matrix and browser-operable validation.
- Deterministic effect replay checks using local facts and fixtures.
- UI/UX audit notes tied to actual screens and real workflows.
- Release/readiness wording updates.

Out of scope:

- Broker interfaces, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, or real DB overwrite.
- Investment return promises or proof that the system improves future investment returns.
- Paid/login/authorized/Level2/high-frequency sources.
- Physical second-machine package repeat.
- Post-P72 package refresh unless explicitly opened as a separate stage.

## Expected Outcome

P73 may claim product-effectiveness/UX validation only if:

- Real UX tasks pass without critical confusion, page errors, unexpected API failures, or forbidden affordances.
- Discipline effect replay shows unsafe/unsupported scenarios are blocked or gap-qualified.
- Evidence/rule/LLM/UI relationships are understandable from the UI and artifacts.
- Remaining gaps are explicitly recorded.

P73 must not claim that investment returns are improved, future market direction is predicted, public sources/LLM providers will remain available, or all possible real-world investment scenarios have been exhausted.
