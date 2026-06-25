# Tasks: P116 多基金交易账本复杂场景验收

## 1. Governance

- [x] 1.1 阅读 `docs/GOVERNANCE.md` 与 `openspec/project.md`。
- [x] 1.2 创建 `p116-multi-fund-transaction-ledger-acceptance` change。
- [x] 1.3 更新 `docs/GOVERNANCE.md`、`openspec/PROGRESS.md`、`openspec/project.md`，标记 P116 为并行验收型 active change。
- [x] 1.4 `openspec validate p116-multi-fund-transaction-ledger-acceptance --strict`。

## 2. Scenario Matrix

- [x] 2.1 新增 P116 多基金复杂交易账本验收矩阵。
- [x] 2.2 覆盖 L01-L16，列出入口、操作、API/browser 证据、SQLite readback、下游联动、安全负证据。
- [x] 2.3 明确 P116 使用 `local_seeded_linkage`，不声称真实券商成交、外部 provider 或 fresh real LLM。

## 3. Runner

- [x] 3.1 新增 `scripts/p116-multi-fund-transaction-ledger-acceptance.sh`，使用临时 config、临时 SQLite、动态端口、backend/Vite lifecycle 和 cleanup trap。
- [x] 3.2 新增 `scripts/p116_multi_fund_transaction_ledger_acceptance.py`，实现 API/SQLite runner 与 merge-only 模式。
- [x] 3.3 新增 `web/e2e/p116-multi-fund-transaction-ledger-acceptance.spec.ts`，覆盖核心浏览器路径和截图。
- [x] 3.4 Runner 输出 API/SQLite summary、browser summary、final merged summary 和截图证据。

## 4. Complex Ledger Coverage

- [x] 4.1 验收多基金初始组合：`510300`、`159915`、`588000`、`512000`、`110022`。
- [x] 4.2 验收多日期线下 buy/sell/reduce 记录、费用、现金变化和持仓变化。
- [x] 4.3 验收混合批量导入：holding row、transaction row、合法行、非法行；validate 不写库，confirm 只写合法数据。
- [x] 4.4 验收非法交易拒绝：现金不足、超持仓卖出、未来执行时间、负费用、缺 symbol、非法 position_state。
- [x] 4.5 验收持仓编辑/移除、本地事实修正和季度再平衡。
- [x] 4.6 验收决策详情人工执行确认、marked_error、决策闭环、审计、复盘联动。
- [x] 4.7 验收风险、通知、数据质量处置、首页、工作台、日报/复盘/审计聚合读回。
- [x] 4.8 验收桌面与 390px 移动端组合页面。

## 5. Safety

- [x] 5.1 SQLite broker/order/push 相关表不存在或计数为 0。
- [x] 5.2 自动确认记录为 0。
- [x] 5.3 自动规则应用审计事件为 0。
- [x] 5.4 前端无自动交易、一键交易、代下单、外部推送、收益承诺 affordance。
- [x] 5.5 首层 UI 无敏感 key、raw prompt、本机路径泄露。

## 6. Evidence

- [x] 6.1 执行 P116 runner 并生成 summary JSON、截图、SQLite readback。
- [x] 6.2 新增 P116 acceptance record。
- [x] 6.3 记录 P93 stale 边界，不伪称 fresh P93 pass。

## 7. Regression Gates

- [x] 7.1 `bash scripts/p116-multi-fund-transaction-ledger-acceptance.sh`。
- [x] 7.2 `openspec validate p116-multi-fund-transaction-ledger-acceptance --strict`。
- [x] 7.3 `go test ./...`。
- [x] 7.4 `go vet ./...`。
- [x] 7.5 `npm --prefix web test -- --run`。
- [x] 7.6 `npm --prefix web run build`。
- [x] 7.7 `openspec validate --all --strict`。
- [x] 7.8 `python3 scripts/p92_final_requirement_audit.py --check`。
- [x] 7.9 `python3 scripts/p93_code_reality_audit.py --check`，结果为 stale：`docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`。
- [x] 7.10 `git diff --check`。

## 8. Archive

- [x] 8.1 用户确认后再 archive；本轮不自动归档。
