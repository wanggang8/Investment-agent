package wiring

import (
	"net/http"
	"strings"
	"time"

	"investment-agent/internal/application/service"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	"investment-agent/internal/infrastructure/llm/deepseek"
)

// NewWorkflowDependencies 按配置组装生产工作流依赖；默认本地 stub 不依赖公网。
func NewWorkflowDependencies(cfg *config.Config, repos repository.Repositories, tx repository.Transactor) workflow.WorkflowDependencies {
	deps := workflow.NewWorkflowDependencies(repos, tx)
	if cfg == nil {
		return deps
	}
	deps.DailyAutoRunConfig = cfg.DailyAutoRun
	if cfg.DataSources.UseStub {
		deps.MarketDataSource = workflow.StubMarketDataSource{}
		deps.IntelligenceSource = workflow.StubIntelligenceSource{}
	} else {
		client := &http.Client{Timeout: 10 * time.Second}
		marketSource := workflow.ConfiguredMarketDataSource{Enabled: cfg.DataSources.Enabled, MarketEndpoint: cfg.DataSources.MarketEndpoint, HTTPClient: client}
		if cfg.DataSources.MarketCollectors.Enabled {
			collectors := []workflow.MarketDataSource{}
			for _, source := range cfg.DataSources.MarketCollectors.Sources {
				switch strings.TrimSpace(source) {
				case "csindex":
					collectors = append(collectors, workflow.CsindexCollector{BaseURL: cfg.DataSources.MarketCollectors.CSIndexBaseURL, HTTPClient: client, IncludeExtended: true})
				case "eastmoney_fund":
					collectors = append(collectors, workflow.EastmoneyFundCollector{BaseURL: cfg.DataSources.MarketCollectors.EastmoneyFundBaseURL, HTTPClient: client, IncludeExtended: true})
				case "p89_structured_public":
					collectors = append(collectors, workflow.P89StructuredPublicCollector{HTTPClient: client})
				}
			}
			if cfg.DataSources.MarketEndpoint != "" {
				collectors = append(collectors, workflow.ConfiguredMarketDataSource{Enabled: cfg.DataSources.Enabled, MarketEndpoint: cfg.DataSources.MarketEndpoint, HTTPClient: client})
			}
			deps.MarketDataSource = workflow.CompositeMarketDataCollector{Collectors: collectors}
		} else {
			deps.MarketDataSource = marketSource
		}
		deps.IntelligenceSource = workflow.ConfiguredIntelligenceSource{Enabled: cfg.DataSources.Enabled, IntelligenceEndpoint: cfg.DataSources.IntelligenceEndpoint, HTTPClient: client}
	}
	if strings.TrimSpace(cfg.DeepSeek.APIKey) != "" {
		timeout := 15 * time.Second
		if cfg.DeepSeek.TimeoutSeconds > 0 {
			timeout = time.Duration(cfg.DeepSeek.TimeoutSeconds) * time.Second
		}
		deps.AnalystService = deepseek.NewClient(deepseek.Config{APIKey: cfg.DeepSeek.APIKey, BaseURL: cfg.DeepSeek.BaseURL, Model: cfg.DeepSeek.Model, TimeoutSeconds: cfg.DeepSeek.TimeoutSeconds}, &http.Client{Timeout: timeout})
	}
	index := service.NewFileVectorIndex(cfg.VecLite.Path)
	deps.RetrievalService = service.NewRetrievalAdapter(tx, index)
	deps.VectorIndexWriter = index
	return deps
}
