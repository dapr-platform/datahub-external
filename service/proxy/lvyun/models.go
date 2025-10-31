package lvyun

// ReservationRequest 预订单请求参数
type ReservationRequest struct {
	HotelGroupCode string `json:"hotel_group_code"` // 集团代码
	HotelCode      string `json:"hotel_code"`       // 酒店代码
	StartDate      string `json:"start_date"`       // 开始日期 (格式: 2025-10-20 00:00:00)
	EndDate        string `json:"end_date"`         // 结束日期 (格式: 2025-10-21 00:00:00)
}

// ReservationResponse 预订单响应
type ReservationResponse struct {
	Hotelcode      string `json:"Hotelcode"`      // 酒店代码
	HotelName      string `json:"HotelName"`      // 酒店名称
	Id             string `json:"Id"`             // 账号
	Name           string `json:"Name"`           // 姓名
	RmType         string `json:"RmType"`         // 房型
	Rmno           string `json:"Rmno"`           // 房号
	Rmnum          int    `json:"Rmnum"`          // 房数
	Adult          int    `json:"Adult"`          // 成人
	ArrDate        string `json:"ArrDate"`        // 入住日期
	DepDate        string `json:"DepDate"`        // 离店日期
	RsvClass       string `json:"RsvClass"`       // 客户类型
	RsvNo          string `json:"RsvNo"`          // 预定单号
	Packages       string `json:"Packages"`       // 包价
	CreateUser     string `json:"CreateUser"`     // 创建用户
	CreateDatetime string `json:"CreateDatetime"` // 预定时间
	Mobile         string `json:"Mobile"`         // 手机号
	Sta            string `json:"Sta"`            // 订单状态 (R代表预定状态)
	Remark         string `json:"Remark"`         // 备注
}

// RegistrationRequest 登记单请求参数
type RegistrationRequest struct {
	HotelGroupCode string `json:"hotel_group_code"` // 集团代码
	HotelCode      string `json:"hotel_code"`       // 酒店代码
	StartDate      string `json:"start_date"`       // 开始日期 (格式: 2025-10-20 00:00:00)
	EndDate        string `json:"end_date"`         // 结束日期 (格式: 2025-10-21 00:00:00)
}

// RegistrationResponse 登记单响应
type RegistrationResponse struct {
	Hotelcode      string `json:"Hotelcode"`      // 酒店代码
	HotelName      string `json:"HotelName"`      // 酒店名称
	Id             string `json:"Id"`             // 账号
	Name           string `json:"Name"`           // 姓名
	RmType         string `json:"RmType"`         // 房型
	Rmno           string `json:"Rmno"`           // 房号
	Rmnum          int    `json:"Rmnum"`          // 房数
	Adult          int    `json:"Adult"`          // 成人
	ArrDate        string `json:"ArrDate"`        // 入住日期
	DepDate        string `json:"DepDate"`        // 离店日期
	RsvClass       string `json:"RsvClass"`       // 客户类型
	RsvNo          string `json:"RsvNo"`          // 预定单号
	Packages       string `json:"Packages"`       // 包价
	CreateUser     string `json:"CreateUser"`     // 创建用户
	CreateDatetime string `json:"CreateDatetime"` // 预定时间
	Mobile         string `json:"Mobile"`         // 手机号
	Sta            string `json:"Sta"`            // 订单状态 (R代表预定状态)
	MasterId       string `json:"MasterId"`       // 主单id
	Remark         string `json:"Remark"`         // 备注
}

// CheckoutRequest 结账单请求参数
type CheckoutRequest struct {
	HotelGroupCode string `json:"hotel_group_code"` // 集团代码
	HotelCode      string `json:"hotel_code"`       // 酒店代码
	BizDate        string `json:"biz_date"`         // 营业日期 (格式: 2025-10-20 00:00:00)
}

// CheckoutResponse 结账单响应
type CheckoutResponse struct {
	Hotelcode string  `json:"Hotelcode"` // 酒店代码
	HotelName string  `json:"HotelName"` // 酒店名称
	Id        string  `json:"Id"`        // 结账流水号
	BizDate   string  `json:"BizDate"`   // 日期
	Accnt     string  `json:"Accnt"`     // 账号
	Arrange   string  `json:"Arrange"`   // 入账类型
	TaCode    string  `json:"TaCode"`    // 入账代码
	TaDesc    string  `json:"TaDesc"`    // 入账描述
	Amount    float64 `json:"Amount"`    // 金额
	Rmno      string  `json:"Rmno"`      // 房间号
	NAME      string  `json:"NAME"`      // 姓名
	ArrDep    string  `json:"ArrDep"`    // 到达日期
	DepDate   string  `json:"DepDate"`   // 离开日期
	Mobile    string  `json:"Mobile"`    // 手机号
}

// BusinessReportRequest 营业报表请求参数
type BusinessReportRequest struct {
	HotelGroupCode string `json:"hotel_group_code"` // 集团代码
	HotelCode      string `json:"hotel_code"`       // 酒店代码
	StartDate      string `json:"start_date"`       // 开始日期 (格式: 2022-10-20 00:00:00)
	EndDate        string `json:"end_date"`         // 结束日期 (格式: 2022-10-21 00:00:00)
}

// BusinessReportResponse 营业报表响应
type BusinessReportResponse struct {
	Hotelcode   string  `json:"Hotelcode"`   // 酒店代码
	HotelName   string  `json:"HotelName"`   // 酒店名称
	BizDate     string  `json:"BizDate"`     // 日期
	Code        string  `json:"Code"`        // 代码
	Descript    string  `json:"Descript"`    // 描述
	Day         float64 `json:"Day"`         // 本日发生
	Month       float64 `json:"Month"`       // 本月发生
	Year        float64 `json:"Year"`        // 本年发生
	DayRebate   float64 `json:"DayRebate"`   // 本日发生-rebate
	MonthRebate float64 `json:"MonthRebate"` // 本月发生-rebate
	YearRebate  float64 `json:"YearRebate"`  // 本年发生-rebate
	TaxDay      float64 `json:"TaxDay"`      // 本日发生-税
	TaxMonth    float64 `json:"TaxMonth"`    // 本月发生-税
	TaxYear     float64 `json:"TaxYear"`     // 本年发生-税
}
