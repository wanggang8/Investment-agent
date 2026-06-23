package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 是后端运行配置的总入口，对应 configs/config.example.yaml。
// 当前只保存本地服务、SQLite、VecLite、DeepSeek 与日志配置。
type Config struct {
	Server       ServerConfig       `yaml:"server"`
	SQLite       SQLiteConfig       `yaml:"sqlite"`
	VecLite      VecLiteConfig      `yaml:"veclite"`
	DeepSeek     DeepSeekConfig     `yaml:"deepseek"`
	DataSources  DataSourceConfig   `yaml:"data_sources"`
	DailyAutoRun DailyAutoRunConfig `yaml:"daily_auto_run"`
	Log          LogConfig          `yaml:"log"`
}

// ServerConfig 描述本地 HTTP 服务监听地址。
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// SQLiteConfig 保存本地事实数据库文件路径。
type SQLiteConfig struct {
	Path string `yaml:"path"`
}

// VecLiteConfig 保存可重建检索索引的本地路径。
type VecLiteConfig struct {
	Path string `yaml:"path"`
}

// DeepSeekConfig 保存分析模型调用所需配置；最终裁决仍由领域规则完成。
type DeepSeekConfig struct {
	APIKey         string `yaml:"api_key"`
	BaseURL        string `yaml:"base_url"`
	Model          string `yaml:"model"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

// DataSourceConfig 控制真实数据源和本地 stub 的启用方式。
type DataSourceConfig struct {
	Enabled              []string                    `yaml:"enabled"`
	UseStub              bool                        `yaml:"use_stub"`
	MarketEndpoint       string                      `yaml:"market_endpoint"`
	IntelligenceEndpoint string                      `yaml:"intelligence_endpoint"`
	PublicEvidence       PublicEvidenceSourceConfig  `yaml:"public_evidence"`
	MarketCollectors     MarketCollectorSourceConfig `yaml:"market_collectors"`
}

// PublicEvidenceSourceConfig controls P26 public evidence collectors.
type PublicEvidenceSourceConfig struct {
	Enabled       bool              `yaml:"enabled"`
	Sources       []string          `yaml:"sources"`
	CNInfoBaseURL string            `yaml:"cninfo_base_url"`
	CNInfoOrgIDs  map[string]string `yaml:"cninfo_org_ids"`
	SZSEBaseURL   string            `yaml:"szse_base_url"`
	CSRCBaseURL   string            `yaml:"csrc_base_url"`
}

// MarketCollectorSourceConfig controls P27 read-only market data collectors.
type MarketCollectorSourceConfig struct {
	Enabled              bool     `yaml:"enabled"`
	Sources              []string `yaml:"sources"`
	CSIndexBaseURL       string   `yaml:"csindex_base_url"`
	EastmoneyFundBaseURL string   `yaml:"eastmoney_fund_base_url"`
}

// LogConfig 控制 slog 日志等级。
type LogConfig struct {
	Level string `yaml:"level"`
}

// DailyAutoRunConfig controls the local-only daily scheduler.
type DailyAutoRunConfig struct {
	Enabled        bool   `yaml:"enabled"`
	RunTime        string `yaml:"run_time"`
	Timezone       string `yaml:"timezone"`
	Scope          string `yaml:"scope"`
	Retry          int    `yaml:"retry"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	MaxSymbols     int    `yaml:"max_symbols"`
}

// Addr 返回 HTTP 服务监听地址；缺省值面向本地开发。
func (c ServerConfig) Addr() string {
	host := c.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := c.Port
	if port == 0 {
		port = 8080
	}
	return fmt.Sprintf("%s:%d", host, port)
}

// Load 读取 YAML 配置，并叠加环境变量覆盖项。
// path 为空时依次使用 INVESTMENT_AGENT_CONFIG 与示例配置文件。
func Load(path string) (*Config, error) {
	if path == "" {
		path = os.Getenv("INVESTMENT_AGENT_CONFIG")
	}
	if path == "" {
		path = "configs/config.example.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	applyEnvOverrides(&cfg)
	applyDefaults(&cfg)
	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if strings.TrimSpace(cfg.DeepSeek.BaseURL) == "" {
		cfg.DeepSeek.BaseURL = "https://api.deepseek.com"
	}
	if strings.TrimSpace(cfg.DeepSeek.Model) == "" {
		cfg.DeepSeek.Model = "deepseek-chat"
	}
	if cfg.DeepSeek.TimeoutSeconds <= 0 {
		cfg.DeepSeek.TimeoutSeconds = 15
	}
}

// applyEnvOverrides 支持本地部署时用环境变量覆盖敏感或环境相关配置。
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("INVESTMENT_AGENT_SERVER_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("INVESTMENT_AGENT_SQLITE_PATH"); v != "" {
		cfg.SQLite.Path = v
	}
	if v := os.Getenv("INVESTMENT_AGENT_VECLITE_PATH"); v != "" {
		cfg.VecLite.Path = v
	}
	if v := os.Getenv("DEEPSEEK_API_KEY"); v != "" {
		cfg.DeepSeek.APIKey = v
	}
	if v := os.Getenv("DEEPSEEK_BASE_URL"); v != "" {
		cfg.DeepSeek.BaseURL = v
	}
	if v := os.Getenv("DEEPSEEK_MODEL"); v != "" {
		cfg.DeepSeek.Model = v
	}
	if v := os.Getenv("DEEPSEEK_TIMEOUT_SECONDS"); v != "" {
		if seconds, err := strconv.Atoi(v); err == nil {
			cfg.DeepSeek.TimeoutSeconds = seconds
		}
	}
	if v := os.Getenv("INVESTMENT_AGENT_DATA_SOURCES"); v != "" {
		cfg.DataSources.Enabled = splitCSV(v)
	}
	if v := os.Getenv("INVESTMENT_AGENT_MARKET_DATA_ENDPOINT"); v != "" {
		cfg.DataSources.MarketEndpoint = v
	}
	if v := os.Getenv("INVESTMENT_AGENT_INTELLIGENCE_ENDPOINT"); v != "" {
		cfg.DataSources.IntelligenceEndpoint = v
	}
	if v := os.Getenv("INVESTMENT_AGENT_USE_STUB_DATA"); v != "" {
		cfg.DataSources.UseStub = v == "true" || v == "1"
	}
	if v := os.Getenv("INVESTMENT_AGENT_LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
}

// Validate 检查本地运行所需的关键配置，返回聚合错误便于 CLI 诊断。
func (c Config) Validate() error {
	var problems []string
	if c.Server.Port < 0 || c.Server.Port > 65535 {
		problems = append(problems, "server.port must be between 0 and 65535")
	}
	if strings.TrimSpace(c.SQLite.Path) == "" {
		problems = append(problems, "sqlite.path is required")
	}
	if strings.TrimSpace(c.VecLite.Path) == "" {
		problems = append(problems, "veclite.path is required")
	}
	if !c.DataSources.UseStub {
		if strings.TrimSpace(c.DataSources.MarketEndpoint) == "" {
			if !c.DataSources.MarketCollectors.Enabled {
				problems = append(problems, "data_sources.market_endpoint is required when use_stub is false and market_collectors is disabled")
			}
		} else if !validHTTPURL(c.DataSources.MarketEndpoint) {
			problems = append(problems, "data_sources.market_endpoint must be http or https URL")
		}
		if requiresIntelligenceEndpoint(c.DataSources.Enabled) {
			if strings.TrimSpace(c.DataSources.IntelligenceEndpoint) == "" {
				problems = append(problems, "data_sources.intelligence_endpoint is required when intelligence sources are enabled")
			} else if !validHTTPURL(c.DataSources.IntelligenceEndpoint) {
				problems = append(problems, "data_sources.intelligence_endpoint must be http or https URL")
			}
		}
	}
	if c.DataSources.PublicEvidence.Enabled {
		if len(c.DataSources.PublicEvidence.Sources) == 0 {
			problems = append(problems, "data_sources.public_evidence.sources is required when public evidence is enabled")
		}
		for _, source := range c.DataSources.PublicEvidence.Sources {
			switch strings.TrimSpace(source) {
			case "cninfo":
				if !validHTTPURL(c.DataSources.PublicEvidence.CNInfoBaseURL) {
					problems = append(problems, "data_sources.public_evidence.cninfo_base_url must be http or https URL")
				}
			case "szse":
				if !validHTTPURL(c.DataSources.PublicEvidence.SZSEBaseURL) {
					problems = append(problems, "data_sources.public_evidence.szse_base_url must be http or https URL")
				}
			case "csrc":
				if !validHTTPURL(c.DataSources.PublicEvidence.CSRCBaseURL) {
					problems = append(problems, "data_sources.public_evidence.csrc_base_url must be http or https URL")
				}
			case "csindex_index", "eastmoney_fund":
			default:
				problems = append(problems, "data_sources.public_evidence.sources must contain only cninfo, szse, csrc, csindex_index or eastmoney_fund")
			}
		}
	}
	if c.DataSources.MarketCollectors.Enabled {
		if len(c.DataSources.MarketCollectors.Sources) == 0 {
			problems = append(problems, "data_sources.market_collectors.sources is required when market collectors are enabled")
		}
		for _, source := range c.DataSources.MarketCollectors.Sources {
			switch strings.TrimSpace(source) {
			case "csindex":
				if !validHTTPURL(c.DataSources.MarketCollectors.CSIndexBaseURL) {
					problems = append(problems, "data_sources.market_collectors.csindex_base_url must be http or https URL")
				}
			case "eastmoney_fund":
				if !validHTTPURL(c.DataSources.MarketCollectors.EastmoneyFundBaseURL) {
					problems = append(problems, "data_sources.market_collectors.eastmoney_fund_base_url must be http or https URL")
				}
			case "p89_structured_public":
			default:
				problems = append(problems, "data_sources.market_collectors.sources must contain only csindex, eastmoney_fund or p89_structured_public")
			}
		}
	}
	if c.DailyAutoRun.Enabled {
		if _, err := time.Parse("15:04", c.DailyAutoRun.RunTime); err != nil {
			problems = append(problems, "daily_auto_run.run_time must use HH:MM local time")
		}
		if strings.TrimSpace(c.DailyAutoRun.Timezone) == "" {
			problems = append(problems, "daily_auto_run.timezone is required when daily auto-run is enabled")
		} else if _, err := time.LoadLocation(c.DailyAutoRun.Timezone); err != nil {
			problems = append(problems, "daily_auto_run.timezone must be a valid IANA timezone")
		}
		switch strings.TrimSpace(c.DailyAutoRun.Scope) {
		case "holdings":
		default:
			problems = append(problems, "daily_auto_run.scope must be holdings")
		}
		if c.DailyAutoRun.Retry < 0 || c.DailyAutoRun.Retry > 3 {
			problems = append(problems, "daily_auto_run.retry must be between 0 and 3")
		}
		if c.DailyAutoRun.TimeoutSeconds <= 0 || c.DailyAutoRun.TimeoutSeconds > 3600 {
			problems = append(problems, "daily_auto_run.timeout_seconds must be between 1 and 3600")
		}
		if c.DailyAutoRun.MaxSymbols <= 0 || c.DailyAutoRun.MaxSymbols > 200 {
			problems = append(problems, "daily_auto_run.max_symbols must be between 1 and 200")
		}
	}
	switch strings.TrimSpace(c.Log.Level) {
	case "", "debug", "info", "warn", "error":
	default:
		problems = append(problems, "log.level must be debug, info, warn or error")
	}
	if len(problems) > 0 {
		return fmt.Errorf("invalid config: %s", strings.Join(problems, "; "))
	}
	return nil
}

func validHTTPURL(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && parsed.Host != "" && (parsed.Scheme == "http" || parsed.Scheme == "https")
}

func requiresIntelligenceEndpoint(enabled []string) bool {
	for _, source := range enabled {
		switch strings.TrimSpace(source) {
		case "http", "public_http", "official", "exchange", "news", "intelligence":
			return true
		}
	}
	return false
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}
