// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
// Trong Clean Architecture, package này thuộc tầng domain - tầng trung tâm nhất,
// độc lập với các framework, UI, database và các yếu tố bên ngoài khác.
package model

// KetQua đại diện cho kết quả xử lý của một tác vụ
// MaTacVu là ID của tác vụ đã được xử lý
// TrangThai là trạng thái kết quả (completed, error, processing)
// ChiTiet chứa thông tin chi tiết về kết quả xử lý
// LoiXayRa chứa thông báo lỗi nếu có
// Đây là một domain entity thể hiện kết quả xử lý, là phần cốt lõi của business rules
type KetQua struct {
	MaTacVu   string // MaTacVu: mã tác vụ - định danh của tác vụ gốc
	TrangThai string // TrangThai: trạng thái - kết quả xử lý (thành công, lỗi...)
	ChiTiet   string // ChiTiet: chi tiết - thông tin chi tiết về kết quả
	LoiXayRa  string // LoiXayRa: lỗi xảy ra - mô tả lỗi nếu có
}

// KetQuaTacVu là một phiên bản đơn giản hơn của KetQua
// Được sử dụng trong một số trường hợp cụ thể
// Đây là một ví dụ về việc sử dụng các entity đơn giản hơn khi cần thiết
type KetQuaTacVu struct {
	MaTacVu   string // MaTacVu: mã tác vụ - định danh của tác vụ
	TrangThai string // TrangThai: trạng thái - kết quả xử lý đơn giản
}
