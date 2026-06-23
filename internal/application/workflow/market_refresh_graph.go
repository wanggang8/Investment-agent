package workflow

import (
	"context"
	"encoding/json"
	"strings"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// MarketRefreshInput 是市场刷新工作流输入。
type MarketRefreshInput struct {
	RequestID    string
	Symbol       string
	PEPercentile float64
	PBPercentile float64
}

// MarketRefreshOutput 是市场刷新后允许写出的事实。
type MarketRefreshOutput struct {
	MarketSnapshot model.MarketSnapshot
	AuditEvents    []model.AuditEvent
}

// MarketRefreshGraph 独立负责市场数据标准化和快照写入。
type MarketRefreshGraph struct {
	auditWriter AuditWriter
	deps        WorkflowDependencies
}

// NewMarketRefreshGraph 创建市场刷新工作流。
func NewMarketRefreshGraph(writer AuditWriter) *MarketRefreshGraph {
	if writer == nil {
		writer = &MemoryAuditWriter{}
	}
	return &MarketRefreshGraph{auditWriter: writer}
}

// NewMarketRefreshGraphWithDependencies 创建带 SQLite 写入能力的市场刷新工作流。
func NewMarketRefreshGraphWithDependencies(deps WorkflowDependencies) *MarketRefreshGraph {
	return &MarketRefreshGraph{auditWriter: NewRepositoryAuditWriter(deps.AuditRepo), deps: deps}
}

// Run 读取外部市场输入，标准化为 MarketSnapshot，并保存刷新审计。
func (g *MarketRefreshGraph) Run(ctx context.Context, in MarketRefreshInput) (MarketRefreshOutput, error) {
	wf := WorkflowContext{RequestID: in.RequestID, WorkflowType: WorkflowMarketRefresh, Symbol: in.Symbol, RuleVersion: workflowRuleVersion(ctx, g.deps.RuleRepo)}
	now := workflowNowRFC3339()
	point, err := g.deps.marketDataSource().FetchMarketData(ctx, in.Symbol)
	if err != nil {
		appErr := apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "市场数据源不可用", err)
		out, auditErr := g.writeFailureAudit(ctx, &wf, in.Symbol, string(apperr.CodeDataSourceUnavailable))
		if auditErr != nil {
			return out, auditErr
		}
		return out, appErr
	}
	if point.Stale {
		appErr := apperr.New(apperr.CodeDataStale, apperr.CategoryInvalidState, "市场数据已过期")
		out, auditErr := g.writeFailureAudit(ctx, &wf, in.Symbol, string(apperr.CodeDataStale))
		if auditErr != nil {
			return out, auditErr
		}
		return out, appErr
	}
	structuredFields := P88NormalizeStructuredDataMetadata(point.Metadata)
	if !structuredFields.Empty() {
		if point.Metadata == nil {
			point.Metadata = map[string]any{}
		}
		point.Metadata["p88_structured_fields"] = structuredFields
		point.Metadata["p88_structured_fields_provenance"] = "collector_metadata_normalized"
	}
	liquidity := point.LiquidityState
	if liquidity == "" {
		liquidity = model.LiquidityNormal
	}
	sentiment := point.SentimentState
	if sentiment == "" {
		sentiment = model.SentimentNeutral
	}
	pe := point.PEPercentile
	pb := point.PBPercentile
	if point.SourceName == "" {
		if pe == 0 {
			pe = in.PEPercentile
		}
		if pb == 0 {
			pb = in.PBPercentile
		}
	}
	metrics := map[string]any{}
	if point.ClosePrice != 0 {
		metrics["close_price"] = point.ClosePrice
	}
	if point.TurnoverRate != 0 {
		metrics["turnover_rate"] = point.TurnoverRate
	}
	if !structuredFields.Empty() {
		metrics["p88_structured_fields"] = structuredFields
		metrics["p88_structured_fields_provenance"] = "collector_metadata_normalized"
	}
	if point.SourceName != "" {
		metrics["source_name"] = point.SourceName
	}
	if point.SourceLevel != "" {
		metrics["source_level"] = string(point.SourceLevel)
	}
	if point.SourceType != "" {
		metrics["source_type"] = point.SourceType
	}
	if point.TradeDate != "" {
		metrics["trade_date"] = point.TradeDate
	}
	if point.CapturedAt != "" {
		metrics["captured_at"] = point.CapturedAt
	}
	if point.ContentHash != "" {
		metrics["content_hash"] = point.ContentHash
	}
	attachRequestID := metadataHasP34SourceHealth(point.Metadata) && strings.TrimSpace(in.RequestID) != ""
	if attachRequestID {
		metrics["request_id"] = strings.TrimSpace(in.RequestID)
	}
	if len(point.Metadata) > 0 {
		if attachRequestID {
			metrics["metadata"] = sourceMetadataWithRequestID(point.Metadata, in.RequestID)
		} else {
			metrics["metadata"] = point.Metadata
		}
	}
	metricsJSON := "{}"
	if len(metrics) > 0 {
		buf, _ := json.Marshal(metrics)
		metricsJSON = string(buf)
	}
	dedupeSnapshot := point.SourceName != "" && point.SourceType != "" && point.TradeDate != ""
	snapshotID := workflowID("market")
	if dedupeSnapshot {
		snapshotID = workflowStableID("market", stableHash("market", point.SourceName, in.Symbol, point.TradeDate, point.SourceType))
	}
	// 市场刷新只写行情事实和审计事件，不生成交易动作。
	marginBalance := 0.0
	marginBalanceChange := 0.0
	if structuredFields.MarginFinancing != nil {
		marginBalance = structuredFields.MarginFinancing.MarginBalance
		marginBalanceChange = structuredFields.MarginFinancing.BalanceChangeRate
	}
	snapshot := model.MarketSnapshot{MarketSnapshotID: snapshotID, Symbol: in.Symbol, TradeDate: point.TradeDate, ClosePrice: point.ClosePrice, TurnoverRate: point.TurnoverRate, MarginBalance: marginBalance, MarginBalanceChange: marginBalanceChange, PEPercentile: pe, PBPercentile: pb, VolumePercentile: point.VolumePercentile, VolatilityPercentile: point.VolatilityPercentile, LiquidityState: liquidity, SentimentState: sentiment, MarketMetricsJSON: metricsJSON}
	result := NodeResult{Status: StatusSuccess, Audit: AuditFragment{Action: string(model.AuditActionRefreshMarketData), NodeName: "MarketRefreshGraph", NodeAction: "refresh_market", Status: StatusSuccess, InputRefType: "symbol", InputRef: in.Symbol, OutputRefType: "market_snapshot", OutputRef: snapshot.MarketSnapshotID}}
	if g.deps.Transactor != nil && g.deps.MarketRepo != nil {
		err := g.deps.Transactor.WithinTx(ctx, func(txCtx context.Context, repos repository.Repositories) error {
			exists := false
			if dedupeSnapshot {
				var err error
				exists, err = repos.MarketRepo.MarketSnapshotExists(txCtx, snapshot.MarketSnapshotID)
				if err != nil {
					return apperr.Wrap(apperr.CodeMarketSnapshotWriteFailed, apperr.CategoryInternal, "市场快照读取失败", err)
				}
			}
			if !exists {
				if err := repos.MarketRepo.SaveMarketSnapshot(txCtx, snapshot, now); err != nil {
					return apperr.Wrap(apperr.CodeMarketSnapshotWriteFailed, apperr.CategoryInternal, "市场快照写入失败", err)
				}
			}
			return writeAuditEvent(txCtx, repos.AuditRepo, &wf, result)
		})
		if err != nil {
			out, auditErr := g.writeFailureAudit(ctx, &wf, in.Symbol, string(apperr.CodeMarketSnapshotWriteFailed))
			if auditErr != nil {
				return out, auditErr
			}
			return out, err
		}
		return MarketRefreshOutput{MarketSnapshot: snapshot, AuditEvents: wf.AuditEvents}, nil
	}
	if g.deps.MarketRepo != nil {
		exists := false
		if dedupeSnapshot {
			var err error
			exists, err = g.deps.MarketRepo.MarketSnapshotExists(ctx, snapshot.MarketSnapshotID)
			if err != nil {
				out, auditErr := g.writeFailureAudit(ctx, &wf, in.Symbol, string(apperr.CodeMarketSnapshotWriteFailed))
				if auditErr != nil {
					return out, auditErr
				}
				return out, apperr.Wrap(apperr.CodeMarketSnapshotWriteFailed, apperr.CategoryInternal, "市场快照读取失败", err)
			}
		}
		if !exists {
			if err := g.deps.MarketRepo.SaveMarketSnapshot(ctx, snapshot, now); err != nil {
				out, auditErr := g.writeFailureAudit(ctx, &wf, in.Symbol, string(apperr.CodeMarketSnapshotWriteFailed))
				if auditErr != nil {
					return out, auditErr
				}
				return out, apperr.Wrap(apperr.CodeMarketSnapshotWriteFailed, apperr.CategoryInternal, "市场快照写入失败", err)
			}
		}
	}
	if err := g.auditWriter.Write(ctx, &wf, result); err != nil {
		return MarketRefreshOutput{}, err
	}
	return MarketRefreshOutput{MarketSnapshot: snapshot, AuditEvents: wf.AuditEvents}, nil
}

func metadataHasP34SourceHealth(metadata map[string]any) bool {
	if len(metadata) == 0 {
		return false
	}
	_, ok := metadata["p34_source_health"]
	return ok
}

func sourceMetadataWithRequestID(metadata map[string]any, requestID string) map[string]any {
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return metadata
	}
	out := make(map[string]any, len(metadata)+1)
	for key, value := range metadata {
		if key != "p34_source_health" {
			out[key] = value
			continue
		}
		health, ok := value.(map[string]any)
		if !ok {
			out[key] = value
			continue
		}
		healthOut := make(map[string]any, len(health))
		for category, rawItem := range health {
			item, ok := rawItem.(map[string]any)
			if !ok {
				healthOut[category] = rawItem
				continue
			}
			itemOut := make(map[string]any, len(item)+1)
			for itemKey, itemValue := range item {
				itemOut[itemKey] = itemValue
			}
			itemOut["request_id"] = requestID
			healthOut[category] = itemOut
		}
		out[key] = healthOut
	}
	out["request_id"] = requestID
	return out
}

func (g *MarketRefreshGraph) writeFailureAudit(ctx context.Context, wf *WorkflowContext, symbol, code string) (MarketRefreshOutput, error) {
	result := NodeResult{Status: StatusFailed, ErrorCode: code, Audit: AuditFragment{Action: string(model.AuditActionRefreshMarketData), NodeName: "MarketRefreshGraph", NodeAction: "refresh_market", Status: StatusFailed, InputRefType: "symbol", InputRef: symbol, OutputRefType: "market_snapshot", OutputRef: "", ErrorCode: code}}
	event := buildDomainAuditEvent(wf, result)
	if g.deps.Transactor != nil && g.deps.AuditRepo != nil {
		err := g.deps.Transactor.WithinTx(ctx, func(txCtx context.Context, repos repository.Repositories) error {
			return writeAuditEvent(txCtx, repos.AuditRepo, wf, result)
		})
		if err != nil {
			return MarketRefreshOutput{AuditEvents: []model.AuditEvent{event}}, err
		}
		return MarketRefreshOutput{AuditEvents: []model.AuditEvent{event}}, nil
	}
	if err := g.auditWriter.Write(ctx, wf, result); err != nil {
		return MarketRefreshOutput{AuditEvents: wf.AuditEvents}, err
	}
	return MarketRefreshOutput{AuditEvents: wf.AuditEvents}, nil
}
