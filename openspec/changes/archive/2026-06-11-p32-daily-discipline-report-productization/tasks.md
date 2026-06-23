# P32 每日纪律报告产品化任务

## 1. OpenSpec 与范围

- [x] 1.1 确认 P32 是 P31 后的报告产品化层，只聚合和展示 daily workflow / auto-run 结果，不改变交易、安全或数据源边界。
- [x] 1.2 定义报告状态、幂等键、缺前提语义和安全边界：本地低频、人工复核、不交易、不外推、不新增登录/付费/高频源。

## 2. 后端模型与持久化

- [x] 2.1 新增 `daily_discipline_reports` migration，包含本地日期、scope、status、summary、missing prerequisites、关联记录、idempotency key 和时间戳。
- [x] 2.2 实现 report repository，支持按本地日期/scope upsert 或 reuse、today 查询、列表分页和详情查询。
- [x] 2.3 完成 repository wiring、transactor 集成和持久化测试。

## 3. 聚合 API

- [x] 3.1 定义 daily discipline report DTO，覆盖 today/list/detail、状态、摘要、缺前提、关联 decision/evidence/audit/notification 链接。
- [x] 3.2 实现 today API，成功时返回今日报告，缺前提时返回结构化 missing prerequisites 状态。
- [x] 3.3 实现 list/detail API，支持历史报告列表和详情回看。
- [x] 3.4 增加 handler/service tests，覆盖 successful report、missing prerequisites、history/detail 和幂等复用。

## 4. 前端产品化

- [x] 4.1 新增前端 types/service，封装 today/list/detail API。
- [x] 4.2 升级今日纪律页，展示今日报告状态、摘要、缺前提、关联材料和人工复核边界。
- [x] 4.3 新增历史报告列表页和报告详情页。
- [x] 4.4 更新路由导航，使今日纪律、历史报告和详情入口可达。
- [x] 4.5 增加前端 tests，覆盖成功报告、缺前提、历史列表和详情。

## 5. Smoke 与文档

- [x] 5.1 更新 smoke seed，提供今日报告、历史报告和 P31 缺前提状态样例数据。
- [x] 5.2 扩展 E2E smoke，验证本地 UI 可展示今日纪律报告、历史报告列表和详情。
- [x] 5.3 更新 docs/progress，记录 P32 范围、验收命令和安全边界。

## 6. 验收

- [x] 6.1 运行 `go test ./...`。
- [x] 6.2 运行 `npm --prefix web test -- --run`。
- [x] 6.3 运行 `npm --prefix web run build`。
- [x] 6.4 运行 P32 E2E smoke 命令。
- [x] 6.5 运行 `openspec validate p32-daily-discipline-report-productization --strict`。
- [x] 6.6 运行 `git status --short`，确认只包含预期修改且无临时产物。
- [x] 6.7 执行只读复审，确认无 Critical / Important 问题；复审 findings 已修复并重新验证。
