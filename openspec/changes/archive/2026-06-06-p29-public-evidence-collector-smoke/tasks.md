# P29 公开证据 collector 真实采集验收修复任务

## 1. 前置复核

- [x] 复核当前 P26 collector 与入库链路
  - 确认 `public-evidence-refresh` 入口、配置开关和默认 disabled/stub 边界。
  - 确认 `PublicEvidenceIngestionService` 写入 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications`、`audit_events`。
  - 确认现有测试覆盖 fixture 入库、幂等、partial repair、多源验证和失败审计。

- [x] 复核安全边界
  - 不接券商交易 API。
  - 不自动交易，不创建订单、交易、确认或外部通知。
  - 不登录、不绕权、不使用付费或授权源。
  - 不高频抓取；真实 smoke 只做少量请求。

## 2. 真实源接口核验

- [x] 核验巨潮资讯当前接口
  - 找到至少一个近期存在公告的可采样标的或 ETF/基金相关公开记录。
  - 验证请求参数、分页、时间窗口、响应字段和附件 URL。
  - 区分无公告与接口不可用。

- [x] 核验深交所当前接口
  - 找到至少一个近期存在公告的可采样标的或 ETF/基金相关公开记录。
  - 验证请求参数、分页、时间窗口、响应字段和附件 URL。
  - 区分无公告与接口不可用。

- [x] 核验证监会当前接口
  - 重新确认当前公开搜索 API 是否仍可用，若原 `/searchList` 已失效，改用当前可公开访问的只读入口。
  - 若证监会无法稳定提供按标的可采集数据，标记为可降级背景源，不阻塞其他 A 级公开公告源入库。

## 3. 后端修复

- [x] 修复 collector 请求/解析
  - 按核验结果修正巨潮/深交所/证监会 collector 的 URL、参数或响应解析。
  - 保持 `PublicEvidencePayload` 标准 shape。
  - 无数据时返回可诊断的 no_data/not_found 类错误，不伪造成接口失败。
  - 接口不可用或解析失败时保留 source-specific audit。

- [x] 修复真实 smoke 可配置性
  - 支持通过 CLI `--symbol`、`--start-date YYYY-MM-DD`、`--end-date YYYY-MM-DD` 指定 smoke 标的和窗口。
  - 保持默认 `public_evidence.enabled=false`，真实采集必须显式开启。

- [x] 保持入库幂等与 partial repair
  - 重复采集同一 source record 不重复写入 summary/chunk。
  - 已满足的多源验证不得被后续单源降级覆盖。
  - 全部启用源均返回 `no_data` 时，作为成功空刷新写入 `success count=0` 审计，并保留各源 `degraded` 诊断，不误报 source_unavailable 或任务失败。

- [x] 支持 CNInfo orgId 配置映射
  - 通过 `data_sources.public_evidence.cninfo_org_ids` 配置 `symbol -> orgId`。
  - 未配置时保留少量内置公开标的映射；也保留直接传入完整 `symbol,orgId` stock 参数的测试边界。

## 4. 测试与验收

- [x] 增加/更新 collector 单元测试
  - 覆盖当前真实响应 shape 的解析。
  - 覆盖 no_data 与 source_unavailable 的错误分类。
  - 覆盖附件 URL 归一化、分页和时间窗口过滤。

- [x] 增加入库 smoke 测试或脚本化验收
  - 使用临时 SQLite 和显式真实源配置。
  - 运行 `public-evidence-refresh` 后检查至少一条真实源数据写入 `intelligence_items`、`intelligence_summary`、`rag_chunks`。
  - 检查 `source_verifications` 与 `audit_events` 状态。

- [x] 运行后端测试
  - `go test ./internal/application/workflow -run 'Test.*Collector|TestPublicEvidenceIngestion'`
  - `go test ./internal/infrastructure/wiring`
  - `go test ./...`

- [x] 运行 OpenSpec 校验
  - `openspec validate --all --strict`：18 passed, 0 failed。

## 5. 文档同步

- [x] 更新 `docs/configuration.md`
  - 说明真实公开证据采集默认关闭，开启方式和安全边界。

- [x] 更新 `docs/workflow.md` / `docs/data-model.md`
  - 说明真实 collector 的 no data、source unavailable、parse error 语义和入库表。

- [x] 更新 `docs/development-plan.md`、`openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`AGENTS.md`
  - P29 已归档到 `openspec/changes/archive/2026-06-06-p29-public-evidence-collector-smoke/`，阶段状态已标记为 `done`，active change 已清空。
