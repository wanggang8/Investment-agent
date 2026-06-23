package service

import (
	"context"
	"encoding/json"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

// MarketService handles market refresh auxiliary writes.
type MarketService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

// NewMarketService creates a market service.
func NewMarketService(tx repository.Transactor) *MarketService {
	return &MarketService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

// MarketSnapshotExists checks for duplicate market snapshot IDs.
func (s *MarketService) MarketSnapshotExists(ctx context.Context, snapshotID string) (bool, error) {
	var exists bool
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		got, err := repos.MarketRepo.MarketSnapshotExists(ctx, snapshotID)
		if err != nil {
			return err
		}
		exists = got
		return nil
	}); err != nil {
		return false, err
	}
	return exists, nil
}

// AppendRefreshAudit persists an audit event for market refresh failures or degraded results.
func (s *MarketService) AppendRefreshAudit(ctx context.Context, requestID, errorCode, outputRef string, failures []dto.MarketRefreshFailure) (string, error) {
	status := string(model.AuditStatusDegraded)
	if errorCode != "" {
		status = string(model.AuditStatusFailed)
	}
	if outputRef == "" && len(failures) > 0 {
		buf, _ := json.Marshal(failures)
		outputRef = string(buf)
	}
	auditID := s.ids.New("audit")
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorSystem), Action: string(model.AuditActionRefreshMarketData), Status: status, ErrorCode: errorCode, OutputRefType: "failed_symbols", OutputRef: outputRef, CreatedAt: s.clk.NowRFC3339()})
	})
	if err != nil {
		return "", err
	}
	return auditID, nil
}
