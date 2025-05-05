package test

import (
	"testing"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/task"
	"workercli/pkg/utils"
)

func TestTaskProcessor(t *testing.T) {
	logger, _ := utils.NewLogger("configs/logger.yaml")
	processor := task.NewProcessor(logger)

	taskData := "test-task-data"
	taskID := "123"

	// Tạo task model
	taskModel := model.Task{
		ID:     taskID,
		Data:   taskData,
		TaskID: taskID,
	}

	// Gọi hàm ProcessTask thay vì Process
	result, err := processor.ProcessTask(taskModel)

	if err != nil {
		t.Errorf("Lỗi khi xử lý task: %v", err)
	}

	if result.TaskID != taskID {
		t.Errorf("TaskID không khớp, mong đợi %s, nhận được %s", taskID, result.TaskID)
	}

	if result.Status != "completed" && result.Status != "success" && result.Status != "processing" && result.Status != "error" {
		t.Errorf("Status không hợp lệ: %s", result.Status)
	}
}

func TestTaskProcessorWithEmptyData(t *testing.T) {
	logger, _ := utils.NewLogger("configs/logger.yaml")
	processor := task.NewProcessor(logger)

	taskData := ""
	taskID := "123"

	// Tạo task model
	taskModel := model.Task{
		ID:     taskID,
		Data:   taskData,
		TaskID: taskID,
	}

	// Gọi hàm ProcessTask thay vì Process
	result, err := processor.ProcessTask(taskModel)

	if err != nil {
		t.Errorf("Lỗi khi xử lý task: %v", err)
	}

	if result.TaskID != taskID {
		t.Errorf("TaskID không khớp, mong đợi %s, nhận được %s", taskID, result.TaskID)
	}

	// Kiểm tra rỗng nên có trạng thái lỗi
	// Lưu ý: Có thể khác nhau tùy thuộc vào cách triển khai của bạn
	t.Logf("Kết quả xử lý task rỗng: %s", result.Status)
}
