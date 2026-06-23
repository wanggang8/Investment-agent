package service

import (
	"context"
	"regexp"
	"strings"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

const decisionLoopSafetyNote = "只读解释链，仅展示本地事实和导航，不改变事实状态。"

var (
	completePromptPattern = regexp.MustCompile(`(?i)完整\s*` + `prompt`)
	httpRequestPattern    = regexp.MustCompile(`(?i)\b(?:GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\s+\S+\s+HTTP/\d(?:\.\d)?`)
	httpResponsePattern   = regexp.MustCompile(`(?i)HTTP/\d(?:\.\d)?\s+\d{3}[^\n，。；;]*`)
	privateKeyPattern     = regexp.MustCompile(`(?is)-{5}BEGIN [^-]*(?:PRIVATE|OPENSSH|RSA)[^-]*-{5}.*?-{5}END [^-]+-{5}`)
	privatePathRegex      = regexp.MustCompile(`/` + `Users/[^\s,;，；。、"'“”]+`)
	promptLabelPattern    = regexp.MustCompile(`(?i)\bprompt\s*:`)
	rawHTTPPattern        = regexp.MustCompile(`(?i)\braw\s+http\b`)
	secretPattern         = regexp.MustCompile(`(?i)\bsk-[A-Za-z0-9][A-Za-z0-9_-]{8,}\b`)
	sqlSelectAllPattern   = regexp.MustCompile(`(?i)\bselect\s+\*\s+from\b`)
)

// DecisionLoopListFilter bounds the read-only decision loop list query.
type DecisionLoopListFilter struct {
	Symbol string
	Limit  int
}

// DecisionLoopService builds read-only explanations from existing local facts.
type DecisionLoopService struct {
	repos repository.Repositories
}

// NewDecisionLoopService creates the P47 read-only decision loop aggregation service.
func NewDecisionLoopService(repos repository.Repositories) *DecisionLoopService {
	return &DecisionLoopService{repos: repos}
}

// ListDecisionLoops returns recent decision loop explanations without mutating local facts.
func (s *DecisionLoopService) ListDecisionLoops(ctx context.Context, filter DecisionLoopListFilter) (dto.DecisionLoopListResponse, error) {
	limit := boundedDecisionLoopLimit(filter.Limit)
	decisions, err := s.repos.DecisionRepo.ListDecisionRecords(ctx)
	if err != nil {
		return dto.DecisionLoopListResponse{}, err
	}
	errors, audits, risks, err := s.loopFacts(ctx)
	if err != nil {
		return dto.DecisionLoopListResponse{}, err
	}
	items := make([]dto.DecisionLoopItem, 0, limit)
	symbol := strings.TrimSpace(filter.Symbol)
	for _, listed := range decisions {
		if symbol != "" && listed.Symbol != symbol {
			continue
		}
		full, _, err := s.repos.DecisionRepo.GetDecisionRecord(ctx, listed.DecisionID)
		if err != nil {
			return dto.DecisionLoopListResponse{}, err
		}
		item, err := s.buildLoop(ctx, full, errors, audits, risks)
		if err != nil {
			return dto.DecisionLoopListResponse{}, err
		}
		items = append(items, item)
		if len(items) >= limit {
			break
		}
	}
	return dto.DecisionLoopListResponse{Items: items, Total: len(items), SafetyNote: decisionLoopSafetyNote}, nil
}

// GetDecisionLoop returns one read-only decision loop explanation.
func (s *DecisionLoopService) GetDecisionLoop(ctx context.Context, decisionID string) (dto.DecisionLoopItem, error) {
	if strings.TrimSpace(decisionID) == "" {
		return dto.DecisionLoopItem{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "decision_id is required")
	}
	decision, _, err := s.repos.DecisionRepo.GetDecisionRecord(ctx, decisionID)
	if err != nil {
		return dto.DecisionLoopItem{}, err
	}
	errors, audits, risks, err := s.loopFacts(ctx)
	if err != nil {
		return dto.DecisionLoopItem{}, err
	}
	return s.buildLoop(ctx, decision, errors, audits, risks)
}

func (s *DecisionLoopService) loopFacts(ctx context.Context) ([]repository.ErrorCase, []repository.AuditEvent, []repository.RiskAlert, error) {
	var errorCases []repository.ErrorCase
	if s.repos.DecisionRepo != nil {
		items, err := s.repos.DecisionRepo.ListErrorCases(ctx)
		if err != nil {
			return nil, nil, nil, err
		}
		errorCases = items
	}
	var audits []repository.AuditEvent
	if s.repos.AuditRepo != nil {
		items, err := s.repos.AuditRepo.ListAuditEvents(ctx)
		if err != nil {
			return nil, nil, nil, err
		}
		audits = items
	}
	var risks []repository.RiskAlert
	if s.repos.RiskAlertRepo != nil {
		items, err := s.repos.RiskAlertRepo.ListRiskAlerts(ctx, repository.RiskAlertFilter{})
		if err != nil {
			return nil, nil, nil, err
		}
		risks = items
	}
	return errorCases, audits, risks, nil
}

func (s *DecisionLoopService) buildLoop(ctx context.Context, decision repository.DecisionRecord, allErrors []repository.ErrorCase, allAudits []repository.AuditEvent, allRisks []repository.RiskAlert) (dto.DecisionLoopItem, error) {
	confirmations, err := s.repos.DecisionRepo.ListOperationConfirmations(ctx, decision.DecisionID)
	if err != nil {
		return dto.DecisionLoopItem{}, err
	}
	txsByConfirmation := map[string][]repository.PositionTransaction{}
	confirmationIDs := map[string]bool{}
	txIDs := map[string]bool{}
	for _, confirmation := range confirmations {
		confirmationIDs[confirmation.ConfirmationID] = true
		txs, err := s.repos.DecisionRepo.ListPositionTransactionsByConfirmation(ctx, confirmation.ConfirmationID)
		if err != nil {
			return dto.DecisionLoopItem{}, err
		}
		txsByConfirmation[confirmation.ConfirmationID] = txs
		for _, tx := range txs {
			txIDs[tx.TransactionID] = true
		}
	}
	errorCases := filterDecisionLoopErrors(allErrors, decision.DecisionID, confirmationIDs)
	errorIDs := map[string]bool{}
	for _, item := range errorCases {
		errorIDs[item.ErrorCaseID] = true
	}
	riskLinks := decisionLoopRiskLinks(decision, allRisks)
	reviewLinks := decisionLoopReviewLinks(errorCases)
	auditLinks := decisionLoopAuditLinks(decision.DecisionID, confirmationIDs, errorIDs, allAudits)
	manualActions := decisionLoopManualActions(confirmations, txsByConfirmation)
	stages, missing := decisionLoopStages(decision, confirmations, txIDs, riskLinks, reviewLinks, auditLinks)
	return normalizeDecisionLoopItem(dto.DecisionLoopItem{
		DecisionID:         decision.DecisionID,
		Symbol:             decision.Symbol,
		GeneratedAt:        decision.CreatedAt,
		FinalVerdictStatus: decision.FinalVerdictStatus,
		FinalVerdictText:   sanitizeDecisionLoopText(decision.FinalVerdictText),
		ConfirmationStatus: decision.ConfirmationStatus,
		LoopStatus:         decisionLoopStatus(decision, stages, txIDs, reviewLinks, auditLinks),
		Stages:             stages,
		ManualActions:      manualActions,
		RiskLinks:          riskLinks,
		ReviewLinks:        reviewLinks,
		AuditLinks:         auditLinks,
		MissingLinks:       missing,
		SafetyNote:         decisionLoopSafetyNote,
	}), nil
}

func boundedDecisionLoopLimit(limit int) int {
	if limit <= 0 {
		return 10
	}
	if limit > 50 {
		return 50
	}
	return limit
}

func normalizeDecisionLoopItem(item dto.DecisionLoopItem) dto.DecisionLoopItem {
	if item.Stages == nil {
		item.Stages = []dto.DecisionLoopStage{}
	}
	if item.ManualActions == nil {
		item.ManualActions = []dto.DecisionLoopManualAction{}
	}
	if item.RiskLinks == nil {
		item.RiskLinks = []dto.DecisionLoopLink{}
	}
	if item.ReviewLinks == nil {
		item.ReviewLinks = []dto.DecisionLoopLink{}
	}
	if item.AuditLinks == nil {
		item.AuditLinks = []dto.DecisionLoopLink{}
	}
	if item.MissingLinks == nil {
		item.MissingLinks = []string{}
	}
	return item
}

func filterDecisionLoopErrors(items []repository.ErrorCase, decisionID string, confirmationIDs map[string]bool) []repository.ErrorCase {
	out := make([]repository.ErrorCase, 0)
	for _, item := range items {
		if item.DecisionID == decisionID || confirmationIDs[item.ConfirmationID] {
			out = append(out, item)
		}
	}
	return out
}

func decisionLoopRiskLinks(decision repository.DecisionRecord, risks []repository.RiskAlert) []dto.DecisionLoopLink {
	out := make([]dto.DecisionLoopLink, 0)
	for _, risk := range risks {
		if risk.RelatedDecisionID != decision.DecisionID && (risk.RelatedDecisionID != "" || risk.Symbol != decision.Symbol) {
			continue
		}
		out = append(out, dto.DecisionLoopLink{
			Type:   "risk_alert",
			ID:     risk.AlertID,
			Label:  sanitizeDecisionLoopText(risk.TriggerSummary),
			Href:   "/risk-alerts/" + risk.AlertID,
			Status: string(risk.SOPStatus),
		})
	}
	return out
}

func decisionLoopReviewLinks(errors []repository.ErrorCase) []dto.DecisionLoopLink {
	out := make([]dto.DecisionLoopLink, 0, len(errors))
	for _, item := range errors {
		label := "错误案例"
		if item.RootCauseTag != "" {
			label += " · " + item.RootCauseTag
		}
		out = append(out, dto.DecisionLoopLink{
			Type:   "error_case",
			ID:     item.ErrorCaseID,
			Label:  label,
			Href:   "/review#error_case-" + item.ErrorCaseID,
			Status: "reviewed",
		})
	}
	return out
}

func decisionLoopAuditLinks(decisionID string, confirmationIDs map[string]bool, errorIDs map[string]bool, audits []repository.AuditEvent) []dto.DecisionLoopLink {
	out := make([]dto.DecisionLoopLink, 0)
	for _, audit := range audits {
		if audit.DecisionID != decisionID && !confirmationIDs[audit.ConfirmationID] && !errorIDs[audit.ErrorCaseID] {
			continue
		}
		label := "审计事件"
		if audit.Action != "" {
			label += " · " + audit.Action
		}
		out = append(out, dto.DecisionLoopLink{
			Type:   "audit_event",
			ID:     audit.AuditEventID,
			Label:  label,
			Href:   "/audit#audit-" + audit.AuditEventID,
			Status: audit.Status,
		})
	}
	return out
}

func decisionLoopManualActions(confirmations []repository.OperationConfirmation, txsByConfirmation map[string][]repository.PositionTransaction) []dto.DecisionLoopManualAction {
	out := make([]dto.DecisionLoopManualAction, 0, len(confirmations))
	for _, confirmation := range confirmations {
		transactions := txsByConfirmation[confirmation.ConfirmationID]
		ids := make([]string, 0, len(transactions))
		for _, tx := range transactions {
			ids = append(ids, tx.TransactionID)
		}
		out = append(out, dto.DecisionLoopManualAction{
			ConfirmationID:   confirmation.ConfirmationID,
			ConfirmationType: confirmation.ConfirmationType,
			OperationType:    confirmation.OperationType,
			Symbol:           confirmation.Symbol,
			Quantity:         confirmation.Quantity,
			Price:            confirmation.Price,
			Fees:             confirmation.Fees,
			ExecutedAt:       confirmation.ExecutedAt,
			TransactionIDs:   ids,
			NotePreview:      sanitizeDecisionLoopText(confirmation.Note),
		})
	}
	return out
}

func decisionLoopStages(decision repository.DecisionRecord, confirmations []repository.OperationConfirmation, txIDs map[string]bool, riskLinks, reviewLinks, auditLinks []dto.DecisionLoopLink) ([]dto.DecisionLoopStage, []string) {
	stages := []dto.DecisionLoopStage{{
		Stage:   "recommendation",
		Status:  "complete",
		Label:   "建议生成",
		Summary: sanitizeDecisionLoopText(decision.FinalVerdictText),
		RefType: "decision",
		RefID:   decision.DecisionID,
		At:      decision.CreatedAt,
	}}
	var missing []string
	confirmationStage := confirmationDecisionLoopStage(decision, confirmations)
	stages = append(stages, confirmationStage)
	if confirmationStage.Status == "missing" {
		missing = append(missing, "缺少用户确认记录")
	}
	manualStage := manualRecordDecisionLoopStage(decision, confirmations, txIDs)
	stages = append(stages, manualStage)
	if manualStage.Status == "missing" {
		missing = append(missing, "缺少线下记录")
	}
	riskStage := riskDecisionLoopStage(decision, riskLinks)
	stages = append(stages, riskStage)
	if riskStage.Status == "pending" || riskStage.Status == "degraded" {
		missing = append(missing, "缺少风险线索")
	}
	reviewStage := reviewDecisionLoopStage(decision, reviewLinks, auditLinks)
	stages = append(stages, reviewStage)
	if reviewStage.Status == "pending" {
		missing = append(missing, "缺少复盘线索")
	}
	return stages, missing
}

func confirmationDecisionLoopStage(decision repository.DecisionRecord, confirmations []repository.OperationConfirmation) dto.DecisionLoopStage {
	if decision.ConfirmationStatus == string(model.ConfirmationNotRequired) {
		return dto.DecisionLoopStage{Stage: "confirmation", Status: "not_required", Label: "用户记录", Summary: "无需记录用户处理", RefType: "decision", RefID: decision.DecisionID, At: decision.CreatedAt}
	}
	if len(confirmations) > 0 {
		first := confirmations[0]
		return dto.DecisionLoopStage{Stage: "confirmation", Status: "complete", Label: "用户记录", Summary: "已记录线下处理：" + first.ConfirmationType, RefType: "confirmation", RefID: first.ConfirmationID, At: first.CreatedAt}
	}
	if decision.ConfirmationStatus == string(model.ConfirmationPending) {
		return dto.DecisionLoopStage{Stage: "confirmation", Status: "pending", Label: "用户记录", Summary: "等待用户线下处理记录", RefType: "decision", RefID: decision.DecisionID, At: decision.CreatedAt}
	}
	return dto.DecisionLoopStage{Stage: "confirmation", Status: "missing", Label: "用户记录", Summary: "缺少用户确认记录", RefType: "decision", RefID: decision.DecisionID, At: decision.CreatedAt}
}

func manualRecordDecisionLoopStage(decision repository.DecisionRecord, confirmations []repository.OperationConfirmation, txIDs map[string]bool) dto.DecisionLoopStage {
	if len(txIDs) > 0 {
		return dto.DecisionLoopStage{Stage: "manual_record", Status: "complete", Label: "线下记录", Summary: "已关联本地流水记录", RefType: "transaction", RefID: firstDecisionLoopMapKey(txIDs)}
	}
	if decision.ConfirmationStatus == string(model.ConfirmationPlanned) || decision.ConfirmationStatus == string(model.ConfirmationWatch) || decision.ConfirmationStatus == string(model.ConfirmationMarkedError) || decision.ConfirmationStatus == string(model.ConfirmationNotRequired) {
		return dto.DecisionLoopStage{Stage: "manual_record", Status: "not_required", Label: "线下记录", Summary: "当前处理状态无需本地流水"}
	}
	for _, confirmation := range confirmations {
		if confirmation.ConfirmationType == string(model.ConfirmationTypePlanned) || confirmation.ConfirmationType == string(model.ConfirmationTypeWatch) || confirmation.ConfirmationType == string(model.ConfirmationTypeMarkedError) {
			return dto.DecisionLoopStage{Stage: "manual_record", Status: "not_required", Label: "线下记录", Summary: "当前确认类型无需本地流水", RefType: "confirmation", RefID: confirmation.ConfirmationID, At: confirmation.CreatedAt}
		}
	}
	if decision.ConfirmationStatus == string(model.ConfirmationExecutedManually) {
		return dto.DecisionLoopStage{Stage: "manual_record", Status: "missing", Label: "线下记录", Summary: "缺少本地流水记录"}
	}
	return dto.DecisionLoopStage{Stage: "manual_record", Status: "pending", Label: "线下记录", Summary: "等待确认后关联本地流水"}
}

func riskDecisionLoopStage(decision repository.DecisionRecord, links []dto.DecisionLoopLink) dto.DecisionLoopStage {
	if len(links) > 0 {
		return dto.DecisionLoopStage{Stage: "risk_review", Status: "complete", Label: "风险线索", Summary: "已关联风险线索", RefType: links[0].Type, RefID: links[0].ID}
	}
	if decision.WorkflowStatus == string(model.WorkflowDegraded) || decision.SourceVerificationStatus == string(model.VerificationFailed) || decision.SourceVerificationStatus == string(model.VerificationBackgroundOnly) || decision.FinalVerdictStatus == string(model.VerdictInsufficientData) || decision.FinalVerdictStatus == string(model.VerdictFrozenWatch) || decision.DashboardState == string(model.DashboardHighRisk) {
		return dto.DecisionLoopStage{Stage: "risk_review", Status: "pending", Label: "风险线索", Summary: "存在降级或高风险状态，尚未关联风险线索"}
	}
	return dto.DecisionLoopStage{Stage: "risk_review", Status: "not_required", Label: "风险线索", Summary: "暂无额外风险线索"}
}

func reviewDecisionLoopStage(decision repository.DecisionRecord, reviewLinks, auditLinks []dto.DecisionLoopLink) dto.DecisionLoopStage {
	if len(reviewLinks) > 0 {
		return dto.DecisionLoopStage{Stage: "review", Status: "complete", Label: "复盘线索", Summary: "已关联复盘线索", RefType: reviewLinks[0].Type, RefID: reviewLinks[0].ID}
	}
	if len(auditLinks) > 0 {
		return dto.DecisionLoopStage{Stage: "review", Status: "complete", Label: "复盘线索", Summary: "已关联审计线索", RefType: auditLinks[0].Type, RefID: auditLinks[0].ID}
	}
	if decision.ConfirmationStatus == string(model.ConfirmationNotRequired) || decision.ConfirmationStatus == string(model.ConfirmationPlanned) || decision.ConfirmationStatus == string(model.ConfirmationWatch) {
		return dto.DecisionLoopStage{Stage: "review", Status: "not_required", Label: "复盘线索", Summary: "当前状态暂无复盘要求"}
	}
	return dto.DecisionLoopStage{Stage: "review", Status: "pending", Label: "复盘线索", Summary: "尚未关联复盘或审计线索"}
}

func decisionLoopStatus(decision repository.DecisionRecord, stages []dto.DecisionLoopStage, txIDs map[string]bool, reviewLinks, auditLinks []dto.DecisionLoopLink) string {
	for _, stage := range stages {
		if stage.Status == "missing" {
			return "incomplete"
		}
	}
	if len(reviewLinks) > 0 || len(auditLinks) > 0 || decision.ConfirmationStatus == string(model.ConfirmationNotRequired) || decision.ConfirmationStatus == string(model.ConfirmationMarkedError) {
		return "reviewed"
	}
	if len(txIDs) > 0 {
		return "recorded"
	}
	if decision.ConfirmationStatus == string(model.ConfirmationPlanned) {
		return "planned"
	}
	return "open"
}

func firstDecisionLoopMapKey(items map[string]bool) string {
	for key := range items {
		return key
	}
	return ""
}

func sanitizeDecisionLoopText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	text = privateKeyPattern.ReplaceAllString(text, "[PRIVATE_BLOCK]")
	text = sqlSelectAllPattern.ReplaceAllString(text, "SELECT [REDACTED] FROM")
	text = promptLabelPattern.ReplaceAllString(text, "[PROMPT]")
	text = completePromptPattern.ReplaceAllString(text, "[PROMPT]")
	text = rawHTTPPattern.ReplaceAllString(text, "[HTTP_RAW]")
	text = httpRequestPattern.ReplaceAllString(text, "[HTTP_REQUEST]")
	text = httpResponsePattern.ReplaceAllString(text, "[HTTP_RESPONSE]")
	text = privatePathRegex.ReplaceAllString(text, "[LOCAL_PATH]")
	text = secretPattern.ReplaceAllString(text, "[REDACTED]")
	if len([]rune(text)) <= 120 {
		return text
	}
	runes := []rune(text)
	return string(runes[:120]) + "..."
}
