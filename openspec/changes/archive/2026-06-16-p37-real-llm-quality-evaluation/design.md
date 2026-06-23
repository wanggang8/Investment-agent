## Context

当前 DeepSeek client 已能调用 `/chat/completions`，但模型固定为 `deepseek-chat`，错误均映射为 `ANALYST_UNAVAILABLE`，缺少 prompt 版本、输入/输出摘要、质量状态和真实 smoke 入口。P37 只增强 LLM 分析材料质量，不改变规则裁决、守门人审计或交易边界。

## Goals / Non-Goals

**Goals:**

- 通过本地配置文件支持真实 base URL、API key 与模型名 `gpt-5.4-mini`。
- 建立真实调用 smoke，确认 LLM 可生成分析材料。
- 对 LLM 失败原因进行稳定分类，便于 API、审计和前端展示。
- 对输出质量进行 fixture 评估，拦截收益承诺、确定涨跌、交易指令、覆盖最终裁决等越权内容。
- 记录 prompt version、输入摘要、输出摘要、解析状态、质量状态和脱敏审计信息。

**Non-Goals:**

- 不让 LLM 写最终 verdict。
- 不用 LLM 触发交易、确认、订单、券商 API 或外部推送。
- 不引入登录、付费、授权、Level2 或高频源。
- 不把 API key 明文写入响应、审计事件或需提交的文档。

## Decisions

### 1. 配置文件优先，环境变量仍兼容

P37 按用户要求使用配置文件承载真实 key。实现上扩展 `deepseek.model` 与 timeout/smoke 相关配置，继续兼容现有环境变量覆盖以免破坏旧部署。真实 key 放入本地未跟踪配置文件，不进入提交。

### 2. LLM client 输出结构化调用结果

DeepSeek client 保持 `Analyze` 接口兼容，同时新增或内部维护调用 metadata：model、prompt_version、input_summary、output_summary、parse_status、quality_status、error_category。服务层和测试读取这些稳定字段，避免前端解析原始错误文本。

### 3. 质量评估先用规则化 fixture

质量评估不依赖另一个 LLM，先用本地规则检查输出文本是否包含收益承诺、确定性预测、交易执行语言、最终裁决覆盖等越权信号。这样测试可复现，也符合“LLM 不覆盖最终规则裁决”的边界。

### 4. 真实 smoke 显式但作为本阶段验收

真实 smoke 使用本地配置文件和用户提供的小额度临时 key 调用代理 base URL。smoke 不写账户、持仓、确认、订单或规则版本；只验证 analyst material、质量状态和脱敏审计摘要。

## Risks / Trade-offs

- 真实代理或 key 暂时不可用会导致 smoke 失败 → 保留错误分类，失败可诊断；本阶段仍以真实 smoke 为目标验收。
- 输出质量检查可能误伤中文表达 → 先拦截高风险明确短语，测试覆盖边界，不做复杂语义判断。
- key 写配置文件有泄漏风险 → 本地配置进入 `.gitignore`；提交只包含 example 和文档占位。

## Migration Plan

1. 新增 P37 OpenSpec change 和 delta。
2. 扩展配置模型与 example，创建本地 `configs/config.local.yaml`。
3. 以 TDD 增加 client、质量评估、workflow 降级和真实 smoke 测试。
4. 更新文档与进度；真实 smoke 验收通过后再归档。
