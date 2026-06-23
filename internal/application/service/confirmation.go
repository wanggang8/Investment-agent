package service

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
)

// ConfirmationService 处理用户确认写入。
type ConfirmationService struct {
	tx  repository.Transactor
	clk clock.Clock
	ids idgen.Generator
}

// NewConfirmationService 创建确认写入服务。
func NewConfirmationService(tx repository.Transactor) *ConfirmationService {
	return &ConfirmationService{tx: tx, clk: clock.SystemClock{}, ids: idgen.NewGenerator()}
}

// Confirm 写入确认、相关事实和审计。
func (s *ConfirmationService) Confirm(ctx context.Context, requestID, decisionID string, req dto.ConfirmationRequest) (dto.ConfirmationResponse, error) {
	var recordType, currentStatus string
	if err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		var err error
		recordType, currentStatus, err = repos.DecisionRepo.GetDecisionConfirmationState(ctx, decisionID)
		return err
	}); err != nil {
		return dto.ConfirmationResponse{}, err
	}
	return s.confirmWithStatus(ctx, requestID, decisionID, recordType, currentStatus, req)
}

func (s *ConfirmationService) confirmWithStatus(ctx context.Context, requestID, decisionID, recordType, currentStatus string, req dto.ConfirmationRequest) (dto.ConfirmationResponse, error) {
	if recordType != "formal_trade_advice" || currentStatus == string(model.ConfirmationNotRequired) {
		return dto.ConfirmationResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "该决策不需要用户确认")
	}
	if currentStatus == string(model.ConfirmationExecutedManually) || currentStatus == string(model.ConfirmationMarkedError) || (currentStatus == string(model.ConfirmationWatch) && req.ConfirmationType == string(model.ConfirmationTypeExecutedManually)) {
		return dto.ConfirmationResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "终态确认不能重复提交")
	}
	if req.Fees < 0 {
		return dto.ConfirmationResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "fees 不能为负数")
	}
	confirmationType := model.ConfirmationType(req.ConfirmationType)
	if !confirmationType.Valid() {
		return dto.ConfirmationResponse{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "confirmation_type 不合法")
	}
	if err := validateConfirmationRequest(req, s.clk.Now()); err != nil {
		return dto.ConfirmationResponse{}, err
	}
	payload, _ := json.Marshal(req)
	now := s.clk.NowRFC3339()
	confirmationID := s.ids.New("confirm")
	auditID := s.ids.New("audit")
	errorCaseID := ""
	transactionID := ""
	snapshotID := ""
	if req.ConfirmationType == string(model.ConfirmationTypeMarkedError) {
		errorCaseID = s.ids.New("err")
	}
	err := s.tx.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		updated, err := repos.DecisionRepo.UpdateDecisionConfirmationStatusIfCurrent(ctx, decisionID, currentStatus, req.ConfirmationType)
		if err != nil {
			return err
		}
		if !updated {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "确认状态已变化，不能重复提交")
		}
		confirmation := repository.OperationConfirmation{ConfirmationID: confirmationID, DecisionID: decisionID, ConfirmationType: req.ConfirmationType, OperationType: req.OperationType, Symbol: req.Symbol, Quantity: req.Quantity, Price: req.Price, Fees: req.Fees, ExecutedAt: req.ExecutedAt, ErrorCaseID: errorCaseID, PayloadJSON: string(payload), Note: req.Note, CreatedAt: now}
		if err := repos.DecisionRepo.SaveOperationConfirmation(ctx, confirmation); err != nil {
			return err
		}
		if req.ConfirmationType == string(model.ConfirmationTypeExecutedManually) {
			var err error
			transactionID, snapshotID, err = s.persistManualExecution(ctx, repos, confirmationID, req, now)
			if err != nil {
				return err
			}
		}
		if req.ConfirmationType == string(model.ConfirmationTypeMarkedError) {
			if err := repos.DecisionRepo.SaveErrorCase(ctx, repository.ErrorCase{ErrorCaseID: errorCaseID, DecisionID: decisionID, ConfirmationID: confirmationID, ActualOutcome: req.ActualOutcome, RootCauseTag: req.RootCauseTag, LessonLearned: req.LessonLearned, CreatedAt: now}); err != nil {
				return err
			}
		}
		action := string(model.AuditActionConfirmOperation)
		if req.ConfirmationType == string(model.ConfirmationTypeMarkedError) {
			action = string(model.AuditActionMarkError)
		}
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{AuditEventID: auditID, RequestID: requestID, DecisionID: decisionID, Actor: string(model.AuditActorUser), Action: action, ConfirmationID: confirmationID, ErrorCaseID: errorCaseID, Status: string(model.AuditStatusSuccess), BeforeState: currentStatus, AfterState: req.ConfirmationType, SnapshotID: snapshotID, CreatedAt: now})
	})
	if err != nil {
		return dto.ConfirmationResponse{}, err
	}
	transactionIDs := []string{}
	if transactionID != "" {
		transactionIDs = append(transactionIDs, transactionID)
	}
	return dto.ConfirmationResponse{ConfirmationID: confirmationID, DecisionID: decisionID, ConfirmationStatus: req.ConfirmationType, ErrorCaseID: errorCaseID, TransactionIDs: transactionIDs, SnapshotID: snapshotID, AuditEventIDs: []string{auditID}}, nil
}

func validateConfirmationRequest(req dto.ConfirmationRequest, now time.Time) error {
	switch req.ConfirmationType {
	case string(model.ConfirmationTypeExecutedManually):
		if req.Symbol == "" || req.Quantity <= 0 || req.Price <= 0 || req.ExecutedAt == "" {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "手工执行确认缺少交易字段")
		}
		if !model.OperationType(req.OperationType).Valid() {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "operation_type 不合法")
		}
		executedAt, err := time.Parse(time.RFC3339, req.ExecutedAt)
		if err != nil {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "executed_at 必须是 RFC3339 时间")
		}
		if executedAt.After(now) {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "executed_at 不能晚于当前时间")
		}
	case string(model.ConfirmationTypeMarkedError):
		if req.ActualOutcome == "" || req.RootCauseTag == "" || req.LessonLearned == "" {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "错误标记确认缺少复盘字段")
		}
		if !model.RootCauseTag(req.RootCauseTag).Valid() {
			return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "root_cause_tag 不合法")
		}
	}
	return nil
}

func (s *ConfirmationService) persistManualExecution(ctx context.Context, repos repository.Repositories, confirmationID string, req dto.ConfirmationRequest, now string) (string, string, error) {
	transactionID := s.ids.New("tx")
	snapshotID := s.ids.New("snap_confirm")
	positions, err := repos.PortfolioRepo.ListPositions(ctx)
	if err != nil {
		return "", "", err
	}
	latestSnapshot, err := repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
	if err != nil {
		latestSnapshot = repository.PortfolioSnapshot{}
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
	tradeAmount := req.Price * req.Quantity
	cashImpact := tradeAmount + req.Fees
	if model.OperationType(req.OperationType) == model.OperationBuy && latestSnapshot.Cash < cashImpact {
		return "", "", apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "现金不足，不能记录买入")
	}
	after, err := applyManualOperation(before, req, now, s.ids.New)
	if err != nil {
		return "", "", err
	}
	beforeJSON := jsonStringOrEmpty(before)
	afterJSON := jsonStringOrEmpty(after)
	if err := repos.DecisionRepo.SavePositionTransaction(ctx, repository.PositionTransaction{TransactionID: transactionID, ConfirmationID: confirmationID, Symbol: req.Symbol, OperationType: req.OperationType, Quantity: req.Quantity, Price: req.Price, Fees: req.Fees, OccurredAt: valueOr(req.ExecutedAt, now), BeforePositionJSON: beforeJSON, AfterPositionJSON: afterJSON, CreatedAt: now}); err != nil {
		return "", "", err
	}
	if after.Quantity > 0 {
		if err := repos.PortfolioRepo.SavePosition(ctx, after); err != nil {
			return "", "", err
		}
		if idx >= 0 {
			positions[idx] = after
		} else {
			positions = append(positions, after)
		}
	} else if idx >= 0 {
		if err := repos.PortfolioRepo.DeletePosition(ctx, before.PositionID); err != nil {
			return "", "", err
		}
		positions = append(positions[:idx], positions[idx+1:]...)
	}
	active := make([]repository.Position, 0, len(positions))
	for _, item := range positions {
		if item.Quantity > 0 {
			active = append(active, item)
		}
	}
	cash := latestSnapshot.Cash
	switch model.OperationType(req.OperationType) {
	case model.OperationBuy:
		cash -= cashImpact
	case model.OperationSell, model.OperationReduce:
		cash += tradeAmount - req.Fees
	}
	snapshots := make([]repository.PositionSnapshot, 0, len(active))
	positionValue := 0.0
	highRiskAssets := 0.0
	for _, item := range active {
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
	if err := repos.PortfolioRepo.SavePortfolioSnapshot(ctx, repository.PortfolioSnapshot{SnapshotID: snapshotID, SnapshotTime: now, Cash: cash, TotalAssets: totalAssets, CashRatio: cashRatio, HighRiskRatio: highRiskRatio, PositionCount: len(active), Source: "manual", CreatedAt: now}, snapshots); err != nil {
		return "", "", err
	}
	return transactionID, snapshotID, nil
}

func applyManualOperation(before repository.Position, req dto.ConfirmationRequest, now string, newID func(string) string) (repository.Position, error) {
	after := before
	if after.PositionID == "" {
		after.PositionID = newID("pos")
		after.Symbol = req.Symbol
		after.Name = req.Symbol
		after.CostPrice = req.Price
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

func jsonStringOrEmpty(v repository.Position) string {
	if v.PositionID == "" {
		return "{}"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func valueOr(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
