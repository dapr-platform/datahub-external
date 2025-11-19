package controllers

import (
	"net/http"

	"github.com/go-chi/render"
)

// APIResponse 统一响应格式
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Status: 200,
		Data:   data,
	}
}

// SuccessResponseWithMessage 带消息的成功响应
func SuccessResponseWithMessage(data interface{}, message string) *APIResponse {
	return &APIResponse{
		Status:  200,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(status int, message string) *APIResponse {
	return &APIResponse{
		Status:  status,
		Message: message,
		Error:   http.StatusText(status),
	}
}

// RespondSuccess 发送成功响应
func RespondSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, SuccessResponse(data))
}

// RespondError 发送错误响应
func RespondError(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, ErrorResponse(status, message))
}



















