## Why

P15 focuses on evidence quality. Existing code already preserves source level, role, source counts, timestamps, and verification status, but the next stage needs clearer quality semantics: freshness weight, formal/background boundary, independent source counts, and separation between structured facts and analyst materials.

## What Changes

- Enrich evidence quality rules around time weight, source role, and independent source counts.
- Ensure C-level sources remain background-only.
- Ensure fewer than two A/S independent sources cannot satisfy major-event verification.
- Distinguish structured facts used by rules from analyst materials produced by LLM or local analysis.
- Keep entity extraction, event classification expansion, and complex relevance scoring as later work.

## Capabilities

### New Capabilities
- `evidence-quality-enrichment`: Defines evidence quality metadata, formal/background boundaries, structured-fact vs analyst-material separation, and acceptance criteria for source quality.

### Modified Capabilities
- `real-data-integration`: Refines evidence metadata and retrieval requirements around quality fields and analyst-material boundaries.
- `product-completeness`: Continues the product-grade evidence baseline without changing rule-first arbitration.

## Impact

- Evidence normalization and retrieval services.
- Evidence DTOs and decision evidence references.
- Tests for C-level background-only behavior, A/S independent source counts, and quality metadata propagation.
- No automatic trading, active recommendation, return guarantee, or rule-arbitration change.
