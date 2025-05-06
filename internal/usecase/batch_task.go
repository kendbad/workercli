package usecase

import (
	"fmt"
	"workercli/internal/adapter/input"
	"workercli/internal/adapter/worker"
	"workercli/internal/domain/model"
	"workercli/internal/domain/service"
	"workercli/pkg/utils"

	"github.com/sirupsen/logrus"
)

type XuLyLoDongTacVu struct {
	boDocDuLieuVao input.Reader
	boXuLy         service.BoXuLyTacVu
	nhomXuLy       *worker.NhomXuLy
	boGhiNhatKy    *utils.Logger
}

func TaoBoXuLyLoDongTacVu(boDocDuLieuVao input.Reader, boXuLy service.BoXuLyTacVu, soLuongXuLy int, boGhiNhatKy *utils.Logger) *XuLyLoDongTacVu {
	nhomXuLy := worker.TaoNhomXuLy(soLuongXuLy, boXuLy, boGhiNhatKy)
	return &XuLyLoDongTacVu{
		boDocDuLieuVao: boDocDuLieuVao,
		boXuLy:         boXuLy,
		nhomXuLy:       nhomXuLy,
		boGhiNhatKy:    boGhiNhatKy,
	}
}

func (uc *XuLyLoDongTacVu) ThucThi(duongDanFileVao string) ([]model.KetQua, error) {
	uc.boGhiNhatKy.Info(fmt.Sprintf("Bắt đầu xử lý file đầu vào: %s", duongDanFileVao))

	danhSachTacVu, err := uc.boDocDuLieuVao.ReadTasks(duongDanFileVao)
	if err != nil {
		uc.boGhiNhatKy.Errorf("Lỗi đọc file đầu vào: %v", err)
		return nil, err
	}

	uc.nhomXuLy.BatDau()
	for _, tacVu := range danhSachTacVu {
		uc.nhomXuLy.NopTacVu(tacVu)
	}

	ketQua := make([]model.KetQua, 0, len(danhSachTacVu))
	for i := 0; i < len(danhSachTacVu); i++ {
		ketQuaDon := <-uc.nhomXuLy.KetQua()
		ketQua = append(ketQua, ketQuaDon)
	}

	uc.nhomXuLy.Dung()

	uc.boGhiNhatKy.WithFields(logrus.Fields{
		"duongDanFileVao": duongDanFileVao,
		"soLuongTacVu":    len(ketQua),
	}).Info("Hoàn thành xử lý tác vụ")

	return ketQua, nil
}
