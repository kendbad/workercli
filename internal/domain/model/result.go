// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
package model

// Result đại diện cho kết quả xử lý của một task
// TaskID là ID của task đã được xử lý
// Status là trạng thái kết quả (completed, error, processing)
// Details chứa thông tin chi tiết về kết quả xử lý
// Error chứa thông báo lỗi nếu có
type Result struct {
	TaskID  string
	Status  string
	Details string
	Error   string
}

// TaskResult là một phiên bản đơn giản hơn của Result
// Được sử dụng trong một số trường hợp cụ thể
type TaskResult struct {
	TaskID string
	Status string
}
