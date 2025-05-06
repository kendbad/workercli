package usecase_adapter_test

import (
	"os"
	"testing"
	"workercli/internal/adapter/proxy"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// BoKiemTraProxyMock là một mockup cho BoKiemTra interface
// Cung cấp kết quả đã cấu hình trước để kiểm tra
type BoKiemTraProxyMock struct {
	DiaChi    string
	TrangThai string
	LoiXayRa  error
}

// KiemTraProxy triển khai phương thức KiemTraProxy
func (m *BoKiemTraProxyMock) KiemTraProxy(_ model.Proxy, _ string) (string, string, error) {
	return m.DiaChi, m.TrangThai, m.LoiXayRa
}

// TestProxyCheckUsecaseAdapterSimple kiểm tra đơn giản giữa usecase và adapter
// Phiên bản đơn giản hơn và tập trung vào tích hợp giữa adapter và usecase
func TestProxyCheckUsecaseAdapterSimple(t *testing.T) {
	// Chuẩn bị dữ liệu kiểm tra
	logger := utils.NewTestLogger()
	tepTamThoi, err := os.CreateTemp("", "proxies.txt")
	if err != nil {
		t.Fatalf("Không thể tạo tệp tạm thời: %v", err)
	}
	defer os.Remove(tepTamThoi.Name())

	// Viết dữ liệu kiểm tra
	tepTamThoi.WriteString("http://192.168.1.1:8080\n")
	tepTamThoi.WriteString("https://10.0.0.1:443\n")
	tepTamThoi.Close()

	// Thiết lập mock
	mockKiemTra := &BoKiemTraProxyMock{
		DiaChi:    "123.45.67.89",
		TrangThai: "Thành công",
		LoiXayRa:  nil,
	}

	// Tạo adapter với mock
	boDocTep := &BoDocTepMock{DuongDan: tepTamThoi.Name()}
	boDocProxy := &proxy.BoDocProxy{
		BoDocMock: boDocTep,
	}

	boKiemTra := proxy.TaoBoKiemTraProxy(logger, nil)
	boKiemTra.BoKiemTraMock = mockKiemTra

	// Kiểm tra đọc proxy thông qua adapter
	danhSachProxy, err := boDocProxy.ReadProxies(tepTamThoi.Name())
	if err != nil {
		t.Errorf("Không mong đợi lỗi khi đọc proxy: %v", err)
	}

	if len(danhSachProxy) != 2 {
		t.Fatalf("Mong đợi 2 proxy, nhận được %d", len(danhSachProxy))
	}

	// Kiểm tra proxy đầu tiên
	if danhSachProxy[0].GiaoDien != "http" || danhSachProxy[0].DiaChi != "192.168.1.1" || danhSachProxy[0].Cong != "8080" {
		t.Errorf("Proxy đầu tiên không khớp với mong đợi")
	}

	// Kiểm tra kiểm tra proxy thông qua adapter
	diaChi, trangThai, err := boKiemTra.CheckProxy(danhSachProxy[0], "http://ip-api.com/json")
	if err != nil {
		t.Errorf("Không mong đợi lỗi khi kiểm tra proxy: %v", err)
	}

	if diaChi != "123.45.67.89" {
		t.Errorf("Địa chỉ không khớp với mong đợi: %s", diaChi)
	}

	if trangThai != "Thành công" {
		t.Errorf("Trạng thái không khớp với mong đợi: %s", trangThai)
	}
}

// BoDocTepMock là mockup cho Reader interface
type BoDocTepMock struct {
	DuongDan string
}

// ReadProxies triển khai phương thức ReadProxies
func (m *BoDocTepMock) ReadProxies(nguon string) ([]model.Proxy, error) {
	if nguon != m.DuongDan {
		return nil, os.ErrNotExist
	}

	return []model.Proxy{
		{GiaoDien: "http", DiaChi: "192.168.1.1", Cong: "8080"},
		{GiaoDien: "https", DiaChi: "10.0.0.1", Cong: "443"},
	}, nil
}
