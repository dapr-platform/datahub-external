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
	FamilyMemberCron string // 家族成员调度表达式（每天凌晨2点执行）
	PageSize         int    // 每页大小
	MaxPages         int    // 最大页数（防止无限循环）
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
	case "all":
		slog.Info("开始异步执行所有数据同步")
		go s.syncFamilyMembers()
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

