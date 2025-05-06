package input

import (
	"bufio"
	"os"
	"strconv"
	"workercli/internal/domain/model"
	"workercli/pkg/utils"
)

type Reader interface {
	ReadTasks(duongDanTep string) ([]model.TacVu, error)
}

type FileReader struct {
	boGhiNhatKy *utils.Logger
}

func NewFileReader(boGhiNhatKy *utils.Logger) *FileReader {
	return &FileReader{boGhiNhatKy: boGhiNhatKy}
}

func (r *FileReader) ReadTasks(duongDanTep string) ([]model.TacVu, error) {
	duongDanTep = utils.AutoPath(duongDanTep)
	r.boGhiNhatKy.Debugf("Đọc file đầu vào: %s", duongDanTep)
	tep, err := os.Open(duongDanTep)
	if err != nil {
		r.boGhiNhatKy.Errorf("Lỗi mở file %s: %v", duongDanTep, err)
		return nil, err
	}
	defer tep.Close()

	var danhSachTacVu []model.TacVu
	scanner := bufio.NewScanner(tep)
	for i := 1; scanner.Scan(); i++ {
		dong := scanner.Text()
		danhSachTacVu = append(danhSachTacVu, model.TacVu{
			MaTacVu: "task-" + strconv.Itoa(i),
			DuLieu:  dong,
		})
	}

	if err := scanner.Err(); err != nil {
		r.boGhiNhatKy.Errorf("Lỗi đọc file %s: %v", duongDanTep, err)
		return nil, err
	}

	r.boGhiNhatKy.Infof("Đọc được %d task từ file", len(danhSachTacVu))
	return danhSachTacVu, nil
}
