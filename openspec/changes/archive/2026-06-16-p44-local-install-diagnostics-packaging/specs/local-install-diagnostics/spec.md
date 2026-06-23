## ADDED Requirements

### Requirement: Local install diagnostics and packaging page SHALL exist

系统 **SHALL** 新增本地运维上手入口 `/local-install`，用于汇总本地安装检查、配置草稿和诊断打包。

#### Scenario: User checks install guidance and commands

**WHEN** 用户打开 `/local-install`
**THEN** 页面展示本地安装步骤、关键命令模板和安全边界提示。
**AND** 页面不包含交易、自动确认、自动规则应用、外部推送、收益承诺等动作入口。

#### Scenario: User generates startup configuration draft

**WHEN** 用户在配置向导中填写 host、port、数据库/索引路径等字段
**THEN** 页面展示可粘贴的启动配置草稿。
**AND** 页面不直接写入本地文件。

#### Scenario: User imports diagnostics summary

**WHEN** 用户上传 `local-install-diagnostics.sh` 产出的诊断摘要 JSON 文件
**THEN** 页面仅只读展示步骤状态、失败项、摘要路径与生成时间。
**AND** 页面不显示数据库路径、完整密钥、SQL 文本或 raw HTTP 响应。

## REMOVED Requirements

本阶段不移除既有契约条款。
