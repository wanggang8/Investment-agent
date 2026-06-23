package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	appknowledge "investment-agent/internal/application/knowledge"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

// EastmoneyFundEvidenceCollector turns read-only public fund facts into formal
// fund-profile evidence. It does not place orders or connect to brokers.
type EastmoneyFundEvidenceCollector struct {
	HTTPClient *http.Client
	BaseURL    string
}

func (c EastmoneyFundEvidenceCollector) FetchPublicEvidence(ctx context.Context, symbol string, start, end time.Time) ([]PublicEvidencePayload, error) {
	_ = start
	_ = end
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "基金代码不能为空")
	}
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = "https://fund.eastmoney.com"
	}
	raw, err := fetchReadonlyURL(ctx, c.HTTPClient, fmt.Sprintf("%s/pingzhongdata/%s.js", base, symbol))
	if err != nil {
		return nil, PublicEvidenceError{SourceName: "eastmoney_fund", ErrorCode: "source_unavailable", Count: 0, Err: err}
	}
	name := eastmoneyStringVar(raw, "fS_name")
	code := eastmoneyStringVar(raw, "fS_code")
	trend, err := eastmoneyNetWorthTrend(raw)
	if err != nil {
		return nil, PublicEvidenceError{SourceName: "eastmoney_fund", ErrorCode: "parse_error", Count: 0, Err: err}
	}
	if code == "" {
		code = symbol
	}
	capturedAt := time.Now().UTC()
	text := fmt.Sprintf("东方财富基金公开数据：基金代码 %s，基金名称 %s，最新单位净值 %.4f，净值日期 %s。该事实仅用于本地分析和人工复核。",
		code, valueOrText(name, symbol), trend.NAV, time.UnixMilli(trend.Timestamp).UTC().Format(time.DateOnly))
	return NormalizePublicEvidenceItems([]PublicEvidencePayload{{
		SourceName:     "eastmoney_fund",
		SourceLevel:    model.SourceLevelB,
		SourceType:     "fund_profile",
		EvidenceRole:   "formal",
		Symbol:         symbol,
		SourceRecordID: "eastmoney_fund_profile_" + symbol,
		Title:          fundProfileEvidenceTitle(symbol),
		Text:           text,
		URL:            fmt.Sprintf("%s/%s.html", base, symbol),
		PublishedAt:    capturedAt.Format(time.RFC3339),
		CapturedAt:     capturedAt,
		Raw:            map[string]any{"fund_name": name, "fund_code": code, "nav": trend.NAV, "nav_timestamp": trend.Timestamp},
	}})
}

// CsindexIndexEvidenceCollector maps a fund to its tracked index and records
// the official CSIndex index profile as formal public evidence.
type CsindexIndexEvidenceCollector struct {
	HTTPClient    *http.Client
	BaseURL       string
	SymbolToIndex map[string]string
}

type csindexIndexBasicResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		IndexFullNameCn  string  `json:"indexFullNameCn"`
		IndexShortNameCn string  `json:"indexShortNameCn"`
		IndexCode        string  `json:"indexCode"`
		BasicDate        string  `json:"basicDate"`
		BasicIndex       float64 `json:"basicIndex"`
		PublishDate      string  `json:"publishDate"`
		CurrencyCn       string  `json:"currencyCn"`
		IndexType        string  `json:"indexType"`
		IndexCnDesc      string  `json:"indexCnDesc"`
	} `json:"data"`
}

func (c CsindexIndexEvidenceCollector) FetchPublicEvidence(ctx context.Context, symbol string, start, end time.Time) ([]PublicEvidencePayload, error) {
	_ = start
	_ = end
	symbol = strings.TrimSpace(symbol)
	indexCode := strings.TrimSpace(c.SymbolToIndex[symbol])
	if indexCode == "" {
		indexCode = defaultTrackedIndex(symbol)
	}
	if indexCode == "" {
		return nil, PublicEvidenceError{SourceName: "csindex_index", ErrorCode: "no_data", Count: 0, Err: apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "未配置跟踪指数")}
	}
	base := strings.TrimRight(c.BaseURL, "/")
	if base == "" {
		base = "https://www.csindex.com.cn"
	}
	raw, err := fetchReadonlyURL(ctx, c.HTTPClient, fmt.Sprintf("%s/csindex-home/indexInfo/index-basic-info/%s", base, indexCode))
	if err != nil {
		return nil, PublicEvidenceError{SourceName: "csindex_index", ErrorCode: "source_unavailable", Count: 0, Err: err}
	}
	var body csindexIndexBasicResponse
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, PublicEvidenceError{SourceName: "csindex_index", ErrorCode: "parse_error", Count: 0, Err: apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "中证指数基础信息解析失败", err)}
	}
	if body.Code != "200" || strings.TrimSpace(body.Data.IndexCode) == "" {
		return nil, PublicEvidenceError{SourceName: "csindex_index", ErrorCode: "source_unavailable", Count: 0, Err: apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "中证指数基础信息不可用")}
	}
	capturedAt := time.Now().UTC()
	text := fmt.Sprintf("中证指数官方数据：基金 %s 跟踪指数 %s（%s），指数代码 %s，发布日期 %s，基日 %s，基点 %.2f。%s",
		symbol, body.Data.IndexFullNameCn, body.Data.IndexShortNameCn, body.Data.IndexCode, body.Data.PublishDate, body.Data.BasicDate, body.Data.BasicIndex, body.Data.IndexCnDesc)
	return NormalizePublicEvidenceItems([]PublicEvidencePayload{{
		SourceName:     "csindex_index",
		SourceLevel:    model.SourceLevelA,
		SourceType:     "fund_profile",
		EvidenceRole:   "formal",
		Symbol:         symbol,
		SourceRecordID: "csindex_index_profile_" + body.Data.IndexCode,
		Title:          fundProfileEvidenceTitle(symbol),
		Text:           text,
		URL:            fmt.Sprintf("%s/#/indices/family/detail?indexCode=%s", base, body.Data.IndexCode),
		PublishedAt:    capturedAt.Format(time.RFC3339),
		CapturedAt:     capturedAt,
		Raw:            map[string]any{"tracked_index": body.Data.IndexCode, "index_name": body.Data.IndexFullNameCn, "publish_date": body.Data.PublishDate},
	}})
}

func fundProfileEvidenceTitle(symbol string) string {
	return strings.TrimSpace(symbol) + " 产品与跟踪指数事实"
}

func defaultTrackedIndex(symbol string) string {
	return appknowledge.DefaultTrackedIndexSymbol(symbol)
}

func valueOrText(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}
