## Context

Evidence quality is already represented through `source_level`, `evidence_role`, `time_weight`, `relevance_score`, independent source counts, and verification status. P15 makes these semantics explicit and strengthens tests around their propagation and boundaries.

## Goals / Non-Goals

**Goals:**
- Preserve time weight and relevance score from source normalization through DTOs and decision evidence refs.
- Enforce C-level source background-only behavior across refresh, retrieval, and decision evidence paths.
- Enforce major-event satisfaction only when at least two A/S independent sources exist.
- Distinguish structured facts used by rules from analyst materials used for explanation.

**Non-Goals:**
- No new entity extraction pipeline.
- No complex relevance scoring model.
- No expanded event classification taxonomy beyond current event types.
- No changes to final rule arbitration authority.
- No frontend redesign; display improvements can be handled by P16.

## Decisions

1. Treat structured facts as rule inputs.
   - `market_snapshots`, `intelligence_summary`, `source_verifications`, and `evidence_refs` provide rule-readable facts.
   - Analyst outputs remain materials and cannot set final verdicts.

2. Keep role enforcement in domain and retrieval paths.
   - C-level evidence cannot become formal in retrieval fallback or decision records.
   - Major events need at least two high-grade independent sources.

3. Keep optional enrichment out of this change.
   - Entity extraction, richer event typing, and complex relevance scoring need separate data and UI decisions.

## Risks / Trade-offs

- Stronger quality checks may classify more evidence as background-only.
- Some metadata may be absent from older rows; APIs should preserve known values without inventing placeholders.
- P16 may improve visibility, but P15 must keep backend and DTO semantics stable first.
