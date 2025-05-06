package testutil

import (
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

// NhomXuLyMock là mock cho NhomXuLy trong worker
// Theo Clean Architecture, đây là một mock cho tầng adapter, giúp kiểm tra usecase
// mà không phụ thuộc vào triển khai thực tế của infrastructure
type NhomXuLyMock struct {
	KetQuaList  []model.KetQua
	StartCalled bool
	StopCalled  bool
	kenhKetQua  chan model.KetQua
}

// TaoNhomXuLyMock tạo một NhomXuLyMock mới với danh sách kết quả cho trước
func TaoNhomXuLyMock(ketQuaList []model.KetQua) *NhomXuLyMock {
	kenhKetQua := make(chan model.KetQua, len(ketQuaList))

	mock := &NhomXuLyMock{
		KetQuaList: ketQuaList,
		kenhKetQua: kenhKetQua,
	}

	// Chuẩn bị kênh kết quả
	for _, kq := range ketQuaList {
		kenhKetQua <- kq
	}

	return mock
}

// BatDau mock phương thức BatDau
func (m *NhomXuLyMock) BatDau() {
	m.StartCalled = true
}

// NopTacVu mock phương thức NopTacVu
func (m *NhomXuLyMock) NopTacVu(tacVu model.TacVu) {
	// Không cần triển khai gì, vì chúng ta đã chuẩn bị kết quả trước
}

// KetQua mock phương thức KetQua
func (m *NhomXuLyMock) KetQua() <-chan model.KetQua {
	return m.kenhKetQua
}

// Dung mock phương thức Dung
func (m *NhomXuLyMock) Dung() {
	m.StopCalled = true
	close(m.kenhKetQua)
}

// ParseProxyMock là mock cho ParseProxy trong infrastructure/proxy
// Clean Architecture: Tách biệt infrastructure từ business logic
// cho phép ta dễ dàng mock các hàm infrastructure trong test
func ParseProxyMock(chuoiProxy string) (model.Proxy, error) {
	return model.Proxy{
		GiaoDien: "http",
		DiaChi:   "192.168.1.1",
		Cong:     "8080",
	}, nil
}

// NewMockLogger tạo một mock logger cho test
func NewMockLogger() *utils.Logger {
	return utils.NewTestLogger()
}

// BoDocMock là một mock đơn giản cho Reader interface
// Theo Clean Architecture, interface này nằm ở tầng adapter
// giúp tách biệt usecase khỏi triển khai cụ thể của việc đọc dữ liệu
type BoDocMock struct {
	DanhSachProxy []model.Proxy
	LoiXayRa      error
}

// ReadProxies triển khai phương thức ReadProxies của Reader interface
func (m *BoDocMock) ReadProxies(_ string) ([]model.Proxy, error) {
	return m.DanhSachProxy, m.LoiXayRa
}

// BoKiemTraMock là một mock đơn giản cho BoKiemTra interface
// Theo Clean Architecture, interface này nằm ở tầng adapter
// giúp tách biệt usecase khỏi triển khai cụ thể của việc kiểm tra proxy
type BoKiemTraMock struct {
	DiaChi    string
	TrangThai string
	LoiXayRa  error
}

// KiemTraProxy triển khai phương thức KiemTraProxy của BoKiemTra interface
func (m *BoKiemTraMock) KiemTraProxy(_ model.Proxy, _ string) (string, string, error) {
	return m.DiaChi, m.TrangThai, m.LoiXayRa
}

// TaoProxyTest tạo một proxy test với các giá trị mặc định
// Hàm tiện ích để tạo nhanh proxy cho mục đích test
func TaoProxyTest() model.Proxy {
	return model.Proxy{
		GiaoDien: "http",
		DiaChi:   "192.168.1.1",
		Cong:     "8080",
	}
}
