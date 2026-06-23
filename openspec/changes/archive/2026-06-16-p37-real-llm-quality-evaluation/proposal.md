## Why

P7 已接入 DeepSeek 分析服务并确保 LLM 不写最终裁决，但当前能力仍偏“可调用/可降级”：模型名硬编码、错误分类粗、prompt 与输出缺少结构化质量记录，真实调用 smoke 也没有产品化。P37 需要把 LLM 路径从“有 client”升级为“真实配置可验证、失败可诊断、输出质量可评估、敏感输入可审计”。

本轮按用户要求以真实调用作为 P37 验收路径：使用本地配置文件中的 base URL、API key 与模型 `gpt-5.4-mini` 执行一次低成本 smoke；无 key、超时、空响应、格式错误和质量不足仍必须有降级测试，确保真实源不可用时规则裁决继续运行。

## What Changes

- 扩展 DeepSeek/LLM 配置：支持 `api_key`、`base_url`、`model`、timeout 和 smoke 开关；真实测试配置写入本地未跟踪配置文件。
- 扩展 LLM client：模型名可配置，错误分类区分 missing_key、timeout、http_error、empty_response、parse_error、quality_failed 等。
- 增加 prompt 版本、输入摘要、输出摘要、解析状态和质量检查结果；审计记录不得包含明文 API key 或不必要敏感信息。
- 增加 LLM 输出质量评估 fixture：检查收益承诺、确定性涨跌预测、直接交易指令、覆盖最终裁决等越权内容。
- 增加真实 LLM smoke：显式读取本地配置，调用 `gpt-5.4-mini` 生成分析材料，并确认最终裁决仍由规则引擎负责。
- 更新 OpenSpec 与 L1 文档 delta，归档时同步到 `docs/api.md`、`docs/data-model.md`、`docs/workflow.md`、`docs/frontend-contract.md` 与配置说明。

## Capabilities

### New Capabilities
- `real-llm-quality-evaluation`: 真实 LLM 配置 smoke、prompt/输出摘要审计、质量评估 fixture、错误分类和敏感信息边界。

### Modified Capabilities
- `real-data-integration`: DeepSeek analyst material 增强为可配置模型、真实调用 smoke、结构化错误分类、prompt 版本与质量记录。
- `product-completeness`: 继续保证 LLM 输出不能写最终裁决、不能触发交易动作。

## Impact

- 后端：扩展 config、DeepSeek client、analyst service 或 workflow 审计字段；新增 smoke task 或测试入口。
- 数据模型/API：记录 LLM 调用摘要、错误分类、质量状态和安全边界；不得返回或审计明文 key。
- 文档：更新配置说明与 P37 契约，明确本地真实 key 配置文件边界。
- 验收：`go test ./...`、前端测试/构建、OpenSpec 校验，以及真实 LLM smoke。
