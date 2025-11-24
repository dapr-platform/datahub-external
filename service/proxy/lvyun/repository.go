package lvyun

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Repository 数据仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建数据仓库
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// AutoMigrate 自动迁移数据库
func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(
		&LvyunReservation{},
		&LvyunRegistration{},
		&LvyunCheckout{},
		&LvyunBusinessReport{},
	)
}

// SaveReservations 保存预订单数据（使用UPSERT）
func (r *Repository) SaveReservations(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	syncTime := time.Now()
	records := make([]LvyunReservation, 0, len(data))

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化预订单数据失败", "error", err)
			continue
		}

		var resp ReservationResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析预订单数据失败", "error", err, "data", string(jsonData))
			continue
		}

		// 过滤频率限制错误响应
		if isRateLimitError(resp.Remark) {
			slog.Warn("检测到频率限制错误响应，跳过保存", "remark", resp.Remark)
			continue
		}

		record := LvyunReservation{
			HotelCode:  resp.HotelCode,
			HotelName:  resp.HotelName,
			LvyunID:    resp.Id,
			GuestName:  resp.Name,
			RoomType:   resp.RmType,
			RoomNo:     resp.Rmno,
			RoomNum:    resp.Rmnum,
			Adult:      resp.Adult,
			RsvClass:   resp.RsvClass,
			RsvNo:      resp.RsvNo,
			Packages:   resp.Packages,
			CreateUser: resp.CreateUser,
			Mobile:     resp.Mobile,
			Status:     resp.Sta,
			Remark:     resp.Remark,
			SyncTime:   syncTime,
			UniqueKey:  fmt.Sprintf("%s_%d", resp.HotelCode, resp.Id),
		}

		// 解析时间字段
		if resp.ArrDate != "" {
			if t, err := parseDateTime(resp.ArrDate); err == nil {
				record.ArrDate = t
			}
		}
		if resp.DepDate != "" {
			if t, err := parseDateTime(resp.DepDate); err == nil {
				record.DepDate = t
			}
		}
		if resp.CreateDatetime != "" {
			if t, err := parseDateTime(resp.CreateDatetime); err == nil {
				record.CreateDatetime = t
			}
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的预订单数据需要保存")
		return nil
	}

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertReservations(ctx, records)
}

// SaveRegistrations 保存登记单数据（使用UPSERT）
func (r *Repository) SaveRegistrations(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	syncTime := time.Now()
	records := make([]LvyunRegistration, 0, len(data))

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化登记单数据失败", "error", err)
			continue
		}

		var resp RegistrationResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析登记单数据失败", "error", err, "data", string(jsonData))
			continue
		}

		// 过滤频率限制错误响应
		if isRateLimitError(resp.Remark) {
			slog.Warn("检测到频率限制错误响应，跳过保存", "remark", resp.Remark)
			continue
		}

		record := LvyunRegistration{
			HotelCode:  resp.HotelCode,
			HotelName:  resp.HotelName,
			LvyunID:    resp.Id,
			GuestName:  resp.Name,
			RoomType:   resp.RmType,
			RoomNo:     resp.Rmno,
			RoomNum:    resp.Rmnum,
			Adult:      resp.Adult,
			RsvClass:   resp.RsvClass,
			RsvNo:      resp.RsvNo,
			Packages:   resp.Packages,
			CreateUser: resp.CreateUser,
			Mobile:     resp.Mobile,
			Status:     resp.Sta,
			MasterID:   resp.MasterId,
			Remark:     resp.Remark,
			SyncTime:   syncTime,
			UniqueKey:  fmt.Sprintf("%s_%d", resp.HotelCode, resp.Id),
		}

		// 解析时间字段
		if resp.ArrDate != "" {
			if t, err := parseDateTime(resp.ArrDate); err == nil {
				record.ArrDate = t
			}
		}
		if resp.DepDate != "" {
			if t, err := parseDateTime(resp.DepDate); err == nil {
				record.DepDate = t
			}
		}
		if resp.CreateDatetime != "" {
			if t, err := parseDateTime(resp.CreateDatetime); err == nil {
				record.CreateDatetime = t
			}
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的登记单数据需要保存")
		return nil
	}

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertRegistrations(ctx, records)
}

// SaveCheckouts 保存结账单数据（使用UPSERT）
func (r *Repository) SaveCheckouts(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	syncTime := time.Now()
	records := make([]LvyunCheckout, 0, len(data))

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化结账单数据失败", "error", err)
			continue
		}

		var resp CheckoutResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析结账单数据失败", "error", err, "data", string(jsonData))
			continue
		}

		// 过滤频率限制错误响应
		if isRateLimitError(resp.Remark) {
			slog.Warn("检测到频率限制错误响应，跳过保存", "remark", resp.Remark)
			continue
		}

		record := LvyunCheckout{
			HotelCode: resp.HotelCode,
			HotelName: resp.HotelName,
			LvyunID:   resp.Id,
			Accnt:     resp.Accnt,
			Arrange:   resp.Arrange,
			TaCode:    resp.TaCode,
			TaDesc:    resp.TaDesc,
			Amount:    resp.Amount,
			RoomNo:    resp.Rmno,
			GuestName: resp.NAME,
			Mobile:    resp.Mobile,
			SyncTime:  syncTime,
			UniqueKey: fmt.Sprintf("%s_%s_%d", resp.HotelCode, resp.BizDate, resp.Id),
		}

		// 解析时间字段
		if resp.BizDate != "" {
			if t, err := parseDateTime(resp.BizDate); err == nil {
				record.BizDate = t
			}
		}
		if resp.ArrDep != "" {
			if t, err := parseDateTime(resp.ArrDep); err == nil {
				record.ArrDep = t
			}
		}
		if resp.DepDate != "" {
			if t, err := parseDateTime(resp.DepDate); err == nil {
				record.DepDate = t
			}
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的结账单数据需要保存")
		return nil
	}

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertCheckouts(ctx, records)
}

// SaveBusinessReports 保存营业报表数据（使用UPSERT）
func (r *Repository) SaveBusinessReports(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	syncTime := time.Now()
	records := make([]LvyunBusinessReport, 0, len(data))

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化营业报表数据失败", "error", err)
			continue
		}

		var resp BusinessReportResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析营业报表数据失败", "error", err, "data", string(jsonData))
			continue
		}

		// 过滤频率限制错误响应
		if isRateLimitError(resp.Remark) {
			slog.Warn("检测到频率限制错误响应，跳过保存", "remark", resp.Remark)
			continue
		}

		record := LvyunBusinessReport{
			HotelCode:   resp.HotelCode,
			HotelName:   resp.HotelName,
			Code:        resp.Code,
			Descript:    resp.Descript,
			Day:         resp.Day,
			Month:       resp.Month,
			Year:        resp.Year,
			DayRebate:   resp.DayRebate,
			MonthRebate: resp.MonthRebate,
			YearRebate:  resp.YearRebate,
			TaxDay:      resp.TaxDay,
			TaxMonth:    resp.TaxMonth,
			TaxYear:     resp.TaxYear,
			SyncTime:    syncTime,
			UniqueKey:   fmt.Sprintf("%s_%s_%s", resp.HotelCode, resp.BizDate, resp.Code),
		}

		// 解析时间字段
		if resp.BizDate != "" {
			if t, err := parseDateTime(resp.BizDate); err == nil {
				record.BizDate = t
			}
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的营业报表数据需要保存")
		return nil
	}

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertBusinessReports(ctx, records)
}

// parseDateTime 解析日期时间字符串，支持多种格式
func parseDateTime(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("空日期字符串")
	}

	// 常见的日期格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析日期: %s", dateStr)
}

// isRateLimitError 检查是否为频率限制错误
func isRateLimitError(remark string) bool {
	if remark == "" {
		return false
	}
	// 检查是否包含频率限制的关键字
	return strings.Contains(remark, "请不要频繁调用接口") ||
		strings.Contains(remark, "请10分钟后再查询") ||
		(strings.Contains(remark, "请在") && strings.Contains(remark, "之后重试"))
}

