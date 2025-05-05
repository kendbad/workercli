// Package model chứa các mô hình dữ liệu cốt lõi của ứng dụng
package model

// Proxy đại diện cho một proxy server với thông tin kết nối
// Protocol là loại proxy (http, https, socks5)
// IP là địa chỉ IP của proxy server
// Port là cổng của proxy server
type Proxy struct {
	Protocol string
	IP       string
	Port     string
}

// ProxyResult chứa kết quả kiểm tra một proxy
// Proxy là thông tin proxy được kiểm tra
// IP là địa chỉ IP được trả về khi sử dụng proxy này
// Status là trạng thái kiểm tra (ok, error, timeout)
// Error là thông báo lỗi chi tiết nếu có
type ProxyResult struct {
	Proxy  Proxy
	IP     string
	Status string
	Error  string
}
