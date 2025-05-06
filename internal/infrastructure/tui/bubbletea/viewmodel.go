package bubbletea

import (
	"workercli/internal/domain/model"
)

type ViewModel struct {
	KetQua            []model.KetQua
	DanhSachTrungGian []model.TrungGian
}

func (vm *ViewModel) UpdateResults(ketQua []model.KetQua) {
	vm.KetQua = ketQua
}

func (vm *ViewModel) UpdateProxies(danhSachTrungGian []model.TrungGian) {
	vm.DanhSachTrungGian = danhSachTrungGian
}
