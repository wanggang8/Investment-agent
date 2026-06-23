-- P1 默认规则种子数据。
-- 说明：只初始化本地可用的 active 规则版本和信源等级，不代表后续规则不可演进。

-- 默认 active 规则版本，对应现有需求文档中的根本规则快照。
INSERT OR IGNORE INTO rule_versions (
  rule_version,
  status,
  rules_json,
  effective_at,
  created_from_proposal_id,
  created_at
) VALUES (
  'v3.0',
  'active',
  '{"version":"v3.0","source":"docs/requirements.md","auto_trade":false}',
  CURRENT_TIMESTAMP,
  NULL,
  CURRENT_TIMESTAMP
);

-- 信源等级配置：S/A/B 可作为正式证据，C 级只能作为背景材料。
INSERT OR IGNORE INTO source_level_configs (source_level, description, formal_allowed, created_at) VALUES
  ('S', 'Official filings, exchanges, regulators, and primary company disclosures', 1, CURRENT_TIMESTAMP),
  ('A', 'Established financial media and reputable data providers', 1, CURRENT_TIMESTAMP),
  ('B', 'Secondary but attributable sources', 1, CURRENT_TIMESTAMP),
  ('C', 'Background-only low-confidence sources', 0, CURRENT_TIMESTAMP);
