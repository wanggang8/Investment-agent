package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/pkg/apperr"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := config.Load("")
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}
	store, err := appsqlite.Open(cfg.SQLite.Path)
	if err != nil {
		return err
	}
	defer store.Close()
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		return err
	}
	ctx := context.Background()
	now := time.Now().UTC().Format(time.RFC3339)
	decisionRepo := appsqlite.NewDecisionRepository(store.DB)
	intelligenceRepo := appsqlite.NewIntelligenceRepository(store.DB)
	auditRepo := appsqlite.NewAuditRepository(store.DB)
	dailyAutoRunRepo := appsqlite.NewDailyAutoRunRepository(store.DB)
	dailyReportRepo := appsqlite.NewDailyDisciplineReportRepository(store.DB)
	marketRepo := appsqlite.NewMarketRepository(store.DB)
	riskAlertRepo := appsqlite.NewRiskAlertRepository(store.DB)
	ruleRepo := appsqlite.NewRuleRepository(store.DB)
	ruleEffectRepo := appsqlite.NewRuleEffectRepository(store.DB)
	settingsRepo := appsqlite.NewSettingsRepository(store.DB)
	decisionExists := false
	if _, _, err := decisionRepo.GetDecisionRecord(ctx, "decision_smoke_p30"); err == nil {
		decisionExists = true
	}
	if !decisionExists {
		if err := decisionRepo.SaveDecisionRecord(ctx, repository.DecisionRecord{
			DecisionID:                  "decision_smoke_p30",
			RequestID:                   "req_smoke_p30",
			WorkflowType:                "consultation",
			Symbol:                      "510300",
			Question:                    "P30 本地 E2E smoke 决策",
			WorkflowStatus:              "completed",
			RecordType:                  "formal_trade_advice",
			DashboardState:              "normal",
			CapabilityStatus:            "in_scope",
			CapabilityReason:            "P30 smoke fixture",
			SourceVerificationStatus:    "satisfied",
			TriggeredRulesJSON:          `[{"rule_id":"rule_smoke_p30","rule_name":"P30 smoke 规则","severity":"info","description":"只读 smoke fixture"}]`,
			ErrorsJSON:                  `[]`,
			FinalVerdictStatus:          "hold",
			FinalVerdictText:            "继续持有",
			ProhibitedActionsJSON:       `[]`,
			OptionalActionsJSON:         `["人工复核"]`,
			ConfirmationStatus:          "pending",
			RuleVersion:                 "rule_version_smoke_p30",
			AnalystReportsJSON:          `[{"agent_name":"P30SmokeAnalyst","conclusion":"本地 smoke fixture 可渲染。","key_reasons":[],"risk_warnings":[],"confidence":"qualitative","evidence_ids":[]}]`,
			ExpectedReturnScenariosJSON: `{"precision_status":"insufficient","reason":"P30 smoke 样本不足，仅展示定性情景。","sample_count":8,"sample_window":"2026-01-01/2026-06-01","screening_condition":"local smoke fixture","scenarios":[{"scenario":"base","return_range":"0%~3%","probability":null,"trigger":"人工复核"}],"sell_evaluation":{"status":"review_needed","triggers":["base_midpoint_downshift"],"prompts":["人工复核预期收益边界"],"actions":["记录人工计划"],"non_trading_disclaimer":"卖出评估仅用于人工复核，不会自动交易。"},"reassessment_trigger":{"reason":"base_midpoint_downshift","boundary":"base midpoint","current_value":0.012}}`,
			ArbitrationChainJSON:        `[{"priority":1,"rule_id":"rule_smoke_p30","result":"P30 smoke 只读裁决链"}]`,
			ContextSnapshotJSON:         `{"retrieval_quality_summary":{"query_summary":"510300","top_k":1,"status":"degraded","index_health":"missing","index_freshness":"unknown","fallback_source":"sqlite_summary","source_consistency_status":"checked","degraded_reason":"veclite index missing in local fixture"}}`,
			CreatedAt:                   now,
		}, []repository.EvidenceRef{{
			EvidenceRefID:                   "eref_smoke_p30",
			EvidenceID:                      "summary_smoke_p30",
			DecisionID:                      "decision_smoke_p30",
			SummaryID:                       "summary_smoke_p30",
			SourceName:                      "P30SmokeSource",
			SourceLevel:                     "A",
			EvidenceRole:                    "formal",
			PublishedAt:                     now,
			CapturedAt:                      now,
			OriginalURL:                     "https://example.invalid/p30-smoke",
			Summary:                         "P30 smoke 证据摘要",
			ContentHash:                     "hash_smoke_p30",
			TimeWeight:                      1,
			RelevanceScore:                  0.9,
			IndependentSourceCount:          1,
			HighGradeIndependentSourceCount: 1,
			CreatedAt:                       now,
		}}); err != nil {
			return err
		}
		if err := intelligenceRepo.SaveIntelligenceItem(ctx, repository.IntelligenceItem{
			IntelligenceID: "intel_smoke_p30",
			SourceName:     "P30SmokeSource",
			SourceLevel:    "A",
			OriginalURL:    "https://example.invalid/p30-smoke",
			PublishedAt:    now,
			CapturedAt:     now,
			ContentHash:    "hash_smoke_p30",
			RawTitle:       "P30 smoke 证据摘要",
			RawTextRef:     "local-smoke",
			CreatedAt:      now,
		}); err != nil {
			return err
		}
		if err := intelligenceRepo.SaveIntelligenceSummary(ctx, repository.IntelligenceSummary{
			SummaryID:                       "summary_smoke_p30",
			IntelligenceID:                  "intel_smoke_p30",
			Symbol:                          "510300",
			Entity:                          "510300",
			EventType:                       "smoke",
			ImpactDirection:                 "neutral",
			Summary:                         "P30 smoke 证据摘要",
			SourceLevel:                     "A",
			EvidenceRole:                    "formal",
			TimeWeight:                      1,
			RelevanceScore:                  0.9,
			VerificationGroupID:             "verification_group_smoke_p30",
			VerificationStatus:              "satisfied",
			IndependentSourceCount:          1,
			HighGradeIndependentSourceCount: 1,
			SourceName:                      "P30SmokeSource",
			OriginalURL:                     "https://example.invalid/p30-smoke",
			PublishedAt:                     now,
			CapturedAt:                      now,
			ContentHash:                     "hash_smoke_p30",
			CreatedAt:                       now,
		}, nil); err != nil {
			return err
		}
		if err := intelligenceRepo.SaveSourceVerification(ctx, repository.SourceVerification{
			VerificationID:                  "verification_smoke_p30",
			VerificationGroupID:             "verification_group_smoke_p30",
			EventID:                         "event_smoke_p30",
			Symbol:                          "510300",
			EventType:                       "smoke",
			EvidenceRole:                    "formal",
			VerificationStatus:              "satisfied",
			IndependentSourceCount:          1,
			HighGradeIndependentSourceCount: 1,
			HighestSourceLevel:              "A",
			LatestPublishedAt:               now,
			EvidenceIDsJSON:                 `["summary_smoke_p30"]`,
			CreatedAt:                       now,
		}); err != nil {
			return err
		}
		if err := appendAuditEventIfMissing(ctx, auditRepo, repository.AuditEvent{
			AuditEventID:  "audit_smoke_p30",
			RequestID:     "req_smoke_p30",
			DecisionID:    "decision_smoke_p30",
			WorkflowType:  "local_e2e_smoke",
			NodeName:      "P30SmokeSeed",
			Actor:         "system",
			Action:        "run_local_task",
			NodeAction:    "seed_local_e2e_smoke",
			Status:        "success",
			InputRefType:  "openspec_change",
			InputRef:      "p30-real-e2e-smoke",
			OutputRefType: "decision_id",
			OutputRef:     "decision_smoke_p30",
			CreatedAt:     now,
		}); err != nil {
			return err
		}
	}
	localDate, err := configuredLocalDate(cfg.DailyAutoRun.Timezone)
	if err != nil {
		return err
	}
	scope := strings.TrimSpace(cfg.DailyAutoRun.Scope)
	if scope == "" {
		scope = "holdings"
	}
	symbolSetHash := "e3b0c44298fc1c14"
	autoRunKey := fmt.Sprintf("%s:%s:%s:v1", localDate, scope, symbolSetHash)
	if err := dailyAutoRunRepo.UpsertDailyAutoRunState(ctx, repository.DailyAutoRunState{RunID: "auto_run_smoke_p31", IdempotencyKey: autoRunKey, LocalDate: localDate, Scope: scope, SymbolSetHash: symbolSetHash, Status: "failed", LastRunAt: now, NextRunAt: now, FailureCode: "missing_prerequisites", FailureReason: "缺少本地持仓", CreatedAt: now, UpdatedAt: now}); err != nil {
		return err
	}
	reportSymbolSetHash := "p32smokereport"
	reportSourceID := "auto_run_smoke_p32_success"
	if err := dailyReportRepo.UpsertDailyDisciplineReport(ctx, repository.DailyDisciplineReport{ReportID: "daily_report_smoke_p32", LocalDate: localDate, Scope: scope, SymbolSetHash: reportSymbolSetHash, SourceType: "auto_run", SourceID: reportSourceID, DecisionID: "decision_smoke_p30", Status: "success", Summary: "P32 smoke 今日纪律报告已生成", CreatedAt: now, UpdatedAt: now}); err != nil {
		return err
	}
	if err := appsqlite.NewNotificationRepository(store.DB).SaveNotification(ctx, repository.Notification{NotificationID: "notif_smoke_p31_auto_run", Type: "daily_auto_run_failed", Severity: "warning", Title: "每日自动运行未完成", Message: "缺少本地持仓", SourceType: "daily_auto_run", SourceID: autoRunKey, CreatedAt: now}); err != nil {
		return err
	}
	if err := appendAuditEventIfMissing(ctx, auditRepo, repository.AuditEvent{AuditEventID: "audit_smoke_p31_auto_run", RequestID: "auto_run_smoke_p31", WorkflowType: "daily_auto_run", NodeName: "DailyAutoRunner", Actor: "system", Action: "run_local_task", NodeAction: "daily_auto_run", Status: "failed", ErrorCode: "missing_prerequisites", InputRefType: "idempotency_key", InputRef: autoRunKey, OutputRefType: "diagnostic", OutputRef: "status=failed;step=prerequisites;safety=no_auto_trading;code=missing_prerequisites;reason=缺少本地持仓", CreatedAt: now}); err != nil {
		return err
	}
	if err := seedP39DecisionFixtures(ctx, decisionRepo, now); err != nil {
		return err
	}
	if err := seedP73EffectivenessFixtures(ctx, decisionRepo, intelligenceRepo, auditRepo, now); err != nil {
		return err
	}
	if err := seedP39LocalSettings(ctx, settingsRepo, now); err != nil {
		return err
	}
	return seedP39JourneyFacts(ctx, marketRepo, riskAlertRepo, ruleRepo, ruleEffectRepo, auditRepo, now)
}

func seedP39DecisionFixtures(ctx context.Context, decisionRepo *appsqlite.DecisionRepository, now string) error {
	if _, _, err := decisionRepo.GetDecisionRecord(ctx, "decision_smoke_p39_out_of_scope"); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := decisionRepo.SaveDecisionRecord(ctx, repository.DecisionRecord{
			DecisionID:                  "decision_smoke_p39_out_of_scope",
			RequestID:                   "req_smoke_p39_out_of_scope",
			WorkflowType:                "consultation",
			Symbol:                      "159915",
			Question:                    "P39 能力圈外降级 fixture",
			WorkflowStatus:              "completed",
			RecordType:                  "rejection_record",
			DashboardState:              "frozen_watch",
			CapabilityStatus:            "out_of_scope",
			CapabilityReason:            "P39 fixture excluded symbol",
			SourceVerificationStatus:    "background_only",
			TriggeredRulesJSON:          `[{"rule_id":"capability_scope","rule_name":"能力圈边界","severity":"warning","description":"能力圈外只能拒绝交易类建议"}]`,
			ErrorsJSON:                  `[]`,
			FinalVerdictStatus:          "rejected",
			FinalVerdictText:            "能力圈外，拒绝交易类建议",
			ProhibitedActionsJSON:       `["自动交易","收益承诺","自动规则应用"]`,
			OptionalActionsJSON:         `["人工复核能力圈配置"]`,
			ConfirmationStatus:          "not_required",
			RuleVersion:                 "rule_version_smoke_p30",
			AnalystReportsJSON:          `[]`,
			ExpectedReturnScenariosJSON: `{"precision_status":"insufficient","reason":"能力圈外，不生成收益判断。","sample_count":0,"missing_categories":["capability_scope"],"scenarios":[]}`,
			ArbitrationChainJSON:        `[{"priority":1,"rule_id":"capability_scope","result":"out_of_scope rejected"}]`,
			ContextSnapshotJSON:         `{"retrieval_quality_summary":{"query_summary":"159915","top_k":0,"status":"empty","index_health":"missing","index_freshness":"unknown","fallback_source":"none","source_consistency_status":"unknown","degraded_reason":"out_of_scope"}}`,
			CreatedAt:                   now,
		}, nil); err != nil {
			return err
		}
	}
	if _, _, err := decisionRepo.GetDecisionRecord(ctx, "decision_smoke_p39_llm_degraded"); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := decisionRepo.SaveDecisionRecord(ctx, repository.DecisionRecord{
			DecisionID:                  "decision_smoke_p39_llm_degraded",
			RequestID:                   "req_smoke_p39_llm_degraded",
			WorkflowType:                "consultation",
			Symbol:                      "510300",
			Question:                    "P39 LLM 降级 fixture",
			WorkflowStatus:              "degraded",
			RecordType:                  "non_trade_record",
			DashboardState:              "insufficient_data",
			CapabilityStatus:            "in_scope",
			CapabilityReason:            "P39 fixture",
			SourceVerificationStatus:    "failed",
			TriggeredRulesJSON:          `[{"rule_id":"analyst_degraded","rule_name":"分析材料降级","severity":"warning","description":"LLM 不可用时只展示规则与已有事实"}]`,
			ErrorsJSON:                  `["ANALYST_UNAVAILABLE"]`,
			FinalVerdictStatus:          "insufficient_data",
			FinalVerdictText:            "LLM 降级，暂停交易类建议",
			ProhibitedActionsJSON:       `["自动交易","收益承诺","自动规则应用"]`,
			OptionalActionsJSON:         `["人工复核已有事实"]`,
			ConfirmationStatus:          "not_required",
			RuleVersion:                 "rule_version_smoke_p30",
			AnalystReportsJSON:          `[{"agent_name":"P39LLMDegraded","conclusion":"分析服务暂不可用，页面仅展示规则与已有数据。","key_reasons":["deepseek api key missing in local fixture"],"risk_warnings":["不得用降级材料生成交易执行语义"],"confidence":"low","evidence_ids":[]}]`,
			ExpectedReturnScenariosJSON: `{"precision_status":"insufficient","reason":"LLM 降级时不展示精确收益概率。","sample_count":0,"missing_categories":["llm_analysis"],"scenarios":[]}`,
			ArbitrationChainJSON:        `[{"priority":1,"rule_id":"analyst_degraded","result":"insufficient_data"}]`,
			ContextSnapshotJSON:         `{"retrieval_quality_summary":{"query_summary":"510300","top_k":0,"status":"degraded","index_health":"missing","index_freshness":"unknown","fallback_source":"sqlite_summary","source_consistency_status":"unknown","degraded_reason":"llm_unavailable_fixture"}}`,
			CreatedAt:                   now,
		}, nil); err != nil {
			return err
		}
	}
	return nil
}

func seedP73EffectivenessFixtures(ctx context.Context, decisionRepo *appsqlite.DecisionRepository, intelligenceRepo *appsqlite.IntelligenceRepository, auditRepo *appsqlite.AuditRepository, now string) error {
	if _, err := intelligenceRepo.GetIntelligenceItem(ctx, "intel_smoke_p73_background"); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := intelligenceRepo.SaveIntelligenceItem(ctx, repository.IntelligenceItem{
			IntelligenceID: "intel_smoke_p73_background",
			SourceName:     "P73BackgroundSource",
			SourceLevel:    "C",
			OriginalURL:    "https://example.invalid/p73-background",
			PublishedAt:    now,
			CapturedAt:     now,
			ContentHash:    "hash_smoke_p73_background",
			RawTitle:       "P73 C 级背景材料",
			RawTextRef:     "local-smoke-p73",
			CreatedAt:      now,
		}); err != nil {
			return err
		}
		if err := intelligenceRepo.SaveIntelligenceSummary(ctx, repository.IntelligenceSummary{
			SummaryID:                       "summary_smoke_p73_background",
			IntelligenceID:                  "intel_smoke_p73_background",
			Symbol:                          "510300",
			Entity:                          "510300",
			EventType:                       "background_note",
			ImpactDirection:                 "neutral",
			Summary:                         "P73 背景材料：用户从社群听到短期观点，仅能作为背景，不能作为正式裁决依据。",
			SourceLevel:                     "C",
			EvidenceRole:                    "background",
			TimeWeight:                      0.2,
			RelevanceScore:                  0.4,
			VerificationGroupID:             "verification_group_smoke_p73_background",
			VerificationStatus:              "background_only",
			IndependentSourceCount:          0,
			HighGradeIndependentSourceCount: 0,
			SourceName:                      "P73BackgroundSource",
			OriginalURL:                     "https://example.invalid/p73-background",
			PublishedAt:                     now,
			CapturedAt:                      now,
			ContentHash:                     "hash_smoke_p73_background",
			CreatedAt:                       now,
		}, []repository.RAGChunk{{
			ChunkID:     "chunk_smoke_p73_background",
			SummaryID:   "summary_smoke_p73_background",
			Symbol:      "510300",
			ChunkText:   "P73 背景材料：C 级社群观点仅能作为背景，不能作为正式裁决依据。",
			ChunkHash:   "chunk_hash_smoke_p73_background",
			IndexStatus: "pending",
			CreatedAt:   now,
		}}); err != nil {
			return err
		}
		if err := intelligenceRepo.SaveSourceVerification(ctx, repository.SourceVerification{
			VerificationID:                  "verification_smoke_p73_background",
			VerificationGroupID:             "verification_group_smoke_p73_background",
			EventID:                         "event_smoke_p73_background",
			Symbol:                          "510300",
			EventType:                       "background_note",
			EvidenceRole:                    "background",
			VerificationStatus:              "background_only",
			IndependentSourceCount:          0,
			HighGradeIndependentSourceCount: 0,
			HighestSourceLevel:              "C",
			LatestPublishedAt:               now,
			EvidenceIDsJSON:                 `["summary_smoke_p73_background"]`,
			CreatedAt:                       now,
		}); err != nil {
			return err
		}
	}
	if _, _, err := decisionRepo.GetDecisionRecord(ctx, "decision_smoke_p73_background_only"); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := decisionRepo.SaveDecisionRecord(ctx, repository.DecisionRecord{
			DecisionID:                  "decision_smoke_p73_background_only",
			RequestID:                   "req_smoke_p73_background_only",
			WorkflowType:                "consultation",
			Symbol:                      "510300",
			Question:                    "P73 C 级背景材料是否足够支持交易类建议？",
			WorkflowStatus:              "completed",
			RecordType:                  "non_trade_record",
			DashboardState:              "insufficient_data",
			CapabilityStatus:            "in_scope",
			CapabilityReason:            "P73 effectiveness fixture",
			SourceVerificationStatus:    "background_only",
			TriggeredRulesJSON:          `[{"rule_id":"formal_evidence_required","rule_name":"正式证据门禁","severity":"warning","description":"C 级背景材料不得作为正式裁决依据"}]`,
			ErrorsJSON:                  `[]`,
			FinalVerdictStatus:          "insufficient_data",
			FinalVerdictText:            "仅有背景材料，不能生成交易类建议",
			ProhibitedActionsJSON:       `["自动交易","收益承诺","把 C 级材料作为正式证据"]`,
			OptionalActionsJSON:         `["补充 A/B/S 级正式证据","人工复核"]`,
			ConfirmationStatus:          "not_required",
			RuleVersion:                 "rule_version_smoke_p30",
			AnalystReportsJSON:          `[]`,
			ExpectedReturnScenariosJSON: `{"precision_status":"unavailable","reason":"仅有 C 级背景材料，不生成收益区间。","sample_count":0,"missing_categories":["formal_evidence"],"scenarios":[]}`,
			ArbitrationChainJSON:        `[{"priority":1,"rule_id":"formal_evidence_required","result":"background_only_insufficient_data"}]`,
			ContextSnapshotJSON:         `{"retrieval_quality_summary":{"query_summary":"510300 background","top_k":1,"status":"hit","index_health":"missing","index_freshness":"unknown","fallback_source":"sqlite_summary","source_consistency_status":"checked","degraded_reason":"formal evidence unavailable"}}`,
			CreatedAt:                   now,
		}, nil); err != nil {
			return err
		}
	}
	return appendAuditEventIfMissing(ctx, auditRepo, repository.AuditEvent{AuditEventID: "audit_smoke_p73_effectiveness", RequestID: "req_smoke_p73_background_only", DecisionID: "decision_smoke_p73_background_only", WorkflowType: "effectiveness_replay", NodeName: "P73EffectivenessSeed", Actor: "system", Action: "run_local_task", NodeAction: "seed_effectiveness_fixture", Status: "success", InputRefType: "openspec_change", InputRef: "p73-product-effectiveness-ux-validation", OutputRefType: "fixture", OutputRef: "background_only=blocked;no_auto_trading;no_return_promise", CreatedAt: now})
}

func seedP39LocalSettings(ctx context.Context, settingsRepo *appsqlite.SettingsRepository, now string) error {
	if _, err := settingsRepo.GetLatestCapabilityConfig(ctx); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := settingsRepo.SaveCapabilityConfig(ctx, repository.CapabilityConfig{
			CapabilityID:        "cap_smoke_p39",
			AssetTypesJSON:      `["ETF"]`,
			SymbolsJSON:         `["510300","NO_MARKET"]`,
			ExcludedSymbolsJSON: `["159915"]`,
			StrategyScopeJSON:   `["consultation","discipline_review"]`,
			UpdatedAt:           now,
		}); err != nil {
			return err
		}
	}
	if _, err := settingsRepo.GetLatestSystemSettings(ctx); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		return settingsRepo.SaveSystemSettings(ctx, repository.SystemSettings{
			SettingsID:             "settings_smoke_p39",
			NotificationConfigJSON: `{"enabled":false,"page_preference":"local_e2e"}`,
			DataSourcesJSON:        `["stub","p39_fixture"]`,
			UpdatedAt:              now,
		})
	}
	return nil
}

func seedP39JourneyFacts(ctx context.Context, marketRepo *appsqlite.MarketRepository, riskAlertRepo *appsqlite.RiskAlertRepository, ruleRepo *appsqlite.RuleRepository, ruleEffectRepo *appsqlite.RuleEffectRepository, auditRepo *appsqlite.AuditRepository, now string) error {
	if exists, err := marketRepo.MarketSnapshotExists(ctx, "market_smoke_p39"); err != nil {
		return err
	} else if !exists {
		if err := marketRepo.SaveMarketSnapshot(ctx, model.MarketSnapshot{
			MarketSnapshotID:  "market_smoke_p39",
			Symbol:            "510300",
			TradeDate:         now[:10],
			ClosePrice:        3.88,
			TurnoverRate:      0.42,
			PEPercentile:      0.82,
			PBPercentile:      0.76,
			LiquidityState:    model.LiquidityWarning,
			SentimentState:    model.SentimentHot,
			MarketMetricsJSON: `{"source_name":"csindex_extended","source_level":"A","source_type":"public_market","metadata":{"p34_data_categories":["index_valuation_files"],"p34_source_health":{"index_valuation_files":{"source_name":"csindex_extended","source_level":"A","source_type":"public_market","freshness":"stale","data_date":"2026-06-05","last_success_at":"2026-06-05T15:00:00Z","last_failure_at":"2026-06-16T08:00:00Z","failure_category":"stale","affected_symbols":["510300"]}}}}`,
		}, now); err != nil {
			return err
		}
	}

	if _, err := ruleRepo.GetRuleProposal(ctx, "prop_smoke_p39"); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := ruleRepo.SaveRuleProposal(ctx, repository.RuleProposal{
			ProposalID:            "prop_smoke_p39",
			ProposalType:          "threshold",
			Status:                "pending_final_confirm",
			SourceErrorCaseID:     "err_smoke_p39",
			Title:                 "P39 E2E 规则提案",
			ProposalVersion:       "v_p39_draft",
			BeforeRuleJSON:        `{"content":"旧阈值：数据源过期时继续普通观察"}`,
			AfterRuleJSON:         `{"content":"新阈值：数据源过期时进入人工复核队列"}`,
			Reason:                "季度阈值复盘发现 source health stale 需要人工复核。",
			ImpactScopeJSON:       `{"scope":"review_and_risk_alert","auto_apply":false}`,
			RiskNotesJSON:         `{"safety":"仍需守门人审计和最终确认，不自动应用规则"}`,
			SampleCount:           6,
			RelatedErrorCasesJSON: `["err_smoke_p39"]`,
			CreatedAt:             now,
		}); err != nil {
			return err
		}
	}
	if _, err := ruleRepo.GetGatekeeperAudit(ctx, "gk_smoke_p39"); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := ruleRepo.SaveGatekeeperAudit(ctx, repository.GatekeeperAudit{
			GatekeeperAuditID:   "gk_smoke_p39",
			ProposalID:          "prop_smoke_p39",
			AuditResult:         "approved",
			AuditReason:         "P39 fixture：守门人通过，但仍需用户最终确认。",
			BacktestMetricsJSON: `{"sample_count":6,"missed":0}`,
			AllowApply:          true,
			AuditedRuleVersion:  "v_p39_draft",
			CreatedAt:           now,
		}); err != nil {
			return err
		}
	}
	if _, err := ruleEffectRepo.GetRuleEffectValidation(ctx, "val_smoke_p39"); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := ruleEffectRepo.SaveRuleEffectValidation(ctx, repository.RuleEffectValidation{
			ValidationID:             "val_smoke_p39",
			ProposalID:               "prop_smoke_p39",
			CandidateRuleVersion:     "v_p39_draft",
			ValidationStatus:         model.RuleEffectValidationPassed,
			SampleCount:              6,
			SampleWindow:             "2026-Q2",
			RepresentativenessStatus: model.RuleEffectValidationPassed,
			OverfitRisk:              model.RuleEffectOverfitLow,
			ReplayResult:             model.RuleEffectReplayPassed,
			GuardrailDecision:        model.RuleEffectGuardrailPassed,
			SourceExplanationJSON:    `{"source_case_count":3,"related_error_case_ids":["err_smoke_p39"],"related_decision_ids":["decision_smoke_p30"],"related_risk_alert_ids":["risk_smoke_p39"]}`,
			MetricsJSON:              `{"hit_count":6,"misjudgment_count":0,"missing_evidence_count":1,"degraded_count":1,"risk_alert_count":1}`,
			RiskNotesJSON:            `["E2E fixture 只用于浏览器验收，不承诺收益"]`,
			RelatedErrorCasesJSON:    `["err_smoke_p39"]`,
			RelatedDecisionIDsJSON:   `["decision_smoke_p30"]`,
			RelatedRiskAlertIDsJSON:  `["risk_smoke_p39"]`,
			RelatedAuditEventIDsJSON: `["audit_smoke_p39"]`,
			SafetyNote:               "规则效果验证只读展示，不自动应用规则。",
			CreatedAt:                now,
			UpdatedAt:                now,
		}); err != nil {
			return err
		}
	}
	if _, err := ruleEffectRepo.GetRuleEffectTracking(ctx, "track_smoke_p39"); err != nil {
		if !apperr.IsCode(err, apperr.CodeNotFound) {
			return err
		}
		if err := ruleEffectRepo.SaveRuleEffectTracking(ctx, repository.RuleEffectTracking{
			TrackingID:               "track_smoke_p39",
			AppliedRuleVersion:       "v_p39_observed",
			ProposalID:               "prop_smoke_p39",
			Period:                   "monthly",
			HitCount:                 6,
			MisjudgmentCount:         0,
			MissingEvidenceCount:     1,
			DegradedCount:            1,
			RiskAlertCount:           1,
			TrendDirection:           model.RuleEffectTrendFlat,
			MetricsJSON:              `{"hit_count":6}`,
			RelatedProposalIDsJSON:   `["prop_smoke_p39"]`,
			RelatedAuditEventIDsJSON: `["audit_smoke_p39"]`,
			RelatedRiskAlertIDsJSON:  `["risk_smoke_p39"]`,
			SafetyNote:               "只读追踪",
			CreatedAt:                now,
			UpdatedAt:                now,
		}); err != nil {
			return err
		}
	}
	if err := riskAlertRepo.UpsertRiskAlert(ctx, repository.RiskAlert{
		AlertID:               "risk_smoke_p39",
		RiskType:              model.RiskTypeDataDegraded,
		Severity:              model.RiskSeverityWarning,
		SOPStatus:             model.RiskSOPActive,
		Symbol:                "510300",
		TriggerSummary:        "P39 source health stale 触发数据降级风险",
		TriggerContextJSON:    `{"source_name":"csindex_extended","freshness":"stale"}`,
		ProhibitedActionsJSON: `["自动交易","外部推送"]`,
		SuggestedActionsJSON:  `["人工复核数据源状态","查看关联决策和审计"]`,
		RelatedDecisionID:     "decision_smoke_p30",
		RelatedReportID:       "daily_report_smoke_p32",
		RelatedAuditEventID:   "audit_smoke_p39",
		LastTriggeredAt:       now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}); err != nil {
		return err
	}
	return appendAuditEventIfMissing(ctx, auditRepo, repository.AuditEvent{AuditEventID: "audit_smoke_p39", RequestID: "req_smoke_p39", DecisionID: "decision_smoke_p30", ProposalID: "prop_smoke_p39", WorkflowType: "local_e2e_smoke", NodeName: "P39SmokeSeed", Actor: "system", Action: "run_local_task", NodeAction: "seed_full_user_journey", Status: "success", InputRefType: "openspec_change", InputRef: "p39-frontend-full-user-journey-e2e", OutputRefType: "fixture", OutputRef: "risk=active;rule=pending_final_confirm;source_health=stale;no_auto_trading", CreatedAt: now})
}

func configuredLocalDate(timezone string) (string, error) {
	name := strings.TrimSpace(timezone)
	if name == "" {
		name = "UTC"
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return "", fmt.Errorf("load daily auto-run timezone: %w", err)
	}
	return time.Now().In(loc).Format(time.DateOnly), nil
}

func appendAuditEventIfMissing(ctx context.Context, repo *appsqlite.AuditRepository, event repository.AuditEvent) error {
	if _, err := repo.GetAuditEvent(ctx, event.AuditEventID); err == nil {
		return nil
	}
	return repo.AppendAuditEvent(ctx, event)
}
