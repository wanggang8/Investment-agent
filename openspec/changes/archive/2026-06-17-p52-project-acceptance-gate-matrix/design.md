# P52 设计：项目验收门禁矩阵

## 设计目标

P52 把已有验证入口整理为发布前验收矩阵，解决“怎么保证项目验收通过”的判断口径。它不运行完整验收，也不新增自动化；P53 发布候选材料必须引用 P52 的门禁结果。

## 验收矩阵结构

新增 `docs/project-acceptance-gate-matrix.md`，建议结构：

1. **总览**
   - 适用范围。
   - 与 P51/P53 的关系。
   - 发布阻断规则。
2. **门禁分层**
   - G0 治理和文档一致性。
   - G1 Go 单元测试。
   - G2 Go 集成和工作流测试。
   - G3 前端单元/页面测试与构建。
   - G4 E2E / Playwright smoke。
   - G5 本地 fixture/current smoke。
   - G6 真实公开源 opt-in 测试。
   - G7 真实 LLM opt-in 测试。
   - G8 本地安装、备份恢复和发布升级检查。
   - G9 安全边界和脱敏检查。
3. **门禁表字段**
   - Gate ID。
   - 目标。
   - 命令或入口。
   - 前置条件。
   - 通过标准。
   - 允许降级。
   - 产物位置。
   - 是否阻断发布。
4. **真实测试策略**
   - 必须显式 opt-in。
   - 必须使用临时配置或明确测试 key。
   - 必须分类失败：网络、限流、认证/凭证、源 schema 变化、无数据、解析失败、模型不可用、质量不达标。
   - 不得把真实测试通过解释为收益承诺、未来可用性承诺或交易能力。
5. **验收记录格式**
   - 建议每次验收输出 `docs/release/acceptance/YYYY-MM-DD-<label>.md` 或等价 release notes 附件。
   - P52 只定义格式，不创建实际验收结果。

## 建议门禁

| Gate | 分类 | 示例命令 | 发布阻断 |
| --- | --- | --- | --- |
| G0 | 治理和 OpenSpec | `openspec validate --all --strict`、`git diff --check` | 是 |
| G1 | Go 单元 | `go test ./...` | 是 |
| G2 | Go 聚焦集成 | `go test ./cmd/agent ./cmd/server ./internal/application/workflow ./internal/application/handler ./internal/infrastructure/persistence/sqlite` | 是 |
| G3 | 前端测试和构建 | `npm --prefix web test -- --run`、`npm --prefix web run build` | 是 |
| G4 | E2E smoke | `bash scripts/e2e-smoke.sh` | 是，除非缺本机浏览器依赖且记录为环境跳过 |
| G5 | 本地 smoke | `bash scripts/recovery-smoke.sh`、`go run ./cmd/agent --task data-source-quality-regression --source fixture --symbol 000300` | 是 |
| G6 | 真实公开源 opt-in | `go run ./cmd/agent --task public-evidence-refresh --symbol <symbol> --start-date YYYY-MM-DD --end-date YYYY-MM-DD` | 条件阻断：失败需分类，非本地代码问题可降级但必须记录 |
| G7 | 真实 LLM opt-in | `go run ./cmd/agent --task llm-smoke --symbol 510300` | 条件阻断：质量失败阻断 LLM 发布声明，网络/额度/模型不可用需分类 |
| G8 | 本地安装与升级 | `bash scripts/local-install-diagnostics.sh --include-release-upgrade`、`bash scripts/local-release-upgrade-check.sh --target-version vNEXT --output-dir <tmp>/release-upgrade` | 是 |
| G9 | 安全/脱敏 | `rg` 禁止能力和敏感信息扫描 + 人工复核 | 是 |

## 验收状态口径

- `pass`：命令通过，产物存在，未发现阻断问题。
- `degraded`：真实源、LLM 或环境依赖不稳定，但失败已分类，且不影响本地核心能力；是否阻断由矩阵定义。
- `blocked`：发布阻断门禁失败，或出现安全边界问题。
- `skipped`：仅允许非必需/显式 opt-in 门禁；必须记录理由。

## 安全边界

P52 不新增运行时能力。矩阵必须继续禁止券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
