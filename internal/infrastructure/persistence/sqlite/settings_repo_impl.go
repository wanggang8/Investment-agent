package sqlite

import (
	"context"
	"database/sql"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// SettingsRepository is the SQLite implementation for user and capability settings.
type SettingsRepository struct{ db dbtx }

// NewSettingsRepository creates a settings repository instance.
func NewSettingsRepository(db *sql.DB) *SettingsRepository { return &SettingsRepository{db: db} }

// SaveSystemSettings saves ordinary local preferences.
func (r *SettingsRepository) SaveSystemSettings(ctx context.Context, s repository.SystemSettings) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO user_settings (settings_id,notification_config_json,data_sources_json,updated_at) VALUES (?,?,?,?)`, s.SettingsID, nullString(s.NotificationConfigJSON), nullString(s.DataSourcesJSON), s.UpdatedAt)
	return apperr.FromRepositoryError(err)
}

func (r *SettingsRepository) GetLatestSystemSettings(ctx context.Context) (repository.SystemSettings, error) {
	var s repository.SystemSettings
	err := r.db.QueryRowContext(ctx, `SELECT settings_id,COALESCE(notification_config_json,'{}'),COALESCE(data_sources_json,'[]'),updated_at FROM user_settings ORDER BY updated_at DESC LIMIT 1`).Scan(&s.SettingsID, &s.NotificationConfigJSON, &s.DataSourcesJSON, &s.UpdatedAt)
	return s, apperr.FromRepositoryError(err)
}

// SaveCapabilityConfig saves capability scope settings.
func (r *SettingsRepository) SaveCapabilityConfig(ctx context.Context, c repository.CapabilityConfig) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO capability_configs (capability_id,asset_types_json,symbols_json,excluded_symbols_json,strategy_scope_json,updated_at) VALUES (?,?,?,?,?,?)`, c.CapabilityID, nullString(c.AssetTypesJSON), nullString(c.SymbolsJSON), nullString(c.ExcludedSymbolsJSON), nullString(c.StrategyScopeJSON), c.UpdatedAt)
	return apperr.FromRepositoryError(err)
}

// GetLatestCapabilityConfig reads the latest capability scope settings.
func (r *SettingsRepository) GetLatestCapabilityConfig(ctx context.Context) (repository.CapabilityConfig, error) {
	var c repository.CapabilityConfig
	err := r.db.QueryRowContext(ctx, `SELECT capability_id,COALESCE(asset_types_json,'[]'),COALESCE(symbols_json,'[]'),COALESCE(excluded_symbols_json,'[]'),COALESCE(strategy_scope_json,'[]'),updated_at FROM capability_configs ORDER BY updated_at DESC LIMIT 1`).Scan(&c.CapabilityID, &c.AssetTypesJSON, &c.SymbolsJSON, &c.ExcludedSymbolsJSON, &c.StrategyScopeJSON, &c.UpdatedAt)
	return c, apperr.FromRepositoryError(err)
}
