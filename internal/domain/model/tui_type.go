package model

// LoaiHienThi định nghĩa các loại hiển thị khác nhau trong ứng dụng
type LoaiHienThi string

const (
	// LoaiKiemTraProxy chế độ hiển thị cho kiểm tra proxy
	LoaiKiemTraProxy LoaiHienThi = "proxy"

	// LoaiXuLyTacVu chế độ hiển thị cho xử lý tác vụ
	LoaiXuLyTacVu LoaiHienThi = "task"
)
