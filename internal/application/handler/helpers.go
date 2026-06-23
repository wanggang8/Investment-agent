package handler

import (
	"encoding/json"
	"net/http"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/pkg/httputil"
)

func nowRFC3339() string { return clock.SystemClock{}.NowRFC3339() }

func nullStringLocal(v string) any {
	if v == "" {
		return nil
	}
	return v
}

func decodeJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return nil
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return apperr.Wrap(apperr.CodeBadRequest, apperr.CategoryBadRequest, "请求 JSON 格式错误", err)
	}
	return nil
}

func decodeOptionalJSON(r *http.Request, v any) error {
	if r.Body == nil || r.ContentLength == 0 {
		return nil
	}
	return decodeJSON(r, v)
}

func writeOK(w http.ResponseWriter, requestID string, data any) {
	httputil.WriteSuccess(w, requestID, data)
}

func portfolioDTO(s repository.PortfolioSnapshot) dto.PortfolioSnapshotDTO {
	return dto.PortfolioSnapshotDTO{SnapshotID: s.SnapshotID, SnapshotTime: s.SnapshotTime, Cash: s.Cash, TotalAssets: s.TotalAssets, CashRatio: s.CashRatio, HighRiskRatio: s.HighRiskRatio, PositionCount: s.PositionCount}
}

func positionDTO(p repository.Position) dto.PositionDTO {
	return dto.PositionDTO{PositionID: p.PositionID, Symbol: p.Symbol, Name: p.Name, Quantity: p.Quantity, CostPrice: p.CostPrice, CurrentPrice: p.CurrentPrice, MarketValue: p.MarketValue, UnrealizedProfitRatio: p.UnrealizedProfitRatio, PositionState: p.PositionState, BuyDate: p.BuyDate, BuyReason: p.BuyReason, AssetTag: p.AssetTag}
}

func evidenceDTO(e repository.EvidenceRef) dto.EvidenceDTO {
	return dto.EvidenceDTO{EvidenceID: e.EvidenceID, SourceName: e.SourceName, SourceLevel: e.SourceLevel, EvidenceRole: e.EvidenceRole, PublishedAt: e.PublishedAt, CapturedAt: e.CapturedAt, OriginalURL: e.OriginalURL, Summary: e.Summary, ContentHash: e.ContentHash, TimeWeight: e.TimeWeight, RelevanceScore: e.RelevanceScore, IndependentSourceCount: e.IndependentSourceCount, HighGradeIndependentSourceCount: e.HighGradeIndependentSourceCount}
}

func splitJSONStrings(raw string) []string {
	if raw == "" {
		return nil
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil
	}
	return out
}

func parseJSONAny(raw string) any {
	if raw == "" {
		return nil
	}
	var out any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return raw
	}
	return out
}
