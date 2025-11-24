package controllers

import (
	"net/http"
	"strconv"

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

// QueryTaskRecords 查询任务执行记录
// @Summary 查询任务执行记录
// @Description 查询定期任务的执行记录，支持按数据源、任务类型、状态过滤
// @Tags 数据源
// @Accept json
// @Produce json
// @Param data_source query string false "数据源名称（lvyun/ps）"
// @Param task_type query string false "任务类型（如：reservations, positions-inc等）"
// @Param status query string false "执行状态（running/success/failed）"
// @Param limit query int false "返回记录数（默认100）"
// @Success 200 {object} APIResponse{data=[]proxy.TaskRecord}
// @Security ApiKeyAuth
// @Router /task-records [get]
func (c *ProxyController) QueryTaskRecords(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	dataSource := r.URL.Query().Get("data_source")
	taskType := r.URL.Query().Get("task_type")
	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")

	// 解析limit参数
	limit := 100
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// 查询记录
	taskRecordService := proxy.GetGlobalTaskRecordService()
	records := taskRecordService.QueryRecords(dataSource, taskType, status, limit)

	RespondSuccess(w, r, records)
}

// GetTaskStatistics 获取任务统计信息
// @Summary 获取任务统计信息
// @Description 获取任务执行的统计信息（总数、成功数、失败数、运行中）
// @Tags 数据源
// @Accept json
// @Produce json
// @Param data_source query string false "数据源名称（lvyun/ps）"
// @Success 200 {object} APIResponse{data=map[string]interface{}}
// @Security ApiKeyAuth
// @Router /task-statistics [get]
func (c *ProxyController) GetTaskStatistics(w http.ResponseWriter, r *http.Request) {
	// 获取查询参数
	dataSource := r.URL.Query().Get("data_source")

	// 获取统计信息
	taskRecordService := proxy.GetGlobalTaskRecordService()
	stats := taskRecordService.GetStatistics(dataSource)

	RespondSuccess(w, r, stats)
}
