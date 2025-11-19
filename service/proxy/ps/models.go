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

// PSPaginationData 分页数据（家族成员）
type PSPaginationData struct {
	TotalNum int              `json:"totalNum"` // 总记录数
	PageNum  int              `json:"pageNum"`  // 当前页
	PageSize int              `json:"pageSize"` // 页大小
	Rows     []FamilyResponse `json:"rows"`     // 记录列表
}

// PositionResponse 岗位信息响应
type PositionResponse struct {
	PositionNbr      string  `json:"position_nbr"`       // 职位编码
	PosnDescr        string  `json:"posn_descr"`         // 职位描述
	BusinessUnit     string  `json:"business_unit"`      // 单位名称
	BusinessDescr    string  `json:"business_descr"`     // 单位描述
	DeptID           string  `json:"deptid"`             // 部门ID
	DeptDescr        string  `json:"dept_descr"`         // 部门描述
	CDeptType        string  `json:"c_dept_type"`        // 部门类型
	Jobcode          string  `json:"jobcode"`            // 岗位编码
	JobcodeDescr     string  `json:"jobcode_descr"`      // 岗位名称
	CCountryDesc     string  `json:"c_country_desc"`     // 国家
	CStateDesc       string  `json:"c_state_desc"`       // 省份
	CCityDesc        string  `json:"c_city_desc"`        // 城市
	Address1         string  `json:"address1"`           // 所在地点
	ReportsTo        string  `json:"reports_to"`         // 汇报岗位ID
	Descr100         string  `json:"descr100"`           // 汇报岗位名称
	EmplID           string  `json:"emplid"`             // 员工编码
	Name             string  `json:"name"`               // 姓名
	CSupvLvl         string  `json:"c_supv_lvl"`         // 员工级别
	CSupvLvlDesc     string  `json:"c_supv_lvl_desc"`    // 员工级别描述
	BgnDt            string  `json:"bgn_dt"`             // 开始日期
	EndDt            *string `json:"end_dt"`             // 结束日期
	Descr1           string  `json:"descr1"`             // 岗位级别
	Descr2           string  `json:"descr2"`             // 级别描述
	HrsRowAddDttm    string  `json:"hrs_row_add_dttm"`   // 添加时间
	HrsRowAddOprid   string  `json:"hrs_row_add_oprid"`  // 添加人ID
	HrsRowUpdDttm    string  `json:"hrs_row_upd_dttm"`   // 修改时间
	HrsRowUpdOprid   string  `json:"hrs_row_upd_oprid"`  // 修改人ID
	COldPosnNbr      string  `json:"c_old_posn_nbr"`     // 旧岗位编码
	Createdttm       string  `json:"createdttm"`         // 创建时间
	CIntFlag         *string `json:"c_int_flag"`         // 标识（增量接口）
	Effdt            *string `json:"effdt"`              // 生效时间（增量接口）
	AdbDate          string  `json:"adb_date"`           // 业务数据变动时间
	CSequenceID      string  `json:"C_SEQUENCE_ID"`      // 序列代码
	CSequenceDescr   string  `json:"C_SEQUENCE_DESCR"`   // 序列描述
	CSubsequenceID   string  `json:"C_SUBSEQUENCE_ID"`   // 子序列代码
	CSubsequenceDesc string  `json:"C_SUBSEQUENCE_DESC"` // 子序描述
}

// OrganizationResponse 组织信息响应
type OrganizationResponse struct {
	Setid            string  `json:"setid"`              // 集合ID
	DeptID           string  `json:"deptid"`             // 部门ID
	TreeLevelNum     int64   `json:"tree_level_num"`     // 组织层级
	TreeNodeNum      int64   `json:"tree_node_num"`      // 组织节点
	DeptDescr        string  `json:"dept_descr"`         // 部门描述
	CDeptType        string  `json:"c_dept_type"`        // 部门类型
	CDeptTypeDescr   string  `json:"c_dept_type_descr"`  // 部门类型描述
	CCountryDesc     string  `json:"c_country_desc"`     // 国家
	CStateDesc       string  `json:"c_state_desc"`       // 省份
	CCityDesc        string  `json:"c_city_desc"`        // 城市
	Address1         string  `json:"address1"`           // 工作地址
	Parentname       string  `json:"parentname"`         // 父节点名称
	CDeptDescr       string  `json:"c_dept_descr"`       // 部门名称
	EmplID           string  `json:"emplid"`             // 员工ID
	Name             string  `json:"name"`               // 姓名
	BgnDt            string  `json:"bgn_dt"`             // 开始日期
	EndDt            *string `json:"end_dt"`             // 结束日期
	HrsRowAddDttm    string  `json:"hrs_row_add_dttm"`   // 添加时间
	HrsRowAddOprid   string  `json:"hrs_row_add_oprid"`  // 添加人ID
	HrsRowUpdDttm    string  `json:"hrs_row_upd_dttm"`   // 修改时间
	HrsRowUpdOprid   string  `json:"hrs_row_upd_oprid"`  // 修改人ID
	COldDeptID       string  `json:"c_old_deptid"`       // 旧部门ID
	CIntFlag         *string `json:"c_int_flag"`         // 标识（增量接口）
	Effdt            *string `json:"effdt"`              // 生效时间（增量接口）
	Createdttm       string  `json:"createdttm"`         // 创建时间
	CCompany         string  `json:"c_company"`          // 公司ID
	CCompanyDesc     string  `json:"c_company_desc"`     // 公司名称
	AdbDate          string  `json:"adb_date"`           // 业务数据变动时间
}

// EmployeeResponse 员工信息响应
type EmployeeResponse struct {
	HrsRowAddDttm      string  `json:"hrs_row_add_dttm"`       // 插入数据时间
	HrsRowAddOprid     string  `json:"hrs_row_add_oprid"`      // 插入数据ID
	HrsRowUpdDttm      string  `json:"hrs_row_upd_dttm"`       // 数据修改时间
	HrsRowUpdOprid     string  `json:"hrs_row_upd_oprid"`      // 数据修改操作人ID
	Createdttm         string  `json:"createdttm"`             // 创建时间
	EmplID             string  `json:"emplid"`                 // 职员编码
	Name               string  `json:"name"`                   // 员工姓名
	Sex                string  `json:"sex"`                    // 性别编码
	CSexDescr          string  `json:"c_sex_descr"`            // 性别
	MarStatus          string  `json:"mar_status"`             // 婚姻状态编码
	CMarDescr          string  `json:"c_mar_descr"`            // 婚姻状态
	Birthdate          string  `json:"birthdate"`              // 出生日期（身份证号截取）
	CBirthType         string  `json:"c_birth_type"`           // 生日类型编码
	CBirtypeDescr      string  `json:"c_birtype_descr"`        // 生日类型（阴历/阳历）
	CBirthdate         string  `json:"c_birthdate"`            // 生日
	CBirthdate1        string  `json:"c_birthdate1"`           // 今年生日是哪一天
	CCountryDesc1      string  `json:"c_country_desc1"`        // 国家编码
	CCountryDesc2      string  `json:"c_country_desc2"`        // 国家
	CStateDesc1        string  `json:"c_state_desc1"`          // 籍贯编码
	CStateDesc2        string  `json:"c_state_desc2"`          // 籍贯
	CEthnicGrpDesc     string  `json:"c_ethnic_grp_desc"`      // 民族拼写
	CEthnicGrpDesc1    string  `json:"c_ethnic_grp_desc1"`     // 民族
	CPersPolityDesc    string  `json:"c_pers_polity_desc"`     // 政治面貌
	Address1           string  `json:"address1"`               // 现居地
	Address3           string  `json:"address3"`               // 身份证件地址
	CNidTypeDesc       string  `json:"c_nid_type_desc"`        // 证件类型
	NationalID         string  `json:"national_id"`            // 证件号
	Phone              string  `json:"phone"`                  // 联系电话
	EmailAddr          string  `json:"email_addr"`             // 邮箱
	EducationLvlAchv   string  `json:"education_lvl_achv"`     // 学历编码
	CEducationDescr    string  `json:"c_education_descr"`      // 学历
	SchoolDescr        string  `json:"school_descr"`           // 毕业学校
	MajorDescr         string  `json:"major_descr"`            // 专业
	CHireDate          string  `json:"c_hire_date"`            // 入职日期
	CLeaveDate         string  `json:"c_leave_date"`           // 离职日期
	RehireDt           string  `json:"rehire_dt"`              // 再次入职日期
	EmplClass          string  `json:"empl_class"`             // 员工类型编码
	CEmplclsDescr      string  `json:"c_emplcls_descr"`        // 员工类型描述
	BusinessDescr      string  `json:"business_descr"`         // 业务单位名称
	DeptID             string  `json:"deptid"`                 // 部门编码
	TreeNodeNum        int64   `json:"tree_node_num"`          // 组织节点编码
	DeptDescr          string  `json:"dept_descr"`             // 部门名称
	Jobcode            string  `json:"jobcode"`                // 岗位编码
	JobcodeDescr       string  `json:"jobcode_descr"`          // 岗位名称
	PositionNbr        string  `json:"position_nbr"`           // 职位编码
	PosnDescr          string  `json:"posn_descr"`             // 职位描述
	CSupvLvl           string  `json:"c_supv_lvl"`             // 岗位级别
	CSupvLvlDesc       string  `json:"c_supv_lvl_desc"`        // 岗位描述
	CCountryDesc       string  `json:"c_country_desc"`         // 国家
	CStateDesc         string  `json:"c_state_desc"`           // 省（市）
	CCityDesc          string  `json:"c_city_desc"`            // 城市
	Address2           string  `json:"address2"`               // 办公地址
	CIntFlag           *string `json:"c_int_flag"`             // 初始化标识（增量接口）
	AdbDate            string  `json:"adb_date"`               // 业务数据变动时间
}

// EmployeeHonorResponse 员工荣誉信息响应
type EmployeeHonorResponse struct {
	EmplID           string  `json:"emplid"`             // 工号
	BeginDt          string  `json:"begin_dt"`           // 开始日期
	EffStatus        string  `json:"eff_status"`         // 荣誉状态
	EndDt            *string `json:"end_dt"`             // 失效日期
	Descr254         string  `json:"descr254"`           // 荣誉原因
	Descr254a        string  `json:"descr254a"`          // 失效原因
	Comments256      string  `json:"comments_256"`       // 备注信息(256字符)
	DeptID           string  `json:"deptid"`             // 发文部门ID
	DeptDescr        string  `json:"dept_descr"`         // 发文部门描述
	HrsRowAddDttm    string  `json:"hrs_row_add_dttm"`   // 添加时间
	HrsRowAddOprid   string  `json:"hrs_row_add_oprid"`  // 添加用户ID
	HrsRowUpdDttm    string  `json:"hrs_row_upd_dttm"`   // 数据修改时间
	HrsRowUpdOprid   string  `json:"hrs_row_upd_oprid"`  // 更新操作人
	AdbDate          string  `json:"adb_date"`           // 业务数据变动时间
}

// FamilyMainResponse 家族父表信息响应
type FamilyMainResponse struct {
	CFamilyID        int     `json:"c_family_id"`        // 家族编号
	EmplID           string  `json:"emplid"`             // 家族长工号
	CFamilyName      string  `json:"c_family_name"`      // 家族番号
	Date1            string  `json:"date1"`              // 成立日期
	EffStatus        string  `json:"eff_status"`         // 生效状态
	EffdtTo          *string `json:"effdt_to"`           // 失效日期
	HrsRowAddDttm    string  `json:"hrs_row_add_dttm"`   // 添加时间
	HrsRowAddOprid   string  `json:"hrs_row_add_oprid"`  // 添加用户ID
	HrsRowUpdDttm    string  `json:"hrs_row_upd_dttm"`   // 数据修改时间
	HrsRowUpdOprid   string  `json:"hrs_row_upd_oprid"`  // 更新操作人
	AdbDate          string  `json:"adb_date"`           // 业务数据变动时间
}

// PositionPaginationData 岗位分页数据
type PositionPaginationData struct {
	TotalNum int                `json:"totalNum"` // 总记录数
	PageNum  int                `json:"pageNum"`  // 当前页
	PageSize int                `json:"pageSize"` // 页大小
	Rows     []PositionResponse `json:"rows"`     // 记录列表
}

// OrganizationPaginationData 组织分页数据
type OrganizationPaginationData struct {
	TotalNum int                      `json:"totalNum"` // 总记录数
	PageNum  int                      `json:"pageNum"`  // 当前页
	PageSize int                      `json:"pageSize"` // 页大小
	Rows     []OrganizationResponse   `json:"rows"`     // 记录列表
}

// EmployeePaginationData 员工分页数据
type EmployeePaginationData struct {
	TotalNum int                `json:"totalNum"` // 总记录数
	PageNum  int                `json:"pageNum"`  // 当前页
	PageSize int                `json:"pageSize"` // 页大小
	Rows     []EmployeeResponse `json:"rows"`     // 记录列表
}

// EmployeeHonorPaginationData 员工荣誉分页数据
type EmployeeHonorPaginationData struct {
	TotalNum int                     `json:"totalNum"` // 总记录数
	PageNum  int                     `json:"pageNum"`  // 当前页
	PageSize int                     `json:"pageSize"` // 页大小
	Rows     []EmployeeHonorResponse `json:"rows"`     // 记录列表
}

// FamilyMainPaginationData 家族父表分页数据
type FamilyMainPaginationData struct {
	TotalNum int                  `json:"totalNum"` // 总记录数
	PageNum  int                  `json:"pageNum"`  // 当前页
	PageSize int                  `json:"pageSize"` // 页大小
	Rows     []FamilyMainResponse `json:"rows"`     // 记录列表
}

// PositionPaginationResponse 岗位分页响应
type PositionPaginationResponse struct {
	ErrCode   int                    `json:"errCode"`   // 错误码（0表示成功）
	ErrMsg    string                 `json:"errMsg"`    // 错误消息
	Data      PositionPaginationData `json:"data"`      // 分页数据
	RequestID string                 `json:"requestId"` // 请求ID
	APILog    interface{}            `json:"apiLog"`    // API日志
}

// OrganizationPaginationResponse 组织分页响应
type OrganizationPaginationResponse struct {
	ErrCode   int                        `json:"errCode"`   // 错误码（0表示成功）
	ErrMsg    string                     `json:"errMsg"`    // 错误消息
	Data      OrganizationPaginationData `json:"data"`      // 分页数据
	RequestID string                     `json:"requestId"` // 请求ID
	APILog    interface{}                `json:"apiLog"`    // API日志
}

// EmployeePaginationResponse 员工分页响应
type EmployeePaginationResponse struct {
	ErrCode   int                    `json:"errCode"`   // 错误码（0表示成功）
	ErrMsg    string                 `json:"errMsg"`    // 错误消息
	Data      EmployeePaginationData `json:"data"`      // 分页数据
	RequestID string                 `json:"requestId"` // 请求ID
	APILog    interface{}            `json:"apiLog"`    // API日志
}

// EmployeeHonorPaginationResponse 员工荣誉分页响应
type EmployeeHonorPaginationResponse struct {
	ErrCode   int                         `json:"errCode"`   // 错误码（0表示成功）
	ErrMsg    string                      `json:"errMsg"`    // 错误消息
	Data      EmployeeHonorPaginationData `json:"data"`      // 分页数据
	RequestID string                      `json:"requestId"` // 请求ID
	APILog    interface{}                 `json:"apiLog"`    // API日志
}

// FamilyMainPaginationResponse 家族父表分页响应
type FamilyMainPaginationResponse struct {
	ErrCode   int                      `json:"errCode"`   // 错误码（0表示成功）
	ErrMsg    string                   `json:"errMsg"`    // 错误消息
	Data      FamilyMainPaginationData `json:"data"`      // 分页数据
	RequestID string                   `json:"requestId"` // 请求ID
	APILog    interface{}              `json:"apiLog"`    // API日志
}

