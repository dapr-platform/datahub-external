package ps

import (
	"time"
)

// PSFamilyMember PS家族成员信息表
type PSFamilyMember struct {
	ID               uint       `gorm:"primarykey" json:"id"`
	CFamilyID        int        `gorm:"column:c_family_id;type:int;index" json:"c_family_id"`                 // 家族编号
	EmplID           string     `gorm:"column:emplid;type:varchar(50);index" json:"emplid"`                   // 家族长工号
	CFmyRelationship string     `gorm:"column:c_fmy_relationship;type:varchar(50)" json:"c_fmy_relationship"` // 家族关系
	CFamilyName      string     `gorm:"column:c_family_name;type:varchar(200)" json:"c_family_name"`          // 家族番号
	Date1            *time.Time `gorm:"column:date1;type:timestamp" json:"date1"`                             // 成立日期
	EffStatus        string     `gorm:"column:eff_status;type:varchar(20)" json:"eff_status"`                 // 生效状态
	HrsRowAddDttm    *time.Time `gorm:"column:hrs_row_add_dttm;type:timestamp" json:"hrs_row_add_dttm"`       // 添加时间（可为null）
	HrsRowAddOprid   string     `gorm:"column:hrs_row_add_oprid;type:varchar(50)" json:"hrs_row_add_oprid"`   // 添加用户ID
	HrsRowUpdDttm    *time.Time `gorm:"column:hrs_row_upd_dttm;type:timestamp" json:"hrs_row_upd_dttm"`       // 数据修改时间（可为null）
	HrsRowUpdOprid   string     `gorm:"column:hrs_row_upd_oprid;type:varchar(50)" json:"hrs_row_upd_oprid"`   // 更新操作人
	AdbDate          string     `gorm:"column:adb_date;type:varchar(20)" json:"adb_date"`                     // 业务数据变动
	DS               string     `gorm:"column:ds;type:varchar(20);index" json:"ds"`                           // 分区字段（YYYYMMDD）
	SyncTime         time.Time  `gorm:"column:sync_time;type:timestamp;index" json:"sync_time"`               // 同步时间
	UniqueKey        string     `gorm:"column:unique_key;type:varchar(200);uniqueIndex" json:"unique_key"`    // 唯一键（c_family_id_emplid_ds）
	CreatedAt        time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"`    // 创建时间
	UpdatedAt        time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at"`    // 更新时间
}

// TableName 指定表名
func (PSFamilyMember) TableName() string {
	return "ps_family_members"
}

// PSPosition PS岗位信息表
type PSPosition struct {
	ID               uint       `gorm:"primarykey" json:"id"`
	PositionNbr      string     `gorm:"column:position_nbr;type:varchar(50);index" json:"position_nbr"`     // 职位编码
	PosnDescr        string     `gorm:"column:posn_descr;type:varchar(200)" json:"posn_descr"`             // 职位描述
	BusinessUnit     string     `gorm:"column:business_unit;type:varchar(50)" json:"business_unit"`        // 单位名称
	BusinessDescr    string     `gorm:"column:business_descr;type:varchar(200)" json:"business_descr"`     // 单位描述
	DeptID           string     `gorm:"column:deptid;type:varchar(50);index" json:"deptid"`                // 部门ID
	DeptDescr        string     `gorm:"column:dept_descr;type:varchar(200)" json:"dept_descr"`             // 部门描述
	CDeptType        string     `gorm:"column:c_dept_type;type:varchar(20)" json:"c_dept_type"`            // 部门类型
	Jobcode          string     `gorm:"column:jobcode;type:varchar(50)" json:"jobcode"`                    // 岗位编码
	JobcodeDescr     string     `gorm:"column:jobcode_descr;type:varchar(200)" json:"jobcode_descr"`       // 岗位名称
	CCountryDesc     string     `gorm:"column:c_country_desc;type:varchar(100)" json:"c_country_desc"`     // 国家
	CStateDesc       string     `gorm:"column:c_state_desc;type:varchar(100)" json:"c_state_desc"`         // 省份
	CCityDesc        string     `gorm:"column:c_city_desc;type:varchar(100)" json:"c_city_desc"`           // 城市
	Address1         string     `gorm:"column:address1;type:varchar(500)" json:"address1"`                 // 所在地点
	ReportsTo        string     `gorm:"column:reports_to;type:varchar(50)" json:"reports_to"`              // 汇报岗位ID
	Descr100         string     `gorm:"column:descr100;type:varchar(200)" json:"descr100"`                 // 汇报岗位名称
	EmplID           string     `gorm:"column:emplid;type:varchar(50);index" json:"emplid"`                // 员工编码
	Name             string     `gorm:"column:name;type:varchar(100)" json:"name"`                         // 姓名
	CSupvLvl         string     `gorm:"column:c_supv_lvl;type:varchar(20)" json:"c_supv_lvl"`              // 员工级别
	CSupvLvlDesc     string     `gorm:"column:c_supv_lvl_desc;type:varchar(100)" json:"c_supv_lvl_desc"`   // 员工级别描述
	BgnDt            *time.Time `gorm:"column:bgn_dt;type:timestamp" json:"bgn_dt"`                        // 开始日期
	EndDt            *time.Time `gorm:"column:end_dt;type:timestamp" json:"end_dt"`                        // 结束日期
	Descr1           string     `gorm:"column:descr1;type:varchar(100)" json:"descr1"`                     // 岗位级别
	Descr2           string     `gorm:"column:descr2;type:varchar(100)" json:"descr2"`                     // 级别描述
	HrsRowAddDttm    *time.Time `gorm:"column:hrs_row_add_dttm;type:timestamp" json:"hrs_row_add_dttm"`    // 添加时间
	HrsRowAddOprid   string     `gorm:"column:hrs_row_add_oprid;type:varchar(50)" json:"hrs_row_add_oprid"` // 添加人ID
	HrsRowUpdDttm    *time.Time `gorm:"column:hrs_row_upd_dttm;type:timestamp" json:"hrs_row_upd_dttm"`    // 修改时间
	HrsRowUpdOprid   string     `gorm:"column:hrs_row_upd_oprid;type:varchar(50)" json:"hrs_row_upd_oprid"` // 修改人ID
	COldPosnNbr      string     `gorm:"column:c_old_posn_nbr;type:varchar(50)" json:"c_old_posn_nbr"`      // 旧岗位编码
	Createdttm       *time.Time `gorm:"column:createdttm;type:timestamp" json:"createdttm"`                // 创建时间
	CIntFlag         string     `gorm:"column:c_int_flag;type:varchar(20)" json:"c_int_flag"`              // 标识
	Effdt            *time.Time `gorm:"column:effdt;type:timestamp" json:"effdt"`                          // 生效时间
	AdbDate          string     `gorm:"column:adb_date;type:varchar(20)" json:"adb_date"`                  // 业务数据变动时间
	CSequenceID      string     `gorm:"column:c_sequence_id;type:varchar(50)" json:"c_sequence_id"`        // 序列代码
	CSequenceDescr   string     `gorm:"column:c_sequence_descr;type:varchar(200)" json:"c_sequence_descr"` // 序列描述
	CSubsequenceID   string     `gorm:"column:c_subsequence_id;type:varchar(50)" json:"c_subsequence_id"`  // 子序列代码
	CSubsequenceDesc string     `gorm:"column:c_subsequence_desc;type:varchar(200)" json:"c_subsequence_desc"` // 子序描述
	DS               string     `gorm:"column:ds;type:varchar(20);index" json:"ds"`                        // 分区字段
	SyncTime         time.Time  `gorm:"column:sync_time;type:timestamp;index" json:"sync_time"`            // 同步时间
	UniqueKey        string     `gorm:"column:unique_key;type:varchar(200);uniqueIndex" json:"unique_key"` // 唯一键（position_nbr_ds）
	CreatedAt        time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt        time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (PSPosition) TableName() string {
	return "ps_positions"
}

// PSOrganization PS组织信息表
type PSOrganization struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	Setid          string     `gorm:"column:setid;type:varchar(50)" json:"setid"`                        // 集合ID
	DeptID         string     `gorm:"column:deptid;type:varchar(50);index" json:"deptid"`                // 部门ID
	TreeLevelNum   int64      `gorm:"column:tree_level_num;type:bigint" json:"tree_level_num"`           // 组织层级
	TreeNodeNum    int64      `gorm:"column:tree_node_num;type:bigint;index" json:"tree_node_num"`       // 组织节点
	DeptDescr      string     `gorm:"column:dept_descr;type:varchar(200)" json:"dept_descr"`             // 部门描述
	CDeptType      string     `gorm:"column:c_dept_type;type:varchar(20)" json:"c_dept_type"`            // 部门类型
	CDeptTypeDescr string     `gorm:"column:c_dept_type_descr;type:varchar(100)" json:"c_dept_type_descr"` // 部门类型描述
	CCountryDesc   string     `gorm:"column:c_country_desc;type:varchar(100)" json:"c_country_desc"`     // 国家
	CStateDesc     string     `gorm:"column:c_state_desc;type:varchar(100)" json:"c_state_desc"`         // 省份
	CCityDesc      string     `gorm:"column:c_city_desc;type:varchar(100)" json:"c_city_desc"`           // 城市
	Address1       string     `gorm:"column:address1;type:varchar(500)" json:"address1"`                 // 工作地址
	Parentname     string     `gorm:"column:parentname;type:varchar(50)" json:"parentname"`              // 父节点名称
	CDeptDescr     string     `gorm:"column:c_dept_descr;type:varchar(200)" json:"c_dept_descr"`         // 部门名称
	EmplID         string     `gorm:"column:emplid;type:varchar(50);index" json:"emplid"`                // 员工ID
	Name           string     `gorm:"column:name;type:varchar(100)" json:"name"`                         // 姓名
	BgnDt          *time.Time `gorm:"column:bgn_dt;type:timestamp" json:"bgn_dt"`                        // 开始日期
	EndDt          *time.Time `gorm:"column:end_dt;type:timestamp" json:"end_dt"`                        // 结束日期
	HrsRowAddDttm  *time.Time `gorm:"column:hrs_row_add_dttm;type:timestamp" json:"hrs_row_add_dttm"`    // 添加时间
	HrsRowAddOprid string     `gorm:"column:hrs_row_add_oprid;type:varchar(50)" json:"hrs_row_add_oprid"` // 添加人ID
	HrsRowUpdDttm  *time.Time `gorm:"column:hrs_row_upd_dttm;type:timestamp" json:"hrs_row_upd_dttm"`    // 修改时间
	HrsRowUpdOprid string     `gorm:"column:hrs_row_upd_oprid;type:varchar(50)" json:"hrs_row_upd_oprid"` // 修改人ID
	COldDeptID     string     `gorm:"column:c_old_deptid;type:varchar(50)" json:"c_old_deptid"`          // 旧部门ID
	CIntFlag       string     `gorm:"column:c_int_flag;type:varchar(20)" json:"c_int_flag"`              // 标识
	Effdt          *time.Time `gorm:"column:effdt;type:timestamp" json:"effdt"`                          // 生效时间
	Createdttm     *time.Time `gorm:"column:createdttm;type:timestamp" json:"createdttm"`                // 创建时间
	CCompany       string     `gorm:"column:c_company;type:varchar(50)" json:"c_company"`                // 公司ID
	CCompanyDesc   string     `gorm:"column:c_company_desc;type:varchar(200)" json:"c_company_desc"`     // 公司名称
	AdbDate        string     `gorm:"column:adb_date;type:varchar(20)" json:"adb_date"`                  // 业务数据变动时间
	DS             string     `gorm:"column:ds;type:varchar(20);index" json:"ds"`                        // 分区字段
	SyncTime       time.Time  `gorm:"column:sync_time;type:timestamp;index" json:"sync_time"`            // 同步时间
	UniqueKey      string     `gorm:"column:unique_key;type:varchar(200);uniqueIndex" json:"unique_key"` // 唯一键（deptid_ds）
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (PSOrganization) TableName() string {
	return "ps_organizations"
}

// PSEmployee PS员工信息表
type PSEmployee struct {
	ID                 uint       `gorm:"primarykey" json:"id"`
	HrsRowAddDttm      *time.Time `gorm:"column:hrs_row_add_dttm;type:timestamp" json:"hrs_row_add_dttm"`    // 插入数据时间
	HrsRowAddOprid     string     `gorm:"column:hrs_row_add_oprid;type:varchar(50)" json:"hrs_row_add_oprid"` // 插入数据ID
	HrsRowUpdDttm      *time.Time `gorm:"column:hrs_row_upd_dttm;type:timestamp" json:"hrs_row_upd_dttm"`    // 数据修改时间
	HrsRowUpdOprid     string     `gorm:"column:hrs_row_upd_oprid;type:varchar(50)" json:"hrs_row_upd_oprid"` // 数据修改操作人ID
	Createdttm         *time.Time `gorm:"column:createdttm;type:timestamp" json:"createdttm"`                // 创建时间
	EmplID             string     `gorm:"column:emplid;type:varchar(50);index" json:"emplid"`                // 职员编码
	Name               string     `gorm:"column:name;type:varchar(100)" json:"name"`                         // 员工姓名
	Sex                string     `gorm:"column:sex;type:varchar(10)" json:"sex"`                            // 性别编码
	CSexDescr          string     `gorm:"column:c_sex_descr;type:varchar(20)" json:"c_sex_descr"`            // 性别
	MarStatus          string     `gorm:"column:mar_status;type:varchar(10)" json:"mar_status"`              // 婚姻状态编码
	CMarDescr          string     `gorm:"column:c_mar_descr;type:varchar(50)" json:"c_mar_descr"`            // 婚姻状态
	Birthdate          string     `gorm:"column:birthdate;type:varchar(50)" json:"birthdate"`                // 出生日期
	CBirthType         string     `gorm:"column:c_birth_type;type:varchar(10)" json:"c_birth_type"`          // 生日类型编码
	CBirtypeDescr      string     `gorm:"column:c_birtype_descr;type:varchar(50)" json:"c_birtype_descr"`    // 生日类型
	CBirthdate         string     `gorm:"column:c_birthdate;type:varchar(50)" json:"c_birthdate"`            // 生日
	CBirthdate1        string     `gorm:"column:c_birthdate1;type:varchar(50)" json:"c_birthdate1"`          // 今年生日
	CCountryDesc1      string     `gorm:"column:c_country_desc1;type:varchar(100)" json:"c_country_desc1"`   // 国家编码
	CCountryDesc2      string     `gorm:"column:c_country_desc2;type:varchar(100)" json:"c_country_desc2"`   // 国家
	CStateDesc1        string     `gorm:"column:c_state_desc1;type:varchar(100)" json:"c_state_desc1"`       // 籍贯编码
	CStateDesc2        string     `gorm:"column:c_state_desc2;type:varchar(100)" json:"c_state_desc2"`       // 籍贯
	CEthnicGrpDesc     string     `gorm:"column:c_ethnic_grp_desc;type:varchar(100)" json:"c_ethnic_grp_desc"` // 民族拼写
	CEthnicGrpDesc1    string     `gorm:"column:c_ethnic_grp_desc1;type:varchar(100)" json:"c_ethnic_grp_desc1"` // 民族
	CPersPolityDesc    string     `gorm:"column:c_pers_polity_desc;type:varchar(100)" json:"c_pers_polity_desc"` // 政治面貌
	Address1           string     `gorm:"column:address1;type:varchar(500)" json:"address1"`                 // 现居地
	Address3           string     `gorm:"column:address3;type:varchar(500)" json:"address3"`                 // 身份证件地址
	CNidTypeDesc       string     `gorm:"column:c_nid_type_desc;type:varchar(50)" json:"c_nid_type_desc"`    // 证件类型
	NationalID         string     `gorm:"column:national_id;type:varchar(100)" json:"national_id"`           // 证件号
	Phone              string     `gorm:"column:phone;type:varchar(50)" json:"phone"`                        // 联系电话
	EmailAddr          string     `gorm:"column:email_addr;type:varchar(200)" json:"email_addr"`             // 邮箱
	EducationLvlAchv   string     `gorm:"column:education_lvl_achv;type:varchar(50)" json:"education_lvl_achv"` // 学历编码
	CEducationDescr    string     `gorm:"column:c_education_descr;type:varchar(100)" json:"c_education_descr"` // 学历
	SchoolDescr        string     `gorm:"column:school_descr;type:varchar(200)" json:"school_descr"`         // 毕业学校
	MajorDescr         string     `gorm:"column:major_descr;type:varchar(200)" json:"major_descr"`           // 专业
	CHireDate          string     `gorm:"column:c_hire_date;type:varchar(50)" json:"c_hire_date"`            // 入职日期
	CLeaveDate         string     `gorm:"column:c_leave_date;type:varchar(50)" json:"c_leave_date"`          // 离职日期
	RehireDt           string     `gorm:"column:rehire_dt;type:varchar(50)" json:"rehire_dt"`                // 再次入职日期
	EmplClass          string     `gorm:"column:empl_class;type:varchar(50)" json:"empl_class"`              // 员工类型编码
	CEmplclsDescr      string     `gorm:"column:c_emplcls_descr;type:varchar(100)" json:"c_emplcls_descr"`   // 员工类型描述
	BusinessDescr      string     `gorm:"column:business_descr;type:varchar(200)" json:"business_descr"`     // 业务单位名称
	DeptID             string     `gorm:"column:deptid;type:varchar(50);index" json:"deptid"`                // 部门编码
	TreeNodeNum        int64      `gorm:"column:tree_node_num;type:bigint;index" json:"tree_node_num"`       // 组织节点编码
	DeptDescr          string     `gorm:"column:dept_descr;type:varchar(200)" json:"dept_descr"`             // 部门名称
	Jobcode            string     `gorm:"column:jobcode;type:varchar(50)" json:"jobcode"`                    // 岗位编码
	JobcodeDescr       string     `gorm:"column:jobcode_descr;type:varchar(200)" json:"jobcode_descr"`       // 岗位名称
	PositionNbr        string     `gorm:"column:position_nbr;type:varchar(50);index" json:"position_nbr"`    // 职位编码
	PosnDescr          string     `gorm:"column:posn_descr;type:varchar(200)" json:"posn_descr"`             // 职位描述
	CSupvLvl           string     `gorm:"column:c_supv_lvl;type:varchar(20)" json:"c_supv_lvl"`              // 岗位级别
	CSupvLvlDesc       string     `gorm:"column:c_supv_lvl_desc;type:varchar(100)" json:"c_supv_lvl_desc"`   // 岗位描述
	CCountryDesc       string     `gorm:"column:c_country_desc;type:varchar(100)" json:"c_country_desc"`     // 国家
	CStateDesc         string     `gorm:"column:c_state_desc;type:varchar(100)" json:"c_state_desc"`         // 省（市）
	CCityDesc          string     `gorm:"column:c_city_desc;type:varchar(100)" json:"c_city_desc"`           // 城市
	Address2           string     `gorm:"column:address2;type:varchar(500)" json:"address2"`                 // 办公地址
	CIntFlag           string     `gorm:"column:c_int_flag;type:varchar(20)" json:"c_int_flag"`              // 初始化标识
	AdbDate            string     `gorm:"column:adb_date;type:varchar(20)" json:"adb_date"`                  // 业务数据变动时间
	DS                 string     `gorm:"column:ds;type:varchar(20);index" json:"ds"`                        // 分区字段
	SyncTime           time.Time  `gorm:"column:sync_time;type:timestamp;index" json:"sync_time"`            // 同步时间
	UniqueKey          string     `gorm:"column:unique_key;type:varchar(200);uniqueIndex" json:"unique_key"` // 唯一键（emplid_ds）
	CreatedAt          time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt          time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (PSEmployee) TableName() string {
	return "ps_employees"
}

// PSEmployeeHonor PS员工荣誉信息表
type PSEmployeeHonor struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	EmplID         string     `gorm:"column:emplid;type:varchar(50);index" json:"emplid"`                // 工号
	BeginDt        *time.Time `gorm:"column:begin_dt;type:timestamp;index" json:"begin_dt"`              // 开始日期
	EffStatus      string     `gorm:"column:eff_status;type:varchar(20)" json:"eff_status"`             // 荣誉状态
	EndDt          *time.Time `gorm:"column:end_dt;type:timestamp" json:"end_dt"`                       // 失效日期
	Descr254       string     `gorm:"column:descr254;type:varchar(500)" json:"descr254"`                // 荣誉原因
	Descr254a      string     `gorm:"column:descr254a;type:varchar(500)" json:"descr254a"`              // 失效原因
	Comments256    string     `gorm:"column:comments_256;type:varchar(500)" json:"comments_256"`        // 备注信息
	DeptID         string     `gorm:"column:deptid;type:varchar(50)" json:"deptid"`                     // 发文部门ID
	DeptDescr      string     `gorm:"column:dept_descr;type:varchar(200)" json:"dept_descr"`            // 发文部门描述
	HrsRowAddDttm  *time.Time `gorm:"column:hrs_row_add_dttm;type:timestamp" json:"hrs_row_add_dttm"`   // 添加时间
	HrsRowAddOprid string     `gorm:"column:hrs_row_add_oprid;type:varchar(50)" json:"hrs_row_add_oprid"` // 添加用户ID
	HrsRowUpdDttm  *time.Time `gorm:"column:hrs_row_upd_dttm;type:timestamp" json:"hrs_row_upd_dttm"`   // 数据修改时间
	HrsRowUpdOprid string     `gorm:"column:hrs_row_upd_oprid;type:varchar(50)" json:"hrs_row_upd_oprid"` // 更新操作人
	AdbDate        string     `gorm:"column:adb_date;type:varchar(20)" json:"adb_date"`                 // 业务数据变动时间
	DS             string     `gorm:"column:ds;type:varchar(20);index" json:"ds"`                       // 分区字段
	SyncTime       time.Time  `gorm:"column:sync_time;type:timestamp;index" json:"sync_time"`           // 同步时间
	UniqueKey      string     `gorm:"column:unique_key;type:varchar(200);uniqueIndex" json:"unique_key"` // 唯一键（emplid_begin_dt_ds）
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (PSEmployeeHonor) TableName() string {
	return "ps_employee_honors"
}

// PSFamilyMain PS家族父表信息表
type PSFamilyMain struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	CFamilyID      int        `gorm:"column:c_family_id;type:int;index" json:"c_family_id"`          // 家族编号
	EmplID         string     `gorm:"column:emplid;type:varchar(50);index" json:"emplid"`            // 家族长工号
	CFamilyName    string     `gorm:"column:c_family_name;type:varchar(200)" json:"c_family_name"`   // 家族番号
	Date1          *time.Time `gorm:"column:date1;type:timestamp" json:"date1"`                      // 成立日期
	EffStatus      string     `gorm:"column:eff_status;type:varchar(20)" json:"eff_status"`          // 生效状态
	EffdtTo        *time.Time `gorm:"column:effdt_to;type:timestamp" json:"effdt_to"`                // 失效日期
	HrsRowAddDttm  *time.Time `gorm:"column:hrs_row_add_dttm;type:timestamp" json:"hrs_row_add_dttm"` // 添加时间
	HrsRowAddOprid string     `gorm:"column:hrs_row_add_oprid;type:varchar(50)" json:"hrs_row_add_oprid"` // 添加用户ID
	HrsRowUpdDttm  *time.Time `gorm:"column:hrs_row_upd_dttm;type:timestamp" json:"hrs_row_upd_dttm"` // 数据修改时间
	HrsRowUpdOprid string     `gorm:"column:hrs_row_upd_oprid;type:varchar(50)" json:"hrs_row_upd_oprid"` // 更新操作人
	AdbDate        string     `gorm:"column:adb_date;type:varchar(20)" json:"adb_date"`              // 业务数据变动时间
	DS             string     `gorm:"column:ds;type:varchar(20);index" json:"ds"`                    // 分区字段
	SyncTime       time.Time  `gorm:"column:sync_time;type:timestamp;index" json:"sync_time"`        // 同步时间
	UniqueKey      string     `gorm:"column:unique_key;type:varchar(200);uniqueIndex" json:"unique_key"` // 唯一键（c_family_id_ds）
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp;autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp;autoUpdateTime" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (PSFamilyMain) TableName() string {
	return "ps_family_mains"
}
