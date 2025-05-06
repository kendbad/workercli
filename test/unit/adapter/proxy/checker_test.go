package proxy_test

import (
	"testing"
	"workercli/internal/adapter/proxy"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// BoKiemTraProxyTest là một mockup cho BoKiemTra interface
// Đây là một ví dụ về việc sử dụng mock/stub trong unit test để cô lập các thành phần
type BoKiemTraProxyTest struct {
	KetQua    string
	TrangThai string
	LoiXayRa  error
}

// KiemTraProxy triển khai phương thức KiemTraProxy của interface BoKiemTra
// Trả về kết quả đã được cấu hình trước cho mục đích kiểm tra
func (m *BoKiemTraProxyTest) KiemTraProxy(_ model.Proxy, _ string) (string, string, error) {
	return m.KetQua, m.TrangThai, m.LoiXayRa
}

// TestBoKiemTraProxy kiểm tra adapter BoKiemTraProxy
// Trong Clean Architecture, tầng Adapter chịu trách nhiệm chuyển đổi dữ liệu
// giữa tầng Use Cases và thế giới bên ngoài
func TestBoKiemTraProxy(t *testing.T) {
	// Chuẩn bị dữ liệu kiểm tra
	logger := utils.NewTestLogger()
	mockKiemTra := &BoKiemTraProxyTest{
		KetQua:    "123.45.67.89",
		TrangThai: "Thành công",
		LoiXayRa:  nil,
	}

	// Thực hiện kiểm tra
	boKiemTra := proxy.TaoBoKiemTraProxy(logger, mockKiemTra)
	diaChi, trangThai, err := boKiemTra.CheckProxy(model.Proxy{
		GiaoDien: "http",
		DiaChi:   "192.168.1.1",
		Cong:     "8080",
	}, "http://ip-api.com/json")

	// Xác nhận kết quả
	if err != nil {
		t.Errorf("Kỳ vọng không có lỗi, nhận được: %v", err)
	}

	if diaChi != "123.45.67.89" {
		t.Errorf("Địa chỉ không khớp với mong đợi: %s", diaChi)
	}

	if trangThai != "Thành công" {
		t.Errorf("Trạng thái không khớp với mong đợi: %s", trangThai)
	}
}
