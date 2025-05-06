// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
// Trong Clean Architecture, model là phần trung tâm nhất của hệ thống,
// chứa các business entity và quy tắc nghiệp vụ cốt lõi.
package model

// TacVu đại diện cho một tác vụ cần được xử lý
// MaTacVu là định danh duy nhất của tác vụ
// DuLieu chứa thông tin cần thiết để xử lý tác vụ
// Meta chứa thông tin bổ sung về tác vụ
// Đây là một entity trong Clean Architecture - thể hiện khái niệm nghiệp vụ cốt lõi
type TacVu struct {
	MaTacVu string            // MaTacVu: mã tác vụ - định danh duy nhất
	DuLieu  string            // DuLieu: dữ liệu - nội dung cần xử lý
	Meta    map[string]string // Meta: thông tin phụ - các thông tin bổ sung
}
