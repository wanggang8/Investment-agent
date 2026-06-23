# P71 Real Product Acceptance True Pass

## Why

P70 closed the previous milestone as `release_ready_limited_current_data_scope`. That status was honest but limited: P66 current-data strict gate still reports `policy=blocked` / `gate=block`, P67 excludes current local data health from clean release claims, and P63 recorded a non-blocking `VECTOR_INDEX_UNAVAILABLE` degradation during the real UI consultation journey.

The user now requires full real product acceptance: all product functionality must pass real local UI operation and current evidence must not rely on mocks, scope exclusion, or non-blocking retrieval degradation. P71 converts those remaining caveats into blocking acceptance gates and refreshes the final package only after the strict gates pass.

## What Changes

- Establish P71 as the strict full-product acceptance stage after P70.
- Resolve current local data quality to a true P66 pass: `policy=passed` / `gate=pass` without a P67 scope exclusion.
- Harden acceptance setup so VecLite/RAG index health is prepared before real UI consultation and `VECTOR_INDEX_UNAVAILABLE` blocks the P71 result.
- Re-run full real UI route regression and real LLM-backed consultation against the local backend/frontend without using frontend mocks as pass evidence.
- Refresh release materials and produce a post-P70 package only if strict current-data, VecLite, real UI, real LLM, safety, and packaging gates pass.

## In Scope

- OpenSpec change, tasks, release acceptance records, release candidate/handoff wording, progress/governance updates, and final package refresh evidence.
- Runtime, script, fixture, local seed, collector, or acceptance-test changes needed to make real current data, VecLite preparation, and UI acceptance deterministic and truthfully verifiable.
- Current data verification using `data-source-quality-regression --source current --strict-quality-gate`.
- Browser validation through the real backend and real Vite frontend, including primary routes, key user operations, real consultation, decision detail, data-quality, evidence/index rebuild, settings/refresh paths, mobile and desktop reflow, console/page/API failure classification, and forbidden capability scans.

## Out of Scope

- No broker interface, trading execution, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair promise, automatic migration, automatic upgrade, automatic restore, automatic overwrite of real user databases, return promise, login-gated source, paid source, authorization-gated source, Level2 data, or high-frequency source.
- No claim that future public websites, model providers, or market data will remain available.
- No use of P67 scope exclusion, frontend mocks, fixture-only data, or documented waivers as the pass criterion for P71 full product acceptance.
- No package refresh before strict acceptance gates pass.

## Impact

P71 may modify scripts, acceptance tests, data-source handling, VecLite preparation, local seed/diagnostic paths, release documentation, and packaging artifacts. If current public sources or real LLM providers are unavailable, P71 must record `release_blocked` or a narrower non-pass status rather than relabeling the run as full real acceptance.
