package lvyun

import (
	"time"

	"gorm.io/gorm"
)

// LvyunReservation 预订单数据模型
type LvyunReservation struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `gorm:"type:timestamp(0)" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:timestamp(0)" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	HotelCode      string         `gorm:"size:50;index" json:"hotel_code"`          // 酒店代码
	HotelName      string         `gorm:"size:200" json:"hotel_name"`               // 酒店名称
	LvyunID        int            `gorm:"index" json:"lvyun_id"`                    // 绿云系统ID(业务主键)
	GuestName      string         `gorm:"size:200" json:"guest_name"`               // 姓名
	RoomType       string         `gorm:"size:50" json:"room_type"`                 // 房型
	RoomNo         string         `gorm:"size:50" json:"room_no"`                   // 房号
	RoomNum        int            `json:"room_num"`                                 // 房数
	Adult          int            `json:"adult"`                                    // 成人
	ArrDate        time.Time      `gorm:"type:timestamp(0);index" json:"arr_date"`  // 入住日期
	DepDate        time.Time      `gorm:"type:timestamp(0);index" json:"dep_date"`  // 离店日期
	RsvClass       string         `gorm:"size:50" json:"rsv_class"`                 // 客户类型
	RsvNo          string         `gorm:"size:50;index" json:"rsv_no"`              // 预定单号
	Packages       string         `gorm:"size:100" json:"packages"`                 // 包价
	CreateUser     string         `gorm:"size:100" json:"create_user"`              // 创建用户
	CreateDatetime time.Time      `gorm:"type:timestamp(0)" json:"create_datetime"` // 预定时间
	Mobile         string         `gorm:"size:50" json:"mobile"`                    // 手机号
	Status         string         `gorm:"size:10;index" json:"status"`              // 订单状态
	Remark         string         `gorm:"type:text" json:"remark,omitempty"`        // 备注
	SyncTime       time.Time      `gorm:"type:timestamp(0);index" json:"sync_time"` // 同步时间
	UniqueKey      string         `gorm:"size:200;uniqueIndex" json:"unique_key"`   // 唯一键
}

func (LvyunReservation) TableName() string {
	return "lvyun_reservations"
}

// LvyunRegistration 登记单数据模型
type LvyunRegistration struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `gorm:"type:timestamp(0)" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:timestamp(0)" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	HotelCode      string         `gorm:"size:50;index" json:"hotel_code"`          // 酒店代码
	HotelName      string         `gorm:"size:200" json:"hotel_name"`               // 酒店名称
	LvyunID        int            `gorm:"index" json:"lvyun_id"`                    // 绿云系统ID(业务主键)
	GuestName      string         `gorm:"size:200" json:"guest_name"`               // 姓名
	RoomType       string         `gorm:"size:50" json:"room_type"`                 // 房型
	RoomNo         string         `gorm:"size:50" json:"room_no"`                   // 房号
	RoomNum        int            `json:"room_num"`                                 // 房数
	Adult          int            `json:"adult"`                                    // 成人
	ArrDate        time.Time      `gorm:"type:timestamp(0);index" json:"arr_date"`  // 入住日期
	DepDate        time.Time      `gorm:"type:timestamp(0);index" json:"dep_date"`  // 离店日期
	RsvClass       string         `gorm:"size:50" json:"rsv_class"`                 // 客户类型
	RsvNo          string         `gorm:"size:50;index" json:"rsv_no"`              // 预定单号
	Packages       string         `gorm:"size:100" json:"packages"`                 // 包价
	CreateUser     string         `gorm:"size:100" json:"create_user"`              // 创建用户
	CreateDatetime time.Time      `gorm:"type:timestamp(0)" json:"create_datetime"` // 预定时间
	Mobile         string         `gorm:"size:50" json:"mobile"`                    // 手机号
	Status         string         `gorm:"size:10;index" json:"status"`              // 订单状态
	MasterID       string         `gorm:"size:50;index" json:"master_id"`           // 主单id
	Remark         string         `gorm:"type:text" json:"remark,omitempty"`        // 备注
	SyncTime       time.Time      `gorm:"type:timestamp(0);index" json:"sync_time"` // 同步时间
	UniqueKey      string         `gorm:"size:200;uniqueIndex" json:"unique_key"`   // 唯一键
}

func (LvyunRegistration) TableName() string {
	return "lvyun_registrations"
}

// LvyunCheckout 结账单数据模型
type LvyunCheckout struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"type:timestamp(0)" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp(0)" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	HotelCode string         `gorm:"size:50;index" json:"hotel_code"`          // 酒店代码
	HotelName string         `gorm:"size:200" json:"hotel_name"`               // 酒店名称
	LvyunID   int            `gorm:"index" json:"lvyun_id"`                    // 绿云系统ID(业务主键)
	BizDate   time.Time      `gorm:"type:timestamp(0);index" json:"biz_date"`  // 日期
	Accnt     string         `gorm:"size:50" json:"accnt"`                     // 账号
	Arrange   string         `gorm:"size:50" json:"arrange"`                   // 入账类型
	TaCode    string         `gorm:"size:50" json:"ta_code"`                   // 入账代码
	TaDesc    string         `gorm:"size:200" json:"ta_desc"`                  // 入账描述
	Amount    float64        `json:"amount"`                                   // 金额
	RoomNo    string         `gorm:"size:50" json:"room_no"`                   // 房间号
	GuestName string         `gorm:"size:200" json:"guest_name"`               // 姓名
	ArrDep    time.Time      `gorm:"type:timestamp(0)" json:"arr_dep"`         // 到达日期
	DepDate   time.Time      `gorm:"type:timestamp(0)" json:"dep_date"`        // 离开日期
	Mobile    string         `gorm:"size:50" json:"mobile"`                    // 手机号
	SyncTime  time.Time      `gorm:"type:timestamp(0);index" json:"sync_time"` // 同步时间
	UniqueKey string         `gorm:"size:200;uniqueIndex" json:"unique_key"`   // 唯一键
}

func (LvyunCheckout) TableName() string {
	return "lvyun_checkouts"
}

// LvyunBusinessReport 营业报表数据模型
type LvyunBusinessReport struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `gorm:"type:timestamp(0)" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:timestamp(0)" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	HotelCode   string         `gorm:"size:50;index" json:"hotel_code"`          // 酒店代码
	HotelName   string         `gorm:"size:200" json:"hotel_name"`               // 酒店名称
	BizDate     time.Time      `gorm:"type:timestamp(0);index" json:"biz_date"`  // 日期
	Code        string         `gorm:"size:50;index" json:"code"`                // 代码
	Descript    string         `gorm:"size:200" json:"descript"`                 // 描述
	Day         float64        `json:"day"`                                      // 本日发生
	Month       float64        `json:"month"`                                    // 本月发生
	Year        float64        `json:"year"`                                     // 本年发生
	DayRebate   float64        `json:"day_rebate"`                               // 本日发生-rebate
	MonthRebate float64        `json:"month_rebate"`                             // 本月发生-rebate
	YearRebate  float64        `json:"year_rebate"`                              // 本年发生-rebate
	TaxDay      float64        `json:"tax_day"`                                  // 本日发生-税
	TaxMonth    float64        `json:"tax_month"`                                // 本月发生-税
	TaxYear     float64        `json:"tax_year"`                                 // 本年发生-税
	SyncTime    time.Time      `gorm:"type:timestamp(0);index" json:"sync_time"` // 同步时间
	UniqueKey   string         `gorm:"size:200;uniqueIndex" json:"unique_key"`   // 唯一键
}

func (LvyunBusinessReport) TableName() string {
	return "lvyun_business_reports"
}
