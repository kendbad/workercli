package tview

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

func (vm *ViewModel) GetResults() []model.KetQua {
	return vm.KetQua
}

func (vm *ViewModel) GetProxies() []model.TrungGian {
	return vm.DanhSachTrungGian
}

func (vm *ViewModel) GetResult(index int) model.KetQua {
	return vm.KetQua[index]
}

func (vm *ViewModel) GetProxy(index int) model.TrungGian {
	return vm.DanhSachTrungGian[index]
}

func (vm *ViewModel) GetResultCount() int {
	return len(vm.KetQua)
}

func (vm *ViewModel) GetProxyCount() int {
	return len(vm.DanhSachTrungGian)
}
