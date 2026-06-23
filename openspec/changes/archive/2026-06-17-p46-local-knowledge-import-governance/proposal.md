# P46: 本地知识库与数据导入治理

## Summary

新增本地知识导入治理能力，让用户把本地笔记、研究摘录或 CSV/表格化材料先经过只读校验、脱敏预览和索引重建计划检查，再由用户显式确认写入本地证据事实与 RAG 文本块。P46 只处理本地用户提供的文本/表格材料，不接外部授权源、不接券商、不自动交易、不自动应用规则。

## Why

P38 已加固 RAG/VecLite 检索质量，P43 已提供数据质量可观测，P44 已补齐本地诊断打包。下一步需要允许用户安全纳入自己的研究材料，但不能让“导入文件”变成静默污染知识库、暴露私钥/完整 key/SQL/本地隐私路径，或绕过证据分层和规则裁决的入口。

## What Changes

- 新增后端本地知识导入校验 API：`POST /api/v1/local-knowledge/imports/validate`。
  - 输入为页面粘贴/上传解析后的本地行数据，不由后端直接读取任意文件路径。
  - 输出行级校验结果、脱敏预览、风险项、内容 hash、索引重建计划和 `import_batch_id`。
- 新增后端本地知识导入确认 API：`POST /api/v1/local-knowledge/imports/confirm`。
  - 仅在校验无阻塞错误且用户显式确认时写入本地 SQLite 事实。
  - 复用现有 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 和 `audit_events`。
  - 默认写入为 `source_level=C`、`evidence_role=background`，不得伪装成 A/S 正式独立信源。
- 新增前端 `/local-knowledge` 页面。
  - 允许用户粘贴 JSON/CSV 风格本地材料并先校验。
  - 展示脱敏预览、阻塞项、警告项、索引重建计划和确认结果。
  - 不展示完整 key、完整私有路径、原始 SQL、完整 prompt 或 raw HTTP 响应。
- 更新文档、前端契约和 E2E smoke，确保页面可达且无高风险入口。

## Scope

- 本地手工导入治理：validate -> preview -> explicit confirm -> local facts。
- 只接收请求体中的本地材料行，不接收服务端文件路径，不扫描本机目录。
- 导入材料只能作为背景知识或用户研究材料；正式裁决仍由规则、证据验证和人工复核链路处理。

## Out of Scope

- 券商接口、自动交易、一键交易、代下单。
- 外部推送、自动确认、自动规则应用、自动修复承诺。
- 收益承诺、确定性涨跌预测。
- 登录源、付费源、授权源、Level2、高频源。
- 后端读取任意本地文件路径、目录扫描、浏览器爬虫、云同步。
- 将本地用户材料标记为 A/S 高等级正式信源。

## Validation

- `go test ./...`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `bash scripts/e2e-smoke.sh`
- P46 安全扫描（见 `tasks.md` 7.8）
- `openspec validate p46-local-knowledge-import-governance --strict`
- `openspec validate --all --strict`
- `git diff --check`
