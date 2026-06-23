## Context

P10 is marked done in `openspec/PROGRESS.md`, but the post-P10 roadmap now splits remaining work into P11-P18. The first step is governance cleanup: verify no unexpected active changes remain, align phase tracking, and encode the execution rule that each change is proposed, implemented, verified, reviewed by a readonly subagent, then archived before the next change begins.

## Goals / Non-Goals

**Goals:**
- Make the repository's OpenSpec state match the post-P10 roadmap.
- Add an explicit P11-P18 phase progression to governance and progress documents.
- Require readonly subagent review before each archive.
- Preserve the existing L1 contract governance: changes are made through OpenSpec delta and merged during archive.

**Non-Goals:**
- No backend feature implementation.
- No frontend feature implementation.
- No change to automatic trading, broker API, active stock-picking, or return-guarantee boundaries.
- No implementation of P12-P18 work inside P11.

## Decisions

1. Treat P11 as a governance-only change.
   - Rationale: Later changes need a clean baseline before feature work starts.
   - Alternative considered: Fold governance cleanup into P12. Rejected because it would mix state cleanup with data-provider implementation.

2. Keep `openspec/PROGRESS.md` as the machine-readable phase pointer.
   - Rationale: Existing project instructions already use it for phase progression.
   - Alternative considered: Track only in the external roadmap plan. Rejected because the plan is not a governance source of truth.

3. Make subagent review a pre-archive gate.
   - Rationale: The P10 process relied on repeated readonly reviews to catch governance and implementation drift.
   - Alternative considered: Run reviews only at the end of all P11-P18 work. Rejected because each change must stay independently reviewable and archivable.

## Risks / Trade-offs

- Governance-only changes can look light → Mitigation: include concrete checks for active changes, OpenSpec validation, and progress alignment.
- Updating phase tables may drift later → Mitigation: require archive-time update after every change.
- Subagent review can slow each phase → Mitigation: keep reviews readonly and scoped to the current change.
