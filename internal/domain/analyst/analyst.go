package analyst

import "context"

// Request 是分析服务的有界输入。
type Request struct {
	AgentName               string
	Symbol                  string
	EvidenceSummary         string
	PositionContext         string
	RuleBoundary            string
	KnowledgeContextSummary string
}

// Response 只承载分析材料，不承载最终裁决。
type Response struct {
	Reports  map[string]string
	Metadata map[string]string
}

// Service 封装外部或本地分析材料生成。
type Service interface {
	Analyze(ctx context.Context, req Request) (Response, error)
}
