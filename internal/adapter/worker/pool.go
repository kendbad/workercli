package worker

import (
	"sync"
	"workercli/internal/domain/model"
	"workercli/internal/domain/service"
	"workercli/pkg/utils"
)

// Pool quản lý một tập hợp các worker
type Pool struct {
	tasks     chan model.Task
	results   chan model.Result
	processor service.TaskProcessor
	workers   int
	logger    *utils.Logger
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

// NewPool tạo một pool worker mới
func NewPool(workers int, processor service.TaskProcessor, logger *utils.Logger) *Pool {
	return &Pool{
		tasks:     make(chan model.Task, 1000), // Queue lớn cho số lượng task lớn
		results:   make(chan model.Result, 1000),
		processor: processor,
		workers:   workers,
		logger:    logger,
		stopCh:    make(chan struct{}),
	}
}

// Start khởi động tất cả worker
func (p *Pool) Start() {
	p.logger.Infof("Khởi động pool với %d worker", p.workers)
	p.wg.Add(p.workers)
	for i := 0; i < p.workers; i++ {
		worker := NewWorker(i+1, p.tasks, p.results, p.processor, p.logger)
		go func() {
			defer p.wg.Done()
			worker.Run(p.stopCh)
		}()
	}
}

// Submit gửi một task vào queue
func (p *Pool) Submit(task model.Task) {
	p.tasks <- task
}

// Results trả về channel kết quả
func (p *Pool) Results() <-chan model.Result {
	return p.results
}

// Stop dừng tất cả worker
func (p *Pool) Stop() {
	p.logger.Info("Dừng pool worker")
	close(p.stopCh)
	close(p.tasks)
	p.wg.Wait()
	close(p.results)
}
