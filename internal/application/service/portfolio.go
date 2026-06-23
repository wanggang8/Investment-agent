package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

// PortfolioService 处理账户快照写入。
type PortfolioService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

// NewPortfolioService 创建账户写入服务。
func NewPortfolioService(tx repository.Transactor) *PortfolioService {
	return &PortfolioService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

// WriteSnapshot 写入账户快照与持仓快照。
func (s *PortfolioService) WriteSnapshot(ctx context.Context, requestID string, req dto.PortfolioInitRequest, source string) (dto.PortfolioWriteResponse, error) {
	return s.writeSnapshot(ctx, requestID, req, source, "")
}

func (s *PortfolioService) writeSnapshot(ctx context.Context, requestID string, req dto.PortfolioInitRequest, source string, adjustReason string) (dto.PortfolioWriteResponse, error) {
	if err := validatePortfolioSnapshotRequest(req); err != nil {
		return dto.PortfolioWriteResponse{}, err
	}
	now := s.clk.NowRFC3339()
	snapshotID := s.ids.New("snap")
	auditID := s.ids.New("audit")
	positions := make([]repository.PositionSnapshot, 0, len(req.Positions))
	currentPositions := make([]repository.Position, 0, len(req.Positions))
	highRiskAssets := 0.0
	for _, p := range req.Positions {
		positionID := s.ids.New("pos")
		positionSnapshotID := s.ids.New("ps")
		marketValue := p.Quantity * p.CurrentPrice
		positionState := strings.TrimSpace(p.PositionState)
		if positionState == "" {
			positionState = string(model.PositionNormal)
		}
		if positionState == string(model.PositionSellOnly) || positionState == string(model.PositionFrozenWatch) {
			highRiskAssets += marketValue
		}
		currentPositions = append(currentPositions, repository.Position{
			PositionID:            positionID,
			Symbol:                p.Symbol,
			Name:                  p.Name,
			Quantity:              p.Quantity,
			CostPrice:             p.CostPrice,
			CurrentPrice:          p.CurrentPrice,
			MarketValue:           marketValue,
			UnrealizedProfitRatio: profitRatio(p.CostPrice, p.CurrentPrice),
			PositionState:         positionState,
			BuyDate:               p.BuyDate,
			BuyReason:             p.BuyReason,
			AssetTag:              p.AssetTag,
			UpdatedAt:             now,
		})
		positions = append(positions, repository.PositionSnapshot{
			PositionSnapshotID:    positionSnapshotID,
			SnapshotID:            snapshotID,
			PositionID:            positionID,
			Symbol:                p.Symbol,
			Name:                  p.Name,
			Quantity:              p.Quantity,
			CostPrice:             p.CostPrice,
			CurrentPrice:          p.CurrentPrice,
			MarketValue:           marketValue,
			UnrealizedProfitRatio: profitRatio(p.CostPrice, p.CurrentPrice),
			PositionState:         positionState,
			BuyDate:               p.BuyDate,
			BuyReason:             p.BuyReason,
			AssetTag:              p.AssetTag,
			CreatedAt:             now,
		})
	}
	highRiskRatio := 0.0
	if req.TotalAssets > 0 {
		highRiskRatio = highRiskAssets / req.TotalAssets
	}
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.PortfolioRepo.SavePortfolioSnapshot(ctx, repository.PortfolioSnapshot{SnapshotID: snapshotID, SnapshotTime: now, Cash: req.Cash, TotalAssets: req.TotalAssets, CashRatio: req.Cash / req.TotalAssets, HighRiskRatio: highRiskRatio, PositionCount: len(req.Positions), Source: source, CreatedAt: now}, positions); err != nil {
			return err
		}
		if err := repos.PortfolioRepo.ReplacePositions(ctx, currentPositions); err != nil {
			return err
		}
		audit := repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionUpdateSettings), Status: string(model.AuditStatusSuccess), SnapshotID: snapshotID, CreatedAt: now}
		if adjustReason != "" {
			audit.InputRefType = "adjust_reason"
			audit.InputRef = adjustReason
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, audit)
	}); err != nil {
		return dto.PortfolioWriteResponse{}, err
	}
	return dto.PortfolioWriteResponse{SnapshotID: snapshotID, PositionCount: len(req.Positions), PositionSnapshotCount: len(req.Positions), AuditEventIDs: []string{auditID}}, nil
}

// WritePortfolioSnapshot is a compatibility wrapper for the handler.
func (s *PortfolioService) WritePortfolioSnapshot(ctx context.Context, requestID string, req dto.PortfolioInitRequest, source string) (dto.PortfolioWriteResponse, error) {
	return s.WriteSnapshot(ctx, requestID, req, source)
}

func (s *PortfolioService) WriteAdjustment(ctx context.Context, requestID string, req dto.PortfolioAdjustmentRequest, source string) (dto.PortfolioWriteResponse, error) {
	return s.writeSnapshot(ctx, requestID, dto.PortfolioInitRequest{Cash: req.Cash, TotalAssets: req.TotalAssets, Positions: req.Positions}, source, req.AdjustReason)
}

func (s *PortfolioService) EditHolding(ctx context.Context, requestID string, req dto.HoldingEditRequest) (dto.LocalFactWriteResponse, error) {
	if strings.TrimSpace(req.Reason) == "" || strings.TrimSpace(req.Confirmation) == "" {
		return dto.LocalFactWriteResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "reason 和 confirmation 不能为空")
	}
	if err := validatePositionInput(req.Position); err != nil {
		return dto.LocalFactWriteResponse{}, err
	}
	now := s.clk.NowRFC3339()
	snapshotID := s.ids.New("snap")
	auditID := s.ids.New("audit")
	positionID := strings.TrimSpace(req.PositionID)
	positionInput := req.Position
	var out dto.LocalFactWriteResponse
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		latestSnapshot, err := repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
		if err != nil {
			if !apperr.IsCode(err, apperr.CodeNotFound) {
				return err
			}
			latestSnapshot = repository.PortfolioSnapshot{}
		}
		positions, err := repos.PortfolioRepo.ListPositions(ctx)
		if err != nil {
			return err
		}
		idx := -1
		for i, item := range positions {
			if item.PositionID == positionID || (positionID == "" && item.Symbol == positionInput.Symbol) {
				idx = i
				positionID = item.PositionID
				break
			}
		}
		if positionID == "" {
			positionID = s.ids.New("pos")
		}
		position := positionFromInput(positionID, positionInput, now)
		if idx >= 0 {
			positions[idx] = position
		} else {
			positions = append(positions, position)
		}
		snapshot, positionSnapshots := s.snapshotFromPositions(snapshotID, now, positions, latestSnapshot.Cash)
		if err := repos.PortfolioRepo.SavePortfolioSnapshot(ctx, snapshot, positionSnapshots); err != nil {
			return err
		}
		if err := repos.PortfolioRepo.ReplacePositions(ctx, positions); err != nil {
			return err
		}
		if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionUpdateSettings), Status: string(model.AuditStatusSuccess), SnapshotID: snapshotID, InputRefType: "holding_edit_reason", InputRef: req.Reason, CreatedAt: now}); err != nil {
			return err
		}
		out = dto.LocalFactWriteResponse{SnapshotID: snapshotID, PositionID: positionID, AuditEventIDs: []string{auditID}, SafetyStatement: localFactSafetyStatement}
		return nil
	})
	return out, err
}

func (s *PortfolioService) RemoveHolding(ctx context.Context, requestID string, req dto.HoldingRemoveRequest) (dto.LocalFactWriteResponse, error) {
	if strings.TrimSpace(req.PositionID) == "" || strings.TrimSpace(req.Reason) == "" || strings.TrimSpace(req.Confirmation) == "" {
		return dto.LocalFactWriteResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "position_id、reason 和 confirmation 不能为空")
	}
	now := s.clk.NowRFC3339()
	snapshotID := s.ids.New("snap")
	auditID := s.ids.New("audit")
	var out dto.LocalFactWriteResponse
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		latestSnapshot, err := repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
		if err != nil {
			if !apperr.IsCode(err, apperr.CodeNotFound) {
				return err
			}
			latestSnapshot = repository.PortfolioSnapshot{}
		}
		positions, err := repos.PortfolioRepo.ListPositions(ctx)
		if err != nil {
			return err
		}
		idx := -1
		for i, item := range positions {
			if item.PositionID == req.PositionID {
				idx = i
				break
			}
		}
		if idx < 0 {
			return apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "position not found")
		}
		positions = append(positions[:idx], positions[idx+1:]...)
		snapshot, positionSnapshots := s.snapshotFromPositions(snapshotID, now, positions, latestSnapshot.Cash)
		if err := repos.PortfolioRepo.SavePortfolioSnapshot(ctx, snapshot, positionSnapshots); err != nil {
			return err
		}
		if err := repos.PortfolioRepo.ReplacePositions(ctx, positions); err != nil {
			return err
		}
		if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionUpdateSettings), Status: string(model.AuditStatusSuccess), SnapshotID: snapshotID, InputRefType: "holding_remove_reason", InputRef: req.Reason, CreatedAt: now}); err != nil {
			return err
		}
		out = dto.LocalFactWriteResponse{SnapshotID: snapshotID, PositionID: req.PositionID, AuditEventIDs: []string{auditID}, SafetyStatement: localFactSafetyStatement}
		return nil
	})
	return out, err
}

func (s *PortfolioService) RecordOfflineTransaction(ctx context.Context, requestID string, req dto.OfflineTransactionRequest) (dto.LocalFactWriteResponse, error) {
	if err := validateOfflineTransactionRequest(req, s.clk.Now()); err != nil {
		return dto.LocalFactWriteResponse{}, err
	}
	now := s.clk.NowRFC3339()
	confirmationID := s.ids.New("confirm")
	transactionID := s.ids.New("tx")
	snapshotID := s.ids.New("snap")
	auditID := s.ids.New("audit")
	var out dto.LocalFactWriteResponse
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		latestSnapshot, err := repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
		if err != nil {
			if !apperr.IsCode(err, apperr.CodeNotFound) {
				return err
			}
			latestSnapshot = repository.PortfolioSnapshot{}
		}
		positions, err := repos.PortfolioRepo.ListPositions(ctx)
		if err != nil {
			return err
		}
		idx := -1
		for i, item := range positions {
			if item.Symbol == req.Symbol {
				idx = i
				break
			}
		}
		before := repository.Position{}
		if idx >= 0 {
			before = positions[idx]
		}
		cash := latestSnapshot.Cash
		tradeAmount := req.Price * req.Quantity
		switch model.OperationType(req.OperationType) {
		case model.OperationBuy:
			cash -= tradeAmount + req.Fees
		case model.OperationSell, model.OperationReduce:
			cash += tradeAmount - req.Fees
		}
		if cash < -0.01 {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "现金不足，不能记录买入")
		}
		after, err := applyOfflineTransaction(before, req, now, s.ids.New)
		if err != nil {
			return err
		}
		payload, _ := json.Marshal(req)
		confirmation := repository.OperationConfirmation{ConfirmationID: confirmationID, DecisionID: "", ConfirmationType: string(model.ConfirmationTypeExecutedManually), OperationType: req.OperationType, Symbol: req.Symbol, Quantity: req.Quantity, Price: req.Price, Fees: req.Fees, ExecutedAt: req.ExecutedAt, PayloadJSON: string(payload), Note: req.Note, CreatedAt: now}
		if err := repos.DecisionRepo.SaveOperationConfirmation(ctx, confirmation); err != nil {
			return err
		}
		if err := repos.DecisionRepo.SavePositionTransaction(ctx, repository.PositionTransaction{TransactionID: transactionID, ConfirmationID: confirmationID, Symbol: req.Symbol, OperationType: req.OperationType, Quantity: req.Quantity, Price: req.Price, Fees: req.Fees, OccurredAt: req.ExecutedAt, BeforePositionJSON: jsonStringOrEmpty(before), AfterPositionJSON: jsonStringOrEmpty(after), CreatedAt: now}); err != nil {
			return err
		}
		if after.Quantity > 0 {
			if idx >= 0 {
				positions[idx] = after
			} else {
				positions = append(positions, after)
			}
		} else if idx >= 0 {
			positions = append(positions[:idx], positions[idx+1:]...)
		}
		snapshot, positionSnapshots := s.snapshotFromPositions(snapshotID, now, positions, cash)
		if err := repos.PortfolioRepo.SavePortfolioSnapshot(ctx, snapshot, positionSnapshots); err != nil {
			return err
		}
		if err := repos.PortfolioRepo.ReplacePositions(ctx, positions); err != nil {
			return err
		}
		if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionConfirmOperation), ConfirmationID: confirmationID, Status: string(model.AuditStatusSuccess), SnapshotID: snapshotID, InputRefType: "offline_transaction_note", InputRef: req.Note, CreatedAt: now}); err != nil {
			return err
		}
		out = dto.LocalFactWriteResponse{SnapshotID: snapshotID, PositionID: after.PositionID, TransactionID: transactionID, AuditEventIDs: []string{auditID}, SafetyStatement: localFactSafetyStatement}
		return nil
	})
	return out, err
}

func (s *PortfolioService) ValidateImport(ctx context.Context, requestID string, req dto.BatchImportValidationRequest) (dto.BatchImportValidationResponse, error) {
	batchID := s.ids.New("import")
	rows := make([]dto.BatchImportRowResult, 0, len(req.Rows))
	validCount := 0
	for _, row := range req.Rows {
		result := validateImportRow(row, s.clk.Now())
		if result.Valid {
			validCount++
		}
		rows = append(rows, result)
	}
	summary := dto.BatchImportValidationSummary{RowCount: len(req.Rows), ValidCount: validCount, InvalidCount: len(req.Rows) - validCount}
	summaryBytes, _ := json.Marshal(summary)
	rowsHash := importRowsHash(req.Rows)
	now := s.clk.NowRFC3339()
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		return repos.PortfolioRepo.SaveLocalAccountImportBatch(ctx, repository.LocalAccountImportBatch{ImportBatchID: batchID, RequestID: requestID, Status: "validated", RowCount: summary.RowCount, ValidCount: summary.ValidCount, InvalidCount: summary.InvalidCount, ValidationSummaryJSON: string(summaryBytes), RowsHash: rowsHash, CreatedAt: now})
	}); err != nil {
		return dto.BatchImportValidationResponse{}, err
	}
	return dto.BatchImportValidationResponse{ImportBatchID: batchID, Summary: summary, Rows: rows}, nil
}

func (s *PortfolioService) ConfirmImport(ctx context.Context, requestID string, req dto.BatchImportConfirmRequest) (dto.LocalFactWriteResponse, error) {
	if strings.TrimSpace(req.ConfirmReason) == "" {
		return dto.LocalFactWriteResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "confirm_reason 不能为空")
	}
	if strings.TrimSpace(req.ImportBatchID) == "" {
		return dto.LocalFactWriteResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "import_batch_id 不能为空")
	}
	validCount := 0
	for _, row := range req.Rows {
		result := validateImportRow(row, s.clk.Now())
		if !result.Valid {
			return dto.LocalFactWriteResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "批量导入存在校验错误")
		}
		validCount++
	}
	summary := dto.BatchImportValidationSummary{RowCount: len(req.Rows), ValidCount: validCount, InvalidCount: 0}
	rowsHash := importRowsHash(req.Rows)
	now := s.clk.NowRFC3339()
	batchID := strings.TrimSpace(req.ImportBatchID)
	snapshotID := s.ids.New("snap")
	auditID := s.ids.New("audit")
	var out dto.LocalFactWriteResponse
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		batch, err := repos.PortfolioRepo.GetLocalAccountImportBatch(ctx, batchID)
		if err != nil {
			return err
		}
		if batch.Status != "validated" || batch.InvalidCount != 0 {
			return apperr.New(apperr.CodeConflict, apperr.CategoryConflict, "批量导入必须先完成无错误校验")
		}
		if batch.RowsHash == "" || batch.RowsHash != rowsHash {
			return apperr.New(apperr.CodeConflict, apperr.CategoryConflict, "批量导入内容与已校验批次不一致")
		}
		latestSnapshot, err := repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
		if err != nil {
			if !apperr.IsCode(err, apperr.CodeNotFound) {
				return err
			}
			latestSnapshot = repository.PortfolioSnapshot{}
		}
		positions, err := repos.PortfolioRepo.ListPositions(ctx)
		if err != nil {
			return err
		}
		cash := latestSnapshot.Cash
		for _, row := range req.Rows {
			switch row.RowType {
			case "holding":
				positionID := s.ids.New("pos")
				positions = append(positions, positionFromInput(positionID, dto.PositionInput{Symbol: row.Symbol, Name: row.Name, Quantity: row.Quantity, CostPrice: row.CostPrice, CurrentPrice: row.CurrentPrice, BuyDate: row.BuyDate, BuyReason: row.BuyReason, PositionState: row.PositionState, AssetTag: row.AssetTag}, now))
			case "transaction":
				txReq := dto.OfflineTransactionRequest{OperationType: row.OperationType, Symbol: row.Symbol, Name: row.Name, Quantity: row.Quantity, Price: row.Price, Fees: row.Fees, ExecutedAt: row.OccurredAt, BuyReason: row.BuyReason, AssetTag: row.AssetTag, Note: req.ConfirmReason}
				idx := -1
				for i, item := range positions {
					if item.Symbol == txReq.Symbol {
						idx = i
						break
					}
				}
				before := repository.Position{}
				if idx >= 0 {
					before = positions[idx]
				}
				tradeAmount := txReq.Price * txReq.Quantity
				switch model.OperationType(txReq.OperationType) {
				case model.OperationBuy:
					cash -= tradeAmount + txReq.Fees
				case model.OperationSell, model.OperationReduce:
					cash += tradeAmount - txReq.Fees
				}
				if cash < -0.01 {
					return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "现金不足，不能记录买入")
				}
				after, err := applyOfflineTransaction(before, txReq, now, s.ids.New)
				if err != nil {
					return err
				}
				confirmationID := s.ids.New("confirm")
				transactionID := s.ids.New("tx")
				payload, _ := json.Marshal(txReq)
				confirmation := repository.OperationConfirmation{ConfirmationID: confirmationID, DecisionID: "", ConfirmationType: string(model.ConfirmationTypeExecutedManually), OperationType: txReq.OperationType, Symbol: txReq.Symbol, Quantity: txReq.Quantity, Price: txReq.Price, Fees: txReq.Fees, ExecutedAt: txReq.ExecutedAt, PayloadJSON: string(payload), Note: txReq.Note, CreatedAt: now}
				if err := repos.DecisionRepo.SaveOperationConfirmation(ctx, confirmation); err != nil {
					return err
				}
				if err := repos.DecisionRepo.SavePositionTransaction(ctx, repository.PositionTransaction{TransactionID: transactionID, ConfirmationID: confirmationID, Symbol: txReq.Symbol, OperationType: txReq.OperationType, Quantity: txReq.Quantity, Price: txReq.Price, Fees: txReq.Fees, OccurredAt: txReq.ExecutedAt, BeforePositionJSON: jsonStringOrEmpty(before), AfterPositionJSON: jsonStringOrEmpty(after), CreatedAt: now}); err != nil {
					return err
				}
				if after.Quantity > 0 {
					if idx >= 0 {
						positions[idx] = after
					} else {
						positions = append(positions, after)
					}
				} else if idx >= 0 {
					positions = append(positions[:idx], positions[idx+1:]...)
				}
			}
		}
		snapshot, positionSnapshots := s.snapshotFromPositions(snapshotID, now, positions, cash)
		if err := repos.PortfolioRepo.SavePortfolioSnapshot(ctx, snapshot, positionSnapshots); err != nil {
			return err
		}
		if err := repos.PortfolioRepo.ReplacePositions(ctx, positions); err != nil {
			return err
		}
		summaryBytes, _ := json.Marshal(summary)
		if err := repos.PortfolioRepo.SaveLocalAccountImportBatch(ctx, repository.LocalAccountImportBatch{ImportBatchID: batchID, RequestID: requestID, Status: "committed", RowCount: summary.RowCount, ValidCount: summary.ValidCount, InvalidCount: summary.InvalidCount, ValidationSummaryJSON: string(summaryBytes), RowsHash: rowsHash, CreatedAt: now, CommittedAt: now}); err != nil {
			return err
		}
		if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionUpdateSettings), Status: string(model.AuditStatusSuccess), SnapshotID: snapshotID, InputRefType: "import_batch", InputRef: batchID, CreatedAt: now}); err != nil {
			return err
		}
		out = dto.LocalFactWriteResponse{SnapshotID: snapshotID, ImportBatchID: batchID, AuditEventIDs: []string{auditID}, SafetyStatement: localFactSafetyStatement}
		return nil
	})
	return out, err
}

func (s *PortfolioService) CorrectFact(ctx context.Context, requestID string, req dto.CorrectionRequest) (dto.LocalFactWriteResponse, error) {
	if strings.TrimSpace(req.TargetType) == "" || strings.TrimSpace(req.TargetID) == "" || strings.TrimSpace(req.BeforeJSON) == "" || strings.TrimSpace(req.AfterJSON) == "" || strings.TrimSpace(req.CorrectionReason) == "" {
		return dto.LocalFactWriteResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "修正请求缺少必填字段")
	}
	if !validCorrectionTargetType(req.TargetType) {
		return dto.LocalFactWriteResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "target_type 不合法")
	}
	now := s.clk.NowRFC3339()
	correctionID := s.ids.New("corr")
	auditID := s.ids.New("audit")
	var out dto.LocalFactWriteResponse
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.PortfolioRepo.SaveLocalAccountCorrection(ctx, repository.LocalAccountCorrection{CorrectionID: correctionID, TargetType: req.TargetType, TargetID: req.TargetID, BeforeJSON: req.BeforeJSON, AfterJSON: req.AfterJSON, CorrectionReason: req.CorrectionReason, AuditEventID: auditID, CreatedAt: now}); err != nil {
			return err
		}
		if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionUpdateSettings), Status: string(model.AuditStatusSuccess), BeforeState: req.BeforeJSON, AfterState: req.AfterJSON, InputRefType: "correction_reason", InputRef: req.CorrectionReason, CreatedAt: now}); err != nil {
			return err
		}
		out = dto.LocalFactWriteResponse{CorrectionID: correctionID, AuditEventIDs: []string{auditID}, SafetyStatement: localFactSafetyStatement}
		return nil
	})
	return out, err
}

func (s *PortfolioService) ReviewQuarterlyRebalance(ctx context.Context, requestID string, req dto.RebalanceReviewRequest) (dto.RebalanceReviewResponse, error) {
	if err := validateRebalanceReviewRequest(req); err != nil {
		return dto.RebalanceReviewResponse{}, err
	}
	now := s.clk.NowRFC3339()
	reviewID := s.ids.New("rebalance")
	auditID := s.ids.New("audit")
	threshold := req.DriftThreshold
	if threshold == 0 {
		threshold = 0.15
	}
	reviewDate := strings.TrimSpace(req.ReviewDate)
	if reviewDate == "" {
		reviewDate = now[:10]
	}
	var out dto.RebalanceReviewResponse
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		snapshot, err := repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
		if err != nil {
			return err
		}
		positions, err := repos.PortfolioRepo.ListPositions(ctx)
		if err != nil {
			return err
		}
		out = calculateRebalanceReview(reviewID, reviewDate, snapshot, positions, req, threshold)
		out.AuditEventIDs = []string{auditID}
		if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, Actor: string(model.AuditActorUser), Action: string(model.AuditActionRunLocalTask), Status: string(model.AuditStatusSuccess), SnapshotID: snapshot.SnapshotID, InputRefType: "rebalance_review", InputRef: reviewID, OutputRefType: "manual_recommendations", OutputRef: strconv.Itoa(len(out.Items)), CreatedAt: now}); err != nil {
			return err
		}
		return nil
	})
	return out, err
}

func calculateRebalanceReview(reviewID, reviewDate string, snapshot repository.PortfolioSnapshot, positions []repository.Position, req dto.RebalanceReviewRequest, threshold float64) dto.RebalanceReviewResponse {
	values := map[string]float64{"cash": snapshot.Cash, "core": 0, "satellite": 0}
	for _, position := range positions {
		bucket := rebalanceBucket(position.AssetTag)
		values[bucket] += position.MarketValue
	}
	totalAssets := snapshot.TotalAssets
	if totalAssets == 0 {
		totalAssets = values["cash"] + values["core"] + values["satellite"]
	}
	targets := map[string]float64{"core": req.TargetCoreRatio, "satellite": req.TargetSatelliteRatio, "cash": req.TargetCashRatio}
	items := make([]dto.RebalanceReviewItem, 0, 3)
	for _, bucket := range []string{"core", "satellite", "cash"} {
		actualValue := values[bucket]
		actualRatio := 0.0
		if totalAssets > 0 {
			actualRatio = actualValue / totalAssets
		}
		targetValue := targets[bucket] * totalAssets
		drift := actualRatio - targets[bucket]
		recommendation := "hold"
		amount := 0.0
		if math.Abs(drift) >= threshold {
			amount = math.Abs(actualValue - targetValue)
			if drift > 0 {
				recommendation = "manual_sell_or_reduce"
			} else {
				recommendation = "manual_buy_or_add"
			}
			if bucket == "cash" && drift < 0 {
				recommendation = "manual_raise_cash"
			}
		}
		items = append(items, dto.RebalanceReviewItem{Bucket: bucket, TargetRatio: targets[bucket], ActualRatio: actualRatio, DriftRatio: drift, TargetValue: targetValue, ActualValue: actualValue, Recommendation: recommendation, ManualAmount: amount})
	}
	return dto.RebalanceReviewResponse{ReviewID: reviewID, ReviewDate: reviewDate, TotalAssets: totalAssets, DriftThreshold: threshold, Items: items, SafetyStatement: "季度再平衡仅生成人工计划金额，不连接券商、不自动交易、不创建订单。"}
}

func rebalanceBucket(assetTag string) string {
	switch strings.TrimSpace(assetTag) {
	case "satellite":
		return "satellite"
	case "cash":
		return "cash"
	default:
		return "core"
	}
}

func importRowsHash(rows []dto.BatchImportRow) string {
	payload, _ := json.Marshal(rows)
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func validCorrectionTargetType(targetType string) bool {
	switch targetType {
	case "portfolio_snapshot", "position", "position_snapshot", "position_transaction", "import_batch":
		return true
	default:
		return false
	}
}

func validateRebalanceReviewRequest(req dto.RebalanceReviewRequest) error {
	sum := req.TargetCoreRatio + req.TargetSatelliteRatio + req.TargetCashRatio
	if math.Abs(sum-1) > 0.0001 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "再平衡目标比例之和必须等于 1")
	}
	if req.TargetCoreRatio < 0 || req.TargetSatelliteRatio < 0 || req.TargetCashRatio < 0 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "再平衡目标比例不能为负")
	}
	if req.DriftThreshold < 0 || req.DriftThreshold > 1 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "再平衡偏离阈值必须在 0 到 1 之间")
	}
	return nil
}

func validatePortfolioSnapshotRequest(req dto.PortfolioInitRequest) error {
	if req.Cash < 0 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "cash 必须大于等于 0")
	}
	if req.TotalAssets <= 0 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "total_assets 必须大于 0")
	}
	positionsValue := 0.0
	for _, p := range req.Positions {
		if err := validatePositionInput(p); err != nil {
			return err
		}
		positionsValue += p.Quantity * p.CurrentPrice
	}
	if math.Abs(req.TotalAssets-(req.Cash+positionsValue)) > 0.01 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "total_assets 与现金和持仓市值不一致")
	}
	return nil
}

func validatePositionInput(p dto.PositionInput) error {
	if strings.TrimSpace(p.Symbol) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "symbol 不能为空")
	}
	if strings.TrimSpace(p.Name) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "name 不能为空")
	}
	if p.Quantity < 0 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "quantity 必须大于等于 0")
	}
	if p.CostPrice <= 0 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "cost_price 必须大于 0")
	}
	if p.CurrentPrice < 0 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "current_price 必须大于等于 0")
	}
	if strings.TrimSpace(p.BuyReason) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "buy_reason 不能为空")
	}
	if p.BuyDate != "" {
		if _, err := time.Parse("2006-01-02", p.BuyDate); err != nil {
			if _, err := time.Parse(time.RFC3339, p.BuyDate); err != nil {
				return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "buy_date 格式不合法")
			}
		}
	}
	if strings.TrimSpace(p.PositionState) != "" && !model.PositionState(p.PositionState).Valid() {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "position_state 不合法")
	}
	return nil
}

func validateOfflineTransactionRequest(req dto.OfflineTransactionRequest, now time.Time) error {
	if !model.OperationType(req.OperationType).Valid() {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "operation_type 不合法")
	}
	if strings.TrimSpace(req.Symbol) == "" || req.Quantity <= 0 || req.Price <= 0 || req.ExecutedAt == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "线下交易记录缺少必填字段")
	}
	if req.Fees < 0 {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "fees 不能为负数")
	}
	executedAt, err := time.Parse(time.RFC3339, req.ExecutedAt)
	if err != nil {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "executed_at 必须是 RFC3339 时间")
	}
	if executedAt.After(now) {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "executed_at 不能晚于当前时间")
	}
	if req.OperationType == string(model.OperationBuy) && strings.TrimSpace(req.BuyReason) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "buy_reason 不能为空")
	}
	return nil
}

func validateImportRow(row dto.BatchImportRow, now time.Time) dto.BatchImportRowResult {
	rowNumber := row.RowNumber
	if rowNumber == 0 {
		rowNumber = 1
	}
	result := dto.BatchImportRowResult{RowNumber: rowNumber, Valid: true}
	addErr := func(msg string) {
		result.Valid = false
		result.Errors = append(result.Errors, msg)
	}
	switch row.RowType {
	case "holding":
		if err := validatePositionInput(dto.PositionInput{Symbol: row.Symbol, Name: row.Name, Quantity: row.Quantity, CostPrice: row.CostPrice, CurrentPrice: row.CurrentPrice, BuyDate: row.BuyDate, BuyReason: row.BuyReason, PositionState: row.PositionState, AssetTag: row.AssetTag}); err != nil {
			addErr(err.Error())
		}
	case "transaction":
		if err := validateOfflineTransactionRequest(dto.OfflineTransactionRequest{OperationType: row.OperationType, Symbol: row.Symbol, Name: row.Name, Quantity: row.Quantity, Price: row.Price, Fees: row.Fees, ExecutedAt: row.OccurredAt, BuyReason: row.BuyReason, AssetTag: row.AssetTag}, now); err != nil {
			addErr(err.Error())
		}
	default:
		addErr("row_type 不合法")
	}
	return result
}

func positionFromInput(positionID string, p dto.PositionInput, now string) repository.Position {
	marketValue := p.Quantity * p.CurrentPrice
	positionState := strings.TrimSpace(p.PositionState)
	if positionState == "" {
		positionState = string(model.PositionNormal)
	}
	return repository.Position{PositionID: positionID, Symbol: p.Symbol, Name: p.Name, Quantity: p.Quantity, CostPrice: p.CostPrice, CurrentPrice: p.CurrentPrice, MarketValue: marketValue, UnrealizedProfitRatio: profitRatio(p.CostPrice, p.CurrentPrice), PositionState: positionState, BuyDate: p.BuyDate, BuyReason: p.BuyReason, AssetTag: p.AssetTag, UpdatedAt: now}
}

func (s *PortfolioService) snapshotFromPositions(snapshotID, now string, positions []repository.Position, cash float64) (repository.PortfolioSnapshot, []repository.PositionSnapshot) {
	snapshots := make([]repository.PositionSnapshot, 0, len(positions))
	positionValue := 0.0
	highRiskAssets := 0.0
	for _, item := range positions {
		positionValue += item.MarketValue
		if item.PositionState == string(model.PositionSellOnly) || item.PositionState == string(model.PositionFrozenWatch) {
			highRiskAssets += item.MarketValue
		}
		snapshots = append(snapshots, repository.PositionSnapshot{PositionSnapshotID: s.ids.New("ps"), SnapshotID: snapshotID, PositionID: item.PositionID, Symbol: item.Symbol, Name: item.Name, Quantity: item.Quantity, CostPrice: item.CostPrice, CurrentPrice: item.CurrentPrice, MarketValue: item.MarketValue, UnrealizedProfitRatio: item.UnrealizedProfitRatio, PositionState: item.PositionState, BuyDate: item.BuyDate, BuyReason: item.BuyReason, AssetTag: item.AssetTag, CreatedAt: now})
	}
	totalAssets := cash + positionValue
	highRiskRatio := 0.0
	cashRatio := 0.0
	if totalAssets > 0 {
		highRiskRatio = highRiskAssets / totalAssets
		cashRatio = cash / totalAssets
	}
	return repository.PortfolioSnapshot{SnapshotID: snapshotID, SnapshotTime: now, Cash: cash, TotalAssets: totalAssets, CashRatio: cashRatio, HighRiskRatio: highRiskRatio, PositionCount: len(positions), Source: "manual", CreatedAt: now}, snapshots
}

func applyOfflineTransaction(before repository.Position, req dto.OfflineTransactionRequest, now string, newID func(string) string) (repository.Position, error) {
	after := before
	if after.PositionID == "" {
		after.PositionID = newID("pos")
		after.Symbol = req.Symbol
		after.Name = valueOr(req.Name, req.Symbol)
		after.CostPrice = req.Price
		after.BuyReason = req.BuyReason
		after.AssetTag = req.AssetTag
	}
	switch model.OperationType(req.OperationType) {
	case model.OperationBuy:
		totalCost := after.CostPrice*after.Quantity + req.Price*req.Quantity + req.Fees
		after.Quantity += req.Quantity
		if after.Quantity > 0 {
			after.CostPrice = totalCost / after.Quantity
		}
	case model.OperationSell, model.OperationReduce:
		if before.PositionID == "" || before.Quantity < req.Quantity {
			return repository.Position{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "持仓数量不足，不能记录卖出或减仓")
		}
		after.Quantity = math.Max(0, after.Quantity-req.Quantity)
	}
	after.CurrentPrice = req.Price
	after.MarketValue = after.Quantity * req.Price
	if after.CostPrice > 0 {
		after.UnrealizedProfitRatio = (after.CurrentPrice - after.CostPrice) / after.CostPrice
	}
	if after.PositionState == "" {
		after.PositionState = string(model.PositionNormal)
	}
	after.UpdatedAt = now
	return after, nil
}

const localFactSafetyStatement = "仅记录用户已在线下完成或手动确认的本地事实，不连接券商、不自动交易。"

func profitRatio(cost, current float64) float64 {
	if cost == 0 {
		return 0
	}
	return (current - cost) / cost
}
