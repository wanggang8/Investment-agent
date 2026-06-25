# Tasks: P117 连续产品可用性验收

## 1. Governance

- [x] 1.1 阅读 `docs/GOVERNANCE.md` 与 `openspec/project.md`。
- [x] 1.2 创建 `p117-continuous-product-usability-acceptance` change。
- [x] 1.3 更新 `docs/GOVERNANCE.md`、`openspec/PROGRESS.md`、`openspec/project.md`，标记 P117 为并行验收型 active change。
- [x] 1.4 `openspec validate p117-continuous-product-usability-acceptance --strict`。

## 2. Scenario Matrix

- [x] 2.1 新增 P117 七天连续使用验收矩阵。
- [x] 2.2 覆盖 U01-U18，列出每日操作、入口、API/browser 证据、SQLite readback、解释性结论和安全负证据。
- [x] 2.3 明确 P117 是 local seeded continuous-use acceptance，不声称券商成交、外部 provider、fresh real LLM 或发布包验收。

## 3. Runner

- [x] 3.1 新增 `scripts/p117-continuous-product-usability-acceptance.sh`，使用临时 config、临时 SQLite、动态端口、backend/Vite lifecycle、backend restart 和 cleanup trap。
- [x] 3.2 新增 `scripts/p117_continuous_product_usability_acceptance.py`，实现 API/SQLite runner、restart persistence probe 和 merge-only 模式。
- [x] 3.3 新增 `web/e2e/p117-continuous-product-usability-acceptance.spec.ts`，覆盖关键浏览器路径和截图。
- [x] 3.4 Runner 输出 API/SQLite summary、restart summary、browser summary、final merged usability report 和截图证据。

## 4. Seven-Day Coverage

- [x] 4.1 Day 0 冷启动：空组合、空状态、无假数据。
- [x] 4.2 Day 1 入门：录入账户/持仓，读回组合事实。
- [x] 4.3 Day 2 日常：今日纪律、工作台、复盘、审计聚合读回。
- [x] 4.4 Day 3 交易补记：多基金线下 buy/sell/reduce 和风险/通知处理。
- [x] 4.5 Day 4 错误恢复：坏批量导入、非法交易拒绝、修正审计，不产生半写入。
- [x] 4.6 Day 5 数据质量：降级 gate resolution 创建/退役，文案不冒充 clean pass。
- [x] 4.7 Day 6 决策闭环：人工执行确认、marked_error、错误复盘。
- [x] 4.8 Day 7 收口：跨页面一致性、审计、决策闭环、重启后数据仍可读。

## 5. Usability Interpretation

- [x] 5.1 输出任务完成率与阻断项。
- [x] 5.2 输出用户可理解性结论：下一步是否清楚、错误是否可恢复、解释是否可追溯。
- [x] 5.3 输出跨页面一致性结论：同一事实在持仓/今日纪律/工作台/复盘/审计/闭环中不矛盾。
- [x] 5.4 输出上线边界：本地工具可用性与不承诺范围。

## 6. Safety

- [x] 6.1 SQLite broker/order/push 相关表不存在或计数为 0。
- [x] 6.2 自动确认记录为 0。
- [x] 6.3 自动规则应用审计事件为 0。
- [x] 6.4 前端无自动交易、一键交易、代下单、外部推送、收益承诺 affordance。
- [x] 6.5 首层 UI 无敏感 key、raw prompt 泄露。

## 7. Evidence

- [x] 7.1 执行 P117 runner 并生成 summary JSON、截图、SQLite readback。
- [x] 7.2 新增 P117 acceptance record。
- [x] 7.3 记录 P93 stale 边界，不伪称 fresh P93 pass。

## 8. Regression Gates

- [x] 8.1 `bash scripts/p117-continuous-product-usability-acceptance.sh`。
- [x] 8.2 `openspec validate p117-continuous-product-usability-acceptance --strict`。
- [x] 8.3 `go test ./...`。
- [x] 8.4 `go vet ./...`。
- [x] 8.5 `npm --prefix web test -- --run`。
- [x] 8.6 `npm --prefix web run build`。
- [x] 8.7 `openspec validate --all --strict`。
- [x] 8.8 `python3 scripts/p92_final_requirement_audit.py --check`。
- [x] 8.9 `python3 scripts/p93_code_reality_audit.py --check`，结果为 stale：`docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`。
- [x] 8.10 `git diff --check`。

## 9. Archive

- [x] 9.1 用户确认后再 archive；本轮不自动归档。
