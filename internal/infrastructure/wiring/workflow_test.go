package wiring

import (
	"testing"

	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
)

func TestNewWorkflowDependenciesUsesDeepSeekWhenKeyConfigured(t *testing.T) {
	cfg := &config.Config{DeepSeek: config.DeepSeekConfig{APIKey: "test-key", BaseURL: "https://example.invalid"}, DataSources: config.DataSourceConfig{UseStub: true}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	if _, ok := deps.AnalystService.(workflow.StaticAnalystService); ok {
		t.Fatal("expected non-static DeepSeek analyst service when key is configured")
	}
	if deps.RetrievalService == nil {
		t.Fatal("expected retrieval service wiring")
	}
}

func TestNewWorkflowDependenciesUsesVecLitePathWhenConfigured(t *testing.T) {
	cfg := &config.Config{VecLite: config.VecLiteConfig{Path: "/tmp/investment-agent.veclite"}, DataSources: config.DataSourceConfig{UseStub: true}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	adapter, ok := deps.RetrievalService.(interface{ VectorIndexPath() string })
	if !ok {
		t.Fatalf("expected retrieval adapter to expose vector index path, got %T", deps.RetrievalService)
	}
	if adapter.VectorIndexPath() != "/tmp/investment-agent.veclite" {
		t.Fatalf("expected configured veclite path, got %q", adapter.VectorIndexPath())
	}
}

func TestNewWorkflowDependenciesKeepsStubWithoutDeepSeekKey(t *testing.T) {
	cfg := &config.Config{DataSources: config.DataSourceConfig{UseStub: true}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	if _, ok := deps.AnalystService.(workflow.StaticAnalystService); !ok {
		t.Fatalf("expected static analyst fallback, got %T", deps.AnalystService)
	}
	if _, ok := deps.MarketDataSource.(workflow.StubMarketDataSource); !ok {
		t.Fatalf("expected stub market source, got %T", deps.MarketDataSource)
	}
}

func TestNewWorkflowDependenciesUsesConfiguredRealSourcePlaceholders(t *testing.T) {
	cfg := &config.Config{DataSources: config.DataSourceConfig{Enabled: []string{"official", "manual"}, UseStub: false}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	market, ok := deps.MarketDataSource.(workflow.ConfiguredMarketDataSource)
	if !ok {
		t.Fatalf("expected configured market source, got %T", deps.MarketDataSource)
	}
	if len(market.Enabled) != 2 || market.Enabled[0] != "official" {
		t.Fatalf("expected configured market sources preserved: %+v", market.Enabled)
	}
	intel, ok := deps.IntelligenceSource.(workflow.ConfiguredIntelligenceSource)
	if !ok {
		t.Fatalf("expected configured intelligence source, got %T", deps.IntelligenceSource)
	}
	if len(intel.Enabled) != 2 || intel.Enabled[1] != "manual" {
		t.Fatalf("expected configured intelligence sources preserved: %+v", intel.Enabled)
	}
}

func TestNewWorkflowDependenciesDoesNotFallbackToStubForRealMarketEndpoint(t *testing.T) {
	cfg := &config.Config{DataSources: config.DataSourceConfig{Enabled: []string{"official"}, UseStub: false, MarketEndpoint: "https://example.invalid/market"}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	market, ok := deps.MarketDataSource.(workflow.ConfiguredMarketDataSource)
	if !ok {
		t.Fatalf("expected configured market source, got %T", deps.MarketDataSource)
	}
	if market.Fallback != nil {
		t.Fatalf("real market endpoint must not use stub fallback, got %T", market.Fallback)
	}
}

func TestNewWorkflowDependenciesUsesMarketCollectorsWhenEnabled(t *testing.T) {
	cfg := &config.Config{DataSources: config.DataSourceConfig{Enabled: []string{"official"}, UseStub: false, MarketEndpoint: "https://example.invalid/market", IntelligenceEndpoint: "https://example.invalid/news", MarketCollectors: config.MarketCollectorSourceConfig{Enabled: true, Sources: []string{"csindex", "eastmoney_fund"}, CSIndexBaseURL: "https://www.csindex.com.cn", EastmoneyFundBaseURL: "https://fund.eastmoney.com"}}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	collector, ok := deps.MarketDataSource.(workflow.CompositeMarketDataCollector)
	if !ok {
		t.Fatalf("expected composite market collector, got %T", deps.MarketDataSource)
	}
	if len(collector.Collectors) != 3 {
		t.Fatalf("expected two P27 collectors and configured endpoint, got %+v", collector.Collectors)
	}
	if _, ok := collector.Collectors[0].(workflow.CsindexCollector); !ok {
		t.Fatalf("expected csindex collector first, got %T", collector.Collectors[0])
	}
	if _, ok := collector.Collectors[1].(workflow.EastmoneyFundCollector); !ok {
		t.Fatalf("expected eastmoney collector second, got %T", collector.Collectors[1])
	}
	for _, source := range collector.Collectors {
		if _, ok := source.(workflow.StubMarketDataSource); ok {
			t.Fatal("real market collector configuration must not append stub fallback")
		}
	}
}

func TestNewWorkflowDependenciesUsesP89StructuredPublicCollectorWhenEnabled(t *testing.T) {
	cfg := &config.Config{DataSources: config.DataSourceConfig{Enabled: []string{"official"}, UseStub: false, MarketCollectors: config.MarketCollectorSourceConfig{Enabled: true, Sources: []string{"p89_structured_public"}}}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	collector, ok := deps.MarketDataSource.(workflow.CompositeMarketDataCollector)
	if !ok {
		t.Fatalf("expected composite market collector, got %T", deps.MarketDataSource)
	}
	if len(collector.Collectors) != 1 {
		t.Fatalf("expected only P89 structured collector, got %+v", collector.Collectors)
	}
	if _, ok := collector.Collectors[0].(workflow.P89StructuredPublicCollector); !ok {
		t.Fatalf("expected P89 structured collector, got %T", collector.Collectors[0])
	}
}

func TestNewWorkflowDependenciesAllowsStubFallbackOnlyWhenUseStubEnabled(t *testing.T) {
	cfg := &config.Config{DataSources: config.DataSourceConfig{UseStub: true, MarketCollectors: config.MarketCollectorSourceConfig{Enabled: true, Sources: []string{"eastmoney_fund"}, EastmoneyFundBaseURL: "https://fund.eastmoney.com"}}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	if _, ok := deps.MarketDataSource.(workflow.StubMarketDataSource); !ok {
		t.Fatalf("expected explicit stub mode to use stub source, got %T", deps.MarketDataSource)
	}
}

func TestConfiguredRealSourcesDoNotExposeTradingCapabilities(t *testing.T) {
	cfg := &config.Config{DataSources: config.DataSourceConfig{Enabled: []string{"official", "exchange"}, UseStub: false}}

	deps := NewWorkflowDependencies(cfg, repository.Repositories{}, nil)
	if _, ok := any(deps.MarketDataSource).(interface{ PlaceOrder() error }); ok {
		t.Fatal("market data source must not expose trading capability")
	}
	if _, ok := any(deps.IntelligenceSource).(interface{ PlaceOrder() error }); ok {
		t.Fatal("intelligence source must not expose trading capability")
	}
}
