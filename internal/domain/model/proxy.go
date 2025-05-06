// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
package model

// Proxy đại diện cho một proxy server với thông tin kết nối
// GiaoDien là loại proxy (http, https, socks5)
// DiaChi là địa chỉ IP của proxy server
// Cong là cổng của proxy server
type Proxy struct {
	GiaoDien string
	DiaChi   string
	Cong     string
}

// KetQuaProxy chứa kết quả kiểm tra một proxy
// Proxy là thông tin proxy được kiểm tra
// DiaChi là địa chỉ IP được trả về khi sử dụng proxy này
// TrangThai là trạng thái kiểm tra (ok, error, timeout)
// LoiXayRa là thông báo lỗi chi tiết nếu có
type KetQuaProxy struct {
	Proxy     Proxy
	DiaChi    string
	TrangThai string
	LoiXayRa  string
}
