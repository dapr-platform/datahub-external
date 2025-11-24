package proxy

import (
	"sync"
	"time"
)

// TaskStatus 任务执行状态
type TaskStatus string

const (
	TaskStatusRunning TaskStatus = "running" // 执行中
	TaskStatusSuccess TaskStatus = "success" // 成功
	TaskStatusFailed  TaskStatus = "failed"  // 失败
)

// TaskRecord 任务执行记录
type TaskRecord struct {
	ID          string     `json:"id"`           // 记录ID（使用时间戳生成）
	DataSource  string     `json:"data_source"`  // 数据源名称（lvyun/ps）
	TaskType    string     `json:"task_type"`    // 任务类型（reservations/positions-inc等）
	Status      TaskStatus `json:"status"`       // 执行状态
	StartTime   time.Time  `json:"start_time"`   // 开始时间
	EndTime     *time.Time `json:"end_time"`     // 结束时间
	Duration    int64      `json:"duration"`     // 执行时长（毫秒）
	RecordCount int        `json:"record_count"` // 处理记录数
	ErrorMsg    string     `json:"error_msg"`    // 错误信息
}

// TaskRecordService 任务记录服务
type TaskRecordService struct {
	records []TaskRecord
	mu      sync.RWMutex
	maxSize int // 最大保存记录数
}

// NewTaskRecordService 创建任务记录服务
func NewTaskRecordService(maxSize int) *TaskRecordService {
	if maxSize <= 0 {
		maxSize = 1000 // 默认保存最近1000条记录
	}
	return &TaskRecordService{
		records: make([]TaskRecord, 0),
		maxSize: maxSize,
	}
}

// StartTask 开始任务
func (s *TaskRecordService) StartTask(dataSource, taskType string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := time.Now().Format("20060102150405.000")
	record := TaskRecord{
		ID:         id,
		DataSource: dataSource,
		TaskType:   taskType,
		Status:     TaskStatusRunning,
		StartTime:  time.Now(),
	}

	s.records = append(s.records, record)

	// 保持记录数不超过maxSize
	if len(s.records) > s.maxSize {
		s.records = s.records[len(s.records)-s.maxSize:]
	}

	return id
}

// FinishTask 完成任务
func (s *TaskRecordService) FinishTask(id string, recordCount int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 从后往前查找（最近的记录在后面）
	for i := len(s.records) - 1; i >= 0; i-- {
		if s.records[i].ID == id {
			endTime := time.Now()
			s.records[i].EndTime = &endTime
			s.records[i].Duration = endTime.Sub(s.records[i].StartTime).Milliseconds()
			s.records[i].RecordCount = recordCount

			if err != nil {
				s.records[i].Status = TaskStatusFailed
				s.records[i].ErrorMsg = err.Error()
			} else {
				s.records[i].Status = TaskStatusSuccess
			}
			break
		}
	}
}

// QueryRecords 查询任务记录
func (s *TaskRecordService) QueryRecords(dataSource, taskType, status string, limit int) []TaskRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit <= 0 {
		limit = 100 // 默认返回最近100条
	}

	result := make([]TaskRecord, 0)

	// 从后往前遍历（最新的记录在后面）
	for i := len(s.records) - 1; i >= 0 && len(result) < limit; i-- {
		record := s.records[i]

		// 过滤条件
		if dataSource != "" && record.DataSource != dataSource {
			continue
		}
		if taskType != "" && record.TaskType != taskType {
			continue
		}
		if status != "" && string(record.Status) != status {
			continue
		}

		result = append(result, record)
	}

	return result
}

// GetStatistics 获取统计信息
func (s *TaskRecordService) GetStatistics(dataSource string) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := map[string]interface{}{
		"total":   0,
		"success": 0,
		"failed":  0,
		"running": 0,
	}

	for _, record := range s.records {
		if dataSource != "" && record.DataSource != dataSource {
			continue
		}

		stats["total"] = stats["total"].(int) + 1
		switch record.Status {
		case TaskStatusSuccess:
			stats["success"] = stats["success"].(int) + 1
		case TaskStatusFailed:
			stats["failed"] = stats["failed"].(int) + 1
		case TaskStatusRunning:
			stats["running"] = stats["running"].(int) + 1
		}
	}

	return stats
}

// 全局任务记录服务
var globalTaskRecordService = NewTaskRecordService(1000)

// GetGlobalTaskRecordService 获取全局任务记录服务
func GetGlobalTaskRecordService() *TaskRecordService {
	return globalTaskRecordService
}

