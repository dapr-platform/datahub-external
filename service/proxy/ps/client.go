package ps

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
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

	"github.com/google/uuid"
)

const (
	// ContentTypeForm 表单类型
	ContentTypeForm = "application/x-www-form-urlencoded"
	// ContentTypeJSON JSON类型
	ContentTypeJSON = "application/json"
	// ContentTypeText 文本类型
	ContentTypeText = "text/plain"
	// DefaultStage 默认stage
	DefaultStage = "RELEASE"
	// MaxRetries 最大重试次数
	MaxRetries = 3
	// RetryDelay 重试延迟
	RetryDelay = 2 * time.Second
)

// PSClient PS系统API网关客户端
type PSClient struct {
	appKey     string
	appSecret  []byte
	stage      string
	baseURL    string
	httpClient *http.Client
	mu         sync.RWMutex
	status     string
	repository *Repository
	scheduler  *Scheduler
	stopChan   chan struct{}
}

// RequestOptions 请求选项
type RequestOptions struct {
	Query       map[string]string // 查询参数
	Data        interface{}       // 请求体数据
	Headers     map[string]string // 自定义headers
	SignHeaders map[string]string // 需要参与签名的headers
	Timeout     time.Duration     // 请求超时时间
}

// Response 响应结构
type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

// PSConfig PS客户端配置
type PSConfig struct {
	AppKey              string
	AppSecret           string
	Stage               string
	BaseURL             string
	EnableScheduler     bool
	FamilyMemberCron    string
	PositionIncCron     string
	PositionAllCron     string
	OrganizationIncCron string
	OrganizationAllCron string
	EmployeeIncCron     string
	EmployeeAllCron     string
	EmployeeHonorCron   string
	FamilyMainCron      string
	PageSize            int
	MaxPages            int
}

// NewPSClient 创建PS系统客户端（简单版本，不带数据库）
// appKey: 应用标识
// appSecret: 应用密钥
// stage: 环境标识(RELEASE/PRE/TEST)
func NewPSClient(appKey, appSecret string, stage ...string) *PSClient {
	stageValue := DefaultStage
	if len(stage) > 0 && stage[0] != "" {
		stageValue = stage[0]
	}

	slog.Info("创建PS系统客户端",
		"app_key", appKey,
		"stage", stageValue)

	client := &PSClient{
		appKey:    appKey,
		appSecret: []byte(appSecret),
		stage:     stageValue,
		baseURL:   "https://datadisclose.hdlapis.com",
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // 增加到5分钟，处理大数据量请求
		},
		status:   "active",
		stopChan: make(chan struct{}),
	}

	return client
}

// NewPSClientWithConfig 创建PS系统客户端（带配置和数据库支持）
func NewPSClientWithConfig(cfg PSConfig, db interface{}) *PSClient {
	slog.Info("创建PS系统客户端（完整配置）",
		"app_key", cfg.AppKey,
		"stage", cfg.Stage,
		"base_url", cfg.BaseURL,
		"enable_scheduler", cfg.EnableScheduler)

	client := &PSClient{
		appKey:    cfg.AppKey,
		appSecret: []byte(cfg.AppSecret),
		stage:     cfg.Stage,
		baseURL:   cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // 增加到5分钟，处理大数据量请求
		},
		status:   "active",
		stopChan: make(chan struct{}),
	}

	// 初始化数据仓库（如果提供了数据库连接）
	if db != nil {
		// 尝试将db断言为*gorm.DB
		type GormDB interface {
			AutoMigrate(...interface{}) error
		}
		if gormDB, ok := db.(GormDB); ok {
			// 由于我们无法直接使用接口类型，这里需要Repository自己处理
			// 通过反射或者类型断言来处理
			client.repository = NewRepository(db)
			if err := client.repository.AutoMigrate(); err != nil {
				slog.Error("数据库自动迁移失败", "error", err)
			} else {
				slog.Info("数据库自动迁移成功")
			}
			_ = gormDB // 避免未使用警告
		}
	}

	// 启动调度器（如果启用）
	if client.repository != nil && cfg.EnableScheduler {
		client.startScheduler(cfg)
	}

	return client
}

// startScheduler 启动调度器
func (c *PSClient) startScheduler(cfg PSConfig) {
	schedulerConfig := SchedulerConfig{
		FamilyMemberCron:    cfg.FamilyMemberCron,
		PositionIncCron:     cfg.PositionIncCron,
		PositionAllCron:     cfg.PositionAllCron,
		OrganizationIncCron: cfg.OrganizationIncCron,
		OrganizationAllCron: cfg.OrganizationAllCron,
		EmployeeIncCron:     cfg.EmployeeIncCron,
		EmployeeAllCron:     cfg.EmployeeAllCron,
		EmployeeHonorCron:   cfg.EmployeeHonorCron,
		FamilyMainCron:      cfg.FamilyMainCron,
		PageSize:            cfg.PageSize,
		MaxPages:            cfg.MaxPages,
	}

	c.scheduler = NewScheduler(schedulerConfig, c, c.repository)
	if err := c.scheduler.Start(); err != nil {
		slog.Error("启动调度器失败", "error", err)
	}
}

// Stop 停止客户端
func (c *PSClient) Stop() {
	if c.scheduler != nil {
		c.scheduler.Stop()
	}
	close(c.stopChan)
}

// GetScheduler 获取调度器
func (c *PSClient) GetScheduler() *Scheduler {
	return c.scheduler
}

// GetName 获取客户端名称
func (c *PSClient) GetName() string {
	return "ps"
}

// GetRoutePrefix 获取路由前缀
func (c *PSClient) GetRoutePrefix() string {
	return "/ps"
}

// GetDescription 获取客户端描述
func (c *PSClient) GetDescription() string {
	return "PS系统API网关接口"
}

// GetStatus 获取客户端状态
func (c *PSClient) GetStatus() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.status
}

// HandleRequest 处理请求
func (c *PSClient) HandleRequest(ctx context.Context, path string, params map[string]string) (interface{}, error) {
	// 根据路径路由到不同的方法
	switch path {
	case "/family-members":
		return c.queryFamilyMembersWithParams(ctx, params)
	case "/positions-inc":
		return c.queryPositionsIncWithParams(ctx, params)
	case "/positions-all":
		return c.queryPositionsAllWithParams(ctx, params)
	case "/organizations-inc":
		return c.queryOrganizationsIncWithParams(ctx, params)
	case "/organizations-all":
		return c.queryOrganizationsAllWithParams(ctx, params)
	case "/employees-inc":
		return c.queryEmployeesIncWithParams(ctx, params)
	case "/employees-all":
		return c.queryEmployeesAllWithParams(ctx, params)
	case "/employee-honors":
		return c.queryEmployeeHonorsWithParams(ctx, params)
	case "/family-main":
		return c.queryFamilyMainWithParams(ctx, params)
	default:
		return nil, fmt.Errorf("不支持的路径: %s", path)
	}
}

// queryFamilyMembersWithParams 根据参数查询家族成员
func (c *PSClient) queryFamilyMembersWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]
	emplID := params["emplid"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryFamilyMembers(ctx, pageNum, pageSize, ds, emplID)
}

// QueryFamilyMembers 查询家族成员信息
func (c *PSClient) QueryFamilyMembers(ctx context.Context, pageNum, pageSize int, ds, emplID string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/ps/ps_c_family_mbr_t", c.baseURL)

	// 构建查询参数（作为URL参数，不是body）
	queryParams := map[string]string{
		"pageNum":  fmt.Sprintf("%d", pageNum),
		"pageSize": fmt.Sprintf("%d", pageSize),
		"ds":       ds,
	}

	if emplID != "" {
		queryParams["emplid"] = emplID
	}

	// 详细日志已在doRequest中记录，这里只记录查询参数

	// 发送POST请求，参数作为query参数
	result, err := c.Post(ctx, url, &RequestOptions{
		Query: queryParams,
	})

	if err != nil {
		slog.Error("查询家族成员失败", "error", err, "url", url)
		return nil, err
	}

	slog.Debug("收到API响应", "result_type", fmt.Sprintf("%T", result))

	// 解析响应
	var apiResp PSPaginationResponse
	jsonData, err := json.Marshal(result)
	if err != nil {
		slog.Error("序列化响应失败", "error", err)
		return nil, fmt.Errorf("序列化响应失败: %v", err)
	}

	slog.Debug("响应JSON数据", "json", string(jsonData))

	if err := json.Unmarshal(jsonData, &apiResp); err != nil {
		slog.Error("解析响应失败", "error", err, "json", string(jsonData))
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	slog.Debug("解析API响应",
		"err_code", apiResp.ErrCode,
		"err_msg", apiResp.ErrMsg,
		"request_id", apiResp.RequestID,
		"total_num", apiResp.Data.TotalNum)

	if apiResp.ErrCode != 0 {
		slog.Error("API返回错误",
			"err_code", apiResp.ErrCode,
			"err_msg", apiResp.ErrMsg,
			"request_id", apiResp.RequestID)
		return nil, fmt.Errorf("API返回错误: %s (errCode: %d, requestId: %s)",
			apiResp.ErrMsg, apiResp.ErrCode, apiResp.RequestID)
	}

	// 转换为interface{}数组
	records := make([]interface{}, 0, len(apiResp.Data.Rows))
	for i, record := range apiResp.Data.Rows {
		records = append(records, record)
		// 打印第一条记录
		if i == 0 {
			firstRecordJSON, _ := json.Marshal(record)
			slog.Info("第一条家族成员记录",
				"record", string(firstRecordJSON),
				"c_family_id", record.CFamilyID,
				"emplid", record.EmplID,
				"c_family_name", record.CFamilyName)
		}
	}

	slog.Info("查询家族成员成功",
		"count", len(records),
		"total_num", apiResp.Data.TotalNum,
		"page_num", pageNum,
		"page_size", pageSize,
		"request_id", apiResp.RequestID)

	return records, nil
}

// SaveFamilyMembersToDB 保存家族成员数据到数据库
func (c *PSClient) SaveFamilyMembersToDB(ctx interface{}, data interface{}, ds string) error {
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

	return c.repository.SaveFamilyMembers(ctxTyped, dataArray, ds)
}

// queryPositionsIncWithParams 根据参数查询岗位增量数据
func (c *PSClient) queryPositionsIncWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryPositions(ctx, pageNum, pageSize, ds, true)
}

// queryPositionsAllWithParams 根据参数查询岗位全量数据
func (c *PSClient) queryPositionsAllWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryPositions(ctx, pageNum, pageSize, ds, false)
}

// QueryPositions 查询岗位信息
func (c *PSClient) QueryPositions(ctx context.Context, pageNum, pageSize int, ds string, isIncremental bool) ([]interface{}, error) {
	var apiPath string
	if isIncremental {
		apiPath = "/ps/ps_int_pos_inc"
	} else {
		apiPath = "/ps/ps_int_pos_all"
	}
	url := fmt.Sprintf("%s%s", c.baseURL, apiPath)

	queryParams := map[string]string{
		"pageNum":  fmt.Sprintf("%d", pageNum),
		"pageSize": fmt.Sprintf("%d", pageSize),
		"ds":       ds,
	}

	// 详细日志已在doRequest中记录

	result, err := c.Post(ctx, url, &RequestOptions{
		Query: queryParams,
	})

	if err != nil {
		slog.Error("查询岗位信息失败", "error", err, "url", url)
		return nil, err
	}

	var apiResp PositionPaginationResponse
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化响应失败: %v", err)
	}

	if err := json.Unmarshal(jsonData, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if apiResp.ErrCode != 0 {
		return nil, fmt.Errorf("API返回错误: %s (errCode: %d)", apiResp.ErrMsg, apiResp.ErrCode)
	}

	records := make([]interface{}, 0, len(apiResp.Data.Rows))
	for _, record := range apiResp.Data.Rows {
		records = append(records, record)
	}

	slog.Info("查询岗位信息成功",
		"count", len(records),
		"total_num", apiResp.Data.TotalNum,
		"page_num", pageNum,
		"request_id", apiResp.RequestID)

	return records, nil
}

// SavePositionsToDB 保存岗位数据到数据库
func (c *PSClient) SavePositionsToDB(ctx interface{}, data interface{}, ds string, isIncremental bool) error {
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

	return c.repository.SavePositions(ctxTyped, dataArray, ds)
}

// queryOrganizationsIncWithParams 根据参数查询组织增量数据
func (c *PSClient) queryOrganizationsIncWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryOrganizations(ctx, pageNum, pageSize, ds, true)
}

// queryOrganizationsAllWithParams 根据参数查询组织全量数据
func (c *PSClient) queryOrganizationsAllWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryOrganizations(ctx, pageNum, pageSize, ds, false)
}

// QueryOrganizations 查询组织信息
func (c *PSClient) QueryOrganizations(ctx context.Context, pageNum, pageSize int, ds string, isIncremental bool) ([]interface{}, error) {
	var apiPath string
	if isIncremental {
		apiPath = "/ps/ps_int_org_inc"
	} else {
		apiPath = "/ps/ps_int_org_all"
	}
	url := fmt.Sprintf("%s%s", c.baseURL, apiPath)

	queryParams := map[string]string{
		"pageNum":  fmt.Sprintf("%d", pageNum),
		"pageSize": fmt.Sprintf("%d", pageSize),
		"ds":       ds,
	}

	// 详细日志已在doRequest中记录

	result, err := c.Post(ctx, url, &RequestOptions{
		Query: queryParams,
	})

	if err != nil {
		slog.Error("查询组织信息失败", "error", err, "url", url)
		return nil, err
	}

	var apiResp OrganizationPaginationResponse
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化响应失败: %v", err)
	}

	if err := json.Unmarshal(jsonData, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if apiResp.ErrCode != 0 {
		return nil, fmt.Errorf("API返回错误: %s (errCode: %d)", apiResp.ErrMsg, apiResp.ErrCode)
	}

	records := make([]interface{}, 0, len(apiResp.Data.Rows))
	for _, record := range apiResp.Data.Rows {
		records = append(records, record)
	}

	slog.Info("查询组织信息成功",
		"count", len(records),
		"total_num", apiResp.Data.TotalNum,
		"page_num", pageNum,
		"request_id", apiResp.RequestID)

	return records, nil
}

// SaveOrganizationsToDB 保存组织数据到数据库
func (c *PSClient) SaveOrganizationsToDB(ctx interface{}, data interface{}, ds string, isIncremental bool) error {
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

	return c.repository.SaveOrganizations(ctxTyped, dataArray, ds)
}

// queryEmployeesIncWithParams 根据参数查询员工增量数据
func (c *PSClient) queryEmployeesIncWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]
	emplID := params["emplid"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryEmployees(ctx, pageNum, pageSize, ds, emplID, true)
}

// queryEmployeesAllWithParams 根据参数查询员工全量数据
func (c *PSClient) queryEmployeesAllWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]
	emplID := params["emplid"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryEmployees(ctx, pageNum, pageSize, ds, emplID, false)
}

// QueryEmployees 查询员工信息
func (c *PSClient) QueryEmployees(ctx context.Context, pageNum, pageSize int, ds, emplID string, isIncremental bool) ([]interface{}, error) {
	var apiPath string
	if isIncremental {
		apiPath = "/ps/ps_int_emp_inc_yq"
	} else {
		apiPath = "/ps/ps_int_emp_all_yq"
	}
	url := fmt.Sprintf("%s%s", c.baseURL, apiPath)

	queryParams := map[string]string{
		"pageNum":  fmt.Sprintf("%d", pageNum),
		"pageSize": fmt.Sprintf("%d", pageSize),
		"ds":       ds,
	}

	if emplID != "" {
		queryParams["emplid"] = emplID
	}

	// 详细日志已在doRequest中记录

	result, err := c.Post(ctx, url, &RequestOptions{
		Query: queryParams,
	})

	if err != nil {
		slog.Error("查询员工信息失败", "error", err, "url", url)
		return nil, err
	}

	var apiResp EmployeePaginationResponse
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化响应失败: %v", err)
	}

	if err := json.Unmarshal(jsonData, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if apiResp.ErrCode != 0 {
		return nil, fmt.Errorf("API返回错误: %s (errCode: %d)", apiResp.ErrMsg, apiResp.ErrCode)
	}

	records := make([]interface{}, 0, len(apiResp.Data.Rows))
	for _, record := range apiResp.Data.Rows {
		records = append(records, record)
	}

	slog.Info("查询员工信息成功",
		"count", len(records),
		"total_num", apiResp.Data.TotalNum,
		"page_num", pageNum,
		"request_id", apiResp.RequestID)

	return records, nil
}

// SaveEmployeesToDB 保存员工数据到数据库
func (c *PSClient) SaveEmployeesToDB(ctx interface{}, data interface{}, ds string, isIncremental bool) error {
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

	return c.repository.SaveEmployees(ctxTyped, dataArray, ds)
}

// queryEmployeeHonorsWithParams 根据参数查询员工荣誉数据
func (c *PSClient) queryEmployeeHonorsWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]
	emplID := params["emplid"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryEmployeeHonors(ctx, pageNum, pageSize, ds, emplID)
}

// QueryEmployeeHonors 查询员工荣誉信息
func (c *PSClient) QueryEmployeeHonors(ctx context.Context, pageNum, pageSize int, ds, emplID string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/ps/ps_c_honor_emp_tbl", c.baseURL)

	queryParams := map[string]string{
		"pageNum":  fmt.Sprintf("%d", pageNum),
		"pageSize": fmt.Sprintf("%d", pageSize),
		"ds":       ds,
	}

	if emplID != "" {
		queryParams["emplid"] = emplID
	}

	// 详细日志已在doRequest中记录

	result, err := c.Post(ctx, url, &RequestOptions{
		Query: queryParams,
	})

	if err != nil {
		slog.Error("查询员工荣誉信息失败", "error", err, "url", url)
		return nil, err
	}

	var apiResp EmployeeHonorPaginationResponse
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化响应失败: %v", err)
	}

	if err := json.Unmarshal(jsonData, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if apiResp.ErrCode != 0 {
		return nil, fmt.Errorf("API返回错误: %s (errCode: %d)", apiResp.ErrMsg, apiResp.ErrCode)
	}

	records := make([]interface{}, 0, len(apiResp.Data.Rows))
	for _, record := range apiResp.Data.Rows {
		records = append(records, record)
	}

	slog.Info("查询员工荣誉信息成功",
		"count", len(records),
		"total_num", apiResp.Data.TotalNum,
		"page_num", pageNum,
		"request_id", apiResp.RequestID)

	return records, nil
}

// SaveEmployeeHonorsToDB 保存员工荣誉数据到数据库
func (c *PSClient) SaveEmployeeHonorsToDB(ctx interface{}, data interface{}, ds string) error {
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

	return c.repository.SaveEmployeeHonors(ctxTyped, dataArray, ds)
}

// queryFamilyMainWithParams 根据参数查询家族父表数据
func (c *PSClient) queryFamilyMainWithParams(ctx context.Context, params map[string]string) (interface{}, error) {
	pageNum := 1
	pageSize := 100
	ds := params["ds"]
	emplID := params["emplid"]

	if params["page_num"] != "" {
		fmt.Sscanf(params["page_num"], "%d", &pageNum)
	}
	if params["page_size"] != "" {
		fmt.Sscanf(params["page_size"], "%d", &pageSize)
	}

	return c.QueryFamilyMain(ctx, pageNum, pageSize, ds, emplID)
}

// QueryFamilyMain 查询家族父表信息
func (c *PSClient) QueryFamilyMain(ctx context.Context, pageNum, pageSize int, ds, emplID string) ([]interface{}, error) {
	url := fmt.Sprintf("%s/ps/ps_c_family_main_t", c.baseURL)

	queryParams := map[string]string{
		"pageNum":  fmt.Sprintf("%d", pageNum),
		"pageSize": fmt.Sprintf("%d", pageSize),
		"ds":       ds,
	}

	if emplID != "" {
		queryParams["emplid"] = emplID
	}

	// 详细日志已在doRequest中记录

	result, err := c.Post(ctx, url, &RequestOptions{
		Query: queryParams,
	})

	if err != nil {
		slog.Error("查询家族父表信息失败", "error", err, "url", url)
		return nil, err
	}

	var apiResp FamilyMainPaginationResponse
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化响应失败: %v", err)
	}

	if err := json.Unmarshal(jsonData, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if apiResp.ErrCode != 0 {
		return nil, fmt.Errorf("API返回错误: %s (errCode: %d)", apiResp.ErrMsg, apiResp.ErrCode)
	}

	records := make([]interface{}, 0, len(apiResp.Data.Rows))
	for _, record := range apiResp.Data.Rows {
		records = append(records, record)
	}

	slog.Info("查询家族父表信息成功",
		"count", len(records),
		"total_num", apiResp.Data.TotalNum,
		"page_num", pageNum,
		"request_id", apiResp.RequestID)

	return records, nil
}

// SaveFamilyMainToDB 保存家族父表数据到数据库
func (c *PSClient) SaveFamilyMainToDB(ctx interface{}, data interface{}, ds string) error {
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

	return c.repository.SaveFamilyMain(ctxTyped, dataArray, ds)
}

// HealthCheck 健康检查
func (c *PSClient) HealthCheck(ctx context.Context) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.status != "active" {
		return fmt.Errorf("客户端状态异常: %s", c.status)
	}

	return nil
}

// Get 发送GET请求
func (c *PSClient) Get(ctx context.Context, urlStr string, opts *RequestOptions) (interface{}, error) {
	if opts == nil {
		opts = &RequestOptions{}
	}

	// GET请求将query或data合并到URL查询参数
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("解析URL失败: %v", err)
	}

	// 合并查询参数
	query := parsedURL.Query()
	if opts.Query != nil {
		for k, v := range opts.Query {
			query.Set(k, v)
		}
	}
	if opts.Data != nil {
		if dataMap, ok := opts.Data.(map[string]string); ok {
			for k, v := range dataMap {
				query.Set(k, v)
			}
		}
	}

	parsedURL.RawQuery = query.Encode()
	opts.Data = nil
	opts.Query = nil

	return c.request(ctx, "GET", parsedURL.String(), opts)
}

// Post 发送POST请求
func (c *PSClient) Post(ctx context.Context, urlStr string, opts *RequestOptions) (interface{}, error) {
	if opts == nil {
		opts = &RequestOptions{}
	}

	// 如果有Query参数，附加到URL上
	if len(opts.Query) > 0 {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return nil, fmt.Errorf("解析URL失败: %v", err)
		}

		query := parsedURL.Query()
		for k, v := range opts.Query {
			query.Set(k, v)
		}
		parsedURL.RawQuery = query.Encode()
		urlStr = parsedURL.String()

		slog.Debug("POST请求附加query参数", "url", urlStr)
	}

	// 设置默认Content-Type（仅当有Data时）
	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}
	if opts.Data != nil && opts.Headers["content-type"] == "" && opts.Headers["Content-Type"] == "" {
		opts.Headers["content-type"] = ContentTypeJSON
	}

	return c.request(ctx, "POST", urlStr, opts)
}

// Put 发送PUT请求
func (c *PSClient) Put(ctx context.Context, urlStr string, opts *RequestOptions) (interface{}, error) {
	if opts == nil {
		opts = &RequestOptions{}
	}

	// 设置默认Content-Type
	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}
	if opts.Headers["content-type"] == "" && opts.Headers["Content-Type"] == "" {
		opts.Headers["content-type"] = ContentTypeJSON
	}

	return c.request(ctx, "PUT", urlStr, opts)
}

// Delete 发送DELETE请求
func (c *PSClient) Delete(ctx context.Context, urlStr string, opts *RequestOptions) (interface{}, error) {
	if opts == nil {
		opts = &RequestOptions{}
	}

	// DELETE请求将query或data合并到URL查询参数
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("解析URL失败: %v", err)
	}

	// 合并查询参数
	query := parsedURL.Query()
	if opts.Query != nil {
		for k, v := range opts.Query {
			query.Set(k, v)
		}
	}
	if opts.Data != nil {
		if dataMap, ok := opts.Data.(map[string]string); ok {
			for k, v := range dataMap {
				query.Set(k, v)
			}
		}
	}

	parsedURL.RawQuery = query.Encode()
	opts.Data = nil
	opts.Query = nil

	return c.request(ctx, "DELETE", parsedURL.String(), opts)
}

// request 发送HTTP请求
func (c *PSClient) request(ctx context.Context, method, urlStr string, opts *RequestOptions) (interface{}, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("解析URL失败: %v", err)
	}

	// 处理查询参数
	if opts.Query != nil {
		query := parsedURL.Query()
		for k, v := range opts.Query {
			query.Set(k, v)
		}
		parsedURL.RawQuery = query.Encode()
	}

	// 构建headers
	headers := c.buildHeaders(opts.Headers, opts.SignHeaders)

	// 处理请求体
	var bodyData []byte
	var originData interface{} = opts.Data

	if opts.Data != nil && (method == "POST" || method == "PUT") {
		contentType := c.getContentType(headers)
		bodyData, err = c.encodeBody(opts.Data, contentType)
		if err != nil {
			return nil, fmt.Errorf("编码请求体失败: %v", err)
		}

		// 如果不是form类型，计算Content-MD5
		if !strings.HasPrefix(contentType, ContentTypeForm) {
			md5Hash := c.md5(bodyData)
			headers["content-md5"] = md5Hash
		}
	}

	// 获取需要签名的header keys
	signHeaderKeys := c.getSignHeaderKeys(headers, opts.SignHeaders)
	headers["x-ca-signature-headers"] = strings.Join(signHeaderKeys, ",")

	// 构建签名字符串
	signedHeadersStr := c.getSignedHeadersString(signHeaderKeys, headers)
	stringToSign := c.buildStringToSign(method, headers, signedHeadersStr, parsedURL, originData)

	// 签名
	signature := c.sign(stringToSign)
	headers["x-ca-signature"] = signature

	// 添加User-Agent
	headers["user-agent"] = "ps-api-gateway-go-sdk/1.0.0"

	slog.Debug("PS系统API请求",
		"method", method,
		"url", urlStr,
		"headers", headers)

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, method, parsedURL.String(), bytes.NewReader(bodyData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 设置超时
	client := c.httpClient
	if opts.Timeout > 0 {
		client = &http.Client{
			Timeout: opts.Timeout,
		}
	}

	// 记录请求开始时间
	startTime := time.Now()

	// 发送请求（带重试）
	var resp *http.Response
	var respBody []byte
	var lastErr error

	for attempt := 1; attempt <= MaxRetries; attempt++ {
		// 发送请求
		resp, err = client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < MaxRetries {
				slog.Warn("HTTP请求失败，准备重试",
					"error", err,
					"url", urlStr,
					"attempt", attempt,
					"max_retries", MaxRetries)
				time.Sleep(RetryDelay * time.Duration(attempt)) // 递增延迟
				continue
			}
			slog.Error("HTTP请求失败（已达最大重试次数）",
				"error", err,
				"url", urlStr,
				"attempts", MaxRetries,
				"duration", time.Since(startTime))
			return nil, fmt.Errorf("HTTP请求失败: %v", err)
		}

		// 读取响应
		respBody, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		
		if err != nil {
			lastErr = err
			if attempt < MaxRetries {
				slog.Warn("读取响应失败，准备重试",
					"error", err,
					"url", urlStr,
					"attempt", attempt,
					"max_retries", MaxRetries)
				time.Sleep(RetryDelay * time.Duration(attempt))
				continue
			}
			slog.Error("读取响应失败（已达最大重试次数）",
				"error", err,
				"url", urlStr,
				"attempts", MaxRetries,
				"duration", time.Since(startTime))
			return nil, fmt.Errorf("读取响应失败: %v", err)
		}

		// 成功，记录统计信息
		slog.Info("API请求成功",
			"url", urlStr,
			"method", method,
			"status_code", resp.StatusCode,
			"response_size", len(respBody),
			"duration_ms", time.Since(startTime).Milliseconds(),
			"attempt", attempt)
		break
	}

	if lastErr != nil {
		return nil, lastErr
	}

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorMsg := resp.Header.Get("x-ca-error-message")
		requestID := resp.Header.Get("x-ca-request-id")

		slog.Error("API请求失败",
			"method", method,
			"url", urlStr,
			"status_code", resp.StatusCode,
			"request_id", requestID,
			"error_message", errorMsg,
			"response", string(respBody))

		return nil, fmt.Errorf("%s %s 请求失败, 状态码: %d, 请求ID: %s, 错误信息: %s",
			method, urlStr, resp.StatusCode, requestID, errorMsg)
	}

	// 解析响应
	contentType := resp.Header.Get("content-type")
	if strings.HasPrefix(contentType, ContentTypeJSON) {
		var result interface{}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, fmt.Errorf("解析JSON响应失败: %v", err)
		}
		return result, nil
	}

	// 返回原始字符串
	return string(respBody), nil
}

// buildHeaders 构建请求headers
func (c *PSClient) buildHeaders(headers, signHeaders map[string]string) map[string]string {
	result := map[string]string{
		"x-ca-timestamp": fmt.Sprintf("%d", time.Now().UnixMilli()),
		"x-ca-key":       c.appKey,
		"x-ca-nonce":     uuid.New().String(),
		"x-ca-stage":     c.stage,
		"accept":         ContentTypeJSON,
	}

	// 合并自定义headers（转小写）
	for k, v := range headers {
		result[strings.ToLower(k)] = v
	}

	// 合并签名headers（转小写）
	for k, v := range signHeaders {
		result[strings.ToLower(k)] = v
	}

	return result
}

// buildStringToSign 构建待签名字符串
func (c *PSClient) buildStringToSign(method string, headers map[string]string, signedHeadersStr string, parsedURL *url.URL, data interface{}) string {
	var builder strings.Builder

	// HTTP方法
	builder.WriteString(method)
	builder.WriteString("\n")

	// Accept
	if accept := headers["accept"]; accept != "" {
		builder.WriteString(accept)
	}
	builder.WriteString("\n")

	// Content-MD5
	if contentMD5 := headers["content-md5"]; contentMD5 != "" {
		builder.WriteString(contentMD5)
	}
	builder.WriteString("\n")

	// Content-Type
	if contentType := headers["content-type"]; contentType != "" {
		builder.WriteString(contentType)
	}
	builder.WriteString("\n")

	// Date
	if date := headers["date"]; date != "" {
		builder.WriteString(date)
	}
	builder.WriteString("\n")

	// 签名的headers
	if signedHeadersStr != "" {
		builder.WriteString(signedHeadersStr)
		builder.WriteString("\n")
	}

	// URL路径和查询参数
	contentType := headers["content-type"]
	if strings.HasPrefix(contentType, ContentTypeForm) && data != nil {
		// form类型需要将data参数也加入签名
		builder.WriteString(c.buildURLWithFormData(parsedURL, data))
	} else {
		builder.WriteString(c.buildURL(parsedURL))
	}

	return builder.String()
}

// buildURL 构建URL路径和查询参数字符串
func (c *PSClient) buildURL(parsedURL *url.URL) string {
	result := parsedURL.Path

	query := parsedURL.Query()
	if len(query) > 0 {
		// 按key排序
		keys := make([]string, 0, len(query))
		for k := range query {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// 构建查询字符串
		params := make([]string, 0, len(keys))
		for _, k := range keys {
			v := query.Get(k)
			if v != "" {
				params = append(params, fmt.Sprintf("%s=%s", k, v))
			} else {
				params = append(params, k)
			}
		}
		result += "?" + strings.Join(params, "&")
	}

	return result
}

// buildURLWithFormData 构建包含form数据的URL
func (c *PSClient) buildURLWithFormData(parsedURL *url.URL, data interface{}) string {
	allParams := make(map[string]string)

	// 添加查询参数
	query := parsedURL.Query()
	for k := range query {
		allParams[k] = query.Get(k)
	}

	// 添加form数据
	if dataMap, ok := data.(map[string]interface{}); ok {
		for k, v := range dataMap {
			allParams[k] = fmt.Sprintf("%v", v)
		}
	} else if dataMap, ok := data.(map[string]string); ok {
		for k, v := range dataMap {
			allParams[k] = v
		}
	}

	result := parsedURL.Path
	if len(allParams) > 0 {
		// 按key排序
		keys := make([]string, 0, len(allParams))
		for k := range allParams {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// 构建查询字符串
		params := make([]string, 0, len(keys))
		for _, k := range keys {
			v := allParams[k]
			if v != "" {
				params = append(params, fmt.Sprintf("%s=%s", k, v))
			} else {
				params = append(params, k)
			}
		}
		result += "?" + strings.Join(params, "&")
	}

	return result
}

// getSignHeaderKeys 获取需要签名的header keys
func (c *PSClient) getSignHeaderKeys(headers map[string]string, signHeaders map[string]string) []string {
	keySet := make(map[string]bool)

	for k := range headers {
		// x-ca- 开头的header或者在signHeaders中指定的header
		if strings.HasPrefix(k, "x-ca-") {
			keySet[k] = true
		} else if signHeaders != nil {
			if _, exists := signHeaders[k]; exists {
				keySet[k] = true
			}
		}
	}

	// 转换为切片并排序
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

// getSignedHeadersString 获取签名的headers字符串
func (c *PSClient) getSignedHeadersString(signHeaderKeys []string, headers map[string]string) string {
	if len(signHeaderKeys) == 0 {
		return ""
	}

	parts := make([]string, 0, len(signHeaderKeys))
	for _, key := range signHeaderKeys {
		parts = append(parts, fmt.Sprintf("%s:%s", key, headers[key]))
	}

	return strings.Join(parts, "\n")
}

// sign 使用HMAC-SHA256签名
func (c *PSClient) sign(stringToSign string) string {
	h := hmac.New(sha256.New, c.appSecret)
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// md5 计算MD5
func (c *PSClient) md5(data []byte) string {
	hash := md5.Sum(data)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// getContentType 获取Content-Type
func (c *PSClient) getContentType(headers map[string]string) string {
	if ct := headers["content-type"]; ct != "" {
		return ct
	}
	if ct := headers["Content-Type"]; ct != "" {
		return ct
	}
	return ""
}

// encodeBody 编码请求体
func (c *PSClient) encodeBody(data interface{}, contentType string) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	// 如果已经是字节数组或字符串，直接返回
	if b, ok := data.([]byte); ok {
		return b, nil
	}
	if s, ok := data.(string); ok {
		return []byte(s), nil
	}

	// 根据Content-Type编码
	if strings.HasPrefix(contentType, ContentTypeForm) {
		// application/x-www-form-urlencoded
		values := url.Values{}
		if dataMap, ok := data.(map[string]interface{}); ok {
			for k, v := range dataMap {
				values.Set(k, fmt.Sprintf("%v", v))
			}
		} else if dataMap, ok := data.(map[string]string); ok {
			for k, v := range dataMap {
				values.Set(k, v)
			}
		} else {
			return nil, fmt.Errorf("form类型数据必须是map[string]interface{}或map[string]string")
		}
		return []byte(values.Encode()), nil
	}

	// 默认使用JSON编码
	return json.Marshal(data)
}
