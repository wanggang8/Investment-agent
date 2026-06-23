## 1. OpenSpec 与范围

- [x] 1.1 确认 P34 只覆盖真实公开数据覆盖扩展、数据源健康、失败分类、工作流输入上下文和前端状态展示。
- [x] 1.2 确认 P34 不接券商 API、不自动交易、不外部推送、不登录/付费/授权/Level2/高频源、不承诺收益、不预测确定涨跌。
- [x] 1.3 对齐 P26/P27/P29 collector 基线、P33 本地账户事实、`docs/api.md`、`docs/data-model.md`、`docs/workflow.md`、`docs/frontend-contract.md` 的现有契约。

## 2. 数据源调研与 source contract

- [x] 2.1 梳理中证指数样本、权重、估值文件的当前公开 endpoint shape、字段、日期、分页/文件边界和失败类型。
- [x] 2.2 评估成分股财务、资金流向、融资融券或可替代情绪指标的公开可用源，记录 source level、刷新频率和安全边界。
- [x] 2.3 为每类 P34 数据定义 normalized payload：source_name、source_level、source_type、data_category、symbol/index、data_date、captured_at、content_hash、metrics、raw metadata。
- [x] 2.4 定义 freshness 与失败分类：fresh、stale、missing、no_data、source_unavailable、parse_error、disabled、stubbed。

## 3. 后端 collector 与持久化

- [x] 3.1 扩展或新增 P34 collector 接口和实现，覆盖首批可验证的指数样本/权重/估值及至少一个财务、资金、融资融券或情绪替代指标类别。
- [x] 3.2 为新增 collector 增加 fixture/stub，默认本地验收不依赖公网。
- [x] 3.3 将标准化数据写入现有 market/evidence/审计路径；若现有模型不足，新增轻量 migration delta 保存 source health 或扩展指标 metadata。
- [x] 3.4 实现幂等去重：按 source identity、symbol/index、data_date、file identity 或 content_hash 避免重复事实。
- [x] 3.5 实现 source health 记录：最近成功/失败时间、失败类别、数据日期、影响标的、source level。
- [x] 3.6 增加 repository/service tests，覆盖成功写入、no_data、source_unavailable、parse_error、stale、stub fallback 和幂等。

## 4. 工作流与任务入口

- [x] 4.1 扩展 `cmd/agent` 或本地任务入口，支持显式触发 P34 扩展数据 refresh，包含 source、symbol/index 和日期窗口参数。
- [x] 4.2 将 P34 数据摘要接入 DailyDisciplineGraph 输入上下文，缺失或过期时保留 missing/stale 诊断。
- [x] 4.3 将 P34 数据摘要接入 ExpectedReturnNode 上下文，保留样本限制、source level 和 freshness，不输出收益承诺。
- [x] 4.4 确认 P34 输出不会写 positions、portfolio snapshots、operation confirmations、broker state、orders、external notifications 或 rule_versions。

## 5. API 与前端状态展示

- [x] 5.1 扩展数据源健康 API 或复用 settings/ops DTO，返回 P34 source category、freshness、last success/failure、failure category、data_date 和 affected symbols。
- [x] 5.2 更新 Settings / Ops / Daily Discipline 相关页面，展示 P34 数据 fresh、stale、missing、unavailable、parse-error、disabled、stubbed 状态。
- [x] 5.3 在每日纪律报告中展示 P34 数据覆盖状态，缺失或降级时说明哪些类别不足。
- [x] 5.4 增加前端 tests，覆盖健康状态展示、缺失/过期诊断和禁止自动交易文案。

## 6. 文档与验收

- [x] 6.1 在 P34 delta 中记录待归档合并到 `docs/api.md` 的 refresh/API/DTO/错误分类和事务边界。
- [x] 6.2 在 P34 delta 中记录待归档合并到 `docs/data-model.md` 的新增或复用数据模型、source health、freshness 和失败分类。
- [x] 6.3 在 P34 delta 中记录待归档合并到 `docs/workflow.md` 与 `docs/frontend-contract.md` 的工作流输入和前端状态展示。
- [x] 6.4 更新 `docs/development-plan.md` 和 `openspec/PROGRESS.md` 的 P34 状态与验收命令。
- [x] 6.5 运行 `go test ./...`。
- [x] 6.6 运行 `npm --prefix web test -- --run`。
- [x] 6.7 运行 `npm --prefix web run build`。
- [x] 6.8 运行 P34 fixture/stub smoke 和至少一个显式启用的真实公开源 smoke，真实源不可用时按分类输出可解释结果。
- [x] 6.9 运行 `openspec validate p34-real-data-coverage-expansion --strict`。
- [x] 6.10 运行 `openspec validate --all --strict`。
- [x] 6.11 运行 `git status --short`，确认只包含预期修改且无临时产物。
