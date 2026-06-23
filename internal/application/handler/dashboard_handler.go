package handler

import (
	"net/http"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/pkg/apperr"
)

// GetDashboardToday 聚合最新账户、市场和决策摘要，供驾驶舱首屏渲染。
func (a *App) GetDashboardToday(w http.ResponseWriter, r *http.Request) {
	requestID := RequestID(r)
	snapshot, err := a.QuerySvc.LatestPortfolioSnapshot(r.Context())
	if err != nil {
		WriteHandlerError(w, requestID, apperr.New(apperr.CodeDataRequired, apperr.CategoryConflict, "需要先录入账户和持仓"))
		return
	}
	data := dto.DashboardTodayResponse{
		DashboardState:   "normal",
		DisciplineStatus: "观察",
		DataUpdatedAt:    snapshot.SnapshotTime,
		PortfolioSummary: dto.PortfolioSummary{TotalAssets: snapshot.TotalAssets, CashRatio: snapshot.CashRatio, HighRiskRatio: snapshot.HighRiskRatio, PositionCount: snapshot.PositionCount},
		MarketSummary:    dto.MarketSummary{SentimentState: "neutral", LiquidityState: "normal"},
		TriggeredRules:   []dto.TriggeredRuleDTO{},
		DecisionSummary:  dto.DecisionSummary{Verdict: "暂无正式建议", FinalVerdictStatus: "hold", ProhibitedActions: []string{}, OptionalActions: []string{"查看证据"}, ActionRequired: false, ConfirmationStatus: "not_required"},
	}
	writeOK(w, requestID, data)
}
