package repository

import "context"

// Repositories groups repository interfaces available inside a transaction.
type Repositories struct {
	PortfolioRepo                 PortfolioRepository
	DecisionRepo                  DecisionRepository
	AuditRepo                     AuditRepository
	RuleRepo                      RuleRepository
	MarketRepo                    MarketRepository
	SettingsRepo                  SettingsRepository
	IntelligenceRepo              IntelligenceRepository
	NotificationRepo              NotificationRepository
	DailyAutoRunRepo              DailyAutoRunRepository
	DailyDisciplineReportRepo     DailyDisciplineReportRepository
	RiskAlertRepo                 RiskAlertRepository
	RuleEffectRepo                RuleEffectRepository
	DataQualityGateResolutionRepo DataQualityGateResolutionRepository
}

// Transactor coordinates multi-repository writes as one atomic unit.
type Transactor interface {
	WithinTx(ctx context.Context, fn func(context.Context, Repositories) error) error
}
