package usecase

import (
	"fmt"
	"workercli/internal/adapter/input"
	"workercli/internal/adapter/worker"
	"workercli/internal/domain/model"
	"workercli/internal/domain/service"
	"workercli/pkg/utils"

	"github.com/sirupsen/logrus"
)

type BatchTaskUseCase struct {
	inputReader input.Reader
	processor   service.TaskProcessor
	workerPool  *worker.Pool
	logger      *utils.Logger
}

func NewBatchTaskUseCase(reader input.Reader, processor service.TaskProcessor, workers int, logger *utils.Logger) *BatchTaskUseCase {
	pool := worker.NewPool(workers, processor, logger)
	return &BatchTaskUseCase{
		inputReader: reader,
		processor:   processor,
		workerPool:  pool,
		logger:      logger,
	}
}

func (uc *BatchTaskUseCase) Execute(inputFile string) ([]model.Result, error) {
	uc.logger.Info(fmt.Sprintf("Bắt đầu xử lý file đầu vào: %s", inputFile))

	tasks, err := uc.inputReader.ReadTasks(inputFile)
	if err != nil {
		uc.logger.Errorf("Lỗi đọc file đầu vào: %v", err)
		return nil, err
	}

	uc.workerPool.Start()
	for _, task := range tasks {
		uc.workerPool.Submit(task)
	}

	results := make([]model.Result, 0, len(tasks))
	for i := 0; i < len(tasks); i++ {
		result := <-uc.workerPool.Results()
		results = append(results, result)
	}

	uc.workerPool.Stop()

	uc.logger.WithFields(logrus.Fields{
		"inputFile": inputFile,
		"taskCount": len(results),
	}).Info("Hoàn thành xử lý task")

	return results, nil
}
