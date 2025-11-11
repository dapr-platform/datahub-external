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
