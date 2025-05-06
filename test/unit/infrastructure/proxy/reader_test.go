package proxy_test

import (
	"os"
	"testing"
	"workercli/internal/infrastructure/proxy"
	"workercli/pkg/utils"
)

// TestFileReader kiểm tra chức năng đọc tệp proxy
// Đây là một ví dụ về unit test cho tầng Infrastructure
// Trong Clean Architecture, tầng Infrastructure chứa code giao tiếp với thế giới bên ngoài
// (cơ sở dữ liệu, tệp, API, UI...)
func TestFileReader(t *testing.T) {
	// Chuẩn bị dữ liệu kiểm tra
	logger := utils.NewTestLogger()
	tepTamThoi, err := os.CreateTemp("", "proxies.txt")
	if err != nil {
		t.Fatalf("Không thể tạo tệp tạm thời: %v", err)
	}
	defer os.Remove(tepTamThoi.Name())

	// Viết dữ liệu kiểm tra
	tepTamThoi.WriteString("http://192.168.1.1:8080\n")
	tepTamThoi.WriteString("socks5://10.0.0.1:1080\n")
	tepTamThoi.Close()

	// Thực hiện kiểm tra
	boDoc := proxy.NewFileReader(logger)
	danhSachProxy, err := boDoc.ReadProxies(tepTamThoi.Name())

	// Xác nhận kết quả
	if err != nil {
		t.Errorf("Lỗi khi đọc proxy: %v", err)
	}

	if len(danhSachProxy) != 2 {
		t.Errorf("Kỳ vọng 2 proxy, nhận được %d", len(danhSachProxy))
	}

	if danhSachProxy[0].GiaoDien != "http" || danhSachProxy[0].DiaChi != "192.168.1.1" || danhSachProxy[0].Cong != "8080" {
		t.Errorf("Proxy đầu tiên không khớp với mong đợi")
	}

	if danhSachProxy[1].GiaoDien != "socks5" || danhSachProxy[1].DiaChi != "10.0.0.1" || danhSachProxy[1].Cong != "1080" {
		t.Errorf("Proxy thứ hai không khớp với mong đợi")
	}
}
