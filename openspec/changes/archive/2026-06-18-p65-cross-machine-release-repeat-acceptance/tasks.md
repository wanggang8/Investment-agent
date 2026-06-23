# Tasks: P65 Cross-Machine Release Repeat Acceptance

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 P63/P64 release handoff、P64 package docs、P64 package script、P52 gate matrix 和现有 smoke 脚本。
- [x] 1.3 创建 `p65-cross-machine-release-repeat-acceptance` OpenSpec change。
- [x] 1.4 写明 P65 是跨机器等价的本地隔离复验阶段，不新增运行时业务能力、API、SQLite schema、Eino workflow、LLM 能力或前端产品能力。
- [x] 1.5 更新当前进度文档，标记 P65 active。
- [x] 1.6 运行 `openspec validate p65-cross-machine-release-repeat-acceptance --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 1.7 子 agent 方案复审无 Critical / Important 后执行。

## 2. 隔离复验脚本

- [x] 2.1 新增 `scripts/local-release-repeat-acceptance.sh`，支持 `--archive PATH`、`--output-dir PATH`、`--skip-install`、`--skip-e2e`。
- [x] 2.2 限制输出目录必须位于当前仓库 `tmp/` 下，不写真实用户数据库、私有配置目录、home 目录或 active repo 源码目录。
- [x] 2.3 在解包前调用 `scripts/local-release-package.sh --verify <archive>`，要求 archive 旁存在 sidecar `release-manifest.json`。
- [x] 2.4 解包到 `tmp/local-release-repeat/<timestamp>/workspace/`，从 archive listing 自动识别 package root。
- [x] 2.5 在 extracted package root 内执行 OpenSpec、Go、frontend install/test/build 和 E2E smoke，不依赖 active repo 工作树。
- [x] 2.6 使用临时端口和 extracted package workspace 的临时 SQLite/VecLite 路径执行 smoke，避免污染真实库。
- [x] 2.7 写入 `repeat-summary.json`，记录 package basename、sha、release label、source commit、source status、命令、状态、耗时、summary path、`skip_install`、`skip_e2e` 和安全边界，并脱敏绝对路径。
- [x] 2.8 确认脚本不调用真实 public providers、LLM providers、迁移、恢复、修复、升级、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、登录源、付费源、授权源、Level2 或高频源。

## 3. 发布材料与治理文档

- [x] 3.1 新增 `docs/release/acceptance/2026-06-18-p65-cross-machine-repeat.md`，记录隔离复验结果、命令矩阵、包 SHA、release label、source commit、source status、summary path、`skip_e2e=false`、caveats 和 Not Claimed。
- [x] 3.2 更新 `docs/release/README.md`，指向 P65 repeat acceptance。
- [x] 3.3 更新 `docs/release/release-handoff-2026-06-18.md`，把 P65 作为当前 package handoff 的复验入口。
- [x] 3.4 更新 `docs/development-plan.md`、`docs/README.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 3.5 确认发布材料不承诺物理第二机器已执行、远程发布、Git tag、自动升级、自动迁移、自动恢复、自动修复、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、登录源、付费源、授权源、Level2、高频源、未来 provider 可用性或收益。

## 4. 测试与验证

- [x] 4.1 运行 `bash scripts/local-release-package.sh --release-label p65-rc --output-dir tmp/p65-release`，生成由 P64 package workflow 创建的 P65 candidate archive。
- [x] 4.2 运行 package verify against generated P65 archive。
- [x] 4.3 运行 `bash scripts/local-release-repeat-acceptance.sh --archive <generated-archive> --output-dir tmp/p65-repeat`，主验收不得使用 `--skip-e2e`。
- [x] 4.4 扫描 repeat summary、manifest 和 archive listing，确认无完整 key、私有路径、SQLite DB、trace、log、raw prompt、raw vendor payload 或 local config。
- [x] 4.5 运行 `go test ./...`。
- [x] 4.6 运行 `npm --prefix web test`。
- [x] 4.7 运行 `npm --prefix web run build`。
- [x] 4.8 运行 `bash scripts/e2e-smoke.sh`。
- [x] 4.9 运行 `openspec validate p65-cross-machine-release-repeat-acceptance --strict`、`openspec validate --all --strict`、`git diff --check`。

## 5. 复审、归档与提交

- [x] 5.1 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 5.2 执行 OpenSpec archive。
- [x] 5.3 archive 后确认无活跃 change，并规划 P65 后下一步。
- [x] 5.4 提交前子 agent 复审无 Critical / Important。
- [ ] 5.5 提交 P65。
