package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// Adapter layer
	"workercli/internal/adapter/input"
	proxyiface "workercli/internal/adapter/proxy"

	// Config & Domain layers
	"workercli/internal/config"
	"workercli/internal/domain/model"

	// Infrastructure layer
	"workercli/internal/infrastructure/proxy"
	"workercli/internal/infrastructure/task"
	tuiinfra "workercli/internal/infrastructure/tui"

	// Usecase layer
	"workercli/internal/usecase"

	// Utilities
	"workercli/pkg/utils"
)

// Application represents the fully configured application
type Application struct {
	config      *config.Config
	logger      *utils.Logger
	tuiUseCase  *tuiinfra.TUIUseCase
	stopChannel chan struct{}
}

// NewApplication creates and configures the application
func NewApplication(configDir string) (*Application, error) {
	// 1. Load configuration
	cfg, err := config.Load(configDir)
	if err != nil {
		return nil, err
	}

	// 2. Initialize logger
	logger, err := utils.NewLogger("configs/logger.yaml")
	if err != nil {
		return nil, err
	}

	app := &Application{
		config:      cfg,
		logger:      logger,
		stopChannel: make(chan struct{}),
	}

	return app, nil
}

// setupProxyChecking creates and configures all the components needed for proxy checking
func (app *Application) setupProxyChecking(clientType, tuiMode string) (*usecase.ProxyCheckUseCase, error) {
	// Setup proxy infrastructure
	fileReader := proxy.NewFileReader(app.logger)
	proxyReader := proxyiface.NewProxyReader(app.logger, fileReader)

	// Create IP checker via proxy implementation that implements the Checker interface
	proxyIPChecker := proxy.NewIPChecker(app.logger, clientType)
	proxyChecker := proxyiface.NewProxyChecker(app.logger, proxyIPChecker)

	// Create proxy check usecase
	proxyUsecase := usecase.NewProxyCheckUseCase(
		proxyReader,
		proxyChecker,
		app.config.Proxy.CheckURL,
		app.config.Worker.Workers,
		app.logger,
	)

	// Setup TUI if needed
	if tuiMode != "" {
		factory := tuiinfra.NewRendererFactory(app.logger, tuiMode)
		renderer := factory.CreateProxyRenderer(
			app.logger,
			&[]model.ProxyResult{},
			&sync.Mutex{},
			make(chan model.ProxyResult, 100),
			app.stopChannel,
		)
		app.tuiUseCase = tuiinfra.NewTUIUseCase(app.logger, tuiMode, renderer)
	}

	return proxyUsecase, nil
}

// setupTaskProcessing creates and configures all the components needed for task processing
func (app *Application) setupTaskProcessing(tuiMode string) (*usecase.BatchTaskUseCase, error) {
	// Task infrastructure
	processor := task.NewProcessor(app.logger)
	inputReader := input.NewFileReader(app.logger)

	// Create batch task usecase
	batchTask := usecase.NewBatchTaskUseCase(
		inputReader,
		processor,
		app.config.Worker.Workers,
		app.logger,
	)

	// Setup TUI if needed
	if tuiMode != "" {
		factory := tuiinfra.NewRendererFactory(app.logger, tuiMode)
		renderer := factory.CreateTaskRenderer(
			app.logger,
			&[]model.Result{},
			&sync.Mutex{},
			make(chan model.Result, 100),
			app.stopChannel,
		)
		app.tuiUseCase = tuiinfra.NewTUIUseCase(app.logger, tuiMode, renderer)
	}

	return batchTask, nil
}

// ExecuteProxyChecking performs the proxy checking workflow
func (app *Application) ExecuteProxyChecking(clientType, tuiMode string) error {
	proxyUsecase, err := app.setupProxyChecking(clientType, tuiMode)
	if err != nil {
		return err
	}

	// Create log file for TUI mode
	if tuiMode != "" {
		logFile, err := utils.CreateLogFile()
		if err != nil {
			return err
		}
		defer logFile.Close()
		app.logger.SetOutput(logFile)

		// Start TUI
		if err := app.tuiUseCase.Start(); err != nil {
			app.logger.Errorf("Unable to start TUI: %v", err)
			return err
		}
		defer app.tuiUseCase.Close()
	}

	// Execute the proxy check usecase
	results, err := proxyUsecase.Execute(app.config.Proxy.FilePath)
	if err != nil {
		app.logger.Errorf("Error checking proxies: %v", err)
		return err
	}

	// Display or send results
	if tuiMode != "" {
		for _, result := range results {
			app.logger.Infof("Sending proxy result to TUI: %v", result)
			app.tuiUseCase.AddProxyResult(result)
		}
	} else {
		for _, result := range results {
			status := result.Status
			if result.Error != "" {
				status += " (" + result.Error + ")"
			}
			log.Printf("Proxy %s://%s:%s, IP: %s, Status: %s\n",
				result.Proxy.Protocol, result.Proxy.IP, result.Proxy.Port, result.IP, status)
		}
	}

	app.logger.Infof("Proxy check completed! Total proxies: %d", len(results))
	return nil
}

// ExecuteTaskProcessing performs the task processing workflow
func (app *Application) ExecuteTaskProcessing(tuiMode string) error {
	batchTask, err := app.setupTaskProcessing(tuiMode)
	if err != nil {
		return err
	}

	// Create log file for TUI mode
	if tuiMode != "" {
		logFile, err := utils.CreateLogFile()
		if err != nil {
			return err
		}
		defer logFile.Close()
		app.logger.SetOutput(logFile)

		// Start TUI
		if err := app.tuiUseCase.Start(); err != nil {
			app.logger.Errorf("Unable to start TUI: %v", err)
			return err
		}
		defer app.tuiUseCase.Close()
	}

	// Execute the task processing usecase
	results, err := batchTask.Execute(app.config.Input.FilePath)
	if err != nil {
		app.logger.Errorf("Error processing tasks: %v", err)
		return err
	}

	// Display or send results
	if tuiMode != "" {
		for _, result := range results {
			app.logger.Infof("Sending task result to TUI: %v", result)
			app.tuiUseCase.AddTaskResult(result)
		}
	} else {
		for _, result := range results {
			log.Printf("Task %s: %s\n", result.TaskID, result.Status)
		}
	}

	app.logger.Infof("Task processing completed! Total tasks: %d", len(results))
	return nil
}

// Shutdown gracefully shuts down the application
func (app *Application) Shutdown() {
	app.logger.Info("Shutting down application...")
	if app.tuiUseCase != nil {
		app.tuiUseCase.Close()
	}
	close(app.stopChannel)
}

func main() {
	// Parse command line flags
	tuiMode := flag.String("tui", "", "TUI type: tview, bubbletea, termui")
	checkProxy := flag.Bool("proxy", false, "Check proxies from proxy.txt")
	checkTask := flag.Bool("task", false, "Process tasks from tasks.txt")
	clientType := flag.String("client", "nethttp", "HTTP client: fasthttp, nethttp")
	flag.Parse()

	// Create application
	app, err := NewApplication("configs/")
	if err != nil {
		log.Fatalf("Unable to initialize application: %v", err)
	}

	app.logger.Info("WorkerCLI application starting")

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		app.Shutdown()
		os.Exit(0)
	}()

	// Execute workflows based on flags
	if *checkProxy {
		if err := app.ExecuteProxyChecking(*clientType, *tuiMode); err != nil {
			log.Fatalf("Proxy checking failed: %v", err)
		}
	} else if *checkTask {
		if err := app.ExecuteTaskProcessing(*tuiMode); err != nil {
			log.Fatalf("Task processing failed: %v", err)
		}
	} else {
		app.logger.Info("No option selected (-proxy, -task)")
		log.Println("No option selected. Use -proxy to check proxies or -task to process tasks.")
	}
}
