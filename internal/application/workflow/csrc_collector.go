package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/pkg/apperr"
)

// CsrcCollector 实现证监会监管信息采集（P26 首批）。
type CsrcCollector struct {
	HTTPClient *http.Client
	BaseURL    string
}

type csrcResponse struct {
	Data csrcData `json:"data"`
}

type csrcData struct {
	Page      int          `json:"page"`
	Rows      int          `json:"rows"`
	ChannelID string       `json:"channelId"`
	Total     int          `json:"total"`
	Results   []csrcResult `json:"results"`
}

type csrcResult struct {
	Title            string   `json:"title"`
	Content          string   `json:"content"`
	ContentHTML      string   `json:"contentHtml"`
	Memo             string   `json:"memo"`
	URL              string   `json:"url"`
	PublishedTime    string   `json:"publishedTime"`
	PublishedTimeStr string   `json:"publishedTimeStr"`
	ChannelName      string   `json:"channelName"`
	ChannelCodeName  string   `json:"channelCodeName"`
	ManuscriptID     string   `json:"manuscriptId"`
	ResList          []string `json:"resList"`
	DomainMetaList   []any    `json:"domainMetaList"`
}

func (c *CsrcCollector) FetchPublicEvidence(ctx context.Context, symbol string, start, end time.Time) ([]PublicEvidencePayload, error) {
	if c.BaseURL == "" {
		c.BaseURL = "https://www.csrc.gov.cn"
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 30 * time.Second}
	}
	start, end = evidenceDateRange(start, end)

	payloads := []PublicEvidencePayload{}
	pages := 1
	for page := 1; page <= pages; page++ {
		apiResp, err := c.fetchPage(ctx, symbol, page)
		if err != nil {
			return nil, err
		}
		pages = totalPages(apiResp.Data.Total, publicEvidencePageSize)
		payloads = append(payloads, c.mapResults(symbol, start, end, apiResp.Data.Results)...)
	}
	if len(payloads) == 0 {
		return nil, PublicEvidenceError{SourceName: "csrc", ErrorCode: "no_data", Count: 0, Err: apperr.New(apperr.CodeNotFound, apperr.CategoryNotFound, "未找到监管信息")}
	}
	return NormalizePublicEvidenceItems(payloads)
}

func (c *CsrcCollector) fetchPage(ctx context.Context, symbol string, page int) (csrcResponse, error) {
	params := url.Values{}
	params.Set("_isAgg", "true")
	params.Set("_isJson", "true")
	params.Set("keyword", symbol)
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(publicEvidencePageSize))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/searchList?"+params.Encode(), nil)
	if err != nil {
		return csrcResponse{}, apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "构建请求失败", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; investment-agent/1.0)")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return csrcResponse{}, PublicEvidenceError{SourceName: "csrc", ErrorCode: "source_unavailable", Count: 0, Err: apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "请求失败", err)}
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return csrcResponse{}, PublicEvidenceError{SourceName: "csrc", ErrorCode: "source_unavailable", Count: 0, Err: apperr.New(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)))}
	}
	var apiResp csrcResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return csrcResponse{}, PublicEvidenceError{SourceName: "csrc", ErrorCode: "parse_error", Count: 0, Err: apperr.Wrap(apperr.CodeDataSourceUnavailable, apperr.CategoryInternal, "解析响应失败", err)}
	}
	return apiResp, nil
}

func (c *CsrcCollector) mapResults(symbol string, start, end time.Time, results []csrcResult) []PublicEvidencePayload {
	capturedAt := time.Now().UTC()
	payloads := []PublicEvidencePayload{}
	for _, result := range results {
		if !evidenceInRange(result.PublishedTime, start, end) {
			continue
		}
		text := result.Content
		if text == "" {
			text = result.Memo
		}
		if text == "" {
			text = result.ContentHTML
		}
		evidenceRole := "background"
		channelName := result.ChannelName
		if channelName == "" {
			channelName = result.ChannelCodeName
		}
		if result.ChannelCodeName == "rules" || result.ChannelCodeName == "punishment" || result.ChannelCodeName == "market_ban" {
			evidenceRole = "formal"
		}
		payloads = append(payloads, PublicEvidencePayload{
			SourceName:     "csrc",
			SourceLevel:    model.SourceLevelA,
			SourceType:     "regulatory_disclosure",
			EvidenceRole:   evidenceRole,
			Symbol:         symbol,
			SourceRecordID: result.ManuscriptID,
			Title:          result.Title,
			Text:           text,
			URL:            resolveEvidenceURL(c.BaseURL, result.URL),
			PublishedAt:    result.PublishedTime,
			CapturedAt:     capturedAt,
			Raw:            map[string]any{"channel_name": channelName, "channel_code_name": result.ChannelCodeName, "published_time_str": result.PublishedTimeStr, "res_list": result.ResList},
		})
	}
	return payloads
}
