package worker

import (
	"workercli/internal/domain/model"
	"workercli/internal/domain/service"
	"workercli/pkg/utils"
)

// NguoiXuLy đại diện cho một bộ xử lý riêng lẻ trong nhóm
type NguoiXuLy struct {
	id             int
	danhSachTacVu  <-chan model.TacVu
	danhSachKetQua chan<- model.KetQua
	boXuLy         service.BoXuLyTacVu
	boGhiNhatKy    *utils.Logger
}

// TaoNguoiXuLy tạo một bộ xử lý mới
func TaoNguoiXuLy(id int, danhSachTacVu <-chan model.TacVu, danhSachKetQua chan<- model.KetQua, boXuLy service.BoXuLyTacVu, boGhiNhatKy *utils.Logger) *NguoiXuLy {
	return &NguoiXuLy{
		id:             id,
		danhSachTacVu:  danhSachTacVu,
		danhSachKetQua: danhSachKetQua,
		boXuLy:         boXuLy,
		boGhiNhatKy:    boGhiNhatKy,
	}
}

// Chay bắt đầu vòng lặp xử lý tác vụ của bộ xử lý
func (w *NguoiXuLy) Chay(kenhDung <-chan struct{}) {
	w.boGhiNhatKy.Debugf("Bộ xử lý %d khởi động", w.id)
	for {
		select {
		case tacVu, ok := <-w.danhSachTacVu:
			if !ok {
				w.boGhiNhatKy.Debugf("Bộ xử lý %d dừng do kênh tác vụ đóng", w.id)
				return
			}
			w.boGhiNhatKy.Debugf("Bộ xử lý %d nhận tác vụ %s", w.id, tacVu.MaTacVu)
			ketQua, err := w.boXuLy.XuLyTacVu(tacVu)
			if err != nil {
				w.boGhiNhatKy.Errorf("Bộ xử lý %d lỗi xử lý tác vụ %s: %v", w.id, tacVu.MaTacVu, err)
				ketQua = model.KetQua{
					MaTacVu:   tacVu.MaTacVu,
					TrangThai: "Thất bại",
					ChiTiet:   "Lỗi: " + err.Error(),
				}
			}
			w.danhSachKetQua <- ketQua
			w.boGhiNhatKy.Debugf("Bộ xử lý %d hoàn thành tác vụ %s với trạng thái %s", w.id, tacVu.MaTacVu, ketQua.TrangThai)

		case <-kenhDung:
			w.boGhiNhatKy.Debugf("Bộ xử lý %d dừng do tín hiệu dừng", w.id)
			return
		}
	}
}
