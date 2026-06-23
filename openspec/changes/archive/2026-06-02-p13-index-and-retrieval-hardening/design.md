## Context

The current local index is a JSON file adapter around `rag_chunks`. It is intentionally rebuildable from SQLite and not the source of truth. P13 improves operational visibility rather than replacing it with a real VecLite API.

## Goals / Non-Goals

**Goals:**
- Detect missing, corrupted, and incompatible local index files.
- Report index health, chunk count, rebuild count, last rebuild time, and degradation reason.
- Keep SQLite summaries as the stable fallback when local index search fails or is empty.
- Keep frontend/API consumers able to distinguish healthy, missing, corrupted, incompatible, rebuilding, and degraded states.

**Non-Goals:**
- No real VecLite API dependency.
- No semantic embedding model integration.
- No P14-P18 workflow or UI feature expansion.
- No trading behavior.

## Decisions

1. Add metadata to the JSON file index envelope.
   - Rationale: A versioned envelope makes incompatible files distinguishable from plain corruption.
   - Alternative considered: infer compatibility from chunk fields only. Rejected because missing fields are ambiguous.

2. Keep rebuild source as SQLite `rag_chunks`.
   - Rationale: SQLite remains the factual source; index files are derived artifacts.
   - Alternative considered: recover from the index file itself. Rejected because damaged indexes cannot be trusted.

3. Expose health through application DTOs rather than direct file access.
   - Rationale: frontend must not read local files.
   - Alternative considered: document path inspection. Rejected because it violates frontend boundaries.

## Risks / Trade-offs

- Existing JSON index files may be legacy arrays → classify as incompatible and rebuild from SQLite.
- More status fields add API surface → keep fields descriptive and read-only.
- Health checks must not fail user-facing pages → return degraded status and reason instead of panicking.
