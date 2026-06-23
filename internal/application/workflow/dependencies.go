package workflow

import (
	"context"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
)

// VectorIndexWriter 定义证据刷新时写入本地 VecLite 索引的边界。
type VectorIndexWriter interface {
	Upsert(ctx context.Context, chunk repository.RAGChunk) error
}

// WorkflowDependencies 汇总工作流写事实表需要的仓储。
// P3 的生产路径必须传入这些依赖，避免只停留在内存审计。
type WorkflowDependencies struct {
	DecisionRepo              repository.DecisionRepository
	AuditRepo                 repository.AuditRepository
	RuleRepo                  repository.RuleRepository
	MarketRepo                repository.MarketRepository
	SettingsRepo              repository.SettingsRepository
	IntelligenceRepo          repository.IntelligenceRepository
	NotificationRepo          repository.NotificationRepository
	DailyAutoRunRepo          repository.DailyAutoRunRepository
	DailyDisciplineReportRepo repository.DailyDisciplineReportRepository
	RiskAlertRepo             repository.RiskAlertRepository
	DailyAutoRunConfig        config.DailyAutoRunConfig
	PortfolioRepo             repository.PortfolioRepository
	Transactor                repository.Transactor
	MarketDataSource          MarketDataSource
	IntelligenceSource        IntelligenceSource
	AnalystService            AnalystService
	RetrievalService          RetrievalService
	VectorIndexWriter         VectorIndexWriter
}

// NewWorkflowDependencies 基于仓储接口创建工作流依赖。
func NewWorkflowDependencies(repos repository.Repositories, tx ...repository.Transactor) WorkflowDependencies {
	var transactor repository.Transactor
	if len(tx) > 0 {
		transactor = tx[0]
	}
	return WorkflowDependencies{
		DecisionRepo:              repos.DecisionRepo,
		AuditRepo:                 repos.AuditRepo,
		RuleRepo:                  repos.RuleRepo,
		MarketRepo:                repos.MarketRepo,
		SettingsRepo:              repos.SettingsRepo,
		IntelligenceRepo:          repos.IntelligenceRepo,
		NotificationRepo:          repos.NotificationRepo,
		DailyAutoRunRepo:          repos.DailyAutoRunRepo,
		DailyDisciplineReportRepo: repos.DailyDisciplineReportRepo,
		RiskAlertRepo:             repos.RiskAlertRepo,
		PortfolioRepo:             repos.PortfolioRepo,
		Transactor:                transactor,
		MarketDataSource:          StubMarketDataSource{},
		IntelligenceSource:        StubIntelligenceSource{},
		AnalystService:            StaticAnalystService{},
	}
}

func (d WorkflowDependencies) repositories() repository.Repositories {
	return repository.Repositories{
		DecisionRepo:              d.DecisionRepo,
		AuditRepo:                 d.AuditRepo,
		RuleRepo:                  d.RuleRepo,
		MarketRepo:                d.MarketRepo,
		SettingsRepo:              d.SettingsRepo,
		IntelligenceRepo:          d.IntelligenceRepo,
		NotificationRepo:          d.NotificationRepo,
		DailyAutoRunRepo:          d.DailyAutoRunRepo,
		DailyDisciplineReportRepo: d.DailyDisciplineReportRepo,
		RiskAlertRepo:             d.RiskAlertRepo,
		PortfolioRepo:             d.PortfolioRepo,
	}
}

func (d WorkflowDependencies) marketDataSource() MarketDataSource {
	if d.MarketDataSource == nil {
		return StubMarketDataSource{}
	}
	return d.MarketDataSource
}

func (d WorkflowDependencies) intelligenceSource() IntelligenceSource {
	if d.IntelligenceSource == nil {
		return StubIntelligenceSource{}
	}
	return d.IntelligenceSource
}

func (d WorkflowDependencies) analystService() AnalystService {
	if d.AnalystService == nil {
		return StaticAnalystService{}
	}
	return d.AnalystService
}

func (d WorkflowDependencies) retrievalService() RetrievalService {
	return d.RetrievalService
}
