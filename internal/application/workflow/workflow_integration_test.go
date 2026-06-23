package workflow

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/idgen"

	_ "modernc.org/sqlite"
)

func workflowTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", t.TempDir()+"/workflow.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := appsqlite.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	return db
}

func workflowDependenciesForDB(db *sql.DB) WorkflowDependencies {
	repos := repository.Repositories{
		DecisionRepo:              appsqlite.NewDecisionRepository(db),
		AuditRepo:                 appsqlite.NewAuditRepository(db),
		RuleRepo:                  appsqlite.NewRuleRepository(db),
		MarketRepo:                appsqlite.NewMarketRepository(db),
		IntelligenceRepo:          appsqlite.NewIntelligenceRepository(db),
		NotificationRepo:          appsqlite.NewNotificationRepository(db),
		PortfolioRepo:             appsqlite.NewPortfolioRepository(db),
		DailyAutoRunRepo:          appsqlite.NewDailyAutoRunRepository(db),
		DailyDisciplineReportRepo: appsqlite.NewDailyDisciplineReportRepository(db),
	}
	return NewWorkflowDependencies(repos, appsqlite.NewTransactor(db))
}

func TestDailyGraphPersistsDecisionEvidenceAndAudit(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	out, err := NewDailyDisciplineGraphWithDependencies(deps).Run(context.Background(), sampleWorkflowContext(20))
	if err != nil {
		t.Fatalf("run daily graph: %v", err)
	}

	var decisionCount, evidenceCount, auditCount int
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM decision_records WHERE decision_id=? AND expected_return_scenarios_json IS NOT NULL`, out.DecisionID).Scan(&decisionCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM evidence_refs WHERE decision_id=?`, out.DecisionID).Scan(&evidenceCount); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM audit_events WHERE request_id=? AND node_name='DecisionRecordNode' AND input_ref_type<>'' AND output_ref_type<>''`, out.RequestID).Scan(&auditCount); err != nil {
		t.Fatal(err)
	}
	if decisionCount != 1 || evidenceCount == 0 || auditCount == 0 {
		t.Fatalf("expected persisted decision/evidence/audit, got decision=%d evidence=%d audit=%d", decisionCount, evidenceCount, auditCount)
	}
}

func TestDailyGraphStopsOnMissingState(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := sampleWorkflowContext(20)
	ctx.PortfolioSnapshot.SnapshotID = ""
	out, err := NewDailyDisciplineGraphWithDependencies(deps).Run(context.Background(), ctx)
	if err != nil {
		t.Fatalf("run daily graph: %v", err)
	}
	var decisionCount int
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM decision_records WHERE request_id=?`, out.RequestID).Scan(&decisionCount); err != nil {
		t.Fatal(err)
	}
	if decisionCount != 0 || !hasString(out.Errors, ErrCodeDataRequired) {
		t.Fatalf("expected no decision and DATA_REQUIRED, got decisions=%d errors=%+v", decisionCount, out.Errors)
	}
}

func TestSupportGraphsPersistFacts(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()

	if _, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(ctx, EvidenceVerificationInput{RequestID: "req_ev_db", Symbol: "510300", Sources: []string{"official", "exchange"}}); err != nil {
		t.Fatalf("evidence graph: %v", err)
	}
	if _, err := NewMarketRefreshGraphWithDependencies(deps).Run(ctx, MarketRefreshInput{RequestID: "req_mk_db", Symbol: "510300", PEPercentile: 45, PBPercentile: 40}); err != nil {
		t.Fatalf("market graph: %v", err)
	}
	evoOut, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_db", ErrorCaseID: "err_1", SampleCount: 5})
	if err != nil {
		t.Fatalf("evolution graph: %v", err)
	}
	if err := deps.RuleRepo.UpdateRuleProposalStatus(ctx, evoOut.RuleProposal.ProposalID, string(model.ProposalUnderGatekeeperAudit)); err != nil {
		t.Fatalf("mark under audit: %v", err)
	}
	if _, err := NewGatekeeperAuditGraphWithDependencies(deps).Run(ctx, GatekeeperAuditInput{RequestID: "req_gate_db", ProposalID: evoOut.RuleProposal.ProposalID, Approved: true}); err != nil {
		t.Fatalf("gatekeeper graph: %v", err)
	}

	counts := map[string]int{}
	for _, table := range []string{"intelligence_items", "intelligence_summary", "rag_chunks", "source_verifications", "market_snapshots", "rule_proposals", "gatekeeper_audits"} {
		var count int
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+table).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		counts[table] = count
	}
	for table, count := range counts {
		if count == 0 {
			t.Fatalf("expected %s to have rows, counts=%+v", table, counts)
		}
	}
}

func TestEvolutionProposalGraphCreatesUserConfirmableProposal(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()

	out, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_confirmable", ErrorCaseID: "err_confirmable", SampleCount: 5})
	if err != nil {
		t.Fatalf("evolution graph: %v", err)
	}
	if out.RuleProposal.Status != model.ProposalPendingUserConfirm {
		t.Fatalf("expected pending user confirm output, got %s", out.RuleProposal.Status)
	}
	var status string
	if err := db.QueryRowContext(ctx, `SELECT status FROM rule_proposals WHERE proposal_id=?`, out.RuleProposal.ProposalID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(model.ProposalPendingUserConfirm) {
		t.Fatalf("expected persisted pending user confirm, got %s", status)
	}
	var notificationType, sourceType, sourceID string
	if err := db.QueryRowContext(ctx, `SELECT type,COALESCE(source_type,''),COALESCE(source_id,'') FROM notifications WHERE read_at IS NULL ORDER BY created_at DESC LIMIT 1`).Scan(&notificationType, &sourceType, &sourceID); err != nil {
		t.Fatal(err)
	}
	if notificationType != "rule_proposal_pending" || sourceType != "rule_proposal" || sourceID != out.RuleProposal.ProposalID {
		t.Fatalf("expected pending rule proposal notification, got type=%s source=%s/%s", notificationType, sourceType, sourceID)
	}
}

func TestEvolutionProposalGraphPersistsConcreteRuleChangePayload(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()

	out, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_payload", ErrorCaseID: "err_payload", ReviewPeriod: "monthly:2026-05", SupportingDecisionIDs: []string{"dec_payload"}, SupportingAuditEventIDs: []string{"audit_payload"}, SampleCount: 5})
	if err != nil {
		t.Fatalf("evolution graph: %v", err)
	}

	var proposalType, beforeRule, afterRule, reason, impactScope, riskNotes string
	if err := db.QueryRowContext(ctx, `SELECT proposal_type,before_rule_json,after_rule_json,reason,COALESCE(impact_scope_json,''),COALESCE(risk_notes_json,'') FROM rule_proposals WHERE proposal_id=?`, out.RuleProposal.ProposalID).Scan(&proposalType, &beforeRule, &afterRule, &reason, &impactScope, &riskNotes); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"source_verification", "target_rule", "required_high_grade_sources", "err_payload", "dec_payload", "audit_payload"} {
		if !strings.Contains(proposalType+beforeRule+afterRule+reason+impactScope+riskNotes, want) {
			t.Fatalf("expected concrete proposal payload to contain %q, got type=%s before=%s after=%s reason=%s impact=%s risk=%s", want, proposalType, beforeRule, afterRule, reason, impactScope, riskNotes)
		}
	}
}

func TestEvolutionProposalGraphSupportsP75ProposalTypes(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	SetWorkflowIDGenerator(idgen.NewFixedGenerator(map[string][]string{
		"proposal": {"proposal_threshold", "proposal_sop", "proposal_master", "proposal_behavior"},
		"notif":    {"notif_threshold", "notif_sop", "notif_master", "notif_behavior"},
		"audit":    {"audit_threshold", "audit_sop", "audit_master", "audit_behavior"},
	}))
	defer SetWorkflowIDGenerator(idgen.NewGenerator())

	cases := []struct {
		name         string
		proposalType string
		wantDBType   string
		targetRule   string
		wantTitle    string
	}{
		{name: "threshold adjustment", proposalType: "threshold_adjustment", wantDBType: "threshold", targetRule: "valuation_threshold", wantTitle: "阈值调整提案"},
		{name: "sop addition", proposalType: "sop_addition", wantDBType: "sop", targetRule: "risk_sop.panic_sell", wantTitle: "SOP 增补提案"},
		{name: "master weight", proposalType: "master_weight_adjustment", wantDBType: "capability", targetRule: "master.graham.margin_of_safety", wantTitle: "大师权重调整提案"},
		{name: "behavior alert", proposalType: "behavior_pattern_alert", wantDBType: "risk_rule", targetRule: "emotion_bias", wantTitle: "个人行为模式预警提案"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_" + tc.proposalType, ErrorCaseID: "err_" + tc.proposalType, ProposalType: tc.proposalType, TargetRule: tc.targetRule, SampleCount: 5})
			if err != nil {
				t.Fatalf("evolution graph: %v", err)
			}
			var proposalType, title, afterRule string
			if err := db.QueryRowContext(ctx, `SELECT proposal_type,title,after_rule_json FROM rule_proposals WHERE proposal_id=?`, out.RuleProposal.ProposalID).Scan(&proposalType, &title, &afterRule); err != nil {
				t.Fatal(err)
			}
			if proposalType != tc.wantDBType || title != tc.wantTitle || !strings.Contains(afterRule, tc.proposalType) || !strings.Contains(afterRule, tc.targetRule) {
				t.Fatalf("unexpected proposal type payload: type=%s title=%s after=%s", proposalType, title, afterRule)
			}
		})
	}
}

func TestEvolutionProposalGraphPersistsReviewSourceMetadataAndAuditRefs(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()

	out, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_review_source", ErrorCaseID: "err_review", ReviewPeriod: "monthly:2026-05", SupportingDecisionIDs: []string{"dec_1", "dec_2"}, SupportingAuditEventIDs: []string{"audit_1"}, SampleCount: 5})
	if err != nil {
		t.Fatalf("evolution graph: %v", err)
	}

	var sourceErrorCaseID, relatedJSON string
	if err := db.QueryRowContext(ctx, `SELECT source_error_case_id,related_error_cases_json FROM rule_proposals WHERE proposal_id=?`, out.RuleProposal.ProposalID).Scan(&sourceErrorCaseID, &relatedJSON); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"err_review", "monthly:2026-05", "dec_1", "dec_2", "audit_1"} {
		if !strings.Contains(sourceErrorCaseID+relatedJSON, want) {
			t.Fatalf("expected proposal source metadata %q in source=%s related=%s", want, sourceErrorCaseID, relatedJSON)
		}
	}

	var inputRefType, inputRef, outputRefType, outputRef string
	if err := db.QueryRowContext(ctx, `SELECT input_ref_type,input_ref,output_ref_type,output_ref FROM audit_events WHERE request_id=? AND action=?`, "req_evo_review_source", string(model.AuditActionCreateProposal)).Scan(&inputRefType, &inputRef, &outputRefType, &outputRef); err != nil {
		t.Fatal(err)
	}
	if inputRefType != "review_summary" || !strings.Contains(inputRef, "monthly:2026-05") || outputRefType != "rule_proposal" || outputRef != out.RuleProposal.ProposalID {
		t.Fatalf("unexpected audit refs inputType=%s input=%s outputType=%s output=%s", inputRefType, inputRef, outputRefType, outputRef)
	}
}

func TestEvolutionProposalGraphDoesNotCreateActiveRuleVersionForReviewProposal(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()

	if _, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_no_apply", ErrorCaseID: "err_no_apply", ReviewPeriod: "quarterly:2026-Q2", SampleCount: 5}); err != nil {
		t.Fatalf("evolution graph: %v", err)
	}

	var activeFromProposalCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM rule_versions WHERE status='active' AND COALESCE(created_from_proposal_id,'')<>''`).Scan(&activeFromProposalCount); err != nil {
		t.Fatal(err)
	}
	if activeFromProposalCount != 0 {
		t.Fatalf("proposal generation must not create active rule version, got %d", activeFromProposalCount)
	}
}

func TestEvolutionProposalGraphMarksInsufficientSourceAsDraftOnly(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()

	out, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_low_source", ReviewPeriod: "monthly:2026-05", SampleCount: 1})
	if err != nil {
		t.Fatalf("evolution graph: %v", err)
	}

	var status, riskNotes string
	if err := db.QueryRowContext(ctx, `SELECT status,risk_notes_json FROM rule_proposals WHERE proposal_id=?`, out.RuleProposal.ProposalID).Scan(&status, &riskNotes); err != nil {
		t.Fatal(err)
	}
	if status != string(model.ProposalDraft) || !strings.Contains(riskNotes, "INSUFFICIENT_SAMPLE") {
		t.Fatalf("expected draft proposal with insufficient sample note, status=%s notes=%s", status, riskNotes)
	}
}

func TestEvolutionProposalGraphMissingSourceCreatesAuditedDraft(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()

	out, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_missing_source", SampleCount: 5})
	if err != nil {
		t.Fatalf("evolution graph: %v", err)
	}

	var status, inputRefType, inputRef string
	if err := db.QueryRowContext(ctx, `SELECT status FROM rule_proposals WHERE proposal_id=?`, out.RuleProposal.ProposalID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if err := db.QueryRowContext(ctx, `SELECT input_ref_type,input_ref FROM audit_events WHERE request_id=? AND action=?`, "req_evo_missing_source", string(model.AuditActionCreateProposal)).Scan(&inputRefType, &inputRef); err != nil {
		t.Fatal(err)
	}
	if status != string(model.ProposalDraft) || inputRefType != "missing_source" || inputRef == "" {
		t.Fatalf("expected audited draft for missing source, status=%s inputType=%s input=%s", status, inputRefType, inputRef)
	}
	var activeFromProposalCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM rule_versions WHERE status='active' AND COALESCE(created_from_proposal_id,'')<>''`).Scan(&activeFromProposalCount); err != nil {
		t.Fatal(err)
	}
	if activeFromProposalCount != 0 {
		t.Fatalf("missing source must not create active rule version, got %d", activeFromProposalCount)
	}
}

func TestEvidenceFailurePersistsFailedVerification(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	if _, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_ev_fail", Symbol: "510500", Sources: []string{"single"}}); err != nil {
		t.Fatalf("evidence graph: %v", err)
	}
	var status, role string
	if err := db.QueryRowContext(context.Background(), `SELECT verification_status,evidence_role FROM source_verifications WHERE symbol='510500'`).Scan(&status, &role); err != nil {
		t.Fatal(err)
	}
	if status != string(model.VerificationFailed) || role != string(model.EvidenceBackground) {
		t.Fatalf("unexpected verification status=%s role=%s", status, role)
	}
}

func TestEvidenceSourceFailureDoesNotPersistPlaceholderFacts(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	deps.IntelligenceSource = testIntelligenceSource{err: apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "source failed")}

	out, err := NewEvidenceVerificationGraphWithDependencies(deps).Run(context.Background(), EvidenceVerificationInput{RequestID: "req_ev_source_failed", Symbol: "510800", Sources: []string{"official", "exchange"}})
	if err != nil {
		t.Fatalf("evidence graph: %v", err)
	}
	if len(out.IntelligenceItems) != 0 || len(out.RAGChunks) != 0 {
		t.Fatalf("source failure must not persist placeholder facts: %+v", out)
	}
	var itemCount, summaryCount, chunkCount int
	for table, dest := range map[string]*int{"intelligence_items": &itemCount, "intelligence_summary": &summaryCount, "rag_chunks": &chunkCount} {
		if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM `+table+` WHERE rowid > 0`).Scan(dest); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
	}
	if itemCount != 0 || summaryCount != 0 || chunkCount != 0 {
		t.Fatalf("expected no placeholder facts, got items=%d summaries=%d chunks=%d", itemCount, summaryCount, chunkCount)
	}
	var status string
	if err := db.QueryRowContext(context.Background(), `SELECT verification_status FROM source_verifications WHERE symbol='510800'`).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(model.VerificationFailed) {
		t.Fatalf("expected failed verification, got %s", status)
	}
	var auditStatus string
	if err := db.QueryRowContext(context.Background(), `SELECT status FROM audit_events WHERE request_id='req_ev_source_failed' AND node_name='NewsFetchNode'`).Scan(&auditStatus); err != nil {
		t.Fatal(err)
	}
	if auditStatus != string(model.AuditStatusFailed) {
		t.Fatalf("expected failed fetch audit, got %s", auditStatus)
	}
}

func TestGatekeeperPersistsProposalStatus(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	out, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_status", ErrorCaseID: "err_status", SampleCount: 5})
	if err != nil {
		t.Fatalf("evolution graph: %v", err)
	}
	if err := deps.RuleRepo.UpdateRuleProposalStatus(ctx, out.RuleProposal.ProposalID, string(model.ProposalUnderGatekeeperAudit)); err != nil {
		t.Fatalf("mark under audit: %v", err)
	}
	if _, err := db.ExecContext(ctx, `UPDATE rule_proposals SET after_rule_json=? WHERE proposal_id=?`, `{"rule":"updated"}`, out.RuleProposal.ProposalID); err != nil {
		t.Fatalf("mark effective change: %v", err)
	}
	if _, err := NewGatekeeperAuditGraphWithDependencies(deps).Run(ctx, GatekeeperAuditInput{RequestID: "req_gate_status", ProposalID: out.RuleProposal.ProposalID, Approved: true}); err != nil {
		t.Fatalf("gatekeeper graph: %v", err)
	}
	var status string
	if err := db.QueryRowContext(ctx, `SELECT status FROM rule_proposals WHERE proposal_id=?`, out.RuleProposal.ProposalID).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(model.ProposalPendingFinalConfirm) {
		t.Fatalf("status=%s", status)
	}
}

func TestGatekeeperAuditPersistsDetailedReviewReason(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	out, err := NewEvolutionProposalGraphWithDependencies(deps).Run(ctx, EvolutionProposalInput{RequestID: "req_evo_reason", ErrorCaseID: "err_reason", SampleCount: 5})
	if err != nil {
		t.Fatalf("evolution graph: %v", err)
	}
	if err := deps.RuleRepo.UpdateRuleProposalStatus(ctx, out.RuleProposal.ProposalID, string(model.ProposalUnderGatekeeperAudit)); err != nil {
		t.Fatalf("mark under audit: %v", err)
	}
	if _, err := NewGatekeeperAuditGraphWithDependencies(deps).Run(ctx, GatekeeperAuditInput{RequestID: "req_gate_reason", ProposalID: out.RuleProposal.ProposalID, Approved: true}); err != nil {
		t.Fatalf("gatekeeper graph: %v", err)
	}
	var reason, metrics string
	if err := db.QueryRowContext(ctx, `SELECT audit_reason,backtest_metrics_json FROM gatekeeper_audits WHERE proposal_id=?`, out.RuleProposal.ProposalID).Scan(&reason, &metrics); err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{"FundamentalRuleCheck", "ConflictCheck", "Backtest", "AuditDecision"} {
		if !strings.Contains(reason, part) {
			t.Fatalf("expected %s in audit reason %q", part, reason)
		}
	}
	if !strings.Contains(metrics, `"sample_count":5`) || !strings.Contains(metrics, `"passed":true`) {
		t.Fatalf("unexpected backtest metrics %s", metrics)
	}
}

func TestMarketRefreshGraphRollsBackSnapshotWhenAuditWriteFails(t *testing.T) {
	db := workflowTestDB(t)
	transactor := appsqlite.NewTransactor(db)
	deps := NewWorkflowDependencies(repository.Repositories{
		AuditRepo:  appsqlite.NewAuditRepository(db),
		MarketRepo: appsqlite.NewMarketRepository(db),
	}, transactor)
	oldAuditGen := auditIDGen
	oldWorkflowIDGen := workflowIDGen
	defer func() {
		auditIDGen = oldAuditGen
		workflowIDGen = oldWorkflowIDGen
	}()
	auditIDGen = idgen.NewFixedGenerator(map[string][]string{"audit": {"duplicate_audit"}})
	workflowIDGen = idgen.NewFixedGenerator(map[string][]string{"market": {"market_one"}})

	if _, err := db.ExecContext(context.Background(), `INSERT INTO audit_events (audit_event_id,actor,action,status,created_at) VALUES ('duplicate_audit','system','refresh_market_data','success','2026-05-29T00:00:00Z')`); err != nil {
		t.Fatalf("seed audit event: %v", err)
	}
	_, err := NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: "req_market_audit_failed", Symbol: "510300", PEPercentile: 45, PBPercentile: 40})
	if err == nil {
		t.Fatal("expected audit failure")
	}
	var snapshotCount int
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM market_snapshots WHERE symbol='510300'`).Scan(&snapshotCount); err != nil {
		t.Fatal(err)
	}
	if snapshotCount != 0 {
		t.Fatalf("expected market snapshot rollback, got %d", snapshotCount)
	}
}

func TestGatekeeperGraphRejectsInvalidProposalState(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `INSERT INTO rule_proposals (proposal_id,proposal_type,status,title,proposal_version,before_rule_json,after_rule_json,sample_count,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "prop_draft_gate", "risk_rule", string(model.ProposalDraft), "草稿提案", "draft", "{}", "{}", 5, "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed draft proposal: %v", err)
	}

	_, err := NewGatekeeperAuditGraphWithDependencies(deps).Run(ctx, GatekeeperAuditInput{RequestID: "req_gate_draft", ProposalID: "prop_draft_gate", Approved: true})
	if err == nil {
		t.Fatal("expected invalid proposal state error")
	}
	var status string
	if err := db.QueryRowContext(ctx, `SELECT status FROM rule_proposals WHERE proposal_id='prop_draft_gate'`).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(model.ProposalDraft) {
		t.Fatalf("expected draft status unchanged, got %s", status)
	}
	var auditCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gatekeeper_audits WHERE proposal_id='prop_draft_gate'`).Scan(&auditCount); err != nil {
		t.Fatal(err)
	}
	if auditCount != 0 {
		t.Fatalf("expected no gatekeeper audit for invalid state, got %d", auditCount)
	}
}

func TestGatekeeperGraphRejectsInsufficientSamplesBeforeWritingAudit(t *testing.T) {
	db := workflowTestDB(t)
	deps := workflowDependenciesForDB(db)
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `INSERT INTO rule_proposals (proposal_id,proposal_type,status,title,proposal_version,before_rule_json,after_rule_json,sample_count,created_at) VALUES (?,?,?,?,?,?,?,?,?)`, "prop_low_sample_gate", "risk_rule", string(model.ProposalUnderGatekeeperAudit), "样本不足提案", "draft", "{}", "{}", 2, "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed low sample proposal: %v", err)
	}

	_, err := NewGatekeeperAuditGraphWithDependencies(deps).Run(ctx, GatekeeperAuditInput{RequestID: "req_gate_low_sample", ProposalID: "prop_low_sample_gate", Approved: true})
	if err == nil {
		t.Fatal("expected insufficient samples error")
	}
	var status string
	if err := db.QueryRowContext(ctx, `SELECT status FROM rule_proposals WHERE proposal_id='prop_low_sample_gate'`).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(model.ProposalUnderGatekeeperAudit) {
		t.Fatalf("expected status unchanged, got %s", status)
	}
	var auditCount int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gatekeeper_audits WHERE proposal_id='prop_low_sample_gate'`).Scan(&auditCount); err != nil {
		t.Fatal(err)
	}
	if auditCount != 0 {
		t.Fatalf("expected no gatekeeper audit for low sample proposal, got %d", auditCount)
	}
}

func TestGatekeeperGraphRollsBackAuditWhenProposalStatusUpdateFails(t *testing.T) {
	db := workflowTestDB(t)
	transactor := appsqlite.NewTransactor(db)
	deps := NewWorkflowDependencies(repository.Repositories{
		AuditRepo: appsqlite.NewAuditRepository(db),
		RuleRepo:  appsqlite.NewRuleRepository(db),
	}, transactor)

	_, err := NewGatekeeperAuditGraphWithDependencies(deps).Run(context.Background(), GatekeeperAuditInput{RequestID: "req_gate_missing_proposal", ProposalID: "missing_proposal", Approved: true})
	if err == nil {
		t.Fatal("expected status update failure")
	}
	var auditCount int
	if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM gatekeeper_audits WHERE proposal_id='missing_proposal'`).Scan(&auditCount); err != nil {
		t.Fatal(err)
	}
	if auditCount != 0 {
		t.Fatalf("expected gatekeeper audit rollback, got %d", auditCount)
	}
}
