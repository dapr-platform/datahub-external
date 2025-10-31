package lvyun

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// SaveReservations 保存预订单数据
func (r *Repository) SaveReservations(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	records := make([]LvyunReservation, 0, len(data))
	syncTime := time.Now()

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

		record := LvyunReservation{
			HotelCode:  resp.Hotelcode,
			HotelName:  resp.HotelName,
			RecordID:   resp.Id,
			GuestName:  resp.Name,
			RoomType:   resp.RmType,
			RoomNo:     resp.Rmno,
			RsvClass:   resp.RsvClass,
			RsvNo:      resp.RsvNo,
			Packages:   resp.Packages,
			CreateUser: resp.CreateUser,
			Mobile:     resp.Mobile,
			Status:     resp.Sta,
			Remark:     resp.Remark,
			SyncTime:   syncTime,
			UniqueKey:  fmt.Sprintf("%s_%s_%s", resp.Hotelcode, resp.RsvNo, resp.Id),
		}

		// 解析数值字段
		if resp.Rmnum != "" {
			fmt.Sscanf(resp.Rmnum, "%d", &record.RoomNum)
		}
		if resp.Adult != "" {
			fmt.Sscanf(resp.Adult, "%d", &record.Adult)
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
		return nil
	}

	// 使用UPSERT操作
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "unique_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"hotel_name", "guest_name", "room_type", "room_no", "room_num", "adult", "arr_date", "dep_date", "rsv_class", "packages", "create_user", "mobile", "status", "remark", "sync_time", "updated_at"}),
	}).Create(&records)

	if result.Error != nil {
		return fmt.Errorf("保存预订单数据失败: %v", result.Error)
	}

	slog.Info("保存预订单数据成功", "count", len(records))
	return nil
}

// SaveRegistrations 保存登记单数据
func (r *Repository) SaveRegistrations(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	records := make([]LvyunRegistration, 0, len(data))
	syncTime := time.Now()

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

		record := LvyunRegistration{
			HotelCode:  resp.Hotelcode,
			HotelName:  resp.HotelName,
			RecordID:   resp.Id,
			GuestName:  resp.Name,
			RoomType:   resp.RmType,
			RoomNo:     resp.Rmno,
			RsvClass:   resp.RsvClass,
			RsvNo:      resp.RsvNo,
			Packages:   resp.Packages,
			CreateUser: resp.CreateUser,
			Mobile:     resp.Mobile,
			Status:     resp.Sta,
			MasterID:   resp.MasterId,
			Remark:     resp.Remark,
			SyncTime:   syncTime,
			UniqueKey:  fmt.Sprintf("%s_%s_%s", resp.Hotelcode, resp.RsvNo, resp.Id),
		}

		// 解析数值字段
		if resp.Rmnum != "" {
			fmt.Sscanf(resp.Rmnum, "%d", &record.RoomNum)
		}
		if resp.Adult != "" {
			fmt.Sscanf(resp.Adult, "%d", &record.Adult)
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
		return nil
	}

	// 使用UPSERT操作
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "unique_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"hotel_name", "guest_name", "room_type", "room_no", "room_num", "adult", "arr_date", "dep_date", "rsv_class", "packages", "create_user", "mobile", "status", "master_id", "remark", "sync_time", "updated_at"}),
	}).Create(&records)

	if result.Error != nil {
		return fmt.Errorf("保存登记单数据失败: %v", result.Error)
	}

	slog.Info("保存登记单数据成功", "count", len(records))
	return nil
}

// SaveCheckouts 保存结账单数据
func (r *Repository) SaveCheckouts(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	records := make([]LvyunCheckout, 0, len(data))
	syncTime := time.Now()

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

		record := LvyunCheckout{
			HotelCode: resp.Hotelcode,
			HotelName: resp.HotelName,
			RecordID:  resp.Id,
			Accnt:     resp.Accnt,
			Arrange:   resp.Arrange,
			TaCode:    resp.TaCode,
			TaDesc:    resp.TaDesc,
			Amount:    resp.Amount,
			RoomNo:    resp.Rmno,
			GuestName: resp.NAME,
			Mobile:    resp.Mobile,
			SyncTime:  syncTime,
			UniqueKey: fmt.Sprintf("%s_%s_%s", resp.Hotelcode, resp.BizDate, resp.Id),
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
		return nil
	}

	// 使用UPSERT操作
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "unique_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"hotel_name", "accnt", "arrange", "ta_code", "ta_desc", "amount", "room_no", "guest_name", "arr_dep", "dep_date", "mobile", "sync_time", "updated_at"}),
	}).Create(&records)

	if result.Error != nil {
		return fmt.Errorf("保存结账单数据失败: %v", result.Error)
	}

	slog.Info("保存结账单数据成功", "count", len(records))
	return nil
}

// SaveBusinessReports 保存营业报表数据
func (r *Repository) SaveBusinessReports(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	records := make([]LvyunBusinessReport, 0, len(data))
	syncTime := time.Now()

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

		record := LvyunBusinessReport{
			HotelCode:   resp.Hotelcode,
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
			UniqueKey:   fmt.Sprintf("%s_%s_%s", resp.Hotelcode, resp.BizDate, resp.Code),
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
		return nil
	}

	// 使用UPSERT操作
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "unique_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"hotel_name", "descript", "day", "month", "year", "day_rebate", "month_rebate", "year_rebate", "tax_day", "tax_month", "tax_year", "sync_time", "updated_at"}),
	}).Create(&records)

	if result.Error != nil {
		return fmt.Errorf("保存营业报表数据失败: %v", result.Error)
	}

	slog.Info("保存营业报表数据成功", "count", len(records))
	return nil
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
