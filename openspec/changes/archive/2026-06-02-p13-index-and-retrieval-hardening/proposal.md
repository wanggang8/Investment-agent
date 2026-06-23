## Why

P10/P12 已保留 JSON 文件索引作为本地可重建索引，但缺少显式健康状态、重建统计和检索降级可观测性。P13 需要让索引缺失、损坏、不兼容和 SQLite 摘要降级都能被 API 与前端理解。

## What Changes

- Keep the JSON file index as the local rebuildable index adapter.
- Add index health status, version/compatibility checks, rebuild statistics, and degradation reasons.
- Make retrieval fallback to SQLite summaries explicit and auditable.
- Expose index health through service/API DTOs for frontend display.
- Document the future replacement boundary for a real VecLite API.

## Capabilities

### New Capabilities
- `index-and-retrieval-hardening`: Covers local index health, rebuild stats, retrieval fallback observability, and future VecLite replacement boundaries.

### Modified Capabilities
- `real-data-integration`: Refines RAG/VecLite retrieval behavior to include local JSON index health, corruption/incompatibility handling, and observable SQLite fallback.

## Impact

- `internal/application/service` vector index and retrieval adapter.
- Evidence/index handler DTOs and tests.
- Frontend-visible API fields if existing endpoint already returns index status.
- Documentation for configuration and recovery.
- No external VecLite API integration in this change.
