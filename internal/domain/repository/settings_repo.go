package repository

import "context"

// SystemSettings stores local non-rule preferences.
type SystemSettings struct {
	SettingsID             string
	NotificationConfigJSON string
	DataSourcesJSON        string
	UpdatedAt              string
}

// CapabilityConfig stores capability scope configuration.
type CapabilityConfig struct {
	CapabilityID        string
	AssetTypesJSON      string
	SymbolsJSON         string
	ExcludedSymbolsJSON string
	StrategyScopeJSON   string
	UpdatedAt           string
}

// SettingsRepository defines persistence for ordinary settings and capability scope.
type SettingsRepository interface {
	SaveSystemSettings(ctx context.Context, settings SystemSettings) error
	GetLatestSystemSettings(ctx context.Context) (SystemSettings, error)
	SaveCapabilityConfig(ctx context.Context, config CapabilityConfig) error
	GetLatestCapabilityConfig(ctx context.Context) (CapabilityConfig, error)
}
