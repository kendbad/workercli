package task

import (
	"time"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type BoXuLy struct {
	boGhiNhatKy *utils.Logger
}

func NewProcessor(boGhiNhatKy *utils.Logger) *BoXuLy {
	return &BoXuLy{boGhiNhatKy: boGhiNhatKy}
}

func (p *BoXuLy) XuLyTacVu(tacVu model.TacVu) (model.KetQua, error) {
	p.boGhiNhatKy.Debugf("Xử lý tác vụ %s với dữ liệu: %s", tacVu.MaTacVu, tacVu.DuLieu)
	time.Sleep(10 * time.Millisecond) // Giả lập xử lý nhanh
	ketQua := model.KetQua{
		MaTacVu:   tacVu.MaTacVu,
		TrangThai: "Thành công",
		ChiTiet:   "Tác vụ hoàn thành: " + tacVu.DuLieu,
	}
	p.boGhiNhatKy.Infof("Kết quả tác vụ %s: %s", tacVu.MaTacVu, ketQua.TrangThai)
	return ketQua, nil
}
