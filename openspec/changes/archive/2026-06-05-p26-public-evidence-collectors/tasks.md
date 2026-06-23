# P26 公告与证据源 collector 任务

## 1. 前置确认

- [x] 复核 P25 验证结论
  - 阅读 `openspec/changes/archive/2026-06-05-p25-real-public-data-sources/verification.md`。
  - 确认 P26 首批范围仅包含巨潮资讯、深交所、证监会；AMAC 行业统计/自律栏目作为二线背景源。
  - 确认上交所、AMAC 机构/产品/人员查询、东方财富基金、中证指数、新浪财经不进入 P26 首批实现。

- [x] 复核法律声明、robots 和访问频率
  - 对巨潮资讯、深交所、证监会、AMAC 公开页面补充只读、低频、公开访问边界检查。
  - 若发现禁止自动化采集、登录限制、验证码或付费/授权限制，将对应源移出首批实现范围。

## 2. Collector 设计

- [x] 定义公开证据 collector 接口
  - 支持按 `symbol`、日期范围、分页参数执行只读抓取。
  - 返回标准公告/证据 JSON，包含 `source_name`、`source_level`、`source_type`、`evidence_role`、`symbol`、`title`、`text`、`url`、`attachment_url`、`published_at`、`captured_at`、`content_hash`、`raw`。

- [x] 设计去重与幂等策略
  - 使用 `source_name + source_record_id` 或 `content_hash` 去重。
  - 重复抓取不得重复写入 `intelligence_items` 或 `rag_chunks`。
  - 内容更新时保留审计记录，不覆盖历史事实。

- [x] 设计失败与降级策略
  - 源不可用、分页失败、附件失败、字段缺失、解析失败分别映射错误码。
  - 失败写入 `audit_events`，不阻塞其他源和本地应用启动。
  - 少于 2 个 A/S 独立信源时保持 `insufficient_data` 或 `frozen_watch`，不得生成正式高置信结论。

## 3. 首批源实现范围

- [x] 实现巨潮资讯公告 collector
  - 使用 P25 验证的 `hisAnnouncement/query` 公开请求。
  - 支持最近 90 天补采和低频增量刷新。
  - 解析公告标题、证券代码、证券名称、公告时间、附件 URL、公告类型和分页信息。
  - 写入 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。

- [x] 实现深交所公告 collector
  - 使用 P25 验证的 `api/disc/announcement/detailinfo` 与 `searchQuery` 公开请求。
  - 支持最近 90 天补采和低频增量刷新。
  - 解析证券代码、证券名称、公告标题、附件路径、附件格式、公告 ID、分类、发布时间和分页信息。
  - 写入 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。

- [x] 实现证监会监管信息 collector
  - 使用 P25 验证的 `getChannelList` 与 `searchList` 公开 JSON 请求。
  - 覆盖监管公告、规则、处罚、市场禁入、监管措施等适合公开采集的栏目。
  - 解析标题、正文/摘要、URL、发布时间、栏目、文号、发布机构和附件。
  - 写入 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。

- [x] 评估 AMAC 二线背景源是否暂缓
  - AMAC 不作为 P26 首批必需 collector 或 runtime dependency。
  - 仅当法律声明、robots、稳定字段和访问边界均满足时，才可实现 P25 已验证的行业统计 JSON 或稳定自律栏目。
  - 不实现机构/产品/人员查询、考试、培训、报送、登录业务系统。
  - P26 首批实现中已明确暂缓 AMAC，只保留为后续二线背景候选。

## 4. 测试与验收

- [x] 增加 collector 单元测试
  - 使用 fixture/httptest 模拟公开 JSON 响应、分页、空数据、字段缺失和失败响应。
  - 覆盖去重、source_level、evidence_role、content_hash 和错误码。

- [x] 增加入库集成测试
  - 验证 collector 输出可写入 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。
  - 验证重复执行不会重复写入事实。

- [x] 增加本地任务或 API 验收路径
  - 明确通过 `cmd/agent` 或现有 evidence refresh 触发 P26 collector 的方式。
  - 失败时输出可解释降级状态。

- [x] 运行后端测试
  - `go test ./...`

- [x] 运行 OpenSpec 校验
  - `openspec validate --all --strict`

## 5. 文档同步

- [x] 更新 `docs/development-plan.md`
  - 将 P26 范围、首批源、暂缓源和验收目标写清楚。

- [x] 更新 `docs/configuration.md`
  - 补充 P26 collector 配置项或说明；不得写真实密钥或登录凭证。

- [x] 按需更新 `docs/data-model.md`、`docs/api.md`、`docs/frontend-contract.md`
  - 仅当 P26 新增字段、API 或前端契约时更新。

- [x] 更新 `openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`AGENTS.md`
  - 实现完成和归档前同步阶段状态。
