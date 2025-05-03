package usecase

import (
	"fmt"
	"os"
	"strings"
	"workercli/internal/domain/model"
	"workercli/internal/domain/service"
	"workercli/internal/interface/proxy"
	"workercli/internal/interface/worker"
	"workercli/pkg/utils"
)

type ProxyCheckUseCase struct {
	reader     proxy.Reader
	checker    service.ProxyChecker
	checkURL   string
	workerPool *worker.Pool
	logger     *utils.Logger
}

func NewProxyCheckUseCase(reader proxy.Reader, checker service.ProxyChecker, checkURL string, workers int, logger *utils.Logger) *ProxyCheckUseCase {
	processor := &ProxyTaskProcessor{
		checker:  checker,
		checkURL: checkURL,
		logger:   logger,
	}
	pool := worker.NewPool(workers, processor, logger)
	return &ProxyCheckUseCase{
		reader:     reader,
		checker:    checker,
		checkURL:   checkURL,
		workerPool: pool,
		logger:     logger,
	}
}

// ProxyTaskProcessor processes model.Task for proxy checking
type ProxyTaskProcessor struct {
	checker  service.ProxyChecker
	checkURL string
	logger   *utils.Logger
}

// ProxyResult chứa thông tin về kết quả kiểm tra proxy
type ProxyResult struct {
	Proxy  Proxy  `json:"proxy"`
	IP     string `json:"ip"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Proxy định nghĩa thông tin proxy
type Proxy struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

func (p *ProxyTaskProcessor) ProcessTask(task model.Task) (model.Result, error) {
	// Parse TaskID (e.g., "http://192.168.1.1:8080") and Data (protocol)
	parts := strings.SplitN(task.TaskID, "://", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid proxy format: %s", task.TaskID)
		p.logger.Errorf("Invalid proxy format: %s", task.TaskID)
		return model.Result{
			TaskID: task.TaskID,
			Status: "Failed",
			Error:  err.Error(),
		}, err
	}
	protocol := parts[0]
	addrParts := strings.Split(parts[1], ":")
	if len(addrParts) != 2 {
		err := fmt.Errorf("invalid proxy address: %s", parts[1])
		p.logger.Errorf("Invalid proxy address: %s", parts[1])
		return model.Result{
			TaskID: task.TaskID,
			Status: "Failed",
			Error:  err.Error(),
		}, err
	}

	proxy := model.Proxy{
		Protocol: protocol,
		IP:       addrParts[0],
		Port:     addrParts[1],
	}

	result, err := p.checker.CheckProxy(proxy, p.checkURL)
	if err != nil {
		p.logger.Errorf("Lỗi kiểm tra proxy %s: %v", task.TaskID, err)
		return model.Result{
			TaskID: task.TaskID,
			Status: "Failed",
			Error:  err.Error(),
		}, err
	}

	p.logger.Infof("Proxy %s trả về IP: %s", task.TaskID, result.IP)
	return model.Result{
		TaskID: task.TaskID,
		Status: "Success",
	}, nil
}

func (uc *ProxyCheckUseCase) Execute(proxyFile string) ([]model.ProxyResult, error) {
	proxies, err := uc.reader.ReadProxies(proxyFile)
	if err != nil {
		uc.logger.Errorf("Không thể đọc proxy: %v", err)
		return nil, err
	}

	uc.workerPool.Start()
	results := make([]model.ProxyResult, 0, len(proxies))
	proxyResultCh := make(chan model.ProxyResult, len(proxies))

	// Submit tasks to worker pool
	for _, p := range proxies {
		proxy := p
		task := model.Task{
			TaskID: fmt.Sprintf("%s://%s:%s", proxy.Protocol, proxy.IP, proxy.Port),
			Data:   proxy.Protocol,
		}
		uc.workerPool.Submit(task)

		// Collect results in a goroutine
		go func(proxy model.Proxy) {
			taskResult := <-uc.workerPool.Results()
			proxyResult := model.ProxyResult{
				Proxy:  proxy,
				Status: taskResult.Status,
				Error:  taskResult.Error,
			}
			if taskResult.Status == "Success" {
				// Re-run CheckProxy to get IP
				pr, err := uc.checker.CheckProxy(proxy, uc.checkURL)
				if err == nil {
					proxyResult.IP = pr.IP
				}
			}
			proxyResultCh <- proxyResult
		}(proxy)
	}

	// Collect results
	for i := 0; i < len(proxies); i++ {
		results = append(results, <-proxyResultCh)
	}

	uc.workerPool.Stop()

	// Save results to file
	outputFile := utils.AutoPath("output/proxy_results.txt")
	file, err := os.Create(outputFile)
	if err != nil {
		uc.logger.Errorf("Không thể tạo file kết quả: %v", err)
		return results, err
	}
	defer file.Close()

	for _, result := range results {
		status := result.Status
		if result.Error != "" {
			status += " (" + result.Error + ")"
		}
		_, err := fmt.Fprintf(file, "Proxy: %s://%s:%s, IP: %s, Status: %s\n",
			result.Proxy.Protocol, result.Proxy.IP, result.Proxy.Port, result.IP, status)
		if err != nil {
			uc.logger.Errorf("Lỗi ghi kết quả: %v", err)
		}
	}

	uc.logger.Infof("Đã lưu kết quả vào %s", outputFile)
	return results, nil
}
