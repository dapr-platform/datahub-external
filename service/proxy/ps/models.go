package ps

// FamilyRequest 家族信息请求参数
type FamilyRequest struct {
	PageNum  int    `json:"pageNum"`  // 页编号
	PageSize int    `json:"pageSize"` // 页大小（最大2000）
	DS       string `json:"ds"`       // 分区字段（示例：20251013，格式：YYYYMMDD）
	EmplID   string `json:"emplid"`   // 员工ID（可选）
}

// FamilyResponse 家族信息响应
type FamilyResponse struct {
	CFamilyID         int     `json:"c_family_id"`          // 家族编号
	EmplID            string  `json:"emplid"`               // 家族长工号
	CFmyRelationship  string  `json:"c_fmy_relationship"`   // 家族关系
	CFamilyName       string  `json:"c_family_name"`        // 家族番号
	Date1             string  `json:"date1"`                // 成立日期
	EffStatus         string  `json:"eff_status"`           // 生效状态
	HrsRowAddDttm     *string `json:"hrs_row_add_dttm"`     // 添加时间（可能为null）
	HrsRowAddOprid    string  `json:"hrs_row_add_oprid"`    // 添加用户ID
	HrsRowUpdDttm     *string `json:"hrs_row_upd_dttm"`     // 数据修改时间（可能为null）
	HrsRowUpdOprid    string  `json:"hrs_row_upd_oprid"`    // 更新操作人
	AdbDate           string  `json:"adb_date"`             // 业务数据变动
}

// PSAPIResponse PS API统一响应格式
type PSAPIResponse struct {
	ErrCode   int         `json:"errCode"`   // 错误码（0表示成功）
	ErrMsg    string      `json:"errMsg"`    // 错误消息
	Data      interface{} `json:"data"`      // 响应数据
	RequestID string      `json:"requestId"` // 请求ID
	APILog    interface{} `json:"apiLog"`    // API日志
}

// PSPaginationResponse PS API分页响应
type PSPaginationResponse struct {
	ErrCode   int              `json:"errCode"`   // 错误码（0表示成功）
	ErrMsg    string           `json:"errMsg"`    // 错误消息
	Data      PSPaginationData `json:"data"`      // 分页数据
	RequestID string           `json:"requestId"` // 请求ID
	APILog    interface{}      `json:"apiLog"`    // API日志
}

// PSPaginationData 分页数据
type PSPaginationData struct {
	TotalNum int              `json:"totalNum"` // 总记录数
	PageNum  int              `json:"pageNum"`  // 当前页
	PageSize int              `json:"pageSize"` // 页大小
	Rows     []FamilyResponse `json:"rows"`     // 记录列表
}

