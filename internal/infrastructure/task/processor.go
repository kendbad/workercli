package task

import (
	"time"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type Processor struct {
	logger *utils.Logger
}

func NewProcessor(logger *utils.Logger) *Processor {
	return &Processor{logger: logger}
}

func (p *Processor) ProcessTask(task model.Task) (model.Result, error) {
	p.logger.Debugf("Xử lý task %s với dữ liệu: %s", task.ID, task.Data)
	time.Sleep(10 * time.Millisecond) // Giả lập xử lý nhanh
	result := model.Result{
		TaskID:  task.ID,
		Status:  "success",
		Details: "Task hoàn thành: " + task.Data,
	}
	p.logger.Infof("Kết quả task %s: %s", task.ID, result.Status)
	return result, nil
}
