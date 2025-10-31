package lvyun

import (
	"context"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
	ReservationCron    string // 预订单调度表达式
	RegistrationCron   string // 登记单调度表达式
	CheckoutCron       string // 结账单调度表达式
	BusinessReportCron string // 营业报表调度表达式
	HotelCode          string // 酒店代码
	QueryDays          int    // 查询天数（默认查询最近N天的数据）
	BusinessReportDays int    // 营业报表查询天数
}

// Scheduler 调度器
type Scheduler struct {
	config     SchedulerConfig
	client     *LvyunClient
	repository *Repository
	cron       *cron.Cron
}

// NewScheduler 创建调度器
func NewScheduler(config SchedulerConfig, client *LvyunClient, repository *Repository) *Scheduler {
	// 设置默认值
	if config.QueryDays == 0 {
		config.QueryDays = 7 // 默认查询最近7天
	}
	if config.BusinessReportDays == 0 {
		config.BusinessReportDays = 1 // 默认查询最近1天
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
	slog.Info("启动绿云数据调度器")

	// 添加预订单调度任务
	if s.config.ReservationCron != "" {
		_, err := s.cron.AddFunc(s.config.ReservationCron, func() {
			s.syncReservations()
		})
		if err != nil {
			slog.Error("添加预订单调度任务失败", "error", err)
			return err
		}
		slog.Info("添加预订单调度任务", "cron", s.config.ReservationCron)
	}

	// 添加登记单调度任务
	if s.config.RegistrationCron != "" {
		_, err := s.cron.AddFunc(s.config.RegistrationCron, func() {
			s.syncRegistrations()
		})
		if err != nil {
			slog.Error("添加登记单调度任务失败", "error", err)
			return err
		}
		slog.Info("添加登记单调度任务", "cron", s.config.RegistrationCron)
	}

	// 添加结账单调度任务
	if s.config.CheckoutCron != "" {
		_, err := s.cron.AddFunc(s.config.CheckoutCron, func() {
			s.syncCheckouts()
		})
		if err != nil {
			slog.Error("添加结账单调度任务失败", "error", err)
			return err
		}
		slog.Info("添加结账单调度任务", "cron", s.config.CheckoutCron)
	}

	// 添加营业报表调度任务
	if s.config.BusinessReportCron != "" {
		_, err := s.cron.AddFunc(s.config.BusinessReportCron, func() {
			s.syncBusinessReports()
		})
		if err != nil {
			slog.Error("添加营业报表调度任务失败", "error", err)
			return err
		}
		slog.Info("添加营业报表调度任务", "cron", s.config.BusinessReportCron)
	}

	// 启动调度器
	s.cron.Start()
	slog.Info("绿云数据调度器已启动")

	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	slog.Info("停止绿云数据调度器")
	s.cron.Stop()
}

// syncReservations 同步预订单数据
func (s *Scheduler) syncReservations() {
	slog.Info("开始同步预订单数据", "hotel_code", s.config.HotelCode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 计算查询日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -s.config.QueryDays)

	params := map[string]string{
		"hotel_code": s.config.HotelCode,
		"start_date": startDate.Format("2006-01-02 00:00:00"),
		"end_date":   endDate.Format("2006-01-02 23:59:59"),
	}

	// 查询数据
	result, err := s.client.queryReservations(ctx, params)
	if err != nil {
		slog.Error("同步预订单数据失败", "error", err)
		return
	}

	// 转换为数组
	data, ok := result.([]interface{})
	if !ok {
		slog.Error("预订单数据格式错误")
		return
	}

	// 保存到数据库
	if err := s.repository.SaveReservations(ctx, data); err != nil {
		slog.Error("保存预订单数据失败", "error", err)
		return
	}

	slog.Info("同步预订单数据完成", "count", len(data))
}

// syncRegistrations 同步登记单数据
func (s *Scheduler) syncRegistrations() {
	slog.Info("开始同步登记单数据", "hotel_code", s.config.HotelCode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 计算查询日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -s.config.QueryDays)

	params := map[string]string{
		"hotel_code": s.config.HotelCode,
		"start_date": startDate.Format("2006-01-02 00:00:00"),
		"end_date":   endDate.Format("2006-01-02 23:59:59"),
	}

	// 查询数据
	result, err := s.client.queryRegistrations(ctx, params)
	if err != nil {
		slog.Error("同步登记单数据失败", "error", err)
		return
	}

	// 转换为数组
	data, ok := result.([]interface{})
	if !ok {
		slog.Error("登记单数据格式错误")
		return
	}

	// 保存到数据库
	if err := s.repository.SaveRegistrations(ctx, data); err != nil {
		slog.Error("保存登记单数据失败", "error", err)
		return
	}

	slog.Info("同步登记单数据完成", "count", len(data))
}

// syncCheckouts 同步结账单数据
func (s *Scheduler) syncCheckouts() {
	slog.Info("开始同步结账单数据", "hotel_code", s.config.HotelCode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 查询今天的数据
	bizDate := time.Now().Format("2006-01-02 00:00:00")

	params := map[string]string{
		"hotel_code": s.config.HotelCode,
		"biz_date":   bizDate,
	}

	// 查询数据
	result, err := s.client.queryCheckouts(ctx, params)
	if err != nil {
		slog.Error("同步结账单数据失败", "error", err)
		return
	}

	// 转换为数组
	data, ok := result.([]interface{})
	if !ok {
		slog.Error("结账单数据格式错误")
		return
	}

	// 保存到数据库
	if err := s.repository.SaveCheckouts(ctx, data); err != nil {
		slog.Error("保存结账单数据失败", "error", err)
		return
	}

	slog.Info("同步结账单数据完成", "count", len(data))
}

// syncBusinessReports 同步营业报表数据
func (s *Scheduler) syncBusinessReports() {
	slog.Info("开始同步营业报表数据", "hotel_code", s.config.HotelCode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 计算查询日期范围
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -s.config.BusinessReportDays)

	params := map[string]string{
		"hotel_code": s.config.HotelCode,
		"start_date": startDate.Format("2006-01-02 00:00:00"),
		"end_date":   endDate.Format("2006-01-02 23:59:59"),
	}

	// 查询数据
	result, err := s.client.queryBusinessReport(ctx, params)
	if err != nil {
		slog.Error("同步营业报表数据失败", "error", err)
		return
	}

	// 转换为数组
	data, ok := result.([]interface{})
	if !ok {
		slog.Error("营业报表数据格式错误")
		return
	}

	// 保存到数据库
	if err := s.repository.SaveBusinessReports(ctx, data); err != nil {
		slog.Error("保存营业报表数据失败", "error", err)
		return
	}

	slog.Info("同步营业报表数据完成", "count", len(data))
}

// TriggerSync 手动触发同步
func (s *Scheduler) TriggerSync(dataType string) error {
	switch dataType {
	case "reservations":
		go s.syncReservations()
	case "registrations":
		go s.syncRegistrations()
	case "checkouts":
		go s.syncCheckouts()
	case "business-reports":
		go s.syncBusinessReports()
	case "all":
		go s.syncReservations()
		go s.syncRegistrations()
		go s.syncCheckouts()
		go s.syncBusinessReports()
	default:
		return nil
	}
	return nil
}
