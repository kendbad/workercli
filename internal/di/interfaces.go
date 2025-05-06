package di

import (
	"workercli/internal/domain/model"
)

// ProxyReaderInterface định nghĩa giao diện cho việc đọc danh sách proxy từ nguồn dữ liệu
type ProxyReaderInterface interface {
	ReadProxies(duongDan string) ([]model.Proxy, error)
}

// ProxyCheckerInterface định nghĩa giao diện cho việc kiểm tra proxy
type ProxyCheckerInterface interface {
	CheckProxy(proxy model.Proxy, duongDanKiemTra string) (diaChi string, trangThai string, err error)
}

// TaskProcessorInterface định nghĩa giao diện cho bộ xử lý tác vụ
type TaskProcessorInterface interface {
	XuLyTacVu(tacVu model.TacVu) (model.KetQua, error)
}

// InputReaderInterface định nghĩa giao diện cho việc đọc dữ liệu đầu vào
type InputReaderInterface interface {
	ReadTasks(duongDan string) ([]model.TacVu, error)
}

// TUIRendererInterface định nghĩa giao diện cho bộ hiển thị giao diện người dùng
type TUIRendererInterface interface {
	Start() error
	AddTaskResult(ketQua model.KetQua)
	AddProxyResult(ketQua model.KetQuaProxy)
	Close()
}
