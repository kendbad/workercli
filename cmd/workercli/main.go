package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	// Tầng Adapter: Kết nối các thành phần bên ngoài với các use case của ứng dụng
	// Chịu trách nhiệm chuyển đổi dữ liệu qua lại giữa các định dạng bên ngoài và các cấu trúc nội bộ

	// Tầng Cấu hình & Miền: Chứa các quy tắc nghiệp vụ và định nghĩa các thực thể của ứng dụng
	// Đây là phần cốt lõi độc lập nhất của ứng dụng, không phụ thuộc vào các thành phần bên ngoài

	"workercli/internal/di"

	// Tầng Cơ sở hạ tầng: Triển khai giao diện và cung cấp các dịch vụ cụ thể
	// Đây là tầng thực hiện việc giao tiếp với thế giới bên ngoài (file, database, API...)

	// Tầng Trường hợp sử dụng: Điều phối luồng dữ liệu và xử lý logic ứng dụng
	// Chứa các logic nghiệp vụ cụ thể, sử dụng các giao diện được định nghĩa trong domain

	// Tiện ích: Các công cụ hỗ trợ chung cho ứng dụng
	"workercli/pkg/utils"
)

// UngDung đại diện cho ứng dụng đã được cấu hình đầy đủ
// Application trong Clean Architecture là điểm khởi đầu, nơi phối hợp tất cả các thành phần
type UngDung struct {
	diContainer *di.Container // Container DI để quản lý các dependency
}

// TaoUngDung tạo và cấu hình ứng dụng
// NewApplication là điểm khởi tạo các thành phần cần thiết theo nguyên tắc Dependency Injection
func TaoUngDung(thuMucCauHinh string) (*UngDung, error) {
	// Tạo container DI
	container := di.NewContainer()
	if err := container.KhoiTao(thuMucCauHinh); err != nil {
		return nil, err
	}

	ungDung := &UngDung{
		diContainer: container,
	}

	return ungDung, nil
}

// ThucThiKiemTraProxy thực hiện quy trình kiểm tra proxy
// Đây là một điểm vào của ứng dụng, nơi phối hợp luồng công việc của các use case
func (app *UngDung) ThucThiKiemTraProxy(loaiKetNoi, kieuGiaoDien string) error {
	boGhiNhatKy := app.diContainer.LayBoGhiNhatKy()

	// Thiết lập kiểm tra proxy từ container
	boKiemTra, err := app.diContainer.ThietLapKiemTraProxy(loaiKetNoi, kieuGiaoDien)
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
		boGhiNhatKy.SetOutput(tepGhiNhatKy)

		// Khởi động TUI
		boGiaoDien := app.diContainer.LayBoGiaoDien()
		if boGiaoDien != nil {
			if err := boGiaoDien.Start(); err != nil {
				boGhiNhatKy.Errorf("Không thể khởi động TUI: %v", err)
				return err
			}
			defer boGiaoDien.Close()
		}
	}

	// Thực thi usecase kiểm tra proxy - Gọi use case là cách để kích hoạt một luồng nghiệp vụ
	ketQua, err := boKiemTra.ThucThi(app.diContainer.LayCauHinh().Proxy.FilePath) // ketQua: kết quả
	if err != nil {
		boGhiNhatKy.Errorf("Lỗi kiểm tra proxy: %v", err)
		return err
	}

	// Hiển thị hoặc gửi kết quả - Đầu ra có thể thông qua giao diện hoặc log
	if kieuGiaoDien != "" {
		boGiaoDien := app.diContainer.LayBoGiaoDien()
		if boGiaoDien != nil {
			for _, kq := range ketQua {
				boGhiNhatKy.Infof("Gửi kết quả proxy vào TUI: %v", kq)
				boGiaoDien.AddProxyResult(kq)
			}
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

	boGhiNhatKy.Infof("Kiểm tra proxy hoàn tất! Tổng số proxy: %d", len(ketQua))
	return nil
}

// ThucThiXuLyTacVu thực hiện quy trình xử lý tác vụ
// Phối hợp luồng xử lý tác vụ, đóng vai trò điều phối workflow
func (app *UngDung) ThucThiXuLyTacVu(kieuGiaoDien string) error {
	boGhiNhatKy := app.diContainer.LayBoGhiNhatKy()

	// Thiết lập xử lý tác vụ từ container
	boXuLyHangLoat, err := app.diContainer.ThietLapXuLyTacVu(kieuGiaoDien)
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
		boGhiNhatKy.SetOutput(tepGhiNhatKy)

		// Khởi động TUI
		boGiaoDien := app.diContainer.LayBoGiaoDien()
		if boGiaoDien != nil {
			if err := boGiaoDien.Start(); err != nil {
				boGhiNhatKy.Errorf("Không thể khởi động TUI: %v", err)
				return err
			}
			defer boGiaoDien.Close()
		}
	}

	// Thực thi usecase xử lý tác vụ - Use case định nghĩa các luồng nghiệp vụ chính
	ketQua, err := boXuLyHangLoat.ThucThi(app.diContainer.LayCauHinh().Input.FilePath)
	if err != nil {
		boGhiNhatKy.Errorf("Lỗi xử lý các tác vụ: %v", err)
		return err
	}

	// Hiển thị hoặc gửi kết quả
	if kieuGiaoDien != "" {
		boGiaoDien := app.diContainer.LayBoGiaoDien()
		if boGiaoDien != nil {
			for _, kq := range ketQua {
				boGhiNhatKy.Infof("Gửi kết quả tác vụ vào TUI: %v", kq)
				boGiaoDien.AddTaskResult(kq)
			}
		}
	} else {
		for _, kq := range ketQua {
			log.Printf("Tác vụ %s: %s\n", kq.MaTacVu, kq.TrangThai)
		}
	}

	boGhiNhatKy.Infof("Xử lý tác vụ hoàn tất! Tổng số tác vụ: %d", len(ketQua))
	return nil
}

// DungUngDung đóng ứng dụng một cách an toàn
// Xử lý việc dọn dẹp tài nguyên khi ứng dụng kết thúc
func (app *UngDung) DungUngDung() {
	boGhiNhatKy := app.diContainer.LayBoGhiNhatKy()
	boGhiNhatKy.Info("Đang dừng ứng dụng...")
	app.diContainer.DungContainer()
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

	ungDung.diContainer.LayBoGhiNhatKy().Info("Ứng dụng WorkerCLI đang khởi động")

	// Thiết lập xử lý tín hiệu để tắt ứng dụng một cách an toàn
	// Xử lý sự kiện hệ thống là một phần của tầng Infrastructure trong Clean Architecture
	tatDauHieu := make(chan os.Signal, 1) // tatDauHieu: tắt dấu hiệu
	signal.Notify(tatDauHieu, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		dau := <-tatDauHieu // dau: dấu
		ungDung.diContainer.LayBoGhiNhatKy().Infof("Nhận tín hiệu: %v", dau)
		ungDung.DungUngDung()
		os.Exit(0)
	}()

	// Thực thi quy trình dựa trên cờ - Luồng điều khiển chính của ứng dụng
	// Clean Architecture cho phép các use case được gọi từ nhiều nguồn khác nhau
	if *kiemTraProxy {
		if err := ungDung.ThucThiKiemTraProxy(*loaiKetNoi, *kieuGiaoDien); err != nil {
			ungDung.diContainer.LayBoGhiNhatKy().Errorf("Lỗi khi kiểm tra proxy: %v", err)
		}
	} else if *xuLyTacVu {
		if err := ungDung.ThucThiXuLyTacVu(*kieuGiaoDien); err != nil {
			ungDung.diContainer.LayBoGhiNhatKy().Errorf("Lỗi khi xử lý tác vụ: %v", err)
		}
	} else {
		log.Println("Vui lòng chỉ định -proxy hoặc -task")
	}

	ungDung.DungUngDung()
}
