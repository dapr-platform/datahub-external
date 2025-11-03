package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/render"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health 基础健康检查
// @Summary 健康检查
// @Description 检查服务是否运行
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, map[string]interface{}{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// Ready 就绪检查
// @Summary 就绪检查
// @Description 检查服务是否就绪(验证配置)
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /ready [get]
func (c *HealthController) Ready(w http.ResponseWriter, r *http.Request) {
	// 检查必要的环境变量
	issues := []string{}

	if os.Getenv("API_KEY") == "" {
		issues = append(issues, "API_KEY未设置")
	}

	if os.Getenv("LVYUN_BASE_URL") == "" {
		issues = append(issues, "LVYUN_BASE_URL未设置")
	}

	if len(issues) > 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		render.JSON(w, r, map[string]interface{}{
			"status": "not_ready",
			"issues": issues,
		})
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"status": "ready",
		"time":   time.Now().Format(time.RFC3339),
	})
}





