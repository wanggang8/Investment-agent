package dto

// SystemSettingsDTO 表示允许通过普通设置接口保存的本地偏好。
type SystemSettingsDTO struct {
	NotificationEnabled bool     `json:"notification_enabled"`
	PagePreference      string   `json:"page_preference,omitempty"`
	DataSources         []string `json:"data_sources"`
}

// CapabilitySettingsDTO 表示能力圈设置，必须与普通系统设置分离。
type CapabilitySettingsDTO struct {
	CapabilityID      string   `json:"capability_id,omitempty"`
	AssetTypes        []string `json:"asset_types,omitempty"`
	Symbols           []string `json:"symbols,omitempty"`
	ExcludedSymbols   []string `json:"excluded_symbols,omitempty"`
	StrategyScope     []string `json:"strategy_scope,omitempty"`
	UpdatedAt         string   `json:"updated_at,omitempty"`
	AllowedAssetTypes []string `json:"allowed_asset_types,omitempty"`
	AllowedSymbols    []string `json:"allowed_symbols,omitempty"`
	Notes             string   `json:"notes,omitempty"`
}

type SystemStatusDTO struct {
	SQLiteStatus   string   `json:"sqlite_status"`
	SQLitePath     string   `json:"sqlite_path,omitempty"`
	VecLiteStatus  string   `json:"veclite_status"`
	VecLitePath    string   `json:"veclite_path,omitempty"`
	DeepSeekStatus string   `json:"deepseek_status"`
	DataSources    []string `json:"data_sources"`
	LogLevel       string   `json:"log_level"`
}
