## 1. OpenSpec 与范围

- [x] 1.1 确认当前无活跃 change，P45 已归档，P46 为下一功能候选。
- [x] 1.2 确认 P46 聚焦本地知识导入治理：validate、脱敏预览、显式确认、背景材料入库、索引重建计划。
- [x] 1.3 确认 P46 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、收益承诺、登录/付费/授权/Level2/高频源。

## 2. 后端 DTO 与服务

- [x] 2.1 新增 `internal/application/dto/local_knowledge.go`，定义 validate/confirm request/response、row result、preview、index plan；confirm request 必须包含 `source_label` 和 `default_symbol` 以便重算 batch id。
- [x] 2.2 新增 `internal/application/service/local_knowledge.go`，实现本地知识导入校验、脱敏、hash、chunk 估算、确认写入。
- [x] 2.3 confirm 必须重新校验 rows，重算 `import_batch_id`，并拒绝与规范化 `source_label/default_symbol/rows/content_hashes` 不匹配的确认请求。
- [x] 2.4 确认写入复用 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications`、`audit_events`，默认 `source_level=C`、`evidence_role=background`、`index_status=pending`。
- [x] 2.5 添加 `internal/application/service/local_knowledge_test.go`，覆盖成功校验、secret/SQL/private path 风险、confirm 重新校验、batch id mismatch 拒绝、写入失败回滚、写入计数和 C/background 边界。

## 3. 后端 Handler 与 API

- [x] 3.1 新增 `internal/application/handler/local_knowledge_handler.go`，提供 `ValidateLocalKnowledgeImport` 与 `ConfirmLocalKnowledgeImport`。
- [x] 3.2 在 `internal/application/handler/app.go` 注册：
  - `POST /api/v1/local-knowledge/imports/validate`
  - `POST /api/v1/local-knowledge/imports/confirm`
- [x] 3.3 添加 handler 测试，覆盖成功 validate/confirm、blocking 行拒绝 confirm、响应不泄露完整 key/SQL/私有路径。

## 4. 前端页面与服务

- [x] 4.1 新增 `web/src/types/localKnowledge.ts` 与 `web/src/services/localKnowledge.ts`。
- [x] 4.2 新增 `web/src/pages/LocalKnowledgePage.tsx`，提供本地材料输入、校验结果、脱敏预览、索引计划和确认结果。
- [x] 4.3 新增页面测试，覆盖 validate 调用、blocking 展示、确认按钮启用条件、确认结果和禁止项扫描。
- [x] 4.4 在 `web/src/App.tsx` 注册 `/local-knowledge`，在 `web/src/app/AppLayout.tsx` 增加导航入口。
- [x] 4.5 更新 `web/e2e/local-smoke.spec.ts`，覆盖 `/local-knowledge` 可达与安全文本扫描。

## 5. 文档与契约

- [x] 5.1 更新 `docs/api.md`，新增 P46 本地知识导入 validate/confirm API。
- [x] 5.2 更新 `docs/data-model.md`，说明 P46 复用证据事实表与 C/background 边界。
- [x] 5.3 更新 `docs/frontend-contract.md`，新增 `/local-knowledge` 页面契约。
- [x] 5.4 更新 `docs/development-plan.md`、`openspec/project.md`、`openspec/PROGRESS.md`、`docs/GOVERNANCE.md` 和 `AGENTS.md` 当前阶段状态。
- [x] 5.5 在 OpenSpec delta 中记录 P46 行为要求。

## 6. 执行前复审

- [x] 6.1 计划完成后执行只读子 agent 复审，确认无 Critical / Important。
- [x] 6.2 复审通过后再执行实现任务。

## 7. 验收

- [x] 7.1 运行 `go test ./...`。
- [x] 7.2 运行 `npm --prefix web test -- --run`。
- [x] 7.3 运行 `npm --prefix web run build`。
- [x] 7.4 运行 `bash scripts/e2e-smoke.sh`。
- [x] 7.5 运行 `openspec validate p46-local-knowledge-import-governance --strict`。
- [x] 7.6 运行 `openspec validate --all --strict`。
- [x] 7.7 运行 `git diff --check`。
- [x] 7.8 运行安全扫描：`rg -n 'sk-[A-Za-z0-9]{12,}|BEGIN (RSA|OPENSSH|PRIVATE) KEY|/Users/[^[:space:]]+|SELECT \\* FROM|raw HTTP|完整 prompt|券商接口|自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|收益承诺|Level2|高频源' web/src/pages/LocalKnowledgePage.tsx web/src/services/localKnowledge.ts web/src/types/localKnowledge.ts internal/application/dto/local_knowledge.go internal/application/service/local_knowledge.go internal/application/handler/local_knowledge_handler.go docs/api.md docs/frontend-contract.md docs/data-model.md`，人工复核命中项，确认不存在未脱敏敏感内容或高风险操作入口；允许安全边界说明文本命中。

## 8. 归档前复审

- [x] 8.1 执行完成后再次只读子 agent 复审，确认无 Critical / Important。
- [x] 8.2 复审通过后执行 archive，并将 P46 归档。
