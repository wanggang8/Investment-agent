package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadExampleConfig(t *testing.T) {
	root := findModuleRoot(t)
	path := filepath.Join(root, "configs", "config.example.yaml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("port = %d, want 8080", cfg.Server.Port)
	}
	if cfg.SQLite.Path == "" {
		t.Error("sqlite path empty")
	}
	if !cfg.DataSources.UseStub || len(cfg.DataSources.Enabled) == 0 {
		t.Fatalf("data source stub config missing: %+v", cfg.DataSources)
	}
	if cfg.DataSources.PublicEvidence.CNInfoOrgIDs["510300"] != "9900000091" {
		t.Fatalf("cninfo orgId mapping missing: %+v", cfg.DataSources.PublicEvidence.CNInfoOrgIDs)
	}
	if cfg.DeepSeek.Model == "" {
		t.Fatal("deepseek model must have a default/example value")
	}
	if cfg.DeepSeek.TimeoutSeconds != 60 {
		t.Fatalf("deepseek timeout_seconds = %d, want 60", cfg.DeepSeek.TimeoutSeconds)
	}
}

func TestLoadDefaultsToLocalConfigYAMLWhenPresent(t *testing.T) {
	dir := t.TempDir()
	configsDir := filepath.Join(dir, "configs")
	if err := os.MkdirAll(configsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfigFile(t, filepath.Join(configsDir, "config.yaml"), 18080, "local-model")
	writeConfigFile(t, filepath.Join(configsDir, "config.example.yaml"), 28080, "example-model")
	withWorkingDir(t, dir)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Server.Port != 18080 {
		t.Fatalf("port=%d, want local config port", cfg.Server.Port)
	}
	if cfg.DeepSeek.Model != "local-model" {
		t.Fatalf("model=%q, want local-model", cfg.DeepSeek.Model)
	}
}

func TestLoadFallsBackToExampleConfigWhenLocalConfigMissing(t *testing.T) {
	dir := t.TempDir()
	configsDir := filepath.Join(dir, "configs")
	if err := os.MkdirAll(configsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfigFile(t, filepath.Join(configsDir, "config.example.yaml"), 28080, "example-model")
	withWorkingDir(t, dir)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Server.Port != 28080 {
		t.Fatalf("port=%d, want example config port", cfg.Server.Port)
	}
	if cfg.DeepSeek.Model != "example-model" {
		t.Fatalf("model=%q, want example-model", cfg.DeepSeek.Model)
	}
}

func TestLoadEnvConfigOverridesDefaultLocalConfig(t *testing.T) {
	dir := t.TempDir()
	configsDir := filepath.Join(dir, "configs")
	if err := os.MkdirAll(configsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeConfigFile(t, filepath.Join(configsDir, "config.yaml"), 18080, "local-model")
	envConfigPath := filepath.Join(dir, "override.yaml")
	writeConfigFile(t, envConfigPath, 38080, "env-model")
	t.Setenv("INVESTMENT_AGENT_CONFIG", envConfigPath)
	withWorkingDir(t, dir)

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Server.Port != 38080 {
		t.Fatalf("port=%d, want env config port", cfg.Server.Port)
	}
	if cfg.DeepSeek.Model != "env-model" {
		t.Fatalf("model=%q, want env-model", cfg.DeepSeek.Model)
	}
}

func TestLoadDeepSeekModelAndTimeoutFromConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `server:
  host: 127.0.0.1
  port: 8080
sqlite:
  path: ./data/agent.db
veclite:
  path: ./data/veclite
deepseek:
  api_key: test-key
  base_url: http://example.invalid
  model: gpt-5.4-mini
  timeout_seconds: 7
data_sources:
  enabled:
    - stub
  use_stub: true
log:
  level: info
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.DeepSeek.Model != "gpt-5.4-mini" {
		t.Fatalf("model = %q, want gpt-5.4-mini", cfg.DeepSeek.Model)
	}
	if cfg.DeepSeek.TimeoutSeconds != 7 {
		t.Fatalf("timeout_seconds = %d, want 7", cfg.DeepSeek.TimeoutSeconds)
	}
}

func writeConfigFile(t *testing.T, path string, port int, model string) {
	t.Helper()
	content := fmt.Sprintf(`server:
  host: 127.0.0.1
  port: %d
sqlite:
  path: ./data/agent.db
veclite:
  path: ./data/veclite
deepseek:
  api_key: ""
  base_url: https://api.deepseek.com
  model: %s
  timeout_seconds: 15
data_sources:
  enabled:
    - stub
  use_stub: true
log:
  level: info
`, port, model)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("restore working dir: %v", err)
		}
	})
}

func TestLoadDeepSeekDeploymentEnvOverrides(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "test-deployment-key")
	t.Setenv("DEEPSEEK_BASE_URL", "https://llm.example.invalid")
	t.Setenv("DEEPSEEK_MODEL", "deployment-model")
	t.Setenv("DEEPSEEK_TIMEOUT_SECONDS", "42")
	root := findModuleRoot(t)
	path := filepath.Join(root, "configs", "config.example.yaml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.DeepSeek.APIKey != "test-deployment-key" {
		t.Fatalf("api key override missing")
	}
	if cfg.DeepSeek.BaseURL != "https://llm.example.invalid" {
		t.Fatalf("base url = %q, want env override", cfg.DeepSeek.BaseURL)
	}
	if cfg.DeepSeek.Model != "deployment-model" {
		t.Fatalf("model = %q, want env override", cfg.DeepSeek.Model)
	}
	if cfg.DeepSeek.TimeoutSeconds != 42 {
		t.Fatalf("timeout = %d, want env override", cfg.DeepSeek.TimeoutSeconds)
	}
}

func TestLoadDeepSeekAPIKeyFromFileWhenEnvKeyMissing(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "deepseek_api_key")
	if err := os.WriteFile(keyPath, []byte(" file-secret-key \n"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DEEPSEEK_API_KEY_FILE", keyPath)
	root := findModuleRoot(t)
	path := filepath.Join(root, "configs", "config.example.yaml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.DeepSeek.APIKey != "file-secret-key" {
		t.Fatalf("api key = %q, want file-secret-key", cfg.DeepSeek.APIKey)
	}
}

func TestLoadDeepSeekAPIKeyEnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	keyPath := filepath.Join(dir, "deepseek_api_key")
	if err := os.WriteFile(keyPath, []byte("file-secret-key"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DEEPSEEK_API_KEY", "env-secret-key")
	t.Setenv("DEEPSEEK_API_KEY_FILE", keyPath)
	root := findModuleRoot(t)
	path := filepath.Join(root, "configs", "config.example.yaml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.DeepSeek.APIKey != "env-secret-key" {
		t.Fatalf("api key = %q, want env-secret-key", cfg.DeepSeek.APIKey)
	}
}

func TestLoadDataSourceEnvOverrides(t *testing.T) {
	t.Setenv("INVESTMENT_AGENT_DATA_SOURCES", "official,manual")
	t.Setenv("INVESTMENT_AGENT_MARKET_DATA_ENDPOINT", "https://example.invalid/market")
	t.Setenv("INVESTMENT_AGENT_INTELLIGENCE_ENDPOINT", "https://example.invalid/news")
	t.Setenv("INVESTMENT_AGENT_USE_STUB_DATA", "false")
	root := findModuleRoot(t)
	path := filepath.Join(root, "configs", "config.example.yaml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DataSources.UseStub {
		t.Fatal("expected stub data disabled by env")
	}
	if len(cfg.DataSources.Enabled) != 2 || cfg.DataSources.Enabled[0] != "official" || cfg.DataSources.Enabled[1] != "manual" {
		t.Fatalf("unexpected data sources: %+v", cfg.DataSources.Enabled)
	}
	if cfg.DataSources.MarketEndpoint != "https://example.invalid/market" || cfg.DataSources.IntelligenceEndpoint != "https://example.invalid/news" {
		t.Fatalf("unexpected endpoints: %+v", cfg.DataSources)
	}
}

func TestValidateRejectsInvalidLocalRuntimeConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Port: -1}, SQLite: SQLiteConfig{}, VecLite: VecLiteConfig{}, DataSources: DataSourceConfig{Enabled: []string{"public-http"}, UseStub: false}, Log: LogConfig{Level: "verbose"}}

	err := cfg.Validate()

	if err == nil {
		t.Fatal("expected invalid config to fail validation")
	}
	for _, want := range []string{"server.port", "sqlite.path", "veclite.path", "market_endpoint", "log.level"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected validation error to contain %q, got %v", want, err)
		}
	}
}

func TestValidateAcceptsStubLocalRuntimeConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true}, Log: LogConfig{Level: "info"}}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateRejectsReleaseRuntimeWithStubData(t *testing.T) {
	cfg := Config{Runtime: RuntimeConfig{Mode: "release"}, Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true}, Log: LogConfig{Level: "info"}}

	err := cfg.Validate()

	if err == nil {
		t.Fatal("expected release runtime with stub data to fail validation")
	}
	if !strings.Contains(err.Error(), "runtime.mode=release") || !strings.Contains(err.Error(), "data_sources.use_stub") {
		t.Fatalf("expected validation error to mention release runtime and stub data, got %v", err)
	}
}

func TestValidateAcceptsReleaseRuntimeWithStructuredPublicCollector(t *testing.T) {
	cfg := Config{
		Runtime: RuntimeConfig{Mode: "release"},
		Server:  ServerConfig{Host: "127.0.0.1", Port: 8080},
		SQLite:  SQLiteConfig{Path: "./data/agent.db"},
		VecLite: VecLiteConfig{Path: "./data/veclite"},
		DataSources: DataSourceConfig{
			UseStub: false,
			MarketCollectors: MarketCollectorSourceConfig{
				Enabled: true,
				Sources: []string{"p89_structured_public"},
			},
		},
		Log: LogConfig{Level: "info"},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateRejectsInvalidPublicEvidenceConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true, PublicEvidence: PublicEvidenceSourceConfig{Enabled: true, Sources: []string{"cninfo", "cnifno"}, CNInfoBaseURL: "not a url"}}, Log: LogConfig{Level: "info"}}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected invalid public evidence config")
	}
	for _, want := range []string{"public_evidence.sources", "cninfo_base_url"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected validation error to contain %q, got %v", want, err)
		}
	}
}

func TestValidateAcceptsPublicEvidenceConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true, PublicEvidence: PublicEvidenceSourceConfig{Enabled: true, Sources: []string{"cninfo", "szse", "csrc"}, CNInfoBaseURL: "https://www.cninfo.com.cn", SZSEBaseURL: "https://www.szse.cn", CSRCBaseURL: "https://www.csrc.gov.cn"}}, Log: LogConfig{Level: "info"}}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateRejectsInvalidMarketCollectorConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true, MarketCollectors: MarketCollectorSourceConfig{Enabled: true, Sources: []string{"csindex", "sina"}, CSIndexBaseURL: "not a url"}}, Log: LogConfig{Level: "info"}}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected invalid market collector config")
	}
	for _, want := range []string{"market_collectors.sources", "csindex_base_url"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected validation error to contain %q, got %v", want, err)
		}
	}
}

func TestValidateAcceptsMarketCollectorConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true, MarketCollectors: MarketCollectorSourceConfig{Enabled: true, Sources: []string{"csindex", "eastmoney_fund"}, CSIndexBaseURL: "https://www.csindex.com.cn", EastmoneyFundBaseURL: "https://fund.eastmoney.com"}}, Log: LogConfig{Level: "info"}}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateAcceptsMarketCollectorsWithoutGenericEndpoints(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"market_collectors"}, UseStub: false, MarketCollectors: MarketCollectorSourceConfig{Enabled: true, Sources: []string{"csindex", "eastmoney_fund"}, CSIndexBaseURL: "https://www.csindex.com.cn", EastmoneyFundBaseURL: "https://fund.eastmoney.com"}}, Log: LogConfig{Level: "info"}}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateAcceptsSinglePublicEvidenceSourceConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true, PublicEvidence: PublicEvidenceSourceConfig{Enabled: true, Sources: []string{"cninfo"}, CNInfoBaseURL: "https://www.cninfo.com.cn"}}, Log: LogConfig{Level: "info"}}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateAcceptsEastmoneyOnlyMarketCollectorConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"market_collectors"}, UseStub: false, MarketCollectors: MarketCollectorSourceConfig{Enabled: true, Sources: []string{"eastmoney_fund"}, EastmoneyFundBaseURL: "https://fund.eastmoney.com"}}, Log: LogConfig{Level: "info"}}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateAcceptsCSIndexOnlyMarketCollectorConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"market_collectors"}, UseStub: false, MarketCollectors: MarketCollectorSourceConfig{Enabled: true, Sources: []string{"csindex"}, CSIndexBaseURL: "https://www.csindex.com.cn"}}, Log: LogConfig{Level: "info"}}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestLoadExampleConfigKeepsDailyAutoRunDisabled(t *testing.T) {
	root := findModuleRoot(t)
	path := filepath.Join(root, "configs", "config.example.yaml")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.DailyAutoRun.Enabled {
		t.Fatal("daily auto-run must be disabled by default")
	}
	if cfg.DailyAutoRun.RunTime != "08:30" {
		t.Fatalf("run time = %q, want 08:30", cfg.DailyAutoRun.RunTime)
	}
	if cfg.DailyAutoRun.Scope != "holdings" {
		t.Fatalf("scope = %q, want holdings", cfg.DailyAutoRun.Scope)
	}
	if cfg.DailyAutoRun.MaxSymbols != 20 {
		t.Fatalf("max symbols = %d, want 20", cfg.DailyAutoRun.MaxSymbols)
	}
}

func TestValidateRejectsInvalidDailyAutoRunConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true}, DailyAutoRun: DailyAutoRunConfig{Enabled: true, RunTime: "25:99", Timezone: "", Scope: "watchlist", Retry: -1, TimeoutSeconds: 0, MaxSymbols: 0}, Log: LogConfig{Level: "info"}}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected invalid daily auto-run config")
	}
	for _, want := range []string{"daily_auto_run.run_time", "daily_auto_run.timezone", "daily_auto_run.scope", "daily_auto_run.retry", "daily_auto_run.timeout_seconds", "daily_auto_run.max_symbols"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("expected validation error to contain %q, got %v", want, err)
		}
	}
}

func TestValidateAcceptsDailyAutoRunConfig(t *testing.T) {
	cfg := Config{Server: ServerConfig{Host: "127.0.0.1", Port: 8080}, SQLite: SQLiteConfig{Path: "./data/agent.db"}, VecLite: VecLiteConfig{Path: "./data/veclite"}, DataSources: DataSourceConfig{Enabled: []string{"stub"}, UseStub: true}, DailyAutoRun: DailyAutoRunConfig{Enabled: true, RunTime: "08:30", Timezone: "Asia/Shanghai", Scope: "holdings", Retry: 1, TimeoutSeconds: 900, MaxSymbols: 20}, Log: LogConfig{Level: "info"}}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func findModuleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}
