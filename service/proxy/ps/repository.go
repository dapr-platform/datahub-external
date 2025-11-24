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

// deduplicateFamilyMembers 去重家族成员数据（保留最后一条）
func deduplicateFamilyMembers(records []PSFamilyMember) []PSFamilyMember {
	seen := make(map[string]int)
	result := make([]PSFamilyMember, 0, len(records))
	
	// 第一遍：记录每个unique_key最后出现的位置
	for i, record := range records {
		seen[record.UniqueKey] = i
	}
	
	// 第二遍：只保留最后出现的记录
	for i, record := range records {
		if seen[record.UniqueKey] == i {
			result = append(result, record)
		}
	}
	
	if len(result) < len(records) {
		slog.Warn("发现重复的家族成员数据", "original", len(records), "deduplicated", len(result), "duplicates", len(records)-len(result))
	}
	
	return result
}

// deduplicatePositions 去重岗位数据（保留最后一条）
func deduplicatePositions(records []PSPosition) []PSPosition {
	seen := make(map[string]int)
	result := make([]PSPosition, 0, len(records))
	
	for i, record := range records {
		seen[record.UniqueKey] = i
	}
	
	for i, record := range records {
		if seen[record.UniqueKey] == i {
			result = append(result, record)
		}
	}
	
	if len(result) < len(records) {
		slog.Warn("发现重复的岗位数据", "original", len(records), "deduplicated", len(result), "duplicates", len(records)-len(result))
	}
	
	return result
}

// deduplicateOrganizations 去重组织数据（保留最后一条）
func deduplicateOrganizations(records []PSOrganization) []PSOrganization {
	seen := make(map[string]int)
	result := make([]PSOrganization, 0, len(records))
	
	for i, record := range records {
		seen[record.UniqueKey] = i
	}
	
	for i, record := range records {
		if seen[record.UniqueKey] == i {
			result = append(result, record)
		}
	}
	
	if len(result) < len(records) {
		slog.Warn("发现重复的组织数据", "original", len(records), "deduplicated", len(result), "duplicates", len(records)-len(result))
	}
	
	return result
}

// deduplicateEmployees 去重员工数据（保留最后一条）
func deduplicateEmployees(records []PSEmployee) []PSEmployee {
	seen := make(map[string]int)
	result := make([]PSEmployee, 0, len(records))
	
	for i, record := range records {
		seen[record.UniqueKey] = i
	}
	
	for i, record := range records {
		if seen[record.UniqueKey] == i {
			result = append(result, record)
		}
	}
	
	if len(result) < len(records) {
		slog.Warn("发现重复的员工数据", "original", len(records), "deduplicated", len(result), "duplicates", len(records)-len(result))
	}
	
	return result
}

// deduplicateEmployeeHonors 去重员工荣誉数据（保留最后一条）
func deduplicateEmployeeHonors(records []PSEmployeeHonor) []PSEmployeeHonor {
	seen := make(map[string]int)
	result := make([]PSEmployeeHonor, 0, len(records))
	
	for i, record := range records {
		seen[record.UniqueKey] = i
	}
	
	for i, record := range records {
		if seen[record.UniqueKey] == i {
			result = append(result, record)
		}
	}
	
	if len(result) < len(records) {
		slog.Warn("发现重复的员工荣誉数据", "original", len(records), "deduplicated", len(result), "duplicates", len(records)-len(result))
	}
	
	return result
}

// deduplicateFamilyMain 去重家族父表数据（保留最后一条）
func deduplicateFamilyMain(records []PSFamilyMain) []PSFamilyMain {
	seen := make(map[string]int)
	result := make([]PSFamilyMain, 0, len(records))
	
	for i, record := range records {
		seen[record.UniqueKey] = i
	}
	
	for i, record := range records {
		if seen[record.UniqueKey] == i {
			result = append(result, record)
		}
	}
	
	if len(result) < len(records) {
		slog.Warn("发现重复的家族父表数据", "original", len(records), "deduplicated", len(result), "duplicates", len(records)-len(result))
	}
	
	return result
}

// SaveFamilyMembers 保存家族成员数据
func (r *Repository) SaveFamilyMembers(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有家族成员数据需要保存")
		return nil
	}

	syncTime := time.Now()
	records := make([]PSFamilyMember, 0)

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
			UniqueKey:        fmt.Sprintf("%d_%s", resp.CFamilyID, resp.EmplID), // 不包含ds，只保留最新数据
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

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的家族成员数据需要保存")
		return nil
	}

	// 去重：同一批次中如果有相同unique_key，只保留最后一条
	records = deduplicateFamilyMembers(records)

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertFamilyMembers(ctx, records)
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

// CleanOldDSData 清理指定数据类型的旧DS数据（保留当前DS）
// 用于全量同步成功后，删除旧版本的数据
// dataName: 数据类型名称，如 "family-members", "positions-all" 等
func (r *Repository) CleanOldDSData(ctx context.Context, currentDS string, dataName string) error {
	var model interface{}
	var tableName string
	var typeName string

	// 根据dataName确定要清理的表
	switch dataName {
	case "family-members":
		model = &PSFamilyMember{}
		tableName = "ps_family_members"
		typeName = "家族成员"
	case "positions-inc", "positions-all":
		model = &PSPosition{}
		tableName = "ps_positions"
		typeName = "岗位"
	case "organizations-inc", "organizations-all":
		model = &PSOrganization{}
		tableName = "ps_organizations"
		typeName = "组织"
	case "employees-inc", "employees-all":
		model = &PSEmployee{}
		tableName = "ps_employees"
		typeName = "员工"
	case "employee-honors":
		model = &PSEmployeeHonor{}
		tableName = "ps_employee_honors"
		typeName = "员工荣誉"
	case "family-main":
		model = &PSFamilyMain{}
		tableName = "ps_family_mains"
		typeName = "家族父表"
	default:
		slog.Warn("不支持的数据类型，跳过清理旧DS数据", "data_type", dataName)
		return nil
	}

	slog.Info("开始清理旧DS数据", "data_type", dataName, "table", tableName, "current_ds", currentDS)

	// 清理指定表的旧数据
	result := r.db.WithContext(ctx).Where("ds != ?", currentDS).Delete(model)
	if result.Error != nil {
		return fmt.Errorf("清理%s旧数据失败: %v", typeName, result.Error)
	}

	if result.RowsAffected > 0 {
		slog.Info("清理旧DS数据完成", "data_type", dataName, "table", tableName, "deleted", result.RowsAffected, "current_ds", currentDS)
	} else {
		slog.Info("没有旧DS数据需要清理", "data_type", dataName, "table", tableName, "current_ds", currentDS)
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

// SavePositions 保存岗位数据
func (r *Repository) SavePositions(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有岗位数据需要保存")
		return nil
	}

	syncTime := time.Now()
	records := make([]PSPosition, 0)

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
			UniqueKey:        resp.PositionNbr, // 不包含ds，只保留最新数据
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

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的岗位数据需要保存")
		return nil
	}

	// 去重：同一批次中如果有相同unique_key，只保留最后一条
	records = deduplicatePositions(records)

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertPositions(ctx, records)
}

// SaveOrganizations 保存组织数据
func (r *Repository) SaveOrganizations(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有组织数据需要保存")
		return nil
	}

	syncTime := time.Now()
	records := make([]PSOrganization, 0)

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
			UniqueKey:      resp.DeptID, // 不包含ds，只保留最新数据
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

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的组织数据需要保存")
		return nil
	}

	// 去重：同一批次中如果有相同unique_key，只保留最后一条
	records = deduplicateOrganizations(records)

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertOrganizations(ctx, records)
}

// SaveEmployees 保存员工数据
func (r *Repository) SaveEmployees(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有员工数据需要保存")
		return nil
	}

	syncTime := time.Now()
	records := make([]PSEmployee, 0)

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
			HrsRowAddOprid:   resp.HrsRowAddOprid,
			HrsRowUpdOprid:   resp.HrsRowUpdOprid,
			EmplID:           resp.EmplID,
			Name:             resp.Name,
			Sex:              resp.Sex,
			CSexDescr:        resp.CSexDescr,
			MarStatus:        resp.MarStatus,
			CMarDescr:        resp.CMarDescr,
			Birthdate:        resp.Birthdate,
			CBirthType:       resp.CBirthType,
			CBirtypeDescr:    resp.CBirtypeDescr,
			CBirthdate:       resp.CBirthdate,
			CBirthdate1:      resp.CBirthdate1,
			CCountryDesc1:    resp.CCountryDesc1,
			CCountryDesc2:    resp.CCountryDesc2,
			CStateDesc1:      resp.CStateDesc1,
			CStateDesc2:      resp.CStateDesc2,
			CEthnicGrpDesc:   resp.CEthnicGrpDesc,
			CEthnicGrpDesc1:  resp.CEthnicGrpDesc1,
			CPersPolityDesc:  resp.CPersPolityDesc,
			Address1:         resp.Address1,
			Address3:         resp.Address3,
			CNidTypeDesc:     resp.CNidTypeDesc,
			NationalID:       resp.NationalID,
			Phone:            resp.Phone,
			EmailAddr:        resp.EmailAddr,
			EducationLvlAchv: resp.EducationLvlAchv,
			CEducationDescr:  resp.CEducationDescr,
			SchoolDescr:      resp.SchoolDescr,
			MajorDescr:       resp.MajorDescr,
			CHireDate:        resp.CHireDate,
			CLeaveDate:       resp.CLeaveDate,
			RehireDt:         resp.RehireDt,
			EmplClass:        resp.EmplClass,
			CEmplclsDescr:    resp.CEmplclsDescr,
			BusinessDescr:    resp.BusinessDescr,
			DeptID:           resp.DeptID,
			TreeNodeNum:      resp.TreeNodeNum,
			DeptDescr:        resp.DeptDescr,
			Jobcode:          resp.Jobcode,
			JobcodeDescr:     resp.JobcodeDescr,
			PositionNbr:      resp.PositionNbr,
			PosnDescr:        resp.PosnDescr,
			CSupvLvl:         resp.CSupvLvl,
			CSupvLvlDesc:     resp.CSupvLvlDesc,
			CCountryDesc:     resp.CCountryDesc,
			CStateDesc:       resp.CStateDesc,
			CCityDesc:        resp.CCityDesc,
			Address2:         resp.Address2,
			AdbDate:          resp.AdbDate,
			DS:               ds,
			SyncTime:         syncTime,
			UniqueKey:        resp.EmplID, // 不包含ds，只保留最新数据
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

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的员工数据需要保存")
		return nil
	}

	// 去重：同一批次中如果有相同unique_key，只保留最后一条
	records = deduplicateEmployees(records)

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertEmployees(ctx, records)
}

// SaveEmployeeHonors 保存员工荣誉数据
func (r *Repository) SaveEmployeeHonors(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有员工荣誉数据需要保存")
		return nil
	}

	syncTime := time.Now()
	records := make([]PSEmployeeHonor, 0)

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
				// 使用emplid和begin_dt构建唯一键（不包含ds）
				record.UniqueKey = fmt.Sprintf("%s_%s", resp.EmplID, resp.BeginDt)
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

		// 如果没有begin_dt，使用emplid作为唯一键
		if record.UniqueKey == "" {
			record.UniqueKey = resp.EmplID
		}

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的员工荣誉数据需要保存")
		return nil
	}

	// 去重：同一批次中如果有相同unique_key，只保留最后一条
	records = deduplicateEmployeeHonors(records)

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertEmployeeHonors(ctx, records)
}

// SaveFamilyMain 保存家族父表数据
func (r *Repository) SaveFamilyMain(ctx context.Context, data []interface{}, ds string) error {
	if len(data) == 0 {
		slog.Info("没有家族父表数据需要保存")
		return nil
	}

	syncTime := time.Now()
	records := make([]PSFamilyMain, 0)

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
			UniqueKey:      fmt.Sprintf("%d", resp.CFamilyID), // 不包含ds，只保留最新数据
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

		records = append(records, record)
	}

	if len(records) == 0 {
		slog.Info("没有有效的家族父表数据需要保存")
		return nil
	}

	// 去重：同一批次中如果有相同unique_key，只保留最后一条
	records = deduplicateFamilyMain(records)

	// 使用分批UPSERT，避免PostgreSQL参数限制
	return r.batchUpsertFamilyMain(ctx, records)
}
