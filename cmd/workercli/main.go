package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// Tầng Adapter
	"workercli/internal/adapter/input"
	proxyiface "workercli/internal/adapter/proxy"

	// Tầng Cấu hình & Domain
	"workercli/internal/config"
	"workercli/internal/domain/model"

	// Tầng Cơ sở hạ tầng
	"workercli/internal/infrastructure/proxy"
	"workercli/internal/infrastructure/task"
	tuiinfra "workercli/internal/infrastructure/tui"

	// Tầng Trường hợp sử dụng usecase
	"workercli/internal/usecase"

	// Tiện ích
	"workercli/pkg/utils"
)

// Application đại diện cho ứng dụng đã được cấu hình đầy đủ
type Application struct {
	config      *config.Config
	logger      *utils.Logger
	tuiUseCase  *tuiinfra.TUIUseCase
	stopChannel chan struct{}
}

// NewApplication tạo và cấu hình ứng dụng
func NewApplication(configDir string) (*Application, error) {
	// 1. Tải cấu hình
	cfg, err := config.Load(configDir)
	if err != nil {
		return nil, err
	}

	// 2. Khởi tạo logger
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

// thietLapKiemTraProxy tạo và cấu hình tất cả các thành phần cần thiết cho việc kiểm tra proxy
func (app *Application) thietLapKiemTraProxy(loaiKetNoi, kieuGiaoDien string) (*usecase.KiemTraProxy, error) {
	// Thiết lập cơ sở hạ tầng proxy
	boDocTep := proxy.NewFileReader(app.logger)
	boDocProxy := proxyiface.NewProxyReader(app.logger, boDocTep)

	// Tạo bộ kiểm tra IP qua proxy triển khai giao diện Checker
	boKiemTraIP := proxy.NewIPChecker(app.logger, loaiKetNoi)
	boKiemTraProxy := proxyiface.NewProxyChecker(app.logger, boKiemTraIP)

	// Tạo usecase kiểm tra proxy
	boKiemTra := usecase.TaoBoKiemTraProxy(
		boDocProxy,
		boKiemTraProxy,
		app.config.Proxy.CheckURL,
		app.config.Worker.Workers,
		app.logger,
	)

	// Thiết lập TUI nếu cần
	if kieuGiaoDien != "" {
		nhaXuongGiaoDien := tuiinfra.NewRendererFactory(app.logger, kieuGiaoDien)
		boHienThi := nhaXuongGiaoDien.CreateProxyRenderer(
			app.logger,
			&[]model.KetQuaProxy{},
			&sync.Mutex{},
			make(chan model.KetQuaProxy, 100),
			app.stopChannel,
		)
		app.tuiUseCase = tuiinfra.NewTUIUseCase(app.logger, kieuGiaoDien, boHienThi)
	}

	return boKiemTra, nil
}

// thietLapXuLyTacVu tạo và cấu hình tất cả các thành phần cần thiết cho việc xử lý tác vụ
func (app *Application) thietLapXuLyTacVu(kieuGiaoDien string) (*usecase.XuLyLoDongTacVu, error) {
	// Cơ sở hạ tầng tác vụ
	boXuLy := task.NewProcessor(app.logger)
	boDocDauVao := input.NewFileReader(app.logger)

	// Tạo usecase xử lý tác vụ hàng loạt
	boXuLyLoDong := usecase.TaoBoXuLyLoDongTacVu(
		boDocDauVao,
		boXuLy,
		app.config.Worker.Workers,
		app.logger,
	)

	// Thiết lập TUI nếu cần
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

// ThucThiKiemTraProxy thực hiện quy trình kiểm tra proxy
func (app *Application) ThucThiKiemTraProxy(loaiKetNoi, kieuGiaoDien string) error {
	boKiemTra, err := app.thietLapKiemTraProxy(loaiKetNoi, kieuGiaoDien)
	if err != nil {
		return err
	}

	// Tạo tệp log cho chế độ TUI
	if kieuGiaoDien != "" {
		tepGhiNhatKy, err := utils.CreateLogFile()
		if err != nil {
			return err
		}
		defer tepGhiNhatKy.Close()
		app.logger.SetOutput(tepGhiNhatKy)

		// Khởi động TUI
		if err := app.tuiUseCase.Start(); err != nil {
			app.logger.Errorf("Không thể khởi động TUI: %v", err)
			return err
		}
		defer app.tuiUseCase.Close()
	}

	// Thực thi usecase kiểm tra proxy
	ketQua, err := boKiemTra.ThucThi(app.config.Proxy.FilePath)
	if err != nil {
		app.logger.Errorf("Lỗi kiểm tra proxy: %v", err)
		return err
	}

	// Hiển thị hoặc gửi kết quả
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
				kq.Proxy.GiaoDien, kq.Proxy.DiaChi, kq.Proxy.Cong, kq.DiaChi, trangThai)
		}
	}

	app.logger.Infof("Kiểm tra proxy hoàn tất! Tổng số proxy: %d", len(ketQua))
	return nil
}

// ThucThiXuLyTacVu thực hiện quy trình xử lý tác vụ
func (app *Application) ThucThiXuLyTacVu(kieuGiaoDien string) error {
	boXuLyLoDong, err := app.thietLapXuLyTacVu(kieuGiaoDien)
	if err != nil {
		return err
	}

	// Tạo tệp log cho chế độ TUI
	if kieuGiaoDien != "" {
		tepGhiNhatKy, err := utils.CreateLogFile()
		if err != nil {
			return err
		}
		defer tepGhiNhatKy.Close()
		app.logger.SetOutput(tepGhiNhatKy)

		// Khởi động TUI
		if err := app.tuiUseCase.Start(); err != nil {
			app.logger.Errorf("Không thể khởi động TUI: %v", err)
			return err
		}
		defer app.tuiUseCase.Close()
	}

	// Thực thi usecase xử lý tác vụ
	ketQua, err := boXuLyLoDong.ThucThi(app.config.Input.FilePath)
	if err != nil {
		app.logger.Errorf("Lỗi xử lý các tác vụ: %v", err)
		return err
	}

	// Hiển thị hoặc gửi kết quả
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

// DungUngDung đóng ứng dụng một cách an toàn
func (app *Application) DungUngDung() {
	app.logger.Info("Đang dừng ứng dụng...")
	if app.tuiUseCase != nil {
		app.tuiUseCase.Close()
	}
	close(app.stopChannel)
}

func main() {
	// Phân tích cờ dòng lệnh
	kieuGiaoDien := flag.String("tui", "", "Loại TUI: tview, bubbletea, termui")
	kiemTraProxy := flag.Bool("proxy", false, "Kiểm tra proxy từ tệp proxy.txt")
	xuLyTacVu := flag.Bool("task", false, "Xử lý tác vụ từ tệp tasks.txt")
	loaiKetNoi := flag.String("client", "nethttp", "Loại HTTP client: fasthttp, nethttp")
	flag.Parse()

	// Tạo ứng dụng
	app, err := NewApplication("configs/")
	if err != nil {
		log.Fatalf("Không thể khởi tạo ứng dụng: %v", err)
	}

	app.logger.Info("Ứng dụng WorkerCLI đang khởi động")

	// Thiết lập xử lý tín hiệu để tắt ứng dụng một cách an toàn
	tatDauHieu := make(chan os.Signal, 1)
	signal.Notify(tatDauHieu, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-tatDauHieu
		app.logger.Infof("Nhận tín hiệu: %v", sig)
		app.DungUngDung()
		os.Exit(0)
	}()

	// Thực thi quy trình dựa trên cờ
	if *kiemTraProxy {
		if err := app.ThucThiKiemTraProxy(*loaiKetNoi, *kieuGiaoDien); err != nil {
			app.logger.Errorf("Lỗi khi kiểm tra proxy: %v", err)
		}
	} else if *xuLyTacVu {
		if err := app.ThucThiXuLyTacVu(*kieuGiaoDien); err != nil {
			app.logger.Errorf("Lỗi khi xử lý tác vụ: %v", err)
		}
	} else {
		log.Println("Vui lòng chỉ định -proxy hoặc -task")
	}

	app.DungUngDung()
}
