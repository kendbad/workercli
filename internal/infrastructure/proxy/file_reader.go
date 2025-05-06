package proxy

import (
	"bufio"
	"os"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type FileReader struct {
	boGhiNhatKy *utils.Logger
}

func NewFileReader(boGhiNhatKy *utils.Logger) *FileReader {
	return &FileReader{boGhiNhatKy: boGhiNhatKy}
}

func (r *FileReader) ReadProxies(duongDanTep string) ([]model.TrungGian, error) {
	duongDanTep = utils.AutoPath(duongDanTep)
	tep, err := os.Open(duongDanTep)
	if err != nil {
		r.boGhiNhatKy.Errorf("Không thể mở file proxy %s: %v", duongDanTep, err)
		return nil, err
	}
	defer tep.Close()

	var danhSachTrungGian []model.TrungGian
	danhSach := bufio.NewScanner(tep)
	for danhSach.Scan() {
		dong := danhSach.Text()
		if dong == "" {
			continue
		}
		trungGian, err := ParseProxy(dong)
		if err != nil {
			r.boGhiNhatKy.Warnf("Bỏ qua proxy không hợp lệ %s: %v", dong, err)
			continue
		}
		danhSachTrungGian = append(danhSachTrungGian, trungGian)
	}

	if err := danhSach.Err(); err != nil {
		r.boGhiNhatKy.Errorf("Lỗi khi đọc file proxy %s: %v", duongDanTep, err)
		return nil, err
	}

	return danhSachTrungGian, nil
}
