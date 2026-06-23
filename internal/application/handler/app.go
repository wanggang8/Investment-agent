package handler

import (
	"net/http"

	"investment-agent/internal/application/service"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/repository"
)

// App 汇总 P4 HTTP handler 所需依赖。
type App struct {
	Deps                  workflow.WorkflowDependencies
	ConfirmationSvc       service.ConfirmationWriter
	PortfolioSvc          service.PortfolioWriter
	RuleSvc               service.RuleProposalWriter
	SettingsSvc           *service.SettingsService
	MarketSvc             *service.MarketService
	EvidenceSvc           *service.EvidenceService
	LocalKnowledgeSvc     *service.LocalKnowledgeService
	NotificationSvc       *service.NotificationService
	RiskAlertSvc          *service.RiskAlertService
	RuleEffectSvc         *service.RuleEffectValidationService
	VectorIndex           service.VectorIndex
	QuerySvc              *service.QueryService
	DecisionLoopSvc       *service.DecisionLoopService
	DataSourceQualitySvc  *service.DataSourceQualityService
	KnowledgeReadinessSvc *service.KnowledgeReadinessService
}

// NewApp 创建 HTTP 应用层入口。
func NewApp(deps workflow.WorkflowDependencies, repos repository.Repositories, tx repository.Transactor) *App {
	app := &App{
		Deps:                  deps,
		ConfirmationSvc:       service.NewConfirmationService(tx),
		PortfolioSvc:          service.NewPortfolioService(tx),
		RuleSvc:               service.NewRuleProposalService(tx, deps),
		SettingsSvc:           service.NewSettingsService(tx),
		MarketSvc:             service.NewMarketService(tx),
		EvidenceSvc:           service.NewEvidenceService(tx),
		LocalKnowledgeSvc:     service.NewLocalKnowledgeService(tx),
		NotificationSvc:       service.NewNotificationService(tx),
		RiskAlertSvc:          service.NewRiskAlertService(tx),
		RuleEffectSvc:         service.NewRuleEffectValidationService(tx),
		QuerySvc:              service.NewQueryServiceWithDailyAutoRunConfig(repos, deps.DailyAutoRunConfig),
		DecisionLoopSvc:       service.NewDecisionLoopService(repos),
		DataSourceQualitySvc:  service.NewDataSourceQualityService(repos, tx),
		KnowledgeReadinessSvc: service.NewKnowledgeReadinessService(repos),
	}
	if provider, ok := deps.RetrievalService.(interface{ VectorIndex() service.VectorIndex }); ok {
		app.VectorIndex = provider.VectorIndex()
	}
	return app
}

// RegisterRoutes 注册 P4 业务 API。健康检查仍由 cmd/server 保留简单响应。
func (a *App) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/dashboard/today", a.GetDashboardToday)
	mux.HandleFunc("POST /api/v1/decisions/consult", a.ConsultDecision)
	mux.HandleFunc("GET /api/v1/decisions/{decision_id}", a.GetDecision)
	mux.HandleFunc("GET /api/v1/decisions", a.ListDecisions)
	mux.HandleFunc("GET /api/v1/decision-loops/{decision_id}", a.GetDecisionLoop)
	mux.HandleFunc("GET /api/v1/decision-loops", a.ListDecisionLoops)
	mux.HandleFunc("POST /api/v1/decisions/{decision_id}/confirmations", a.CreateConfirmation)
	mux.HandleFunc("POST /api/v1/portfolio/init", a.InitPortfolio)
	mux.HandleFunc("GET /api/v1/portfolio/current", a.GetPortfolioCurrent)
	mux.HandleFunc("POST /api/v1/portfolio/adjustments", a.AdjustPortfolio)
	mux.HandleFunc("POST /api/v1/portfolio/holdings", a.EditHolding)
	mux.HandleFunc("POST /api/v1/portfolio/holdings/remove", a.RemoveHolding)
	mux.HandleFunc("POST /api/v1/portfolio/offline-transactions", a.RecordOfflineTransaction)
	mux.HandleFunc("POST /api/v1/portfolio/imports/validate", a.ValidatePortfolioImport)
	mux.HandleFunc("POST /api/v1/portfolio/imports/confirm", a.ConfirmPortfolioImport)
	mux.HandleFunc("POST /api/v1/portfolio/corrections", a.CorrectPortfolioFact)
	mux.HandleFunc("POST /api/v1/portfolio/rebalance-review", a.ReviewQuarterlyRebalance)
	mux.HandleFunc("POST /api/v1/evidence/refresh", a.RefreshEvidence)
	mux.HandleFunc("GET /api/v1/evidence", a.ListEvidence)
	mux.HandleFunc("GET /api/v1/evidence/verification", a.GetEvidenceVerification)
	mux.HandleFunc("POST /api/v1/evidence/rebuild-index", a.RebuildEvidenceIndex)
	mux.HandleFunc("POST /api/v1/local-knowledge/imports/validate", a.ValidateLocalKnowledgeImport)
	mux.HandleFunc("POST /api/v1/local-knowledge/imports/confirm", a.ConfirmLocalKnowledgeImport)
	mux.HandleFunc("POST /api/v1/market/refresh", a.RefreshMarket)
	mux.HandleFunc("GET /api/v1/market/source-health", a.GetMarketSourceHealth)
	mux.HandleFunc("GET /api/v1/market/snapshots/latest", a.GetLatestMarketSnapshot)
	mux.HandleFunc("GET /api/v1/data-source-quality/regression", a.GetDataSourceQualityRegression)
	mux.HandleFunc("GET /api/v1/data-source-quality/gate-resolution", a.GetDataQualityGateResolution)
	mux.HandleFunc("GET /api/v1/data-source-quality/resolutions", a.ListDataQualityGateResolutions)
	mux.HandleFunc("POST /api/v1/data-source-quality/resolutions", a.CreateDataQualityGateResolution)
	mux.HandleFunc("POST /api/v1/data-source-quality/resolutions/{resolution_id}/retire", a.RetireDataQualityGateResolution)
	mux.HandleFunc("GET /api/v1/knowledge-readiness", a.GetKnowledgeReadiness)
	mux.HandleFunc("GET /api/v1/rules/current", a.GetCurrentRule)
	mux.HandleFunc("GET /api/v1/rule-proposals", a.ListRuleProposals)
	mux.HandleFunc("POST /api/v1/rule-proposals/sop-addendum", a.CreateSOPAddendumProposal)
	mux.HandleFunc("GET /api/v1/rule-proposals/{proposal_id}/effect-validation", a.GetRuleEffectValidation)
	mux.HandleFunc("POST /api/v1/rule-proposals/{proposal_id}/effect-validation", a.RefreshRuleEffectValidation)
	mux.HandleFunc("GET /api/v1/rule-effect-tracking", a.ListRuleEffectTracking)
	mux.HandleFunc("POST /api/v1/rule-proposals/{proposal_id}/confirm", a.ConfirmRuleProposal)
	mux.HandleFunc("POST /api/v1/rule-proposals/{proposal_id}/final-confirm", a.FinalConfirmRuleProposal)
	mux.HandleFunc("GET /api/v1/settings/system", a.GetSystemSettings)
	mux.HandleFunc("PUT /api/v1/settings", a.UpdateSystemSettings)
	mux.HandleFunc("GET /api/v1/settings/capability", a.GetCapabilitySettings)
	mux.HandleFunc("PUT /api/v1/settings/capability", a.UpdateCapabilitySettings)
	mux.HandleFunc("GET /api/v1/audit-events", a.ListAuditEvents)
	mux.HandleFunc("GET /api/v1/notifications", a.ListNotifications)
	mux.HandleFunc("POST /api/v1/notifications/{notification_id}/read", a.MarkNotificationRead)
	mux.HandleFunc("POST /api/v1/notifications/read-all", a.MarkAllNotificationsRead)
	mux.HandleFunc("GET /api/v1/risk-alerts", a.ListRiskAlerts)
	mux.HandleFunc("GET /api/v1/risk-alerts/{alert_id}", a.GetRiskAlert)
	mux.HandleFunc("POST /api/v1/risk-alerts/{alert_id}/lifecycle", a.UpdateRiskAlertLifecycle)
	mux.HandleFunc("GET /api/v1/daily-auto-run/status", a.GetDailyAutoRunStatus)
	mux.HandleFunc("GET /api/v1/daily-discipline/reports/today", a.GetTodayDailyDisciplineReport)
	mux.HandleFunc("GET /api/v1/daily-discipline/reports/{report_id}", a.GetDailyDisciplineReport)
	mux.HandleFunc("GET /api/v1/daily-discipline/reports", a.ListDailyDisciplineReports)
	mux.HandleFunc("GET /api/v1/review/summary", a.GetReviewSummary)
}
