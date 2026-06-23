package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

type P89StructuredPublicCollector struct {
	EastmoneyPush2BaseURL      string
	EastmoneyH5BaseURL         string
	EastmoneyDatacenterBaseURL string
	SSEQueryBaseURL            string
	HTTPClient                 *http.Client
}

func (c P89StructuredPublicCollector) FetchMarketData(ctx context.Context, symbol string) (MarketDataPoint, error) {
	symbol = strings.TrimSpace(symbol)
	metadata := map[string]any{
		"p89_provider_policy": "no-login/no-paid/no-authorization/no-Level2/no-high-frequency public read-only sources",
	}
	health := map[string]any{}
	var errs []error

	capital, err := c.fetchCapitalFlow(ctx, symbol)
	if err != nil {
		errs = append(errs, err)
		health["capital_flow"] = p34SourceHealthMetadata("eastmoney_push2", string(model.SourceLevelB), "structured_public_fields", "capital_flow", "failed", "", workflowNowRFC3339(), symbol, "")
	} else {
		metadata["capital_flow"] = capital
		sourceName := "eastmoney_push2"
		if strings.Contains(stringFromAny(capital["source"]), "h5") {
			sourceName = "eastmoney_h5_zjlx_public"
		}
		health["capital_flow"] = p34SourceHealthMetadata(sourceName, string(model.SourceLevelB), "structured_public_fields", "capital_flow", "fresh", workflowNowRFC3339(), "", symbol, stringFromAny(capital["date"]))
	}

	margin, err := c.fetchMarginFinancing(ctx)
	if err != nil {
		errs = append(errs, err)
		health["margin_financing"] = p34SourceHealthMetadata("sse_query", string(model.SourceLevelB), "structured_public_fields", "margin_financing", "failed", "", workflowNowRFC3339(), symbol, "")
	} else {
		metadata["margin_financing"] = margin
		health["margin_financing"] = p34SourceHealthMetadata("sse_query", string(model.SourceLevelB), "structured_public_fields", "margin_financing", "fresh", workflowNowRFC3339(), "", symbol, stringFromAny(margin["date"]))
	}

	financial, err := c.fetchConstituentFinancial(ctx, symbol)
	if err != nil {
		errs = append(errs, err)
		health["constituent_financials"] = p34SourceHealthMetadata("eastmoney_datacenter", string(model.SourceLevelB), "structured_public_fields", "constituent_financials", "failed", "", workflowNowRFC3339(), symbol, "")
	} else {
		metadata["constituent_financial"] = financial
		health["constituent_financials"] = p34SourceHealthMetadata("eastmoney_datacenter", string(model.SourceLevelB), "structured_public_fields", "constituent_financials", "fresh", workflowNowRFC3339(), "", symbol, stringFromAny(financial["disclosure_date"]))
	}

	if P88NormalizeStructuredDataMetadata(metadata).Empty() {
		if err := errors.Join(errs...); err != nil {
			return MarketDataPoint{}, err
		}
		return MarketDataPoint{}, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 structured provider returned no eligible fields")
	}
	metadata["p34_source_health"] = health
	tradeDate := firstNonEmpty(stringFromAny(metadataValue(metadata, "capital_flow", "date")), stringFromAny(metadataValue(metadata, "margin_financing", "date")), stringFromAny(metadataValue(metadata, "constituent_financial", "disclosure_date")))
	return MarketDataPoint{SourceName: "p89_public_structured_sources", SourceLevel: model.SourceLevelB, SourceType: "structured_public_fields", TradeDate: tradeDate, CapturedAt: workflowNowRFC3339(), ContentHash: "sha256:" + stableHash("p89_public_structured_sources", symbol, tradeDate), Metadata: metadata}, nil
}

func metadataValue(metadata map[string]any, itemKey string, fieldKey string) any {
	item, ok := metadata[itemKey].(map[string]any)
	if !ok {
		return nil
	}
	return item[fieldKey]
}

func (c P89StructuredPublicCollector) fetchCapitalFlow(ctx context.Context, symbol string) (map[string]any, error) {
	if item, err := c.fetchH5CapitalFlow(ctx, symbol); err == nil {
		return item, nil
	}
	return c.fetchPush2CapitalFlow(ctx, symbol)
}

func (c P89StructuredPublicCollector) fetchH5CapitalFlow(ctx context.Context, symbol string) (map[string]any, error) {
	base := strings.TrimRight(c.EastmoneyH5BaseURL, "/")
	if base == "" {
		base = "https://emdatah5.eastmoney.com"
	}
	secid := "0." + symbol
	if strings.HasPrefix(symbol, "6") {
		secid = "1." + symbol
	}
	q := url.Values{}
	q.Set("secid", secid)
	q.Set("fields1", "f1,f2,f3")
	q.Set("fields2", "f51,f52,f53,f54,f55,f56,f62,f63")
	q.Set("ut", "b2884a393a59ad64002292a3e90d46a5")
	body, err := c.fetch(ctx, base+"/dc/ZJLX/getDBHistoryData?"+q.Encode(), map[string]string{"Referer": base + "/dc/zjlx/stock?fc=" + url.QueryEscape(secid)})
	if err != nil {
		return nil, err
	}
	var payload struct {
		Data struct {
			KLines []string `json:"klines"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || len(payload.Data.KLines) == 0 {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P90 H5 资金流向解析失败", err)
	}
	fields := strings.Split(payload.Data.KLines[len(payload.Data.KLines)-1], ",")
	if len(fields) < 2 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P90 H5 资金流向字段不足")
	}
	rawNetFlow, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P90 H5 资金净流向解析失败", err)
	}
	inflow := rawNetFlow
	outflow := 0.0
	if rawNetFlow < 0 {
		inflow = 0
		outflow = math.Abs(rawNetFlow)
	}
	return map[string]any{"date": fields[0], "net_inflow": inflow, "net_outflow": outflow, "raw_net_flow": rawNetFlow, "source": "eastmoney_h5_zjlx_public"}, nil
}

func (c P89StructuredPublicCollector) fetchPush2CapitalFlow(ctx context.Context, symbol string) (map[string]any, error) {
	base := strings.TrimRight(c.EastmoneyPush2BaseURL, "/")
	if base == "" {
		base = "https://push2.eastmoney.com"
	}
	secid := "0." + symbol
	if strings.HasPrefix(symbol, "6") {
		secid = "1." + symbol
	}
	q := url.Values{}
	q.Set("lmt", "1")
	q.Set("klt", "101")
	q.Set("fields1", "f1,f2,f3,f7")
	q.Set("fields2", "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f62,f63")
	q.Set("secid", secid)
	body, err := c.fetch(ctx, base+"/api/qt/stock/fflow/daykline/get?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Data struct {
			KLines []string `json:"klines"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || len(payload.Data.KLines) == 0 {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 资金流向解析失败", err)
	}
	fields := strings.Split(payload.Data.KLines[0], ",")
	if len(fields) < 4 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 资金流向字段不足")
	}
	inflow, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 资金净流入解析失败", err)
	}
	outflow, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 资金净流出解析失败", err)
	}
	return map[string]any{"date": fields[0], "net_inflow": inflow, "net_outflow": math.Abs(outflow), "source": "eastmoney_push2_public"}, nil
}

func (c P89StructuredPublicCollector) fetchMarginFinancing(ctx context.Context) (map[string]any, error) {
	base := strings.TrimRight(c.SSEQueryBaseURL, "/")
	if base == "" {
		base = "https://query.sse.com.cn"
	}
	q := url.Values{}
	q.Set("isPagination", "true")
	q.Set("sqlId", "COMMON_SSE_SJ_GPSJ_GPHYSJ_MX_L")
	q.Set("pageHelp.pageSize", "2")
	q.Set("pageHelp.pageNo", "1")
	body, err := c.fetch(ctx, base+"/marketdata/tradedata/queryMargin.do?"+q.Encode(), map[string]string{"Referer": "https://www.sse.com.cn/market/othersdata/margin/detail/"})
	if err != nil {
		return nil, err
	}
	var payload struct {
		PageHelp struct {
			Data []struct {
				OpDate string  `json:"opDate"`
				RZYE   float64 `json:"rzye"`
			} `json:"data"`
		} `json:"pageHelp"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || len(payload.PageHelp.Data) < 1 {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 融资融券解析失败", err)
	}
	latest := payload.PageHelp.Data[0]
	changeRate := 0.0
	if len(payload.PageHelp.Data) > 1 && payload.PageHelp.Data[1].RZYE != 0 {
		changeRate = (latest.RZYE - payload.PageHelp.Data[1].RZYE) / payload.PageHelp.Data[1].RZYE
	}
	return map[string]any{"date": yyyymmdd(latest.OpDate), "margin_balance": latest.RZYE, "balance_change_rate": changeRate, "source": "sse_query_margin_public"}, nil
}

func (c P89StructuredPublicCollector) fetchConstituentFinancial(ctx context.Context, symbol string) (map[string]any, error) {
	base := strings.TrimRight(c.EastmoneyDatacenterBaseURL, "/")
	if base == "" {
		base = "https://datacenter-web.eastmoney.com"
	}
	q := url.Values{}
	q.Set("sortColumns", "NOTICE_DATE")
	q.Set("sortTypes", "-1")
	q.Set("pageSize", "1")
	q.Set("pageNumber", "1")
	q.Set("reportName", "RPT_LICO_FN_CPD")
	q.Set("columns", "ALL")
	q.Set("filter", `(SECURITY_CODE="`+symbol+`")`)
	body, err := c.fetch(ctx, base+"/api/data/v1/get?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}
	var payload struct {
		Result struct {
			Data []map[string]any `json:"data"`
		} `json:"result"`
	}
	if err := json.Unmarshal(body, &payload); err != nil || len(payload.Result.Data) == 0 {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 成分财务解析失败", err)
	}
	item := payload.Result.Data[0]
	return map[string]any{
		"revenue":         floatValue(item, "TOTAL_OPERATE_INCOME"),
		"net_profit":      floatValue(item, "PARENT_NETPROFIT"),
		"growth":          floatValue(item, "SJLTZ"),
		"disclosure_date": strings.TrimSuffix(stringValue(item, "NOTICE_DATE"), " 00:00:00"),
		"source":          "eastmoney_datacenter_public",
	}, nil
}

func (c P89StructuredPublicCollector) fetch(ctx context.Context, endpoint string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 structured provider request build failed", err)
	}
	req.Header.Set("User-Agent", "InvestmentAgent-P89/1.0 readonly acceptance")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 structured provider unavailable", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "P89 structured provider returned non-2xx")
	}
	return io.ReadAll(resp.Body)
}

func yyyymmdd(value string) string {
	if len(value) == 8 {
		return value[:4] + "-" + value[4:6] + "-" + value[6:8]
	}
	return value
}
