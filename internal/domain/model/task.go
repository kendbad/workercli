// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
package model

// Task đại diện cho một nhiệm vụ cần được xử lý bởi worker
// ID là định danh duy nhất của task trong hệ thống
// Data là dữ liệu cần xử lý
// TaskID là ID tham chiếu từ hệ thống bên ngoài (nếu có)
type Task struct {
	ID     string
	Data   string
	TaskID string
}
