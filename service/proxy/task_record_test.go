package proxy

import (
	"errors"
	"testing"
	"time"
)

func TestTaskRecordService_StartTask(t *testing.T) {
	service := NewTaskRecordService(10)

	taskID := service.StartTask("lvyun", "reservations")

	if taskID == "" {
		t.Error("任务ID不应为空")
	}

	records := service.QueryRecords("", "", "", 100)
	if len(records) != 1 {
		t.Errorf("期望1条记录，实际%d条", len(records))
	}

	record := records[0]
	if record.DataSource != "lvyun" {
		t.Errorf("数据源期望lvyun，实际%s", record.DataSource)
	}
	if record.TaskType != "reservations" {
		t.Errorf("任务类型期望reservations，实际%s", record.TaskType)
	}
	if record.Status != TaskStatusRunning {
		t.Errorf("状态期望running，实际%s", record.Status)
	}
}

func TestTaskRecordService_FinishTask_Success(t *testing.T) {
	service := NewTaskRecordService(10)

	taskID := service.StartTask("ps", "positions-inc")
	time.Sleep(100 * time.Millisecond) // 模拟任务执行时间

	service.FinishTask(taskID, 100, nil)

	records := service.QueryRecords("", "", "", 100)
	if len(records) != 1 {
		t.Fatalf("期望1条记录，实际%d条", len(records))
	}

	record := records[0]
	if record.Status != TaskStatusSuccess {
		t.Errorf("状态期望success，实际%s", record.Status)
	}
	if record.RecordCount != 100 {
		t.Errorf("记录数期望100，实际%d", record.RecordCount)
	}
	if record.EndTime == nil {
		t.Error("结束时间不应为空")
	}
	if record.Duration <= 0 {
		t.Errorf("执行时长应大于0，实际%d", record.Duration)
	}
	if record.ErrorMsg != "" {
		t.Errorf("错误信息应为空，实际%s", record.ErrorMsg)
	}
}

func TestTaskRecordService_FinishTask_Failed(t *testing.T) {
	service := NewTaskRecordService(10)

	taskID := service.StartTask("lvyun", "checkouts")
	time.Sleep(50 * time.Millisecond)

	testError := errors.New("数据库连接失败")
	service.FinishTask(taskID, 0, testError)

	records := service.QueryRecords("", "", "", 100)
	if len(records) != 1 {
		t.Fatalf("期望1条记录，实际%d条", len(records))
	}

	record := records[0]
	if record.Status != TaskStatusFailed {
		t.Errorf("状态期望failed，实际%s", record.Status)
	}
	if record.ErrorMsg != testError.Error() {
		t.Errorf("错误信息期望%s，实际%s", testError.Error(), record.ErrorMsg)
	}
}

func TestTaskRecordService_QueryRecords_Filter(t *testing.T) {
	service := NewTaskRecordService(100)

	// 创建多条不同的记录
	service.StartTask("lvyun", "reservations")
	service.StartTask("lvyun", "registrations")
	service.StartTask("ps", "positions-inc")
	service.StartTask("ps", "employees-all")

	// 测试按数据源过滤
	records := service.QueryRecords("lvyun", "", "", 100)
	if len(records) != 2 {
		t.Errorf("lvyun数据源期望2条记录，实际%d条", len(records))
	}

	records = service.QueryRecords("ps", "", "", 100)
	if len(records) != 2 {
		t.Errorf("ps数据源期望2条记录，实际%d条", len(records))
	}

	// 测试按任务类型过滤
	records = service.QueryRecords("", "reservations", "", 100)
	if len(records) != 1 {
		t.Errorf("reservations任务类型期望1条记录，实际%d条", len(records))
	}

	// 测试按状态过滤
	records = service.QueryRecords("", "", "running", 100)
	if len(records) != 4 {
		t.Errorf("running状态期望4条记录，实际%d条", len(records))
	}
}

func TestTaskRecordService_QueryRecords_Limit(t *testing.T) {
	service := NewTaskRecordService(100)

	// 创建5条记录
	for i := 0; i < 5; i++ {
		service.StartTask("lvyun", "reservations")
		time.Sleep(10 * time.Millisecond) // 确保时间戳不同
	}

	// 限制返回3条
	records := service.QueryRecords("", "", "", 3)
	if len(records) != 3 {
		t.Errorf("期望3条记录，实际%d条", len(records))
	}

	// 验证返回的是最新的记录（倒序）
	for i := 1; i < len(records); i++ {
		if records[i].StartTime.After(records[i-1].StartTime) {
			t.Error("记录应按时间倒序排列（最新的在前）")
		}
	}
}

func TestTaskRecordService_GetStatistics(t *testing.T) {
	service := NewTaskRecordService(100)

	// 创建不同状态的记录
	task1 := service.StartTask("lvyun", "reservations")
	time.Sleep(10 * time.Millisecond) // 确保ID不同
	task2 := service.StartTask("lvyun", "registrations")
	time.Sleep(10 * time.Millisecond)
	task3 := service.StartTask("ps", "positions-inc")
	time.Sleep(10 * time.Millisecond)
	_ = service.StartTask("ps", "employees-all") // task4 保持 running

	// 完成部分任务
	service.FinishTask(task1, 100, nil)              // success
	service.FinishTask(task2, 0, errors.New("测试错误")) // failed
	service.FinishTask(task3, 50, nil)               // success

	// 获取所有数据源的统计
	stats := service.GetStatistics("")
	if stats["total"] != 4 {
		t.Errorf("总数期望4，实际%d", stats["total"])
	}
	if stats["success"] != 2 {
		t.Errorf("成功数期望2，实际%d", stats["success"])
	}
	if stats["failed"] != 1 {
		t.Errorf("失败数期望1，实际%d", stats["failed"])
	}
	if stats["running"] != 1 {
		t.Errorf("运行中期望1，实际%d", stats["running"])
	}

	// 获取lvyun数据源的统计
	lvyunStats := service.GetStatistics("lvyun")
	if lvyunStats["total"] != 2 {
		t.Errorf("lvyun总数期望2，实际%d", lvyunStats["total"])
	}
	if lvyunStats["success"] != 1 {
		t.Errorf("lvyun成功数期望1，实际%d", lvyunStats["success"])
	}
	if lvyunStats["failed"] != 1 {
		t.Errorf("lvyun失败数期望1，实际%d", lvyunStats["failed"])
	}
}

func TestTaskRecordService_MaxSize(t *testing.T) {
	service := NewTaskRecordService(5) // 最多保存5条记录

	// 创建10条记录
	for i := 0; i < 10; i++ {
		service.StartTask("test", "task")
		time.Sleep(10 * time.Millisecond)
	}

	// 验证只保留最新的5条
	records := service.QueryRecords("", "", "", 100)
	if len(records) != 5 {
		t.Errorf("期望保留5条记录，实际%d条", len(records))
	}
}

func TestNewTaskRecordService_DefaultMaxSize(t *testing.T) {
	service := NewTaskRecordService(0)  // 传入0
	if service.maxSize != 1000 {
		t.Errorf("默认maxSize期望1000，实际%d", service.maxSize)
	}

	service2 := NewTaskRecordService(-1) // 传入负数
	if service2.maxSize != 1000 {
		t.Errorf("默认maxSize期望1000，实际%d", service2.maxSize)
	}
}

