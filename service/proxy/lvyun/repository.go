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

// SaveReservations 保存预订单数据
func (r *Repository) SaveReservations(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]LvyunReservation, 0)
	updateRecords := make([]LvyunReservation, 0)

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

		// 查询数据库中是否已存在该记录
		var existing LvyunReservation
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// 新记录，直接添加
			newRecords = append(newRecords, record)
		} else if err == nil {
			// 记录已存在，比较是否有变化
			if isReservationChanged(existing, record) {
				record.ID = existing.ID
				record.CreatedAt = existing.CreatedAt
				updateRecords = append(updateRecords, record)
			}
		} else {
			slog.Error("查询预订单记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	// 批量插入新记录
	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入预订单数据失败: %v", err)
		}
		slog.Info("插入预订单新记录", "count", len(newRecords))
	}

	// 批量更新已变化的记录
	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新预订单记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新预订单记录", "count", len(updateRecords))
	}

	skipped := len(data) - len(newRecords) - len(updateRecords)
	if skipped > 0 {
		slog.Info("跳过未变化的预订单记录", "count", skipped)
	}

	return nil
}

// SaveRegistrations 保存登记单数据
func (r *Repository) SaveRegistrations(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]LvyunRegistration, 0)
	updateRecords := make([]LvyunRegistration, 0)

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

		// 查询数据库中是否已存在该记录
		var existing LvyunRegistration
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// 新记录，直接添加
			newRecords = append(newRecords, record)
		} else if err == nil {
			// 记录已存在，比较是否有变化
			if isRegistrationChanged(existing, record) {
				record.ID = existing.ID
				record.CreatedAt = existing.CreatedAt
				updateRecords = append(updateRecords, record)
			}
		} else {
			slog.Error("查询登记单记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	// 批量插入新记录
	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入登记单数据失败: %v", err)
		}
		slog.Info("插入登记单新记录", "count", len(newRecords))
	}

	// 批量更新已变化的记录
	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新登记单记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新登记单记录", "count", len(updateRecords))
	}

	skipped := len(data) - len(newRecords) - len(updateRecords)
	if skipped > 0 {
		slog.Info("跳过未变化的登记单记录", "count", skipped)
	}

	return nil
}

// SaveCheckouts 保存结账单数据
func (r *Repository) SaveCheckouts(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]LvyunCheckout, 0)
	updateRecords := make([]LvyunCheckout, 0)

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

		// 查询数据库中是否已存在该记录
		var existing LvyunCheckout
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// 新记录，直接添加
			newRecords = append(newRecords, record)
		} else if err == nil {
			// 记录已存在，比较是否有变化
			if isCheckoutChanged(existing, record) {
				record.ID = existing.ID
				record.CreatedAt = existing.CreatedAt
				updateRecords = append(updateRecords, record)
			}
		} else {
			slog.Error("查询结账单记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	// 批量插入新记录
	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入结账单数据失败: %v", err)
		}
		slog.Info("插入结账单新记录", "count", len(newRecords))
	}

	// 批量更新已变化的记录
	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新结账单记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新结账单记录", "count", len(updateRecords))
	}

	skipped := len(data) - len(newRecords) - len(updateRecords)
	if skipped > 0 {
		slog.Info("跳过未变化的结账单记录", "count", skipped)
	}

	return nil
}

// SaveBusinessReports 保存营业报表数据
func (r *Repository) SaveBusinessReports(ctx context.Context, data []interface{}) error {
	if len(data) == 0 {
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]LvyunBusinessReport, 0)
	updateRecords := make([]LvyunBusinessReport, 0)

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

		// 查询数据库中是否已存在该记录
		var existing LvyunBusinessReport
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// 新记录，直接添加
			newRecords = append(newRecords, record)
		} else if err == nil {
			// 记录已存在，比较是否有变化
			if isBusinessReportChanged(existing, record) {
				record.ID = existing.ID
				record.CreatedAt = existing.CreatedAt
				updateRecords = append(updateRecords, record)
			}
		} else {
			slog.Error("查询营业报表记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	// 批量插入新记录
	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入营业报表数据失败: %v", err)
		}
		slog.Info("插入营业报表新记录", "count", len(newRecords))
	}

	// 批量更新已变化的记录
	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新营业报表记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新营业报表记录", "count", len(updateRecords))
	}

	skipped := len(data) - len(newRecords) - len(updateRecords)
	if skipped > 0 {
		slog.Info("跳过未变化的营业报表记录", "count", skipped)
	}

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

// isReservationChanged 比较预订单记录是否有变化
func isReservationChanged(old, new LvyunReservation) bool {
	return old.HotelName != new.HotelName ||
		old.LvyunID != new.LvyunID ||
		old.GuestName != new.GuestName ||
		old.RoomType != new.RoomType ||
		old.RoomNo != new.RoomNo ||
		old.RoomNum != new.RoomNum ||
		old.Adult != new.Adult ||
		!old.ArrDate.Equal(new.ArrDate) ||
		!old.DepDate.Equal(new.DepDate) ||
		old.RsvClass != new.RsvClass ||
		old.RsvNo != new.RsvNo ||
		old.Packages != new.Packages ||
		old.CreateUser != new.CreateUser ||
		!old.CreateDatetime.Equal(new.CreateDatetime) ||
		old.Mobile != new.Mobile ||
		old.Status != new.Status ||
		old.Remark != new.Remark
}

// isRegistrationChanged 比较登记单记录是否有变化
func isRegistrationChanged(old, new LvyunRegistration) bool {
	return old.HotelName != new.HotelName ||
		old.LvyunID != new.LvyunID ||
		old.GuestName != new.GuestName ||
		old.RoomType != new.RoomType ||
		old.RoomNo != new.RoomNo ||
		old.RoomNum != new.RoomNum ||
		old.Adult != new.Adult ||
		!old.ArrDate.Equal(new.ArrDate) ||
		!old.DepDate.Equal(new.DepDate) ||
		old.RsvClass != new.RsvClass ||
		old.RsvNo != new.RsvNo ||
		old.Packages != new.Packages ||
		old.CreateUser != new.CreateUser ||
		!old.CreateDatetime.Equal(new.CreateDatetime) ||
		old.Mobile != new.Mobile ||
		old.Status != new.Status ||
		old.MasterID != new.MasterID ||
		old.Remark != new.Remark
}

// isCheckoutChanged 比较结账单记录是否有变化
func isCheckoutChanged(old, new LvyunCheckout) bool {
	return old.HotelName != new.HotelName ||
		old.LvyunID != new.LvyunID ||
		!old.BizDate.Equal(new.BizDate) ||
		old.Accnt != new.Accnt ||
		old.Arrange != new.Arrange ||
		old.TaCode != new.TaCode ||
		old.TaDesc != new.TaDesc ||
		old.Amount != new.Amount ||
		old.RoomNo != new.RoomNo ||
		old.GuestName != new.GuestName ||
		!old.ArrDep.Equal(new.ArrDep) ||
		!old.DepDate.Equal(new.DepDate) ||
		old.Mobile != new.Mobile
}

// isBusinessReportChanged 比较营业报表记录是否有变化
func isBusinessReportChanged(old, new LvyunBusinessReport) bool {
	return old.HotelName != new.HotelName ||
		!old.BizDate.Equal(new.BizDate) ||
		old.Descript != new.Descript ||
		old.Day != new.Day ||
		old.Month != new.Month ||
		old.Year != new.Year ||
		old.DayRebate != new.DayRebate ||
		old.MonthRebate != new.MonthRebate ||
		old.YearRebate != new.YearRebate ||
		old.TaxDay != new.TaxDay ||
		old.TaxMonth != new.TaxMonth ||
		old.TaxYear != new.TaxYear
}
