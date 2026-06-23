# P46 Design

## Overview

P46 把“本地知识导入”设计成双阶段治理流程：先校验和脱敏预览，再由用户确认写入。实现上复用现有证据事实表和 RAG chunk 表，避免新增复杂 schema；通过 `source_level=C`、`evidence_role=background` 和审计事件明确这些材料是用户自有背景材料，不是高等级独立信源。

## API Flow

1. `POST /api/v1/local-knowledge/imports/validate`
   - 请求体包含 `source_label`、`default_symbol`、`rows`。
   - 每行支持 `title`、`text`、`symbol`、`as_of_date`、`tags`；请求中若附带 `source_url` 也只参与风险校验和 batch hash，不写入事实表。
   - 服务端校验必填、长度、疑似 secret、SQL、私有路径、raw HTTP/prompt 等风险；`source_label`、`source_url` 与 tags 也纳入风险扫描。
   - 服务端返回 `import_batch_id`、行级状态、脱敏预览、`blocking_count`、`warning_count`、`rag_chunk_count`、`index_status=pending`。
   - validate 不写业务事实；`import_batch_id` 必须由规范化后的 `source_label`、`default_symbol`、行内容 hash 和行顺序生成。

2. `POST /api/v1/local-knowledge/imports/confirm`
   - 请求体包含 `import_batch_id`、`confirm_reason`、`source_label`、`default_symbol`、`rows`。
   - 服务端重新执行同一套校验，防止绕过 validate。
   - 服务端必须重算 `import_batch_id`，并拒绝与请求 `import_batch_id` 不匹配的确认请求。
   - 只在无 blocking 行时写入 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 和 `audit_events`。
   - 写入 metadata 标记 `source_type=local_knowledge_import`、`import_batch_id`、`tags`、`content_hash`。

## Persistence

不新增 migration。导入确认复用：

- `intelligence_items`：保存本地来源标签、C 级 source level、脱敏后的 title/ref 和 content hash；P46 `original_url` 固定为空。
- `intelligence_summary`：保存摘要文本、symbol/entity、`source_level=C`、`evidence_role=background`。
- `rag_chunks`：保存可重建索引的 chunk，`index_status=pending`。
- `source_verifications`：保存 background_only 语义，独立信源计数为 1，高等级计数为 0。
- `audit_events`：记录 `run_local_task` 或等价本地任务审计，输入引用只含 batch id 与行数。

## Frontend

新增 `/local-knowledge` 页面并放入主导航。页面以本地运维/研究材料导入为主，不做营销式布局：

- 输入区：source label、default symbol、JSON/CSV 风格文本。
- 校验结果：行级状态、阻塞/警告统计、脱敏预览、索引计划。
- 确认区：仅校验通过后启用确认按钮；确认成功后展示写入计数和审计 id。
- 安全区：明确不接券商、不外推、不自动应用规则、不承诺收益，且本地材料默认仅作为背景材料。

## Risk & Guardrails

- 后端不读取用户提供的本地文件路径，只处理请求体内容。
- suspected secrets、完整 key、私有路径、SQL、raw HTTP、prompt 进入 blocking 或 warning，响应仅展示脱敏预览；raw HTTP 或完整 prompt 命中时整段预览替换为脱敏占位。
- confirm 必须重新校验请求体，不能只信任 `import_batch_id`。
- `import_batch_id` 必须绑定规范化 payload；用户修改 source、symbol、rows 或行顺序后必须重新 validate。
- 本地材料默认 C/background，不得提升为 formal/A/S，不得直接满足重大事件多源验证。
