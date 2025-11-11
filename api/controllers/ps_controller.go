package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"datahub-external/service/proxy"
	"datahub-external/service/proxy/ps"
)

// PSController PS系统接口控制器
type PSController struct {
	registry *proxy.ClientRegistry
}

// NewPSController 创建PS控制器
func NewPSController(registry *proxy.ClientRegistry) *PSController {
	return &PSController{
		registry: registry,
	}
}

// GetFamilyMembers 获取家族成员数据
// @Summary 家族成员数据
// @Description 获取PS系统家族成员信息，支持分页查询
// @Tags PS接口
// @Accept json
// @Produce json
// @Param ds query string true "分区字段 (格式: YYYYMMDD，示例: 20251013)"
// @Param page_num query int false "页编号 (默认: 1)"
// @Param page_size query int false "页大小 (默认: 100，最大: 2000)"
// @Param emplid query string false "员工ID"
// @Param save_db query string false "是否保存到数据库 (默认: true)"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security ApiKeyAuth
// @Router /ps/family-members [get]
func (c *PSController) GetFamilyMembers(w http.ResponseWriter, r *http.Request) {
	ds := r.URL.Query().Get("ds")
	pageNumStr := r.URL.Query().Get("page_num")
	pageSizeStr := r.URL.Query().Get("page_size")
	emplID := r.URL.Query().Get("emplid")
	saveDB := r.URL.Query().Get("save_db")

	// 验证必需参数
	if ds == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少分区字段参数(ds)")
		return
	}

	// 解析分页参数
	pageNum := 1
	if pageNumStr != "" {
		if n, err := strconv.Atoi(pageNumStr); err == nil && n > 0 {
			pageNum = n
		}
	}

	pageSize := 100
	if pageSizeStr != "" {
		if n, err := strconv.Atoi(pageSizeStr); err == nil && n > 0 {
			if n > 2000 {
				n = 2000 // API最大限制
			}
			pageSize = n
		}
	}

	// 获取PS客户端
	client, err := c.registry.GetClient("ps")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "PS客户端未注册")
		return
	}

	// 构建参数
	params := map[string]string{
		"ds":        ds,
		"page_num":  fmt.Sprintf("%d", pageNum),
		"page_size": fmt.Sprintf("%d", pageSize),
	}
	if emplID != "" {
		params["emplid"] = emplID
	}

	// 查询数据
	slog.Info("PS控制器查询家族成员",
		"ds", ds,
		"page_num", pageNum,
		"page_size", pageSize,
		"emplid", emplID)

	result, err := client.HandleRequest(r.Context(), "/family-members", params)
	if err != nil {
		slog.Error("查询家族成员失败", "error", err)
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// 打印第一条记录用于调试
	if resultArray, ok := result.([]interface{}); ok && len(resultArray) > 0 {
		firstRecord, _ := json.Marshal(resultArray[0])
		slog.Info("查询到家族成员数据",
			"total_count", len(resultArray),
			"first_record", string(firstRecord))
	} else {
		slog.Info("查询到家族成员数据", "result_type", fmt.Sprintf("%T", result))
	}

	// 默认保存到数据库（除非明确指定save_db=false）
	if saveDB != "false" {
		slog.Debug("尝试保存家族成员到数据库", "save_db", saveDB)
		if psClient, ok := client.(interface {
			SaveFamilyMembersToDB(ctx interface{}, data interface{}, ds string) error
		}); ok {
			slog.Info("开始保存家族成员到数据库", "ds", ds)
			if err := psClient.SaveFamilyMembersToDB(r.Context(), result, ds); err != nil {
				// 保存失败只记录日志，不影响返回
				slog.Error("保存家族成员到数据库失败", "error", err)
				RespondSuccess(w, r, map[string]interface{}{
					"data":    result,
					"message": "数据查询成功，但保存到数据库失败: " + err.Error(),
				})
				return
			}
			slog.Info("保存家族成员到数据库成功")
		} else {
			slog.Warn("客户端不支持SaveFamilyMembersToDB方法")
		}
	}

	RespondSuccess(w, r, result)
}

// HealthCheck PS接口健康检查
// @Summary PS接口健康检查
// @Description 检查PS接口连接状态
// @Tags PS接口
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security ApiKeyAuth
// @Router /ps/health [get]
func (c *PSController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	client, err := c.registry.GetClient("ps")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "PS客户端未注册")
		return
	}

	if err := client.HealthCheck(r.Context()); err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondSuccess(w, r, map[string]interface{}{
		"status": "healthy",
		"client": "ps",
	})
}

// GetInfo 获取PS接口信息
// @Summary 获取PS接口信息
// @Description 返回PS接口的基本信息和状态
// @Tags PS接口
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=proxy.ClientInfo}
// @Failure 404 {object} APIResponse
// @Security ApiKeyAuth
// @Router /ps/info [get]
func (c *PSController) GetInfo(w http.ResponseWriter, r *http.Request) {
	client, err := c.registry.GetClient("ps")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "PS客户端未注册")
		return
	}

	info := proxy.ClientInfo{
		Name:        client.GetName(),
		RoutePrefix: client.GetRoutePrefix(),
		Description: client.GetDescription(),
		Status:      client.GetStatus(),
	}

	RespondSuccess(w, r, info)
}

// TriggerSync 手动触发数据同步
// @Summary 手动触发数据同步
// @Description 手动触发PS系统数据同步任务
// @Tags PS接口
// @Accept json
// @Produce json
// @Param data_type query string false "数据类型 (family-members/all，默认: family-members)"
// @Success 200 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security ApiKeyAuth
// @Router /ps/trigger-sync [post]
func (c *PSController) TriggerSync(w http.ResponseWriter, r *http.Request) {
	dataType := r.URL.Query().Get("data_type")
	if dataType == "" {
		dataType = "family-members"
	}

	slog.Info("PS控制器触发手动同步", "data_type", dataType)

	client, err := c.registry.GetClient("ps")
	if err != nil {
		slog.Error("获取PS客户端失败", "error", err)
		RespondError(w, r, http.StatusNotFound, "PS客户端未注册")
		return
	}
	psClient, ok := client.(*ps.PSClient)
	if !ok {
		slog.Error("客户端类型断言失败",
			"client_type", fmt.Sprintf("%T", client),
			"expected", "PSClientInterface")
		RespondError(w, r, http.StatusInternalServerError, "客户端类型断言失败")
		return
	}
	scheduler := psClient.GetScheduler()
	if scheduler == nil {
		slog.Warn("调度器未启用", "scheduler", "nil")
		RespondError(w, r, http.StatusServiceUnavailable, "调度器未启用")
		return
	}
	if err := scheduler.TriggerSync(dataType); err != nil {
		slog.Error("触发同步任务失败", "error", err, "data_type", dataType)
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	RespondSuccess(w, r, map[string]interface{}{
		"message":   "同步任务已触发",
		"data_type": dataType,
	})

}
