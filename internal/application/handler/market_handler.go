package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/service"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// RefreshMarket 标准化市场状态并写入市场快照。部分失败时仍返回 200 和 failed_symbols。
func (a *App) RefreshMarket(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	var req dto.MarketRefreshRequest
	if err := decodeJSON(r, &req); err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	if req.AsOfDate != "" {
		if _, err := parseDateParam(req.AsOfDate); err != nil {
			WriteHandlerError(w, requestID, err)
			return
		}
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "as_of_date 暂不支持指定历史交易日刷新"))
		return
	}
	if len(req.Symbols) == 0 {
		req.Symbols = []string{"market"}
	}
	var ids []string
	var auditIDs []string
	var failures []dto.MarketRefreshFailure
	writeFailed := false
	for _, symbol := range req.Symbols {
		if symbol == "" {
			failures = append(failures, dto.MarketRefreshFailure{Symbol: symbol, Reason: "symbol 不能为空"})
			continue
		}
		out, err := workflow.NewMarketRefreshGraphWithDependencies(a.Deps).Run(r.Context(), workflow.MarketRefreshInput{RequestID: requestID, Symbol: symbol, PEPercentile: 50, PBPercentile: 50})
		if err != nil {
			if apperr.IsCode(err, apperr.CodeMarketSnapshotWriteFailed) {
				writeFailed = true
			}
			failures = append(failures, dto.MarketRefreshFailure{Symbol: symbol, Reason: err.Error()})
			continue
		}
		ids = append(ids, out.MarketSnapshot.MarketSnapshotID)
		for _, audit := range out.AuditEvents {
			auditIDs = append(auditIDs, audit.AuditEventID)
		}
	}
	if len(ids) == 0 {
		if writeFailed {
			WriteHandlerError(w, requestID, apperr.New(apperr.CodeMarketSnapshotWriteFailed, apperr.CategoryInternal, "市场快照写入失败"))
			return
		}
		if _, err := a.MarketSvc.AppendRefreshAudit(r.Context(), requestID, string(apperr.CodeDataSourceUnavailable), "", failures); err != nil {
			WriteHandlerError(w, requestID, err)
			return
		}
		if err := a.NotificationSvc.AppendNotification(r.Context(), repository.Notification{Type: "data_source_failure", Severity: "warning", Title: "市场数据源不可用", Message: "市场刷新未获取到可用行情数据", SourceType: "market_refresh", SourceID: "data_source_unavailable"}); err != nil {
			WriteHandlerError(w, requestID, err)
			return
		}
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "市场数据源不可用"))
		return
	}
	if len(failures) > 0 {
		auditID, err := a.MarketSvc.AppendRefreshAudit(r.Context(), requestID, "", string(mustJSON(failures)), failures)
		if err != nil {
			WriteHandlerError(w, requestID, err)
			return
		}
		if err := a.NotificationSvc.AppendNotification(r.Context(), repository.Notification{Type: "data_source_failure", Severity: "warning", Title: "部分市场数据源不可用", Message: "市场刷新存在部分标的失败", SourceType: "market_refresh", SourceID: "partial_data_source_failure"}); err != nil {
			WriteHandlerError(w, requestID, err)
			return
		}
		writeOK(w, requestID, dto.MarketRefreshResponse{RefreshedCount: len(ids), LatestSnapshotIDs: ids, FailedSymbols: failures, AuditEventIDs: []string{auditID}})
		return
	}
	writeOK(w, requestID, dto.MarketRefreshResponse{RefreshedCount: len(ids), LatestSnapshotIDs: ids, FailedSymbols: failures, AuditEventIDs: auditIDs})
}

func parseDateParam(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "日期参数必须是 YYYY-MM-DD")
	}
	return parsed, nil
}

func parseJSONMap(raw string) map[string]any {
	if raw == "" {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return map[string]any{}
	}
	return out
}

func normalizeDate(value string) string {
	if len(value) >= 10 {
		return value[:10]
	}
	return value
}

func marketDataStatus(market model.MarketSnapshot) string {
	if market.MarketSnapshotID == "" {
		return "missing"
	}
	if market.TradeDate == "" {
		return "unknown"
	}
	return "fresh"
}

func mustJSON(v any) []byte {
	buf, _ := json.Marshal(v)
	return buf
}

// GetMarketSourceHealth returns P34 source health derived from the latest normalized market facts.
func (a *App) GetMarketSourceHealth(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	symbol := r.URL.Query().Get("symbol")
	var market model.MarketSnapshot
	var err error
	if symbol != "" {
		market, err = a.QuerySvc.LatestMarketSnapshotBySymbol(r.Context(), symbol)
	} else {
		market, err = a.QuerySvc.LatestMarketSnapshot(r.Context())
	}
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	writeOK(w, requestID, dto.SourceHealthResponse{Sources: service.SourceHealthFromMarketSnapshot(market)})
}

// GetLatestMarketSnapshot 返回最近一条市场快照。
func (a *App) GetLatestMarketSnapshot(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	symbol := r.URL.Query().Get("symbol")
	var market model.MarketSnapshot
	var err error
	if symbol != "" {
		market, err = a.QuerySvc.LatestMarketSnapshotBySymbol(r.Context(), symbol)
	} else {
		market, err = a.QuerySvc.LatestMarketSnapshot(r.Context())
	}
	if err != nil {
		WriteHandlerError(w, requestID, err)
		return
	}
	out := dto.MarketSnapshotDTO{MarketSnapshotID: market.MarketSnapshotID, Symbol: market.Symbol, TradeDate: normalizeDate(market.TradeDate), DataStatus: marketDataStatus(market), ClosePrice: market.ClosePrice, TurnoverRate: market.TurnoverRate, PEPercentile: market.PEPercentile, PBPercentile: market.PBPercentile, VolumePercentile: market.VolumePercentile, VolatilityPercentile: market.VolatilityPercentile, LiquidityState: string(market.LiquidityState), SentimentState: string(market.SentimentState), MarketMetrics: parseJSONMap(market.MarketMetricsJSON)}
	writeOK(w, requestID, out)
}
