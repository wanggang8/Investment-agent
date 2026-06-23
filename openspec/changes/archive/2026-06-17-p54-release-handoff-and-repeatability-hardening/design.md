# P54 设计：发布交付与可重复性加固

## 设计目标

P54 将 P53 的验收结果转化为交付说明和复验规则。它不改变 `release_ready` 结论，也不重新运行验收；它让后续用户或 agent 能按一致规则复验、解释重试、处理 current degraded，并知道哪些声明仍然禁止。

## 新增文档

### `docs/release/README.md`

职责：

- 列出本次发布材料入口。
- 指向 P53 acceptance run、release candidate、P54 handoff 和 repeatability 文档。
- 明确临时产物不提交，真实 key 不进入 release 文档。

### `docs/release/release-handoff-2026-06-17.md`

职责：

- 面向用户交付本次 release candidate。
- 汇总状态：`release_ready`。
- 列出已经通过的验收区域。
- 列出不阻断但必须知道的事项：
  - G5 current data-source quality degraded。
  - G3/G4/G8 初次本机进程 kill 后原命令 retry pass。
  - G6 临时配置缺前置项后修正 retry pass。
- 提供本地复验命令入口。
- 明确不可声明事项：收益承诺、自动交易、未来源可用性、未来模型可用性、登录/付费/授权/Level2/高频源覆盖。
- 建议下一阶段 P55 或后续仅在需要发放安装包/版本标签时再创建。

### `docs/release/acceptance-repeatability.md`

职责：

- 固化 P53 后的复验规则。
- 说明可重复执行的目录结构，例如 `tmp/acceptance/<label>/`。
- 定义重试规则：
  - 仅允许在本机资源 kill、端口瞬时占用、浏览器服务启动瞬时失败时做一次原命令 retry。
  - retry 必须保留原失败记录和 retry 记录。
  - 第二次失败即 `blocked`。
- 定义 G5 current degraded 规则：
  - fixture pass + current degraded + failed=0 可作为非阻断 degraded。
  - current failed>0 或 fixture failed 必须阻断。
- 定义 G6 配置规则：
  - 真实模式 `use_stub=false` 需要满足 market prerequisite 或启用 market collectors。
  - public evidence 必须使用临时 SQLite。
- 定义 G7 脱敏规则：
  - key 只能来自本地私有配置或临时配置。
  - release 文档只能写模型名、结果摘要和脱敏 audit summary。
- 定义 release status 决策：
  - 任意未 waived 的 G0-G5/G8/G9 blocked -> `release_blocked`。
  - 任意 redaction failure -> `release_blocked`。
  - G6/G7 失败时不能声明对应真实能力通过；是否阻断由 release materials 明确记录。

## 安全边界

P54 是文档和治理加固，不新增运行时能力，不改变验收命令，不新增自动修复或自动交易语义。

## 验证策略

- OpenSpec 当前 change 严格校验。
- OpenSpec 全量严格校验。
- `git diff --check`。
- `rg` 扫描 P54 与 release 文档，不允许完整 key、私有路径、raw payload、完整 prompt 或原始 SQL。
- 子 agent 复审计划、执行结果和提交前状态。
