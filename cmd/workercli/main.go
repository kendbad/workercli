package main

import (
	"flag"
	"log"
	"sync"
	"workercli/internal/adapter/input"
	proxyiface "workercli/internal/adapter/proxy"
	"workercli/internal/config"
	"workercli/internal/domain/model"
	"workercli/internal/infrastructure/proxy"
	"workercli/internal/infrastructure/task"
	"workercli/internal/infrastructure/tui"
	"workercli/internal/usecase"
	"workercli/pkg/utils"
)

func main() {
	tuiMode := flag.String("tui", "tview", "Loại giao diện TUI: tview, bubbletea, termui")
	checkProxy := flag.Bool("proxy", true, "Kiểm tra proxy từ proxy.txt")
	checkTask := flag.Bool("task", false, "Kiểm tra task từ tasks.txt")
	checkEmail := flag.Bool("email", false, "Kiểm tra email từ emails.txt")
	clientType := flag.String("client", "nethttp", "Thư viện gửi request: fasthttp, nethttp")
	flag.Parse()

	cfg, err := config.Load("configs/")
	if err != nil {
		log.Fatalf("Không thể tải cấu hình: %v", err)
	}

	logger, err := utils.NewLogger("configs/logger.yaml")
	if err != nil {
		log.Fatalf("Không thể khởi tạo logger: %v", err)
	}

	logger.Info("Ứng dụng WorkerCLI khởi động")

	logFile, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Không thể tạo file log: %v", err)
	}
	defer logFile.Close()

	if *checkProxy {
		fileReader := proxy.NewFileReader(logger)
		proxyReader := proxyiface.NewProxyReader(logger, fileReader)
		ipChecker := proxy.NewIPChecker(logger, *clientType)
		proxyChecker := proxyiface.NewProxyChecker(logger, ipChecker)
		proxyUsecase := usecase.NewProxyCheckUseCase(proxyReader, proxyChecker, cfg.Proxy.CheckURL, cfg.Worker.Workers, logger)

		if *tuiMode != "" {
			logger.SetOutput(logFile)
			factory := tui.NewRendererFactory(logger, *tuiMode)
			renderer := factory.CreateProxyRenderer(logger, &[]model.ProxyResult{}, &sync.Mutex{}, make(chan model.ProxyResult, 100), make(chan struct{}))
			tuiUsecase := tui.NewTUIUseCase(logger, *tuiMode, renderer)

			if err := tuiUsecase.Start(); err != nil {
				logger.Errorf("Không thể khởi động TUI: %v", err)
				log.Fatalf("Không thể khởi động TUI: %v", err)
			}

			results, err := proxyUsecase.Execute("input/proxy.txt")
			if err != nil {
				logger.Errorf("Lỗi kiểm tra proxy: %v", err)
				log.Fatalf("Lỗi kiểm tra proxy: %v", err)
			}
			for _, result := range results {
				logger.Infof("Sending proxy result to TUI: %v", result)
				tuiUsecase.AddProxyResult(result)
			}
			tuiUsecase.Close()

			logger.Infof("Hoàn thành kiểm tra proxy! Số proxy: %d", len(results))
		} else {
			results, err := proxyUsecase.Execute("input/proxy.txt")
			if err != nil {
				logger.Errorf("Lỗi kiểm tra proxy: %v", err)
				log.Fatalf("Lỗi kiểm tra proxy: %v", err)
			}
			for _, result := range results {
				status := result.Status
				if result.Error != "" {
					status += " (" + result.Error + ")"
				}
				log.Printf("Proxy %s://%s:%s, IP: %s, Status: %s\n",
					result.Proxy.Protocol, result.Proxy.IP, result.Proxy.Port, result.IP, status)
			}
			logger.Infof("Hoàn thành kiểm tra proxy! Số proxy: %d", len(results))
		}
	} else if *checkTask {
		processor := task.NewProcessor(logger)
		inputReader := input.NewFileReader(logger)
		batchTask := usecase.NewBatchTaskUseCase(inputReader, processor, cfg.Worker.Workers, logger)

		if *tuiMode != "" {
			logger.SetOutput(logFile)
			factory := tui.NewRendererFactory(logger, *tuiMode)
			renderer := factory.CreateTaskRenderer(logger, &[]model.Result{}, &sync.Mutex{}, make(chan model.Result, 100), make(chan struct{}))
			tuiUsecase := tui.NewTUIUseCase(logger, *tuiMode, renderer)

			if err := tuiUsecase.Start(); err != nil {
				logger.Errorf("Không thể khởi động TUI: %v", err)
				log.Fatalf("Không thể khởi động TUI: %v", err)
			}

			results, err := batchTask.Execute("input/tasks.txt")
			if err != nil {
				logger.Errorf("Lỗi xử lý tasks: %v", err)
				log.Fatalf("Lỗi xử lý tasks: %v", err)
			}
			for _, result := range results {
				logger.Infof("Sending task result to TUI: %v", result)
				tuiUsecase.AddTaskResult(result)
			}
			tuiUsecase.Close()

			logger.Infof("Hoàn thành xử lý tasks! Số task: %d", len(results))
		} else {
			results, err := batchTask.Execute("input/tasks.txt")
			if err != nil {
				logger.Errorf("Lỗi xử lý tasks: %v", err)
				log.Fatalf("Lỗi xử lý tasks: %v", err)
			}
			for _, result := range results {
				log.Printf("Task %s: %s\n", result.TaskID, result.Status)
			}
			logger.Infof("Hoàn thành xử lý tasks! Số task: %d", len(results))
		}
	} else if *checkEmail {
		// inputReader := input.NewFileReader(logger)
		// emailChecker := emailchecker.NewEmailChecker(*clientType, logger)
		// emailUsecase := usecase.NewEmailCheckUseCase(inputReader, emailChecker, cfg.Email.CheckURL, cfg.Worker.Workers, logger)

		// if *tuiMode != "" {
		// 	logger.SetOutput(logFile)
		// 	factory := tui.NewRendererFactory(logger, *tuiMode)
		// 	renderer := factory.CreateTaskRenderer(logger, &[]model.Result{}, &sync.Mutex{}, make(chan model.Result, 100), make(chan struct{}))
		// 	tuiUsecase := tui.NewTUIUseCase(logger, *tuiMode, renderer)

		// 	if err := tuiUsecase.Start(); err != nil {
		// 		logger.Errorf("Không thể khởi động TUI: %v", err)
		// 		log.Fatalf("Không thể khởi động TUI: %v", err)
		// 	}

		// 	results, err := emailUsecase.Execute("input/emails.txt")
		// 	if err != nil {
		// 		logger.Errorf("Lỗi kiểm tra email: %v", err)
		// 		log.Fatalf("Lỗi kiểm tra email: %v", err)
		// 	}
		// 	for _, result := range results {
		// 		logger.Infof("Sending email result to TUI: %v", result)
		// 		tuiUsecase.AddTaskResult(result)
		// 	}
		// 	tuiUsecase.Close()

		// 	logger.Infof("Hoàn thành kiểm tra email! Số email: %d", len(results))
		// } else {
		// 	results, err := emailUsecase.Execute("input/emails.txt")
		// 	if err != nil {
		// 		logger.Errorf("Lỗi kiểm tra email: %v", err)
		// 		log.Fatalf("Lỗi kiểm tra email: %v", err)
		// 	}
		// 	for _, result := range results {
		// 		log.Printf("Email %s: %s\n", result.TaskID, result.Status)
		// 	}
		// 	logger.Infof("Hoàn thành kiểm tra email! Số email: %d", len(results))
		// }
	} else {
		logger.Info("Không có tùy chọn nào được chọn (-proxy, -task, hoặc -email)")
	}
}
