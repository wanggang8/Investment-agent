package dto

// Meta 是业务 API 响应的通用元信息。
// generated_at 表示本次响应生成时间，rule_version 表示裁决使用的规则版本。
type Meta struct {
	GeneratedAt string `json:"generated_at,omitempty"`
	RuleVersion string `json:"rule_version,omitempty"`
}

// PageResult 是列表 API 的通用分页片段，P4 暂按本地单页查询返回。
type PageResult[T any] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
}
