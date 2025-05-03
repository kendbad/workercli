package worker

import (
    "workercli/internal/domain/model"
    "workercli/internal/domain/service"
    "workercli/pkg/utils"
)

// Worker đại diện cho một worker riêng lẻ trong pool
type Worker struct {
    id        int
    tasks     <-chan model.Task
    results   chan<- model.Result
    processor service.TaskProcessor
    logger    *utils.Logger
}

// NewWorker tạo một worker mới
func NewWorker(id int, tasks <-chan model.Task, results chan<- model.Result, processor service.TaskProcessor, logger *utils.Logger) *Worker {
    return &Worker{
        id:        id,
        tasks:     tasks,
        results:   results,
        processor: processor,
        logger:    logger,
    }
}

// Run bắt đầu vòng lặp xử lý task của worker
func (w *Worker) Run(stopCh <-chan struct{}) {
    w.logger.Debugf("Worker %d khởi động", w.id)
    for {
        select {
        case task, ok := <-w.tasks:
            if !ok {
                w.logger.Debugf("Worker %d dừng do channel tasks đóng", w.id)
                return
            }
            w.logger.Debugf("Worker %d nhận task %s", w.id, task.ID)
            result, err := w.processor.ProcessTask(task)
            if err != nil {
                w.logger.Errorf("Worker %d lỗi xử lý task %s: %v", w.id, task.ID, err)
                result = model.Result{
                    TaskID:  task.ID,
                    Status:  "failed",
                    Details: "Lỗi: " + err.Error(),
                }
            }
            w.results <- result
            w.logger.Debugf("Worker %d hoàn thành task %s với trạng thái %s", w.id, task.ID, result.Status)

        case <-stopCh:
            w.logger.Debugf("Worker %d dừng do tín hiệu dừng", w.id)
            return
        }
    }
}
