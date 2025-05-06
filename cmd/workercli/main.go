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

// thietLapKiemTraTrungGian creates and configures all the components needed for proxy checking
func (app *Application) thietLapKiemTraTrungGian(loaiKetNoi, kieuGiaoDien string) (*usecase.KiemTraTrungGian, error) {
	// Setup proxy infrastructure
	boDocTep := proxy.NewFileReader(app.logger)
	boDocTrungGian := proxyiface.NewProxyReader(app.logger, boDocTep)

	// Create IP checker via proxy implementation that implements the Checker interface
	boKiemTraIP := proxy.NewIPChecker(app.logger, loaiKetNoi)
	boKiemTraTrungGian := proxyiface.NewProxyChecker(app.logger, boKiemTraIP)

	// Create proxy check usecase
	boKiemTra := usecase.TaoBoKiemTraTrungGian(
		boDocTrungGian,
		boKiemTraTrungGian,
		app.config.Proxy.CheckURL,
		app.config.Worker.Workers,
		app.logger,
	)

	// Setup TUI if needed
	if kieuGiaoDien != "" {
		nhaXuongGiaoDien := tuiinfra.NewRendererFactory(app.logger, kieuGiaoDien)
		boHienThi := nhaXuongGiaoDien.CreateProxyRenderer(
			app.logger,
			&[]model.KetQuaTrungGian{},
			&sync.Mutex{},
			make(chan model.KetQuaTrungGian, 100),
			app.stopChannel,
		)
		app.tuiUseCase = tuiinfra.NewTUIUseCase(app.logger, kieuGiaoDien, boHienThi)
	}

	return boKiemTra, nil
}

// thietLapXuLyTacVu creates and configures all the components needed for task processing
func (app *Application) thietLapXuLyTacVu(kieuGiaoDien string) (*usecase.XuLyLoDongTacVu, error) {
	// Task infrastructure
	boXuLy := task.NewProcessor(app.logger)
	boDocDauVao := input.NewFileReader(app.logger)

	// Create batch task usecase
	boXuLyLoDong := usecase.TaoBoXuLyLoDongTacVu(
		boDocDauVao,
		boXuLy,
		app.config.Worker.Workers,
		app.logger,
	)

	// Setup TUI if needed
	if kieuGiaoDien != "" {
		nhaXuongGiaoDien := tuiinfra.NewRendererFactory(app.logger, kieuGiaoDien)
		boHienThi := nhaXuongGiaoDien.CreateTaskRenderer(
			app.logger,
			&[]model.KetQua{},
			&sync.Mutex{},
			make(chan model.KetQua, 100),
			app.stopChannel,
		)
		app.tuiUseCase = tuiinfra.NewTUIUseCase(app.logger, kieuGiaoDien, boHienThi)
	}

	return boXuLyLoDong, nil
}

// ThucThiKiemTraTrungGian performs the proxy checking workflow
func (app *Application) ThucThiKiemTraTrungGian(loaiKetNoi, kieuGiaoDien string) error {
	boKiemTra, err := app.thietLapKiemTraTrungGian(loaiKetNoi, kieuGiaoDien)
	if err != nil {
		return err
	}

	// Create log file for TUI mode
	if kieuGiaoDien != "" {
		tepGhiNhatKy, err := utils.CreateLogFile()
		if err != nil {
			return err
		}
		defer tepGhiNhatKy.Close()
		app.logger.SetOutput(tepGhiNhatKy)

		// Start TUI
		if err := app.tuiUseCase.Start(); err != nil {
			app.logger.Errorf("Không thể khởi động TUI: %v", err)
			return err
		}
		defer app.tuiUseCase.Close()
	}

	// Execute the proxy check usecase
	ketQua, err := boKiemTra.ThucThi(app.config.Proxy.FilePath)
	if err != nil {
		app.logger.Errorf("Lỗi kiểm tra proxy: %v", err)
		return err
	}

	// Display or send results
	if kieuGiaoDien != "" {
		for _, kq := range ketQua {
			app.logger.Infof("Gửi kết quả proxy vào TUI: %v", kq)
			app.tuiUseCase.AddProxyResult(kq)
		}
	} else {
		for _, kq := range ketQua {
			trangThai := kq.TrangThai
			if kq.LoiXayRa != "" {
				trangThai += " (" + kq.LoiXayRa + ")"
			}
			log.Printf("Proxy %s://%s:%s, IP: %s, Trạng thái: %s\n",
				kq.TrungGian.GiaoDien, kq.TrungGian.DiaChi, kq.TrungGian.Cong, kq.DiaChi, trangThai)
		}
	}

	app.logger.Infof("Kiểm tra proxy hoàn tất! Tổng số proxy: %d", len(ketQua))
	return nil
}

// ThucThiXuLyTacVu performs the task processing workflow
func (app *Application) ThucThiXuLyTacVu(kieuGiaoDien string) error {
	boXuLyLoDong, err := app.thietLapXuLyTacVu(kieuGiaoDien)
	if err != nil {
		return err
	}

	// Create log file for TUI mode
	if kieuGiaoDien != "" {
		tepGhiNhatKy, err := utils.CreateLogFile()
		if err != nil {
			return err
		}
		defer tepGhiNhatKy.Close()
		app.logger.SetOutput(tepGhiNhatKy)

		// Start TUI
		if err := app.tuiUseCase.Start(); err != nil {
			app.logger.Errorf("Không thể khởi động TUI: %v", err)
			return err
		}
		defer app.tuiUseCase.Close()
	}

	// Execute the task processing usecase
	ketQua, err := boXuLyLoDong.ThucThi(app.config.Input.FilePath)
	if err != nil {
		app.logger.Errorf("Lỗi xử lý các tác vụ: %v", err)
		return err
	}

	// Display or send results
	if kieuGiaoDien != "" {
		for _, kq := range ketQua {
			app.logger.Infof("Gửi kết quả tác vụ vào TUI: %v", kq)
			app.tuiUseCase.AddTaskResult(kq)
		}
	} else {
		for _, kq := range ketQua {
			log.Printf("Tác vụ %s: %s\n", kq.MaTacVu, kq.TrangThai)
		}
	}

	app.logger.Infof("Xử lý tác vụ hoàn tất! Tổng số tác vụ: %d", len(ketQua))
	return nil
}

// DungUngDung gracefully shuts down the application
func (app *Application) DungUngDung() {
	app.logger.Info("Đang dừng ứng dụng...")
	if app.tuiUseCase != nil {
		app.tuiUseCase.Close()
	}
	close(app.stopChannel)
}

func main() {
	// Parse command line flags
	kieuGiaoDien := flag.String("tui", "", "Loại TUI: tview, bubbletea, termui")
	kiemTraProxy := flag.Bool("proxy", false, "Kiểm tra proxy từ tệp proxy.txt")
	xuLyTacVu := flag.Bool("task", false, "Xử lý tác vụ từ tệp tasks.txt")
	loaiKetNoi := flag.String("client", "nethttp", "Loại HTTP client: fasthttp, nethttp")
	flag.Parse()

	// Create application
	app, err := NewApplication("configs/")
	if err != nil {
		log.Fatalf("Không thể khởi tạo ứng dụng: %v", err)
	}

	app.logger.Info("Ứng dụng WorkerCLI đang khởi động")

	// Setup signal handling for graceful shutdown
	kenhTinHieu := make(chan os.Signal, 1)
	signal.Notify(kenhTinHieu, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-kenhTinHieu
		app.DungUngDung()
		os.Exit(0)
	}()

	// Execute workflows based on flags
	if *kiemTraProxy {
		if err := app.ThucThiKiemTraTrungGian(*loaiKetNoi, *kieuGiaoDien); err != nil {
			log.Fatalf("Kiểm tra proxy thất bại: %v", err)
		}
	} else if *xuLyTacVu {
		if err := app.ThucThiXuLyTacVu(*kieuGiaoDien); err != nil {
			log.Fatalf("Xử lý tác vụ thất bại: %v", err)
		}
	} else {
		app.logger.Info("Không có tùy chọn nào được chọn (-proxy, -task)")
		log.Println("Không có tùy chọn nào được chọn. Sử dụng -proxy để kiểm tra proxy hoặc -task để xử lý tác vụ.")
	}
}
