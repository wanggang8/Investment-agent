-- P1 SQLite 事实基准表结构。
-- 说明：SQLite 是本地唯一事实来源，VecLite 和前端 DTO 都可由这些表重建或派生。

-- 账户快照：保存一次决策读取到账户总览，用于历史复现。
CREATE TABLE IF NOT EXISTS portfolio_snapshots (
  snapshot_id TEXT PRIMARY KEY,
  snapshot_time DATETIME NOT NULL,
  cash REAL NOT NULL,
  total_assets REAL NOT NULL,
  cash_ratio REAL NOT NULL,
  high_risk_ratio REAL NOT NULL,
  position_count INTEGER NOT NULL,
  source TEXT NOT NULL CHECK (source IN ('manual', 'system')),
  created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_portfolio_snapshots_time ON portfolio_snapshots(snapshot_time);

-- 当前持仓：保存聚合后的当前态，历史复现应读取 position_snapshots。
CREATE TABLE IF NOT EXISTS positions (
  position_id TEXT PRIMARY KEY,
  symbol TEXT NOT NULL,
  name TEXT NOT NULL,
  quantity REAL NOT NULL,
  cost_price REAL NOT NULL,
  current_price REAL NOT NULL,
  market_value REAL NOT NULL,
  unrealized_profit_ratio REAL NOT NULL,
  position_state TEXT NOT NULL CHECK (position_state IN ('normal', 'sell_only', 'frozen_watch')),
  buy_date DATE,
  buy_reason TEXT,
  asset_tag TEXT,
  updated_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_positions_symbol ON positions(symbol);
CREATE INDEX IF NOT EXISTS idx_positions_state ON positions(position_state);

-- 持仓时点快照：与 portfolio_snapshots 一起复现当时完整持仓集合。
CREATE TABLE IF NOT EXISTS position_snapshots (
  position_snapshot_id TEXT PRIMARY KEY,
  snapshot_id TEXT NOT NULL,
  position_id TEXT,
  symbol TEXT NOT NULL,
  name TEXT NOT NULL,
  quantity REAL NOT NULL,
  cost_price REAL NOT NULL,
  current_price REAL NOT NULL,
  market_value REAL NOT NULL,
  unrealized_profit_ratio REAL NOT NULL,
  position_state TEXT NOT NULL CHECK (position_state IN ('normal', 'sell_only', 'frozen_watch')),
  buy_date DATE,
  buy_reason TEXT,
  asset_tag TEXT,
  created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_position_snapshots_snapshot ON position_snapshots(snapshot_id);
CREATE INDEX IF NOT EXISTS idx_position_snapshots_symbol ON position_snapshots(symbol);

-- 操作确认：只记录用户线下处理结果，不触发自动交易。
CREATE TABLE IF NOT EXISTS operation_confirmations (
  confirmation_id TEXT PRIMARY KEY,
  decision_id TEXT NOT NULL,
  confirmation_type TEXT NOT NULL CHECK (confirmation_type IN ('planned', 'executed_manually', 'watch', 'marked_error')),
  operation_type TEXT CHECK (operation_type IS NULL OR operation_type IN ('buy', 'sell', 'reduce')),
  symbol TEXT,
  quantity REAL,
  price REAL,
  fees REAL,
  executed_at DATETIME,
  error_case_id TEXT,
  payload_json TEXT,
  note TEXT,
  created_at DATETIME NOT NULL
);

-- 持仓流水：记录 executed_manually 产生的事实变化，当前持仓只是聚合结果。
CREATE TABLE IF NOT EXISTS position_transactions (
  transaction_id TEXT PRIMARY KEY,
  confirmation_id TEXT NOT NULL,
  symbol TEXT NOT NULL,
  operation_type TEXT NOT NULL CHECK (operation_type IN ('buy', 'sell', 'reduce')),
  quantity REAL NOT NULL,
  price REAL NOT NULL,
  fees REAL,
  occurred_at DATETIME NOT NULL,
  before_position_json TEXT,
  after_position_json TEXT,
  created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_position_transactions_confirmation ON position_transactions(confirmation_id);
CREATE INDEX IF NOT EXISTS idx_position_transactions_symbol ON position_transactions(symbol);

-- 市场快照：保存行情、估值、流动性和情绪指标，支撑规则裁决。
CREATE TABLE IF NOT EXISTS market_snapshots (
  market_snapshot_id TEXT PRIMARY KEY,
  symbol TEXT NOT NULL,
  trade_date DATE NOT NULL,
  close_price REAL,
  price_change_ratio REAL,
  volume REAL,
  turnover REAL,
  turnover_rate REAL,
  volatility REAL,
  margin_balance REAL,
  margin_balance_change REAL,
  pe REAL,
  pb REAL,
  pe_percentile REAL,
  pb_percentile REAL,
  volume_percentile REAL,
  volatility_percentile REAL,
  liquidity_state TEXT CHECK (liquidity_state IS NULL OR liquidity_state IN ('normal', 'warning', 'danger')),
  sentiment_state TEXT CHECK (sentiment_state IS NULL OR sentiment_state IN ('cold', 'neutral', 'hot', 'extreme')),
  market_metrics_json TEXT,
  created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_market_snapshots_symbol_date ON market_snapshots(symbol, trade_date);

-- 规则版本：保存正式规则快照；同一时间只能有一个 active 版本。
CREATE TABLE IF NOT EXISTS rule_versions (
  rule_version TEXT PRIMARY KEY,
  status TEXT NOT NULL CHECK (status IN ('active', 'archived')),
  rules_json TEXT NOT NULL,
  effective_at DATETIME NOT NULL,
  created_from_proposal_id TEXT,
  created_at DATETIME NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_rule_versions_single_active ON rule_versions(status) WHERE status = 'active';

-- 能力圈配置：定义系统允许分析的资产、标的和策略范围。
CREATE TABLE IF NOT EXISTS capability_configs (
  capability_id TEXT PRIMARY KEY,
  asset_types_json TEXT,
  symbols_json TEXT,
  excluded_symbols_json TEXT,
  strategy_scope_json TEXT,
  updated_at DATETIME NOT NULL
);

-- 用户设置：保存非规则类偏好，以及规则提案生成时需要的配置快照。
CREATE TABLE IF NOT EXISTS user_settings (
  settings_id TEXT PRIMARY KEY,
  position_limits_json TEXT,
  cash_min_ratio REAL,
  notification_config_json TEXT,
  data_sources_json TEXT,
  updated_at DATETIME NOT NULL
);

-- 原始情报：保存外部信源采集后的原始元信息。
CREATE TABLE IF NOT EXISTS intelligence_items (
  intelligence_id TEXT PRIMARY KEY,
  source_name TEXT NOT NULL,
  source_level TEXT NOT NULL CHECK (source_level IN ('S', 'A', 'B', 'C')),
  original_url TEXT,
  published_at DATETIME,
  captured_at DATETIME NOT NULL,
  content_hash TEXT NOT NULL,
  raw_title TEXT,
  raw_text_ref TEXT,
  created_at DATETIME NOT NULL
);

-- 结构化情报摘要：进入证据链和 RAG 索引的清洗后摘要。
-- 约束：C 级信源只能作为 background，不能进入正式裁决链。
CREATE TABLE IF NOT EXISTS intelligence_summary (
  summary_id TEXT PRIMARY KEY,
  intelligence_id TEXT NOT NULL,
  symbol TEXT,
  entity TEXT,
  event_type TEXT,
  impact_direction TEXT CHECK (impact_direction IS NULL OR impact_direction IN ('positive', 'negative', 'neutral')),
  summary TEXT NOT NULL,
  source_level TEXT NOT NULL CHECK (source_level IN ('S', 'A', 'B', 'C')),
  evidence_role TEXT NOT NULL CHECK (evidence_role IN ('formal', 'background')),
  time_weight REAL,
  relevance_score REAL,
  verification_group_id TEXT,
  created_at DATETIME NOT NULL,
  CHECK (source_level != 'C' OR evidence_role = 'background')
);
CREATE INDEX IF NOT EXISTS idx_intelligence_summary_symbol ON intelligence_summary(symbol);
CREATE INDEX IF NOT EXISTS idx_intelligence_summary_group ON intelligence_summary(verification_group_id);

-- RAG 文本块：保存可重建 VecLite 的文本块和索引元数据。
CREATE TABLE IF NOT EXISTS rag_chunks (
  chunk_id TEXT PRIMARY KEY,
  summary_id TEXT NOT NULL,
  chunk_text TEXT NOT NULL,
  chunk_hash TEXT NOT NULL,
  vector_id TEXT,
  vector_collection TEXT,
  embedding_model TEXT,
  embedding_version TEXT,
  index_version TEXT,
  index_status TEXT NOT NULL CHECK (index_status IN ('pending', 'indexed', 'stale', 'failed')),
  indexed_at DATETIME,
  metadata_json TEXT,
  created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_rag_chunks_summary ON rag_chunks(summary_id);
CREATE INDEX IF NOT EXISTS idx_rag_chunks_vector ON rag_chunks(vector_id);
CREATE INDEX IF NOT EXISTS idx_rag_chunks_status ON rag_chunks(index_status);

-- 决策记录：保存一次工作流的核心输出和最终规则裁决。
-- expected_return_scenarios_json 只用于展示，不得覆盖 final_verdict_status。
CREATE TABLE IF NOT EXISTS decision_records (
  decision_id TEXT PRIMARY KEY,
  request_id TEXT NOT NULL,
  workflow_type TEXT NOT NULL CHECK (workflow_type IN ('daily_discipline', 'consultation', 'evidence_verification', 'evolution_proposal', 'gatekeeper_audit', 'market_refresh')),
  symbol TEXT,
  question TEXT,
  workflow_status TEXT NOT NULL CHECK (workflow_status IN ('completed', 'degraded', 'failed')),
  record_type TEXT NOT NULL CHECK (record_type IN ('formal_trade_advice', 'non_trade_record', 'rejection_record')),
  dashboard_state TEXT NOT NULL CHECK (dashboard_state IN ('first_use', 'normal', 'insufficient_data', 'frozen_watch', 'high_risk')),
  capability_status TEXT CHECK (capability_status IS NULL OR capability_status IN ('in_scope', 'out_of_scope', 'unknown')),
  capability_reason TEXT,
  source_verification_status TEXT CHECK (source_verification_status IS NULL OR source_verification_status IN ('satisfied', 'failed', 'background_only')),
  risk_reason_code TEXT,
  media_heat_summary_json TEXT,
  user_emotion_tags_json TEXT,
  triggered_rules_json TEXT,
  errors_json TEXT,
  final_verdict_status TEXT NOT NULL CHECK (final_verdict_status IN ('buy_allowed', 'hold', 'reduce', 'sell_only', 'frozen_watch', 'rejected', 'insufficient_data')),
  final_verdict_text TEXT NOT NULL,
  prohibited_actions_json TEXT,
  optional_actions_json TEXT,
  confirmation_status TEXT NOT NULL CHECK (confirmation_status IN ('not_required', 'pending', 'planned', 'executed_manually', 'watch', 'marked_error')),
  portfolio_snapshot_id TEXT,
  market_snapshot_id TEXT,
  rule_version TEXT NOT NULL,
  analyst_reports_json TEXT,
  expected_return_scenarios_json TEXT,
  arbitration_chain_json TEXT,
  context_snapshot_json TEXT,
  created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_decision_records_created_at ON decision_records(created_at);
CREATE INDEX IF NOT EXISTS idx_decision_records_symbol ON decision_records(symbol);
CREATE INDEX IF NOT EXISTS idx_decision_records_status ON decision_records(confirmation_status);
CREATE INDEX IF NOT EXISTS idx_decision_records_record_type ON decision_records(record_type);
CREATE INDEX IF NOT EXISTS idx_decision_records_dashboard_state ON decision_records(dashboard_state);

-- 证据引用：保存决策使用过的证据快照，避免历史页面随情报变更而变化。
CREATE TABLE IF NOT EXISTS evidence_refs (
  evidence_ref_id TEXT PRIMARY KEY,
  evidence_id TEXT NOT NULL,
  decision_id TEXT NOT NULL,
  summary_id TEXT NOT NULL,
  source_name TEXT NOT NULL,
  source_level TEXT NOT NULL CHECK (source_level IN ('S', 'A', 'B', 'C')),
  evidence_role TEXT NOT NULL CHECK (evidence_role IN ('formal', 'background')),
  published_at DATETIME,
  captured_at DATETIME,
  original_url TEXT,
  summary TEXT NOT NULL,
  content_hash TEXT,
  time_weight REAL,
  relevance_score REAL,
  independent_source_count INTEGER NOT NULL DEFAULT 0,
  high_grade_independent_source_count INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL,
  CHECK (source_level != 'C' OR evidence_role = 'background')
);
CREATE INDEX IF NOT EXISTS idx_evidence_refs_decision ON evidence_refs(decision_id);
CREATE INDEX IF NOT EXISTS idx_evidence_refs_evidence ON evidence_refs(evidence_id);

-- 多源验证：保存同一事件或标的的独立信源验证结果。
CREATE TABLE IF NOT EXISTS source_verifications (
  verification_id TEXT PRIMARY KEY,
  verification_group_id TEXT NOT NULL,
  event_id TEXT NOT NULL,
  symbol TEXT,
  event_type TEXT,
  evidence_role TEXT NOT NULL CHECK (evidence_role IN ('formal', 'background')),
  verification_status TEXT NOT NULL CHECK (verification_status IN ('satisfied', 'failed', 'background_only')),
  independent_source_count INTEGER NOT NULL,
  high_grade_independent_source_count INTEGER NOT NULL DEFAULT 0,
  highest_source_level TEXT CHECK (highest_source_level IS NULL OR highest_source_level IN ('S', 'A', 'B', 'C')),
  latest_published_at DATETIME,
  evidence_ids_json TEXT,
  created_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_source_verifications_group ON source_verifications(verification_group_id);
CREATE INDEX IF NOT EXISTS idx_source_verifications_symbol ON source_verifications(symbol);
CREATE INDEX IF NOT EXISTS idx_source_verifications_symbol_event ON source_verifications(symbol, event_id);

-- 审计事件：保存工作流节点、用户动作、规则提案和错误状态的可追溯记录。
CREATE TABLE IF NOT EXISTS audit_events (
  audit_event_id TEXT PRIMARY KEY,
  request_id TEXT,
  decision_id TEXT,
  workflow_type TEXT,
  node_name TEXT,
  actor TEXT NOT NULL CHECK (actor IN ('system', 'user', 'gatekeeper')),
  action TEXT NOT NULL CHECK (action IN ('generate_decision', 'confirm_operation', 'mark_error', 'create_proposal', 'audit_rule_change', 'update_rule', 'refresh_market_data', 'update_settings', 'update_capability', 'rebuild_index', 'run_local_task', 'risk_alert')),
  node_action TEXT,
  proposal_id TEXT,
  confirmation_id TEXT,
  error_case_id TEXT,
  status TEXT NOT NULL CHECK (status IN ('success', 'degraded', 'failed')),
  error_code TEXT,
  before_state TEXT,
  after_state TEXT,
  rule_version TEXT,
  snapshot_id TEXT,
  input_ref_type TEXT,
  input_ref TEXT,
  output_ref_type TEXT,
  output_ref TEXT,
  created_at DATETIME NOT NULL,
  CHECK (status != 'failed' OR error_code IS NOT NULL)
);
CREATE INDEX IF NOT EXISTS idx_audit_events_request ON audit_events(request_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_decision ON audit_events(decision_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_created ON audit_events(created_at);

-- 错误案例：用户标记错误后形成的复盘样本，供规则提案生成使用。
CREATE TABLE IF NOT EXISTS error_cases (
  error_case_id TEXT PRIMARY KEY,
  decision_id TEXT NOT NULL,
  confirmation_id TEXT NOT NULL,
  actual_outcome TEXT,
  root_cause_tag TEXT CHECK (root_cause_tag IS NULL OR root_cause_tag IN ('evidence_missed', 'rule_threshold_issue', 'analyst_error', 'user_context_missing', 'market_exception')),
  lesson_learned TEXT,
  created_at DATETIME NOT NULL
);

-- 规则提案：规则演进的草案状态，不直接修改正式规则版本。
CREATE TABLE IF NOT EXISTS rule_proposals (
  proposal_id TEXT PRIMARY KEY,
  proposal_type TEXT NOT NULL CHECK (proposal_type IN ('threshold', 'sop', 'risk_rule', 'capability')),
  status TEXT NOT NULL CHECK (status IN ('draft', 'pending_user_confirm', 'under_gatekeeper_audit', 'pending_final_confirm', 'rejected', 'applied')),
  source_error_case_id TEXT,
  title TEXT NOT NULL,
  proposal_version TEXT NOT NULL,
  before_rule_json TEXT,
  after_rule_json TEXT,
  reason TEXT,
  impact_scope_json TEXT,
  risk_notes_json TEXT,
  sample_count INTEGER NOT NULL DEFAULT 0,
  final_confirmed_at DATETIME,
  final_confirmed_note TEXT,
  applied_rule_version TEXT,
  related_error_cases_json TEXT,
  created_at DATETIME NOT NULL
);

-- 守门人审计：记录规则提案是否违反根本规则、是否存在冲突以及是否允许应用。
CREATE TABLE IF NOT EXISTS gatekeeper_audits (
  gatekeeper_audit_id TEXT PRIMARY KEY,
  proposal_id TEXT NOT NULL,
  audit_result TEXT NOT NULL CHECK (audit_result IN ('approved', 'rejected', 'needs_user_review')),
  audit_reason TEXT,
  required_changes TEXT,
  violates_fundamental_rule INTEGER NOT NULL CHECK (violates_fundamental_rule IN (0, 1)),
  has_rule_conflict INTEGER NOT NULL CHECK (has_rule_conflict IN (0, 1)),
  backtest_metrics_json TEXT,
  allow_apply INTEGER NOT NULL CHECK (allow_apply IN (0, 1)),
  audited_rule_version TEXT NOT NULL,
  created_at DATETIME NOT NULL
);

-- 信源等级配置：定义不同等级信源是否允许进入正式裁决链。
CREATE TABLE IF NOT EXISTS source_level_configs (
  source_level TEXT PRIMARY KEY CHECK (source_level IN ('S', 'A', 'B', 'C')),
  description TEXT NOT NULL,
  formal_allowed INTEGER NOT NULL CHECK (formal_allowed IN (0, 1)),
  created_at DATETIME NOT NULL
);
