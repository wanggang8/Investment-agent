package httputil

import (
	"net/http"
	"time"
)

// NewHTTPClient 创建带超时控制的 HTTP 客户端，供后续外部数据源和模型调用复用。
func NewHTTPClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &http.Client{Timeout: timeout}
}
