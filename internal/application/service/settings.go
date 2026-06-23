package service

import (
	"context"
	"encoding/json"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

// SettingsService handles non-rule system settings and capability settings.
type SettingsService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

// NewSettingsService creates a settings service.
func NewSettingsService(tx repository.Transactor) *SettingsService {
	return &SettingsService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

// UpdateSystemSettings persists ordinary local preferences.
func (s *SettingsService) UpdateSystemSettings(ctx context.Context, requestID string, req dto.SystemSettingsDTO) error {
	sources, _ := json.Marshal(req.DataSources)
	notification, _ := json.Marshal(map[string]any{"enabled": req.NotificationEnabled, "page_preference": req.PagePreference})
	return s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.SettingsRepo.SaveSystemSettings(ctx, repository.SystemSettings{SettingsID: s.ids.New("settings"), NotificationConfigJSON: string(notification), DataSourcesJSON: string(sources), UpdatedAt: s.clk.NowRFC3339()}); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: s.ids.New("audit"), RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionUpdateSettings), Status: string(model.AuditStatusSuccess), CreatedAt: s.clk.NowRFC3339()})
	})
}

// UpdateCapabilitySettings persists capability scope settings.
func (s *SettingsService) UpdateCapabilitySettings(ctx context.Context, requestID string, req dto.CapabilitySettingsDTO) error {
	assetTypes := firstNonEmptySlice(req.AssetTypes, req.AllowedAssetTypes)
	symbols := firstNonEmptySlice(req.Symbols, req.AllowedSymbols)
	assets, _ := json.Marshal(assetTypes)
	symbolsRaw, _ := json.Marshal(symbols)
	excluded, _ := json.Marshal(req.ExcludedSymbols)
	scope, _ := json.Marshal(req.StrategyScope)
	return s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.SettingsRepo.SaveCapabilityConfig(ctx, repository.CapabilityConfig{CapabilityID: s.ids.New("cap"), AssetTypesJSON: string(assets), SymbolsJSON: string(symbolsRaw), ExcludedSymbolsJSON: string(excluded), StrategyScopeJSON: string(scope), UpdatedAt: s.clk.NowRFC3339()}); err != nil {
			return err
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: s.ids.New("audit"), RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionUpdateCapability), Status: string(model.AuditStatusSuccess), CreatedAt: s.clk.NowRFC3339()})
	})
}

// GetCapabilitySettings returns capability scope.
func (s *SettingsService) GetCapabilitySettings(ctx context.Context) (dto.CapabilitySettingsDTO, error) {
	var out dto.CapabilitySettingsDTO
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if repos.SettingsRepo == nil {
			return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "settings repository not configured")
		}
		cfg, err := repos.SettingsRepo.GetLatestCapabilityConfig(ctx)
		if err != nil {
			return err
		}
		out.CapabilityID = cfg.CapabilityID
		out.AssetTypes = splitJSONStrings(cfg.AssetTypesJSON)
		out.Symbols = splitJSONStrings(cfg.SymbolsJSON)
		out.ExcludedSymbols = splitJSONStrings(cfg.ExcludedSymbolsJSON)
		out.StrategyScope = splitJSONStrings(cfg.StrategyScopeJSON)
		out.AllowedAssetTypes = out.AssetTypes
		out.AllowedSymbols = out.Symbols
		out.UpdatedAt = cfg.UpdatedAt
		return nil
	}); err != nil {
		return dto.CapabilitySettingsDTO{}, err
	}
	return out, nil
}

func splitJSONStrings(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	_ = json.Unmarshal([]byte(raw), &out)
	return out
}

func firstNonEmptySlice(primary, fallback []string) []string {
	if len(primary) > 0 {
		return primary
	}
	return fallback
}
