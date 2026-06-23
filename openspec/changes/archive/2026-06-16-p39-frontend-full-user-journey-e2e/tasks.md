## 1. OpenSpec 与范围

- [x] 1.1 确认 P39 只覆盖前端完整用户旅程、浏览器级 E2E、降级展示、窄屏/a11y/console 检查。
- [x] 1.2 确认 P39 不接券商 API、不自动交易、不外部推送、不自动应用规则、不绕过守门人或用户最终确认。
- [x] 1.3 对齐 P33 onboarding、P34 source health、P35 risk alert、P36/P37 governance/LLM、P38 retrieval quality 的既有契约。

## 2. E2E Fixture 与运行入口

- [x] 2.1 建立临时 SQLite + 本地配置 + deterministic seed 的 Playwright fixture。
- [x] 2.2 fixture 覆盖账户、持仓、市场、证据、retrieval quality、daily report、risk alert、rule proposal、audit。
- [x] 2.3 E2E 不依赖公网、不含真实密钥、不读取本地 SQLite/VecLite 文件内容。
- [x] 2.4 明确本地运行命令、端口选择、失败清理和 gitignore 边界。

## 3. 完整用户旅程

- [x] 3.1 覆盖空库首次进入、缺前提引导和账户/持仓初始化入口。
- [x] 3.2 覆盖 Dashboard、Portfolio、Evidence、Daily Discipline Report、Decision detail、Audit、Review、Rules、Risk Alerts、Settings 的关键可达路径。
- [x] 3.3 覆盖主动咨询、用户线下确认记录、错误标注或人工复核入口。
- [x] 3.4 覆盖复盘历史、规则提案状态、守门人/最终确认边界。

## 4. 降级与安全边界

- [x] 4.1 覆盖缺账户、缺市场、证据不足、VecLite/RAG 降级、LLM 降级、能力圈外、规则提案待确认。
- [x] 4.2 断言页面不出现自动下单、一键交易、券商接口、自动规则应用或收益承诺入口。
- [x] 4.3 确认 P34 source health、P35 risk alert、P38 retrieval quality 降级信息只读展示。
- [x] 4.4 关键用户动作仅写本地事实或确认记录，不生成交易执行语义。

## 5. 可用性与稳定性

- [x] 5.1 增加关键页面 console error / unhandled rejection 捕获。
- [x] 5.2 增加窄屏 viewport smoke，确认关键文本和按钮不重叠。
- [x] 5.3 增加基础可访问性检查：可聚焦控件、表单 label、导航 landmark 或等价可用性断言。
- [x] 5.4 保持 Vitest 与 Playwright 分层，避免互相收集或重复启动冲突。

## 6. 文档与验收

- [x] 6.1 在 P39 delta 中记录待归档合并到 `docs/frontend-contract.md` 的完整旅程、降级、安全边界契约。
- [x] 6.2 在 P39 delta 中记录待归档合并到 `docs/configuration.md` / `docs/ops-local-scheduler.md` 的 E2E 运行说明。
- [x] 6.3 更新 `docs/development-plan.md`、`openspec/PROGRESS.md`、`AGENTS.md`、`docs/GOVERNANCE.md` 的 P39 active 状态。
- [x] 6.4 运行 `go test ./...`。
- [x] 6.5 运行 `npm --prefix web test -- --run`。
- [x] 6.6 运行 `npm --prefix web run build`。
- [x] 6.7 运行 Playwright E2E / browser smoke。
- [x] 6.8 运行 archive 前只读子 agent 复审，且无 Critical / Important 问题。
- [x] 6.9 运行 `openspec validate p39-frontend-full-user-journey-e2e --strict`。
- [x] 6.10 运行 `openspec validate --all --strict`。
