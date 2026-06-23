## 1. OpenSpec 与范围

- [x] 1.1 确认 P37 只覆盖真实 LLM 配置、调用 smoke、prompt/输出摘要、错误分类、质量评估和敏感信息边界。
- [x] 1.2 确认 P37 不自动交易、不写最终裁决、不覆盖守门人审计、不自动应用规则、不接券商 API、不外部推送、不承诺收益、不预测确定涨跌。
- [x] 1.3 对齐 P7/P28 analyst reports、expected return、`docs/api.md`、`docs/workflow.md`、`docs/data-model.md`、`docs/frontend-contract.md` 的既有契约。

## 2. 配置与本地真实 key

- [x] 2.1 扩展 DeepSeek/LLM 配置：api_key、base_url、model、timeout 和真实 smoke 开关或等价入口。
- [x] 2.2 更新 `configs/config.example.yaml`，不写真实 key。
- [x] 2.3 新增或更新 `.gitignore`，确保本地真实配置文件不会被提交。
- [x] 2.4 写入本地 `configs/config.local.yaml`，包含用户提供的 base URL、临时 key 和模型 `gpt-5.4-mini`。
- [x] 2.5 配置校验不得在 key 缺失时阻断常规本地运行；真实 smoke 缺 key 时必须给出明确诊断。

## 3. LLM client 与质量评估

- [x] 3.1 以 TDD 扩展 DeepSeek client，支持配置模型名和超时。
- [x] 3.2 以 TDD 增加错误分类：missing_key、timeout、http_error、empty_response、parse_error、quality_failed、unavailable。
- [x] 3.3 以 TDD 增加 prompt version、输入摘要、输出摘要、解析状态和质量状态。
- [x] 3.4 以 TDD 增加输出质量评估 fixture，覆盖收益承诺、确定性涨跌预测、直接交易指令、覆盖最终 verdict 和正常分析材料。
- [x] 3.5 确认响应、日志和审计不得包含明文 API key。

## 4. 工作流、API 与审计接入

- [x] 4.1 将 LLM 调用 metadata 接入 analyst reports 或等价本地事实，不改变最终规则裁决。
- [x] 4.2 将错误分类接入工作流降级状态和 API 错误/状态展示。
- [x] 4.3 将 prompt/version/input/output/quality 摘要写入 audit_events 或等价追踪事实，且脱敏。
- [x] 4.4 增加 workflow/service tests，覆盖无 key、超时/不可用、质量不足和规则裁决仍可运行。

## 5. 真实 LLM smoke

- [x] 5.1 增加真实 LLM smoke 命令或测试入口，显式读取本地配置文件。
- [x] 5.2 使用 `configs/config.local.yaml` 调用真实代理 base URL 和 `gpt-5.4-mini`。
- [x] 5.3 smoke 验证返回 analyst report、prompt version、解析状态、质量状态和安全边界。
- [x] 5.4 smoke 确认不写 positions、portfolio_snapshots、operation_confirmations、position_transactions、rule_versions、orders 或 external notifications。

## 6. 文档与验收

- [x] 6.1 在 P37 delta 中记录待归档合并到 `docs/api.md` 的 LLM 状态/错误分类/DTO/安全边界。
- [x] 6.2 在 P37 delta 中记录待归档合并到 `docs/data-model.md` 的 prompt/输出摘要、质量状态和审计脱敏边界。
- [x] 6.3 在 P37 delta 中记录待归档合并到 `docs/workflow.md`、`docs/frontend-contract.md`、`docs/configuration.md` 的真实调用、降级和 UI 展示。
- [x] 6.4 更新 `docs/development-plan.md`、`openspec/PROGRESS.md`、`AGENTS.md`、`docs/GOVERNANCE.md` 的 P37 active 状态。
- [x] 6.5 运行 `go test ./...`。
- [x] 6.6 运行 `npm --prefix web test -- --run`。
- [x] 6.7 运行 `npm --prefix web run build`。
- [x] 6.8 运行真实 LLM smoke。
- [x] 6.9 运行 `openspec validate p37-real-llm-quality-evaluation --strict`。
- [x] 6.10 运行 `openspec validate --all --strict`。
- [x] 6.11 运行 `git status --short`，确认真实 key 配置文件未被纳入提交。
