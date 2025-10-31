package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/render"
)

// ContextKey API Key上下文键类型
type ContextKey string

const (
	// ApiKeyContextKey API Key在上下文中的键
	ApiKeyContextKey ContextKey = "api_key"
)

// ApiKeyAuthMiddleware API Key认证中间件
type ApiKeyAuthMiddleware struct {
	validApiKey    string
	whitelistPaths []string
}

// NewApiKeyAuthMiddleware 创建API Key认证中间件实例
func NewApiKeyAuthMiddleware() *ApiKeyAuthMiddleware {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "default-api-key-please-change-in-production"
	}

	return &ApiKeyAuthMiddleware{
		validApiKey: apiKey,
		whitelistPaths: []string{
			"/health",
			"/ready",
			"/swagger",
			"/metrics",
		},
	}
}

// Middleware 认证中间件处理函数
func (m *ApiKeyAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 检查是否在白名单中
		if m.isWhitelistPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// 方式1: 从 X-API-Key 头提取
		apiKey := r.Header.Get("X-API-Key")

		// 方式2: 从 Authorization Bearer 头提取
		if apiKey == "" {
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 验证API Key
		if apiKey == "" {
			m.respondUnauthorized(w, r, "缺少API Key")
			return
		}

		if apiKey != m.validApiKey {
			m.respondUnauthorized(w, r, "无效的API Key")
			return
		}

		// 将API Key注入上下文
		ctx := context.WithValue(r.Context(), ApiKeyContextKey, apiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// isWhitelistPath 检查路径是否在白名单中
func (m *ApiKeyAuthMiddleware) isWhitelistPath(path string) bool {
	for _, whitelistPath := range m.whitelistPaths {
		if strings.HasPrefix(path, whitelistPath) {
			return true
		}
	}
	return false
}

// respondUnauthorized 返回401未授权响应
func (m *ApiKeyAuthMiddleware) respondUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	render.JSON(w, r, map[string]interface{}{
		"status":  http.StatusUnauthorized,
		"message": message,
		"error":   "Unauthorized",
	})
}



