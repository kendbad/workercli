package usecase

import (
	"fmt"
	"os"
	"strings"
	"workercli/internal/adapter/proxy"
	"workercli/internal/adapter/worker"
	"workercli/internal/domain/model"
	proxi "workercli/internal/infrastructure/proxy"
	"workercli/pkg/utils"
)

type ProxyCheckUseCase struct {
	reader     *proxy.ProxyReader
	checker    *proxy.ProxyChecker
	checkURL   string
	workerPool *worker.Pool
	logger     *utils.Logger
}

func NewProxyCheckUseCase(reader *proxy.ProxyReader, checker *proxy.ProxyChecker, checkURL string, workers int, logger *utils.Logger) *ProxyCheckUseCase {
	// Đảm bảo checkURL có giao thức
	if !strings.HasPrefix(checkURL, "http://") && !strings.HasPrefix(checkURL, "https://") {
		checkURL = "http://" + checkURL
	}
	processor := &ProxyTaskProcessor{checker, checkURL, logger}
	return &ProxyCheckUseCase{
		reader:     reader,
		checker:    checker,
		checkURL:   checkURL,
		workerPool: worker.NewPool(workers, processor, logger),
		logger:     logger,
	}
}

type ProxyTaskProcessor struct {
	checker  *proxy.ProxyChecker
	checkURL string
	logger   *utils.Logger
}

func (p *ProxyTaskProcessor) ProcessTask(task model.Task) (model.Result, error) {
	proxy, err := proxi.ParseProxy(task.TaskID)
	if err != nil {
		p.logger.Errorf("Invalid proxy format: %s", task.TaskID)
		return model.Result{TaskID: task.TaskID, Status: "Failed", Error: err.Error()}, err
	}

	ip, status, err := p.checker.CheckProxy(proxy, p.checkURL)
	if err != nil {
		p.logger.Errorf("Proxy check failed %s: %v", task.TaskID, err)
		return model.Result{TaskID: task.TaskID, Status: status, Error: err.Error()}, err
	}

	p.logger.Infof("Proxy %s returned IP: %s", task.TaskID, ip)
	return model.Result{TaskID: task.TaskID, Status: status}, nil
}

func (uc *ProxyCheckUseCase) Execute(proxyFile string) ([]model.ProxyResult, error) {
	proxies, err := uc.reader.ReadProxies(proxyFile)
	if err != nil {
		uc.logger.Errorf("Failed to read proxies: %v", err)
		return nil, err
	}

	uc.workerPool.Start()
	results := make([]model.ProxyResult, 0, len(proxies))
	proxyResultCh := make(chan model.ProxyResult, len(proxies))

	for _, p := range proxies {
		task := model.Task{
			TaskID: fmt.Sprintf("%s://%s:%s", p.Protocol, p.IP, p.Port),
			Data:   p.Protocol,
		}
		uc.workerPool.Submit(task)

		go func(pr model.Proxy) {
			taskResult := <-uc.workerPool.Results()
			result := model.ProxyResult{Proxy: pr, Status: taskResult.Status, Error: taskResult.Error}
			if taskResult.Status == "Success" {
				if ip, _, err := uc.checker.CheckProxy(pr, uc.checkURL); err == nil {
					result.IP = ip
				} else {
					result.IP = "Failed Get IP"
				}
			} else {
				result.IP = "Failed Send Task"
			}
			proxyResultCh <- result
		}(p)
	}

	for i := 0; i < len(proxies); i++ {
		results = append(results, <-proxyResultCh)
	}

	uc.workerPool.Stop()

	outputFile := utils.AutoPath("output/proxy_results.txt")
	file, err := os.Create(outputFile)
	if err != nil {
		uc.logger.Errorf("Failed to create result file: %v", err)
		return results, err
	}
	defer file.Close()

	for _, r := range results {
		status := r.Status
		if r.Error != "" {
			status += " (" + r.Error + ")"
		}
		fmt.Fprintf(file, "Proxy: %s://%s:%s, IP: %s, Status: %s\n",
			r.Proxy.Protocol, r.Proxy.IP, r.Proxy.Port, r.IP, status)
	}

	uc.logger.Infof("Results saved to %s", outputFile)
	return results, nil
}
