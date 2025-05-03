package usecase

import (
	"fmt"
	"os"
	"strings"
	"workercli/internal/adapter/proxy"  // Gói xử lý đọc danh sách proxy
	"workercli/internal/adapter/worker" // Gói quản lý pool worker xử lý tác vụ
	"workercli/internal/domain/model"   // Gói định nghĩa các struct dữ liệu
	"workercli/internal/domain/service" // Gói cung cấp dịch vụ kiểm tra proxy
	"workercli/pkg/utils"               // Gói tiện ích (logger, auto path,...)
)

// ProxyCheckUseCase là struct chính chịu trách nhiệm thực hiện kiểm tra proxy
type ProxyCheckUseCase struct {
	reader     proxy.Reader         // Đối tượng đọc danh sách proxy từ file
	checker    service.ProxyChecker // Dịch vụ kiểm tra proxy
	checkURL   string               // URL dùng để kiểm tra proxy
	workerPool *worker.Pool         // Pool worker xử lý tác vụ song song
	logger     *utils.Logger        // Logger ghi log hoạt động
}

// NewProxyCheckUseCase khởi tạo một ProxyCheckUseCase mới
// - reader: Đối tượng đọc proxy
// - checker: Dịch vụ kiểm tra proxy
// - checkURL: URL kiểm tra
// - workers: Số lượng worker trong pool
// - logger: Đối tượng logger
func NewProxyCheckUseCase(reader proxy.Reader, checker service.ProxyChecker, checkURL string, workers int, logger *utils.Logger) *ProxyCheckUseCase {
	// Tạo processor để xử lý tác vụ kiểm tra proxy
	processor := &ProxyTaskProcessor{checker, checkURL, logger}
	// Tạo worker pool với số lượng worker và processor
	return &ProxyCheckUseCase{
		reader:     reader,
		checker:    checker,
		checkURL:   checkURL,
		workerPool: worker.NewPool(workers, processor, logger),
		logger:     logger,
	}
}

// ProxyTaskProcessor xử lý các tác vụ kiểm tra proxy
type ProxyTaskProcessor struct {
	checker  service.ProxyChecker // Dịch vụ kiểm tra proxy
	checkURL string               // URL dùng để kiểm tra
	logger   *utils.Logger        // Logger ghi log
}

// ProcessTask xử lý một tác vụ kiểm tra proxy
// - task: Tác vụ chứa thông tin proxy (TaskID dạng "protocol://ip:port")
// Trả về kết quả kiểm tra hoặc lỗi nếu có
func (p *ProxyTaskProcessor) ProcessTask(task model.Task) (model.Result, error) {
	// Tách TaskID thành protocol và địa chỉ (ip:port)
	parts := strings.SplitN(task.TaskID, "://", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid proxy format: %s", task.TaskID)
		p.logger.Errorf("Invalid proxy format: %s", task.TaskID)
		return model.Result{TaskID: task.TaskID, Status: "Failed", Error: err.Error()}, err
	}

	// Tách địa chỉ thành IP và Port
	addrParts := strings.Split(parts[1], ":")
	if len(addrParts) != 2 {
		err := fmt.Errorf("invalid proxy address: %s", parts[1])
		p.logger.Errorf("Invalid proxy address: %s", parts[1])
		return model.Result{TaskID: task.TaskID, Status: "Failed", Error: err.Error()}, err
	}

	// Tạo struct proxy từ thông tin đã tách
	proxy := model.Proxy{Protocol: parts[0], IP: addrParts[0], Port: addrParts[1]}
	// Kiểm tra proxy bằng checker
	result, err := p.checker.CheckProxy(proxy, p.checkURL)
	if err != nil {
		p.logger.Errorf("Proxy check failed %s: %v", task.TaskID, err)
		return model.Result{TaskID: task.TaskID, Status: "Failed", Error: err.Error()}, err
	}

	// Ghi log thành công và trả về kết quả
	p.logger.Infof("Proxy %s returned IP: %s", task.TaskID, result.IP)
	return model.Result{TaskID: task.TaskID, Status: "Success"}, nil
}

// Execute thực hiện quá trình kiểm tra danh sách proxy từ file
// - proxyFile: Đường dẫn file chứa danh sách proxy
// Trả về danh sách kết quả và lỗi nếu có
func (uc *ProxyCheckUseCase) Execute(proxyFile string) ([]model.ProxyResult, error) {
	// Đọc danh sách proxy từ file
	proxies, err := uc.reader.ReadProxies(proxyFile)
	if err != nil {
		uc.logger.Errorf("Failed to read proxies: %v", err)
		return nil, err
	}

	// Khởi động worker pool
	uc.workerPool.Start()
	results := make([]model.ProxyResult, 0, len(proxies))
	proxyResultCh := make(chan model.ProxyResult, len(proxies)) // Kênh nhận kết quả

	// Gửi từng proxy vào worker pool để xử lý
	for _, proxy := range proxies {
		task := model.Task{
			TaskID: fmt.Sprintf("%s://%s:%s", proxy.Protocol, proxy.IP, proxy.Port),
			Data:   proxy.Protocol,
		}
		uc.workerPool.Submit(task)

		// Goroutine thu thập kết quả từ worker pool
		go func(p model.Proxy) {
			taskResult := <-uc.workerPool.Results() // Nhận kết quả từ pool
			result := model.ProxyResult{Proxy: p, Status: taskResult.Status, Error: taskResult.Error}
			if taskResult.Status == "Success" {
				// Kiểm tra lại proxy để lấy IP
				if pr, err := uc.checker.CheckProxy(p, uc.checkURL); err == nil {
					result.IP = pr.IP
				} else {
					result.IP = "Failed Get IP"
				}
			} else {
				result.IP = "Failed Send Task"
			}
			proxyResultCh <- result // Gửi kết quả vào kênh
		}(proxy)
	}

	// Thu thập tất cả kết quả từ kênh
	for i := 0; i < len(proxies); i++ {
		results = append(results, <-proxyResultCh)
	}

	// Dừng worker pool
	uc.workerPool.Stop()

	// Lưu kết quả vào file
	outputFile := utils.AutoPath("output/proxy_results.txt")
	file, err := os.Create(outputFile)
	if err != nil {
		uc.logger.Errorf("Failed to create result file: %v", err)
		return results, err
	}
	defer file.Close()

	// Ghi từng kết quả vào file
	for _, r := range results {
		status := r.Status
		if r.Error != "" {
			status += " (" + r.Error + ")"
		}
		fmt.Fprintf(file, "Proxy: %s://%s:%s, IP: %s, Status: %s\n",
			r.Proxy.Protocol, r.Proxy.IP, r.Proxy.Port, r.IP, status)
	}

	// Ghi log lưu file thành công
	uc.logger.Infof("Results saved to %s", outputFile)
	return results, nil
}
