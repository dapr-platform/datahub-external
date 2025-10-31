package controllers

import (
	"net/http"

	"datahub-external/service/proxy"
)

// ProxyController 代理控制器
type ProxyController struct {
	registry *proxy.ClientRegistry
}

// NewProxyController 创建代理控制器
func NewProxyController(registry *proxy.ClientRegistry) *ProxyController {
	return &ProxyController{
		registry: registry,
	}
}

// ListDataSources 列出所有第三方数据源
// @Summary 列出所有数据源
// @Description 返回所有已注册的第三方数据源信息
// @Tags 数据源
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=[]proxy.ClientInfo}
// @Security ApiKeyAuth
// @Router /datasources [get]
func (c *ProxyController) ListDataSources(w http.ResponseWriter, r *http.Request) {
	clients := c.registry.ListClients()
	RespondSuccess(w, r, clients)
}
