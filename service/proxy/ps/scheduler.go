package ps

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
	FamilyMemberCron   string // 家族成员调度表达式
	PositionIncCron    string // 岗位增量调度表达式（每天凌晨4:30）
	PositionAllCron    string // 岗位全量调度表达式（每周日凌晨3:30）
	OrganizationIncCron string // 组织增量调度表达式（每天凌晨4:30）
	OrganizationAllCron string // 组织全量调度表达式（每周日凌晨3:30）
	EmployeeIncCron    string // 员工增量调度表达式（每天凌晨4:30）
	EmployeeAllCron    string // 员工全量调度表达式（每周日凌晨3:30）
	EmployeeHonorCron  string // 员工荣誉调度表达式（每周日凌晨3:30）
	FamilyMainCron     string // 家族父表调度表达式（每周日凌晨3:30）
	PageSize           int    // 每页大小
	MaxPages           int    // 最大页数（防止无限循环）
}

// Scheduler PS调度器
type Scheduler struct {
	config     SchedulerConfig
	client     *PSClient
	repository *Repository
	cron       *cron.Cron
}

// NewScheduler 创建调度器
func NewScheduler(config SchedulerConfig, client *PSClient, repository *Repository) *Scheduler {
	// 设置默认值
	if config.PageSize == 0 {
		config.PageSize = 2000 // API文档最大支持2000
	}
	if config.MaxPages == 0 {
		config.MaxPages = 1000 // 默认最多查询1000页
	}

	return &Scheduler{
		config:     config,
		client:     client,
		repository: repository,
		cron:       cron.New(cron.WithSeconds()), // 支持秒级调度
	}
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	slog.Info("启动PS数据调度器")

	// 添加家族成员调度任务
	if s.config.FamilyMemberCron != "" {
		_, err := s.cron.AddFunc(s.config.FamilyMemberCron, func() {
			s.syncFamilyMembers()
		})
		if err != nil {
			slog.Error("添加家族成员调度任务失败", "error", err)
			return err
		}
		slog.Info("添加家族成员调度任务", "cron", s.config.FamilyMemberCron)
	}

	// 添加岗位增量调度任务
	if s.config.PositionIncCron != "" {
		_, err := s.cron.AddFunc(s.config.PositionIncCron, func() {
			s.syncPositionsInc()
		})
		if err != nil {
			slog.Error("添加岗位增量调度任务失败", "error", err)
			return err
		}
		slog.Info("添加岗位增量调度任务", "cron", s.config.PositionIncCron)
	}

	// 添加岗位全量调度任务
	if s.config.PositionAllCron != "" {
		_, err := s.cron.AddFunc(s.config.PositionAllCron, func() {
			s.syncPositionsAll()
		})
		if err != nil {
			slog.Error("添加岗位全量调度任务失败", "error", err)
			return err
		}
		slog.Info("添加岗位全量调度任务", "cron", s.config.PositionAllCron)
	}

	// 添加组织增量调度任务
	if s.config.OrganizationIncCron != "" {
		_, err := s.cron.AddFunc(s.config.OrganizationIncCron, func() {
			s.syncOrganizationsInc()
		})
		if err != nil {
			slog.Error("添加组织增量调度任务失败", "error", err)
			return err
		}
		slog.Info("添加组织增量调度任务", "cron", s.config.OrganizationIncCron)
	}

	// 添加组织全量调度任务
	if s.config.OrganizationAllCron != "" {
		_, err := s.cron.AddFunc(s.config.OrganizationAllCron, func() {
			s.syncOrganizationsAll()
		})
		if err != nil {
			slog.Error("添加组织全量调度任务失败", "error", err)
			return err
		}
		slog.Info("添加组织全量调度任务", "cron", s.config.OrganizationAllCron)
	}

	// 添加员工增量调度任务
	if s.config.EmployeeIncCron != "" {
		_, err := s.cron.AddFunc(s.config.EmployeeIncCron, func() {
			s.syncEmployeesInc()
		})
		if err != nil {
			slog.Error("添加员工增量调度任务失败", "error", err)
			return err
		}
		slog.Info("添加员工增量调度任务", "cron", s.config.EmployeeIncCron)
	}

	// 添加员工全量调度任务
	if s.config.EmployeeAllCron != "" {
		_, err := s.cron.AddFunc(s.config.EmployeeAllCron, func() {
			s.syncEmployeesAll()
		})
		if err != nil {
			slog.Error("添加员工全量调度任务失败", "error", err)
			return err
		}
		slog.Info("添加员工全量调度任务", "cron", s.config.EmployeeAllCron)
	}

	// 添加员工荣誉调度任务
	if s.config.EmployeeHonorCron != "" {
		_, err := s.cron.AddFunc(s.config.EmployeeHonorCron, func() {
			s.syncEmployeeHonors()
		})
		if err != nil {
			slog.Error("添加员工荣誉调度任务失败", "error", err)
			return err
		}
		slog.Info("添加员工荣誉调度任务", "cron", s.config.EmployeeHonorCron)
	}

	// 添加家族父表调度任务
	if s.config.FamilyMainCron != "" {
		_, err := s.cron.AddFunc(s.config.FamilyMainCron, func() {
			s.syncFamilyMain()
		})
		if err != nil {
			slog.Error("添加家族父表调度任务失败", "error", err)
			return err
		}
		slog.Info("添加家族父表调度任务", "cron", s.config.FamilyMainCron)
	}

	// 启动调度器
	s.cron.Start()
	slog.Info("PS数据调度器已启动")

	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	slog.Info("停止PS数据调度器")
	s.cron.Stop()
}

// syncFamilyMembers 同步家族成员数据
func (s *Scheduler) syncFamilyMembers() {
	slog.Info("开始同步PS家族成员数据")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// 计算DS分区字段（当天日期减1天，格式：YYYYMMDD）
	// 根据API文档：2025年3月21日3点后调用传 20250320
	yesterday := time.Now().AddDate(0, 0, -1)
	ds := yesterday.Format("20060102")

	slog.Info("开始同步家族成员数据", "ds", ds)

	totalRecords := 0
	pageNum := 1

	// 分页查询所有数据
	for pageNum <= s.config.MaxPages {
		slog.Info("查询家族成员数据", "ds", ds, "page", pageNum, "pageSize", s.config.PageSize)

		// 查询数据
		result, err := s.client.QueryFamilyMembers(ctx, pageNum, s.config.PageSize, ds, "")
		if err != nil {
			slog.Error("同步家族成员数据失败", "error", err, "page", pageNum)
			return
		}

		// 检查是否有数据
		if len(result) == 0 {
			slog.Info("没有更多数据，同步完成", "page", pageNum)
			break
		}

		// 保存到数据库
		if err := s.repository.SaveFamilyMembers(ctx, result, ds); err != nil {
			slog.Error("保存家族成员数据失败", "error", err, "page", pageNum)
			return
		}

		totalRecords += len(result)

		// 如果返回的记录数小于pageSize，说明已经是最后一页
		if len(result) < s.config.PageSize {
			slog.Info("已到最后一页，同步完成", "page", pageNum, "count", len(result))
			break
		}

		pageNum++

		// 添加延迟，避免API请求过快
		time.Sleep(1 * time.Second)
	}

	slog.Info("同步PS家族成员数据完成", "ds", ds, "total_records", totalRecords, "total_pages", pageNum-1)
}

// TriggerSync 手动触发同步
func (s *Scheduler) TriggerSync(dataType string) error {
	slog.Info("调度器收到同步触发请求", "data_type", dataType)
	
	switch dataType {
	case "family-members":
		slog.Info("开始异步执行家族成员数据同步")
		go s.syncFamilyMembers()
		slog.Info("手动触发家族成员数据同步成功")
	case "positions-inc":
		slog.Info("开始异步执行岗位增量数据同步")
		go s.syncPositionsInc()
		slog.Info("手动触发岗位增量数据同步成功")
	case "positions-all":
		slog.Info("开始异步执行岗位全量数据同步")
		go s.syncPositionsAll()
		slog.Info("手动触发岗位全量数据同步成功")
	case "organizations-inc":
		slog.Info("开始异步执行组织增量数据同步")
		go s.syncOrganizationsInc()
		slog.Info("手动触发组织增量数据同步成功")
	case "organizations-all":
		slog.Info("开始异步执行组织全量数据同步")
		go s.syncOrganizationsAll()
		slog.Info("手动触发组织全量数据同步成功")
	case "employees-inc":
		slog.Info("开始异步执行员工增量数据同步")
		go s.syncEmployeesInc()
		slog.Info("手动触发员工增量数据同步成功")
	case "employees-all":
		slog.Info("开始异步执行员工全量数据同步")
		go s.syncEmployeesAll()
		slog.Info("手动触发员工全量数据同步成功")
	case "employee-honors":
		slog.Info("开始异步执行员工荣誉数据同步")
		go s.syncEmployeeHonors()
		slog.Info("手动触发员工荣誉数据同步成功")
	case "family-main":
		slog.Info("开始异步执行家族父表数据同步")
		go s.syncFamilyMain()
		slog.Info("手动触发家族父表数据同步成功")
	case "all":
		slog.Info("开始异步执行所有数据同步")
		go func() {
			s.syncFamilyMembers()
			s.syncPositionsInc()
			s.syncPositionsAll()
			s.syncOrganizationsInc()
			s.syncOrganizationsAll()
			s.syncEmployeesInc()
			s.syncEmployeesAll()
			s.syncEmployeeHonors()
			s.syncFamilyMain()
		}()
		slog.Info("手动触发所有数据同步成功")
	default:
		slog.Error("不支持的数据类型", "data_type", dataType)
		return fmt.Errorf("不支持的数据类型: %s", dataType)
	}
	return nil
}

// GetNextRun 获取下次运行时间
func (s *Scheduler) GetNextRun() map[string]time.Time {
	entries := s.cron.Entries()
	result := make(map[string]time.Time)

	if len(entries) > 0 {
		result["family_members"] = entries[0].Next
	}

	return result
}

// syncPositionsInc 同步岗位增量数据
func (s *Scheduler) syncPositionsInc() {
	slog.Info("开始同步PS岗位增量数据")
	s.syncData("岗位增量", true, func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error) {
		return s.client.QueryPositions(ctx, pageNum, pageSize, ds, true)
	}, func(ctx context.Context, data []interface{}, ds string) error {
		return s.repository.SavePositions(ctx, data, ds)
	})
}

// syncPositionsAll 同步岗位全量数据
func (s *Scheduler) syncPositionsAll() {
	slog.Info("开始同步PS岗位全量数据")
	s.syncData("岗位全量", false, func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error) {
		return s.client.QueryPositions(ctx, pageNum, pageSize, ds, false)
	}, func(ctx context.Context, data []interface{}, ds string) error {
		return s.repository.SavePositions(ctx, data, ds)
	})
}

// syncOrganizationsInc 同步组织增量数据
func (s *Scheduler) syncOrganizationsInc() {
	slog.Info("开始同步PS组织增量数据")
	s.syncData("组织增量", true, func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error) {
		return s.client.QueryOrganizations(ctx, pageNum, pageSize, ds, true)
	}, func(ctx context.Context, data []interface{}, ds string) error {
		return s.repository.SaveOrganizations(ctx, data, ds)
	})
}

// syncOrganizationsAll 同步组织全量数据
func (s *Scheduler) syncOrganizationsAll() {
	slog.Info("开始同步PS组织全量数据")
	s.syncData("组织全量", false, func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error) {
		return s.client.QueryOrganizations(ctx, pageNum, pageSize, ds, false)
	}, func(ctx context.Context, data []interface{}, ds string) error {
		return s.repository.SaveOrganizations(ctx, data, ds)
	})
}

// syncEmployeesInc 同步员工增量数据
func (s *Scheduler) syncEmployeesInc() {
	slog.Info("开始同步PS员工增量数据")
	s.syncData("员工增量", true, func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error) {
		return s.client.QueryEmployees(ctx, pageNum, pageSize, ds, "", true)
	}, func(ctx context.Context, data []interface{}, ds string) error {
		return s.repository.SaveEmployees(ctx, data, ds)
	})
}

// syncEmployeesAll 同步员工全量数据
func (s *Scheduler) syncEmployeesAll() {
	slog.Info("开始同步PS员工全量数据")
	s.syncData("员工全量", false, func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error) {
		return s.client.QueryEmployees(ctx, pageNum, pageSize, ds, "", false)
	}, func(ctx context.Context, data []interface{}, ds string) error {
		return s.repository.SaveEmployees(ctx, data, ds)
	})
}

// syncEmployeeHonors 同步员工荣誉数据
func (s *Scheduler) syncEmployeeHonors() {
	slog.Info("开始同步PS员工荣誉数据")
	s.syncData("员工荣誉", false, func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error) {
		return s.client.QueryEmployeeHonors(ctx, pageNum, pageSize, ds, "")
	}, func(ctx context.Context, data []interface{}, ds string) error {
		return s.repository.SaveEmployeeHonors(ctx, data, ds)
	})
}

// syncFamilyMain 同步家族父表数据
func (s *Scheduler) syncFamilyMain() {
	slog.Info("开始同步PS家族父表数据")
	s.syncData("家族父表", false, func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error) {
		return s.client.QueryFamilyMain(ctx, pageNum, pageSize, ds, "")
	}, func(ctx context.Context, data []interface{}, ds string) error {
		return s.repository.SaveFamilyMain(ctx, data, ds)
	})
}

// syncData 通用同步数据方法
func (s *Scheduler) syncData(
	dataName string,
	isIncremental bool,
	queryFunc func(ctx context.Context, pageNum, pageSize int, ds string) ([]interface{}, error),
	saveFunc func(ctx context.Context, data []interface{}, ds string) error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// 计算DS分区字段
	// 增量接口（3-4点更新）：使用前一天日期，4:30调度
	// 全量接口（2-3点更新）：使用前一天日期，3:30调度（周日）
	yesterday := time.Now().AddDate(0, 0, -1)
	ds := yesterday.Format("20060102")

	slog.Info(fmt.Sprintf("开始同步%s数据", dataName), "ds", ds, "is_incremental", isIncremental)

	totalRecords := 0
	pageNum := 1

	// 分页查询所有数据
	for pageNum <= s.config.MaxPages {
		slog.Info(fmt.Sprintf("查询%s数据", dataName), "ds", ds, "page", pageNum, "pageSize", s.config.PageSize)

		// 查询数据
		result, err := queryFunc(ctx, pageNum, s.config.PageSize, ds)
		if err != nil {
			slog.Error(fmt.Sprintf("同步%s数据失败", dataName), "error", err, "page", pageNum)
			return
		}

		// 检查是否有数据
		if len(result) == 0 {
			slog.Info(fmt.Sprintf("没有更多%s数据，同步完成", dataName), "page", pageNum)
			break
		}

		// 保存到数据库
		if err := saveFunc(ctx, result, ds); err != nil {
			slog.Error(fmt.Sprintf("保存%s数据失败", dataName), "error", err, "page", pageNum)
			return
		}

		totalRecords += len(result)

		// 如果返回的记录数小于pageSize，说明已经是最后一页
		if len(result) < s.config.PageSize {
			slog.Info(fmt.Sprintf("已到最后一页，%s同步完成", dataName), "page", pageNum, "count", len(result))
			break
		}

		pageNum++

		// 添加延迟，避免API请求过快
		time.Sleep(1 * time.Second)
	}

	slog.Info(fmt.Sprintf("同步%s数据完成", dataName), "ds", ds, "total_records", totalRecords, "total_pages", pageNum-1)
}

