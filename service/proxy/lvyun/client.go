package lvyun

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"datahub-external/service/config"
	"gorm.io/gorm"
)

// LvyunClient 绿云接口客户端
type LvyunClient struct {
	config     config.LvyunConfig
	httpClient *http.Client
	sessionID  string
	mu         sync.RWMutex
	stopChan   chan struct{}
	status     string
	repository *Repository
	scheduler  *Scheduler
}

// LoginResponse 登录响应
type LoginResponse struct {
	ResultCode int    `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	ResultInfo string `json:"resultInfo"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	ErrorToken string `json:"errorToken"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	Solution   string `json:"solution"`
}

// QueryResponse 查询响应
type QueryResponse struct {
	ResultCode int             `json:"resultCode"`
	ResultMsg  string          `json:"resultMsg"`
	ResultInfo json.RawMessage `json:"resultInfo"` // 使用RawMessage以支持不同的返回格式
}

// NewLvyunClient 创建绿云客户端
func NewLvyunClient(cfg config.LvyunConfig, db *gorm.DB) *LvyunClient {
	slog.Info("创建绿云客户端",
		"base_url", cfg.BaseURL,
		"hotel_group_code", cfg.HotelGroupCode,
		"user_code", cfg.UserCode)

	client := &LvyunClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		stopChan: make(chan struct{}),
		status:   "initialized",
	}

	// 初始化数据仓库
	if db != nil {
		client.repository = NewRepository(db)
		if err := client.repository.AutoMigrate(); err != nil {
			slog.Error("数据库自动迁移失败", "error", err)
		} else {
			slog.Info("数据库自动迁移成功")
		}
	}

	// 启动时自动登录
	if err := client.login(); err != nil {
		client.status = fmt.Sprintf("login_failed: %v", err)
		slog.Error("绿云客户端初始化登录失败", "error", err)
	} else {
		client.status = "active"
		slog.Info("绿云客户端初始化登录成功")
		// 启动自动刷新goroutine
		go client.autoRefresh()

		// 启动调度器
		if client.repository != nil && cfg.EnableScheduler {
			client.startScheduler()
		}
	}

	return client
}

// startScheduler 启动调度器
func (c *LvyunClient) startScheduler() {
	schedulerConfig := SchedulerConfig{
		ReservationCron:    c.config.ReservationCron,
		RegistrationCron:   c.config.RegistrationCron,
		CheckoutCron:       c.config.CheckoutCron,
		BusinessReportCron: c.config.BusinessReportCron,
		HotelCode:          c.config.HotelCode,
		QueryDays:          c.config.QueryDays,
		BusinessReportDays: c.config.BusinessReportDays,
	}

	c.scheduler = NewScheduler(schedulerConfig, c, c.repository)
	if err := c.scheduler.Start(); err != nil {
		slog.Error("启动调度器失败", "error", err)
	}
}

// GetName 获取客户端名称
func (c *LvyunClient) GetName() string {
	return "lvyun"
}

// GetRoutePrefix 获取路由前缀
func (c *LvyunClient) GetRoutePrefix() string {
	return "/lvyun"
}

// GetDescription 获取客户端描述
func (c *LvyunClient) GetDescription() string {
	return "绿云酒店管理系统接口"
}

// GetStatus 获取客户端状态
func (c *LvyunClient) GetStatus() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

// HandleRequest 处理请求
func (c *LvyunClient) HandleRequest(ctx context.Context, path string, params map[string]string) (interface{}, error) {
	// 根据路径路由到不同的查询方法
	switch path {
	case "/reservations":
		return c.queryReservations(ctx, params)
	case "/registrations":
		return c.queryRegistrations(ctx, params)
	case "/checkouts":
		return c.queryCheckouts(ctx, params)
	case "/business-report":
		return c.queryBusinessReport(ctx, params)
	default:
		return nil, fmt.Errorf("不支持的路径: %s", path)
	}
}

// HealthCheck 健康检查
func (c *LvyunClient) HealthCheck(ctx context.Context) error {
	c.mu.RLock()
	sessionID := c.sessionID
	c.mu.RUnlock()

	if sessionID == "" {
		return fmt.Errorf("未登录")
	}

	return nil
}

// login 登录获取sessionId
func (c *LvyunClient) login() error {
	params := map[string]string{
		"v":              "3.0",
		"hotelGroupCode": c.config.HotelGroupCode,
		"usercode":       c.config.UserCode,
		"password":       c.config.Password,
		"method":         "user.login",
		"local":          "zh_CN",
		"format":         "json",
		"appKey":         c.config.AppKey,
	}

	slog.Info("绿云客户端登录",
		"method", "user.login",
		"hotel_group_code", c.config.HotelGroupCode,
		"user_code", c.config.UserCode)

	// 计算签名
	sign := c.calculateSign(params)
	params["sign"] = sign

	slog.Debug("登录请求签名", "sign", sign)

	// 发送请求
	resp, err := c.postRequest(c.config.BaseURL+"/ipmsgroup/router", params)
	if err != nil {
		slog.Error("登录请求失败", "error", err)
		return err
	}

	slog.Debug("登录响应", "response", string(resp))

	var loginResp LoginResponse
	if err := json.Unmarshal(resp, &loginResp); err != nil {
		slog.Error("解析登录响应失败", "error", err, "response", string(resp))
		return fmt.Errorf("解析登录响应失败: %v", err)
	}

	if loginResp.ResultCode != 0 {
		slog.Error("登录失败",
			"result_code", loginResp.ResultCode,
			"result_msg", loginResp.ResultMsg)
		return fmt.Errorf("登录失败: %s", loginResp.ResultMsg)
	}

	// 保存sessionId
	c.mu.Lock()
	c.sessionID = loginResp.ResultInfo
	c.mu.Unlock()

	slog.Info("登录成功", "session_id", loginResp.ResultInfo)

	return nil
}

// refresh 刷新sessionId
func (c *LvyunClient) refresh() error {
	c.mu.RLock()
	sessionID := c.sessionID
	c.mu.RUnlock()

	if sessionID == "" {
		slog.Warn("sessionId为空，尝试重新登录")
		return c.login()
	}

	slog.Info("刷新sessionId", "old_session_id", sessionID)

	params := map[string]string{
		"v":         "3.0",
		"sessionId": sessionID,
		"method":    "user.refresh",
		"local":     "zh_CN",
		"format":    "json",
		"appKey":    c.config.AppKey,
	}

	// 计算签名
	sign := c.calculateSign(params)
	params["sign"] = sign

	// 发送请求
	resp, err := c.postRequest(c.config.BaseURL+"/ipmsgroup/router", params)
	if err != nil {
		slog.Error("刷新sessionId请求失败", "error", err)
		return err
	}

	slog.Debug("刷新响应", "response", string(resp))

	var loginResp LoginResponse
	if err := json.Unmarshal(resp, &loginResp); err != nil {
		slog.Error("解析刷新响应失败", "error", err, "response", string(resp))
		return fmt.Errorf("解析刷新响应失败: %v", err)
	}

	if loginResp.ResultCode != 0 {
		// 刷新失败,尝试重新登录
		slog.Warn("刷新sessionId失败，尝试重新登录",
			"result_code", loginResp.ResultCode,
			"result_msg", loginResp.ResultMsg)
		return c.login()
	}

	// 更新sessionId
	c.mu.Lock()
	c.sessionID = loginResp.ResultInfo
	c.mu.Unlock()

	slog.Info("刷新sessionId成功", "new_session_id", loginResp.ResultInfo)

	return nil
}

// autoRefresh 自动刷新sessionId
func (c *LvyunClient) autoRefresh() {
	ticker := time.NewTicker(11 * time.Hour) // 每11小时刷新一次
	defer ticker.Stop()

	slog.Info("启动sessionId自动刷新任务", "interval", "11小时")

	for {
		select {
		case <-ticker.C:
			slog.Info("触发sessionId自动刷新")
			if err := c.refresh(); err != nil {
				c.mu.Lock()
				c.status = fmt.Sprintf("refresh_failed: %v", err)
				c.mu.Unlock()
				slog.Error("自动刷新sessionId失败", "error", err)
			} else {
				c.mu.Lock()
				c.status = "active"
				c.mu.Unlock()
				slog.Info("自动刷新sessionId成功")
			}
		case <-c.stopChan:
			slog.Info("停止sessionId自动刷新任务")
			return
		}
	}
}

// Stop 停止客户端
func (c *LvyunClient) Stop() {
	if c.scheduler != nil {
		c.scheduler.Stop()
	}
	close(c.stopChan)
}

// GetScheduler 获取调度器
func (c *LvyunClient) GetScheduler() *Scheduler {
	return c.scheduler
}

// SaveReservationsToDB 保存预订单数据到数据库
func (c *LvyunClient) SaveReservationsToDB(ctx interface{}, data interface{}) error {
	if c.repository == nil {
		return fmt.Errorf("数据库未配置")
	}
	
	dataArray, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf("数据格式错误")
	}
	
	ctxTyped, ok := ctx.(context.Context)
	if !ok {
		ctxTyped = context.Background()
	}
	
	return c.repository.SaveReservations(ctxTyped, dataArray)
}

// SaveRegistrationsToDB 保存登记单数据到数据库
func (c *LvyunClient) SaveRegistrationsToDB(ctx interface{}, data interface{}) error {
	if c.repository == nil {
		return fmt.Errorf("数据库未配置")
	}
	
	dataArray, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf("数据格式错误")
	}
	
	ctxTyped, ok := ctx.(context.Context)
	if !ok {
		ctxTyped = context.Background()
	}
	
	return c.repository.SaveRegistrations(ctxTyped, dataArray)
}

// SaveCheckoutsToDB 保存结账单数据到数据库
func (c *LvyunClient) SaveCheckoutsToDB(ctx interface{}, data interface{}) error {
	if c.repository == nil {
		return fmt.Errorf("数据库未配置")
	}
	
	dataArray, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf("数据格式错误")
	}
	
	ctxTyped, ok := ctx.(context.Context)
	if !ok {
		ctxTyped = context.Background()
	}
	
	return c.repository.SaveCheckouts(ctxTyped, dataArray)
}

// SaveBusinessReportsToDB 保存营业报表数据到数据库
func (c *LvyunClient) SaveBusinessReportsToDB(ctx interface{}, data interface{}) error {
	if c.repository == nil {
		return fmt.Errorf("数据库未配置")
	}
	
	dataArray, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf("数据格式错误")
	}
	
	ctxTyped, ok := ctx.(context.Context)
	if !ok {
		ctxTyped = context.Background()
	}
	
	return c.repository.SaveBusinessReports(ctxTyped, dataArray)
}

// calculateSign 计算签名
func (c *LvyunClient) calculateSign(params map[string]string) string {
	// 1. 获取所有key并排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 2. 按照 appSecret + key1value1 + key2value2 + ... + appSecret 的格式拼接
	var sb bytes.Buffer
	sb.WriteString(c.config.AppSecret)
	for _, key := range keys {
		sb.WriteString(key)
		sb.WriteString(params[key])
	}
	sb.WriteString(c.config.AppSecret)

	// 3. SHA1加密
	h := sha1.New()
	h.Write(sb.Bytes())
	hash := h.Sum(nil)

	// 4. 转换为大写十六进制字符串
	return strings.ToUpper(hex.EncodeToString(hash))
}

// postRequest 发送POST请求
func (c *LvyunClient) postRequest(urlStr string, params map[string]string) ([]byte, error) {
	// 构建form-encoded请求体
	formData := url.Values{}
	// 创建一个不包含敏感信息的参数副本用于日志
	logParams := make(map[string]string)
	for key, value := range params {
		formData.Set(key, value)
		// 隐藏敏感信息
		if key == "password" || key == "sign" {
			logParams[key] = "***"
		} else {
			logParams[key] = value
		}
	}

	slog.Debug("发送POST请求",
		"url", urlStr,
		"params", logParams)

	req, err := http.NewRequest("POST", urlStr, strings.NewReader(formData.Encode()))
	if err != nil {
		slog.Error("创建HTTP请求失败", "error", err, "url", urlStr)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("HTTP请求执行失败", "error", err, "url", urlStr)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("读取响应体失败", "error", err)
		return nil, err
	}

	slog.Debug("收到HTTP响应",
		"url", urlStr,
		"status_code", resp.StatusCode,
		"body_length", len(body))

	if resp.StatusCode != http.StatusOK {
		slog.Error("HTTP请求返回非200状态码",
			"url", urlStr,
			"status_code", resp.StatusCode,
			"response", string(body))
		return nil, fmt.Errorf("请求失败,状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// query 通用查询方法（支持自动重新登录）
func (c *LvyunClient) query(ctx context.Context, exec string, queryParams map[string]string) (interface{}, error) {
	// 最多尝试2次：第一次使用现有sessionID，如果失败则重新登录后再试一次
	maxRetries := 2
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		c.mu.RLock()
		sessionID := c.sessionID
		c.mu.RUnlock()

		if sessionID == "" {
			slog.Warn("查询时发现未登录，尝试重新登录", "attempt", attempt)
			if err := c.login(); err != nil {
				slog.Error("重新登录失败", "error", err, "attempt", attempt)
				if attempt == maxRetries {
					return nil, fmt.Errorf("重新登录失败: %v", err)
				}
				continue
			}
			// 登录成功后，更新sessionID
			c.mu.RLock()
			sessionID = c.sessionID
			c.mu.RUnlock()
		}

		// 构建查询参数
		params := map[string]string{
			"method":         "crs.kpi",
			"v":              "3.0",
			"format":         "json",
			"local":          "zh_CN",
			"appKey":         c.config.AppKey,
			"sessionId":      sessionID,
			"hotelGroupCode": c.config.HotelGroupCode,
			"exec":           exec,
		}

		// 添加额外参数
		for key, value := range queryParams {
			params[key] = value
		}

		slog.Info("绿云查询请求",
			"exec", exec,
			"hotel_group_code", c.config.HotelGroupCode,
			"params", queryParams,
			"attempt", attempt)

		// 计算签名
		sign := c.calculateSign(params)
		params["sign"] = sign

		slog.Debug("查询请求签名", "sign", sign)

		// 发送请求
		resp, err := c.postRequest(c.config.BaseURL+"/ipmsgroup/router", params)
		if err != nil {
			slog.Error("查询请求失败", "exec", exec, "error", err, "attempt", attempt)
			return nil, err
		}

		slog.Debug("查询响应", "exec", exec, "response", string(resp))

		var queryResp QueryResponse
		if err := json.Unmarshal(resp, &queryResp); err != nil {
			slog.Error("解析查询响应失败",
				"exec", exec,
				"error", err,
				"response", string(resp),
				"attempt", attempt)
			return nil, fmt.Errorf("解析查询响应失败: %v", err)
		}

		if queryResp.ResultCode != 0 {
			// 检查是否为未登录错误
			if strings.Contains(queryResp.ResultMsg, "未登录") || 
			   strings.Contains(queryResp.ResultMsg, "登录") ||
			   strings.Contains(queryResp.ResultMsg, "session") {
				slog.Warn("检测到登录失效",
					"exec", exec,
					"result_code", queryResp.ResultCode,
					"result_msg", queryResp.ResultMsg,
					"attempt", attempt)
				
				// 清空sessionID
				c.mu.Lock()
				c.sessionID = ""
				c.mu.Unlock()
				
				// 如果还有重试机会，继续重试
				if attempt < maxRetries {
					slog.Info("将在下一次尝试中重新登录", "attempt", attempt+1)
					continue
				}
			}
			
			slog.Error("查询失败",
				"exec", exec,
				"result_code", queryResp.ResultCode,
				"result_msg", queryResp.ResultMsg,
				"attempt", attempt)
			return nil, fmt.Errorf("查询失败: %s", queryResp.ResultMsg)
		}

		// 解析resultInfo为数组
		var resultData []interface{}
		if len(queryResp.ResultInfo) > 0 {
			if err := json.Unmarshal(queryResp.ResultInfo, &resultData); err != nil {
				slog.Error("解析ResultInfo失败",
					"exec", exec,
					"error", err,
					"resultInfo", string(queryResp.ResultInfo),
					"attempt", attempt)
				return nil, fmt.Errorf("解析ResultInfo失败: %v", err)
			}
		}

		slog.Info("查询成功",
			"exec", exec,
			"result_count", len(resultData),
			"attempt", attempt)

		return resultData, nil
	}
	
	return nil, fmt.Errorf("查询失败：已达最大重试次数")
}

// queryReservations 预订单数据
func (c *LvyunClient) queryReservations(ctx context.Context, params map[string]string) (interface{}, error) {
	hotelCode := params["hotel_code"]
	startDate := params["start_date"]
	endDate := params["end_date"]

	queryParams := map[string]string{
		"hotelCode": hotelCode,
		"params":    c.config.HotelGroupCode + "," + hotelCode + "," + startDate + "," + endDate,
	}

	return c.query(ctx, "Kpi_Ihotel_Rep_Rsvsrc", queryParams)
}

// queryRegistrations 登记单数据
func (c *LvyunClient) queryRegistrations(ctx context.Context, params map[string]string) (interface{}, error) {
	hotelCode := params["hotel_code"]
	startDate := params["start_date"]
	endDate := params["end_date"]

	queryParams := map[string]string{
		"hotelCode": hotelCode,
		"params":    c.config.HotelGroupCode + "," + hotelCode + "," + startDate + "," + endDate,
	}

	return c.query(ctx, "Kpi_Ihotel_Rep_Master", queryParams)
}

// queryCheckouts 结账单数据
func (c *LvyunClient) queryCheckouts(ctx context.Context, params map[string]string) (interface{}, error) {
	hotelCode := params["hotel_code"]
	bizDate := params["biz_date"]

	queryParams := map[string]string{
		"hotelCode": hotelCode,
		"params":    c.config.HotelGroupCode + "," + hotelCode + "," + bizDate,
	}

	return c.query(ctx, "Kpi_Ihotel_Rep_Account", queryParams)
}

// queryBusinessReport 营业报表数据
func (c *LvyunClient) queryBusinessReport(ctx context.Context, params map[string]string) (interface{}, error) {
	hotelCode := params["hotel_code"]
	startDate := params["start_date"]
	endDate := params["end_date"]

	queryParams := map[string]string{
		"hotelCode": hotelCode,
		"params":    c.config.HotelGroupCode + "," + hotelCode + "," + startDate + "," + endDate,
	}

	return c.query(ctx, "Kpi_Report_Jour", queryParams)
}
