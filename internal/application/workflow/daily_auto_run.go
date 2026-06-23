package workflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

const dailyAutoRunWorkflowType = "daily_auto_run"
const dailyAutoRunTaskVersion = "v1"
const dailyAutoRunMissingPrerequisites = "missing_prerequisites"

// DailyAutoRunOutput 描述一次本地每日自动运行的可查询状态。
type DailyAutoRunOutput struct {
	RunID          string
	IdempotencyKey string
	Status         string
	FailureCode    string
	FailureReason  string
}

// DailyAutoRunner 编排本地每日自动运行；默认关闭，不执行交易或外部推送。
type DailyAutoRunner struct {
	cfg       config.DailyAutoRunConfig
	deps      WorkflowDependencies
	ids       idgen.Generator
	clk       clock.Clock
	nextRunAt string
}

func NewDailyAutoRunner(cfg config.DailyAutoRunConfig, deps WorkflowDependencies) *DailyAutoRunner {
	return &DailyAutoRunner{cfg: cfg, deps: deps, ids: idgen.NewGenerator(), clk: clock.SystemClock{}}
}

// SetNextRunAt configures the scheduler-computed next run timestamp persisted with run states.
func (r *DailyAutoRunner) SetNextRunAt(nextRunAt string) {
	r.nextRunAt = strings.TrimSpace(nextRunAt)
}

func (r *DailyAutoRunner) RunOnce(ctx context.Context, now time.Time) (DailyAutoRunOutput, error) {
	if !r.cfg.Enabled {
		return DailyAutoRunOutput{Status: "disabled"}, nil
	}
	positions, err := r.deps.PortfolioRepo.ListPositions(ctx)
	if err != nil || len(positions) == 0 {
		key := r.idempotencyKey(now, nil)
		if existing, ok, err := r.existingOutput(ctx, key, now); err != nil {
			return DailyAutoRunOutput{}, err
		} else if ok {
			return r.persistReuse(ctx, now, existing)
		}
		return r.persistFailure(ctx, now, nil, dailyAutoRunMissingPrerequisites, "缺少本地持仓，无法生成每日自动运行结果", "prerequisites")
	}
	if len(positions) > r.cfg.MaxSymbols {
		positions = positions[:r.cfg.MaxSymbols]
	}
	key := r.idempotencyKey(now, positions)
	if existing, ok, err := r.existingOutput(ctx, key, now); err != nil {
		return DailyAutoRunOutput{}, err
	} else if ok {
		return r.persistReuse(ctx, now, existing)
	}
	if err := r.persistRunning(ctx, now, positions); err != nil {
		return DailyAutoRunOutput{}, err
	}
	runCtx, cancel := r.runContext(ctx)
	defer cancel()
	for _, position := range positions {
		if err := r.runWithRetry(runCtx, now, key, "market_refresh", func() error {
			_, err := NewMarketRefreshGraphWithDependencies(r.deps).Run(runCtx, MarketRefreshInput{RequestID: stableAutoRunID(key), Symbol: position.Symbol})
			return err
		}); err != nil {
			return r.persistStepFailure(ctx, now, positions, "market_refresh_failed", err, "market_refresh")
		}
		if err := r.runWithRetry(runCtx, now, key, "evidence_refresh", func() error {
			_, err := NewEvidenceVerificationGraphWithDependencies(r.deps).Run(runCtx, EvidenceVerificationInput{RequestID: stableAutoRunID(key), Symbol: position.Symbol, Sources: []string{"official", "exchange"}})
			return err
		}); err != nil {
			return r.persistStepFailure(ctx, now, positions, "evidence_refresh_failed", err, "evidence_refresh")
		}
	}
	var dailyCtx WorkflowContext
	if err := r.runWithRetry(runCtx, now, key, "daily_context", func() error {
		var err error
		dailyCtx, err = r.dailyWorkflowContext(runCtx, stableAutoRunID(key), positions)
		return err
	}); err != nil {
		return r.persistStepFailure(ctx, now, positions, dailyAutoRunMissingPrerequisites, err, "daily_context")
	}
	var dailyOut WorkflowContext
	if err := r.runWithRetry(runCtx, now, key, "daily_discipline", func() error {
		var err error
		dailyOut, err = NewDailyDisciplineGraphWithDependencies(r.deps).Run(runCtx, dailyCtx)
		return err
	}); err != nil {
		return r.persistStepFailure(ctx, now, positions, "daily_discipline_failed", err, "daily_discipline")
	}
	return r.persistSuccess(ctx, now, positions, dailyOut)
}

func (r *DailyAutoRunner) persistRunning(ctx context.Context, now time.Time, positions []repository.Position) error {
	key := r.idempotencyKey(now, positions)
	state := repository.DailyAutoRunState{RunID: stableAutoRunID(key), IdempotencyKey: key, LocalDate: localDate(now, r.cfg.Timezone), Scope: "holdings", SymbolSetHash: symbolSetHash(positions), Status: "running", LastRunAt: now.UTC().Format(time.RFC3339), NextRunAt: r.nextRunAt, CreatedAt: now.UTC().Format(time.RFC3339), UpdatedAt: now.UTC().Format(time.RFC3339)}
	return r.deps.DailyAutoRunRepo.UpsertDailyAutoRunState(ctx, state)
}

func (r *DailyAutoRunner) persistStepFailure(ctx context.Context, now time.Time, positions []repository.Position, code string, err error, step string) (DailyAutoRunOutput, error) {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return r.persistFailure(ctx, now, positions, "timeout", fmt.Sprintf("每日自动运行在步骤 %s 超时", step), "timeout")
	}
	return r.persistFailure(ctx, now, positions, code, err.Error(), step)
}

func (r *DailyAutoRunner) persistFailure(ctx context.Context, now time.Time, positions []repository.Position, code string, reason string, step string) (DailyAutoRunOutput, error) {
	key := r.idempotencyKey(now, positions)
	runID := stableAutoRunID(key)
	state := repository.DailyAutoRunState{RunID: runID, IdempotencyKey: key, LocalDate: localDate(now, r.cfg.Timezone), Scope: "holdings", SymbolSetHash: symbolSetHash(positions), Status: "failed", LastRunAt: now.UTC().Format(time.RFC3339), NextRunAt: r.nextRunAt, FailureCode: code, FailureReason: reason, CreatedAt: now.UTC().Format(time.RFC3339), UpdatedAt: now.UTC().Format(time.RFC3339)}
	if err := r.deps.Transactor.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.DailyAutoRunRepo.UpsertDailyAutoRunState(ctx, state); err != nil {
			return err
		}
		if err := r.upsertDailyDisciplineReport(ctx, repos, state, dailyDisciplineReportStatusForAutoRunFailure(code), "", reason, code, reason); err != nil {
			return err
		}
		diagnostic := dailyAutoRunDiagnostic("failed", step, code, reason)
		if err := repos.NotificationRepo.SaveNotification(ctx, repository.Notification{NotificationID: stableAutoRunNotificationID(key), Type: "daily_auto_run_failed", Severity: "warning", Title: "每日自动运行未完成", Message: reason, SourceType: "daily_auto_run", SourceID: key, CreatedAt: state.UpdatedAt}); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: stableAutoRunAuditID(key), RequestID: runID, WorkflowType: dailyAutoRunWorkflowType, NodeName: "DailyAutoRunner", Actor: string(model.AuditActorSystem), Action: string(model.AuditActionRunLocalTask), NodeAction: "daily_auto_run", Status: string(model.AuditStatusFailed), ErrorCode: code, InputRefType: "idempotency_key", InputRef: key, OutputRefType: "diagnostic", OutputRef: diagnostic, CreatedAt: state.UpdatedAt})
	}); err != nil {
		return DailyAutoRunOutput{}, err
	}
	return DailyAutoRunOutput{RunID: runID, IdempotencyKey: key, Status: "failed", FailureCode: code, FailureReason: reason}, nil
}

func (r *DailyAutoRunner) persistSuccess(ctx context.Context, now time.Time, positions []repository.Position, dailyOut WorkflowContext) (DailyAutoRunOutput, error) {
	key := r.idempotencyKey(now, positions)
	runID := stableAutoRunID(key)
	reportStatus := r.dailyDisciplineReportStatus(ctx, dailyOut)
	state := repository.DailyAutoRunState{RunID: runID, IdempotencyKey: key, LocalDate: localDate(now, r.cfg.Timezone), Scope: "holdings", SymbolSetHash: symbolSetHash(positions), Status: reportStatus, LastRunAt: now.UTC().Format(time.RFC3339), NextRunAt: r.nextRunAt, CreatedAt: now.UTC().Format(time.RFC3339), UpdatedAt: now.UTC().Format(time.RFC3339)}
	if err := r.deps.Transactor.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.DailyAutoRunRepo.UpsertDailyAutoRunState(ctx, state); err != nil {
			return err
		}
		if err := r.upsertDailyDisciplineReport(ctx, repos, state, reportStatus, dailyOut.DecisionID, "今日纪律报告已生成", "", ""); err != nil {
			return err
		}
		if err := repos.NotificationRepo.SaveNotification(ctx, repository.Notification{NotificationID: stableAutoRunNotificationID(key), Type: "daily_auto_run_success", Severity: "info", Title: "每日自动运行已完成", Message: "已完成本地每日刷新、纪律评估、通知和审计记录。", SourceType: "daily_auto_run", SourceID: key, CreatedAt: state.UpdatedAt}); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: stableAutoRunAuditID(key), RequestID: runID, WorkflowType: dailyAutoRunWorkflowType, NodeName: "DailyAutoRunner", Actor: string(model.AuditActorSystem), Action: string(model.AuditActionRunLocalTask), NodeAction: "daily_auto_run", Status: string(model.AuditStatusSuccess), InputRefType: "idempotency_key", InputRef: key, OutputRefType: "safety_boundary", OutputRef: "no_auto_trading", CreatedAt: state.UpdatedAt})
	}); err != nil {
		return DailyAutoRunOutput{}, err
	}
	return DailyAutoRunOutput{RunID: runID, IdempotencyKey: key, Status: reportStatus}, nil
}

func (r *DailyAutoRunner) upsertDailyDisciplineReport(ctx context.Context, repos repository.Repositories, state repository.DailyAutoRunState, status string, decisionID string, summary string, failureCode string, failureReason string) error {
	if repos.DailyDisciplineReportRepo == nil {
		return nil
	}
	if summary == "" {
		summary = failureReason
	}
	if summary == "" {
		summary = status
	}
	report := repository.DailyDisciplineReport{ReportID: stableAutoRunReportID(state.IdempotencyKey), LocalDate: state.LocalDate, Scope: state.Scope, SymbolSetHash: state.SymbolSetHash, SourceType: "auto_run", SourceID: state.IdempotencyKey, DecisionID: decisionID, Status: status, Summary: summary, FailureCode: failureCode, FailureReason: failureReason, CreatedAt: state.CreatedAt, UpdatedAt: state.UpdatedAt}
	return repos.DailyDisciplineReportRepo.UpsertDailyDisciplineReport(ctx, report)
}

func (r *DailyAutoRunner) upsertDailyDisciplineReportFromState(ctx context.Context, repos repository.Repositories, state repository.DailyAutoRunState) error {
	status := state.Status
	if state.Status == "failed" {
		status = dailyDisciplineReportStatusForAutoRunFailure(state.FailureCode)
	}
	summary := "今日纪律报告已生成"
	if state.FailureReason != "" {
		summary = state.FailureReason
	}
	return r.upsertDailyDisciplineReport(ctx, repos, state, status, "", summary, state.FailureCode, state.FailureReason)
}

func dailyDisciplineReportStatusForAutoRunFailure(failureCode string) string {
	switch failureCode {
	case dailyAutoRunMissingPrerequisites:
		return "insufficient_data"
	case "degraded":
		return "degraded"
	default:
		return "failed"
	}
}

// DailyDisciplineReportStatus maps a completed daily discipline workflow output to the report index status.
func DailyDisciplineReportStatus(wf WorkflowContext) string {
	if wf.AnalystUnavailable || workflowStatus(wf) == model.WorkflowDegraded {
		return "degraded"
	}
	return "success"
}

func (r *DailyAutoRunner) dailyDisciplineReportStatus(ctx context.Context, wf WorkflowContext) string {
	if hasDegradedDailyDisciplineAudit(ctx, r.deps.AuditRepo, wf.RequestID) {
		return "degraded"
	}
	if status := DailyDisciplineReportStatus(wf); status == "degraded" {
		return status
	}
	if r.deps.DecisionRepo != nil && wf.DecisionID != "" {
		decision, _, err := r.deps.DecisionRepo.GetDecisionRecord(ctx, wf.DecisionID)
		if err == nil && decision.WorkflowStatus == string(model.WorkflowDegraded) {
			return "degraded"
		}
	}
	return "success"
}

func hasDegradedDailyDisciplineAudit(ctx context.Context, repo repository.AuditRepository, requestID string) bool {
	if repo == nil || requestID == "" {
		return false
	}
	events, err := repo.ListAuditEvents(ctx)
	if err != nil {
		return false
	}
	for _, event := range events {
		if event.RequestID == requestID && event.WorkflowType == WorkflowDailyDiscipline && event.Status == string(model.AuditStatusDegraded) {
			return true
		}
	}
	return false
}

func (r *DailyAutoRunner) dailyWorkflowContext(ctx context.Context, requestID string, positions []repository.Position) (WorkflowContext, error) {
	portfolio, err := r.deps.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
	if err != nil {
		return WorkflowContext{}, fmt.Errorf("缺少账户快照")
	}
	market, err := r.deps.MarketRepo.GetLatestMarketSnapshotBySymbol(ctx, positions[0].Symbol)
	if err != nil {
		return WorkflowContext{}, fmt.Errorf("缺少市场快照")
	}
	rule, err := r.deps.RuleRepo.GetActiveRuleVersion(ctx)
	if err != nil {
		return WorkflowContext{}, fmt.Errorf("缺少生效规则版本")
	}
	dailyPositions := make([]model.Position, 0, len(positions))
	for _, position := range positions {
		dailyPositions = append(dailyPositions, model.Position{PositionID: position.PositionID, Symbol: position.Symbol, Name: position.Name, Quantity: position.Quantity, CostPrice: position.CostPrice, CurrentPrice: position.CurrentPrice, MarketValue: position.MarketValue, UnrealizedProfitRatio: position.UnrealizedProfitRatio, PositionState: model.PositionState(position.PositionState), AssetTag: position.AssetTag})
	}
	return WorkflowContext{RequestID: requestID, WorkflowType: WorkflowDailyDiscipline, Symbol: market.Symbol, RuleVersion: rule.RuleVersion, CapabilityStatus: CapabilityInScope, PortfolioSnapshot: model.PortfolioSnapshot{SnapshotID: portfolio.SnapshotID, Cash: portfolio.Cash, TotalAssets: portfolio.TotalAssets, CashRatio: portfolio.CashRatio, HighRiskRatio: portfolio.HighRiskRatio, PositionCount: portfolio.PositionCount}, PositionSnapshots: dailyPositions, MarketSnapshot: market, AnalystUnavailable: true, ExpectedReturnSampleCount: ExpectedReturnSampleCountFromWorkflowData(dailyPositions, market)}, nil
}

func (r *DailyAutoRunner) existingOutput(ctx context.Context, key string, now time.Time) (DailyAutoRunOutput, bool, error) {
	state, err := r.deps.DailyAutoRunRepo.GetDailyAutoRunState(ctx, key)
	if err != nil {
		if apperr.IsCode(err, apperr.CodeNotFound) {
			return DailyAutoRunOutput{}, false, nil
		}
		return DailyAutoRunOutput{}, false, err
	}
	if state.Status == "running" && r.runningStateIsStale(state, now) {
		return DailyAutoRunOutput{}, false, nil
	}
	return DailyAutoRunOutput{RunID: state.RunID, IdempotencyKey: state.IdempotencyKey, Status: state.Status, FailureCode: state.FailureCode, FailureReason: state.FailureReason}, true, nil
}

func (r *DailyAutoRunner) runningStateIsStale(state repository.DailyAutoRunState, now time.Time) bool {
	updatedAt := strings.TrimSpace(state.UpdatedAt)
	if updatedAt == "" {
		updatedAt = strings.TrimSpace(state.LastRunAt)
	}
	if updatedAt == "" {
		return false
	}
	parsed, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return false
	}
	timeout := time.Duration(r.cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = time.Hour
	}
	return !now.UTC().Before(parsed.UTC().Add(timeout))
}

func (r *DailyAutoRunner) persistReuse(ctx context.Context, now time.Time, existing DailyAutoRunOutput) (DailyAutoRunOutput, error) {
	createdAt := now.UTC().Format(time.RFC3339)
	state, err := r.deps.DailyAutoRunRepo.GetDailyAutoRunState(ctx, existing.IdempotencyKey)
	if err != nil {
		return DailyAutoRunOutput{}, err
	}
	err = r.deps.Transactor.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := r.upsertDailyDisciplineReportFromState(ctx, repos, state); err != nil {
			return err
		}
		err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: stableAutoRunReuseAuditID(existing.IdempotencyKey, createdAt), RequestID: existing.RunID, WorkflowType: dailyAutoRunWorkflowType, NodeName: "DailyAutoRunner", Actor: string(model.AuditActorSystem), Action: string(model.AuditActionRunLocalTask), NodeAction: "daily_auto_run_reuse", Status: string(model.AuditStatusSuccess), InputRefType: "idempotency_key", InputRef: existing.IdempotencyKey, OutputRefType: "diagnostic", OutputRef: dailyAutoRunDiagnostic("reuse", "idempotency", existing.FailureCode, "复用同一日期和范围的已有每日自动运行结果"), CreatedAt: createdAt})
		if apperr.IsCode(err, apperr.CodeConflict) {
			return nil
		}
		return err
	})
	if err != nil {
		return DailyAutoRunOutput{}, err
	}
	return existing, nil
}

func (r *DailyAutoRunner) runContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if r.cfg.TimeoutSeconds <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, time.Duration(r.cfg.TimeoutSeconds)*time.Second)
}

func (r *DailyAutoRunner) runWithRetry(ctx context.Context, now time.Time, key string, step string, fn func() error) error {
	attempts := r.cfg.Retry + 1
	if attempts < 1 {
		attempts = 1
	}
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if errors.Is(lastErr, context.DeadlineExceeded) || errors.Is(lastErr, context.Canceled) || ctx.Err() != nil {
			return lastErr
		}
		if attempt < attempts {
			if err := r.persistRetryAudit(ctx, now, key, step, attempt, lastErr); err != nil {
				return err
			}
		}
	}
	return lastErr
}

func (r *DailyAutoRunner) persistRetryAudit(ctx context.Context, now time.Time, key string, step string, attempt int, err error) error {
	createdAt := now.UTC().Format(time.RFC3339)
	return r.deps.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: stableAutoRunRetryAuditID(key, step, attempt), RequestID: stableAutoRunID(key), WorkflowType: dailyAutoRunWorkflowType, NodeName: "DailyAutoRunner", Actor: string(model.AuditActorSystem), Action: string(model.AuditActionRunLocalTask), NodeAction: "daily_auto_run_retry", Status: string(model.AuditStatusFailed), ErrorCode: "retry", InputRefType: "idempotency_key", InputRef: key, OutputRefType: "diagnostic", OutputRef: dailyAutoRunDiagnostic("retry", step, "retry", err.Error()), CreatedAt: createdAt})
}

func dailyAutoRunDiagnostic(status string, step string, code string, reason string) string {
	parts := []string{"status=" + status, "step=" + step, "safety=no_auto_trading"}
	if code != "" {
		parts = append(parts, "code="+code)
	}
	if reason != "" {
		parts = append(parts, "reason="+reason)
	}
	return strings.Join(parts, ";")
}

func (r *DailyAutoRunner) idempotencyKey(now time.Time, positions []repository.Position) string {
	return fmt.Sprintf("%s:holdings:%s:%s", localDate(now, r.cfg.Timezone), symbolSetHash(positions), dailyAutoRunTaskVersion)
}

func localDate(now time.Time, timezone string) string {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.Local
	}
	return now.In(loc).Format(time.DateOnly)
}

func symbolSetHash(positions []repository.Position) string {
	symbols := make([]string, 0, len(positions))
	for _, position := range positions {
		symbol := strings.TrimSpace(position.Symbol)
		if symbol != "" {
			symbols = append(symbols, symbol)
		}
	}
	sort.Strings(symbols)
	sum := sha256.Sum256([]byte(strings.Join(symbols, ",")))
	return hex.EncodeToString(sum[:8])
}

func stableAutoRunID(key string) string {
	return "auto_run_" + stableHash(key)
}

func stableAutoRunNotificationID(key string) string {
	return "notif_" + stableHash("daily_auto_run", key)
}

func stableAutoRunReportID(key string) string {
	return "report_" + stableHash("daily_auto_run_report", key)
}

func stableAutoRunAuditID(key string) string {
	return "audit_" + stableHash("daily_auto_run", key)
}

func stableAutoRunRetryAuditID(key string, step string, attempt int) string {
	return "audit_" + stableHash("daily_auto_run_retry", key, step, fmt.Sprintf("%d", attempt))
}

func stableAutoRunReuseAuditID(key string, createdAt string) string {
	return "audit_" + stableHash("daily_auto_run_reuse", key, createdAt)
}
