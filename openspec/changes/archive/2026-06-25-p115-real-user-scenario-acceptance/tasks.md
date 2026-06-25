# Tasks: P115 真实用户场景全链路验收

## 1. Governance

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 创建 `p115-real-user-scenario-acceptance` change。
- [x] 1.3 更新 `docs/GOVERNANCE.md`、`openspec/PROGRESS.md` 与 `openspec/project.md`，标记 P115 为与 P114 并行的验收型 active change。
- [x] 1.4 校验 `openspec validate p115-real-user-scenario-acceptance --strict`。

## 2. Scenario Matrix

- [x] 2.1 新增 P115 场景矩阵文档，覆盖 S01-S33。
- [x] 2.2 为每个场景标明入口、操作、API/browser 证据、SQLite readback、下游联动、安全负证据。
- [x] 2.3 标明 expected eligibility 与 actual status；actual status 默认 `pending`，执行后才可变为 `fresh_pass`、`scoped_pass`、`degraded_expected`、`blocked`。
- [x] 2.4 明确 P115 不新增运行时投资能力，不处理发布包/安装器/物理第二机器。

## 3. Runner Expansion

- [x] 3.1 复用 P104 runner 的临时 config、临时 SQLite、backend lifecycle 和 forbidden-table checks，但所有 P104-derived evidence 必须标记为 `local_seeded_linkage`，不得用于 provider/LLM fresh claim。
- [x] 3.2 新增 P115 artifact schema，记录 `scenario_id`、`title`、`status`、`expected_eligibility`、`classification_reason`、`config_mode`、`runtime_mode`、`use_stub`、`provider_mode`、`llm_mode`。
- [x] 3.3 artifact schema 必须记录 API method/path/status/request_id、browser route/viewport/screenshot/DOM/console、SQLite table/field/before/after/row_count、downstream endpoint/page、side-effect ids、redaction result 和 safety counters。
- [x] 3.4 按 `api_sqlite`、`browser`、`degradation` 三层分别输出 log/summary，使用 P115 专属临时目录、动态端口或 P115 env vars、trap cleanup backend/Vite/Playwright。
- [x] 3.5 扩展 API/SQLite 场景：S03-S23、S26-S31。
- [x] 3.6 扩展 browser 场景：S01-S05、S09-S19、S21-S29、S32-S33。
- [x] 3.7 扩展 degradation 场景：S05、S09、S13-S17、S20、S22-S23、S28、S30-S31。

## 4. Functional Coverage

- [x] 4.1 验收首次启动、本地能力边界、空账户引导。
- [x] 4.2 验收组合初始化、持仓新增/编辑/删除、批量导入、线下交易、本地事实修正、季度再平衡。
- [x] 4.3 验收主动咨询、决策详情、人工计划确认、决策错误标注 `marked_error`、决策闭环。
- [x] 4.4 验收证据刷新、证据验证、RAG/VecLite 重建、知识准备度、本地知识导入。
- [x] 4.5 验收市场刷新、source health、数据质量回归、gate resolution 创建/退休。
- [x] 4.6 验收风险预警 SOP、规则提案、规则效果验证、通知、日报、自动运行只读状态、复盘、审计、设置、settings 禁止规则/SOP 直接修改、API 诊断。
- [x] 4.7 验收 390px 移动端核心操作路径。
- [x] 4.8 验收失败/降级/非法输入/不存在 id/缺 key/无证据/索引不可用，并区分 `use_stub=false` provider/LLM 失败和 `local_seeded_linkage` 证据。

## 5. Safety Boundary

- [x] 5.1 检查 SQLite 中 broker/order/push 相关表不存在或记录为 0。
- [x] 5.2 检查自动确认记录为 0。
- [x] 5.3 检查自动规则应用审计事件为 0。
- [x] 5.4 检查前端无自动交易、一键交易、代下单、外部推送、收益承诺 affordance。
- [x] 5.5 检查敏感 key、prompt payload、raw secret、本机路径不在首层 UI 泄露。
- [x] 5.6 检查 settings API 不能直接修改规则阈值或 SOP 配置；拒绝后无 rule version / rule proposal / audit auto-apply 副作用。

## 6. Evidence And Review

- [x] 6.1 执行 P115 runner，生成 summary JSON、runner log、截图和 SQLite evidence。
- [x] 6.2 新增 P115 acceptance record，逐场景列出 status 和证据路径。
- [x] 6.3 对 `blocked` 或 `degraded_expected` 场景记录原因和不扩大声明边界；外部 provider/LLM 缺 key、网络失败或只有 stub/local seed 时不得写成 `fresh_pass`。
- [x] 6.4 复审 P115 结果，确认没有 release-blocking functional fake / backend missing / safety boundary failure。
- [x] 6.5 acceptance record 明确 P93 如仍 stale 则不得声明 fresh P93 pass，只能声明 P115 current-source functional reality evidence。

## 7. Regression Gates

- [x] 7.1 `openspec validate p115-real-user-scenario-acceptance --strict`。
- [x] 7.2 `go test ./...`。
- [x] 7.3 `go vet ./...`。
- [x] 7.4 `npm --prefix web test -- --run`。
- [x] 7.5 `npm --prefix web run build`。
- [x] 7.6 `openspec validate --all --strict`。
- [x] 7.7 `python3 scripts/p92_final_requirement_audit.py --check`。
- [x] 7.8 P114 后代码真实性等价检查；若 P93 仍 stale，记录 stale 原因并新增 P115 current-source functional reality evidence，不伪称 P93 fresh pass。
- [x] 7.9 `git diff --check`。

## 8. Archive

- [x] 8.1 更新 governance/progress，记录 P115 fresh scenario acceptance 结论。
- [x] 8.2 OpenSpec archive。
- [ ] 8.3 最终验证后提交。
