// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
package model

// KetQua đại diện cho kết quả xử lý của một tác vụ
// MaTacVu là ID của tác vụ đã được xử lý
// TrangThai là trạng thái kết quả (completed, error, processing)
// ChiTiet chứa thông tin chi tiết về kết quả xử lý
// LoiXayRa chứa thông báo lỗi nếu có
type KetQua struct {
	MaTacVu   string
	TrangThai string
	ChiTiet   string
	LoiXayRa  string
}

// KetQuaTacVu là một phiên bản đơn giản hơn của KetQua
// Được sử dụng trong một số trường hợp cụ thể
type KetQuaTacVu struct {
	MaTacVu   string
	TrangThai string
}
