package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/pkg/httputil"
)

const internalErrorMessage = "内部错误，请稍后重试。"

// WriteHandlerError 是 P4 handler 的统一错误出口。
// 已知 AppError 使用稳定错误码和状态码；未知错误隐藏底层细节，避免把 SQL、文件路径或外部服务原文暴露给前端。
func WriteHandlerError(w http.ResponseWriter, requestID string, err error) {
	if err == nil {
		return
	}
	if appErr, ok := apperr.AsAppError(err); ok {
		message := appErr.Message
		if message == "" {
			message = string(appErr.Code)
		}
		httputil.WriteError(w, appErr.HTTPStatus, requestID, string(appErr.Code), message, "")
		return
	}
	httputil.WriteError(w, http.StatusInternalServerError, requestID, string(apperr.CodeInternalError), internalErrorMessage, "")
}

// RequestID 从请求头读取 request_id；缺失时生成当前请求内稳定的后备 ID。
// 业务 handler 统一使用该值写响应信封，方便前端、日志和审计串联。
func RequestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err == nil {
		return "req_" + hex.EncodeToString(buf[:])
	}
	return "req_" + clock.SystemClock{}.Now().Format("20060102150405.000000000")
}
