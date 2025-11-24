package lvyun

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

	// 设置合理的范围：最小100，最大1000
	if batchSize < 100 {
		batchSize = 100
	}
	if batchSize > 1000 {
		batchSize = 1000
	}

	return batchSize
}

// batchUpsertReservations 分批UPSERT预订单数据
func (r *Repository) batchUpsertReservations(ctx context.Context, records []LvyunReservation) error {
	if len(records) == 0 {
		return nil
	}

	// 预订单表约18个字段，批次大小约3000条
	batchSize := calculateBatchSize(18)
	totalRecords := len(records)

	slog.Info("开始分批保存预订单数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]

		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"hotel_code", "hotel_name", "lvyun_id", "guest_name", "room_type",
				"room_no", "room_num", "adult", "arr_date", "dep_date", "rsv_class",
				"rsv_no", "packages", "create_user", "create_datetime", "mobile",
				"status", "remark", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存预订单数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存预订单批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

// batchUpsertRegistrations 分批UPSERT登记单数据
func (r *Repository) batchUpsertRegistrations(ctx context.Context, records []LvyunRegistration) error {
	if len(records) == 0 {
		return nil
	}

	// 登记单表约19个字段，批次大小约3000条
	batchSize := calculateBatchSize(19)
	totalRecords := len(records)

	slog.Info("开始分批保存登记单数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]

		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"hotel_code", "hotel_name", "lvyun_id", "guest_name", "room_type",
				"room_no", "room_num", "adult", "arr_date", "dep_date", "rsv_class",
				"rsv_no", "packages", "create_user", "create_datetime", "mobile",
				"status", "master_id", "remark", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存登记单数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存登记单批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

// batchUpsertCheckouts 分批UPSERT结账单数据
func (r *Repository) batchUpsertCheckouts(ctx context.Context, records []LvyunCheckout) error {
	if len(records) == 0 {
		return nil
	}

	// 结账单表约15个字段，批次大小约4000条
	batchSize := calculateBatchSize(15)
	totalRecords := len(records)

	slog.Info("开始分批保存结账单数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]

		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"hotel_code", "hotel_name", "lvyun_id", "biz_date", "accnt",
				"arrange", "ta_code", "ta_desc", "amount", "room_no",
				"guest_name", "arr_dep", "dep_date", "mobile", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存结账单数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存结账单批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

// batchUpsertBusinessReports 分批UPSERT营业报表数据
func (r *Repository) batchUpsertBusinessReports(ctx context.Context, records []LvyunBusinessReport) error {
	if len(records) == 0 {
		return nil
	}

	// 营业报表约17个字段，批次大小约3500条
	batchSize := calculateBatchSize(17)
	totalRecords := len(records)

	slog.Info("开始分批保存营业报表数据", "total", totalRecords, "batch_size", batchSize)

	for i := 0; i < totalRecords; i += batchSize {
		end := i + batchSize
		if end > totalRecords {
			end = totalRecords
		}

		batch := records[i:end]

		result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "unique_key"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"hotel_code", "hotel_name", "biz_date", "code", "descript",
				"day", "month", "year", "day_rebate", "month_rebate",
				"year_rebate", "tax_day", "tax_month", "tax_year", "sync_time",
			}),
		}).Create(&batch)

		if result.Error != nil {
			return fmt.Errorf("保存营业报表数据失败(batch %d-%d): %v", i, end, result.Error)
		}

		slog.Info("保存营业报表批次成功", "batch", fmt.Sprintf("%d-%d/%d", i+1, end, totalRecords), "rows_affected", result.RowsAffected)
	}

	return nil
}

