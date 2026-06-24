package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"investment-agent/internal/domain/analyst"
	"investment-agent/internal/pkg/apperr"
)

// Client 封装 DeepSeek Chat API；调用结果只转换为分析材料。
type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

type Config struct {
	APIKey         string
	BaseURL        string
	Model          string
	TimeoutSeconds int
}

const promptVersion = "p37-analyst-v1"
const userAgent = "investment-agent/llm-openai-compatible"

type classifiedError struct {
	category string
	metadata map[string]string
	err      error
}

func (e classifiedError) Error() string {
	if e.err == nil {
		return e.category
	}
	return e.err.Error()
}

func (e classifiedError) Unwrap() error {
	return e.err
}

func (e classifiedError) Category() string {
	return e.category
}

func (e classifiedError) Metadata() map[string]string {
	out := map[string]string{}
	for k, v := range e.metadata {
		out[k] = v
	}
	return out
}

func classifyWithMetadata(category string, err error, metadata map[string]string) error {
	return classifiedError{category: category, metadata: metadata, err: err}
}

func ErrorCategory(err error) string {
	var classified classifiedError
	if errors.As(err, &classified) {
		return classified.category
	}
	return "unavailable"
}

func ErrorMetadata(err error) map[string]string {
	var classified classifiedError
	if errors.As(err, &classified) {
		return classified.Metadata()
	}
	return map[string]string{}
}

type QualityResult struct {
	Status  string
	Reasons []string
}

func NewClient(cfg Config, httpClient *http.Client) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.deepseek.com"
	}
	if cfg.Model == "" {
		cfg.Model = "deepseek-chat"
	}
	if httpClient == nil {
		timeout := 15 * time.Second
		if cfg.TimeoutSeconds > 0 {
			timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
		}
		httpClient = &http.Client{Timeout: timeout}
	}
	return &Client{apiKey: cfg.APIKey, baseURL: strings.TrimRight(cfg.BaseURL, "/"), model: cfg.Model, httpClient: httpClient}
}

func (c *Client) Analyze(ctx context.Context, req analyst.Request) (analyst.Response, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return analyst.Response{}, classifyWithMetadata("missing_key", apperr.New(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "DeepSeek API Key 未配置"), c.callMetadata(req, "not_started", "not_evaluated", ""))
	}
	resp, err := c.analyzeOnce(ctx, req, false)
	if err == nil {
		return resp, nil
	}
	if ErrorCategory(err) == "timeout" {
		retryResp, retryErr := c.analyzeOnce(ctx, req, false)
		if retryErr != nil {
			return analyst.Response{}, retryErr
		}
		if retryResp.Metadata == nil {
			retryResp.Metadata = map[string]string{}
		}
		retryResp.Metadata["retry"] = "timeout_retry"
		return retryResp, nil
	}
	if ErrorCategory(err) != "quality_failed" {
		return analyst.Response{}, err
	}
	retryResp, retryErr := c.analyzeOnce(ctx, req, true)
	if retryErr != nil {
		return analyst.Response{}, retryErr
	}
	if retryResp.Metadata == nil {
		retryResp.Metadata = map[string]string{}
	}
	retryResp.Metadata["retry"] = "quality_failed_safety_reprompt"
	return retryResp, nil
}

func (c *Client) analyzeOnce(ctx context.Context, req analyst.Request, safetyRetry bool) (analyst.Response, error) {
	body := map[string]any{"model": c.model, "messages": []map[string]string{{"role": "system", "content": systemPrompt(safetyRetry)}, {"role": "user", "content": buildPromptWithSafety(req, safetyRetry)}}}
	buf, _ := json.Marshal(body)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(buf))
	if err != nil {
		return analyst.Response{}, classifyWithMetadata("unavailable", apperr.Wrap(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "构造 DeepSeek 请求失败", err), c.callMetadata(req, "not_started", "not_evaluated", ""))
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", userAgent)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		category := "unavailable"
		var netErr net.Error
		if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &netErr) && netErr.Timeout()) || strings.Contains(strings.ToLower(err.Error()), "timeout") || strings.Contains(strings.ToLower(err.Error()), "deadline exceeded") {
			category = "timeout"
		}
		return analyst.Response{}, classifyWithMetadata(category, apperr.Wrap(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "DeepSeek 调用失败", err), c.callMetadata(req, category, "not_evaluated", ""))
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return analyst.Response{}, classifyWithMetadata("http_error", apperr.New(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "DeepSeek 返回不可用状态"), c.callMetadata(req, "http_error", "not_evaluated", ""))
	}
	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return analyst.Response{}, classifyWithMetadata("parse_error", apperr.Wrap(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "DeepSeek 输出不可解析", err), c.callMetadata(req, "parse_error", "not_evaluated", ""))
	}
	if len(decoded.Choices) == 0 || strings.TrimSpace(decoded.Choices[0].Message.Content) == "" {
		return analyst.Response{}, classifyWithMetadata("empty_response", apperr.New(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "DeepSeek 输出为空"), c.callMetadata(req, "parsed", "empty_response", ""))
	}
	content := decoded.Choices[0].Message.Content
	quality := EvaluateQuality(content)
	if quality.Status != "passed" {
		return analyst.Response{}, classifyWithMetadata("quality_failed", apperr.New(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "DeepSeek 输出质量检查失败"), c.callMetadata(req, "parsed", "failed", content))
	}
	return analyst.Response{Reports: map[string]string{req.AgentName: content}, Metadata: c.callMetadata(req, "parsed", quality.Status, content)}, nil
}

func buildPrompt(req analyst.Request) string {
	return buildPromptWithSafety(req, false)
}

func buildPromptWithSafety(req analyst.Request, safetyRetry bool) string {
	// prompt 只包含允许的证据、持仓上下文和规则边界，避免把 LLM 输出误用成裁决。
	knowledge := strings.TrimSpace(req.KnowledgeContextSummary)
	if knowledge == "" {
		knowledge = "未提供；不得自行假设内置知识或数据准备度已满足。"
	}
	prompt := fmt.Sprintf("标的：%s\n证据：%s\n持仓上下文：%s\n知识与数据准备度：%s\n规则边界：%s\n请给出%s分析材料。", req.Symbol, req.EvidenceSummary, req.PositionContext, knowledge, req.RuleBoundary, req.AgentName)
	if safetyRetry {
		prompt += "\n安全重试：上一次输出未通过本地质量闸。只允许描述估值、风险、证据缺口和人工复核问题；不得使用“建议买入”“建议卖出”“立即买入”“立即卖出”“最终裁决”“保证收益”“确定上涨”等交易指令、最终裁决或收益承诺措辞。"
	}
	return prompt
}

func systemPrompt(safetyRetry bool) string {
	base := "你是投资分析员，只输出分析材料，不输出最终裁决或交易动作。"
	if safetyRetry {
		return base + " 严格避免任何交易指令、最终裁决、确定性预测或收益承诺；用“风险/证据/人工复核问题”表述替代买卖措辞。"
	}
	return base
}

func summarizeInput(req analyst.Request) string {
	return summarizeText(strings.Join([]string{req.AgentName, req.Symbol, req.EvidenceSummary, req.PositionContext, req.KnowledgeContextSummary, req.RuleBoundary}, " "))
}

func (c *Client) callMetadata(req analyst.Request, parseStatus string, qualityStatus string, output string) map[string]string {
	return map[string]string{"model": c.model, "prompt_version": promptVersion, "input_summary": summarizeInput(req), "output_summary": summarizeText(output), "parse_status": parseStatus, "quality_status": qualityStatus}
}

func summarizeText(value string) string {
	value = strings.TrimSpace(strings.Join(strings.Fields(value), " "))
	if len([]rune(value)) <= 80 {
		return value
	}
	return string([]rune(value)[:80])
}

func EvaluateQuality(output string) QualityResult {
	normalized := strings.ToLower(strings.TrimSpace(output))
	reasons := []string{}
	checks := []struct {
		reason   string
		patterns []string
	}{
		{reason: "return_promise", patterns: []string{"保证收益", "稳赚", "收益翻倍", "必赚", "稳定获利", "稳定盈利", "锁定收益", "收益确定"}},
		{reason: "deterministic_prediction", patterns: []string{"必涨", "必跌", "一定上涨", "一定下跌", "一定会涨", "一定会跌", "确定会上涨", "确定会下跌", "必然上涨", "必然下跌", "明天涨", "明天跌"}},
		{reason: "direct_trade_instruction", patterns: []string{"立即买入", "马上买入", "全仓买入", "建议买入", "买入510300", "买入该标的", "立即卖出", "马上卖出", "建议卖出", "卖出该标的", "立即下单", "马上下单", "执行下单", "直接下单", "立即建仓", "马上建仓", "全仓建仓", "建议建仓"}},
		{reason: "final_verdict_override", patterns: []string{"最终裁决：", "最终裁决:", "最终裁决为", "最终裁决是", "最终结论：买入", "最终结论:买入", "final verdict:"}},
	}
	for _, check := range checks {
		for _, pattern := range check.patterns {
			if containsUnsafePattern(normalized, strings.ToLower(pattern)) {
				reasons = append(reasons, check.reason)
				break
			}
		}
	}
	if len(reasons) > 0 {
		return QualityResult{Status: "failed", Reasons: reasons}
	}
	return QualityResult{Status: "passed"}
}

func containsUnsafePattern(text string, pattern string) bool {
	runes := []rune(text)
	target := []rune(pattern)
	if len(target) == 0 || len(runes) < len(target) {
		return false
	}
	for i := 0; i <= len(runes)-len(target); i++ {
		if string(runes[i:i+len(target)]) != pattern {
			continue
		}
		prefixStart := i - 6
		if prefixStart < 0 {
			prefixStart = 0
		}
		prefix := string(runes[prefixStart:i])
		if strings.Contains(prefix, "不") || strings.Contains(prefix, "不得") || strings.Contains(prefix, "禁止") || strings.Contains(prefix, "避免") || strings.Contains(prefix, "不要") || strings.Contains(prefix, "不能") || strings.Contains(prefix, "无") {
			continue
		}
		return true
	}
	return false
}
