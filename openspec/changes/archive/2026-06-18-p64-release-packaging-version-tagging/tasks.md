# Tasks: P64 Release Packaging And Version Tagging

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 P63 release candidate / handoff、P49 release-upgrade check、P54 repeatability docs 和现有 `scripts/local-*.sh`。
- [x] 1.3 创建 `p64-release-packaging-version-tagging` OpenSpec change。
- [x] 1.4 写明 P64 是本地发布工程层，不新增运行时业务能力、API、SQLite schema、Eino workflow、LLM 能力或前端产品能力。
- [x] 1.5 更新当前进度文档，标记 P64 active。
- [x] 1.6 运行 `openspec validate p64-release-packaging-version-tagging --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 1.7 子 agent 方案复审无 Critical / Important 后执行。

## 2. 打包脚本与 manifest

- [x] 2.1 新增 `scripts/local-release-package.sh`，支持 `--release-label VALUE`、`--output-dir PATH`、`--verify ARCHIVE`、`--skip-build`。
- [x] 2.2 默认将产物写入 `tmp/local-release-package/<timestamp>/`，不得写入真实用户数据库或私有配置目录。
- [x] 2.3 生成 staging 目录，只复制 release-safe 文件。
- [x] 2.4 排除 `.git/`、`.cursor/`、`tmp/`、`cmd/agent/tmp/`、`configs/config.local.yaml`、`docs/release/ui-audit-assets/`、`web/node_modules/`、`web/dist/`、日志、trace、SQLite、VecLite 本地索引和 raw payload。
- [x] 2.5 写入 `release-manifest.json`，包含 release label、commit、archive、sha256、included roots、excluded patterns、verification commands、acceptance references、known degradations、Not Claimed 和 safety note。
- [x] 2.6 创建 `.tar.gz` archive，并记录 checksum。
- [x] 2.7 verify 模式检查 manifest、archive checksum、必要入口和 forbidden patterns。

## 3. 测试与验证

- [x] 3.1 为 package manifest 或 verify helper 增加可维护测试；脚本逻辑保持 shell/python inline，由 package build/verify smoke 覆盖。
- [x] 3.2 运行 `bash scripts/local-release-package.sh --release-label p64-rc --output-dir tmp/p64-release`。
- [x] 3.3 运行 package verify against generated archive。
- [x] 3.4 扫描 manifest 和 archive listing，确认无完整 key、私有路径、SQLite DB、trace、log、raw prompt、raw vendor payload。
- [x] 3.5 运行 `go test ./...`。
- [x] 3.6 运行 `npm --prefix web test`。
- [x] 3.7 运行 `npm --prefix web run build`。
- [x] 3.8 运行 `bash scripts/e2e-smoke.sh`。

## 4. 发布材料

- [x] 4.1 新增 `docs/release/release-packaging-2026-06-18.md`，记录 package manifest、archive、checksum、verify 命令、排除项和安全边界。
- [x] 4.2 更新 `docs/release/README.md`，指向 P63 release-ready 和 P64 package materials。
- [x] 4.3 更新 `docs/release/release-handoff-2026-06-18.md`，加入 P64 package handoff 和复验入口。
- [x] 4.4 更新 `docs/development-plan.md`、`docs/README.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 4.5 确认发布材料不承诺收益、未来外部源可用性、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、自动升级、自动迁移、自动恢复、自动覆盖真实库、登录源、付费源、授权源、Level2、高频源或收益承诺。

## 5. 复审、归档与提交

- [x] 5.1 运行 `openspec validate p64-release-packaging-version-tagging --strict` 与 `openspec validate --all --strict`。
- [x] 5.2 运行 `git diff --check`。
- [x] 5.3 运行最终必要验证：package build/verify、sensitive scan、`go test ./...`、`npm --prefix web test`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh`。
- [x] 5.4 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 5.5 执行 OpenSpec archive。
- [x] 5.6 archive 后确认无活跃 change，并规划 P64 后下一步。
- [x] 5.7 提交前子 agent 复审无 Critical / Important。
- [x] 5.8 提交 P64。
