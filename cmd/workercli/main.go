package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// Tầng Adapter: Kết nối các thành phần bên ngoài với các use case của ứng dụng
	// Chịu trách nhiệm chuyển đổi dữ liệu qua lại giữa các định dạng bên ngoài và các cấu trúc nội bộ
	"workercli/internal/adapter/input"
	proxyiface "workercli/internal/adapter/proxy"

	// Tầng Cấu hình & Miền: Chứa các quy tắc nghiệp vụ và định nghĩa các thực thể của ứng dụng
	// Đây là phần cốt lõi độc lập nhất của ứng dụng, không phụ thuộc vào các thành phần bên ngoài
	"workercli/internal/config"
	"workercli/internal/domain/model"

	// Tầng Cơ sở hạ tầng: Triển khai giao diện và cung cấp các dịch vụ cụ thể
	// Đây là tầng thực hiện việc giao tiếp với thế giới bên ngoài (file, database, API...)
	"workercli/internal/infrastructure/proxy"
	"workercli/internal/infrastructure/task"
	tuiinfra "workercli/internal/infrastructure/tui"

	// Tầng Trường hợp sử dụng: Điều phối luồng dữ liệu và xử lý logic ứng dụng
	// Chứa các logic nghiệp vụ cụ thể, sử dụng các giao diện được định nghĩa trong domain
	"workercli/internal/usecase"

	// Tiện ích: Các công cụ hỗ trợ chung cho ứng dụng
	"workercli/pkg/utils"
)

// UngDung đại diện cho ứng dụng đã được cấu hình đầy đủ
// Application trong Clean Architecture là điểm khởi đầu, nơi phối hợp tất cả các thành phần
type UngDung struct {
	cauHinh     *config.Config       // cauHinh: cấu hình
	boGhiNhatKy *utils.Logger        // boGhiNhatKy: bộ ghi nhật ký
	boGiaoDien  *tuiinfra.TUIUseCase // boGiaoDien: bộ giao diện người dùng
	kenhDungLai chan struct{}        // kenhDungLai: kênh dừng lại
}

// TaoUngDung tạo và cấu hình ứng dụng
// NewApplication là điểm khởi tạo các thành phần cần thiết theo nguyên tắc Dependency Injection
func TaoUngDung(thuMucCauHinh string) (*UngDung, error) {
	// 1. Tải cấu hình
	cauHinh, err := config.Load(thuMucCauHinh)
	if err != nil {
		return nil, err
	}

	// 2. Khởi tạo bộ ghi nhật ký
	boGhiNhatKy, err := utils.NewLogger("configs/logger.yaml")
	if err != nil {
		return nil, err
	}

	ungDung := &UngDung{
		cauHinh:     cauHinh,
		boGhiNhatKy: boGhiNhatKy,
		kenhDungLai: make(chan struct{}),
	}

	return ungDung, nil
}

// thietLapKiemTraProxy tạo và cấu hình tất cả các thành phần cần thiết cho việc kiểm tra proxy
// Đây là ví dụ về việc áp dụng Dependency Injection trong Clean Architecture
func (app *UngDung) thietLapKiemTraProxy(loaiKetNoi, kieuGiaoDien string) (*usecase.KiemTraProxy, error) {
	// Thiết lập cơ sở hạ tầng proxy - Infrastructure layer triển khai các interface
	boDocTep := proxy.NewFileReader(app.boGhiNhatKy)                   // boDocTep: bộ đọc tệp
	boDocProxy := proxyiface.NewProxyReader(app.boGhiNhatKy, boDocTep) // boDocProxy: bộ đọc proxy

	// Tạo bộ kiểm tra IP qua proxy triển khai giao diện Checker
	// Adapter pattern: liên kết các triển khai cụ thể với interface trong domain
	boKiemTraIP := proxy.NewIPChecker(app.boGhiNhatKy, loaiKetNoi)               // boKiemTraIP: bộ kiểm tra IP
	boKiemTraProxy := proxyiface.TaoBoKiemTraProxy(app.boGhiNhatKy, boKiemTraIP) // boKiemTraProxy: bộ kiểm tra proxy

	// Tạo usecase kiểm tra proxy - Use Case layer chứa logic nghiệp vụ
	// Usecase là trung tâm điều phối các hoạt động, sử dụng các interface từ domain layer
	boKiemTra := usecase.TaoBoKiemTraProxy( // boKiemTra: bộ kiểm tra
		boDocProxy,
		boKiemTraProxy,
		app.cauHinh.Proxy.CheckURL,
		app.cauHinh.Worker.Workers,
		app.boGhiNhatKy,
	)

	// Thiết lập TUI nếu cần - UI layer là tầng ngoài cùng theo Clean Architecture
	if kieuGiaoDien != "" {
		nhaXuongGiaoDien := tuiinfra.NewRendererFactory(app.boGhiNhatKy, kieuGiaoDien) // nhaXuongGiaoDien: nhà xưởng giao diện
		boHienThi := nhaXuongGiaoDien.CreateProxyRenderer(                             // boHienThi: bộ hiển thị
			app.boGhiNhatKy,
			&[]model.KetQuaProxy{},
			&sync.Mutex{},
			make(chan model.KetQuaProxy, 100),
			app.kenhDungLai,
		)
		app.boGiaoDien = tuiinfra.NewTUIUseCase(app.boGhiNhatKy, kieuGiaoDien, boHienThi)
	}

	return boKiemTra, nil
}

// thietLapXuLyTacVu tạo và cấu hình tất cả các thành phần cần thiết cho việc xử lý tác vụ
// Tuân thủ nguyên tắc Dependency Inversion của Clean Architecture
func (app *UngDung) thietLapXuLyTacVu(kieuGiaoDien string) (*usecase.XuLyHangLoatTacVu, error) {
	// Cơ sở hạ tầng tác vụ - Cung cấp triển khai cụ thể cho các interface
	boXuLy := task.NewProcessor(app.boGhiNhatKy)        // boXuLy: bộ xử lý
	boDocDauVao := input.NewFileReader(app.boGhiNhatKy) // boDocDauVao: bộ đọc đầu vào

	// Tạo usecase xử lý tác vụ hàng loạt - Lớp usecase không trực tiếp phụ thuộc vào triển khai cụ thể
	// mà chỉ phụ thuộc vào các interface (Dependency Inversion Principle)
	boXuLyHangLoat := usecase.TaoBoXuLyHangLoatTacVu( // boXuLyHangLoat: bộ xử lý hàng loạt
		boDocDauVao,
		boXuLy,
		app.cauHinh.Worker.Workers,
		app.boGhiNhatKy,
	)

	// Thiết lập TUI nếu cần - Giao diện người dùng là tầng ngoài cùng trong Clean Architecture
	if kieuGiaoDien != "" {
		nhaXuongGiaoDien := tuiinfra.NewRendererFactory(app.boGhiNhatKy, kieuGiaoDien)
		boHienThi := nhaXuongGiaoDien.CreateTaskRenderer(
			app.boGhiNhatKy,
			&[]model.KetQua{},
			&sync.Mutex{},
			make(chan model.KetQua, 100),
			app.kenhDungLai,
		)
		app.boGiaoDien = tuiinfra.NewTUIUseCase(app.boGhiNhatKy, kieuGiaoDien, boHienThi)
	}

	return boXuLyHangLoat, nil
}

// ThucThiKiemTraProxy thực hiện quy trình kiểm tra proxy
// Đây là một điểm vào của ứng dụng, nơi phối hợp luồng công việc của các use case
func (app *UngDung) ThucThiKiemTraProxy(loaiKetNoi, kieuGiaoDien string) error {
	boKiemTra, err := app.thietLapKiemTraProxy(loaiKetNoi, kieuGiaoDien)
	if err != nil {
		return err
	}

	// Tạo tệp log cho chế độ TUI
	if kieuGiaoDien != "" {
		tepGhiNhatKy, err := utils.CreateLogFile() // tepGhiNhatKy: tệp ghi nhật ký
		if err != nil {
			return err
		}
		defer tepGhiNhatKy.Close()
		app.boGhiNhatKy.SetOutput(tepGhiNhatKy)

		// Khởi động TUI
		if err := app.boGiaoDien.Start(); err != nil {
			app.boGhiNhatKy.Errorf("Không thể khởi động TUI: %v", err)
			return err
		}
		defer app.boGiaoDien.Close()
	}

	// Thực thi usecase kiểm tra proxy - Gọi use case là cách để kích hoạt một luồng nghiệp vụ
	ketQua, err := boKiemTra.ThucThi(app.cauHinh.Proxy.FilePath) // ketQua: kết quả
	if err != nil {
		app.boGhiNhatKy.Errorf("Lỗi kiểm tra proxy: %v", err)
		return err
	}

	// Hiển thị hoặc gửi kết quả - Đầu ra có thể thông qua giao diện hoặc log
	if kieuGiaoDien != "" {
		for _, kq := range ketQua {
			app.boGhiNhatKy.Infof("Gửi kết quả proxy vào TUI: %v", kq)
			app.boGiaoDien.AddProxyResult(kq)
		}
	} else {
		for _, kq := range ketQua {
			trangThai := kq.TrangThai // trangThai: trạng thái
			if kq.LoiXayRa != "" {
				trangThai += " (" + kq.LoiXayRa + ")"
			}
			log.Printf("Proxy %s://%s:%s, IP: %s, Trạng thái: %s\n",
				kq.Proxy.GiaoDien, kq.Proxy.DiaChi, kq.Proxy.Cong, kq.DiaChi, trangThai)
		}
	}

	app.boGhiNhatKy.Infof("Kiểm tra proxy hoàn tất! Tổng số proxy: %d", len(ketQua))
	return nil
}

// ThucThiXuLyTacVu thực hiện quy trình xử lý tác vụ
// Phối hợp luồng xử lý tác vụ, đóng vai trò điều phối workflow
func (app *UngDung) ThucThiXuLyTacVu(kieuGiaoDien string) error {
	boXuLyHangLoat, err := app.thietLapXuLyTacVu(kieuGiaoDien)
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
		app.boGhiNhatKy.SetOutput(tepGhiNhatKy)

		// Khởi động TUI
		if err := app.boGiaoDien.Start(); err != nil {
			app.boGhiNhatKy.Errorf("Không thể khởi động TUI: %v", err)
			return err
		}
		defer app.boGiaoDien.Close()
	}

	// Thực thi usecase xử lý tác vụ - Use case định nghĩa các luồng nghiệp vụ chính
	ketQua, err := boXuLyHangLoat.ThucThi(app.cauHinh.Input.FilePath)
	if err != nil {
		app.boGhiNhatKy.Errorf("Lỗi xử lý các tác vụ: %v", err)
		return err
	}

	// Hiển thị hoặc gửi kết quả
	if kieuGiaoDien != "" {
		for _, kq := range ketQua {
			app.boGhiNhatKy.Infof("Gửi kết quả tác vụ vào TUI: %v", kq)
			app.boGiaoDien.AddTaskResult(kq)
		}
	} else {
		for _, kq := range ketQua {
			log.Printf("Tác vụ %s: %s\n", kq.MaTacVu, kq.TrangThai)
		}
	}

	app.boGhiNhatKy.Infof("Xử lý tác vụ hoàn tất! Tổng số tác vụ: %d", len(ketQua))
	return nil
}

// DungUngDung đóng ứng dụng một cách an toàn
// Xử lý việc dọn dẹp tài nguyên khi ứng dụng kết thúc
func (app *UngDung) DungUngDung() {
	app.boGhiNhatKy.Info("Đang dừng ứng dụng...")
	if app.boGiaoDien != nil {
		app.boGiaoDien.Close()
	}
	close(app.kenhDungLai)
}

func main() {
	// Phân tích cờ dòng lệnh - Command Line Interfaces là một trong những cổng vào (ports)
	// trong Clean Architecture, nằm ở tầng giao diện (Interface Adapters)
	kieuGiaoDien := flag.String("tui", "", "Loại TUI: tview, bubbletea, termui")          // kieuGiaoDien: kiểu giao diện
	kiemTraProxy := flag.Bool("proxy", false, "Kiểm tra proxy từ tệp proxy.txt")          // kiemTraProxy: kiểm tra proxy
	xuLyTacVu := flag.Bool("task", false, "Xử lý tác vụ từ tệp tasks.txt")                // xuLyTacVu: xử lý tác vụ
	loaiKetNoi := flag.String("client", "nethttp", "Loại HTTP client: fasthttp, nethttp") // loaiKetNoi: loại kết nối
	flag.Parse()

	// Tạo ứng dụng - Điểm khởi đầu của Clean Architecture, nơi tất cả các thành phần được kết nối
	ungDung, err := TaoUngDung("configs/")
	if err != nil {
		log.Fatalf("Không thể khởi tạo ứng dụng: %v", err)
	}

	ungDung.boGhiNhatKy.Info("Ứng dụng WorkerCLI đang khởi động")

	// Thiết lập xử lý tín hiệu để tắt ứng dụng một cách an toàn
	// Xử lý sự kiện hệ thống là một phần của tầng Infrastructure trong Clean Architecture
	tatDauHieu := make(chan os.Signal, 1) // tatDauHieu: tắt dấu hiệu
	signal.Notify(tatDauHieu, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		dau := <-tatDauHieu // dau: dấu
		ungDung.boGhiNhatKy.Infof("Nhận tín hiệu: %v", dau)
		ungDung.DungUngDung()
		os.Exit(0)
	}()

	// Thực thi quy trình dựa trên cờ - Luồng điều khiển chính của ứng dụng
	// Clean Architecture cho phép các use case được gọi từ nhiều nguồn khác nhau
	if *kiemTraProxy {
		if err := ungDung.ThucThiKiemTraProxy(*loaiKetNoi, *kieuGiaoDien); err != nil {
			ungDung.boGhiNhatKy.Errorf("Lỗi khi kiểm tra proxy: %v", err)
		}
	} else if *xuLyTacVu {
		if err := ungDung.ThucThiXuLyTacVu(*kieuGiaoDien); err != nil {
			ungDung.boGhiNhatKy.Errorf("Lỗi khi xử lý tác vụ: %v", err)
		}
	} else {
		log.Println("Vui lòng chỉ định -proxy hoặc -task")
	}

	ungDung.DungUngDung()
}
