package httputil

import (
	"encoding/json"
	"net/http"
)

// Envelope 是业务 API 的统一响应信封。
// request_id 用于串联前端请求、后端日志和审计事件。
type Envelope struct {
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	Error     *APIError   `json:"error,omitempty"`
}

// APIError 是统一错误体，code 与 API 契约中的错误码保持一致。
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// WriteJSON 写入 JSON 响应，并统一设置 UTF-8 Content-Type。
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteSuccess 写入成功信封。
func WriteSuccess(w http.ResponseWriter, requestID string, data interface{}) {
	WriteJSON(w, http.StatusOK, Envelope{
		RequestID: requestID,
		Data:      data,
	})
}

// WriteError 写入失败信封，保留 request_id 便于问题追踪。
func WriteError(w http.ResponseWriter, status int, requestID, code, message, detail string) {
	WriteJSON(w, status, Envelope{
		RequestID: requestID,
		Error: &APIError{
			Code:    code,
			Message: message,
			Detail:  detail,
		},
	})
}
