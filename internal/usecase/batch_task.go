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

// XuLyHangLoatTacVu là usecase để xử lý nhiều tác vụ cùng lúc.
// Trong Clean Architecture, đây là một tầng usecase, nơi chứa
// các quy tắc nghiệp vụ ứng dụng và điều phối các thành phần khác.
type XuLyHangLoatTacVu struct {
	boDocDuLieuVao input.Reader        // boDocDuLieuVao: bộ đọc dữ liệu vào
	boXuLy         service.BoXuLyTacVu // boXuLy: bộ xử lý tác vụ từ domain service
	nhomXuLy       *worker.NhomXuLy    // nhomXuLy: nhóm xử lý - quản lý các worker
	boGhiNhatKy    *utils.Logger       // boGhiNhatKy: bộ ghi nhật ký
}

// TaoBoXuLyHangLoatTacVu tạo một usecase mới để xử lý nhiều tác vụ cùng lúc.
// Áp dụng nguyên tắc Dependency Injection của Clean Architecture,
// các dependency được truyền vào thay vì khởi tạo bên trong.
func TaoBoXuLyHangLoatTacVu(boDocDuLieuVao input.Reader, boXuLy service.BoXuLyTacVu, soLuongXuLy int, boGhiNhatKy *utils.Logger) *XuLyHangLoatTacVu {
	nhomXuLy := worker.TaoNhomXuLy(soLuongXuLy, boXuLy, boGhiNhatKy)
	return &XuLyHangLoatTacVu{
		boDocDuLieuVao: boDocDuLieuVao,
		boXuLy:         boXuLy,
		nhomXuLy:       nhomXuLy,
		boGhiNhatKy:    boGhiNhatKy,
	}
}

// ThucThi thực hiện việc xử lý tất cả tác vụ từ một tệp tin.
// Đây là phương thức chính của usecase, điều phối toàn bộ luồng công việc
// từ việc đọc dữ liệu, xử lý, đến trả về kết quả.
func (uc *XuLyHangLoatTacVu) ThucThi(duongDanFileVao string) ([]model.KetQua, error) {
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
