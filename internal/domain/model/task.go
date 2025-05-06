// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
package model

// TacVu đại diện cho một tác vụ cần được xử lý
// MaTacVu là định danh duy nhất của tác vụ
// DuLieu chứa thông tin cần thiết để xử lý tác vụ
// Meta chứa thông tin bổ sung về tác vụ
type TacVu struct {
	MaTacVu string
	DuLieu  string
	Meta    map[string]string
}
