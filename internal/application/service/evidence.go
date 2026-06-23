package service

import (
	"context"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

// EvidenceService handles evidence refresh auxiliary writes and index maintenance.
type EvidenceService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

// NewEvidenceService creates an evidence service.
func NewEvidenceService(tx repository.Transactor) *EvidenceService {
	return &EvidenceService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

// MarkChunksIndexStatus updates RAG chunks as part of evidence refresh.
func (s *EvidenceService) MarkChunksIndexStatus(ctx context.Context, chunkIDs []string, status string) error {
	return s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		return repos.IntelligenceRepo.UpdateRAGChunksIndexStatus(ctx, chunkIDs, status)
	})
}

// CountRAGChunks counts local RAG chunks.
func (s *EvidenceService) CountRAGChunks(ctx context.Context) (int, error) {
	var count int
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		got, err := repos.IntelligenceRepo.CountRAGChunks(ctx)
		if err != nil {
			return err
		}
		count = got
		return nil
	}); err != nil {
		return 0, err
	}
	return count, nil
}

// AppendRebuildAudit records a local index rebuild audit.
func (s *EvidenceService) AppendRebuildAudit(ctx context.Context, requestID string) (string, error) {
	auditID := s.ids.New("audit")
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionRebuildIndex), Status: string(model.AuditStatusSuccess), CreatedAt: s.clk.NowRFC3339()})
	})
	if err != nil {
		return "", err
	}
	return auditID, nil
}
