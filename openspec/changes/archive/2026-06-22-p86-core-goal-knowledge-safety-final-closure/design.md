# P86 Design

## Evidence Strategy

P86 is the final integrator. It should not hide gaps left by P81-P85 or P87. It starts from the P87 matrix, owns the exact 137 full-release-required rows still not `real_pass`, reconciles all evidence into one row-level matrix, and runs real end-to-end scenarios that test the product goal:

- The user can set up and maintain a local portfolio.
- The product uses formal evidence and bounded built-in knowledge.
- Data readiness explains missing or degraded information.
- Consultation and expected return analysis are traceable and safe.
- Risk alerts, SOPs, manual confirmations, reviews, and audits create coherent readback.
- LLM/RAG enriches context without becoming final decision authority.
- The UI remains understandable and safe across desktop/mobile.

## Execution Plan

P86 has four execution tracks. Each track must produce command output, an evidence artifact, and row-level conclusions.

1. Inventory and plan gate: prove the P87 matrix has exactly 137 remaining full-release-required non-`real_pass` rows and that P86 owns all of them.
2. Integrated product runner: operate the real local UI through setup/portfolio, data readiness, knowledge/RAG, consultation, expected return, risk/SOP, manual confirmation, review, audit, and release-safety surfaces.
3. Row matrix closure: generate a P86 matrix that only upgrades a row when the evidence directly proves that row; otherwise preserve `partial` and write a concrete blocker or future implementation plan.
4. Governance closure: update release/governance materials, run OpenSpec/test/build checks, request subagent review, fix Critical/Important findings, then archive.

## Evidence Standard

`real_pass` requires all applicable layers: real browser UI action, API/readback, read-only SQLite evidence, workflow metadata or deterministic calculation proof, and safety-boundary verification. For broad goal rows, P86 may use cumulative P81-P87 evidence only when the row-level text is fully covered by named evidence references. P86 must not upgrade rows based only on seeded SQLite decisions, route smoke, screenshots, fixture/mock/stub data, or historical documentation claims.

## Expected Outcomes

The desired outcome is full original-requirement pass. If any row cannot honestly be upgraded, P86 must produce a residual blocker list and a concrete next-plan recommendation instead of claiming completion.

## Final Claim Rule

P86 may claim full original-requirement pass only if the final matrix shows no full-release-required row remains partial, blocked, scoped-only, or unsupported. Otherwise P86 must report the exact remaining rows and blockers.
