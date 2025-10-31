package controllers

import (
	"net/http"

	"datahub-external/service/proxy"
)

// LvyunController 绿云接口控制器
type LvyunController struct {
	registry *proxy.ClientRegistry
}

// NewLvyunController 创建绿云控制器
func NewLvyunController(registry *proxy.ClientRegistry) *LvyunController {
	return &LvyunController{
		registry: registry,
	}
}

// GetReservations 获取预订单数据
// @Summary 预订单数据
// @Description 获取酒店预订数据，包括账号、姓名、房型、房号、入住日期、离店日期等
// @Tags 绿云接口
// @Accept json
// @Produce json
// @Param hotel_code query string true "酒店代码"
// @Param start_date query string true "开始日期 (格式: 2025-10-20 00:00:00)"
// @Param end_date query string true "结束日期 (格式: 2025-10-21 00:00:00)"
// @Success 200 {object} APIResponse{data=[]lvyun.ReservationResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security ApiKeyAuth
// @Router /lvyun/reservations [get]
func (c *LvyunController) GetReservations(w http.ResponseWriter, r *http.Request) {
	hotelCode := r.URL.Query().Get("hotel_code")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	saveDB := r.URL.Query().Get("save_db") // 是否保存到数据库，默认true

	if hotelCode == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少酒店代码参数")
		return
	}
	if startDate == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少开始日期参数")
		return
	}
	if endDate == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少结束日期参数")
		return
	}

	client, err := c.registry.GetClient("lvyun")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "绿云客户端未注册")
		return
	}

	params := map[string]string{
		"hotel_code": hotelCode,
		"start_date": startDate,
		"end_date":   endDate,
	}

	result, err := client.HandleRequest(r.Context(), "/reservations", params)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// 默认保存到数据库（除非明确指定save_db=false）
	if saveDB != "false" {
		// 尝试保存到数据库
		if lvyunClient, ok := client.(interface {
			SaveReservationsToDB(ctx interface{}, data interface{}) error
		}); ok {
			if err := lvyunClient.SaveReservationsToDB(r.Context(), result); err != nil {
				// 保存失败只记录日志，不影响返回
				RespondSuccess(w, r, map[string]interface{}{
					"data":    result,
					"message": "数据查询成功，但保存到数据库失败: " + err.Error(),
				})
				return
			}
		}
	}

	RespondSuccess(w, r, result)
}

// GetRegistrations 获取登记单数据
// @Summary 登记单数据
// @Description 获取酒店登记单数据，包括账号、姓名、房型、房号、入住日期、离店日期、主单id等
// @Tags 绿云接口
// @Accept json
// @Produce json
// @Param hotel_code query string true "酒店代码"
// @Param start_date query string true "开始日期 (格式: 2025-10-20 00:00:00)"
// @Param end_date query string true "结束日期 (格式: 2025-10-21 00:00:00)"
// @Success 200 {object} APIResponse{data=[]lvyun.RegistrationResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security ApiKeyAuth
// @Router /lvyun/registrations [get]
func (c *LvyunController) GetRegistrations(w http.ResponseWriter, r *http.Request) {
	hotelCode := r.URL.Query().Get("hotel_code")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	saveDB := r.URL.Query().Get("save_db") // 是否保存到数据库，默认true

	if hotelCode == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少酒店代码参数")
		return
	}
	if startDate == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少开始日期参数")
		return
	}
	if endDate == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少结束日期参数")
		return
	}

	client, err := c.registry.GetClient("lvyun")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "绿云客户端未注册")
		return
	}

	params := map[string]string{
		"hotel_code": hotelCode,
		"start_date": startDate,
		"end_date":   endDate,
	}

	result, err := client.HandleRequest(r.Context(), "/registrations", params)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// 默认保存到数据库（除非明确指定save_db=false）
	if saveDB != "false" {
		if lvyunClient, ok := client.(interface {
			SaveRegistrationsToDB(ctx interface{}, data interface{}) error
		}); ok {
			if err := lvyunClient.SaveRegistrationsToDB(r.Context(), result); err != nil {
				RespondSuccess(w, r, map[string]interface{}{
					"data":    result,
					"message": "数据查询成功，但保存到数据库失败: " + err.Error(),
				})
				return
			}
		}
	}

	RespondSuccess(w, r, result)
}

// GetCheckouts 获取结账单数据
// @Summary 结账单数据
// @Description 获取结账单流水数据，包括结账流水号、账号、入账类型、入账代码、金额等
// @Tags 绿云接口
// @Accept json
// @Produce json
// @Param hotel_code query string true "酒店代码"
// @Param biz_date query string true "营业日期 (格式: 2025-10-20 00:00:00)"
// @Success 200 {object} APIResponse{data=[]lvyun.CheckoutResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security ApiKeyAuth
// @Router /lvyun/checkouts [get]
func (c *LvyunController) GetCheckouts(w http.ResponseWriter, r *http.Request) {
	hotelCode := r.URL.Query().Get("hotel_code")
	bizDate := r.URL.Query().Get("biz_date")
	saveDB := r.URL.Query().Get("save_db") // 是否保存到数据库，默认true

	if hotelCode == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少酒店代码参数")
		return
	}
	if bizDate == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少营业日期参数")
		return
	}

	client, err := c.registry.GetClient("lvyun")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "绿云客户端未注册")
		return
	}

	params := map[string]string{
		"hotel_code": hotelCode,
		"biz_date":   bizDate,
	}

	result, err := client.HandleRequest(r.Context(), "/checkouts", params)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// 默认保存到数据库（除非明确指定save_db=false）
	if saveDB != "false" {
		if lvyunClient, ok := client.(interface {
			SaveCheckoutsToDB(ctx interface{}, data interface{}) error
		}); ok {
			if err := lvyunClient.SaveCheckoutsToDB(r.Context(), result); err != nil {
				RespondSuccess(w, r, map[string]interface{}{
					"data":    result,
					"message": "数据查询成功，但保存到数据库失败: " + err.Error(),
				})
				return
			}
		}
	}

	RespondSuccess(w, r, result)
}

// GetBusinessReport 获取营业报表数据
// @Summary 营业报表数据
// @Description 获取营业日报表数据，包括代码、描述、本日发生、本月发生、本年发生等
// @Tags 绿云接口
// @Accept json
// @Produce json
// @Param hotel_code query string true "酒店代码"
// @Param start_date query string true "开始日期 (格式: 2022-10-20 00:00:00)"
// @Param end_date query string true "结束日期 (格式: 2022-10-21 00:00:00)"
// @Success 200 {object} APIResponse{data=[]lvyun.BusinessReportResponse}
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security ApiKeyAuth
// @Router /lvyun/business-report [get]
func (c *LvyunController) GetBusinessReport(w http.ResponseWriter, r *http.Request) {
	hotelCode := r.URL.Query().Get("hotel_code")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	saveDB := r.URL.Query().Get("save_db") // 是否保存到数据库，默认true

	if hotelCode == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少酒店代码参数")
		return
	}
	if startDate == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少开始日期参数")
		return
	}
	if endDate == "" {
		RespondError(w, r, http.StatusBadRequest, "缺少结束日期参数")
		return
	}

	client, err := c.registry.GetClient("lvyun")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "绿云客户端未注册")
		return
	}

	params := map[string]string{
		"hotel_code": hotelCode,
		"start_date": startDate,
		"end_date":   endDate,
	}

	result, err := client.HandleRequest(r.Context(), "/business-report", params)
	if err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	// 默认保存到数据库（除非明确指定save_db=false）
	if saveDB != "false" {
		if lvyunClient, ok := client.(interface {
			SaveBusinessReportsToDB(ctx interface{}, data interface{}) error
		}); ok {
			if err := lvyunClient.SaveBusinessReportsToDB(r.Context(), result); err != nil {
				RespondSuccess(w, r, map[string]interface{}{
					"data":    result,
					"message": "数据查询成功，但保存到数据库失败: " + err.Error(),
				})
				return
			}
		}
	}

	RespondSuccess(w, r, result)
}

// HealthCheck 绿云接口健康检查
// @Summary 绿云接口健康检查
// @Description 检查绿云接口连接状态
// @Tags 绿云接口
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Security ApiKeyAuth
// @Router /lvyun/health [get]
func (c *LvyunController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	client, err := c.registry.GetClient("lvyun")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "绿云客户端未注册")
		return
	}

	if err := client.HealthCheck(r.Context()); err != nil {
		RespondError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	RespondSuccess(w, r, map[string]interface{}{
		"status": "healthy",
		"client": "lvyun",
	})
}

// GetInfo 获取绿云接口信息
// @Summary 获取绿云接口信息
// @Description 返回绿云接口的基本信息和状态
// @Tags 绿云接口
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=proxy.ClientInfo}
// @Failure 404 {object} APIResponse
// @Security ApiKeyAuth
// @Router /lvyun/info [get]
func (c *LvyunController) GetInfo(w http.ResponseWriter, r *http.Request) {
	client, err := c.registry.GetClient("lvyun")
	if err != nil {
		RespondError(w, r, http.StatusNotFound, "绿云客户端未注册")
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
