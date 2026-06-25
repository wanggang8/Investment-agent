# Tasks: P118 产品可用性边界场景验收

## 1. Governance

- [x] 1.1 阅读 `docs/GOVERNANCE.md` 与 `openspec/project.md`。
- [x] 1.2 创建 `p118-product-usability-edge-scenario-acceptance` change。
- [x] 1.3 更新 `docs/GOVERNANCE.md`、`openspec/PROGRESS.md`、`openspec/project.md`，标记 P118 为并行验收型 active change。
- [x] 1.4 `openspec validate p118-product-usability-edge-scenario-acceptance --strict`。

## 2. Scenario Matrix

- [x] 2.1 新增 P118 产品可用性边界场景验收矩阵。
- [x] 2.2 覆盖 E01-E18，列出操作、入口、API/browser 证据、SQLite readback、解释性结论和安全负证据。
- [x] 2.3 明确 P118 排除发布/安装/升级，不声称券商成交、外部 provider、fresh real LLM 或收益准确性。

## 3. Runner

- [x] 3.1 新增 `scripts/p118-product-usability-edge-scenario-acceptance.sh`，使用临时 config、临时 SQLite、动态端口、backend/Vite lifecycle、backend restart 和 cleanup trap。
- [x] 3.2 新增 `scripts/p118_product_usability_edge_scenario_acceptance.py`，实现 API/SQLite runner、长期数据 seed、restart persistence probe 和 merge-only 模式。
- [x] 3.3 新增 `web/e2e/p118-product-usability-edge-scenario-acceptance.spec.ts`，覆盖关键浏览器路径和截图。
- [x] 3.4 Runner 输出 API/SQLite summary、restart summary、browser summary、final merged product-usability report 和截图证据。

## 4. Edge Scenario Coverage

- [x] 4.1 30 天本地使用耐久：日报、审计、通知、风险、交易历史积累后仍可读。
- [x] 4.2 异常输入恢复：坏导入、非法交易、重复/冲突交易不产生半写入。
- [x] 4.3 数据源波动：stale/missing/恢复状态以 scoped resolution 表达，不冒充 clean pass。
- [x] 4.4 决策质量解释：上涨、下跌、震荡三类上下文建议有差异且可追溯。
- [x] 4.5 多账户/家庭账本：多个本地账户标签、多基金、多现金/货币基金事实跨页一致。
- [x] 4.6 长列表/历史数据：多交易、多审计、多通知、多报告页面可打开且无 5xx/console/page errors。
- [x] 4.7 移动端核心路径：390px 下组合、工作台、决策闭环关键页面可用。

## 5. Safety

- [x] 5.1 SQLite broker/order/push 相关表不存在或计数为 0。
- [x] 5.2 自动确认记录为 0。
- [x] 5.3 自动规则应用审计事件为 0。
- [x] 5.4 前端无自动交易、一键交易、代下单、外部推送、收益承诺 affordance。
- [x] 5.5 首层 UI 无敏感 key、raw prompt 泄露。

## 6. Evidence

- [x] 6.1 执行 P118 runner 并生成 summary JSON、截图、SQLite readback。
- [x] 6.2 新增 P118 acceptance record。
- [x] 6.3 记录 P93 stale 边界，不伪称 fresh P93 pass。

## 7. Regression Gates

- [x] 7.1 `bash scripts/p118-product-usability-edge-scenario-acceptance.sh`。
- [x] 7.2 `openspec validate p118-product-usability-edge-scenario-acceptance --strict`。
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
