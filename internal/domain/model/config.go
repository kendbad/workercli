// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
package model

// CauHinh mô tả một cấu hình trong hệ thống
// Ten là tên cấu hình để nhận dạng
// GiaTri chứa giá trị của cấu hình
type CauHinh struct {
	Ten    string
	GiaTri string
}

// Config holds business-level configuration for TUI display
type Config struct {
	MaxTasks     int // Maximum number of concurrent tasks
	ProxyTimeout int // Timeout for proxy checks (in seconds)
}
