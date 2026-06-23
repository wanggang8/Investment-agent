package service

import (
	"context"

	"investment-agent/internal/application/dto"
)

// ConfirmationWriter 定义确认写入能力。
type ConfirmationWriter interface {
	Confirm(ctx context.Context, requestID, decisionID string, req dto.ConfirmationRequest) (dto.ConfirmationResponse, error)
}

// PortfolioWriter 定义账户快照写入能力。
type PortfolioWriter interface {
	WriteSnapshot(ctx context.Context, requestID string, req dto.PortfolioInitRequest, source string) (dto.PortfolioWriteResponse, error)
	WriteAdjustment(ctx context.Context, requestID string, req dto.PortfolioAdjustmentRequest, source string) (dto.PortfolioWriteResponse, error)
	EditHolding(ctx context.Context, requestID string, req dto.HoldingEditRequest) (dto.LocalFactWriteResponse, error)
	RemoveHolding(ctx context.Context, requestID string, req dto.HoldingRemoveRequest) (dto.LocalFactWriteResponse, error)
	RecordOfflineTransaction(ctx context.Context, requestID string, req dto.OfflineTransactionRequest) (dto.LocalFactWriteResponse, error)
	ValidateImport(ctx context.Context, requestID string, req dto.BatchImportValidationRequest) (dto.BatchImportValidationResponse, error)
	ConfirmImport(ctx context.Context, requestID string, req dto.BatchImportConfirmRequest) (dto.LocalFactWriteResponse, error)
	CorrectFact(ctx context.Context, requestID string, req dto.CorrectionRequest) (dto.LocalFactWriteResponse, error)
	ReviewQuarterlyRebalance(ctx context.Context, requestID string, req dto.RebalanceReviewRequest) (dto.RebalanceReviewResponse, error)
}

// RuleProposalWriter 定义规则提案确认能力。
type RuleProposalWriter interface {
	ConfirmProposal(ctx context.Context, requestID, proposalID string, req dto.RuleProposalConfirmRequest, final bool) (dto.RuleProposalConfirmResponse, error)
	GenerateSOPAddendumProposal(ctx context.Context, requestID string, req dto.SOPAddendumProposalRequest) (dto.SOPAddendumProposalResponse, error)
}
