package usecase_test

import (
	"errors"
	"testing"
	"workercli/internal/adapter/proxy"
	"workercli/internal/domain/model"
	"workercli/test/testutil"
)

// TestMockNhomXuLy kiểm tra các thành phần riêng lẻ
func TestMockNhomXuLy(t *testing.T) {
	// Tạo mock và kiểm tra các phương thức cơ bản
	mockNhomXuLy := testutil.TaoNhomXuLyMock([]model.KetQua{
		{MaTacVu: "test", TrangThai: "Thành công"},
	})

	mockNhomXuLy.BatDau()
	if !mockNhomXuLy.StartCalled {
		t.Error("Phương thức BatDau không gọi được")
	}

	mockNhomXuLy.NopTacVu(model.TacVu{MaTacVu: "test"})

	// Đọc kết quả từ kênh
	ketQua := <-mockNhomXuLy.KetQua()
	if ketQua.TrangThai != "Thành công" {
		t.Errorf("Trạng thái kết quả không đúng, mong đợi 'Thành công', nhận được: %s", ketQua.TrangThai)
	}

	mockNhomXuLy.Dung()
	if !mockNhomXuLy.StopCalled {
		t.Error("Phương thức Dung không gọi được")
	}
}

// TestKiemTraProxyUsecase kiểm tra usecase KiemTraProxy với mock worker
// Clean Architecture: Test với các mock objects cho phép kiểm tra usecase
// mà không phụ thuộc vào infrastructure thực tế
func TestKiemTraProxyVoiMockRiengLe(t *testing.T) {
	// Chuẩn bị dữ liệu kiểm tra
	logger := testutil.NewMockLogger()

	// Trường hợp thành công
	t.Run("Thành công", func(t *testing.T) {
		// Chuẩn bị mock data
		mockBoDoc := &testutil.BoDocMock{
			DanhSachProxy: []model.Proxy{
				testutil.TaoProxyTest(),
			},
			LoiXayRa: nil,
		}

		mockBoKiemTra := &testutil.BoKiemTraMock{
			DiaChi:    "123.45.67.89",
			TrangThai: "Thành công",
			LoiXayRa:  nil,
		}

		// Tạo adapters với mock objects
		boDocProxy := &proxy.BoDocProxy{
			BoDocMock: mockBoDoc,
		}

		boKiemTraProxy := proxy.TaoBoKiemTraProxy(logger, nil)
		boKiemTraProxy.BoKiemTraMock = mockBoKiemTra

		// Đọc proxy và kiểm tra kết quả
		danhSachProxy, err := boDocProxy.ReadProxies("proxy.txt")
		if err != nil {
			t.Errorf("Không mong đợi lỗi khi đọc proxy, nhận được: %v", err)
		}

		if len(danhSachProxy) != 1 {
			t.Errorf("Mong đợi 1 proxy, nhận được %d", len(danhSachProxy))
		}

		// Kiểm tra proxy và kiểm tra kết quả
		diaChi, trangThai, err := boKiemTraProxy.CheckProxy(testutil.TaoProxyTest(), "http://ip-api.com/json")
		if err != nil {
			t.Errorf("Không mong đợi lỗi khi kiểm tra proxy, nhận được: %v", err)
		}

		if diaChi != "123.45.67.89" {
			t.Errorf("Địa chỉ không khớp với mong đợi: %s", diaChi)
		}

		if trangThai != "Thành công" {
			t.Errorf("Trạng thái không khớp với mong đợi: %s", trangThai)
		}
	})

	// Trường hợp thất bại khi kiểm tra proxy
	t.Run("Lỗi kiểm tra proxy", func(t *testing.T) {
		mockBoKiemTra := &testutil.BoKiemTraMock{
			DiaChi:    "",
			TrangThai: "Thất bại",
			LoiXayRa:  errors.New("không thể kết nối"),
		}

		// Tạo adapter với mock object
		boKiemTraProxy := proxy.TaoBoKiemTraProxy(logger, nil)
		boKiemTraProxy.BoKiemTraMock = mockBoKiemTra

		// Kiểm tra lỗi khi kiểm tra proxy
		_, trangThai, err := boKiemTraProxy.CheckProxy(testutil.TaoProxyTest(), "http://ip-api.com/json")

		if err == nil {
			t.Error("Mong đợi lỗi khi kiểm tra proxy, nhưng không nhận được lỗi")
		}

		if trangThai != "Thất bại" {
			t.Errorf("Mong đợi trạng thái 'Thất bại', nhận được: %s", trangThai)
		}
	})
}

// TestKiemTraProxySimple kiểm tra cơ bản usecase KiemTraProxy
// Clean Architecture: Kiểm tra từng adapter độc lập trước,
// sau đó kết hợp chúng lại trong usecase
func TestKiemTraProxySimple(t *testing.T) {
	// Chuẩn bị dữ liệu kiểm tra
	logger := testutil.NewMockLogger()

	// Trường hợp thành công
	t.Run("Kiểm tra từng thành phần", func(t *testing.T) {
		mockBoDoc := &testutil.BoDocMock{
			DanhSachProxy: []model.Proxy{
				testutil.TaoProxyTest(),
			},
			LoiXayRa: nil,
		}

		mockBoKiemTra := &testutil.BoKiemTraMock{
			DiaChi:    "123.45.67.89",
			TrangThai: "Thành công",
			LoiXayRa:  nil,
		}

		// Tạo adapters với mock objects
		boDocProxy := &proxy.BoDocProxy{
			BoDocMock: mockBoDoc,
		}

		boKiemTraProxy := proxy.TaoBoKiemTraProxy(logger, nil)
		boKiemTraProxy.BoKiemTraMock = mockBoKiemTra

		// Cách đơn giản để kiểm tra đọc proxy
		danhSachProxy, err := boDocProxy.ReadProxies("proxy.txt")
		if err != nil {
			t.Errorf("Không mong đợi lỗi khi đọc proxy, nhận được: %v", err)
		}

		if len(danhSachProxy) != 1 {
			t.Errorf("Mong đợi 1 proxy, nhận được %d", len(danhSachProxy))
		}

		// Cách đơn giản để kiểm tra kiểm tra proxy
		diaChi, trangThai, err := boKiemTraProxy.CheckProxy(testutil.TaoProxyTest(), "http://ip-api.com/json")

		if err != nil {
			t.Errorf("Không mong đợi lỗi khi kiểm tra proxy, nhận được: %v", err)
		}

		if diaChi != "123.45.67.89" {
			t.Errorf("Địa chỉ không khớp với mong đợi: %s", diaChi)
		}

		if trangThai != "Thành công" {
			t.Errorf("Trạng thái không khớp với mong đợi: %s", trangThai)
		}
	})

	// Trường hợp thất bại khi kiểm tra proxy
	t.Run("Lỗi kiểm tra proxy", func(t *testing.T) {
		mockBoKiemTra := &testutil.BoKiemTraMock{
			DiaChi:    "",
			TrangThai: "Thất bại",
			LoiXayRa:  errors.New("không thể kết nối"),
		}

		// Tạo adapter với mock object
		boKiemTraProxy := proxy.TaoBoKiemTraProxy(logger, nil)
		boKiemTraProxy.BoKiemTraMock = mockBoKiemTra

		// Kiểm tra lỗi khi kiểm tra proxy
		_, trangThai, err := boKiemTraProxy.CheckProxy(testutil.TaoProxyTest(), "http://ip-api.com/json")

		if err == nil {
			t.Error("Mong đợi lỗi khi kiểm tra proxy, nhưng không nhận được lỗi")
		}

		if trangThai != "Thất bại" {
			t.Errorf("Mong đợi trạng thái 'Thất bại', nhận được: %s", trangThai)
		}
	})
}
