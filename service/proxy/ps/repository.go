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
		&PSPosition{},
		&PSOrganization{},
		&PSEmployee{},
		&PSEmployeeHonor{},
		&PSFamilyMain{},
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

// SavePositions 保存岗位数据
func (r *Repository) SavePositions(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有岗位数据需要保存")
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]PSPosition, 0)
	updateRecords := make([]PSPosition, 0)

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化岗位数据失败", "error", err)
			continue
		}

		var resp PositionResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析岗位数据失败", "error", err, "data", string(jsonData))
			continue
		}

		record := PSPosition{
			PositionNbr:      resp.PositionNbr,
			PosnDescr:        resp.PosnDescr,
			BusinessUnit:     resp.BusinessUnit,
			BusinessDescr:    resp.BusinessDescr,
			DeptID:           resp.DeptID,
			DeptDescr:        resp.DeptDescr,
			CDeptType:        resp.CDeptType,
			Jobcode:          resp.Jobcode,
			JobcodeDescr:     resp.JobcodeDescr,
			CCountryDesc:     resp.CCountryDesc,
			CStateDesc:       resp.CStateDesc,
			CCityDesc:        resp.CCityDesc,
			Address1:         resp.Address1,
			ReportsTo:        resp.ReportsTo,
			Descr100:         resp.Descr100,
			EmplID:           resp.EmplID,
			Name:             resp.Name,
			CSupvLvl:         resp.CSupvLvl,
			CSupvLvlDesc:     resp.CSupvLvlDesc,
			Descr1:           resp.Descr1,
			Descr2:           resp.Descr2,
			HrsRowAddOprid:   resp.HrsRowAddOprid,
			HrsRowUpdOprid:   resp.HrsRowUpdOprid,
			COldPosnNbr:      resp.COldPosnNbr,
			AdbDate:          resp.AdbDate,
			CSequenceID:      resp.CSequenceID,
			CSequenceDescr:   resp.CSequenceDescr,
			CSubsequenceID:   resp.CSubsequenceID,
			CSubsequenceDesc: resp.CSubsequenceDesc,
			DS:               ds,
			SyncTime:         syncTime,
			UniqueKey:        fmt.Sprintf("%s_%s", resp.PositionNbr, ds),
		}

		// 解析时间字段
		if resp.BgnDt != "" {
			if t, err := parseDateTime(resp.BgnDt); err == nil {
				record.BgnDt = &t
			}
		}
		if resp.EndDt != nil && *resp.EndDt != "" {
			if t, err := parseDateTime(*resp.EndDt); err == nil {
				record.EndDt = &t
			}
		}
		if resp.HrsRowAddDttm != "" {
			if t, err := parseDateTime(resp.HrsRowAddDttm); err == nil {
				record.HrsRowAddDttm = &t
			}
		}
		if resp.HrsRowUpdDttm != "" {
			if t, err := parseDateTime(resp.HrsRowUpdDttm); err == nil {
				record.HrsRowUpdDttm = &t
			}
		}
		if resp.Createdttm != "" {
			if t, err := parseDateTime(resp.Createdttm); err == nil {
				record.Createdttm = &t
			}
		}
		if resp.CIntFlag != nil {
			record.CIntFlag = *resp.CIntFlag
		}
		if resp.Effdt != nil && *resp.Effdt != "" {
			if t, err := parseDateTime(*resp.Effdt); err == nil {
				record.Effdt = &t
			}
		}

		var existing PSPosition
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newRecords = append(newRecords, record)
		} else if err == nil {
			record.ID = existing.ID
			record.CreatedAt = existing.CreatedAt
			updateRecords = append(updateRecords, record)
		} else {
			slog.Error("查询岗位记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入岗位数据失败: %v", err)
		}
		slog.Info("插入岗位新记录", "count", len(newRecords))
	}

	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新岗位记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新岗位记录", "count", len(updateRecords))
	}

	return nil
}

// SaveOrganizations 保存组织数据
func (r *Repository) SaveOrganizations(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有组织数据需要保存")
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]PSOrganization, 0)
	updateRecords := make([]PSOrganization, 0)

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化组织数据失败", "error", err)
			continue
		}

		var resp OrganizationResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析组织数据失败", "error", err, "data", string(jsonData))
			continue
		}

		record := PSOrganization{
			Setid:          resp.Setid,
			DeptID:         resp.DeptID,
			TreeLevelNum:   resp.TreeLevelNum,
			TreeNodeNum:    resp.TreeNodeNum,
			DeptDescr:      resp.DeptDescr,
			CDeptType:      resp.CDeptType,
			CDeptTypeDescr: resp.CDeptTypeDescr,
			CCountryDesc:   resp.CCountryDesc,
			CStateDesc:     resp.CStateDesc,
			CCityDesc:      resp.CCityDesc,
			Address1:       resp.Address1,
			Parentname:     resp.Parentname,
			CDeptDescr:     resp.CDeptDescr,
			EmplID:         resp.EmplID,
			Name:           resp.Name,
			HrsRowAddOprid: resp.HrsRowAddOprid,
			HrsRowUpdOprid: resp.HrsRowUpdOprid,
			COldDeptID:     resp.COldDeptID,
			CCompany:       resp.CCompany,
			CCompanyDesc:   resp.CCompanyDesc,
			AdbDate:        resp.AdbDate,
			DS:             ds,
			SyncTime:       syncTime,
			UniqueKey:      fmt.Sprintf("%s_%s", resp.DeptID, ds),
		}

		// 解析时间字段
		if resp.BgnDt != "" {
			if t, err := parseDateTime(resp.BgnDt); err == nil {
				record.BgnDt = &t
			}
		}
		if resp.EndDt != nil && *resp.EndDt != "" {
			if t, err := parseDateTime(*resp.EndDt); err == nil {
				record.EndDt = &t
			}
		}
		if resp.HrsRowAddDttm != "" {
			if t, err := parseDateTime(resp.HrsRowAddDttm); err == nil {
				record.HrsRowAddDttm = &t
			}
		}
		if resp.HrsRowUpdDttm != "" {
			if t, err := parseDateTime(resp.HrsRowUpdDttm); err == nil {
				record.HrsRowUpdDttm = &t
			}
		}
		if resp.Createdttm != "" {
			if t, err := parseDateTime(resp.Createdttm); err == nil {
				record.Createdttm = &t
			}
		}
		if resp.CIntFlag != nil {
			record.CIntFlag = *resp.CIntFlag
		}
		if resp.Effdt != nil && *resp.Effdt != "" {
			if t, err := parseDateTime(*resp.Effdt); err == nil {
				record.Effdt = &t
			}
		}

		var existing PSOrganization
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newRecords = append(newRecords, record)
		} else if err == nil {
			record.ID = existing.ID
			record.CreatedAt = existing.CreatedAt
			updateRecords = append(updateRecords, record)
		} else {
			slog.Error("查询组织记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入组织数据失败: %v", err)
		}
		slog.Info("插入组织新记录", "count", len(newRecords))
	}

	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新组织记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新组织记录", "count", len(updateRecords))
	}

	return nil
}

// SaveEmployees 保存员工数据
func (r *Repository) SaveEmployees(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有员工数据需要保存")
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]PSEmployee, 0)
	updateRecords := make([]PSEmployee, 0)

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化员工数据失败", "error", err)
			continue
		}

		var resp EmployeeResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析员工数据失败", "error", err, "data", string(jsonData))
			continue
		}

		record := PSEmployee{
			HrsRowAddOprid:  resp.HrsRowAddOprid,
			HrsRowUpdOprid:  resp.HrsRowUpdOprid,
			EmplID:          resp.EmplID,
			Name:            resp.Name,
			Sex:             resp.Sex,
			CSexDescr:       resp.CSexDescr,
			MarStatus:       resp.MarStatus,
			CMarDescr:       resp.CMarDescr,
			Birthdate:       resp.Birthdate,
			CBirthType:      resp.CBirthType,
			CBirtypeDescr:   resp.CBirtypeDescr,
			CBirthdate:      resp.CBirthdate,
			CBirthdate1:     resp.CBirthdate1,
			CCountryDesc1:   resp.CCountryDesc1,
			CCountryDesc2:   resp.CCountryDesc2,
			CStateDesc1:     resp.CStateDesc1,
			CStateDesc2:     resp.CStateDesc2,
			CEthnicGrpDesc:  resp.CEthnicGrpDesc,
			CEthnicGrpDesc1: resp.CEthnicGrpDesc1,
			CPersPolityDesc: resp.CPersPolityDesc,
			Address1:        resp.Address1,
			Address3:        resp.Address3,
			CNidTypeDesc:    resp.CNidTypeDesc,
			NationalID:      resp.NationalID,
			Phone:           resp.Phone,
			EmailAddr:       resp.EmailAddr,
			EducationLvlAchv: resp.EducationLvlAchv,
			CEducationDescr: resp.CEducationDescr,
			SchoolDescr:     resp.SchoolDescr,
			MajorDescr:      resp.MajorDescr,
			CHireDate:       resp.CHireDate,
			CLeaveDate:      resp.CLeaveDate,
			RehireDt:        resp.RehireDt,
			EmplClass:       resp.EmplClass,
			CEmplclsDescr:   resp.CEmplclsDescr,
			BusinessDescr:   resp.BusinessDescr,
			DeptID:          resp.DeptID,
			TreeNodeNum:     resp.TreeNodeNum,
			DeptDescr:       resp.DeptDescr,
			Jobcode:         resp.Jobcode,
			JobcodeDescr:    resp.JobcodeDescr,
			PositionNbr:     resp.PositionNbr,
			PosnDescr:       resp.PosnDescr,
			CSupvLvl:        resp.CSupvLvl,
			CSupvLvlDesc:    resp.CSupvLvlDesc,
			CCountryDesc:    resp.CCountryDesc,
			CStateDesc:      resp.CStateDesc,
			CCityDesc:       resp.CCityDesc,
			Address2:        resp.Address2,
			AdbDate:         resp.AdbDate,
			DS:              ds,
			SyncTime:        syncTime,
			UniqueKey:       fmt.Sprintf("%s_%s", resp.EmplID, ds),
		}

		// 解析时间字段
		if resp.HrsRowAddDttm != "" {
			if t, err := parseDateTime(resp.HrsRowAddDttm); err == nil {
				record.HrsRowAddDttm = &t
			}
		}
		if resp.HrsRowUpdDttm != "" {
			if t, err := parseDateTime(resp.HrsRowUpdDttm); err == nil {
				record.HrsRowUpdDttm = &t
			}
		}
		if resp.Createdttm != "" {
			if t, err := parseDateTime(resp.Createdttm); err == nil {
				record.Createdttm = &t
			}
		}
		if resp.CIntFlag != nil {
			record.CIntFlag = *resp.CIntFlag
		}

		var existing PSEmployee
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newRecords = append(newRecords, record)
		} else if err == nil {
			record.ID = existing.ID
			record.CreatedAt = existing.CreatedAt
			updateRecords = append(updateRecords, record)
		} else {
			slog.Error("查询员工记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入员工数据失败: %v", err)
		}
		slog.Info("插入员工新记录", "count", len(newRecords))
	}

	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新员工记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新员工记录", "count", len(updateRecords))
	}

	return nil
}

// SaveEmployeeHonors 保存员工荣誉数据
func (r *Repository) SaveEmployeeHonors(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有员工荣誉数据需要保存")
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]PSEmployeeHonor, 0)
	updateRecords := make([]PSEmployeeHonor, 0)

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化员工荣誉数据失败", "error", err)
			continue
		}

		var resp EmployeeHonorResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析员工荣誉数据失败", "error", err, "data", string(jsonData))
			continue
		}

		record := PSEmployeeHonor{
			EmplID:         resp.EmplID,
			EffStatus:      resp.EffStatus,
			Descr254:       resp.Descr254,
			Descr254a:      resp.Descr254a,
			Comments256:    resp.Comments256,
			DeptID:         resp.DeptID,
			DeptDescr:      resp.DeptDescr,
			HrsRowAddOprid: resp.HrsRowAddOprid,
			HrsRowUpdOprid: resp.HrsRowUpdOprid,
			AdbDate:        resp.AdbDate,
			DS:             ds,
			SyncTime:       syncTime,
		}

		// 解析时间字段
		if resp.BeginDt != "" {
			if t, err := parseDateTime(resp.BeginDt); err == nil {
				record.BeginDt = &t
				// 使用begin_dt构建唯一键
				record.UniqueKey = fmt.Sprintf("%s_%s_%s", resp.EmplID, resp.BeginDt, ds)
			}
		}
		if resp.EndDt != nil && *resp.EndDt != "" {
			if t, err := parseDateTime(*resp.EndDt); err == nil {
				record.EndDt = &t
			}
		}
		if resp.HrsRowAddDttm != "" {
			if t, err := parseDateTime(resp.HrsRowAddDttm); err == nil {
				record.HrsRowAddDttm = &t
			}
		}
		if resp.HrsRowUpdDttm != "" {
			if t, err := parseDateTime(resp.HrsRowUpdDttm); err == nil {
				record.HrsRowUpdDttm = &t
			}
		}

		if record.UniqueKey == "" {
			record.UniqueKey = fmt.Sprintf("%s_%s", resp.EmplID, ds)
		}

		var existing PSEmployeeHonor
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newRecords = append(newRecords, record)
		} else if err == nil {
			record.ID = existing.ID
			record.CreatedAt = existing.CreatedAt
			updateRecords = append(updateRecords, record)
		} else {
			slog.Error("查询员工荣誉记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入员工荣誉数据失败: %v", err)
		}
		slog.Info("插入员工荣誉新记录", "count", len(newRecords))
	}

	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新员工荣誉记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新员工荣誉记录", "count", len(updateRecords))
	}

	return nil
}

// SaveFamilyMain 保存家族父表数据
func (r *Repository) SaveFamilyMain(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有家族父表数据需要保存")
		return nil
	}

	syncTime := time.Now()
	newRecords := make([]PSFamilyMain, 0)
	updateRecords := make([]PSFamilyMain, 0)

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			slog.Error("序列化家族父表数据失败", "error", err)
			continue
		}

		var resp FamilyMainResponse
		if err := json.Unmarshal(jsonData, &resp); err != nil {
			slog.Error("解析家族父表数据失败", "error", err, "data", string(jsonData))
			continue
		}

		record := PSFamilyMain{
			CFamilyID:      resp.CFamilyID,
			EmplID:         resp.EmplID,
			CFamilyName:    resp.CFamilyName,
			EffStatus:      resp.EffStatus,
			HrsRowAddOprid: resp.HrsRowAddOprid,
			HrsRowUpdOprid: resp.HrsRowUpdOprid,
			AdbDate:        resp.AdbDate,
			DS:             ds,
			SyncTime:       syncTime,
			UniqueKey:      fmt.Sprintf("%d_%s", resp.CFamilyID, ds),
		}

		// 解析时间字段
		if resp.Date1 != "" {
			if t, err := parseDateTime(resp.Date1); err == nil {
				record.Date1 = &t
			}
		}
		if resp.EffdtTo != nil && *resp.EffdtTo != "" {
			if t, err := parseDateTime(*resp.EffdtTo); err == nil {
				record.EffdtTo = &t
			}
		}
		if resp.HrsRowAddDttm != "" {
			if t, err := parseDateTime(resp.HrsRowAddDttm); err == nil {
				record.HrsRowAddDttm = &t
			}
		}
		if resp.HrsRowUpdDttm != "" {
			if t, err := parseDateTime(resp.HrsRowUpdDttm); err == nil {
				record.HrsRowUpdDttm = &t
			}
		}

		var existing PSFamilyMain
		err = r.db.WithContext(ctx).Where("unique_key = ?", record.UniqueKey).First(&existing).Error
		if err == gorm.ErrRecordNotFound {
			newRecords = append(newRecords, record)
		} else if err == nil {
			record.ID = existing.ID
			record.CreatedAt = existing.CreatedAt
			updateRecords = append(updateRecords, record)
		} else {
			slog.Error("查询家族父表记录失败", "error", err, "unique_key", record.UniqueKey)
		}
	}

	if len(newRecords) > 0 {
		if err := r.db.WithContext(ctx).Create(&newRecords).Error; err != nil {
			return fmt.Errorf("插入家族父表数据失败: %v", err)
		}
		slog.Info("插入家族父表新记录", "count", len(newRecords))
	}

	if len(updateRecords) > 0 {
		for _, record := range updateRecords {
			if err := r.db.WithContext(ctx).Save(&record).Error; err != nil {
				slog.Error("更新家族父表记录失败", "error", err, "unique_key", record.UniqueKey)
			}
		}
		slog.Info("更新家族父表记录", "count", len(updateRecords))
	}

	return nil
}

