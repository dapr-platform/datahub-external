package ps

import (
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm/clause"
)

// calculateBatchSize 根据字段数量计算批次大小
// PostgreSQL限制: 最多65535个参数
// 安全起见，使用60000作为上限
func calculateBatchSize(fieldCount int) int {
	const maxParams = 60000 // 留一些余量
	batchSize := maxParams / fieldCount
	
	// 设置合理的范围：最小50，最大500
	if batchSize < 50 {
		batchSize = 50
	}
	if batchSize > 500 {
		batchSize = 500
	}
	
	return batchSize
}

// batchUpsertFamilyMembers 分批UPSERT家族成员数据
func (r *Repository) batchUpsertFamilyMembers(ctx context.Context, records []PSFamilyMember) error {
	if len(records) == 0 {
		return nil
	}

	// 家族成员表约15个字段，批次大小约4000条
	batchSize := calculateBatchSize(15)
	totalRecords := len(records)
	
	slog.Info("开始分批保存家族成员数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]
		
		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"c_family_id", "emplid", "c_fmy_relationship", "c_family_name",
				"date1", "eff_status", "hrs_row_add_dttm", "hrs_row_add_oprid",
				"hrs_row_upd_dttm", "hrs_row_upd_oprid", "adb_date", "ds", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存家族成员数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存家族成员批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

// batchUpsertPositions 分批UPSERT岗位数据
func (r *Repository) batchUpsertPositions(ctx context.Context, records []PSPosition) error {
	if len(records) == 0 {
		return nil
	}

	// 岗位表约40个字段，批次大小约1500条
	batchSize := calculateBatchSize(40)
	totalRecords := len(records)
	
	slog.Info("开始分批保存岗位数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]
		
		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"position_nbr", "posn_descr", "business_unit", "business_descr",
				"deptid", "dept_descr", "c_dept_type", "jobcode", "jobcode_descr",
				"c_country_desc", "c_state_desc", "c_city_desc", "address1",
				"reports_to", "descr100", "emplid", "name", "c_supv_lvl", "c_supv_lvl_desc",
				"bgn_dt", "end_dt", "descr1", "descr2", "hrs_row_add_dttm", "hrs_row_add_oprid",
				"hrs_row_upd_dttm", "hrs_row_upd_oprid", "c_old_posn_nbr", "createdttm",
				"c_int_flag", "effdt", "adb_date", "c_sequence_id", "c_sequence_descr",
				"c_subsequence_id", "c_subsequence_desc", "ds", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存岗位数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存岗位批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

// batchUpsertOrganizations 分批UPSERT组织数据
func (r *Repository) batchUpsertOrganizations(ctx context.Context, records []PSOrganization) error {
	if len(records) == 0 {
		return nil
	}

	// 组织表约35个字段，批次大小约1700条
	batchSize := calculateBatchSize(35)
	totalRecords := len(records)
	
	slog.Info("开始分批保存组织数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]
		
		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"setid", "deptid", "tree_level_num", "tree_node_num", "dept_descr",
				"c_dept_type", "c_dept_type_descr", "c_country_desc", "c_state_desc",
				"c_city_desc", "address1", "parentname", "c_dept_descr", "emplid", "name",
				"bgn_dt", "end_dt", "hrs_row_add_dttm", "hrs_row_add_oprid",
				"hrs_row_upd_dttm", "hrs_row_upd_oprid", "c_old_deptid", "createdttm",
				"c_int_flag", "effdt", "c_company", "c_company_desc", "adb_date", "ds", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存组织数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存组织批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

// batchUpsertEmployees 分批UPSERT员工数据
func (r *Repository) batchUpsertEmployees(ctx context.Context, records []PSEmployee) error {
	if len(records) == 0 {
		return nil
	}

	// 员工表约70个字段，批次大小约850条
	batchSize := calculateBatchSize(70)
	totalRecords := len(records)
	
	slog.Info("开始分批保存员工数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]
		
		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"hrs_row_add_dttm", "hrs_row_add_oprid", "hrs_row_upd_dttm", "hrs_row_upd_oprid",
				"createdttm", "emplid", "name", "sex", "c_sex_descr", "mar_status", "c_mar_descr",
				"birthdate", "c_birth_type", "c_birtype_descr", "c_birthdate", "c_birthdate1",
				"c_country_desc1", "c_country_desc2", "c_state_desc1", "c_state_desc2",
				"c_ethnic_grp_desc", "c_ethnic_grp_desc1", "c_pers_polity_desc",
				"address1", "address3", "c_nid_type_desc", "national_id", "phone", "email_addr",
				"education_lvl_achv", "c_education_descr", "school_descr", "major_descr",
				"c_hire_date", "c_leave_date", "rehire_dt", "empl_class", "c_emplcls_descr",
				"business_descr", "deptid", "tree_node_num", "dept_descr", "jobcode", "jobcode_descr",
				"position_nbr", "posn_descr", "c_supv_lvl", "c_supv_lvl_desc",
				"c_country_desc", "c_state_desc", "c_city_desc", "address2", "c_int_flag",
				"adb_date", "ds", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存员工数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存员工批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

// batchUpsertEmployeeHonors 分批UPSERT员工荣誉数据
func (r *Repository) batchUpsertEmployeeHonors(ctx context.Context, records []PSEmployeeHonor) error {
	if len(records) == 0 {
		return nil
	}

	// 员工荣誉表约18个字段，批次大小约3000条
	batchSize := calculateBatchSize(18)
	totalRecords := len(records)
	
	slog.Info("开始分批保存员工荣誉数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]
		
		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"emplid", "begin_dt", "eff_status", "end_dt", "descr254", "descr254a",
				"comments_256", "deptid", "dept_descr", "hrs_row_add_dttm", "hrs_row_add_oprid",
				"hrs_row_upd_dttm", "hrs_row_upd_oprid", "adb_date", "ds", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存员工荣誉数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存员工荣誉批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

// batchUpsertFamilyMain 分批UPSERT家族父表数据
func (r *Repository) batchUpsertFamilyMain(ctx context.Context, records []PSFamilyMain) error {
	if len(records) == 0 {
		return nil
	}

	// 家族父表约14个字段，批次大小约4000条
	batchSize := calculateBatchSize(14)
	totalRecords := len(records)
	
	slog.Info("开始分批保存家族父表数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]
		
		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"c_family_id", "emplid", "c_family_name", "date1", "eff_status",
				"effdt_to", "hrs_row_add_dttm", "hrs_row_add_oprid",
				"hrs_row_upd_dttm", "hrs_row_upd_oprid", "adb_date", "ds", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存家族父表数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存家族父表批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

