package ps

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Repository PS数据仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建数据仓库
func NewRepository(db interface{}) *Repository {
	if gormDB, ok := db.(*gorm.DB); ok {
		return &Repository{db: gormDB}
	}
	return nil
}

// AutoMigrate 自动迁移数据库
func (r *Repository) AutoMigrate() error {
	return r.db.AutoMigrate(
		&PSFamilyMember{},
	)
}

// SaveFamilyMembers 保存家族成员数据
func (r *Repository) SaveFamilyMembers(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有家族成员数据需要保存")
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]PSFamilyMember, 0)
	updateRecords := make([]PSFamilyMember, 0)

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化家族成员数据失败", "error", err)
			continue
		}

		var resp FamilyResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析家族成员数据失败", "error", err, "data", string(jsonData))
			continue
		}

		// 构建数据库记录
		record := PSFamilyMember{
			CFamilyID:        resp.CFamilyID,
			EmplID:           resp.EmplID,
			CFmyRelationship: resp.CFmyRelationship,
			CFamilyName:      resp.CFamilyName,
			EffStatus:        resp.EffStatus,
			HrsRowAddOprid:   resp.HrsRowAddOprid,
			HrsRowUpdOprid:   resp.HrsRowUpdOprid,
			AdbDate:          resp.AdbDate,
			DS:               ds,
			SyncTime:         syncTime,
			UniqueKey:        fmt.Sprintf("%d_%s_%s", resp.CFamilyID, resp.EmplID, ds),
		}

		// 解析时间字段
		if resp.Date1 != "" {
			if t, err := parseDateTime(resp.Date1); err == nil {
				record.Date1 = &t
			} else {
				slog.Warn("解析成立日期失败", "date1", resp.Date1, "error", err)
			}
		}

		if resp.HrsRowAddDttm != nil && *resp.HrsRowAddDttm != "" {
			if t, err := parseDateTime(*resp.HrsRowAddDttm); err == nil {
				record.HrsRowAddDttm = &t
			} else {
				slog.Warn("解析添加时间失败", "hrs_row_add_dttm", *resp.HrsRowAddDttm, "error", err)
			}
		}

		if resp.HrsRowUpdDttm != nil && *resp.HrsRowUpdDttm != "" {
			if t, err := parseDateTime(*resp.HrsRowUpdDttm); err == nil {
				record.HrsRowUpdDttm = &t
			} else {
				slog.Warn("解析更新时间失败", "hrs_row_upd_dttm", *resp.HrsRowUpdDttm, "error", err)
			}
		}

		// 查询数据库中是否已存在该记录
		var existing PSFamilyMember
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			// 新记录，直接添加
			newRecords = append(newRecords, record)
		} else if err == nil {
			// 记录已存在，比较是否有变化
			if isFamilyMemberChanged(existing, record) {
				record.ID = existing.ID
				record.CreatedAt = existing.CreatedAt
				updateRecords = append(updateRecords, record)
			}
		} else {
			slog.Error("查询家族成员记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	// 批量插入新记录
	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入家族成员数据失败: %v", err)
		}
		slog.Info("插入家族成员新记录", "count", len(newRecords))
	}

	// 批量更新已变化的记录
	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新家族成员记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新家族成员记录", "count", len(updateRecords))
	}

	skipped := len(data) - len(newRecords) - len(updateRecords)
	if skipped > 0 {
		slog.Info("跳过未变化的家族成员记录", "count", skipped)
	}

	return nil
}

// GetFamilyMembers 查询家族成员数据
func (r *Repository) GetFamilyMembers(ctx context.Context, ds string, emplID string, pageNum, pageSize int) ([]PSFamilyMember, int64, error) {
	var records []PSFamilyMember
	var total int64

	query := r.db.WithContext(ctx).Model(&PSFamilyMember{})

	// 添加分区字段过滤
	if ds != "" {
		query = query.Where("ds = ?", ds)
	}

	// 添加员工ID过滤
	if emplID != "" {
		query = query.Where("emplid = ?", emplID)
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询总数失败: %v", err)
	}

	// 分页查询
	offset := (pageNum - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("查询记录失败: %v", err)
	}

	return records, total, nil
}

// GetLatestDS 获取最新的分区字段
func (r *Repository) GetLatestDS(ctx context.Context) (string, error) {
	var result struct {
		DS string
	}

	err := r.db.WithContext(ctx).Model(&PSFamilyMember{}).
		Select("ds").
		Order("ds DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", fmt.Errorf("查询最新分区失败: %v", err)
	}

	return result.DS, nil
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
		time.RFC3339,
		time.RFC3339Nano,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析日期: %s", dateStr)
}

// isFamilyMemberChanged 比较家族成员记录是否有变化
func isFamilyMemberChanged(old, new PSFamilyMember) bool {
	// 比较时间字段（处理指针类型）
	date1Changed := false
	if (old.Date1 == nil) != (new.Date1 == nil) {
		date1Changed = true
	} else if old.Date1 != nil && new.Date1 != nil && !old.Date1.Equal(*new.Date1) {
		date1Changed = true
	}

	hrsRowAddDttmChanged := false
	if (old.HrsRowAddDttm == nil) != (new.HrsRowAddDttm == nil) {
		hrsRowAddDttmChanged = true
	} else if old.HrsRowAddDttm != nil && new.HrsRowAddDttm != nil && !old.HrsRowAddDttm.Equal(*new.HrsRowAddDttm) {
		hrsRowAddDttmChanged = true
	}

	hrsRowUpdDttmChanged := false
	if (old.HrsRowUpdDttm == nil) != (new.HrsRowUpdDttm == nil) {
		hrsRowUpdDttmChanged = true
	} else if old.HrsRowUpdDttm != nil && new.HrsRowUpdDttm != nil && !old.HrsRowUpdDttm.Equal(*new.HrsRowUpdDttm) {
		hrsRowUpdDttmChanged = true
	}

	return old.CFamilyID != new.CFamilyID ||
		old.EmplID != new.EmplID ||
		old.CFmyRelationship != new.CFmyRelationship ||
		old.CFamilyName != new.CFamilyName ||
		date1Changed ||
		old.EffStatus != new.EffStatus ||
		hrsRowAddDttmChanged ||
		old.HrsRowAddOprid != new.HrsRowAddOprid ||
		hrsRowUpdDttmChanged ||
		old.HrsRowUpdOprid != new.HrsRowUpdOprid ||
		old.AdbDate != new.AdbDate ||
		old.DS != new.DS
}

