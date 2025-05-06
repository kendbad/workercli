// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
// Trong Clean Architecture, tầng Domain là tầng trung tâm nhất,
// chứa các entity và business rules độc lập với bất kỳ framework hay thư viện nào.
package model

// Proxy đại diện cho một proxy server với thông tin kết nối
// GiaoDien là loại proxy (http, https, socks5)
// DiaChi là địa chỉ IP của proxy server
// Cong là cổng của proxy server
// Đây là một entity trong Clean Architecture, đại diện cho một khái niệm nghiệp vụ cốt lõi
type Proxy struct {
	GiaoDien string // GiaoDien: giao diện - loại giao thức proxy (http, https, socks5)
	DiaChi   string // DiaChi: địa chỉ - địa chỉ IP hoặc hostname của proxy
	Cong     string // Cong: cổng - cổng kết nối của proxy
}

// KetQuaProxy chứa kết quả kiểm tra một proxy
// Proxy là thông tin proxy được kiểm tra
// DiaChi là địa chỉ IP được trả về khi sử dụng proxy này
// TrangThai là trạng thái kiểm tra (ok, error, timeout)
// LoiXayRa là thông báo lỗi chi tiết nếu có
// Đây cũng là một entity nhưng phức tạp hơn, kết hợp các entity khác
type KetQuaProxy struct {
	Proxy     Proxy  // Thông tin về proxy được kiểm tra
	DiaChi    string // DiaChi: địa chỉ - địa chỉ IP thực tế khi sử dụng proxy này
	TrangThai string // TrangThai: trạng thái - kết quả kiểm tra (thành công, thất bại...)
	LoiXayRa  string // LoiXayRa: lỗi xảy ra - thông báo lỗi chi tiết nếu có
}
